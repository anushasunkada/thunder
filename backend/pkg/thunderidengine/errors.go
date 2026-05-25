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

import "errors"

var (
	// ErrInvalidConfig is returned when engine configuration is incomplete or invalid.
	ErrInvalidConfig = errors.New("thunderidengine: invalid config")
	// ErrMissingHostField is returned when a required host provider is nil.
	ErrMissingHostField = errors.New("thunderidengine: missing required host provider")
	// ErrExecutorExists is returned when registering a duplicate executor.
	ErrExecutorExists = errors.New("thunderidengine: executor already registered")
	// ErrExecutorNotFound is returned when an executor lookup fails.
	ErrExecutorNotFound = errors.New("thunderidengine: executor not found")
	// ErrInvalidFlowType is returned when a flow type is not supported.
	ErrInvalidFlowType = errors.New("thunderidengine: invalid flow type")
	// ErrFlowNotFound is returned when a flow definition cannot be resolved.
	ErrFlowNotFound = errors.New("thunderidengine: flow not found")
	// ErrApplicationNotFound is returned when an application cannot be resolved.
	ErrApplicationNotFound = errors.New("thunderidengine: application not found")
	// ErrInboundClientNotFound is returned when an OAuth client cannot be resolved.
	ErrInboundClientNotFound = errors.New("thunderidengine: inbound client not found")
)

// Error is a runtime error surfaced by the engine and providers.
type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Code
}
