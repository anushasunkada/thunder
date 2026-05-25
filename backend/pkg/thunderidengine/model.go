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

// Package thunderidengine provides shared runtime models for the embeddable engine.
//
//nolint:lll,revive // struct tags mirror inboundclient jsonschema descriptions; enum consts omit per-value docs.
package thunderidengine

import (
	"encoding/json"
	"time"

	"github.com/thunder-id/thunderid/internal/cert"
)

// OAuthClient is the runtime view of an OAuth inbound client.
type OAuthClient struct {
	ClientID                           string
	EntityID                           string
	OUID                               string
	Secret                             string
	RedirectURIs                       []string
	GrantTypes                         []GrantType
	ResponseTypes                      []ResponseType
	Scopes                             []string
	AuthFlowID                         string
	RegistrationFlowID                 string
	RecoveryFlowID                     string
	IsRegistrationFlowEnabled          bool
	IsRecoveryFlowEnabled              bool
	Properties                         map[string]interface{}
	TokenEndpointAuthMethod            TokenEndpointAuthMethod
	PKCERequired                       bool
	PublicClient                       bool
	RequirePushedAuthorizationRequests bool
	Token                              *OAuthTokenConfig
	UserInfo                           *UserInfoConfig
	ScopeClaims                        map[string][]string
	AcrValues                          []string
	AccessTokenValiditySeconds         int64
	IDTokenValiditySeconds             int64
	Assertion                          *AssertionConfig
	LoginConsent                       *LoginConsentConfig
	AllowedUserTypes                   []string
	Certificate                        *Certificate
}

// AssertionConfig is the entity-level assertion config; token configs fall back to it.
type AssertionConfig struct {
	ValidityPeriod int64    `json:"validityPeriod,omitempty" yaml:"validity_period,omitempty" jsonschema:"Assertion validity period in seconds."`
	UserAttributes []string `json:"userAttributes,omitempty" yaml:"user_attributes,omitempty" jsonschema:"User attributes to include in the assertion."`
}

// LoginConsentConfig is the login consent configuration.
type LoginConsentConfig struct {
	ValidityPeriod int64 `json:"validityPeriod" yaml:"validity_period" jsonschema:"Consent validity period in seconds. 0 means never expire."`
}

// Certificate is a user-supplied certificate input.
type Certificate struct {
	Type  cert.CertificateType `json:"type,omitempty"  yaml:"type,omitempty"  jsonschema:"Certificate type (PEM, JWK, etc.)."`
	Value string               `json:"value,omitempty" yaml:"value,omitempty" jsonschema:"Certificate value in the format specified by type."`
}

// Application is the runtime view of an application (or agent) entity.
type Application struct {
	ID                        string
	Name                      string
	Description               string
	OUID                      string
	LogoURL                   string
	URL                       string
	TosURI                    string
	PolicyURI                 string
	IsRegistrationFlowEnabled bool
	IsRecoveryFlowEnabled     bool
	Properties                map[string]interface{}
	AuthFlowID                string
	RegistrationFlowID        string
	RecoveryFlowID            string
}

// EntityGroup represents group membership for an entity.
type EntityGroup struct {
	ID   string
	Name string
}

// AuthResult is returned from AuthenticateUser.
type AuthResult struct {
	UserID     string
	Attributes map[string]interface{}
}

// Attributes holds requested user attributes.
type Attributes struct {
	Values map[string]interface{}
}

// Resource is a protected resource indicator target.
type Resource struct {
	ID  string
	URI string
}

// OU is an organizational unit.
type OU struct {
	ID              string
	Handle          string
	Name            string
	Description     string
	LogoURL         string
	TosURI          string
	PolicyURI       string
	CookiePolicyURI string
}

// IDP is an identity provider configuration.
type IDP struct {
	ID     string
	Name   string
	Type   string
	Config json.RawMessage
}

// FlowDefinition is a complete flow definition used to build an execution graph.
type FlowDefinition struct {
	ID       string
	Handle   string
	FlowType string
	Nodes    []FlowNodeDefinition
	Edges    []FlowEdgeDefinition
}

// FlowNodeDefinition describes a node in a flow graph.
type FlowNodeDefinition struct {
	ID         string
	Type       string
	Properties map[string]interface{}
}

