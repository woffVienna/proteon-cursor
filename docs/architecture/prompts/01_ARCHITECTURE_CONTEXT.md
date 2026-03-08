# 01 -- Proteon Architecture Context

This document is the **evolved architecture context** for the Proteon
platform.

It builds upon the original **00 -- Proteon Architecture Context** and
incorporates the additional architecture intent captured in:

-   INTEGRATION_CONTRACTS.md
-   EVENT_MODEL.md
-   EVENT_OPERATIONS.md
-   API_GATEWAY.md
-   CONTRACT_GOVERNANCE.md
-   DOMAIN_BOUNDARIES.md

This file serves as a **single entry-point architecture overview** for
humans and AI tools working on the repository.

Detailed behaviour and rules remain documented in the individual files
under:

docs/architecture/

------------------------------------------------------------------------

# 1. Proteon System Overview

Proteon is a **microservice platform built in a Go monorepo**.

Key characteristics:

-   independent microservices
-   strict service ownership
-   explicit integration boundaries
-   event-driven communication
-   HTTP APIs for synchronous interactions
-   Kubernetes-based runtime environment

Primary technologies:

-   Go
-   OpenAPI
-   Postgres
-   NATS JetStream
-   Kubernetes (k3d locally)
-   Helm

------------------------------------------------------------------------

# 2. Repository Structure

The repository is organized into several top-level areas.

    services/
    libs/
    contracts/
    infra/
    tools/
    docs/

Responsibilities:

  Area                Purpose
  ------------------- ---------------------------------
  services            independent microservices
  libs/platform       shared technical utilities
  contracts           shared HTTP and event contracts
  infra               infrastructure and deployment
  tools               developer utilities
  docs/architecture   architecture documentation

------------------------------------------------------------------------

# 3. Core Architecture Constraints

Proteon follows strict architectural rules.

-   Each service lives under `services/<service>/`
-   Each service is its own Go module
-   Cross-service imports are forbidden
-   Services communicate only via HTTP APIs or events
-   Shared technical code lives in `libs/platform`
-   `libs/platform` must not contain business logic
-   `libs/platform` must not import services

These constraints preserve **clear service boundaries** and prevent the
monorepo from turning into a distributed monolith.

------------------------------------------------------------------------

# 4. Internal Service Architecture

Each service follows a strict layering model:

    adapters → application → domain

Layer responsibilities:

Adapters

-   HTTP servers
-   database integrations
-   messaging integrations
-   external systems

Application

-   use cases
-   orchestration logic
-   service interfaces
-   DTOs

Domain

-   pure business rules
-   invariants
-   domain concepts

Rules:

-   dependency direction must never be reversed
-   domain must remain framework independent

------------------------------------------------------------------------

# 5. Service Communication

Services interact using two mechanisms.

## HTTP APIs

Used for synchronous interactions.

APIs are defined using OpenAPI:

    services/<service>/api/openapi.yml

Generated artifacts include:

-   service server code
-   shared client contracts

Shared HTTP clients are generated into:

    contracts/http/<service>/

Consumers must import **contracts**, never service internals.

## Events

Used for asynchronous communication.

Events represent domain facts such as:

    identity.user.created
    matchmaking.match.created
    events.session.started

Event schemas live in:

    contracts/events/<domain>/

------------------------------------------------------------------------

# 6. Event Architecture

Proteon uses **event choreography**.

Services publish domain events and other services react independently.

Characteristics:

-   no central orchestration by default
-   multiple consumers may react
-   asynchronous workflows remain loosely coupled

Event delivery assumes:

-   at-least-once delivery
-   retryable processing
-   idempotent consumers

Operational behaviour (retry, DLQ, replay, monitoring) is defined in:

EVENT_OPERATIONS.md

------------------------------------------------------------------------

# 7. Integration Contracts

Proteon uses explicit shared contracts.

Shared artifacts live under:

    contracts/
      http/
      events/

Principles:

-   services own their contracts
-   consumers depend on contracts only
-   contract models are integration-layer types
-   domain models remain service-local

Contract evolution rules are defined in:

CONTRACT_GOVERNANCE.md

------------------------------------------------------------------------

# 8. API Gateway

Proteon is expected to evolve toward the topology:

    client → api-gateway → services

The API gateway is an **edge service** responsible for:

-   authentication entry checks
-   request routing
-   rate limiting
-   external API exposure
-   limited response aggregation

The gateway must **not contain core domain logic**.

Domain behaviour always belongs in domain services.

Detailed responsibilities are defined in:

API_GATEWAY.md

------------------------------------------------------------------------

# 9. Service Roles

Services typically fall into one of several categories.

Domain Services

-   own business logic
-   own persistence
-   publish domain events

Examples:

-   identity
-   matchmaking
-   events

Worker Services

-   consume events
-   perform background tasks

Examples:

-   analytics processors
-   notification services

Edge Services

-   expose platform entry points

Examples:

-   api-gateway

Service boundary rules are defined in:

DOMAIN_BOUNDARIES.md

------------------------------------------------------------------------

# 10. Configuration Model

Configuration is resolved **once at service startup**.

Sources:

-   environment variables
-   service configuration files
-   shared configuration utilities

Rules:

-   configuration is immutable at runtime
-   invalid configuration must fail fast
-   shared configuration helpers live in `libs/platform`

------------------------------------------------------------------------

# 11. Observability

Event-driven platforms require strong observability.

Recommended signals:

-   HTTP request metrics
-   event publication metrics
-   consumer lag metrics
-   retry and DLQ metrics
-   service latency metrics

Event metadata should support cross-service tracing.

------------------------------------------------------------------------

# 12. Engineering Workflow

Proteon uses an AI-assisted workflow.

Architecture reasoning → ChatGPT

Repository planning → Cursor Plan

Implementation → Cursor Agent

Architecture intent must always be persisted under:

    docs/architecture/

Chat history must **never be treated as the source of truth**.

------------------------------------------------------------------------

# 13. Architectural Principles

Proteon architecture is based on:

-   independent services
-   explicit service boundaries
-   strict layering
-   explicit integration contracts
-   event-driven communication
-   incremental architectural evolution

These principles ensure the system remains:

-   loosely coupled
-   scalable
-   operationally predictable
-   evolvable over time

------------------------------------------------------------------------

# 14. Relationship to Detailed Documents

This document provides the **architecture overview**.

Detailed behaviour is defined in:

-   INTEGRATION_CONTRACTS.md
-   EVENT_MODEL.md
-   EVENT_OPERATIONS.md
-   API_GATEWAY.md
-   CONTRACT_GOVERNANCE.md
-   DOMAIN_BOUNDARIES.md

Those documents represent the **durable architecture rules** for the
Proteon platform.
