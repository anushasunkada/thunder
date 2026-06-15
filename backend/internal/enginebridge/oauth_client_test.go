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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thunder-id/thunderid/internal/cert"
	oauth2const "github.com/thunder-id/thunderid/internal/oauth/oauth2/constants"
)

func TestToOAuthClientMapsCertificate(t *testing.T) {
	jwks := `{"keys":[{"kty":"RSA","kid":"key-1"}]}`
	client := &InboundClient{
		ClientID:                "client-1",
		EntityID:                "entity-1",
		TokenEndpointAuthMethod: string(oauth2const.TokenEndpointAuthMethodPrivateKeyJWT),
		Certificate: &Certificate{
			Type:  string(cert.CertificateTypeJWKS),
			Value: jwks,
		},
	}

	oauthClient := toOAuthClient(client)
	require.NotNil(t, oauthClient)
	require.NotNil(t, oauthClient.Certificate)
	require.Equal(t, cert.CertificateTypeJWKS, oauthClient.Certificate.Type)
	require.Equal(t, jwks, oauthClient.Certificate.Value)
}

func TestToOAuthProfileMapsCertificate(t *testing.T) {
	jwksURI := "https://example.com/jwks"
	client := &InboundClient{
		TokenEndpointAuthMethod: string(oauth2const.TokenEndpointAuthMethodPrivateKeyJWT),
		Certificate: &Certificate{
			Type:  string(cert.CertificateTypeJWKSURI),
			Value: jwksURI,
		},
	}

	profile := toOAuthProfile(client)
	require.NotNil(t, profile)
	require.NotNil(t, profile.Certificate)
	require.Equal(t, cert.CertificateTypeJWKSURI, profile.Certificate.Type)
	require.Equal(t, jwksURI, profile.Certificate.Value)
}

func TestToOAuthClientNilCertificate(t *testing.T) {
	oauthClient := toOAuthClient(&InboundClient{ClientID: "client-1"})
	require.NotNil(t, oauthClient)
	require.Nil(t, oauthClient.Certificate)
}
