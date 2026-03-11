# Service Document: identity

------------------------------------------------------------------------

## 1. Service Identity

- **Service name**: `services/identity/`
- **Service role**: domain
- **Owning domain**: Identity
- **Owning team**: platform

------------------------------------------------------------------------

## 2. Responsibilities

### Primary responsibilities

- Reduced Proteon platform identity (resolve or create)
- External identity linkage (mapping customer-side user identity to a
  Proteon platform identity)
- Auth exchange endpoint for customer backends
- Token issuance semantics (JWT creation, signing, claims, TTL)
- Identity domain API (user profile lookup, account state)
- Future refresh and revocation decisions (deferred, owned here when needed)

### Non-responsibilities (out of scope)

- End-user authentication (owned by the customer platform)
- JWT validation at the edge (owned by api-gateway)
- Request routing or edge protection
- Business rules for other domains (social, matchmaking, sessions)

Reference documents:

- `system/SERVICE_TYPES.md`
- `system/DOMAIN_BOUNDARIES.md`
- `product/PRODUCT_CONTEXT.md`
- `01_PRINCIPLES.md`

------------------------------------------------------------------------

## 3. Architecture Shape

- **Service type**: domain
- **Internal layering**: `adapters -> application -> domain`

Domain layer contains the reduced identity model, identity linkage rules,
and token issuance invariants. Application layer contains use cases for
auth exchange, identity resolution, and profile queries. Adapters contain
HTTP transport, persistence, and generated server code.

### Canonical structure

    services/identity/

      api/
        openapi.yml
        oapi-codegen.server.yml
        oapi-codegen.client.yml

      internal/
        adapters/
          db/
          http/
            generated/server/
        application/
          dto/
          interfaces/
          services/
        domain/
          model/
          rules/
        platform/

      cmd/identity/main.go

------------------------------------------------------------------------

## 4. Integration Contracts

### 4.1 HTTP APIs

- Exposes an HTTP API: **Yes**
- OpenAPI source of truth: `services/identity/api/openapi.yml`
- Generated server code:
  `services/identity/internal/adapters/http/generated/server/`
- Shared HTTP client artifacts:
  `contracts/http/identity/`

The API surface has two distinct areas:

- **Auth exchange endpoints**: called by customer backends to assert
  external identity and receive a Proteon access JWT.
- **Identity domain endpoints**: called by other Proteon services (via
  `contracts/http/identity/`) or through the gateway for client-facing
  needs. Provides user profile lookup, identity resolution, and account
  state queries.

### 4.2 Events

- Events published (future, not baseline):
  - `identity.user.created`
  - `identity.user.deleted`
  - Contract location: `contracts/events/identity/`

- Events consumed: none at baseline

Reference documents:

- `system/EVENT_MODEL.md`
- `system/INTEGRATION_CONTRACTS.md`

------------------------------------------------------------------------

## 5. Dependencies and Boundaries

### 5.1 Allowed dependencies

- `libs/platform` for logging, configuration, observability, middleware
- Postgres for identity persistence
- Cryptographic libraries for JWT signing

### 5.2 Forbidden coupling

This service must not:

- Import other services directly
- Read or write another service's database
- Place identity business logic or token issuance rules into `libs/platform`
- Expose internal domain types as integration contracts
- Assume knowledge of which edge service is calling

------------------------------------------------------------------------

## 6. Auth Exchange Flow

1. Customer backend calls the identity auth exchange endpoint
2. Identity validates the external identity assertion
3. Identity resolves an existing platform identity or creates a new one
4. Identity issues a short-lived access JWT with minimal claims
5. Identity returns the JWT to the customer backend

The customer backend is responsible for forwarding the JWT to the customer
frontend. Identity does not interact with the end user directly.

------------------------------------------------------------------------

## 7. Operational Expectations

- **Scale characteristics**: must handle auth exchange volume from all
  integrating customers; identity lookups from downstream services
- **Latency sensitivity**: moderate for auth exchange (not on every
  request path); identity lookups should be fast
- **Persistence**: Postgres — owns identity data exclusively
- **Helm chart**: `infra/k8s/charts/identity/` (when created)
