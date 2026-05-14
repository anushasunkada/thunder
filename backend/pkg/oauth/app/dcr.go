// Package app holds portable data shapes for OAuth host integration (DCR).
package app

// Certificate is a PEM/JWK/JWKS URI material reference used when registering an application.
type Certificate struct {
	Type  string
	Value string
}

// AccessTokenConfig mirrors access token settings used during DCR registration.
type AccessTokenConfig struct {
	ValidityPeriod int64
	UserAttributes []string
}

// IDTokenConfig mirrors ID token settings used during DCR registration.
type IDTokenConfig struct {
	ValidityPeriod int64
	UserAttributes []string
	ResponseType   string
	EncryptionAlg  string
	EncryptionEnc  string
}

// OAuthTokenConfig groups token-related registration fields.
type OAuthTokenConfig struct {
	AccessToken *AccessTokenConfig
	IDToken     *IDTokenConfig
}

// UserInfoConfig mirrors userinfo endpoint registration fields.
type UserInfoConfig struct {
	ResponseType   string
	UserAttributes []string
	SigningAlg     string
	EncryptionAlg  string
	EncryptionEnc  string
}

// OAuthRegistration is the OAuth client configuration block for DCR-driven application creation.
type OAuthRegistration struct {
	ClientID                           string
	ClientSecret                       string
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

// ApplicationCreate is the portable create payload produced by Thunder's DCR
// layer before delegating to a host DCRApplication implementation.
type ApplicationCreate struct {
	ID          string
	OUID        string
	Name        string
	URL         string
	LogoURL     string
	TosURI      string
	PolicyURI   string
	Contacts    []string
	Certificate *Certificate
	OAuth       OAuthRegistration
}

// ApplicationCreated is returned after a successful DCR create. It mirrors the
// fields Thunder needs to build the DCR HTTP response.
type ApplicationCreated struct {
	ID          string
	OUID        string
	Name        string
	URL         string
	LogoURL     string
	TosURI      string
	PolicyURI   string
	Contacts    []string
	Certificate *Certificate
	OAuth       OAuthRegistration
}
