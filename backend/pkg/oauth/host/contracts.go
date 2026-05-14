// Package host defines OAuth host-facing contracts that can be implemented by
// applications outside Thunder's internal/ tree (separate Go modules).
package host

import (
	"context"

	"github.com/thunder-id/thunderid/pkg/oauth/app"
	"github.com/thunder-id/thunderid/pkg/oauth/oauthclient"
)

// DCRApplication is the application registry surface required for OAuth Dynamic
// Client Registration (DCR). Implementations map Thunder's DCR payloads to the
// host's persistence model.
type DCRApplication interface {
	CreateApplication(ctx context.Context, create *app.ApplicationCreate) (*app.ApplicationCreated, error)
	DeleteApplication(ctx context.Context, appID string) error
}

// InboundOAuth is the inbound surface required for OAuth runtime (client lookup
// by public client_id). Hosts that only override OAuth client resolution should
// implement this interface; Thunder's default wiring adapts the full inbound
// service via NewThunderInbound in internal/oauth/hostbridge.
type InboundOAuth interface {
	GetOAuthClientByClientID(ctx context.Context, clientID string) (*oauthclient.Client, error)
}
