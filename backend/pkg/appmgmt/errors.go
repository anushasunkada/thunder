/*
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
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

package appmgmt

import "errors"

// Sentinel errors returned by application services. Callers should use
// errors.Is for comparison.
var (
	// ErrNotFound is returned when a requested resource does not exist.
	ErrNotFound = errors.New("application: resource not found")

	// ErrConflict is returned when an operation cannot be completed due to a
	// conflict with the current state of the target resource (e.g. duplicate
	// name).
	ErrConflict = errors.New("application: resource conflict")

	// ErrInvalidInput is returned when the caller provides invalid input data.
	ErrInvalidInput = errors.New("application: invalid input")
)
