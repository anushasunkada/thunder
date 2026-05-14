package oauth

import (
	"fmt"
	"net/http"

	internaloauth "github.com/thunder-id/thunderid/internal/oauth"
	oauthdeps "github.com/thunder-id/thunderid/pkg/oauth/deps"
)

type (
	ApplicationService    = oauthdeps.ApplicationService
	InboundClientService  = oauthdeps.InboundClient
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

// InitializeWithDependencies is the low-level composable API that accepts a single
// dependency object and initializes all OAuth routes and services.
func InitializeWithDependencies(mux *http.ServeMux, deps Dependencies) error {
	if err := deps.Validate(); err != nil {
		return err
	}
	return internaloauth.Initialize(
		mux,
		deps.ApplicationService,
		deps.InboundClient,
		deps.AuthnProvider,
		deps.JWTService,
		deps.JWEService,
		deps.FlowExecService,
		deps.ObservabilitySvc,
		deps.PKIService,
		deps.OUService,
		deps.AttributeCacheSvc,
		deps.AuthzService,
		deps.EntityProvider,
		deps.ResourceService,
		deps.I18nService,
		deps.IDPService,
	)
}

func Initialize(
	mux *http.ServeMux,
	applicationService oauthdeps.ApplicationService,
	inboundClient oauthdeps.InboundClient,
	authnProvider oauthdeps.AuthnProviderManager,
	jwtService oauthdeps.JWTService,
	jweService oauthdeps.JWEService,
	flowExecService oauthdeps.FlowExecService,
	observabilitySvc oauthdeps.ObservabilityService,
	pkiService oauthdeps.PKIService,
	ouService oauthdeps.OUService,
	attributeCacheSvc oauthdeps.AttributeCacheService,
	authzService oauthdeps.AuthorizationService,
	entityProvider oauthdeps.EntityProvider,
	resourceService oauthdeps.ResourceService,
	i18nService oauthdeps.I18nService,
	idpService oauthdeps.IDPService,
) error {
	return InitializeWithDependencies(mux, Dependencies{
		ApplicationService: applicationService,
		InboundClient:      inboundClient,
		AuthnProvider:      authnProvider,
		JWTService:         jwtService,
		JWEService:         jweService,
		FlowExecService:    flowExecService,
		ObservabilitySvc:   observabilitySvc,
		PKIService:         pkiService,
		OUService:          ouService,
		AttributeCacheSvc:  attributeCacheSvc,
		AuthzService:       authzService,
		EntityProvider:     entityProvider,
		ResourceService:    resourceService,
		I18nService:        i18nService,
		IDPService:         idpService,
	})
}
