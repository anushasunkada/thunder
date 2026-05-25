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

//nolint:lll // compact assertion helpers in table-style tests.
package thunderidengine_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	sysconfig "github.com/thunder-id/thunderid/internal/system/config"
	"github.com/thunder-id/thunderid/pkg/thunderidengine"
)

var (
	noOAuthPolicy       = thunderidengine.OAuthPolicy{}
	wildcardOAuthPolicy = thunderidengine.OAuthPolicy{AllowWildcardRedirectURI: true}
	parOAuthPolicy      = thunderidengine.OAuthPolicy{RequirePAR: true}
)

const (
	errRedirectURIFragment          = "redirect URI must not contain a fragment component"
	errRedirectURINotRegistered     = "your application's redirect URL does not match with the registered redirect URLs"
	errRedirectURIRequired          = "redirect URI is required in the authorization request"
	errRedirectURINotFullyQualified = "registered redirect URI is not fully qualified"
)

type OAuthClientTestSuite struct {
	suite.Suite
}

func TestOAuthClientTestSuite(t *testing.T) {
	suite.Run(t, new(OAuthClientTestSuite))
}

func (suite *OAuthClientTestSuite) SetupTest() {
	sysconfig.ResetServerRuntime()
	suite.Require().NoError(sysconfig.InitializeServerRuntime("/tmp/test", &sysconfig.Config{}))
}

func (suite *OAuthClientTestSuite) TestIsAllowedGrantType_AuthorizationCode() {
	c := &thunderidengine.OAuthClient{
		GrantTypes: []thunderidengine.GrantType{
			thunderidengine.GrantTypeAuthorizationCode,
			thunderidengine.GrantTypeRefreshToken,
		},
	}

	suite.True(c.IsAllowedGrantType(thunderidengine.GrantTypeAuthorizationCode))
}

func (suite *OAuthClientTestSuite) TestIsAllowedGrantType_ClientCredentials() {
	c := &thunderidengine.OAuthClient{
		GrantTypes: []thunderidengine.GrantType{
			thunderidengine.GrantTypeClientCredentials,
		},
	}

	suite.True(c.IsAllowedGrantType(thunderidengine.GrantTypeClientCredentials))
}

func (suite *OAuthClientTestSuite) TestIsAllowedGrantType_RefreshToken() {
	c := &thunderidengine.OAuthClient{
		GrantTypes: []thunderidengine.GrantType{
			thunderidengine.GrantTypeRefreshToken,
		},
	}

	suite.True(c.IsAllowedGrantType(thunderidengine.GrantTypeRefreshToken))
}

func (suite *OAuthClientTestSuite) TestIsAllowedGrantType_TokenExchange() {
	c := &thunderidengine.OAuthClient{
		GrantTypes: []thunderidengine.GrantType{
			thunderidengine.GrantTypeTokenExchange,
		},
	}

	suite.True(c.IsAllowedGrantType(thunderidengine.GrantTypeTokenExchange))
}

func (suite *OAuthClientTestSuite) TestIsAllowedGrantType_NotAllowed() {
	c := &thunderidengine.OAuthClient{
		GrantTypes: []thunderidengine.GrantType{
			thunderidengine.GrantTypeAuthorizationCode,
		},
	}

	suite.False(c.IsAllowedGrantType(thunderidengine.GrantTypeClientCredentials))
}

func (suite *OAuthClientTestSuite) TestIsAllowedGrantType_EmptyGrantType() {
	c := &thunderidengine.OAuthClient{
		GrantTypes: []thunderidengine.GrantType{
			thunderidengine.GrantTypeAuthorizationCode,
		},
	}

	suite.False(c.IsAllowedGrantType(""))
}

func (suite *OAuthClientTestSuite) TestIsAllowedGrantType_EmptyGrantTypesList() {
	c := &thunderidengine.OAuthClient{
		GrantTypes: []thunderidengine.GrantType{},
	}

	suite.False(c.IsAllowedGrantType(thunderidengine.GrantTypeAuthorizationCode))
}

