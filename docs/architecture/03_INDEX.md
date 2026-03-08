# Proteon Architecture Index

This document is the navigation map for the architecture documents.

If you are new to the repository, start with:

1. `00_CONTEXT.md`
2. `01_PRINCIPLES.md`
3. `02_WORKFLOW.md`

Then use the topic-specific documents below.

------------------------------------------------------------------------

# Core Documents

| Document | Purpose |
| --- | --- |
| `00_CONTEXT.md` | single entry point for system context |
| `01_PRINCIPLES.md` | stable architecture guardrails |
| `02_WORKFLOW.md` | engineering and AI-assisted workflow |
| `03_INDEX.md` | navigation map for architecture documents |

------------------------------------------------------------------------

# System Architecture Documents

| Document | Purpose |
| --- | --- |
| `system/DOMAIN_BOUNDARIES.md` | service ownership and boundary rules |
| `system/SHARED_LIBRARIES.md` | role and limits of shared platform code |
| `system/INTEGRATION_CONTRACTS.md` | HTTP and event contract placement and usage |
| `system/CONTRACT_GOVERNANCE.md` | contract evolution and versioning rules |
| `system/EVENT_MODEL.md` | event-driven architecture semantics |
| `system/EVENT_OPERATIONS.md` | runtime behaviour for event processing |
| `system/SERVICE_TYPES.md` | allowed service roles and expectations |
| `system/API_GATEWAY.md` | role and limits of the future API gateway |

------------------------------------------------------------------------

# Prompt Assets

Prompt assets live under `prompts/`.

They are not the source of truth. They are convenience documents that point
back to the canonical architecture docs.

| Document | Purpose |
| --- | --- |
| `prompts/ARCHITECT_PROMPT.md` | reusable prompt for architecture reasoning |
| `prompts/CHAT_RESET_PROMPT.md` | reset prompt for new chats |
| `prompts/00_ARCHITECTURE_CONTEXT.md` | compact bootstrap prompt |

------------------------------------------------------------------------

# Templates

| Document | Purpose |
| --- | --- |
| `templates/ARCHITECTURE_BRIEF_TEMPLATE.md` | template for solution briefs |
| `templates/ADR_TEMPLATE.md` | template for architecture decisions |

------------------------------------------------------------------------

# Relationship to Other Repository Files

| File | Purpose |
| --- | --- |
| `ENGINEERING.md` | engineering guide for repository work |
| `DEV.md` | development workflow and commands |
| `.cursorrules` | Cursor execution guardrails |

Architecture documents define durable design intent. The root documents define
engineering conventions and execution constraints.

------------------------------------------------------------------------

# Documentation Rule

When architecture changes materially:

- update the relevant topic document
- update this index if the structure changed
- add or update an architecture brief or ADR where appropriate
