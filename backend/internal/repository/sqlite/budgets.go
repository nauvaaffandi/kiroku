package sqlite

import (
	"context"
	"database/sql"

	"kiroku/backend/internal/db"
	"kiroku/backend/internal/domain"
)

type BudgetsRepo struct{ db *db.DB }

func NewBudgetsRepo(database *db.DB) *BudgetsRepo { return &BudgetsRepo{db: database} }

func (r *BudgetsRepo) UserIDsWithBudgets(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT DISTINCT user_id FROM budgets WHERE deleted_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		out = append(out, uid)
	}
	return out, rows.Err()
}

func (r *BudgetsRepo) List(ctx context.Context, userID string) ([]domain.Budget, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, name, period, amount, currency, account_id,
		       trending_over, over_budget, starts_at, ends_at,
		       deleted_at, created_at, updated_at
		FROM budgets
		WHERE user_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Budget
	for rows.Next() {
		b, err := r.scanBudget(rows)
		if err != nil {
			return nil, err
		}
		if err := r.loadRelations(ctx, &b); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	if out == nil {
		out = []domain.Budget{}
	}
	return out, rows.Err()
}

func (r *BudgetsRepo) Get(ctx context.Context, id, userID string) (*domain.Budget, error) {
	var b domain.Budget
	var amount int64
	var startsAt, endsAt, deletedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, name, period, amount, currency, account_id,
		       trending_over, over_budget, starts_at, ends_at,
		       deleted_at, created_at, updated_at
		FROM budgets
		WHERE id = ? AND user_id = ? AND deleted_at IS NULL
	`, id, userID).Scan(
		&b.ID, &b.UserID, &b.Name, &b.Period, &amount, &b.Currency, &b.AccountID,
		&b.TrendingOver, &b.OverBudget, &startsAt, &endsAt,
		&deletedAt, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	b.Amount = amount
	b.AmountStr = domain.FormatCents(amount)
	if startsAt.Valid {
		b.StartsAt = &startsAt.Time
	}
	if endsAt.Valid {
		b.EndsAt = &endsAt.Time
	}
	if deletedAt.Valid {
		t := deletedAt.Time
		b.DeletedAt = &t
	}
	if err := r.loadRelations(ctx, &b); err != nil {
		return nil, err
	}
	return &b, nil
}
