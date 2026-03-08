# Event Operations

This document defines the operational rules for event delivery,
consumption, failure handling, and observability in Proteon.

It complements `EVENT_MODEL.md` by focusing on runtime behaviour and
operational expectations rather than only architectural principles.

------------------------------------------------------------------------

# 1. Purpose

Proteon uses events for asynchronous communication between independent
services.

This document standardizes how events are operated in practice so that
services behave consistently under normal load, retries, failures, and
reprocessing scenarios.

The goals are:

- predictable event delivery behaviour
- safe consumer processing
- operational consistency across services
- explicit handling of retries and dead-letter scenarios
- strong observability for event-driven flows

------------------------------------------------------------------------

# 2. Scope

This document applies to:

- event publication
- event consumption
- JetStream subject usage
- retries
- dead-letter handling
- replay considerations
- monitoring and observability

This document does not redefine event ownership or event contract
placement. Those concerns are defined in:

- `EVENT_MODEL.md`
- `INTEGRATION_CONTRACTS.md`

------------------------------------------------------------------------

# 3. Subject Conventions

Proteon uses domain-oriented JetStream subjects.

Subjects must follow the canonical event naming convention:

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

- use lowercase segments
- use dot-separated naming
- represent meaningful domain facts
- avoid infrastructure-oriented or generic names

Avoid subjects such as:

```
events
messages
service-events
user-updates
```

These obscure meaning and reduce routing clarity.

------------------------------------------------------------------------

# 4. Stream Layout

JetStream streams should group related subjects by bounded domain area.

Examples:

```
identity.>
matchmaking.>
events.>
```

Guidance:

- prefer domain-bounded streams over one global catch-all stream
- keep stream grouping aligned with service or domain ownership
- avoid overly broad stream definitions unless explicitly justified

This supports:

- clearer operational ownership
- simpler debugging
- bounded replay and retention strategy
- reduced accidental coupling

------------------------------------------------------------------------

# 5. Delivery Semantics

Proteon must assume **at-least-once delivery**.

Implications:

- messages may be delivered more than once
- retries may occur after consumer failure
- redelivery may happen after timeout or acknowledgement problems

Consumers must never assume exactly-once processing.

Exactly-once semantics must not be assumed in service design.

------------------------------------------------------------------------

# 6. Consumer Processing Rules

Consumers must be designed for safe repeated handling.

Minimum expectations:

- handlers must tolerate duplicate delivery
- processing must be idempotent where relevant
- event handling must remain bounded to the consumer’s own responsibility
- handlers must avoid side effects that cannot be safely retried

Consumer logic must not rely on:

- exclusive single delivery
- global ordering across unrelated events
- producer internals
- direct access to producer-owned persistence

------------------------------------------------------------------------

# 7. Acknowledgement Strategy

A consumer should acknowledge an event only after its required processing
has completed successfully.

Guidance:

- do not acknowledge before durable consumer-side state changes complete
- do not acknowledge before required follow-up actions are safely recorded
- if processing fails, allow retry according to broker behaviour

This reduces message loss caused by premature acknowledgement.

------------------------------------------------------------------------

# 8. Idempotency

Idempotency is a required design concern for event consumers.

Consumers should use one or more of the following strategies:

- persist processed event IDs
- detect duplicate business actions
- use natural business keys where appropriate
- check whether target state already exists before reapplying logic

Examples:

```
if eventId already processed -> skip
```

or

```
if target entity already exists -> skip creation
```

or

```
if state transition already applied -> do not apply again
```

Idempotency must be owned by the consuming service.

------------------------------------------------------------------------

# 9. Ordering

Strict ordering must not be assumed across unrelated event types.

Consumers must tolerate:

- delayed events
- retried events
- occasional out-of-order delivery
- parallel consumer execution where configured

If a workflow depends on strict ordering for a specific entity or event
family, that requirement must be designed explicitly and documented.

Ordering-sensitive logic should be the exception, not the default.

------------------------------------------------------------------------

# 10. Retry Strategy

Retries are the normal recovery mechanism for transient failures.

Retryable failures may include:

- temporary downstream unavailability
- transient database errors
- short-lived network failures
- temporary resource contention

Rules:

- retry handling must be safe under duplicate execution
- retries must not assume the previous attempt left no partial effects
- retry behaviour should be observable through metrics and logs

Consumers should distinguish between:

- transient failures worth retrying
- terminal failures that require dead-letter handling

------------------------------------------------------------------------

# 11. Dead Letter Handling

Events that cannot be processed successfully after repeated attempts
must be moved to a Dead Letter Queue (DLQ) or equivalent dead-letter
subject.

Dead-letter handling exists to:

- isolate poison messages
- avoid infinite retry loops
- preserve failing events for investigation
- enable manual or controlled replay

Rules:

- dead-lettered events must not be silently dropped
- dead-letter paths must be monitorable
- dead-lettered events must remain inspectable
- replay must be deliberate and controlled

Operational ownership for DLQ inspection and replay should be explicit.

------------------------------------------------------------------------

# 12. Replay

Replay is useful for recovery, backfill, or consumer bootstrap scenarios.

Replay must be treated as a first-class operational concern.

Consumers must therefore be safe under replay conditions:

- duplicate processing must remain safe
- old events must not corrupt current state
- replay should be scoped to known subject or stream boundaries
- replay procedures should be deliberate and auditable

Large-scale replay should not be performed casually in production.

------------------------------------------------------------------------

# 13. Observability

Event-driven systems require explicit observability.

Minimum recommended signals:

- publish success/failure counters
- consumer processing success/failure counters
- retry counters
- dead-letter counts
- consumer lag or backlog visibility
- processing latency metrics

Logs should include enough context to trace event handling across
services.

Recommended event metadata for tracing:

- event ID
- event type
- producer
- consumer
- timestamp
- correlation or causation identifiers where available

------------------------------------------------------------------------

# 14. Operational Logging

Consumers and producers should log key lifecycle stages consistently.

Typical useful log points:

- event published
- event received
- event processing succeeded
- event processing failed
- event dead-lettered
- event replay started or completed

Logs must not leak sensitive business or personal data.

------------------------------------------------------------------------

# 15. Failure Classification

Event handling failures should be treated as either:

## 15.1 Transient Failures

Examples:

- temporary dependency outage
- network timeout
- temporary lock contention

Expected action:

```
retry
```

## 15.2 Terminal Failures

Examples:

- invalid payload for known schema
- unsupported event version
- permanently inconsistent business preconditions
- unrecoverable consumer logic issue

Expected action:

```
dead-letter
investigate
correct and replay if appropriate
```

------------------------------------------------------------------------

# 16. Operational Ownership

Each consuming service owns:

- its consumer logic
- its idempotency guarantees
- its retry safety
- its dead-letter handling procedures
- its monitoring and alerting for event consumption

Each producing service owns:

- valid event publication
- event schema correctness
- producer-side observability
- maintaining the meaning of published events

------------------------------------------------------------------------

# 17. Summary

Proteon standardizes the following event operations rules:

- use domain-oriented JetStream subjects
- assume at-least-once delivery
- design consumers for idempotent safe retry
- avoid relying on strict ordering by default
- dead-letter repeatedly failing events
- make replay deliberate and safe
- require monitoring for publication, consumption, retries, and DLQ state

These rules are intentionally conservative to keep the event system
reliable as more services and consumers are introduced.