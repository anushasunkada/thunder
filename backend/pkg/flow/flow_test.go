package flow

import "testing"

func TestPublicFlowAPIsAreExposed(t *testing.T) {
	var _ = InitializeCore
	var _ = InitializeExecution
	var _ = InitializeExecutionWithDependencies
	var _ FlowFactory
	var _ GraphCache
	var _ FlowExecService
	var _ FlowMgtService
	var _ ExecutorRegistry
	var _ ExecutionDependencies
}
