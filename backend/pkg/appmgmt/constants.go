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

// InboundAuthType represents the type of inbound authentication.
type InboundAuthType string

const (
	// OAuthInboundAuthType represents the OAuth 2.0 inbound authentication type.
	OAuthInboundAuthType InboundAuthType = "oauth2"
)

// UserInfoResponseType represents the response format of the UserInfo endpoint.
type UserInfoResponseType string

const (
	// UserInfoResponseTypeJSON represents the JSON userinfo response type.
	UserInfoResponseTypeJSON UserInfoResponseType = "JSON"

	// UserInfoResponseTypeJWS represents the JWS userinfo response type.
	UserInfoResponseTypeJWS UserInfoResponseType = "JWS"
)

// ApplicationNotFoundError is the error returned when an application is not found.
var ApplicationNotFoundError error = errors.New("application not found")

// ApplicationDataCorruptedError is the error returned when application data is corrupted.
var ApplicationDataCorruptedError error = errors.New("application data is corrupted")

// Constants for MCP tool defaults
var (
	// DefaultUserAttributes are the standard user attributes for application templates.
	DefaultUserAttributes = []string{
		"email", "name", "given_name", "family_name",
		"picture", "phone_number", "address", "created_at",
	}
	// DefaultScopes are the standard OAuth scopes for application templates.
	DefaultScopes = []string{"openid", "profile", "email"}
)

// CertificateType represents the type of certificates in the system.
type CertificateType string

const (
	// CertificateTypeJWKS represents a JSON Web Key Set (JWKS) certificate.
	CertificateTypeJWKS CertificateType = "JWKS"
	// CertificateTypeJWKSURI represents a JWKS URI certificate.
	CertificateTypeJWKSURI CertificateType = "JWKS_URI"
)

// GrantType defines a type for OAuth2 grant types.
type GrantType string

const (
	// GrantTypeAuthorizationCode represents the authorization code grant type.
	GrantTypeAuthorizationCode GrantType = "authorization_code"
	// GrantTypeClientCredentials represents the client credentials grant type.
	GrantTypeClientCredentials GrantType = "client_credentials"
	// GrantTypeRefreshToken represents the refresh token grant type.
	GrantTypeRefreshToken GrantType = "refresh_token"
	// GrantTypeTokenExchange represents the token exchange grant type.
	GrantTypeTokenExchange GrantType = "urn:ietf:params:oauth:grant-type:token-exchange" //nolint:gosec
)

// supportedGrantTypes is the single source of truth for all supported grant types.
var supportedGrantTypes = []GrantType{
	GrantTypeAuthorizationCode,
	GrantTypeClientCredentials,
	GrantTypeRefreshToken,
	GrantTypeTokenExchange,
}

// IsValid checks if the GrantType is valid.
func (gt GrantType) IsValid() bool {
	for _, valid := range supportedGrantTypes {
		if gt == valid {
			return true
		}
	}
	return false
}

// ResponseType defines a type for OAuth2 response types.
type ResponseType string

const (
	// ResponseTypeCode represents the authorization code response type.
	ResponseTypeCode ResponseType = "code"
	// ResponseTypeIDToken represents the id token response type.
	ResponseTypeIDToken ResponseType = "id_token"
)

// supportedResponseTypes is the single source of truth for all supported response types.
var supportedResponseTypes = []ResponseType{
	ResponseTypeCode,
}

// IsValid checks if the ResponseType is valid.
func (rt ResponseType) IsValid() bool {
	for _, valid := range supportedResponseTypes {
		if rt == valid {
			return true
		}
	}
	return false
}

// TokenEndpointAuthMethod defines a type for token endpoint authentication methods.
type TokenEndpointAuthMethod string

