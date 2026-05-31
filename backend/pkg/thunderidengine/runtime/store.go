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
	"time"
)

// Store holds ephemeral engine state (flow contexts, OAuth codes, PAR, JTI, attribute cache).
type Store interface {
	StoreFlowContext(ctx context.Context, executionID string, data []byte, expiry time.Time) error
	GetFlowContext(ctx context.Context, executionID string) ([]byte, error)
	UpdateFlowContext(ctx context.Context, executionID string, data []byte) error
	DeleteFlowContext(ctx context.Context, executionID string) error

	StoreAuthCode(ctx context.Context, code string, data []byte, expiry time.Time) error
	GetAuthCode(ctx context.Context, code string) ([]byte, error)
	DeleteAuthCode(ctx context.Context, code string) error

	StoreAuthRequest(ctx context.Context, requestID string, data []byte, expiry time.Time) error
	GetAuthRequest(ctx context.Context, requestID string) ([]byte, error)
	DeleteAuthRequest(ctx context.Context, requestID string) error

	StorePAR(ctx context.Context, requestURI string, data []byte, expiry time.Time) error
	GetPAR(ctx context.Context, requestURI string) ([]byte, error)
	DeletePAR(ctx context.Context, requestURI string) error

	StoreJTI(ctx context.Context, jti string, expiry time.Time) error
	ExistsJTI(ctx context.Context, jti string) (bool, error)

	StoreAttributeCache(ctx context.Context, id string, data []byte, expiry time.Time) error
	GetAttributeCache(ctx context.Context, id string) ([]byte, error)
	ExtendAttributeCacheExpiry(ctx context.Context, id string, expiry time.Time) error
	DeleteAttributeCache(ctx context.Context, id string) error
}
