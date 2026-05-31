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

package authz

import (
	"context"
	"encoding/json"
	"time"

	"github.com/thunder-id/thunderid/internal/system/utils"
)

// RuntimeAuthCodeBackend stores serialized authorization codes.
type RuntimeAuthCodeBackend interface {
	StoreAuthCode(ctx context.Context, code string, data []byte, expiry time.Time) error
	GetAuthCode(ctx context.Context, code string) ([]byte, error)
	DeleteAuthCode(ctx context.Context, code string) error
}

// RuntimeAuthRequestBackend stores serialized authorization requests.
type RuntimeAuthRequestBackend interface {
	StoreAuthRequest(ctx context.Context, requestID string, data []byte, expiry time.Time) error
	GetAuthRequest(ctx context.Context, requestID string) ([]byte, error)
	DeleteAuthRequest(ctx context.Context, requestID string) error
}

// NewAuthorizationCodeStoreFromRuntime adapts a runtime backend to AuthorizationCodeStoreInterface.
func NewAuthorizationCodeStoreFromRuntime(backend RuntimeAuthCodeBackend) AuthorizationCodeStoreInterface {
	return &runtimeAuthCodeStore{backend: backend}
}

// NewAuthorizationRequestStoreFromRuntime adapts a runtime backend to AuthorizationRequestStore.
func NewAuthorizationRequestStoreFromRuntime(backend RuntimeAuthRequestBackend) AuthorizationRequestStore {
	return &runtimeAuthRequestStore{backend: backend}
}

type runtimeAuthCodeStore struct {
	backend RuntimeAuthCodeBackend
}

func (s *runtimeAuthCodeStore) InsertAuthorizationCode(ctx context.Context, authzCode AuthorizationCode) error {
	data, err := json.Marshal(authzCode)
	if err != nil {
		return err
	}
	return s.backend.StoreAuthCode(ctx, authzCode.Code, data, authzCode.ExpiryTime)
}

func (s *runtimeAuthCodeStore) GetAuthorizationCode(ctx context.Context, authCode string) (*AuthorizationCode, error) {
	data, err := s.backend.GetAuthCode(ctx, authCode)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var code AuthorizationCode
	if err := json.Unmarshal(data, &code); err != nil {
		return nil, err
	}
	return &code, nil
}

func (s *runtimeAuthCodeStore) ConsumeAuthorizationCode(ctx context.Context, authCode string) (bool, error) {
	data, err := s.backend.GetAuthCode(ctx, authCode)
	if err != nil {
		return false, err
	}
	if data == nil {
		return false, nil
	}
	if err := s.backend.DeleteAuthCode(ctx, authCode); err != nil {
		return false, err
	}
	return true, nil
}

type runtimeAuthRequestStore struct {
	backend RuntimeAuthRequestBackend
}

func (s *runtimeAuthRequestStore) AddRequest(ctx context.Context, value AuthRequestContext) (string, error) {
	key, err := utils.GenerateUUIDv7()
	if err != nil {
		return "", err
	}
	data, err := MarshalAuthRequestContext(value)
	if err != nil {
		return "", err
	}
	expiry := time.Now().Add(10 * time.Minute)
	if err := s.backend.StoreAuthRequest(ctx, key, data, expiry); err != nil {
		return "", err
	}
	return key, nil
}

func (s *runtimeAuthRequestStore) GetRequest(ctx context.Context, key string) (bool, AuthRequestContext, error) {
	if key == "" {
		return false, AuthRequestContext{}, nil
	}
	data, err := s.backend.GetAuthRequest(ctx, key)
	if err != nil {
		return false, AuthRequestContext{}, err
	}
	if data == nil {
		return false, AuthRequestContext{}, nil
	}
	value, err := UnmarshalAuthRequestContext(data)
	if err != nil {
		return false, AuthRequestContext{}, err
	}
	return true, value, nil
}

func (s *runtimeAuthRequestStore) ClearRequest(ctx context.Context, key string) error {
	return s.backend.DeleteAuthRequest(ctx, key)
}
