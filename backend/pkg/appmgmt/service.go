/*
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package appmgmt provides portable, dependency-free types and interfaces for
// managing OAuth/OIDC applications. It is intentionally kept free of internal
// server dependencies so that it can be consumed by any layer of the stack
// (e.g. oauth, consent, flow) without creating import cycles.
package appmgmt

import "context"

// ApplicationServiceInterface defines the read/write contract for application management.
// Implementations live in internal/application; this interface lives here so that
// external packages can depend on it without depending on internal/.
type ApplicationServiceInterface interface {
	// CreateApplication registers a new application and returns the persisted DTO.
	// Returns ErrConflict if an application with the same name or client-ID already exists.
	CreateApplication(ctx context.Context, app *ApplicationDTO) (*ApplicationDTO, error)

	// GetApplication retrieves a single application by its unique ID.
	// Returns ErrNotFound if no matching application exists.
	GetApplication(ctx context.Context, appID string) (*Application, error)

	// GetApplicationList returns a paginated list of all registered applications.
	GetApplicationList(ctx context.Context) (*ApplicationListResponse, error)

	// GetOAuthApplication returns the processed OAuth configuration for the application
	// identified by clientID.
	// Returns ErrNotFound if no matching application exists.
	GetOAuthApplication(ctx context.Context, clientID string) (*OAuthAppConfigProcessedDTO, error)

	// UpdateApplication fully replaces the application identified by appID with the
	// provided DTO and returns the updated DTO.
	// Returns ErrNotFound if no matching application exists.
	// Returns ErrConflict if the update would violate a uniqueness constraint.
	UpdateApplication(ctx context.Context, appID string, app *ApplicationDTO) (*ApplicationDTO, error)

	// DeleteApplication permanently removes the application identified by appID.
	// Returns ErrNotFound if no matching application exists.
	DeleteApplication(ctx context.Context, appID string) error

	// ValidateApplication validates the provided ApplicationDTO and returns its processed
	// form together with the resolved inbound-auth config.
	// Returns ErrInvalidInput with a descriptive message when validation fails.
	ValidateApplication(
		ctx context.Context,
		app *ApplicationDTO,
	) (*ApplicationProcessedDTO, *InboundAuthConfigProcessedDTO, error)
}
