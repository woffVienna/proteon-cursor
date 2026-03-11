# Architecture Brief: Public Auth Baseline

## Status

Accepted. Implemented: identity auth exchange and api-gateway with JWT
validation and reverse proxy to identity.

## Goal

Establish the baseline public authentication model for Proteon and define
the ownership split between api-gateway and identity.

## Context

Proteon is a backend platform (PaaS), not a consumer-facing product. It
provides platform capabilities such as identity, social, matchmaking,
sessions, and real-time event distribution to integrating applications.

The API gateway is defined as an edge service and must not become a hidden
orchestration or business-logic layer. Identity is a domain service that
owns identity lifecycle and token issuance semantics.

Reference documents:

- `product/PRODUCT_CONTEXT.md`
- `system/API_GATEWAY.md`
- `system/SERVICE_TYPES.md`

## Constraints

- Services are independent and must not import each other directly
- Integration must happen through HTTP APIs or events
- Gateway remains an edge service
- Identity remains the domain owner for identity semantics
- Shared libraries must not contain service-specific business logic
- Architecture intent must be captured in repository documents

Reference documents:

- `01_PRINCIPLES.md`
- `00_CONTEXT.md`
- `system/INTEGRATION_CONTRACTS.md`

## Clarifications

- The customer platform authenticates the end user
- The customer backend performs the Proteon auth exchange
- Proteon issues a short-lived pure access JWT
- The customer frontend uses that JWT against api-gateway
- Realtime and streaming auth are explicitly deferred from this decision

------------------------------------------------------------------------

## Architecture Overview

Proteon uses a customer-authenticated, platform-token model.

The customer backend exchanges or asserts user identity with Proteon. The
identity service resolves or creates a reduced platform identity and issues
a short-lived Proteon access JWT. The customer frontend then calls
api-gateway directly using that JWT. The gateway validates the JWT and
forwards verified identity context to downstream services.

This preserves a clean edge/domain split where:

- api-gateway owns validation, routing, and context forwarding
- identity owns identity lifecycle, linkage, and token issuance
- downstream services receive verified claims and act on them

### Why JWT

JWT is the correct primitive for this topology because:

- **Stateless edge validation.** The gateway validates tokens using a
  public key without calling identity on every request. This preserves
  service independence and avoids making identity a scalability bottleneck.
- **Self-contained claims.** Verified identity context travels with the
  token. The gateway extracts claims and forwards them downstream without
  an introspection call or shared session state.
- **No server-side session state.** Proteon does not own the user session.
  The customer does. A short-lived access JWT is the clean expression of a
  point-in-time platform identity assertion.
- **Industry standard.** This is the established model for platform tokens
  in PaaS and BaaS systems.

Alternatives considered and rejected:

- Opaque tokens with introspection: creates synchronous gateway-to-identity
  coupling on every request. Violates loose coupling.
- Server-side sessions: requires shared state between gateway and identity.
  Violates service independence.
- API keys: not suitable for end-user-scoped access.

### Token Claims

Claims should be minimal:

- Platform user ID
- Tenant or customer ID
- Token expiry
- Issuer and audience

Claims must not carry profile data, permissions lists, or account metadata.
The platform user ID is the correlation key. Downstream services that need
richer identity data call the identity service HTTP API directly.

------------------------------------------------------------------------

## Key Components

- `services/api-gateway/`
- `services/identity/`
- `services/api-gateway/api/openapi.yml`
- `services/identity/api/openapi.yml`
- `contracts/http/identity/` for reusable client artifacts consumed by
  other services
- `contracts/events/identity/` when identity lifecycle events are later
  published
- `libs/platform` for shared technical JWT utilities only (token parsing,
  middleware). No issuance logic or identity business rules.

------------------------------------------------------------------------

## Data / Event Flow

1. User authenticates on the customer platform
2. Customer backend calls Proteon identity auth exchange endpoint
3. Identity resolves or creates a reduced platform identity
4. Identity issues a short-lived access JWT
5. Customer backend forwards the JWT to the customer frontend
6. Frontend calls api-gateway with the JWT
7. Gateway validates the JWT and extracts claims
8. Gateway forwards verified identity context to downstream services
   via trusted internal headers

------------------------------------------------------------------------

## Implementation Plan

1. Define the reduced identity model in the identity service domain
2. Define the auth exchange contract in `services/identity/api/openapi.yml`
3. Define token claims structure and TTL policy
4. Implement identity linkage and token issuance in the identity service
5. Implement JWT validation and context forwarding in the gateway
6. Add observability and failure metrics at both services
7. Document ownership boundaries in service documents

------------------------------------------------------------------------

## Risks and Design Characteristics

### Architectural risk

- **Gateway absorbing identity logic.** The gateway must validate tokens
  and forward claims. It must not accumulate identity checks, user lookups,
  or issuance logic over time. This risk is guarded by the forbidden
  responsibilities in `system/API_GATEWAY.md` and the service documents.

### Known design characteristics

- **Customer session and Proteon token are independent.** The customer
  session may outlive the Proteon token. When the token expires, the
  customer backend re-exchanges. This is working as designed. TTL is
  the control mechanism.
- **Client-side token exposure is intentional.** JWT is a bearer token
  designed to be carried by the client. The security model relies on
  short TTL, HTTPS-only transport, minimal claims, and audience/issuer
  validation.

### Deferred scope

- Token refresh and revocation semantics
- Realtime and streaming authentication
- Admin or backoffice access surfaces

------------------------------------------------------------------------

## Documentation Impact

- Update `system/API_GATEWAY.md` to reflect the chosen public auth
  baseline (section 6 becomes definitive)
- Add `services/api-gateway.md` as a service document
- Add `services/identity.md` as a service document
- Update `03_INDEX.md` with new document entries
