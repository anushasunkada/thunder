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

package discovery

import (
	"net/http"
	"strings"

	kmprovider "github.com/thunder-id/thunderid/internal/system/kmprovider/common"
)

// EngineConfig holds discovery settings for engine mode.
type EngineConfig struct {
	Issuer      string
	PARRequired bool
	DPoPAlgs    []string
}

// InitializeForEngine initializes discovery using explicit issuer configuration.
func InitializeForEngine(
	mux *http.ServeMux,
	cfg EngineConfig,
	cryptoProvider kmprovider.RuntimeCryptoProvider,
) DiscoveryServiceInterface {
	baseURL := strings.TrimRight(cfg.Issuer, "/")
	discoveryService := &discoveryService{
		baseURL:        baseURL,
		cryptoProvider: cryptoProvider,
		issuerOverride: cfg.Issuer,
		parRequired:    cfg.PARRequired,
		dpopAlgs:       cfg.DPoPAlgs,
	}
	discoveryHandler := newDiscoveryHandler(discoveryService)
	registerRoutes(mux, discoveryHandler)
	return discoveryService
}
