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

package role

import (
	"fmt"

	oupkg "github.com/thunder-id/thunderid/internal/ou"
	resourcepkg "github.com/thunder-id/thunderid/internal/resource"
)

// InitializeProvider creates a role service without registering HTTP routes.
func InitializeProvider(
	ouService oupkg.OrganizationUnitServiceInterface,
	resourceService resourcepkg.ResourceServiceInterface,
) (RoleServiceInterface, error) {
	roleStore, transactioner, fileStore, dbStore, err := initializeStore()
	if err != nil {
		return nil, err
	}

	roleService := newRoleService(roleStore, nil, nil, ouService, resourceService, transactioner)
	if fileStore != nil {
		if err := loadDeclarativeResources(fileStore, dbStore, roleService); err != nil {
			return nil, fmt.Errorf("failed to load declarative roles: %w", err)
		}
	}

	return roleService, nil
}
