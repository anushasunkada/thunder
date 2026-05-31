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

package thunderidengineinit

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/thunder-id/thunderid/internal/attributecache"
	"github.com/thunder-id/thunderid/internal/design/resolve"
	"github.com/thunder-id/thunderid/internal/enginebridge"
	"github.com/thunder-id/thunderid/internal/flow/core"
	"github.com/thunder-id/thunderid/internal/flow/executor"
	"github.com/thunder-id/thunderid/internal/flow/flowbuilder"
	"github.com/thunder-id/thunderid/internal/flow/flowexec"
	"github.com/thunder-id/thunderid/internal/flow/flowmeta"
	flowmgt "github.com/thunder-id/thunderid/internal/flow/mgt"
	hostdeclarative "github.com/thunder-id/thunderid/internal/hostadapters/declarative"
	"github.com/thunder-id/thunderid/internal/oauth"
	"github.com/thunder-id/thunderid/internal/ou"
	"github.com/thunder-id/thunderid/internal/role"
	"github.com/thunder-id/thunderid/internal/system/cache"
	"github.com/thunder-id/thunderid/internal/system/config"
	i18nmgt "github.com/thunder-id/thunderid/internal/system/i18n/mgt"
)

// Engine wraps initialized engine services.
type Engine struct {
	flowExec flowexec.FlowExecServiceInterface
}

// FlowExec returns the flow execution service.
func (e *Engine) FlowExec() flowexec.FlowExecServiceInterface {
	return e.flowExec
}

// Initialize wires the embeddable ThunderID engine and registers public routes.
func Initialize(mux *http.ServeMux, cfg Config) (*Engine, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	flowSource, err := resolveFlowSource(cfg)
	if err != nil {
		return nil, err
	}

	cacheManager := cache.Initialize(config.CacheConfig{}, "engine")
	flowFactory := core.Initialize()
	entityProvider := enginebridge.NewEntityProvider(cfg.Actors)
	inboundClientService := enginebridge.NewInboundClientService(cfg.Actors)
	authzService := enginebridge.NewAuthzService(cfg.Authorization)

	cryptoProvider, err := initCryptoProvider(cfg.SigningKeyPath)
	if err != nil {
		return nil, err
	}

	jwtService, err := enginebridge.NewJWTService(cryptoProvider)
	if err != nil {
		return nil, err
	}

	attrCacheSvc := attributecache.NewServiceFromRuntimeStore(cfg.Runtime)

	var declarativeSvc *hostdeclarative.Services
	if cfg.DataDir != "" {
		declarativeSvc, err = hostdeclarative.InitializeServices(
			cfg.DataDir, cacheManager, inboundClientService, entityProvider)
		if err != nil {
			return nil, err
		}
	}

	execRegistry, err := executor.RegisterFromEngineDeps(executor.EngineDeps{
		FlowFactory:       flowFactory,
		EntityProvider:    entityProvider,
		AuthnProvider:     cfg.Authn,
		AuthZService:      authzService,
		JWTService:        jwtService,
		AuthAssertGen:     enginebridge.NewAuthAssertGenerator(),
		ConsentEnforcer:   enginebridge.NewConsentEnforcer(cfg.Consent),
		AttributeCacheSvc: attrCacheSvc,
		ExecutorNames:     cfg.Flow.Executors,
		OUService:         declarativeOU(declarativeSvc),
		RoleService:       declarativeRole(declarativeSvc),
	})
	if err != nil {
		return nil, err
	}

	graphBuilder := flowbuilder.Initialize(cacheManager, flowFactory, execRegistry)

	flowExec, err := flowexec.InitializeForEngine(mux, flowexec.EngineDeps{
		FlowProvider:         enginebridge.NewFlowExecProvider(flowSource),
		GraphBuilder:         graphBuilder,
		FlowContextStore:     enginebridge.NewFlowContextStore(cfg.Runtime),
		InboundClientService: inboundClientService,
		EntityProvider:       entityProvider,
		ExecutorRegistry:     execRegistry,
		CryptoSvc:            cryptoProvider,
	})
	if err != nil {
		return nil, err
	}

	flowmeta.InitializeForEngine(mux, flowmeta.EngineDeps{
		InboundClientService: inboundClientService,
		EntityProvider:       entityProvider,
		OUService:            declarativeOU(declarativeSvc),
		DesignResolve:        declarativeDesign(declarativeSvc),
		I18nService:          declarativeI18n(declarativeSvc),
	})

	oauthDeps := oauth.EngineDeps{
		RuntimeStore:         cfg.Runtime,
		InboundClientService: inboundClientService,
		AuthnProvider:        cfg.Authn,
		AuthzService:         authzService,
		EntityProvider:       entityProvider,
		FlowExecService:      flowExec,
		AttributeCache:       attrCacheSvc,
		CryptoProvider:       cryptoProvider,
		Config: oauth.EngineConfig{
			Issuer: cfg.Issuer,
		},
	}
	if declarativeSvc != nil {
		oauthDeps.ResourceService = declarativeSvc.Resource
		oauthDeps.OUService = declarativeSvc.OU
		oauthDeps.IDPService = declarativeSvc.IDP
	}
	if err := oauth.InitializeForEngine(mux, oauthDeps); err != nil {
		return nil, err
	}

	return &Engine{flowExec: flowExec}, nil
}

func validateConfig(cfg Config) error {
	if cfg.Actors == nil {
		return errors.New("Actors provider is required")
	}
	if cfg.Runtime == nil {
		return errors.New("Runtime store is required")
	}
	if cfg.Authn == nil {
		return errors.New("Authn provider is required")
	}
	if cfg.Authorization == nil {
		return errors.New("Authorization provider is required")
	}
	if cfg.Consent == nil {
		return errors.New("Consent enforcer is required")
	}
	if cfg.Issuer == "" {
		return errors.New("Issuer is required")
	}
	if cfg.FlowSource == nil && cfg.FlowStore.StoreMode == "" {
		return errors.New("FlowProvider or FlowStore config is required")
	}
	if cfg.DataDir == "" {
		return errors.New("DataDir is required for declarative engine services")
	}
	return nil
}

func resolveFlowSource(cfg Config) (enginebridge.FlowSource, error) {
	if cfg.FlowSource != nil {
		return cfg.FlowSource, nil
	}
	cacheManager := cache.Initialize(config.CacheConfig{}, "engine")
	flowFactory := core.Initialize()
	graphBuilder := flowbuilder.Initialize(cacheManager, flowFactory, nil)
	service, err := flowmgt.InitializeFlowProvider(cacheManager, graphBuilder, flowmgt.FlowProviderConfig{
		StoreMode:       cfg.FlowStore.StoreMode,
		DefinitionsPath: cfg.FlowStore.DefinitionsPath,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize flow provider: %w", err)
	}
	return enginebridge.NewFlowSource(service), nil
}

func declarativeOU(svc *hostdeclarative.Services) ou.OrganizationUnitServiceInterface {
	if svc == nil {
		return nil
	}
	return svc.OU
}

func declarativeDesign(svc *hostdeclarative.Services) resolve.DesignResolveServiceInterface {
	if svc == nil {
		return nil
	}
	return svc.DesignResolve
}

func declarativeI18n(svc *hostdeclarative.Services) i18nmgt.I18nServiceInterface {
	if svc == nil {
		return nil
	}
	return svc.I18n
}

func declarativeRole(svc *hostdeclarative.Services) role.RoleServiceInterface {
	if svc == nil {
		return nil
	}
	return svc.Role
}
