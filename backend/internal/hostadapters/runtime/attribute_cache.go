/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
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

// Package runtime adapts ThunderID server runtime stores for the engine Store interface.
package runtime

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/thunder-id/thunderid/internal/attributecache"
	"github.com/thunder-id/thunderid/internal/enginebridge"
)

type attributeCacheBridge struct {
	inner attributecache.PersistStore
}

func newAttributeCacheBridge(inner attributecache.PersistStore) enginebridge.AttributeCacheStore {
	return &attributeCacheBridge{inner: inner}
}

func (b *attributeCacheBridge) CreateAttributeCache(
	ctx context.Context, id string, data []byte, ttlSeconds int,
) error {
	var cache attributecache.AttributeCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return err
	}
	cache.ID = id
	cache.TTLSeconds = ttlSeconds
	return b.inner.CreateAttributeCache(ctx, cache)
}

func (b *attributeCacheBridge) GetAttributeCache(ctx context.Context, id string) ([]byte, int, error) {
	cache, err := b.inner.GetAttributeCache(ctx, id)
	if err != nil {
		if errors.Is(err, attributecache.ErrStoreNotFound()) {
			return nil, 0, enginebridge.ErrNotFound
		}
		return nil, 0, err
	}
	data, err := json.Marshal(cache)
	if err != nil {
		return nil, 0, err
	}
	return data, cache.TTLSeconds, nil
}

func (b *attributeCacheBridge) ExtendAttributeCacheTTL(ctx context.Context, id string, ttlSeconds int) error {
	err := b.inner.ExtendAttributeCacheTTL(ctx, id, ttlSeconds)
	if errors.Is(err, attributecache.ErrStoreNotFound()) {
		return enginebridge.ErrNotFound
	}
	return err
}

func (b *attributeCacheBridge) DeleteAttributeCache(ctx context.Context, id string) error {
	err := b.inner.DeleteAttributeCache(ctx, id)
	if errors.Is(err, attributecache.ErrStoreNotFound()) {
		return enginebridge.ErrNotFound
	}
	return err
}
