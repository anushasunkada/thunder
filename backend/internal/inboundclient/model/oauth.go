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

// Package model defines OAuth-related types for inbound client configuration.
//
//nolint:lll
package model

import (
	"github.com/thunder-id/thunderid/internal/system/jose/jwe"
	"github.com/thunder-id/thunderid/internal/system/jose/jws"
	"github.com/thunder-id/thunderid/pkg/thunderidengine"
)

// InboundAuthType identifies the kind of inbound authentication configured for an entity.
type InboundAuthType string

const (
	// OAuthInboundAuthType is the OAuth 2.0 inbound authentication type.
	OAuthInboundAuthType InboundAuthType = "oauth2"
)

type (
	// OAuthTokenConfig wraps access and ID token configs.
	OAuthTokenConfig = thunderidengine.OAuthTokenConfig
	// AccessTokenConfig is the access token configuration.
	AccessTokenConfig = thunderidengine.AccessTokenConfig
	// IDTokenConfig is the ID token configuration.
	IDTokenConfig = thunderidengine.IDTokenConfig
	// UserInfoConfig is the userinfo endpoint configuration.
	UserInfoConfig = thunderidengine.UserInfoConfig
	// GrantType is an OAuth 2.0 grant type.
	GrantType = thunderidengine.GrantType
	// ResponseType is an OAuth 2.0 authorization response type.
	ResponseType = thunderidengine.ResponseType
	// TokenEndpointAuthMethod is a token endpoint client authentication method.
	TokenEndpointAuthMethod = thunderidengine.TokenEndpointAuthMethod
	// IDTokenResponseType is the response format of the ID token.
	IDTokenResponseType = thunderidengine.IDTokenResponseType
	// UserInfoResponseType is the response format of the UserInfo endpoint.
	UserInfoResponseType = thunderidengine.UserInfoResponseType
)

//nolint:revive // re-exported thunderidengine OAuth enum values.
const (
	GrantTypeAuthorizationCode               = thunderidengine.GrantTypeAuthorizationCode
	GrantTypeClientCredentials               = thunderidengine.GrantTypeClientCredentials
	GrantTypeRefreshToken                    = thunderidengine.GrantTypeRefreshToken
	GrantTypeTokenExchange                   = thunderidengine.GrantTypeTokenExchange
	ResponseTypeCode                         = thunderidengine.ResponseTypeCode
	ResponseTypeIDToken                      = thunderidengine.ResponseTypeIDToken
	TokenEndpointAuthMethodClientSecretBasic = thunderidengine.TokenEndpointAuthMethodClientSecretBasic
	TokenEndpointAuthMethodClientSecretPost  = thunderidengine.TokenEndpointAuthMethodClientSecretPost
	TokenEndpointAuthMethodPrivateKeyJWT     = thunderidengine.TokenEndpointAuthMethodPrivateKeyJWT
	TokenEndpointAuthMethodNone              = thunderidengine.TokenEndpointAuthMethodNone
	IDTokenResponseTypeJWT                   = thunderidengine.IDTokenResponseTypeJWT
	IDTokenResponseTypeJWE                   = thunderidengine.IDTokenResponseTypeJWE
	IDTokenResponseTypeNESTEDJWT             = thunderidengine.IDTokenResponseTypeNestedJWT //nolint:gosec
	UserInfoResponseTypeJSON                 = thunderidengine.UserInfoResponseTypeJSON
	UserInfoResponseTypeJWS                  = thunderidengine.UserInfoResponseTypeJWS
	UserInfoResponseTypeJWE                  = thunderidengine.UserInfoResponseTypeJWE
	UserInfoResponseTypeNESTEDJWT            = thunderidengine.UserInfoResponseTypeNestedJWT
)

// Supported JOSE algorithms for userinfo responses.
var (
	SupportedUserInfoSigningAlgs    = []string{string(jws.RS256), string(jws.RS512), string(jws.PS256), string(jws.ES256), string(jws.ES384), string(jws.ES512), string(jws.EdDSA)}
	SupportedUserInfoEncryptionAlgs = []string{string(jwe.RSAOAEP), string(jwe.RSAOAEP256)}
	SupportedUserInfoEncryptionEncs = []string{string(jwe.A128CBCHS256), string(jwe.A256GCM)}
)

