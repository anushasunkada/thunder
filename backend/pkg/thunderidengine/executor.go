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
	"github.com/thunder-id/thunderid/internal/flow/core"
	"github.com/thunder-id/thunderid/internal/flow/executor"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/flow"
)

// ExecutorRegistry registers custom flow executors alongside built-in executors.
type ExecutorRegistry struct {
	inner executor.ExecutorRegistryInterface
}

// Register adds a custom executor to the registry.
func (r ExecutorRegistry) Register(name string, exec flow.Executor) {
	r.inner.RegisterExecutor(name, exec)
}

// IsRegistered reports whether an executor with the given name is registered.
func (r ExecutorRegistry) IsRegistered(name string) bool {
	return r.inner.IsRegistered(name)
}

// FlowFactory creates flow graph components for custom executor implementations.
type FlowFactory struct {
	inner core.FlowFactoryInterface
}

// CreateExecutor creates a base executor with the given name, type, inputs, and prerequisites.
// Host executors typically embed the returned value and override Execute.
func (f FlowFactory) CreateExecutor(name string, executorType flow.ExecutorType,
	defaultInputs, prerequisites []flow.Input) flow.Executor {
	return f.inner.CreateExecutor(name, executorType, defaultInputs, prerequisites)
}

// NewExecutorRegistry wraps an internal executor registry for host registration.
func NewExecutorRegistry(inner executor.ExecutorRegistryInterface) ExecutorRegistry {
	return ExecutorRegistry{inner: inner}
}

// NewFlowFactory wraps an internal flow factory for custom executor construction.
func NewFlowFactory(inner core.FlowFactoryInterface) FlowFactory {
	return FlowFactory{inner: inner}
}

// NewDefaultFlowFactory returns the default flow factory for custom executor construction.
func NewDefaultFlowFactory() FlowFactory {
	return NewFlowFactory(core.Initialize())
}

// NewEmptyExecutorRegistry returns an empty executor registry for host registration.
func NewEmptyExecutorRegistry() ExecutorRegistry {
	return NewExecutorRegistry(executor.NewRegistry())
}

func wrapExecutorRegistry(reg executor.ExecutorRegistryInterface) ExecutorRegistry {
	return ExecutorRegistry{inner: reg}
}

func wrapFlowFactory(factory core.FlowFactoryInterface) FlowFactory {
	return FlowFactory{inner: factory}
}
