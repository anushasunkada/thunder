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

import (
	"context"
	"encoding/json"
)

// AuthnMetadata carries contextual information for authentication.
type AuthnMetadata struct {
	ApplicationID string
	OUID          string
	FlowID        string
	// AppMetadata is the full application metadata map from the flow (tenant_id, client_ids, oauth_client_id, etc.).
	AppMetadata map[string]interface{}
	// RuntimeMetadata is the full runtime metadata map from the flow (custom_metadata, etc.).
	RuntimeMetadata map[string]interface{}
}

// AuthnResult is returned after successful authentication.
type AuthnResult struct {
	Authenticated   bool
	UserID          string
	AuthToken       string
	Attributes      json.RawMessage
	AuthenticatorID string
}

// RequestedAttributes describes attributes requested after authentication.
type RequestedAttributes struct {
	AttributeNames []string
	Scopes         []string
}

// GetAttributesMetadata carries context for attribute retrieval.
type GetAttributesMetadata struct {
	ApplicationID string
	OUID          string
	Locale        string
	// AppMetadata is the full application metadata map from the flow (tenant_id, client_ids, oauth_client_id, etc.).
	AppMetadata map[string]interface{}
	// RuntimeMetadata is the full runtime metadata map from the flow (custom_metadata, etc.).
	RuntimeMetadata map[string]interface{}
}

// GetAttributesResult is returned from attribute retrieval.
type GetAttributesResult struct {
	Attributes json.RawMessage
}

// AuthnProvider verifies credentials and retrieves user attributes.
type AuthnProvider interface {
	Authenticate(ctx context.Context, identifiers, credentials map[string]interface{},
		metadata *AuthnMetadata) (*AuthnResult, error)
	GetAttributes(ctx context.Context, token string, requested *RequestedAttributes,
		metadata *GetAttributesMetadata) (*GetAttributesResult, error)
}

// AuthnProviderManager is an opaque handle produced by InitializeAuthnProviderManager.
type AuthnProviderManager interface {
	isAuthnProviderManager()
}
