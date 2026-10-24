package httpdelivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"kiroku/backend/internal/domain"
)

type SettingsHandler struct {
	repo domain.SettingsRepository
}

func NewSettingsHandler(repo domain.SettingsRepository) *SettingsHandler {
	return &SettingsHandler{repo: repo}
}

func (h *SettingsHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.GetAll)
	r.Put("/", h.UpdateAll)
	r.Get("/{key}", h.GetOne)
	r.Put("/{key}", h.SetOne)
	return r
}

func (h *SettingsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	settings, err := h.repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defaults := map[string]string{
		"theme":                  "system",
		"language":               "ja",
		"base_currency":          "JPY",
		"date_format":            "YYYY-MM-DD",
		"number_format":          "ja-JP",
		"week_start":             "1",
		"backup_passphrase_hash": "",
		"totp_enabled":           "false",
	}
	for k, v := range defaults {
		if _, ok := settings[k]; !ok {
			settings[k] = v
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func (h *SettingsHandler) UpdateAll(w http.ResponseWriter, r *http.Request) {
	var input map[string]string
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := h.repo.UpdateAll(r.Context(), input); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *SettingsHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	val, err := h.repo.GetOne(r.Context(), key)
	if err != nil {
		defaults := map[string]string{
			"theme": "system", "language": "ja", "base_currency": "JPY",
			"date_format": "YYYY-MM-DD", "number_format": "ja-JP", "week_start": "1",
		}
		if v, ok := defaults[key]; ok {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"key": key, "value": v})
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"key": key, "value": val})
}

func (h *SettingsHandler) SetOne(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	var req struct{ Value string `json:"value"` }
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := h.repo.SetOne(r.Context(), key, req.Value); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
