package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionKind string

const (
	KindIncome  TransactionKind = "income"
	KindExpense TransactionKind = "expense"
)

type Transaction struct {
	ID          uuid.UUID       `json:"id"`
	UserID      uuid.UUID       `json:"user_id"`
	CategoryID  *uuid.UUID      `json:"category_id,omitempty"`
	Kind        TransactionKind `json:"kind"`
	Amount      float64         `json:"amount"`
	Currency    string          `json:"currency"`
	Description string          `json:"description,omitempty"`
	OccurredAt  time.Time       `json:"occurred_at"`
}
