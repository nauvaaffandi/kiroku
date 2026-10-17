package categories

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/models"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateCategoryReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}

	if req.Kind != nil && req.ParentID != nil {
		// Check if parent exists and is leaf
		var childCount int
		err := h.db.QueryRow("SELECT COUNT(*) FROM categories WHERE parent_id = ?", *req.ParentID).Scan(&childCount)
		if err != nil || childCount > 0 {
			http.Error(w, "kind can only be set on leaf nodes (no children)", http.StatusBadRequest)
			return
		}
	}

	maxSortOrder := 0
	h.db.QueryRow("SELECT COALESCE(MAX(sort_order), 0) FROM categories WHERE parent_id IS ?", req.ParentID).Scan(&maxSortOrder)

	id := uuid.New().String()
	color := req.Color
	if color == "" {
		color = randomPastel()
	}

	_, err := h.db.Exec(
		`INSERT INTO categories (id, parent_id, name, icon, color, kind, sort_order) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, req.ParentID, req.Name, req.Icon, color, req.Kind, maxSortOrder+1,
	)
	if err != nil {
		log.Error().Err(err).Msg("categories: create")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	query := `SELECT id, parent_id, name, icon, color, kind, sort_order FROM categories WHERE id = ?`
	var c models.Category
	var parentID sql.NullString
	var kind sql.NullString
	err := h.db.QueryRow(query, id).Scan(&c.ID, &parentID, &c.Name, &c.Icon, &c.Color, &kind, &c.SortOrder)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("categories: get")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if parentID.Valid {
		c.ParentID = &parentID.String
	}
	if kind.Valid {
		c.Kind = &kind.String
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	var req models.UpdateCategoryReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.ParentID != nil {
		// Check for cycle: cannot move to self or descendant
		if *req.ParentID == id {
			http.Error(w, "cannot set parent to self", http.StatusBadRequest)
			return
		}
		if h.isDescendant(id, *req.ParentID) {
			http.Error(w, "cannot set parent to descendant", http.StatusBadRequest)
			return
		}
	}

	if req.Kind != nil {
		if req.ParentID != nil {
			var childCount int
			err := h.db.QueryRow("SELECT COUNT(*) FROM categories WHERE parent_id = ?", *req.ParentID).Scan(&childCount)
			if err != nil || childCount > 0 {
				http.Error(w, "kind can only be set on leaf nodes", http.StatusBadRequest)
				return
			}
		} else {
			var childCount int
			err := h.db.QueryRow("SELECT COUNT(*) FROM categories WHERE parent_id = ?", id).Scan(&childCount)
			if err != nil || childCount > 0 {
				http.Error(w, "kind can only be set on leaf nodes", http.StatusBadRequest)
				return
			}
		}
	}

	setParts := []string{}
	args := []interface{}{}
	if req.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *req.Name)
	}
	if req.ParentID != nil {
		setParts = append(setParts, "parent_id = ?")
		args = append(args, *req.ParentID)
	}
	if req.Icon != nil {
		setParts = append(setParts, "icon = ?")
		args = append(args, *req.Icon)
	}
	if req.Color != nil {
		setParts = append(setParts, "color = ?")
		args = append(args, *req.Color)
	}
	if req.Kind != nil {
		setParts = append(setParts, "kind = ?")
		args = append(args, *req.Kind)
	}
	if req.Archived != nil {
		setParts = append(setParts, "archived = ?")
		if *req.Archived {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
	}

	if len(setParts) == 0 {
		http.Error(w, "no fields to update", http.StatusBadRequest)
		return
	}

	args = append(args, id)
	_, err := h.db.Exec("UPDATE categories SET "+join(setParts, ", ")+" WHERE id = ?", args...)
	if err != nil {
		log.Error().Err(err).Msg("categories: update")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	// Archive this category and all descendants
	_, err := h.db.Exec(`
		UPDATE categories SET archived = 1
		WHERE id = ? OR id IN (
			WITH RECURSIVE descendants AS (
				SELECT id FROM categories WHERE parent_id = ?
				UNION ALL
				SELECT c.id FROM categories c
				JOIN descendants d ON c.parent_id = d.id
			)
			SELECT id FROM descendants
		)
	`, id, id)
	if err != nil {
		log.Error().Err(err).Msg("categories: delete")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) Reorder(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.ReorderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		log.Error().Err(err).Msg("categories: reorder tx")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	for _, item := range req.Items {
		_, err := tx.Exec("UPDATE categories SET sort_order = ? WHERE id = ?", item.SortOrder, item.ID)
		if err != nil {
			tx.Rollback()
			log.Error().Err(err).Msg("categories: reorder item")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("categories: reorder commit")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
