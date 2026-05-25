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

import "fmt"

// GateClientOptions holds Gate UI redirect settings for the authorization endpoint.
type GateClientOptions struct {
	Scheme    string
	Hostname  string
	Port      int
	LoginPath string
	ErrorPath string
}

// FlowOptions holds flow defaults used when resolving inbound OAuth client profiles.
type FlowOptions struct {
	DefaultAuthFlowHandle string
	AutoInferRegistration bool
}

// Options holds OAuth runtime configuration passed explicitly at initialization.
type Options struct {
	Issuer                    string
	Audience                  string
	ValidityPeriod            int64
	Leeway                    int64
	DeploymentID              string
	BaseURL                   string
	RequirePAR                bool
	PARExpiresIn              int64
	AllowWildcardRedirectURI  bool
	AuthorizationCodeValidity int64
	RefreshTokenValidity      int64
	RefreshTokenRenewOnGrant  bool
	AcrAMR                    map[string][]string
	RuntimeStoreType          string
	GateClient                GateClientOptions
	DCRInsecure               bool
	Flow                      FlowOptions
}

// OAuthPolicy returns server-level OAuth policy flags derived from these options.
func (o Options) OAuthPolicy() OAuthPolicy {
	return OAuthPolicy{
		RequirePAR:               o.RequirePAR,
		AllowWildcardRedirectURI: o.AllowWildcardRedirectURI,
	}
}

// Validate checks required fields for OAuth initialization.
func (o Options) Validate() error {
	if o.Issuer == "" {
		return fmt.Errorf("oauth options: issuer is required")
	}
	if o.DeploymentID == "" {
		return fmt.Errorf("oauth options: deployment ID is required")
	}
	if o.BaseURL == "" {
		return fmt.Errorf("oauth options: base URL is required")
	}
	return nil
}
