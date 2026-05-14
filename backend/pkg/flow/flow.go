package flow

import (
	"net/http"

	"github.com/thunder-id/thunderid/internal/flow/core"
	flowexec "github.com/thunder-id/thunderid/internal/flow/flowexec"
	flowdeps "github.com/thunder-id/thunderid/pkg/flow/deps"
)

type (
	FlowFactory           = core.FlowFactoryInterface
	GraphCache            = core.GraphCacheInterface
	FlowExecService       = flowexec.FlowExecServiceInterface
	FlowMgtService        = flowdeps.FlowMgtService
	InboundClient         = flowdeps.InboundClient
	EntityProvider        = flowdeps.EntityProvider
	ExecutorRegistry      = flowdeps.ExecutorRegistry
	ObservabilityService  = flowdeps.ObservabilityService
	RuntimeCryptoProvider = flowdeps.RuntimeCryptoProvider
	CacheManager          = flowdeps.CacheManager
)

type ExecutionDependencies = flowdeps.ExecutionDependencies

func InitializeCore(cacheManager flowdeps.CacheManager) (FlowFactory, GraphCache) {
	return core.Initialize(cacheManager)
}

// InitializeExecutionWithDependencies wires flow execution using a single dependency object.
func InitializeExecutionWithDependencies(mux *http.ServeMux, deps ExecutionDependencies) (FlowExecService, error) {
	if err := deps.Validate(); err != nil {
		return nil, err
	}
	return flowexec.Initialize(
		mux,
		deps.FlowMgtService,
		deps.InboundClient,
		deps.EntityProvider,
		deps.ExecutorRegistry,
		deps.ObservabilitySvc,
		deps.CryptoSvc,
	)
}

func InitializeExecution(
	mux *http.ServeMux,
	flowMgtService flowdeps.FlowMgtService,
	inboundClientService flowdeps.InboundClient,
	entityProvider flowdeps.EntityProvider,
	executorRegistry flowdeps.ExecutorRegistry,
	observabilitySvc flowdeps.ObservabilityService,
	cryptoSvc flowdeps.RuntimeCryptoProvider,
) (FlowExecService, error) {
	return InitializeExecutionWithDependencies(mux, ExecutionDependencies{
		FlowMgtService:   flowMgtService,
		InboundClient:    inboundClientService,
		EntityProvider:   entityProvider,
		ExecutorRegistry: executorRegistry,
		ObservabilitySvc: observabilitySvc,
		CryptoSvc:        cryptoSvc,
	})
}
