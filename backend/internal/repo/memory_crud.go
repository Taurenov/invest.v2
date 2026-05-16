package repo

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/fin-helper/backend/internal/domain"
	"github.com/google/uuid"
)

type memoryUserStore struct {
	mu    sync.RWMutex
	users map[string]struct {
		id           uuid.UUID
		passwordHash string
		displayName  string
	}
	byID map[uuid.UUID]string
}

func NewMemoryUserStore() *memoryUserStore {
	return &memoryUserStore{
		users: make(map[string]struct {
			id           uuid.UUID
			passwordHash string
			displayName  string
		}),
		byID: make(map[uuid.UUID]string),
	}
}

func (m *memoryUserStore) CreateUser(_ context.Context, email, passwordHash, displayName string) (*domain.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.users[email]; ok {
		return nil, errors.New("email exists")
	}
	u := domain.User{ID: uuid.New(), Email: email, DisplayName: displayName, CreatedAt: time.Now()}
	m.users[email] = struct {
		id           uuid.UUID
		passwordHash string
		displayName  string
	}{u.ID, passwordHash, displayName}
	m.byID[u.ID] = email
	return &u, nil
}

func (m *memoryUserStore) GetByEmail(_ context.Context, email string) (uuid.UUID, string, string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.users[email]
	if !ok {
		return uuid.Nil, "", "", errors.New("not found")
	}
	return u.id, u.passwordHash, u.displayName, nil
}

func (m *memoryUserStore) GetByID(_ context.Context, id uuid.UUID) (*domain.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	email, ok := m.byID[id]
	if !ok {
		return nil, errors.New("not found")
	}
	u := m.users[email]
	return &domain.User{ID: id, Email: email, DisplayName: u.displayName, CreatedAt: time.Now()}, nil
}

func (m *memoryUserStore) EnsureDefaultCategories(context.Context, uuid.UUID) error { return nil }

func (s *MemoryStore) CreateTransaction(_ context.Context, userID uuid.UUID, in CreateTransactionInput) (*domain.Transaction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := domain.Transaction{
		ID: uuid.New(), UserID: userID, CategoryID: in.CategoryID, Kind: in.Kind,
		Amount: in.Amount, Currency: in.Currency, Description: in.Description, OccurredAt: in.OccurredAt,
	}
	s.data[userID] = append(s.data[userID], t)
	return &t, nil
}

func (s *MemoryStore) UpdateTransaction(_ context.Context, userID, txID uuid.UUID, in UpdateTransactionInput) (*domain.Transaction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	list := s.data[userID]
	for i, t := range list {
		if t.ID != txID {
			continue
		}
		if in.Kind != nil {
			list[i].Kind = *in.Kind
		}
		if in.Amount != nil {
			list[i].Amount = *in.Amount
		}
		if in.Currency != nil {
			list[i].Currency = *in.Currency
		}
		if in.Description != nil {
			list[i].Description = *in.Description
		}
		if in.OccurredAt != nil {
			list[i].OccurredAt = *in.OccurredAt
		}
		if in.CategoryID != nil {
			list[i].CategoryID = in.CategoryID
		}
		out := list[i]
		return &out, nil
	}
	return nil, errors.New("not found")
}

func (s *MemoryStore) DeleteTransaction(_ context.Context, userID, txID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	list := s.data[userID]
	for i, t := range list {
		if t.ID == txID {
			s.data[userID] = append(list[:i], list[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

type MemoryCategoryStore struct {
	mu   sync.RWMutex
	data map[uuid.UUID][]domain.Category
}

func NewMemoryCategoryStore() *MemoryCategoryStore {
	return &MemoryCategoryStore{data: make(map[uuid.UUID][]domain.Category)}
}

func (m *MemoryCategoryStore) ListCategoriesByUser(_ context.Context, userID uuid.UUID) ([]domain.Category, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]domain.Category{}, m.data[userID]...), nil
}

func (m *MemoryCategoryStore) CreateCategory(_ context.Context, userID uuid.UUID, in CreateCategoryInput) (*domain.Category, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	uid := userID
	c := domain.Category{ID: uuid.New(), UserID: &uid, Name: in.Name, Kind: in.Kind, Icon: in.Icon, Color: in.Color}
	m.data[userID] = append(m.data[userID], c)
	return &c, nil
}

func (m *MemoryGoalStore) CreateGoal(_ context.Context, userID uuid.UUID, in CreateGoalInput) (*domain.Goal, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	g := domain.Goal{
		ID: uuid.New(), UserID: userID, Title: in.Title, GoalType: in.GoalType,
		TargetAmount: in.TargetAmount, CurrentAmount: 0, Currency: in.Currency, Deadline: in.Deadline,
	}
	m.goals = append(m.goals, g)
	return &g, nil
}

func (m *MemoryGoalStore) Contribute(_ context.Context, userID, goalID uuid.UUID, amount float64, _ string) (*domain.Goal, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, g := range m.goals {
		if g.ID == goalID && g.UserID == userID {
			m.goals[i].CurrentAmount += amount
			out := m.goals[i]
			return &out, nil
		}
	}
	return nil, errors.New("not found")
}
