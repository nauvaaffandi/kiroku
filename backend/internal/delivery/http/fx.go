package httpdelivery

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"kiroku/backend/internal/domain"
)

type FXHandler struct {
	repo domain.FXRepository
}

func NewFXHandler(repo domain.FXRepository) *FXHandler {
	return &FXHandler{repo: repo}
}

func (h *FXHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/rates", h.GetRates)
	r.Get("/convert", h.Convert)
	return r
}

func (h *FXHandler) GetRates(w http.ResponseWriter, r *http.Request) {
	base := r.URL.Query().Get("base")
	if base == "" {
		base = "JPY"
	}
	dateStr := r.URL.Query().Get("date")
	date := time.Now()
	if dateStr != "" && dateStr != "latest" {
		if t, err := time.Parse("2006-01-02", dateStr); err == nil {
			date = t
		}
	}
	quotes, err := h.repo.GetQuotes(r.Context(), base, date)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	ratesMap := map[string]float64{base: 1.0}
	for _, q := range quotes {
		ratesMap[q.Quote] = q.Rate
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"base":  base,
		"date":  date.Format("2006-01-02"),
		"rates": ratesMap,
	})
}

func (h *FXHandler) Convert(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	if err != nil || from == "" || to == "" {
		http.Error(w, "invalid parameters", http.StatusBadRequest)
		return
	}
	if from == to {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"from": from, "to": to, "amount": amount, "converted": amount,
		})
		return
	}
	rFrom, _ := h.repo.GetRate(r.Context(), "JPY", from, time.Now())
	rTo, _ := h.repo.GetRate(r.Context(), "JPY", to, time.Now())
	rate := 1.0
	if rFrom != nil && rTo != nil && *rFrom > 0 {
		rate = *rTo / *rFrom
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"from":      from,
		"to":        to,
		"amount":    amount,
		"converted": amount * rate,
	})
}
