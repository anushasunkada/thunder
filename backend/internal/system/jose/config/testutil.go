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

package config

import sysconfig "github.com/thunder-id/thunderid/internal/system/config"

// InitTestFromServerRuntime reloads JOSE config from the initialized server runtime.
// For use in tests after config.InitializeServerRuntime.
func InitTestFromServerRuntime() {
	Reset()
	cfg := sysconfig.GetServerRuntime().Config
	Set(FromSystemConfig(cfg, BuildOptionsForServer(&cfg)))
}

// InitTestServerRuntime initializes server runtime and JOSE config for tests.
func InitTestServerRuntime(serverHome string, cfg *sysconfig.Config) error {
	sysconfig.ResetServerRuntime()
	Reset()
	if err := sysconfig.InitializeServerRuntime(serverHome, cfg); err != nil {
		return err
	}
	Set(FromSystemConfig(*cfg, BuildOptionsForServer(cfg)))
	return nil
}
