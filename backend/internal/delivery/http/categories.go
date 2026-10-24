package httpdelivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"kiroku/backend/internal/domain"
)

type CategoriesHandler struct {
	repo domain.CategoryRepository
}

func NewCategoriesHandler(repo domain.CategoryRepository) *CategoriesHandler {
	return &CategoriesHandler{repo: repo}
}

func (h *CategoriesHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Get("/tree", h.Tree)
	r.Post("/", h.Create)
	r.Put("/reorder", h.Reorder)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

func (h *CategoriesHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.List(r.Context(), false)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if list == nil {
		list = []domain.Category{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *CategoriesHandler) Tree(w http.ResponseWriter, r *http.Request) {
	tree, err := h.repo.Tree(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tree)
}

func (h *CategoriesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateCategoryInput
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	if req.Kind != nil && req.ParentID != nil {
		children, _ := h.repo.CountChildren(r.Context(), *req.ParentID)
		if children > 0 {
			http.Error(w, "kind can only be set on leaf nodes", http.StatusBadRequest)
			return
		}
	}
	cat, err := h.repo.Create(r.Context(), req)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": cat.ID})
}

func (h *CategoriesHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	cat, err := h.repo.Get(r.Context(), id)
	if err != nil || cat == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cat)
}

func (h *CategoriesHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req domain.UpdateCategoryInput
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

func (h *CategoriesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *CategoriesHandler) Reorder(w http.ResponseWriter, r *http.Request) {
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
