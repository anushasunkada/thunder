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

//nolint:revive // OAuth enum constants mirror wire values; per-const doc comments are omitted intentionally.
package thunderidengine

// GrantType is an OAuth 2.0 grant type.
type GrantType string

const (
	GrantTypeAuthorizationCode GrantType = "authorization_code"
	GrantTypeClientCredentials GrantType = "client_credentials"
	GrantTypeRefreshToken      GrantType = "refresh_token"
	GrantTypeTokenExchange     GrantType = "urn:ietf:params:oauth:grant-type:token-exchange" //nolint:gosec
)

var supportedGrantTypes = []GrantType{
	GrantTypeAuthorizationCode,
	GrantTypeClientCredentials,
	GrantTypeRefreshToken,
	GrantTypeTokenExchange,
}

// IsValid reports whether gt is a supported grant type.
func (gt GrantType) IsValid() bool {
	for _, valid := range supportedGrantTypes {
		if gt == valid {
			return true
		}
	}
	return false
}

// ResponseType is an OAuth 2.0 authorization response type.
type ResponseType string

const (
	ResponseTypeCode    ResponseType = "code"
	ResponseTypeIDToken ResponseType = "id_token"
)

var supportedResponseTypes = []ResponseType{
	ResponseTypeCode,
}

// IsValid reports whether rt is a supported response type.
func (rt ResponseType) IsValid() bool {
	for _, valid := range supportedResponseTypes {
		if rt == valid {
			return true
		}
	}
	return false
}

// TokenEndpointAuthMethod is a token endpoint client authentication method.
type TokenEndpointAuthMethod string

const (
	TokenEndpointAuthMethodClientSecretBasic TokenEndpointAuthMethod = "client_secret_basic"
	TokenEndpointAuthMethodClientSecretPost  TokenEndpointAuthMethod = "client_secret_post"
	TokenEndpointAuthMethodPrivateKeyJWT     TokenEndpointAuthMethod = "private_key_jwt" //nolint:gosec
	TokenEndpointAuthMethodNone              TokenEndpointAuthMethod = "none"
)

var supportedTokenEndpointAuthMethods = []TokenEndpointAuthMethod{
	TokenEndpointAuthMethodClientSecretBasic,
	TokenEndpointAuthMethodClientSecretPost,
	TokenEndpointAuthMethodPrivateKeyJWT,
	TokenEndpointAuthMethodNone,
}

// IsValid reports whether tam is a supported token endpoint authentication method.
func (tam TokenEndpointAuthMethod) IsValid() bool {
	for _, valid := range supportedTokenEndpointAuthMethods {
		if tam == valid {
			return true
		}
	}
	return false
}

// IDTokenResponseType is the response format of the ID token.
type IDTokenResponseType string

const (
	IDTokenResponseTypeJWT       IDTokenResponseType = "JWT"
	IDTokenResponseTypeJWE       IDTokenResponseType = "JWE"
	IDTokenResponseTypeNestedJWT IDTokenResponseType = "NESTED_JWT" //nolint:gosec
)

// UserInfoResponseType is the response format of the UserInfo endpoint.
type UserInfoResponseType string

const (
	UserInfoResponseTypeJSON      UserInfoResponseType = "JSON"
	UserInfoResponseTypeJWS       UserInfoResponseType = "JWS"
	UserInfoResponseTypeJWE       UserInfoResponseType = "JWE"
	UserInfoResponseTypeNestedJWT UserInfoResponseType = "NESTED_JWT"
)

// SupportedResponseTypeStrings returns supported OAuth2 response types as strings.
func SupportedResponseTypeStrings() []string {
	result := make([]string, len(supportedResponseTypes))
	for i, rt := range supportedResponseTypes {
		result[i] = string(rt)
	}
	return result
}

// SupportedGrantTypeStrings returns supported OAuth2 grant types as strings.
func SupportedGrantTypeStrings() []string {
	result := make([]string, len(supportedGrantTypes))
	for i, gt := range supportedGrantTypes {
		result[i] = string(gt)
	}
	return result
}

// SupportedTokenEndpointAuthMethodStrings returns supported token endpoint auth methods as strings.
func SupportedTokenEndpointAuthMethodStrings() []string {
	result := make([]string, len(supportedTokenEndpointAuthMethods))
	for i, tam := range supportedTokenEndpointAuthMethods {
		result[i] = string(tam)
	}
	return result
}
