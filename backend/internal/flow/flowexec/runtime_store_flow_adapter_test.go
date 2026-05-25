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

package flowexec

import (
	"context"
	"errors"
	"testing"

	"github.com/thunder-id/thunderid/pkg/thunderidengine"
)

var errUnexpectedRuntimeStoreCall = errors.New("unexpected runtime store call in flowexec test")

// runtimeStoreFlowAdapter implements thunderidengine.RuntimeStore for tests by delegating
// flow persistence to FlowStoreInterface and rejecting OAuth store calls.
type runtimeStoreFlowAdapter struct {
	flow FlowStoreInterface
}

func newTestFlowRuntimeStore(t *testing.T) (thunderidengine.RuntimeStore, *flowStoreInterfaceMock) {
	t.Helper()
	mock := newFlowStoreInterfaceMock(t)
	return &runtimeStoreFlowAdapter{flow: mock}, mock
}

func (a *runtimeStoreFlowAdapter) Store(
	_ context.Context, _ thunderidengine.PARRequest, _ int64,
) (string, error) {
	return "", errUnexpectedRuntimeStoreCall
}

func (a *runtimeStoreFlowAdapter) Consume(
	_ context.Context, _ string,
) (thunderidengine.PARRequest, bool, error) {
	return thunderidengine.PARRequest{}, false, errUnexpectedRuntimeStoreCall
}

func (a *runtimeStoreFlowAdapter) AddRequest(
	_ context.Context, _ thunderidengine.AuthRequestContext,
) (string, error) {
	return "", errUnexpectedRuntimeStoreCall
}

func (a *runtimeStoreFlowAdapter) GetRequest(
	_ context.Context, _ string,
) (bool, thunderidengine.AuthRequestContext, error) {
	return false, thunderidengine.AuthRequestContext{}, errUnexpectedRuntimeStoreCall
}

func (a *runtimeStoreFlowAdapter) ClearRequest(_ context.Context, _ string) error {
	return errUnexpectedRuntimeStoreCall
}

func (a *runtimeStoreFlowAdapter) InsertAuthorizationCode(
	_ context.Context, _ thunderidengine.AuthorizationCode,
) error {
	return errUnexpectedRuntimeStoreCall
}

func (a *runtimeStoreFlowAdapter) ConsumeAuthorizationCode(_ context.Context, _ string) (bool, error) {
	return false, errUnexpectedRuntimeStoreCall
}

func (a *runtimeStoreFlowAdapter) GetAuthorizationCode(
	_ context.Context, _ string,
) (*thunderidengine.AuthorizationCode, error) {
	return nil, errUnexpectedRuntimeStoreCall
}

func (a *runtimeStoreFlowAdapter) StoreFlowContext(
	ctx context.Context, dbModel thunderidengine.FlowContextDB, expirySeconds int64,
) error {
	return a.flow.StoreFlowContext(ctx, FlowContextDB(dbModel), expirySeconds)
}

func (a *runtimeStoreFlowAdapter) GetFlowContext(
	ctx context.Context, executionID string,
) (*thunderidengine.FlowContextDB, error) {
	dbModel, err := a.flow.GetFlowContext(ctx, executionID)
	if dbModel == nil {
		return nil, err
	}
	engineModel := thunderidengine.FlowContextDB(*dbModel)
	return &engineModel, err
}

func (a *runtimeStoreFlowAdapter) UpdateFlowContext(
	ctx context.Context, dbModel thunderidengine.FlowContextDB,
) error {
	return a.flow.UpdateFlowContext(ctx, FlowContextDB(dbModel))
}

func (a *runtimeStoreFlowAdapter) DeleteFlowContext(ctx context.Context, executionID string) error {
	return a.flow.DeleteFlowContext(ctx, executionID)
}
