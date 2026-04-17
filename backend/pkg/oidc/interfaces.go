package oidc

import (
	"context"
	"net/http"
)

type AuthorizeRequest struct {
	ClientID            string
	Scope               string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
	RequestURI          string
	Nonce               string
	Metadata            map[string]any // for inter-hook communication; initialized by the server before hooks run
}

type TokenRequest struct {
	ClientID            string
	CodeVerifier        string
	Code                string
	ClientAssertion     string
	ClientAssertionType string
}

type User struct {
	ID       string
	Username string
	Email    string
	Claims   map[string]any
}

type Capability interface {
	Name() string
}

type Ordered interface {
	Order() int
}

// AuthorizationResolver is implemented by capabilities that need to enrich or
// transform an AuthorizeRequest before validation hooks run (Phase 1).
// Example: PAR resolves a request_uri to its stored parameters.
type AuthorizationResolver interface {
	ResolveAuthorize(ctx context.Context, req *AuthorizeRequest) error
}

// AuthorizationHook is implemented by capabilities that validate or process
// a fully-resolved AuthorizeRequest (Phase 2).
type AuthorizationHook interface {
	BeforeAuthorize(ctx context.Context, req *AuthorizeRequest) error
}

type TokenHook interface {
	BeforeToken(ctx context.Context, req *TokenRequest) error
}

type UserinfoHook interface {
	AfterUserinfo(ctx context.Context, res string) error
}

type EndpointRegistrar interface {
	RegisterEndpoints(mux *http.ServeMux)
}

type DiscoveryContributor interface {
	DiscoveryMetadata() map[string]any
}

type Provider interface {
	Name() string
	FetchUser(ctx context.Context, identifier string) (*User, error)
}
