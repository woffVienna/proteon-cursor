package interfaces

import (
	"context"
	"crypto/ed25519"
	"time"

	"github.com/woffVienna/proteon-cursor/services/identity/internal/domain"
)

// CredentialValidator validates login credentials and returns user info.
// Implemented by adapters (e.g. DB, demo).
type CredentialValidator interface {
	Validate(ctx context.Context, login, password string) (domain.UserInfo, error)
}

// RefreshTokenStore stores and retrieves refresh tokens.
// Implemented by adapters (e.g. in-memory, Postgres).
type RefreshTokenStore interface {
	Store(ctx context.Context, token string, info domain.SessionInfo) error
	Get(ctx context.Context, token string) (domain.SessionInfo, bool, error)
	Delete(ctx context.Context, token string) error
}

// TokenIssuer issues signed access tokens (JWTs).
// Implemented by adapters (e.g. Ed25519 JWT issuer).
type TokenIssuer interface {
	Issue(ctx context.Context, userID, tenant string, ttl time.Duration) (string, error)
	PublicKey() ed25519.PublicKey
	Kid() string
}
