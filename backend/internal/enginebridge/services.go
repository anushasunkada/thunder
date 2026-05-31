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
	"fmt"

	"github.com/thunder-id/thunderid/internal/authn/assert"
	consentauthn "github.com/thunder-id/thunderid/internal/authn/consent"
	authnprovidercm "github.com/thunder-id/thunderid/internal/authnprovider/common"
	"github.com/thunder-id/thunderid/internal/authz"
	"github.com/thunder-id/thunderid/internal/cert"
	"github.com/thunder-id/thunderid/internal/consent"
	flowcommon "github.com/thunder-id/thunderid/internal/flow/common"
	flowmgt "github.com/thunder-id/thunderid/internal/flow/mgt"
	"github.com/thunder-id/thunderid/internal/inboundclient"
	inboundmodel "github.com/thunder-id/thunderid/internal/inboundclient/model"
	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
	"github.com/thunder-id/thunderid/internal/system/jose/jwt"
	kmprovider "github.com/thunder-id/thunderid/internal/system/kmprovider/common"
)

type authzAdapter struct {
	provider AuthorizationSource
}

// NewAuthzService adapts AuthorizationSource to authz.AuthorizationServiceInterface.
func NewAuthzService(provider AuthorizationSource) authz.AuthorizationServiceInterface {
	return &authzAdapter{provider: provider}
}

func (a *authzAdapter) GetAuthorizedPermissions(ctx context.Context, request authz.GetAuthorizedPermissionsRequest) (
	*authz.GetAuthorizedPermissionsResponse, *serviceerror.ServiceError) {
	resp, err := a.provider.GetAuthorizedPermissions(ctx, AuthorizationRequest{
		EntityID:             request.EntityID,
		RequestedPermissions: request.RequestedPermissions,
	})
	if err != nil {
		return nil, &serviceerror.InternalServerError
	}
	return &authz.GetAuthorizedPermissionsResponse{AuthorizedPermissions: resp.AuthorizedPermissions}, nil
}

type consentAdapter struct {
	enforcer ConsentSource
}

// NewConsentEnforcer adapts ConsentSource to consentauthn.ConsentEnforcerServiceInterface.
func NewConsentEnforcer(enforcer ConsentSource) consentauthn.ConsentEnforcerServiceInterface {
	return &consentAdapter{enforcer: enforcer}
}

func (c *consentAdapter) ResolveConsent(ctx context.Context, ouID, appID, appName, userID string,
	essentialAttributes, optionalAttributes, authorizedPermissions []string,
	availableAttributes *authnprovidercm.AttributesResponse,
) (*consentauthn.ConsentPromptData, *serviceerror.ServiceError) {
	resolution, err := c.enforcer.ResolveConsent(ctx, ouID, appID, appName, userID, essentialAttributes)
	if err != nil {
		return nil, &serviceerror.InternalServerError
	}
	if resolution == nil || !resolution.Required {
		return nil, nil
	}
	prompt := &consentauthn.ConsentPromptData{Purposes: []consentauthn.ConsentPurposePrompt{}}
	for _, item := range resolution.Items {
		element := consentauthn.PromptElement{Name: item.ID}
		purpose := consentauthn.ConsentPurposePrompt{PurposeName: appID, PurposeID: item.ID}
		if item.Required {
			purpose.Essential = []consentauthn.PromptElement{element}
		} else {
			purpose.Optional = []consentauthn.PromptElement{element}
		}
		prompt.Purposes = append(prompt.Purposes, purpose)
	}
	return prompt, nil
}

func (c *consentAdapter) RecordConsent(ctx context.Context, ouID, appID, userID string,
	decisions *consentauthn.ConsentDecisions, sessionToken string, validityPeriod int64,
) (*consent.Consent, *serviceerror.ServiceError) {
	if decisions == nil {
		return nil, nil
	}
	publicDecisions := make([]ConsentDecision, 0)
	for _, purpose := range decisions.Purposes {
		for _, element := range purpose.Elements {
			publicDecisions = append(publicDecisions, ConsentDecision{
				ID:      element.Name,
				Granted: element.Approved,
			})
		}
	}
	if err := c.enforcer.RecordConsent(ctx, ouID, appID, userID, publicDecisions); err != nil {
		return nil, &serviceerror.InternalServerError
	}
	return nil, nil
}

type inboundClientAdapter struct {
	actors ActorSource
}

// NewInboundClientService adapts ActorSource to inboundclient.InboundClientServiceInterface.
func NewInboundClientService(actors ActorSource) inboundclient.InboundClientServiceInterface {
	return &inboundClientAdapter{actors: actors}
}

func (a *inboundClientAdapter) GetInboundClientByEntityID(
	ctx context.Context, entityID string,
) (*inboundmodel.InboundClient, error) {
	client, err := a.actors.GetInboundClientByEntityID(ctx, entityID)
	if err != nil {
		return nil, err
	}
	return &inboundmodel.InboundClient{ID: client.ClientID}, nil
}

