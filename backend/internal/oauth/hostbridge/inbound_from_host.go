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

	"github.com/thunder-id/thunderid/internal/cert"
	"github.com/thunder-id/thunderid/internal/inboundclient"
	inboundmodel "github.com/thunder-id/thunderid/internal/inboundclient/model"
	"github.com/thunder-id/thunderid/pkg/oauth/host"
)

// InboundFromHost adapts pkg/oauth host inbound lookup to Thunder's inbound client interface.
func InboundFromHost(h host.InboundOAuth) inboundclient.InboundClientServiceInterface {
	if h == nil {
		return nil
	}
	if th, ok := h.(*ThunderInboundOAuth); ok {
		return th.InboundService
	}
	return &inboundHostBridge{host: h}
}

type inboundHostBridge struct {
	host host.InboundOAuth
}

func (b *inboundHostBridge) CreateInboundClient(ctx context.Context, client *inboundmodel.InboundClient,
	appCert *inboundmodel.Certificate, oauthProfile *inboundmodel.OAuthProfile, hasClientSecret bool, entityName string,
) error {
	return fmt.Errorf("inbound host bridge: CreateInboundClient is not supported; provide Thunder inbound for admin APIs")
}

func (b *inboundHostBridge) GetInboundClientByEntityID(ctx context.Context, entityID string) (*inboundmodel.InboundClient, error) {
	return nil, fmt.Errorf("inbound host bridge: GetInboundClientByEntityID is not supported; provide Thunder inbound for admin APIs")
}

func (b *inboundHostBridge) GetInboundClientList(ctx context.Context) ([]inboundmodel.InboundClient, error) {
	return nil, fmt.Errorf("inbound host bridge: GetInboundClientList is not supported; provide Thunder inbound for admin APIs")
}

func (b *inboundHostBridge) UpdateInboundClient(ctx context.Context, client *inboundmodel.InboundClient,
	appCert *inboundmodel.Certificate, oauthProfile *inboundmodel.OAuthProfile,
	hasClientSecret bool, oauthClientID string, entityName string,
) error {
	return fmt.Errorf("inbound host bridge: UpdateInboundClient is not supported; provide Thunder inbound for admin APIs")
}

func (b *inboundHostBridge) DeleteInboundClient(ctx context.Context, entityID string) error {
	return fmt.Errorf("inbound host bridge: DeleteInboundClient is not supported; provide Thunder inbound for admin APIs")
}

func (b *inboundHostBridge) Validate(ctx context.Context, client *inboundmodel.InboundClient,
	oauthProfile *inboundmodel.OAuthProfile, hasClientSecret bool,
) error {
	return fmt.Errorf("inbound host bridge: Validate is not supported; provide Thunder inbound for admin APIs")
}

func (b *inboundHostBridge) GetOAuthProfileByEntityID(ctx context.Context, entityID string) (*inboundmodel.OAuthProfile, error) {
	return nil, fmt.Errorf("inbound host bridge: GetOAuthProfileByEntityID is not supported; provide Thunder inbound for admin APIs")
}

func (b *inboundHostBridge) GetOAuthClientByClientID(ctx context.Context, clientID string) (*inboundmodel.OAuthClient, error) {
	if b == nil || b.host == nil {
		return nil, fmt.Errorf("inbound host is not configured")
	}
	c, err := b.host.GetOAuthClientByClientID(ctx, clientID)
	if err != nil || c == nil {
		return nil, err
	}
	return oauthClientToInbound(c)
}

func (b *inboundHostBridge) IsDeclarative(ctx context.Context, entityID string) bool {
	return false
}

func (b *inboundHostBridge) LoadDeclarativeResources(ctx context.Context, cfg inboundmodel.DeclarativeLoaderConfig) error {
	return fmt.Errorf("inbound host bridge: LoadDeclarativeResources is not supported; provide Thunder inbound for declarative mode")
}

func (b *inboundHostBridge) GetCertificate(ctx context.Context, refType cert.CertificateReferenceType, refID string) (
	*inboundmodel.Certificate, *inboundclient.CertOperationError,
) {
	underlying := dcrHostServiceError("inbound host bridge: GetCertificate is not supported; provide Thunder inbound for admin APIs")
	return nil, &inboundclient.CertOperationError{
		Operation:  inboundclient.CertOpRetrieve,
		RefType:    refType,
		Underlying: underlying,
	}
}
