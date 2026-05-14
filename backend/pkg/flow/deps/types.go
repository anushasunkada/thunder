package flowdeps

import (
	"github.com/thunder-id/thunderid/internal/entityprovider"
	"github.com/thunder-id/thunderid/internal/flow/executor"
	flowmgt "github.com/thunder-id/thunderid/internal/flow/mgt"
	"github.com/thunder-id/thunderid/internal/inboundclient"
	"github.com/thunder-id/thunderid/internal/system/cache"
	"github.com/thunder-id/thunderid/internal/system/kmprovider"
	"github.com/thunder-id/thunderid/internal/system/observability"
)

type (
	CacheManager          = cache.CacheManagerInterface
	FlowMgtService        = flowmgt.FlowMgtServiceInterface
	InboundClient         = inboundclient.InboundClientServiceInterface
	EntityProvider        = entityprovider.EntityProviderInterface
	ExecutorRegistry      = executor.ExecutorRegistryInterface
	ObservabilityService  = observability.ObservabilityServiceInterface
	RuntimeCryptoProvider = kmprovider.RuntimeCryptoProvider
)
