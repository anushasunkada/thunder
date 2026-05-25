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

// Package cryptolab provides pure cryptographic primitives: signing, encryption,
// decryption, hashing, and secure token utilities. Algorithm types are defined
// in pkg/thunderidengine; this package implements operations on those types.
package cryptolab

import (
	gocrypto "crypto"

	"github.com/thunder-id/thunderid/pkg/thunderidengine"
)

type (
	// Algorithm is a thunderidengine.Algorithm alias for use within cryptolab call sites.
	Algorithm = thunderidengine.Algorithm
	// SignAlgorithm is a thunderidengine.SignAlgorithm alias for use within cryptolab call sites.
	SignAlgorithm = thunderidengine.SignAlgorithm
	// AlgorithmParams is a thunderidengine.AlgorithmParams alias.
	AlgorithmParams = thunderidengine.AlgorithmParams
	// RSAOAEP256Params is a thunderidengine.RSAOAEP256Params alias.
	RSAOAEP256Params = thunderidengine.RSAOAEP256Params
	// RSAOAEPParams is a thunderidengine.RSAOAEPParams alias.
	RSAOAEPParams = thunderidengine.RSAOAEPParams
	// AESKWParams is a thunderidengine.AESKWParams alias.
	AESKWParams = thunderidengine.AESKWParams
	// AESGCMKWParams is a thunderidengine.AESGCMKWParams alias.
	AESGCMKWParams = thunderidengine.AESGCMKWParams
	// ECDHESParams is a thunderidengine.ECDHESParams alias.
	ECDHESParams = thunderidengine.ECDHESParams
)

// Algorithm constants alias thunderidengine algorithm identifiers.
const (
	AlgorithmRS256        = thunderidengine.AlgorithmRS256
	AlgorithmRS512        = thunderidengine.AlgorithmRS512
	AlgorithmPS256        = thunderidengine.AlgorithmPS256
	AlgorithmES256        = thunderidengine.AlgorithmES256
	AlgorithmES384        = thunderidengine.AlgorithmES384
	AlgorithmES512        = thunderidengine.AlgorithmES512
	AlgorithmEdDSA        = thunderidengine.AlgorithmEdDSA
	AlgorithmRSAOAEP256   = thunderidengine.AlgorithmRSAOAEP256
	AlgorithmECDHES       = thunderidengine.AlgorithmECDHES
	AlgorithmECDHESA128KW = thunderidengine.AlgorithmECDHESA128KW
	AlgorithmECDHESA256KW = thunderidengine.AlgorithmECDHESA256KW
	AlgorithmAESGCM       = thunderidengine.AlgorithmAESGCM
	AlgorithmRSAOAEP      = thunderidengine.AlgorithmRSAOAEP
	AlgorithmECDHESA192KW = thunderidengine.AlgorithmECDHESA192KW
	AlgorithmA128KW       = thunderidengine.AlgorithmA128KW
	AlgorithmA192KW       = thunderidengine.AlgorithmA192KW
	AlgorithmA256KW       = thunderidengine.AlgorithmA256KW
	AlgorithmA128GCMKW    = thunderidengine.AlgorithmA128GCMKW
	AlgorithmA192GCMKW    = thunderidengine.AlgorithmA192GCMKW
	AlgorithmA256GCMKW    = thunderidengine.AlgorithmA256GCMKW

	RSASHA256    = thunderidengine.RSASHA256
	RSASHA512    = thunderidengine.RSASHA512
	RSAPSSSHA256 = thunderidengine.RSAPSSSHA256
	ECDSASHA256  = thunderidengine.ECDSASHA256
	ECDSASHA384  = thunderidengine.ECDSASHA384
	ECDSASHA512  = thunderidengine.ECDSASHA512
	ED25519      = thunderidengine.ED25519
)

// CryptoDetails carries algorithm-specific outputs from an Encrypt operation.
// EPK is the generated ephemeral public key for ECDH-ES variants, to be embedded in the JWE header.
// CEK is the content encryption key generated or derived during key establishment.
// Nil CryptoDetails is returned for algorithms that produce no extra output (e.g. AES-GCM).
// For RSA-OAEP, RSA-OAEP-256 and ECDH-ES variants, both EPK (where applicable) and CEK are populated.
// CEK is nil for AES-GCM; EPK is nil for RSA-OAEP and RSA-OAEP-256 (no ephemeral key is generated).
// IV and Tag are set only for AES-GCM Key Wrap (A128GCMKW etc.) and must be embedded in the JWE protected header.
type CryptoDetails struct {
	EPK gocrypto.PublicKey // ECDH-ES variants only; nil for RSA-OAEP, RSA-OAEP-256 and AES-GCM
	CEK []byte             // Generated or derived CEK; nil for AES-GCM
	IV  []byte             // AES-GCM Key Wrap only: IV used to wrap the CEK
	Tag []byte             // AES-GCM Key Wrap only: authentication tag from wrapping the CEK
}
