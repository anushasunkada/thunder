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

// Package thunderidengine provides OAuth runtime models for the embeddable engine.
//
//nolint:lll // struct tags mirror inboundclient jsonschema descriptions.
package thunderidengine

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/thunder-id/thunderid/internal/system/log"
	"github.com/thunder-id/thunderid/internal/system/utils"
)

// OAuthTokenConfig wraps access and ID token configs.
type OAuthTokenConfig struct {
	AccessToken *AccessTokenConfig `json:"accessToken,omitempty" yaml:"access_token,omitempty" jsonschema:"Access token configuration."`
	IDToken     *IDTokenConfig     `json:"idToken,omitempty"    yaml:"id_token,omitempty"     jsonschema:"ID token configuration."`
}

// AccessTokenConfig is the access token configuration.
type AccessTokenConfig struct {
	ValidityPeriod int64    `json:"validityPeriod,omitempty" yaml:"validity_period,omitempty" jsonschema:"Access token validity period in seconds."`
	UserAttributes []string `json:"userAttributes,omitempty" yaml:"user_attributes,omitempty" jsonschema:"User attributes to embed in the access token."`
}

// IDTokenConfig is the ID token configuration.
type IDTokenConfig struct {
	ValidityPeriod int64               `json:"validityPeriod,omitempty" yaml:"validity_period,omitempty" jsonschema:"ID token validity period in seconds."`
	UserAttributes []string            `json:"userAttributes,omitempty" yaml:"user_attributes,omitempty" jsonschema:"User attributes to embed in the ID token."`
	ResponseType   IDTokenResponseType `json:"responseType,omitempty"   yaml:"response_type,omitempty"   jsonschema:"ID token response type (JWT, JWE, NESTED_JWT). Defaults to JWT."`
	EncryptionAlg  string              `json:"encryptionAlg,omitempty"  yaml:"encryption_alg,omitempty"  jsonschema:"JWE key-management algorithm. Required when responseType is JWE or NESTED_JWT."`
	EncryptionEnc  string              `json:"encryptionEnc,omitempty"  yaml:"encryption_enc,omitempty"  jsonschema:"JWE content-encryption algorithm. Required when responseType is JWE or NESTED_JWT."`
}

// UserInfoConfig is the userinfo endpoint configuration.
type UserInfoConfig struct {
	ResponseType   UserInfoResponseType `json:"responseType,omitempty"   yaml:"response_type,omitempty"   jsonschema:"UserInfo response type (JSON, JWS, JWE, NESTED_JWT). Required algorithm fields must match the selected response type."`
	UserAttributes []string             `json:"userAttributes,omitempty" yaml:"user_attributes,omitempty" jsonschema:"User attributes to include in the userinfo response."`
	SigningAlg     string               `json:"signingAlg,omitempty"     yaml:"signing_alg,omitempty"     jsonschema:"JWS algorithm for signed userinfo responses (e.g. RS256)."`
	EncryptionAlg  string               `json:"encryptionAlg,omitempty"  yaml:"encryption_alg,omitempty"  jsonschema:"JWE key-management algorithm for encrypted userinfo responses (e.g. RSA-OAEP-256)."`
	EncryptionEnc  string               `json:"encryptionEnc,omitempty"  yaml:"encryption_enc,omitempty"  jsonschema:"JWE content-encryption algorithm (e.g. A256GCM). Required when encryptionAlg is set."`
}

// FlowApplication is the runtime view used when assembling flow execution context.
type FlowApplication struct {
	ID                        string
	Name                      string
	Assertion                 *AssertionConfig
	LoginConsent              *LoginConsentConfig
	AllowedUserTypes          []string
	Metadata                  map[string]interface{}
	OAuthClientID             string
	AuthFlowID                string
	RegistrationFlowID        string
	RecoveryFlowID            string
	IsRegistrationFlowEnabled bool
	IsRecoveryFlowEnabled     bool
}

// IsAllowedGrantType reports whether the given grant type is allowed for this client.
func (o *OAuthClient) IsAllowedGrantType(grantType GrantType) bool {
	if grantType == "" {
		return false
	}
	return slices.Contains(o.GrantTypes, grantType)
}

// IsAllowedResponseType reports whether the given response type is allowed for this client.
func (o *OAuthClient) IsAllowedResponseType(responseType ResponseType) bool {
	if responseType == "" {
		return false
	}
	return slices.Contains(o.ResponseTypes, responseType)
}

// IsAllowedResponseTypeString reports whether the wire response type is allowed for this client.
func (o *OAuthClient) IsAllowedResponseTypeString(responseType string) bool {
	return o.IsAllowedResponseType(ResponseType(responseType))
}

// IsAllowedTokenEndpointAuthMethod reports whether the given auth method matches this client.
func (o *OAuthClient) IsAllowedTokenEndpointAuthMethod(method TokenEndpointAuthMethod) bool {
	return o.TokenEndpointAuthMethod == method
}

// ValidateRedirectURI validates the given redirect URI against this client's registered URIs.
func (o *OAuthClient) ValidateRedirectURI(redirectURI string, policy OAuthPolicy) error {
	return ValidateRedirectURI(o.RedirectURIs, redirectURI, policy)
}

// RequiresPKCE reports whether PKCE is required for this client.
func (o *OAuthClient) RequiresPKCE() bool {
	return o.PKCERequired || o.PublicClient
}

// RequiresPAR reports whether pushed authorization requests are required for this client.
func (o *OAuthClient) RequiresPAR(policy OAuthPolicy) bool {
	return o.RequirePushedAuthorizationRequests || policy.RequirePAR
}

// ValidateRedirectURI validates the provided redirect URI against the registered list.
func ValidateRedirectURI(redirectURIs []string, redirectURI string, policy OAuthPolicy) error {
	logger := log.GetLogger()

	if redirectURI == "" {
		if len(redirectURIs) != 1 {
			return fmt.Errorf("redirect URI is required in the authorization request")
		}
		if strings.Contains(redirectURIs[0], "*") {
			return fmt.Errorf("redirect URI is required in the authorization request")
		}
		parsed, err := url.Parse(redirectURIs[0])
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return fmt.Errorf("registered redirect URI is not fully qualified")
		}
		return nil
	}

	if !matchAnyRedirectURIPattern(redirectURIs, redirectURI, policy) {
		return fmt.Errorf("your application's redirect URL does not match with the registered redirect URLs")
	}

	parsedRedirectURI, err := utils.ParseURL(redirectURI)
	if err != nil {
		logger.Error("Failed to parse redirect URI", log.Error(err))
		return fmt.Errorf("invalid redirect URI: %s", err.Error())
	}
	if parsedRedirectURI.Fragment != "" {
		return fmt.Errorf("redirect URI must not contain a fragment component")
	}

	return nil
}

func matchAnyRedirectURIPattern(patterns []string, redirectURI string, policy OAuthPolicy) bool {
	wildcardEnabled := policy.AllowWildcardRedirectURI
	for _, pattern := range patterns {
		if !wildcardEnabled || !strings.Contains(pattern, "*") {
			if pattern == redirectURI {
				return true
			}
			continue
		}
		matched, err := utils.MatchURIPattern(pattern, redirectURI)
		if err != nil {
			continue
		}
		if matched {
			return true
		}
	}
	return false
}
