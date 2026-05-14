// Package embed wires the same HTTP handlers as the Thunder server binary onto a
// caller-owned [net/http.ServeMux], for embedding Thunder in another process
// (separate Go module) that cannot construct [pkg/oauth/deps.Dependencies] by hand.
package embed

import (
	"net/http"

	"github.com/thunder-id/thunderid/internal/serverbootstrap"
	"github.com/thunder-id/thunderid/internal/system/cache"
	"github.com/thunder-id/thunderid/internal/system/security"
	pkgauthn "github.com/thunder-id/thunderid/pkg/authnprovider"
)

// WireThunder loads Thunder configuration from serverHome, initializes the global
// runtime, and registers OAuth, flow execution, and the rest of Thunder's HTTP
// surface on mux. Register esignet-specific routes on mux before calling this if
// they must win over Thunder's patterns.
//
// customAuthn: when non-nil, replaces Thunder's default authn provider manager for
// OAuth and flow (see [github.com/thunder-id/thunderid/pkg/authnprovider/wire.NewManager]).
// When nil, Thunder's built-in passkey/OTP/federated manager is used.
func WireThunder(mux *http.ServeMux, serverHome string, customAuthn pkgauthn.AuthnProvider) error {
	cfg, err := serverbootstrap.InitializeRuntime(serverHome)
	if err != nil {
		return err
	}
	security.InitSystemPermissions(cfg.Resource.SystemResourceServer.Handle)
	cm := cache.Initialize()
	serverbootstrap.RegisterServices(mux, cm, customAuthn)
	return nil
}

// ShutdownThunder runs the same teardown as the Thunder server binary's graceful
// shutdown path for components registered via [WireThunder] (observability, etc.).
func ShutdownThunder() {
	serverbootstrap.UnregisterServices()
}
