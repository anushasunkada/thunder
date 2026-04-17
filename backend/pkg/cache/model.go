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

// CacheType identifies the backing store implementation.
type CacheType string

const (
	// TypeRedis uses a Redis server as the backing store.
	TypeRedis CacheType = "redis"

	// TypeInMemory uses an in-process LRU map with TTL eviction.
	TypeInMemory CacheType = "inmemory"
)

// EvictionPolicy controls how entries are selected for removal when the
// in-memory cache reaches its maximum capacity.
type EvictionPolicy string

const (
	// EvictionLRU removes the least-recently-used entry first.
	EvictionLRU EvictionPolicy = "lru"

	// EvictionLFU removes the least-frequently-used entry first.
	// Note: for simplicity the current in-memory implementation approximates
	// LFU using access counters; it falls back to LRU for equal frequencies.
	EvictionLFU EvictionPolicy = "lfu"
)

// Stats is a point-in-time snapshot of a provider's operational counters.
// All fields are cumulative since the provider was created.
type Stats struct {
	// Hits is the number of Get calls that found a valid entry.
	Hits int64

	// Misses is the number of Get calls that found no entry (or an expired one).
	Misses int64

	// Sets is the number of successful Set calls.
	Sets int64

	// Deletes is the number of keys removed by Delete or Clear.
	Deletes int64
}

// HitRate returns the cache hit ratio in the range [0, 1].
// Returns 0 if no Get calls have been made yet.
func (s Stats) HitRate() float64 {
	total := s.Hits + s.Misses
	if total == 0 {
		return 0
	}
	return float64(s.Hits) / float64(total)
}
