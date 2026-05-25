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

// Package flowgraphbridge wires thunderidengine flow graph types to Thunder runtime implementations.
package flowgraphbridge

import (
	"context"
	"errors"

	"github.com/thunder-id/thunderid/internal/flow/common"
	flowcore "github.com/thunder-id/thunderid/internal/flow/core"
	"github.com/thunder-id/thunderid/internal/flow/executor"
	flowmgt "github.com/thunder-id/thunderid/internal/flow/mgt"
	"github.com/thunder-id/thunderid/internal/system/cache"
	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
	"github.com/thunder-id/thunderid/pkg/thunderidengine"
)

type runtimeFlowGraph struct {
	graph flowcore.GraphInterface
}

func (g *runtimeFlowGraph) GetID() string {
	return g.graph.GetID()
}

func (g *runtimeFlowGraph) FlowType() string {
	return string(g.graph.GetType())
}

// WrapCoreGraph wraps an internal graph as a FlowGraph for tests and host wiring.
func WrapCoreGraph(graph flowcore.GraphInterface) thunderidengine.FlowGraph {
	return &runtimeFlowGraph{graph: graph}
}

// CoreGraphFromFlowGraph returns the internal graph for flow execution.
func CoreGraphFromFlowGraph(fg thunderidengine.FlowGraph) (flowcore.GraphInterface, bool) {
	if g, ok := fg.(*runtimeFlowGraph); ok {
		return g.graph, true
	}
	return nil, false
}

type graphCacheHolder struct {
	cache flowcore.GraphCacheInterface
}

func (h *graphCacheHolder) inner() flowcore.GraphCacheInterface {
	return h.cache
}

// NewInMemoryGraphCache returns a graph cache backed by the host in-memory cache manager.
func NewInMemoryGraphCache() thunderidengine.GraphCache {
	cm := cache.Initialize()
	_, graphCache := flowcore.Initialize(cm)
	return &graphCacheHolder{cache: graphCache}
}

type graphBuilderBridge struct {
	inner flowmgt.GraphBuilder
}

func (b *graphBuilderBridge) GetGraph(ctx context.Context, flow *thunderidengine.FlowDefinition) (
	thunderidengine.FlowGraph, error) {
	complete, err := toCompleteFlowDefinition(flow)
	if err != nil {
		return nil, err
	}
	graph, svcErr := b.inner.GetGraph(ctx, complete)
	if svcErr != nil {
		return nil, asError(svcErr)
	}
	return &runtimeFlowGraph{graph: graph}, nil
}

func (b *graphBuilderBridge) InvalidateCache(ctx context.Context, flowID string) {
	b.inner.InvalidateCache(ctx, flowID)
}

// NewGraphBuilderFromRuntime creates a GraphBuilder using Thunder runtime flow dependencies.
func NewGraphBuilderFromRuntime(
	flowFactory flowcore.FlowFactoryInterface,
	executorRegistry executor.ExecutorRegistryInterface,
	graphCache flowcore.GraphCacheInterface,
) thunderidengine.GraphBuilder {
	return &graphBuilderBridge{
		inner: flowmgt.NewGraphBuilder(flowFactory, executorRegistry, graphCache),
	}
}

// NewGraphBuilder creates a GraphBuilder for an embedder-supplied executor registry.
func NewGraphBuilder(opts thunderidengine.GraphBuilderOptions) (thunderidengine.GraphBuilder, error) {
	if opts.Executors == nil {
		return nil, thunderidengine.ErrInvalidConfig
	}
	var graphCache flowcore.GraphCacheInterface
	if opts.Cache != nil {
		holder, ok := opts.Cache.(*graphCacheHolder)
		if !ok {
			return nil, thunderidengine.ErrInvalidConfig
		}
		graphCache = holder.inner()
	}
	if graphCache == nil {
		holder, ok := NewInMemoryGraphCache().(*graphCacheHolder)
		if !ok {
			return nil, thunderidengine.ErrInvalidConfig
		}
		graphCache = holder.inner()
	}
	return NewGraphBuilderFromRuntime(
		flowcore.NewFlowFactory(),
		newExecutorRegistryBridge(opts.Executors),
		graphCache,
	), nil
}

