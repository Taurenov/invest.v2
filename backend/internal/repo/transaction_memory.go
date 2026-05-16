package repo

import (
	"context"
	"sync"

	"github.com/fin-helper/backend/internal/domain"
	"github.com/google/uuid"
)

// MemoryStore — заглушка для примера; в prod — PostgreSQL через pgx/sqlc.
type MemoryStore struct {
	mu   sync.RWMutex
	data map[uuid.UUID][]domain.Transaction
}

func NewMemoryStore(seed []domain.Transaction) *MemoryStore {
	s := &MemoryStore{data: make(map[uuid.UUID][]domain.Transaction)}
	for _, t := range seed {
		s.data[t.UserID] = append(s.data[t.UserID], t)
	}
	return s
}

func (s *MemoryStore) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Transaction, error) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Transaction, len(s.data[userID]))
	copy(out, s.data[userID])
	return out, nil
}
