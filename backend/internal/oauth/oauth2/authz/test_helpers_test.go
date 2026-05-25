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

package authz

import (
	"time"

	"github.com/thunder-id/thunderid/internal/oauth/oauth2/tokenservice"
	"github.com/thunder-id/thunderid/internal/system/config"
	"github.com/thunder-id/thunderid/pkg/thunderidengine"
)

func testHandlerOptions() Options {
	return Options{
		Issuer: "https://localhost:8090",
		GateClient: GateClientOptions{
			Scheme:    "https",
			Hostname:  "localhost",
			Port:      3000,
			LoginPath: "/login",
			ErrorPath: "/error",
		},
		AuthorizationCodeValidity: 600,
	}
}

func testInitOptions() Options {
	cfg := config.GetServerRuntime().Config
	return Options{
		Issuer:                    cfg.JWT.Issuer,
		DeploymentID:              cfg.Server.Identifier,
		RuntimeStoreType:          cfg.Database.Runtime.Type,
		AuthorizationCodeValidity: cfg.OAuth.AuthorizationCode.ValidityPeriod,
		GateClient: GateClientOptions{
			Scheme:    cfg.GateClient.Scheme,
			Hostname:  cfg.GateClient.Hostname,
			Port:      cfg.GateClient.Port,
			LoginPath: cfg.GateClient.LoginPath,
			ErrorPath: cfg.GateClient.ErrorPath,
		},
		OAuthPolicy: thunderidengine.OAuthPolicy{},
		TokenDefaults: tokenservice.Options{
			Issuer:               cfg.JWT.Issuer,
			ValidityPeriod:       cfg.JWT.ValidityPeriod,
			RefreshTokenValidity: cfg.OAuth.RefreshToken.ValidityPeriod,
		},
	}
}

func testTokenDefaultsFromRuntime() tokenservice.Options {
	cfg := config.GetServerRuntime().Config
	return tokenservice.Options{
		Issuer:               cfg.JWT.Issuer,
		ValidityPeriod:       cfg.JWT.ValidityPeriod,
		RefreshTokenValidity: cfg.OAuth.RefreshToken.ValidityPeriod,
	}
}

func testAuthCodeTTLFromRuntime() int64 {
	return config.GetServerRuntime().Config.OAuth.AuthorizationCode.ValidityPeriod
}

func testAuthCodeValidity() int64 {
	return testHandlerOptions().AuthorizationCodeValidity
}

func createTestAuthorizationCode(
	authRequestCtx *authRequestContext, clms *assertionClaims, authTime time.Time,
) (AuthorizationCode, error) {
	return createAuthorizationCode(authRequestCtx, clms, authTime, testAuthCodeValidity())
}

func resolveAttrCacheTTLForTest(app *thunderidengine.OAuthClient) int64 {
	return resolveUserAttributesCacheTTL(
		app, testTokenDefaultsFromRuntime(), testAuthCodeTTLFromRuntime(),
	)
}
