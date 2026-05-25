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
	"fmt"

	appmodel "github.com/thunder-id/thunderid/internal/application/model"
	authncm "github.com/thunder-id/thunderid/internal/authn/common"
	authnprovidercm "github.com/thunder-id/thunderid/internal/authnprovider/common"
	authnprovidermgr "github.com/thunder-id/thunderid/internal/authnprovider/manager"
	"github.com/thunder-id/thunderid/internal/flow/common"
	flowcore "github.com/thunder-id/thunderid/internal/flow/core"
	"github.com/thunder-id/thunderid/internal/flow/executor"
	"github.com/thunder-id/thunderid/pkg/thunderidengine"
)

type executorRegistryBridge struct {
	registry thunderidengine.ExecutorRegistry
}

func newExecutorRegistryBridge(registry thunderidengine.ExecutorRegistry) executor.ExecutorRegistryInterface {
	return &executorRegistryBridge{registry: registry}
}

func (r *executorRegistryBridge) GetExecutor(name string) (flowcore.ExecutorInterface, error) {
	ex, ok := r.registry.GetExecutor(name)
	if !ok {
		return nil, fmt.Errorf("executor %q not found", name)
	}
	return &executorBridge{engine: ex}, nil
}

func (r *executorRegistryBridge) RegisterExecutor(name string, ex flowcore.ExecutorInterface) {
	if ex == nil {
		return
	}
	_ = r.registry.Register(name, &internalExecutorBridge{inner: ex})
}

func (r *executorRegistryBridge) IsRegistered(name string) bool {
	return r.registry.IsRegistered(name)
}

type executorBridge struct {
	engine thunderidengine.ExecutorInterface
}

func (e *executorBridge) Execute(ctx *flowcore.NodeContext) (*common.ExecutorResponse, error) {
	resp, err := e.engine.Execute(toEngineNodeContext(ctx))
	if err != nil {
		return nil, err
	}
	return toCommonExecutorResponse(resp), nil
}

func (e *executorBridge) GetName() string {
	return e.engine.GetName()
}

func (e *executorBridge) GetType() common.ExecutorType {
	return common.ExecutorType(e.engine.GetType())
}

func (e *executorBridge) GetDefaultInputs() []common.Input {
	return toCommonInputs(e.engine.GetDefaultInputs())
}

func (e *executorBridge) GetPrerequisites() []common.Input {
	return toCommonInputs(e.engine.GetPrerequisites())
}

func (e *executorBridge) HasRequiredInputs(ctx *flowcore.NodeContext, execResp *common.ExecutorResponse) bool {
	return e.engine.HasRequiredInputs(toEngineNodeContext(ctx), toEngineExecutorResponse(execResp))
}

func (e *executorBridge) ValidatePrerequisites(ctx *flowcore.NodeContext, execResp *common.ExecutorResponse) bool {
	return e.engine.ValidatePrerequisites(toEngineNodeContext(ctx), toEngineExecutorResponse(execResp))
}

func (e *executorBridge) GetUserIDFromContext(ctx *flowcore.NodeContext) string {
	return e.engine.GetUserIDFromContext(toEngineNodeContext(ctx))
}

func (e *executorBridge) GetRequiredInputs(ctx *flowcore.NodeContext) []common.Input {
	return toCommonInputs(e.engine.GetRequiredInputs(toEngineNodeContext(ctx)))
}

func (e *executorBridge) GetExecutionPolicy(mode string) *flowcore.ExecutionPolicy {
	policy := e.engine.GetExecutionPolicy(mode)
	if policy == nil {
		return nil
	}
	return &flowcore.ExecutionPolicy{
		SkipChallengeValidation: policy.SkipChallengeValidation,
		AllowSegmentRestart:     policy.AllowSegmentRestart,
	}
}

type internalExecutorBridge struct {
	inner flowcore.ExecutorInterface
}

func (e *internalExecutorBridge) Execute(ctx *thunderidengine.NodeContext) (*thunderidengine.ExecutorResponse, error) {
	resp, err := e.inner.Execute(toCoreNodeContext(ctx))
	if err != nil {
		return nil, err
	}
	return toEngineExecutorResponse(resp), nil
}

