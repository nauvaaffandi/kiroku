package auth

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

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
