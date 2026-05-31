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
	"context"

	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
	"github.com/thunder-id/thunderid/internal/system/security"
	"github.com/thunder-id/thunderid/internal/system/sysauthz"
)

type engineSysAuthz struct{}

func (e *engineSysAuthz) IsActionAllowed(ctx context.Context, action security.Action,
	actionCtx *sysauthz.ActionContext) (bool, *serviceerror.ServiceError) {
	_ = ctx
	_ = action
	_ = actionCtx
	return true, nil
}

func (e *engineSysAuthz) GetAccessibleResources(ctx context.Context, action security.Action,
	resourceType security.ResourceType) (*sysauthz.AccessibleResources, *serviceerror.ServiceError) {
	_ = ctx
	_ = action
	_ = resourceType
	return &sysauthz.AccessibleResources{AllAllowed: true}, nil
}

func (e *engineSysAuthz) SetOUHierarchyResolver(resolver sysauthz.OUHierarchyResolver) {
	_ = resolver
}
