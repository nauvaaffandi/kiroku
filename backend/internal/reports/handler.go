package reports

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/config"
	"kiroku/backend/internal/db"
	"kiroku/backend/internal/models"
)

type Handler struct {
	db  *db.DB
	cfg *config.Config
}

func NewHandler(database *db.DB, cfg *config.Config) *Handler {
	return &Handler{db: database, cfg: cfg}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/spending", h.Spending)
	r.Get("/trend", h.Trend)
	return r
}

func (h *Handler) Spending(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "this_month"
	}

	dateFrom, dateTo := h.periodToDates(period)
	if df := r.URL.Query().Get("date_from"); df != "" {
		dateFrom = df
	}
	if dt := r.URL.Query().Get("date_to"); dt != "" {
		dateTo = dt
	}

	groupBy := r.URL.Query().Get("group_by")
	if groupBy == "" {
		groupBy = "category"
	}

	accountID := r.URL.Query().Get("account_id")

	baseCurrency := h.cfg.FXBaseCurrency
	if baseCurrency == "" {
		baseCurrency = "JPY"
	}

	args := []interface{}{dateFrom, dateTo}
	if accountID != "" {
		args = append(args, accountID, accountID)
	}

	whereClause := "WHERE t.deleted_at IS NULL AND t.date_time >= ? AND t.date_time <= ?"
	if accountID != "" {
		whereClause += " AND (t.from_account_id = ? OR t.to_account_id = ?)"
	}
	whereClause += " AND t.type = 'expense'"

	var query string
	if groupBy == "parent_category" {
		query = `
			SELECT 
				COALESCE(pc.id, c.id) as group_id,
				COALESCE(pc.name, c.name) as group_name,
				COALESCE(pc.color, c.color) as group_color,
				SUM(t.amount * COALESCE(r.rate, 1)) as total_amount,
				COUNT(*) as count,
				? as currency
			FROM transactions t
			LEFT JOIN categories c ON t.category_id = c.id
			LEFT JOIN categories pc ON c.parent_id = pc.id
			LEFT JOIN fx_rates r ON t.currency = r.base AND r.quote = ?
			` + whereClause + `
			GROUP BY group_id, group_name, group_color
			ORDER BY total_amount DESC
		`
		args = append(args, baseCurrency, baseCurrency)
	} else {
		query = `
			SELECT 
				c.id as group_id,
				c.name as group_name,
				c.color as group_color,
				c.kind as kind,
				SUM(t.amount * COALESCE(r.rate, 1)) as total_amount,
				COUNT(*) as count,
				? as currency
			FROM transactions t
			LEFT JOIN categories c ON t.category_id = c.id
			LEFT JOIN fx_rates r ON t.currency = r.base AND r.quote = ?
			` + whereClause + `
			GROUP BY group_id, group_name, group_color, kind
			ORDER BY total_amount DESC
		`
		args = append(args, baseCurrency, baseCurrency)
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		log.Error().Err(err).Msg("reports: spending query")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var g models.Group
		var kind sql.NullString
		var color sql.NullString
		if err := rows.Scan(&g.ID, &g.Name, &color, &kind, &g.Amount, &g.Count, &g.Currency); err != nil {
			log.Error().Err(err).Msg("reports: spending scan")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if color.Valid {
			g.Color = color.String
		}
		if kind.Valid {
			g.Kind = kind.String
		}
		groups = append(groups, g)
	}

	if groupBy == "parent_category" {
		// Fetch children for each parent
		for i := range groups {
			h.fetchChildren(&groups[i], dateFrom, dateTo, accountID, baseCurrency)
		}
	}

	response := map[string]interface{}{
		"period": period,
		"groups": groups,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Trend(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	months := 12
	if m := r.URL.Query().Get("months"); m != "" {
		if v, err := strconv.Atoi(m); err == nil && v > 0 && v <= 60 {
			months = v
		}
	}

	accountID := r.URL.Query().Get("account_id")
	categoryID := r.URL.Query().Get("category_id")

	baseCurrency := h.cfg.FXBaseCurrency
	if baseCurrency == "" {
		baseCurrency = "JPY"
	}

	rates := h.getFxRates(baseCurrency)

	// Generate last N months
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	start = start.AddDate(0, -months+1, 0)

	var series []models.MonthData

	for i := 0; i < months; i++ {
		monthStart := start.AddDate(0, i, 0)
		monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Second)

		monthStr := monthStart.Format("2006-01")

		income, expense := h.getMonthTotals(monthStart, monthEnd, accountID, categoryID, baseCurrency, rates)

		series = append(series, models.MonthData{
			Month:    monthStr,
			Income:   income,
			Expense:  expense,
			Net:      income - expense,
			Currency: baseCurrency,
		})
	}

	response := map[string]interface{}{
		"series": series,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
