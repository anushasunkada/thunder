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

package declarative

import (
	"fmt"

	"github.com/thunder-id/thunderid/internal/application"
	layoutmgt "github.com/thunder-id/thunderid/internal/design/layout/mgt"
	"github.com/thunder-id/thunderid/internal/design/resolve"
	thememgt "github.com/thunder-id/thunderid/internal/design/theme/mgt"
	"github.com/thunder-id/thunderid/internal/entityprovider"
	"github.com/thunder-id/thunderid/internal/idp"
	"github.com/thunder-id/thunderid/internal/inboundclient"
	"github.com/thunder-id/thunderid/internal/ou"
	"github.com/thunder-id/thunderid/internal/resource"
	"github.com/thunder-id/thunderid/internal/role"
	"github.com/thunder-id/thunderid/internal/system/cache"
	"github.com/thunder-id/thunderid/internal/system/config"
	i18nmgt "github.com/thunder-id/thunderid/internal/system/i18n/mgt"
)

// Services holds declarative configuration-backed services for engine initialization.
type Services struct {
	Resource      resource.ResourceServiceInterface
	OU            ou.OrganizationUnitServiceInterface
	I18n          i18nmgt.I18nServiceInterface
	DesignResolve resolve.DesignResolveServiceInterface
	IDP           idp.IDPServiceInterface
	Role          role.RoleServiceInterface
}

// InitializeServices loads declarative YAML from dataDir/repository/resources and wires provider services.
func InitializeServices(
	dataDir string,
	cacheManager cache.CacheManagerInterface,
	inboundClient inboundclient.InboundClientServiceInterface,
	entityProvider entityprovider.EntityProviderInterface,
) (*Services, error) {
	if dataDir == "" {
		return nil, fmt.Errorf("data directory is required for declarative services")
	}
	if err := BootstrapDataDir(dataDir); err != nil {
		return nil, fmt.Errorf("bootstrap declarative data dir: %w", err)
	}

	ouService, err := ou.InitializeProviderWithAuthz(cacheManager, &engineSysAuthz{})
	if err != nil {
		return nil, fmt.Errorf("initialize OU service: %w", err)
	}

	resourceService, err := resource.InitializeProvider(ouService, nil)
	if err != nil {
		return nil, fmt.Errorf("initialize resource service: %w", err)
	}

	i18nService, err := i18nmgt.InitializeProvider(config.GetServerRuntime().Config.Translation)
	if err != nil {
		return nil, fmt.Errorf("initialize i18n service: %w", err)
	}

	idpService, err := idp.InitializeProvider(cacheManager)
	if err != nil {
		return nil, fmt.Errorf("initialize IDP service: %w", err)
	}

	themeService, err := thememgt.InitializeProvider()
	if err != nil {
		return nil, fmt.Errorf("initialize theme service: %w", err)
	}

	layoutService, err := layoutmgt.InitializeProvider()
	if err != nil {
		return nil, fmt.Errorf("initialize layout service: %w", err)
	}

	appService := application.InitializeProvider(inboundClient, entityProvider, ouService, i18nService)
	designResolve := resolve.InitializeProvider(themeService, layoutService, appService)

	roleService, err := role.InitializeProvider(ouService, resourceService)
	if err != nil {
		return nil, fmt.Errorf("initialize role service: %w", err)
	}

	return &Services{
		Resource:      resourceService,
		OU:            ouService,
		I18n:          i18nService,
		DesignResolve: designResolve,
		IDP:           idpService,
		Role:          roleService,
	}, nil
}
