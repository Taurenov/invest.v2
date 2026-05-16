package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/fin-helper/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (p *Postgres) CreateTransaction(ctx context.Context, userID uuid.UUID, in CreateTransactionInput) (*domain.Transaction, error) {
	var t domain.Transaction
	var catID *uuid.UUID
	err := p.pool.QueryRow(ctx, `
		INSERT INTO transactions (user_id, category_id, kind, amount, currency, description, occurred_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, category_id, kind, amount::float8, currency, COALESCE(description,''), occurred_at
	`, userID, in.CategoryID, in.Kind, in.Amount, in.Currency, in.Description, in.OccurredAt).Scan(
		&t.ID, &t.UserID, &catID, &t.Kind, &t.Amount, &t.Currency, &t.Description, &t.OccurredAt,
	)
	if err != nil {
		return nil, err
	}
	t.CategoryID = catID
	return &t, nil
}

func (p *Postgres) UpdateTransaction(ctx context.Context, userID, txID uuid.UUID, in UpdateTransactionInput) (*domain.Transaction, error) {
	cur, err := p.getTx(ctx, userID, txID)
	if err != nil {
		return nil, err
	}
	if in.Kind != nil {
		cur.Kind = *in.Kind
	}
	if in.Amount != nil {
		cur.Amount = *in.Amount
	}
	if in.Currency != nil {
		cur.Currency = *in.Currency
	}
	if in.Description != nil {
		cur.Description = *in.Description
	}
	if in.OccurredAt != nil {
		cur.OccurredAt = *in.OccurredAt
	}
	if in.CategoryID != nil {
		cur.CategoryID = in.CategoryID
	}

	var t domain.Transaction
	var catID *uuid.UUID
	err = p.pool.QueryRow(ctx, `
		UPDATE transactions SET category_id=$1, kind=$2, amount=$3, currency=$4, description=$5, occurred_at=$6, updated_at=now()
		WHERE id=$7 AND user_id=$8
		RETURNING id, user_id, category_id, kind, amount::float8, currency, COALESCE(description,''), occurred_at
	`, cur.CategoryID, cur.Kind, cur.Amount, cur.Currency, cur.Description, cur.OccurredAt, txID, userID).Scan(
		&t.ID, &t.UserID, &catID, &t.Kind, &t.Amount, &t.Currency, &t.Description, &t.OccurredAt,
	)
	if err != nil {
		return nil, err
	}
	t.CategoryID = catID
	return &t, nil
}

