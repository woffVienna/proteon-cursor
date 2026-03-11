package auth

import (
	"context"
	"sync"
	"time"

	"github.com/woffVienna/proteon-cursor/services/identity/internal/domain"
)

// linkageKey uniquely identifies an external identity.
type linkageKey struct {
	Provider       string
	ExternalUserID string
}

// MemoryIdentityStore is an in-memory implementation of both IdentityResolver
// and IdentityLookup. For production, replace with a Postgres-backed store.
type MemoryIdentityStore struct {
	mu       sync.Mutex
	linkages map[linkageKey]domain.PlatformIdentity
	byID     map[string]domain.PlatformIdentity
	idGen    func() string
}

// NewMemoryIdentityStore creates an in-memory identity store.
// idGen provides platform user IDs (e.g. UUID generator).
func NewMemoryIdentityStore(idGen func() string) *MemoryIdentityStore {
	return &MemoryIdentityStore{
		linkages: make(map[linkageKey]domain.PlatformIdentity),
		byID:     make(map[string]domain.PlatformIdentity),
		idGen:    idGen,
	}
}

// Resolve implements interfaces.IdentityResolver.
// Returns an existing platform identity or creates a new one.
func (s *MemoryIdentityStore) Resolve(_ context.Context, provider, externalUserID, tenant string) (domain.PlatformIdentity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := linkageKey{Provider: provider, ExternalUserID: externalUserID}

	if identity, ok := s.linkages[key]; ok {
		return identity, nil
	}

	identity := domain.PlatformIdentity{
		PlatformUserID: s.idGen(),
		Provider:       provider,
		ExternalUserID: externalUserID,
		Tenant:         tenant,
		CreatedAt:      time.Now(),
	}
	s.linkages[key] = identity
	s.byID[identity.PlatformUserID] = identity

	return identity, nil
}

// GetByPlatformUserID implements interfaces.IdentityLookup.
func (s *MemoryIdentityStore) GetByPlatformUserID(_ context.Context, platformUserID string) (domain.PlatformIdentity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	identity, ok := s.byID[platformUserID]
	if !ok {
		return domain.PlatformIdentity{}, domain.ErrIdentityNotFound
	}
	return identity, nil
}
