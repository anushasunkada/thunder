package par

import "context"

// StoredRequest holds the authorization parameters pushed via the PAR endpoint.
type StoredRequest struct {
	ClientID            string
	Scope               string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
	Nonce               string
	ExpiresAt           int64 // Unix timestamp; enforced by PARStore implementations
}

// PARStore is the storage contract for pushed authorization requests.
type PARStore interface {
	// Store saves a pushed authorization request under the given request_uri.
	// Implementations use StoredRequest.ExpiresAt to set an appropriate TTL;
	// a zero ExpiresAt falls back to the provider's configured default TTL.
	Store(ctx context.Context, requestURI string, req *StoredRequest) error

	// GetByRequestURI retrieves a stored request by its request_uri.
	// Returns an error if the URI is unknown or expired.
	GetByRequestURI(ctx context.Context, requestURI string) (*StoredRequest, error)

	// Delete removes a stored request. Called after successful retrieval
	// to enforce the one-time-use requirement of RFC 9126 §2.2.
	Delete(ctx context.Context, requestURI string) error
}
