package market

import (
	"context"
	"fmt"
	"time"

	"github.com/fin-helper/backend/internal/domain"
)

type JSONCache interface {
	GetJSON(ctx context.Context, key string, dest any) (bool, error)
	SetJSON(ctx context.Context, key string, v any, ttl time.Duration) error
}

type Cached struct {
	Inner Provider
	Cache JSONCache
	TTL   time.Duration
}

func (c *Cached) Quote(ctx context.Context, symbol, exchange string) (*domain.Quote, error) {
	key := fmt.Sprintf("quote:%s:%s", exchange, symbol)
	var q domain.Quote
	if ok, err := c.Cache.GetJSON(ctx, key, &q); err != nil {
		return nil, err
	} else if ok {
		return &q, nil
	}

	qp, err := c.Inner.Quote(ctx, symbol, exchange)
	if err != nil {
		return nil, err
	}
	_ = c.Cache.SetJSON(ctx, key, qp, c.TTL)
	return qp, nil
}
