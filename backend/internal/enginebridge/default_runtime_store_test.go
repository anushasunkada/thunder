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

package enginebridge

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/thunder-id/thunderid/internal/flow/flowexec"
	oauthauthz "github.com/thunder-id/thunderid/internal/oauth/oauth2/authz"
	oauth2model "github.com/thunder-id/thunderid/internal/oauth/oauth2/model"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/par"
	"github.com/thunder-id/thunderid/pkg/thunderidengine"
)

func newDefaultRuntimeStoreForTest(host thunderidengine.RuntimeStore) thunderidengine.RuntimeStore {
	return &defaultRuntimeStore{stores: RuntimeStores{
		PAR:         par.NewStoreFromRuntime(host),
		AuthCode:    oauthauthz.NewCodeStoreFromRuntime(host),
		AuthRequest: oauthauthz.NewRequestStoreFromRuntime(host),
		FlowContext: flowexec.NewContextStoreFromRuntime(host),
	}}
}

func TestDefaultRuntimeStorePARRoundTrip(t *testing.T) {
	store := newDefaultRuntimeStoreForTest(newTestMemoryRuntimeStore())
	params := oauth2model.OAuthParameters{ClientID: "client-1", State: "st"}
	raw, err := json.Marshal(params)
	require.NoError(t, err)

	requestURI, err := store.Store(context.Background(), thunderidengine.PushedAuthorizationRequest{
		ClientID:        "client-1",
		OAuthParameters: thunderidengine.OAuthParameters(raw),
	}, 600)
	require.NoError(t, err)
	require.Contains(t, requestURI, parRequestURIPrefix)

	got, found, err := store.Consume(context.Background(), requestURI)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, "client-1", got.ClientID)
}

func TestDefaultRuntimeStoreAuthorizationCodeRoundTrip(t *testing.T) {
	store := newDefaultRuntimeStoreForTest(newTestMemoryRuntimeStore())
	now := time.Now().UTC()
	code := thunderidengine.AuthorizationCode{
		CodeID: "id-1", Code: "code-1", ClientID: "client-1",
		TimeCreated: now, ExpiryTime: now.Add(time.Minute),
	}
	require.NoError(t, store.InsertAuthorizationCode(context.Background(), code))

	got, err := store.GetAuthorizationCode(context.Background(), "code-1")
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "code-1", got.Code)

	consumed, err := store.ConsumeAuthorizationCode(context.Background(), "code-1")
	require.NoError(t, err)
	require.True(t, consumed)
}

func TestDefaultRuntimeStoreFlowContextRoundTrip(t *testing.T) {
	store := newDefaultRuntimeStoreForTest(newTestMemoryRuntimeStore())
	flow := thunderidengine.FlowContext{ExecutionID: "exec-1", Context: `{"graphId":"g1"}`}
	require.NoError(t, store.StoreFlowContext(context.Background(), flow, 120))

	got, err := store.GetFlowContext(context.Background(), "exec-1")
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, `{"graphId":"g1"}`, got.Context)

	flow.Context = `{"graphId":"g2"}`
	require.NoError(t, store.UpdateFlowContext(context.Background(), flow))
	got, err = store.GetFlowContext(context.Background(), "exec-1")
	require.NoError(t, err)
	require.Equal(t, `{"graphId":"g2"}`, got.Context)

	require.NoError(t, store.DeleteFlowContext(context.Background(), "exec-1"))
	got, err = store.GetFlowContext(context.Background(), "exec-1")
	require.NoError(t, err)
	require.Nil(t, got)
}
