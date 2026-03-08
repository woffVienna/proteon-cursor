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

Because contracts are shared across service boundaries, changes to them
must be governed carefully.

The goals of contract governance are:

- preserve explicit ownership
- make contract changes visible
- prevent accidental breaking changes
- support independent service evolution
- reduce integration drift across the monorepo

------------------------------------------------------------------------

# 2. Scope

This document applies to:

- service OpenAPI specifications
- generated shared HTTP client artifacts
- canonical event schemas
- versioned contract artifacts under `contracts/`

This document does not redefine the basic structure of `contracts/`.
That is defined in `INTEGRATION_CONTRACTS.md`.

------------------------------------------------------------------------

# 3. Ownership Model

## 3.1 Service Ownership

A service owns the semantics of its contracts.

This includes:

- its HTTP API definition
- the events it publishes
- the meaning of its request, response, and event payloads

Ownership remains with the service even when generated or shared
artifacts are published into `contracts/`.

## 3.2 Source of Truth

For HTTP APIs, the source of truth is:

```
services/<service>/api/openapi.yml
```

For events, the source of truth is the canonical event contract and its
producing service’s defined semantics, published under:

```
contracts/events/<service-or-domain>/
```

Consumers do not own producer contracts.

------------------------------------------------------------------------

# 4. Change Visibility

Contract changes must be treated as explicit interface changes.

A contract change must never be hidden inside unrelated implementation
work.

Changes to the following must be visible and reviewable:

- OpenAPI specs
- generated shared HTTP client artifacts
- event schemas
- event version changes
- deprecations
- removals

Contract changes should be easy to identify in pull requests and
architecture discussions.

------------------------------------------------------------------------

# 5. Versioning Principles

Versioning exists to manage change safely across service boundaries.

General principles:

- compatible additive changes are preferred
- incompatible changes must be explicit
- existing consumers must not be broken silently
- contracts must not be repurposed under the same version

------------------------------------------------------------------------

# 6. HTTP Contract Governance

## 6.1 Preferred Change Style

For HTTP APIs, prefer additive changes where possible.

Examples of typically compatible changes:

- adding optional fields to responses
- adding new endpoints
- adding optional request parameters where safe

Examples of potentially breaking changes:

- removing fields
- renaming fields
- changing field meaning
- changing required request structure
- changing response shape incompatibly
- removing endpoints

Breaking changes must be deliberate and visible.

## 6.2 Generated Client Artifacts

Generated shared client artifacts under:

```
contracts/http/<service>/
```

must be regenerated from the service-owned OpenAPI specification.

Shared client artifacts must not be manually edited.

The generated output is a derived contract artifact, not an independent
source of truth.

------------------------------------------------------------------------

# 7. Event Contract Governance

## 7.1 Event Compatibility

Event contracts should be evolved conservatively.

Preferred compatible changes may include:

- adding optional payload fields
- extending metadata in a backward-compatible way

Breaking changes include:

- removing required fields
- changing field meaning incompatibly
- changing payload structure incompatibly
- reusing an existing event name for a different semantic meaning

## 7.2 Explicit Event Versioning

When an event changes incompatibly, publish a new version.

Example:

```
user-created.v1.json
user-created.v2.json
```

Do not silently replace the meaning of an existing event version.

If compatibility matters for multiple consumers, old and new versions may
need to coexist during a transition period.

------------------------------------------------------------------------

# 8. Breaking Change Policy

Breaking changes are allowed only when they are:

- intentional
- reviewed
- visible
- coordinated with affected consumers where necessary

Breaking changes must not be introduced accidentally through routine
refactoring.

When a breaking change is required, the preferred approach is one of:

- introduce a new version
- add a parallel endpoint or event version
- support a bounded migration window
- remove the old version only after explicit migration

------------------------------------------------------------------------

# 9. Deprecation Policy

Deprecation should be explicit and time-bounded where relevant.

A deprecation should identify:

- what is deprecated
- what replaces it
- whether migration is required
- when removal is expected, if known

Deprecation allows consumers to move in a controlled way rather than
absorbing sudden breakage.

------------------------------------------------------------------------

# 10. Consumer Expectations

Consumers are responsible for integrating against explicit contracts, not
producer internals.

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

# 11. Review Expectations

Contract changes should receive a higher level of scrutiny than purely
internal implementation changes.

Reviewers should assess:

- ownership clarity
- backward compatibility
- downstream impact
- versioning correctness
- whether a new version is required
- whether deprecation should be introduced

For non-trivial contract changes, the architectural intent should be
reflected in `docs/architecture/` where appropriate.

------------------------------------------------------------------------

# 12. Generation Discipline

Contract generation must remain deterministic and standardized.

Rules:

- generate from the canonical source of truth
- do not hand-edit generated artifacts
- keep generation wired into normal workflows
- ensure generated outputs stay aligned with source contracts

This protects the repository from contract drift.

------------------------------------------------------------------------

# 13. Governance Boundaries

`libs/platform` must not become a home for service-owned contracts.

Contract governance applies to:

```
services/<service>/api/openapi.yml
contracts/http/...
contracts/events/...
```

It does not change the rule that `libs/platform` is reserved for
technical cross-cutting concerns only.

------------------------------------------------------------------------

# 14. Summary

Proteon standardizes the following contract governance rules:

- services own the meaning of their contracts
- HTTP specs and event schemas are explicit shared boundaries
- contract changes must be visible and reviewable
- additive evolution is preferred
- breaking changes must be deliberate
- incompatible event changes require versioning
- generated artifacts must be derived, not hand-maintained
- consumers depend on contracts, not service internals

These rules are intended to keep service integration explicit, safe, and
evolvable as the platform grows.