func (suite *OAuthClientTestSuite) TestIsAllowedGrantType_NilGrantTypesList() {
	c := &thunderidengine.OAuthClient{
		GrantTypes: nil,
	}

	suite.False(c.IsAllowedGrantType(thunderidengine.GrantTypeAuthorizationCode))
}

func (suite *OAuthClientTestSuite) TestIsAllowedGrantType_MultipleGrantTypes() {
	c := &thunderidengine.OAuthClient{
		GrantTypes: []thunderidengine.GrantType{
			thunderidengine.GrantTypeAuthorizationCode,
			thunderidengine.GrantTypeClientCredentials,
			thunderidengine.GrantTypeRefreshToken,
			thunderidengine.GrantTypeTokenExchange,
		},
	}

	suite.True(c.IsAllowedGrantType(thunderidengine.GrantTypeAuthorizationCode))
	suite.True(c.IsAllowedGrantType(thunderidengine.GrantTypeClientCredentials))
	suite.True(c.IsAllowedGrantType(thunderidengine.GrantTypeRefreshToken))
	suite.True(c.IsAllowedGrantType(thunderidengine.GrantTypeTokenExchange))
}

func (suite *OAuthClientTestSuite) TestIsAllowedResponseType_Code() {
	c := &thunderidengine.OAuthClient{
		ResponseTypes: []thunderidengine.ResponseType{
			thunderidengine.ResponseTypeCode,
		},
	}

	suite.True(c.IsAllowedResponseTypeString("code"))
}

func (suite *OAuthClientTestSuite) TestIsAllowedResponseType_NotAllowed() {
	c := &thunderidengine.OAuthClient{
		ResponseTypes: []thunderidengine.ResponseType{
			thunderidengine.ResponseTypeCode,
		},
	}

	suite.False(c.IsAllowedResponseTypeString("token"))
}

func (suite *OAuthClientTestSuite) TestIsAllowedResponseType_EmptyResponseType() {
	c := &thunderidengine.OAuthClient{
		ResponseTypes: []thunderidengine.ResponseType{
			thunderidengine.ResponseTypeCode,
		},
	}

	suite.False(c.IsAllowedResponseTypeString(""))
}

func (suite *OAuthClientTestSuite) TestIsAllowedResponseType_EmptyResponseTypesList() {
	c := &thunderidengine.OAuthClient{
		ResponseTypes: []thunderidengine.ResponseType{},
	}

	suite.False(c.IsAllowedResponseTypeString("code"))
}

func (suite *OAuthClientTestSuite) TestIsAllowedResponseType_NilResponseTypesList() {
	c := &thunderidengine.OAuthClient{
		ResponseTypes: nil,
	}

	suite.False(c.IsAllowedResponseTypeString("code"))
}

func (suite *OAuthClientTestSuite) TestIsAllowedResponseType_MultipleResponseTypes() {
	c := &thunderidengine.OAuthClient{
		ResponseTypes: []thunderidengine.ResponseType{
			thunderidengine.ResponseTypeCode,
			"token",
			"id_token",
		},
	}

	suite.True(c.IsAllowedResponseTypeString("code"))
	suite.True(c.IsAllowedResponseTypeString("token"))
	suite.True(c.IsAllowedResponseTypeString("id_token"))
}

func (suite *OAuthClientTestSuite) TestIsAllowedTokenEndpointAuthMethod_ClientSecretBasic() {
	c := &thunderidengine.OAuthClient{
		TokenEndpointAuthMethod: thunderidengine.TokenEndpointAuthMethodClientSecretBasic,
	}

	suite.True(c.IsAllowedTokenEndpointAuthMethod(thunderidengine.TokenEndpointAuthMethodClientSecretBasic))
}

