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

// Package config holds JOSE settings injected at initialize time.
package config

import "sync"

const defaultEngineJWTValiditySeconds int64 = 3600

// SystemConfig holds JWT/JWE service settings.
type SystemConfig struct {
	PreferredKeyID string
	Issuer         string
	ValidityPeriod int64
	Leeway         int64
	JWKSCacheTTL   int64
}

var (
	runtimeConfig *SystemConfig
	once          sync.Once
)

// Set stores JOSE configuration for the process. It may only be called once.
func Set(cfg SystemConfig) {
	once.Do(func() {
		c := cfg
		runtimeConfig = &c
	})
}

// Get returns the injected JOSE configuration.
func Get() SystemConfig {
	if runtimeConfig == nil {
		panic("JOSE config is not initialized")
	}
	return *runtimeConfig
}

// Reset clears injected JOSE configuration. For tests only.
func Reset() {
	runtimeConfig = nil
	once = sync.Once{}
}
