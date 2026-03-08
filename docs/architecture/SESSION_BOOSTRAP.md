# Proteon Session Bootstrap

Proteon is a Go monorepo containing independent microservices.

## Core Architecture Constraints

- Each service lives under `services/<service>/`
- Each service is its own Go module
- Cross-service imports are forbidden
- Services communicate via HTTP APIs or asynchronous events only
- Shared technical code lives in `libs/platform`
- `libs/platform` must not contain business logic
- `libs/platform` must not import services

## Internal Service Architecture

Each service follows this dependency direction:

`adapters -> application -> domain`

Rules:

- adapters contain transport, persistence, and framework integrations
- application contains use cases, orchestration, DTOs, and interfaces
- domain contains pure business concepts and rules
- dependency direction must never be reversed
- domain must remain pure

## Configuration Model

Configuration is resolved once at startup.

- base configuration comes from environment
- service-specific typed configuration is assembled inside the service
- shared configuration loading/orchestration lives in `libs/platform`
- no live runtime mutation
- changes require restart/redeploy

## Local Development Standard

- local runtime standard is `k3d + Helm`
- prefer Make targets over ad-hoc commands
- do not introduce docker-compose orchestration

## AI Workflow

- ChatGPT is used for architecture reasoning and design intent
- Cursor Plan is used for repo-specific planning
- Cursor Agent is used for bounded execution
- architecture/context documents in `docs/architecture/` are the shared intent layer for both tools

## Interaction Rule

Always ask clarifying questions before proposing a solution for non-trivial architecture or design problems.