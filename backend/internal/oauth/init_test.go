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

package oauth

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	yaml "gopkg.in/yaml.v3"

	"github.com/thunder-id/thunderid/internal/oauth/oauth2/clientprovidertest"
	"github.com/thunder-id/thunderid/internal/system/config"
	"github.com/thunder-id/thunderid/internal/system/cors"
	"github.com/thunder-id/thunderid/pkg/thunderidengine"
	"github.com/thunder-id/thunderid/tests/mocks/attributecachemock"
	"github.com/thunder-id/thunderid/tests/mocks/authnprovider/managermock"
	rbacauthzmock "github.com/thunder-id/thunderid/tests/mocks/authzmock"
	"github.com/thunder-id/thunderid/tests/mocks/crypto/cryptomock"
	"github.com/thunder-id/thunderid/tests/mocks/flow/flowexecmock"
	"github.com/thunder-id/thunderid/tests/mocks/idp/idpmock"
	"github.com/thunder-id/thunderid/tests/mocks/jose/jwtmock"
	"github.com/thunder-id/thunderid/tests/mocks/observability/observabilitymock"
	"github.com/thunder-id/thunderid/tests/mocks/oumock"
	"github.com/thunder-id/thunderid/tests/mocks/resourcemock"
)

type InitTestSuite struct {
	suite.Suite
	mockClientProvider        *clientprovidertest.ClientProviderMock
	mockAuthnProvider         *managermock.AuthnProviderManagerInterfaceMock
	mockJWTService            *jwtmock.JWTServiceMock
	mockFlowExecService       *flowexecmock.FlowExecServiceInterfaceMock
	mockObservabilityService  *observabilitymock.ObservabilityServiceInterfaceMock
	mockRuntimeCrypto         *cryptomock.RuntimeCryptoProviderMock
	mockOUService             *oumock.OrganizationUnitServiceInterfaceMock
	mockAttributeCacheService *attributecachemock.AttributeCacheServiceInterfaceMock
	mockRBACAuthzService      *rbacauthzmock.AuthorizationServiceInterfaceMock
	mockResourceService       *resourcemock.ResourceServiceInterfaceMock
	mockIDPService            *idpmock.IDPServiceInterfaceMock
}

func TestInitTestSuite(t *testing.T) {
	suite.Run(t, new(InitTestSuite))
}

func (suite *InitTestSuite) SetupTest() {
	var allowedOrigins cors.OriginEntries
	suite.Require().NoError(yaml.Unmarshal([]byte(`
- https://example.com
`), &allowedOrigins))
	testConfig := &config.Config{
		Server: config.ServerConfig{
			Identifier: "test",
			Hostname:   "localhost",
			Port:       8090,
		},
		JWT: config.JWTConfig{
			Issuer:         "https://localhost:8090",
			ValidityPeriod: 3600,
		},
		OAuth: config.OAuthConfig{
			AuthorizationCode: config.AuthorizationCodeConfig{ValidityPeriod: 600},
			RefreshToken:      config.RefreshTokenConfig{ValidityPeriod: 86400},
		},
		Database: config.DatabaseConfig{
			Config: config.DataSource{
				Type:   "sqlite",
				SQLite: config.SQLiteDataSource{Path: "test.db"},
			},
			Runtime: config.DataSource{
				Type:   "sqlite",
				SQLite: config.SQLiteDataSource{Path: "test.db"},
			},
		},
		GateClient: config.GateClientConfig{
			Scheme:    "https",
			Hostname:  "localhost",
			Port:      3000,
			LoginPath: "/login",
			ErrorPath: "/error",
		},
		CORS: config.CORSConfig{AllowedOrigins: allowedOrigins},
	}
	suite.Require().NoError(cors.InitializeMatcher(testConfig.CORS.AllowedOrigins))
	_ = config.InitializeServerRuntime("", testConfig)

	suite.mockClientProvider = clientprovidertest.NewClientProviderMock(suite.T())
	suite.mockAuthnProvider = managermock.NewAuthnProviderManagerInterfaceMock(suite.T())
	suite.mockJWTService = jwtmock.NewJWTServiceMock(suite.T())
	suite.mockFlowExecService = flowexecmock.NewFlowExecServiceInterfaceMock(suite.T())
	suite.mockObservabilityService = observabilitymock.NewObservabilityServiceInterfaceMock(suite.T())
	suite.mockRuntimeCrypto = cryptomock.NewRuntimeCryptoProviderMock(suite.T())
	suite.mockOUService = oumock.NewOrganizationUnitServiceInterfaceMock(suite.T())
	suite.mockAttributeCacheService = attributecachemock.NewAttributeCacheServiceInterfaceMock(suite.T())
	suite.mockRBACAuthzService = rbacauthzmock.NewAuthorizationServiceInterfaceMock(suite.T())
	suite.mockResourceService = resourcemock.NewResourceServiceInterfaceMock(suite.T())
	suite.mockIDPService = idpmock.NewIDPServiceInterfaceMock(suite.T())
}

