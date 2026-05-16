package domain

import "github.com/google/uuid"

type Holding struct {
	ID           uuid.UUID `json:"id"`
	InstrumentID uuid.UUID `json:"instrument_id"`
	Symbol       string    `json:"symbol"`
	Exchange     string    `json:"exchange"`
	Name         string    `json:"name"`
	Quantity     float64   `json:"quantity"`
	AvgCost      float64   `json:"avg_cost"`
	Currency     string    `json:"currency"`
	CurrentPrice float64   `json:"current_price,omitempty"`
	MarketValue  float64   `json:"market_value,omitempty"`
	CostBasis    float64   `json:"cost_basis,omitempty"`
	PnL          float64   `json:"pnl,omitempty"`
	PnLPercent   float64   `json:"pnl_percent,omitempty"`
}

type PortfolioView struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Holdings  []Holding `json:"holdings"`
	TotalCost float64   `json:"total_cost"`
	TotalValue float64  `json:"total_value"`
	TotalPnL  float64   `json:"total_pnl"`
}
