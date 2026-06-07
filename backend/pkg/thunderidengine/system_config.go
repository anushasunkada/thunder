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
	oauth2config "github.com/thunder-id/thunderid/internal/oauth/oauth2/config"
	sysconfig "github.com/thunder-id/thunderid/internal/system/config"
)

// BuildSystemConfig maps embedder EngineConfig into the shared system configuration shape.
func BuildSystemConfig(cfg EngineConfig) sysconfig.Config {
	audience := cfg.Audience
	if audience == "" {
		audience = cfg.JWT.Audience
	}

	sysCfg := sysconfig.Config{
		JWT: sysconfig.JWTConfig{
			Issuer:         cfg.JWT.Issuer,
			Audience:       audience,
			Leeway:         cfg.JWT.Leeway,
			ValidityPeriod: int64(cfg.OAuth.AccessTokenLifetimeSeconds),
		},
		OAuth: sysconfig.OAuthConfig{
			AuthorizationCode: sysconfig.AuthorizationCodeConfig{
				ValidityPeriod: int64(cfg.OAuth.AuthorizationCodeLifetimeSeconds),
			},
			RefreshToken: sysconfig.RefreshTokenConfig{
				RenewOnGrant:   cfg.OAuth.RefreshTokenRenewOnGrant,
				ValidityPeriod: int64(cfg.OAuth.RefreshTokenLifetimeSeconds),
			},
			PAR: sysconfig.PARConfig{
				RequirePAR: cfg.OAuth.PARRequired,
				ExpiresIn:  int64(cfg.OAuth.PARExpirySeconds),
			},
			DPoP: sysconfig.DPoPConfig{
				Required:     cfg.OAuth.DPoPRequired,
				IatWindow:    cfg.OAuth.DPoPIatWindow,
				Leeway:       cfg.OAuth.DPoPLeeway,
				AllowedAlgs:  append([]string(nil), cfg.OAuth.DPoPAllowedAlgs...),
				MaxJTILength: cfg.OAuth.DPoPMaxJTILength,
			},
		},
	}
	oauth2config.ApplyEngineDefaults(&sysCfg, cfg.Issuer, audience)
	return sysCfg
}