func (suite *OAuthClientTestSuite) TestIsAllowedTokenEndpointAuthMethod_ClientSecretPost() {
	c := &thunderidengine.OAuthClient{
		TokenEndpointAuthMethod: thunderidengine.TokenEndpointAuthMethodClientSecretPost,
	}

	suite.True(c.IsAllowedTokenEndpointAuthMethod(thunderidengine.TokenEndpointAuthMethodClientSecretPost))
}

func (suite *OAuthClientTestSuite) TestIsAllowedTokenEndpointAuthMethod_None() {
	c := &thunderidengine.OAuthClient{
		TokenEndpointAuthMethod: thunderidengine.TokenEndpointAuthMethodNone,
	}

	suite.True(c.IsAllowedTokenEndpointAuthMethod(thunderidengine.TokenEndpointAuthMethodNone))
}

func (suite *OAuthClientTestSuite) TestIsAllowedTokenEndpointAuthMethod_NotAllowed() {
	c := &thunderidengine.OAuthClient{
		TokenEndpointAuthMethod: thunderidengine.TokenEndpointAuthMethodClientSecretBasic,
	}

	suite.False(c.IsAllowedTokenEndpointAuthMethod(thunderidengine.TokenEndpointAuthMethodClientSecretPost))
}

func (suite *OAuthClientTestSuite) TestIsAllowedTokenEndpointAuthMethod_Empty() {
	c := &thunderidengine.OAuthClient{
		TokenEndpointAuthMethod: thunderidengine.TokenEndpointAuthMethodClientSecretBasic,
	}

	suite.False(c.IsAllowedTokenEndpointAuthMethod(""))
}

func (suite *OAuthClientTestSuite) TestValidateRedirectURI_ValidWithSingleRegisteredURI() {
	c := &thunderidengine.OAuthClient{
		RedirectURIs: []string{"https://example.com/callback"},
	}

	suite.NoError(c.ValidateRedirectURI("https://example.com/callback", noOAuthPolicy))
}

func (suite *OAuthClientTestSuite) TestValidateRedirectURI_ValidHTTPLocalhostWithPort() {
	c := &thunderidengine.OAuthClient{
		RedirectURIs: []string{"http://localhost:3000/callback"},
	}

	suite.NoError(c.ValidateRedirectURI("http://localhost:3000/callback", noOAuthPolicy))
}

func (suite *OAuthClientTestSuite) TestValidateRedirectURI_ValidHTTPSWithPath() {
	c := &thunderidengine.OAuthClient{
		RedirectURIs: []string{"https://app.example.com/oauth/callback"},
	}

	suite.NoError(c.ValidateRedirectURI("https://app.example.com/oauth/callback", noOAuthPolicy))
}

func (suite *OAuthClientTestSuite) TestValidateRedirectURI_ValidCustomScheme() {
	c := &thunderidengine.OAuthClient{
		RedirectURIs: []string{"myapp://callback"},
	}

	suite.NoError(c.ValidateRedirectURI("myapp://callback", noOAuthPolicy))
}

func (suite *OAuthClientTestSuite) TestValidateRedirectURI_ValidWithQueryParameters() {
	c := &thunderidengine.OAuthClient{
		RedirectURIs: []string{"https://example.com/callback?param=value"},
	}

	suite.NoError(c.ValidateRedirectURI("https://example.com/callback?param=value", noOAuthPolicy))
}

func (suite *OAuthClientTestSuite) TestValidateRedirectURI_InvalidWithFragment() {
	c := &thunderidengine.OAuthClient{
		RedirectURIs: []string{"https://example.com/callback#fragment"},
	}

	err := c.ValidateRedirectURI("https://example.com/callback#fragment", noOAuthPolicy)
	suite.EqualError(err, errRedirectURIFragment)
}