// OAuthProfile is the persistence shape (OAUTH_PROFILE JSONB column).
type OAuthProfile struct {
	RedirectURIs                       []string            `json:"redirectUris"`
	GrantTypes                         []string            `json:"grantTypes"`
	ResponseTypes                      []string            `json:"responseTypes"`
	TokenEndpointAuthMethod            string              `json:"tokenEndpointAuthMethod"`
	PKCERequired                       bool                `json:"pkceRequired"`
	PublicClient                       bool                `json:"publicClient"`
	RequirePushedAuthorizationRequests bool                `json:"requirePushedAuthorizationRequests"`
	Token                              *OAuthTokenConfig   `json:"token,omitempty"`
	Scopes                             []string            `json:"scopes,omitempty"`
	UserInfo                           *UserInfoConfig     `json:"userInfo,omitempty"`
	ScopeClaims                        map[string][]string `json:"scopeClaims,omitempty"`
	Certificate                        *Certificate        `json:"certificate,omitempty"`
	AcrValues                          []string            `json:"acrValues,omitempty"`
}

// OAuthConfigWithSecret is the wire input shape and the create/update echo response shape.
// Carries ClientSecret (omitempty) so it appears only when freshly issued.
type OAuthConfigWithSecret struct {
	ClientID                           string                  `json:"clientId,omitempty"                          yaml:"client_id,omitempty"                          jsonschema:"OAuth client ID (auto-generated if not provided)"`
	ClientSecret                       string                  `json:"clientSecret,omitempty"                      yaml:"client_secret,omitempty"                      jsonschema:"OAuth client secret (auto-generated if not provided)"`
	RedirectURIs                       []string                `json:"redirectUris,omitempty"                      yaml:"redirect_uris,omitempty"                      jsonschema:"Allowed redirect URIs. Required for Public (SPA/Mobile) and Confidential (Server) clients. Omit for M2M."`
	GrantTypes                         []GrantType             `json:"grantTypes,omitempty"                        yaml:"grant_types,omitempty"                        jsonschema:"OAuth grant types. Common: [authorization_code, refresh_token] for user apps, [client_credentials] for M2M."`
	ResponseTypes                      []ResponseType          `json:"responseTypes,omitempty"                     yaml:"response_types,omitempty"                     jsonschema:"OAuth response types. Common: [code] for user apps. Omit for M2M."`
	TokenEndpointAuthMethod            TokenEndpointAuthMethod `json:"tokenEndpointAuthMethod,omitempty"           yaml:"token_endpoint_auth_method,omitempty"         jsonschema:"Client authentication method. Use 'none' for Public clients, 'client_secret_basic' for Confidential/M2M."`
	PKCERequired                       bool                    `json:"pkceRequired"                                yaml:"pkce_required"                                jsonschema:"Require PKCE for security. Recommended for all user-interactive flows."`
	PublicClient                       bool                    `json:"publicClient"                                yaml:"public_client"                                jsonschema:"Identify if client is public (cannot store secrets). Set true for SPA/Mobile."`
	RequirePushedAuthorizationRequests bool                    `json:"requirePushedAuthorizationRequests"          yaml:"require_pushed_authorization_requests"        jsonschema:"Require Pushed Authorization Requests (PAR) per RFC 9126."`
	Token                              *OAuthTokenConfig       `json:"token,omitempty"                             yaml:"token,omitempty"                              jsonschema:"Token configuration for access tokens and ID tokens"`
	Scopes                             []string                `json:"scopes,omitempty"                            yaml:"scopes,omitempty"                             jsonschema:"Allowed OAuth scopes. Add custom scopes as needed for your application."`
	UserInfo                           *UserInfoConfig         `json:"userInfo,omitempty"                          yaml:"user_info,omitempty"                          jsonschema:"UserInfo endpoint configuration. Configure user attributes returned from the OIDC userinfo endpoint."`
	ScopeClaims                        map[string][]string     `json:"scopeClaims,omitempty"                       yaml:"scope_claims,omitempty"                       jsonschema:"Scope-to-claims mapping. Maps OAuth scopes to user claims for both ID token and userinfo."`
	Certificate                        *Certificate            `json:"certificate,omitempty"                       yaml:"certificate,omitempty"                        jsonschema:"Application certificate. Optional. For certificate-based authentication or JWT validation."`
	AcrValues                          []string                `json:"acrValues,omitempty"                         yaml:"acr_values,omitempty"                         jsonschema:"Default ACR values applied when the request does not specify acr_values."`
}

