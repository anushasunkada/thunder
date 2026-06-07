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

import sysconfig "github.com/thunder-id/thunderid/internal/system/config"

func gateClientFromConfig(cfg EngineConfig) *sysconfig.GateClientConfig {
	if cfg.GateClient.Hostname == "" {
		return nil
	}

	loginPath := cfg.GateClient.LoginPath
	if loginPath == "" {
		loginPath = "/signin"
	}
	errorPath := cfg.GateClient.ErrorPath
	if errorPath == "" {
		errorPath = "/error"
	}

	return &sysconfig.GateClientConfig{
		Scheme:    cfg.GateClient.Scheme,
		Hostname:  cfg.GateClient.Hostname,
		Port:      cfg.GateClient.Port,
		LoginPath: loginPath,
		ErrorPath: errorPath,
	}
}
