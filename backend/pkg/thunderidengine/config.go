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
	"github.com/thunder-id/thunderid/pkg/thunderidengine/host"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/runtime"
)

// StoreMode selects flow definition storage behavior.
type StoreMode string

const (
	// StoreModeMutable uses a mutable flow definition store.
	StoreModeMutable StoreMode = "mutable"
	// StoreModeDeclarative uses a declarative flow definition store.
	StoreModeDeclarative StoreMode = "declarative"
	// StoreModeComposite uses a composite flow definition store.
	StoreModeComposite StoreMode = "composite"
)

// FlowProviderConfig configures flow definition storage for InitializeFlowProvider.
type FlowProviderConfig struct {
	StoreMode       StoreMode
	DefinitionsPath string
}

// OAuthConfig holds OAuth2/OIDC authorization server settings.
type OAuthConfig struct {
	AuthorizationCodeLifetimeSeconds int
	AccessTokenLifetimeSeconds       int
	RefreshTokenLifetimeSeconds      int
	PARExpirySeconds                 int
	DPoPRequired                     bool
}

// JWTConfig holds JWT signing and validation settings.
type JWTConfig struct {
	Issuer   string
	Audience string
	Leeway   int
}

// CryptoConfig holds cryptographic settings for the engine.
type CryptoConfig struct {
	SigningKeyPath string
}

// FlowConfig holds flow executor registration settings for the engine.
type FlowConfig struct {
	Executors []string
}

// EngineConfig configures the embeddable ThunderID engine.
type EngineConfig struct {
	Issuer   string
	Audience string
	JWKSPath string
	OAuth    OAuthConfig
	JWT      JWTConfig
	Crypto   CryptoConfig

	DataDir string

	FlowProvider host.FlowProvider
	FlowStore    FlowProviderConfig
	Flow         FlowConfig

	Actors        host.ActorProvider
	Runtime       runtime.Store
	Authn         host.AuthnProvider
	Authorization host.AuthorizationProvider
	Consent       host.ConsentEnforcer
}
