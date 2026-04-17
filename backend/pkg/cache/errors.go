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

import "errors"

// Sentinel errors returned by CacheProvider implementations.
// Callers should use errors.Is for comparison.
var (
	// ErrKeyNotFound is returned by Expire when the target key does not exist.
	ErrKeyNotFound = errors.New("cache: key not found")

	// ErrProviderClosed is returned when an operation is attempted after Close
	// has been called.
	ErrProviderClosed = errors.New("cache: provider is closed")

	// ErrSerialization is returned when a value cannot be marshalled to or
	// unmarshalled from the backing store's wire format.
	ErrSerialization = errors.New("cache: serialization error")

	// ErrInvalidConfig is returned by New or backend constructors when the
	// supplied Config is missing required fields or contains invalid values.
	ErrInvalidConfig = errors.New("cache: invalid configuration")
)
