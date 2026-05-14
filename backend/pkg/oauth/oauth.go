package oauth

import (
	"fmt"
	"net/http"

	internaloauth "github.com/thunder-id/thunderid/internal/oauth"
	oauthdeps "github.com/thunder-id/thunderid/pkg/oauth/deps"
)

type (
	DCRApplication        = oauthdeps.DCRApplication
	InboundOAuth          = oauthdeps.InboundOAuth
	AuthnProviderManager  = oauthdeps.AuthnProviderManager
	JWTService            = oauthdeps.JWTService
	JWEService            = oauthdeps.JWEService
	FlowExecService       = oauthdeps.FlowExecService
	ObservabilityService  = oauthdeps.ObservabilityService
	PKIService            = oauthdeps.PKIService
	OUService             = oauthdeps.OUService
	AttributeCacheService = oauthdeps.AttributeCacheService
	AuthorizationService  = oauthdeps.AuthorizationService
	EntityProvider        = oauthdeps.EntityProvider
	ResourceService       = oauthdeps.ResourceService
	I18nService           = oauthdeps.I18nService
	IDPService            = oauthdeps.IDPService
	Transactioner         = oauthdeps.Transactioner
	DBProvider            = oauthdeps.DBProvider
	RedisProvider         = oauthdeps.RedisProvider
)

type Dependencies = oauthdeps.Dependencies

// RegisterRoutes is the high-level entrypoint to register OAuth routes for callers
// that already have all OAuth dependencies prepared.
func RegisterRoutes(mux *http.ServeMux, deps Dependencies) error {
	if mux == nil {
		return fmt.Errorf("oauth dependency is required: mux")
	}
	return InitializeWithDependencies(mux, deps)
}

// InitializeWithDependencies wires OAuth using a single dependency object.
func InitializeWithDependencies(mux *http.ServeMux, deps Dependencies) error {
	if err := deps.Validate(); err != nil {
		return err
	}
	return internaloauth.Initialize(mux, deps)
}
