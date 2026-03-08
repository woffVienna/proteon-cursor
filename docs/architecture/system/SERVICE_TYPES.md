# Service Types

This document defines the different service roles that may exist in the
Proteon platform.

Clearly defining service roles prevents uncontrolled microservice growth,
architectural drift, and unclear responsibilities.

Each service in the platform should fall into one of the defined categories.

------------------------------------------------------------------------

# 1. Purpose

Proteon is built as a distributed system composed of independent services.

Without clear service roles, systems tend to drift toward:

- unclear boundaries
- duplicated responsibilities
- accidental coupling
- infrastructure misuse
- uncontrolled service proliferation

To avoid this, Proteon defines explicit **service types**.

Each service type has:

- a clear responsibility
- expected dependency boundaries
- typical runtime characteristics
- a predictable lifecycle

These definitions guide both architectural decisions and service design.

------------------------------------------------------------------------

# 2. Service Types

Proteon defines three primary service types:

- Edge Services
- Domain Services
- Worker Services

Each service in the system should fit one of these roles.

------------------------------------------------------------------------

# 3. Edge Services

Edge services represent **system entry points**.

They expose APIs to external clients or upstream systems and translate
external requests into internal service interactions.

Edge services typically do **not contain core business logic**.

Instead they coordinate requests toward domain services.

## Responsibilities

Edge services typically handle:

- HTTP API exposure
- request validation
- authentication and authorization integration
- request routing
- aggregation of multiple domain calls
- protocol translation (HTTP → events or internal APIs)

They act as **system boundaries**.

## Characteristics

Typical characteristics:

- stateless
- horizontally scalable
- minimal domain logic
- high request throughput
- strong observability requirements

Edge services should remain thin.

Complex business logic must reside in domain services.

## Allowed Dependencies

Edge services may depend on:

```
contracts/http/<service>
contracts/events/...
libs/platform
```

Edge services must **never import domain services directly**.

All communication must occur via:

- HTTP APIs
- asynchronous events

## Typical Examples

Examples of edge services include:

```
api-gateway
public-api
admin-api
webhook-receiver
```

------------------------------------------------------------------------

# 4. Domain Services

Domain services implement the **core business capabilities of the system**.

They contain the actual domain logic and maintain ownership over
their respective domain data.

Domain services represent the **primary building blocks** of the platform.

## Responsibilities

Domain services are responsible for:

- implementing business rules
- managing domain entities
- owning domain data stores
- exposing domain APIs
- publishing domain events

Domain services must preserve **clear domain ownership**.

They should not take responsibility for unrelated domains.

## Characteristics

Typical characteristics:

- encapsulated business logic
- persistent data ownership
- stable domain APIs
- event producers
- event consumers

Domain services should remain **independent and autonomous**.

## Allowed Dependencies

Domain services may depend on:

```
libs/platform
contracts/http/<service>
contracts/events/...
```

Domain services must **not depend on other services directly**.

Communication must occur via:

- HTTP contracts
- events

Domain services must never access another service’s database.

## Typical Examples

Examples of domain services include:

```
identity-service
matchmaking-service
profile-service
wallet-service
session-service
```

------------------------------------------------------------------------

# 5. Worker Services

Worker services perform **background or asynchronous processing**.

They typically react to events or execute scheduled workloads.

Workers usually do not expose external APIs.

## Responsibilities

Worker services typically handle:

- event-driven processing
- background workflows
- scheduled jobs
- data processing pipelines
- integration with external systems
- heavy compute workloads

Workers allow domain services to remain responsive and focused.

## Characteristics

Typical characteristics:

- event consumers
- background processing
- potentially long-running tasks
- minimal or no public API
- often horizontally scalable

Workers should remain **task-focused**.

They should not evolve into full domain services.

## Allowed Dependencies

Worker services may depend on:

```
contracts/events/...
contracts/http/<service>
libs/platform
```

Workers must never depend directly on service code.

Communication with other services must occur via:

- events
- HTTP APIs

## Typical Examples

Examples include:

```
notification-worker
analytics-worker
email-delivery-worker
reconciliation-worker
reporting-worker
```

------------------------------------------------------------------------

# 6. Dependency Expectations

Regardless of service type, the following dependency rules apply.

Forbidden dependencies:

- importing another service’s internal code
- accessing another service’s database
- placing business logic inside shared libraries

Allowed dependencies:

```
libs/platform
contracts/http/<service>
contracts/events/...
```

Communication between services must occur through explicit integration
contracts.

Service types may differ in **how they use these contracts**, but the
dependency rules remain consistent.

------------------------------------------------------------------------

# 7. Lifecycle Expectations

Different service types tend to have different lifecycle characteristics.

## Edge Services

Lifecycle expectations:

- stable external API surface
- backward compatibility requirements
- frequent scaling
- careful change management

## Domain Services

Lifecycle expectations:

- long-lived ownership of a business capability
- stable domain model evolution
- event version management
- strong data ownership boundaries

Domain services often become **core platform components**.

## Worker Services

Lifecycle expectations:

- may be temporary or evolving
- often added to support new workflows
- can be replaced or consolidated over time
- easier to scale horizontally

Worker services should remain **specialized and replaceable**.

------------------------------------------------------------------------

# 8. Examples

Example repository structure illustrating different service types:

```
services/
  api-gateway/            (edge service)
  public-api/             (edge service)

  identity-service/       (domain service)
  matchmaking-service/    (domain service)
  profile-service/        (domain service)

  notification-worker/    (worker service)
  analytics-worker/       (worker service)
  reconciliation-worker/  (worker service)
```

Each service still follows the standard internal architecture:

```
adapters → application → domain
```

Service type determines **system role**, not internal architecture.

------------------------------------------------------------------------

# 9. Decision Summary

Proteon standardizes the following service roles:

- **Edge Services**  
  system entry points that expose APIs and route requests

- **Domain Services**  
  core business services that own domain logic and data

- **Worker Services**  
  background processors handling asynchronous workloads

All services must:

- respect integration boundaries
- communicate via HTTP or events
- avoid cross-service code dependencies
- maintain clear domain ownership

These rules ensure a scalable and maintainable service architecture.

------------------------------------------------------------------------

# 10. Consequences

Benefits:

- prevents uncontrolled microservice proliferation
- enforces architectural clarity
- simplifies onboarding for new services
- provides predictable service responsibilities
- improves long-term maintainability

Costs:

- requires architectural discipline
- requires careful service classification
- encourages thoughtful service creation

These costs are intentional and preferable to architectural drift.