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

	"github.com/thunder-id/thunderid/internal/application"
	applicationmodel "github.com/thunder-id/thunderid/internal/application/model"
	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
)

type thunderDCRPartner struct {
	inner application.ApplicationServiceInterface
}

func (t *thunderDCRPartner) CreateApplication(
	ctx context.Context, app *applicationmodel.ApplicationDTO,
) (*applicationmodel.ApplicationDTO, *serviceerror.ServiceError) {
	if t == nil || t.inner == nil {
		return nil, dcrHostServiceError("application service is not configured")
	}
	return t.inner.CreateApplication(ctx, app)
}

func (t *thunderDCRPartner) DeleteApplication(ctx context.Context, appID string) *serviceerror.ServiceError {
	if t == nil || t.inner == nil {
		return dcrHostServiceError("application service is not configured")
	}
	return t.inner.DeleteApplication(ctx, appID)
}
