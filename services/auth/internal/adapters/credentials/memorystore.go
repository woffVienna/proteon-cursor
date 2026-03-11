package credentials

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/woffVienna/proteon-cursor/services/auth/internal/domain"
)

// MemoryStore is an in-memory credential store for backoffice users.
type MemoryStore struct {
	// For now we keep a single hard-coded user.
	username string
	hash     []byte
	userID   string
	tenant   string
}

// NewMemoryStore creates a store with a single user.
func NewMemoryStore() (*MemoryStore, error) {
	// username: robert, password: proteon
	hash, err := bcrypt.GenerateFromPassword([]byte("proteon"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &MemoryStore{
		username: "robert",
		hash:     hash,
		// Stable dummy platform user id (UUID-format string).
		userID: "00000000-0000-0000-0000-000000000001",
		tenant: "proteon",
	}, nil
}

// Validate checks credentials and returns user identity information.
func (s *MemoryStore) Validate(_ context.Context, username, password string) (string, string, string, error) {
	if username != s.username {
		return "", "", "", domain.ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword(s.hash, []byte(password)); err != nil {
		return "", "", "", domain.ErrInvalidCredentials
	}
	// subjectType: operator (backoffice user from Proteon)
	return s.userID, "operator", s.tenant, nil
}
