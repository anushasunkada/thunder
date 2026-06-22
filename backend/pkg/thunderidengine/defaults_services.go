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

package thunderidengine

import (
	"fmt"
	"net/http"
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/thunder-id/thunderid/internal/application"
	"github.com/thunder-id/thunderid/internal/attributecache"
	attributecacheconfig "github.com/thunder-id/thunderid/internal/attributecache/config"
	"github.com/thunder-id/thunderid/internal/authz"
	"github.com/thunder-id/thunderid/internal/cert"
	certconfig "github.com/thunder-id/thunderid/internal/cert/config"
	"github.com/thunder-id/thunderid/internal/consent"
	layoutmgt "github.com/thunder-id/thunderid/internal/design/layout/mgt"
	"github.com/thunder-id/thunderid/internal/design/resolve"
	thememgt "github.com/thunder-id/thunderid/internal/design/theme/mgt"
	"github.com/thunder-id/thunderid/internal/entity"
	"github.com/thunder-id/thunderid/internal/entityprovider"
	"github.com/thunder-id/thunderid/internal/entitytype"
	flowcore "github.com/thunder-id/thunderid/internal/flow/core"
	"github.com/thunder-id/thunderid/internal/flow/executor"
	flowmgt "github.com/thunder-id/thunderid/internal/flow/mgt"
	"github.com/thunder-id/thunderid/internal/group"
	"github.com/thunder-id/thunderid/internal/idp"
	"github.com/thunder-id/thunderid/internal/inboundclient"
	"github.com/thunder-id/thunderid/internal/notification"
	"github.com/thunder-id/thunderid/internal/ou"
	"github.com/thunder-id/thunderid/internal/resource"
	"github.com/thunder-id/thunderid/internal/role"
	"github.com/thunder-id/thunderid/internal/system/cache"
	"github.com/thunder-id/thunderid/internal/system/config"
	"github.com/thunder-id/thunderid/internal/system/cryptolib"
	dbprovider "github.com/thunder-id/thunderid/internal/system/database/provider"
	i18nmgt "github.com/thunder-id/thunderid/internal/system/i18n/mgt"
	"github.com/thunder-id/thunderid/internal/system/jose/jwt"
	sysmcp "github.com/thunder-id/thunderid/internal/system/mcp"
	"github.com/thunder-id/thunderid/internal/system/sysauthz"
	"github.com/thunder-id/thunderid/internal/system/template"
	"github.com/thunder-id/thunderid/internal/user"
)

// declarativeBase carries the intermediate system-of-record services built in
// buildDeclarativeBaseServices that buildDeclarativeFlowAndDesign needs to finish the graph
// (after the flow executor registry has been built). The engine-required providers
// (ou/resource/idp/authz/attributecache) are recorded directly on engineConfig.
type declarativeBase struct {
	mcpServer      *mcpsdk.Server
	entityProvider entityprovider.EntityProviderInterface
	entityService  entity.EntityServiceInterface
	entityType     entitytype.EntityTypeServiceInterface
	consentService consent.ConsentServiceInterface
	themeService   thememgt.ThemeMgtServiceInterface
	layoutService  layoutmgt.LayoutMgtServiceInterface
	certService    cert.CertificateServiceInterface
}

