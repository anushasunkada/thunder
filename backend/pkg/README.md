# Thunder backend `pkg` — public API for embedding

This directory holds **importable Go packages** for the ThunderID backend module (`github.com/thunder-id/thunderid`). Code under `internal/` is not importable from **other** repositories; `pkg/` is the supported surface for composing OAuth, flow execution, and related behavior in your own process.

This document describes **what lives here today**, **how an external application should use it**, and **how Thunder wires the same packages inside this repository**.

---

## Layout

| Path | Role |
|------|------|
| [`pkg/oauth`](oauth) | OAuth route registration: `RegisterRoutes` / `InitializeWithDependencies` and type aliases for dependency types. |
| [`pkg/oauth/deps`](oauth/deps) | `Dependencies` struct and validation; bundles everything needed to register OAuth HTTP handlers. |
| [`pkg/oauth/host`](oauth/host) | **Host contracts** you can implement in another module: DCR application registry and OAuth client lookup (no `internal/` types in method signatures). |
| [`pkg/oauth/app`](oauth/app) | Portable structs for **DCR** create/created payloads (`ApplicationCreate`, `ApplicationCreated`, nested OAuth registration fields). |
| [`pkg/oauth/oauthclient`](oauth/oauthclient) | Portable **resolved OAuth client** view returned by `host.InboundOAuth`. |
| [`pkg/flow`](flow) | Flow factory/cache (`InitializeCore`) and flow **execution** wiring (`InitializeExecution` / `InitializeExecutionWithDependencies`). |
| [`pkg/flow/deps`](flow/deps) | `ExecutionDependencies` for flow execution. |
| [`pkg/flow/host`](flow/host) | **Host contract** for entity inbound profile resolution used by flow execution (`InboundFlow`, `EntityInboundProfile`, not-found sentinel). |
| [`pkg/authnprovider`](authnprovider) | Public authentication-provider contract and models (used by broader auth; not all OAuth deps are migrated to `pkg` yet). |
| [`pkg/authnprovider/wire`](authnprovider/wire) | Small wiring helpers for authn provider construction in Thunder. |
| [`pkg/embed`](embed) | **`WireThunder`**: registers the same HTTP surface as the Thunder server binary (OAuth, flow, and the rest) on a caller `*http.ServeMux`, after loading config from a Thunder deployment directory (`THUNDER_HOME` layout). Optional `pkg/authnprovider.AuthnProvider` overrides the default authn manager. |

Supporting implementation for Thunder defaults lives under **`internal/oauth/hostbridge`** and **`internal/flow/hostbridge`** (adapters from Thunder `internal` services to the `pkg/.../host` contracts). External applications **do not** import those packages; they implement the `host` interfaces directly.

For **full** Thunder routing inside another binary without hand-building `oauthdeps.Dependencies`, use **`pkg/embed.WireThunder`** (see [`pkg/embed/embed.go`](embed/embed.go)).

---

## OAuth: entrypoints and dependencies

### Entry API

Import: `github.com/thunder-id/thunderid/pkg/oauth`

- **`oauth.RegisterRoutes(mux, deps)`** — same as `InitializeWithDependencies`; registers OAuth routes on an `*http.ServeMux`.
- **`oauth.InitializeWithDependencies(mux, deps)`** — validates `deps` then delegates to `internal/oauth`.

### Dependency bag

Import: `github.com/thunder-id/thunderid/pkg/oauth/deps` (often referenced as `oauthdeps` in Thunder code).

Struct: **`oauthdeps.Dependencies`**. Required fields include (non-exhaustive; see `Validate()` in [`oauth/deps/deps.go`](oauth/deps/deps.go)):

- **`Application`** — type **`oauthdeps.DCRApplication`** = `pkg/oauth/host.DCRApplication`: DCR create/delete using **`pkg/oauth/app`** types.
- **`Inbound`** — type **`oauthdeps.InboundOAuth`** = `pkg/oauth/host.InboundOAuth`: `GetOAuthClientByClientID` → **`pkg/oauth/oauthclient.Client`**.
- **`Transactioner`**, **`DBProvider`**, **`DeploymentID`**, **`DatabaseRuntimeType`**, and when the runtime type is Redis, **`RedisProvider`**.
- **`AuthnProvider`**, **`JWTService`**, **`JWEService`**, **`FlowExecService`**, and other collaborators — **still expressed as aliases to Thunder `internal/` interfaces** in [`oauth/deps/types.go`](oauth/deps/types.go). A standalone module **cannot implement those interfaces** (they are not in `pkg`). For full OAuth in another repo you currently either depend on Thunder-instantiated concrete services (if you stay inside the same module / a fork) or extend `pkg` further to narrow additional contracts.

### Host contracts (external-friendly)

Import: `github.com/thunder-id/thunderid/pkg/oauth/host`

- **`DCRApplication`** — `CreateApplication(ctx, *app.ApplicationCreate) (*app.ApplicationCreated, error)`, `DeleteApplication(ctx, appID) error`.
- **`InboundOAuth`** — `GetOAuthClientByClientID(ctx, clientID) (*oauthclient.Client, error)`.

Portable models: **`github.com/thunder-id/thunderid/pkg/oauth/app`**, **`github.com/thunder-id/thunderid/pkg/oauth/oauthclient`**.

---

## Flow: entrypoints and dependencies

### Core (graphs, factory)

Import: `github.com/thunder-id/thunderid/pkg/flow`

- **`flow.InitializeCore(cacheManager)`** — returns flow factory and graph cache (uses `internal/flow/core`).

### Execution (HTTP `/flow/execute`)

