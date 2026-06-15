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
	"github.com/thunder-id/thunderid/internal/cert"
	inboundmodel "github.com/thunder-id/thunderid/internal/inboundclient/model"
	oauth2const "github.com/thunder-id/thunderid/internal/oauth/oauth2/constants"
)

func toInboundModelClient(client *InboundClient) *inboundmodel.InboundClient {
	if client == nil {
		return nil
	}
	entityID := client.EntityID
	if entityID == "" {
		entityID = client.ApplicationID
	}
	appID := client.ApplicationID
	if appID == "" {
		appID = entityID
	}
	return &inboundmodel.InboundClient{
		ID:                        appID,
		AuthFlowID:                client.AuthFlowID,
		RegistrationFlowID:        client.RegistrationFlowID,
		IsRegistrationFlowEnabled: client.IsRegistrationFlowEnabled,
		RecoveryFlowID:            client.RecoveryFlowID,
		IsRecoveryFlowEnabled:     client.IsRecoveryFlowEnabled,
		Properties: map[string]interface{}{
			"clientId": client.ClientID,
			"name":     client.ApplicationID,
		},
	}
}

func toOAuthClient(client *InboundClient) *inboundmodel.OAuthClient {
	if client == nil {
		return nil
	}
	entityID := client.EntityID
	if entityID == "" {
		entityID = client.ApplicationID
	}
	return &inboundmodel.OAuthClient{
		ID:                                 entityID,
		OUID:                               client.OUID,
		ClientID:                           client.ClientID,
		RedirectURIs:                       append([]string(nil), client.RedirectURIs...),
		GrantTypes:                         toGrantTypes(client.GrantTypes),
		ResponseTypes:                      toResponseTypes(client.ResponseTypes),
		TokenEndpointAuthMethod:            oauth2const.TokenEndpointAuthMethod(client.TokenEndpointAuthMethod),
		PKCERequired:                       client.PKCERequired,
		PublicClient:                       client.PublicClient,
		RequirePushedAuthorizationRequests: client.RequirePushedAuthorizationRequests,
		Certificate:                        toInboundModelCertificate(client.Certificate),
	}
}

func toOAuthProfile(client *InboundClient) *inboundmodel.OAuthProfile {
	if client == nil {
		return nil
	}
	return &inboundmodel.OAuthProfile{
		RedirectURIs:                       append([]string(nil), client.RedirectURIs...),
		GrantTypes:                         append([]string(nil), client.GrantTypes...),
		ResponseTypes:                      append([]string(nil), client.ResponseTypes...),
		TokenEndpointAuthMethod:            client.TokenEndpointAuthMethod,
		PKCERequired:                       client.PKCERequired,
		PublicClient:                       client.PublicClient,
		RequirePushedAuthorizationRequests: client.RequirePushedAuthorizationRequests,
		Certificate:                        toInboundModelCertificate(client.Certificate),
	}
}

func toInboundModelCertificate(c *Certificate) *inboundmodel.Certificate {
	if c == nil {
		return nil
	}
	return &inboundmodel.Certificate{
		Type:  cert.CertificateType(c.Type),
		Value: c.Value,
	}
}

func toGrantTypes(values []string) []oauth2const.GrantType {
	if len(values) == 0 {
		return nil
	}
	out := make([]oauth2const.GrantType, len(values))
	for i, value := range values {
		out[i] = oauth2const.GrantType(value)
	}
	return out
}

func toResponseTypes(values []string) []oauth2const.ResponseType {
	if len(values) == 0 {
		return nil
	}
	out := make([]oauth2const.ResponseType, len(values))
	for i, value := range values {
		out[i] = oauth2const.ResponseType(value)
	}
	return out
}
