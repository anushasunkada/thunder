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

package adapter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/thunder-id/thunderid/internal/entityprovider"
	"github.com/thunder-id/thunderid/internal/inboundclient"
	inboundmodel "github.com/thunder-id/thunderid/internal/inboundclient/model"
	oauth2const "github.com/thunder-id/thunderid/internal/oauth/oauth2/constants"
	"github.com/thunder-id/thunderid/pkg/thunderidengine"
	"github.com/thunder-id/thunderid/tests/mocks/entityprovidermock"
	"github.com/thunder-id/thunderid/tests/mocks/inboundclientmock"
)

type ClientProviderTestSuite struct {
	suite.Suite
	mockEntityProvider *entityprovidermock.EntityProviderInterfaceMock
	mockInboundClient  *inboundclientmock.InboundClientServiceInterfaceMock
	provider           thunderidengine.ClientProvider
	ctx                context.Context
}

func TestClientProviderTestSuite(t *testing.T) {
	suite.Run(t, new(ClientProviderTestSuite))
}

func (suite *ClientProviderTestSuite) SetupTest() {
	suite.mockEntityProvider = entityprovidermock.NewEntityProviderInterfaceMock(suite.T())
	suite.mockInboundClient = inboundclientmock.NewInboundClientServiceInterfaceMock(suite.T())
	suite.provider = NewClientProvider(suite.mockEntityProvider, suite.mockInboundClient)
	suite.ctx = context.Background()
}

func (suite *ClientProviderTestSuite) TestGetOAuthClientByClientID_Success() {
	oauth := &inboundmodel.OAuthClient{
		ID:                      "app-1",
		ClientID:                "client-1",
		RedirectURIs:            []string{"https://example.com/callback"},
		GrantTypes:              []oauth2const.GrantType{oauth2const.GrantTypeAuthorizationCode},
		TokenEndpointAuthMethod: oauth2const.TokenEndpointAuthMethodClientSecretBasic,
		Scopes:                  []string{"openid"},
		Token: &inboundmodel.OAuthTokenConfig{
			AccessToken: &inboundmodel.AccessTokenConfig{ValidityPeriod: 3600},
			IDToken:     &inboundmodel.IDTokenConfig{ValidityPeriod: 7200},
		},
	}
	inbound := &inboundmodel.InboundClient{
		ID:                        "app-1",
		AuthFlowID:                "auth-flow",
		RegistrationFlowID:        "reg-flow",
		RecoveryFlowID:            "rec-flow",
		IsRegistrationFlowEnabled: true,
		IsRecoveryFlowEnabled:     false,
		Properties:                map[string]interface{}{"logo_url": "https://example.com/logo.png"},
	}

	suite.mockInboundClient.EXPECT().GetOAuthClientByClientID(mock.Anything, "client-1").Return(oauth, nil)
	suite.mockInboundClient.EXPECT().GetInboundClientByEntityID(mock.Anything, "app-1").Return(inbound, nil)

	got, err := suite.provider.GetOAuthClientByClientID(suite.ctx, "client-1")

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), got)
	assert.Equal(suite.T(), "client-1", got.ClientID)
	assert.Equal(suite.T(), "app-1", got.EntityID)
	assert.Equal(suite.T(), []thunderidengine.GrantType{thunderidengine.GrantTypeAuthorizationCode}, got.GrantTypes)
	assert.Equal(suite.T(), thunderidengine.TokenEndpointAuthMethodClientSecretBasic, got.TokenEndpointAuthMethod)
	assert.Equal(suite.T(), int64(3600), got.AccessTokenValiditySeconds)
	assert.Equal(suite.T(), int64(7200), got.IDTokenValiditySeconds)
	assert.True(suite.T(), got.IsRegistrationFlowEnabled)
	assert.Equal(suite.T(), "auth-flow", got.AuthFlowID)
}

func (suite *ClientProviderTestSuite) TestGetOAuthClientByClientID_NotFound() {
	suite.mockInboundClient.EXPECT().GetOAuthClientByClientID(mock.Anything, "missing").Return(nil, nil)

	got, err := suite.provider.GetOAuthClientByClientID(suite.ctx, "missing")

	assert.ErrorIs(suite.T(), err, thunderidengine.ErrInboundClientNotFound)
	assert.Nil(suite.T(), got)
}

