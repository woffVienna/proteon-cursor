# API Gateway

This document defines the intended role, responsibilities, and limits of
the future API gateway in Proteon.

It exists to prevent the gateway from becoming a hidden orchestration
layer or a business-logic-heavy god service.

------------------------------------------------------------------------

# 1. Purpose

Proteon is expected to evolve toward an external interaction model of:

```
client → api-gateway → services
```

The API gateway provides a controlled external boundary for clients while
preserving explicit service ownership behind that boundary.

The gateway exists to:

- centralize external entry concerns
- present a stable external access layer
- shield internal services from direct public exposure where appropriate
- support consistent cross-cutting edge behaviour

------------------------------------------------------------------------

# 2. Position in the Architecture

The API gateway is an edge service.

It sits at the external boundary of the platform and forwards or
composes requests toward internal services.

The gateway is **not a domain owner**.

It must not replace domain services or absorb domain logic that belongs
inside them.

------------------------------------------------------------------------

# 3. Core Responsibilities

The gateway may own the following responsibilities.

## 3.1 External Request Entry

The gateway is the main public entry point for client-facing API traffic.

This may include:

- public HTTP endpoints
- client-facing API routing
- external request normalization

## 3.2 Authentication Entry Checks

The gateway may validate authentication credentials or tokens before
forwarding requests into the platform.

Examples:

- JWT validation
- token parsing
- claim extraction
- forwarding verified identity context

## 3.3 Authorization Pre-Checks

The gateway may perform coarse-grained edge authorization where useful.

Examples:

- rejecting unauthenticated access
- applying route-level access checks
- enforcing basic client or tenant scoping

Fine-grained domain authorization remains the responsibility of the
relevant domain service.

## 3.4 Routing

The gateway routes requests to the appropriate internal service.

Routing must remain explicit and understandable.

The gateway must not create hidden dependency paths between services.

## 3.5 Rate Limiting and Edge Protection

The gateway is a suitable place for edge-level protection such as:

- rate limiting
- abuse protection
- request size limits
- basic traffic shaping

## 3.6 API Composition

The gateway may aggregate multiple internal calls where doing so serves
an external client need.

This should be used selectively.

Composition is appropriate when:

- it simplifies the external client contract
- it avoids pushing platform-internal topology to clients
- the aggregation remains shallow and understandable

------------------------------------------------------------------------

# 4. Forbidden Responsibilities

The gateway must not take on the following responsibilities.

## 4.1 Core Domain Logic

The gateway must not own business rules that belong to domain services.

Examples of forbidden behaviour:

- pricing rules
- eligibility decisions
- identity lifecycle rules
- matchmaking rules
- business workflow state

## 4.2 Persistence Ownership

The gateway must not own domain persistence for other services.

It must not become a shared database owner or a hidden read model store
for unrelated domains without an explicit design decision.

## 4.3 Hidden Cross-Service Orchestration

The gateway must not become a generic internal workflow engine.

It must not coordinate multi-step domain processes by default.

If a business workflow requires orchestration, that concern should be
designed explicitly in the appropriate service or workflow component.

## 4.4 Service Internal Leakage

The gateway must not expose service-internal transport details,
implementation quirks, or unstable internal APIs directly to clients
without deliberate contract design.

------------------------------------------------------------------------

# 5. External vs Internal APIs

The gateway may expose external APIs that are different from internal
service-to-service APIs.

This is desirable when it helps preserve clear boundaries.

Implications:

- internal service contracts remain owned by services
- external gateway-facing contracts may be curated for client use
- the gateway may translate between external and internal models
- external API stability does not require leaking internal topology

This separation helps avoid coupling clients directly to internal
service design.

------------------------------------------------------------------------

# 6. Relationship to Services

Domain services remain the authoritative owners of:

- domain logic
- domain persistence
- internal APIs
- domain events

The gateway is a consumer of service APIs, not a replacement for them.

Internal services must still be designed as proper standalone services,
not as thin implementation fragments of the gateway.

------------------------------------------------------------------------

# 7. Relationship to Identity

Identity handling is expected to evolve over time.

For now, the gateway should be designed to support an identity service
without assuming the gateway itself becomes the owner of identity
business rules.

A reasonable direction is:

- gateway performs token validation and request admission checks
- identity service owns identity lifecycle and token issuance semantics
- downstream services receive verified identity context as needed

This keeps concerns separated while allowing a clean external entry
layer.

------------------------------------------------------------------------

# 8. Request Propagation

When forwarding requests, the gateway may propagate contextual metadata
required by downstream services.

Examples may include:

- authenticated subject identifier
- tenant identifier
- request correlation ID
- tracing headers

The gateway must propagate only intentional context.

It must not create implicit domain coupling through arbitrary header
forwarding or internal transport leakage.

------------------------------------------------------------------------

# 9. Aggregation Guidance

Gateway aggregation should remain limited and deliberate.

Appropriate examples:

- combining a small number of domain reads for one client screen
- presenting a client-friendly response assembled from multiple services
- reducing unnecessary client round trips

Inappropriate examples:

- embedding business workflows in the gateway
- building broad orchestration trees
- centralizing domain decision making
- making the gateway the default place for multi-service logic

If gateway aggregation becomes deep, stateful, or business-heavy, the
design should be reconsidered.

------------------------------------------------------------------------

# 10. Observability

The gateway is a critical edge component and must provide strong
observability.

Recommended capabilities:

- request volume metrics
- route-level latency metrics
- error rate metrics
- upstream dependency metrics
- authentication rejection metrics
- rate-limit metrics
- request tracing and correlation

Because the gateway sits at the platform boundary, it is often the first
place to detect systemic issues.

------------------------------------------------------------------------

# 11. Failure Behaviour

Gateway failure behaviour should be explicit.

The gateway should:

- fail clearly when required upstream dependencies are unavailable
- avoid masking domain errors incorrectly
- preserve useful error semantics where appropriate
- avoid introducing hidden retry storms toward internal services

Retries at the gateway should be conservative and deliberate.

Blind retries can amplify load and create cascading failure patterns.

------------------------------------------------------------------------

# 12. Evolution Guidance

The gateway should evolve incrementally.

Recommended approach:

- begin with routing and basic edge concerns
- add authentication checks and context propagation
- add selective aggregation only where justified
- introduce external contract shaping deliberately

Avoid starting with an overly ambitious gateway that tries to own every
cross-cutting and business concern at once.

------------------------------------------------------------------------

# 13. Summary

Proteon standardizes the following gateway intent:

- the API gateway is an **edge service**
- it owns **external entry concerns**, not core business logic
- it may perform **auth checks, routing, rate limiting, and selective aggregation**
- it must not become a **domain owner or orchestration god service**
- internal services remain authoritative for domain behaviour
- external contracts may differ from internal service contracts

These rules keep the gateway useful without allowing it to erode service
boundaries.