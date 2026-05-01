# Thunder – Architecture Reference

Go IAM server (`github.com/asgardeo/thunder`). Single binary serving a REST API + two React SPAs (`/gate`, `/console`).

## Structure

```text
backend/cmd/server/
  main.go               # startup
  servicemanager.go     # calls every internal/*/init.go to register routes
  bootstrap/flows/      # JSON auth/registration flow definitions (auto-seeded)
  repository/           # configdb.db · runtimedb.db · userdb.db created at runtime in the configured data directory (SQLite or Postgres)
backend/internal/
  authn/                # credential / OTP / passkey / social login
  oauth/                # OAuth 2.0 + OIDC server (authorize, token, introspect, userinfo, JWKS, DCR)
  flow/flowexec/        # flow execution engine  →  POST /flow/execute
  flow/executor/        # one file per executor; all names in constants.go; registered in init.go
  flow/core/            # ExecutorInterface, node/graph model
  flow/mgt/             # flow CRUD API
  consent/ application/ user/ group/ role/ ou/ idp/   # management domains
  system/               # config · database · cache · jose/jwt · security · mcp · log · i18n
frontend/apps/
  gate/         # login/registration SPA  (@asgardeo/react — app-native mode)
  thunder-console/      # admin SPA               (@asgardeo/react — redirect mode)
frontend/packages/      # @thunder/contexts · design · hooks · i18n · utils · types · logger
samples/apps/           # react-sdk-sample · react-api-based-sample · react-vanilla-sample
```

## Backend patterns

- Each domain package: `handler → service → store`, single `Initialize(mux, …)` in `init.go`.
- Public paths (no JWT): `/auth/**`, `/flow/execute/**`, `/oauth2/**`, `/.well-known/openid-configuration/**`, `/.well-known/oauth-authorization-server/**`, `/.well-known/oauth-protected-resource`, `/gate/**`, `/console/**`, `/mcp/**` — full list in `system/security/permissions.go`.
- Errors: `serviceerror.ServiceError` internally; `sysutils.WriteErrorResponse(w, status, errConst)` for HTTP.

## Flow engine

Authentication/registration are JSON node graphs (`START → PROMPT → TASK → DECISION → COMPLETE`). The engine steps through nodes, persisting state in `runtimedb` across requests. Each `TASK` node names an executor (e.g. `"BasicAuthExecutor"`). To add one: implement `core.ExecutorInterface`, add name to `executor/constants.go`, register in `executor/init.go`.

## Asgardeo React SDK

| Mode | `AsgardeoProvider` props | Used in |
|------|--------------------------|---------|
| Redirect (Thunder-hosted login) | `clientId` + `baseUrl` + `platform="AsgardeoV2"` | `Console`, `react-sdk-sample` |
| App-native (Flow API) | `applicationId` + `baseUrl` + `platform="AsgardeoV2"` | `Gate`, `react-api-based-sample` |

`clientId` vs `applicationId` is the critical distinction. Common primitives: `useAsgardeo()`, `<SignedIn/Out>`, `<SignInButton/SignOutButton>`, `<ProtectedRoute>` (`@asgardeo/react-router@2.0`).

## Auth Flow 

### (Mode 1)

```text
Client → GET /oauth2/authorize → 302 /gate?executionId=…
Gate SPA → POST /flow/execute (loop) → 302 redirect_uri?code=…
Client → POST /oauth2/token → { access_token, id_token }
```

### Mode 2

Client posts directly to `POST /flow/execute` with `applicationId` (first call) then `executionId` until `status: COMPLETE`.
