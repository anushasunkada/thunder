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
	"encoding/json"
	"errors"
	"time"

	"github.com/thunder-id/thunderid/internal/flow/flowexec"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/authz"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/jti"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/par"
)

// LegacyRuntimeStores groups typed runtime stores used by ThunderID host adapters.
type LegacyRuntimeStores struct {
	FlowStore      flowexec.RuntimeFlowContextStore
	AuthCode       authz.AuthorizationCodeStoreInterface
	AuthReq        authz.AuthorizationRequestStore
	PAR            par.PARStore
	JTI            jti.JTIStoreInterface
	AttributeCache AttributeCacheStore
}

// NewRuntimeStore adapts typed legacy stores to the RuntimeStore interface.
func NewRuntimeStore(stores LegacyRuntimeStores) RuntimeStore {
	return &runtimeStoreAdapter{stores: stores}
}

type runtimeStoreAdapter struct {
	stores LegacyRuntimeStores
}

func (a *runtimeStoreAdapter) StoreFlowContext(
	ctx context.Context, executionID string, data []byte, expiry time.Time,
) error {
	var model flowexec.FlowContextDB
	if err := json.Unmarshal(data, &model); err != nil {
		return err
	}
	model.ExecutionID = executionID
	seconds := int64(time.Until(expiry).Seconds())
	if seconds < 1 {
		seconds = 1
	}
	return a.stores.FlowStore.StoreFlowContext(ctx, model, seconds)
}

func (a *runtimeStoreAdapter) GetFlowContext(ctx context.Context, executionID string) ([]byte, error) {
	model, err := a.stores.FlowStore.GetFlowContext(ctx, executionID)
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, ErrNotFound
	}
	return json.Marshal(model)
}

func (a *runtimeStoreAdapter) UpdateFlowContext(ctx context.Context, executionID string, data []byte) error {
	var model flowexec.FlowContextDB
	if err := json.Unmarshal(data, &model); err != nil {
		return err
	}
	model.ExecutionID = executionID
	return a.stores.FlowStore.UpdateFlowContext(ctx, model)
}

func (a *runtimeStoreAdapter) DeleteFlowContext(ctx context.Context, executionID string) error {
	return a.stores.FlowStore.DeleteFlowContext(ctx, executionID)
}

func (a *runtimeStoreAdapter) StoreAuthCode(ctx context.Context, code string, data []byte, expiry time.Time) error {
	var authCode authz.AuthorizationCode
	if err := json.Unmarshal(data, &authCode); err != nil {
		return err
	}
	authCode.Code = code
	authCode.ExpiryTime = expiry
	return a.stores.AuthCode.InsertAuthorizationCode(ctx, authCode)
}

func (a *runtimeStoreAdapter) GetAuthCode(ctx context.Context, code string) ([]byte, error) {
	authCode, err := a.stores.AuthCode.GetAuthorizationCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if authCode == nil {
		return nil, ErrNotFound
	}
	return json.Marshal(authCode)
}

func (a *runtimeStoreAdapter) DeleteAuthCode(ctx context.Context, code string) error {
	consumed, err := a.stores.AuthCode.ConsumeAuthorizationCode(ctx, code)
	if err != nil {
		return err
	}
	if !consumed {
		return ErrNotFound
	}
	return nil
}

func (a *runtimeStoreAdapter) StoreAuthRequest(
	ctx context.Context, requestID string, data []byte, expiry time.Time,
) error {
	req, err := authz.UnmarshalAuthRequestContext(data)
	if err != nil {
		return err
	}
	key, err := a.stores.AuthReq.AddRequest(ctx, req)
	if err != nil {
		return err
	}
	if key != requestID {
		return errors.New("auth request key mismatch")
	}
	return nil
}

func (a *runtimeStoreAdapter) GetAuthRequest(ctx context.Context, requestID string) ([]byte, error) {
	found, req, err := a.stores.AuthReq.GetRequest(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrNotFound
	}
	return authz.MarshalAuthRequestContext(req)
}

func (a *runtimeStoreAdapter) DeleteAuthRequest(ctx context.Context, requestID string) error {
	return a.stores.AuthReq.ClearRequest(ctx, requestID)
}

func (a *runtimeStoreAdapter) StorePAR(ctx context.Context, requestURI string, data []byte, expiry time.Time) error {
	req, err := par.UnmarshalPushedAuthorizationRequest(data)
	if err != nil {
		return err
	}
	seconds := int64(time.Until(expiry).Seconds())
	if seconds < 1 {
		seconds = 1
	}
	key, err := a.stores.PAR.Store(ctx, req, seconds)
	if err != nil {
		return err
	}
	if key != requestURI {
		return errors.New("par key mismatch")
	}
	return nil
}

func (a *runtimeStoreAdapter) GetPAR(ctx context.Context, requestURI string) ([]byte, error) {
	req, ok, err := a.stores.PAR.Consume(ctx, requestURI)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotFound
	}
	return par.MarshalPushedAuthorizationRequest(req)
}

func (a *runtimeStoreAdapter) DeletePAR(ctx context.Context, requestURI string) error {
	_, ok, err := a.stores.PAR.Consume(ctx, requestURI)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

func (a *runtimeStoreAdapter) StoreJTI(ctx context.Context, jtiValue string, expiry time.Time) error {
	inserted, err := a.stores.JTI.RecordJTI(ctx, "engine", jtiValue, expiry)
	if err != nil {
		return err
	}
	if !inserted {
		return errors.New("jti replay detected")
	}
	return nil
}

func (a *runtimeStoreAdapter) ExistsJTI(ctx context.Context, jtiValue string) (bool, error) {
	inserted, err := a.stores.JTI.RecordJTI(ctx, "engine-check", jtiValue, time.Now().Add(time.Minute))
	if err != nil {
		return false, err
	}
	return !inserted, nil
}
