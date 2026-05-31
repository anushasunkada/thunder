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

package host

import "context"

// ConsentItem describes a consent prompt item.
type ConsentItem struct {
	ID          string
	DisplayName string
	Description string
	Required    bool
}

// ConsentResolution describes whether consent is required and which items to prompt.
type ConsentResolution struct {
	Required bool
	Items    []ConsentItem
}

// ConsentDecision records a user's consent choice.
type ConsentDecision struct {
	ID      string
	Granted bool
}

// ConsentEnforcer resolves and records user consent for applications.
type ConsentEnforcer interface {
	ResolveConsent(ctx context.Context, ouID, appID, agentID, userID string,
		requestedScopes []string) (*ConsentResolution, error)
	RecordConsent(ctx context.Context, ouID, appID, userID string,
		decisions []ConsentDecision) error
}