func (p *Postgres) DeleteTransaction(ctx context.Context, userID, txID uuid.UUID) error {
	tag, err := p.pool.Exec(ctx, `DELETE FROM transactions WHERE id=$1 AND user_id=$2`, txID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (p *Postgres) getTx(ctx context.Context, userID, txID uuid.UUID) (*domain.Transaction, error) {
	var t domain.Transaction
	var catID *uuid.UUID
	err := p.pool.QueryRow(ctx, `
		SELECT id, user_id, category_id, kind, amount::float8, currency, COALESCE(description,''), occurred_at
		FROM transactions WHERE id=$1 AND user_id=$2
	`, txID, userID).Scan(&t.ID, &t.UserID, &catID, &t.Kind, &t.Amount, &t.Currency, &t.Description, &t.OccurredAt)
	if err != nil {
		return nil, err
	}
	t.CategoryID = catID
	return &t, nil
}

func (p *Postgres) ListCategoriesByUser(ctx context.Context, userID uuid.UUID) ([]domain.Category, error) {
	rows, err := p.pool.Query(ctx, `
		SELECT id, user_id, name, kind, COALESCE(icon,''), COALESCE(color,''), is_system
		FROM categories
		WHERE user_id IS NULL OR user_id = $1
		ORDER BY is_system DESC, kind, name
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Kind, &c.Icon, &c.Color, &c.IsSystem); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (p *Postgres) CreateCategory(ctx context.Context, userID uuid.UUID, in CreateCategoryInput) (*domain.Category, error) {
	var c domain.Category
	var uid *uuid.UUID = &userID
	err := p.pool.QueryRow(ctx, `
		INSERT INTO categories (user_id, name, kind, icon, color)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, name, kind, COALESCE(icon,''), COALESCE(color,''), is_system
	`, userID, in.Name, in.Kind, in.Icon, in.Color).Scan(
		&c.ID, &uid, &c.Name, &c.Kind, &c.Icon, &c.Color, &c.IsSystem,
	)
	if err != nil {
		return nil, err
	}
	c.UserID = uid
	return &c, nil
}

func (p *Postgres) CreateGoal(ctx context.Context, userID uuid.UUID, in CreateGoalInput) (*domain.Goal, error) {
	var g domain.Goal
	err := p.pool.QueryRow(ctx, `
		INSERT INTO goals (user_id, title, goal_type, target_amount, currency, deadline)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, title, goal_type, target_amount::float8, current_amount::float8, currency, deadline
	`, userID, in.Title, in.GoalType, in.TargetAmount, in.Currency, in.Deadline).Scan(
		&g.ID, &g.UserID, &g.Title, &g.GoalType, &g.TargetAmount, &g.CurrentAmount, &g.Currency, &g.Deadline,
	)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (p *Postgres) Contribute(ctx context.Context, userID, goalID uuid.UUID, amount float64, note string) (*domain.Goal, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var owner uuid.UUID
	err = tx.QueryRow(ctx, `SELECT user_id FROM goals WHERE id=$1`, goalID).Scan(&owner)
	if err != nil {
		return nil, err
	}
	if owner != userID {
		return nil, errors.New("forbidden")
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO goal_contributions (goal_id, amount, note) VALUES ($1, $2, $3)
	`, goalID, amount, note)
	if err != nil {
		return nil, err
	}

	var g domain.Goal
	err = tx.QueryRow(ctx, `
		UPDATE goals SET current_amount = current_amount + $1
		WHERE id = $2
		RETURNING id, user_id, title, goal_type, target_amount::float8, current_amount::float8, currency, deadline
	`, amount, goalID).Scan(
		&g.ID, &g.UserID, &g.Title, &g.GoalType, &g.TargetAmount, &g.CurrentAmount, &g.Currency, &g.Deadline,
	)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &g, nil
}

func (p *Postgres) CreateUser(ctx context.Context, email, passwordHash, displayName string) (*domain.User, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var u domain.User
	err = tx.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, display_name)
		VALUES ($1, $2, $3)
		RETURNING id, email, display_name, created_at
	`, email, passwordHash, displayName).Scan(&u.ID, &u.Email, &u.DisplayName, &u.CreatedAt)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO user_settings (user_id) VALUES ($1)
	`, u.ID)
	if err != nil {
		return nil, err
	}

	if err := seedCategoriesTx(ctx, tx, u.ID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &u, nil
}

func (p *Postgres) GetByEmail(ctx context.Context, email string) (uuid.UUID, string, string, error) {
	var id uuid.UUID
	var hash, name string
	err := p.pool.QueryRow(ctx, `
		SELECT id, password_hash, display_name FROM users WHERE email=$1
	`, email).Scan(&id, &hash, &name)
	return id, hash, name, err
}

func (p *Postgres) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var u domain.User
	err := p.pool.QueryRow(ctx, `
		SELECT id, email, display_name, created_at FROM users WHERE id=$1
	`, id).Scan(&u.ID, &u.Email, &u.DisplayName, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (p *Postgres) EnsureDefaultCategories(ctx context.Context, userID uuid.UUID) error {
	return seedCategoriesTx(ctx, p.pool, userID)
}

type pgExecer interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

func seedCategoriesTx(ctx context.Context, tx pgExecer, userID uuid.UUID) error {
	type cat struct{ name, kind, icon, color string }
	defaults := []cat{
		{"Зарплата", "income", "💼", "#22c55e"},
		{"Подработка", "income", "📈", "#3b82f6"},
		{"Продукты", "expense", "🛒", "#f59e0b"},
		{"Транспорт", "expense", "🚌", "#8b5cf6"},
		{"Жильё", "expense", "🏠", "#ef4444"},
		{"Развлечения", "expense", "🎬", "#ec4899"},
	}
	for _, c := range defaults {
		_, err := tx.Exec(ctx, `
			INSERT INTO categories (user_id, name, kind, icon, color, is_system)
			SELECT $1, $2, $3, $4, $5, true
			WHERE NOT EXISTS (
				SELECT 1 FROM categories WHERE user_id = $1 AND name = $2 AND kind = $3
			)
		`, userID, c.name, c.kind, c.icon, c.color)
		if err != nil {
			return fmt.Errorf("seed category %s: %w", c.name, err)
		}
	}
	return nil
}
