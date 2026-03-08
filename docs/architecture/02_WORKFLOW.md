# Proteon Engineering and AI Workflow

This document describes the intended engineering workflow for development and
AI-assisted work on Proteon.

------------------------------------------------------------------------

# 1. Canonical Development Workflow

Prefer Make targets over ad-hoc commands.

Core repository targets:

- `make setup`
- `make generate`
- `make test`
- `make check`
- `make stack-up`
- `make stack-down`

Typical service targets:

- `make generate`
- `make dev`
- `make run`
- `make test`
- `make build`
- `make containerise`

------------------------------------------------------------------------

# 2. Local Runtime Standard

Local development runtime standard is:

- k3d
- kubectl
- Helm

Do not introduce docker-compose orchestration unless explicitly agreed.

------------------------------------------------------------------------

# 3. AI Workflow Split

Use ChatGPT for:

- architecture reasoning
- design intent
- tradeoff discussion
- option analysis
- architecture briefs

Use Cursor Plan for:

- repository-aware planning
- affected file identification
- mapping architecture to repo changes
- bounded implementation steps

Use Cursor Agent for:

- bounded code changes
- incremental implementation
- minimal diffs
- execution aligned with repository architecture rules

------------------------------------------------------------------------

# 4. Shared Intent Layer

The shared intent layer for humans and AI tools is:

- `ENGINEERING.md`
- `DEV.md`
- `.cursorrules`
- `docs/architecture/`

`docs/architecture/` contains the durable architecture context and design
intent for the platform.

------------------------------------------------------------------------

# 5. Working Style for Larger Changes

For non-trivial changes:

1. clarify the problem
2. discuss architecture first
3. write or update the relevant architecture document
4. create an architecture brief or ADR where needed
5. let Cursor Plan map the design to repository changes
6. let Cursor Agent execute in bounded steps
7. update documentation as part of the change

------------------------------------------------------------------------

# 6. Interaction Rule for Architecture Discussions

For non-trivial architecture or design topics:

- ask clarifying questions first
- make assumptions explicit
- present tradeoffs
- prefer bounded plans over broad refactors

Preferred response structure:

1. Clarifying Questions
2. Architecture Overview
3. Key Components
4. Data or Event Flow
5. Implementation Plan
6. Risks and Tradeoffs

------------------------------------------------------------------------

# 7. Documentation Discipline

Do not rely on chat history as the source of truth.

Persist important outcomes in repository documents.

If implementation diverges from architecture intent, either:

- bring the implementation back into alignment, or
- update the architecture document deliberately
