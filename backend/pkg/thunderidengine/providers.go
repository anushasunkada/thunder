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

import "context"

// Host aggregates all provider implementations supplied by the embedding process.
type Host struct {
	ClientProvider         ClientProvider
	AuthnProvider          AuthnProvider
	AuthzProvider          AuthzProvider
	ResourceProvider       ResourceProvider
	OUProvider             OUProvider
	IDPProvider            IDPProvider
	FlowDefinitionProvider FlowDefinitionProvider
	ObservabilityProvider  ObservabilityProvider
	RuntimeStore           RuntimeStore
	ConsentProvider        ConsentProvider
	DesignProvider         DesignProvider
	I18n                   I18nProvider
	Crypto                 RuntimeCryptoProvider
}

// ClientProvider resolves OAuth clients, applications, and entity groups.
type ClientProvider interface {
	GetOAuthClientByClientID(ctx context.Context, clientID string) (*OAuthClient, error)
	GetTransitiveEntityGroups(ctx context.Context, entityID string) ([]EntityGroup, error)
	GetApplicationByID(ctx context.Context, appID string) (*Application, error)
	GetFlowApplicationByID(ctx context.Context, appID string) (*FlowApplication, error)
}

// AuthnProvider performs user authentication and attribute lookup.
type AuthnProvider interface {
	AuthenticateUser(ctx context.Context, credentials map[string]interface{}) (*AuthResult, error)
	GetUserAttributes(ctx context.Context, userID string, attrs []string) (*Attributes, error)
}

// AuthzProvider checks authorization for a subject/action/resource.
type AuthzProvider interface {
	IsAuthorized(ctx context.Context, subjectID, action, resourceID string) bool
}

// ResourceProvider resolves protected resources by URI.
type ResourceProvider interface {
	GetResource(ctx context.Context, resourceURI string) (*Resource, error)
}

// OUProvider resolves organizational units.
type OUProvider interface {
	GetOU(ctx context.Context, ouID string) (*OU, error)
	GetOUAncestors(ctx context.Context, ouID string) ([]OU, error)
}

// IDPProvider resolves identity provider configuration.
type IDPProvider interface {
	GetIDPByID(ctx context.Context, id string) (*IDP, error)
	GetIDPByName(ctx context.Context, name string) (*IDP, error)
}

// FlowDefinitionProvider supplies flow definitions for execution.
type FlowDefinitionProvider interface {
	GetFlowByID(ctx context.Context, id string) (*FlowDefinition, error)
	GetFlowByHandle(ctx context.Context, appID, handle string) (*FlowDefinition, error)
}

// ObservabilityProvider publishes runtime events.
type ObservabilityProvider interface {
	IsEnabled() bool
	PublishEvent(ctx context.Context, event ObservabilityEvent) error
}

// RuntimeStore persists OAuth and flow runtime state.
// PAR Store returns an opaque random key; the request_uri URN prefix is added by the service layer.
// PAR Consume takes that random key, not a full request_uri.
type RuntimeStore interface {
	Store(ctx context.Context, request PARRequest, expirySeconds int64) (string, error)
	Consume(ctx context.Context, randomKey string) (PARRequest, bool, error)
	AddRequest(ctx context.Context, value AuthRequestContext) (string, error)
	GetRequest(ctx context.Context, key string) (bool, AuthRequestContext, error)
	ClearRequest(ctx context.Context, key string) error
	InsertAuthorizationCode(ctx context.Context, authzCode AuthorizationCode) error
	ConsumeAuthorizationCode(ctx context.Context, authCode string) (bool, error)
	GetAuthorizationCode(ctx context.Context, authCode string) (*AuthorizationCode, error)
	StoreFlowContext(ctx context.Context, dbModel FlowContextDB, expirySeconds int64) error
	GetFlowContext(ctx context.Context, executionID string) (*FlowContextDB, error)
	UpdateFlowContext(ctx context.Context, dbModel FlowContextDB) error
	DeleteFlowContext(ctx context.Context, executionID string) error
}

// ConsentProvider resolves and records end-user consent.
type ConsentProvider interface {
	ResolveConsent(ctx context.Context, req ConsentResolveRequest) (*ConsentPromptData, error)
	RecordConsent(ctx context.Context, req ConsentRecordRequest) (*Consent, error)
}

// DesignProvider resolves theme and layout for an application or OU.
type DesignProvider interface {
	ResolveDesign(ctx context.Context, resolveType DesignResolveType, id string) (*DesignResponse, error)
}

// I18nProvider resolves translations and lists supported languages.
type I18nProvider interface {
	ResolveTranslations(ctx context.Context, language, namespace string) (map[string]map[string]string, error)
	ListLanguages(ctx context.Context) ([]string, error)
}

// RuntimeCryptoProvider performs runtime signing and encryption via the host KMS.
type RuntimeCryptoProvider interface {
	Encrypt(ctx context.Context, keyRef *KeyRef, params AlgorithmParams, content []byte) ([]byte, error)
	Decrypt(ctx context.Context, keyRef *KeyRef, params AlgorithmParams, content []byte) ([]byte, error)
	Sign(ctx context.Context, keyRef KeyRef, algorithm SignAlgorithm, content []byte) ([]byte, error)
	GetPublicKeys(ctx context.Context, filter PublicKeyFilter) ([]PublicKeyInfo, error)
	GetTLSMaterial(ctx context.Context, keyRef *KeyRef) (*TLSMaterial, error)
}
