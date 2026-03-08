# Event Operations

This document defines the operational rules for event delivery, consumption,
failure handling, and observability in Proteon.

It complements `EVENT_MODEL.md` by focusing on runtime behaviour and
operational expectations.

------------------------------------------------------------------------

# 1. Subject Conventions

Proteon uses domain-oriented JetStream subjects.

Subjects follow the canonical event naming convention:

`<domain>.<entity>.<event>`

Examples:

- `identity.user.created`
- `identity.user.deleted`
- `matchmaking.match.created`
- `events.session.started`

Avoid generic subjects that obscure domain meaning.

------------------------------------------------------------------------

# 2. Stream Layout

JetStream streams should group related subjects by bounded domain area.

Examples:

- `identity.>`
- `matchmaking.>`
- `events.>`

Prefer domain-bounded streams over one global catch-all stream.

------------------------------------------------------------------------

# 3. Delivery Semantics

Proteon assumes at-least-once delivery.

Implications:

- messages may be delivered more than once
- retries may occur after consumer failure
- redelivery may happen after timeout or acknowledgement problems

Consumers must never assume exactly-once processing.

------------------------------------------------------------------------

# 4. Consumer Processing Rules

Consumers must be designed for safe repeated handling.

Minimum expectations:

- handlers tolerate duplicate delivery
- processing is idempotent where relevant
- event handling remains bounded to the consumer’s own responsibility
- handlers avoid side effects that cannot be safely retried

------------------------------------------------------------------------

# 5. Acknowledgement Strategy

A consumer should acknowledge an event only after required processing has
completed successfully.

Do not acknowledge before:

- durable consumer-side state changes complete
- required follow-up actions are safely recorded

If processing fails, allow retry according to broker behaviour.

------------------------------------------------------------------------

# 6. Idempotency

Idempotency is a required design concern for event consumers.

Common strategies include:

- persisting processed event IDs
- detecting duplicate business actions
- using natural business keys where appropriate
- checking whether target state already exists before reapplying logic

Idempotency belongs to the consuming service.

------------------------------------------------------------------------

# 7. Ordering

Strict ordering must not be assumed across unrelated event types.

Consumers must tolerate:

- delayed events
- retried events
- occasional out-of-order delivery
- parallel consumer execution where configured

If strict ordering is required for a specific entity or workflow, that
requirement must be explicit and documented.

------------------------------------------------------------------------

# 8. Retry Strategy

Retries are the normal recovery mechanism for transient failures.

Retry-safe processing is mandatory.

Consumers should distinguish between:

- transient failures worth retrying
- terminal failures that should be dead-lettered

------------------------------------------------------------------------

# 9. Dead Letter Handling

Events that cannot be processed successfully after repeated attempts must be
moved to a dead-letter queue or equivalent dead-letter subject.

Rules:

- dead-lettered events must not be silently dropped
- dead-letter paths must be monitorable
- dead-lettered events must remain inspectable
- replay must be deliberate and controlled

------------------------------------------------------------------------

# 10. Replay

Replay is useful for recovery, backfill, or consumer bootstrap scenarios.

Consumers must therefore be safe under replay conditions:

- duplicate processing remains safe
- old events do not corrupt current state
- replay is scoped to known subject or stream boundaries

Large-scale replay should not be performed casually in production.

------------------------------------------------------------------------

# 11. Observability

Minimum recommended signals:

- publish success and failure counters
- consumer processing success and failure counters
- retry counters
- dead-letter counts
- consumer lag or backlog visibility
- processing latency metrics

Recommended event metadata for tracing:

- event ID
- event type
- producer
- consumer
- timestamp
- correlation or causation identifiers where available
