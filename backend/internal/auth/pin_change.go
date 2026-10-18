package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"kiroku/backend/internal/models"
)

func (h *Handler) ChangePin(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.ChangePINReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := ValidatePIN(req.CurrentPIN); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ValidatePIN(req.NewPIN); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var currentHash string
	err := h.db.QueryRow("SELECT pin_hash FROM users WHERE id = ?", userID).Scan(&currentHash)
	if err != nil {
		log.Error().Err(err).Msg("auth: get current pin hash")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	ok, err := VerifyPin(req.CurrentPIN, currentHash)
	if err != nil || !ok {
		http.Error(w, "invalid current pin", http.StatusUnauthorized)
		return
	}

	newHash, err := HashPin(req.NewPIN, h.cfg.Argon2Time, h.cfg.Argon2Memory, h.cfg.Argon2Threads)
	if err != nil {
		log.Error().Err(err).Msg("auth: hash new pin")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	_, err = h.db.Exec("UPDATE users SET pin_hash = ? WHERE id = ?", newHash, userID)
	if err != nil {
		log.Error().Err(err).Msg("auth: update pin")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) RequestPinChange(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.PinChangeRequestReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := ValidatePIN(req.NewPIN); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var chatID string
	err := h.db.QueryRow(`
		SELECT telegram_chat_id FROM telegram_bindings WHERE user_id = ?
	`, userID).Scan(&chatID)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Telegram not linked. Link your Telegram in settings first.", http.StatusBadRequest)
		} else {
			log.Error().Err(err).Msg("auth: check telegram binding")
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	token := uuid.New().String()
	expiresAt := time.Now().Add(5 * time.Minute)
	_, err = h.db.Exec(`
		INSERT INTO pin_change_tokens (id, user_id, token, type, expires_at)
		VALUES (?, ?, ?, 'changepassword', ?)
	`, uuid.New().String(), userID, token, expiresAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Error().Err(err).Msg("auth: create changepassword token")
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	if h.tgService != nil {
		msg := fmt.Sprintf(
			"🔑 PIN Change Token (valid 5 minutes):\n\n`%s`\n\nPaste this token in Kiroku web settings along with your new PIN.",
			token,
		)
		h.tgService.SendMessage(chatID, msg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "token_sent"})
}

func (h *Handler) VerifyPinChange(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.VerifyPinChangeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Token == "" {
		http.Error(w, "token required", http.StatusBadRequest)
		return
	}
	if err := ValidatePIN(req.NewPIN); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var dbUserID string
	err := h.db.QueryRow(`
		SELECT user_id FROM pin_change_tokens
		WHERE token = ? AND type = 'changepassword' AND expires_at > datetime('now')
	`, req.Token).Scan(&dbUserID)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		} else {
			log.Error().Err(err).Msg("auth: verify changepassword token")
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	if dbUserID != userID {
		http.Error(w, "Token does not belong to this user", http.StatusForbidden)
		return
	}

	newHash, err := HashPin(req.NewPIN, h.cfg.Argon2Time, h.cfg.Argon2Memory, h.cfg.Argon2Threads)
	if err != nil {
		log.Error().Err(err).Msg("auth: hash new pin")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	_, err = h.db.Exec("UPDATE users SET pin_hash = ? WHERE id = ?", newHash, userID)
	if err != nil {
		log.Error().Err(err).Msg("auth: update pin")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	_, err = h.db.Exec(`DELETE FROM pin_change_tokens WHERE token = ?`, req.Token)
	if err != nil {
		log.Warn().Err(err).Msg("auth: cleanup changepassword token")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) RequestPinReset(w http.ResponseWriter, r *http.Request) {
	var req models.RequestPinResetReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "email required", http.StatusBadRequest)
		return
	}

	var userID, chatID string
	err := h.db.QueryRow(`
		SELECT u.id, tb.telegram_chat_id
		FROM users u
		JOIN telegram_bindings tb ON u.id = tb.user_id
		WHERE u.email = ?
	`, req.Email).Scan(&userID, &chatID)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No account found with that email or Telegram not linked", http.StatusBadRequest)
		} else {
			log.Error().Err(err).Msg("auth: forgotpassword lookup")
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	token := uuid.New().String()
	expiresAt := time.Now().Add(10 * time.Minute)
	_, err = h.db.Exec(`
		INSERT INTO pin_change_tokens (id, user_id, token, type, expires_at)
		VALUES (?, ?, ?, 'forgotpassword', ?)
	`, uuid.New().String(), userID, token, expiresAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Error().Err(err).Msg("auth: create forgotpassword token")
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	if h.tgService != nil {
		msg := fmt.Sprintf(
			"🔐 PIN Reset Token (valid 10 minutes):\n\n`%s`\n\nUse this token in Kiroku web to reset your PIN.",
			token,
		)
		h.tgService.SendMessage(chatID, msg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "token_sent"})
}

func (h *Handler) VerifyPinReset(w http.ResponseWriter, r *http.Request) {
	var req models.VerifyPinResetReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Token == "" {
		http.Error(w, "token required", http.StatusBadRequest)
		return
	}
	if err := ValidatePIN(req.NewPIN); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userID string
	err := h.db.QueryRow(`
		SELECT user_id FROM pin_change_tokens
		WHERE token = ? AND type = 'forgotpassword' AND expires_at > datetime('now')
	`, req.Token).Scan(&userID)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		} else {
			log.Error().Err(err).Msg("auth: verify forgotpassword token")
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	newHash, err := HashPin(req.NewPIN, h.cfg.Argon2Time, h.cfg.Argon2Memory, h.cfg.Argon2Threads)
	if err != nil {
		log.Error().Err(err).Msg("auth: hash new pin")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	_, err = h.db.Exec("UPDATE users SET pin_hash = ? WHERE id = ?", newHash, userID)
	if err != nil {
		log.Error().Err(err).Msg("auth: update pin")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	_, err = h.db.Exec(`DELETE FROM pin_change_tokens WHERE token = ?`, req.Token)
	if err != nil {
		log.Warn().Err(err).Msg("auth: cleanup forgotpassword token")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
