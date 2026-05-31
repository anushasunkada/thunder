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
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/thunder-id/thunderid/internal/enginebridge"
	flowcommon "github.com/thunder-id/thunderid/internal/flow/common"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/host"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/runtime"
)

type actorBridge struct {
	inner host.ActorProvider
}

func (a *actorBridge) IdentifyEntity(filters map[string]interface{}) (*string, error) {
	return a.inner.IdentifyEntity(filters)
}

func (a *actorBridge) GetEntity(entityID string) (*enginebridge.Actor, error) {
	actor, err := a.inner.GetEntity(entityID)
	if err != nil {
		return nil, err
	}
	return &enginebridge.Actor{
		ID:         actor.ID,
		EntityType: actor.EntityType,
		Attributes: actor.Attributes,
	}, nil
}

func (a *actorBridge) SearchEntities(filters map[string]interface{}) ([]*enginebridge.Actor, error) {
	actors, err := a.inner.SearchEntities(filters)
	if err != nil {
		return nil, err
	}
	out := make([]*enginebridge.Actor, 0, len(actors))
	for _, actor := range actors {
		if actor == nil {
			continue
		}
		out = append(out, &enginebridge.Actor{
			ID:         actor.ID,
			EntityType: actor.EntityType,
			Attributes: actor.Attributes,
		})
	}
	return out, nil
}

func (a *actorBridge) GetApplication(ctx context.Context, appID string) (*enginebridge.Application, error) {
	app, err := a.inner.GetApplication(ctx, appID)
	if err != nil {
		return nil, err
	}
	return &enginebridge.Application{
		ID:       app.ID,
		Name:     app.Name,
		OUID:     app.OUID,
		EntityID: app.EntityID,
	}, nil
}

func (a *actorBridge) GetInboundClientByEntityID(
	ctx context.Context, entityID string,
) (*enginebridge.InboundClient, error) {
	client, err := a.inner.GetInboundClientByEntityID(ctx, entityID)
	if err != nil {
		return nil, err
	}
	return mapInboundClient(client), nil
}

func (a *actorBridge) GetInboundClientByClientID(
	ctx context.Context, clientID string,
) (*enginebridge.InboundClient, error) {
	client, err := a.inner.GetInboundClientByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	return mapInboundClient(client), nil
}

func mapInboundClient(client *host.InboundClient) *enginebridge.InboundClient {
	if client == nil {
		return nil
	}
	return &enginebridge.InboundClient{
		ClientID:                           client.ClientID,
		EntityID:                           client.EntityID,
		ApplicationID:                      client.ApplicationID,
		OUID:                               client.OUID,
		Secret:                             client.Secret,
		GrantTypes:                         client.GrantTypes,
		RedirectURIs:                       client.RedirectURIs,
		ResponseTypes:                      client.ResponseTypes,
		TokenEndpointAuthMethod:            client.TokenEndpointAuthMethod,
		PKCERequired:                       client.PKCERequired,
		PublicClient:                       client.PublicClient,
		RequirePushedAuthorizationRequests: client.RequirePushedAuthorizationRequests,
		AuthFlowID:                         client.AuthFlowID,
		RegistrationFlowID:                 client.RegistrationFlowID,
		IsRegistrationFlowEnabled:          client.IsRegistrationFlowEnabled,
		RecoveryFlowID:                     client.RecoveryFlowID,
		IsRecoveryFlowEnabled:              client.IsRecoveryFlowEnabled,
	}
}

func (a *actorBridge) GetEntityType(ctx context.Context, typeID string) (*enginebridge.EntityType, error) {
	entityType, err := a.inner.GetEntityType(ctx, typeID)
	if err != nil {
		return nil, err
	}
	return &enginebridge.EntityType{
		ID:         entityType.ID,
		Name:       entityType.Name,
		Attributes: entityType.Attributes,
	}, nil
}

type runtimeBridge struct {
	inner runtime.Store
}

func (r *runtimeBridge) StoreFlowContext(ctx context.Context, executionID string, data []byte, expiry time.Time) error {
	return r.inner.StoreFlowContext(ctx, executionID, data, expiry)
}

func (r *runtimeBridge) GetFlowContext(ctx context.Context, executionID string) ([]byte, error) {
	data, err := r.inner.GetFlowContext(ctx, executionID)
	if err != nil {
		if errors.Is(err, runtime.ErrNotFound) {
			return nil, enginebridge.ErrNotFound
		}
		return nil, err
	}
	return data, nil
}

