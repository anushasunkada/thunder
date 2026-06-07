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
	sysconfig "github.com/thunder-id/thunderid/internal/system/config"
	"github.com/thunder-id/thunderid/internal/system/kmprovider/defaultkm/pki"
)

// BuildOptions overlays signing bootstrap fields when mapping system config into JOSE config.
type BuildOptions struct {
	SigningKeyPath string
	PreferredKeyID string
	JWKSCacheTTL   int64
}

// BuildOptionsForServer derives build options from a loaded server configuration.
func BuildOptionsForServer(cfg *sysconfig.Config) BuildOptions {
	return BuildOptions{
		PreferredKeyID: cfg.JWT.PreferredKeyID,
		JWKSCacheTTL:   int64(cfg.Server.SecurityConfig.JWKSCacheTTL),
	}
}

// FromSystemConfig maps system configuration into JOSE configuration.
func FromSystemConfig(cfg sysconfig.Config, opts BuildOptions) SystemConfig {
	preferredKeyID := opts.PreferredKeyID
	if preferredKeyID == "" {
		preferredKeyID = cfg.JWT.PreferredKeyID
	}
	if preferredKeyID == "" && opts.SigningKeyPath != "" {
		preferredKeyID = pki.DefaultEngineKeyID
	}

	validityPeriod := cfg.JWT.ValidityPeriod
	if validityPeriod == 0 {
		validityPeriod = defaultEngineJWTValiditySeconds
	}

	jwksCacheTTL := int64(cfg.Server.SecurityConfig.JWKSCacheTTL)
	if opts.JWKSCacheTTL != 0 {
		jwksCacheTTL = opts.JWKSCacheTTL
	}

	return SystemConfig{
		PreferredKeyID: preferredKeyID,
		Issuer:         cfg.JWT.Issuer,
		ValidityPeriod: validityPeriod,
		Leeway:         cfg.JWT.Leeway,
		JWKSCacheTTL:   jwksCacheTTL,
	}
}
