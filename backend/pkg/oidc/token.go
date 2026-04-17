package oidc

import (
	"fmt"
	"net/http"
)

func (s *OIDCServer) Token(w http.ResponseWriter, r *http.Request) {
	req := &TokenRequest{
		ClientID:     r.FormValue("client_id"),
		CodeVerifier: r.FormValue("code_verifier"),
		Code:         r.FormValue("code"),
	}

	// Execute token hooks (e.g. for client authentication, PKCE verification, etc.)
	for _, h := range s.tokenHooks {
		if err := h.BeforeToken(r.Context(), req); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	}

	//TODO : validate client_id, code_verifier, code, etc.
	//TODO : generate and return ID token, access token, etc.
	fmt.Println("Handler : Token")
	w.Write([]byte("token issued"))
}
