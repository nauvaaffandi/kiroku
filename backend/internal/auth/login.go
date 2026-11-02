package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"kiroku/backend/internal/config"
	"kiroku/backend/internal/models"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := ValidatePIN(string(req.PIN)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userID, hash string
	err := h.db.QueryRow("SELECT id, pin_hash FROM users LIMIT 1").Scan(&userID, &hash)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		log.Error().Err(err).Msg("auth: get user")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	ok, err := VerifyPin(string(req.PIN), hash)
	if err != nil || !ok {
		http.Error(w, "invalid pin", http.StatusUnauthorized)
		return
	}

	accessTTL := config.JWTAccessTTL()
	accessToken, err := GenerateAccessToken(userID, config.JWTSecret(), accessTTL)
	if err != nil {
		log.Error().Err(err).Msg("auth: generate access token")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	refreshTTL := config.JWTRefreshTTL()
	refreshToken, err := GenerateRefreshToken(userID, config.JWTSecret(), refreshTTL)
	if err != nil {
		log.Error().Err(err).Msg("auth: generate refresh token")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	setAuthCookies(w, accessToken, refreshToken, accessTTL, refreshTTL)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	clearAuthCookies(w)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
