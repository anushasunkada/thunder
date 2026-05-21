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
	"testing"

	"github.com/stretchr/testify/require"

	authnprovidermgr "github.com/thunder-id/thunderid/internal/authnprovider/manager"
	"github.com/thunder-id/thunderid/pkg/thunderidengine"
)

type stubAuthnProvider struct {
	result *thunderidengine.AuthnResult
	err    error
	attrs  map[string]interface{}
}

func (s *stubAuthnProvider) AuthenticateUser(
	_ context.Context, _ thunderidengine.Credentials,
) (*thunderidengine.AuthnResult, error) {
	return s.result, s.err
}

func (s *stubAuthnProvider) GetUserAttributes(
	_ context.Context, _ string, _ []string,
) (map[string]interface{}, error) {
	return s.attrs, nil
}

func TestAuthnBridgeAuthenticateUser(t *testing.T) {
	bridge := newAuthnBridge(&stubAuthnProvider{result: &thunderidengine.AuthnResult{
		UserID: "user-1", OUID: "ou-1", EntityType: "user",
	}})
	user, result, svcErr := bridge.AuthenticateUser(
		context.Background(), nil, nil, nil, nil, authnprovidermgr.AuthUser{},
	)
	require.Nil(t, svcErr)
	require.NotNil(t, result)
	require.Equal(t, "user-1", result.UserID)
	require.Equal(t, "user-1", user.UserID())
}

func TestAuthnBridgeGetUserAttributes(t *testing.T) {
	bridge := newAuthnBridge(&stubAuthnProvider{attrs: map[string]interface{}{"email": "a@b.c"}})
	user := authnprovidermgr.AuthUser{}
	user.ApplyIdentity("user-1", "user", "ou-1")
	_, attrs, svcErr := bridge.GetUserAttributes(context.Background(), nil, nil, user)
	require.Nil(t, svcErr)
	require.Equal(t, "a@b.c", attrs.Attributes["email"].Value)
}
