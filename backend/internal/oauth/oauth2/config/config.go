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

// Package config holds OAuth2/OIDC authorization-server settings injected at initialize time.
package config

import (
	"sync"

	sysconfig "github.com/thunder-id/thunderid/internal/system/config"
)

const (
	defaultEngineJWTValiditySeconds      int64 = 3600
	defaultEngineAuthCodeValiditySeconds int64 = 600
	defaultEnginePARExpirySeconds        int64 = 60
)

// Config holds OAuth2/OIDC settings consumed by oauth2 subpackages.
type Config struct {
	BaseURL          string
	Issuer           string
	DeploymentID     string
	RuntimeStoreType string

	JWT sysconfig.JWTConfig

	AuthorizationCode sysconfig.AuthorizationCodeConfig
	RefreshToken      sysconfig.RefreshTokenConfig
	PAR               sysconfig.PARConfig
	DPoP              sysconfig.DPoPConfig
	DCR               sysconfig.DCRConfig
	AuthClass         sysconfig.AuthClassConfig

	GateClient *sysconfig.GateClientConfig
}

var (
	runtimeConfig *Config
	once          sync.Once
)

// Set stores OAuth2 configuration for the process. It may only be called once.
func Set(cfg Config) {
	once.Do(func() {
		c := cfg
		runtimeConfig = &c
	})
}

// Get returns the injected OAuth2 configuration.
func Get() Config {
	if runtimeConfig == nil {
		panic("OAuth2 config is not initialized")
	}
	return *runtimeConfig
}

// Reset clears injected OAuth2 configuration. For tests only.
func Reset() {
	runtimeConfig = nil
	once = sync.Once{}
}
