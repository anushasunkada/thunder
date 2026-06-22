# thunderidengine

An embeddable ThunderID identity engine. It mounts the flow-metadata (`GET /flow/meta`),
flow-execution (`POST /flow/execute`), and OAuth2/OIDC (`/oauth2/*`) endpoint groups onto a
caller-supplied `http.ServeMux`. Short-lived runtime state is persisted in a caller-supplied Redis
connection; the engine never opens a SQL database at runtime.

## Quickstart

```go
rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

cfg, err := thunderidengine.LoadConfig("/etc/thunderid")
if err != nil {
    log.Fatal(err)
}

eng, err := thunderidengine.New(
    thunderidengine.WithRedis(rdb, "thunderid:"),
    thunderidengine.WithConfig("/etc/thunderid", cfg),
    thunderidengine.WithPKIKey("default", "certs/server.crt", "certs/server.key"),
    thunderidengine.WithHostActorProvider(myActorProvider{}),
    thunderidengine.WithHostAuthnProvider(myAuthnProvider{}),
    thunderidengine.WithExecutorDependencies(thunderidengine.ExecutorDependencies{}),
    thunderidengine.WithEnabledExecutors(
        "CredentialsAuthExecutor", "AuthorizationExecutor",
        "AuthAssertExecutor", "ConsentExecutor",
    ),
)
if err != nil {
    log.Fatal(err)
}
defer eng.Shutdown(context.Background())

handler, _ := eng.Handler()
log.Fatal(http.ListenAndServe(":9443", handler))
```

A complete, compiling version (host providers and a custom executor) is in
[`example_test.go`](example_test.go).

## What you supply

The engine is built with functional options. The dependency-free SDK contract lives in the
[`host`](host) and [`runtime`](runtime) packages, so an external application implements the engine's
identity source without importing any `internal/*` type.

| Concern | Option(s) | Required |
|---------|-----------|----------|
| Runtime state store | `WithRedis` | Yes |
| Server configuration (GateClient, OAuth, flow, crypto, declarative resources, ...) | `WithConfig` / `LoadConfig` / `LoadConfigFromPaths` | Yes |
| Crypto / JWT / JWE | `WithPKIKey` (derives them) or `WithRuntimeCrypto` (+ `WithJWTService` / `WithJWEService`) | Yes |
| Identity source (entities, applications, inbound clients) | `WithHostActorProvider` (external) or `WithActorProvider` (in-tree) | Yes |
| Authentication | `WithHostAuthnProvider` (external) or `WithAuthnProvider` (in-tree) | Yes |
| Consent enforcement | `ConsentEnforcer` field on `ExecutorDependencies` | When `ConsentExecutor` is enabled |
| Executors | `WithExecutorDependencies` + `WithEnabledExecutors`, or `WithExecutorRegistry`; plus `WithCustomExecutors` | Yes |
| Observability | `WithObservability` | No |
| System-of-record services (OU, resource, IDP, authz, attribute cache, design, flow, i18n) | `WithOUService`, `WithResourceService`, ... | No — fall back to declarative |

## Declarative fallback (file-based storage)

When declarative mode is enabled in the configuration (`declarative_resources.enabled: true`), any
system-of-record service you do not inject is built read-only from declarative file-based resources.
The fallback is all-or-nothing for that set: if any of OU, resource, IDP, authz, attribute-cache,
design, or flow is missing, the full declarative graph is built. Management REST routes the
underlying services register are mounted on a throwaway mux and are never exposed on your mux.

## Executors

You can enable a subset of the built-in executors and add your own, mixed together:

- `WithExecutorDependencies(...)` + `WithEnabledExecutors("CredentialsAuthExecutor", ...)` — the
  engine builds the registry and registers the named built-ins (an empty list registers all).
- `WithCustomExecutors(map[string]thunderidengine.ExecutorInterface{...})` — registers your
  executors on top of that registry, so they run alongside the enabled built-ins. A custom executor
  whose name matches a built-in overrides it. This also layers onto a registry supplied via
  `WithExecutorRegistry`.

Author a custom executor from the public package without importing `internal/*`: implement
`thunderidengine.ExecutorInterface`, embedding a `thunderidengine.NewBaseExecutor(...)` value to
inherit the boilerplate methods and overriding only `Execute`. The executor name you register is
the name a flow `TASK` node references.

```go
type greetExecutor struct {
    thunderidengine.ExecutorInterface
}

func (*greetExecutor) Execute(
    *thunderidengine.ExecutorNodeContext,
) (*thunderidengine.ExecutorResponse, error) {
    return &thunderidengine.ExecutorResponse{Status: thunderidengine.ExecComplete}, nil
}
```

## Notes and constraints

- The server runtime configuration is a process singleton (first initialization wins), so only one
  engine instance per process is supported.
- The dependency graph still links the SQL drivers (`lib/pq`, `modernc.org/sqlite`) even though they
  are unused at runtime, so an embedding application must not also blank-import those drivers
  (`database/sql` panics on duplicate registration).
- Dynamic Client Registration is not part of the engine.
- The Redis connection is owned by the caller and is not closed by `Shutdown`.
