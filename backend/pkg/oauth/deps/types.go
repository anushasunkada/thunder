package oauthdeps

import (
	"github.com/thunder-id/thunderid/internal/application"
	"github.com/thunder-id/thunderid/internal/attributecache"
	authnprovidermgr "github.com/thunder-id/thunderid/internal/authnprovider/manager"
	"github.com/thunder-id/thunderid/internal/authz"
	"github.com/thunder-id/thunderid/internal/entityprovider"
	"github.com/thunder-id/thunderid/internal/flow/flowexec"
	"github.com/thunder-id/thunderid/internal/idp"
	"github.com/thunder-id/thunderid/internal/inboundclient"
	"github.com/thunder-id/thunderid/internal/ou"
	"github.com/thunder-id/thunderid/internal/resource"
	i18nmgt "github.com/thunder-id/thunderid/internal/system/i18n/mgt"
	"github.com/thunder-id/thunderid/internal/system/jose/jwe"
	"github.com/thunder-id/thunderid/internal/system/jose/jwt"
	"github.com/thunder-id/thunderid/internal/system/kmprovider/defaultkm/pkiservice"
	"github.com/thunder-id/thunderid/internal/system/observability"
)

type (
	ApplicationService    = application.ApplicationServiceInterface
	InboundClient         = inboundclient.InboundClientServiceInterface
	AuthnProviderManager  = authnprovidermgr.AuthnProviderManagerInterface
	JWTService            = jwt.JWTServiceInterface
	JWEService            = jwe.JWEServiceInterface
	FlowExecService       = flowexec.FlowExecServiceInterface
	ObservabilityService  = observability.ObservabilityServiceInterface
	PKIService            = pkiservice.PKIServiceInterface
	OUService             = ou.OrganizationUnitServiceInterface
	AttributeCacheService = attributecache.AttributeCacheServiceInterface
	AuthorizationService  = authz.AuthorizationServiceInterface
	EntityProvider        = entityprovider.EntityProviderInterface
	ResourceService       = resource.ResourceServiceInterface
	I18nService           = i18nmgt.I18nServiceInterface
	IDPService            = idp.IDPServiceInterface
)
