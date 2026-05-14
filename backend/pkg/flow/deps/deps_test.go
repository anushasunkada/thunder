package flowdeps

import "testing"

func TestExecutionDependenciesValidate_empty(t *testing.T) {
	if err := (ExecutionDependencies{}).Validate(); err == nil {
		t.Fatal("expected dependency validation error")
	}
}
