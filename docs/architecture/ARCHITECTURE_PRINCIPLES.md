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

`adapters -> application -> domain`

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