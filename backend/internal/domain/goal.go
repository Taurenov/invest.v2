package domain

import (
	"time"

	"github.com/google/uuid"
)

type Goal struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	Title         string    `json:"title"`
	GoalType      string    `json:"goal_type"`
	TargetAmount  float64   `json:"target_amount"`
	CurrentAmount float64   `json:"current_amount"`
	Currency      string    `json:"currency"`
	Deadline      *time.Time `json:"deadline,omitempty"`
}
