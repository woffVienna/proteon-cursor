00 – Proteon Architecture Context

This chat defines the authoritative architecture context for the Proteon platform.

All architecture discussions in this project should treat the following
documents as the primary source of truth.

If a design question depends on details not present here, ask for the
relevant repository document instead of making assumptions.

When discussing architecture:

Always ask clarifying questions before proposing a solution.

Preferred response structure:

1. Clarifying Questions
2. Architecture Overview
3. Key Components
4. Data or Event Flow
5. Implementation Plan
6. Risks and Tradeoffs


---------------------------------------------------------------------

SESSION_BOOTSTRAP.md

# Proteon Session Bootstrap

Proteon is a Go monorepo containing independent microservices.

## Core Architecture Constraints

- Each service lives under `services/<service>/`
- Each service is its own Go module
- Cross-service imports are forbidden
- Services communicate via HTTP APIs or asynchronous events only
- Shared technical code lives in `libs/platform`
- `libs/platform` must not contain business logic
- `libs/platform` must not import services

## Internal Service Architecture

Each service follows this dependency direction:

adapters → application → domain

Rules:

- adapters contain transport, persistence, and framework integrations
- application contains use cases, orchestration, DTOs, and interfaces
- domain contains pure business concepts and rules
- dependency direction must never be reversed
- domain must remain pure

## Configuration Model

Configuration is resolved once at startup.

- base configuration comes from environment
- service-specific typed configuration is assembled inside the service
- shared configuration loading/orchestration lives in `libs/platform`
- no live runtime mutation
- changes require restart/redeploy

## Local Development Standard

- local runtime standard is `k3d + Helm`
- prefer Make targets over ad-hoc commands
- do not introduce docker-compose orchestration

## AI Workflow

- ChatGPT is used for architecture reasoning and design intent
- Cursor Plan is used for repo-specific planning
- Cursor Agent is used for bounded execution
- architecture/context documents in `docs/architecture/` are the shared intent layer for both tools

## Interaction Rule

Always ask clarifying questions before proposing a solution for non-trivial architecture or design problems.


---------------------------------------------------------------------

PROTEON_CONTEXT.md

# Proteon System Context

Proteon is a Go monorepo built around independent microservices.

## System Type

- microservice platform
- monorepo
- independently deployable services
- explicit integration boundaries

## Technology Context

Current core stack includes:

- Go
- OpenAPI
- Postgres
- NATS with JetStream
- Kubernetes for local/runtime orchestration
- Helm for deployment packaging

## Repository Structure

Key top-level areas:

- `services/`
- `libs/platform`
- `infra/`
- `tools/`
- `docs/`

## Service Model

Each service:

- lives under `services/<service>/`
- has its own Go module
- owns its own domain logic
- exposes explicit integration boundaries
- must not be imported by another service

## Communication Model

Service-to-service interaction happens via:

- HTTP APIs
- asynchronous events

Avoid hidden coupling or implicit dependency paths.

## Shared Technical Layer

Shared technical concerns live in:

`libs/platform`

This layer exists for technical reuse, not business logic centralization.

## Configuration Model

Configuration is resolved at startup.

- shared loading/orchestration in `libs/platform`
- service-specific typed config in the service
- restart required for config changes

## Development Model

Proteon uses:

- Make-based workflows
- OpenAPI-first server generation
- k3d + Helm for local stack orchestration

## AI Usage Model

Proteon uses a split AI workflow:

- ChatGPT for architecture/design intent
- Cursor Plan for repo-specific planning
- Cursor Agent for bounded code execution


---------------------------------------------------------------------

ARCHITECTURE_PRINCIPLES.md

# Proteon Architecture Principles

This document defines the stable architectural principles of the Proteon monorepo.

## 1. Monorepo with Independent Services

Proteon is a monorepo containing multiple independent Go services.

Each service:

- lives under `services/<service>/`
- is a standalone Go module
- is independently evolvable
- must not import another service directly
- communicates with other services via HTTP or asynchronous events only

Cross-service imports are forbidden.

## 2. Shared Libraries

Shared Go code lives under:

`libs/platform`

Allowed responsibilities:

- logging abstractions
- configuration loading/orchestration
- error primitives
- observability helpers
- middleware and technical cross-cutting helpers

Forbidden responsibilities:

- business logic
- service-specific logic
- cross-service orchestration
- domain ownership that belongs to a service

`libs/platform` must not import services.

## 3. Internal Service Layering

Each service follows a strict internal layering model:

adapters → application → domain

Interpretation:

- adapters: transport, database, messaging, generated server code, external integrations
- application: use cases, orchestration, interfaces, DTOs
- domain: pure business concepts, rules, invariants

Rules:

- never reverse dependency direction
- keep transport types inside adapters
- keep domain free of framework concerns
- keep orchestration out of domain

## 4. Configuration

Proteon uses startup-resolved configuration.

Rules:

- configuration is resolved once at startup
- shared configuration loading/orchestration belongs in `libs/platform`
- service-specific typed config belongs inside each service
- no live runtime mutation
- invalid configuration must fail fast

## 5. API and Event Boundaries

Service integration happens through explicit boundaries only:

- HTTP APIs
- asynchronous events

Avoid hidden coupling through:

- direct code imports between services
- shared service-owned database access
- leaking internal types across boundaries

## 6. Incremental Evolution

Prefer:

- incremental migrations
- bounded changes
- minimal diffs
- preserving established patterns unless intentionally changed

Avoid:

- broad refactors without a clear plan
- introducing new patterns where existing patterns already fit
- silent architectural drift

## 7. Documentation and Intent

Architecture intent must be documented in `docs/architecture/`.

When architecture changes materially, the corresponding documents should be updated.


---------------------------------------------------------------------

ENGINEERING_WORKFLOW.md

# Proteon Engineering Workflow

This document describes the intended workflow for development and AI-assisted work.

## 1. Canonical Development Workflow

Prefer Make targets over ad-hoc commands.

Core root targets:

- `make setup`
- `make generate`
- `make test`
- `make check`
- `make stack-up`
- `make stack-down`

Core service targets:

- `make generate`
- `make dev`
- `make run`
- `make test`
- `make build`
- `make containerise`

## 2. Local Runtime Standard

Local development runtime standard is:

- `k3d`
- `kubectl`
- `helm`

Do not introduce docker-compose orchestration unless explicitly agreed.

## 3. AI-Assisted Workflow

Architecture reasoning → ChatGPT  
Repository planning → Cursor Plan  
Implementation → Cursor Agent

## 4. Shared Intent Layer

The shared intent layer for ChatGPT and Cursor is:

- `ENGINEERING.md`
- `DEV.md`
- `.cursorrules`
- `docs/architecture/`

`docs/architecture/` contains the durable architecture context and design intent.

## 5. Working Style for Larger Changes

For non-trivial changes:

1. clarify the problem
2. discuss architecture in ChatGPT
3. produce/update an architecture brief
4. persist the result in the repository
5. let Cursor Plan map it to repo changes
6. let Cursor Agent execute in bounded steps
7. update documentation if needed

## 6. Documentation Discipline

Do not rely on chat history as the source of truth.

Persist important outcomes in repository documents.

## 7. UNDERSTANDING

Please confirm you understand the architecture context and ask any clarifying questions before we start discussing new architecture topics.