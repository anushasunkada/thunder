# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

@AGENTS.md

## Build & Run Commands

All commands run from the repo root unless noted.

| Task | Command |
|------|---------|
| Build everything | `make build` |
| Build backend only | `make build_backend` |
| Build frontend only | `make build_frontend` |
| Run (backend + frontend) | `make run` |
| Run backend only | `make run_backend` |
| Run frontend only | `make run_frontend` |
| Lint backend | `make lint_backend` |
| Lint frontend | `make lint_frontend` |
| Lint all | `make lint` |
| Run all unit tests | `make test_unit` |
| Run integration tests | `make test_integration` |
| Regenerate mocks | `make mockery` |
| Regenerate i18n defaults | `make generate_i18n` |

### Run a single Go test / package

```bash
cd backend
go test ./internal/<package>/... -run TestFunctionName -v
```

### Run a single frontend test

```bash
cd frontend
pnpm test --filter=<package-name>        # e.g. thunder-gate
# or within the app directory:
pnpm exec vitest run path/to/Component.test.tsx
```

## Architecture

See [ARCHITECTURE.md](ARCHITECTURE.md) for the authoritative reference. Key points:

**Single binary** serving REST API + two React SPAs (`/gate` login UI, `/console` admin UI), backed by three SQLite/Postgres databases: `configdb`, `runtimedb`, `userdb`.

### Backend (`backend/internal/`)

Domain packages follow a strict `handler → service → store` layering, each with a single `Initialize(mux, deps)` in `init.go` that wires dependencies and registers routes. The `cmd/server/servicemanager.go` orchestrates all `Initialize` calls in dependency order.

- `flow/flowexec` — flow execution engine; auth/registration are JSON node graphs stepped via `POST /flow/execute`
- `flow/executor/` — one file per executor; add a new one by implementing `core.ExecutorInterface`, naming it in `constants.go`, and registering in `init.go`
- `oauth/` — full OAuth 2.0 + OIDC server (authorize, token, introspect, userinfo, JWKS, DCR)
- `authn/` — credential / OTP / passkey / social login mechanisms
- `system/` — shared infrastructure: config, database clients, cache, JWT/JOSE, security middleware, logging, i18n

**Database rules:** every table has `DEPLOYMENT_ID` as the last column in composite PKs; always the last parameter in queries. Use `DBClient` + `DBQuery` from `internal/system/database`.

**Errors:** `serviceerror.ServiceError` internally; `sysutils.WriteErrorResponse` for HTTP. Never expose 5xx details in responses.

**Public paths** (no JWT required): `/auth/**`, `/flow/execute/**`, `/oauth2/**`, `/.well-known/**`, `/gate/**`, `/console/**`, `/mcp/**` — full list in `system/security/permissions.go`.

### Frontend (`frontend/`)

pnpm workspace managed by Nx. Two apps: `thunder-gate` (app-native Flow API mode) and `thunder-console` (redirect OAuth mode). Shared packages under `frontend/packages/` (`@thunder/shared-contexts`, `design`, `hooks`, `i18n`, `utils`, `types`, `logger`).

Tests use Vitest + `@testing-library/react`. Place tests in `__tests__/` next to source; manual mocks in `__mocks__/`.

### References
 - Refer [feature and design documents](feature-docs) for more details.