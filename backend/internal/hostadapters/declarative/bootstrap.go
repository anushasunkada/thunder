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
		GateClient: config.GateClientConfig{
			Scheme:    "http",
			Hostname:  "localhost",
			Port:      8080,
			LoginPath: "/signin",
			ErrorPath: "/error",
		},
		OAuth: config.OAuthConfig{
			PAR: config.PARConfig{
				RequirePAR: false,
			},
		},
		Crypto: config.CryptoConfig{
			Encryption: config.EncryptionConfig{
				Key: "0579f866ac7c9273580d0ff163fa01a7b2401a7ff3ddc3e3b14ae3136fa6025e",
			},
			PasswordHashing: config.PasswordHashingConfig{
				Algorithm: "PBKDF2",
			},
		},
	}
	return config.InitializeServerRuntime(dataDir, cfg)
}
