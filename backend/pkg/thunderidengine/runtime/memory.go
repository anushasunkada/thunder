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

package runtime

import (
	"context"
	"sync"
	"time"
)

type memoryEntry struct {
	data   []byte
	expiry time.Time
}

// MemoryRuntimeStore is an in-process Store for development and tests.
type MemoryRuntimeStore struct {
	mu        sync.RWMutex
	flow      map[string]memoryEntry
	codes     map[string]memoryEntry
	reqs      map[string]memoryEntry
	par       map[string]memoryEntry
	jti       map[string]memoryEntry
	attrCache map[string]memoryEntry
}

// NewMemoryRuntimeStore creates an empty in-memory runtime store.
func NewMemoryRuntimeStore() *MemoryRuntimeStore {
	return &MemoryRuntimeStore{
		flow:      make(map[string]memoryEntry),
		codes:     make(map[string]memoryEntry),
		reqs:      make(map[string]memoryEntry),
		par:       make(map[string]memoryEntry),
		jti:       make(map[string]memoryEntry),
		attrCache: make(map[string]memoryEntry),
	}
}

// StoreFlowContext implements Store.
func (s *MemoryRuntimeStore) StoreFlowContext(
	ctx context.Context, executionID string, data []byte, expiry time.Time,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.flow[executionID] = memoryEntry{data: append([]byte(nil), data...), expiry: expiry}
	return nil
}

// GetFlowContext implements Store.
func (s *MemoryRuntimeStore) GetFlowContext(ctx context.Context, executionID string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.flow[executionID]
	if !ok || time.Now().After(entry.expiry) {
		return nil, ErrNotFound
	}
	return append([]byte(nil), entry.data...), nil
}

// UpdateFlowContext implements Store.
func (s *MemoryRuntimeStore) UpdateFlowContext(
	ctx context.Context, executionID string, data []byte,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.flow[executionID]
	if !ok {
		return ErrNotFound
	}
	entry.data = append([]byte(nil), data...)
	s.flow[executionID] = entry
	return nil
}

// DeleteFlowContext implements Store.
func (s *MemoryRuntimeStore) DeleteFlowContext(ctx context.Context, executionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.flow, executionID)
	return nil
}

// StoreAuthCode implements Store.
func (s *MemoryRuntimeStore) StoreAuthCode(
	ctx context.Context, code string, data []byte, expiry time.Time,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.codes[code] = memoryEntry{data: append([]byte(nil), data...), expiry: expiry}
	return nil
}

// GetAuthCode implements Store.
func (s *MemoryRuntimeStore) GetAuthCode(ctx context.Context, code string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.codes[code]
	if !ok || time.Now().After(entry.expiry) {
		return nil, ErrNotFound
	}
	return append([]byte(nil), entry.data...), nil
}

// DeleteAuthCode implements Store.
func (s *MemoryRuntimeStore) DeleteAuthCode(ctx context.Context, code string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.codes, code)
	return nil
}

// StoreAuthRequest implements Store.
func (s *MemoryRuntimeStore) StoreAuthRequest(
	ctx context.Context, requestID string, data []byte, expiry time.Time,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reqs[requestID] = memoryEntry{data: append([]byte(nil), data...), expiry: expiry}
	return nil
}

// GetAuthRequest implements Store.
func (s *MemoryRuntimeStore) GetAuthRequest(ctx context.Context, requestID string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.reqs[requestID]
	if !ok || time.Now().After(entry.expiry) {
		return nil, ErrNotFound
	}
	return append([]byte(nil), entry.data...), nil
}

// DeleteAuthRequest implements Store.
func (s *MemoryRuntimeStore) DeleteAuthRequest(ctx context.Context, requestID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.reqs, requestID)
	return nil
}

// StorePAR implements Store.
func (s *MemoryRuntimeStore) StorePAR(
	ctx context.Context, requestURI string, data []byte, expiry time.Time,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.par[requestURI] = memoryEntry{data: append([]byte(nil), data...), expiry: expiry}
	return nil
}

// GetPAR implements Store.
func (s *MemoryRuntimeStore) GetPAR(ctx context.Context, requestURI string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.par[requestURI]
	if !ok || time.Now().After(entry.expiry) {
		return nil, ErrNotFound
	}
	return append([]byte(nil), entry.data...), nil
}

// DeletePAR implements Store.
func (s *MemoryRuntimeStore) DeletePAR(ctx context.Context, requestURI string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.par, requestURI)
	return nil
}

// StoreJTI implements Store.
func (s *MemoryRuntimeStore) StoreJTI(ctx context.Context, jti string, expiry time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jti[jti] = memoryEntry{expiry: expiry}
	return nil
}

// ExistsJTI implements Store.
func (s *MemoryRuntimeStore) ExistsJTI(ctx context.Context, jti string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.jti[jti]
	if !ok || time.Now().After(entry.expiry) {
		return false, nil
	}
	return true, nil
}

// StoreAttributeCache implements Store.
func (s *MemoryRuntimeStore) StoreAttributeCache(
	ctx context.Context, id string, data []byte, expiry time.Time,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.attrCache[id] = memoryEntry{data: append([]byte(nil), data...), expiry: expiry}
	return nil
}

// GetAttributeCache implements Store.
func (s *MemoryRuntimeStore) GetAttributeCache(ctx context.Context, id string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.attrCache[id]
	if !ok || time.Now().After(entry.expiry) {
		return nil, ErrNotFound
	}
	return append([]byte(nil), entry.data...), nil
}

// ExtendAttributeCacheExpiry implements Store.
func (s *MemoryRuntimeStore) ExtendAttributeCacheExpiry(
	ctx context.Context, id string, expiry time.Time,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.attrCache[id]
	if !ok || time.Now().After(entry.expiry) {
		return ErrNotFound
	}
	entry.expiry = expiry
	s.attrCache[id] = entry
	return nil
}

// DeleteAttributeCache implements Store.
func (s *MemoryRuntimeStore) DeleteAttributeCache(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.attrCache, id)
	return nil
}
