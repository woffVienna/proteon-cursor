# Event Model

This document defines the event-driven communication model used in Proteon.

It describes the principles, responsibilities, and operational rules for
publishing and consuming events between services.

Events are the primary mechanism for asynchronous service integration.

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

Events are **integration contracts**, not shared domain models.

------------------------------------------------------------------------

# 3. Responsibilities

## 3.1 Producer Responsibilities

Producers must:

- emit clear and intentional domain events
- preserve domain ownership
- publish events only after the state change they represent is durable
- maintain event schema compatibility
- avoid leaking internal implementation details

Producers are responsible for the **meaning and lifecycle of the event**.

## 3.2 Consumer Responsibilities

Consumers must:

- process events safely
- assume events may be delivered more than once
- implement idempotent handling where required
- tolerate temporary producer or broker failures
- restrict event handling to the consumer’s own responsibility

Consumers must never assume exclusive ownership of an event.

------------------------------------------------------------------------

# 4. Event Naming Conventions

Events must represent **domain facts**.

Event names follow the structure:

```
<domain>.<entity>.<event>
```

Examples:

```
identity.user.created
identity.user.deleted
matchmaking.match.created
events.session.started
```

Rules:

- use lowercase names
- use dot-separated segments
- use past-tense verbs where possible
- avoid technical or infrastructure-oriented names

Avoid names such as:

```
identity.user.update_event
user-event-stream
internal-user-change
```

These hide domain meaning.

------------------------------------------------------------------------

# 5. Event Envelope

All events should follow a consistent envelope structure.

Example conceptual structure:

```json
{
  "eventId": "uuid",
  "eventType": "identity.user.created",
  "eventVersion": "v1",
  "timestamp": "2026-01-01T12:00:00Z",
  "producer": "identity-service",
  "payload": { }
}
```

Field meanings:

| Field | Description |
|------|-------------|
| eventId | unique identifier for the event instance |
| eventType | canonical event name |
| eventVersion | schema version |
| timestamp | time of event creation |
| producer | service that emitted the event |
| payload | event-specific data |

The payload structure is defined by the **event contract schema**.

------------------------------------------------------------------------

# 6. Event Contracts

Event schemas are stored in the shared contracts directory:

```
contracts/events/<service-or-domain>/
```

Example:

```
contracts/events/identity/user-created.v1.json
```

Event schemas represent the **canonical integration contract** between
producers and consumers.

Rules:

- producers must publish events conforming to the schema
- consumers depend on the schema, not on producer code
- schemas must be versioned when incompatible changes occur

Event schemas must **not live inside services** once they become shared.

------------------------------------------------------------------------

# 7. JetStream Subject Conventions

Proteon uses **NATS JetStream** for event delivery.

Subjects should follow the event naming convention directly.

Example:

```
identity.user.created
matchmaking.match.created
events.session.started
```

JetStream streams may group subjects using wildcards.

Example stream configuration conceptually:

```
identity.>
matchmaking.>
events.>
```

This allows services to subscribe to relevant event families.

Avoid overly generic subjects such as:

```
events
messages
service-events
```

These obscure domain meaning.

------------------------------------------------------------------------

# 8. Delivery Guarantees

Event delivery must assume **at-least-once semantics**.

Implications:

- events may be delivered multiple times
- events may be retried by the broker
- consumers must implement safe processing

Exactly-once delivery must **not be assumed**.

------------------------------------------------------------------------

# 9. Idempotency

Consumers must implement idempotent processing when handling events.

Common strategies include:

- storing processed event IDs
- deduplicating using event keys
- detecting existing domain state before applying changes

Example patterns:

```
if eventId already processed -> ignore
```

or

```
if entity already exists -> skip creation
```

Idempotency logic belongs in the **consumer service**.

------------------------------------------------------------------------

# 10. Ordering

Ordering guarantees should not be assumed across different event types.

Within a single subject or entity stream, ordering may be preserved
depending on JetStream configuration.

Consumers must tolerate:

- delayed events
- out-of-order events
- retries after failure

Business logic should not rely on strict ordering unless explicitly
designed for it.

------------------------------------------------------------------------

# 11. Retry Strategy

Event processing failures should trigger retries.

Typical retry flow:

```
consumer fails
        ↓
broker redelivers
        ↓
consumer retries processing
```

Retries must assume the event may already have been partially processed.

Therefore:

- processing must be idempotent
- retries must be safe

Retry behavior is primarily managed by the message broker.

------------------------------------------------------------------------

# 12. Dead Letter Handling

Events that repeatedly fail processing should be moved to a
Dead Letter Queue (DLQ).

DLQ usage allows:

- inspection of failing events
- debugging consumer logic
- safe isolation of poison messages

DLQ handling must never silently discard events.

Operational procedures should exist for inspecting and replaying
dead-lettered messages.

------------------------------------------------------------------------

# 13. Event Choreography

Proteon defaults to **event choreography**.

Example flow:

```
identity.user.created
        ↓
profile-service reacts
        ↓
notification-service reacts
```

Each service independently reacts to events.

Benefits:

- loose coupling
- scalable fan-out
- service autonomy

Explicit orchestration services should only be introduced when a
workflow requires centralized coordination.

------------------------------------------------------------------------

# 14. When to Use Events vs HTTP

Prefer events when:

- communication is asynchronous
- temporal decoupling is beneficial
- multiple downstream consumers may exist
- workflows should avoid direct request coupling

Prefer HTTP when:

- immediate response is required
- the caller needs the result synchronously
- the operation is request-driven

Both patterns coexist in the system.

------------------------------------------------------------------------

# 15. Observability

Event systems require strong observability.

Recommended capabilities include:

- event emission metrics
- consumer lag monitoring
- retry and failure tracking
- DLQ monitoring

Event IDs and timestamps should allow tracing across services.

------------------------------------------------------------------------

# 16. Future Expansion

This document may evolve to include:

- a producer and consumer catalog
- event schema governance processes
- event compatibility rules
- automated schema validation
- replay tooling
- event debugging guidelines