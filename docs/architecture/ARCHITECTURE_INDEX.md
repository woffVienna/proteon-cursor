# Proteon Architecture Index

If you are an AI assistant reading this repository, start with `SESSION_BOOTSTRAP.md`.

This document is the entry point for understanding the Proteon platform architecture.

It provides a navigation map to the architecture documents and describes the
overall system structure.

Proteon is implemented as a Go monorepo containing independent microservices.

Key architecture constraints:

- services live under `services/<service>/`
- each service is its own Go module
- cross-service imports are forbidden
- services communicate via HTTP APIs or asynchronous events
- shared technical code lives in `libs/platform`
- service internals follow `adapters → application → domain`

---

# Core Architecture Context

These documents define the fundamental architecture of the platform.

| Document | Purpose |
|--------|--------|
| SESSION_BOOTSTRAP.md | Quick orientation and core architecture rules |
| ARCHITECTURE_PRINCIPLES.md | Stable architectural guardrails |
| ENGINEERING_WORKFLOW.md | Engineering and AI-assisted workflow |
| system/PROTEON_CONTEXT.md | High-level system overview |

These files together define the **architecture intent of the platform**.

---

# System Architecture

Detailed architecture descriptions are stored under:

`system/`

| Document | Purpose |
|--------|--------|
| PROTEON_CONTEXT.md | System overview and platform structure |
| DOMAIN_BOUNDARIES.md | Service ownership and domain boundaries |
| EVENT_MODEL.md | Event-driven architecture guidelines |
| SHARED_LIBRARIES.md | Responsibilities of shared platform code |

These documents describe how the platform is structured internally.

---

# Architecture Prompts

Reusable prompts used for architecture discussions.

`prompts/`

| Document | Purpose |
|--------|--------|
| ARCHITECT_PROMPT.md | Guidance for architecture reasoning |
| CHAT_RESET_PROMPT.md | Reset prompt for starting new architecture chats |

These help ensure consistent architecture discussions.

---

# Architecture Templates

Reusable templates for documenting architecture decisions.

`templates/`

| Document | Purpose |
|--------|--------|
| ARCHITECTURE_BRIEF_TEMPLATE.md | Template for architecture design documents |
| ADR_TEMPLATE.md | Template for Architecture Decision Records |

Use these templates when documenting new architecture work.

---

# Relationship to Other Repository Documents

Architecture context works together with the following repository files:

| File | Purpose |
|----|----|
| ENGINEERING.md | Canonical engineering guide |
| DEV.md | Development workflow and commands |
| .cursorrules | Cursor agent guardrails |

Architecture documents define **intent**, while these files define
**engineering rules and execution constraints**.

---

# How to Use This Index

### For Humans

Start here when exploring the platform architecture.

Follow links to the documents relevant to your question.

### For AI Tools

Use this file as the entry point for architecture context.

Primary architecture documents are located under:

`docs/architecture/`

When discussing system design or proposing changes, consult these documents
before proposing architectural modifications.

---

# Documentation Discipline

Architecture decisions should be persisted in the repository.

Do not rely on chat history as the source of truth.

When architecture changes materially:

1. update the relevant architecture document
2. update this index if the structure changed
3. document the decision using an architecture brief or ADR

---

# Future Expansion

As the platform grows, this index may expand to include:

- service architecture summaries
- event catalogs
- deployment/runtime topology
- observability architecture
- security and identity model