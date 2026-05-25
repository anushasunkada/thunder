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

package adapter

import (
	"context"
	"fmt"

	"github.com/thunder-id/thunderid/internal/flow/flowexec"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/authz"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/par"
	dbprovider "github.com/thunder-id/thunderid/internal/system/database/provider"
	"github.com/thunder-id/thunderid/internal/system/log"
	"github.com/thunder-id/thunderid/internal/system/transaction"
	"github.com/thunder-id/thunderid/pkg/thunderidengine"
)

type runtimeStore struct {
	transactioner transaction.Transactioner
	parStore      par.PARStoreInterface
	authReqStore  authz.AuthorizationRequestStoreInterface
	authCodeStore authz.AuthorizationCodeStoreInterface
	flowStore     flowexec.FlowStoreInterface
}

// NewRuntimeStore returns a thunderidengine.RuntimeStore backed by host runtime persistence.
func NewRuntimeStore(runtimeStoreType, deploymentID string) (thunderidengine.RuntimeStore, error) {
	var transactioner transaction.Transactioner
	if runtimeStoreType == dbprovider.DataSourceTypeRedis {
		transactioner = transaction.NewNoOpTransactioner()
	} else {
		dbProvider := dbprovider.GetDBProvider()
		var err error
		transactioner, err = dbProvider.GetRuntimeDBTransactioner()
		if err != nil {
			return nil, err
		}
	}
	return &runtimeStore{
		transactioner: transactioner,
		parStore:      par.NewPARRequestStore(runtimeStoreType, deploymentID),
		authReqStore:  authz.NewAuthorizationRequestStore(runtimeStoreType, deploymentID),
		authCodeStore: authz.NewAuthorizationCodeStore(runtimeStoreType, deploymentID),
		flowStore:     flowexec.NewFlowStore(runtimeStoreType),
	}, nil
}

func (s *runtimeStore) Store(
	ctx context.Context, request thunderidengine.PARRequest, expirySeconds int64,
) (string, error) {
	return s.parStore.Store(ctx, request, expirySeconds)
}

func (s *runtimeStore) Consume(
	ctx context.Context, randomKey string,
) (thunderidengine.PARRequest, bool, error) {
	return s.parStore.Consume(ctx, randomKey)
}

func (s *runtimeStore) AddRequest(
	ctx context.Context, value thunderidengine.AuthRequestContext,
) (string, error) {
	return s.authReqStore.AddRequest(ctx, value)
}

func (s *runtimeStore) GetRequest(
	ctx context.Context, key string,
) (bool, thunderidengine.AuthRequestContext, error) {
	return s.authReqStore.GetRequest(ctx, key)
}

func (s *runtimeStore) ClearRequest(ctx context.Context, key string) error {
	return s.authReqStore.ClearRequest(ctx, key)
}

func (s *runtimeStore) InsertAuthorizationCode(
	ctx context.Context, authzCode thunderidengine.AuthorizationCode,
) error {
	return s.authCodeStore.InsertAuthorizationCode(ctx, authzCode)
}

func (s *runtimeStore) ConsumeAuthorizationCode(ctx context.Context, authCode string) (bool, error) {
	return s.authCodeStore.ConsumeAuthorizationCode(ctx, authCode)
}

func (s *runtimeStore) GetAuthorizationCode(
	ctx context.Context, authCode string,
) (*thunderidengine.AuthorizationCode, error) {
	return s.authCodeStore.GetAuthorizationCode(ctx, authCode)
}

func (s *runtimeStore) StoreFlowContext(
	ctx context.Context, dbModel thunderidengine.FlowContextDB, expirySeconds int64,
) error {
	txErr := s.transactioner.Transact(ctx, func(txCtx context.Context) error {
		return s.flowStore.StoreFlowContext(txCtx, flowexec.FlowContextDB(dbModel), expirySeconds)
	})
	if txErr != nil {
		return fmt.Errorf("failed to store flow context in database: %w", txErr)
	}

	log.GetLogger().Debug("Flow context stored successfully in database",
		log.String("executionID", dbModel.ExecutionID), log.Error(txErr))
	return nil
}

func (s *runtimeStore) GetFlowContext(
	ctx context.Context, executionID string,
) (*thunderidengine.FlowContextDB, error) {
	dbModel, err := s.flowStore.GetFlowContext(ctx, executionID)
	if dbModel == nil {
		return nil, err
	}
	engineModel := thunderidengine.FlowContextDB(*dbModel)
	return &engineModel, err
}

func (s *runtimeStore) UpdateFlowContext(
	ctx context.Context, dbModel thunderidengine.FlowContextDB,
) error {
	txErr := s.transactioner.Transact(ctx, func(txCtx context.Context) error {
		return s.flowStore.UpdateFlowContext(ctx, flowexec.FlowContextDB(dbModel))
	})
	if txErr != nil {
		return fmt.Errorf("failed to update flow context in database: %w", txErr)
	}

	log.GetLogger().Debug("Flow context updated successfully in database",
		log.String("executionID", dbModel.ExecutionID), log.Error(txErr))
	return nil
}

func (s *runtimeStore) DeleteFlowContext(ctx context.Context, executionID string) error {
	txErr := s.transactioner.Transact(ctx, func(txCtx context.Context) error {
		return s.flowStore.DeleteFlowContext(ctx, executionID)
	})
	if txErr != nil {
		return fmt.Errorf("failed to remove flow context from database: %w", txErr)
	}

	log.GetLogger().Debug("Flow context removed successfully from database",
		log.String("executionID", executionID), log.Error(txErr))
	return nil
}
