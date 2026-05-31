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

package attributecache

import (
	"context"
	"encoding/json"
	"time"
)

// PersistStore persists attribute cache entries.
type PersistStore interface {
	CreateAttributeCache(ctx context.Context, cache AttributeCache) error
	GetAttributeCache(ctx context.Context, id string) (AttributeCache, error)
	ExtendAttributeCacheTTL(ctx context.Context, id string, ttlSeconds int) error
	DeleteAttributeCache(ctx context.Context, id string) error
}

// RuntimeStoreAttributeCache is the subset of RuntimeStore used for attribute cache persistence.
type RuntimeStoreAttributeCache interface {
	StoreAttributeCache(ctx context.Context, id string, data []byte, expiry time.Time) error
	GetAttributeCache(ctx context.Context, id string) ([]byte, error)
	ExtendAttributeCacheExpiry(ctx context.Context, id string, expiry time.Time) error
	DeleteAttributeCache(ctx context.Context, id string) error
}

type runtimePersistStore struct {
	runtime RuntimeStoreAttributeCache
}

// NewPersistStoreFromRuntimeStore adapts RuntimeStore attribute-cache methods to PersistStore.
func NewPersistStoreFromRuntimeStore(runtime RuntimeStoreAttributeCache) PersistStore {
	return &runtimePersistStore{runtime: runtime}
}

// NewServiceFromRuntimeStore creates an attribute cache service backed by RuntimeStore.
func NewServiceFromRuntimeStore(runtime RuntimeStoreAttributeCache) AttributeCacheServiceInterface {
	return newAttributeCacheService(NewPersistStoreFromRuntimeStore(runtime))
}

func (s *runtimePersistStore) CreateAttributeCache(ctx context.Context, cache AttributeCache) error {
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}
	expiry := time.Now().Add(time.Duration(cache.TTLSeconds) * time.Second)
	return s.runtime.StoreAttributeCache(ctx, cache.ID, data, expiry)
}

func (s *runtimePersistStore) GetAttributeCache(ctx context.Context, id string) (AttributeCache, error) {
	data, err := s.runtime.GetAttributeCache(ctx, id)
	if err != nil {
		if isRuntimeStoreNotFound(err) {
			return AttributeCache{}, errAttributeCacheNotFound
		}
		return AttributeCache{}, err
	}
	var cache AttributeCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return AttributeCache{}, err
	}
	cache.ID = id
	return cache, nil
}

func (s *runtimePersistStore) ExtendAttributeCacheTTL(ctx context.Context, id string, ttlSeconds int) error {
	expiry := time.Now().Add(time.Duration(ttlSeconds) * time.Second)
	err := s.runtime.ExtendAttributeCacheExpiry(ctx, id, expiry)
	if isRuntimeStoreNotFound(err) {
		return errAttributeCacheNotFound
	}
	return err
}

func (s *runtimePersistStore) DeleteAttributeCache(ctx context.Context, id string) error {
	err := s.runtime.DeleteAttributeCache(ctx, id)
	if isRuntimeStoreNotFound(err) {
		return errAttributeCacheNotFound
	}
	return err
}
