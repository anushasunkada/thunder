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

package par

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	yaml "gopkg.in/yaml.v3"

	"github.com/thunder-id/thunderid/internal/oauth/oauth2/clientprovidertest"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/discovery"
	"github.com/thunder-id/thunderid/internal/system/config"
	"github.com/thunder-id/thunderid/internal/system/cors"
	"github.com/thunder-id/thunderid/pkg/thunderidengine"
	"github.com/thunder-id/thunderid/tests/mocks/authnprovider/managermock"
	"github.com/thunder-id/thunderid/tests/mocks/jose/jwtmock"
	"github.com/thunder-id/thunderid/tests/mocks/oauth/oauth2/discoverymock"
	"github.com/thunder-id/thunderid/tests/mocks/resourcemock"
)

type InitTestSuite struct {
	suite.Suite
	mockClientProvider   *clientprovidertest.ClientProviderMock
	mockAuthnProvider    *managermock.AuthnProviderManagerInterfaceMock
	mockJWTService       *jwtmock.JWTServiceMock
	mockDiscoveryService *discoverymock.DiscoveryServiceInterfaceMock
	mockResourceService  *resourcemock.ResourceServiceInterfaceMock
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
			Hostname: "localhost",
			Port:     8090,
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
		CORS: config.CORSConfig{AllowedOrigins: allowedOrigins},
	}
	suite.Require().NoError(cors.InitializeMatcher(testConfig.CORS.AllowedOrigins))
	_ = config.InitializeServerRuntime("", testConfig)

	suite.mockClientProvider = clientprovidertest.NewClientProviderMock(suite.T())
	suite.mockAuthnProvider = managermock.NewAuthnProviderManagerInterfaceMock(suite.T())
	suite.mockJWTService = jwtmock.NewJWTServiceMock(suite.T())
	suite.mockDiscoveryService = discoverymock.NewDiscoveryServiceInterfaceMock(suite.T())
	suite.mockDiscoveryService.On("GetOAuth2AuthorizationServerMetadata", mock.Anything).
		Return(&discovery.OAuth2AuthorizationServerMetadata{
			PushedAuthorizationRequestEndpoint: "https://localhost:8090/oauth2/par",
		})
	suite.mockResourceService = resourcemock.NewResourceServiceInterfaceMock(suite.T())
}

func (suite *InitTestSuite) TearDownTest() {
	config.ResetServerRuntime()
}

func (suite *InitTestSuite) testInitOptions() Options {
	cfg := config.GetServerRuntime().Config
	return Options{
		DeploymentID:     cfg.Server.Identifier,
		RuntimeStoreType: cfg.Database.Runtime.Type,
		PARExpiresIn:     cfg.OAuth.PAR.ExpiresIn,
		OAuthPolicy:      thunderidengine.OAuthPolicy{},
	}
}

func (suite *InitTestSuite) TestInitialize() {
	mux := http.NewServeMux()

	service := Initialize(
		mux, suite.mockClientProvider, suite.mockAuthnProvider, suite.mockJWTService,
		suite.mockDiscoveryService, suite.mockResourceService, suite.testInitOptions(),
	)

	assert.NotNil(suite.T(), service)
	assert.Implements(suite.T(), (*PARServiceInterface)(nil), service)
}

func (suite *InitTestSuite) TestInitialize_RegistersRoutes() {
	mux := http.NewServeMux()

	Initialize(
		mux, suite.mockClientProvider, suite.mockAuthnProvider, suite.mockJWTService,
		suite.mockDiscoveryService, suite.mockResourceService, suite.testInitOptions(),
	)

	_, pattern := mux.Handler(&http.Request{Method: "POST", URL: &url.URL{Path: "/oauth2/par"}})
	assert.Contains(suite.T(), pattern, "/oauth2/par")

	req := httptest.NewRequest(http.MethodPost, "/oauth2/par", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	assert.NotEqual(suite.T(), http.StatusNotFound, rec.Code)
}
