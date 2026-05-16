package market

import (
	"context"
	"math"
	"time"

	"github.com/fin-helper/backend/internal/domain"
)

// Mock — fallback без сети.
type Mock struct{}

func (Mock) Quote(_ context.Context, symbol, exchange string) (*domain.Quote, error) {
	t := float64(time.Now().Unix() % 3600)
	price := 280 + math.Sin(t/120)*5
	return &domain.Quote{
		Symbol:    symbol,
		Exchange:  exchange,
		Price:     price,
		ChangePct: math.Sin(t/200) * 2,
	}, nil
}
