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

// SetOption is a functional option for the CacheProvider.Set method.
// Implementations call ApplySetOptions to materialise the final values.
type SetOption func(*SetOptions)

// SetOptions holds per-call overrides accumulated from a slice of SetOption
// values. It is exported so that backend implementations in sub-packages
// can read the resolved values without import cycles.
type SetOptions struct {
	// TTL overrides the provider's DefaultTTL for a single Set call.
	// nil means "use the provider's configured default".
	TTL *time.Duration
}

// ApplySetOptions merges opts into a fresh SetOptions and returns it.
// Backend implementations should call this inside their Set method.
func ApplySetOptions(opts []SetOption) SetOptions {
	var o SetOptions
	for _, fn := range opts {
		fn(&o)
	}
	return o
}

// WithTTL overrides the provider's DefaultTTL for a single Set call.
// Pass 0 to store the entry without expiry regardless of the provider default.
func WithTTL(ttl time.Duration) SetOption {
	return func(o *SetOptions) {
		o.TTL = &ttl
	}
}
