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

import (
	"fmt"
	"net/url"
	"slices"
)

// IsAllowedGrantType checks if the provided grant type is allowed.
func (o *OAuthAppConfigDTO) IsAllowedGrantType(grantType GrantType) bool {
	return containsGrantType(o.GrantTypes, grantType)
}

// IsAllowedResponseType checks if the provided response type is allowed.
func (o *OAuthAppConfigDTO) IsAllowedResponseType(responseType string) bool {
	return containsResponseType(o.ResponseTypes, responseType)
}

// IsAllowedTokenEndpointAuthMethod checks if the provided token endpoint authentication method is allowed.
func (o *OAuthAppConfigDTO) IsAllowedTokenEndpointAuthMethod(method TokenEndpointAuthMethod) bool {
	return o.TokenEndpointAuthMethod == method
}

// ValidateRedirectURI validates the provided redirect URI against the registered redirect URIs.
func (o *OAuthAppConfigDTO) ValidateRedirectURI(redirectURI string) error {
	return validateRedirectURI(o.RedirectURIs, redirectURI)
}

// IsAllowedGrantType checks if the provided grant type is allowed.
func (o *OAuthAppConfigProcessedDTO) IsAllowedGrantType(grantType GrantType) bool {
	return containsGrantType(o.GrantTypes, grantType)
}

// IsAllowedResponseType checks if the provided response type is allowed.
func (o *OAuthAppConfigProcessedDTO) IsAllowedResponseType(responseType string) bool {
	return containsResponseType(o.ResponseTypes, responseType)
}

// IsAllowedTokenEndpointAuthMethod checks if the provided token endpoint authentication method is allowed.
func (o *OAuthAppConfigProcessedDTO) IsAllowedTokenEndpointAuthMethod(method TokenEndpointAuthMethod) bool {
	return o.TokenEndpointAuthMethod == method
}

// ValidateRedirectURI validates the provided redirect URI against the registered redirect URIs.
func (o *OAuthAppConfigProcessedDTO) ValidateRedirectURI(redirectURI string) error {
	return validateRedirectURI(o.RedirectURIs, redirectURI)
}

// RequiresPKCE checks if PKCE is required for this application.
func (o *OAuthAppConfigProcessedDTO) RequiresPKCE() bool {
	return o.PKCERequired || o.PublicClient
}

// --- internal helpers ---

// containsGrantType checks if the provided grant type is in the allowed list.
func containsGrantType(grantTypes []GrantType, grantType GrantType) bool {
	if grantType == "" {
		return false
	}
	return slices.Contains(grantTypes, grantType)
}

// containsResponseType checks if the provided response type is in the allowed list.
func containsResponseType(responseTypes []ResponseType, responseType string) bool {
	if responseType == "" {
		return false
	}
	return slices.Contains(responseTypes, ResponseType(responseType))
}

// validateRedirectURI checks if the provided redirect URI is valid against the registered redirect URIs.
func validateRedirectURI(redirectURIs []string, redirectURI string) error {
	// Allow empty redirectURI only when exactly one redirect URI is registered and it is fully qualified.
	if redirectURI == "" {
		if len(redirectURIs) != 1 {
			return fmt.Errorf("redirect URI is required in the authorization request")
		}
		parsed, err := url.Parse(redirectURIs[0])
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return fmt.Errorf("registered redirect URI is not fully qualified")
		}
		return nil
	}

	if !slices.Contains(redirectURIs, redirectURI) {
		return fmt.Errorf("your application's redirect URL does not match with the registered redirect URLs")
	}

	parsed, err := url.Parse(redirectURI)
	if err != nil {
		return fmt.Errorf("redirect URI is not a valid URL: %w", err)
	}
	if parsed.Fragment != "" {
		return fmt.Errorf("redirect URI must not contain a fragment component")
	}
	return nil
}
