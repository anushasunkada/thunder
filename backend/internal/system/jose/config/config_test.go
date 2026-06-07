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

package config

import (
	"testing"

	sysconfig "github.com/thunder-id/thunderid/internal/system/config"
	"github.com/thunder-id/thunderid/internal/system/kmprovider/defaultkm/pki"
)

func TestFromSystemConfigEngineDefaults(t *testing.T) {
	t.Cleanup(Reset)

	sysCfg := sysconfig.Config{
		JWT: sysconfig.JWTConfig{Issuer: "https://as.example.com"},
	}

	cfg := FromSystemConfig(sysCfg, BuildOptions{
		SigningKeyPath: "/keys/signing.pem",
	})

	if cfg.PreferredKeyID != pki.DefaultEngineKeyID {
		t.Fatalf("PreferredKeyID = %q", cfg.PreferredKeyID)
	}
	if cfg.ValidityPeriod != defaultEngineJWTValiditySeconds {
		t.Fatalf("ValidityPeriod = %d", cfg.ValidityPeriod)
	}
	if cfg.Issuer != "https://as.example.com" {
		t.Fatalf("Issuer = %q", cfg.Issuer)
	}
}

func TestFromServerRuntime(t *testing.T) {
	t.Cleanup(func() {
		Reset()
		sysconfig.ResetServerRuntime()
	})

	serverCfg := &sysconfig.Config{
		JWT: sysconfig.JWTConfig{
			PreferredKeyID: "key-1",
			Issuer:         "https://issuer.example.com",
			ValidityPeriod: 1800,
			Leeway:         15,
		},
		Server: sysconfig.ServerConfig{
			SecurityConfig: sysconfig.SecurityConfig{JWKSCacheTTL: 120},
		},
	}
	if err := sysconfig.InitializeServerRuntime("", serverCfg); err != nil {
		t.Fatalf("InitializeServerRuntime: %v", err)
	}

	cfg := FromServerRuntime()
	if cfg.PreferredKeyID != "key-1" || cfg.JWKSCacheTTL != 120 {
		t.Fatalf("cfg = %+v", cfg)
	}
}
