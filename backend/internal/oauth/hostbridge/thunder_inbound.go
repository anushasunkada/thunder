/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 */

package hostbridge

import (
	"context"
	"fmt"

	"github.com/thunder-id/thunderid/internal/inboundclient"
	"github.com/thunder-id/thunderid/pkg/oauth/host"
	"github.com/thunder-id/thunderid/pkg/oauth/oauthclient"
)

// ThunderInboundOAuth implements host.InboundOAuth using Thunder's inbound client service.
type ThunderInboundOAuth struct {
	InboundService inboundclient.InboundClientServiceInterface
}

// GetOAuthClientByClientID implements host.InboundOAuth.
func (t *ThunderInboundOAuth) GetOAuthClientByClientID(ctx context.Context, clientID string) (*oauthclient.Client, error) {
	if t == nil || t.InboundService == nil {
		return nil, fmt.Errorf("inbound service is required")
	}
	o, err := t.InboundService.GetOAuthClientByClientID(ctx, clientID)
	if err != nil || o == nil {
		return nil, err
	}
	return oauthClientFromInbound(o), nil
}

var _ host.InboundOAuth = (*ThunderInboundOAuth)(nil)
