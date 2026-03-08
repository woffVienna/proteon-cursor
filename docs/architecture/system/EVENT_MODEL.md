# Event Model

This document defines the event-driven communication model used in Proteon.

It describes the principles, responsibilities, and conceptual rules for
publishing and consuming events between services.

------------------------------------------------------------------------

# 1. Purpose

Events are used for asynchronous communication between independent services.

They support:

- loose coupling between services
- independent service evolution
- scalable fan-out to multiple consumers
- temporal decoupling between producers and consumers

Events must always respect service ownership and explicit integration
boundaries.

------------------------------------------------------------------------

# 2. Principles

The following principles govern event-driven communication in Proteon.

- events represent meaningful domain facts or state transitions
- event producers own the event emission
- consumers react without introducing tight coupling
- event processing must assume duplicate delivery is possible
- consumers must implement idempotent processing where required
- events must never expose internal-only implementation details

Events are integration contracts, not shared domain models.

------------------------------------------------------------------------

# 3. Responsibilities

## 3.1 Producer Responsibilities

Producers must:

- emit clear and intentional domain events
- preserve domain ownership
- publish events only after the represented state change is durable
- maintain event schema compatibility
- avoid leaking internal implementation details

## 3.2 Consumer Responsibilities

Consumers must:

- process events safely
- assume events may be delivered more than once
- implement idempotent handling where required
- tolerate temporary producer or broker failures
- restrict event handling to the consumer’s own responsibility

------------------------------------------------------------------------

# 4. Event Naming

Events must represent meaningful domain facts.

Event names follow:

`<domain>.<entity>.<event>`

Examples:

- `identity.user.created`
- `identity.user.deleted`
- `matchmaking.match.created`
- `events.session.started`

Use lowercase dot-separated names and prefer past-tense verbs where
possible.

------------------------------------------------------------------------

# 5. Event Envelope

All events should follow a consistent conceptual envelope.

Example shape:

- `eventId`
- `eventType`
- `eventVersion`
- `timestamp`
- `producer`
- `payload`

The payload structure is defined by the event contract schema.

------------------------------------------------------------------------

# 6. Event Contracts

Event schemas are stored in:

`contracts/events/<service-or-domain>/`

Example:

`contracts/events/identity/user-created.v1.json`

These schemas are the canonical integration contracts between producers and
consumers.

Rules:

- producers must publish events conforming to the schema
- consumers depend on the schema, not producer code
- schemas must be versioned when incompatible changes occur

------------------------------------------------------------------------

# 7. Choreography Default

Proteon defaults to event choreography for cross-service asynchronous
workflows.

This means:

- one service emits a domain event
- other interested services consume and react independently
- no central coordinator is introduced unless explicitly justified

Use explicit orchestration only when a business workflow genuinely requires
central coordination.

------------------------------------------------------------------------

# 8. Events vs HTTP

Prefer events when:

- communication is asynchronous
- temporal decoupling is beneficial
- multiple downstream consumers may exist
- workflows should avoid direct request coupling

Prefer HTTP when:

- immediate response is required
- the caller needs the result synchronously
- the operation is request-driven
