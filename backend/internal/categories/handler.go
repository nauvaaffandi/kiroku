package categories

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/db"
	"kiroku/backend/internal/models"
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
	r.Get("/tree", h.Tree)
	r.Post("/", h.Create)
	r.Put("/reorder", h.Reorder)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	query := `
		SELECT id, parent_id, name, icon, color, kind, sort_order
		FROM categories
		WHERE archived = 0
		ORDER BY sort_order, name
	`

	rows, err := h.db.Query(query)
	if err != nil {
		log.Error().Err(err).Msg("categories: list")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		var parentID sql.NullString
		var kind sql.NullString
		if err := rows.Scan(&c.ID, &parentID, &c.Name, &c.Icon, &c.Color, &kind, &c.SortOrder); err != nil {
			log.Error().Err(err).Msg("categories: scan")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if parentID.Valid {
			c.ParentID = &parentID.String
		}
		if kind.Valid {
			c.Kind = &kind.String
		}
		categories = append(categories, c)
	}

	// Build depth and path
	catMap := make(map[string]*models.Category)
	for i := range categories {
		catMap[categories[i].ID] = &categories[i]
	}

	for i := range categories {
		h.computeDepthAndPath(&categories[i], catMap)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *Handler) Tree(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	query := `
		SELECT id, parent_id, name, icon, color, kind, sort_order
		FROM categories
		WHERE archived = 0
		ORDER BY sort_order, name
	`

	rows, err := h.db.Query(query)
	if err != nil {
		log.Error().Err(err).Msg("categories: tree")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	catMap := make(map[string]*models.Category)
	var roots []*models.Category

	for rows.Next() {
		var c models.Category
		var parentID sql.NullString
		var kind sql.NullString
		if err := rows.Scan(&c.ID, &parentID, &c.Name, &c.Icon, &c.Color, &kind, &c.SortOrder); err != nil {
			log.Error().Err(err).Msg("categories: tree scan")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if parentID.Valid {
			c.ParentID = &parentID.String
		}
		if kind.Valid {
			c.Kind = &kind.String
		}
		c.Children = []models.Category{}
		catMap[c.ID] = &c
	}

	for _, c := range catMap {
		if c.ParentID == nil {
			roots = append(roots, c)
		} else if parent, ok := catMap[*c.ParentID]; ok {
			parent.Children = append(parent.Children, *c)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roots)
}
