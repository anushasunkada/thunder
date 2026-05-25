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

package thunderidengine_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thunder-id/thunderid/pkg/thunderidengine"
)

type stubExecutor struct{}

func (stubExecutor) Execute(_ *thunderidengine.NodeContext) (*thunderidengine.ExecutorResponse, error) {
	return nil, nil
}
func (stubExecutor) GetName() string { return "stub" }
func (stubExecutor) GetType() thunderidengine.ExecutorType {
	return thunderidengine.ExecutorTypeAuthentication
}
func (stubExecutor) GetDefaultInputs() []thunderidengine.Input { return nil }
func (stubExecutor) GetPrerequisites() []thunderidengine.Input { return nil }
func (stubExecutor) HasRequiredInputs(_ *thunderidengine.NodeContext, _ *thunderidengine.ExecutorResponse) bool {
	return true
}
func (stubExecutor) ValidatePrerequisites(_ *thunderidengine.NodeContext, _ *thunderidengine.ExecutorResponse) bool {
	return true
}
func (stubExecutor) GetUserIDFromContext(_ *thunderidengine.NodeContext) string { return "" }
func (stubExecutor) GetRequiredInputs(_ *thunderidengine.NodeContext) []thunderidengine.Input {
	return nil
}
func (stubExecutor) GetExecutionPolicy(_ string) *thunderidengine.ExecutionPolicy { return nil }

func TestExecutorRegistry_RegisterAndGet(t *testing.T) {
	reg := thunderidengine.NewExecutorRegistry()
	err := reg.Register("StubExecutor", stubExecutor{})
	require.NoError(t, err)
	ex, ok := reg.GetExecutor("StubExecutor")
	require.True(t, ok)
	require.Equal(t, "stub", ex.GetName())
}

func TestExecutorRegistry_DuplicateRegister(t *testing.T) {
	reg := thunderidengine.NewExecutorRegistry()
	require.NoError(t, reg.Register("StubExecutor", stubExecutor{}))
	err := reg.Register("StubExecutor", stubExecutor{})
	require.ErrorIs(t, err, thunderidengine.ErrExecutorExists)
}
