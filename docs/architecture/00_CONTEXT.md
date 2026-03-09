# Proteon Architecture Context

This document is the primary entry point for understanding the Proteon
architecture.

It provides the high-level system context, core constraints, and the reading
path for the more specific architecture documents in this directory.

Proteon is a Go monorepo containing independent microservices.

------------------------------------------------------------------------

# 1. System Overview

Proteon is designed as a microservice platform with explicit service
boundaries.

Core characteristics:

- each service lives under `services/<service>/`
- each service is its own Go module
- services are independently evolvable and deployable
- cross-service imports are forbidden
- services communicate via HTTP APIs or asynchronous events only
- shared technical code lives in `libs/platform`
- reusable integration artifacts live in `contracts/`

Primary technology context:

- Go
- OpenAPI-first API design
- Postgres
- NATS with JetStream
- Kubernetes
- Helm
- k3d for local runtime

Key top-level repository areas:

- `services/`
- `libs/platform`
- `contracts/`
- `infra/`
- `tools/`
- `docs/`

------------------------------------------------------------------------

# 2. Core Architectural Constraints

The following constraints are non-negotiable.

- services must not import each other directly
- services must not access each other’s databases
- business logic must not be centralized in `libs/platform`
- integration must happen through explicit HTTP or event boundaries
- configuration is resolved once at startup
- runtime configuration mutation is not allowed

These constraints exist to preserve loose coupling and prevent the monorepo
from collapsing into a distributed monolith.

------------------------------------------------------------------------

# 3. Internal Service Architecture

Each service follows this dependency direction:

`adapters -> application -> domain`

Layer responsibilities:

Adapters:

- transport
- persistence
- messaging
- generated server code
- external integrations

Application:

- use cases
- orchestration
- DTOs
- service interfaces

Domain:

- business concepts
- rules
- invariants

Rules:

- dependency direction must never be reversed
- domain must remain framework independent
- orchestration belongs in application, not domain
- transport and persistence concerns remain in adapters

Canonical service directory layout:

    services/<svc>/
      api/
        openapi.yml
        oapi-codegen.server.yml
        oapi-codegen.client.yml

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

Each service exposes its entrypoint under `cmd/<svc>/main.go`.
The `<svc>` name must match the service folder name.

------------------------------------------------------------------------

# 4. Communication Model

Proteon uses two communication styles.

## HTTP APIs

Use HTTP when:

- the caller needs an immediate result
- request/response is naturally synchronous
- the operation is client-driven

HTTP contracts are owned by the producing service and specified in:

`services/<service>/api/openapi.yml`

Generated reusable HTTP client artifacts are published under:

`contracts/http/<service>/`

## Events

Use events when:

- communication is asynchronous
- temporal decoupling is useful
- multiple consumers may react independently
- the producer should not depend on downstream runtime availability

Event contracts are published under:

`contracts/events/<service-or-domain>/`

Proteon defaults to domain events and event choreography.

------------------------------------------------------------------------

# 5. Shared Code Model

Shared technical code lives in:

`libs/platform`

Allowed examples:

- logging
- configuration loading and orchestration
- observability
- middleware
- technical helpers

Forbidden examples:

- business logic
- service-specific logic
- shared domain ownership
- cross-service orchestration
- imports from services

`contracts/` exists separately from `libs/platform` because integration
contracts and technical utilities are different concerns.

------------------------------------------------------------------------

# 6. Runtime and Development Model

Configuration is resolved once at startup using a two-layer model.

CoreConfig:

- source: environment variables (and `.env.local` for local development)
- immutable at runtime
- loaded at service startup
- fail-fast validation
- changes require restart or redeploy

RuntimeConfig:

- source: Postgres (`service_runtime_settings`)
- defaults applied first, DB overrides applied second
- validation executed at startup
- changes require service restart
- no live mutation

Rules:

- service-specific typed configuration belongs inside the service
- shared configuration loading/orchestration belongs in `libs/platform`
- invalid configuration must fail fast

Local development standard:

- k3d
- Helm
- Make-based workflows

Do not introduce docker-compose orchestration unless explicitly agreed.

------------------------------------------------------------------------

# 7. AI Workflow

Proteon uses a split AI workflow for architecture and implementation:

- ChatGPT for architecture reasoning, design intent, and tradeoff discussion
- Cursor Plan for repository-aware planning and mapping design to repo changes
- Cursor Agent for bounded code changes and incremental implementation

Details and interaction rules are in `02_WORKFLOW.md`.

------------------------------------------------------------------------

# 8. Service Roles

Proteon recognizes three primary service roles.

- edge services
- domain services
- worker services

Examples:

- `api-gateway` is an edge service
- `identity` is a domain service
- background processors are worker services

The detailed definitions live in `system/SERVICE_TYPES.md`.

------------------------------------------------------------------------

# 9. Reading Path

Read the architecture documents in this order:

1. `00_CONTEXT.md`
2. `01_PRINCIPLES.md`
3. `02_WORKFLOW.md`
4. `03_INDEX.md`

Then use the topic-specific documents under `system/` as needed.

------------------------------------------------------------------------

# 10. Deferred Decisions

The following decisions are explicitly unresolved:

- (none)

Do not design against assumptions in these areas without explicit discussion.

------------------------------------------------------------------------

# 11. Source of Truth Rule

Architecture intent must be persisted in:

`docs/architecture/`

Chat history is not the source of truth.

When architecture changes materially:

- update the relevant architecture document
- update `03_INDEX.md` if the structure changes
- add an architecture brief or ADR when appropriate