func (suite *OAuthClientTestSuite) TestValidateRedirectURI_NotRegistered() {
	c := &thunderidengine.OAuthClient{
		RedirectURIs: []string{"https://example.com/callback"},
	}

	err := c.ValidateRedirectURI("https://different.com/callback", noOAuthPolicy)
	suite.EqualError(err, errRedirectURINotRegistered)
}

func (suite *OAuthClientTestSuite) TestValidateRedirectURI_EmptyWithSingleFullyQualifiedURI() {
	c := &thunderidengine.OAuthClient{
		RedirectURIs: []string{"https://example.com/callback"},
	}

	suite.NoError(c.ValidateRedirectURI("", noOAuthPolicy))
}

func (suite *OAuthClientTestSuite) TestValidateRedirectURI_EmptyWithMultipleURIs() {
	c := &thunderidengine.OAuthClient{
		RedirectURIs: []string{
			"https://example.com/callback",
			"https://example.com/callback2",
		},
	}

	err := c.ValidateRedirectURI("", noOAuthPolicy)
	suite.EqualError(err, errRedirectURIRequired)
}

func (suite *OAuthClientTestSuite) TestValidateRedirectURI_EmptyWithPartialRegisteredURI() {
	c := &thunderidengine.OAuthClient{
		RedirectURIs: []string{"/callback"},
	}

	err := c.ValidateRedirectURI("", noOAuthPolicy)
	suite.EqualError(err, errRedirectURINotFullyQualified)
}

func (suite *OAuthClientTestSuite) TestValidateRedirectURI_EmptyWithInvalidRegisteredURI() {
	c := &thunderidengine.OAuthClient{
		RedirectURIs: []string{"://invalid"},
	}

	err := c.ValidateRedirectURI("", noOAuthPolicy)
	suite.EqualError(err, errRedirectURINotFullyQualified)
}

func (suite *OAuthClientTestSuite) TestValidateRedirectURI_EmptyRedirectURIsList() {
	c := &thunderidengine.OAuthClient{
		RedirectURIs: []string{},
	}

	err := c.ValidateRedirectURI("", noOAuthPolicy)
	suite.EqualError(err, errRedirectURIRequired)
}

func (suite *OAuthClientTestSuite) TestValidateRedirectURI_NilRedirectURIsList() {
	c := &thunderidengine.OAuthClient{
		RedirectURIs: nil,
	}

	err := c.ValidateRedirectURI("", noOAuthPolicy)
	suite.EqualError(err, errRedirectURIRequired)
}

func (suite *OAuthClientTestSuite) TestRequiresPKCE_PKCERequiredTrue() {
	c := &thunderidengine.OAuthClient{PKCERequired: true, PublicClient: false}
	suite.True(c.RequiresPKCE())
}

func (suite *OAuthClientTestSuite) TestRequiresPKCE_PublicClientTrue() {
	c := &thunderidengine.OAuthClient{PKCERequired: false, PublicClient: true}
	suite.True(c.RequiresPKCE())
}

func (suite *OAuthClientTestSuite) TestRequiresPKCE_BothTrue() {
	c := &thunderidengine.OAuthClient{PKCERequired: true, PublicClient: true}
	suite.True(c.RequiresPKCE())
}

func (suite *OAuthClientTestSuite) TestRequiresPKCE_BothFalse() {
	c := &thunderidengine.OAuthClient{PKCERequired: false, PublicClient: false}
	suite.False(c.RequiresPKCE())
}

type OAuthHelperTestSuite struct {
	suite.Suite
}

func TestOAuthHelperTestSuite(t *testing.T) {
	suite.Run(t, new(OAuthHelperTestSuite))
}

func (suite *OAuthHelperTestSuite) SetupTest() {
	sysconfig.ResetServerRuntime()
	suite.Require().NoError(sysconfig.InitializeServerRuntime("/tmp/test", &sysconfig.Config{}))
}

