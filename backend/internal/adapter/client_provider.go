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

// Package adapter provides host implementations of thunderidengine provider interfaces.
package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/thunder-id/thunderid/internal/entityprovider"
	"github.com/thunder-id/thunderid/internal/inboundclient"
	inboundmodel "github.com/thunder-id/thunderid/internal/inboundclient/model"
	"github.com/thunder-id/thunderid/pkg/thunderidengine"
)

const (
	sysAttrName        = "name"
	sysAttrDescription = "description"
	sysAttrClientID    = "clientId"
	propURL            = "url"
	propLogoURL        = "logo_url"
	propTosURI         = "tos_uri"
	propPolicyURI      = "policy_uri"
	propMetadata       = "metadata"
)

type clientProvider struct {
	entityProvider entityprovider.EntityProviderInterface
	inboundClient  inboundclient.InboundClientServiceInterface
}

// NewClientProvider returns a thunderidengine.ClientProvider backed by the host entity and inbound client services.
func NewClientProvider(
	entityProvider entityprovider.EntityProviderInterface,
	inboundClient inboundclient.InboundClientServiceInterface,
) thunderidengine.ClientProvider {
	return &clientProvider{
		entityProvider: entityProvider,
		inboundClient:  inboundClient,
	}
}

func (p *clientProvider) GetOAuthClientByClientID(
	ctx context.Context, clientID string,
) (*thunderidengine.OAuthClient, error) {
	oauth, err := p.inboundClient.GetOAuthClientByClientID(ctx, clientID)
	if err != nil {
		if errors.Is(err, inboundclient.ErrInboundClientNotFound) {
			return nil, thunderidengine.ErrInboundClientNotFound
		}
		return nil, err
	}
	if oauth == nil {
		return nil, thunderidengine.ErrInboundClientNotFound
	}

	inbound, err := p.inboundClient.GetInboundClientByEntityID(ctx, oauth.ID)
	if err != nil && !errors.Is(err, inboundclient.ErrInboundClientNotFound) {
		return nil, err
	}

	return toEngineOAuthClient(oauth, inbound), nil
}

func (p *clientProvider) GetTransitiveEntityGroups(
	_ context.Context, entityID string,
) ([]thunderidengine.EntityGroup, error) {
	groups, epErr := p.entityProvider.GetTransitiveEntityGroups(entityID)
	if epErr != nil {
		return nil, fmt.Errorf("get transitive entity groups: %w", epErr)
	}

	result := make([]thunderidengine.EntityGroup, len(groups))
	for i := range groups {
		result[i] = thunderidengine.EntityGroup{
			ID:   groups[i].ID,
			Name: groups[i].Name,
		}
	}
	return result, nil
}

func (p *clientProvider) GetApplicationByID(
	ctx context.Context, appID string,
) (*thunderidengine.Application, error) {
	entity, epErr := p.entityProvider.GetEntity(appID)
	if epErr != nil && epErr.Code != entityprovider.ErrorCodeEntityNotFound {
		return nil, fmt.Errorf("get entity: %w", epErr)
	}

	inbound, err := p.inboundClient.GetInboundClientByEntityID(ctx, appID)
	if err != nil {
		if errors.Is(err, inboundclient.ErrInboundClientNotFound) {
			return nil, thunderidengine.ErrApplicationNotFound
		}
		return nil, err
	}
	if inbound == nil {
		return nil, thunderidengine.ErrApplicationNotFound
	}

	return toEngineApplication(appID, entity, inbound), nil
}

func (p *clientProvider) GetFlowApplicationByID(
	ctx context.Context, appID string,
) (*thunderidengine.FlowApplication, error) {
	inbound, err := p.inboundClient.GetInboundClientByEntityID(ctx, appID)
	if err != nil {
		if errors.Is(err, inboundclient.ErrInboundClientNotFound) {
			return nil, thunderidengine.ErrApplicationNotFound
		}
		return nil, err
	}
	if inbound == nil {
		return nil, thunderidengine.ErrApplicationNotFound
	}

	entity, epErr := p.entityProvider.GetEntity(appID)
	if epErr != nil && epErr.Code != entityprovider.ErrorCodeEntityNotFound {
		return nil, fmt.Errorf("get entity: %w", epErr)
	}

	return toEngineFlowApplication(appID, entity, inbound), nil
}

