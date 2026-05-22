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
)

func TestLoadEngineConfigRequiresConfigPath(t *testing.T) {
	cfg, err := loadEngineConfig("")
	require.Error(t, err)
	require.Nil(t, cfg)
	require.Contains(t, err.Error(), "ConfigPath is required")

	cfg, err = loadEngineConfig("   ")
	require.Error(t, err)
	require.Nil(t, cfg)
	require.Contains(t, err.Error(), "ConfigPath is required")
}

func TestLoadEngineConfigDeploymentNotFound(t *testing.T) {
	cfg, err := loadEngineConfig("/nonexistent-server-home")
	require.Error(t, err)
	require.Nil(t, cfg)
	require.Contains(t, err.Error(), "deployment config not found")
}

func TestLoadEngineConfigSuccess(t *testing.T) {
	serverHome := filepath.Join("..", "..", "cmd", "server")
	deployment := filepath.Join(serverHome, "repository", "conf", "deployment.yaml")
	if _, err := os.Stat(deployment); err != nil {
		t.Skip("cmd/server deployment config not available")
	}
	ensureServerTestAssets(t, serverHome)

	cfg, err := loadEngineConfig(serverHome)
	require.NoError(t, err)
	require.NotNil(t, cfg)
}

func TestResolveServerHomeAndDeploymentFromServerHomeDir(t *testing.T) {
	serverHome := filepath.Join("tmp", "thunder-server")
	serverHome, deployment := resolveServerHomeAndDeployment(serverHome)
	require.Equal(t, filepath.Join("tmp", "thunder-server"), serverHome)
	require.Equal(t, filepath.Join("tmp", "thunder-server", "repository", "conf", "deployment.yaml"), deployment)
}

func TestResolveServerHomeAndDeploymentFromDeploymentYAML(t *testing.T) {
	deployment := filepath.Join("tmp", "thunder-server", "repository", "conf", "deployment.yaml")
	serverHome, resolved := resolveServerHomeAndDeployment(deployment)
	require.Equal(t, filepath.Join("tmp", "thunder-server"), serverHome)
	require.Equal(t, deployment, resolved)
}