type flowMgtGraphSource interface {
	GetGraph(ctx context.Context, flowID string) (flowcore.GraphInterface, *serviceerror.ServiceError)
	GetFlowByHandle(ctx context.Context, handle string, flowType common.FlowType) (
		*flowmgt.CompleteFlowDefinition, *serviceerror.ServiceError)
}

type flowGraphProviderFromMgt struct {
	source flowMgtGraphSource
}

// NewFlowGraphProviderFromMgt returns a FlowGraphProvider backed by flow management.
func NewFlowGraphProviderFromMgt(source flowMgtGraphSource) thunderidengine.FlowGraphProvider {
	return &flowGraphProviderFromMgt{source: source}
}

func (p *flowGraphProviderFromMgt) GetGraph(ctx context.Context, flowID string) (thunderidengine.FlowGraph, error) {
	graph, svcErr := p.source.GetGraph(ctx, flowID)
	if svcErr != nil {
		return nil, asError(svcErr)
	}
	return &runtimeFlowGraph{graph: graph}, nil
}

func (p *flowGraphProviderFromMgt) GetFlowIDByHandle(ctx context.Context, handle, flowType string) (string, error) {
	flow, svcErr := p.source.GetFlowByHandle(ctx, handle, common.FlowType(flowType))
	if svcErr != nil {
		return "", asError(svcErr)
	}
	return flow.ID, nil
}

type serviceErrorCarrier struct {
	err *serviceerror.ServiceError
}

func (e serviceErrorCarrier) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.ErrorDescription.DefaultValue
}

// AsServiceError converts a service error into a standard error for FlowGraphProvider returns.
func AsServiceError(svcErr *serviceerror.ServiceError) error {
	return asError(svcErr)
}

func asError(svcErr *serviceerror.ServiceError) error {
	if svcErr == nil {
		return nil
	}
	return serviceErrorCarrier{err: svcErr}
}

// ServiceErrorFromErr unwraps a service error returned by this package.
func ServiceErrorFromErr(err error) (*serviceerror.ServiceError, bool) {
	var carrier serviceErrorCarrier
	if errors.As(err, &carrier) {
		return carrier.err, true
	}
	return nil, false
}

type flowGraphProvider struct {
	builder thunderidengine.GraphBuilder
	defs    thunderidengine.FlowDefinitionProvider
}

// NewFlowGraphProvider composes flow definition lookup and graph building.
func NewFlowGraphProvider(builder thunderidengine.GraphBuilder, defs thunderidengine.FlowDefinitionProvider) thunderidengine.FlowGraphProvider {
	return &flowGraphProvider{builder: builder, defs: defs}
}

func (p *flowGraphProvider) GetGraph(ctx context.Context, flowID string) (thunderidengine.FlowGraph, error) {
	def, err := p.defs.GetFlowByID(ctx, flowID)
	if err != nil {
		return nil, err
	}
	return p.builder.GetGraph(ctx, def)
}

func (p *flowGraphProvider) GetFlowIDByHandle(ctx context.Context, handle, flowType string) (string, error) {
	def, err := p.defs.GetFlowByHandleAndType(ctx, handle, flowType)
	if err != nil {
		return "", err
	}
	return def.ID, nil
}

func toCompleteFlowDefinition(flow *thunderidengine.FlowDefinition) (*flowmgt.CompleteFlowDefinition, error) {
	if flow == nil {
		return nil, errors.New("flow definition is nil")
	}
	nodes := make([]flowmgt.NodeDefinition, 0, len(flow.Nodes))
	for _, n := range flow.Nodes {
		nodes = append(nodes, flowmgt.NodeDefinition{
			ID:         n.ID,
			Type:       n.Type,
			Properties: n.Properties,
		})
	}
	name := flow.Handle
	if name == "" {
		name = flow.ID
	}
	return &flowmgt.CompleteFlowDefinition{
		ID:       flow.ID,
		Handle:   flow.Handle,
		Name:     name,
		FlowType: common.FlowType(flow.FlowType),
		Nodes:    nodes,
	}, nil
}
