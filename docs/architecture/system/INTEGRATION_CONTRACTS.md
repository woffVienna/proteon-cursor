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

Services must **never import each other directly**. Integration between
services must occur through explicit boundaries only:

- HTTP APIs
- asynchronous events

To support this cleanly, Proteon uses a dedicated top-level directory:

```
contracts/
```

This directory contains reusable **integration artifacts** derived from
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

```
services/
```

Services contain:

- domain logic
- application use cases
- adapters
- API definitions
- event producers

Services **own the meaning of their APIs and events**.

## 2.2 Platform Library

```
libs/platform
```

`libs/platform` is reserved for **technical cross-cutting concerns**.

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

```
contracts/
```

The `contracts/` directory contains **shared integration artifacts**
derived from service-owned contracts.

These artifacts exist only to allow other services to integrate with
a service **without importing its code**.

------------------------------------------------------------------------

# 3. Contract Categories

The `contracts/` directory contains two main categories.

```
contracts/
  http/
  events/
```

## 3.1 HTTP Contracts

```
contracts/http/<service>/
```

Contains **generated client artifacts** derived from a service’s
OpenAPI specification.

Typical contents include:

- generated HTTP clients
- generated request/response models
- supporting generated types

These artifacts allow other services to call the service without
importing its internal code.

## 3.2 Event Contracts

```
contracts/events/<service-or-domain>/
```

Contains canonical **event schemas** and optionally generated bindings
for event payloads.

Typical contents include:

- versioned event schemas
- optional generated Go bindings
- schema metadata

These artifacts allow services to publish and consume events using a
shared, explicit contract.

------------------------------------------------------------------------

# 4. HTTP Contract Model

Proteon uses an **OpenAPI-first model** for HTTP APIs.

## 4.1 API Ownership

Each service owns its API specification:

```
services/<service>/api/openapi.yml
```

This file is the **source of truth** for the service HTTP interface.

## 4.2 Generated Server Code

Server code generated from the OpenAPI specification belongs inside
the service adapters layer.

Example:

```
services/<service>/internal/adapters/http/generated/server/openapi.gen.go
```

This keeps transport concerns isolated in adapters.

## 4.3 Generated Shared Client Contracts

Client artifacts for other services are generated into:

```
contracts/http/<service>/
```

These artifacts are generated from the service OpenAPI specification.

They allow consuming services to call the API **without importing
service code**.

## 4.4 Consumption Rules

Consuming services may import:

```
contracts/http/<service>
```

But must never import:

```
services/<service>/...
```

Generated HTTP clients may only be used inside the **adapters layer**
of the consuming service.

## 4.5 Mapping Rule

Generated request/response models are **integration contract types**.

They must **not be treated as domain models**.

Services must map integration types into their own application or
domain types where appropriate.

------------------------------------------------------------------------

# 5. Event Contract Model

Proteon uses **domain events** and **event choreography** as the default
model for asynchronous integration.

## 5.1 Event Naming

Events must represent meaningful domain facts.

Examples:

```
identity.user.created
identity.user.deleted
matchmaking.match.created
```

Avoid generic stream names that hide domain meaning.

## 5.2 Canonical Event Schemas

Reusable event schemas are stored under:

```
contracts/events/<service-or-domain>/
```

Example:

```
contracts/events/identity/user-created.v1.json
```

Schemas may be defined using JSON Schema or another agreed format.

Generated language bindings may be added if useful.

## 5.3 Producer Rule

The producing service publishes events according to the canonical schema.

The producing service **owns the event semantics and version lifecycle**.

## 5.4 Consumer Rule

Consumers depend on the explicit event contract artifact in `contracts/`.

Consumers must **not import producer internal code**.

Event payloads must be mapped into service-local application or domain
models before use.

Event payloads are **integration contracts, not shared domain models**.

------------------------------------------------------------------------

# 6. Event Choreography

Proteon defaults to **event choreography** for cross-service workflows.

Example flow:

```
identity.user.created
        ↓
profile-service reacts
        ↓
notification-service reacts
```

Each service reacts independently.

There is **no central coordinator**.

This model supports:

- loose coupling
- independent service evolution
- simpler service boundaries

Explicit orchestration services may be introduced only where a
business workflow requires central coordination.

------------------------------------------------------------------------

# 7. Dependency Rules

## 7.1 Forbidden Dependencies

The following are prohibited:

- importing one service from another service
- storing service-owned contracts inside `libs/platform`
- treating generated integration models as domain models
- leaking transport types into domain
- accessing another service’s database

## 7.2 Allowed Dependencies

Allowed examples include:

- services importing `contracts/http/<service>`
- services importing event bindings from `contracts/events/...`
- services importing `libs/platform` for technical concerns

## 7.3 Layering Rules

Within a service, contract usage must respect service layering:

```
adapters → application → domain
```

Implications:

- HTTP clients live in adapters
- event consumers/producers live in adapters
- mapping into application/domain models happens before business logic
- domain remains framework and contract independent

------------------------------------------------------------------------

# 8. Generation Model

Integration contracts are generated as part of service development
workflows.

Generation runs from the service `make generate` process.

## 8.1 HTTP Generation

From:

```
services/<service>/api/openapi.yml
```

Generate:

```
services/<service>/internal/adapters/http/generated/server/
contracts/http/<service>/
```

## 8.2 Event Generation

Event schemas and optional generated bindings are placed under:

```
contracts/events/<service-or-domain>/
```

------------------------------------------------------------------------

# 9. Versioning

## 9.1 HTTP APIs

Changes to a service OpenAPI specification represent **interface changes**.

Breaking changes must be introduced deliberately.

## 9.2 Event Contracts

Event schemas must be versioned when incompatible changes occur.

Example:

```
user-created.v1.json
user-created.v2.json
```

Existing versions must not be silently modified.

------------------------------------------------------------------------

# 10. API Gateway Consideration

Proteon is expected to evolve toward the topology:

```
client → api-gateway → services
```

This does not change the integration contract model.

Implications:

- services still own their APIs and events
- `contracts/` remains the internal integration layer
- external gateway APIs may introduce additional external contracts
- the gateway must not replace explicit service contracts

------------------------------------------------------------------------

# 11. Example Repository Layout

```
services/
  identity/
    api/
      openapi.yml
    internal/
      adapters/
        http/
          generated/
            server/
              openapi.gen.go

contracts/
  http/
    identity/
      client.gen.go
      models.gen.go
  events/
    identity/
      user-created.v1.json
      user-deleted.v1.json
```

In this structure:

- the service owns the API and event semantics
- generated server code remains inside the service
- shared HTTP client contracts live under `contracts/http/`
- shared event schemas live under `contracts/events/`

------------------------------------------------------------------------

# 12. Decision Summary

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

------------------------------------------------------------------------

# 13. Consequences

## Benefits

- explicit service boundaries
- strong architectural isolation
- reusable contracts without coupling services
- consistent model for HTTP and events
- monorepo scale without architectural drift

## Costs

- contract generation workflows must be maintained
- contract versioning requires discipline
- services must perform explicit model mapping

These costs are intentional and preferable to hidden coupling.