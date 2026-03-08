# Proteon System Context

Proteon is a Go monorepo built around independent microservices.

## System Type

- microservice platform
- monorepo
- independently deployable services
- explicit integration boundaries

## Technology Context

Current core stack includes:

- Go
- OpenAPI
- Postgres
- NATS with JetStream
- Kubernetes for local/runtime orchestration
- Helm for deployment packaging

## Repository Structure

Key top-level areas:

- `services/`
- `libs/platform`
- `infra/`
- `tools/`
- `docs/`

## Service Model

Each service:

- lives under `services/<service>/`
- has its own Go module
- owns its own domain logic
- exposes explicit integration boundaries
- must not be imported by another service

## Communication Model

Service-to-service interaction happens via:

- HTTP APIs
- asynchronous events

Avoid hidden coupling or implicit dependency paths.

## Shared Technical Layer

Shared technical concerns live in:

`libs/platform`

This layer exists for technical reuse, not business logic centralization.

## Configuration Model

Configuration is resolved at startup.

- shared loading/orchestration in `libs/platform`
- service-specific typed config in the service
- restart required for config changes

## Development Model

Proteon uses:

- Make-based workflows
- OpenAPI-first server generation
- k3d + Helm for local stack orchestration

## AI Usage Model

Proteon uses a split AI workflow:

- ChatGPT for architecture/design intent
- Cursor Plan for repo-specific planning
- Cursor Agent for bounded code execution