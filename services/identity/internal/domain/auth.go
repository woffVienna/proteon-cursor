package domain

import (
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

// UserInfo identifies an authenticated user.
type UserInfo struct {
	ID string
}
