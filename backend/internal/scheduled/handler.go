package scheduled

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
	r.Post("/", h.Create)
	r.Post("/preset", h.CreatePreset)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	r.Post("/{id}/toggle", h.Toggle)
	return r
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	query := `SELECT id, rrule, max_occurrences, transaction_template_json, enabled, next_run, occurrences_done, created_at, updated_at FROM scheduled ORDER BY created_at DESC`
	rows, err := h.db.Query(query)
	if err != nil {
		log.Error().Err(err).Msg("scheduled: list")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var scheduled []models.Scheduled
	for rows.Next() {
		var s models.Scheduled
		var maxOcc sql.NullInt64
		var nextRun sql.NullTime
		if err := rows.Scan(&s.ID, &s.RRule, &maxOcc, &s.TransactionTemplateJSON, &s.Enabled, &nextRun, &s.OccurrencesDone, &s.CreatedAt, &s.UpdatedAt); err != nil {
			log.Error().Err(err).Msg("scheduled: scan")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if maxOcc.Valid {
			v := int(maxOcc.Int64)
			s.MaxOccurrences = &v
		}
		if nextRun.Valid {
			s.NextRun = &nextRun.Time
		}
		scheduled = append(scheduled, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scheduled)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	query := `SELECT id, rrule, max_occurrences, transaction_template_json, enabled, next_run, occurrences_done, created_at, updated_at FROM scheduled WHERE id = ?`
	var s models.Scheduled
	var maxOcc sql.NullInt64
	var nextRun sql.NullTime
	err := h.db.QueryRow(query, id).Scan(&s.ID, &s.RRule, &maxOcc, &s.TransactionTemplateJSON, &s.Enabled, &nextRun, &s.OccurrencesDone, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("scheduled: get")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if maxOcc.Valid {
		v := int(maxOcc.Int64)
		s.MaxOccurrences = &v
	}
	if nextRun.Valid {
		s.NextRun = &nextRun.Time
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	_, err := h.db.Exec("DELETE FROM scheduled WHERE id = ?", id)
	if err != nil {
		log.Error().Err(err).Msg("scheduled: delete")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) Toggle(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	_, err := h.db.Exec("UPDATE scheduled SET enabled = NOT enabled WHERE id = ?", id)
	if err != nil {
		log.Error().Err(err).Msg("scheduled: toggle")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
