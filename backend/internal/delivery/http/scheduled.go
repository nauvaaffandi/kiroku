package httpdelivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"kiroku/backend/internal/domain"
)

type ScheduledHandler struct {
	repo domain.ScheduledRepository
}

func NewScheduledHandler(repo domain.ScheduledRepository) *ScheduledHandler {
	return &ScheduledHandler{repo: repo}
}

func (h *ScheduledHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	r.Post("/{id}/toggle", h.Toggle)
	return r
}

func (h *ScheduledHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *ScheduledHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateScheduledInput
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	s, err := h.repo.Create(r.Context(), req)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": s.ID})
}

func (h *ScheduledHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	s, err := h.repo.Get(r.Context(), id)
	if err != nil || s == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func (h *ScheduledHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req domain.UpdateScheduledInput
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	_, err := h.repo.Update(r.Context(), id, req)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *ScheduledHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *ScheduledHandler) Toggle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.repo.Toggle(r.Context(), id, true)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
