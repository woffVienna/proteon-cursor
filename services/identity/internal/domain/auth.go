package domain

import (
	"errors"
	"time"
)

var (
	ErrIdentityNotFound = errors.New("platform identity not found")
	ErrInvalidAssertion = errors.New("invalid external identity assertion")
)

// PlatformIdentity represents a reduced Proteon platform identity.
// It maps an external provider + external user ID to a stable platform user ID.
type PlatformIdentity struct {
	PlatformUserID string
	Provider       string
	ExternalUserID string
	Tenant         string
	CreatedAt      time.Time
}

// TokenResult is the result of a successful auth exchange.
type TokenResult struct {
	AccessToken    string
	PlatformUserID string
	ExpiresIn      int32
}
