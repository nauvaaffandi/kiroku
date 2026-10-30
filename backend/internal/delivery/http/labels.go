package httpdelivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"kiroku/backend/internal/domain"
)

type LabelsHandler struct {
	repo domain.LabelRepository
}

func NewLabelsHandler(repo domain.LabelRepository) *LabelsHandler {
	return &LabelsHandler{repo: repo}
}

func (h *LabelsHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Put("/reorder", h.Reorder)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

func (h *LabelsHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if list == nil {
		list = []domain.Label{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *LabelsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateLabelInput
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	label, err := h.repo.Create(r.Context(), req)
	if err != nil {
		http.Error(w, "internal error or duplicate", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(label)
}

func (h *LabelsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	label, err := h.repo.Get(r.Context(), id)
	if err != nil || label == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(label)
}

func (h *LabelsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req domain.UpdateLabelInput
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	label, err := h.repo.Update(r.Context(), id, req)
	if err != nil || label == nil {
		http.Error(w, "internal error or not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(label)
}

func (h *LabelsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *LabelsHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Items []domain.ReorderItem `json:"items"`
	}
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := h.repo.Reorder(r.Context(), req.Items); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
