# Architecture Brief: Backoffice Auth Baseline

## Status

Accepted. Baseline only; implementation deferred.

## Goal

Establish the baseline authentication model for the backoffice: separate
edge (backoffice-gateway), auth methods and credentials in the auth
service, Identity for user records and token issuance. Two types of
backoffice users: operators (Proteon staff) and tenant users (the
tenant's staff). See `GLOSSARY.md` for persona terms.

## Context

Proteon has a backoffice application used by operators (control plane)
and by tenant users (tenant self-service: settings, analytics). This
auth model is separate from player-facing auth (see
`briefs/player-auth-baseline.md`).

Reference documents:

- `product/PRODUCT_CONTEXT.md`
- `system/API_GATEWAY.md`
- `system/SERVICE_TYPES.md`
- `GLOSSARY.md`

## Constraints

- Services are independent; integration via HTTP only
- Backoffice-gateway is the single entry point for backoffice traffic
- Credentials for backoffice users live in the auth service; Identity
  holds user records only (link: `user_id`)
- Login and registration are not JWT-secured; they are protected by
  app-key from the backoffice app

Reference documents:

- `01_PRINCIPLES.md`
- `00_CONTEXT.md`
- `system/DOMAIN_BOUNDARIES.md`

------------------------------------------------------------------------

## Architecture Overview

- **backoffice-gateway:** Edge service. Exposes login and register
  routes protected by app-key. All other routes require JWT. Validates
  JWT using Identity's JWKS; forwards claims to downstream services.
- **auth:** Owns authentication methods and credential storage for
  backoffice users. Validates credentials (or runs OAuth flow), then
  calls Identity to issue a JWT for the authenticated `user_id`. Does
  not store user records; does not issue tokens.
- **Identity:** Holds user records for all platform users (players and
  backoffice users: operators, tenant users). Issues JWTs on request
  from auth (internal API). No credentials stored in Identity for
  backoffice users.

Backoffice users are either **operators** (Proteon staff) or **tenant
users** (the tenant's staff). Token claims carry subject type and tenant
scope so downstream services can enforce access.

------------------------------------------------------------------------

## Key Components

- `services/backoffice-gateway/`
- `services/auth/`
- `services/identity/`
- `briefs/player-auth-baseline.md` (player auth; separate flow)
- `GLOSSARY.md` (personas)

------------------------------------------------------------------------

## Data Flow

1. Backoffice app calls backoffice-gateway with app-key for login/register
2. Gateway forwards to auth service
3. Auth validates credentials (or completes OAuth); looks up or creates
   credential binding keyed by `user_id`
4. Auth calls Identity (internal) to issue JWT for that `user_id`
5. Identity returns JWT; auth returns it to gateway; gateway to app
6. App uses JWT for all other requests to backoffice-gateway
7. Gateway validates JWT, extracts claims, forwards to downstream services

------------------------------------------------------------------------

## Deferred

- Detailed implementation plan
- App-key lifecycle and storage
- OAuth, SSO, MFA for backoffice
- Token refresh and revocation for backoffice
