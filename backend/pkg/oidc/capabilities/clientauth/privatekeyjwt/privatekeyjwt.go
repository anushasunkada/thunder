package privatekeyjwt

import (
	"context"
	"fmt"

	"github.com/asgardeo/thunder/pkg/oidc"
)

type PrivateKeyJWT struct{}

func (c *PrivateKeyJWT) Name() string { return "private-key-jwt" }

func (c *PrivateKeyJWT) Order() int {
	return 1
}

func (c *PrivateKeyJWT) BeforeToken(_ context.Context, r *oidc.TokenRequest) error {
	fmt.Println("BeforeToken Hook : PrivateKeyJWT")
	return nil
}

func (p *PrivateKeyJWT) DiscoveryMetadata() map[string]any {
	fmt.Println("DiscoveryMetadata Hook : PrivateKeyJWT")
	return map[string]any{
		"token_endpoint_auth_methods_supported": []string{"private_key_jwt"},
	}
}
