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
	"errors"
	"fmt"

	"github.com/thunder-id/thunderid/internal/inboundclient"
	flowhost "github.com/thunder-id/thunderid/pkg/flow/host"
)

// ThunderInboundFlow adapts Thunder's inbound client service to pkg/flow/host.InboundFlow.
type ThunderInboundFlow struct {
	InboundService inboundclient.InboundClientServiceInterface
}

// NewThunderInboundFlow returns host.InboundFlow backed by Thunder's inbound service.
func NewThunderInboundFlow(svc inboundclient.InboundClientServiceInterface) flowhost.InboundFlow {
	if svc == nil {
		return nil
	}
	return &ThunderInboundFlow{InboundService: svc}
}

// GetInboundClientByEntityID implements flowhost.InboundFlow.
func (t *ThunderInboundFlow) GetInboundClientByEntityID(ctx context.Context, entityID string) (*flowhost.EntityInboundProfile, error) {
	if t == nil || t.InboundService == nil {
		return nil, fmt.Errorf("inbound service is required")
	}
	c, err := t.InboundService.GetInboundClientByEntityID(ctx, entityID)
	if err != nil {
		if errors.Is(err, inboundclient.ErrInboundClientNotFound) {
			return nil, fmt.Errorf("%w", flowhost.ErrEntityInboundNotFound)
		}
		return nil, err
	}
	return entityProfileFromInbound(c), nil
}

var _ flowhost.InboundFlow = (*ThunderInboundFlow)(nil)
