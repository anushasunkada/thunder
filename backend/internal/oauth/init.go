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

	"github.com/thunder-id/thunderid/internal/attributecache"
	authnprovidermgr "github.com/thunder-id/thunderid/internal/authnprovider/manager"
	"github.com/thunder-id/thunderid/internal/authz"
	"github.com/thunder-id/thunderid/internal/flow/flowexec"
	"github.com/thunder-id/thunderid/internal/idp"
	"github.com/thunder-id/thunderid/internal/oauth/jwks"
	oauth2authz "github.com/thunder-id/thunderid/internal/oauth/oauth2/authz"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/discovery"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/granthandlers"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/introspect"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/jwksresolver"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/par"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/token"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/tokenservice"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/userinfo"
	"github.com/thunder-id/thunderid/internal/oauth/scope"
	"github.com/thunder-id/thunderid/internal/ou"
	"github.com/thunder-id/thunderid/internal/resource"
	"github.com/thunder-id/thunderid/internal/system/database/provider"
	syshttp "github.com/thunder-id/thunderid/internal/system/http"
	"github.com/thunder-id/thunderid/internal/system/observability"
	"github.com/thunder-id/thunderid/pkg/thunderidengine"
)

// Initialize initializes all OAuth-related services and registers their routes.
func Initialize(
	mux *http.ServeMux,
	clientProvider thunderidengine.ClientProvider,
	authnProvider authnprovidermgr.AuthnProviderManagerInterface,
	jwtService thunderidengine.JWTService,
	jweService thunderidengine.JWEService,
	flowExecService flowexec.FlowExecServiceInterface,
	observabilitySvc observability.ObservabilityServiceInterface,
	runtimeCrypto thunderidengine.RuntimeCryptoProvider,
	ouService ou.OrganizationUnitServiceInterface,
	attributeCacheSvc attributecache.AttributeCacheServiceInterface,
	authzService authz.AuthorizationServiceInterface,
	resourceService resource.ResourceServiceInterface,
	idpService idp.IDPServiceInterface,
	opts thunderidengine.Options,
) error {
	if err := opts.Validate(); err != nil {
		return err
	}

	transactioner, err := provider.GetDBProvider().GetRuntimeDBTransactioner()
	if err != nil {
		return err
	}

	tokenOpts := tokenservice.Options{
		Issuer:               opts.Issuer,
		Audience:             opts.Audience,
		ValidityPeriod:       opts.ValidityPeriod,
		Leeway:               opts.Leeway,
		RefreshTokenValidity: opts.RefreshTokenValidity,
	}
	discoveryOpts := discovery.Options{
		Issuer:     opts.Issuer,
		BaseURL:    opts.BaseURL,
		RequirePAR: opts.RequirePAR,
		AcrAMR:     opts.AcrAMR,
	}
	parOpts := par.Options{
		DeploymentID:     opts.DeploymentID,
		RuntimeStoreType: opts.RuntimeStoreType,
		PARExpiresIn:     opts.PARExpiresIn,
		OAuthPolicy:      opts.OAuthPolicy(),
	}
	authzOpts := oauth2authz.Options{
		Issuer:                    opts.Issuer,
		DeploymentID:              opts.DeploymentID,
		RuntimeStoreType:          opts.RuntimeStoreType,
		AuthorizationCodeValidity: opts.AuthorizationCodeValidity,
		GateClient: oauth2authz.GateClientOptions{
			Scheme:    opts.GateClient.Scheme,
			Hostname:  opts.GateClient.Hostname,
			Port:      opts.GateClient.Port,
			LoginPath: opts.GateClient.LoginPath,
			ErrorPath: opts.GateClient.ErrorPath,
		},
		OAuthPolicy:   opts.OAuthPolicy(),
		TokenDefaults: tokenOpts,
	}
	userInfoOpts := userinfo.Options{
		Issuer:         opts.Issuer,
		ValidityPeriod: opts.ValidityPeriod,
	}

	jwks.Initialize(mux, runtimeCrypto)
	httpClient := syshttp.NewHTTPClientWithCheckRedirect(func(req *http.Request, _ []*http.Request) error {
		return syshttp.IsSSRFSafeURL(req.URL.String())
	})
	resolver := jwksresolver.Initialize(httpClient)
	tokenBuilder, tokenValidator := tokenservice.Initialize(jwtService, jweService, resolver, idpService, tokenOpts)
	scopeValidator := scope.Initialize()
	discoveryService := discovery.Initialize(mux, runtimeCrypto, discoveryOpts)
	parService := par.Initialize(mux, clientProvider, authnProvider, jwtService, discoveryService,
		resourceService, parOpts)
	grantHandlerProvider, err := granthandlers.Initialize(
		mux, jwtService, clientProvider, flowExecService, tokenBuilder, tokenValidator,
		attributeCacheSvc, ouService, authzService, resourceService, parService,
		authzOpts, opts.RefreshTokenRenewOnGrant, tokenOpts)
	if err != nil {
		return err
	}
	token.Initialize(mux, jwtService, clientProvider, authnProvider, grantHandlerProvider,
		scopeValidator, observabilitySvc, discoveryService, transactioner)
	introspect.Initialize(mux, jwtService, clientProvider, authnProvider, discoveryService)
	userinfo.Initialize(mux, jwtService, jweService, resolver,
		tokenValidator, clientProvider, ouService, attributeCacheSvc, transactioner, userInfoOpts)
	return nil
}
