# pkg/oidc

A framework-agnostic, extensible OIDC core library. It implements the minimal OIDC/OAuth2 server
lifecycle and exposes clean hook points so that any additional/extension RFC — PAR, PKCE, DPoP, JAR, etc. — can be
layered on top without touching the core.

---

## Design Philosophy

The core package owns the core OIDC HTTP handlers and the lifecycle skeleton. It only implements
core OIDC behaviour itself. It defines a set of hook interfaces. Capabilities
(separate packages that each implement one RFC or auth method) attach to those hooks and run
in a well-defined order.

This means:

- Adding a new RFC = adding a new package that implements one or more hook interfaces.
- Removing a capability = not registering it — no core changes needed.
- Different deployments (a minimal server vs. a full-featured one) simply register different
  capability sets at wiring time.

---

## Hook Interfaces

```
interfaces.go
```

| Interface | Method | Phase | Purpose |
|---|---|---|---|
| `AuthorizationResolver` | `ResolveAuthorize` | Authorize Phase 1 | Enrich / transform the request before any validation runs |
| `AuthorizationHook` | `BeforeAuthorize` | Authorize Phase 2 | Validate or process the fully-resolved request |
| `TokenHook` | `BeforeToken` | Token | Validate or process the token request |
| `UserinfoHook` | `AfterUserinfo` | Userinfo | Post-process the userinfo response |
| `EndpointRegistrar` | `RegisterEndpoints` | Server init | Mount additional HTTP endpoints (e.g. `/par`) |
| `DiscoveryContributor` | `DiscoveryMetadata` | Discovery | Contribute fields to the OpenID configuration document |

A capability may implement any combination of these interfaces. The server detects which ones
apply via Go type assertions at construction time.

### Execution ordering

Capabilities that implement `Ordered` (`Order() int`) are sorted within their phase before any
request is processed. Lower values run first. Capabilities without `Ordered` sort after all
ordered ones, preserving registration order among themselves.

---

## Two-Phase Authorize Flow

The authorize endpoint uses two distinct phases. This is the key architectural decision that
makes RFC interactions safe.

```
HTTP GET /authorize
        │
        ▼
┌───────────────────────────────┐
│  Phase 1 — Resolve            │  All AuthorizationResolver hooks run in order.
│                               │  Each hook may enrich / transform AuthorizeRequest.
│  e.g. PAR resolves            │  By the end of this phase the request is fully populated,
│       request_uri → params    │  regardless of how the client submitted it.
└───────────────────────────────┘
        │
        ▼
┌───────────────────────────────┐
│  Phase 2 — Validate           │  All AuthorizationHook hooks run in order.
│                               │  Each hook validates one concern against the resolved request.
│  e.g. PKCE validates          │  Hooks are guaranteed to see the complete, resolved params.
│       code_challenge          │
└───────────────────────────────┘
        │
        ▼
     Response
```

### Why two phases?

Consider PAR (RFC 9126). A client pushes authorization parameters to `/par` and receives a
`request_uri`. It then calls `/authorize?client_id=x&request_uri=urn:par:abc`. At that point,
`scope`, `code_challenge`, `nonce`, etc. are **absent from the HTTP request** — they live in
the stored PAR record.

If PKCE's validator ran before PAR had a chance to resolve the `request_uri`, it would see an
empty `code_challenge` and fail. Making PAR run earlier via `Order()` alone would work, but
provides no structural guarantee: a future developer could add a validator with `Order() = 0`
and silently break PAR flows.

The two-phase split is a **type-level guarantee**:

- `AuthorizationResolver` implementations **always** run before `AuthorizationHook` implementations.
- No ordering number bridges the two phases.
- Validators never need to know PAR exists; they simply see a fully-populated `AuthorizeRequest`.

The same pattern applies to any future RFC that needs to enrich the request before validation
(e.g. JAR — JWT-Secured Authorization Requests, RFC 9101).

---

## AuthorizeRequest — Metadata field

```go
type AuthorizeRequest struct {
    ClientID, Scope, State, CodeChallenge, CodeChallengeMethod, RequestURI, Nonce string
    Metadata map[string]any  // for inter-hook communication
}
```

`Metadata` is initialised to an empty map by the server before any hook runs. Resolver hooks
can write values here to signal context to downstream hooks without coupling them directly.

Example: PAR sets `Metadata["par_resolved"] = true`. A hypothetical strict-mode validator
hook could read this to adjust its behaviour. Most hooks ignore `Metadata` entirely.

---

## Capability Interface

Every capability must implement `Capability`:

```go
type Capability interface {
    Name() string
}
```

Beyond that, a capability opts into hooks by implementing the relevant interfaces. The server
introspects each capability at construction time using type assertions.

```go
type Ordered interface {
    Order() int  // lower = runs first within its phase
}
```

## PAR capability in detail

PAR participates in **Phase 1** as an `AuthorizationResolver`:

1. If `request_uri` is absent → no-op (plain authorize flows are unaffected).
2. Fetch stored params from `PARStore` by `request_uri`. Error → `invalid_request`.
3. Validate `client_id` matches the stored record (RFC 9126 §4).
4. Populate `AuthorizeRequest` fields from stored params (scope, code_challenge, nonce, etc.).
5. Set `Metadata["par_resolved"] = true`.
6. Delete the `request_uri` from storage — one-time use (RFC 9126 §2.2).

After step 6, Phase 2 hooks (PKCE, etc.) see a fully-populated `AuthorizeRequest` and run
their validation normally without any knowledge of PAR.

**`PARStore`** is an interface with three methods — `Store`, `GetByRequestURI`, `Delete` —
so any backing store can be used. `CacheStore` is the provided implementation, which wraps
`pkg/cache.CacheProvider[*StoredRequest]` and supports both Redis and in-memory backends:

```go
// Redis-backed (production)
store, err := par.NewCacheStoreFromConfig(cache.Config{
    Type:       cache.TypeRedis,
    Name:       "par",
    DefaultTTL: 90 * time.Second,
    Redis:      cache.RedisConfig{Address: "localhost:6379"},
})

// In-memory (local dev / tests)
store, err := par.NewCacheStoreFromConfig(cache.Config{
    Type:       cache.TypeInMemory,
    Name:       "par",
    DefaultTTL: 90 * time.Second,
})

parCapability := par.NewPAR(store)
```

---

## Server Setup

```go
oidcServer := oidc.NewOIDCServer(
    authnProvider,          // implements AuthnProvider interface
    parCapability,         // AuthorizationResolver, EndpointRegistrar, DiscoveryContributor
    pkceCapability,        // AuthorizationHook, TokenHook, DiscoveryContributor
    clientSecretCapability, // TokenHook, DiscoveryContributor
)

mux := http.NewServeMux()
mux.HandleFunc("/authorize", oidcServer.Authorize)
mux.HandleFunc("/token",     oidcServer.Token)
mux.HandleFunc("/userinfo",  oidcServer.Userinfo)
mux.HandleFunc("/.well-known/openid-configuration", oidcServer.Discovery)

// Mounts capability-defined endpoints (e.g. /par)
oidcServer.RegisterCapabilities(mux,
    parCapability
)
```

---

## Adding a New Capability

1. Create a package under `capabilities/` (or anywhere — the server takes `Capability` values,
   not import paths).
2. Implement `Capability` (`Name() string`).
3. Implement whichever hook interfaces apply.
4. Optionally implement `Ordered` to control execution position within a phase.
5. Register it with `oidc.NewOIDCServer(...)`.

No changes to `pkg/oidc` core are needed.
