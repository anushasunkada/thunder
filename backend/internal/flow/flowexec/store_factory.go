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

package flowexec

import (
	"github.com/thunder-id/thunderid/internal/system/database/provider"
)

// NewSQLFlowStore creates a SQL-backed flow context store with explicit dependencies.
func NewSQLFlowStore(dbProvider provider.DBProviderInterface, deploymentID string) RuntimeFlowContextStore {
	return asRuntimeFlowContextStore(&flowStore{
		dbProvider:   dbProvider,
		deploymentID: deploymentID,
	})
}

// NewRedisFlowStore creates a Redis-backed flow context store with explicit dependencies.
func NewRedisFlowStore(p provider.RedisProviderInterface, deploymentID string) RuntimeFlowContextStore {
	return asRuntimeFlowContextStore(&redisFlowStore{
		client:       p.GetRedisClient(),
		keyPrefix:    p.GetKeyPrefix(),
		deploymentID: deploymentID,
	})
}
