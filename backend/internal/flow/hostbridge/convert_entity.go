/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 */

package hostbridge

import (
	inboundmodel "github.com/thunder-id/thunderid/internal/inboundclient/model"
	flowhost "github.com/thunder-id/thunderid/pkg/flow/host"
)

func entityProfileFromInbound(c *inboundmodel.InboundClient) *flowhost.EntityInboundProfile {
	if c == nil {
		return nil
	}
	p := &flowhost.EntityInboundProfile{
		ID:                        c.ID,
		AuthFlowID:                c.AuthFlowID,
		RegistrationFlowID:        c.RegistrationFlowID,
		IsRegistrationFlowEnabled: c.IsRegistrationFlowEnabled,
		RecoveryFlowID:            c.RecoveryFlowID,
		IsRecoveryFlowEnabled:     c.IsRecoveryFlowEnabled,
		AllowedUserTypes:          append([]string(nil), c.AllowedUserTypes...),
	}
	if c.Assertion != nil {
		p.Assertion = &flowhost.AssertionConfig{
			ValidityPeriod: c.Assertion.ValidityPeriod,
			UserAttributes: append([]string(nil), c.Assertion.UserAttributes...),
		}
	}
	if c.LoginConsent != nil {
		p.LoginConsent = &flowhost.LoginConsentConfig{ValidityPeriod: c.LoginConsent.ValidityPeriod}
	}
	if c.Properties != nil {
		p.Properties = make(map[string]interface{}, len(c.Properties))
		for k, v := range c.Properties {
			p.Properties[k] = v
		}
	}
	return p
}
