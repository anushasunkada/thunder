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

package runtime

import (
	"github.com/thunder-id/thunderid/internal/attributecache"
	"github.com/thunder-id/thunderid/internal/enginebridge"
	"github.com/thunder-id/thunderid/internal/flow/flowexec"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/authz"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/jti"
	"github.com/thunder-id/thunderid/internal/oauth/oauth2/par"
	"github.com/thunder-id/thunderid/internal/system/database/provider"
	"github.com/thunder-id/thunderid/internal/system/transaction"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/runtime"
)

// StoreConfig configures ThunderID runtime store selection.
type StoreConfig struct {
	StoreType    string
	DeploymentID string
}

// InitializeRuntimeStore selects Redis or SQL runtime storage for ThunderID server use.
func InitializeRuntimeStore(cfg StoreConfig, dbProvider provider.DBProviderInterface,
	redisProvider provider.RedisProviderInterface) (runtime.Store, transaction.Transactioner, error) {
	if cfg.StoreType == provider.DataSourceTypeRedis {
		stores := enginebridge.LegacyRuntimeStores{
			FlowStore: flowexec.NewRedisFlowStore(redisProvider, cfg.DeploymentID),
			AuthCode:  authz.NewRedisAuthorizationCodeStore(redisProvider, cfg.DeploymentID),
			AuthReq:   authz.NewRedisAuthRequestStore(redisProvider, cfg.DeploymentID),
			PAR:       par.NewRedisPARStore(redisProvider, cfg.DeploymentID),
			JTI:       jti.NewRedisJTIStore(redisProvider, cfg.DeploymentID),
			AttributeCache: newAttributeCacheBridge(
				attributecache.NewRedisPersistStore(redisProvider, cfg.DeploymentID)),
		}
		return &runtimeStoreWrapper{inner: enginebridge.NewRuntimeStore(stores)},
			transaction.NewNoOpTransactioner(), nil
	}
	transactioner, err := dbProvider.GetRuntimeDBTransactioner()
	if err != nil {
		return nil, nil, err
	}
	stores := enginebridge.LegacyRuntimeStores{
		FlowStore:      flowexec.NewSQLFlowStore(dbProvider, cfg.DeploymentID),
		AuthCode:       authz.NewSQLAuthorizationCodeStore(dbProvider, cfg.DeploymentID),
		AuthReq:        authz.NewSQLAuthRequestStore(dbProvider, cfg.DeploymentID),
		PAR:            par.NewSQLPARStore(dbProvider, cfg.DeploymentID),
		JTI:            jti.NewSQLJTIStore(dbProvider, cfg.DeploymentID),
		AttributeCache: newAttributeCacheBridge(attributecache.NewSQLPersistStore(dbProvider, cfg.DeploymentID)),
	}
	return &runtimeStoreWrapper{inner: enginebridge.NewRuntimeStore(stores)}, transactioner, nil
}