func (a *inboundClientAdapter) CreateInboundClient(ctx context.Context, client *inboundmodel.InboundClient,
	appCert *inboundmodel.Certificate, oauthProfile *inboundmodel.OAuthProfile,
	hasClientSecret bool, entityName string) error {
	return fmt.Errorf("not supported in engine mode")
}

func (a *inboundClientAdapter) GetInboundClientList(ctx context.Context) ([]inboundmodel.InboundClient, error) {
	return nil, fmt.Errorf("not supported in engine mode")
}

func (a *inboundClientAdapter) UpdateInboundClient(ctx context.Context, client *inboundmodel.InboundClient,
	appCert *inboundmodel.Certificate, oauthProfile *inboundmodel.OAuthProfile,
	hasClientSecret bool, oauthClientID string, entityName string) error {
	return fmt.Errorf("not supported in engine mode")
}

func (a *inboundClientAdapter) DeleteInboundClient(ctx context.Context, entityID string) error {
	return fmt.Errorf("not supported in engine mode")
}

func (a *inboundClientAdapter) Validate(ctx context.Context, client *inboundmodel.InboundClient,
	oauthProfile *inboundmodel.OAuthProfile, hasClientSecret bool) error {
	return fmt.Errorf("not supported in engine mode")
}

func (a *inboundClientAdapter) GetCertificate(
	ctx context.Context, refType cert.CertificateReferenceType, refID string,
) (*inboundmodel.Certificate, *inboundclient.CertOperationError) {
	return nil, &inboundclient.CertOperationError{
		Operation:  inboundclient.CertOpRetrieve,
		RefType:    refType,
		Underlying: &serviceerror.InternalServerError,
	}
}

func (a *inboundClientAdapter) ResolveInboundAuthProfileHandles(ctx context.Context,
	profile *inboundmodel.InboundAuthProfile) error {
	return fmt.Errorf("not supported in engine mode")
}

func (a *inboundClientAdapter) GetOAuthProfileByEntityID(
	ctx context.Context, entityID string,
) (*inboundmodel.OAuthProfile, error) {
	return nil, fmt.Errorf("not supported in engine mode")
}

func (a *inboundClientAdapter) GetOAuthClientByClientID(
	ctx context.Context, clientID string,
) (*inboundmodel.OAuthClient, error) {
	client, err := a.actors.GetInboundClientByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	return &inboundmodel.OAuthClient{ID: client.ClientID, ClientID: client.ClientID}, nil
}

func (a *inboundClientAdapter) IsDeclarative(ctx context.Context, entityID string) bool {
	return false
}

func (a *inboundClientAdapter) LoadDeclarativeResources(
	ctx context.Context, cfg inboundmodel.DeclarativeLoaderConfig,
) error {
	return fmt.Errorf("not supported in engine mode")
}

type flowMgtAdapter struct {
	service flowmgt.FlowMgtServiceInterface
}

// NewFlowSource adapts FlowMgtServiceInterface to FlowSource.
func NewFlowSource(service flowmgt.FlowMgtServiceInterface) FlowSource {
	return &flowMgtAdapter{service: service}
}

func (a *flowMgtAdapter) GetFlow(ctx context.Context, flowID string) (*FlowDefinition, error) {
	flow, svcErr := a.service.GetFlow(ctx, flowID)
	if svcErr != nil {
		return nil, fmt.Errorf("%s", svcErr.ErrorDescription.DefaultValue)
	}
	return &FlowDefinition{
		ID:       flow.ID,
		Handle:   flow.Handle,
		FlowType: string(flow.FlowType),
	}, nil
}

func (a *flowMgtAdapter) GetFlowByHandle(ctx context.Context, handle, flowType string) (*FlowDefinition, error) {
	flow, svcErr := a.service.GetFlowByHandle(ctx, handle, flowcommon.FlowType(flowType))
	if svcErr != nil {
		return nil, fmt.Errorf("%s", svcErr.ErrorDescription.DefaultValue)
	}
	return &FlowDefinition{
		ID:       flow.ID,
		Handle:   flow.Handle,
		FlowType: string(flow.FlowType),
	}, nil
}

// NewJWTService creates a JWT service from a runtime crypto provider.
func NewJWTService(crypto kmprovider.RuntimeCryptoProvider) (jwt.JWTServiceInterface, error) {
	if crypto == nil {
		return nil, fmt.Errorf("crypto provider is required")
	}
	return jwt.Initialize(crypto)
}

// NewAuthAssertGenerator creates an auth assertion generator for engine mode.
func NewAuthAssertGenerator() assert.AuthAssertGeneratorInterface {
	return assert.Initialize()
}

// NewCryptoProvider returns the provided runtime crypto provider.
func NewCryptoProvider(crypto kmprovider.RuntimeCryptoProvider) kmprovider.RuntimeCryptoProvider {
	return crypto
}
