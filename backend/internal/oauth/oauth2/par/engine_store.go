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

package par

import (
	"context"
	"time"
)

// RuntimePARBackend stores serialized pushed authorization requests.
type RuntimePARBackend interface {
	StorePAR(ctx context.Context, requestURI string, data []byte, expiry time.Time) error
	GetPAR(ctx context.Context, requestURI string) ([]byte, error)
	DeletePAR(ctx context.Context, requestURI string) error
}

// NewPARStoreFromRuntime adapts a runtime backend to PARStore.
func NewPARStoreFromRuntime(backend RuntimePARBackend) PARStore {
	return &runtimePARStore{backend: backend}
}

type runtimePARStore struct {
	backend RuntimePARBackend
}

func (s *runtimePARStore) Store(
	ctx context.Context, request PushedAuthorizationRequest, expirySeconds int64,
) (string, error) {
	randomKey, err := generateRandomKey()
	if err != nil {
		return "", err
	}
	data, err := MarshalPushedAuthorizationRequest(request)
	if err != nil {
		return "", err
	}
	expiry := time.Now().UTC().Add(time.Duration(expirySeconds) * time.Second)
	if err := s.backend.StorePAR(ctx, randomKey, data, expiry); err != nil {
		return "", err
	}
	return randomKey, nil
}

func (s *runtimePARStore) Consume(ctx context.Context, randomKey string) (PushedAuthorizationRequest, bool, error) {
	data, err := s.backend.GetPAR(ctx, randomKey)
	if err != nil {
		return PushedAuthorizationRequest{}, false, err
	}
	if data == nil {
		return PushedAuthorizationRequest{}, false, nil
	}
	if err := s.backend.DeletePAR(ctx, randomKey); err != nil {
		return PushedAuthorizationRequest{}, false, err
	}
	value, err := UnmarshalPushedAuthorizationRequest(data)
	if err != nil {
		return PushedAuthorizationRequest{}, false, err
	}
	return value, true, nil
}
