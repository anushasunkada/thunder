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

// Package authorization adapts ThunderID authorization services for engine host wiring.
package authorization

import (
	"context"
	"fmt"

	"github.com/thunder-id/thunderid/internal/authz"
	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/host"
)

type thunderAuthorizationProvider struct {
	authz authz.AuthorizationServiceInterface
}

// InitializeAuthorizationProvider creates a ThunderID AuthorizationProvider adapter.
func InitializeAuthorizationProvider(authzService authz.AuthorizationServiceInterface) host.AuthorizationProvider {
	return &thunderAuthorizationProvider{authz: authzService}
}

func (p *thunderAuthorizationProvider) GetAuthorizedPermissions(ctx context.Context,
	req host.GetAuthorizedPermissionsRequest) (*host.GetAuthorizedPermissionsResponse, error) {
	resp, svcErr := p.authz.GetAuthorizedPermissions(ctx, authz.GetAuthorizedPermissionsRequest{
		EntityID:             req.EntityID,
		RequestedPermissions: req.RequestedPermissions,
	})
	if svcErr != nil {
		return nil, asError(svcErr)
	}
	return &host.GetAuthorizedPermissionsResponse{
		AuthorizedPermissions: resp.AuthorizedPermissions,
	}, nil
}

func asError(svcErr *serviceerror.ServiceError) error {
	if svcErr == nil {
		return nil
	}
	return fmt.Errorf("%s", svcErr.ErrorDescription.DefaultValue)
}