// OAuthConfig is the wire output shape (GET responses). ClientSecret is structurally absent.
type OAuthConfig struct {
	ClientID                           string                  `json:"clientId,omitempty"                 yaml:"client_id,omitempty"`
	RedirectURIs                       []string                `json:"redirectUris,omitempty"             yaml:"redirect_uris,omitempty"`
	GrantTypes                         []GrantType             `json:"grantTypes,omitempty"               yaml:"grant_types,omitempty"`
	ResponseTypes                      []ResponseType          `json:"responseTypes,omitempty"            yaml:"response_types,omitempty"`
	TokenEndpointAuthMethod            TokenEndpointAuthMethod `json:"tokenEndpointAuthMethod,omitempty"  yaml:"token_endpoint_auth_method,omitempty"`
	PKCERequired                       bool                    `json:"pkceRequired"                       yaml:"pkce_required"`
	PublicClient                       bool                    `json:"publicClient"                       yaml:"public_client"`
	RequirePushedAuthorizationRequests bool                    `json:"requirePushedAuthorizationRequests" yaml:"require_pushed_authorization_requests"`
	Token                              *OAuthTokenConfig       `json:"token,omitempty"                    yaml:"token,omitempty"`
	Scopes                             []string                `json:"scopes,omitempty"                   yaml:"scopes,omitempty"`
	UserInfo                           *UserInfoConfig         `json:"userInfo,omitempty"                 yaml:"user_info,omitempty"`
	ScopeClaims                        map[string][]string     `json:"scopeClaims,omitempty"              yaml:"scope_claims,omitempty"`
	Certificate                        *Certificate            `json:"certificate,omitempty"              yaml:"certificate,omitempty"`
	AcrValues                          []string                `json:"acrValues,omitempty"                yaml:"acr_values,omitempty"`
}

// SupportedIDTokenEncryptionAlgs lists JWE key-management algorithms supported for ID token encryption.
var SupportedIDTokenEncryptionAlgs = []string{string(jwe.RSAOAEP), string(jwe.RSAOAEP256)}

// SupportedIDTokenEncryptionEncs lists JWE content-encryption algorithms supported for ID token encryption.
var SupportedIDTokenEncryptionEncs = []string{string(jwe.A128CBCHS256), string(jwe.A256GCM)}

// OAuthClient is the inbound-service resolved OAuth profile (entity id in ID). Runtime OAuth behavior
// methods live on thunderidengine.OAuthClient; use adapter.ClientProvider for the merged engine view.
type OAuthClient struct {
	ID                                 string                  `yaml:"id,omitempty"`
	OUID                               string                  `yaml:"ou_id,omitempty"`
	ClientID                           string                  `yaml:"client_id,omitempty"`
	RedirectURIs                       []string                `yaml:"redirect_uris,omitempty"`
	GrantTypes                         []GrantType             `yaml:"grant_types,omitempty"`
	ResponseTypes                      []ResponseType          `yaml:"response_types,omitempty"`
	TokenEndpointAuthMethod            TokenEndpointAuthMethod `yaml:"token_endpoint_auth_method,omitempty"`
	PKCERequired                       bool                    `yaml:"pkce_required,omitempty"`
	PublicClient                       bool                    `yaml:"public_client,omitempty"`
	RequirePushedAuthorizationRequests bool                    `yaml:"require_pushed_authorization_requests,omitempty"`
	Token                              *OAuthTokenConfig       `yaml:"token,omitempty"`
	Scopes                             []string                `yaml:"scopes,omitempty"`
	UserInfo                           *UserInfoConfig         `yaml:"user_info,omitempty"`
	ScopeClaims                        map[string][]string     `yaml:"scope_claims,omitempty"`
	Certificate                        *Certificate            `yaml:"certificate,omitempty"`
	AcrValues                          []string                `yaml:"acr_values,omitempty"`
}

// InboundAuthConfigWithSecret is the wire input wrapper and create/update echo response wrapper.
type InboundAuthConfigWithSecret struct {
	Type        InboundAuthType        `json:"type"             yaml:"type"             jsonschema:"Inbound authentication type. Use 'oauth2' for OAuth/OIDC applications."`
	OAuthConfig *OAuthConfigWithSecret `json:"config,omitempty" yaml:"config,omitempty" jsonschema:"OAuth/OIDC configuration. Required when type is 'oauth2'. Defines OAuth grant types, redirect URIs, client authentication, and PKCE settings."`
}

// InboundAuthConfig is the wire output wrapper (GET responses).
type InboundAuthConfig struct {
	Type        InboundAuthType `json:"type"             yaml:"type"`
	OAuthConfig *OAuthConfig    `json:"config,omitempty" yaml:"config,omitempty"`
}

// InboundAuthConfigProcessed is the runtime wrapper.
type InboundAuthConfigProcessed struct {
	Type        InboundAuthType `json:"type"             yaml:"type,omitempty"`
	OAuthConfig *OAuthClient    `json:"config,omitempty" yaml:"config,omitempty"`
}
