/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 */

// Package hostbridge adapts pkg/oauth host contracts to Thunder internal OAuth wiring.
package hostbridge

import (
	"context"

	applicationmodel "github.com/thunder-id/thunderid/internal/application/model"
	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
)

// DCRPartner is the minimal application capability required by the DCR package.
type DCRPartner interface {
	CreateApplication(
		ctx context.Context, app *applicationmodel.ApplicationDTO,
	) (*applicationmodel.ApplicationDTO, *serviceerror.ServiceError)
	DeleteApplication(ctx context.Context, appID string) *serviceerror.ServiceError
}
