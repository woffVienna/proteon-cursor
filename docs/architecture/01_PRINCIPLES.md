# Proteon Architecture Principles

This document defines the stable architectural guardrails of the Proteon
monorepo.

These principles are intended to be durable and should change rarely.

------------------------------------------------------------------------

# 1. Independent Services in a Monorepo

Proteon is a monorepo containing multiple independent services.

Each service:

- lives under `services/<service>/`
- is its own Go module
- owns its own domain behaviour
- must not import another service directly
- communicates with other services only through HTTP APIs or events

Cross-service imports are forbidden.

------------------------------------------------------------------------

# 2. Explicit Ownership

Each service owns:

- its domain logic
- its application orchestration
- its persistence
- its HTTP API semantics
- the events it publishes

Ownership must remain explicit.

Avoid hidden coupling through:

- shared database access
- leaking internal types
- placing service-owned behaviour into shared libraries

------------------------------------------------------------------------

# 3. Strict Internal Layering

Each service follows:

`adapters -> application -> domain`

Interpretation:

- adapters contain transport, persistence, messaging, and integrations
- application contains use cases, orchestration, interfaces, and DTOs
- domain contains pure business rules and invariants

Rules:

- never reverse dependency direction
- keep framework concerns out of domain
- keep orchestration out of domain
- keep transport and persistence concerns in adapters

------------------------------------------------------------------------

# 4. Shared Code Boundaries

Shared technical code lives in:

`libs/platform`

Allowed responsibilities:

- logging
- configuration loading and orchestration
- middleware
- observability helpers
- shared technical primitives

Forbidden responsibilities:

- business logic
- service-specific logic
- service-owned contracts
- cross-service orchestration
- imports from services

------------------------------------------------------------------------

# 5. Explicit Integration Boundaries

Service integration happens through explicit boundaries only:

- HTTP APIs
- asynchronous events

Avoid hidden coupling through:

- direct service imports
- database-level integration
- non-contractual payload assumptions
- transport leakage into domain logic

Reusable integration artifacts belong in `contracts/`, not in
`libs/platform`.

------------------------------------------------------------------------

# 6. Startup-Resolved Configuration

Configuration is resolved once at startup.

Rules:

- shared loading/orchestration belongs in `libs/platform`
- service-specific typed config belongs inside the service
- runtime mutation is not allowed
- invalid configuration must fail fast

------------------------------------------------------------------------

# 7. Incremental Evolution

Prefer:

- bounded changes
- incremental migration
- preserving established patterns unless intentionally changed
- minimal diffs where possible

Avoid:

- broad refactors without a clear plan
- introducing new patterns without need
- silent architectural drift

------------------------------------------------------------------------

# 8. Documentation Discipline

Architecture intent must be documented in `docs/architecture/`.

Do not treat chat history as the source of truth.

When a change materially affects architecture, update the relevant document
and, where useful, capture the decision as a brief or ADR.
