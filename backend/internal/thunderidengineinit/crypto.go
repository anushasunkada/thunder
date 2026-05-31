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

package thunderidengineinit

import (
	"errors"
	"fmt"

	"github.com/thunder-id/thunderid/internal/system/kmprovider"
	kmprovidercommon "github.com/thunder-id/thunderid/internal/system/kmprovider/common"
	"github.com/thunder-id/thunderid/internal/system/kmprovider/defaultkm/pki"
)

func initCryptoProvider(signingKeyPath string) (kmprovidercommon.RuntimeCryptoProvider, error) {
	if signingKeyPath != "" {
		return initCryptoFromKeyFile(signingKeyPath)
	}
	pkiService, err := pki.Initialize()
	if err != nil {
		return nil, fmt.Errorf(
			"crypto provider required: set EngineConfig.Crypto.SigningKeyPath or run with ThunderID server config: %w",
			err)
	}
	runtimeCrypto, _, err := kmprovider.Initialize(pkiService)
	if err != nil {
		return nil, err
	}
	return runtimeCrypto, nil
}

func initCryptoFromKeyFile(keyPath string) (kmprovidercommon.RuntimeCryptoProvider, error) {
	pkiService, err := pki.InitializeFromFiles(keyPath)
	if err != nil {
		return nil, err
	}
	runtimeCrypto, _, err := kmprovider.Initialize(pkiService)
	if err != nil {
		return nil, err
	}
	if runtimeCrypto == nil {
		return nil, errors.New("failed to initialize runtime crypto provider")
	}
	return runtimeCrypto, nil
}
