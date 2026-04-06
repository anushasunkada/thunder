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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseACRValues_SingleACR(t *testing.T) {
	result := parseACRValues("urn:thunder:acr:password")
	assert.Equal(t, []string{"urn:thunder:acr:password"}, result)
}

func TestParseACRValues_MultipleACRs(t *testing.T) {
	result := parseACRValues("urn:thunder:acr:password urn:thunder:acr:generated-code")
	assert.Equal(t, []string{"urn:thunder:acr:password", "urn:thunder:acr:generated-code"}, result)
}

func TestParseACRValues_DeduplicatesPreservingFirstOccurrence(t *testing.T) {
	result := parseACRValues("urn:thunder:acr:generated-code urn:thunder:acr:generated-code urn:thunder:acr:password")
	assert.Equal(t, []string{"urn:thunder:acr:generated-code", "urn:thunder:acr:password"}, result)
}

func TestParseACRValues_EmptyString(t *testing.T) {
	result := parseACRValues("")
	assert.Empty(t, result)
}

func TestParseACRValues_OnlyWhitespace(t *testing.T) {
	result := parseACRValues("   ")
	assert.Empty(t, result)
}

func TestParseACRValues_ExtraSpacesBetweenACRs(t *testing.T) {
	result := parseACRValues("urn:thunder:acr:password   urn:thunder:acr:generated-code")
	assert.Equal(t, []string{"urn:thunder:acr:password", "urn:thunder:acr:generated-code"}, result)
}

func TestParseACRValues_PreservesOrder(t *testing.T) {
	result := parseACRValues("urn:thunder:acr:biometrics urn:thunder:acr:password urn:thunder:acr:generated-code")
	assert.Equal(t, []string{
		"urn:thunder:acr:biometrics",
		"urn:thunder:acr:password",
		"urn:thunder:acr:generated-code",
	}, result)
}

func TestResolveACRValues_NoRequest_NoDefaults(t *testing.T) {
	assert.Equal(t, "", ResolveACRValues("", nil))
}

func TestResolveACRValues_NoRequest_FallsBackToDefaults(t *testing.T) {
	defaults := []string{"urn:thunder:acr:password", "urn:thunder:acr:generated-code"}
	result := ResolveACRValues("", defaults)
	assert.ElementsMatch(t, defaults, strings.Fields(result))
}

func TestResolveACRValues_AllRequestedInDefaults_PreservesRequestedOrder(t *testing.T) {
	defaults := []string{"urn:thunder:acr:password", "urn:thunder:acr:generated-code"}
	result := ResolveACRValues("urn:thunder:acr:generated-code urn:thunder:acr:password", defaults)
	assert.Equal(t,
		[]string{"urn:thunder:acr:generated-code", "urn:thunder:acr:password"},
		strings.Fields(result))
}

func TestResolveACRValues_SomeNotInDefaults_FiltersOutUnknown(t *testing.T) {
	defaults := []string{"urn:thunder:acr:password", "urn:thunder:acr:generated-code"}
	result := ResolveACRValues("urn:thunder:acr:password urn:thunder:acr:biometrics", defaults)
	assert.Equal(t, []string{"urn:thunder:acr:password"}, strings.Fields(result))
}

func TestResolveACRValues_NoneInDefaults_FallsBackToDefaults(t *testing.T) {
	defaults := []string{"urn:thunder:acr:password", "urn:thunder:acr:generated-code"}
	result := ResolveACRValues("urn:thunder:acr:biometrics urn:thunder:acr:linked-wallet", defaults)
	assert.ElementsMatch(t, defaults, strings.Fields(result))
}

func TestResolveACRValues_DuplicatesDeduped(t *testing.T) {
	defaults := []string{"urn:thunder:acr:password", "urn:thunder:acr:generated-code"}
	result := ResolveACRValues(
		"urn:thunder:acr:password urn:thunder:acr:password urn:thunder:acr:generated-code", defaults)
	assert.Equal(t,
		[]string{"urn:thunder:acr:password", "urn:thunder:acr:generated-code"},
		strings.Fields(result))
}

func TestResolveACRValues_RequestPresent_NoDefaults_ReturnsEmpty(t *testing.T) {
	result := ResolveACRValues("urn:thunder:acr:password urn:thunder:acr:generated-code", nil)
	assert.Equal(t, "", result)
}
