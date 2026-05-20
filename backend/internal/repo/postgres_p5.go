package repo

import (
	"context"
	"time"

	"github.com/fin-helper/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) ListByUserFiltered(ctx context.Context, userID uuid.UUID, q TransactionQuery) ([]domain.Transaction, error) {
	limit := q.Limit
	if limit <= 0 || limit > 500 {
		limit = 200
	}

	// Simple dynamic WHERE (safe: bind params only)
	sql := `
		SELECT DISTINCT t.id, t.user_id, t.category_id, t.kind, t.amount::float8, t.currency,
		       COALESCE(t.description, ''), t.occurred_at
		FROM transactions t
		LEFT JOIN transaction_tags tt ON tt.transaction_id = t.id
		WHERE t.user_id = $1
	`
	args := []any{userID}
	n := 1

	if q.Kind != "" {
		n++
		sql += " AND t.kind = $" + itoa(n)
		args = append(args, q.Kind)
	}
	if q.CategoryID != nil {
		n++
		sql += " AND t.category_id = $" + itoa(n)
		args = append(args, *q.CategoryID)
	}
	if q.Text != "" {
		n++
		sql += " AND (t.description ILIKE $" + itoa(n) + ")"
		args = append(args, "%"+q.Text+"%")
	}
	if q.From != nil {
		n++
		sql += " AND t.occurred_at >= $" + itoa(n)
		args = append(args, *q.From)
	}
	if q.To != nil {
		n++
		sql += " AND t.occurred_at < $" + itoa(n)
		args = append(args, *q.To)
	}
	if len(q.TagIDs) > 0 {
		n++
		sql += " AND tt.tag_id = ANY($" + itoa(n) + ")"
		args = append(args, q.TagIDs)
	}

	n++
	sql += " ORDER BY t.occurred_at DESC LIMIT $" + itoa(n)
	args = append(args, limit)

	rows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		var catID *uuid.UUID
		if err := rows.Scan(&t.ID, &t.UserID, &catID, &t.Kind, &t.Amount, &t.Currency, &t.Description, &t.OccurredAt); err != nil {
			return nil, err
		}
		t.CategoryID = catID
		out = append(out, t)
	}
	return out, rows.Err()
}

