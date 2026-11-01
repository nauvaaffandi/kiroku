package auth

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/config"
	"kiroku/backend/internal/db"
)

// generateRandomPIN returns a random 6-digit PIN string (e.g. "048392").
func generateRandomPIN() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// writeGenesisPIN writes the plaintext PIN to {DATA_DIR}/persistent/genesis-pin.txt.
// It returns an error on failure so the caller can abort user creation instead of
// silently leaving a user with no recoverable PIN (lockout).
func writeGenesisPin(pin string) error {
	dataDir := config.DataDir()
	if dataDir == "" {
		log.Warn().Msg("auth: DATA_DIR not set, skipping genesis-pin.txt write")
		return nil
	}
	dir := config.GenesisDir()
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create genesis dir: %w", err)
	}
	pinPath := config.GenesisPinPath()
	if err := os.WriteFile(pinPath, []byte(pin), 0o600); err != nil {
		return fmt.Errorf("write genesis-pin.txt: %w", err)
	}
	return nil
}

// WarnIfDefaultGenesisPinInUse logs a warning if the genesis PIN file still matches the
// stored user PIN hash. The PIN value itself is never logged.
func WarnIfDefaultGenesisPinInUse(database *db.DB) {
	if config.DataDir() == "" {
		return
	}
	pinBytes, err := os.ReadFile(config.GenesisPinPath())
	if err != nil {
		return
	}
	genesisPin := strings.TrimSpace(string(pinBytes))
	if genesisPin == "" {
		return
	}
	var hash string
	if err := database.QueryRow("SELECT pin_hash FROM users LIMIT 1").Scan(&hash); err != nil {
		return
	}
	if ok, verr := VerifyPin(genesisPin, hash); verr == nil && ok {
		log.Warn().Msg("auth: default genesis PIN still in use; please change it immediately")
	}
}

// getOrCreateUser returns an existing user ID or creates one with a random
// 6-digit PIN. On first-time creation the plaintext PIN is written to
// {DATA_DIR}/persistent/genesis-pin.txt before the user is inserted.
// If the PIN file cannot be written, the user is NOT created (avoids lockout).
func (h *Handler) getOrCreateUser() (string, error) {
	var userID string
	err := h.db.QueryRow("SELECT id FROM users LIMIT 1").Scan(&userID)
	if err == nil {
		return userID, nil
	}
	if err != sql.ErrNoRows {
		return "", err
	}

	pin, err := generateRandomPIN()
	if err != nil {
		return "", err
	}

	hash, err := HashPin(pin,
		config.Argon2Time(),
		config.Argon2Memory(),
		config.Argon2Threads(),
	)
	if err != nil {
		return "", err
	}

	// Persist the recoverable PIN before creating the user. If the genesis file
	// cannot be written (disk full, permissions), abort before inserting the user
	// so no account is left with an unrecoverable PIN (lockout).
	if err := writeGenesisPin(pin); err != nil {
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