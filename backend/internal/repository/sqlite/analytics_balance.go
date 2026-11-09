package sqlite

import (
	"context"
	"time"

	"kiroku/backend/internal/domain"
)

func (r *AnalyticsRepo) BalanceAccounts(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.AccountBalance, error) {
	var startDate, endDate time.Time
	if filter.DateFrom != nil {
		startDate = *filter.DateFrom
	} else {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}
	if filter.DateTo != nil {
		endDate = *filter.DateTo
	} else {
		endDate = startDate.AddDate(0, 1, 0).Add(-time.Second)
	}

	query := `
		SELECT
			a.id AS account_id,
			a.name AS account_name,
			a.currency,
			a.color,
			COALESCE((
				SELECT SUM(
					CASE
						WHEN t.type = 'income' AND t.from_account_id = a.id THEN t.amount
						WHEN t.type = 'income' AND t.to_account_id = a.id THEN t.amount
						WHEN t.type = 'expense' AND t.from_account_id = a.id THEN -t.amount
						ELSE 0
					END
				)
				FROM transactions t
				WHERE t.deleted_at IS NULL
				  AND t.date_time < ?
				  AND (t.from_account_id = a.id OR t.to_account_id = a.id)
			), 0) AS opening_balance,
			COALESCE((
				SELECT SUM(
					CASE
						WHEN t.type = 'income' AND t.from_account_id = a.id THEN t.amount
						WHEN t.type = 'income' AND t.to_account_id = a.id THEN t.amount
						WHEN t.type = 'expense' AND t.from_account_id = a.id THEN -t.amount
						ELSE 0
					END
				)
				FROM transactions t
				WHERE t.deleted_at IS NULL
				  AND t.date_time >= ? AND t.date_time <= ?
				  AND (t.from_account_id = a.id OR t.to_account_id = a.id)
			), 0) AS net_change,
			COALESCE((
				SELECT SUM(
					CASE
						WHEN t.type = 'income' AND t.from_account_id = a.id THEN t.amount
						WHEN t.type = 'income' AND t.to_account_id = a.id THEN t.amount
						WHEN t.type = 'expense' AND t.from_account_id = a.id THEN -t.amount
						ELSE 0
					END
				)
				FROM transactions t
				WHERE t.deleted_at IS NULL
				  AND t.date_time <= ?
				  AND (t.from_account_id = a.id OR t.to_account_id = a.id)
			), 0) AS closing_balance
		FROM accounts a
		WHERE a.archived = 0
		ORDER BY a.name
	`
	rows, err := r.db.QueryContext(ctx, query,
		startDate.Format("2006-01-02"),
		startDate.Format("2006-01-02"), endDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.AccountBalance
	for rows.Next() {
		var ab domain.AccountBalance
		if err := rows.Scan(&ab.AccountID, &ab.AccountName, &ab.Currency, &ab.Color, &ab.OpeningBalance, &ab.NetChange, &ab.ClosingBalance); err != nil {
			return nil, err
		}
		ab.Balance = ab.ClosingBalance
		result = append(result, ab)
	}

	return result, nil
}

func (r *AnalyticsRepo) BalanceCurrencies(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CurrencyBalanceSummary, error) {
	var startDate, endDate time.Time
	if filter.DateFrom != nil {
		startDate = *filter.DateFrom
	} else {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}
	if filter.DateTo != nil {
		endDate = *filter.DateTo
	} else {
		endDate = startDate.AddDate(0, 1, 0).Add(-time.Second)
	}

	query := `
		SELECT
			t.currency AS currency,
			COALESCE((
				SELECT SUM(
					CASE
						WHEN t2.type = 'income' THEN t2.amount
						WHEN t2.type = 'expense' THEN -t2.amount
						ELSE 0
					END
				)
				FROM transactions t2
				WHERE t2.deleted_at IS NULL
				  AND t2.currency = t.currency
				  AND t2.date_time <= ?
			), 0) AS closing_balance
		FROM transactions t
		WHERE t.deleted_at IS NULL
		  AND t.date_time >= ? AND t.date_time <= ?
		GROUP BY t.currency
	`
	rows, err := r.db.QueryContext(ctx, query,
		endDate.Format("2006-01-02"),
		startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.CurrencyBalanceSummary
	for rows.Next() {
		var cbs domain.CurrencyBalanceSummary
		if err := rows.Scan(&cbs.Currency, &cbs.Balance); err != nil {
			return nil, err
		}
		result = append(result, cbs)
	}

	return result, nil
}
