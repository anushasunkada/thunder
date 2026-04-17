package oidc

import (
	"net/http"
	"sort"
)

type OIDCServer struct {
	authResolvers []AuthorizationResolver // Phase 1: enrich/transform the request
	authHooks     []AuthorizationHook     // Phase 2: validate/process the request
	tokenHooks    []TokenHook
	userinfoHooks []UserinfoHook
	discovery     []DiscoveryContributor
	metadata      map[string]any
	userProvider  Provider
}

func NewOIDCServer(userbaseProvider Provider, capabilities ...Capability) *OIDCServer {
	s := &OIDCServer{}

	s.userProvider = userbaseProvider

	//Adding hooks
	for _, c := range capabilities {
		if r, ok := c.(AuthorizationResolver); ok {
			s.authResolvers = append(s.authResolvers, r)
		}
		if h, ok := c.(AuthorizationHook); ok {
			s.authHooks = append(s.authHooks, h)
		}
		if h, ok := c.(TokenHook); ok {
			s.tokenHooks = append(s.tokenHooks, h)
		}
		if h, ok := c.(UserinfoHook); ok {
			s.userinfoHooks = append(s.userinfoHooks, h)
		}
		if d, ok := c.(DiscoveryContributor); ok {
			s.discovery = append(s.discovery, d)
		}
	}

	//Order the hooks
	sortHooks(s.authResolvers)
	sortHooks(s.authHooks)
	sortHooks(s.tokenHooks)
	sortHooks(s.userinfoHooks)

	return s
}

func (s *OIDCServer) RegisterCapabilities(mux *http.ServeMux, caps ...Capability) {
	for _, c := range caps {
		if e, ok := c.(EndpointRegistrar); ok {
			e.RegisterEndpoints(mux)
		}
	}
}

func sortHooks[T any](hooks []T) {
	//SliceStable - if the hooks are with same order -> Registration order is preserved
	sort.SliceStable(hooks, func(i, j int) bool {
		hi, iok := any(hooks[i]).(Ordered)
		hj, jok := any(hooks[j]).(Ordered)

		if !iok && !jok {
			return false
		}
		if !iok {
			return false
		}
		if !jok {
			return true
		}
		return hi.Order() < hj.Order()
	})
}
