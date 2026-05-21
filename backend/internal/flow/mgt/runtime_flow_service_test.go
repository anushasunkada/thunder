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

package flowmgt

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thunder-id/thunderid/internal/flow/common"
	"github.com/thunder-id/thunderid/internal/flow/core"
	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
)

type stubHostFlowProvider struct {
	byID     *HostFlowDefinition
	byHandle *HostFlowDefinition
	err      error
}

func (s *stubHostFlowProvider) GetFlowByID(_ context.Context, _ string) (*HostFlowDefinition, error) {
	return s.byID, s.err
}

func (s *stubHostFlowProvider) GetFlowByHandle(_ context.Context, _, _ string) (*HostFlowDefinition, error) {
	return s.byHandle, s.err
}

type stubGraphBuilder struct{}

func (stubGraphBuilder) GetGraph(
	_ context.Context, _ *CompleteFlowDefinition,
) (core.GraphInterface, *serviceerror.ServiceError) {
	return nil, nil
}

func (stubGraphBuilder) InvalidateCache(context.Context, string) {}

func TestRuntimeFlowDefinitionServiceGetFlow(t *testing.T) {
	svc := NewRuntimeFlowDefinitionService(&stubHostFlowProvider{
		byID: &HostFlowDefinition{
			ID: "flow-1", Handle: "login", Name: "Login",
			FlowType: string(common.FlowTypeAuthentication),
			Nodes:    json.RawMessage(`[]`),
		},
	}, stubGraphBuilder{})
	got, svcErr := svc.GetFlow(context.Background(), "flow-1")
	require.Nil(t, svcErr)
	require.Equal(t, "flow-1", got.ID)
}

func TestRuntimeFlowDefinitionServiceGetFlowByHandle(t *testing.T) {
	svc := NewRuntimeFlowDefinitionService(&stubHostFlowProvider{
		byHandle: &HostFlowDefinition{
			ID: "flow-1", Handle: "login", Name: "Login",
			FlowType: string(common.FlowTypeAuthentication),
			Nodes:    json.RawMessage(`[]`),
		},
	}, stubGraphBuilder{})
	got, svcErr := svc.GetFlowByHandle(context.Background(), "login", common.FlowTypeAuthentication)
	require.Nil(t, svcErr)
	require.Equal(t, "login", got.Handle)
}

func TestRuntimeFlowDefinitionServiceUnsupportedCRUD(t *testing.T) {
	svc := NewRuntimeFlowDefinitionService(&stubHostFlowProvider{}, stubGraphBuilder{})
	_, svcErr := svc.ListFlows(context.Background(), 0, 0, common.FlowTypeAuthentication)
	require.NotNil(t, svcErr)
	_, svcErr = svc.CreateFlow(context.Background(), &FlowDefinition{})
	require.NotNil(t, svcErr)
}
