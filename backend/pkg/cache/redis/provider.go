/*
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package redis provides a Redis-backed implementation of cache.CacheProvider[T].
//
// Values are JSON-serialised before being written to Redis and deserialised on
// retrieval, so T must be JSON-compatible.
//
// All cache keys are namespaced using the pattern:
//
//	<KeyPrefix>:<Name>:<key>
//
// where KeyPrefix and Name come from the RedisConfig / name argument passed to
// NewProvider. This prevents key collisions when a single Redis instance is
// shared across multiple deployments or cache instances.
//
// Usage:
//
//	p, err := redis.NewProvider[MyStruct]("sessions", 30*time.Minute, cache.RedisConfig{
//	    Address:   "localhost:6379",
//	    KeyPrefix: "thunder:prod",
//	})
//	if err != nil { ... }
//	defer p.Close()
//
//	_ = p.Set(ctx, "user:42", myStruct)
//	val, ok, _ := p.Get(ctx, "user:42")
package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/asgardeo/thunder/pkg/cache"
	goredis "github.com/redis/go-redis/v9"
)

// provider is the Redis-backed implementation of cache.CacheProvider[T].
type provider[T any] struct {
	name       string
	client     *goredis.Client
	defaultTTL time.Duration
	keyPrefix  string // pre-built "<KeyPrefix>:<Name>" string

	closed atomic.Bool

	hits    atomic.Int64
	misses  atomic.Int64
	sets    atomic.Int64
	deletes atomic.Int64

	closeOnce sync.Once
}

// NewProvider constructs a Redis CacheProvider[T].
//
// It dials the server described by cfg and returns an error if the initial
// PING fails, so callers detect misconfiguration at start-up rather than at
// the first cache operation.
//
// name must be non-empty. A defaultTTL ≤ 0 means entries never expire.
func NewProvider[T any](name string, defaultTTL time.Duration, cfg cache.RedisConfig) (cache.CacheProvider[T], error) {
	if name == "" {
		return nil, fmt.Errorf("%w: name must not be empty", cache.ErrInvalidConfig)
	}

	addr := cfg.Address
	if addr == "" {
		addr = "localhost:6379"
	}

	client := goredis.NewClient(&goredis.Options{
		Addr:         addr,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("cache/redis: failed to connect to %s: %w", addr, err)
	}

	return &provider[T]{
		name:       name,
		client:     client,
		defaultTTL: defaultTTL,
		keyPrefix:  buildKeyPrefix(cfg.KeyPrefix, name),
	}, nil
}

// ---------------------------------------------------------------------------
// cache.CacheProvider[T] implementation
// ---------------------------------------------------------------------------

// Get retrieves the value stored under key.
//
//   - Cache hit:  (value, true,  nil)
//   - Cache miss: (zero,  false, nil)
//   - Error:      (zero,  false, err)
func (p *provider[T]) Get(ctx context.Context, key string) (T, bool, error) {
	var zero T
	if p.closed.Load() {
		return zero, false, cache.ErrProviderClosed
	}

	data, err := p.client.Get(ctx, p.prefixedKey(key)).Bytes()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			p.misses.Add(1)
			return zero, false, nil
		}
		return zero, false, fmt.Errorf("cache/redis: GET %q: %w", key, err)
	}

	var value T
	if err = json.Unmarshal(data, &value); err != nil {
		return zero, false, fmt.Errorf("%w: unmarshal key %q: %v", cache.ErrSerialization, key, err)
	}

	p.hits.Add(1)
	return value, true, nil
}

// Set stores value under key.
// The provider's defaultTTL is used unless overridden with cache.WithTTL.
// A TTL of zero stores the entry without expiry.
func (p *provider[T]) Set(ctx context.Context, key string, value T, opts ...cache.SetOption) error {
	if p.closed.Load() {
		return cache.ErrProviderClosed
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("%w: marshal key %q: %v", cache.ErrSerialization, key, err)
	}

	ttl := p.resolveTTL(opts)

	if err = p.client.Set(ctx, p.prefixedKey(key), data, ttl).Err(); err != nil {
		return fmt.Errorf("cache/redis: SET %q: %w", key, err)
	}

	p.sets.Add(1)
	return nil
}

// Delete removes one or more keys from the cache. Missing keys are silently
// ignored; this method never returns an error for absent keys.
func (p *provider[T]) Delete(ctx context.Context, keys ...string) error {
	if p.closed.Load() {
		return cache.ErrProviderClosed
	}
	if len(keys) == 0 {
		return nil
	}

	prefixed := make([]string, len(keys))
	for i, k := range keys {
		prefixed[i] = p.prefixedKey(k)
	}

	n, err := p.client.Del(ctx, prefixed...).Result()
	if err != nil {
		return fmt.Errorf("cache/redis: DEL: %w", err)
	}

	p.deletes.Add(n)
	return nil
}

// Exists reports whether key is present in the cache.
func (p *provider[T]) Exists(ctx context.Context, key string) (bool, error) {
	if p.closed.Load() {
		return false, cache.ErrProviderClosed
	}

	n, err := p.client.Exists(ctx, p.prefixedKey(key)).Result()
	if err != nil {
		return false, fmt.Errorf("cache/redis: EXISTS %q: %w", key, err)
	}

	return n > 0, nil
}

// Expire resets the TTL of an existing key.
// Returns cache.ErrKeyNotFound if the key does not exist.
// A ttl of zero removes any expiry, making the key persistent.
func (p *provider[T]) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if p.closed.Load() {
		return cache.ErrProviderClosed
	}

	var (
		ok  bool
		err error
	)

	if ttl <= 0 {
		// ttl == 0 → remove expiry (PERSIST)
		ok, err = p.client.Persist(ctx, p.prefixedKey(key)).Result()
	} else {
		ok, err = p.client.Expire(ctx, p.prefixedKey(key), ttl).Result()
	}

	if err != nil {
		return fmt.Errorf("cache/redis: EXPIRE %q: %w", key, err)
	}
	if !ok {
		return fmt.Errorf("%w: %q", cache.ErrKeyNotFound, key)
	}

	return nil
}

// Clear removes all keys in this provider's namespace using SCAN + DEL to
// avoid blocking the Redis event loop on large key sets.
func (p *provider[T]) Clear(ctx context.Context) error {
	if p.closed.Load() {
		return cache.ErrProviderClosed
	}

	pattern := p.keyPrefix + ":*"
	var (
		cursor  uint64
		deleted int64
	)

	for {
		keys, nextCursor, err := p.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("cache/redis: SCAN during Clear: %w", err)
		}

		if len(keys) > 0 {
			n, err := p.client.Del(ctx, keys...).Result()
			if err != nil {
				return fmt.Errorf("cache/redis: DEL during Clear: %w", err)
			}
			deleted += n
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	p.deletes.Add(deleted)
	return nil
}

// Ping verifies that the Redis server is reachable.
func (p *provider[T]) Ping(ctx context.Context) error {
	if p.closed.Load() {
		return cache.ErrProviderClosed
	}

	if err := p.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("cache/redis: PING: %w", err)
	}

	return nil
}

// Close closes the underlying connection pool.
// It is safe to call Close multiple times; only the first call has effect.
// All subsequent cache operations will return cache.ErrProviderClosed.
func (p *provider[T]) Close() error {
	var closeErr error
	p.closeOnce.Do(func() {
		p.closed.Store(true)
		closeErr = p.client.Close()
	})
	return closeErr
}

// Name returns the logical identifier given at construction time.
func (p *provider[T]) Name() string { return p.name }

// Stats returns a point-in-time snapshot of operational counters.
func (p *provider[T]) Stats() cache.Stats {
	return cache.Stats{
		Hits:    p.hits.Load(),
		Misses:  p.misses.Load(),
		Sets:    p.sets.Load(),
		Deletes: p.deletes.Load(),
	}
}

// ---------------------------------------------------------------------------
// internal helpers
// ---------------------------------------------------------------------------

// prefixedKey returns the fully-qualified Redis key for the given logical key.
func (p *provider[T]) prefixedKey(key string) string {
	return p.keyPrefix + ":" + key
}

// resolveTTL returns the TTL that should be passed to Redis SET.
// Priority: per-call WithTTL > provider defaultTTL > 0 (no expiry).
func (p *provider[T]) resolveTTL(opts []cache.SetOption) time.Duration {
	o := cache.ApplySetOptions(opts)
	if o.TTL != nil {
		return *o.TTL
	}
	if p.defaultTTL > 0 {
		return p.defaultTTL
	}
	return 0
}

// buildKeyPrefix returns "<keyPrefix>:<name>" or just "<name>" when keyPrefix
// is empty.
func buildKeyPrefix(keyPrefix, name string) string {
	keyPrefix = strings.TrimSpace(keyPrefix)
	if keyPrefix == "" {
		return name
	}
	return keyPrefix + ":" + name
}
