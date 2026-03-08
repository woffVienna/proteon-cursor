# Service Types

This document defines the allowed service roles in Proteon.

Clearly defining service roles prevents architectural drift, duplicated
responsibilities, and uncontrolled service proliferation.

------------------------------------------------------------------------

# 1. Service Categories

Proteon defines three primary service types:

- edge services
- domain services
- worker services

Each service should fit one of these roles.

------------------------------------------------------------------------

# 2. Edge Services

Edge services are system entry points.

They expose APIs to external clients or upstream systems and translate
external requests into internal service interactions.

Typical responsibilities:

- external HTTP API exposure
- request validation
- authentication and authorization integration
- request routing
- limited aggregation of multiple domain calls
- protocol translation

Characteristics:

- stateless
- horizontally scalable
- minimal domain logic
- strong observability requirements

Examples:

- `api-gateway`
- `public-api`
- `admin-api`

------------------------------------------------------------------------

# 3. Domain Services

Domain services implement core business capabilities.

They contain actual domain logic and own domain data.

Typical responsibilities:

- business rules
- domain entities
- persistence
- domain APIs
- domain event publication

Characteristics:

- clear business ownership
- persistent data ownership
- stable service contracts
- event producers and consumers

Examples:

- `identity-service`
- `matchmaking-service`
- `profile-service`
- `wallet-service`

------------------------------------------------------------------------

# 4. Worker Services

Worker services perform background or asynchronous processing.

They usually react to events or execute scheduled workloads.

Typical responsibilities:

- event-driven processing
- scheduled jobs
- data processing pipelines
- integration with external systems
- heavy or long-running background tasks

Characteristics:

- event consumers
- background processing
- minimal or no public API
- task-focused lifecycle

Examples:

- notification worker
- analytics processor
- backfill processor

------------------------------------------------------------------------

# 5. Dependency Expectations

All service types must respect the same architectural boundaries.

Allowed dependencies include:

- `libs/platform`
- `contracts/http/<service>`
- `contracts/events/...`

Forbidden dependencies include:

- direct imports of other services
- direct database access to another service
- hiding service coupling in shared libraries

Communication between services must still happen through HTTP APIs or events.
