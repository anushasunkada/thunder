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

package executor

import (
	"github.com/thunder-id/thunderid/internal/attributecache"
	"github.com/thunder-id/thunderid/internal/authn/assert"
	consentauthn "github.com/thunder-id/thunderid/internal/authn/consent"
	authnprovidermgr "github.com/thunder-id/thunderid/internal/authnprovider/manager"
	"github.com/thunder-id/thunderid/internal/authz"
	"github.com/thunder-id/thunderid/internal/entityprovider"
	"github.com/thunder-id/thunderid/internal/flow/core"
	"github.com/thunder-id/thunderid/internal/ou"
	"github.com/thunder-id/thunderid/internal/role"
	"github.com/thunder-id/thunderid/internal/system/jose/jwt"
)

// EngineDeps holds dependencies for engine-mode executor registration.
type EngineDeps struct {
	FlowFactory       core.FlowFactoryInterface
	EntityProvider    entityprovider.EntityProviderInterface
	AuthnProvider     authnprovidermgr.AuthnProviderManagerInterface
	AuthZService      authz.AuthorizationServiceInterface
	JWTService        jwt.JWTServiceInterface
	AuthAssertGen     assert.AuthAssertGeneratorInterface
	ConsentEnforcer   consentauthn.ConsentEnforcerServiceInterface
	OUService         ou.OrganizationUnitServiceInterface
	AttributeCacheSvc attributecache.AttributeCacheServiceInterface
	RoleService       role.RoleServiceInterface
	ExecutorNames     []string
}

// RegisterFromEngineDeps registers executors from engine dependencies.
func RegisterFromEngineDeps(deps EngineDeps) (ExecutorRegistryInterface, error) {
	return InitializeForEngine(RegisterDeps{
		FlowFactory:       deps.FlowFactory,
		EntityProvider:    deps.EntityProvider,
		AuthnProvider:     deps.AuthnProvider,
		AuthZService:      deps.AuthZService,
		JWTService:        deps.JWTService,
		AuthAssertGen:     deps.AuthAssertGen,
		ConsentEnforcer:   deps.ConsentEnforcer,
		OUService:         deps.OUService,
		AttributeCacheSvc: deps.AttributeCacheSvc,
		RoleService:       deps.RoleService,
	}, deps.ExecutorNames)
}
