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
	"strings"

	sysconfig "github.com/thunder-id/thunderid/internal/system/config"
)

var defaultEngineDPoPAlgs = []string{"ES256", "PS256", "ES384", "ES512", "EdDSA", "RS256"}

// BuildOptions overlays deployment-specific fields when mapping system config into OAuth2 config.
type BuildOptions struct {
	BaseURL          string
	DeploymentID     string
	RuntimeStoreType string
	GateClient       *sysconfig.GateClientConfig
}

// BuildOptionsForServer derives build options from a loaded server configuration.
func BuildOptionsForServer(cfg *sysconfig.Config) BuildOptions {
	gateClient := cfg.GateClient
	return BuildOptions{
		BaseURL:          sysconfig.GetServerURL(&cfg.Server),
		DeploymentID:     cfg.Server.Identifier,
		RuntimeStoreType: cfg.Database.Runtime.Type,
		GateClient:       &gateClient,
	}
}

// EngineBuildOptions derives build options for embeddable engine initialization.
func EngineBuildOptions(issuer string) BuildOptions {
	return BuildOptions{
		BaseURL: strings.TrimRight(strings.TrimSpace(issuer), "/"),
	}
}

// FromSystemConfig maps system configuration into OAuth2 configuration.
func FromSystemConfig(cfg sysconfig.Config, opts BuildOptions) Config {
	issuer := cfg.JWT.Issuer
	if opts.BaseURL != "" && issuer == "" {
		issuer = opts.BaseURL
	}

	baseURL := opts.BaseURL
	if baseURL == "" {
		baseURL = sysconfig.GetServerURL(&cfg.Server)
	}
	if baseURL == "" {
		baseURL = strings.TrimRight(issuer, "/")
	}

	return Config{
		BaseURL:           baseURL,
		Issuer:            issuer,
		DeploymentID:      opts.DeploymentID,
		RuntimeStoreType:  opts.RuntimeStoreType,
		JWT:               cfg.JWT,
		AuthorizationCode: cfg.OAuth.AuthorizationCode,
		RefreshToken:      cfg.OAuth.RefreshToken,
		PAR:               cfg.OAuth.PAR,
		DPoP:              cfg.OAuth.DPoP,
		DCR:               cfg.OAuth.DCR,
		AuthClass:         cfg.OAuth.AuthClass,
		GateClient:        opts.GateClient,
	}
}

// ApplyEngineDefaults fills unset OAuth/JWT fields for embeddable engine initialization.
func ApplyEngineDefaults(cfg *sysconfig.Config, issuer, audience string) {
	issuer = strings.TrimSpace(issuer)
	if cfg.JWT.Issuer == "" {
		cfg.JWT.Issuer = issuer
	}
	if cfg.JWT.Audience == "" {
		cfg.JWT.Audience = audience
	}
	if cfg.JWT.ValidityPeriod == 0 {
		cfg.JWT.ValidityPeriod = defaultEngineJWTValiditySeconds
	}
	if cfg.OAuth.AuthorizationCode.ValidityPeriod == 0 {
		cfg.OAuth.AuthorizationCode.ValidityPeriod = defaultEngineAuthCodeValiditySeconds
	}
	if cfg.OAuth.RefreshToken.ValidityPeriod == 0 {
		cfg.OAuth.RefreshToken.ValidityPeriod = 86400
	}
	if cfg.OAuth.PAR.ExpiresIn == 0 {
		cfg.OAuth.PAR.ExpiresIn = defaultEnginePARExpirySeconds
	}
	if cfg.OAuth.DPoP.IatWindow == 0 {
		cfg.OAuth.DPoP.IatWindow = 60
	}
	if cfg.OAuth.DPoP.Leeway == 0 {
		cfg.OAuth.DPoP.Leeway = 5
	}
	if cfg.OAuth.DPoP.MaxJTILength == 0 {
		cfg.OAuth.DPoP.MaxJTILength = 256
	}
	if len(cfg.OAuth.DPoP.AllowedAlgs) == 0 {
		cfg.OAuth.DPoP.AllowedAlgs = append([]string(nil), defaultEngineDPoPAlgs...)
	}
}
