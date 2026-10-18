package auth

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"kiroku/backend/internal/models"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := ValidatePIN(req.PIN); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    userID, err := h.getOrCreateUser()
    if err != nil {
        log.Error().Err(err).Msg("auth: get/create user")
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    var hash string
    err = h.db.QueryRow("SELECT pin_hash FROM users WHERE id = ?", userID).Scan(&hash)
    if err != nil {
        log.Error().Err(err).Msg("auth: get pin hash")
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    ok, err := VerifyPin(req.PIN, hash)
    if err != nil || !ok {
        http.Error(w, "invalid pin", http.StatusUnauthorized)
        return
    }

    accessToken, err := GenerateAccessToken(userID, h.cfg.JWTSecret, h.cfg.JWTAccessTTL)
    if err != nil {
        log.Error().Err(err).Msg("auth: generate access token")
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    refreshToken, err := GenerateRefreshToken(userID, h.cfg.JWTSecret, h.cfg.JWTRefreshTTL)
    if err != nil {
        log.Error().Err(err).Msg("auth: generate refresh token")
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    setAuthCookies(w, accessToken, refreshToken, h.cfg.JWTAccessTTL, h.cfg.JWTRefreshTTL)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
    clearAuthCookies(w)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}