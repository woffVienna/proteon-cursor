# Proteon Architecture Context

Use `docs/architecture/00_CONTEXT.md` as the primary entry point for
understanding the Proteon architecture.

Canonical architecture documents live under:

```
docs/architecture/
```

The recommended reading order is:

- `00_CONTEXT.md`
- `01_PRINCIPLES.md`
- `02_WORKFLOW.md`
- `03_INDEX.md`

Topic-specific architecture rules are defined under:

```
docs/architecture/system/
```

Key documents include:

- `DOMAIN_BOUNDARIES.md`
- `SHARED_LIBRARIES.md`
- `INTEGRATION_CONTRACTS.md`
- `CONTRACT_GOVERNANCE.md`
- `EVENT_MODEL.md`
- `EVENT_OPERATIONS.md`
- `SERVICE_TYPES.md`
- `API_GATEWAY.md`

------------------------------------------------------------------------

# Core Architecture Constraints

The following constraints define the structural foundation of Proteon:

- each service lives under `services/<service>/`
- each service is its own Go module
- cross-service imports are forbidden
- services communicate via HTTP APIs or asynchronous events only
- shared technical code lives in `libs/platform`
- reusable integration artifacts live in `contracts/`
- services follow the internal layering:

```
adapters → application → domain
```

Additional principles:

- services own their domain logic and persistence
- contracts define cross-service integration boundaries
- event choreography is the default asynchronous interaction model
- configuration is resolved once at service startup
- architectural evolution should remain incremental and bounded

------------------------------------------------------------------------

# Working Style for Architecture Discussions

When discussing non-trivial architecture topics:

1. Ask clarifying questions first.
2. Make assumptions explicit.
3. Preserve explicit service boundaries.
4. Prefer incremental evolution over large refactors.
5. Reference the canonical documents instead of restating them.

Use the topic-specific documents under `docs/architecture/system/`
when deeper rules or operational guidance are required.

------------------------------------------------------------------------

# Guiding Intent

Proteon aims to maintain a platform architecture that is:

- loosely coupled
- service-owned
- contract-driven
- event-oriented
- operationally predictable
- evolvable over time

Architecture documentation in `docs/architecture/` is the **durable
source of truth** for these principles.