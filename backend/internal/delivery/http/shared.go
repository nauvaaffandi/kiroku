package httpdelivery

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"kiroku/backend/internal/repository/sqlite"
)

type SharedHandler struct {
	repo *sqlite.SharedRepo
}

func NewSharedHandler(repo *sqlite.SharedRepo) *SharedHandler {
	return &SharedHandler{repo: repo}
}

func (h *SharedHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.CreateShare)
	r.Get("/", h.ListShares)
	r.Delete("/{id}", h.RevokeShare)
	r.Get("/public/{token}", h.GetSharedView)
	return r
}

func (h *SharedHandler) CreateShare(w http.ResponseWriter, r *http.Request) {
	userID := "system"
	var req struct {
		AccountID string  `json:"account_id"`
		ExpiresAt *string `json:"expires_at"`
	}
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.AccountID == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	var expiresAt *time.Time
	if req.ExpiresAt != nil {
		if t, err := time.Parse(time.RFC3339, *req.ExpiresAt); err == nil {
			expiresAt = &t
		}
	}
	share, err := h.repo.Create(r.Context(), userID, req.AccountID, expiresAt)
	if err != nil || share == nil {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":        share.ID,
		"token":     share.Token,
		"share_url": "/shared/" + share.Token,
	})
}

func (h *SharedHandler) ListShares(w http.ResponseWriter, r *http.Request) {
	userID := "system"
	shares, err := h.repo.List(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shares)
}

func (h *SharedHandler) RevokeShare(w http.ResponseWriter, r *http.Request) {
	userID := "system"
	id := chi.URLParam(r, "id")
	revoked, err := h.repo.Revoke(r.Context(), userID, id)
	if err != nil || !revoked {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "revoked"})
}

func (h *SharedHandler) GetSharedView(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	share, err := h.repo.GetByToken(r.Context(), token)
	if err != nil || share == nil {
		http.Error(w, "share not found", http.StatusNotFound)
		return
	}
	acc, _ := h.repo.GetAccountInfo(r.Context(), share.AccountID)
	txns, _ := h.repo.GetTransactions(r.Context(), share.AccountID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"account":      acc,
		"transactions": txns,
		"attribution":  "Recorded by you",
	})
}
