package transactions

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/models"
)

func (h *Handler) Batch(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.BatchReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if len(req.IDs) == 0 {
		http.Error(w, "ids required", http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		log.Error().Err(err).Msg("transactions: batch tx")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	affected := 0
	var errors []string

	for _, id := range req.IDs {
		switch req.Action {
		case "delete":
			res, err := tx.Exec("UPDATE transactions SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL", id)
			if err != nil {
				errors = append(errors, id+": "+err.Error())
				continue
			}
			n, _ := res.RowsAffected()
			if n > 0 {
				affected++
			}
		case "recategorize":
			res, err := tx.Exec("UPDATE transactions SET category_id = ? WHERE id = ? AND deleted_at IS NULL", req.CategoryID, id)
			if err != nil {
				errors = append(errors, id+": "+err.Error())
				continue
			}
			n, _ := res.RowsAffected()
			if n > 0 {
				affected++
			}
		case "relabel":
			labels := ""
			if req.Labels != nil {
				labels = *req.Labels
			}
			res, err := tx.Exec("UPDATE transactions SET labels = ? WHERE id = ? AND deleted_at IS NULL", labels, id)
			if err != nil {
				errors = append(errors, id+": "+err.Error())
				continue
			}
			n, _ := res.RowsAffected()
			if n > 0 {
				affected++
			}
		default:
			tx.Rollback()
			http.Error(w, "unknown action", http.StatusBadRequest)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("transactions: batch commit")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"affected": affected,
		"errors":   errors,
	})
}

func join(slice []string, sep string) string {
	result := ""
	for i, s := range slice {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