- **`flow.InitializeExecutionWithDependencies(mux, deps)`** — validates `flowdeps.ExecutionDependencies`, then `internal/flow/flowexec.Initialize`.
- **`flow.InitializeExecution(mux, flowMgt, inboundFlow, entityProvider, executorRegistry, observabilitySvc, cryptoSvc)`** — convenience wrapper building `ExecutionDependencies`.

### Execution dependency bag

Import: `github.com/thunder-id/thunderid/pkg/flow/deps`

Struct: **`flowdeps.ExecutionDependencies`**. Notable field:

- **`Inbound`** — type **`flowdeps.InboundFlow`** = `pkg/flow/host.InboundFlow`. Implement **`GetInboundClientByEntityID(ctx, entityID) (*host.EntityInboundProfile, error)`**.

Not-found semantics: wrap or return **`github.com/thunder-id/thunderid/pkg/flow/host.ErrEntityInboundNotFound`** so flow execution treats “missing profile” like Thunder’s inbound store.

Other fields (**`FlowMgtService`**, **`EntityProvider`**, **`ExecutorRegistry`**, **`RuntimeCryptoProvider`**, optional **`ObservabilitySvc`**) remain **aliases to `internal/` types** today (see [`flow/deps/types.go`](flow/deps/types.go)).

---

## How Thunder uses `pkg` (this repository)

Thunder’s HTTP server composes the same packages you would import, then fills dependencies with **internal** services.

| Concern | Where it is wired |
|--------|-------------------|
| OAuth `Dependencies` | [`cmd/server/servicemanager.go`](../cmd/server/servicemanager.go) — e.g. `Application: hostbridge.NewThunderApplication(applicationService)`, `Inbound: hostbridge.NewThunderInbound(inboundClientService)` with **`internal/oauth/hostbridge`**. |
| Flow execution | Same file — `flow.InitializeExecution(..., flowhostbridge.NewThunderInboundFlow(inboundClientService), ...)` using **`internal/flow/hostbridge`**. |
| OAuth + flow together | OAuth `Dependencies` includes **`FlowExecService`** from the initialized flow execution service so authorization/token paths can invoke flows. |

So: **Thunder = your app + pre-built internal services + `hostbridge` adapters**. **External app = your app + your implementations of `pkg/.../host` + whatever Thunder later exposes as `pkg` contracts for the remaining deps.**

---

## How an external Go application should use `pkg`

### 1. Add the module dependency

In your `go.mod`:

```text
require github.com/thunder-id/thunderid vX.Y.Z
```

Use the version or replace directive your team publishes (fork, pseudo-version, etc.).

### 2. Own the HTTP mux

```go
mux := http.NewServeMux()
```

### 3. Implement the host contracts you need

**OAuth — custom DCR / inbound only**

Implement:

- `github.com/thunder-id/thunderid/pkg/oauth/host`.`DCRApplication` using `pkg/oauth/app` types for create/created payloads.
- `github.com/thunder-id/thunderid/pkg/oauth/host`.`InboundOAuth` returning `pkg/oauth/oauthclient.Client` from `GetOAuthClientByClientID`.

Pass them as **`oauth.Dependencies.Application`** and **`oauth.Dependencies.Inbound`**.

**Flow execution — custom entity inbound resolution**

Implement:

- `github.com/thunder-id/thunderid/pkg/flow/host`.`InboundFlow` (entity profile for auth / registration / recovery graph selection).

Pass as **`flow.ExecutionDependencies.Inbound`**.

### 4. Satisfy remaining `Dependencies` / `ExecutionDependencies`

Until those fields are migrated to `pkg`-local interfaces, you must supply implementations **visible inside the Thunder module** (for example by contributing additional `pkg` contracts, or by running your binary inside a fork of Thunder that constructs internal services). Plan ahead for **database schema**, **config**, **crypto**, and **PAR** expectations documented elsewhere in the project.

### 5. Register routes

```go
import (
    "github.com/thunder-id/thunderid/pkg/flow"
    "github.com/thunder-id/thunderid/pkg/oauth"
)

// After building deps and flowDeps (and resolving FlowExecService for OAuth if needed):
if err := flow.InitializeExecutionWithDependencies(mux, flowDeps); err != nil { ... }
if err := oauth.InitializeWithDependencies(mux, oauthDeps); err != nil { ... }
```

Order can depend on whether OAuth needs an already-built `FlowExecService`; mirror Thunder’s `servicemanager` order if in doubt.

---

## `pkg/authnprovider` (summary)

`github.com/thunder-id/thunderid/pkg/authnprovider` exposes **`AuthnProvider`** and related result/metadata types for authentication providers. OAuth’s **`oauthdeps.Dependencies.AuthnProvider`** is still typed as Thunder’s internal manager interface in `oauth/deps/types.go`; integrating custom authn at the OAuth boundary may require additional `pkg` work or internal adapters depending on your deployment.

---

## Related reading

- Repository root [`README.md`](../../README.md) and [`Makefile`](../../Makefile) for build and run.
- [`AGENTS.md`](../../AGENTS.md) for contributor rules.
- In-tree OAuth/flow implementation: `internal/oauth`, `internal/flow/flowexec` (not importable from other modules’ `internal` rules, but useful when working inside Thunder).

---

## Stability note

The **`pkg/oauth/host`**, **`pkg/oauth/app`**, **`pkg/oauth/oauthclient`**, and **`pkg/flow/host`** contracts are intended as the **extension surface** for embedders. Other `pkg/oauth/deps` and `pkg/flow/deps` aliases may still change as more collaborators move behind `pkg`-local interfaces. When upgrading Thunder, re-run tests and check `deps.Validate()` / `ExecutionDependencies.Validate()` expectations.
