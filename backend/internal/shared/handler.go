package shared

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

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
	r.Post("/", h.CreateShare)
	r.Get("/", h.ListShares)
	r.Delete("/{id}", h.RevokeShare)
	r.Get("/public/{token}", h.GetSharedView)
	return r
}

func (h *Handler) CreateShare(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req models.CreateShareReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.AccountID == "" {
		http.Error(w, "account_id required", http.StatusBadRequest)
		return
	}

	// Verify account belongs to user
	var count int
	err := h.db.QueryRow("SELECT COUNT(*) FROM accounts WHERE id = ? AND user_id = ?", req.AccountID, userID).Scan(&count)
	if err != nil || count == 0 {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	var expiresAt sql.NullTime
	if req.ExpiresAt != nil {
		t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			http.Error(w, "invalid expires_at format (RFC3339)", http.StatusBadRequest)
			return
		}
		expiresAt = sql.NullTime{Time: t, Valid: true}
	}

	shareID := uuid.New().String()
	token := uuid.New().String()

	_, err = h.db.Exec(
		`INSERT INTO shares (id, account_id, token, expires_at) VALUES (?, ?, ?, ?)`,
		shareID, req.AccountID, token, expiresAt,
	)
	if err != nil {
		log.Error().Err(err).Msg("create share")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":        shareID,
		"token":     token,
		"share_url": "/shared/" + token,
	})
}

func (h *Handler) ListShares(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	rows, err := h.db.Query(`
		SELECT s.id, s.account_id, a.name, s.token, s.expires_at, s.created_at, s.last_accessed_at
		FROM shares s
		JOIN accounts a ON s.account_id = a.id
		WHERE a.user_id = ?
		ORDER BY s.created_at DESC
	`, userID)
	if err != nil {
		log.Error().Err(err).Msg("list shares")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var shares []models.Share
	for rows.Next() {
		var s models.Share
		var token string
		var expiresAt, lastAccessedAt sql.NullTime
		if err := rows.Scan(&s.ID, &s.AccountID, &s.AccountName, &token, &expiresAt, &s.CreatedAt, &lastAccessedAt); err != nil {
			log.Error().Err(err).Msg("scan share")
			continue
		}
		// Mask token for security
		if len(token) > 8 {
			s.Token = token[:8] + "..."
		} else {
			s.Token = token
		}
		if expiresAt.Valid {
			s.ExpiresAt = &expiresAt.Time
		}
		if lastAccessedAt.Valid {
			s.LastAccessedAt = &lastAccessedAt.Time
		}
		shares = append(shares, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shares)
}

func (h *Handler) RevokeShare(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	shareID := chi.URLParam(r, "id")

	// Verify share belongs to user's account
	result, err := h.db.Exec(`
		DELETE FROM shares
		WHERE id = ? AND account_id IN (SELECT id FROM accounts WHERE user_id = ?)
	`, shareID, userID)
	if err != nil {
		log.Error().Err(err).Msg("revoke share")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "share not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "revoked"})
}

func (h *Handler) GetSharedView(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	var shareID, accountID string
	var expiresAt sql.NullTime
	err := h.db.QueryRow(`
		SELECT id, account_id, expires_at FROM shares WHERE token = ?
	`, token).Scan(&shareID, &accountID, &expiresAt)
	if err == sql.ErrNoRows {
		http.Error(w, "share not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("get share by token")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if expiresAt.Valid && expiresAt.Time.Before(time.Now()) {
		http.Error(w, "share expired", http.StatusGone)
		return
	}

	// Update last accessed
	_, _ = h.db.Exec("UPDATE shares SET last_accessed_at = CURRENT_TIMESTAMP WHERE id = ?", shareID)

	var response models.SharedViewResponse

	// Get account info
	var balance float64
	err = h.db.QueryRow(`
		SELECT id, name, currency,
		       COALESCE((SELECT SUM(CASE WHEN type = 'income' THEN amount ELSE -amount END) FROM transactions WHERE account_id = a.id), 0) as balance
		FROM accounts a WHERE id = ?
	`, accountID).Scan(&response.Account.ID, &response.Account.Name, &response.Account.Currency, &balance)
	if err != nil {
		log.Error().Err(err).Msg("get shared account")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	response.Account.Balance = balance

	// Get transactions
	rows, err := h.db.Query(`
		SELECT t.id, t.date, t.amount, t.type, t.note, t.payee,
		       c.name, c.icon, c.color
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.account_id = ?
		ORDER BY t.date DESC, t.created_at DESC
		LIMIT 100
	`, accountID)
	if err != nil {
		log.Error().Err(err).Msg("get shared transactions")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tx struct {
			ID            string  `json:"id"`
			Date          string  `json:"date"`
			Amount        float64 `json:"amount"`
			Type          string  `json:"type"`
			Note          string  `json:"note"`
			Payee         string  `json:"payee"`
			Category      string  `json:"category"`
			CategoryIcon  string  `json:"category_icon"`
			CategoryColor string  `json:"category_color"`
		}
		var catName, catIcon, catColor sql.NullString
		if err := rows.Scan(&tx.ID, &tx.Date, &tx.Amount, &tx.Type, &tx.Note, &tx.Payee, &catName, &catIcon, &catColor); err != nil {
			continue
		}
		if catName.Valid {
			tx.Category = catName.String
		}
		if catIcon.Valid {
			tx.CategoryIcon = catIcon.String
		}
		if catColor.Valid {
			tx.CategoryColor = catColor.String
		}
		response.Transactions = append(response.Transactions, tx)
	}

	response.Attribution = "Recorded by you"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
