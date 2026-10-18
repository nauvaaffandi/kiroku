package auth

import (
    "database/sql"
    "net/http"
    "time"

    "github.com/google/uuid"
    "github.com/rs/zerolog/log"
)

// getOrCreateUser returns an existing user ID or creates one with a default PIN hash
func (h *Handler) getOrCreateUser() (string, error) {
    var userID string
    err := h.db.QueryRow("SELECT id FROM users LIMIT 1").Scan(&userID)
    if err == nil {
        return userID, nil
    }
    if err != sql.ErrNoRows {
        return "", err
    }

    hash, err := HashPin("123456", h.cfg.Argon2Time, h.cfg.Argon2Memory, h.cfg.Argon2Threads)
    if err != nil {
        return "", err
    }

    userID = uuid.New().String()
    _, err = h.db.Exec(
        "INSERT INTO users (id, email, pin_hash) VALUES (?, ?, ?)",
        userID, "user@kiroku.local", hash,
    )
    if err != nil {
        return "", err
    }

    return userID, nil
}

// setAuthCookies sets the access and refresh cookies on the response
func setAuthCookies(w http.ResponseWriter, accessToken, refreshToken string, accessTTL, refreshTTL time.Duration) {
    http.SetCookie(w, &http.Cookie{
        Name:     "access_token",
        Value:    accessToken,
        Path:     "/",
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteNoneMode,
        MaxAge:   int(accessTTL.Seconds()),
    })
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    refreshToken,
        Path:     "/",
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteNoneMode,
        MaxAge:   int(refreshTTL.Seconds()),
    })
}

// clearAuthCookies clears the access and refresh cookies
func clearAuthCookies(w http.ResponseWriter) {
    http.SetCookie(w, &http.Cookie{
        Name:     "access_token",
        Value:    "",
        Path:     "/",
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteNoneMode,
        MaxAge:   -1,
    })
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    "",
        Path:     "/",
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteNoneMode,
        MaxAge:   -1,
    })
}

// ensure log is used to avoid unused import
var _ = log.Error

// CheckAndResetPinHash verifies the default user's PIN hash on startup.
// If verification fails (e.g., due to Argon2 config change), it resets the hash.
func (h *Handler) CheckAndResetPinHash() {
	var userID, hash string
	err := h.db.QueryRow("SELECT id, pin_hash FROM users LIMIT 1").Scan(&userID, &hash)
	if err != nil {
		return // No user exists yet, will be created on first login
	}

	ok, err := VerifyPin("123456", hash)
	if err == nil && ok {
		return // PIN hash is valid
	}

	log.Warn().Msg("auth: default PIN hash invalid, resetting")

	newHash, err := HashPin("123456", h.cfg.Argon2Time, h.cfg.Argon2Memory, h.cfg.Argon2Threads)
	if err != nil {
		log.Error().Err(err).Msg("auth: failed to hash new PIN")
		return
	}

	_, err = h.db.Exec("UPDATE users SET pin_hash = ? WHERE id = ?", newHash, userID)
	if err != nil {
		log.Error().Err(err).Msg("auth: failed to reset PIN hash")
		return
	}

	log.Info().Msg("auth: PIN hash reset successfully")
}