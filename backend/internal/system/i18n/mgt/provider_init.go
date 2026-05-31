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

package mgt

import (
	"github.com/thunder-id/thunderid/internal/system/config"
	serverconst "github.com/thunder-id/thunderid/internal/system/constants"
)

// InitializeProvider creates an i18n service without registering HTTP routes.
func InitializeProvider(translationConfig config.TranslationConfig) (I18nServiceInterface, error) {
	var store i18nStoreInterface

	storeMode := getI18nStoreMode(translationConfig)
	switch storeMode {
	case serverconst.StoreModeDeclarative:
		fileStore := newFileBasedStore()
		if err := loadDeclarativeResources(fileStore); err != nil {
			return nil, err
		}
		store = fileStore
	case serverconst.StoreModeComposite:
		fileStore := newFileBasedStore()
		if err := loadDeclarativeResources(fileStore); err != nil {
			return nil, err
		}
		store = newCompositeI18nStore(fileStore, newI18nStore())
	default:
		store = newI18nStore()
	}

	return newI18nService(store), nil
}
