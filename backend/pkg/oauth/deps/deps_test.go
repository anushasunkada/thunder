package oauthdeps

import "testing"

func TestDependenciesValidate(t *testing.T) {
	if err := (Dependencies{}).Validate(); err == nil {
		t.Fatal("expected dependency validation error")
	}
}
