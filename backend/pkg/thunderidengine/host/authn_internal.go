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
	"context"

	authnprovidercm "github.com/thunder-id/thunderid/internal/authnprovider/common"
	authnprovidermgr "github.com/thunder-id/thunderid/internal/authnprovider/manager"
	"github.com/thunder-id/thunderid/internal/authnprovider/provider"
	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
	"github.com/thunder-id/thunderid/internal/system/i18n/core"
)

type internalAuthnProvider struct {
	host AuthnProvider
}

func (p *internalAuthnProvider) Authenticate(ctx context.Context, identifiers, credentials map[string]interface{},
	metadata *authnprovidercm.AuthnMetadata) (*authnprovidercm.AuthnResult, *serviceerror.ServiceError) {
	result, err := p.host.Authenticate(ctx, identifiers, credentials, toPublicAuthnMetadata(metadata))
	if err != nil {
		return nil, serviceError(err)
	}
	return toInternalAuthnResult(result), nil
}

func (p *internalAuthnProvider) GetAttributes(ctx context.Context, token string,
	requested *authnprovidercm.RequestedAttributes,
	metadata *authnprovidercm.GetAttributesMetadata,
) (*authnprovidercm.GetAttributesResult, *serviceerror.ServiceError) {
	result, err := p.host.GetAttributes(ctx, token, toPublicRequestedAttributes(requested),
		toPublicGetAttributesMetadata(metadata))
	if err != nil {
		return nil, serviceError(err)
	}
	return toInternalGetAttributesResult(result), nil
}

func serviceError(err error) *serviceerror.ServiceError {
	return serviceerror.CustomServiceError(serviceerror.InternalServerError, core.I18nMessage{
		Key:          "error.internal",
		DefaultValue: err.Error(),
	})
}

func toPublicAuthnMetadata(m *authnprovidercm.AuthnMetadata) *AuthnMetadata {
	if m == nil {
		return nil
	}
	meta := &AuthnMetadata{AppMetadata: m.AppMetadata, RuntimeMetadata: m.RuntimeMetadata}
	if m.AppMetadata != nil {
		if appID, ok := m.AppMetadata["applicationId"].(string); ok {
			meta.ApplicationID = appID
		}
		if ouID, ok := m.AppMetadata["ouId"].(string); ok {
			meta.OUID = ouID
		}
	}
	meta.RuntimeMetadata = m.RuntimeMetadata
	return meta
}

func toPublicGetAttributesMetadata(m *authnprovidercm.GetAttributesMetadata) *GetAttributesMetadata {
	if m == nil {
		return nil
	}
	meta := &GetAttributesMetadata{
		AppMetadata: m.AppMetadata,
		Locale:      m.Locale,
	}
	if m.AppMetadata != nil {
		if appID, ok := m.AppMetadata["applicationId"].(string); ok {
			meta.ApplicationID = appID
		}
		if ouID, ok := m.AppMetadata["ouId"].(string); ok {
			meta.OUID = ouID
		}
	}
	meta.RuntimeMetadata = m.RuntimeMetadata
	return meta
}

func toPublicRequestedAttributes(r *authnprovidercm.RequestedAttributes) *RequestedAttributes {
	if r == nil {
		return nil
	}
	names := make([]string, 0, len(r.Attributes))
	for name := range r.Attributes {
		names = append(names, name)
	}
	return &RequestedAttributes{AttributeNames: names}
}

func toInternalAuthnResult(r *AuthnResult) *authnprovidercm.AuthnResult {
	if r == nil {
		return nil
	}
	return &authnprovidercm.AuthnResult{
		EntityID:       r.UserID,
		UserID:         r.UserID,
		Token:          r.AuthToken,
		IsExistingUser: r.Authenticated,
	}
}

func toInternalGetAttributesResult(r *GetAttributesResult) *authnprovidercm.GetAttributesResult {
	if r == nil {
		return nil
	}
	return &authnprovidercm.GetAttributesResult{}
}

var _ provider.AuthnProviderInterface = (*internalAuthnProvider)(nil)

type authnProviderManager struct {
	mgr authnprovidermgr.AuthnProviderManagerInterface
}

func (a *authnProviderManager) isAuthnProviderManager() {}

// InitializeAuthnProviderManager wraps a host AuthnProvider with manager-layer semantics.
func InitializeAuthnProviderManager(provider AuthnProvider) AuthnProviderManager {
	return &authnProviderManager{
		mgr: authnprovidermgr.InitializeFromProvider(&internalAuthnProvider{host: provider}),
	}
}

// InternalAuthnManager extracts the internal authn manager from a public handle.
func InternalAuthnManager(m AuthnProviderManager) authnprovidermgr.AuthnProviderManagerInterface {
	return m.(*authnProviderManager).mgr
}
