# Service Document: api-gateway

------------------------------------------------------------------------

## 1. Service Identity

- **Service name**: `services/api-gateway/`
- **Service role**: edge
- **Owning domain**: Platform Edge
- **Owning team**: platform

------------------------------------------------------------------------

## 2. Responsibilities

### Primary responsibilities

- Public HTTP entry points for client traffic
- JWT validation and token parsing
- Claim extraction and verified identity context forwarding
- Request routing to downstream services
- Coarse route-level access checks
- Rate limiting and edge protection
- Request-level observability (volume, latency, errors, auth rejections)

### Non-responsibilities (out of scope)

- Identity lifecycle or token issuance semantics
- User resolution, profile lookups, or identity linkage
- Core business rules for any domain
- Domain persistence
- Multi-step workflow orchestration
- Hidden cross-service dependency paths

Reference documents:

- `system/API_GATEWAY.md`
- `system/SERVICE_TYPES.md`
- `system/DOMAIN_BOUNDARIES.md`

------------------------------------------------------------------------

## 3. Architecture Shape

- **Service type**: edge
- **Internal layering**: `adapters -> application -> domain`

Domain layer is expected to be minimal or empty for an edge service.
Application layer contains routing orchestration and request admission
logic. Adapters contain HTTP transport, downstream service clients, and
generated server code.

### Canonical structure

    services/api-gateway/

      api/
        openapi.yml
        oapi-codegen.server.yml

      internal/
        adapters/
          http/
            generated/server/
        application/
        domain/
        platform/

      cmd/api-gateway/main.go

------------------------------------------------------------------------

## 4. Integration Contracts

### 4.1 HTTP APIs

- Exposes an HTTP API: **Yes**
- OpenAPI source of truth: `services/api-gateway/api/openapi.yml`
- Generated server code:
  `services/api-gateway/internal/adapters/http/generated/server/`
- This service does not publish a shared HTTP client for other services.
  It is the external entry point, not a service consumed internally.

### 4.2 Events

- Events published: none at baseline
- Events consumed: none at baseline

------------------------------------------------------------------------

## 5. Dependencies and Boundaries

### 5.1 Allowed dependencies

- `libs/platform` for logging, configuration, observability, middleware,
  and technical JWT utilities (token parsing, claim extraction)
- `contracts/http/<service>/` for downstream service clients when the
  gateway needs to forward or compose requests

### 5.2 Forbidden coupling

This service must not:

- Import other services directly
- Read or write another service's database
- Place service-specific logic into `libs/platform`
- Contain identity lifecycle logic, token issuance, or user resolution
- Treat integration contracts as shared domain models
- Accumulate business rules that belong to downstream domain services

------------------------------------------------------------------------

## 6. Auth Behaviour

The gateway performs the following on each authenticated request:

1. Extract JWT from the request (Authorization header)
2. Validate signature, expiry, issuer, and audience
3. Extract claims (platform user ID, tenant ID)
4. Forward verified identity context to downstream services via trusted
   internal headers
5. Reject requests with invalid, expired, or missing tokens

The gateway does not:

- Issue tokens
- Resolve or create identities
- Perform user lookups
- Make runtime calls to the identity service for token validation

Fine-grained business authorization remains the responsibility of the
relevant downstream domain service.

------------------------------------------------------------------------

## 7. Operational Expectations

- **Scale characteristics**: stateless, horizontally scalable
- **Latency sensitivity**: high — sits on the critical path for all
  client requests
- **Observability requirements**: request volume, route-level latency,
  error rates, upstream dependency metrics, auth rejection metrics,
  rate-limit metrics, request tracing and correlation
- **Helm chart**: `infra/k8s/charts/api-gateway/` (when created)