func (r *runtimeBridge) UpdateFlowContext(ctx context.Context, executionID string, data []byte) error {
	return r.inner.UpdateFlowContext(ctx, executionID, data)
}

func (r *runtimeBridge) DeleteFlowContext(ctx context.Context, executionID string) error {
	return r.inner.DeleteFlowContext(ctx, executionID)
}

func (r *runtimeBridge) StoreAuthCode(ctx context.Context, code string, data []byte, expiry time.Time) error {
	return r.inner.StoreAuthCode(ctx, code, data, expiry)
}

func (r *runtimeBridge) GetAuthCode(ctx context.Context, code string) ([]byte, error) {
	data, err := r.inner.GetAuthCode(ctx, code)
	if errors.Is(err, runtime.ErrNotFound) {
		return nil, enginebridge.ErrNotFound
	}
	return data, err
}

func (r *runtimeBridge) DeleteAuthCode(ctx context.Context, code string) error {
	err := r.inner.DeleteAuthCode(ctx, code)
	if errors.Is(err, runtime.ErrNotFound) {
		return enginebridge.ErrNotFound
	}
	return err
}

func (r *runtimeBridge) StoreAuthRequest(ctx context.Context, requestID string, data []byte, expiry time.Time) error {
	return r.inner.StoreAuthRequest(ctx, requestID, data, expiry)
}

func (r *runtimeBridge) GetAuthRequest(ctx context.Context, requestID string) ([]byte, error) {
	data, err := r.inner.GetAuthRequest(ctx, requestID)
	if errors.Is(err, runtime.ErrNotFound) {
		return nil, enginebridge.ErrNotFound
	}
	return data, err
}

func (r *runtimeBridge) DeleteAuthRequest(ctx context.Context, requestID string) error {
	return r.inner.DeleteAuthRequest(ctx, requestID)
}

func (r *runtimeBridge) StorePAR(ctx context.Context, requestURI string, data []byte, expiry time.Time) error {
	return r.inner.StorePAR(ctx, requestURI, data, expiry)
}

func (r *runtimeBridge) GetPAR(ctx context.Context, requestURI string) ([]byte, error) {
	data, err := r.inner.GetPAR(ctx, requestURI)
	if errors.Is(err, runtime.ErrNotFound) {
		return nil, enginebridge.ErrNotFound
	}
	return data, err
}

func (r *runtimeBridge) DeletePAR(ctx context.Context, requestURI string) error {
	err := r.inner.DeletePAR(ctx, requestURI)
	if errors.Is(err, runtime.ErrNotFound) {
		return enginebridge.ErrNotFound
	}
	return err
}

func (r *runtimeBridge) StoreJTI(ctx context.Context, jti string, expiry time.Time) error {
	return r.inner.StoreJTI(ctx, jti, expiry)
}

func (r *runtimeBridge) ExistsJTI(ctx context.Context, jti string) (bool, error) {
	return r.inner.ExistsJTI(ctx, jti)
}

func (r *runtimeBridge) StoreAttributeCache(ctx context.Context, id string, data []byte, expiry time.Time) error {
	return r.inner.StoreAttributeCache(ctx, id, data, expiry)
}

func (r *runtimeBridge) GetAttributeCache(ctx context.Context, id string) ([]byte, error) {
	data, err := r.inner.GetAttributeCache(ctx, id)
	if errors.Is(err, runtime.ErrNotFound) {
		return nil, enginebridge.ErrNotFound
	}
	return data, err
}

func (r *runtimeBridge) ExtendAttributeCacheExpiry(ctx context.Context, id string, expiry time.Time) error {
	err := r.inner.ExtendAttributeCacheExpiry(ctx, id, expiry)
	if errors.Is(err, runtime.ErrNotFound) {
		return enginebridge.ErrNotFound
	}
	return err
}

func (r *runtimeBridge) DeleteAttributeCache(ctx context.Context, id string) error {
	err := r.inner.DeleteAttributeCache(ctx, id)
	if errors.Is(err, runtime.ErrNotFound) {
		return enginebridge.ErrNotFound
	}
	return err
}

type authzBridge struct {
	inner host.AuthorizationProvider
}

