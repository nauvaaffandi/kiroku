package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"kiroku/backend/internal/domain"
)

func (r *BudgetsRepo) Calculate(ctx context.Context, budget *domain.Budget) (*domain.BudgetResult, error) {
	return r.budgetCalc(ctx, budget)
}

func (r *BudgetsRepo) budgetCalc(ctx context.Context, b *domain.Budget) (*domain.BudgetResult, error) {
	periodStart, periodEnd, err := periodWindow(b)
	if err != nil {
		return nil, err
	}

	accountCond := ""
	args := []interface{}{periodStart, periodEnd}
	if b.AccountID != nil {
		accountCond = `AND t.from_account_id = ?`
		args = append(args, *b.AccountID)
	}

	categoryCond := ""
	if len(b.Categories) > 0 {
		catPlaceholders := make([]string, len(b.Categories))
		for i, c := range b.Categories {
			catPlaceholders[i] = "?"
			args = append(args, c.ID)
		}
		categoryCond = `AND EXISTS (SELECT 1 FROM budget_categories bc WHERE bc.budget_id = ? AND bc.category_id = t.category_id)`
		args = append(args, b.ID)
	}

	labelCond := ""
	if len(b.Labels) > 0 {
		labelCond = `AND EXISTS (SELECT 1 FROM transaction_labels tl JOIN budget_labels bl ON bl.label_id = tl.label_id AND bl.budget_id = ? WHERE tl.transaction_id = t.id)`
		args = append(args, b.ID)
	}

	query := `
		SELECT t.amount, t.currency, t.date_time
		FROM transactions t
		WHERE t.type = 'expense'
		  AND t.deleted_at IS NULL
		  AND t.date_time >= ?
		  AND t.date_time < ?
		  ` + accountCond + categoryCond + labelCond

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	totalCents := int64(0)
	for rows.Next() {
		var amt int64
		var currency string
		var dt time.Time
		if err := rows.Scan(&amt, &currency, &dt); err != nil {
			return nil, err
		}
		if currency == b.Currency {
			totalCents += amt
			continue
		}
		rate, err := latestRate(ctx, r.db, currency, b.Currency, dt)
		if err != nil {
			return nil, err
		}
		if rate == nil {
			return nil, fmt.Errorf("missing exchange rate for %s to %s on %s", currency, b.Currency, dt.Format("2006-01-02"))
		}
		totalCents += int64(math.Round(float64(amt) * *rate))
	}

	spentStr := domain.FormatCents(totalCents)
	remainingCents := b.Amount - totalCents
	remainingStr := domain.FormatCents(remainingCents)
	pct := 0
	if b.Amount > 0 {
		pct = int(math.Round(float64(totalCents) / float64(b.Amount) * 100))
	}

	return &domain.BudgetResult{
		BudgetID:       b.ID,
		Name:           b.Name,
		Period:         b.Period,
		Amount:         b.AmountStr,
		Currency:       b.Currency,
		Spent:          spentStr,
		Remaining:      remainingStr,
		PercentageUsed: pct,
		IsOverBudget:   totalCents >= b.Amount,
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
	}, nil
}

func periodWindow(b *domain.Budget) (time.Time, time.Time, error) {
	now := time.Now()
	loc := now.Location()
	switch b.Period {
	case domain.BudgetPeriodWeekly:
		weekday := int(now.Weekday())
		daysSinceMonday := (weekday + 6) % 7
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, -daysSinceMonday)
		return start, start.AddDate(0, 0, 7), nil
	case domain.BudgetPeriodMonthly:
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		return start, start.AddDate(0, 1, 0), nil
	case domain.BudgetPeriodYearly:
		start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, loc)
		return start, start.AddDate(1, 0, 0), nil
	case domain.BudgetPeriodOneTime:
		if b.StartsAt == nil || b.EndsAt == nil || b.StartsAt.IsZero() || b.EndsAt.IsZero() {
			return time.Time{}, time.Time{}, nil
		}
		return b.StartsAt.In(loc), b.EndsAt.In(loc), nil
	}
	return time.Time{}, time.Time{}, fmt.Errorf("unsupported period: %s", b.Period)
}

func latestRate(ctx context.Context, db interface {
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}, base, quote string, date time.Time) (*float64, error) {
	dateStr := date.Format("2006-01-02")
	var rate float64
	err := db.QueryRowContext(ctx, `
		SELECT rate FROM fx_rates
		WHERE base = ? AND quote = ? AND date <= ?
		ORDER BY date DESC LIMIT 1
	`, base, quote, dateStr).Scan(&rate)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rate, nil
}
