package repo

import (
	"context"
	"time"

	"github.com/fin-helper/backend/internal/domain"
	"github.com/google/uuid"
)

type TagStore interface {
	ListTags(ctx context.Context, userID uuid.UUID) ([]domain.Tag, error)
	CreateTag(ctx context.Context, userID uuid.UUID, name, color string) (*domain.Tag, error)
	DeleteTag(ctx context.Context, userID, tagID uuid.UUID) error
	SetTransactionTags(ctx context.Context, userID, txID uuid.UUID, tagIDs []uuid.UUID) error
}

type BudgetStore interface {
	ListBudgets(ctx context.Context, userID uuid.UUID) ([]domain.Budget, error)
	BudgetStatusThisMonth(ctx context.Context, userID uuid.UUID) ([]domain.BudgetStatus, error)
	UpsertBudget(ctx context.Context, userID, categoryID uuid.UUID, amount float64, currency string) (*domain.Budget, error)
	DeleteBudget(ctx context.Context, userID, budgetID uuid.UUID) error
}

type RecurringStore interface {
	ListRecurring(ctx context.Context, userID uuid.UUID) ([]domain.RecurringTransaction, error)
	CreateRecurring(ctx context.Context, userID uuid.UUID, in domain.RecurringTransaction) (*domain.RecurringTransaction, error)
	ToggleRecurring(ctx context.Context, userID, id uuid.UUID, active bool) error
	DeleteRecurring(ctx context.Context, userID, id uuid.UUID) error
	RunRecurringDue(ctx context.Context, now time.Time) (int, error)
}

