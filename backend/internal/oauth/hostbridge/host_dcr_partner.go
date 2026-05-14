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

	applicationmodel "github.com/thunder-id/thunderid/internal/application/model"
	"github.com/thunder-id/thunderid/internal/system/error/serviceerror"
	"github.com/thunder-id/thunderid/internal/system/i18n/core"
	"github.com/thunder-id/thunderid/pkg/oauth/host"
)

type hostDCRPartner struct {
	host host.DCRApplication
}

func (h *hostDCRPartner) CreateApplication(
	ctx context.Context, app *applicationmodel.ApplicationDTO,
) (*applicationmodel.ApplicationDTO, *serviceerror.ServiceError) {
	if h == nil || h.host == nil {
		return nil, dcrHostServiceError("application host is not configured")
	}
	create := applicationCreateFromDTO(app)
	if create == nil {
		return nil, dcrHostServiceError("invalid application payload for DCR")
	}
	created, err := h.host.CreateApplication(ctx, create)
	if err != nil {
		return nil, dcrHostServiceError(err.Error())
	}
	dto, convErr := applicationCreatedToDTO(created)
	if convErr != nil {
		return nil, dcrHostServiceError(convErr.Error())
	}
	return dto, nil
}

func (h *hostDCRPartner) DeleteApplication(ctx context.Context, appID string) *serviceerror.ServiceError {
	if h == nil || h.host == nil {
		return dcrHostServiceError("application host is not configured")
	}
	if err := h.host.DeleteApplication(ctx, appID); err != nil {
		return dcrHostServiceError(err.Error())
	}
	return nil
}

func dcrHostServiceError(msg string) *serviceerror.ServiceError {
	return serviceerror.CustomServiceError(serviceerror.InternalServerError, core.I18nMessage{
		Key:          "error.oauth.host_application",
		DefaultValue: msg,
	})
}
