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
	"path/filepath"

	"github.com/thunder-id/thunderid/internal/system/config"
)

// LoadConfig loads the engine configuration from a server home directory, reading
// <serverHome>/deployment.yaml with defaults from <serverHome>/config/default.json. The returned
// *Config is passed to WithConfig. It is an SDK convenience so that embedders can seed the engine
// the same way the standalone server does without importing internal configuration packages.
func LoadConfig(serverHome string) (*Config, error) {
	return config.LoadConfig(
		filepath.Join(serverHome, "deployment.yaml"),
		filepath.Join(serverHome, "config", "default.json"),
		serverHome,
	)
}

// LoadConfigFromPaths loads the engine configuration from explicit deployment and default config
// file paths, resolving relative resource paths against serverHome. Use this when the embedder's
// configuration files are not laid out under the standard <serverHome>/deployment.yaml and
// <serverHome>/config/default.json convention.
func LoadConfigFromPaths(configPath, defaultPath, serverHome string) (*Config, error) {
	return config.LoadConfig(configPath, defaultPath, serverHome)
}
