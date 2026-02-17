# Proteon Engineering Guide

This document defines the architectural principles, structure, and
development workflow for the Proteon monorepo.

It is intentionally opinionated. Deviations require explicit discussion.

------------------------------------------------------------------------

# 1. Architectural Principles

## 1.1 Monorepo, Independent Services

Proteon is a monorepo containing multiple independent Go services.

Each service:

- Lives under `services/<service>/`
- Is a fully independent Go module (`go.mod` per service)
- Must not import other services directly
- Communicates with other services via HTTP or asynchronous events
    only

Cross-service imports are forbidden.

------------------------------------------------------------------------

## 1.2 Shared Libraries (libs/)

Shared Go code lives under:

    libs/
      platform/

`libs/platform` is a standalone Go module used by services for shared
technical concerns.

It may contain:

- Logging abstractions
- Configuration loading
- Error primitives
- Observability helpers
- Cross-cutting middleware utilities

Rules:

- Services may import `libs/platform`.
- `libs/platform` must not import any service.
- `libs/platform` must not contain business logic.

------------------------------------------------------------------------

## 1.3 Layered Clean Architecture

Each service follows a strict internal layering model:

    internal/
      domain/
      application/
      adapters/
      platform/

------------------------------------------------------------------------

# 2. Service Structure

Canonical layout:

    services/<svc>/
      api/
        openapi.yml
        oapi-codegen.server.yml

      internal/
        adapters/
          db/
          http/
            generated/server/
          nats/
        application/
          dto/
          interfaces/
          services/
        domain/
          model/
          rules/
        platform/

      cmd/<svc>/
        main.go

      .build/

      Makefile
      go.mod

------------------------------------------------------------------------

# 3. OpenAPI & Code Generation

- `api/openapi.yml` is the source of truth.
- Shared schemas live in `libs/api/openapi/`.
- Bundled spec: `.build/generated/openapi.bundle.yml` (ignored).
- Generated server stubs:
    `internal/adapters/http/generated/server/openapi.gen.go`
    (committed).

------------------------------------------------------------------------

# 4. Tooling

- Node tooling: `tools/node/`
- Go tooling: `tools/bin/`

------------------------------------------------------------------------

# 5. Build Artifacts

    .build/
      generated/openapi.bundle.yml
      bin/<service>

Ignored via `.gitignore`.

------------------------------------------------------------------------

# 6. Makefile Responsibilities

This section defines the canonical target names referenced by `DEV.md`.

Root Makefile:

- setup
- tooling
- tooling-node
- tooling-go
- work
- create-service
- generate
- verify-generated
- check
- test
- build
- clean
- deps-up
- deps-down
- stack-up
- stack-refresh-up
- stack-down

Service Makefile:

- tidy
- fmt
- lint
- generate
- test
- build
- dev
- run
- clean
- containerise

------------------------------------------------------------------------

# 7. Local Development Workflow

From repo root:

    make setup
    make generate
    make test

From a service folder:

    make generate
    make dev
    make run

After adding new imports:

    go mod tidy

------------------------------------------------------------------------

# 8. CI Contract

Minimum CI pipeline:

    make setup
    make verify-generated
    make test

------------------------------------------------------------------------

# 9. Git Policy

Ignored:

- services/\*/.build/
- tools/node/node_modules/
- tools/bin/

Committed:

- Generated Go server code
- go.work

------------------------------------------------------------------------

# 10. Deferred Decisions

- SDK / client codegen location
- Cross-service client strategy
- Event schema versioning strategy

------------------------------------------------------------------------

This document evolves with the system.
