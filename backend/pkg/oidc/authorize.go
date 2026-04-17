package oidc

import (
	"net/http"
)

func (s *OIDCServer) Authorize(w http.ResponseWriter, r *http.Request) {
	req := &AuthorizeRequest{
		ClientID:            r.URL.Query().Get("client_id"),
		Scope:               r.URL.Query().Get("scope"),
		State:               r.URL.Query().Get("state"),
		CodeChallenge:       r.URL.Query().Get("code_challenge"),
		CodeChallengeMethod: r.URL.Query().Get("code_challenge_method"),
		RequestURI:          r.URL.Query().Get("request_uri"),
		Nonce:               r.URL.Query().Get("nonce"),
		Metadata:            make(map[string]any),
	}

	// Phase 1: Resolvers enrich/transform the request (e.g. PAR resolves request_uri).
	// All resolvers complete before any validator runs.
	for _, resolver := range s.authResolvers {
		if err := resolver.ResolveAuthorize(r.Context(), req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Phase 2: Validators/processors operate on the fully-resolved request.
	for _, h := range s.authHooks {
		if err := h.BeforeAuthorize(r.Context(), req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	//TODO : validate client_id, scope, code_challenge_method, etc.
	//TODO : generate and store flow state
	//TODO : redirect to gate application with flow state
	w.Write([]byte("authorization successful"))
}
