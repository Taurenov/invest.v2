package repo

import (
	"context"
	"time"

	"github.com/fin-helper/backend/internal/domain"
	"github.com/google/uuid"
)

type TransactionStore interface {
	ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Transaction, error)
	ListByUserFiltered(ctx context.Context, userID uuid.UUID, q TransactionQuery) ([]domain.Transaction, error)
	CreateTransaction(ctx context.Context, userID uuid.UUID, t CreateTransactionInput) (*domain.Transaction, error)
	UpdateTransaction(ctx context.Context, userID, txID uuid.UUID, t UpdateTransactionInput) (*domain.Transaction, error)
	DeleteTransaction(ctx context.Context, userID, txID uuid.UUID) error
}

type TransactionQuery struct {
	Text      string
	Kind      string
	From      *time.Time
	To        *time.Time
	TagIDs    []uuid.UUID
	CategoryID *uuid.UUID
	Limit     int
}

type CreateTransactionInput struct {
	CategoryID  *uuid.UUID
	Kind        domain.TransactionKind
	Amount      float64
	Currency    string
	Description string
	OccurredAt  time.Time
}

type UpdateTransactionInput struct {
	CategoryID  *uuid.UUID
	Kind        *domain.TransactionKind
	Amount      *float64
	Currency    *string
	Description *string
	OccurredAt  *time.Time
}

type CategoryStore interface {
	ListCategoriesByUser(ctx context.Context, userID uuid.UUID) ([]domain.Category, error)
	CreateCategory(ctx context.Context, userID uuid.UUID, c CreateCategoryInput) (*domain.Category, error)
}

type CreateCategoryInput struct {
	Name  string
	Kind  string
	Icon  string
	Color string
}

type GoalStore interface {
	ListGoalsByUser(ctx context.Context, userID uuid.UUID) ([]domain.Goal, error)
	CreateGoal(ctx context.Context, userID uuid.UUID, g CreateGoalInput) (*domain.Goal, error)
	Contribute(ctx context.Context, userID, goalID uuid.UUID, amount float64, note string) (*domain.Goal, error)
}

type CreateGoalInput struct {
	Title        string
	GoalType     string
	TargetAmount float64
	Currency     string
	Deadline     *time.Time
}

type UserStore interface {
	CreateUser(ctx context.Context, email, passwordHash, displayName string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (id uuid.UUID, passwordHash, displayName string, err error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	EnsureDefaultCategories(ctx context.Context, userID uuid.UUID) error
}
