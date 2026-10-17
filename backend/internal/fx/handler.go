package fx

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/config"
	"kiroku/backend/internal/db"
)

type Handler struct {
	fetcher *Fetcher
}

func NewHandler(database *db.DB, cfg *config.Config) *Handler {
	base := cfg.FXBaseCurrency
	if base == "" {
		base = "JPY"
	}
	return &Handler{
		fetcher: NewFetcher(database, base),
	}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/rates", h.GetRates)
	r.Get("/convert", h.Convert)
	return r
}

func (h *Handler) GetRates(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	base := r.URL.Query().Get("base")
	date := r.URL.Query().Get("date")

	rates, actualDate, err := h.fetcher.GetRates(base, date)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if base == "" {
		base = h.fetcher.base
	}

	response := map[string]interface{}{
		"base":  base,
		"date":  actualDate,
		"rates": rates,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Convert(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	amountStr := r.URL.Query().Get("amount")

	if from == "" || to == "" || amountStr == "" {
		http.Error(w, "from, to, amount required", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	converted, err := h.fetcher.Convert(from, to, amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"from":      from,
		"to":        to,
		"amount":    amount,
		"converted": converted,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
