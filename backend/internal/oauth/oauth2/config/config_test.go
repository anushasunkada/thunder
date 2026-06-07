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
	"github.com/thunder-id/thunderid/internal/system/database/provider"
)

func TestFromSystemConfigEngineDefaults(t *testing.T) {
	t.Cleanup(Reset)

	sysCfg := sysconfig.Config{}
	ApplyEngineDefaults(&sysCfg, "https://as.example.com/", "my-aud")

	cfg := FromSystemConfig(sysCfg, EngineBuildOptions("https://as.example.com/"))

	if cfg.BaseURL != "https://as.example.com" {
		t.Fatalf("BaseURL = %q, want https://as.example.com", cfg.BaseURL)
	}
	if cfg.Issuer != "https://as.example.com/" {
		t.Fatalf("Issuer = %q", cfg.Issuer)
	}
	if cfg.JWT.ValidityPeriod != defaultEngineJWTValiditySeconds {
		t.Fatalf("JWT validity = %d, want %d", cfg.JWT.ValidityPeriod, defaultEngineJWTValiditySeconds)
	}
	if cfg.AuthorizationCode.ValidityPeriod != defaultEngineAuthCodeValiditySeconds {
		t.Fatalf("auth code TTL = %d", cfg.AuthorizationCode.ValidityPeriod)
	}
	if cfg.JWT.Audience != "my-aud" {
		t.Fatalf("audience = %q", cfg.JWT.Audience)
	}
	if cfg.PAR.ExpiresIn != defaultEnginePARExpirySeconds {
		t.Fatalf("PAR expires = %d", cfg.PAR.ExpiresIn)
	}
	if len(cfg.DPoP.AllowedAlgs) == 0 {
		t.Fatal("expected default DPoP algs")
	}
	if cfg.GateClient != nil {
		t.Fatal("engine config should not set gate client")
	}
}

func TestSetGetReset(t *testing.T) {
	t.Cleanup(Reset)

	want := Config{Issuer: "https://issuer.example.com", DeploymentID: "dep-1"}
	Set(want)
	got := Get()
	if got.Issuer != want.Issuer || got.DeploymentID != want.DeploymentID {
		t.Fatalf("Get() = %+v, want %+v", got, want)
	}
}

func TestFromServerRuntime(t *testing.T) {
	t.Cleanup(func() {
		Reset()
		sysconfig.ResetServerRuntime()
	})

	serverCfg := &sysconfig.Config{
		Server: sysconfig.ServerConfig{
			Hostname:   "localhost",
			Port:       8090,
			HTTPOnly:   true,
			Identifier: "test-dep",
		},
		Database: sysconfig.DatabaseConfig{
			Runtime: sysconfig.DataSource{Type: provider.DataSourceTypeRedis},
		},
		JWT: sysconfig.JWTConfig{
			Issuer:         "https://localhost:8090",
			ValidityPeriod: 7200,
			Audience:       "app",
			Leeway:         10,
		},
		OAuth: sysconfig.OAuthConfig{
			AuthorizationCode: sysconfig.AuthorizationCodeConfig{ValidityPeriod: 300},
			PAR:               sysconfig.PARConfig{RequirePAR: true, ExpiresIn: 90},
		},
		GateClient: sysconfig.GateClientConfig{
			Scheme:    "https",
			Hostname:  "gate.local",
			Port:      443,
			LoginPath: "/login",
		},
	}
	if err := sysconfig.InitializeServerRuntime("", serverCfg); err != nil {
		t.Fatalf("InitializeServerRuntime: %v", err)
	}

	cfg := FromServerRuntime()
	if cfg.DeploymentID != "test-dep" {
		t.Fatalf("DeploymentID = %q", cfg.DeploymentID)
	}
	if cfg.RuntimeStoreType != provider.DataSourceTypeRedis {
		t.Fatalf("RuntimeStoreType = %q", cfg.RuntimeStoreType)
	}
	if cfg.JWT.ValidityPeriod != 7200 {
		t.Fatalf("JWT validity = %d", cfg.JWT.ValidityPeriod)
	}
	if cfg.PAR.RequirePAR != true {
		t.Fatal("expected PAR required")
	}
	if cfg.GateClient == nil || cfg.GateClient.LoginPath != "/login" {
		t.Fatalf("GateClient = %+v", cfg.GateClient)
	}
}
