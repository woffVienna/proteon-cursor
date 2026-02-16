# Identity Service – Architecture Review

Review against ENGINEERING.md and .cursorrules.

---

## Summary

| Category | Status | Notes |
|----------|--------|-------|
| Cross-service imports | ✅ OK | No violations |
| Domain purity | ✅ Fixed | Interfaces moved to application/interfaces |
| Application layer | ✅ Fixed | Uses application/interfaces |
| Adapters | ✅ OK | Implement interfaces, own mapping |
| OpenAPI paths | ✅ Fixed | Bundle path corrected to .build/generated/openapi.bundle.yml |
| Build artifacts | ⚠️ Minor | Bundle in .build/generated/ vs ENGINEERING .build/ |
| Canonical layout | ⚠️ Optional | Flat structure acceptable for small service |

---

## 1. Boundary Interfaces ✅ DONE

**Rule:** "Boundary interfaces live in internal/application/interfaces"

**Change applied:** Moved `CredentialValidator`, `RefreshTokenStore`, `TokenIssuer` to `internal/application/interfaces/auth.go`. Domain now keeps only pure types (`TokenPair`, `SessionInfo`, `UserInfo`) and errors.

---

## 2. OpenAPI Bundle Path ✅ DONE

**Change applied:** Updated `OpenAPIBundlePath` to `.build/generated/openapi.bundle.yml`. Swagger UI and OpenAPI spec endpoints now serve the correct bundle.

---

## 3. Bundle Path vs ENGINEERING.md

**ENGINEERING.md:** `Bundled spec: .build/openapi.bundle.yml (ignored)`

**Identity Makefile:** `OPENAPI_BUNDLE := $(GEN_DIR)/openapi.bundle.yml` where `GEN_DIR := $(BUILD_DIR)/generated` → `.build/generated/openapi.bundle.yml`

**Recommendation:** Align with ENGINEERING: use `.build/openapi.bundle.yml` (no `generated/` subdir). Update Makefile and server config.

---

## 4. Duplicate Helper Functions ✅ DONE

**Change applied:** Consolidated in `internal/adapters/auth/env.go` as `IssuerFromEnv()` and `AudienceFromEnv()`. Removed `adapters/http/helpers.go`. JWT issuer and HTTP server both use auth package.

---

## 5. Makefile Clean Target ✅ DONE

**Change applied:** Removed `@rm -f $(HTTP_GEN_FILE)` from clean. Clean now only removes `.build/` (build artifacts). Generated Go code is committed per ENGINEERING.

---

## 6. Optional: Canonical Layout

ENGINEERING.md canonical layout includes:
- `domain/model/` – domain entities
- `domain/rules/` – domain rules
- `application/dto/` – request/response DTOs
- `internal/platform/` – health, shutdown, buildinfo

For a small service, flat structure is acceptable. Consider migrating when the service grows.

---

## 7. Unused Adapters

`internal/adapters/db/db.go` and `internal/adapters/nats/nats.go` are stubs. They are fine as placeholders for future use.
