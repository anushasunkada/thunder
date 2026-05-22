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

package enginebridge

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thunder-id/thunderid/internal/flow/common"
	"github.com/thunder-id/thunderid/internal/flow/executor"
	"github.com/thunder-id/thunderid/internal/system/cache"
	"github.com/thunder-id/thunderid/pkg/thunderidengine"
	"github.com/thunder-id/thunderid/tests/mocks/flow/coremock"
)

func TestInitializeHostOnlyRequiresProviders(t *testing.T) {
	hostOnly := true
	err := initializeHostOnly(thunderidengine.EngineConfig{
		HostOnly:   &hostOnly,
		Providers:  thunderidengine.Providers{Client: &testClientProvider{}},
		ConfigPath: "/nonexistent",
	}, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "host-only mode requires providers")
}

func TestHostOnlyEnabledAutoDetectsCompleteProviders(t *testing.T) {
	cfg := thunderidengine.EngineConfig{Providers: testProviders()}
	require.True(t, cfg.HostOnlyEnabled())
}

func TestInitializeHostOnlyRequiresConfigPath(t *testing.T) {
	hostOnly := true
	err := initializeHostOnly(thunderidengine.EngineConfig{
		HostOnly:   &hostOnly,
		Providers:  testProviders(),
		ConfigPath: "",
	}, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "ConfigPath is required")
}

func TestBuildHostExecutorRegistryUsesCustomRegistry(t *testing.T) {
	customReg := executor.NewExecutorRegistry()
	customEx := coremock.NewExecutorInterfaceMock(t)
	customEx.On("GetName").Return("CustomExecutor").Maybe()

	reg := buildHostExecutorRegistry(
		thunderidengine.ExecutorConfig{
			CustomRegistry: customReg,
			InjectCustom:   []thunderidengine.ExecutorInterface{customEx},
		},
		testProviders(),
		cache.Initialize(),
		nil,
		nil,
	)

	require.True(t, reg.IsRegistered("CustomExecutor"))
	_, err := reg.GetExecutor("CustomExecutor")
	require.NoError(t, err)
	require.True(t, customReg.IsRegistered("CustomExecutor"))
}

func TestBuildHostExecutorRegistryRegistersDefaultExecutors(t *testing.T) {
	cacheManager, infra, runtime := hostExecutorTestDeps(t)

	reg := buildHostExecutorRegistry(
		thunderidengine.ExecutorConfig{
			Names: []string{"BasicAuthExecutor", "AuthAssertExecutor"},
		},
		testProviders(),
		cacheManager,
		infra,
		runtime,
	)

	require.True(t, reg.IsRegistered("BasicAuthExecutor"))
	require.True(t, reg.IsRegistered("AuthAssertExecutor"))
	require.False(t, reg.IsRegistered("AuthorizationExecutor"))
}

func TestBuildHostExecutorRegistryInjectsCustomOnDefaultRegistry(t *testing.T) {
	cacheManager, infra, runtime := hostExecutorTestDeps(t)
	customEx := coremock.NewExecutorInterfaceMock(t)
	customEx.On("GetName").Return("InjectedExecutor").Maybe()
	customEx.On("GetType").Return(common.ExecutorTypeUtility).Maybe()

	reg := buildHostExecutorRegistry(
		thunderidengine.ExecutorConfig{
			Names:        []string{"BasicAuthExecutor"},
			InjectCustom: []thunderidengine.ExecutorInterface{customEx},
		},
		testProviders(),
		cacheManager,
		infra,
		runtime,
	)

	require.True(t, reg.IsRegistered("BasicAuthExecutor"))
	require.True(t, reg.IsRegistered("InjectedExecutor"))
}

func hostExecutorTestDeps(t *testing.T) (cache.CacheManagerInterface, *leanInfra, *runtimeServices) {
	t.Helper()

	serverHome := filepath.Join("..", "..", "cmd", "server")
	deployment := filepath.Join(serverHome, "repository", "conf", "deployment.yaml")
	if _, err := os.Stat(deployment); err != nil {
		t.Skip("cmd/server deployment config not available")
	}
	ensureServerTestAssets(t, serverHome)

	thunderCfg, err := loadEngineConfig(serverHome)
	require.NoError(t, err)

	cacheManager, err := initPlatform(thunderCfg)
	require.NoError(t, err)

	infra, err := initLeanInfra()
	require.NoError(t, err)

	runtime, err := buildRuntimeFromHostProviders(testProviders(), infra)
	require.NoError(t, err)

	return cacheManager, infra, runtime
}
