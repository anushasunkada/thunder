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

// Package inmemory provides an in-process implementation of cache.CacheProvider[T].
//
// The implementation is backed by a Go map guarded by a sync.RWMutex.
// It supports:
//   - Per-entry TTL with lazy expiry on read and periodic background cleanup.
//   - Bounded capacity with configurable eviction policies (LRU or LFU).
//   - Concurrent access by multiple goroutines.
//
// Because all state lives in process memory, it is lost on restart and not
// shared across multiple instances. Use the Redis backend for distributed
// or persistent caching.
//
// Usage:
//
//	p, err := inmemory.NewProvider[MyStruct]("sessions", 30*time.Minute,
//	    cache.InMemoryConfig{MaxSize: 1000, EvictionPolicy: cache.EvictionLRU})
//	if err != nil { ... }
//	defer p.Close()
//
//	_ = p.Set(ctx, "key", myStruct)
//	val, ok, _ := p.Get(ctx, "key")
package inmemory

import (
	"container/list"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/asgardeo/thunder/pkg/cache"
)

const (
	defaultCleanupInterval = time.Minute
)

// entry holds a single cached item along with its metadata.
type entry[T any] struct {
	key       string
	value     T
	expiresAt time.Time // zero ⇒ never expires
	freq      int64     // access counter used by LFU eviction
}

// isExpired reports whether the entry has passed its expiry time.
func (e *entry[T]) isExpired() bool {
	return !e.expiresAt.IsZero() && time.Now().After(e.expiresAt)
}

// provider is the in-process implementation of cache.CacheProvider[T].
type provider[T any] struct {
	name       string
	defaultTTL time.Duration
	maxSize    int
	policy     cache.EvictionPolicy

	mu      sync.RWMutex
	items   map[string]*list.Element // key → list element containing *entry[T]
	lruList *list.List               // front = most recently used

	closed atomic.Bool

	hits    atomic.Int64
	misses  atomic.Int64
	sets    atomic.Int64
	deletes atomic.Int64

	stopCleanup chan struct{}
	wg          sync.WaitGroup
}

// NewProvider constructs an in-memory CacheProvider[T].
//
// name must be non-empty. A defaultTTL ≤ 0 means entries never expire.
// cfg.MaxSize ≤ 0 means no capacity limit.
// cfg.CleanupInterval ≤ 0 disables background expiry scanning.
func NewProvider[T any](name string, defaultTTL time.Duration, cfg cache.InMemoryConfig) (cache.CacheProvider[T], error) {
	if name == "" {
		return nil, fmt.Errorf("%w: name must not be empty", cache.ErrInvalidConfig)
	}

	policy := cfg.EvictionPolicy
	if policy == "" {
		policy = cache.EvictionLRU
	}
	if policy != cache.EvictionLRU && policy != cache.EvictionLFU {
		return nil, fmt.Errorf("%w: unknown eviction policy %q", cache.ErrInvalidConfig, policy)
	}

	cleanupInterval := cfg.CleanupInterval
	if cleanupInterval == 0 {
		cleanupInterval = defaultCleanupInterval
	}

	p := &provider[T]{
		name:        name,
		defaultTTL:  defaultTTL,
		maxSize:     cfg.MaxSize,
		policy:      policy,
		items:       make(map[string]*list.Element),
		lruList:     list.New(),
		stopCleanup: make(chan struct{}),
	}

	if cleanupInterval > 0 {
		p.wg.Add(1)
		go p.cleanupLoop(cleanupInterval)
	}

	return p, nil
}

// ---------------------------------------------------------------------------
// cache.CacheProvider[T] implementation
// ---------------------------------------------------------------------------

