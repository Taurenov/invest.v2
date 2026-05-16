package repo

import (
	"context"
	"sync"

	"github.com/fin-helper/backend/internal/domain"
	"github.com/google/uuid"
)

type MemoryGoalStore struct {
	mu    sync.RWMutex
	goals []domain.Goal
}

func NewMemoryGoalStore() *MemoryGoalStore {
	uid := uuid.MustParse("00000000-0000-4000-8000-000000000001")
	return &MemoryGoalStore{goals: []domain.Goal{
		{
			ID: uuid.New(), UserID: uid, Title: "Подушка безопасности",
			GoalType: "savings", TargetAmount: 300000, CurrentAmount: 200000, Currency: "RUB",
		},
	}}
}

func (m *MemoryGoalStore) ListGoalsByUser(_ context.Context, userID uuid.UUID) ([]domain.Goal, error) {
	var out []domain.Goal
	for _, g := range m.goals {
		if g.UserID == userID {
			out = append(out, g)
		}
	}
	return out, nil
}
