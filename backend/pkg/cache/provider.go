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

// Package cache defines a generic, pluggable interface for distributed cache operations.
//
// The central abstraction is CacheProvider[T], a type-parameterised interface that
// hides the backing technology (Redis, in-memory, …) behind a uniform API. Two
// built-in implementations are provided in sub-packages:
//
//   - cache/redis    – backed by Redis via go-redis/v9
//   - cache/inmemory – backed by an in-process LRU map with TTL eviction
//
// Quick start:
//
//	// Build with the factory (picks the backend from Config.Type):
//	p, err := cache.New[MyStruct](cache.Config{
//	    Type:       cache.TypeRedis,
//	    Name:       "sessions",
//	    DefaultTTL: 30 * time.Minute,
//	    Redis: cache.RedisConfig{Address: "localhost:6379"},
//	})
//
//	// Or directly instantiate a backend:
//	p, err := redis.NewProvider[MyStruct](cache.Config{...})
//
//	// Use it:
//	_ = p.Set(ctx, "key", myValue)
//	val, ok, _ := p.Get(ctx, "key")
package cache

import (
	"context"
	"time"
)

// CacheProvider is a generic, type-safe interface for cache operations.
// All implementations must be safe for concurrent use by multiple goroutines.
//
// The type parameter T is the value type stored and retrieved from the cache.
// Implementations are responsible for serialising T to/from the underlying
// storage format (e.g. JSON for Redis, direct pointer storage for in-memory).
type CacheProvider[T any] interface {
	// Get retrieves the cached value for key.
	//
	// Return semantics:
	//   (value, true,  nil) – cache hit; value is valid.
	//   (zero,  false, nil) – cache miss; key does not exist or has expired.
	//   (zero,  false, err) – a backend error occurred.
	Get(ctx context.Context, key string) (T, bool, error)

	// Set stores value under key.
	//
	// The default TTL from the provider's Config is used unless overridden
	// with the WithTTL option. A TTL of zero means the entry never expires.
	Set(ctx context.Context, key string, value T, opts ...SetOption) error

	// Delete removes the given keys from the cache.
	// Missing keys are silently ignored; the operation never returns an error
	// for keys that do not exist.
	Delete(ctx context.Context, keys ...string) error

	// Exists reports whether key is present in the cache without fetching its
	// value. This is cheaper than Get when only presence is needed.
	Exists(ctx context.Context, key string) (bool, error)

	// Expire resets the remaining TTL of an existing key.
	// Returns ErrKeyNotFound if the key does not exist.
	// A ttl of zero removes any expiry, making the key persistent.
	Expire(ctx context.Context, key string, ttl time.Duration) error

	// Clear removes every entry from this cache's namespace.
	// For keyed-prefix providers (e.g. Redis with a KeyPrefix) only keys
	// belonging to this provider's namespace are removed.
	Clear(ctx context.Context) error

	// Ping verifies that the backing store is reachable and responsive.
	// Returns nil on success.
	Ping(ctx context.Context) error

	// Close releases any resources held by the provider (connections, goroutines,
	// etc.). The provider must not be used after Close returns.
	Close() error

	// Name returns the logical identifier that was given to this cache instance
	// via Config.Name. It is used to namespace keys and in log output.
	Name() string

	// Stats returns a snapshot of operational counters accumulated since the
	// provider was created or last reset.
	Stats() Stats
}
