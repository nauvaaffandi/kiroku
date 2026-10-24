package sqlite

import (
	"context"
	"database/sql"
	"time"

	"kiroku/backend/internal/domain"
	"kiroku/backend/internal/db"
)

type ReportsRepo struct{ db *db.DB }

func NewReportsRepo(database *db.DB) *ReportsRepo { return &ReportsRepo{db: database} }

func (r *ReportsRepo) Spending(ctx context.Context, filter domain.SpendingFilter, groupBy string) ([]domain.SpendingGroup, error) {
	period := filter.Period
	if period == "" {
		period = "this_month"
	}
	dateFrom, dateTo := periodToDates(period)
	if filter.DateFrom != nil {
		dateFrom = filter.DateFrom.Format(time.RFC3339)
	}
	if filter.DateTo != nil {
		dateTo = filter.DateTo.Format(time.RFC3339)
	}
	baseCurrency := filter.Currency
	if baseCurrency == "" {
		baseCurrency = "JPY"
	}

	args := []interface{}{dateFrom, dateTo}
	whereClause := "WHERE t.deleted_at IS NULL AND t.date_time >= ? AND t.date_time <= ? AND t.type = 'expense'"
	if filter.AccountID != nil && *filter.AccountID != "" {
		whereClause += " AND (t.from_account_id = ? OR t.to_account_id = ?)"
		args = append(args, *filter.AccountID, *filter.AccountID)
	}

	var query string
	if groupBy == "parent_category" {
		query = `SELECT COALESCE(pc.id, c.id), COALESCE(pc.name, c.name), COALESCE(pc.color, c.color), SUM(t.amount * COALESCE(r.rate, 1)), COUNT(*), ? FROM transactions t LEFT JOIN categories c ON t.category_id = c.id LEFT JOIN categories pc ON c.parent_id = pc.id LEFT JOIN fx_rates r ON t.currency = r.base AND r.quote = ? ` + whereClause + ` GROUP BY group_id, group_name, group_color ORDER BY total_amount DESC`
		args = append(args, baseCurrency, baseCurrency)
	} else {
		query = `SELECT c.id, c.name, c.color, c.kind, SUM(t.amount * COALESCE(r.rate, 1)), COUNT(*), ? FROM transactions t LEFT JOIN categories c ON t.category_id = c.id LEFT JOIN fx_rates r ON t.currency = r.base AND r.quote = ? ` + whereClause + ` GROUP BY group_id, group_name, group_color, kind ORDER BY total_amount DESC`
		args = append(args, baseCurrency, baseCurrency)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []domain.SpendingGroup
	for rows.Next() {
		var g domain.SpendingGroup
		var kind sql.NullString
		var color sql.NullString
		if err := rows.Scan(&g.ID, &g.Name, &color, &kind, &g.Amount, &g.Count, &g.Currency); err != nil {
			return nil, err
		}
		if color.Valid {
			g.Color = color.String
		}
		if kind.Valid {
			str := kind.String
			g.Kind = &str
		}
		groups = append(groups, g)
	}
	return groups, rows.Err()
}

func (r *ReportsRepo) Trend(ctx context.Context, filter domain.TrendFilter) ([]domain.TrendPoint, error) {
	months := filter.Months
	if months <= 0 || months > 60 {
		months = 12
	}
	baseCurrency := filter.Currency
	if baseCurrency == "" {
		baseCurrency = "JPY"
	}

	rates := r.getFxRates(ctx, baseCurrency)

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	start = start.AddDate(0, -months+1, 0)

	var series []domain.TrendPoint
	for i := 0; i < months; i++ {
		monthStart := start.AddDate(0, i, 0)
		monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Second)
		accountIDStr := ""
		if filter.AccountID != nil {
			accountIDStr = *filter.AccountID
		}
		income, expense := r.getMonthTotals(ctx, monthStart, monthEnd, accountIDStr, "", baseCurrency, rates)
		series = append(series, domain.TrendPoint{
			Month:   monthStart.Format("2006-01"),
			Income:  income,
			Expense: expense,
			Net:     income - expense,
		})
	}
	return series, nil
}

func (r *ReportsRepo) getFxRates(ctx context.Context, baseCurrency string) map[string]float64 {
	rates := map[string]float64{}
	rows, err := r.db.QueryContext(ctx, "SELECT base, rate FROM fx_rates WHERE quote = ?", baseCurrency)
	if err != nil {
		rates[baseCurrency] = 1.0
		return rates
	}
	defer rows.Close()
	for rows.Next() {
		var from string
		var rate float64
		if rows.Scan(&from, &rate) == nil {
			rates[from] = rate
		}
	}
	rates[baseCurrency] = 1.0
	return rates
}

func (r *ReportsRepo) getMonthTotals(ctx context.Context, monthStart, monthEnd time.Time, accountID, categoryID, baseCurrency string, rates map[string]float64) (float64, float64) {
	args := []interface{}{monthStart.Format(time.RFC3339), monthEnd.Format(time.RFC3339)}
	whereClause := "WHERE t.deleted_at IS NULL AND t.date_time >= ? AND t.date_time <= ?"
	if accountID != "" {
		whereClause += " AND (t.from_account_id = ? OR t.to_account_id = ?)"
		args = append(args, accountID, accountID)
	}
	if categoryID != "" {
		whereClause += " AND t.category_id = ?"
		args = append(args, categoryID)
	}
	query := `SELECT t.type, t.currency, SUM(t.amount) as total FROM transactions t ` + whereClause + ` GROUP BY t.type, t.currency`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, 0
	}
	defer rows.Close()

	var income, expense float64
	for rows.Next() {
		var typ, currency string
		var total float64
		if rows.Scan(&typ, &currency, &total) != nil {
			continue
		}
		rate := rates[currency]
		if rate == 0 {
			rate = 1
		}
		converted := total * rate
		if typ == "income" {
			income += converted
		} else if typ == "expense" {
			expense += converted
		}
	}
	return income, expense
}

func periodToDates(period string) (string, string) {
	now := time.Now()
	var start, end time.Time
	switch period {
	case "this_week":
		start = now.AddDate(0, 0, -int(now.Weekday()))
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 7).Add(-time.Second)
	case "this_month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, 0).Add(-time.Second)
	case "this_year":
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(1, 0, 0).Add(-time.Second)
	case "last_month":
		start = time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, 0).Add(-time.Second)
	case "last_year":
		start = time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(1, 0, 0).Add(-time.Second)
	default:
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, 0).Add(-time.Second)
	}
	return start.Format(time.RFC3339), end.Format(time.RFC3339)
}
