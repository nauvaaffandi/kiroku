package auth

import (
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"

    "kiroku/backend/internal/config"
    "kiroku/backend/internal/db"
    "kiroku/backend/internal/telegram"
)

type Handler struct {
    db        *db.DB
    cfg       *config.Config
    tgService *telegram.TelegramService
}

func NewHandler(database *db.DB, cfg *config.Config, tgService *telegram.TelegramService) *Handler {
    return &Handler{
        db:        database,
        cfg:       cfg,
        tgService: tgService,
    }
}

func (h *Handler) Routes() chi.Router {
    r := chi.NewRouter()
    r.Post("/login", h.Login)
    r.Post("/logout", h.Logout)
    r.Get("/me", h.Me)
    r.Post("/password", h.ChangePassword)
    r.Post("/password/change-request", h.RequestPasswordChange)
    r.Post("/password/verify", h.VerifyPasswordChange)
    r.Post("/password/reset-request", h.RequestPasswordReset)
    r.Post("/password/reset-verify", h.VerifyPasswordReset)
    r.Post("/telegram/register-token", h.GenerateRegisterToken)
    r.Get("/telegram/status", h.GetTelegramStatus)
    return r
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
    userID := GetUserID(r)
    if userID == "" {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"user_id": userID})
}