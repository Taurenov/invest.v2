package domain

import (
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Budget struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	CategoryID uuid.UUID `json:"category_id"`
	Period     string    `json:"period"`
	Amount     float64   `json:"amount"`
	Currency   string    `json:"currency"`
	CreatedAt  time.Time `json:"created_at"`
}

type BudgetStatus struct {
	Budget    Budget  `json:"budget"`
	Spent     float64 `json:"spent"`
	Remaining float64 `json:"remaining"`
	Percent   float64 `json:"percent"`
}

type RecurringTransaction struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Kind        string    `json:"kind"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Description string    `json:"description,omitempty"`
	Schedule    string    `json:"schedule"` // daily/weekly/monthly
	DayOfMonth  *int      `json:"day_of_month,omitempty"`
	DayOfWeek   *int      `json:"day_of_week,omitempty"`
	NextRunAt   string    `json:"next_run_at"` // YYYY-MM-DD
	IsActive    bool      `json:"is_active"`
}

