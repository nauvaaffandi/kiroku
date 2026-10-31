package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"time"

	"kiroku/backend/internal/db"
	"kiroku/backend/internal/domain"
)

type AnalyticsRepo struct {
	db *db.DB
}

func NewAnalyticsRepo(database *db.DB) *AnalyticsRepo {
	return &AnalyticsRepo{db: database}
}

func (r *AnalyticsRepo) CashFlowProjection(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.ProjectionPoint, []domain.ScheduledItem, error) {
	return nil, nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) Calendar(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CalendarDay, error) {
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

	query := `SELECT t.id, COALESCE(t.payee_payer, t.note, ''), t.type, t.amount, t.currency, t.date_time, c.name FROM transactions t LEFT JOIN categories c ON t.category_id = c.id WHERE t.deleted_at IS NULL AND t.date_time >= ? AND t.date_time <= ?`
	args := []interface{}{startDate.Format(time.RFC3339), endDate.Format(time.RFC3339)}

	if filter.AccountID != nil && *filter.AccountID != "" {
		query += ` AND (t.from_account_id = ? OR t.to_account_id = ?)`
		args = append(args, *filter.AccountID, *filter.AccountID)
	}
	query += ` ORDER BY t.date_time ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	daysMap := make(map[string][]domain.CalendarTxn)
	dayTotals := make(map[string]struct{ income, expense, net float64 })

	current := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDay := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())

	for !current.After(endDay) {
		dateStr := current.Format("2006-01-02")
		daysMap[dateStr] = []domain.CalendarTxn{}
		dayTotals[dateStr] = struct{ income, expense, net float64 }{}
		current = current.AddDate(0, 0, 1)
	}

	for rows.Next() {
		var id, name, typ, currency, dateTime string
		var amount float64
		var cat sql.NullString
		if err := rows.Scan(&id, &name, &typ, &amount, &currency, &dateTime, &cat); err != nil {
			return nil, err
		}

		tTime, err := time.Parse(time.RFC3339, dateTime)
		if err != nil {
			if parsed, parseErr := time.Parse("2006-01-02 15:04:05", dateTime); parseErr == nil {
				tTime = parsed
			} else {
				tTime, _ = time.Parse("2006-01-02", dateTime[:10])
			}
		}
		dateStr := tTime.Format("2006-01-02")

		categoryName := ""
		if cat.Valid {
			categoryName = cat.String
		}

		ctxn := domain.CalendarTxn{
			ID:       id,
			Name:     name,
			Amount:   amount,
			Type:     typ,
			Category: categoryName,
			Currency: currency,
		}

		if _, ok := daysMap[dateStr]; !ok {
			daysMap[dateStr] = []domain.CalendarTxn{}
		}
		daysMap[dateStr] = append(daysMap[dateStr], ctxn)

		tot := dayTotals[dateStr]
		if typ == "income" {
			tot.income += amount
			tot.net += amount
		} else if typ == "expense" {
			tot.expense += amount
			tot.net -= amount
		}
		dayTotals[dateStr] = tot
	}

	// Sort dates chronologically to ensure deterministic iteration and running balance calculation
	var sortedDates []string
	for dateStr := range daysMap {
		sortedDates = append(sortedDates, dateStr)
	}
	sort.Strings(sortedDates)

	var result []domain.CalendarDay
	for _, dateStr := range sortedDates {
		txns := daysMap[dateStr]
		tot := dayTotals[dateStr]
		result = append(result, domain.CalendarDay{
			Date:         dateStr,
			Income:       tot.income,
			Expense:      tot.expense,
			Net:          tot.net,
			Transactions: txns,
		})
	}

	return result, nil
}

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

	// Get opening balance
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

func (r *AnalyticsRepo) OutlookCategory(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CategoryStats, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) OutlookTimeline(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.TimelinePoint, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) OutlookForecast(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.ForecastPoint, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) BudgetSummary(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.BudgetItem, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) BudgetVsActual(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CategoryPlanned, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) GoalsProgress(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.GoalItem, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) GoalsProjection(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.ProjectionPoint, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) StatsSummary(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CurrencyStats, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) StatsByCategory(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CategoryStats, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) StatsByCurrency(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CurrencyStats, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) StatsByAccount(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.AccountStats, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) StatsTopMerchants(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.MerchantStats, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) StatsCashflow(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CashFlowSeriesPoint, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) StatsNetWorth(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.NetWorthAccount, error) {
	return nil, fmt.Errorf("not implemented")
}