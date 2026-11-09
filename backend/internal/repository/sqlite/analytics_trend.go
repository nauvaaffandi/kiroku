package sqlite

import (
	"context"
	"time"

	"kiroku/backend/internal/domain"
)

func (r *AnalyticsRepo) BalanceTrend(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.BalanceTrendPoint, error) {
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
		WITH daily_net AS (
			SELECT
				date(t.date_time) AS date_key,
				SUM(
					CASE
						WHEN t.type = 'income' THEN t.amount / 1.0
						WHEN t.type = 'expense' THEN -t.amount / 1.0
						ELSE 0
					END
				) AS net_change
			FROM transactions t
			WHERE t.deleted_at IS NULL
			  AND t.date_time >= ? AND t.date_time <= ?
			  AND t.type IN ('income', 'expense')
			GROUP BY date(t.date_time)
		),
		opening AS (
			SELECT
				SUM(
					CASE
						WHEN t.type = 'income' THEN t.amount / 1.0
						WHEN t.type = 'expense' THEN -t.amount / 1.0
						ELSE 0
					END
				) AS opening_balance
			FROM transactions t
			WHERE t.deleted_at IS NULL
			  AND t.date_time < ?
		)
		SELECT date_key, net_change
		FROM daily_net
		ORDER BY date_key ASC
	`
	rows, err := r.db.QueryContext(ctx, query,
		startDate.Format("2006-01-02"), endDate.Format("2006-01-02"),
		startDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dailyTotals := make(map[string]float64)
	for rows.Next() {
		var dateKey string
		var netChange float64
		if err := rows.Scan(&dateKey, &netChange); err != nil {
			return nil, err
		}
		dailyTotals[dateKey] = netChange
	}

	var result []domain.BalanceTrendPoint
	var running float64

	openingRow := r.db.QueryRowContext(ctx, `
		SELECT SUM(
			CASE
				WHEN t.type = 'income' THEN t.amount / 1.0
				WHEN t.type = 'expense' THEN -t.amount / 1.0
				ELSE 0
			END
		)
		FROM transactions t
		WHERE t.deleted_at IS NULL AND t.date_time < ?
	`, startDate.Format("2006-01-02"))
	if err := openingRow.Scan(&running); err != nil {
		running = 0
	}

	current := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDay := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())

	for !current.After(endDay) {
		dateKey := current.Format("2006-01-02")
		running += dailyTotals[dateKey]
		result = append(result, domain.BalanceTrendPoint{
			Date:    time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, time.Local),
			Balance: running,
		})
		current = current.AddDate(0, 0, 1)
	}

	return result, nil
}
