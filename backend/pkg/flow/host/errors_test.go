/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 */

package host

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrEntityInboundNotFound_wrapped(t *testing.T) {
	w := fmt.Errorf("outer: %w", ErrEntityInboundNotFound)
	if !errors.Is(w, ErrEntityInboundNotFound) {
		t.Fatal("expected errors.Is to match wrapped sentinel")
	}
}
