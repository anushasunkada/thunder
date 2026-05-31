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

// Package flow exposes flow executor types for host applications that register custom executors.
package flow

import (
	"github.com/thunder-id/thunderid/internal/flow/common"
	"github.com/thunder-id/thunderid/internal/flow/core"
)

type (
	// Executor is a flow node executor implementation.
	Executor = core.ExecutorInterface
	// NodeContext carries flow execution state for an executor invocation.
	NodeContext = core.NodeContext
	// Input describes a user or system input required by an executor.
	Input = common.Input
	// ExecutorResponse is the result returned from an executor run.
	ExecutorResponse = common.ExecutorResponse
	// ExecutorType categorizes an executor (authentication, utility, etc.).
	ExecutorType = common.ExecutorType
	// ExecutorStatus reports completion state for an executor response.
	ExecutorStatus = common.ExecutorStatus
)

const (
	// ExecComplete indicates the executor finished successfully.
	ExecComplete = common.ExecComplete
	// ExecFailure indicates the executor failed.
	ExecFailure = common.ExecFailure
	// ExecUserInputRequired indicates the executor needs additional user input.
	ExecUserInputRequired = common.ExecUserInputRequired
	// ExecutorTypeAuthentication marks authentication executors.
	ExecutorTypeAuthentication = common.ExecutorTypeAuthentication
	// ExecutorTypeUtility marks utility executors.
	ExecutorTypeUtility = common.ExecutorTypeUtility
	// InputTypeText marks a text input field.
	InputTypeText = common.InputTypeText
	// InputTypeOTP marks a one-time password input field.
	InputTypeOTP = common.InputTypeOTP
	// InputTypePhone marks a phone number input field.
	InputTypePhone = common.InputTypePhone
)
