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

	"github.com/thunder-id/thunderid/internal/application"
	"github.com/thunder-id/thunderid/pkg/oauth/app"
	"github.com/thunder-id/thunderid/pkg/oauth/host"
)

// ThunderDCRApplication implements host.DCRApplication by delegating to Thunder's application service.
type ThunderDCRApplication struct {
	ApplicationService application.ApplicationServiceInterface
}

// CreateApplication implements host.DCRApplication.
func (t *ThunderDCRApplication) CreateApplication(ctx context.Context, in *app.ApplicationCreate) (*app.ApplicationCreated, error) {
	if t == nil || t.ApplicationService == nil {
		return nil, fmt.Errorf("application service is required")
	}
	dto, err := applicationDTOFromAppCreate(in)
	if err != nil {
		return nil, err
	}
	created, se := t.ApplicationService.CreateApplication(ctx, dto)
	if se != nil {
		return nil, serviceErrorAsError(se)
	}
	return applicationCreatedFromDTO(created), nil
}

// DeleteApplication implements host.DCRApplication.
func (t *ThunderDCRApplication) DeleteApplication(ctx context.Context, appID string) error {
	if t == nil || t.ApplicationService == nil {
		return fmt.Errorf("application service is required")
	}
	if se := t.ApplicationService.DeleteApplication(ctx, appID); se != nil {
		return serviceErrorAsError(se)
	}
	return nil
}

var _ host.DCRApplication = (*ThunderDCRApplication)(nil)
