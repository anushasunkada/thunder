/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 */

package hostbridge

import (
	"github.com/thunder-id/thunderid/internal/application"
	"github.com/thunder-id/thunderid/internal/inboundclient"
	"github.com/thunder-id/thunderid/pkg/oauth/host"
)

// NewThunderApplication wraps Thunder's application service for pkg/oauth Dependencies.Application.
func NewThunderApplication(svc application.ApplicationServiceInterface) host.DCRApplication {
	if svc == nil {
		return nil
	}
	return &ThunderDCRApplication{ApplicationService: svc}
}

// NewThunderInbound wraps Thunder's inbound client service for pkg/oauth Dependencies.Inbound.
func NewThunderInbound(svc inboundclient.InboundClientServiceInterface) host.InboundOAuth {
	if svc == nil {
		return nil
	}
	return &ThunderInboundOAuth{InboundService: svc}
}
