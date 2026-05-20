package repo

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/fin-helper/backend/internal/domain"
	"github.com/google/uuid"
)

type MemoryTagStore struct {
	mu   sync.RWMutex
	data map[uuid.UUID][]domain.Tag
}

func NewMemoryTagStore() *MemoryTagStore {
	return &MemoryTagStore{data: make(map[uuid.UUID][]domain.Tag)}
}

func (m *MemoryTagStore) ListTags(_ context.Context, userID uuid.UUID) ([]domain.Tag, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]domain.Tag{}, m.data[userID]...), nil
}

func (m *MemoryTagStore) CreateTag(_ context.Context, userID uuid.UUID, name, color string) (*domain.Tag, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t := domain.Tag{ID: uuid.New(), UserID: userID, Name: name, Color: color, CreatedAt: time.Now()}
	m.data[userID] = append(m.data[userID], t)
	return &t, nil
}

func (m *MemoryTagStore) DeleteTag(_ context.Context, userID, tagID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	list := m.data[userID]
	for i, t := range list {
		if t.ID == tagID {
			m.data[userID] = append(list[:i], list[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *MemoryTagStore) SetTransactionTags(context.Context, uuid.UUID, uuid.UUID, []uuid.UUID) error {
	return nil
}

type MemoryBudgetStore struct {
	mu   sync.RWMutex
	data map[uuid.UUID][]domain.Budget
}

func NewMemoryBudgetStore() *MemoryBudgetStore {
	return &MemoryBudgetStore{data: make(map[uuid.UUID][]domain.Budget)}
}

func (m *MemoryBudgetStore) ListBudgets(_ context.Context, userID uuid.UUID) ([]domain.Budget, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]domain.Budget{}, m.data[userID]...), nil
}

func (m *MemoryBudgetStore) BudgetStatusThisMonth(ctx context.Context, userID uuid.UUID) ([]domain.BudgetStatus, error) {
	items, err := m.ListBudgets(ctx, userID)
	if err != nil {
		return nil, err
	}
	var out []domain.BudgetStatus
	for _, b := range items {
		out = append(out, domain.BudgetStatus{Budget: b, Spent: 0, Remaining: b.Amount, Percent: 0})
	}
	return out, nil
}

func (m *MemoryBudgetStore) UpsertBudget(_ context.Context, userID, categoryID uuid.UUID, amount float64, currency string) (*domain.Budget, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.data[userID] {
		if m.data[userID][i].CategoryID == categoryID {
			m.data[userID][i].Amount = amount
			m.data[userID][i].Currency = currency
			out := m.data[userID][i]
			return &out, nil
		}
	}
	b := domain.Budget{ID: uuid.New(), UserID: userID, CategoryID: categoryID, Period: "monthly", Amount: amount, Currency: currency, CreatedAt: time.Now()}
	m.data[userID] = append(m.data[userID], b)
	return &b, nil
}

func (m *MemoryBudgetStore) DeleteBudget(_ context.Context, userID, budgetID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	list := m.data[userID]
	for i, b := range list {
		if b.ID == budgetID {
			m.data[userID] = append(list[:i], list[i+1:]...)
			return nil
		}
	}
	return nil
}

type MemoryRecurringStore struct {
	mu   sync.RWMutex
	data map[uuid.UUID][]domain.RecurringTransaction
}

func NewMemoryRecurringStore() *MemoryRecurringStore {
	return &MemoryRecurringStore{data: make(map[uuid.UUID][]domain.RecurringTransaction)}
}

func (m *MemoryRecurringStore) ListRecurring(_ context.Context, userID uuid.UUID) ([]domain.RecurringTransaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]domain.RecurringTransaction{}, m.data[userID]...), nil
}

func (m *MemoryRecurringStore) CreateRecurring(_ context.Context, userID uuid.UUID, in domain.RecurringTransaction) (*domain.RecurringTransaction, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	in.ID = uuid.New()
	in.UserID = userID
	if in.NextRunAt == "" {
		in.NextRunAt = time.Now().Format("2006-01-02")
	}
	m.data[userID] = append(m.data[userID], in)
	return &in, nil
}

func (m *MemoryRecurringStore) ToggleRecurring(_ context.Context, userID, id uuid.UUID, active bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.data[userID] {
		if m.data[userID][i].ID == id {
			m.data[userID][i].IsActive = active
			return nil
		}
	}
	return errors.New("not found")
}

func (m *MemoryRecurringStore) DeleteRecurring(_ context.Context, userID, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	list := m.data[userID]
	for i, r := range list {
		if r.ID == id {
			m.data[userID] = append(list[:i], list[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *MemoryRecurringStore) RunRecurringDue(context.Context, time.Time) (int, error) {
	return 0, nil
}

func (s *MemoryStore) ListByUserFiltered(ctx context.Context, userID uuid.UUID, q TransactionQuery) ([]domain.Transaction, error) {
	// minimal in-memory filter (text + kind)
	items, _ := s.ListByUser(ctx, userID)
	var out []domain.Transaction
	for _, t := range items {
		if q.Kind != "" && string(t.Kind) != q.Kind {
			continue
		}
		if q.Text != "" && t.Description != q.Text {
			// simple contains fallback without strings import
			continue
		}
		out = append(out, t)
	}
	return out, nil
}