func (suite *OAuthHelperTestSuite) TestIsAllowedGrantType_ValidGrantType() {
	grantTypes := []thunderidengine.GrantType{
		thunderidengine.GrantTypeAuthorizationCode,
		thunderidengine.GrantTypeRefreshToken,
	}

	suite.True((&thunderidengine.OAuthClient{GrantTypes: grantTypes}).IsAllowedGrantType(thunderidengine.GrantTypeAuthorizationCode))
}

func (suite *OAuthHelperTestSuite) TestIsAllowedGrantType_InvalidGrantType() {
	grantTypes := []thunderidengine.GrantType{
		thunderidengine.GrantTypeAuthorizationCode,
	}

	suite.False((&thunderidengine.OAuthClient{GrantTypes: grantTypes}).IsAllowedGrantType(thunderidengine.GrantTypeClientCredentials))
}

func (suite *OAuthHelperTestSuite) TestIsAllowedGrantType_EmptyGrantType() {
	grantTypes := []thunderidengine.GrantType{
		thunderidengine.GrantTypeAuthorizationCode,
	}

	suite.False((&thunderidengine.OAuthClient{GrantTypes: grantTypes}).IsAllowedGrantType(""))
}

func (suite *OAuthHelperTestSuite) TestIsAllowedGrantType_EmptyList() {
	suite.False((&thunderidengine.OAuthClient{GrantTypes: []thunderidengine.GrantType{}}).IsAllowedGrantType(thunderidengine.GrantTypeAuthorizationCode))
}

func (suite *OAuthHelperTestSuite) TestIsAllowedGrantType_NilList() {
	suite.False((&thunderidengine.OAuthClient{GrantTypes: nil}).IsAllowedGrantType(thunderidengine.GrantTypeAuthorizationCode))
}

func (suite *OAuthHelperTestSuite) TestIsAllowedResponseType_ValidResponseType() {
	responseTypes := []thunderidengine.ResponseType{
		thunderidengine.ResponseTypeCode,
		"token",
	}

	suite.True((&thunderidengine.OAuthClient{ResponseTypes: responseTypes}).IsAllowedResponseTypeString("code"))
}

func (suite *OAuthHelperTestSuite) TestIsAllowedResponseType_InvalidResponseType() {
	responseTypes := []thunderidengine.ResponseType{
		thunderidengine.ResponseTypeCode,
	}

	suite.False((&thunderidengine.OAuthClient{ResponseTypes: responseTypes}).IsAllowedResponseTypeString("token"))
}

func (suite *OAuthHelperTestSuite) TestIsAllowedResponseType_EmptyResponseType() {
	responseTypes := []thunderidengine.ResponseType{
		thunderidengine.ResponseTypeCode,
	}

	suite.False((&thunderidengine.OAuthClient{ResponseTypes: responseTypes}).IsAllowedResponseTypeString(""))
}

func (suite *OAuthHelperTestSuite) TestIsAllowedResponseType_EmptyList() {
	suite.False((&thunderidengine.OAuthClient{ResponseTypes: []thunderidengine.ResponseType{}}).IsAllowedResponseTypeString("code"))
}

func (suite *OAuthHelperTestSuite) TestIsAllowedResponseType_NilList() {
	suite.False((&thunderidengine.OAuthClient{ResponseTypes: nil}).IsAllowedResponseTypeString("code"))
}

func (suite *OAuthHelperTestSuite) TestValidateRedirectURI_ValidSingleURI() {
	err := thunderidengine.ValidateRedirectURI([]string{"https://example.com/callback"}, "https://example.com/callback", noOAuthPolicy)
	suite.NoError(err)
}

func (suite *OAuthHelperTestSuite) TestValidateRedirectURI_ValidMultipleURIs() {
	redirectURIs := []string{
		"https://example.com/callback",
		"https://example.com/callback2",
	}

	err := thunderidengine.ValidateRedirectURI(redirectURIs, "https://example.com/callback2", noOAuthPolicy)
	suite.NoError(err)
}

