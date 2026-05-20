package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ForecastRecord struct {
	CreatedAt          time.Time `json:"created_at"`
	HorizonDays        int       `json:"horizon_days"`
	PredictedChangePct float64   `json:"predicted_change_pct"`
	Confidence         float64   `json:"confidence"`
	ModelVersion       string    `json:"model_version"`
}

func (p *Postgres) SaveForecast(
	ctx context.Context,
	instrumentID uuid.UUID,
	userID *uuid.UUID,
	horizon int,
	changePct float64,
	confidence float64,
	narrative string,
	modelVersion string,
) error {
	_, err := p.pool.Exec(ctx, `
		INSERT INTO ai_forecasts (
			instrument_id, user_id, horizon_days, predicted_change_pct,
			confidence, narrative, model_version
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, instrumentID, userID, horizon, changePct, confidence, narrative, modelVersion)
	return err
}

func (p *Postgres) ListForecastHistory(
	ctx context.Context,
	instrumentID uuid.UUID,
	userID *uuid.UUID,
	limit int,
) ([]ForecastRecord, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var rows pgRows
	var err error
	if userID != nil {
		rows, err = p.pool.Query(ctx, `
			SELECT created_at, horizon_days, predicted_change_pct::float8, confidence::float8, model_version
			FROM ai_forecasts
			WHERE instrument_id = $1 AND (user_id = $2 OR user_id IS NULL)
			ORDER BY created_at DESC
			LIMIT $3
		`, instrumentID, *userID, limit)
	} else {
		rows, err = p.pool.Query(ctx, `
			SELECT created_at, horizon_days, predicted_change_pct::float8, confidence::float8, model_version
			FROM ai_forecasts
			WHERE instrument_id = $1
			ORDER BY created_at DESC
			LIMIT $2
		`, instrumentID, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ForecastRecord
	for rows.Next() {
		var r ForecastRecord
		if err := rows.Scan(&r.CreatedAt, &r.HorizonDays, &r.PredictedChangePct, &r.Confidence, &r.ModelVersion); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

type PricePoint struct {
	Time  time.Time `json:"time"`
	Close float64   `json:"close"`
}

func (p *Postgres) ListPricePointsByInstrument(ctx context.Context, instrumentID uuid.UUID, points int) ([]PricePoint, error) {
	if points <= 0 || points > 1000 {
		points = 200
	}
	rows, err := p.pool.Query(ctx, `
		SELECT time, close::float8
		FROM (
			SELECT time, close
			FROM market_prices
			WHERE instrument_id = $1
			ORDER BY time DESC
			LIMIT $2
		) q
		ORDER BY time ASC
	`, instrumentID, points)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []PricePoint
	for rows.Next() {
		var pnt PricePoint
		if err := rows.Scan(&pnt.Time, &pnt.Close); err != nil {
			return nil, err
		}
		out = append(out, pnt)
	}
	return out, rows.Err()
}

type pgRows interface {
	Next() bool
	Scan(dest ...any) error
	Close()
	Err() error
}
