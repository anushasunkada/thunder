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

package resource

import (
	"fmt"

	"github.com/thunder-id/thunderid/internal/consent"
	oupkg "github.com/thunder-id/thunderid/internal/ou"
	serverconst "github.com/thunder-id/thunderid/internal/system/constants"
)

// InitializeProvider creates a resource service without registering HTTP routes.
func InitializeProvider(
	ouService oupkg.OrganizationUnitServiceInterface,
	consentService consent.ConsentServiceInterface,
) (ResourceServiceInterface, error) {
	resourceStore, transactioner, err := initializeStore()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize resource store: %w", err)
	}

	resourceService, err := newResourceService(ouService, consentService, resourceStore, transactioner)
	if err != nil {
		return nil, err
	}

	storeMode := getResourceStoreMode()
	if storeMode == serverconst.StoreModeDeclarative || storeMode == serverconst.StoreModeComposite {
		if err := loadDeclarativeResources(resourceStore, resourceService); err != nil {
			return nil, fmt.Errorf("failed to load declarative resources: %w", err)
		}
	}

	return resourceService, nil
}
