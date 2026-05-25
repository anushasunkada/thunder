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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveJOSEOptions(t *testing.T) {
	cfg := ResolveJOSEOptions([]Option{WithPreferredKeyID("kid-1")})
	assert.Equal(t, "kid-1", cfg.PreferredKeyID)
}

func TestInitializeWithoutRegistration(t *testing.T) {
	prev := registerJOSEInitializer
	t.Cleanup(func() { registerJOSEInitializer = prev })
	registerJOSEInitializer = nil

	jwtSvc, jweSvc, err := InitializeJose(nil)
	assert.Error(t, err)
	assert.Nil(t, jwtSvc)
	assert.Nil(t, jweSvc)
}
