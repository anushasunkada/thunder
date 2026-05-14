package flowdeps

import "fmt"

// ExecutionDependencies bundles collaborators required to register flow execution HTTP routes.
// ObservabilitySvc may be nil; flow execution will not publish observability events in that case.
type ExecutionDependencies struct {
	FlowMgtService   FlowMgtService
	InboundClient    InboundClient
	EntityProvider   EntityProvider
	ExecutorRegistry ExecutorRegistry
	ObservabilitySvc ObservabilityService
	CryptoSvc        RuntimeCryptoProvider
}

// Validate returns an error if any required dependency is nil.
// ObservabilitySvc is optional and may be nil.
func (d ExecutionDependencies) Validate() error {
	switch {
	case d.FlowMgtService == nil:
		return fmt.Errorf("flow dependency is required: FlowMgtService")
	case d.InboundClient == nil:
		return fmt.Errorf("flow dependency is required: InboundClient")
	case d.EntityProvider == nil:
		return fmt.Errorf("flow dependency is required: EntityProvider")
	case d.ExecutorRegistry == nil:
		return fmt.Errorf("flow dependency is required: ExecutorRegistry")
	case d.CryptoSvc == nil:
		return fmt.Errorf("flow dependency is required: CryptoSvc")
	default:
		return nil
	}
}
