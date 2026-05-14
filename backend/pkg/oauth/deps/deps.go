package oauthdeps

import "fmt"

// Dependencies bundles collaborators required to register OAuth HTTP routes.
type Dependencies struct {
	ApplicationService ApplicationService
	InboundClient      InboundClient
	AuthnProvider      AuthnProviderManager
	JWTService         JWTService
	JWEService         JWEService
	FlowExecService    FlowExecService
	ObservabilitySvc   ObservabilityService
	PKIService         PKIService
	OUService          OUService
	AttributeCacheSvc  AttributeCacheService
	AuthzService       AuthorizationService
	EntityProvider     EntityProvider
	ResourceService    ResourceService
	I18nService        I18nService
	IDPService         IDPService
}

// Validate returns an error if any required dependency is nil.
func (d Dependencies) Validate() error {
	switch {
	case d.ApplicationService == nil:
		return fmt.Errorf("oauth dependency is required: ApplicationService")
	case d.InboundClient == nil:
		return fmt.Errorf("oauth dependency is required: InboundClient")
	case d.AuthnProvider == nil:
		return fmt.Errorf("oauth dependency is required: AuthnProvider")
	case d.JWTService == nil:
		return fmt.Errorf("oauth dependency is required: JWTService")
	case d.JWEService == nil:
		return fmt.Errorf("oauth dependency is required: JWEService")
	case d.FlowExecService == nil:
		return fmt.Errorf("oauth dependency is required: FlowExecService")
	case d.ObservabilitySvc == nil:
		return fmt.Errorf("oauth dependency is required: ObservabilitySvc")
	case d.PKIService == nil:
		return fmt.Errorf("oauth dependency is required: PKIService")
	case d.OUService == nil:
		return fmt.Errorf("oauth dependency is required: OUService")
	case d.AttributeCacheSvc == nil:
		return fmt.Errorf("oauth dependency is required: AttributeCacheSvc")
	case d.AuthzService == nil:
		return fmt.Errorf("oauth dependency is required: AuthzService")
	case d.EntityProvider == nil:
		return fmt.Errorf("oauth dependency is required: EntityProvider")
	case d.ResourceService == nil:
		return fmt.Errorf("oauth dependency is required: ResourceService")
	case d.I18nService == nil:
		return fmt.Errorf("oauth dependency is required: I18nService")
	case d.IDPService == nil:
		return fmt.Errorf("oauth dependency is required: IDPService")
	default:
		return nil
	}
}
