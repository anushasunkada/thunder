# ThunderID Engine (`pkg/thunderidengine`)

Embeddable OIDC authorization server: flow execution, OAuth2/OIDC endpoints (excluding DCR), and declarative configuration loading.

## Imports

| Package | Use for |
|---------|---------|
| `github.com/thunder-id/thunderid/pkg/thunderidengine` | `Initialize`, `EngineConfig`, `Engine`, `FlowExec` |
| `github.com/thunder-id/thunderid/pkg/thunderidengine/runtime` | `Store`, `NewMemoryRuntimeStore`, `ErrNotFound` |
| `github.com/thunder-id/thunderid/pkg/thunderidengine/host` | `ActorProvider`, `AuthnProvider`, `AuthorizationProvider`, `ConsentEnforcer`, `FlowProvider` |

## Quick start

```go
import (
    "net/http"

    "github.com/thunder-id/thunderid/pkg/thunderidengine"
    "github.com/thunder-id/thunderid/pkg/thunderidengine/host"
    "github.com/thunder-id/thunderid/pkg/thunderidengine/runtime"
)

mux := http.NewServeMux()
_, err := thunderidengine.Initialize(mux, thunderidengine.EngineConfig{
    Issuer:  "https://as.example.com",
    DataDir: "/path/to/data",
    Runtime: runtime.NewMemoryRuntimeStore(),
    Crypto:  thunderidengine.CryptoConfig{SigningKeyPath: "/path/to/signing.key"},
    FlowStore: thunderidengine.FlowProviderConfig{
        StoreMode: thunderidengine.StoreModeDeclarative,
    },
    Actors:        myActors,        // host.ActorProvider
    Authn:         myAuthn,         // host.AuthnProvider
    Authorization: myAuthz,         // host.AuthorizationProvider
    Consent:         myConsent,       // host.ConsentEnforcer
})
```

## Required configuration

- **Issuer** — OIDC issuer URL
- **DataDir** — server home; declarative YAML is read from `DataDir/repository/resources/`
- **Runtime** — `runtime.Store` (in-memory for dev; ThunderID server uses Redis/SQL via host adapters)
- **Crypto.SigningKeyPath** — PEM private key for JWT signing
- **Actors, Authn, Authorization, Consent** — host implementations

Optional:

- **FlowProvider** — custom flow source; if nil, **FlowStore** loads flows from declarative files under DataDir
- **Flow.Executors** — built-in executor names to register; when empty, defaults to `BasicAuthExecutor`, `AuthAssertExecutor`, `ConsentExecutor`

## Data directory layout

```text
data/
  repository/resources/   OU, IDP, themes, layouts, roles, translations, resource servers, …
  flows/                  optional; use FlowStore.DefinitionsPath when needed
```

## Runtime storage

`runtime.Store` holds ephemeral state: flow contexts, authorization codes, auth requests, PAR, JTI replay cache, and attribute cache (session claims via JWT `aci`).

For production on ThunderID server, use Redis or SQL runtime via `internal/hostadapters/runtime` (returns `runtime.Store`).

## Further reading

See repository [README.md](../../../README.md) and [ARCHITECTURE.md](../../../ARCHITECTURE.md) for the full product.
