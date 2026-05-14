/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 */

package host

import "errors"

// ErrEntityInboundNotFound is returned when no inbound profile exists for the entity ID.
var ErrEntityInboundNotFound = errors.New("flow host: inbound profile not found for entity")
