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

package authz

import (
	"time"

	"github.com/thunder-id/thunderid/internal/system/database/provider"
)

// AuthorizationRequestStore stores OAuth authorization request context.
type AuthorizationRequestStore = authorizationRequestStoreInterface

// NewSQLAuthorizationCodeStore creates a SQL authorization code store with explicit dependencies.
func NewSQLAuthorizationCodeStore(
	dbProvider provider.DBProviderInterface, deploymentID string,
) AuthorizationCodeStoreInterface {
	return &authorizationCodeStore{dbProvider: dbProvider, deploymentID: deploymentID}
}

// NewRedisAuthorizationCodeStore creates a Redis authorization code store with explicit dependencies.
func NewRedisAuthorizationCodeStore(
	p provider.RedisProviderInterface, deploymentID string,
) AuthorizationCodeStoreInterface {
	return &redisAuthorizationCodeStore{
		client:       p.GetRedisClient(),
		keyPrefix:    p.GetKeyPrefix(),
		deploymentID: deploymentID,
	}
}

// NewSQLAuthRequestStore creates a SQL authorization request store with explicit dependencies.
func NewSQLAuthRequestStore(dbProvider provider.DBProviderInterface, deploymentID string) AuthorizationRequestStore {
	return &authorizationRequestStore{
		dbProvider:     dbProvider,
		validityPeriod: 10 * time.Minute,
		deploymentID:   deploymentID,
	}
}

// NewRedisAuthRequestStore creates a Redis authorization request store with explicit dependencies.
func NewRedisAuthRequestStore(p provider.RedisProviderInterface, deploymentID string) AuthorizationRequestStore {
	return &redisAuthorizationRequestStore{
		client:       p.GetRedisClient(),
		keyPrefix:    p.GetKeyPrefix(),
		deploymentID: deploymentID,
	}
}
