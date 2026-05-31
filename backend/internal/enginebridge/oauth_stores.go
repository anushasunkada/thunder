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
	"errors"
	"time"

	"github.com/thunder-id/thunderid/internal/oauth/oauth2/authz"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/jti"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/par"
)

// OAuthStores groups OAuth runtime stores derived from a RuntimeStore.
type OAuthStores struct {
	AuthCode authz.AuthorizationCodeStoreInterface
	AuthReq  authz.AuthorizationRequestStore
	PAR      par.PARStore
	JTI      jti.JTIStoreInterface
}

type runtimeOAuthBackend struct {
	store RuntimeStore
}

func (b *runtimeOAuthBackend) StoreAuthCode(ctx context.Context, code string, data []byte, expiry time.Time) error {
	return b.store.StoreAuthCode(ctx, code, data, expiry)
}

func (b *runtimeOAuthBackend) GetAuthCode(ctx context.Context, code string) ([]byte, error) {
	data, err := b.store.GetAuthCode(ctx, code)
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	}
	return data, err
}

func (b *runtimeOAuthBackend) DeleteAuthCode(ctx context.Context, code string) error {
	return b.store.DeleteAuthCode(ctx, code)
}

func (b *runtimeOAuthBackend) StoreAuthRequest(
	ctx context.Context, requestID string, data []byte, expiry time.Time,
) error {
	return b.store.StoreAuthRequest(ctx, requestID, data, expiry)
}

func (b *runtimeOAuthBackend) GetAuthRequest(ctx context.Context, requestID string) ([]byte, error) {
	data, err := b.store.GetAuthRequest(ctx, requestID)
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	}
	return data, err
}

func (b *runtimeOAuthBackend) DeleteAuthRequest(ctx context.Context, requestID string) error {
	return b.store.DeleteAuthRequest(ctx, requestID)
}

func (b *runtimeOAuthBackend) StorePAR(ctx context.Context, requestURI string, data []byte, expiry time.Time) error {
	return b.store.StorePAR(ctx, requestURI, data, expiry)
}

func (b *runtimeOAuthBackend) GetPAR(ctx context.Context, requestURI string) ([]byte, error) {
	data, err := b.store.GetPAR(ctx, requestURI)
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	}
	return data, err
}

func (b *runtimeOAuthBackend) DeletePAR(ctx context.Context, requestURI string) error {
	return b.store.DeletePAR(ctx, requestURI)
}

func (b *runtimeOAuthBackend) StoreJTI(ctx context.Context, jtiValue string, expiry time.Time) error {
	return b.store.StoreJTI(ctx, jtiValue, expiry)
}

func (b *runtimeOAuthBackend) ExistsJTI(ctx context.Context, jtiValue string) (bool, error) {
	return b.store.ExistsJTI(ctx, jtiValue)
}

// NewOAuthStores adapts a RuntimeStore to typed OAuth store interfaces.
func NewOAuthStores(store RuntimeStore) OAuthStores {
	backend := &runtimeOAuthBackend{store: store}
	return OAuthStores{
		AuthCode: authz.NewAuthorizationCodeStoreFromRuntime(backend),
		AuthReq:  authz.NewAuthorizationRequestStoreFromRuntime(backend),
		PAR:      par.NewPARStoreFromRuntime(backend),
		JTI:      jti.NewJTIStoreFromRuntime(backend),
	}
}
