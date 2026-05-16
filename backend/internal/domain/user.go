package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
}

type Category struct {
	ID       uuid.UUID `json:"id"`
	UserID   *uuid.UUID `json:"user_id,omitempty"`
	Name     string    `json:"name"`
	Kind     string    `json:"kind"`
	Icon     string    `json:"icon,omitempty"`
	Color    string    `json:"color,omitempty"`
	IsSystem bool      `json:"is_system"`
}
