package oidc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *OIDCServer) Discovery(w http.ResponseWriter, _ *http.Request) {
	meta := map[string]any{
		"issuer":         "http://localhost:8080",
		"authorize_url":  "http://localhost:8080/authorize",
		"token_endpoint": "http://localhost:8080/token",
	}

	//Execute if there are any Discovery hooks that want to modify the discovery metadata before returning
	//Each hook can add its own metadata to the discovery response by returning a map[string]any from the DiscoveryMetadata method of the hook interface
	for _, d := range s.discovery {
		for k, v := range d.DiscoveryMetadata() {
			meta[k] = v
		}
	}

	fmt.Println("Handler : Discovery")
	json.NewEncoder(w).Encode(meta)
}
