package wire

import (
	mgr "github.com/thunder-id/thunderid/internal/authnprovider/manager"
	pkgauthn "github.com/thunder-id/thunderid/pkg/authnprovider"
)

// NewManager wraps an AuthnProvider for use with Thunder OAuth and flow layers.
func NewManager(p pkgauthn.AuthnProvider) mgr.AuthnProviderManagerInterface {
	return mgr.InitializeAuthnProviderManagerWithProvider(p)
}
