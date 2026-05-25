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

package thunderidengine

import (
	gocrypto "crypto"
	"crypto/tls"
)

// KeyRef identifies a cryptographic key by its ID.
type KeyRef struct {
	KeyID string
}

// Algorithm identifies a signing or encryption algorithm (JWA-aligned, RFC 7518).
type Algorithm string

// JWA algorithm identifiers used for signing and encryption.
const (
	AlgorithmRS256        Algorithm = "RS256"
	AlgorithmRS512        Algorithm = "RS512"
	AlgorithmPS256        Algorithm = "PS256"
	AlgorithmES256        Algorithm = "ES256"
	AlgorithmES384        Algorithm = "ES384"
	AlgorithmES512        Algorithm = "ES512"
	AlgorithmEdDSA        Algorithm = "EdDSA"
	AlgorithmRSAOAEP256   Algorithm = "RSA-OAEP-256"
	AlgorithmRSAOAEP      Algorithm = "RSA-OAEP"
	AlgorithmECDHES       Algorithm = "ECDH-ES"
	AlgorithmECDHESA128KW Algorithm = "ECDH-ES+A128KW"
	AlgorithmECDHESA192KW Algorithm = "ECDH-ES+A192KW"
	AlgorithmECDHESA256KW Algorithm = "ECDH-ES+A256KW"
	AlgorithmAESGCM       Algorithm = "AES-GCM"
	AlgorithmA128KW       Algorithm = "A128KW"
	AlgorithmA192KW       Algorithm = "A192KW"
	AlgorithmA256KW       Algorithm = "A256KW"
	AlgorithmA128GCMKW    Algorithm = "A128GCMKW"
	AlgorithmA192GCMKW    Algorithm = "A192GCMKW"
	AlgorithmA256GCMKW    Algorithm = "A256GCMKW"
)

// SignAlgorithm identifies a signing algorithm.
type SignAlgorithm string

// Sign algorithm identifiers used for credential signing.
const (
	RSASHA256    SignAlgorithm = "RSA-SHA256"
	RSASHA512    SignAlgorithm = "RSA-SHA512"
	RSAPSSSHA256 SignAlgorithm = "RSA-PSS-SHA256"
	ECDSASHA256  SignAlgorithm = "ECDSA-SHA256"
	ECDSASHA384  SignAlgorithm = "ECDSA-SHA384"
	ECDSASHA512  SignAlgorithm = "ECDSA-SHA512"
	ED25519      SignAlgorithm = "ED25519"
)

// AlgorithmParams carries algorithm-specific parameters for encrypt/decrypt.
type AlgorithmParams struct {
	Algorithm  Algorithm
	RSAOAEP256 RSAOAEP256Params
	RSAOAEP    RSAOAEPParams
	ECDHES     ECDHESParams
	AESKW      AESKWParams
	AESGCMKW   AESGCMKWParams
}

// RSAOAEP256Params carries RSA-OAEP-256-specific inputs.
type RSAOAEP256Params struct {
	ContentEncryptionAlgorithm Algorithm
}

// RSAOAEPParams carries RSA-OAEP (SHA-1)-specific inputs.
type RSAOAEPParams struct {
	ContentEncryptionAlgorithm Algorithm
}

// AESKWParams carries AES Key Wrap-specific inputs.
type AESKWParams struct {
	ContentEncryptionAlgorithm Algorithm
}

// AESGCMKWParams carries AES-GCM Key Wrap-specific inputs.
type AESGCMKWParams struct {
	ContentEncryptionAlgorithm Algorithm
	IV                         []byte
	Tag                        []byte
}

// ECDHESParams carries ECDH-ES-specific inputs.
type ECDHESParams struct {
	EPK                        gocrypto.PublicKey
	ContentEncryptionAlgorithm Algorithm
}

// PublicKeyFilter specifies criteria for filtering public keys.
type PublicKeyFilter struct {
	KeyID     string
	Algorithm Algorithm
}

// PublicKeyInfo describes a public key returned by GetPublicKeys.
type PublicKeyInfo struct {
	KeyID          string
	Algorithm      Algorithm
	PublicKey      gocrypto.PublicKey
	Thumbprint     string
	CertificateDER []byte
}

// TLSMaterial holds TLS certificate material for a key reference.
type TLSMaterial struct {
	Certificate tls.Certificate
}