func (e *internalExecutorBridge) GetName() string {
	return e.inner.GetName()
}

func (e *internalExecutorBridge) GetType() thunderidengine.ExecutorType {
	return thunderidengine.ExecutorType(e.inner.GetType())
}

func (e *internalExecutorBridge) GetDefaultInputs() []thunderidengine.Input {
	return toEngineInputs(e.inner.GetDefaultInputs())
}

func (e *internalExecutorBridge) GetPrerequisites() []thunderidengine.Input {
	return toEngineInputs(e.inner.GetPrerequisites())
}

func (e *internalExecutorBridge) HasRequiredInputs(ctx *thunderidengine.NodeContext, execResp *thunderidengine.ExecutorResponse) bool {
	return e.inner.HasRequiredInputs(toCoreNodeContext(ctx), toCommonExecutorResponse(execResp))
}

func (e *internalExecutorBridge) ValidatePrerequisites(ctx *thunderidengine.NodeContext, execResp *thunderidengine.ExecutorResponse) bool {
	return e.inner.ValidatePrerequisites(toCoreNodeContext(ctx), toCommonExecutorResponse(execResp))
}

func (e *internalExecutorBridge) GetUserIDFromContext(ctx *thunderidengine.NodeContext) string {
	return e.inner.GetUserIDFromContext(toCoreNodeContext(ctx))
}

func (e *internalExecutorBridge) GetRequiredInputs(ctx *thunderidengine.NodeContext) []thunderidengine.Input {
	return toEngineInputs(e.inner.GetRequiredInputs(toCoreNodeContext(ctx)))
}

func (e *internalExecutorBridge) GetExecutionPolicy(mode string) *thunderidengine.ExecutionPolicy {
	policy := e.inner.GetExecutionPolicy(mode)
	if policy == nil {
		return nil
	}
	return &thunderidengine.ExecutionPolicy{
		SkipChallengeValidation: policy.SkipChallengeValidation,
		AllowSegmentRestart:     policy.AllowSegmentRestart,
	}
}

func toCoreNodeContext(ctx *thunderidengine.NodeContext) *flowcore.NodeContext {
	if ctx == nil {
		return nil
	}
	return &flowcore.NodeContext{
		Context:           ctx.Context,
		ExecutionID:       ctx.ExecutionID,
		FlowType:          common.FlowType(ctx.FlowType),
		EntityID:          ctx.EntityID,
		Verbose:           ctx.Verbose,
		CurrentAction:     ctx.CurrentAction,
		CurrentNodeID:     ctx.CurrentNodeID,
		ExecutorMode:      ctx.ExecutorMode,
		NodeProperties:    ctx.NodeProperties,
		NodeInputs:        toCommonInputs(ctx.NodeInputs),
		UserInputs:        ctx.UserInputs,
		RuntimeData:       ctx.RuntimeData,
		ForwardedData:     ctx.ForwardedData,
		Application:       toAppModelApplication(ctx.Application),
		AuthenticatedUser: toAuthnAuthenticatedUser(ctx.AuthenticatedUser),
		AuthUser:          authnprovidermgr.AuthUser{},
		ExecutionHistory:  toCommonExecutionHistory(ctx.ExecutionHistory),
	}
}

func toEngineNodeContext(ctx *flowcore.NodeContext) *thunderidengine.NodeContext {
	if ctx == nil {
		return nil
	}
	return &thunderidengine.NodeContext{
		Context:           ctx.Context,
		ExecutionID:       ctx.ExecutionID,
		FlowType:          thunderidengine.FlowType(ctx.FlowType),
		EntityID:          ctx.EntityID,
		Verbose:           ctx.Verbose,
		CurrentAction:     ctx.CurrentAction,
		CurrentNodeID:     ctx.CurrentNodeID,
		ExecutorMode:      ctx.ExecutorMode,
		NodeProperties:    ctx.NodeProperties,
		NodeInputs:        toEngineInputs(ctx.NodeInputs),
		UserInputs:        ctx.UserInputs,
		RuntimeData:       ctx.RuntimeData,
		ForwardedData:     ctx.ForwardedData,
		Application:       toEngineApplication(ctx.Application),
		AuthenticatedUser: toEngineAuthenticatedUser(ctx.AuthenticatedUser),
		AuthUser:          thunderidengine.AuthUser{},
		ExecutionHistory:  toEngineExecutionHistory(ctx.ExecutionHistory),
	}
}

