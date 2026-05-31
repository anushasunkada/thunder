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

	"github.com/thunder-id/thunderid/internal/flow/common"
	"github.com/thunder-id/thunderid/internal/flow/flowexec"
	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
)

type flowProviderAdapter struct {
	provider FlowSource
}

// NewFlowExecProvider adapts a FlowSource to flowexec.FlowProvider.
func NewFlowExecProvider(provider FlowSource) flowexec.FlowProvider {
	return &flowProviderAdapter{provider: provider}
}

func (a *flowProviderAdapter) GetFlow(
	ctx context.Context, flowID string,
) (*common.CompleteFlowDefinition, *serviceerror.ServiceError) {
	flow, err := a.provider.GetFlow(ctx, flowID)
	if err != nil {
		return nil, toServiceError(err)
	}
	return toCompleteFlowDefinition(flow), nil
}

func (a *flowProviderAdapter) GetFlowByHandle(ctx context.Context, handle string, flowType common.FlowType) (
	*common.CompleteFlowDefinition, *serviceerror.ServiceError) {
	flow, err := a.provider.GetFlowByHandle(ctx, handle, string(flowType))
	if err != nil {
		return nil, toServiceError(err)
	}
	return toCompleteFlowDefinition(flow), nil
}

func toCompleteFlowDefinition(flow *FlowDefinition) *common.CompleteFlowDefinition {
	if flow == nil {
		return nil
	}
	return &common.CompleteFlowDefinition{
		ID:       flow.ID,
		Handle:   flow.Handle,
		Name:     flow.Name,
		FlowType: common.FlowType(flow.FlowType),
		Nodes:    flow.Nodes,
	}
}

func toServiceError(err error) *serviceerror.ServiceError {
	if err == nil {
		return nil
	}
	return &serviceerror.InternalServerError
}
