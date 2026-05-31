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

package enginebridge

import (
	"context"
	"time"
)

// AttributeCacheStore persists attribute cache entries for legacy runtime composition.
type AttributeCacheStore interface {
	CreateAttributeCache(ctx context.Context, id string, data []byte, ttlSeconds int) error
	GetAttributeCache(ctx context.Context, id string) (data []byte, ttlSeconds int, err error)
	ExtendAttributeCacheTTL(ctx context.Context, id string, ttlSeconds int) error
	DeleteAttributeCache(ctx context.Context, id string) error
}

func (a *runtimeStoreAdapter) StoreAttributeCache(ctx context.Context, id string, data []byte, expiry time.Time) error {
	if a.stores.AttributeCache == nil {
		return ErrNotFound
	}
	ttlSeconds := int(time.Until(expiry).Seconds())
	if ttlSeconds < 1 {
		ttlSeconds = 1
	}
	return a.stores.AttributeCache.CreateAttributeCache(ctx, id, data, ttlSeconds)
}

func (a *runtimeStoreAdapter) GetAttributeCache(ctx context.Context, id string) ([]byte, error) {
	if a.stores.AttributeCache == nil {
		return nil, ErrNotFound
	}
	data, _, err := a.stores.AttributeCache.GetAttributeCache(ctx, id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (a *runtimeStoreAdapter) ExtendAttributeCacheExpiry(ctx context.Context, id string, expiry time.Time) error {
	if a.stores.AttributeCache == nil {
		return ErrNotFound
	}
	ttlSeconds := int(time.Until(expiry).Seconds())
	if ttlSeconds < 1 {
		ttlSeconds = 1
	}
	return a.stores.AttributeCache.ExtendAttributeCacheTTL(ctx, id, ttlSeconds)
}

func (a *runtimeStoreAdapter) DeleteAttributeCache(ctx context.Context, id string) error {
	if a.stores.AttributeCache == nil {
		return ErrNotFound
	}
	return a.stores.AttributeCache.DeleteAttributeCache(ctx, id)
}
