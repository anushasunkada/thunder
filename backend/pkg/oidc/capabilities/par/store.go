package par

import (
	"context"
	"fmt"
	"time"

	"github.com/asgardeo/thunder/pkg/cache"
)

// CacheStore implements PARStore using a cache.CacheProvider.
//
// Expiry is enforced at the cache layer via TTL: when Store is called,
// the TTL is derived from StoredRequest.ExpiresAt. If ExpiresAt is zero,
// the provider's configured default TTL is used. A Get on an expired entry
// returns (nil, false, nil) from the provider, which this implementation
// maps to ErrRequestURINotFound.
//
// Use NewCacheStore to construct one, or use NewCacheStoreFromConfig to let
// the factory build the underlying provider from a cache.Config.
type CacheStore struct {
	provider cache.CacheProvider[*StoredRequest]
}

// NewCacheStore creates a PARStore backed by an existing cache provider.
// The caller is responsible for the provider's lifecycle (Close, Ping, etc.).
func NewCacheStore(provider cache.CacheProvider[*StoredRequest]) PARStore {
	return &CacheStore{provider: provider}
}

// NewCacheStoreFromConfig creates a PARStore by building a new cache provider
// from cfg. Returns an error if the provider cannot be initialised (e.g. Redis
// unreachable).
//
// Example (Redis):
//
//	store, err := par.NewCacheStoreFromConfig(cache.Config{
//	    Type:       cache.TypeRedis,
//	    Name:       "par",
//	    DefaultTTL: 90 * time.Second,
//	    Redis:      cache.RedisConfig{Address: "localhost:6379"},
//	})
//
// Example (in-memory, useful for testing):
//
//	store, err := par.NewCacheStoreFromConfig(cache.Config{
//	    Type:       cache.TypeInMemory,
//	    Name:       "par",
//	    DefaultTTL: 90 * time.Second,
//	})
func NewCacheStoreFromConfig(cfg cache.Config) (PARStore, error) {
	provider, err := cache.New[*StoredRequest](cfg)
	if err != nil {
		return nil, fmt.Errorf("par: failed to create cache provider: %w", err)
	}
	return &CacheStore{provider: provider}, nil
}

// Store saves req under requestURI with a TTL derived from req.ExpiresAt.
// If req.ExpiresAt is zero the provider's default TTL applies.
func (s *CacheStore) Store(ctx context.Context, requestURI string, req *StoredRequest) error {
	var opts []cache.SetOption
	if req.ExpiresAt != 0 {
		ttl := time.Until(time.Unix(req.ExpiresAt, 0))
		if ttl <= 0 {
			// Already expired before it could be stored; reject immediately.
			return fmt.Errorf("par: pushed request has already expired")
		}
		opts = append(opts, cache.WithTTL(ttl))
	}
	if err := s.provider.Set(ctx, requestURI, req, opts...); err != nil {
		return fmt.Errorf("par: failed to store pushed request: %w", err)
	}
	return nil
}

// GetByRequestURI retrieves the stored pushed authorization request.
// Returns ErrRequestURINotFound if the key is absent or has expired.
func (s *CacheStore) GetByRequestURI(ctx context.Context, requestURI string) (*StoredRequest, error) {
	stored, found, err := s.provider.Get(ctx, requestURI)
	if err != nil {
		return nil, fmt.Errorf("par: cache get error: %w", err)
	}
	if !found {
		return nil, ErrRequestURINotFound
	}
	return stored, nil
}

// Delete removes the pushed request for requestURI.
// Missing keys are silently ignored (the underlying provider guarantees this).
func (s *CacheStore) Delete(ctx context.Context, requestURI string) error {
	if err := s.provider.Delete(ctx, requestURI); err != nil {
		return fmt.Errorf("par: cache delete error: %w", err)
	}
	return nil
}
