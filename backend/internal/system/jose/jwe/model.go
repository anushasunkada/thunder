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

package jwe

import "github.com/thunder-id/thunderid/pkg/thunderidengine"

// KeyEncAlgorithm is a thunderidengine.KeyEncAlgorithm alias.
type KeyEncAlgorithm = thunderidengine.KeyEncAlgorithm

// Key management algorithms (JWE alg header).
const (
	RSAOAEP      = thunderidengine.RSAOAEP
	RSAOAEP256   = thunderidengine.RSAOAEP256
	A128KW       = thunderidengine.A128KW
	A192KW       = thunderidengine.A192KW
	A256KW       = thunderidengine.A256KW
	ECDHES       = thunderidengine.ECDHES
	ECDHESA128KW = thunderidengine.ECDHESA128KW
	ECDHESA192KW = thunderidengine.ECDHESA192KW
	ECDHESA256KW = thunderidengine.ECDHESA256KW
	A128GCMKW    = thunderidengine.A128GCMKW
	A192GCMKW    = thunderidengine.A192GCMKW
	A256GCMKW    = thunderidengine.A256GCMKW
)

// ContentEncAlgorithm is a thunderidengine.ContentEncAlgorithm alias.
type ContentEncAlgorithm = thunderidengine.ContentEncAlgorithm

// Content encryption algorithms (JWE enc header).
const (
	A128CBCHS256 = thunderidengine.A128CBCHS256
	A192CBCHS384 = thunderidengine.A192CBCHS384
	A256CBCHS512 = thunderidengine.A256CBCHS512
	A128GCM      = thunderidengine.A128GCM
	A192GCM      = thunderidengine.A192GCM
	A256GCM      = thunderidengine.A256GCM
)
