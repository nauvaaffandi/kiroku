package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"kiroku/backend/internal/domain"
)

func (r *BudgetsRepo) Create(ctx context.Context, input domain.CreateBudgetInput, userID string) (*domain.Budget, error) {
	id := uuid.NewString()
	now := time.Now()
	amount := domain.Cents(input.Amount)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	startsAt, endsAt, err := oneTimeWindow(input.Period, input.StartsAt, input.EndsAt)
	if err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO budgets (id, user_id, name, period, amount, currency, account_id,
		           trending_over, over_budget, starts_at, ends_at, deleted_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, id, userID, input.Name, input.Period, amount, input.Currency, input.AccountID,
		input.TrendingOver, input.OverBudget, startsAt, endsAt, nil, now, now); err != nil {
		return nil, err
	}

	if err := insertRelations(tx, ctx, id, input.CategoryIDs, `budget_categories`, `category_id`); err != nil {
		return nil, err
	}
	if err := insertRelations(tx, ctx, id, input.LabelIDs, `budget_labels`, `label_id`); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.Get(ctx, id, userID)
}

func (r *BudgetsRepo) Update(ctx context.Context, id, userID string, input domain.UpdateBudgetInput) (*domain.Budget, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var owned bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM budgets WHERE id = ? AND user_id = ? AND deleted_at IS NULL)`, id, userID).Scan(&owned); err != nil {
		return nil, err
	}
	if !owned {
		return nil, sql.ErrNoRows
	}

	setParts := []string{}
	args := []interface{}{}

	if input.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *input.Name)
	}
	if input.Period != nil {
		if *input.Period == domain.BudgetPeriodOneTime {
			if input.StartsAt == nil || input.EndsAt == nil {
				tx.Rollback()
				return nil, fmt.Errorf("one_time budget requires starts_at and ends_at")
			}
			s, e, err := oneTimeWindow(*input.Period, input.StartsAt, input.EndsAt)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			setParts = append(setParts, "period = ?, starts_at = ?, ends_at = ?")
			args = append(args, *input.Period, s, e)
		} else {
			setParts = append(setParts, "period = ?, starts_at = NULL, ends_at = NULL")
			args = append(args, *input.Period)
		}
	}
	if input.Amount != nil {
		setParts = append(setParts, "amount = ?")
		args = append(args, domain.Cents(*input.Amount))
	}
	if input.Currency != nil {
		setParts = append(setParts, "currency = ?")
		args = append(args, *input.Currency)
	}
	if input.AccountID != nil {
		setParts = append(setParts, "account_id = ?")
		args = append(args, *input.AccountID)
	}
	if input.TrendingOver != nil {
		setParts = append(setParts, "trending_over = ?")
		args = append(args, *input.TrendingOver)
	}
	if input.OverBudget != nil {
		setParts = append(setParts, "over_budget = ?")
		args = append(args, *input.OverBudget)
	}

	if len(setParts) > 0 {
		args = append(args, id, userID)
		if _, err := tx.ExecContext(ctx, "UPDATE budgets SET "+strings.Join(setParts, ", ")+" WHERE id = ? AND user_id = ?", args...); err != nil {
			return nil, err
		}
	}

	if input.CategoryIDs != nil {
		if _, err := tx.ExecContext(ctx, `DELETE FROM budget_categories WHERE budget_id = ? AND budget_id IN (SELECT id FROM budgets WHERE id = ? AND user_id = ?)`, id, id, userID); err != nil {
			return nil, err
		}
		if err := insertRelations(tx, ctx, id, input.CategoryIDs, `budget_categories`, `category_id`); err != nil {
			return nil, err
		}
	}
	if input.LabelIDs != nil {
		if _, err := tx.ExecContext(ctx, `DELETE FROM budget_labels WHERE budget_id = ? AND budget_id IN (SELECT id FROM budgets WHERE id = ? AND user_id = ?)`, id, id, userID); err != nil {
			return nil, err
		}
		if err := insertRelations(tx, ctx, id, input.LabelIDs, `budget_labels`, `label_id`); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.Get(ctx, id, userID)
}

func (r *BudgetsRepo) Delete(ctx context.Context, id, userID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE budgets SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND user_id = ?`, id, userID)
	return err
}

var _ interface {
	Create(context.Context, domain.CreateBudgetInput, string) (*domain.Budget, error)
	Update(context.Context, string, string, domain.UpdateBudgetInput) (*domain.Budget, error)
	Delete(context.Context, string, string) error
} = (*BudgetsRepo)(nil)
