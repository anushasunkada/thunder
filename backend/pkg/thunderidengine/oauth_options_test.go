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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionsValidate_RequiresIssuerDeploymentIDBaseURL(t *testing.T) {
	err := Options{}.Validate()
	assert.Error(t, err)
	assert.ErrorContains(t, err, "issuer")

	err = Options{
		Issuer:  "https://id.example.com",
		BaseURL: "https://id.example.com",
	}.Validate()
	assert.ErrorContains(t, err, "deployment ID")

	err = Options{
		Issuer:       "https://id.example.com",
		DeploymentID: "prod",
	}.Validate()
	assert.ErrorContains(t, err, "base URL")

	err = Options{
		Issuer:       "https://id.example.com",
		DeploymentID: "prod",
		BaseURL:      "https://id.example.com",
	}.Validate()
	assert.NoError(t, err)
}

func TestOptionsOAuthPolicy(t *testing.T) {
	opts := Options{
		RequirePAR:               true,
		AllowWildcardRedirectURI: true,
	}
	policy := opts.OAuthPolicy()
	assert.True(t, policy.RequirePAR)
	assert.True(t, policy.AllowWildcardRedirectURI)
}
