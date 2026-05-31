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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thunder-id/thunderid/internal/flow/common"
	"github.com/thunder-id/thunderid/internal/flow/core"
	"github.com/thunder-id/thunderid/internal/flow/executor"
	"github.com/thunder-id/thunderid/tests/mocks/flow/coremock"
)

func TestRegisterCustomExecutorsCallback(t *testing.T) {
	flowFactory := core.Initialize()
	reg, err := executor.RegisterFromEngineDeps(executor.EngineDeps{
		FlowFactory:   flowFactory,
		ExecutorNames: []string{executor.ExecutorNameBasicAuth},
	})
	require.NoError(t, err)

	custom := coremock.NewExecutorInterfaceMock(t)
	custom.On("GetName").Return("HostExecutor").Maybe()
	custom.On("GetType").Return(common.ExecutorTypeAuthentication).Maybe()

	reg.RegisterExecutor("HostExecutor", custom)

	assert.True(t, reg.IsRegistered("HostExecutor"))
	assert.True(t, reg.IsRegistered(executor.ExecutorNameBasicAuth))
}

func TestRegisterCustomExecutorsCallback_ErrorPropagates(t *testing.T) {
	flowFactory := core.Initialize()
	reg, err := executor.RegisterFromEngineDeps(executor.EngineDeps{
		FlowFactory: flowFactory,
	})
	require.NoError(t, err)

	err = func(_ executor.ExecutorRegistryInterface, _ core.FlowFactoryInterface) error {
		return fmt.Errorf("custom registration failed")
	}(reg, flowFactory)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "custom registration failed")
}
