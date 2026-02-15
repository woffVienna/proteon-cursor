package auth

import (
	"context"

	"github.com/woffVienna/proteon-cursor/services/identity/internal/domain"
)

const (
	demoLogin    = "demo@proteon.dev"
	demoPassword = "demo"
	demoUserID   = "00000000-0000-0000-0000-000000000001"
)

// DemoValidator validates credentials against a hardcoded demo user.
// For production, replace with a DB-backed validator.
type DemoValidator struct{}

// NewDemoValidator creates a demo credential validator.
func NewDemoValidator() *DemoValidator {
	return &DemoValidator{}
}

// Validate implements domain.CredentialValidator.
func (v *DemoValidator) Validate(ctx context.Context, login, password string) (domain.UserInfo, error) {
	if login != demoLogin || password != demoPassword {
		return domain.UserInfo{}, domain.ErrInvalidCredentials
	}
	return domain.UserInfo{ID: demoUserID}, nil
}
