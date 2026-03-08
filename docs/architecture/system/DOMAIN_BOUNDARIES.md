# Domain Boundaries

This document defines the intended domain and ownership boundaries in Proteon.

## Primary Boundary

The primary boundary in Proteon is the service boundary.

Each service owns:

- its domain logic
- its application orchestration
- its persistence rules
- its API/event contract

## Forbidden Coupling

Services must not couple through:

- direct code imports
- shared service-owned database access
- leaking internal implementation types
- hidden runtime assumptions

## Allowed Integration Styles

Services integrate through:

- HTTP APIs
- asynchronous events

These boundaries must remain explicit.

## Internal Layer Boundaries

Within a service:

- adapters depend on application
- application depends on domain
- domain depends on nothing inward from outer layers

## Shared Technical Code Boundary

`libs/platform` may support services technically but does not own domain behavior.

It must not become a place where service/business behavior is centralized.

## Boundary Change Rule

If a proposed change alters service ownership or responsibility, it should be treated as an architectural change and documented explicitly.