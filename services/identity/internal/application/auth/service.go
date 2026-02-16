package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/woffVienna/proteon-cursor/services/identity/internal/application/interfaces"
	"github.com/woffVienna/proteon-cursor/services/identity/internal/domain"
)

const (
	accessTokenTTL   = 10 * time.Minute
	refreshTokenTTL  = 7 * 24 * time.Hour
	refreshTokenSize = 32
)

// Service implements the auth use cases.
type Service struct {
	validator interfaces.CredentialValidator
	store     interfaces.RefreshTokenStore
	issuer    interfaces.TokenIssuer
}

// NewService creates an auth service with the given dependencies.
func NewService(
	validator interfaces.CredentialValidator,
	store interfaces.RefreshTokenStore,
	issuer interfaces.TokenIssuer,
) *Service {
	return &Service{
		validator: validator,
		store:     store,
		issuer:    issuer,
	}
}

// Login validates credentials, creates a session, and returns a token pair.
func (s *Service) Login(ctx context.Context, login, password, tenant string) (*domain.TokenPair, error) {
	if tenant == "" {
		tenant = "proteon.dev"
	}

	user, err := s.validator.Validate(ctx, login, password)
	if err != nil {
		return nil, err
	}

	refreshToken, err := newOpaqueToken(refreshTokenSize)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(refreshTokenTTL)
	if err := s.store.Store(ctx, refreshToken, domain.SessionInfo{
		UserID:    user.ID,
		Tenant:    tenant,
		ExpiresAt: expiresAt,
	}); err != nil {
		return nil, err
	}

	accessToken, err := s.issuer.Issue(ctx, user.ID, tenant, accessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int32(accessTokenTTL.Seconds()),
	}, nil
}

// Refresh validates the refresh token, rotates it, and returns a new token pair.
func (s *Service) Refresh(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	info, ok, err := s.store.Get(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, domain.ErrInvalidRefreshToken
	}
	if time.Now().After(info.ExpiresAt) {
		return nil, domain.ErrRefreshTokenExpired
	}

	// Rotate: delete old token
	_ = s.store.Delete(ctx, refreshToken)

	newRefreshToken, err := newOpaqueToken(refreshTokenSize)
	if err != nil {
		return nil, err
	}

	if err := s.store.Store(ctx, newRefreshToken, info); err != nil {
		return nil, err
	}

	accessToken, err := s.issuer.Issue(ctx, info.UserID, info.Tenant, accessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int32(accessTokenTTL.Seconds()),
	}, nil
}

// Logout revokes the refresh token. Idempotent.
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	return s.store.Delete(ctx, refreshToken)
}

func newOpaqueToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
