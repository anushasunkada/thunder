# cache

A generic, pluggable distributed cache package for Go.

The central abstraction is `CacheProvider[T]` — a type-parameterised interface that hides
the backing technology behind a uniform API. Two built-in backends are provided:

| Backend | Package | Use case |
|---------|---------|----------|
| **Redis** | `pkg/cache/redis` | Distributed / multi-instance deployments |
| **In-memory** | `pkg/cache/inmemory` | Single-process, testing, development |

Swapping backends is a one-line config change; every call site stays the same.

---

## Table of contents

- [Requirements](#requirements)
- [Package layout](#package-layout)
- [Quick start](#quick-start)
- [Interface reference](#interface-reference)
- [Configuration](#configuration)
- [Backends](#backends)
  - [Redis](#redis-backend)
  - [In-memory](#in-memory-backend)
- [Per-call options](#per-call-options)
- [Statistics](#statistics)
- [Error handling](#error-handling)
- [Adding a new backend](#adding-a-new-backend)

---

## Requirements

- **Go 1.26+** (generic type parameters are used throughout)
- **Redis backend only:** [`github.com/redis/go-redis/v9 v9.18.0`](https://github.com/redis/go-redis) — already declared in the module's `go.mod`

---

## Package layout

```
pkg/cache/
├── provider.go        # CacheProvider[T] interface
├── config.go          # Config, RedisConfig, InMemoryConfig
├── model.go           # Stats, CacheType, EvictionPolicy
├── options.go         # SetOption, SetOptions, WithTTL, ApplySetOptions
├── errors.go          # Sentinel error variables
├── cache.go           # New[T] factory function
├── redis/
│   └── provider.go    # Redis implementation
└── inmemory/
    └── provider.go    # In-memory LRU/LFU + TTL implementation
```

---

## Quick start

### Using the factory (recommended)

```go
import "github.com/asgardeo/thunder/pkg/cache"

// Redis backend
p, err := cache.New[UserSession](cache.Config{
    Type:       cache.TypeRedis,
    Name:       "sessions",
    DefaultTTL: 30 * time.Minute,
    Redis: cache.RedisConfig{
        Address:   "localhost:6379",
        KeyPrefix: "thunder:prod",
    },
})
if err != nil {
    log.Fatal(err)
}
defer p.Close()

// Store a value
err = p.Set(ctx, "user:42", session)

// Retrieve a value
val, ok, err := p.Get(ctx, "user:42")
if err != nil {
    // backend error
}
if !ok {
    // cache miss
}

// Override TTL for a single call
err = p.Set(ctx, "user:42", session, cache.WithTTL(5*time.Minute))

// Delete one or more keys
err = p.Delete(ctx, "user:42", "user:99")

// Check existence without fetching the value
exists, err := p.Exists(ctx, "user:42")

// Reset TTL on an existing key
err = p.Expire(ctx, "user:42", 10*time.Minute)

// Remove all entries in this cache's namespace
err = p.Clear(ctx)

// Health check
err = p.Ping(ctx)
```

### Using a backend directly

The sub-packages can be used directly, which is useful when you want to avoid importing
the factory (and its transitive dependencies):

```go
import (
    "github.com/asgardeo/thunder/pkg/cache"
    "github.com/asgardeo/thunder/pkg/cache/redis"
    "github.com/asgardeo/thunder/pkg/cache/inmemory"
)

// Redis
rp, err := redis.NewProvider[UserSession]("sessions", 30*time.Minute, cache.RedisConfig{
    Address:   "localhost:6379",
    KeyPrefix: "thunder:prod",
})

// In-memory (e.g. in tests or single-node deployments)
ip, err := inmemory.NewProvider[UserSession]("sessions", 30*time.Minute, cache.InMemoryConfig{
    MaxSize:        1000,
    EvictionPolicy: cache.EvictionLRU,
})
```

### Dependency injection pattern

Accept `CacheProvider[T]` as an interface in your service so the backend can be swapped
without changing any business logic:

```go
type UserService struct {
    cache cache.CacheProvider[User]
    // ...
}

func NewUserService(c cache.CacheProvider[User]) *UserService {
    return &UserService{cache: c}
}
```

---

## Interface reference

`CacheProvider[T]` is defined in `provider.go`. All implementations must be safe for
concurrent use by multiple goroutines.

```go
type CacheProvider[T any] interface {
    // Get retrieves the cached value for key.
    //   (value, true,  nil) → cache hit
    //   (zero,  false, nil) → cache miss (key absent or expired)
    //   (zero,  false, err) → backend error
    Get(ctx context.Context, key string) (T, bool, error)

    // Set stores value under key.
    // Uses Config.DefaultTTL unless overridden with WithTTL.
    // TTL == 0 means no expiry.
    Set(ctx context.Context, key string, value T, opts ...SetOption) error

    // Delete removes the given keys. Missing keys are silently ignored.
    Delete(ctx context.Context, keys ...string) error

    // Exists reports whether key is present without fetching its value.
    Exists(ctx context.Context, key string) (bool, error)

    // Expire resets the TTL of an existing key.
    // Returns ErrKeyNotFound if the key does not exist.
    // TTL == 0 removes any expiry (makes the key persistent).
    Expire(ctx context.Context, key string, ttl time.Duration) error

    // Clear removes every entry from this cache's namespace.
    Clear(ctx context.Context) error

    // Ping verifies connectivity to the backing store.
    Ping(ctx context.Context) error

    // Close releases resources held by the provider.
    // The provider must not be used after Close returns.
    Close() error

    // Name returns the logical identifier given at construction time.
    Name() string

    // Stats returns a snapshot of operational counters since creation.
    Stats() Stats
}
```

---

## Configuration

### `Config`

Top-level configuration passed to `cache.New[T]`.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Type` | `CacheType` | `TypeInMemory` | Selects the backend (`"redis"` or `"inmemory"`) |
| `Name` | `string` | — | **Required.** Logical name; used to namespace keys |
| `DefaultTTL` | `time.Duration` | `0` | Default expiry for `Set` calls. `0` = never expire |
| `Redis` | `RedisConfig` | — | Redis-specific settings (ignored for in-memory) |
| `InMemory` | `InMemoryConfig` | — | In-memory-specific settings (ignored for Redis) |

### `RedisConfig`

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Address` | `string` | `"localhost:6379"` | Redis server `host:port` |
| `Username` | `string` | `""` | ACL username (Redis 6+) |
| `Password` | `string` | `""` | AUTH password |
| `DB` | `int` | `0` | Logical database index |
| `KeyPrefix` | `string` | `""` | Prepended to all keys as `<KeyPrefix>:<Name>:<key>` |
| `PoolSize` | `int` | go-redis default | Max connections in the pool |
| `DialTimeout` | `time.Duration` | go-redis default (5 s) | Connection timeout |
| `ReadTimeout` | `time.Duration` | go-redis default (3 s) | Socket read timeout |
| `WriteTimeout` | `time.Duration` | go-redis default | Socket write timeout |

### `InMemoryConfig`

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `MaxSize` | `int` | `0` (unlimited) | Max number of entries before eviction |
| `EvictionPolicy` | `EvictionPolicy` | `EvictionLRU` | `"lru"` or `"lfu"` |
| `CleanupInterval` | `time.Duration` | `1 minute` | How often expired entries are swept. `< 0` disables background cleanup |

---

## Backends

### Redis backend

```
pkg/cache/redis/provider.go
```

**Key behaviour:**

- Values are **JSON-serialised** before storage; `T` must be JSON-compatible.
- Keys are namespaced as `<KeyPrefix>:<Name>:<key>`. This lets multiple apps or
  environments share a single Redis instance without collisions.
- `Clear` uses `SCAN + DEL` in batches of 100 to avoid blocking the Redis event loop.
- `Expire` with `ttl == 0` issues a `PERSIST` command, removing any expiry.
- A `PING` is issued at construction time; the constructor returns an error if Redis
  is unreachable, so misconfiguration is caught at start-up.
- Hit/miss/set/delete counters use `sync/atomic` for lock-free concurrent updates.
- `Close` is idempotent; subsequent calls to any method return `ErrProviderClosed`.

**Key namespacing example:**

```
KeyPrefix = "thunder:prod"
Name      = "sessions"
key       = "user:42"

→ Redis key: "thunder:prod:sessions:user:42"
```

### In-memory backend

```
pkg/cache/inmemory/provider.go
```

**Key behaviour:**

- All state is in-process; it is **not shared** between instances and is **lost on restart**.
- Uses a `map[string]*list.Element` + `container/list` for O(1) LRU tracking.
- TTL expiry is **lazy** (checked on `Get`) plus **periodic** (background goroutine).
- `Ping` always returns `nil`; the in-memory store is always available.
- `Close` stops the background cleanup goroutine and clears all entries.
- `Close` is idempotent.

**Eviction policies:**

| Policy | Behaviour |
|--------|-----------|
| `EvictionLRU` | Removes the least-recently-used entry when `MaxSize` is reached |
| `EvictionLFU` | Removes the entry with the lowest access count; ties broken by recency |

LFU uses per-entry access counters. Both policies only activate when `MaxSize > 0`.

---

## Per-call options

`Set` accepts variadic `SetOption` values to override defaults on a per-call basis.

### `WithTTL(ttl time.Duration)`

Overrides `Config.DefaultTTL` for a single `Set` call.

```go
// Store with a specific TTL
p.Set(ctx, "token:xyz", tok, cache.WithTTL(15*time.Minute))

// Store without expiry, even if DefaultTTL is set
p.Set(ctx, "permanent:key", val, cache.WithTTL(0))
```

### Implementing custom options

`SetOption` is `func(*SetOptions)`. Backend authors apply options by calling
`cache.ApplySetOptions(opts)` and reading the result:

```go
func (p *myProvider[T]) Set(ctx context.Context, key string, value T, opts ...cache.SetOption) error {
    o := cache.ApplySetOptions(opts)
    ttl := p.defaultTTL
    if o.TTL != nil {
        ttl = *o.TTL
    }
    // ...
}
```

---

## Statistics

`Stats()` returns a `cache.Stats` snapshot accumulated since the provider was created.

```go
type Stats struct {
    Hits    int64 // Get calls that found a valid entry
    Misses  int64 // Get calls that found nothing (miss or expired)
    Sets    int64 // Successful Set calls
    Deletes int64 // Keys removed via Delete or Clear
}

// Derived metric
func (s Stats) HitRate() float64
```

```go
s := p.Stats()
fmt.Printf("hit rate: %.1f%%  (hits=%d  misses=%d)\n",
    s.HitRate()*100, s.Hits, s.Misses)
```

---

## Error handling

All sentinel errors are declared in `errors.go` and should be compared with `errors.Is`.

| Error | Returned by | Meaning |
|-------|-------------|---------|
| `ErrKeyNotFound` | `Expire` | The target key does not exist in the cache |
| `ErrProviderClosed` | All methods | `Close` has already been called |
| `ErrSerialization` | `Get`, `Set` (Redis) | JSON marshal/unmarshal failed |
| `ErrInvalidConfig` | `New`, `NewProvider` | Config is missing required fields or has invalid values |

```go
val, ok, err := p.Get(ctx, "key")
if errors.Is(err, cache.ErrProviderClosed) {
    // provider was already shut down
}

err = p.Expire(ctx, "key", time.Minute)
if errors.Is(err, cache.ErrKeyNotFound) {
    // key didn't exist; handle gracefully
}
```

---

## Adding a new backend

1. Create a new sub-package, e.g. `pkg/cache/memcached/`.
2. Implement every method of `cache.CacheProvider[T]`:

   ```go
   package memcached

   import (
       "context"
       "time"
       "github.com/asgardeo/thunder/pkg/cache"
   )

   type provider[T any] struct { /* ... */ }

   func NewProvider[T any](name string, defaultTTL time.Duration, cfg MyConfig) (cache.CacheProvider[T], error) {
       // connect, ping, return &provider[T]{...}
   }

   func (p *provider[T]) Get(ctx context.Context, key string) (T, bool, error)             { /* ... */ }
   func (p *provider[T]) Set(ctx context.Context, key string, value T, opts ...cache.SetOption) error { /* ... */ }
   func (p *provider[T]) Delete(ctx context.Context, keys ...string) error                 { /* ... */ }
   func (p *provider[T]) Exists(ctx context.Context, key string) (bool, error)             { /* ... */ }
   func (p *provider[T]) Expire(ctx context.Context, key string, ttl time.Duration) error  { /* ... */ }
   func (p *provider[T]) Clear(ctx context.Context) error                                  { /* ... */ }
   func (p *provider[T]) Ping(ctx context.Context) error                                   { /* ... */ }
   func (p *provider[T]) Close() error                                                      { /* ... */ }
   func (p *provider[T]) Name() string                                                      { /* ... */ }
   func (p *provider[T]) Stats() cache.Stats                                                { /* ... */ }
   ```

3. Add a new `CacheType` constant in `model.go`:

   ```go
   const TypeMemcached CacheType = "memcached"
   ```

4. Add a `case` in the `New[T]` factory in `cache.go`:

   ```go
   case TypeMemcached:
       return memcached.NewProvider[T](cfg.Name, cfg.DefaultTTL, cfg.Memcached)
   ```

5. Add the corresponding config struct to `config.go` and a `Memcached` field to `Config`.

The existing call sites need no changes.
