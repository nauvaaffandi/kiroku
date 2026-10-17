package auth

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/google/uuid"
    "github.com/rs/zerolog/log"
)

func (h *Handler) GenerateRegisterToken(w http.ResponseWriter, r *http.Request) {
    userID := GetUserID(r)
    if userID == "" {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    // Check if already bound
    var existingChatID string
    err := h.db.QueryRow(`
        SELECT telegram_chat_id FROM telegram_bindings WHERE user_id = ?
    `, userID).Scan(&existingChatID)

    if err == nil && existingChatID != "" {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "status":  "already_linked",
            "chat_id": existingChatID,
        })
        return
    }

    token := uuid.New().String()
    expiresAt := time.Now().Add(10 * time.Minute)
    _, err = h.db.Exec(`
        INSERT INTO password_change_tokens (id, user_id, token, type, expires_at)
        VALUES (?, ?, ?, 'register', ?)
    `, uuid.New().String(), userID, token, expiresAt.Format("2006-01-02 15:04:05"))
    if err != nil {
        log.Error().Err(err).Msg("auth: create register token")
        http.Error(w, "failed to generate token", http.StatusInternalServerError)
        return
    }

    botUsername := "nauva_kiroku_bot"
    if h.tgService != nil {
        botUsername = h.tgService.GetBotUsername()
    }

    telegramURL := fmt.Sprintf("https://t.me/%s?start=register_%s", botUsername, token)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "token":        token,
        "telegram_url": telegramURL,
        "bot_username": botUsername,
    })
}

func (h *Handler) GetTelegramStatus(w http.ResponseWriter, r *http.Request) {
    userID := GetUserID(r)
    if userID == "" {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    var chatID string
    err := h.db.QueryRow(`
        SELECT telegram_chat_id FROM telegram_bindings WHERE user_id = ?
    `, userID).Scan(&chatID)

    linked := err == nil && chatID != ""
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "linked":  linked,
        "chat_id": chatID,
    })
}