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

package requestvalidator

import (
	"slices"
	"strings"
)

// ResolveACRValues resolves the effective acr_values for an authorization request.
//
// When requestedAcrValues is present and appAcrValues is non-empty, requested ACRs are filtered
// to only those included in appAcrValues. Requested ACRs not in the app's list are
// silently ignored. If filtering removes all ACRs (or no ACRs were requested), the full
// appAcrValues list is returned. The result is a space-separated, order-preserving,
// deduplicated string suitable for forwarding to the authentication flow engine.
func ResolveACRValues(requestedAcrValues string, appAcrValues []string) string {
	requested := parseACRValues(requestedAcrValues)
	filtered := make([]string, 0, len(requested))
	for _, acr := range requested {
		if slices.Contains(appAcrValues, acr) {
			filtered = append(filtered, acr)
		}
	}
	if len(filtered) == 0 {
		return strings.Join(appAcrValues, " ")
	}
	return strings.Join(filtered, " ")
}

// parseACRValues splits a space-separated acr_values string into a deduplicated, order-preserving slice.
func parseACRValues(acrValues string) []string {
	parts := strings.Fields(acrValues)
	seen := make(map[string]bool, len(parts))
	result := make([]string, 0, len(parts))
	for _, acr := range parts {
		if !seen[acr] {
			seen[acr] = true
			result = append(result, acr)
		}
	}
	return result
}
