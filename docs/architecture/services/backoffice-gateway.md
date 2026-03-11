# Service Document: backoffice-gateway

------------------------------------------------------------------------

## 1. Service Identity

- **Service name**: `services/backoffice-gateway/`
- **Service role**: edge
- **Owning domain**: Platform Edge
- **Owning team**: platform

------------------------------------------------------------------------

## 2. Responsibilities

### Primary responsibilities

- Single HTTP entry point for backoffice traffic (control-plane and
  tenant self-service)
- App-key validation for unauthenticated auth routes (login, register)
- JWT validation and claim extraction for all other backoffice routes
- Request routing to auth service and downstream backoffice APIs
- Coarse route-level access checks
- Rate limiting and edge protection
- Request-level observability

### Non-responsibilities (out of scope)

- Identity lifecycle or token issuance
- Credential storage or authentication method logic (owned by auth service)
- Core business rules for any domain
- Domain persistence
- Hidden cross-service dependency paths

Reference documents:

- `system/API_GATEWAY.md`
- `system/SERVICE_TYPES.md`
- `briefs/backoffice-auth-baseline.md`
- `GLOSSARY.md`

------------------------------------------------------------------------

## 3. Architecture Shape

- **Service type**: edge
- **Internal layering**: `adapters -> application -> domain`

Domain layer is minimal or empty. Adapters contain HTTP transport,
reverse proxy, app-key middleware for auth routes, JWT middleware for
secured routes, and routing to auth and downstream services.

### Canonical structure

    services/backoffice-gateway/

      api/
        openapi.yml

      internal/
        adapters/
          http/
            middleware/
            proxy/
        application/
        domain/
        platform/

      cmd/backoffice-gateway/main.go

------------------------------------------------------------------------

## 4. Integration Contracts

### 4.1 HTTP APIs

- Exposes an HTTP API: **Yes**
- OpenAPI source of truth: `services/backoffice-gateway/api/openapi.yml`
- Auth routes (e.g. login, register) are exposed on the gateway but
  protected by app-key only (no JWT). All other routes require JWT.
- This service does not publish a shared HTTP client; it is the
  external entry point for the backoffice.

### 4.2 Events

- Events published: none at baseline
- Events consumed: none at baseline

------------------------------------------------------------------------

## 5. Dependencies and Boundaries

### 5.1 Allowed dependencies

- `libs/platform` for logging, configuration, observability, middleware,
  and technical JWT utilities
- `contracts/http/<service>/` for downstream service clients (auth,
  identity, and other backoffice backends)

### 5.2 Forbidden coupling

This service must not:

- Import other services directly
- Read or write another service's database
- Place service-specific logic into `libs/platform`
- Contain credential validation, token issuance, or user resolution
- Treat integration contracts as shared domain models
- Accumulate business rules that belong to downstream services

------------------------------------------------------------------------

## 6. Auth Behaviour

- **Auth routes (login, register):** Require a valid app-key from the
  backoffice app. No JWT. Gateway proxies to auth service.
- **All other routes:** Require a valid JWT (issued by Identity).
  Gateway validates JWT, extracts claims (e.g. user ID, tenant ID,
  subject type for operator vs tenant user), forwards verified context
  to downstream services.

The gateway does not issue tokens or resolve identities. Fine-grained
authorization remains the responsibility of downstream services.

------------------------------------------------------------------------

## 7. Operational Expectations

- **Scale characteristics**: stateless, horizontally scalable
- **Latency sensitivity**: moderate — backoffice traffic is not
  player-facing
- **Observability**: request volume, route-level latency, error rates,
  auth rejection metrics, rate-limit metrics, request tracing
- **Helm chart**: `infra/k8s/charts/backoffice-gateway/` (when created)
