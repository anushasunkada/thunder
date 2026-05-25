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

package main

import (
	"github.com/thunder-id/thunderid/internal/system/config"
	"github.com/thunder-id/thunderid/pkg/thunderidengine"
)

func oauthOptionsFromConfig(cfg *config.Config) thunderidengine.Options {
	return thunderidengine.Options{
		Issuer:                    cfg.JWT.Issuer,
		Audience:                  cfg.JWT.Audience,
		ValidityPeriod:            cfg.JWT.ValidityPeriod,
		Leeway:                    cfg.JWT.Leeway,
		DeploymentID:              cfg.Server.Identifier,
		BaseURL:                   config.GetServerURL(&cfg.Server),
		RequirePAR:                cfg.OAuth.PAR.RequirePAR,
		PARExpiresIn:              cfg.OAuth.PAR.ExpiresIn,
		AllowWildcardRedirectURI:  cfg.OAuth.AllowWildcardRedirectURI,
		AuthorizationCodeValidity: cfg.OAuth.AuthorizationCode.ValidityPeriod,
		RefreshTokenValidity:      cfg.OAuth.RefreshToken.ValidityPeriod,
		RefreshTokenRenewOnGrant:  cfg.OAuth.RefreshToken.RenewOnGrant,
		AcrAMR:                    cfg.OAuth.AuthClass.AcrAMR,
		RuntimeStoreType:          cfg.Database.Runtime.Type,
		GateClient: thunderidengine.GateClientOptions{
			Scheme:    cfg.GateClient.Scheme,
			Hostname:  cfg.GateClient.Hostname,
			Port:      cfg.GateClient.Port,
			LoginPath: cfg.GateClient.LoginPath,
			ErrorPath: cfg.GateClient.ErrorPath,
		},
		DCRInsecure: cfg.OAuth.DCR.Insecure,
		Flow: thunderidengine.FlowOptions{
			DefaultAuthFlowHandle: cfg.Flow.DefaultAuthFlowHandle,
			AutoInferRegistration: cfg.Flow.AutoInferRegistration,
		},
	}
}
