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
	"fmt"
	"net/http"

	"github.com/thunder-id/thunderid/internal/attributecache"
	authnprovidermgr "github.com/thunder-id/thunderid/internal/authnprovider/manager"
	"github.com/thunder-id/thunderid/internal/authz"
	"github.com/thunder-id/thunderid/internal/enginebridge"
	"github.com/thunder-id/thunderid/internal/entityprovider"
	"github.com/thunder-id/thunderid/internal/flow/flowexec"
	"github.com/thunder-id/thunderid/internal/idp"
	"github.com/thunder-id/thunderid/internal/inboundclient"
	"github.com/thunder-id/thunderid/internal/oauth/jwks"
	oauth2authz "github.com/thunder-id/thunderid/internal/oauth/oauth2/authz"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/discovery"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/dpop"
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
	syshttp "github.com/thunder-id/thunderid/internal/system/http"
	"github.com/thunder-id/thunderid/internal/system/jose/jwe"
	"github.com/thunder-id/thunderid/internal/system/jose/jwt"
	kmprovider "github.com/thunder-id/thunderid/internal/system/kmprovider/common"
	"github.com/thunder-id/thunderid/internal/system/observability"
	"github.com/thunder-id/thunderid/internal/system/transaction"
)

// EngineConfig holds OAuth settings for engine mode.
type EngineConfig struct {
	Issuer      string
	PARRequired bool
	DPoPAlgs    []string
}

// EngineDeps holds dependencies for engine-mode OAuth initialization.
type EngineDeps struct {
	RuntimeStore         enginebridge.RuntimeStore
	InboundClientService inboundclient.InboundClientServiceInterface
	AuthnProvider        authnprovidermgr.AuthnProviderManagerInterface
	AuthzService         authz.AuthorizationServiceInterface
	EntityProvider       entityprovider.EntityProviderInterface
	FlowExecService      flowexec.FlowExecServiceInterface
	ResourceService      resource.ResourceServiceInterface
	AttributeCache       attributecache.AttributeCacheServiceInterface
	OUService            ou.OrganizationUnitServiceInterface
	IDPService           idp.IDPServiceInterface
	CryptoProvider       kmprovider.RuntimeCryptoProvider
	Observability        observability.ObservabilityServiceInterface
	Config               EngineConfig
	Transactioner        transaction.Transactioner
}

// InitializeForEngine registers OAuth2/OIDC routes (excluding DCR) using injected runtime storage.
func InitializeForEngine(mux *http.ServeMux, deps EngineDeps) error {
	oauthStores := enginebridge.NewOAuthStores(deps.RuntimeStore)
	transactioner := deps.Transactioner
	if transactioner == nil {
		transactioner = transaction.NewNoOpTransactioner()
	}

	jwks.Initialize(mux, deps.CryptoProvider)
	httpClient := syshttp.NewHTTPClientWithCheckRedirect(func(req *http.Request, _ []*http.Request) error {
		return syshttp.IsSSRFSafeURL(req.URL.String())
	})
	resolver := jwksresolver.Initialize(httpClient)

	jwtService, err := jwt.Initialize(deps.CryptoProvider)
	if err != nil {
		return err
	}
	jweService, err := jwe.Initialize(deps.CryptoProvider)
	if err != nil {
		return err
	}

	idpService := deps.IDPService
	if idpService == nil {
		return fmt.Errorf("IDP service is required")
	}

	tokenBuilder, tokenValidator := tokenservice.Initialize(jwtService, jweService, resolver, idpService)
	scopeValidator := scope.Initialize()
	discoveryService := discovery.InitializeForEngine(mux, discovery.EngineConfig{
		Issuer:      deps.Config.Issuer,
		PARRequired: deps.Config.PARRequired,
		DPoPAlgs:    deps.Config.DPoPAlgs,
	}, deps.CryptoProvider)
	dpopVerifier := dpop.Initialize(oauthStores.JTI)

	resourceService := deps.ResourceService

	parService := par.InitializeWithStore(mux, deps.InboundClientService, deps.AuthnProvider, jwtService,
		discoveryService, resourceService, dpopVerifier, oauthStores.PAR)

	oauthAuthzService, err := oauth2authz.InitializeWithStores(
		mux, deps.InboundClientService, resourceService, jwtService,
		deps.FlowExecService, parService, oauthStores.AuthCode, oauthStores.AuthReq, transactioner)
	if err != nil {
		return err
	}

	attrCache := deps.AttributeCache
	if attrCache == nil {
		return fmt.Errorf("attribute cache service is required")
	}
	ouService := deps.OUService
	if ouService == nil {
		return fmt.Errorf("organization unit service is required")
	}

	grantHandlerProvider := granthandlers.InitializeForEngine(
		jwtService, oauthAuthzService, tokenBuilder, tokenValidator,
		attrCache, ouService, deps.AuthzService, deps.EntityProvider, resourceService,
	)

	token.Initialize(mux, jwtService, deps.InboundClientService, deps.AuthnProvider, grantHandlerProvider,
		scopeValidator, deps.Observability, discoveryService, dpopVerifier)
	introspect.Initialize(mux, jwtService, deps.InboundClientService, deps.AuthnProvider, discoveryService)
	userinfo.Initialize(mux, jwtService, jweService, resolver,
		tokenValidator, deps.InboundClientService, attrCache,
		discoveryService, dpopVerifier)
	return nil
}
