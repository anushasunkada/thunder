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

package host

import (
	"testing"

	"github.com/stretchr/testify/suite"

	authnprovidercm "github.com/thunder-id/thunderid/internal/authnprovider/common"
)

type AuthnInternalTestSuite struct {
	suite.Suite
}

func TestAuthnInternalSuite(t *testing.T) {
	suite.Run(t, new(AuthnInternalTestSuite))
}

func (suite *AuthnInternalTestSuite) TestToPublicAuthnMetadata_PassesFullAppMetadata() {
	internal := &authnprovidercm.AuthnMetadata{
		AppMetadata: map[string]interface{}{
			"applicationId":   "app-1",
			"ouId":            "ou-1",
			"tenant_id":       "tenant-1",
			"client_ids":      []string{"client-a", "client-b"},
			"oauth_client_id": "client-a",
		},
	}

	meta := toPublicAuthnMetadata(internal)

	suite.Equal("app-1", meta.ApplicationID)
	suite.Equal("ou-1", meta.OUID)
	suite.Equal(internal.AppMetadata, meta.AppMetadata)
	suite.Equal("tenant-1", meta.AppMetadata["tenant_id"])
	suite.Equal("client-a", meta.AppMetadata["oauth_client_id"])
}

func (suite *AuthnInternalTestSuite) TestToPublicAuthnMetadata_Nil() {
	suite.Nil(toPublicAuthnMetadata(nil))
}

func (suite *AuthnInternalTestSuite) TestToPublicGetAttributesMetadata_PassesFullAppMetadataAndLocale() {
	internal := &authnprovidercm.GetAttributesMetadata{
		AppMetadata: map[string]interface{}{
			"applicationId":   "app-1",
			"ouId":            "ou-1",
			"oauth_client_id": "flow-client",
			"client_ids":      []string{"flow-client"},
		},
		Locale: "en-US",
	}

	meta := toPublicGetAttributesMetadata(internal)

	suite.Equal("app-1", meta.ApplicationID)
	suite.Equal("ou-1", meta.OUID)
	suite.Equal("en-US", meta.Locale)
	suite.Equal(internal.AppMetadata, meta.AppMetadata)
	suite.Equal("flow-client", meta.AppMetadata["oauth_client_id"])
}

func (suite *AuthnInternalTestSuite) TestToPublicGetAttributesMetadata_Nil() {
	suite.Nil(toPublicGetAttributesMetadata(nil))
}
