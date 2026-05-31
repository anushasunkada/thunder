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

package jti

import (
	"context"
	"time"
)

// RuntimeJTIBackend stores JTI replay markers.
type RuntimeJTIBackend interface {
	StoreJTI(ctx context.Context, jti string, expiry time.Time) error
	ExistsJTI(ctx context.Context, jti string) (bool, error)
}

// NewJTIStoreFromRuntime adapts a runtime backend to JTIStoreInterface.
func NewJTIStoreFromRuntime(backend RuntimeJTIBackend) JTIStoreInterface {
	return &runtimeJTIStore{backend: backend}
}

type runtimeJTIStore struct {
	backend RuntimeJTIBackend
}

func (s *runtimeJTIStore) RecordJTI(ctx context.Context, namespace, jtiValue string, expiry time.Time) (bool, error) {
	key := namespace + ":" + jtiValue
	exists, err := s.backend.ExistsJTI(ctx, key)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}
	if err := s.backend.StoreJTI(ctx, key, expiry); err != nil {
		return false, err
	}
	return true, nil
}
