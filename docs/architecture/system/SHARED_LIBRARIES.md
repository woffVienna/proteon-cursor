# Shared Libraries

This document defines the role and limits of shared libraries in Proteon.

------------------------------------------------------------------------

# 1. Shared Technical Module

Shared technical code lives in:

`libs/platform`

This module exists to provide common technical building blocks across
services.

------------------------------------------------------------------------

# 2. Allowed Responsibilities

Examples of allowed shared concerns include:

- logging abstractions
- configuration loading and orchestration
- error primitives
- observability helpers
- middleware
- technical utilities that do not own domain behaviour

------------------------------------------------------------------------

# 3. Forbidden Responsibilities

Examples of forbidden contents include:

- business logic
- service-specific use cases
- domain ownership that belongs to a service
- service-owned HTTP clients or event schemas
- cross-service orchestration that should live in a dedicated service
- imports from service modules

------------------------------------------------------------------------

# 4. Design Intent

The goal of `libs/platform` is technical reuse without collapsing service
boundaries.

It should reduce duplication, not centralize business behaviour.

Integration contracts belong in `contracts/`, not in `libs/platform`.

------------------------------------------------------------------------

# 5. Decision Rule

When deciding whether code belongs in `libs/platform`, ask:

- is it technical rather than business logic?
- is it reusable across multiple services?
- does it avoid taking ownership of a service concern?
- can it exist without importing any service?

If not, it likely belongs inside a service instead.
