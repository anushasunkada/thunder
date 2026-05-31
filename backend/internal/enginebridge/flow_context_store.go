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
	"encoding/json"
	"errors"
	"time"

	"github.com/thunder-id/thunderid/internal/flow/flowexec"
)

type flowContextStoreAdapter struct {
	store RuntimeStore
}

// NewFlowContextStore adapts a RuntimeStore to flowexec's flow context store.
func NewFlowContextStore(store RuntimeStore) flowexec.RuntimeFlowContextStore {
	return &flowContextStoreAdapter{store: store}
}

func (a *flowContextStoreAdapter) StoreFlowContext(
	ctx context.Context, dbModel flowexec.FlowContextDB, expirySeconds int64,
) error {
	data, err := json.Marshal(dbModel)
	if err != nil {
		return err
	}
	expiry := dbModel.ExpiryTime
	if expiry.IsZero() && expirySeconds > 0 {
		expiry = time.Now().UTC().Add(time.Duration(expirySeconds) * time.Second)
	}
	if expiry.IsZero() && !dbModel.CreatedAt.IsZero() {
		expiry = dbModel.CreatedAt
	}
	return a.store.StoreFlowContext(ctx, dbModel.ExecutionID, data, expiry)
}

func (a *flowContextStoreAdapter) GetFlowContext(
	ctx context.Context, executionID string,
) (*flowexec.FlowContextDB, error) {
	data, err := a.store.GetFlowContext(ctx, executionID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	var model flowexec.FlowContextDB
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, err
	}
	return &model, nil
}

func (a *flowContextStoreAdapter) UpdateFlowContext(ctx context.Context, dbModel flowexec.FlowContextDB) error {
	data, err := json.Marshal(dbModel)
	if err != nil {
		return err
	}
	return a.store.UpdateFlowContext(ctx, dbModel.ExecutionID, data)
}

func (a *flowContextStoreAdapter) DeleteFlowContext(ctx context.Context, executionID string) error {
	return a.store.DeleteFlowContext(ctx, executionID)
}
