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

package flowexec

import (
	"errors"
	"net/http"

	"github.com/thunder-id/thunderid/internal/entityprovider"
	"github.com/thunder-id/thunderid/internal/flow/executor"
	"github.com/thunder-id/thunderid/internal/flow/flowbuilder"
	"github.com/thunder-id/thunderid/internal/inboundclient"
	kmprovider "github.com/thunder-id/thunderid/internal/system/kmprovider/common"
	"github.com/thunder-id/thunderid/internal/system/observability"
	"github.com/thunder-id/thunderid/internal/system/transaction"
)

// EngineDeps holds dependencies for engine-mode flow execution initialization.
type EngineDeps struct {
	FlowProvider         FlowProvider
	GraphBuilder         flowbuilder.GraphBuilderInterface
	FlowContextStore     RuntimeFlowContextStore
	InboundClientService inboundclient.InboundClientServiceInterface
	EntityProvider       entityprovider.EntityProviderInterface
	ExecutorRegistry     executor.ExecutorRegistryInterface
	ObservabilitySvc     observability.ObservabilityServiceInterface
	CryptoSvc            kmprovider.RuntimeCryptoProvider
	Transactioner        transaction.Transactioner
}

// InitializeForEngine creates flow execution services using injected runtime storage.
func InitializeForEngine(mux *http.ServeMux, deps EngineDeps) (FlowExecServiceInterface, error) {
	if deps.FlowContextStore == nil {
		return nil, errors.New("flow context store is required")
	}
	flowStore := asFlowStore(deps.FlowContextStore)
	transactioner := deps.Transactioner
	if transactioner == nil {
		transactioner = transaction.NewNoOpTransactioner()
	}
	flowEngine := newFlowEngine(deps.ExecutorRegistry, deps.ObservabilitySvc)
	flowExecService := newFlowExecService(deps.FlowProvider, deps.GraphBuilder, flowStore, flowEngine,
		deps.InboundClientService, deps.EntityProvider, deps.ObservabilitySvc, transactioner, deps.CryptoSvc)
	handler := newFlowExecutionHandler(flowExecService)
	registerRoutes(mux, handler)
	return flowExecService, nil
}

func asFlowStore(store RuntimeFlowContextStore) flowStoreInterface {
	return store
}
