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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thunder-id/thunderid/internal/flow/common"
	"github.com/thunder-id/thunderid/internal/flow/core"
	"github.com/thunder-id/thunderid/internal/flow/executor"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/flow"
	"github.com/thunder-id/thunderid/tests/mocks/flow/coremock"
)

func TestExecutorRegistry_RegisterCustomExecutor(t *testing.T) {
	innerReg, err := executor.InitializeForEngine(executor.RegisterDeps{
		FlowFactory: core.Initialize(),
	}, []string{executor.ExecutorNameBasicAuth})
	require.NoError(t, err)

	custom := coremock.NewExecutorInterfaceMock(t)
	custom.On("GetName").Return("CustomExecutor").Maybe()
	custom.On("GetType").Return(common.ExecutorTypeUtility).Maybe()

	reg := wrapExecutorRegistry(innerReg)
	reg.Register("CustomExecutor", custom)

	assert.True(t, reg.IsRegistered("CustomExecutor"))
	assert.True(t, reg.IsRegistered(executor.ExecutorNameBasicAuth))
	assert.False(t, reg.IsRegistered(executor.ExecutorNameOAuth))
}

func TestFlowFactory_CreateExecutor(t *testing.T) {
	factory := wrapFlowFactory(core.Initialize())
	base := factory.CreateExecutor("TestExecutor", flow.ExecutorTypeUtility, nil, nil)
	require.NotNil(t, base)
	assert.Equal(t, "TestExecutor", base.GetName())
	assert.Equal(t, common.ExecutorTypeUtility, base.GetType())
}

func TestRegisterCustomCallbackBridge(t *testing.T) {
	innerReg, err := executor.InitializeForEngine(executor.RegisterDeps{
		FlowFactory: core.Initialize(),
	}, []string{executor.ExecutorNameAuthAssert})
	require.NoError(t, err)

	custom := coremock.NewExecutorInterfaceMock(t)
	custom.On("GetName").Return("HostExecutor").Maybe()
	custom.On("GetType").Return(common.ExecutorTypeAuthentication).Maybe()

	reg := wrapExecutorRegistry(innerReg)
	factory := wrapFlowFactory(core.Initialize())
	reg.Register("HostExecutor", custom)
	require.NotNil(t, factory.CreateExecutor("unused", flow.ExecutorTypeUtility, nil, nil))
	assert.True(t, innerReg.IsRegistered("HostExecutor"))
}
