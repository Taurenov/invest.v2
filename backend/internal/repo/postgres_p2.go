package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fin-helper/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) EnsureInstrument(ctx context.Context, symbol, exchange, name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := p.pool.QueryRow(ctx, `
		INSERT INTO instruments (symbol, exchange, name, currency)
		VALUES ($1, $2, $3, 'RUB')
		ON CONFLICT (symbol, exchange) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, symbol, exchange, name).Scan(&id)
	return id, err
}

func (p *Postgres) defaultPortfolioID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	err := p.pool.QueryRow(ctx, `
		SELECT id FROM portfolios WHERE user_id = $1 ORDER BY created_at LIMIT 1
	`, userID).Scan(&id)
	if err == pgx.ErrNoRows {
		err = p.pool.QueryRow(ctx, `
			INSERT INTO portfolios (user_id, name) VALUES ($1, 'Основной') RETURNING id
		`, userID).Scan(&id)
	}
	return id, err
}

func (p *Postgres) defaultWatchlistID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	err := p.pool.QueryRow(ctx, `
		SELECT id FROM watchlists WHERE user_id = $1 ORDER BY created_at LIMIT 1
	`, userID).Scan(&id)
	if err == pgx.ErrNoRows {
		err = p.pool.QueryRow(ctx, `
			INSERT INTO watchlists (user_id, name) VALUES ($1, 'Основной') RETURNING id
		`, userID).Scan(&id)
	}
	return id, err
}

func (p *Postgres) GetPortfolio(ctx context.Context, userID uuid.UUID) (*domain.PortfolioView, error) {
	pid, err := p.defaultPortfolioID(ctx, userID)
	if err != nil {
		return nil, err
	}
	var name string
	_ = p.pool.QueryRow(ctx, `SELECT name FROM portfolios WHERE id=$1`, pid).Scan(&name)

	rows, err := p.pool.Query(ctx, `
		SELECT h.id, h.instrument_id, i.symbol, i.exchange, i.name,
		       h.quantity::float8, h.avg_cost::float8, h.currency
		FROM portfolio_holdings h
		JOIN instruments i ON i.id = h.instrument_id
		WHERE h.portfolio_id = $1
	`, pid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	view := &domain.PortfolioView{ID: pid, Name: name}
	for rows.Next() {
		var h domain.Holding
		if err := rows.Scan(&h.ID, &h.InstrumentID, &h.Symbol, &h.Exchange, &h.Name,
			&h.Quantity, &h.AvgCost, &h.Currency); err != nil {
			return nil, err
		}
		h.CostBasis = h.Quantity * h.AvgCost
		view.Holdings = append(view.Holdings, h)
		view.TotalCost += h.CostBasis
	}
	return view, rows.Err()
}

func (p *Postgres) AddHolding(ctx context.Context, userID uuid.UUID, symbol, exchange string, qty, avgCost float64) (*domain.Holding, error) {
	pid, err := p.defaultPortfolioID(ctx, userID)
	if err != nil {
		return nil, err
	}
	instID, err := p.EnsureInstrument(ctx, symbol, exchange, symbol)
	if err != nil {
		return nil, err
	}
	var h domain.Holding
	err = p.pool.QueryRow(ctx, `
		INSERT INTO portfolio_holdings (portfolio_id, instrument_id, quantity, avg_cost, currency)
		VALUES ($1, $2, $3, $4, 'RUB')
		ON CONFLICT (portfolio_id, instrument_id) DO UPDATE
		SET quantity = portfolio_holdings.quantity + EXCLUDED.quantity,
		    avg_cost = (portfolio_holdings.avg_cost + EXCLUDED.avg_cost) / 2,
		    updated_at = now()
		RETURNING id, instrument_id, quantity::float8, avg_cost::float8, currency
	`, pid, instID, qty, avgCost).Scan(&h.ID, &h.InstrumentID, &h.Quantity, &h.AvgCost, &h.Currency)
	if err != nil {
		return nil, err
	}
	h.Symbol = symbol
	h.Exchange = exchange
	h.Name = symbol
	h.CostBasis = h.Quantity * h.AvgCost
	return &h, nil
}

func (p *Postgres) RemoveHolding(ctx context.Context, userID, holdingID uuid.UUID) error {
	pid, err := p.defaultPortfolioID(ctx, userID)
	if err != nil {
		return err
	}
	tag, err := p.pool.Exec(ctx, `
		DELETE FROM portfolio_holdings WHERE id=$1 AND portfolio_id=$2
	`, holdingID, pid)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (p *Postgres) ListWatchlist(ctx context.Context, userID uuid.UUID) ([]domain.WatchlistItem, error) {
	wid, err := p.defaultWatchlistID(ctx, userID)
	if err != nil {
		return nil, err
	}
	rows, err := p.pool.Query(ctx, `
		SELECT i.id, i.symbol, i.exchange, i.name, wi.sort_order
		FROM watchlist_items wi
		JOIN instruments i ON i.id = wi.instrument_id
		WHERE wi.watchlist_id = $1
		ORDER BY wi.sort_order, i.symbol
	`, wid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.WatchlistItem
	for rows.Next() {
		var w domain.WatchlistItem
		if err := rows.Scan(&w.InstrumentID, &w.Symbol, &w.Exchange, &w.Name, &w.SortOrder); err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	return out, rows.Err()
}

func (p *Postgres) AddToWatchlist(ctx context.Context, userID uuid.UUID, symbol, exchange string) (*domain.WatchlistItem, error) {
	wid, err := p.defaultWatchlistID(ctx, userID)
	if err != nil {
		return nil, err
	}
	instID, err := p.EnsureInstrument(ctx, symbol, exchange, symbol)
	if err != nil {
		return nil, err
	}
	_, err = p.pool.Exec(ctx, `
		INSERT INTO watchlist_items (watchlist_id, instrument_id, sort_order)
		VALUES ($1, $2, (SELECT COALESCE(MAX(sort_order),0)+1 FROM watchlist_items WHERE watchlist_id=$1))
		ON CONFLICT DO NOTHING
	`, wid, instID)
	if err != nil {
		return nil, err
	}
	return &domain.WatchlistItem{InstrumentID: instID, Symbol: symbol, Exchange: exchange, Name: symbol}, nil
}

func (p *Postgres) RemoveFromWatchlist(ctx context.Context, userID uuid.UUID, instrumentID uuid.UUID) error {
	wid, err := p.defaultWatchlistID(ctx, userID)
	if err != nil {
		return err
	}
	_, err = p.pool.Exec(ctx, `DELETE FROM watchlist_items WHERE watchlist_id=$1 AND instrument_id=$2`, wid, instrumentID)
	return err
}

func (p *Postgres) GetSettings(ctx context.Context, userID uuid.UUID) (*domain.UserSettings, error) {
	var s domain.UserSettings
	err := p.pool.QueryRow(ctx, `
		SELECT user_id, locale, base_currency, theme, timezone FROM user_settings WHERE user_id=$1
	`, userID).Scan(&s.UserID, &s.Locale, &s.BaseCurrency, &s.Theme, &s.Timezone)
	if err == pgx.ErrNoRows {
		s = domain.UserSettings{UserID: userID, Locale: "ru", BaseCurrency: "RUB", Theme: "dark", Timezone: "Europe/Moscow"}
		_, _ = p.pool.Exec(ctx, `INSERT INTO user_settings (user_id) VALUES ($1) ON CONFLICT DO NOTHING`, userID)
		return &s, nil
	}
	return &s, err
}

func (p *Postgres) UpdateSettings(ctx context.Context, userID uuid.UUID, in domain.UserSettings) (*domain.UserSettings, error) {
	var s domain.UserSettings
	err := p.pool.QueryRow(ctx, `
		INSERT INTO user_settings (user_id, locale, base_currency, theme, timezone)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE SET
			locale = EXCLUDED.locale,
			base_currency = EXCLUDED.base_currency,
			theme = EXCLUDED.theme,
			timezone = EXCLUDED.timezone
		RETURNING user_id, locale, base_currency, theme, timezone
	`, userID, in.Locale, in.BaseCurrency, in.Theme, in.Timezone).Scan(
		&s.UserID, &s.Locale, &s.BaseCurrency, &s.Theme, &s.Timezone,
	)
	return &s, err
}

func (p *Postgres) Analytics(ctx context.Context, userID uuid.UUID, from, to time.Time) (*domain.AnalyticsReport, error) {
	report := &domain.AnalyticsReport{
		From: from.Format("2006-01-02"),
		To:   to.Format("2006-01-02"),
	}

	rows, err := p.pool.Query(ctx, `
		SELECT to_char(date_trunc('month', occurred_at), 'YYYY-MM') AS m,
		       COALESCE(SUM(CASE WHEN kind='income' THEN amount ELSE 0 END), 0)::float8,
		       COALESCE(SUM(CASE WHEN kind='expense' THEN amount ELSE 0 END), 0)::float8
		FROM transactions
		WHERE user_id=$1 AND occurred_at >= $2 AND occurred_at < $3
		GROUP BY 1 ORDER BY 1
	`, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var b domain.AnalyticsBucket
		if err := rows.Scan(&b.Label, &b.Income, &b.Expense); err != nil {
			return nil, err
		}
		report.ByMonth = append(report.ByMonth, b)
	}

	rows2, err := p.pool.Query(ctx, `
		SELECT COALESCE(c.name, t.description, 'Прочее'), t.kind, SUM(t.amount)::float8
		FROM transactions t
		LEFT JOIN categories c ON c.id = t.category_id
		WHERE t.user_id=$1 AND t.occurred_at >= $2 AND t.occurred_at < $3
		GROUP BY 1, 2 ORDER BY 3 DESC
	`, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()
	for rows2.Next() {
		var c domain.CategoryStat
		if err := rows2.Scan(&c.Name, &c.Kind, &c.Total); err != nil {
			return nil, err
		}
		report.ByCategory = append(report.ByCategory, c)
	}
	return report, nil
}

func (p *Postgres) GetOrCreateSummary(ctx context.Context, symbol, exchange string) (*domain.CompanySummary, error) {
	instID, err := p.EnsureInstrument(ctx, symbol, exchange, symbol)
	if err != nil {
		return nil, err
	}

	var text string
	var metricsJSON []byte
	var expires time.Time
	err = p.pool.QueryRow(ctx, `
		SELECT summary_text, key_metrics, expires_at FROM company_summaries WHERE instrument_id=$1
	`, instID).Scan(&text, &metricsJSON, &expires)
	if err == nil && time.Now().Before(expires) {
		var metrics map[string]any
		_ = json.Unmarshal(metricsJSON, &metrics)
		return &domain.CompanySummary{Symbol: symbol, Exchange: exchange, SummaryText: text, KeyMetrics: metrics}, nil
	}

	metrics := map[string]any{
		"sector":   "Финансы",
		"exchange": exchange,
		"symbol":   symbol,
	}
	text = fmt.Sprintf(
		"%s (%s): краткий обзор. Компания торгуется на %s. Данные обновляются из публичных источников. "+
			"Для инвестиционных решений изучайте официальную отчётность эмитента.",
		symbol, exchange, exchange,
	)
	b, _ := json.Marshal(metrics)
	expires = time.Now().Add(24 * time.Hour)
	_, _ = p.pool.Exec(ctx, `
		INSERT INTO company_summaries (instrument_id, summary_text, key_metrics, expires_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (instrument_id) DO UPDATE SET summary_text=$2, key_metrics=$3, expires_at=$4, fetched_at=now()
	`, instID, text, b, expires)

	return &domain.CompanySummary{Symbol: symbol, Exchange: exchange, SummaryText: text, KeyMetrics: metrics}, nil
}
