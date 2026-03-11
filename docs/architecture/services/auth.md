# Service Document: auth

------------------------------------------------------------------------

## 1. Service Identity

- **Service name**: `services/auth/`
- **Service role**: domain
- **Owning domain**: Authentication (backoffice)
- **Owning team**: platform

------------------------------------------------------------------------

## 2. Responsibilities

### Primary responsibilities

- Authentication methods and flows for backoffice users (operators and
  tenant users; see `GLOSSARY.md`)
- Credential storage (password hashes; later OAuth links, MFA) for
  backoffice users only
- Registration and login endpoints (exposed via backoffice-gateway)
- Exchange with Identity service to obtain JWTs after successful
  authentication
- Future: OAuth, SSO, MFA (owned here as auth methods)

### Non-responsibilities (out of scope)

- User record storage or token issuance (owned by Identity; link is
  `user_id`)
- JWT validation at the edge (owned by backoffice-gateway)
- Player authentication (player identity and auth exchange are the
  tenant's application and Identity concern; see `briefs/player-auth-baseline.md`)
- Request routing or edge protection

Reference documents:

- `system/SERVICE_TYPES.md`
- `system/DOMAIN_BOUNDARIES.md`
- `briefs/backoffice-auth-baseline.md`
- `services/identity.md`
- `GLOSSARY.md`

------------------------------------------------------------------------

## 3. Architecture Shape

- **Service type**: domain
- **Internal layering**: `adapters -> application -> domain`

Domain layer contains authentication method rules and credential
validation invariants. Application layer contains use cases for
registration, login, and token exchange. Adapters contain HTTP
transport, persistence for credentials (keyed by `user_id`), and
client calls to Identity for token issuance.

### Canonical structure

    services/auth/

      api/
        openapi.yml
        oapi-codegen.server.yml
        oapi-codegen.client.yml (if needed for internal callers)

      internal/
        adapters/
          db/
          http/
            generated/server/
        application/
        domain/
        platform/

      cmd/auth/main.go

------------------------------------------------------------------------

## 4. Integration Contracts

### 4.1 HTTP APIs

- Exposes an HTTP API: **Yes**
- OpenAPI source of truth: `services/auth/api/openapi.yml`
- Called by backoffice-gateway (login, register, and any future auth
  method endpoints). Not called directly by the backoffice app;
  traffic goes through the gateway.
- Consumes Identity service: internal API to request token issuance for
  a given `user_id` (after credential validation). Identity holds user
  records only; auth holds credentials. Link between them is `user_id`.

### 4.2 Events

- Events published: none at baseline
- Events consumed: none at baseline

------------------------------------------------------------------------

## 5. Dependencies and Boundaries

### 5.1 Allowed dependencies

- `libs/platform` for logging, configuration, observability
- `contracts/http/identity/` for token issuance calls to Identity
- Postgres (or equivalent) for credential storage, keyed by `user_id`

### 5.2 Forbidden coupling

This service must not:

- Import other services directly
- Read or write another service's database (except its own credential
  store)
- Place auth business logic into `libs/platform`
- Issue tokens (Identity issues; auth requests)
- Own user records (Identity owns; auth stores only credentials linked
  by `user_id`)

------------------------------------------------------------------------

## 6. Operational Expectations

- **Scale characteristics**: must handle login/register volume from
  backoffice; not on the critical path for every backoffice API call
- **Latency sensitivity**: moderate for auth flows
- **Persistence**: owns credential data only; user records in Identity
- **Helm chart**: `infra/k8s/charts/auth/` (when created)
