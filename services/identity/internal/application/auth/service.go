package auth

import (
	"context"
	"time"

	"github.com/woffVienna/proteon-cursor/services/identity/internal/application/interfaces"
	"github.com/woffVienna/proteon-cursor/services/identity/internal/domain"
)

const accessTokenTTL = 10 * time.Minute
const backofficeTokenTTL = 10 * time.Minute

// Service implements the auth exchange use case.
type Service struct {
	resolver interfaces.IdentityResolver
	lookup   interfaces.IdentityLookup
	issuer   interfaces.TokenIssuer
}

// NewService creates an auth service with the given dependencies.
func NewService(
	resolver interfaces.IdentityResolver,
	lookup interfaces.IdentityLookup,
	issuer interfaces.TokenIssuer,
) *Service {
	return &Service{
		resolver: resolver,
		lookup:   lookup,
		issuer:   issuer,
	}
}

// Exchange processes an external identity assertion from a customer backend.
// It resolves or creates the platform identity and issues a short-lived access JWT.
func (s *Service) Exchange(ctx context.Context, provider, externalUserID, tenant string) (*domain.TokenResult, error) {
	if provider == "" || externalUserID == "" {
		return nil, domain.ErrInvalidAssertion
	}

	identity, err := s.resolver.Resolve(ctx, provider, externalUserID, tenant)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.issuer.Issue(ctx, identity.PlatformUserID, identity.Tenant, accessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &domain.TokenResult{
		AccessToken:    accessToken,
		PlatformUserID: identity.PlatformUserID,
		ExpiresIn:      int32(accessTokenTTL.Seconds()),
	}, nil
}

// IssueBackofficeToken issues a backoffice access token for a known platform user.
func (s *Service) IssueBackofficeToken(ctx context.Context, userID, subjectType, tenant, audience string) (*domain.TokenResult, error) {
	accessToken, err := s.issuer.IssueBackoffice(ctx, userID, subjectType, tenant, audience, backofficeTokenTTL)
	if err != nil {
		return nil, err
	}
	return &domain.TokenResult{
		AccessToken:    accessToken,
		PlatformUserID: userID,
		ExpiresIn:      int32(backofficeTokenTTL.Seconds()),
	}, nil
}

// GetIdentity retrieves a platform identity by platform user ID.
func (s *Service) GetIdentity(ctx context.Context, platformUserID string) (*domain.PlatformIdentity, error) {
	identity, err := s.lookup.GetByPlatformUserID(ctx, platformUserID)
	if err != nil {
		return nil, err
	}
	return &identity, nil
}
