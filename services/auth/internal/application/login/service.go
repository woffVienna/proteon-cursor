package login

import (
	"context"

	"github.com/woffVienna/proteon-cursor/services/auth/internal/application/interfaces"
	"github.com/woffVienna/proteon-cursor/services/auth/internal/domain"
)

// Service implements the backoffice login use case.
type Service struct {
	creds  interfaces.CredentialStore
	idAuth interfaces.IdentityTokenClient
}

// NewService creates a login service with the given dependencies.
func NewService(creds interfaces.CredentialStore, idAuth interfaces.IdentityTokenClient) *Service {
	return &Service{
		creds:  creds,
		idAuth: idAuth,
	}
}

// Login authenticates the user and returns a backoffice access token.
func (s *Service) Login(ctx context.Context, username, password string) (domain.LoginResult, error) {
	userID, subjectType, tenant, err := s.creds.Validate(ctx, username, password)
	if err != nil {
		return domain.LoginResult{}, err
	}

	return s.idAuth.IssueBackofficeToken(ctx, userID, subjectType, tenant)
}
