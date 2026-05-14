/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 */

package hostbridge

import "testing"

func TestNewThunderInboundFlow_nil(t *testing.T) {
	if NewThunderInboundFlow(nil) != nil {
		t.Fatal("expected nil")
	}
}
