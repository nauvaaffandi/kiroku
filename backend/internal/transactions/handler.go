package transactions

import (
	"github.com/go-chi/chi/v5"

	"kiroku/backend/internal/db"
)

type Handler struct {
	db *db.DB
}

func NewHandler(database *db.DB) *Handler {
	return &Handler{db: database}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Post("/batch", h.Batch)
	r.Get("/search", h.Search)
	r.Get("/autocomplete", h.Autocomplete)
	r.Get("/summary", h.Summary)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}
