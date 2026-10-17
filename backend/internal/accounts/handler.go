package accounts

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	r.Get("/summary", h.Summary)
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
		SELECT
			a.id, a.name, a.currency, a.type, a.icon, a.color, a.archived,
			COALESCE(SUM(CASE WHEN t.to_account_id = a.id THEN t.amount ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN t.from_account_id = a.id THEN t.amount ELSE 0 END), 0) AS balance
		FROM accounts a
		LEFT JOIN transactions t ON (t.from_account_id = a.id OR t.to_account_id = a.id) AND t.deleted_at IS NULL
		WHERE a.archived = FALSE
		GROUP BY a.id
		ORDER BY a.created_at
	`

	rows, err := h.db.Query(query)
	if err != nil {
		log.Error().Err(err).Msg("accounts: list")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	accounts := []models.Account{}
	for rows.Next() {
		var a models.Account
		if err := rows.Scan(&a.ID, &a.Name, &a.Currency, &a.Type, &a.Icon, &a.Color, &a.Archived, &a.Balance); err != nil {
			log.Error().Err(err).Msg("accounts: scan")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		accounts = append(accounts, a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateAccountReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Currency == "" || req.Type == "" {
		http.Error(w, "name, currency, type required", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	_, err := h.db.Exec(
		`INSERT INTO accounts (id, name, currency, type, icon, color, archived) VALUES (?, ?, ?, ?, ?, ?, 0)`,
		id, req.Name, req.Currency, req.Type, req.Icon, req.Color,
	)
	if err != nil {
		log.Error().Err(err).Msg("accounts: create")
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
	query := `
		SELECT
			a.id, a.name, a.currency, a.type, a.icon, a.color, a.archived,
			COALESCE(SUM(CASE WHEN t.to_account_id = a.id THEN t.amount ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN t.from_account_id = a.id THEN t.amount ELSE 0 END), 0) AS balance
		FROM accounts a
		LEFT JOIN transactions t ON (t.from_account_id = a.id OR t.to_account_id = a.id) AND t.deleted_at IS NULL
		WHERE a.id = ? AND a.archived = FALSE
		GROUP BY a.id
	`

	var a models.Account
	err := h.db.QueryRow(query, id).Scan(&a.ID, &a.Name, &a.Currency, &a.Type, &a.Icon, &a.Color, &a.Archived, &a.Balance)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("accounts: get")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	var req models.UpdateAccountReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Currency != nil {
		var count int
		err := h.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE (from_account_id = ? OR to_account_id = ?) AND deleted_at IS NULL", id, id).Scan(&count)
		if err != nil {
			log.Error().Err(err).Msg("accounts: check currency change")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if count > 0 {
			http.Error(w, "cannot change currency: account has transactions", http.StatusBadRequest)
			return
		}
	}

	setParts := []string{}
	args := []interface{}{}
	if req.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *req.Name)
	}
	if req.Currency != nil {
		setParts = append(setParts, "currency = ?")
		args = append(args, *req.Currency)
	}
	if req.Type != nil {
		setParts = append(setParts, "type = ?")
		args = append(args, *req.Type)
	}
	if req.Icon != nil {
		setParts = append(setParts, "icon = ?")
		args = append(args, *req.Icon)
	}
	if req.Color != nil {
		setParts = append(setParts, "color = ?")
		args = append(args, *req.Color)
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
	_, err := h.db.Exec("UPDATE accounts SET "+join(setParts, ", ")+" WHERE id = ?", args...)
	if err != nil {
		log.Error().Err(err).Msg("accounts: update")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
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

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	_, err := h.db.Exec("UPDATE accounts SET archived = 1 WHERE id = ?", id)
	if err != nil {
		log.Error().Err(err).Msg("accounts: delete")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	query := `
		SELECT a.currency,
			COALESCE(SUM(CASE WHEN t.to_account_id = a.id THEN t.amount ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN t.from_account_id = a.id THEN t.amount ELSE 0 END), 0) AS balance
		FROM accounts a
		LEFT JOIN transactions t ON (t.from_account_id = a.id OR t.to_account_id = a.id) AND t.deleted_at IS NULL
		WHERE a.archived = FALSE
		GROUP BY a.currency
	`

	rows, err := h.db.Query(query)
	if err != nil {
		log.Error().Err(err).Msg("accounts: summary")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	summary := make(map[string]int64)
	for rows.Next() {
		var currency string
		var balance int64
		if err := rows.Scan(&currency, &balance); err != nil {
			log.Error().Err(err).Msg("accounts: summary scan")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		summary[currency] = balance
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
