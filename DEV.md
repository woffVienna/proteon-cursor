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

Prerequisites for local stack orchestration:

-   `k3d`
-   `kubectl`
-   `helm`

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
-   Starts with `RUNTIME_MODE=local` by default
-   Loads `.env.local` through the shared config loader
-   Runs `go run ./cmd/<service>`
-   Prevents running with stale OpenAPI stubs

### Local service env file

Example (`services/<service>/.env.local`):

    SERVICE_NAME=identity-service
    ENV=dev
    MARKET=AT
    VERSION=dev
    PORT=8081
    PUBLIC_BASE_URL=http://localhost:8081
    DB_DSN=postgres://proteon:proteon@localhost:5432/proteon?sslmode=disable

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

# 4. Local Dependencies and Stack (k3d + Helm)

Local runtime standard is Kubernetes via k3d.

Core targets:

-   `make cluster-up` / `make cluster-down`
-   `make ns-up` / `make ns-down`
-   `make deps-install` / `make deps-uninstall`
-   `make deploy SERVICE=identity`
-   `make wait-deps` / `make wait-services` / `make wait-ingress`
-   `make stack-up` / `make stack-down`

Service Helm charts are stored under:

    infra/k8s/charts/<service>/

------------------------------------------------------------------------

## 4.1 Start Dependencies Only

``` bash
make deps-install
```

This:

-   Ensures namespace `proteon-dev` exists
-   Installs Postgres and NATS via Helm
-   Enables JetStream for NATS

Stop dependencies:

``` bash
make deps-uninstall
```

------------------------------------------------------------------------

## 4.2 Start Full Local Stack

``` bash
make stack-up
```

This runs:

-   k3d cluster provisioning
-   namespace setup (`proteon-dev`)
-   Helm dependency install (`postgresql`, `nats`)
-   dependency readiness checks
-   service image build + k3d image import
-   Helm deploy for charted services
-   deployment + ingress readiness checks

Stop full stack:

``` bash
make stack-down
```

Validate local ingress route:

``` bash
curl http://localhost:8080/
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

`containerise` first runs `generate` so `.build/generated/openapi.bundle.yml`
exists before image build. Docker runtime image includes this bundled spec so
`/swagger` and `/openapi.yaml` work inside containers as well.

------------------------------------------------------------------------

## Container Networking Rules

When running inside Kubernetes:

-   Postgres hostname is `postgresql`
-   DSN example:

    postgres://proteon:proteon@postgresql:5432/proteon?sslmode=disable

When running on host (`go run`):

-   Postgres hostname is `localhost`

This difference is expected.

------------------------------------------------------------------------

# 6. Configuration Model

Proteon services resolve configuration once at startup from environment.

Shared loader:

-   `libs/platform/config/env.go` resolves `RUNTIME_MODE` and loads
    `.env.<mode>` from the service directory
-   `libs/platform/config/loader.go` provides a typed loader with defaults
    for common fields (`ServiceName`, `ENV`, `MARKET`, `VERSION`, `HTTP`)
-   Each service assembles bespoke typed fields in
    `services/<service>/internal/platform/config`

------------------------------------------------------------------------

## 6.1 Runtime Modes and Env Files

`RUNTIME_MODE` values:

-   `local`  -> `.env.local`
-   `docker` -> `.env.docker`
-   `cloud`  -> `.env.cloud` (optional; runtime env injection is preferred)

Key names stay the same across all modes (for example `PORT`, `DB_DSN`,
`JWT_ISSUER`).

------------------------------------------------------------------------

## 6.2 Resolution Order and Typed Service Config

Resolution order per key:

1. Existing process env (shell/docker/cloud injected)
2. Selected `.env.<mode>` file value
3. Built-in default (if defined)

Service-specific typed config is parsed into `Config.Service` by the service
loader callback (for example `Config.Service.JWT` in identity).

Configuration is resolved once at startup.\
Changes require restart/redeploy.

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
