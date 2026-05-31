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

// Package consent adapts ThunderID consent services for engine host wiring.
package consent

import (
	"context"
	"fmt"

	consentauthn "github.com/thunder-id/thunderid/internal/authn/consent"
	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/host"
)

type thunderConsentEnforcer struct {
	enforcer consentauthn.ConsentEnforcerServiceInterface
}

// InitializeConsentEnforcer creates a ThunderID ConsentEnforcer adapter.
func InitializeConsentEnforcer(enforcer consentauthn.ConsentEnforcerServiceInterface) host.ConsentEnforcer {
	return &thunderConsentEnforcer{enforcer: enforcer}
}

func (c *thunderConsentEnforcer) ResolveConsent(ctx context.Context, ouID, appID, agentID, userID string,
	requestedScopes []string) (*host.ConsentResolution, error) {
	prompt, svcErr := c.enforcer.ResolveConsent(ctx, ouID, appID, appID, userID,
		requestedScopes, nil, nil, nil)
	if svcErr != nil {
		return nil, asError(svcErr)
	}
	if prompt == nil {
		return &host.ConsentResolution{Required: false}, nil
	}
	items := make([]host.ConsentItem, 0)
	for _, purpose := range prompt.Purposes {
		for _, element := range purpose.Essential {
			items = append(items, host.ConsentItem{
				ID:          element.Name,
				DisplayName: element.Name,
				Required:    true,
			})
		}
		for _, element := range purpose.Optional {
			items = append(items, host.ConsentItem{
				ID:          element.Name,
				DisplayName: element.Name,
				Required:    false,
			})
		}
	}
	return &host.ConsentResolution{Required: len(items) > 0, Items: items}, nil
}

func (c *thunderConsentEnforcer) RecordConsent(ctx context.Context, ouID, appID, userID string,
	decisions []host.ConsentDecision) error {
	purposeDecisions := make([]consentauthn.PurposeDecision, 0, len(decisions))
	elements := make([]consentauthn.ElementDecision, 0, len(decisions))
	for _, decision := range decisions {
		elements = append(elements, consentauthn.ElementDecision{
			Name:     decision.ID,
			Approved: decision.Granted,
		})
	}
	if len(elements) > 0 {
		purposeDecisions = append(purposeDecisions, consentauthn.PurposeDecision{
			PurposeName: appID,
			Approved:    true,
			Elements:    elements,
		})
	}
	_, svcErr := c.enforcer.RecordConsent(ctx, ouID, appID, userID,
		&consentauthn.ConsentDecisions{Purposes: purposeDecisions}, "", 0)
	return asError(svcErr)
}

func asError(svcErr *serviceerror.ServiceError) error {
	if svcErr == nil {
		return nil
	}
	return fmt.Errorf("%s", svcErr.ErrorDescription.DefaultValue)
}