func (p *Postgres) ListTags(ctx context.Context, userID uuid.UUID) ([]domain.Tag, error) {
	rows, err := p.pool.Query(ctx, `
		SELECT id, user_id, name, COALESCE(color,''), created_at
		FROM tags WHERE user_id = $1
		ORDER BY name
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Tag
	for rows.Next() {
		var t domain.Tag
		if err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.Color, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (p *Postgres) CreateTag(ctx context.Context, userID uuid.UUID, name, color string) (*domain.Tag, error) {
	var t domain.Tag
	err := p.pool.QueryRow(ctx, `
		INSERT INTO tags (user_id, name, color)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, name, COALESCE(color,''), created_at
	`, userID, name, color).Scan(&t.ID, &t.UserID, &t.Name, &t.Color, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (p *Postgres) DeleteTag(ctx context.Context, userID, tagID uuid.UUID) error {
	_, err := p.pool.Exec(ctx, `DELETE FROM tags WHERE id=$1 AND user_id=$2`, tagID, userID)
	return err
}

func (p *Postgres) SetTransactionTags(ctx context.Context, userID, txID uuid.UUID, tagIDs []uuid.UUID) error {
	// ensure tx belongs to user
	var owner uuid.UUID
	if err := p.pool.QueryRow(ctx, `SELECT user_id FROM transactions WHERE id=$1`, txID).Scan(&owner); err != nil {
		return err
	}
	if owner != userID {
		return pgx.ErrNoRows
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM transaction_tags WHERE transaction_id=$1`, txID); err != nil {
		return err
	}
	for _, id := range tagIDs {
		_, err := tx.Exec(ctx, `
			INSERT INTO transaction_tags (transaction_id, tag_id)
			SELECT $1, id FROM tags WHERE id=$2 AND user_id=$3
		`, txID, id, userID)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (p *Postgres) ListBudgets(ctx context.Context, userID uuid.UUID) ([]domain.Budget, error) {
	rows, err := p.pool.Query(ctx, `
		SELECT id, user_id, category_id, period, amount::float8, currency, created_at
		FROM budgets
		WHERE user_id=$1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Budget
	for rows.Next() {
		var b domain.Budget
		if err := rows.Scan(&b.ID, &b.UserID, &b.CategoryID, &b.Period, &b.Amount, &b.Currency, &b.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (p *Postgres) BudgetStatusThisMonth(ctx context.Context, userID uuid.UUID) ([]domain.BudgetStatus, error) {
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	to := from.AddDate(0, 1, 0)

	rows, err := p.pool.Query(ctx, `
		SELECT b.id, b.user_id, b.category_id, b.period, b.amount::float8, b.currency, b.created_at,
		       COALESCE(SUM(t.amount), 0)::float8 AS spent
		FROM budgets b
		LEFT JOIN transactions t
		  ON t.user_id = b.user_id
		 AND t.kind = 'expense'
		 AND t.category_id = b.category_id
		 AND t.occurred_at >= $2 AND t.occurred_at < $3
		WHERE b.user_id = $1 AND b.period = 'monthly'
		GROUP BY b.id
		ORDER BY spent DESC
	`, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.BudgetStatus
	for rows.Next() {
		var b domain.Budget
		var spent float64
		if err := rows.Scan(&b.ID, &b.UserID, &b.CategoryID, &b.Period, &b.Amount, &b.Currency, &b.CreatedAt, &spent); err != nil {
			return nil, err
		}
		rem := b.Amount - spent
		var pct float64
		if b.Amount > 0 {
			pct = (spent / b.Amount) * 100
		}
		out = append(out, domain.BudgetStatus{
			Budget:    b,
			Spent:     spent,
			Remaining: rem,
			Percent:   pct,
		})
	}
	return out, rows.Err()
}

func (p *Postgres) UpsertBudget(ctx context.Context, userID, categoryID uuid.UUID, amount float64, currency string) (*domain.Budget, error) {
	if currency == "" {
		currency = "RUB"
	}
	var b domain.Budget
	err := p.pool.QueryRow(ctx, `
		INSERT INTO budgets (user_id, category_id, period, amount, currency)
		VALUES ($1, $2, 'monthly', $3, $4)
		ON CONFLICT (user_id, category_id, period)
		DO UPDATE SET amount=EXCLUDED.amount, currency=EXCLUDED.currency
		RETURNING id, user_id, category_id, period, amount::float8, currency, created_at
	`, userID, categoryID, amount, currency).Scan(&b.ID, &b.UserID, &b.CategoryID, &b.Period, &b.Amount, &b.Currency, &b.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (p *Postgres) DeleteBudget(ctx context.Context, userID, budgetID uuid.UUID) error {
	_, err := p.pool.Exec(ctx, `DELETE FROM budgets WHERE id=$1 AND user_id=$2`, budgetID, userID)
	return err
}

func (p *Postgres) ListRecurring(ctx context.Context, userID uuid.UUID) ([]domain.RecurringTransaction, error) {
	rows, err := p.pool.Query(ctx, `
		SELECT id, user_id, kind, category_id, amount::float8, currency, COALESCE(description,''),
		       schedule, day_of_month, day_of_week, next_run_at::text, is_active
		FROM recurring_transactions
		WHERE user_id=$1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.RecurringTransaction
	for rows.Next() {
		var r domain.RecurringTransaction
		if err := rows.Scan(&r.ID, &r.UserID, &r.Kind, &r.CategoryID, &r.Amount, &r.Currency, &r.Description,
			&r.Schedule, &r.DayOfMonth, &r.DayOfWeek, &r.NextRunAt, &r.IsActive); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (p *Postgres) CreateRecurring(ctx context.Context, userID uuid.UUID, in domain.RecurringTransaction) (*domain.RecurringTransaction, error) {
	if in.Currency == "" {
		in.Currency = "RUB"
	}
	if in.NextRunAt == "" {
		in.NextRunAt = time.Now().Format("2006-01-02")
	}
	var r domain.RecurringTransaction
	err := p.pool.QueryRow(ctx, `
		INSERT INTO recurring_transactions (
			user_id, kind, category_id, amount, currency, description,
			schedule, day_of_month, day_of_week, next_run_at, is_active
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING id, user_id, kind, category_id, amount::float8, currency, COALESCE(description,''),
		          schedule, day_of_month, day_of_week, next_run_at::text, is_active
	`, userID, in.Kind, in.CategoryID, in.Amount, in.Currency, in.Description, in.Schedule, in.DayOfMonth, in.DayOfWeek, in.NextRunAt, in.IsActive).Scan(
		&r.ID, &r.UserID, &r.Kind, &r.CategoryID, &r.Amount, &r.Currency, &r.Description, &r.Schedule, &r.DayOfMonth, &r.DayOfWeek, &r.NextRunAt, &r.IsActive,
	)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (p *Postgres) ToggleRecurring(ctx context.Context, userID, id uuid.UUID, active bool) error {
	_, err := p.pool.Exec(ctx, `UPDATE recurring_transactions SET is_active=$1 WHERE id=$2 AND user_id=$3`, active, id, userID)
	return err
}

func (p *Postgres) DeleteRecurring(ctx context.Context, userID, id uuid.UUID) error {
	_, err := p.pool.Exec(ctx, `DELETE FROM recurring_transactions WHERE id=$1 AND user_id=$2`, id, userID)
	return err
}

func (p *Postgres) RunRecurringDue(ctx context.Context, now time.Time) (int, error) {
	today := now.Format("2006-01-02")

	rows, err := p.pool.Query(ctx, `
		SELECT id, user_id, kind, category_id, amount::float8, currency, COALESCE(description,''),
		       schedule, day_of_month, day_of_week, next_run_at::text, is_active
		FROM recurring_transactions
		WHERE is_active = true AND next_run_at <= $1
		ORDER BY next_run_at ASC
		LIMIT 200
	`, today)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	type item struct {
		ID         uuid.UUID
		UserID     uuid.UUID
		Kind       string
		CategoryID *uuid.UUID
		Amount     float64
		Currency   string
		Desc       string
		Schedule   string
		DayOfMonth *int
		DayOfWeek  *int
		NextRunAt  string
	}
	var due []item
	for rows.Next() {
		var it item
		var isActive bool
		if err := rows.Scan(&it.ID, &it.UserID, &it.Kind, &it.CategoryID, &it.Amount, &it.Currency, &it.Desc,
			&it.Schedule, &it.DayOfMonth, &it.DayOfWeek, &it.NextRunAt, &isActive); err != nil {
			return 0, err
		}
		due = append(due, it)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	if len(due) == 0 {
		return 0, nil
	}

	applied := 0
	for _, it := range due {
		// create transaction at noon local time
		occ := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
		kind := domain.TransactionKind(it.Kind)
		_, err := p.CreateTransaction(ctx, it.UserID, CreateTransactionInput{
			CategoryID:  it.CategoryID,
			Kind:        kind,
			Amount:      it.Amount,
			Currency:    it.Currency,
			Description: it.Desc,
			OccurredAt:  occ,
		})
		if err != nil {
			continue
		}

		next := computeNextDate(it.Schedule, it.NextRunAt, it.DayOfMonth, it.DayOfWeek)
		_, _ = p.pool.Exec(ctx, `UPDATE recurring_transactions SET next_run_at=$1 WHERE id=$2`, next, it.ID)
		applied++
	}

	return applied, nil
}

func computeNextDate(schedule, current string, dom, dow *int) string {
	t, err := time.Parse("2006-01-02", current)
	if err != nil {
		return time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	}
	switch schedule {
	case "daily":
		return t.AddDate(0, 0, 1).Format("2006-01-02")
	case "weekly":
		return t.AddDate(0, 0, 7).Format("2006-01-02")
	case "monthly":
		nt := t.AddDate(0, 1, 0)
		if dom != nil && *dom >= 1 && *dom <= 28 {
			nt = time.Date(nt.Year(), nt.Month(), *dom, 0, 0, 0, 0, nt.Location())
		}
		_ = dow // reserved
		return nt.Format("2006-01-02")
	default:
		return t.AddDate(0, 0, 1).Format("2006-01-02")
	}
}

func itoa(i int) string {
	// tiny helper, avoids strconv import in hot path
	if i == 0 {
		return "0"
	}
	var b [12]byte
	pos := len(b)
	n := i
	for n > 0 {
		pos--
		b[pos] = byte('0' + (n % 10))
		n /= 10
	}
	return string(b[pos:])
}