func (suite *OAuthHelperTestSuite) TestValidateRedirectURI_InvalidNotRegistered() {
	err := thunderidengine.ValidateRedirectURI([]string{"https://example.com/callback"}, "https://different.com/callback", noOAuthPolicy)
	suite.EqualError(err, errRedirectURINotRegistered)
}

func (suite *OAuthHelperTestSuite) TestValidateRedirectURI_InvalidWithFragment() {
	uri := "https://example.com/callback#fragment"
	err := thunderidengine.ValidateRedirectURI([]string{uri}, uri, noOAuthPolicy)
	suite.EqualError(err, errRedirectURIFragment)
}

func (suite *OAuthHelperTestSuite) TestValidateRedirectURI_EmptyURIWithSingleFullyQualified() {
	err := thunderidengine.ValidateRedirectURI([]string{"https://example.com/callback"}, "", noOAuthPolicy)
	suite.NoError(err)
}

func (suite *OAuthHelperTestSuite) TestValidateRedirectURI_EmptyURIWithMultiple() {
	redirectURIs := []string{
		"https://example.com/callback",
		"https://example.com/callback2",
	}

	err := thunderidengine.ValidateRedirectURI(redirectURIs, "", noOAuthPolicy)
	suite.EqualError(err, errRedirectURIRequired)
}

func (suite *OAuthHelperTestSuite) TestValidateRedirectURI_EmptyURIWithPartialRegistered() {
	err := thunderidengine.ValidateRedirectURI([]string{"/callback"}, "", noOAuthPolicy)
	suite.EqualError(err, errRedirectURINotFullyQualified)
}

func (suite *OAuthHelperTestSuite) TestValidateRedirectURI_EmptyURIWithNoScheme() {
	err := thunderidengine.ValidateRedirectURI([]string{"example.com/callback"}, "", noOAuthPolicy)
	suite.EqualError(err, errRedirectURINotFullyQualified)
}

func (suite *OAuthHelperTestSuite) TestValidateRedirectURI_EmptyURIWithNoHost() {
	err := thunderidengine.ValidateRedirectURI([]string{"https:///callback"}, "", noOAuthPolicy)
	suite.EqualError(err, errRedirectURINotFullyQualified)
}

func (suite *OAuthHelperTestSuite) TestValidateRedirectURI_EmptyURIList() {
	err := thunderidengine.ValidateRedirectURI([]string{}, "", noOAuthPolicy)
	suite.EqualError(err, errRedirectURIRequired)
}

func (suite *OAuthHelperTestSuite) TestValidateRedirectURI_NilList() {
	err := thunderidengine.ValidateRedirectURI(nil, "", noOAuthPolicy)
	suite.EqualError(err, errRedirectURIRequired)
}

func (suite *OAuthHelperTestSuite) TestValidateRedirectURI_CustomScheme() {
	err := thunderidengine.ValidateRedirectURI([]string{"myapp://callback"}, "myapp://callback", noOAuthPolicy)
	suite.NoError(err)
}

func (suite *OAuthHelperTestSuite) TestValidateRedirectURI_LocalhostHTTP() {
	err := thunderidengine.ValidateRedirectURI([]string{"http://localhost:3000/callback"}, "http://localhost:3000/callback", noOAuthPolicy)
	suite.NoError(err)
}

func (suite *OAuthHelperTestSuite) TestValidateRedirectURI_WithQueryParams() {
	uri := "https://example.com/callback?foo=bar"
	suite.NoError(thunderidengine.ValidateRedirectURI([]string{uri}, uri, noOAuthPolicy))
}

func (suite *OAuthHelperTestSuite) TestValidateRedirectURI_IPAddress() {
	err := thunderidengine.ValidateRedirectURI([]string{"https://192.168.1.1/callback"}, "https://192.168.1.1/callback", noOAuthPolicy)
	suite.NoError(err)
}

