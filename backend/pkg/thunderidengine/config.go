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

// Config configures a thunderidengine instance.
type Config struct {
	Host         Host
	Executors    ExecutorRegistry
	Issuer       string
	DeploymentID string
}

//nolint:unused // reserved for engine bootstrap in a follow-up change.
func (c *Config) validate() error {
	if c.Issuer == "" {
		return ErrInvalidConfig
	}
	if c.Executors == nil {
		return ErrInvalidConfig
	}
	if c.Host.ClientProvider == nil ||
		c.Host.AuthnProvider == nil ||
		c.Host.AuthzProvider == nil ||
		c.Host.ResourceProvider == nil ||
		c.Host.OUProvider == nil ||
		c.Host.IDPProvider == nil ||
		c.Host.FlowDefinitionProvider == nil ||
		c.Host.RuntimeStore == nil ||
		c.Host.ConsentProvider == nil ||
		c.Host.DesignProvider == nil ||
		c.Host.I18n == nil ||
		c.Host.Crypto == nil {
		return ErrMissingHostField
	}
	return nil
}
