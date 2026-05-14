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

// Package oauth provides centralized initialization for all OAuth-related services.
package oauth

import (
	"net/http"

	"github.com/thunder-id/thunderid/internal/oauth/hostbridge"
	"github.com/thunder-id/thunderid/internal/oauth/jwks"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/dcr"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/discovery"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/granthandlers"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/introspect"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/jwksresolver"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/par"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/token"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/tokenservice"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/userinfo"
	"github.com/thunder-id/thunderid/internal/oauth/scope"
	syshttp "github.com/thunder-id/thunderid/internal/system/http"
	oauthdeps "github.com/thunder-id/thunderid/pkg/oauth/deps"
)

// Initialize initializes all OAuth-related services and registers their routes.
func Initialize(mux *http.ServeMux, deps oauthdeps.Dependencies) error {
	transactioner := deps.Transactioner
	inboundInternal := hostbridge.InboundFromHost(deps.Inbound)
	dcrPartner := hostbridge.DCRPartnerFromApplication(deps.Application)

	jwks.Initialize(mux, deps.PKIService)
	httpClient := syshttp.NewHTTPClientWithCheckRedirect(func(req *http.Request, _ []*http.Request) error {
		return syshttp.IsSSRFSafeURL(req.URL.String())
	})
	resolver := jwksresolver.Initialize(httpClient)
	tokenBuilder, tokenValidator := tokenservice.Initialize(deps.JWTService, deps.JWEService, resolver, deps.IDPService)
	scopeValidator := scope.Initialize()
	discoveryService := discovery.Initialize(mux, deps.PKIService)
	parService := par.Initialize(mux, inboundInternal, deps.AuthnProvider, deps.JWTService, discoveryService,
		deps.ResourceService, deps.DBProvider, deps.RedisProvider, deps.DeploymentID, deps.DatabaseRuntimeType)
	grantHandlerProvider, err := granthandlers.Initialize(
		mux, deps.JWTService, inboundInternal, deps.FlowExecService, tokenBuilder, tokenValidator,
		deps.AttributeCacheSvc, deps.OUService, deps.AuthzService, deps.EntityProvider, deps.ResourceService, parService)
	if err != nil {
		return err
	}
	token.Initialize(mux, deps.JWTService, inboundInternal, deps.AuthnProvider, grantHandlerProvider,
		scopeValidator, deps.ObservabilitySvc, discoveryService, transactioner)
	introspect.Initialize(mux, deps.JWTService, inboundInternal, deps.AuthnProvider, discoveryService)
	userinfo.Initialize(mux, deps.JWTService, deps.JWEService, resolver,
		tokenValidator, inboundInternal, deps.OUService, deps.AttributeCacheSvc, transactioner)
	dcr.Initialize(mux, dcrPartner, deps.OUService, deps.I18nService, transactioner)
	return nil
}
