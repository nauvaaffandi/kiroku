package httpdelivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"kiroku/backend/internal/domain"
	"kiroku/backend/internal/repository/sqlite"
)

type ReportsHandler struct {
	repo *sqlite.ReportsRepo
}

func NewReportsHandler(repo *sqlite.ReportsRepo) *ReportsHandler {
	return &ReportsHandler{repo: repo}
}

func (h *ReportsHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/spending", h.Spending)
	r.Get("/trend", h.Trend)
	return r
}

func (h *ReportsHandler) Spending(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := domain.SpendingFilter{
		Period:  q.Get("period"),
		GroupBy: q.Get("group_by"),
	}
	if aid := q.Get("account_id"); aid != "" {
		filter.AccountID = &aid
	}
	groups, err := h.repo.Spending(r.Context(), filter, filter.GroupBy)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"period": filter.Period,
		"groups": groups,
	})
}

func (h *ReportsHandler) Trend(w http.ResponseWriter, r *http.Request) {
	months, _ := strconv.Atoi(r.URL.Query().Get("months"))
	filter := domain.TrendFilter{
		Months: months,
	}
	if aid := r.URL.Query().Get("account_id"); aid != "" {
		filter.AccountID = &aid
	}
	series, err := h.repo.Trend(r.Context(), filter)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"series": series,
	})
}
