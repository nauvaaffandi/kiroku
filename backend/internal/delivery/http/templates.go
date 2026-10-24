package httpdelivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"kiroku/backend/internal/domain"
)

type TemplatesHandler struct {
	repo domain.TemplateRepository
}

func NewTemplatesHandler(repo domain.TemplateRepository) *TemplatesHandler {
	return &TemplatesHandler{repo: repo}
}

func (h *TemplatesHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	r.Post("/{id}/apply", h.Apply)
	return r
}

func (h *TemplatesHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *TemplatesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateTemplateInput
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	t, err := h.repo.Create(r.Context(), req)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": t.ID})
}

func (h *TemplatesHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	t, err := h.repo.Get(r.Context(), id)
	if err != nil || t == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (h *TemplatesHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req domain.UpdateTemplateInput
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

func (h *TemplatesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *TemplatesHandler) Apply(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	t, err := h.repo.Get(r.Context(), id)
	if err != nil || t == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}
