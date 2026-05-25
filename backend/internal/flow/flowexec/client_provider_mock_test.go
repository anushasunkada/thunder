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

package flowexec

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/thunder-id/thunderid/pkg/thunderidengine"
)

type clientProviderMock struct {
	mock.Mock
}

func (m *clientProviderMock) GetOAuthClientByClientID(
	ctx context.Context, clientID string,
) (*thunderidengine.OAuthClient, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*thunderidengine.OAuthClient), args.Error(1)
}

func (m *clientProviderMock) GetTransitiveEntityGroups(
	ctx context.Context, entityID string,
) ([]thunderidengine.EntityGroup, error) {
	args := m.Called(ctx, entityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]thunderidengine.EntityGroup), args.Error(1)
}

func (m *clientProviderMock) GetApplicationByID(
	ctx context.Context, appID string,
) (*thunderidengine.Application, error) {
	args := m.Called(ctx, appID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*thunderidengine.Application), args.Error(1)
}

func (m *clientProviderMock) GetFlowApplicationByID(
	ctx context.Context, appID string,
) (*thunderidengine.FlowApplication, error) {
	args := m.Called(ctx, appID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*thunderidengine.FlowApplication), args.Error(1)
}
