# API Gateway

This document defines the intended role, responsibilities, and limits of
the API gateway in Proteon.

It exists to prevent the gateway from becoming a hidden orchestration layer
or a business-logic-heavy god service.

The platform has two edge gateways: **api-gateway** (player-facing;
customer platform traffic) and **backoffice-gateway** (backoffice
traffic: operators and tenant users). This document describes the
gateway role in general. Backoffice-specific behaviour (app-key for
auth routes, JWT for rest) is in `services/backoffice-gateway.md`.

------------------------------------------------------------------------

# 1. Purpose

Proteon is expected to evolve toward an external interaction model of:

`client -> api-gateway -> services`

(or, for backoffice: `backoffice app -> backoffice-gateway -> services`).

The API gateway provides a controlled external boundary for clients while
preserving explicit service ownership behind that boundary.

------------------------------------------------------------------------

# 2. Position in the Architecture

The API gateway is an edge service.

It sits at the external boundary of the platform and forwards or composes
requests toward internal services.

The gateway is not a domain owner.

------------------------------------------------------------------------

# 3. Core Responsibilities

The gateway may own the following responsibilities:

- public HTTP entry points
- client-facing API routing
- request normalization
- JWT validation and token parsing
- claim extraction and verified identity context forwarding
- coarse-grained route-level access checks
- rate limiting and edge protection
- selective API composition where it simplifies client interaction

Fine-grained business authorization remains the responsibility of the
relevant domain service.

------------------------------------------------------------------------

# 4. Forbidden Responsibilities

The gateway must not own:

- core business rules
- domain persistence for other services
- generic multi-step workflow orchestration by default
- hidden cross-service dependency paths
- direct exposure of unstable internal service APIs without deliberate design

------------------------------------------------------------------------

# 5. External vs Internal APIs

The gateway may expose external APIs that differ from internal
service-to-service APIs.

This is desirable when it helps preserve clear boundaries.

Implications:

- internal service contracts remain owned by services
- external gateway-facing contracts may be curated for client use
- the gateway may translate between external and internal models
- external API stability does not require leaking internal topology

------------------------------------------------------------------------

# 6. Relationship to Identity

The public auth baseline establishes the following ownership split:

- gateway validates JWTs and performs request admission checks
- gateway extracts claims and forwards verified identity context to
  downstream services via trusted internal headers
- identity service owns identity lifecycle, external identity linkage,
  and token issuance semantics
- downstream services receive verified identity context and call identity
  directly when they need richer identity data

The gateway does not issue tokens, resolve identities, or call the
identity service on the request path for validation. Token validation
is stateless using public key verification.

See `briefs/player-auth-baseline.md` for the player-facing auth decision
and rationale. See `briefs/backoffice-auth-baseline.md` for backoffice
auth. See `services/api-gateway.md`, `services/backoffice-gateway.md`,
`services/identity.md`, and `services/auth.md` for ownership details.

------------------------------------------------------------------------

# 7. Aggregation Guidance

Gateway aggregation should remain limited and deliberate.

Appropriate uses:

- combining a small number of domain reads for one client need
- reducing unnecessary client round trips
- presenting a client-friendly response across a few services

Inappropriate uses:

- embedding business workflows in the gateway
- centralizing domain decision making
- making the gateway the default place for multi-service logic

------------------------------------------------------------------------

# 8. Observability and Failure Behaviour

The gateway must provide strong observability:

- request volume metrics
- route-level latency metrics
- error rate metrics
- upstream dependency metrics
- authentication rejection metrics
- rate-limit metrics
- request tracing and correlation

Retries at the gateway should be conservative and deliberate. Blind retries
can amplify load and create cascading failure patterns.

------------------------------------------------------------------------

# 9. Summary

Proteon standardizes the following gateway intent:

- the API gateway is an edge service
- it owns external entry concerns, not core business logic
- it may perform auth checks, routing, rate limiting, and selective aggregation
- it must not become a domain owner or orchestration god service
- internal services remain authoritative for domain behaviour
