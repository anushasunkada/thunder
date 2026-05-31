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

package ou

import (
	"github.com/thunder-id/thunderid/internal/system/cache"
	"github.com/thunder-id/thunderid/internal/system/sysauthz"
)

// InitializeProvider creates an OU service without registering HTTP routes.
func InitializeProvider(cacheManager cache.CacheManagerInterface) (OrganizationUnitServiceInterface, error) {
	ouStore, transactioner, err := initializeStore(cacheManager)
	if err != nil {
		return nil, err
	}
	return newOrganizationUnitService(nil, ouStore, transactioner), nil
}

// InitializeProviderWithAuthz creates an OU service with optional system authorization.
func InitializeProviderWithAuthz(
	cacheManager cache.CacheManagerInterface,
	authzService sysauthz.SystemAuthorizationServiceInterface,
) (OrganizationUnitServiceInterface, error) {
	ouStore, transactioner, err := initializeStore(cacheManager)
	if err != nil {
		return nil, err
	}
	return newOrganizationUnitService(authzService, ouStore, transactioner), nil
}
