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
	"github.com/thunder-id/thunderid/internal/system/database/provider"
	"github.com/thunder-id/thunderid/internal/system/transaction"
)

// NewAuthorizationCodeStore selects the authorization code store for the configured runtime store type.
func NewAuthorizationCodeStore(runtimeStoreType, deploymentID string) AuthorizationCodeStoreInterface {
	if runtimeStoreType == provider.DataSourceTypeRedis {
		return newRedisAuthorizationCodeStore(provider.GetRedisProvider(), deploymentID)
	}
	return newAuthorizationCodeStore(deploymentID)
}

// NewAuthorizationRequestStore selects the authorization request store for the configured runtime store type.
func NewAuthorizationRequestStore(runtimeStoreType, deploymentID string) AuthorizationRequestStoreInterface {
	if runtimeStoreType == provider.DataSourceTypeRedis {
		return newRedisAuthorizationRequestStore(provider.GetRedisProvider(), deploymentID)
	}
	return newAuthorizationRequestStore(deploymentID)
}

// NewAuthorizationStores creates authorization code and request stores for the configured runtime store type.
// For Redis, the returned transactioner is a no-op; for SQL, it is the runtime DB transactioner.
func NewAuthorizationStores(
	runtimeStoreType, deploymentID string,
) (AuthorizationCodeStoreInterface, AuthorizationRequestStoreInterface, transaction.Transactioner, error) {
	if runtimeStoreType == provider.DataSourceTypeRedis {
		redisProvider := provider.GetRedisProvider()
		return newRedisAuthorizationCodeStore(redisProvider, deploymentID),
			newRedisAuthorizationRequestStore(redisProvider, deploymentID),
			transaction.NewNoOpTransactioner(),
			nil
	}
	dbProvider := provider.GetDBProvider()
	transactioner, err := dbProvider.GetRuntimeDBTransactioner()
	if err != nil {
		return nil, nil, nil, err
	}
	return newAuthorizationCodeStore(deploymentID),
		newAuthorizationRequestStore(deploymentID), transactioner, nil
}
