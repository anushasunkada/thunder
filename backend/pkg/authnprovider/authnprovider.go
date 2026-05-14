package authnprovider

import "context"

// AuthnProvider is the contract for authentication providers (public API).
type AuthnProvider interface {
	Authenticate(
		ctx context.Context,
		identifiers, credentials map[string]interface{},
		metadata *AuthnMetadata,
	) (*AuthnResult, *ServiceError)
	GetAttributes(
		ctx context.Context,
		token string,
		requestedAttributes *RequestedAttributes,
		metadata *GetAttributesMetadata,
	) (*GetAttributesResult, *ServiceError)
}
