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
	"context"
	"crypto"
	"errors"

	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
)

// JWTService defines JWT signing and verification operations.
type JWTService interface {
	GenerateJWT(ctx context.Context, sub, iss string, validityPeriod int64,
		claims map[string]interface{}, typ, alg string) (string, int64, *serviceerror.ServiceError)
	VerifyJWT(jwtToken string, expectedAud, expectedIss string) *serviceerror.ServiceError
	VerifyJWTWithPublicKey(jwtToken string, jwtPublicKey crypto.PublicKey, expectedAud,
		expectedIss string) *serviceerror.ServiceError
	VerifyJWTWithJWKS(jwtToken, jwksURL, expectedAud, expectedIss string) *serviceerror.ServiceError
	VerifyJWTSignature(jwtToken string) *serviceerror.ServiceError
	VerifyJWTSignatureWithPublicKey(jwtToken string, jwtPublicKey crypto.PublicKey) *serviceerror.ServiceError
	VerifyJWTSignatureWithJWKS(jwtToken string, jwksURL string) *serviceerror.ServiceError
}

// JWEService defines JWE encryption and decryption operations.
type JWEService interface {
	Encrypt(payload []byte, recipientPublicKey crypto.PublicKey,
		alg KeyEncAlgorithm, enc ContentEncAlgorithm, cty string, kid string) (string, *serviceerror.ServiceError)
	Decrypt(jweToken string) ([]byte, *serviceerror.ServiceError)
}

// KeyEncAlgorithm represents the JWE key management algorithm (alg header parameter).
type KeyEncAlgorithm string

// JWE key management algorithm identifiers (alg header parameter).
const (
	RSAOAEP      KeyEncAlgorithm = "RSA-OAEP"
	RSAOAEP256   KeyEncAlgorithm = "RSA-OAEP-256"
	A128KW       KeyEncAlgorithm = "A128KW"
	A192KW       KeyEncAlgorithm = "A192KW"
	A256KW       KeyEncAlgorithm = "A256KW"
	ECDHES       KeyEncAlgorithm = "ECDH-ES"
	ECDHESA128KW KeyEncAlgorithm = "ECDH-ES+A128KW"
	ECDHESA192KW KeyEncAlgorithm = "ECDH-ES+A192KW"
	ECDHESA256KW KeyEncAlgorithm = "ECDH-ES+A256KW"
	A128GCMKW    KeyEncAlgorithm = "A128GCMKW"
	A192GCMKW    KeyEncAlgorithm = "A192GCMKW"
	A256GCMKW    KeyEncAlgorithm = "A256GCMKW"
)

// ContentEncAlgorithm represents the JWE content encryption algorithm (enc header parameter).
type ContentEncAlgorithm string

// JWE content encryption algorithm identifiers (enc header parameter).
const (
	A128CBCHS256 ContentEncAlgorithm = "A128CBC-HS256"
	A192CBCHS384 ContentEncAlgorithm = "A192CBC-HS384"
	A256CBCHS512 ContentEncAlgorithm = "A256CBC-HS512"
	A128GCM      ContentEncAlgorithm = "A128GCM"
	A192GCM      ContentEncAlgorithm = "A192GCM"
	A256GCM      ContentEncAlgorithm = "A256GCM"
)

// Option configures JOSE service initialization.
type Option func(*joseOptions)

type joseOptions struct {
	preferredKeyID string
}

// WithPreferredKeyID sets the signing and decryption key ID used by JWT and JWE services.
func WithPreferredKeyID(keyID string) Option {
	return func(o *joseOptions) {
		o.preferredKeyID = keyID
	}
}

// JOSEInitConfig holds resolved JOSE initialization settings.
type JOSEInitConfig struct {
	PreferredKeyID string
}

// ResolveJOSEOptions applies the given options.
func ResolveJOSEOptions(opts []Option) JOSEInitConfig {
	options := joseOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	return JOSEInitConfig{PreferredKeyID: options.preferredKeyID}
}

type joseInitializer func(RuntimeCryptoProvider, ...Option) (JWTService, JWEService, error)

var registerJOSEInitializer joseInitializer

// RegisterJOSEInitializer registers the default JOSE service factory. Called from the internal jose package.
func RegisterJOSEInitializer(initializer joseInitializer) {
	registerJOSEInitializer = initializer
}

// InitializeJose constructs JWT and JWE service instances backed by the host crypto provider.
func InitializeJose(crypto RuntimeCryptoProvider, opts ...Option) (JWTService, JWEService, error) {
	if registerJOSEInitializer == nil {
		return nil, nil, errors.New(
			"JOSE initializer not registered; import github.com/thunder-id/thunderid/internal/system/jose",
		)
	}
	return registerJOSEInitializer(crypto, opts...)
}
