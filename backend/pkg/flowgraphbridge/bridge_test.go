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

package flowgraphbridge

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/thunder-id/thunderid/internal/flow/common"
	flowcore "github.com/thunder-id/thunderid/internal/flow/core"
	flowmgt "github.com/thunder-id/thunderid/internal/flow/mgt"
	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
)

type flowMgtGraphSourceMock struct {
	mock.Mock
}

func (m *flowMgtGraphSourceMock) GetGraph(ctx context.Context, flowID string) (
	flowcore.GraphInterface, *serviceerror.ServiceError) {
	args := m.Called(ctx, flowID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*serviceerror.ServiceError)
	}
	return args.Get(0).(flowcore.GraphInterface), nil
}

func (m *flowMgtGraphSourceMock) GetFlowByHandle(ctx context.Context, handle string, flowType common.FlowType) (
	*flowmgt.CompleteFlowDefinition, *serviceerror.ServiceError) {
	args := m.Called(ctx, handle, flowType)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*serviceerror.ServiceError)
	}
	return args.Get(0).(*flowmgt.CompleteFlowDefinition), nil
}

func TestNewFlowGraphProviderFromMgt_GetGraph(t *testing.T) {
	factory := flowcore.NewFlowFactory()
	graph := factory.CreateGraph("flow-1", common.FlowTypeAuthentication)
	source := &flowMgtGraphSourceMock{}
	source.On("GetGraph", mock.Anything, "flow-1").Return(graph, nil)

	provider := NewFlowGraphProviderFromMgt(source)
	fg, err := provider.GetGraph(context.Background(), "flow-1")
	assert.NoError(t, err)
	assert.Equal(t, "flow-1", fg.GetID())

	core, ok := CoreGraphFromFlowGraph(fg)
	assert.True(t, ok)
	assert.Equal(t, graph, core)
}

func TestServiceErrorFromErr(t *testing.T) {
	svcErr := &serviceerror.ServiceError{
		Code: flowmgt.ErrorFlowNotFound.Code,
		Type: flowmgt.ErrorFlowNotFound.Type,
	}
	err := AsServiceError(svcErr)
	unwrapped, ok := ServiceErrorFromErr(err)
	assert.True(t, ok)
	assert.Equal(t, svcErr.Code, unwrapped.Code)

	_, ok = ServiceErrorFromErr(errors.New("other"))
	assert.False(t, ok)
}
