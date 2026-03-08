# Proteon Integration Contracts

This document defines how reusable integration contracts are modeled,
generated, and consumed across the Proteon monorepo.

It complements the core architecture documents by specifying the rules for:

- HTTP client contract generation
- event schema ownership and sharing
- integration artifact placement
- allowed dependency directions for contract consumption

------------------------------------------------------------------------

# 1. Purpose

Proteon is a monorepo containing independent services.

Services must never import each other directly. Integration between
services must occur through explicit boundaries only:

- HTTP APIs
- asynchronous events

To support this cleanly, Proteon uses a dedicated top-level directory:

`contracts/`

This directory contains reusable integration artifacts derived from
service-owned contracts.

This approach ensures:

- explicit integration boundaries
- no cross-service imports
- no duplicated integration models across services
- no architectural drift into hidden coupling

------------------------------------------------------------------------

# 2. Separation of Responsibilities

Proteon uses three clearly separated layers of code ownership.

## 2.1 Services

`services/`

Services contain:

- domain logic
- application use cases
- adapters
- API definitions
- event producers

Services own the meaning of their APIs and events.

## 2.2 Platform Library

`libs/platform`

`libs/platform` is reserved for technical cross-cutting concerns.

Allowed responsibilities include:

- logging
- configuration loading
- observability
- middleware
- shared technical helpers

Forbidden responsibilities include:

- business logic
- service-owned HTTP clients
- service-owned event schemas
- cross-service orchestration logic

## 2.3 Contracts

`contracts/`

The `contracts/` directory contains shared integration artifacts derived from
service-owned contracts.

These artifacts exist only to allow other services to integrate with a
service without importing its code.

------------------------------------------------------------------------

# 3. Contract Categories

The `contracts/` directory contains two main categories.

`contracts/http/`
`contracts/events/`

## 3.1 HTTP Contracts

`contracts/http/<service>/`

Contains generated client artifacts derived from a service’s OpenAPI
specification.

Typical contents include:

- generated HTTP clients
- generated request and response models
- supporting generated types

## 3.2 Event Contracts

`contracts/events/<service-or-domain>/`

Contains canonical event schemas and optionally generated bindings for event
payloads.

Typical contents include:

- versioned event schemas
- optional generated Go bindings
- schema metadata

------------------------------------------------------------------------

# 4. HTTP Contract Flow

Each service owns its API specification at:

`services/<service>/api/openapi.yml`

Generated server code belongs inside the service adapters layer.

Example:

`services/<service>/internal/adapters/http/generated/server/openapi.gen.go`

Generated shared client artifacts for other services are produced under:

`contracts/http/<service>/`

Consumption rules:

- consuming services may import `contracts/http/<service>`
- consuming services must never import `services/<service>/...`
- generated HTTP clients may only be used in adapters
- integration models must be mapped into service-local models before use

------------------------------------------------------------------------

# 5. Event Contract Flow

Proteon uses domain events and event choreography as the default model for
asynchronous integration.

Event names should represent meaningful domain facts, for example:

- `identity.user.created`
- `identity.user.deleted`
- `matchmaking.match.created`

Canonical reusable event contracts live under:

`contracts/events/<service-or-domain>/`

Producers publish according to the canonical schema.

Consumers depend on explicit event contract artifacts, not on producer
internal packages.

Event payloads are integration contracts, not shared domain models.

------------------------------------------------------------------------

# 6. Dependency Rules

Forbidden:

- importing one service from another service
- storing service-owned contracts inside `libs/platform`
- treating generated integration models as domain models
- leaking transport types into domain
- accessing another service’s database

Allowed:

- services importing `contracts/http/<service>`
- services importing event bindings from `contracts/events/...`
- services importing `libs/platform` for technical concerns

Within a service, contract usage must respect:

`adapters -> application -> domain`

------------------------------------------------------------------------

# 7. Generation Model

Integration contracts are generated as part of service development workflows.

Generation runs from the service `make generate` process.

From:

`services/<service>/api/openapi.yml`

Generate:

- `services/<service>/internal/adapters/http/generated/server/`
- `contracts/http/<service>/`

Event schemas and optional generated bindings belong under:

`contracts/events/<service-or-domain>/`

------------------------------------------------------------------------

# 8. Summary

Proteon standardizes the following rules:

- use a dedicated top-level `contracts/` directory
- keep `libs/platform` limited to technical cross-cutting concerns
- derive shared HTTP client artifacts from service OpenAPI specifications
- keep generated server code inside the producing service
- store canonical event schemas under `contracts/events/`
- use domain events for asynchronous communication
- default to event choreography across services
- prohibit cross-service imports
- treat integration models as adapter-layer types, not domain models
