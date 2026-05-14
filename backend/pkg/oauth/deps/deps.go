package oauthdeps

import (
	"fmt"
	"strings"

	"github.com/thunder-id/thunderid/internal/system/database/provider"
)

// Dependencies bundles collaborators required to register OAuth HTTP routes.
type Dependencies struct {
	Application         DCRApplication
	Inbound             InboundOAuth
	AuthnProvider       AuthnProviderManager
	JWTService          JWTService
	JWEService          JWEService
	FlowExecService     FlowExecService
	ObservabilitySvc    ObservabilityService
	PKIService          PKIService
	OUService           OUService
	AttributeCacheSvc   AttributeCacheService
	AuthzService        AuthorizationService
	EntityProvider      EntityProvider
	ResourceService     ResourceService
	I18nService         I18nService
	IDPService          IDPService
	Transactioner       Transactioner
	DBProvider          DBProvider
	RedisProvider       RedisProvider
	DeploymentID        string
	DatabaseRuntimeType string
}

// Validate returns an error if any required dependency is nil.
func (d Dependencies) Validate() error {
	switch {
	case d.Application == nil:
		return fmt.Errorf("oauth dependency is required: Application")
	case d.Inbound == nil:
		return fmt.Errorf("oauth dependency is required: Inbound")
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
	case d.Transactioner == nil:
		return fmt.Errorf("oauth dependency is required: Transactioner")
	case d.DBProvider == nil:
		return fmt.Errorf("oauth dependency is required: DBProvider")
	case strings.TrimSpace(d.DeploymentID) == "":
		return fmt.Errorf("oauth dependency is required: DeploymentID")
	case strings.TrimSpace(d.DatabaseRuntimeType) == "":
		return fmt.Errorf("oauth dependency is required: DatabaseRuntimeType")
	case d.DatabaseRuntimeType == provider.DataSourceTypeRedis && d.RedisProvider == nil:
		return fmt.Errorf("oauth dependency is required: RedisProvider when DatabaseRuntimeType is redis")
	default:
		return nil
	}
}
