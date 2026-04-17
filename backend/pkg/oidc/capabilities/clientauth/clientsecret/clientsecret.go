package clientsecret

import (
	"context"
	"fmt"

	"github.com/asgardeo/thunder/pkg/oidc"
)

type ClientSecret struct{}

func (c *ClientSecret) Name() string { return "client-secret-auth" }

func (c *ClientSecret) Order() int {
	return 1
}

func (c *ClientSecret) BeforeToken(_ context.Context, r *oidc.TokenRequest) error {
	fmt.Println("BeforeToken Hook : ClientSecret")
	return nil
}

func (p *ClientSecret) DiscoveryMetadata() map[string]any {
	fmt.Println("DiscoveryMetadata Hook : ClientSecret")
	return map[string]any{
		"token_endpoint_auth_methods_supported": []string{"client_secret_basic", "client_secret_post"},
	}
}