func (suite *OAuthHelperTestSuite) TestValidateRedirectURI_Localhost127() {
	err := thunderidengine.ValidateRedirectURI([]string{"http://127.0.0.1:8080/callback"}, "http://127.0.0.1:8080/callback", noOAuthPolicy)
	suite.NoError(err)
}

func (suite *OAuthHelperTestSuite) TestValidateRedirectURI_InvalidURLFormat() {
	uri := "http://example.com/callback\x00invalid"
	err := thunderidengine.ValidateRedirectURI([]string{uri}, uri, noOAuthPolicy)
	suite.Error(err)
	assert.Contains(suite.T(), err.Error(), "invalid redirect URI")
}

func (suite *OAuthClientTestSuite) TestRequiresPAR_GlobalConfigEnabled() {
	c := &thunderidengine.OAuthClient{RequirePushedAuthorizationRequests: false}
	suite.True(c.RequiresPAR(parOAuthPolicy))
}

func (suite *OAuthClientTestSuite) TestRequiresPAR_PerClientEnabled() {
	c := &thunderidengine.OAuthClient{RequirePushedAuthorizationRequests: true}
	suite.True(c.RequiresPAR(noOAuthPolicy))
}

func (suite *OAuthClientTestSuite) TestRequiresPAR_BothFalse() {
	c := &thunderidengine.OAuthClient{RequirePushedAuthorizationRequests: false}
	suite.False(c.RequiresPAR(noOAuthPolicy))
}

func (suite *OAuthHelperTestSuite) TestMatchAnyRedirectURIPattern_WildcardEnabled_Matches() {
	err := thunderidengine.ValidateRedirectURI(
		[]string{"https://app.example.com/*"},
		"https://app.example.com/cb",
		wildcardOAuthPolicy,
	)
	suite.NoError(err)
}

func (suite *OAuthHelperTestSuite) TestMatchAnyRedirectURIPattern_WildcardDisabled_NoMatch() {
	err := thunderidengine.ValidateRedirectURI(
		[]string{"https://app.example.com/*"},
		"https://app.example.com/cb",
		noOAuthPolicy,
	)
	suite.Error(err)
}

func (suite *OAuthHelperTestSuite) TestMatchAnyRedirectURIPattern_HostWildcardEnabled_Matches() {
	err := thunderidengine.ValidateRedirectURI(
		[]string{"https://tenant-app-*-*.gateway.example.com/cb"},
		"https://tenant-app-019dfc78-f19ab4f2.gateway.example.com/cb",
		wildcardOAuthPolicy,
	)
	suite.NoError(err)
}

func (suite *OAuthHelperTestSuite) TestMatchAnyRedirectURIPattern_HostWildcardEnabled_NonMatchingDynamicPart() {
	// Hyphen inside the dynamic part is not in [0-9a-zA-Z]+, so this must fail.
	err := thunderidengine.ValidateRedirectURI(
		[]string{"https://app-*-prod.example.com/cb"},
		"https://app-foo-bar-prod.example.com/cb",
		wildcardOAuthPolicy,
	)
	suite.Error(err)
}

func (suite *OAuthHelperTestSuite) TestMatchAnyRedirectURIPattern_HostWildcardDisabled_NoMatch() {
	// Default: AllowWildcardRedirectURI = false. Note the pattern would never have made it
	// past registration with the flag off, but we still verify the matcher returns no match.
	err := thunderidengine.ValidateRedirectURI(
		[]string{"https://app-*.example.com/cb"},
		"https://app-prod.example.com/cb",
		noOAuthPolicy,
	)
	suite.Error(err)
}

func (suite *OAuthHelperTestSuite) TestMatchAnyRedirectURIPattern_HostWildcardDoesNotCrossDot() {
	err := thunderidengine.ValidateRedirectURI(
		[]string{"https://app-*.example.com/cb"},
		"https://app-foo.evil.example.com/cb",
		wildcardOAuthPolicy,
	)
	suite.Error(err)
}
