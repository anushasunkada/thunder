package wire_test

import (
	"context"
	"testing"

	pkgauthn "github.com/thunder-id/thunderid/pkg/authnprovider"
	"github.com/thunder-id/thunderid/pkg/authnprovider/wire"
)

type testProvider struct{}

func (t *testProvider) Authenticate(
	ctx context.Context,
	identifiers, credentials map[string]interface{},
	metadata *pkgauthn.AuthnMetadata,
) (*pkgauthn.AuthnResult, *pkgauthn.ServiceError) {
	return &pkgauthn.AuthnResult{}, nil
}

func (t *testProvider) GetAttributes(
	ctx context.Context,
	token string,
	requestedAttributes *pkgauthn.RequestedAttributes,
	metadata *pkgauthn.GetAttributesMetadata,
) (*pkgauthn.GetAttributesResult, *pkgauthn.ServiceError) {
	return &pkgauthn.GetAttributesResult{}, nil
}

func TestNewManager(t *testing.T) {
	manager := wire.NewManager(&testProvider{})
	if manager == nil {
		t.Fatal("expected manager to be initialized")
	}
}
