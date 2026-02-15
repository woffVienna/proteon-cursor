package auth

import (
	"context"
	"sync"

	"github.com/woffVienna/proteon-cursor/services/identity/internal/domain"
)

// MemoryStore is an in-memory implementation of RefreshTokenStore.
// For production, replace with a DB-backed store (e.g. Postgres).
type MemoryStore struct {
	mu   sync.Mutex
	data map[string]domain.SessionInfo
}

// NewMemoryStore creates an in-memory refresh token store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]domain.SessionInfo),
	}
}

// Store implements domain.RefreshTokenStore.
func (s *MemoryStore) Store(_ context.Context, token string, info domain.SessionInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[token] = info
	return nil
}

// Get implements domain.RefreshTokenStore.
// Returns (info, true, nil) if found, (zero, false, nil) if not found.
func (s *MemoryStore) Get(_ context.Context, token string) (domain.SessionInfo, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	info, ok := s.data[token]
	return info, ok, nil
}

// Delete implements domain.RefreshTokenStore.
func (s *MemoryStore) Delete(_ context.Context, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, token)
	return nil
}
