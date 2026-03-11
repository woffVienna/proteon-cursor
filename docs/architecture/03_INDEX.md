# Proteon Architecture Index

This document is the navigation map for the architecture documents.

If you are new to the repository, start with:

1. `00_CONTEXT.md`
2. `product/PRODUCT_CONTEXT.md`
3. `01_PRINCIPLES.md`
4. `02_WORKFLOW.md`

Then use the topic-specific documents below.

------------------------------------------------------------------------

# Core Documents

| Document | Purpose |
| --- | --- |
| `00_CONTEXT.md` | single entry point for system context |
| `product/PRODUCT_CONTEXT.md` | platform definition, boundaries, domains, non-goals |
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
| `system/API_GATEWAY.md` | role and limits of the API gateway |

------------------------------------------------------------------------

# Service Documents

| Document | Purpose |
| --- | --- |
| `services/api-gateway.md` | ownership, responsibilities, and auth behaviour for the API gateway |
| `services/identity.md` | ownership, responsibilities, and auth exchange for the identity service |

------------------------------------------------------------------------

# Architecture Briefs

| Document | Purpose |
| --- | --- |
| `briefs/auth-baseline.md` | public auth baseline decision: JWT model, gateway/identity ownership split |

------------------------------------------------------------------------

# GPT Assets

GPT assets live under `gpt/`.

They are not the source of truth. They are operational guidance for how
ChatGPT should use the canonical architecture docs.

| Document | Purpose |
| --- | --- |
| `gpt/ARCHITECT_PROJECT_INSTRUCTIONS.md` | ChatGPT project instructions for architecture reasoning |

------------------------------------------------------------------------

# Templates

| Document | Purpose |
| --- | --- |
| `templates/ARCHITECTURE_BRIEF_TEMPLATE.md` | template for solution briefs |
| `templates/ADR_TEMPLATE.md` | template for architecture decisions |
| `templates/SERVICE_TEMPLATE.md` | template for per-service architecture/ownership docs |

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
