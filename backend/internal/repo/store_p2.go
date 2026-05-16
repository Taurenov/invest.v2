package repo

import (
	"context"
	"time"

	"github.com/fin-helper/backend/internal/domain"
	"github.com/google/uuid"
)

type PortfolioStore interface {
	GetPortfolio(ctx context.Context, userID uuid.UUID) (*domain.PortfolioView, error)
	AddHolding(ctx context.Context, userID uuid.UUID, symbol, exchange string, qty, avgCost float64) (*domain.Holding, error)
	RemoveHolding(ctx context.Context, userID, holdingID uuid.UUID) error
}

type WatchlistStore interface {
	ListWatchlist(ctx context.Context, userID uuid.UUID) ([]domain.WatchlistItem, error)
	AddToWatchlist(ctx context.Context, userID uuid.UUID, symbol, exchange string) (*domain.WatchlistItem, error)
	RemoveFromWatchlist(ctx context.Context, userID uuid.UUID, instrumentID uuid.UUID) error
}

type SettingsStore interface {
	GetSettings(ctx context.Context, userID uuid.UUID) (*domain.UserSettings, error)
	UpdateSettings(ctx context.Context, userID uuid.UUID, s domain.UserSettings) (*domain.UserSettings, error)
}

type AnalyticsStore interface {
	Analytics(ctx context.Context, userID uuid.UUID, from, to time.Time) (*domain.AnalyticsReport, error)
}

type SummaryStore interface {
	GetOrCreateSummary(ctx context.Context, symbol, exchange string) (*domain.CompanySummary, error)
	EnsureInstrument(ctx context.Context, symbol, exchange, name string) (uuid.UUID, error)
}

type AddHoldingInput struct {
	Symbol   string
	Exchange string
	Quantity float64
	AvgCost  float64
}
