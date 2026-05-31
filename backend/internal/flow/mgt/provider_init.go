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
	"github.com/thunder-id/thunderid/internal/flow/common"
	"github.com/thunder-id/thunderid/internal/flow/flowbuilder"
	"github.com/thunder-id/thunderid/internal/system/cache"
	serverconst "github.com/thunder-id/thunderid/internal/system/constants"
	"github.com/thunder-id/thunderid/internal/system/transaction"
)

// InitializeFlowProvider creates a flow definition provider without HTTP/MCP/export routes.
func InitializeFlowProvider(
	cacheManager cache.CacheManagerInterface,
	graphBuilder flowbuilder.GraphBuilderInterface,
	cfg FlowProviderConfig,
) (FlowMgtServiceInterface, error) {
	store, compositeStore, transactioner, err := initializeStoreFromConfig(cacheManager, cfg)
	if err != nil {
		return nil, err
	}
	inferenceService := newFlowInferenceService()
	service := newFlowMgtService(store, inferenceService, graphBuilder, compositeStore, transactioner)
	return service, nil
}

func initializeStoreFromConfig(cacheManager cache.CacheManagerInterface, cfg FlowProviderConfig) (
	flowStoreInterface, *compositeFlowStore, transaction.Transactioner, error) {
	flowByIDCache := cache.GetCache[*common.CompleteFlowDefinition](cacheManager, "FlowByIDCache")
	flowByHandleCache := cache.GetCache[*common.CompleteFlowDefinition](cacheManager, "FlowByHandleCache")

	switch cfg.StoreMode {
	case serverconst.StoreModeComposite:
		fileStore, _ := newFileBasedStore()
		dbStore := cfg.MutableStore
		var transactioner transaction.Transactioner
		if cfg.MutableStore != nil {
			transactioner = cfg.Transactioner
		} else {
			var err error
			dbStore, transactioner, err = newCacheBackedFlowStore(flowByIDCache, flowByHandleCache)
			if err != nil {
				return nil, nil, nil, err
			}
		}
		compositeStore := newCompositeFlowStore(fileStore, dbStore)
		if err := loadDeclarativeResources(fileStore); err != nil {
			return nil, nil, nil, err
		}
		return compositeStore, compositeStore, transactioner, nil

	case serverconst.StoreModeDeclarative:
		fileStore, transactioner := newFileBasedStore()
		if err := loadDeclarativeResources(fileStore); err != nil {
			return nil, nil, nil, err
		}
		return fileStore, nil, transactioner, nil

	default:
		if cfg.MutableStore != nil {
			transactioner := cfg.Transactioner
			if transactioner == nil {
				transactioner = transaction.NewNoOpTransactioner()
			}
			return cfg.MutableStore, nil, transactioner, nil
		}
		store, transactioner, err := newCacheBackedFlowStore(flowByIDCache, flowByHandleCache)
		if err != nil {
			return nil, nil, nil, err
		}
		return store, nil, transactioner, nil
	}
}