// Get retrieves the value for key. Expired entries are treated as misses and
// lazily removed from the cache.
//
//   - Cache hit:  (value, true,  nil)
//   - Cache miss: (zero,  false, nil)
//   - Error:      (zero,  false, err)
func (p *provider[T]) Get(_ context.Context, key string) (T, bool, error) {
	var zero T
	if p.closed.Load() {
		return zero, false, cache.ErrProviderClosed
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	elem, ok := p.items[key]
	if !ok {
		p.misses.Add(1)
		return zero, false, nil
	}

	e := elem.Value.(*entry[T])

	if e.isExpired() {
		p.removeElement(elem)
		p.misses.Add(1)
		return zero, false, nil
	}

	// Update recency / frequency for eviction accounting.
	e.freq++
	p.lruList.MoveToFront(elem)

	p.hits.Add(1)
	return e.value, true, nil
}

// Set stores value under key using the provider's defaultTTL, unless
// overridden with cache.WithTTL. A TTL of zero stores the entry without expiry.
func (p *provider[T]) Set(_ context.Context, key string, value T, opts ...cache.SetOption) error {
	if p.closed.Load() {
		return cache.ErrProviderClosed
	}

	ttl := p.resolveTTL(opts)

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if elem, exists := p.items[key]; exists {
		// Update in place — move to front to reset LRU position.
		e := elem.Value.(*entry[T])
		e.value = value
		e.expiresAt = expiresAt
		e.freq++
		p.lruList.MoveToFront(elem)
		p.sets.Add(1)
		return nil
	}

	// Enforce capacity limit before inserting a new entry.
	if p.maxSize > 0 && len(p.items) >= p.maxSize {
		p.evict()
	}

	e := &entry[T]{
		key:       key,
		value:     value,
		expiresAt: expiresAt,
		freq:      1,
	}
	elem := p.lruList.PushFront(e)
	p.items[key] = elem

	p.sets.Add(1)
	return nil
}

// Delete removes one or more keys. Missing keys are silently ignored.
func (p *provider[T]) Delete(_ context.Context, keys ...string) error {
	if p.closed.Load() {
		return cache.ErrProviderClosed
	}
	if len(keys) == 0 {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	var removed int64
	for _, k := range keys {
		if elem, ok := p.items[k]; ok {
			p.removeElement(elem)
			removed++
		}
	}

	p.deletes.Add(removed)
	return nil
}

// Exists reports whether key is present and not expired, without fetching
// the value or updating LRU position.
func (p *provider[T]) Exists(_ context.Context, key string) (bool, error) {
	if p.closed.Load() {
		return false, cache.ErrProviderClosed
	}

	p.mu.RLock()
	elem, ok := p.items[key]
	p.mu.RUnlock()

	if !ok {
		return false, nil
	}
	e := elem.Value.(*entry[T])
	return !e.isExpired(), nil
}

// Expire resets the TTL of an existing key.
// Returns cache.ErrKeyNotFound if the key does not exist or has already expired.
// A ttl of zero removes any expiry (makes the key persistent).
func (p *provider[T]) Expire(_ context.Context, key string, ttl time.Duration) error {
	if p.closed.Load() {
		return cache.ErrProviderClosed
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	elem, ok := p.items[key]
	if !ok {
		return fmt.Errorf("%w: %q", cache.ErrKeyNotFound, key)
	}

	e := elem.Value.(*entry[T])
	if e.isExpired() {
		p.removeElement(elem)
		return fmt.Errorf("%w: %q", cache.ErrKeyNotFound, key)
	}

	if ttl <= 0 {
		e.expiresAt = time.Time{} // no expiry
	} else {
		e.expiresAt = time.Now().Add(ttl)
	}

	return nil
}

// Clear removes every entry from the cache.
func (p *provider[T]) Clear(_ context.Context) error {
	if p.closed.Load() {
		return cache.ErrProviderClosed
	}

	p.mu.Lock()
	n := int64(len(p.items))
	p.items = make(map[string]*list.Element)
	p.lruList.Init()
	p.mu.Unlock()

	p.deletes.Add(n)
	return nil
}

// Ping always returns nil; the in-memory backend is always available.
func (p *provider[T]) Ping(_ context.Context) error {
	if p.closed.Load() {
		return cache.ErrProviderClosed
	}
	return nil
}

// Close stops the background cleanup goroutine and clears all entries.
// Subsequent operations return cache.ErrProviderClosed.
func (p *provider[T]) Close() error {
	if p.closed.Swap(true) {
		return nil // already closed
	}
	close(p.stopCleanup)
	p.wg.Wait()

	p.mu.Lock()
	p.items = make(map[string]*list.Element)
	p.lruList.Init()
	p.mu.Unlock()

	return nil
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
// Eviction
// ---------------------------------------------------------------------------

// evict removes one entry according to the configured policy.
// Must be called with p.mu held for writing.
func (p *provider[T]) evict() {
	switch p.policy {
	case cache.EvictionLFU:
		p.evictLFU()
	default: // LRU
		p.evictLRU()
	}
}

// evictLRU removes the least-recently-used (back of list) entry.
func (p *provider[T]) evictLRU() {
	elem := p.lruList.Back()
	if elem != nil {
		p.removeElement(elem)
		p.deletes.Add(1)
	}
}

// evictLFU removes the entry with the lowest access frequency.
// Ties are broken in favour of the entry closest to the LRU tail.
func (p *provider[T]) evictLFU() {
	var victim *list.Element
	var minFreq int64 = -1

	for elem := p.lruList.Back(); elem != nil; elem = elem.Prev() {
		e := elem.Value.(*entry[T])
		if minFreq < 0 || e.freq < minFreq {
			minFreq = e.freq
			victim = elem
		}
	}

	if victim != nil {
		p.removeElement(victim)
		p.deletes.Add(1)
	}
}

// ---------------------------------------------------------------------------
// Background cleanup
// ---------------------------------------------------------------------------

// cleanupLoop runs on a goroutine and removes expired entries periodically.
func (p *provider[T]) cleanupLoop(interval time.Duration) {
	defer p.wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.deleteExpired()
		case <-p.stopCleanup:
			return
		}
	}
}

// deleteExpired scans all entries and removes those that have expired.
func (p *provider[T]) deleteExpired() {
	p.mu.Lock()
	defer p.mu.Unlock()

	var removed int64
	for elem := p.lruList.Back(); elem != nil; {
		prev := elem.Prev()
		e := elem.Value.(*entry[T])
		if e.isExpired() {
			p.removeElement(elem)
			removed++
		}
		elem = prev
	}

	if removed > 0 {
		p.deletes.Add(removed)
	}
}

// ---------------------------------------------------------------------------
// Helpers (must be called with p.mu held)
// ---------------------------------------------------------------------------

// removeElement removes elem from both the LRU list and the items map.
func (p *provider[T]) removeElement(elem *list.Element) {
	e := elem.Value.(*entry[T])
	p.lruList.Remove(elem)
	delete(p.items, e.key)
}

// resolveTTL returns the effective TTL for a Set call.
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
