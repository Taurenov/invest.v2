package market

import (
	"context"

	"github.com/fin-helper/backend/internal/domain"
)

type Provider interface {
	Quote(ctx context.Context, symbol, exchange string) (*domain.Quote, error)
}
