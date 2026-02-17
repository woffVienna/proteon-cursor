# Proteon -- Developer Guide

This document describes the canonical development workflow for the
Proteon monorepo.

Always prefer Make targets over ad-hoc commands.
Canonical target names and responsibilities are defined in
`ENGINEERING.md` (#6 Makefile Responsibilities).

------------------------------------------------------------------------

# 1. One-Time Setup

From repository root:

``` bash
make setup
```

This will:

-   Install Node tooling (`tools/node/`) -- Redocly\
-   Install Go tooling (`tools/bin/`) -- oapi-codegen, golangci-lint\
-   Initialize/update `go.work` to include all services and
    `libs/platform`

No manual Docker setup is required.\
The shared Docker network is created automatically by `make deps-up` or
`make stack-up`.

------------------------------------------------------------------------

# 2. Typical Daily Workflow

## 2.1 Generate OpenAPI Stubs (All Services)

``` bash
make generate
```

Use this after:

-   Modifying `api/openapi.yml`
-   Changing shared schemas in `libs/api/openapi`

------------------------------------------------------------------------

## 2.2 Run Tests (All Services)

``` bash
make test
```

------------------------------------------------------------------------

## 2.3 Full Local Verification (Pre-Push)

``` bash
make check
```

This runs:

-   `make verify-generated`
-   `make test`

If `verify-generated` fails, run:

``` bash
make generate
```

and commit the changes.

------------------------------------------------------------------------

# 3. Working on a Single Service

Navigate into the service:

``` bash
cd services/<service>
```

## 3.1 Generate (Service Only)

``` bash
make generate
```

------------------------------------------------------------------------

## 3.2 Active Development Mode (Golden Path)

``` bash
make dev
```

`dev` contract:

-   Always runs `generate` first
-   Loads environment variables from local `.env` (if present)
-   Runs `go run ./cmd/<service>`
-   Prevents running with stale OpenAPI stubs

### Required local `.env`

Example:

    DB_DSN=postgres://proteon:proteon@localhost:5432/proteon?sslmode=disable
    ENV=dev
    MARKET=AT
    SERVICE_NAME=identity

When running via `go run`:

-   Postgres host is `localhost`
-   Dependencies must be running via Docker

------------------------------------------------------------------------

## 3.3 Run Without Side Effects

``` bash
make run
```

-   Does NOT run `generate`
-   Uses current working tree state

------------------------------------------------------------------------

# 4. Local Dependencies (Docker)

Docker assets are located in:

    tools/docker/
      compose.deps.yml
      compose.services.yml

All containers use a shared Docker network:

    proteon

This network is created automatically by Make targets.

------------------------------------------------------------------------

## 4.1 Start Dependencies Only

``` bash
make deps-up
```

This:

-   Ensures Docker network exists
-   Starts Postgres
-   Exposes Postgres on `localhost:5432`

Stop dependencies:

``` bash
make deps-down
```

------------------------------------------------------------------------

## 4.2 Start Full Local Stack (All Services in Containers)

``` bash
make stack-up
```

This runs:

``` bash
docker compose   -f tools/docker/compose.deps.yml   -f tools/docker/compose.services.yml   up -d
```

Stop full stack:

``` bash
make stack-down
```

------------------------------------------------------------------------

# 5. Containerising a Service

Each service provides:

``` bash
make containerise
```

This builds:

    proteon/<service>:dev

Dockerfile location:

    services/<service>/Dockerfile

Build context is the repository root.

------------------------------------------------------------------------

## Container Networking Rules

When running inside Docker:

-   Postgres hostname is `postgres`
-   DSN example:

    postgres://proteon:proteon@postgres:5432/proteon?sslmode=disable

When running on host (`go run`):

-   Postgres hostname is `localhost`

This difference is expected.

------------------------------------------------------------------------

# 6. Configuration Model

Proteon uses a two-layer configuration model.

------------------------------------------------------------------------

## 6.1 CoreConfig (Immutable)

Source:

-   Environment variables
-   Optional static config files

Used for:

-   DB endpoints
-   Ports
-   Infra wiring
-   ENV
-   MARKET

Behavior:

-   Loaded at startup
-   Strict validation
-   Fail fast on missing required values
-   Changes require restart or redeploy

------------------------------------------------------------------------

## 6.2 RuntimeConfig (DB-Backed)

Source:

-   `service_runtime_settings` table in Postgres

Behavior:

-   Defaults applied first
-   DB overrides applied second
-   Unknown keys ignored with warning
-   Validation executed
-   Warnings logged
-   Changes require service restart
-   No redeploy required

Runtime configuration is resolved once at startup.\
No runtime mutation.

------------------------------------------------------------------------

# 7. Dependency Management (Per Service)

Each service is its own Go module.

After adding new imports (including generated OpenAPI server code):

``` bash
go mod tidy
```

This updates:

-   `go.mod`
-   `go.sum`

Run it inside the service directory.

------------------------------------------------------------------------

# 8. Build Artifacts

Each service writes build output to:

    .build/
      generated/openapi.bundle.yml
      bin/<service>

Bundled spec: `.build/generated/openapi.bundle.yml` (ignored).

`.build/` is ignored and must never be committed.

------------------------------------------------------------------------

# 9. CI Contract

CI runs:

``` bash
make setup
make verify-generated
make test
```

Expectations:

-   Generated server stubs are committed
-   All services compile
-   Tests pass
-   Clean checkout builds deterministically

------------------------------------------------------------------------

# 10. Architectural Rules

All code must follow:

-   `ENGINEERING.md`
-   `.cursorrules`

In particular:

-   No cross-service imports
-   Enforce dependency direction: adapters → application → domain
-   Domain is pure
-   Transport types stay in adapters
-   Configuration resolved once at startup

------------------------------------------------------------------------

This document describes the golden path.\
If you need to deviate, update `ENGINEERING.md` accordingly.
