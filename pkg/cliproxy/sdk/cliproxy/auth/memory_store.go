package auth

import (
	"context"
	"sync"
)

// MemoryStore implements Store interface with in-memory storage.
// This is suitable for desktop applications where persistence is not required.
type MemoryStore struct {
	mu    sync.RWMutex
	auths map[string]*Auth
}

// NewMemoryStore creates a new in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		auths: make(map[string]*Auth),
	}
}

// Save stores an auth entry in memory.
func (s *MemoryStore) Save(ctx context.Context, auth *Auth) (string, error) {
	if auth == nil || auth.ID == "" {
		return "", nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.auths[auth.ID] = auth.Clone()
	return auth.ID, nil
}

// Get retrieves an auth entry by ID.
func (s *MemoryStore) Get(ctx context.Context, id string) (*Auth, error) {
	if id == "" {
		return nil, nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if auth, ok := s.auths[id]; ok {
		return auth.Clone(), nil
	}
	return nil, nil
}

// List returns all auth entries.
func (s *MemoryStore) List(ctx context.Context) ([]*Auth, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]*Auth, 0, len(s.auths))
	for _, auth := range s.auths {
		list = append(list, auth.Clone())
	}
	return list, nil
}

// Delete removes an auth entry by ID.
func (s *MemoryStore) Delete(ctx context.Context, id string) error {
	if id == "" {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.auths, id)
	return nil
}
