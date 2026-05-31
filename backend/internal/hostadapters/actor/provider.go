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

// Package actor adapts ThunderID server actor services for engine host wiring.
package actor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/thunder-id/thunderid/internal/application"
	"github.com/thunder-id/thunderid/internal/entityprovider"
	"github.com/thunder-id/thunderid/internal/entitytype"
	"github.com/thunder-id/thunderid/internal/inboundclient"
	inboundmodel "github.com/thunder-id/thunderid/internal/inboundclient/model"
	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/host"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/runtime"
)

type thunderActorProvider struct {
	entities     entityprovider.EntityProviderInterface
	applications application.ApplicationServiceInterface
	inbound      inboundclient.InboundClientServiceInterface
	entityTypes  entitytype.EntityTypeServiceInterface
}

// InitializeActorProvider creates a ThunderID ActorProvider adapter.
func InitializeActorProvider(
	entities entityprovider.EntityProviderInterface,
	applications application.ApplicationServiceInterface,
	inbound inboundclient.InboundClientServiceInterface,
	entityTypes entitytype.EntityTypeServiceInterface,
) host.ActorProvider {
	return &thunderActorProvider{
		entities:     entities,
		applications: applications,
		inbound:      inbound,
		entityTypes:  entityTypes,
	}
}

func (p *thunderActorProvider) IdentifyEntity(filters map[string]interface{}) (*string, error) {
	id, provErr := p.entities.IdentifyEntity(filters)
	if provErr != nil {
		return nil, errors.New(provErr.Message)
	}
	return id, nil
}

func (p *thunderActorProvider) GetEntity(entityID string) (*host.Actor, error) {
	entity, provErr := p.entities.GetEntity(entityID)
	if provErr != nil {
		return nil, errors.New(provErr.Message)
	}
	if entity == nil {
		return nil, runtime.ErrNotFound
	}
	return &host.Actor{
		ID:         entity.ID,
		EntityType: entity.Type,
		Attributes: entity.Attributes,
	}, nil
}

func (p *thunderActorProvider) SearchEntities(filters map[string]interface{}) ([]*host.Actor, error) {
	entities, provErr := p.entities.SearchEntities(filters)
	if provErr != nil {
		return nil, errors.New(provErr.Message)
	}
	actors := make([]*host.Actor, 0, len(entities))
	for _, entity := range entities {
		if entity == nil {
			continue
		}
		actors = append(actors, &host.Actor{
			ID:         entity.ID,
			EntityType: entity.Type,
			Attributes: entity.Attributes,
		})
	}
	return actors, nil
}

func (p *thunderActorProvider) GetApplication(ctx context.Context, appID string) (*host.Application, error) {
	app, svcErr := p.applications.GetApplication(ctx, appID)
	if svcErr != nil {
		return nil, asError(svcErr)
	}
	return &host.Application{
		ID:       app.ID,
		Name:     app.Name,
		OUID:     app.OUID,
		EntityID: app.ID,
	}, nil
}

func (p *thunderActorProvider) GetInboundClientByEntityID(
	ctx context.Context, entityID string,
) (*host.InboundClient, error) {
	client, err := p.inbound.GetInboundClientByEntityID(ctx, entityID)
	if err != nil {
		return nil, err
	}
	return toInboundClient(client), nil
}

func (p *thunderActorProvider) GetInboundClientByClientID(
	ctx context.Context, clientID string,
) (*host.InboundClient, error) {
	clients, err := p.inbound.GetInboundClientList(ctx)
	if err != nil {
		return nil, err
	}
	for i := range clients {
		if clients[i].ID == clientID {
			return toInboundClient(&clients[i]), nil
		}
	}
	return nil, runtime.ErrNotFound
}

func (p *thunderActorProvider) GetEntityType(ctx context.Context, typeID string) (*host.EntityType, error) {
	entityType, svcErr := p.entityTypes.GetEntityType(ctx, entitytype.TypeCategoryUser, typeID, false)
	if svcErr != nil {
		return nil, asError(svcErr)
	}
	raw, _ := json.Marshal(entityType)
	return &host.EntityType{ID: entityType.ID, Name: entityType.Name, Attributes: raw}, nil
}

func toInboundClient(client *inboundmodel.InboundClient) *host.InboundClient {
	if client == nil {
		return nil
	}
	return &host.InboundClient{
		ClientID:      client.ID,
		EntityID:      client.ID,
		ApplicationID: client.ID,
	}
}

func asError(svcErr *serviceerror.ServiceError) error {
	if svcErr == nil {
		return nil
	}
	return fmt.Errorf("%s", svcErr.ErrorDescription.DefaultValue)
}
