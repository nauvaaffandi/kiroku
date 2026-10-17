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

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.ChangePasswordReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

    if req.OldPassword == "" || req.NewPassword == "" {
        http.Error(w, "old_password and new_password required", http.StatusBadRequest)
        return
    }

    var currentHash string
    err := h.db.QueryRow("SELECT password_hash FROM users WHERE id = ?", userID).Scan(&currentHash)
    if err != nil {
        log.Error().Err(err).Msg("auth: get current password hash")
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    ok, err := VerifyPassword(req.OldPassword, currentHash)
    if err != nil || !ok {
        http.Error(w, "invalid old password", http.StatusUnauthorized)
        return
    }

    newHash, err := HashPassword(req.NewPassword, h.cfg.Argon2Time, h.cfg.Argon2Memory, h.cfg.Argon2Threads)
    if err != nil {
        log.Error().Err(err).Msg("auth: hash new password")
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    _, err = h.db.Exec("UPDATE users SET password_hash = ? WHERE id = ?", newHash, userID)
    if err != nil {
        log.Error().Err(err).Msg("auth: update password")
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

type passwordChangeRequestReq struct {
    NewPassword string `json:"new_password"`
}

func (h *Handler) RequestPasswordChange(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.PasswordChangeRequestReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    if req.NewPassword == "" {
        http.Error(w, "new_password required", http.StatusBadRequest)
        return
    }

    // Check if user has Telegram binding
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
        INSERT INTO password_change_tokens (id, user_id, token, type, expires_at)
        VALUES (?, ?, ?, 'changepassword', ?)
    `, uuid.New().String(), userID, token, expiresAt.Format("2006-01-02 15:04:05"))
    if err != nil {
        log.Error().Err(err).Msg("auth: create changepassword token")
        http.Error(w, "failed to generate token", http.StatusInternalServerError)
        return
    }

    // Send token via Telegram
    if h.tgService != nil {
        msg := fmt.Sprintf(
            "🔑 Password Change Token (valid 5 minutes):\n\n`%s`\n\nPaste this token in Kiroku web settings along with your new password.",
            token,
        )
        h.tgService.SendMessage(chatID, msg)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "token_sent"})
}

func (h *Handler) VerifyPasswordChange(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.VerifyPasswordChangeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    if req.Token == "" || req.NewPassword == "" {
        http.Error(w, "token and new_password required", http.StatusBadRequest)
        return
    }

    var dbUserID string
    err := h.db.QueryRow(`
        SELECT user_id FROM password_change_tokens
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

    newHash, err := HashPassword(req.NewPassword, h.cfg.Argon2Time, h.cfg.Argon2Memory, h.cfg.Argon2Threads)
    if err != nil {
        log.Error().Err(err).Msg("auth: hash new password")
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    _, err = h.db.Exec("UPDATE users SET password_hash = ? WHERE id = ?", newHash, userID)
    if err != nil {
        log.Error().Err(err).Msg("auth: update password")
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    _, err = h.db.Exec(`DELETE FROM password_change_tokens WHERE token = ?`, req.Token)
    if err != nil {
        log.Warn().Err(err).Msg("auth: cleanup changepassword token")
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req models.RequestPasswordResetReq
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
        INSERT INTO password_change_tokens (id, user_id, token, type, expires_at)
        VALUES (?, ?, ?, 'forgotpassword', ?)
    `, uuid.New().String(), userID, token, expiresAt.Format("2006-01-02 15:04:05"))
    if err != nil {
        log.Error().Err(err).Msg("auth: create forgotpassword token")
        http.Error(w, "failed to generate token", http.StatusInternalServerError)
        return
    }

    if h.tgService != nil {
        msg := fmt.Sprintf(
            "🔐 Password Reset Token (valid 10 minutes):\n\n`%s`\n\nUse this token in Kiroku web to reset your password.",
            token,
        )
        h.tgService.SendMessage(chatID, msg)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "token_sent"})
}

func (h *Handler) VerifyPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req models.VerifyPasswordResetReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    if req.Token == "" || req.NewPassword == "" {
        http.Error(w, "token and new_password required", http.StatusBadRequest)
        return
    }

    var userID string
    err := h.db.QueryRow(`
        SELECT user_id FROM password_change_tokens
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

    newHash, err := HashPassword(req.NewPassword, h.cfg.Argon2Time, h.cfg.Argon2Memory, h.cfg.Argon2Threads)
    if err != nil {
        log.Error().Err(err).Msg("auth: hash new password")
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    _, err = h.db.Exec("UPDATE users SET password_hash = ? WHERE id = ?", newHash, userID)
    if err != nil {
        log.Error().Err(err).Msg("auth: update password")
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    _, err = h.db.Exec(`DELETE FROM password_change_tokens WHERE token = ?`, req.Token)
    if err != nil {
        log.Warn().Err(err).Msg("auth: cleanup forgotpassword token")
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}