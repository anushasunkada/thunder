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

// Package declarative bootstraps declarative resource stores for engine mode.
package declarative

import (
	"github.com/thunder-id/thunderid/internal/system/config"
)

// BootstrapDataDir initializes server runtime state so declarative file stores resolve
// under dataDir/repository/resources.
func BootstrapDataDir(dataDir string) error {
	cfg := &config.Config{
		Server: config.ServerConfig{Identifier: "engine"},
		Resource: config.ResourceConfig{
			Store: "declarative",
		},
		OrganizationUnit: config.OrganizationUnitConfig{
			Store: "declarative",
		},
		IdentityProvider: config.IdentityProviderConfig{
			Store: "declarative",
		},
		Theme: config.ThemeConfig{
			Store: "declarative",
		},
		Layout: config.LayoutConfig{
			Store: "declarative",
		},
		Role: config.RoleConfig{
			Store: "declarative",
		},
		Translation: config.TranslationConfig{
			Store: "declarative",
		},
		DeclarativeResources: config.DeclarativeResources{
			Enabled: true,
		},
	}
	return config.InitializeServerRuntime(dataDir, cfg)
}
