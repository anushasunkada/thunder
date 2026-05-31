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

package flowmgt

import (
	serverconst "github.com/thunder-id/thunderid/internal/system/constants"
	"github.com/thunder-id/thunderid/internal/system/transaction"
)

// FlowProviderConfig configures flow definition storage for provider-only initialization.
type FlowProviderConfig struct {
	StoreMode       serverconst.StoreMode
	DefinitionsPath string
	MutableStore    flowStoreInterface
	Transactioner   transaction.Transactioner
}

// ToFlowProviderConfig converts engine-level store config into flowmgt config.
func ToFlowProviderConfig(storeMode string, definitionsPath string,
	mutableStore flowStoreInterface, transactioner transaction.Transactioner) FlowProviderConfig {
	return FlowProviderConfig{
		StoreMode:       serverconst.StoreMode(storeMode),
		DefinitionsPath: definitionsPath,
		MutableStore:    mutableStore,
		Transactioner:   transactioner,
	}
}
