/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 */

package host

// AssertionConfig mirrors assertion settings used when resolving flow context.
type AssertionConfig struct {
	ValidityPeriod int64
	UserAttributes []string
}

// LoginConsentConfig mirrors login consent settings used when resolving flow context.
type LoginConsentConfig struct {
	ValidityPeriod int64
}

// EntityInboundProfile is the portable inbound row shape flow execution needs for an entity.
type EntityInboundProfile struct {
	ID                        string
	AuthFlowID                string
	RegistrationFlowID        string
	IsRegistrationFlowEnabled bool
	RecoveryFlowID            string
	IsRecoveryFlowEnabled     bool
	Assertion                 *AssertionConfig
	LoginConsent              *LoginConsentConfig
	AllowedUserTypes          []string
	Properties                map[string]interface{}
}
