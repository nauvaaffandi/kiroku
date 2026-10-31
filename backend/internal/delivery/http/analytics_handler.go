package httpdelivery

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"kiroku/backend/internal/domain"
)

type AnalyticsHandler struct {
	repo domain.AnalyticsRepository
}

func NewAnalyticsHandler(repo domain.AnalyticsRepository) *AnalyticsHandler {
	return &AnalyticsHandler{repo: repo}
}

func (h *AnalyticsHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/calendar", h.Calendar)
	r.Get("/balance/trend", h.BalanceTrend)
	r.Get("/balance/by-accounts", h.BalanceAccounts)
	r.Get("/balance/by-currencies", h.BalanceCurrencies)
	return r
}

func (h *AnalyticsHandler) Calendar(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	yearStr := q.Get("year")
	monthStr := q.Get("month")

	year := time.Now().Year()
	month := int(time.Now().Month())

	if yearStr != "" {
		y, err := strconv.Atoi(yearStr)
		if err != nil || y <= 0 {
			http.Error(w, "invalid year parameter", http.StatusBadRequest)
			return
		}
		year = y
	}
	if monthStr != "" {
		m, err := strconv.Atoi(monthStr)
		if err != nil || m <= 0 || m > 12 {
			http.Error(w, "invalid month parameter", http.StatusBadRequest)
			return
		}
		month = m
	}

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	filter := domain.AnalyticsFilter{
		DateFrom: &startDate,
		DateTo:   &endDate,
	}

	if aid := q.Get("account_id"); aid != "" {
		filter.AccountID = &aid
	}

	days, err := h.repo.Calendar(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sort.Slice(days, func(i, j int) bool {
		return days[i].Date < days[j].Date
	})

	var backendDays []domain.CalendarDayBackend
	var runningBalance float64

	for _, d := range days {
		var txns []domain.CalendarTxn
		for _, t := range d.Transactions {
			txns = append(txns, domain.CalendarTxn{
				ID:       t.ID,
				Name:     t.Name,
				Type:     t.Type,
				Amount:   t.Amount,
				Category: t.Category,
				Currency: t.Currency,
			})
		}
		runningBalance += d.Net
		backendDays = append(backendDays, domain.CalendarDayBackend{
			Date:             d.Date,
			Transactions:     txns,
			Scheduled:        []domain.CalendarTxn{},
			BalanceDelta:     d.Net,
			ProjectedBalance: runningBalance,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"days": backendDays,
	})
}

func (h *AnalyticsHandler) BalanceTrend(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var startDate, endDate time.Time
	if dateFrom := q.Get("date_from"); dateFrom != "" {
		t, err := time.Parse("2006-01-02", dateFrom)
		if err != nil {
			http.Error(w, "invalid date_from parameter", http.StatusBadRequest)
			return
		}
		startDate = t
	}
	if dateTo := q.Get("date_to"); dateTo != "" {
		t, err := time.Parse("2006-01-02", dateTo)
		if err != nil {
			http.Error(w, "invalid date_to parameter", http.StatusBadRequest)
			return
		}
		endDate = t.AddDate(0, 0, 1).Add(-time.Second)
	} else {
		endDate = time.Now().AddDate(0, 0, 1).Add(-time.Second)
	}
	if startDate.IsZero() {
		startDate = time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, time.Local)
	}

	filter := domain.AnalyticsFilter{
		DateFrom: &startDate,
		DateTo:   &endDate,
	}

	series, err := h.repo.BalanceTrend(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type trendPoint struct {
		Date    string    `json:"date"`
		Balance float64   `json:"balance"`
	}
	points := make([]trendPoint, len(series))
	for i, p := range series {
		points[i] = trendPoint{
			Date:    p.Date.Format("2006-01-02"),
			Balance: p.Balance,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"series": points,
	})
}

func (h *AnalyticsHandler) BalanceAccounts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var startDate, endDate time.Time
	if dateFrom := q.Get("date_from"); dateFrom != "" {
		t, err := time.Parse("2006-01-02", dateFrom)
		if err != nil {
			http.Error(w, "invalid date_from parameter", http.StatusBadRequest)
			return
		}
		startDate = t
	}
	if dateTo := q.Get("date_to"); dateTo != "" {
		t, err := time.Parse("2006-01-02", dateTo)
		if err != nil {
			http.Error(w, "invalid date_to parameter", http.StatusBadRequest)
			return
		}
		endDate = t.AddDate(0, 0, 1).Add(-time.Second)
	} else {
		endDate = time.Now().AddDate(0, 0, 1).Add(-time.Second)
	}
	if startDate.IsZero() {
		startDate = time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, time.Local)
	}

	filter := domain.AnalyticsFilter{
		DateFrom: &startDate,
		DateTo:   &endDate,
	}

	accounts, err := h.repo.BalanceAccounts(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type accountSummary struct {
		ID             string   `json:"id"`
		Name           string   `json:"name"`
		Currency       string   `json:"currency"`
		OpeningBalance float64  `json:"opening_balance"`
		ClosingBalance float64  `json:"closing_balance"`
		NetChange      float64  `json:"net_change"`
		Color          string   `json:"color"`
	}
	summaries := make([]accountSummary, len(accounts))
	totalByCurrency := make(map[string]float64)
	for i, a := range accounts {
		summaries[i] = accountSummary{
			ID:             a.AccountID,
			Name:           a.AccountName,
			Currency:       a.Currency,
			OpeningBalance: a.OpeningBalance,
			ClosingBalance: a.ClosingBalance,
			NetChange:      a.NetChange,
			Color:          a.Color,
		}
		totalByCurrency[a.Currency] += a.ClosingBalance
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"accounts":        summaries,
		"total_by_currency": totalByCurrency,
	})
}

func (h *AnalyticsHandler) BalanceCurrencies(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var startDate, endDate time.Time
	if dateFrom := q.Get("date_from"); dateFrom != "" {
		t, err := time.Parse("2006-01-02", dateFrom)
		if err != nil {
			http.Error(w, "invalid date_from parameter", http.StatusBadRequest)
			return
		}
		startDate = t
	}
	if dateTo := q.Get("date_to"); dateTo != "" {
		t, err := time.Parse("2006-01-02", dateTo)
		if err != nil {
			http.Error(w, "invalid date_to parameter", http.StatusBadRequest)
			return
		}
		endDate = t.AddDate(0, 0, 1).Add(-time.Second)
	} else {
		endDate = time.Now().AddDate(0, 0, 1).Add(-time.Second)
	}
	if startDate.IsZero() {
		startDate = time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, time.Local)
	}

	filter := domain.AnalyticsFilter{
		DateFrom: &startDate,
		DateTo:   &endDate,
	}

	currencies, err := h.repo.BalanceCurrencies(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type currencySummary struct {
		Currency       string  `json:"currency"`
		OpeningBalance float64 `json:"opening_balance"`
		ClosingBalance float64 `json:"closing_balance"`
		NetChange      float64 `json:"net_change"`
		FXRateToBase   float64 `json:"fx_rate_to_base"`
		ConvertedTotal float64 `json:"converted_total"`
	}
	summaries := make([]currencySummary, len(currencies))
	for i, c := range currencies {
		summaries[i] = currencySummary{
			Currency:       c.Currency,
			OpeningBalance: 0,
			ClosingBalance: c.Balance,
			NetChange:      0,
			FXRateToBase:   1.0,
			ConvertedTotal: c.Balance,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"currencies": summaries,
	})
}