func (suite *ClientProviderTestSuite) TestGetTransitiveEntityGroups_Success() {
	suite.mockEntityProvider.EXPECT().GetTransitiveEntityGroups("user-1").Return([]entityprovider.EntityGroup{
		{ID: "g1", Name: "Admins"},
	}, nil)

	got, err := suite.provider.GetTransitiveEntityGroups(suite.ctx, "user-1")

	require.NoError(suite.T(), err)
	require.Len(suite.T(), got, 1)
	assert.Equal(suite.T(), "g1", got[0].ID)
	assert.Equal(suite.T(), "Admins", got[0].Name)
}

func (suite *ClientProviderTestSuite) TestGetApplicationByID_Success() {
	sysAttrs := []byte(`{"name":"Test App","description":"A test app"}`)
	entity := &entityprovider.Entity{
		ID:               "app-1",
		OUID:             "ou-1",
		SystemAttributes: sysAttrs,
	}
	inbound := &inboundmodel.InboundClient{
		ID:                        "app-1",
		IsRegistrationFlowEnabled: true,
		IsRecoveryFlowEnabled:     true,
		AuthFlowID:                "auth-flow",
		Properties: map[string]interface{}{
			"logo_url":   "https://example.com/logo.png",
			"url":        "https://example.com",
			"tos_uri":    "https://example.com/tos",
			"policy_uri": "https://example.com/policy",
		},
	}

	suite.mockEntityProvider.EXPECT().GetEntity("app-1").Return(entity, nil)
	suite.mockInboundClient.EXPECT().GetInboundClientByEntityID(mock.Anything, "app-1").Return(inbound, nil)

	got, err := suite.provider.GetApplicationByID(suite.ctx, "app-1")

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), got)
	assert.Equal(suite.T(), "app-1", got.ID)
	assert.Equal(suite.T(), "ou-1", got.OUID)
	assert.Equal(suite.T(), "Test App", got.Name)
	assert.Equal(suite.T(), "A test app", got.Description)
	assert.Equal(suite.T(), "https://example.com/logo.png", got.LogoURL)
	assert.True(suite.T(), got.IsRegistrationFlowEnabled)
}

func (suite *ClientProviderTestSuite) TestGetFlowApplicationByID_Success() {
	sysAttrs := []byte(`{"name":"Flow App","clientId":"oauth-client-1"}`)
	entity := &entityprovider.Entity{
		ID:               "app-1",
		SystemAttributes: sysAttrs,
	}
	inbound := &inboundmodel.InboundClient{
		ID:               "app-1",
		AuthFlowID:       "auth-flow",
		AllowedUserTypes: []string{"user"},
		Assertion:        &inboundmodel.AssertionConfig{ValidityPeriod: 300},
		Properties: map[string]interface{}{
			"metadata": map[string]interface{}{"tier": "premium"},
		},
	}

	suite.mockInboundClient.EXPECT().GetInboundClientByEntityID(mock.Anything, "app-1").Return(inbound, nil)
	suite.mockEntityProvider.EXPECT().GetEntity("app-1").Return(entity, nil)

	got, err := suite.provider.GetFlowApplicationByID(suite.ctx, "app-1")

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), got)
	assert.Equal(suite.T(), "app-1", got.ID)
	assert.Equal(suite.T(), "Flow App", got.Name)
	assert.Equal(suite.T(), "oauth-client-1", got.OAuthClientID)
	assert.Equal(suite.T(), int64(300), got.Assertion.ValidityPeriod)
	assert.Equal(suite.T(), "premium", got.Metadata["tier"])
}

func (suite *ClientProviderTestSuite) TestGetApplicationByID_NotFound() {
	suite.mockEntityProvider.EXPECT().GetEntity("missing").Return(nil, entityprovider.NewEntityProviderError(
		entityprovider.ErrorCodeEntityNotFound, "Not found", "entity not found"))
	suite.mockInboundClient.EXPECT().GetInboundClientByEntityID(mock.Anything, "missing").
		Return(nil, inboundclient.ErrInboundClientNotFound)

	got, err := suite.provider.GetApplicationByID(suite.ctx, "missing")

	assert.ErrorIs(suite.T(), err, thunderidengine.ErrApplicationNotFound)
	assert.Nil(suite.T(), got)
}
