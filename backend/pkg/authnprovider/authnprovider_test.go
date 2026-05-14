package authnprovider

import (
	"context"
	"testing"
)

func TestPublicTypesCompile(t *testing.T) {
	var _ AuthnProvider
	var _ *AuthnMetadata
	var _ *AuthnResult
	var _ *GetAttributesMetadata
	var _ *GetAttributesResult
	var _ *RequestedAttributes
	var _ *AttributesResponse
}

type testProvider struct{}

func (t *testProvider) Authenticate(
	ctx context.Context,
	identifiers, credentials map[string]interface{},
	metadata *AuthnMetadata,
) (*AuthnResult, *ServiceError) {
	return &AuthnResult{}, nil
}

func (t *testProvider) GetAttributes(
	ctx context.Context,
	token string,
	requestedAttributes *RequestedAttributes,
	metadata *GetAttributesMetadata,
) (*GetAttributesResult, *ServiceError) {
	return &GetAttributesResult{}, nil
}