func (a *authzBridge) GetAuthorizedPermissions(ctx context.Context, req enginebridge.AuthorizationRequest) (
	*enginebridge.AuthorizationResponse, error) {
	resp, err := a.inner.GetAuthorizedPermissions(ctx, host.GetAuthorizedPermissionsRequest{
		EntityID:             req.EntityID,
		RequestedPermissions: req.RequestedPermissions,
	})
	if err != nil {
		return nil, err
	}
	return &enginebridge.AuthorizationResponse{AuthorizedPermissions: resp.AuthorizedPermissions}, nil
}

type consentBridge struct {
	inner host.ConsentEnforcer
}

func (c *consentBridge) ResolveConsent(ctx context.Context, ouID, appID, agentID, userID string,
	requestedScopes []string) (*enginebridge.ConsentResolution, error) {
	resolution, err := c.inner.ResolveConsent(ctx, ouID, appID, agentID, userID, requestedScopes)
	if err != nil {
		return nil, err
	}
	if resolution == nil {
		return &enginebridge.ConsentResolution{}, nil
	}
	items := make([]enginebridge.ConsentItem, 0, len(resolution.Items))
	for _, item := range resolution.Items {
		items = append(items, enginebridge.ConsentItem{
			ID:          item.ID,
			DisplayName: item.DisplayName,
			Description: item.Description,
			Required:    item.Required,
		})
	}
	return &enginebridge.ConsentResolution{Required: resolution.Required, Items: items}, nil
}

func (c *consentBridge) RecordConsent(ctx context.Context, ouID, appID, userID string,
	decisions []enginebridge.ConsentDecision) error {
	publicDecisions := make([]host.ConsentDecision, 0, len(decisions))
	for _, decision := range decisions {
		publicDecisions = append(publicDecisions, host.ConsentDecision{
			ID:      decision.ID,
			Granted: decision.Granted,
		})
	}
	return c.inner.RecordConsent(ctx, ouID, appID, userID, publicDecisions)
}

type flowBridge struct {
	inner host.FlowProvider
}

func (f *flowBridge) GetFlow(ctx context.Context, flowID string) (*enginebridge.FlowDefinition, error) {
	flow, err := f.inner.GetFlow(ctx, flowID)
	if err != nil {
		return nil, err
	}
	return hostFlowDefinitionToEngine(flow)
}

func (f *flowBridge) GetFlowByHandle(
	ctx context.Context, handle, flowType string,
) (*enginebridge.FlowDefinition, error) {
	flow, err := f.inner.GetFlowByHandle(ctx, handle, host.FlowType(flowType))
	if err != nil {
		return nil, err
	}
	return hostFlowDefinitionToEngine(flow)
}

func hostFlowDefinitionToEngine(flow *host.FlowDefinition) (*enginebridge.FlowDefinition, error) {
	if flow == nil {
		return nil, nil
	}
	def := &enginebridge.FlowDefinition{
		ID:       flow.ID,
		Handle:   flow.Handle,
		FlowType: string(flow.FlowType),
	}
	if len(flow.Graph) == 0 {
		return def, nil
	}
	var complete flowcommon.CompleteFlowDefinition
	if err := json.Unmarshal(flow.Graph, &complete); err != nil {
		return nil, err
	}
	def.Name = complete.Name
	if len(complete.Nodes) > 0 {
		def.Nodes = complete.Nodes
	}
	return def, nil
}

// WrapRuntimeStore adapts runtime.Store for internal engine wiring (import-cycle boundary).
func WrapRuntimeStore(store runtime.Store) enginebridge.RuntimeStore {
	return &runtimeBridge{inner: store}
}

// WrapActorProvider adapts host.ActorProvider for internal engine wiring.
func WrapActorProvider(actors host.ActorProvider) enginebridge.ActorSource {
	return &actorBridge{inner: actors}
}

// WrapAuthorization adapts host.AuthorizationProvider for internal engine wiring.
func WrapAuthorization(authz host.AuthorizationProvider) enginebridge.AuthorizationSource {
	return &authzBridge{inner: authz}
}

// WrapConsent adapts host.ConsentEnforcer for internal engine wiring.
func WrapConsent(consent host.ConsentEnforcer) enginebridge.ConsentSource {
	return &consentBridge{inner: consent}
}

// WrapFlowProvider adapts host.FlowProvider for internal engine wiring.
func WrapFlowProvider(flow host.FlowProvider) enginebridge.FlowSource {
	return &flowBridge{inner: flow}
}
