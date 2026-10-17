package auth

import (
    "database/sql"
    "net/http"
    "time"

    "github.com/google/uuid"
    "github.com/rs/zerolog/log"
)

// getOrCreateUser returns an existing user ID or creates one with a default password hash
func (h *Handler) getOrCreateUser() (string, error) {
    var userID string
    err := h.db.QueryRow("SELECT id FROM users LIMIT 1").Scan(&userID)
    if err == nil {
        return userID, nil
    }
    if err != sql.ErrNoRows {
        return "", err
    }

    hash, err := HashPassword("123456", h.cfg.Argon2Time, h.cfg.Argon2Memory, h.cfg.Argon2Threads)
    if err != nil {
        return "", err
    }

    userID = uuid.New().String()
    _, err = h.db.Exec(
        "INSERT INTO users (id, email, password_hash) VALUES (?, ?, ?)",
        userID, "user@kiroku.local", hash,
    )
    if err != nil {
        return "", err
    }

    return userID, nil
}

// setAuthCookies sets the access and refresh cookies on the response
func setAuthCookies(w http.ResponseWriter, accessToken, refreshToken string, accessTTL, refreshTTL time.Duration) {
    secure := false // dev mode
    http.SetCookie(w, &http.Cookie{
        Name:     "access_token",
        Value:    accessToken,
        Path:     "/",
        HttpOnly: true,
        Secure:   secure,
        SameSite: http.SameSiteLaxMode,
        MaxAge:   int(accessTTL.Seconds()),
    })
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    refreshToken,
        Path:     "/",
        HttpOnly: true,
        Secure:   secure,
        SameSite: http.SameSiteLaxMode,
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
        MaxAge:   -1,
    })
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    "",
        Path:     "/",
        HttpOnly: true,
        MaxAge:   -1,
    })
}

// ensure log is used to avoid unused import
var _ = log.Error