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

import "time"

// Config is the unified configuration for any CacheProvider implementation.
// Callers populate only the sub-struct that matches Config.Type; all other
// sub-structs are ignored by the chosen backend.
type Config struct {
	// Type selects the backend. Defaults to TypeInMemory if empty.
	Type CacheType

	// Name is a logical identifier for this cache instance. It is used to
	// namespace keys (e.g. as a key prefix) and appears in log output.
	// Required; must be non-empty.
	Name string

	// DefaultTTL is the expiry applied to Set calls that do not supply a
	// WithTTL option. A value of zero means entries never expire.
	DefaultTTL time.Duration

	// Redis holds configuration for the Redis backend.
	// Only used when Type == TypeRedis.
	Redis RedisConfig

	// InMemory holds configuration for the in-process backend.
	// Only used when Type == TypeInMemory.
	InMemory InMemoryConfig
}

// RedisConfig contains the connection and runtime parameters for the Redis
// backend. It intentionally mirrors the fields exposed by go-redis/v9's
// redis.Options so that callers can share a single config source.
type RedisConfig struct {
	// Address is the host:port of the Redis server.
	// Defaults to "localhost:6379".
	Address string

	// Username for Redis 6+ ACL authentication. Leave empty if not used.
	Username string

	// Password for Redis AUTH. Leave empty if not used.
	Password string

	// DB is the Redis logical database index to select. Defaults to 0.
	DB int

	// KeyPrefix is prepended to every key written by this provider, using
	// the format "<KeyPrefix>:<Name>:<key>". This allows multiple applications
	// or environments to share the same Redis instance without key collisions.
	// Example: "thunder:prod"
	KeyPrefix string

	// PoolSize is the maximum number of socket connections maintained in the
	// connection pool. 0 means use the go-redis default (10 per CPU).
	PoolSize int

	// DialTimeout is the timeout for establishing new connections.
	// 0 means use the go-redis default (5 s).
	DialTimeout time.Duration

	// ReadTimeout is the timeout for socket reads.
	// 0 means use the go-redis default (3 s).
	ReadTimeout time.Duration

	// WriteTimeout is the timeout for socket writes.
	// 0 means use the go-redis default (ReadTimeout).
	WriteTimeout time.Duration
}

// InMemoryConfig controls the behaviour of the in-process cache backend.
type InMemoryConfig struct {
	// MaxSize is the maximum number of entries the cache will hold before
	// eviction kicks in. 0 (the default) means unlimited.
	MaxSize int

	// EvictionPolicy determines which entry is removed when the cache is full.
	// Defaults to EvictionLRU when MaxSize > 0.
	EvictionPolicy EvictionPolicy

	// CleanupInterval is how often the background goroutine scans for and
	// removes expired entries. Defaults to 1 minute.
	// Set to a negative value to disable background cleanup entirely.
	CleanupInterval time.Duration
}
