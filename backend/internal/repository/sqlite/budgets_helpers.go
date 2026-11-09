package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"kiroku/backend/internal/domain"
)

func insertRelations(tx interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}, ctx context.Context, budgetID string, ids []string, table, column string) error {
	for _, id := range ids {
		if _, err := tx.ExecContext(ctx, `INSERT INTO `+table+` (budget_id, `+column+`) VALUES (?, ?)`, budgetID, id); err != nil {
			return err
		}
	}
	return nil
}

func oneTimeWindow(period string, startsAt, endsAt *string) (interface{}, interface{}, error) {
	if period != domain.BudgetPeriodOneTime {
		return nil, nil, nil
	}
	if startsAt == nil || endsAt == nil {
		return nil, nil, fmt.Errorf("one_time budget requires starts_at and ends_at")
	}
	s, err := time.Parse(time.RFC3339, *startsAt)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid starts_at: %w", err)
	}
	e, err := time.Parse(time.RFC3339, *endsAt)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid ends_at: %w", err)
	}
	if !s.Before(e) {
		return nil, nil, fmt.Errorf("starts_at must be before ends_at")
	}
	return s.Format(time.RFC3339), e.Format(time.RFC3339), nil
}
