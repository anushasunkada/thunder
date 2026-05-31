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
	"encoding/json"
	"errors"
	"time"

	"github.com/thunder-id/thunderid/internal/flow/common"
)

// ErrNotFound indicates a runtime store entry was not found or expired.
var ErrNotFound = errors.New("runtime store entry not found")

// RuntimeStore holds ephemeral engine state.
type RuntimeStore interface {
	StoreFlowContext(ctx context.Context, executionID string, data []byte, expiry time.Time) error
	GetFlowContext(ctx context.Context, executionID string) ([]byte, error)
	UpdateFlowContext(ctx context.Context, executionID string, data []byte) error
	DeleteFlowContext(ctx context.Context, executionID string) error
	StoreAuthCode(ctx context.Context, code string, data []byte, expiry time.Time) error
	GetAuthCode(ctx context.Context, code string) ([]byte, error)
	DeleteAuthCode(ctx context.Context, code string) error
	StoreAuthRequest(ctx context.Context, requestID string, data []byte, expiry time.Time) error
	GetAuthRequest(ctx context.Context, requestID string) ([]byte, error)
	DeleteAuthRequest(ctx context.Context, requestID string) error
	StorePAR(ctx context.Context, requestURI string, data []byte, expiry time.Time) error
	GetPAR(ctx context.Context, requestURI string) ([]byte, error)
	DeletePAR(ctx context.Context, requestURI string) error
	StoreJTI(ctx context.Context, jti string, expiry time.Time) error
	ExistsJTI(ctx context.Context, jti string) (bool, error)

	StoreAttributeCache(ctx context.Context, id string, data []byte, expiry time.Time) error
	GetAttributeCache(ctx context.Context, id string) ([]byte, error)
	ExtendAttributeCacheExpiry(ctx context.Context, id string, expiry time.Time) error
	DeleteAttributeCache(ctx context.Context, id string) error
}

// ActorSource supplies identity and client configuration.
type ActorSource interface {
	IdentifyEntity(filters map[string]interface{}) (*string, error)
	GetEntity(entityID string) (*Actor, error)
	SearchEntities(filters map[string]interface{}) ([]*Actor, error)
	GetApplication(ctx context.Context, appID string) (*Application, error)
	GetInboundClientByEntityID(ctx context.Context, entityID string) (*InboundClient, error)
	GetInboundClientByClientID(ctx context.Context, clientID string) (*InboundClient, error)
	GetEntityType(ctx context.Context, typeID string) (*EntityType, error)
}

// Actor represents an identity entity.
type Actor struct {
	ID         string
	EntityType string
	Attributes json.RawMessage
}

// Application represents an OAuth/OIDC application.
type Application struct {
	ID       string
	Name     string
	OUID     string
	EntityID string
}

// InboundClient represents an OAuth inbound client profile.
type InboundClient struct {
	ClientID      string
	EntityID      string
	ApplicationID string
	OUID          string
	Secret        string
	GrantTypes    []string
	RedirectURIs  []string
	ResponseTypes []string
	TokenEndpointAuthMethod string
	PKCERequired  bool
	PublicClient  bool
	RequirePushedAuthorizationRequests bool

	AuthFlowID                string
	RegistrationFlowID        string
	IsRegistrationFlowEnabled bool
	RecoveryFlowID            string
	IsRecoveryFlowEnabled     bool
}

// EntityType describes the schema for an actor type.
type EntityType struct {
	ID         string
	Name       string
	Attributes json.RawMessage
}

// AuthorizationSource checks permissions.
type AuthorizationSource interface {
	GetAuthorizedPermissions(ctx context.Context, req AuthorizationRequest) (*AuthorizationResponse, error)
}

// ConsentSource resolves and records consent.
type ConsentSource interface {
	ResolveConsent(
		ctx context.Context, ouID, appID, agentID, userID string, requestedScopes []string,
	) (*ConsentResolution, error)
	RecordConsent(ctx context.Context, ouID, appID, userID string, decisions []ConsentDecision) error
}

// FlowDefinition carries a flow graph for engine execution.
type FlowDefinition struct {
	ID       string
	Handle   string
	Name     string
	FlowType string
	Nodes    []common.NodeDefinition
}

// FlowSource supplies flow definitions.
type FlowSource interface {
	GetFlow(ctx context.Context, flowID string) (*FlowDefinition, error)
	GetFlowByHandle(ctx context.Context, handle, flowType string) (*FlowDefinition, error)
}

// AuthorizationRequest describes a permission check request.
type AuthorizationRequest struct {
	EntityID             string
	RequestedPermissions []string
}

// AuthorizationResponse lists authorized permissions.
type AuthorizationResponse struct {
	AuthorizedPermissions []string
}

// ConsentItem describes a consent prompt item.
type ConsentItem struct {
	ID          string
	DisplayName string
	Description string
	Required    bool
}

// ConsentResolution describes whether consent is required.
type ConsentResolution struct {
	Required bool
	Items    []ConsentItem
}

// ConsentDecision records a user's consent choice.
type ConsentDecision struct {
	ID      string
	Granted bool
}

// AuthnSource verifies credentials.
type AuthnSource interface {
	Authenticate(
		ctx context.Context, identifiers, credentials map[string]interface{}, appID, ouID string,
	) (userID, token string, authenticated bool, err error)
	GetAttributes(
		ctx context.Context, token string, requested []string, appID, ouID string,
	) (json.RawMessage, error)
}