func toAppModelApplication(app thunderidengine.Application) appmodel.Application {
	return appmodel.Application{
		ID:          app.ID,
		OUID:        app.OUID,
		Name:        app.Name,
		Description: app.Description,
		URL:         app.URL,
		LogoURL:     app.LogoURL,
		TosURI:      app.TosURI,
		PolicyURI:   app.PolicyURI,
	}
}

func toEngineApplication(app appmodel.Application) thunderidengine.Application {
	return thunderidengine.Application{
		ID:                        app.ID,
		Name:                      app.Name,
		Description:               app.Description,
		OUID:                      app.OUID,
		LogoURL:                   app.LogoURL,
		URL:                       app.URL,
		TosURI:                    app.TosURI,
		PolicyURI:                 app.PolicyURI,
		IsRegistrationFlowEnabled: app.IsRegistrationFlowEnabled,
		IsRecoveryFlowEnabled:     app.IsRecoveryFlowEnabled,
		Properties:                app.Metadata,
		AuthFlowID:                app.AuthFlowID,
		RegistrationFlowID:        app.RegistrationFlowID,
		RecoveryFlowID:            app.RecoveryFlowID,
	}
}

func toAuthnAuthenticatedUser(user thunderidengine.AuthenticatedUser) authncm.AuthenticatedUser {
	return authncm.AuthenticatedUser{
		IsAuthenticated:     user.IsAuthenticated,
		UserID:              user.UserID,
		OUID:                user.OUID,
		UserType:            user.UserType,
		Attributes:          user.Attributes,
		AvailableAttributes: toAuthnAttributesResponse(user.AvailableAttributes),
		Token:               user.Token,
	}
}

func toEngineAuthenticatedUser(user authncm.AuthenticatedUser) thunderidengine.AuthenticatedUser {
	return thunderidengine.AuthenticatedUser{
		IsAuthenticated:     user.IsAuthenticated,
		UserID:              user.UserID,
		OUID:                user.OUID,
		UserType:            user.UserType,
		Attributes:          user.Attributes,
		AvailableAttributes: toEngineAttributesResponse(user.AvailableAttributes),
		Token:               user.Token,
	}
}

func toAuthnAttributesResponse(attrs *thunderidengine.AttributesResponse) *authnprovidercm.AttributesResponse {
	if attrs == nil || attrs.Attributes == nil {
		return nil
	}
	return nil
}

func toEngineAttributesResponse(attrs *authnprovidercm.AttributesResponse) *thunderidengine.AttributesResponse {
	if attrs == nil || attrs.Attributes == nil {
		return nil
	}
	out := make(map[string]interface{}, len(attrs.Attributes))
	for k, v := range attrs.Attributes {
		if v != nil {
			out[k] = v.Value
		}
	}
	return &thunderidengine.AttributesResponse{Attributes: out}
}

func toCommonExecutionHistory(history map[string]*thunderidengine.NodeExecutionRecord) map[string]*common.NodeExecutionRecord {
	if history == nil {
		return nil
	}
	out := make(map[string]*common.NodeExecutionRecord, len(history))
	for k, v := range history {
		if v == nil {
			continue
		}
		out[k] = &common.NodeExecutionRecord{
			NodeID:       v.NodeID,
			NodeType:     v.NodeType,
			ExecutorName: v.ExecutorName,
			ExecutorType: common.ExecutorType(v.ExecutorType),
			ExecutorMode: v.ExecutorMode,
			Step:         v.Step,
			Status:       common.FlowStatus(v.Status),
			Executions:   toCommonExecutionAttempts(v.Executions),
			StartTime:    v.StartTime,
			EndTime:      v.EndTime,
		}
	}
	return out
}