func (suite *InitTestSuite) TearDownTest() {
	config.ResetServerRuntime()
}

func (suite *InitTestSuite) testOAuthOptions() thunderidengine.Options {
	cfg := config.GetServerRuntime().Config
	return thunderidengine.Options{
		Issuer:                    cfg.JWT.Issuer,
		Audience:                  cfg.JWT.Audience,
		ValidityPeriod:            cfg.JWT.ValidityPeriod,
		Leeway:                    cfg.JWT.Leeway,
		DeploymentID:              cfg.Server.Identifier,
		BaseURL:                   config.GetServerURL(&cfg.Server),
		RequirePAR:                cfg.OAuth.PAR.RequirePAR,
		PARExpiresIn:              cfg.OAuth.PAR.ExpiresIn,
		AllowWildcardRedirectURI:  cfg.OAuth.AllowWildcardRedirectURI,
		AuthorizationCodeValidity: cfg.OAuth.AuthorizationCode.ValidityPeriod,
		RefreshTokenValidity:      cfg.OAuth.RefreshToken.ValidityPeriod,
		RefreshTokenRenewOnGrant:  cfg.OAuth.RefreshToken.RenewOnGrant,
		AcrAMR:                    cfg.OAuth.AuthClass.AcrAMR,
		RuntimeStoreType:          cfg.Database.Runtime.Type,
		GateClient: thunderidengine.GateClientOptions{
			Scheme:    cfg.GateClient.Scheme,
			Hostname:  cfg.GateClient.Hostname,
			Port:      cfg.GateClient.Port,
			LoginPath: cfg.GateClient.LoginPath,
			ErrorPath: cfg.GateClient.ErrorPath,
		},
	}
}

func (suite *InitTestSuite) TestInitialize_InvalidOptions() {
	mux := http.NewServeMux()

	err := Initialize(
		mux, suite.mockClientProvider, suite.mockAuthnProvider, suite.mockJWTService, nil,
		suite.mockFlowExecService, suite.mockObservabilityService, suite.mockRuntimeCrypto,
		suite.mockOUService, suite.mockAttributeCacheService, suite.mockRBACAuthzService,
		suite.mockResourceService, suite.mockIDPService, thunderidengine.Options{},
	)

	assert.Error(suite.T(), err)
}

func (suite *InitTestSuite) TestInitialize() {
	mux := http.NewServeMux()

	err := Initialize(
		mux, suite.mockClientProvider, suite.mockAuthnProvider, suite.mockJWTService, nil,
		suite.mockFlowExecService, suite.mockObservabilityService, suite.mockRuntimeCrypto,
		suite.mockOUService, suite.mockAttributeCacheService, suite.mockRBACAuthzService,
		suite.mockResourceService, suite.mockIDPService, suite.testOAuthOptions(),
	)

	assert.NoError(suite.T(), err)
}

func (suite *InitTestSuite) TestInitialize_RegistersRoutes() {
	mux := http.NewServeMux()

	err := Initialize(
		mux, suite.mockClientProvider, suite.mockAuthnProvider, suite.mockJWTService, nil,
		suite.mockFlowExecService, suite.mockObservabilityService, suite.mockRuntimeCrypto,
		suite.mockOUService, suite.mockAttributeCacheService, suite.mockRBACAuthzService,
		suite.mockResourceService, suite.mockIDPService, suite.testOAuthOptions(),
	)
	assert.NoError(suite.T(), err)

	routes := []struct {
		method string
		path   string
	}{
		{method: "GET", path: "/oauth2/authorize"},
		{method: "POST", path: "/oauth2/par"},
		{method: "POST", path: "/oauth2/token"},
		{method: "GET", path: "/oauth2/userinfo"},
		{method: "GET", path: "/.well-known/oauth-authorization-server"},
		{method: "GET", path: "/oauth2/jwks"},
	}

	for _, route := range routes {
		suite.Run(route.method+" "+route.path, func() {
			_, pattern := mux.Handler(&http.Request{
				Method: route.method,
				URL:    &url.URL{Path: route.path},
			})
			assert.Contains(suite.T(), pattern, route.path)
		})
	}
}
