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

package pki

import (
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/thunder-id/thunderid/internal/system/log"
)

// DefaultEngineKeyID is the PKI certificate ID used when loading a signing key from files.
const DefaultEngineKeyID = "engine-signing-key"

// InitializeFromFiles loads a single PEM key/certificate pair for engine embed mode.
// keyPath may be a private key file; a certificate is resolved from the same basename with .crt or .pem.
func InitializeFromFiles(keyPath string) (PKIServiceInterface, error) {
	certPath := resolveCertPath(keyPath)
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load signing key pair: %w", err)
	}
	if len(cert.Certificate) == 0 {
		return nil, fmt.Errorf("certificate file %s contains no certificates", certPath)
	}
	algorithm, err := getAlgorithmFromKey(cert.PrivateKey)
	if err != nil {
		return nil, err
	}
	thumbprint, err := getThumbprint(cert)
	if err != nil {
		return nil, err
	}
	return &pkiService{
		certificates: map[string]PKI{
			DefaultEngineKeyID: {
				ID:          DefaultEngineKeyID,
				Algorithm:   algorithm,
				PrivateKey:  cert.PrivateKey,
				Certificate: cert,
				ThumbPrint:  thumbprint,
			},
		},
		logger: log.GetLogger().With(log.String(log.LoggerKeyComponentName, "PKIService")),
	}, nil
}

func resolveCertPath(keyPath string) string {
	dir := filepath.Dir(keyPath)
	base := strings.TrimSuffix(filepath.Base(keyPath), filepath.Ext(keyPath))
	for _, ext := range []string{".crt", ".pem", ".cer"} {
		candidate := filepath.Join(dir, base+ext)
		if candidate != keyPath {
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
		}
	}
	return keyPath
}
