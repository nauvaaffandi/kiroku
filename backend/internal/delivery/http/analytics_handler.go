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
