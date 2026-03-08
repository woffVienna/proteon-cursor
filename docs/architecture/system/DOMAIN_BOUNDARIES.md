# Domain Boundaries

This document defines service ownership boundaries in Proteon.

------------------------------------------------------------------------

# 1. Primary Boundary

The primary boundary in Proteon is the service boundary.

Each service owns:

- its domain logic
- its application orchestration
- its persistence
- its HTTP API semantics
- the events it publishes

------------------------------------------------------------------------

# 2. Forbidden Coupling

Services must not couple through:

- direct code imports
- shared service-owned database access
- leaking internal implementation types
- hidden runtime assumptions
- treating another service’s contracts as shared domain models

------------------------------------------------------------------------

# 3. Allowed Integration Styles

Services integrate through explicit boundaries only:

- HTTP APIs
- asynchronous events

These boundaries must remain contractual and visible.

------------------------------------------------------------------------

# 4. Internal Layer Boundaries

Within a service:

- adapters depend on application
- application depends on domain
- domain depends on nothing from outer layers

Boundary discipline applies both across services and inside services.

------------------------------------------------------------------------

# 5. Boundary Ownership Rule

A service is the authoritative owner of its domain behaviour.

A service must not:

- read or write another service’s data store directly
- centralize another service’s business rules
- rely on another service’s internal package structure
- hide coupling in shared libraries

------------------------------------------------------------------------

# 6. Boundary Change Rule

If a proposed change alters service ownership, service responsibilities, or
cross-service dependency direction, it should be treated as an architectural
change and documented explicitly.

Use an architecture brief or ADR when the change is material.
