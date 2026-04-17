package par

import (
	"context"
	"fmt"
	"net/http"

	"thunder/backend/pkg/oidc"
)

type PAR struct {
	store PARStore
}

func NewPAR(store PARStore) *PAR {
	return &PAR{store: store}
}

func (p *PAR) Name() string { return "par" }

func (p *PAR) Order() int { return 1 }

// ResolveAuthorize implements oidc.AuthorizationResolver (Phase 1).
// If request_uri is present, it fetches the stored PAR request and
// populates the AuthorizeRequest fields before any validator runs.
func (p *PAR) ResolveAuthorize(ctx context.Context, req *oidc.AuthorizeRequest) error {
	if req.RequestURI == "" {
		// Not a PAR request; nothing to do.
		return nil
	}

	stored, err := p.store.GetByRequestURI(ctx, req.RequestURI)
	if err != nil {
		// RFC 9126 §4: invalid or expired request_uri MUST return invalid_request.
		return fmt.Errorf("invalid_request: %w", ErrRequestURINotFound)
	}

	// RFC 9126 §4: client_id in the authorize query MUST match the one in the PAR record.
	if req.ClientID != stored.ClientID {
		return fmt.Errorf("invalid_request: %w", ErrClientIDMismatch)
	}

	// Populate the request from stored PAR params.
	// Any authorization parameters sent alongside request_uri are intentionally
	// overwritten here per RFC 9126 §4 ("SHOULD ignore").
	req.Scope = stored.Scope
	req.State = stored.State
	req.CodeChallenge = stored.CodeChallenge
	req.CodeChallengeMethod = stored.CodeChallengeMethod
	req.Nonce = stored.Nonce

	// Signal to downstream hooks that this request was resolved via PAR.
	req.Metadata["par_resolved"] = true
	req.Metadata["par_request_uri"] = req.RequestURI

	// RFC 9126 §2.2: request_uri MUST be one-time-use.
	_ = p.store.Delete(ctx, req.RequestURI)

	return nil
}

func (p *PAR) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("/par", p.handlePAR)
}

func (p *PAR) handlePAR(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(`{"request_uri":"urn:par:123","expires_in":90}`))
}

func (p *PAR) DiscoveryMetadata() map[string]any {
	return map[string]any{
		"pushed_authorization_request_endpoint": "http://localhost:8080/par",
	}
}