const (
	// TokenEndpointAuthMethodClientSecretBasic represents the client secret basic authentication method.
	TokenEndpointAuthMethodClientSecretBasic TokenEndpointAuthMethod = "client_secret_basic"
	// TokenEndpointAuthMethodClientSecretPost represents the client secret post authentication method.
	TokenEndpointAuthMethodClientSecretPost TokenEndpointAuthMethod = "client_secret_post"
	// TokenEndpointAuthMethodPrivateKeyJWT represents the private key JWT authentication method.
	// #nosec G101 - This is not a hardcoded credential, but a constant representing an authentication method.
	TokenEndpointAuthMethodPrivateKeyJWT TokenEndpointAuthMethod = "private_key_jwt"
	// TokenEndpointAuthMethodNone represents no authentication method.
	TokenEndpointAuthMethodNone TokenEndpointAuthMethod = "none"
)

// supportedTokenEndpointAuthMethods is the single source of truth for all supported token endpoint
// authentication methods.
var supportedTokenEndpointAuthMethods = []TokenEndpointAuthMethod{
	TokenEndpointAuthMethodClientSecretBasic,
	TokenEndpointAuthMethodClientSecretPost,
	TokenEndpointAuthMethodPrivateKeyJWT,
	TokenEndpointAuthMethodNone,
}

// IsValid checks if the TokenEndpointAuthMethod is valid.
func (tam TokenEndpointAuthMethod) IsValid() bool {
	for _, valid := range supportedTokenEndpointAuthMethods {
		if tam == valid {
			return true
		}
	}
	return false
}

// OAuth2 token types.
const (
	TokenTypeBearer = "Bearer"
)

// TokenTypeIdentifier defines a type for RFC 8693 token type identifiers.
type TokenTypeIdentifier string

// RFC 8693 Token Type Identifiers
const (
	//nolint:gosec // Token type identifier, not a credential
	TokenTypeIdentifierAccessToken TokenTypeIdentifier = "urn:ietf:params:oauth:token-type:access_token"
	//nolint:gosec // Token type identifier, not a credential
	TokenTypeIdentifierRefreshToken TokenTypeIdentifier = "urn:ietf:params:oauth:token-type:refresh_token"
	//nolint:gosec // Token type identifier, not a credential
	TokenTypeIdentifierIDToken TokenTypeIdentifier = "urn:ietf:params:oauth:token-type:id_token"
	//nolint:gosec // Token type identifier, not a credential
	TokenTypeIdentifierJWT TokenTypeIdentifier = "urn:ietf:params:oauth:token-type:jwt"
)

// supportedTokenTypeIdentifiers is the single source of truth for all supported token type identifiers.
var supportedTokenTypeIdentifiers = []TokenTypeIdentifier{
	TokenTypeIdentifierAccessToken,
	TokenTypeIdentifierRefreshToken,
	TokenTypeIdentifierIDToken,
	TokenTypeIdentifierJWT,
}

// IsValid checks if the TokenTypeIdentifier is valid.
func (tti TokenTypeIdentifier) IsValid() bool {
	for _, valid := range supportedTokenTypeIdentifiers {
		if tti == valid {
			return true
		}
	}
	return false
}

// GetSupportedGrantTypes returns all supported grant types as strings.
func GetSupportedGrantTypes() []string {
	result := make([]string, len(supportedGrantTypes))
	for i, gt := range supportedGrantTypes {
		result[i] = string(gt)
	}
	return result
}

// GetSupportedResponseTypes returns all supported response types as strings.
func GetSupportedResponseTypes() []string {
	result := make([]string, len(supportedResponseTypes))
	for i, rt := range supportedResponseTypes {
		result[i] = string(rt)
	}
	return result
}

// GetSupportedTokenEndpointAuthMethods returns all supported token endpoint auth methods as strings.
func GetSupportedTokenEndpointAuthMethods() []string {
	result := make([]string, len(supportedTokenEndpointAuthMethods))
	for i, m := range supportedTokenEndpointAuthMethods {
		result[i] = string(m)
	}
	return result
}

// GetSupportedTokenTypeIdentifiers returns all supported token type identifiers as strings.
func GetSupportedTokenTypeIdentifiers() []string {
	result := make([]string, len(supportedTokenTypeIdentifiers))
	for i, t := range supportedTokenTypeIdentifiers {
		result[i] = string(t)
	}
	return result
}
