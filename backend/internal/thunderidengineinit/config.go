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

// Package thunderidengineinit wires ThunderID server services into the public engine API.
package thunderidengineinit

import (
	authnprovidermgr "github.com/thunder-id/thunderid/internal/authnprovider/manager"
	"github.com/thunder-id/thunderid/internal/enginebridge"
	serverconst "github.com/thunder-id/thunderid/internal/system/constants"
)

// FlowProviderConfig configures flow definition storage.
type FlowProviderConfig struct {
	StoreMode       serverconst.StoreMode
	DefinitionsPath string
}

// FlowConfig holds flow executor registration settings.
type FlowConfig struct {
	Executors []string
}

// Config holds engine initialization inputs.
type Config struct {
	Issuer         string
	SigningKeyPath string
	DataDir        string
	FlowSource     enginebridge.FlowSource
	FlowStore      FlowProviderConfig
	Flow           FlowConfig
	Actors         enginebridge.ActorSource
	Runtime        enginebridge.RuntimeStore
	Authn          authnprovidermgr.AuthnProviderManagerInterface
	Authorization  enginebridge.AuthorizationSource
	Consent        enginebridge.ConsentSource
}
