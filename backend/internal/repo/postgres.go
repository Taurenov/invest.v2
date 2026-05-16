package repo

import (
	"context"
	"fmt"

	"github.com/fin-helper/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, databaseURL string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("pgx pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres ping: %w", err)
	}
	return &Postgres{pool: pool}, nil
}

func (p *Postgres) Close() {
	p.pool.Close()
}

func (p *Postgres) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Transaction, error) {
	rows, err := p.pool.Query(ctx, `
		SELECT id, user_id, category_id, kind, amount::float8, currency,
		       COALESCE(description, ''), occurred_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY occurred_at DESC
		LIMIT 200
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		var catID *uuid.UUID
		if err := rows.Scan(
			&t.ID, &t.UserID, &catID, &t.Kind, &t.Amount, &t.Currency,
			&t.Description, &t.OccurredAt,
		); err != nil {
			return nil, err
		}
		t.CategoryID = catID
		out = append(out, t)
	}
	return out, rows.Err()
}

func (p *Postgres) ListGoalsByUser(ctx context.Context, userID uuid.UUID) ([]domain.Goal, error) {
	rows, err := p.pool.Query(ctx, `
		SELECT id, user_id, title, goal_type,
		       target_amount::float8, current_amount::float8, currency, deadline
		FROM goals
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Goal
	for rows.Next() {
		var g domain.Goal
		if err := rows.Scan(
			&g.ID, &g.UserID, &g.Title, &g.GoalType,
			&g.TargetAmount, &g.CurrentAmount, &g.Currency, &g.Deadline,
		); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

func (p *Postgres) GetInstrument(ctx context.Context, symbol, exchange string) (*domain.Instrument, error) {
	var inst domain.Instrument
	err := p.pool.QueryRow(ctx, `
		SELECT id, symbol, exchange, name
		FROM instruments WHERE symbol = $1 AND exchange = $2
	`, symbol, exchange).Scan(&inst.ID, &inst.Symbol, &inst.Exchange, &inst.Name)
	if err != nil {
		return nil, err
	}
	return &inst, nil
}

func (p *Postgres) ListPrices(ctx context.Context, instrumentID uuid.UUID, limit int) ([]float64, error) {
	rows, err := p.pool.Query(ctx, `
		SELECT close::float8 FROM market_prices
		WHERE instrument_id = $1
		ORDER BY time ASC
		LIMIT $2
	`, instrumentID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []float64
	for rows.Next() {
		var v float64
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		prices = append(prices, v)
	}
	return prices, rows.Err()
}
