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

import "sync"

// ExecutorRegistry holds executors registered by the host before engine construction.
type ExecutorRegistry interface {
	Register(name string, executor ExecutorInterface) error
	GetExecutor(name string) (ExecutorInterface, bool)
	IsRegistered(name string) bool
}

type executorRegistry struct {
	mu        sync.RWMutex
	executors map[string]ExecutorInterface
}

// NewExecutorRegistry creates an empty executor registry.
func NewExecutorRegistry() ExecutorRegistry {
	return &executorRegistry{executors: make(map[string]ExecutorInterface)}
}

func (r *executorRegistry) Register(name string, executor ExecutorInterface) error {
	if name == "" || executor == nil {
		return ErrInvalidConfig
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.executors[name]; exists {
		return ErrExecutorExists
	}
	r.executors[name] = executor
	return nil
}

func (r *executorRegistry) GetExecutor(name string) (ExecutorInterface, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ex, ok := r.executors[name]
	return ex, ok
}

func (r *executorRegistry) IsRegistered(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.executors[name]
	return ok
}