// buildDeclarativeBaseServices constructs the system-of-record services from declarative
// resources and records the engine-required ones (ou, resource, idp, authz, attribute cache) on
// the config. All management REST routes the underlying services register are mounted on a
// throwaway mux so they are never exposed on the embedder's mux. It is a DB-free-at-the-boundary
// port of the first half of the standalone server's registerServices: the embedder supplies the
// datasource via WithConfig, and the DB provider initializes lazily from it.
func (c *engineConfig) buildDeclarativeBaseServices(
	cacheManager cache.CacheManagerInterface,
	jwtService jwt.JWTServiceInterface,
) (*declarativeBase, error) {
	mux := http.NewServeMux() // throwaway: discards all management routes
	mcpServer := sysmcp.Initialize(mux, jwtService)

	ouAuthz, err := sysauthz.Initialize()
	if err != nil {
		return nil, fmt.Errorf("thunderidengine: system authorization service: %w", err)
	}
	ouService, ouHierarchyResolver, _, err := ou.Initialize(mux, mcpServer, cacheManager, ouAuthz)
	if err != nil {
		return nil, fmt.Errorf("thunderidengine: organization unit service: %w", err)
	}
	ouAuthz.SetOUHierarchyResolver(ouHierarchyResolver)

	hashCfg, err := buildEngineHashConfig()
	if err != nil {
		return nil, err
	}
	hashService, err := cryptolib.Initialize(hashCfg)
	if err != nil {
		return nil, fmt.Errorf("thunderidengine: hash service: %w", err)
	}

	consentService := consent.Initialize()

	entityTypeService, _, err := entitytype.Initialize(
		mux, mcpServer, cacheManager, ouService, ouAuthz, consentService)
	if err != nil {
		return nil, fmt.Errorf("thunderidengine: entity type service: %w", err)
	}
	entityService, err := entity.Initialize(cacheManager, hashService, entityTypeService, ouService)
	if err != nil {
		return nil, fmt.Errorf("thunderidengine: entity service: %w", err)
	}
	entityProvider := entityprovider.InitializeEntityProvider(entityService)

	_, ouUserResolver, _, err := user.Initialize(mux, entityService, ouService, entityTypeService, ouAuthz)
	if err != nil {
		return nil, fmt.Errorf("thunderidengine: user service: %w", err)
	}
	groupService, ouGroupResolver, _, err := group.Initialize(
		mux, dbprovider.GetDBProvider(), ouService, entityService, entityTypeService, ouAuthz)
	if err != nil {
		return nil, fmt.Errorf("thunderidengine: group service: %w", err)
	}
	ouService.SetOUUserResolver(ouUserResolver)
	ouService.SetOUGroupResolver(ouGroupResolver)

	resourceService, _, err := resource.Initialize(mux, ouService, consentService)
	if err != nil {
		return nil, fmt.Errorf("thunderidengine: resource service: %w", err)
	}
	roleService, _, _, err := role.Initialize(
		mux, entityService, groupService, ouService, resourceService, entityTypeService)
	if err != nil {
		return nil, fmt.Errorf("thunderidengine: role service: %w", err)
	}
	authZService := authz.Initialize(roleService)

	idpService, _, err := idp.Initialize(cacheManager, mux, entityTypeService)
	if err != nil {
		return nil, fmt.Errorf("thunderidengine: idp service: %w", err)
	}

	templateService, err := template.Initialize()
	if err != nil {
		return nil, fmt.Errorf("thunderidengine: template service: %w", err)
	}
	if _, _, _, _, err = notification.Initialize(mux, jwtService, templateService); err != nil {
		return nil, fmt.Errorf("thunderidengine: notification service: %w", err)
	}

	attributeCacheService := attributecache.Initialize(attributecacheconfig.FromServerRuntime())

	themeMgtService, _, err := thememgt.Initialize(mux, mcpServer)
	if err != nil {
		return nil, fmt.Errorf("thunderidengine: theme service: %w", err)
	}
	layoutMgtService, _, err := layoutmgt.Initialize(mux)
	if err != nil {
		return nil, fmt.Errorf("thunderidengine: layout service: %w", err)
	}
	certService, err := cert.Initialize(cacheManager, dbprovider.GetDBProvider(), certconfig.FromServerRuntime())
	if err != nil {
		return nil, fmt.Errorf("thunderidengine: certificate service: %w", err)
	}

	// Record the engine-required providers.
	c.ouService = ouService
	c.resourceService = resourceService
	c.idpService = idpService
	c.authZService = authZService
	c.attributeCacheSvc = attributeCacheService

	return &declarativeBase{
		mcpServer:      mcpServer,
		entityProvider: entityProvider,
		entityService:  entityService,
		entityType:     entityTypeService,
		consentService: consentService,
		themeService:   themeMgtService,
		layoutService:  layoutMgtService,
		certService:    certService,
	}, nil
}

// buildDeclarativeFlowAndDesign finishes the declarative graph once the flow executor registry is
// available: it builds the flow management service (used as the default FlowProvider), the inbound
// client and application services, and the design resolve service. It records FlowProvider and
// DesignResolveService on the config. Management routes are again mounted on a throwaway mux.
func (c *engineConfig) buildDeclarativeFlowAndDesign(
	base *declarativeBase,
	cacheManager cache.CacheManagerInterface,
	flowFactory flowcore.FlowFactoryInterface,
	graphCache flowcore.GraphCacheInterface,
	execRegistry executor.ExecutorRegistryInterface,
	i18nService i18nmgt.I18nServiceInterface,
) error {
	mux := http.NewServeMux() // throwaway

	flowMgtService, _, err := flowmgt.Initialize(
		mux, base.mcpServer, cacheManager, flowFactory, execRegistry, graphCache)
	if err != nil {
		return fmt.Errorf("thunderidengine: flow management service: %w", err)
	}
	c.flowProvider = flowMgtService

	inboundClientService, err := inboundclient.Initialize(
		cacheManager, base.certService, base.entityProvider,
		base.themeService, base.layoutService, flowMgtService, base.entityType, base.consentService)
	if err != nil {
		return fmt.Errorf("thunderidengine: inbound client service: %w", err)
	}

	applicationService, _, err := application.Initialize(
		mux, base.mcpServer, base.entityProvider, base.entityService,
		inboundClientService, c.ouService, i18nService)
	if err != nil {
		return fmt.Errorf("thunderidengine: application service: %w", err)
	}

	c.designResolveService = resolve.Initialize(mux, base.themeService, base.layoutService, applicationService)
	return nil
}

// buildEngineHashConfig constructs a cryptolib.HashConfig from the seeded server runtime
// configuration (mirrors the standalone server's buildHashConfig).
func buildEngineHashConfig() (cryptolib.HashConfig, error) {
	cfg := config.GetServerRuntime().Config.Crypto.PasswordHashing
	alg := cryptolib.CredAlgorithm(strings.ToUpper(cfg.Algorithm))
	switch alg {
	case "", cryptolib.SHA256:
		return cryptolib.HashConfig{Algorithm: cryptolib.SHA256, SaltSize: cfg.SHA256.SaltSize}, nil
	case cryptolib.PBKDF2:
		return cryptolib.HashConfig{Algorithm: alg, SaltSize: cfg.PBKDF2.SaltSize,
			Iterations: cfg.PBKDF2.Iterations, KeySize: cfg.PBKDF2.KeySize}, nil
	case cryptolib.ARGON2ID:
		return cryptolib.HashConfig{Algorithm: alg, SaltSize: cfg.Argon2ID.SaltSize,
			Iterations: cfg.Argon2ID.Iterations, Memory: cfg.Argon2ID.Memory,
			Parallelism: cfg.Argon2ID.Parallelism, KeySize: cfg.Argon2ID.KeySize}, nil
	default:
		return cryptolib.HashConfig{}, fmt.Errorf("thunderidengine: unrecognized password hashing algorithm %q",
			cfg.Algorithm)
	}
}
