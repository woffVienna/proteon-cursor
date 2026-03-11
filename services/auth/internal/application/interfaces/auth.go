package interfaces

import (
	"context"

	"github.com/woffVienna/proteon-cursor/services/auth/internal/domain"
)

// CredentialStore validates credentials for backoffice users.
type CredentialStore interface {
	Validate(ctx context.Context, username, password string) (userID string, subjectType string, tenant string, err error)
}

// IdentityTokenClient issues backoffice tokens via the identity service.
type IdentityTokenClient interface {
	IssueBackofficeToken(ctx context.Context, userID, subjectType, tenant string) (domain.LoginResult, error)
}
