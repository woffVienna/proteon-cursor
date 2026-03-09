# Proteon Engineering Guide

This document defines the repository structure, conventions, and
development workflow for the Proteon monorepo.

Architecture intent and constraints are defined in `docs/architecture/`
(see `00_CONTEXT.md`, `01_PRINCIPLES.md`, `02_WORKFLOW.md`, `03_INDEX.md`).
This document aligns with those and specifies how they apply to the repo.

It is intentionally opinionated. Deviations require explicit discussion.

------------------------------------------------------------------------

# 1. Architectural Principles

## 1.1 Monorepo, Independent Services

Proteon is a monorepo containing multiple independent Go services.

Each service:

-   Lives under `services/<service>/`
-   Is a fully independent Go module (`go.mod` per service)
-   Must not import other services directly
-   Communicates with other services via HTTP or asynchronous events
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

-   Logging abstractions
-   Configuration loading/orchestration
-   Error primitives
-   Observability helpers
-   Cross-cutting middleware utilities

Rules:

-   Services may import `libs/platform`.
-   `libs/platform` must not import any service.
-   `libs/platform` must not contain business logic.

------------------------------------------------------------------------

## 1.3 Layered Clean Architecture

Each service follows a strict internal layering model:

    internal/
      domain/
      application/
      adapters/
      platform/

Dependency direction is strictly:

    adapters → application → domain

Never the reverse.

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

Convention:

-   Each service must expose its entrypoint under `cmd/<svc>/main.go`.
-   The `<svc>` name must match the service folder name.

------------------------------------------------------------------------

# 3. Coding Style Principles

Proteon prefers cohesive, struct-based design.

## 3.1 Preferred Pattern

Use constructor functions and receiver methods:

    type Service struct {
        repo Repository
        logger Logger
    }

    func NewService(repo Repository, logger Logger) *Service {
        return &Service{repo: repo, logger: logger}
    }

    func (s *Service) Execute(ctx context.Context, input Input) error {
        ...
    }

## 3.2 Avoid

-   Large standalone functions that pass many dependencies as
    parameters.
-   "Function soup" orchestration across layers.

## 3.3 When to Use Pure Functions

-   Domain logic
-   Stateless mappers
-   Validation helpers
-   Small utility helpers

Application services, repositories, handlers, schedulers, and clients
must be struct-based.

------------------------------------------------------------------------

# 4. Configuration Model

Proteon uses a two-layer configuration model.

## 4.1 CoreConfig

-   Source: environment variables
-   Immutable at runtime
-   Loaded at service startup
-   Fail-fast validation
-   Changes require restart or redeploy

## 4.2 RuntimeConfig

-   Source: Postgres (`service_runtime_settings`)
-   Defaults applied first
-   DB overrides applied second
-   Validation executed at startup
-   Changes require service restart
-   No live mutation

Shared configuration loading logic lives in `libs/platform`.
Service-specific config structs live inside each service.

Configuration is resolved once at startup.

------------------------------------------------------------------------

# 5. OpenAPI & Code Generation

Each service owns its HTTP API. See `docs/architecture/system/INTEGRATION_CONTRACTS.md` for the full model.

-   Per service: `services/<service>/api/openapi.yml` is the source of truth for that service’s API.
-   Shared schemas may live in `libs/api/openapi/`.
-   Generated server code lives inside the service: `internal/adapters/http/generated/server/` (e.g. `openapi.gen.go`, committed).
-   Shared HTTP client artifacts for other services go under `contracts/http/<service>/` when used.
-   Bundled spec: `.build/generated/openapi.bundle.yml` (ignored).

------------------------------------------------------------------------

# 6. Tooling

-   Node tooling: `tools/node/`
-   Go tooling: `tools/bin/`
-   Local Kubernetes assets: `infra/k8s/local/`
-   Service Helm charts: `infra/k8s/charts/`

------------------------------------------------------------------------

# 7. Build Artifacts

    .build/
      generated/openapi.bundle.yml
      bin/<service>

Ignored via `.gitignore`.

------------------------------------------------------------------------

# 8. Helm Chart Conventions

-   Service Helm charts live in `infra/k8s/charts/<service>/`.
-   Use `values.yaml` for defaults and `values-local.yaml` for local overrides.
-   Namespace is controlled by Helm `--namespace`; do not hardcode
    `metadata.namespace` in chart templates.

------------------------------------------------------------------------

# 9. Makefile Responsibilities

Root Makefile:

-   setup
-   tooling
-   tooling-node
-   tooling-go
-   work
-   create-service
-   generate
-   verify-generated
-   check
-   test
-   build
-   clean
-   cluster-up
-   cluster-down
-   ns-up
-   ns-down
-   deps-install
-   deps-uninstall
-   image-load
-   deploy
-   deploy-all
-   stack-up
-   stack-down

Service Makefile:

-   tidy
-   fmt
-   lint
-   generate
-   test
-   build
-   dev
-   run
-   clean
-   containerise

------------------------------------------------------------------------

# 10. Local Development Workflow

From repo root:

    make setup
    make generate
    make test
    make stack-up

From a service folder:

    make generate
    make dev
    make run

After adding new imports:

    go mod tidy

Local runtime standard is k3d + Helm.
Do not use docker-compose orchestration in this repository.

------------------------------------------------------------------------

# 11. CI Contract

Minimum CI pipeline:

    make setup
    make verify-generated
    make test

------------------------------------------------------------------------

# 12. Git Policy

Ignored:

-   services/\*/.build/
-   services/\*/.env
-   tools/node/node_modules/
-   tools/bin/

Committed:

-   Generated Go server code
-   go.work

------------------------------------------------------------------------

# 13. Deferred Decisions

-   (none)

------------------------------------------------------------------------

This document evolves with the system.
