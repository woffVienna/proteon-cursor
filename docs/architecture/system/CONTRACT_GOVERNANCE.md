# Contract Governance

This document defines how integration contracts are owned, evolved,
versioned, and changed in Proteon.

It complements `INTEGRATION_CONTRACTS.md` by focusing on governance and
change discipline rather than only placement and usage.

------------------------------------------------------------------------

# 1. Purpose

Proteon uses explicit integration contracts for communication between
independent services.

These contracts include:

- HTTP API contracts
- event contracts

Because contracts are shared across service boundaries, changes to them must
be governed carefully.

------------------------------------------------------------------------

# 2. Ownership Model

A service owns the semantics of its contracts.

This includes:

- its HTTP API definition
- the events it publishes
- the meaning of its request, response, and event payloads

For HTTP APIs, the source of truth is:

`services/<service>/api/openapi.yml`

For events, the source of truth is the canonical event contract and its
producing service’s semantics, published under:

`contracts/events/<service-or-domain>/`

Consumers do not own producer contracts.

------------------------------------------------------------------------

# 3. Change Visibility

Contract changes must be treated as explicit interface changes.

A contract change must never be hidden inside unrelated implementation work.

Changes to the following must be visible and reviewable:

- OpenAPI specs
- generated shared HTTP client artifacts
- event schemas
- event version changes
- deprecations
- removals

------------------------------------------------------------------------

# 4. Versioning Principles

General principles:

- compatible additive changes are preferred
- incompatible changes must be explicit
- existing consumers must not be broken silently
- contracts must not be repurposed under the same version

------------------------------------------------------------------------

# 5. HTTP Contract Governance

Prefer additive changes where possible.

Typically compatible changes:

- adding optional fields to responses
- adding new endpoints
- adding optional request parameters where safe

Potentially breaking changes:

- removing fields
- renaming fields
- changing field meaning
- changing required request structure
- changing response shape incompatibly
- removing endpoints

Generated shared client artifacts under `contracts/http/<service>/` must be
regenerated from the service-owned OpenAPI specification and must not be
edited by hand.

------------------------------------------------------------------------

# 6. Event Contract Governance

Event contracts should be evolved conservatively.

Preferred compatible changes may include:

- adding optional payload fields
- extending metadata in a backward-compatible way

Breaking changes include:

- removing required fields
- changing field meaning incompatibly
- changing payload structure incompatibly
- reusing an existing event name for a different semantic meaning

When an event changes incompatibly, publish a new version.

Example:

- `user-created.v1.json`
- `user-created.v2.json`

Do not silently replace the meaning of an existing event version.

------------------------------------------------------------------------

# 7. Breaking Change Policy

Breaking changes are allowed only when they are:

- intentional
- reviewed
- visible
- coordinated with affected consumers where necessary

Preferred approaches:

- introduce a new version
- add a parallel endpoint or event version
- support a bounded migration window
- remove the old version only after explicit migration

------------------------------------------------------------------------

# 8. Deprecation Policy

Deprecation should be explicit and time-bounded where relevant.

A deprecation should identify:

- what is deprecated
- what replaces it
- whether migration is required
- when removal is expected, if known

------------------------------------------------------------------------

# 9. Consumer Expectations

Consumers must:

- depend on published contracts only
- avoid assumptions beyond the contract
- map contract types into local models
- tolerate additive evolution where reasonable
- plan migration when a contract is deprecated or versioned

Consumers must not:

- import service internals
- rely on producer database structure
- infer unsupported semantics from incidental payload shape

------------------------------------------------------------------------

# 10. Summary

Proteon standardizes the following contract governance rules:

- services own the meaning of their contracts
- HTTP specs and event schemas are explicit shared boundaries
- contract changes must be visible and reviewable
- additive evolution is preferred
- breaking changes must be deliberate
- incompatible event changes require versioning
- generated artifacts must be derived, not hand-maintained
- consumers depend on contracts, not service internals
