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

// Package host defines provider interfaces implemented by ThunderID engine hosts.
package host

import (
	"context"
	"encoding/json"
)

// Actor represents an identity entity in the host directory.
type Actor struct {
	ID         string
	EntityType string
	Attributes json.RawMessage
}

// EntityType describes the schema for an actor type.
type EntityType struct {
	ID         string
	Name       string
	Attributes json.RawMessage
}

// Application represents an OAuth/OIDC application registration.
type Application struct {
	ID           string
	Name         string
	OUID         string
	EntityID     string
	RedirectURIs []string
}

// InboundClient represents an OAuth inbound client profile.
type InboundClient struct {
	ClientID                           string
	EntityID                           string
	ApplicationID                      string
	OUID                               string
	Secret                             string
	GrantTypes                         []string
	RedirectURIs                       []string
	ResponseTypes                      []string
	TokenEndpointAuthMethod            string
	PKCERequired                       bool
	PublicClient                       bool
	RequirePushedAuthorizationRequests bool

	AuthFlowID                string
	RegistrationFlowID        string
	IsRegistrationFlowEnabled bool
	RecoveryFlowID            string
	IsRecoveryFlowEnabled     bool
}

// ActorProvider supplies identity, application, and inbound client data to the engine.
type ActorProvider interface {
	IdentifyEntity(filters map[string]interface{}) (*string, error)
	GetEntity(entityID string) (*Actor, error)
	SearchEntities(filters map[string]interface{}) ([]*Actor, error)
	GetApplication(ctx context.Context, appID string) (*Application, error)
	GetInboundClientByEntityID(ctx context.Context, entityID string) (*InboundClient, error)
	GetInboundClientByClientID(ctx context.Context, clientID string) (*InboundClient, error)
	GetEntityType(ctx context.Context, typeID string) (*EntityType, error)
}