func toEngineOAuthClient(
	oauth *inboundmodel.OAuthClient, inbound *inboundmodel.InboundClient,
) *thunderidengine.OAuthClient {
	client := &thunderidengine.OAuthClient{
		ClientID:                           oauth.ClientID,
		EntityID:                           oauth.ID,
		OUID:                               oauth.OUID,
		RedirectURIs:                       oauth.RedirectURIs,
		Scopes:                             oauth.Scopes,
		TokenEndpointAuthMethod:            oauth.TokenEndpointAuthMethod,
		PKCERequired:                       oauth.PKCERequired,
		PublicClient:                       oauth.PublicClient,
		RequirePushedAuthorizationRequests: oauth.RequirePushedAuthorizationRequests,
		UserInfo:                           toEngineUserInfoConfig(oauth.UserInfo),
		ScopeClaims:                        copyScopeClaims(oauth.ScopeClaims),
		AcrValues:                          append([]string(nil), oauth.AcrValues...),
		Certificate:                        toEngineCertificate(oauth.Certificate),
		Properties:                         copyProperties(inboundProperties(inbound)),
	}
	client.GrantTypes = append(client.GrantTypes, oauth.GrantTypes...)
	client.ResponseTypes = append(client.ResponseTypes, oauth.ResponseTypes...)
	if oauth.Token != nil {
		client.Token = toEngineOAuthTokenConfig(oauth.Token)
		if oauth.Token.AccessToken != nil {
			client.AccessTokenValiditySeconds = oauth.Token.AccessToken.ValidityPeriod
		}
		if oauth.Token.IDToken != nil {
			client.IDTokenValiditySeconds = oauth.Token.IDToken.ValidityPeriod
		}
	}
	if inbound != nil {
		client.AuthFlowID = inbound.AuthFlowID
		client.RegistrationFlowID = inbound.RegistrationFlowID
		client.RecoveryFlowID = inbound.RecoveryFlowID
		client.IsRegistrationFlowEnabled = inbound.IsRegistrationFlowEnabled
		client.IsRecoveryFlowEnabled = inbound.IsRecoveryFlowEnabled
		client.Assertion = toEngineAssertionConfig(inbound.Assertion)
		client.LoginConsent = toEngineLoginConsentConfig(inbound.LoginConsent)
		if len(inbound.AllowedUserTypes) > 0 {
			client.AllowedUserTypes = append([]string(nil), inbound.AllowedUserTypes...)
		}
		if len(client.Properties) == 0 && inbound.Properties != nil {
			client.Properties = copyProperties(inbound.Properties)
		}
	}
	return client
}

func toEngineFlowApplication(
	appID string, entity *entityprovider.Entity, inbound *inboundmodel.InboundClient,
) *thunderidengine.FlowApplication {
	app := &thunderidengine.FlowApplication{
		ID:                        appID,
		Assertion:                 toEngineAssertionConfig(inbound.Assertion),
		LoginConsent:              toEngineLoginConsentConfig(inbound.LoginConsent),
		AuthFlowID:                inbound.AuthFlowID,
		RegistrationFlowID:        inbound.RegistrationFlowID,
		RecoveryFlowID:            inbound.RecoveryFlowID,
		IsRegistrationFlowEnabled: inbound.IsRegistrationFlowEnabled,
		IsRecoveryFlowEnabled:     inbound.IsRecoveryFlowEnabled,
	}
	if len(inbound.AllowedUserTypes) > 0 {
		app.AllowedUserTypes = append([]string(nil), inbound.AllowedUserTypes...)
	}
	if inbound.Properties != nil {
		if metadata, ok := inbound.Properties[propMetadata].(map[string]interface{}); ok {
			app.Metadata = copyProperties(metadata)
		}
	}
	entityAttrs := readEntitySystemAttributes(entity)
	if name, ok := entityAttrs[sysAttrName].(string); ok {
		app.Name = name
	}
	if clientID, ok := entityAttrs[sysAttrClientID].(string); ok {
		app.OAuthClientID = clientID
	}
	return app
}

