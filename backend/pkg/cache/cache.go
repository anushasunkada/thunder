/*
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
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

package cache

import (
	"fmt"

	"github.com/asgardeo/thunder/pkg/cache/inmemory"
	"github.com/asgardeo/thunder/pkg/cache/redis"
)

// New creates and returns a CacheProvider[T] for the backend specified in cfg.
//
// Supported backends:
//
//	TypeRedis    – connects to Redis using cfg.Redis; returns an error if the
//	               server is unreachable at construction time.
//	TypeInMemory – creates an in-process LRU/TTL cache; never fails with a
//	               connectivity error.
//
// cfg.Name must be non-empty. If cfg.Type is empty, TypeInMemory is used.
//
// Example:
//
//	p, err := cache.New[UserSession](cache.Config{
//	    Type:       cache.TypeRedis,
//	    Name:       "sessions",
//	    DefaultTTL: 30 * time.Minute,
//	    Redis:      cache.RedisConfig{Address: "localhost:6379"},
//	})
func New[T any](cfg Config) (CacheProvider[T], error) {
	if cfg.Name == "" {
		return nil, fmt.Errorf("%w: Name must not be empty", ErrInvalidConfig)
	}

	if cfg.Type == "" {
		cfg.Type = TypeInMemory
	}

	switch cfg.Type {
	case TypeRedis:
		return redis.NewProvider[T](cfg.Name, cfg.DefaultTTL, cfg.Redis)
	case TypeInMemory:
		return inmemory.NewProvider[T](cfg.Name, cfg.DefaultTTL, cfg.InMemory)
	default:
		return nil, fmt.Errorf("%w: unknown cache type %q", ErrInvalidConfig, cfg.Type)
	}
}
