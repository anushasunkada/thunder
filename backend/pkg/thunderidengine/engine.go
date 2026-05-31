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

// Package thunderidengine provides the embeddable ThunderID OIDC authorization engine.
// Host contracts live in thunderidengine/host; ephemeral storage in thunderidengine/runtime.
// bridge.go maps those types to internal enginebridge packages to avoid import cycles.
package thunderidengine

import (
	"net/http"

	"github.com/thunder-id/thunderid/internal/flow/flowexec"
	serverconst "github.com/thunder-id/thunderid/internal/system/constants"
	"github.com/thunder-id/thunderid/internal/thunderidengineinit"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/host"
)

// Engine exposes initialized engine services.
type Engine = thunderidengineinit.Engine

// Initialize wires the embeddable ThunderID engine and registers public routes.
func Initialize(mux *http.ServeMux, cfg EngineConfig) (*Engine, error) {
	initCfg := thunderidengineinit.Config{
		Issuer:         cfg.Issuer,
		SigningKeyPath: cfg.Crypto.SigningKeyPath,
		DataDir:        cfg.DataDir,
		Actors:         WrapActorProvider(cfg.Actors),
		Runtime:        WrapRuntimeStore(cfg.Runtime),
		Authn:          host.InternalAuthnManager(host.InitializeAuthnProviderManager(cfg.Authn)),
		Authorization:  WrapAuthorization(cfg.Authorization),
		Consent:        WrapConsent(cfg.Consent),
		Flow: thunderidengineinit.FlowConfig{
			Executors: cfg.Flow.Executors,
		},
	}
	if cfg.FlowProvider != nil {
		initCfg.FlowSource = WrapFlowProvider(cfg.FlowProvider)
	} else {
		initCfg.FlowStore = thunderidengineinit.FlowProviderConfig{
			StoreMode:       serverconst.StoreMode(cfg.FlowStore.StoreMode),
			DefinitionsPath: cfg.FlowStore.DefinitionsPath,
		}
	}
	return thunderidengineinit.Initialize(mux, initCfg)
}

// FlowExec returns the flow execution service from an initialized engine.
func FlowExec(engine *Engine) flowexec.FlowExecServiceInterface {
	return engine.FlowExec()
}
