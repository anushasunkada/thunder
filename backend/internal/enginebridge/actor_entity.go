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

// Package enginebridge adapts ThunderID engine host providers to internal services.
package enginebridge

import (
	"encoding/json"

	"github.com/thunder-id/thunderid/internal/entityprovider"
)

type actorEntityProvider struct {
	actors ActorSource
}

// NewEntityProvider adapts ActorSource to entityprovider.EntityProviderInterface.
func NewEntityProvider(actors ActorSource) entityprovider.EntityProviderInterface {
	return &actorEntityProvider{actors: actors}
}

func engineNotSupported() *entityprovider.EntityProviderError {
	return entityprovider.NewEntityProviderError(
		entityprovider.ErrorCodeNotImplemented, "engine", "not supported in engine mode")
}

func (p *actorEntityProvider) IdentifyEntity(
	filters map[string]interface{},
) (*string, *entityprovider.EntityProviderError) {
	id, err := p.actors.IdentifyEntity(filters)
	if err != nil {
		return nil, entityprovider.NewEntityProviderError(
			entityprovider.ErrorCodeSystemError, "engine", err.Error())
	}
	return id, nil
}

func (p *actorEntityProvider) SearchEntities(
	filters map[string]interface{},
) ([]*entityprovider.Entity, *entityprovider.EntityProviderError) {
	actors, err := p.actors.SearchEntities(filters)
	if err != nil {
		return nil, entityprovider.NewEntityProviderError(
			entityprovider.ErrorCodeSystemError, "engine", err.Error())
	}
	entities := make([]*entityprovider.Entity, 0, len(actors))
	for _, actor := range actors {
		if actor == nil {
			continue
		}
		entities = append(entities, &entityprovider.Entity{
			ID:         actor.ID,
			Type:       actor.EntityType,
			Attributes: actor.Attributes,
		})
	}
	return entities, nil
}

func (p *actorEntityProvider) GetEntity(entityID string) (*entityprovider.Entity, *entityprovider.EntityProviderError) {
	actor, err := p.actors.GetEntity(entityID)
	if err != nil {
		return nil, entityprovider.NewEntityProviderError(
			entityprovider.ErrorCodeSystemError, "engine", err.Error())
	}
	return &entityprovider.Entity{
		ID:         actor.ID,
		Type:       actor.EntityType,
		Attributes: actor.Attributes,
	}, nil
}

func (p *actorEntityProvider) CreateEntity(
	*entityprovider.Entity, json.RawMessage,
) (*entityprovider.Entity, *entityprovider.EntityProviderError) {
	return nil, engineNotSupported()
}

func (p *actorEntityProvider) UpdateEntity(
	string, *entityprovider.Entity,
) (*entityprovider.Entity, *entityprovider.EntityProviderError) {
	return nil, engineNotSupported()
}

func (p *actorEntityProvider) DeleteEntity(string) *entityprovider.EntityProviderError {
	return engineNotSupported()
}

func (p *actorEntityProvider) UpdateCredentials(string, json.RawMessage) *entityprovider.EntityProviderError {
	return engineNotSupported()
}

func (p *actorEntityProvider) UpdateAttributes(string, json.RawMessage) *entityprovider.EntityProviderError {
	return engineNotSupported()
}

func (p *actorEntityProvider) UpdateSystemAttributes(string, json.RawMessage) *entityprovider.EntityProviderError {
	return engineNotSupported()
}

func (p *actorEntityProvider) UpdateSystemCredentials(string, json.RawMessage) *entityprovider.EntityProviderError {
	return engineNotSupported()
}

func (p *actorEntityProvider) GetTransitiveEntityGroups(
	string,
) ([]entityprovider.EntityGroup, *entityprovider.EntityProviderError) {
	return nil, engineNotSupported()
}

func (p *actorEntityProvider) ValidateEntityIDs(
	[]string,
) ([]string, *entityprovider.EntityProviderError) {
	return nil, engineNotSupported()
}

func (p *actorEntityProvider) GetEntitiesByIDs(
	[]string,
) ([]entityprovider.Entity, *entityprovider.EntityProviderError) {
	return nil, engineNotSupported()
}

func (p *actorEntityProvider) GetEntityListCount(
	entityprovider.EntityCategory, map[string]interface{},
) (int, *entityprovider.EntityProviderError) {
	return 0, engineNotSupported()
}

func (p *actorEntityProvider) GetEntityList(
	entityprovider.EntityCategory, int, int, map[string]interface{},
) ([]entityprovider.Entity, *entityprovider.EntityProviderError) {
	return nil, engineNotSupported()
}