func toEngineOAuthTokenConfig(src *inboundmodel.OAuthTokenConfig) *thunderidengine.OAuthTokenConfig {
	if src == nil {
		return nil
	}
	dst := &thunderidengine.OAuthTokenConfig{}
	if src.AccessToken != nil {
		dst.AccessToken = &thunderidengine.AccessTokenConfig{
			ValidityPeriod: src.AccessToken.ValidityPeriod,
			UserAttributes: append([]string(nil), src.AccessToken.UserAttributes...),
		}
	}
	if src.IDToken != nil {
		dst.IDToken = &thunderidengine.IDTokenConfig{
			ValidityPeriod: src.IDToken.ValidityPeriod,
			UserAttributes: append([]string(nil), src.IDToken.UserAttributes...),
			ResponseType:   src.IDToken.ResponseType,
			EncryptionAlg:  src.IDToken.EncryptionAlg,
			EncryptionEnc:  src.IDToken.EncryptionEnc,
		}
	}
	return dst
}

func toEngineUserInfoConfig(src *inboundmodel.UserInfoConfig) *thunderidengine.UserInfoConfig {
	if src == nil {
		return nil
	}
	return &thunderidengine.UserInfoConfig{
		ResponseType:   src.ResponseType,
		UserAttributes: append([]string(nil), src.UserAttributes...),
		SigningAlg:     src.SigningAlg,
		EncryptionAlg:  src.EncryptionAlg,
		EncryptionEnc:  src.EncryptionEnc,
	}
}

func toEngineAssertionConfig(src *inboundmodel.AssertionConfig) *thunderidengine.AssertionConfig {
	if src == nil {
		return nil
	}
	return &thunderidengine.AssertionConfig{
		ValidityPeriod: src.ValidityPeriod,
		UserAttributes: append([]string(nil), src.UserAttributes...),
	}
}

func toEngineLoginConsentConfig(src *inboundmodel.LoginConsentConfig) *thunderidengine.LoginConsentConfig {
	if src == nil {
		return nil
	}
	return &thunderidengine.LoginConsentConfig{
		ValidityPeriod: src.ValidityPeriod,
	}
}

func toEngineCertificate(src *inboundmodel.Certificate) *thunderidengine.Certificate {
	if src == nil {
		return nil
	}
	return &thunderidengine.Certificate{
		Type:  src.Type,
		Value: src.Value,
	}
}

func toEngineApplication(
	appID string, entity *entityprovider.Entity, inbound *inboundmodel.InboundClient,
) *thunderidengine.Application {
	app := &thunderidengine.Application{
		ID:                        appID,
		IsRegistrationFlowEnabled: inbound.IsRegistrationFlowEnabled,
		IsRecoveryFlowEnabled:     inbound.IsRecoveryFlowEnabled,
		Properties:                copyProperties(inbound.Properties),
		AuthFlowID:                inbound.AuthFlowID,
		RegistrationFlowID:        inbound.RegistrationFlowID,
		RecoveryFlowID:            inbound.RecoveryFlowID,
	}
	if entity != nil {
		app.OUID = entity.OUID
		entityAttrs := readEntitySystemAttributes(entity)
		if name, ok := entityAttrs[sysAttrName].(string); ok {
			app.Name = name
		}
		if desc, ok := entityAttrs[sysAttrDescription].(string); ok {
			app.Description = desc
		}
	}
	if inbound.Properties != nil {
		if url, ok := inbound.Properties[propURL].(string); ok {
			app.URL = url
		}
		if logoURL, ok := inbound.Properties[propLogoURL].(string); ok {
			app.LogoURL = logoURL
		}
		if tosURI, ok := inbound.Properties[propTosURI].(string); ok {
			app.TosURI = tosURI
		}
		if policyURI, ok := inbound.Properties[propPolicyURI].(string); ok {
			app.PolicyURI = policyURI
		}
	}
	return app
}

func readEntitySystemAttributes(entity *entityprovider.Entity) map[string]interface{} {
	if entity == nil || len(entity.SystemAttributes) == 0 {
		return map[string]interface{}{}
	}
	var attrs map[string]interface{}
	if err := json.Unmarshal(entity.SystemAttributes, &attrs); err != nil || attrs == nil {
		return map[string]interface{}{}
	}
	return attrs
}

func inboundProperties(inbound *inboundmodel.InboundClient) map[string]interface{} {
	if inbound == nil {
		return nil
	}
	return inbound.Properties
}

func copyScopeClaims(src map[string][]string) map[string][]string {
	if src == nil {
		return nil
	}
	dst := make(map[string][]string, len(src))
	for k, v := range src {
		dst[k] = append([]string(nil), v...)
	}
	return dst
}

func copyProperties(src map[string]interface{}) map[string]interface{} {
	if src == nil {
		return nil
	}
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
