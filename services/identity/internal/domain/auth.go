package domain

import (
	"context"
	"crypto/ed25519"
	"errors"
	"time"
)

// Domain errors for auth. Adapters map these to HTTP status codes.
var (
	ErrInvalidCredentials    = errors.New("invalid login or password")
	ErrInvalidRefreshToken   = errors.New("invalid refresh token")
	ErrRefreshTokenExpired   = errors.New("refresh token expired")
	ErrMissingRequestPayload = errors.New("missing request body")
)

// TokenPair is the result of a successful login or refresh.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int32
}

// SessionInfo holds data stored with a refresh token.
type SessionInfo struct {
	UserID    string
	Tenant    string
	ExpiresAt time.Time
}

// CredentialValidator validates login credentials and returns user info.
// Implemented by adapters (e.g. DB, demo).
type CredentialValidator interface {
	Validate(ctx context.Context, login, password string) (UserInfo, error)
}

// UserInfo identifies an authenticated user.
type UserInfo struct {
	ID string
}

// RefreshTokenStore stores and retrieves refresh tokens.
// Implemented by adapters (e.g. in-memory, Postgres).
type RefreshTokenStore interface {
	Store(ctx context.Context, token string, info SessionInfo) error
	Get(ctx context.Context, token string) (SessionInfo, bool, error)
	Delete(ctx context.Context, token string) error
}

// TokenIssuer issues signed access tokens (JWTs).
// Implemented by adapters (e.g. Ed25519 JWT issuer).
type TokenIssuer interface {
	Issue(ctx context.Context, userID, tenant string, ttl time.Duration) (string, error)
	PublicKey() ed25519.PublicKey
	Kid() string
}
