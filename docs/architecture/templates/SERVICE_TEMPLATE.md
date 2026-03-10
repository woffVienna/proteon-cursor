# Service Template

Use this template when creating or materially changing a service under
`services/<service>/`.

It exists to prevent architectural drift and to make ownership, boundaries,
and integration contracts explicit for each service.

---

## 1. Service Identity

- **Service name**: `services/<service>/`
- **Service role**: edge / domain / worker (see `system/SERVICE_TYPES.md`)
- **Owning domain**: e.g. Identity, Social, Sessions, Matchmaking, Events
- **Owning team**: who is accountable for this service over time

---

## 2. Responsibilities

Describe what this service **owns** and is expected to provide.

- **Primary responsibilities**
  - ...
  - ...

- **Non-responsibilities (out of scope)**
  - ...
  - ...

When in doubt, align with:

- `00_CONTEXT.md`
- `system/DOMAIN_BOUNDARIES.md`
- `system/SERVICE_TYPES.md`
- `product/PRODUCT_CONTEXT.md`

---

## 3. Architecture Shape

- **Service type**: edge / domain / worker
- **Internal layering**: confirm the service follows

  `adapters -> application -> domain`

- **Canonical structure**:

      services/<service>/

        api/
          openapi.yml
          oapi-codegen.server.yml
          oapi-codegen.client.yml (if this service publishes a reusable HTTP client)

        internal/
          adapters/
          application/
          domain/
          platform/ (service-local technical helpers only; shared utilities belong in libs/platform)

        cmd/<service>/main.go

Note any intentional, agreed deviations from the canonical layout.

---

## 4. Integration Contracts

### 4.1 HTTP APIs

- Does this service expose an HTTP API? **Yes / No**
- OpenAPI source of truth: `services/<service>/api/openapi.yml`
- Generated server code location:
  `services/<service>/internal/adapters/http/generated/server/`
- Shared HTTP client artifacts (if any):
  `contracts/http/<service>/`

Relevant architecture documents:

- `system/INTEGRATION_CONTRACTS.md`
- `system/CONTRACT_GOVERNANCE.md`

### 4.2 Events

- Events **published** by this service:
  - `contracts/events/<service-or-domain>/...`

  Examples:
  - `identity.user.created`
  - `matchmaking.match.created`

- Events **consumed** by this service:
  - ...

Reference documents:

- `system/EVENT_MODEL.md`
- `system/EVENT_OPERATIONS.md`

---

## 5. Dependencies and Boundaries

### 5.1 Allowed dependencies

List the key dependencies this service uses and confirm they respect the
architecture rules.

Examples:

- `libs/platform/...`
- `contracts/http/<other-service>/...`
- `contracts/events/...`
- External systems (e.g. LiveKit, payment providers)

Describe how external systems are accessed (HTTP, events, SDKs) and through
which adapters.

All dependencies must respect the service boundary and the internal
dependency direction:

`adapters -> application -> domain`

### 5.2 Forbidden coupling (must not do)

Confirm this service does **not**:

- import other services directly
- read or write another service's database
- place service-specific business logic into `libs/platform`
- treat integration contracts as shared domain models

See:

- `01_PRINCIPLES.md`
- `system/DOMAIN_BOUNDARIES.md`

---

## 6. Operational Expectations

Summarize how this service is expected to behave in production.

- **SLOs / critical expectations** (if defined)
  - ...
- **Scale characteristics**
  - expected QPS
  - fan-out patterns
  - latency sensitivity
- **Deployment**
  - Helm chart location and configuration (if applicable):
    `infra/k8s/charts/<service>/`
  - Document key runtime considerations if relevant.

---

## 7. Checklist for New or Changed Services

Before merging a new service or a material service reshaping:

- [ ] `services/<service>/` created via `make create-service <service>`
- [ ] Internal layout matches `00_CONTEXT.md` and `ENGINEERING.md`
- [ ] `api/openapi.yml` defined (or explicitly not needed for this service type)
- [ ] HTTP server generation wired via `oapi-codegen.server.yml`
- [ ] HTTP client generation (if needed) wired to `contracts/http/<service>/`
- [ ] Event contracts placed under `contracts/events/<service-or-domain>/`
- [ ] Dependencies comply with `01_PRINCIPLES.md` and `system/DOMAIN_BOUNDARIES.md`
- [ ] Any architectural changes are captured in an architecture brief or ADR
      using the provided templates
