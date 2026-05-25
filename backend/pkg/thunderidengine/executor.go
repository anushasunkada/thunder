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

package thunderidengine

import (
	"context"
	"slices"
)

// AuthenticatedUser represents an authenticated user in flow execution.
type AuthenticatedUser struct {
	IsAuthenticated     bool
	UserID              string
	OUID                string
	UserType            string
	Attributes          map[string]interface{}
	AvailableAttributes *AttributesResponse
	Token               string
}

// AttributesResponse holds attribute metadata for a user.
type AttributesResponse struct {
	Attributes map[string]interface{}
}

// AuthUser accumulates per-provider authentication state during flow execution.
type AuthUser struct {
	UserID     string
	Attributes map[string]interface{}
	Token      string
}

// Input represents inputs required for a flow step.
type Input struct {
	Ref         string   `json:"ref,omitempty"`
	Identifier  string   `json:"identifier"`
	Type        string   `json:"type"`
	Required    bool     `json:"required"`
	Options     []string `json:"options,omitempty"`
	DisplayName string   `json:"-"`
}

// IsSensitive reports whether this input type is sensitive.
func (i Input) IsSensitive() bool {
	return slices.Contains(sensitiveInputTypes, i.Type)
}

// Action is an action on a flow step.
type Action struct {
	Ref      string `json:"ref,omitempty"`
	Type     string `json:"type,omitempty"`
	NextNode string `json:"nextNode,omitempty"`
}

// Prompt groups inputs with an action for prompt nodes.
type Prompt struct {
	Inputs []Input `json:"inputs,omitempty"`
	Action *Action `json:"action,omitempty"`
}

// ExecutorResponse is returned from ExecutorInterface.Execute.
type ExecutorResponse struct {
	Status            ExecutorStatus         `json:"status"`
	Inputs            []Input                `json:"inputs,omitempty"`
	AdditionalData    map[string]string      `json:"additionalData,omitempty"`
	RedirectURL       string                 `json:"redirectUrl,omitempty"`
	RuntimeData       map[string]string      `json:"runtimeData,omitempty"`
	ForwardedData     map[string]interface{} `json:"forwardedData,omitempty"`
	AuthenticatedUser AuthenticatedUser      `json:"authenticatedUser,omitempty"`
	Assertion         string                 `json:"assertion,omitempty"`
	FailureReason     string                 `json:"failureReason,omitempty"`
	AuthUser          AuthUser               `json:"-"`
}

// NodeExecutionRecord records a node execution in the flow.
type NodeExecutionRecord struct {
	NodeID       string             `json:"nodeId"`
	NodeType     string             `json:"nodeType"`
	ExecutorName string             `json:"executorName,omitempty"`
	ExecutorType ExecutorType       `json:"executorType,omitempty"`
	ExecutorMode string             `json:"executorMode,omitempty"`
	Step         int                `json:"step"`
	Status       FlowStatus         `json:"status"`
	Executions   []ExecutionAttempt `json:"executions"`
	StartTime    int64              `json:"startTime,omitempty"`
	EndTime      int64              `json:"endTime,omitempty"`
}

// ExecutionAttempt is a single node execution attempt.
type ExecutionAttempt struct {
	Attempt   int        `json:"attempt"`
	Timestamp int64      `json:"timestamp"`
	Status    FlowStatus `json:"status"`
	StartTime int64      `json:"startTime"`
	EndTime   int64      `json:"endTime"`
}

// NodeContext is passed to ExecutorInterface.Execute.
type NodeContext struct {
	Context context.Context

	ExecutionID   string
	FlowType      FlowType
	EntityID      string
	Verbose       bool
	CurrentAction string
	CurrentNodeID string
	ExecutorMode  string

	NodeProperties map[string]interface{}
	NodeInputs     []Input
	UserInputs     map[string]string
	RuntimeData    map[string]string
	ForwardedData  map[string]interface{}

	Application       Application
	AuthenticatedUser AuthenticatedUser
	AuthUser          AuthUser
	ExecutionHistory  map[string]*NodeExecutionRecord
}

// ExecutionPolicy configures executor behavior.
type ExecutionPolicy struct {
	SkipChallengeValidation bool
	AllowSegmentRestart     bool
}

// ExecutorInterface is implemented by custom flow step executors.
type ExecutorInterface interface {
	Execute(ctx *NodeContext) (*ExecutorResponse, error)
	GetName() string
	GetType() ExecutorType
	GetDefaultInputs() []Input
	GetPrerequisites() []Input
	HasRequiredInputs(ctx *NodeContext, execResp *ExecutorResponse) bool
	ValidatePrerequisites(ctx *NodeContext, execResp *ExecutorResponse) bool
	GetUserIDFromContext(ctx *NodeContext) string
	GetRequiredInputs(ctx *NodeContext) []Input
	GetExecutionPolicy(mode string) *ExecutionPolicy
}
