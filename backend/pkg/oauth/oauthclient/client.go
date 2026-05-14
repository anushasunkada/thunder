// Package oauthclient holds the portable OAuth client runtime view used by
// pkg/oauth host contracts.
package oauthclient

// Certificate is optional client or JWKS material carried on the resolved client.
type Certificate struct {
	Type  string
	Value string
}

// AccessTokenConfig mirrors resolved access token policy.
type AccessTokenConfig struct {
	ValidityPeriod int64
	UserAttributes []string
}

// IDTokenConfig mirrors resolved ID token policy.
type IDTokenConfig struct {
	ValidityPeriod int64
	UserAttributes []string
	ResponseType   string
	EncryptionAlg  string
	EncryptionEnc  string
}

// OAuthTokenConfig groups resolved token configuration.
type OAuthTokenConfig struct {
	AccessToken *AccessTokenConfig
	IDToken     *IDTokenConfig
}

// UserInfoConfig mirrors resolved userinfo policy.
type UserInfoConfig struct {
	ResponseType   string
	UserAttributes []string
	SigningAlg     string
	EncryptionAlg  string
	EncryptionEnc  string
}

// Client is the resolved OAuth/OIDC client used by Thunder's OAuth stack.
type Client struct {
	ID                                 string
	OUID                               string
	ClientID                           string
	RedirectURIs                       []string
	GrantTypes                         []string
	ResponseTypes                      []string
	TokenEndpointAuthMethod            string
	PKCERequired                       bool
	PublicClient                       bool
	RequirePushedAuthorizationRequests bool
	Token                              *OAuthTokenConfig
	Scopes                             []string
	UserInfo                           *UserInfoConfig
	ScopeClaims                        map[string][]string
	Certificate                        *Certificate
	AcrValues                          []string
}
