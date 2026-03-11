# Architecture Glossary

Canonical terminology for personas and key concepts. Use these terms
consistently in architecture and service documents.

------------------------------------------------------------------------

# Personas

| Term | Definition |
| --- | --- |
| **tenant** | The organization (company or platform) that has a tenancy on Proteon — e.g. a game studio or white-label operator integrating with the platform. Identified by `tenant_id`. Use "tenant" in technical and architecture docs. |
| **players** | The end-users of a tenant's product: people who use the tenant's app (e.g. the customer platform). They authenticate as players; Identity holds player identities. |
| **backoffice user** | Anyone using the backoffice application. Umbrella term for operator and tenant user. |
| **operator** | A backoffice user who is Proteon staff. Uses the backoffice as a control plane (e.g. onboarding tenants, platform config). |
| **tenant user** | A backoffice user who is the tenant's employee. Uses the backoffice to manage their tenant (e.g. settings, analytics). |

------------------------------------------------------------------------

# Documentation Rule

When writing or updating architecture docs (service documents, briefs,
system docs), use the persona terms above. Point to this glossary as the
source of truth.