// FlowEdgeDefinition describes an edge in a flow graph.
type FlowEdgeDefinition struct {
	Source    string
	Target    string
	Condition string
	SegmentID string
}

// ObservabilityEvent is a runtime observability payload.
type ObservabilityEvent struct {
	Name       string
	Properties map[string]string
}

// PARRequest is a pushed authorization request payload.
type PARRequest struct {
	ClientID          string
	OAuthParameters   OAuthParameters
}

// AuthRequestContext holds in-progress authorize request state.
type AuthRequestContext struct {
	OAuthParameters OAuthParameters
}

// AuthorizationCode is an OAuth authorization code.
type AuthorizationCode struct {
	CodeID              string
	Code                string
	ClientID            string
	RedirectURI         string
	AuthorizedUserID    string
	AttributeCacheID    string
	TimeCreated         time.Time
	ExpiryTime          time.Time
	Scopes              string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
	Resources           []string
	ClaimsRequest       *ClaimsRequest
	ClaimsLocales       string
	Nonce               string
	CompletedACR        string
}

// FlowContextDB is persisted flow execution state.
type FlowContextDB struct {
	ExecutionID string
	Context     string
	ExpiryTime  time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// FlowInitContext starts a new flow execution.
type FlowInitContext struct {
	ApplicationID string
	FlowType      string
	RuntimeData   map[string]string
}

// ConsentResolveRequest is input for ResolveConsent.
type ConsentResolveRequest struct {
	OUID                  string
	AppID                 string
	AppName               string
	UserID                string
	EssentialAttributes   []string
	OptionalAttributes    []string
	AuthorizedPermissions []string
}

// ConsentPromptData describes missing consents.
type ConsentPromptData struct {
	SessionToken string
	Purposes     []ConsentPurpose
}

// ConsentPurpose is a consent purpose shown to the user.
type ConsentPurpose struct {
	ID          string
	DisplayName string
}

// ConsentRecordRequest is input for RecordConsent.
type ConsentRecordRequest struct {
	OUID           string
	AppID          string
	UserID         string
	SessionToken   string
	ValidityPeriod int64
	Decisions      map[string]bool
}

// Consent is a recorded consent record.
type Consent struct {
	ID string
}

// DesignResolveType identifies design resolution scope.
type DesignResolveType string

const (
	DesignResolveTypeAPP DesignResolveType = "APP"
	DesignResolveTypeOU  DesignResolveType = "OU"
)

// DesignResponse holds resolved theme and layout JSON.
type DesignResponse struct {
	Theme  json.RawMessage
	Layout json.RawMessage
}

// MetaType is the type query parameter for flow metadata.
type MetaType string

const (
	MetaTypeAPP MetaType = "APP"
	MetaTypeOU  MetaType = "OU"
)

// IsValid reports whether mt is a supported flow metadata type.
func (mt MetaType) IsValid() bool {
	return mt == MetaTypeAPP || mt == MetaTypeOU
}

// ApplicationMetadata is the application section of flow metadata.
type ApplicationMetadata struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	LogoURL     string `json:"logoUrl,omitempty"`
	URL         string `json:"url,omitempty"`
	TosURI      string `json:"tosUri,omitempty"`
	PolicyURI   string `json:"policyUri,omitempty"`
}

// OUMetadata is the OU section of flow metadata.
type OUMetadata struct {
	ID              string `json:"id,omitempty"`
	Handle          string `json:"handle,omitempty"`
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	LogoURL         string `json:"logoUrl,omitempty"`
	TosURI          string `json:"tosUri,omitempty"`
	PolicyURI       string `json:"policyUri,omitempty"`
	CookiePolicyURI string `json:"cookiePolicyUri,omitempty"`
}

// DesignMetadata is theme/layout in flow metadata.
type DesignMetadata struct {
	Theme  json.RawMessage `json:"theme"`
	Layout json.RawMessage `json:"layout"`
}

// I18nMetadata is translations in flow metadata.
type I18nMetadata struct {
	Languages    []string                     `json:"languages"`
	Language     string                       `json:"language"`
	TotalResults int                          `json:"totalResults"`
	Translations map[string]map[string]string `json:"translations"`
}
