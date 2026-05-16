package domain

import "github.com/google/uuid"

type WatchlistItem struct {
	InstrumentID uuid.UUID `json:"instrument_id"`
	Symbol       string    `json:"symbol"`
	Exchange     string    `json:"exchange"`
	Name         string    `json:"name"`
	SortOrder    int       `json:"sort_order"`
}