func toEngineExecutionHistory(history map[string]*common.NodeExecutionRecord) map[string]*thunderidengine.NodeExecutionRecord {
	if history == nil {
		return nil
	}
	out := make(map[string]*thunderidengine.NodeExecutionRecord, len(history))
	for k, v := range history {
		if v == nil {
			continue
		}
		out[k] = &thunderidengine.NodeExecutionRecord{
			NodeID:       v.NodeID,
			NodeType:     v.NodeType,
			ExecutorName: v.ExecutorName,
			ExecutorType: thunderidengine.ExecutorType(v.ExecutorType),
			ExecutorMode: v.ExecutorMode,
			Step:         v.Step,
			Status:       thunderidengine.FlowStatus(v.Status),
			Executions:   toEngineExecutionAttempts(v.Executions),
			StartTime:    v.StartTime,
			EndTime:      v.EndTime,
		}
	}
	return out
}

func toCommonExecutionAttempts(attempts []thunderidengine.ExecutionAttempt) []common.ExecutionAttempt {
	out := make([]common.ExecutionAttempt, len(attempts))
	for i, a := range attempts {
		out[i] = common.ExecutionAttempt{
			Attempt:   a.Attempt,
			Timestamp: a.Timestamp,
			Status:    common.FlowStatus(a.Status),
			StartTime: a.StartTime,
			EndTime:   a.EndTime,
		}
	}
	return out
}

func toEngineExecutionAttempts(attempts []common.ExecutionAttempt) []thunderidengine.ExecutionAttempt {
	out := make([]thunderidengine.ExecutionAttempt, len(attempts))
	for i, a := range attempts {
		out[i] = thunderidengine.ExecutionAttempt{
			Attempt:   a.Attempt,
			Timestamp: a.Timestamp,
			Status:    thunderidengine.FlowStatus(a.Status),
			StartTime: a.StartTime,
			EndTime:   a.EndTime,
		}
	}
	return out
}

func toEngineInputs(inputs []common.Input) []thunderidengine.Input {
	out := make([]thunderidengine.Input, len(inputs))
	for i, in := range inputs {
		out[i] = thunderidengine.Input{
			Ref:         in.Ref,
			Identifier:  in.Identifier,
			Type:        in.Type,
			Required:    in.Required,
			Options:     in.Options,
			DisplayName: in.DisplayName,
		}
	}
	return out
}

func toCommonInputs(inputs []thunderidengine.Input) []common.Input {
	out := make([]common.Input, len(inputs))
	for i, in := range inputs {
		out[i] = common.Input{
			Ref:         in.Ref,
			Identifier:  in.Identifier,
			Type:        in.Type,
			Required:    in.Required,
			Options:     in.Options,
			DisplayName: in.DisplayName,
		}
	}
	return out
}

func toEngineExecutorResponse(resp *common.ExecutorResponse) *thunderidengine.ExecutorResponse {
	if resp == nil {
		return nil
	}
	return &thunderidengine.ExecutorResponse{
		Status:            thunderidengine.ExecutorStatus(resp.Status),
		Inputs:            toEngineInputs(resp.Inputs),
		AdditionalData:    resp.AdditionalData,
		RedirectURL:       resp.RedirectURL,
		RuntimeData:       resp.RuntimeData,
		ForwardedData:     resp.ForwardedData,
		AuthenticatedUser: toEngineAuthenticatedUser(resp.AuthenticatedUser),
		Assertion:         resp.Assertion,
		FailureReason:     resp.FailureReason,
		AuthUser:          thunderidengine.AuthUser{},
	}
}

func toCommonExecutorResponse(resp *thunderidengine.ExecutorResponse) *common.ExecutorResponse {
	if resp == nil {
		return nil
	}
	return &common.ExecutorResponse{
		Status:            common.ExecutorStatus(resp.Status),
		Inputs:            toCommonInputs(resp.Inputs),
		AdditionalData:    resp.AdditionalData,
		RedirectURL:       resp.RedirectURL,
		RuntimeData:       resp.RuntimeData,
		ForwardedData:     resp.ForwardedData,
		AuthenticatedUser: toAuthnAuthenticatedUser(resp.AuthenticatedUser),
		Assertion:         resp.Assertion,
		FailureReason:     resp.FailureReason,
		AuthUser:          authnprovidermgr.AuthUser{},
	}
}
