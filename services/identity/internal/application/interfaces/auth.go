package interfaces

import (
	"context"
	"crypto/ed25519"
	"time"

	"github.com/woffVienna/proteon-cursor/services/identity/internal/domain"
)

// IdentityResolver resolves or creates a platform identity from an external
// identity assertion. Implemented by adapters (e.g. in-memory, Postgres).
type IdentityResolver interface {
	Resolve(ctx context.Context, provider, externalUserID, tenant string) (domain.PlatformIdentity, error)
}

// IdentityLookup retrieves an existing platform identity by platform user ID.
// Implemented by adapters (e.g. in-memory, Postgres).
type IdentityLookup interface {
	GetByPlatformUserID(ctx context.Context, platformUserID string) (domain.PlatformIdentity, error)
}

// TokenIssuer issues signed access tokens (JWTs).
// Implemented by adapters (e.g. Ed25519 JWT issuer).
type TokenIssuer interface {
	Issue(ctx context.Context, platformUserID, tenant string, ttl time.Duration) (string, error)
	IssueBackoffice(ctx context.Context, userID, subjectType, tenant, audience string, ttl time.Duration) (string, error)
	PublicKey() ed25519.PublicKey
	Kid() string
}
