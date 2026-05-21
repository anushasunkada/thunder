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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thunder-id/thunderid/pkg/thunderidengine"
)

type stubFlowRuntimeStore struct {
	flows map[string]thunderidengine.FlowContext
}

func (s *stubFlowRuntimeStore) Store(
	context.Context, thunderidengine.PushedAuthorizationRequest, int64,
) (string, error) {
	return "", nil
}

func (s *stubFlowRuntimeStore) Consume(
	context.Context, string,
) (thunderidengine.PushedAuthorizationRequest, bool, error) {
	return thunderidengine.PushedAuthorizationRequest{}, false, nil
}

func (s *stubFlowRuntimeStore) AddRequest(context.Context, thunderidengine.AuthRequestContext) (string, error) {
	return "", nil
}

func (s *stubFlowRuntimeStore) GetRequest(context.Context, string) (bool, thunderidengine.AuthRequestContext, error) {
	return false, thunderidengine.AuthRequestContext{}, nil
}

func (s *stubFlowRuntimeStore) ClearRequest(context.Context, string) error { return nil }

func (s *stubFlowRuntimeStore) InsertAuthorizationCode(context.Context, thunderidengine.AuthorizationCode) error {
	return nil
}

func (s *stubFlowRuntimeStore) ConsumeAuthorizationCode(context.Context, string) (bool, error) {
	return false, nil
}

func (s *stubFlowRuntimeStore) GetAuthorizationCode(
	context.Context, string,
) (*thunderidengine.AuthorizationCode, error) {
	return nil, nil
}

func (s *stubFlowRuntimeStore) StoreFlowContext(
	_ context.Context, flow thunderidengine.FlowContext, _ int64,
) error {
	if s.flows == nil {
		s.flows = make(map[string]thunderidengine.FlowContext)
	}
	s.flows[flow.ExecutionID] = flow
	return nil
}

func (s *stubFlowRuntimeStore) GetFlowContext(
	_ context.Context, executionID string,
) (*thunderidengine.FlowContext, error) {
	if s.flows == nil {
		return nil, nil
	}
	flow, ok := s.flows[executionID]
	if !ok {
		return nil, nil
	}
	return &flow, nil
}

func (s *stubFlowRuntimeStore) UpdateFlowContext(_ context.Context, flow thunderidengine.FlowContext) error {
	return s.StoreFlowContext(context.Background(), flow, 0)
}

func (s *stubFlowRuntimeStore) DeleteFlowContext(_ context.Context, executionID string) error {
	delete(s.flows, executionID)
	return nil
}

func TestHostContextStoreRoundTrip(t *testing.T) {
	host := &stubFlowRuntimeStore{}
	store := NewContextStoreFromRuntime(host)
	model := FlowContextDB{ExecutionID: "exec-1", Context: `{"graphId":"g1"}`}
	require.NoError(t, store.StoreFlowContext(context.Background(), model, 60))

	got, err := store.GetFlowContext(context.Background(), "exec-1")
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, `{"graphId":"g1"}`, got.Context)

	model.Context = `{"graphId":"g2"}`
	require.NoError(t, store.UpdateFlowContext(context.Background(), model))
	got, err = store.GetFlowContext(context.Background(), "exec-1")
	require.NoError(t, err)
	require.Equal(t, `{"graphId":"g2"}`, got.Context)

	require.NoError(t, store.DeleteFlowContext(context.Background(), "exec-1"))
	got, err = store.GetFlowContext(context.Background(), "exec-1")
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestHostContextStoreNilHost(t *testing.T) {
	store := NewContextStoreFromRuntime(nil)
	err := store.StoreFlowContext(context.Background(), FlowContextDB{ExecutionID: "x"}, 60)
	require.Error(t, err)
}
