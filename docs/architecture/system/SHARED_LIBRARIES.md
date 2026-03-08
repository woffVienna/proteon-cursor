# Shared Libraries

This document defines the role of shared libraries in Proteon.

## Shared Technical Module

Shared technical code lives in:

`libs/platform`

This module exists to provide common technical building blocks across services.

## Allowed Responsibilities

Examples of allowed shared concerns:

- logging abstractions
- configuration loading/orchestration
- error primitives
- observability helpers
- middleware
- technical utilities that do not own domain behavior

## Forbidden Responsibilities

Examples of forbidden contents:

- business logic
- service-specific use cases
- domain ownership that belongs to a service
- cross-service orchestration that should live in a dedicated service
- imports from service modules

## Design Intent

The goal is reuse of technical concerns without collapsing service boundaries.

`libs/platform` should reduce duplication, not centralize business behavior.

## Decision Rule

When deciding whether code belongs in `libs/platform`, ask:

- is it technical rather than business logic?
- is it reusable across multiple services?
- does it avoid service ownership?
- can it exist without importing any service?

If not, it likely belongs inside a service instead.