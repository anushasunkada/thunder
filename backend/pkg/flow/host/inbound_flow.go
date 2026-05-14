/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 */

package host

import "context"

// InboundFlow resolves the inbound profile for an entity (application, agent, etc.)
// for flow graph selection and execution context. Implementations may be backed by
// Thunder's inbound service or a host-specific registry.
type InboundFlow interface {
	GetInboundClientByEntityID(ctx context.Context, entityID string) (*EntityInboundProfile, error)
}
