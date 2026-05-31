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

package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thunder-id/thunderid/internal/flow/common"
	"github.com/thunder-id/thunderid/internal/flow/core"
)

func TestInitializeForEngine_CustomExecutorAfterBuiltInSubset(t *testing.T) {
	flowFactory := core.Initialize()
	reg, err := InitializeForEngine(RegisterDeps{FlowFactory: flowFactory},
		[]string{ExecutorNameBasicAuth, ExecutorNameAuthAssert})
	require.NoError(t, err)

	custom := createMockExecutorForRegistry(t, "MosipOtpExecutor", common.ExecutorTypeAuthentication)
	reg.RegisterExecutor("MosipOtpExecutor", custom)

	assert.True(t, reg.IsRegistered("MosipOtpExecutor"))
	assert.True(t, reg.IsRegistered(ExecutorNameBasicAuth))
	assert.True(t, reg.IsRegistered(ExecutorNameAuthAssert))
	assert.False(t, reg.IsRegistered(ExecutorNameConsent))
}

func TestResolveEngineExecutorNames_EmptyUsesDefaults(t *testing.T) {
	names, err := ResolveEngineExecutorNames(nil)
	require.NoError(t, err)
	assert.Equal(t, EngineDefaultExecutorNames, names)
}
