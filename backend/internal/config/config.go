package config

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Addr string
}

func loadDotEnv() {
	paths := []string{".env", "../.env"}
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(b), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
				if _, exists := os.LookupEnv(key); !exists {
					os.Setenv(key, val)
				}
			}
		}
	}
}

func Load() (*Config, error) {
	loadDotEnv()
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}
	return &Config{
		Addr: ":" + addr,
	}, nil
}

// FXBaseCurrency returns the base currency from FX_BASE_CURRENCY env var, default "JPY".
func FXBaseCurrency() string {
	if v := os.Getenv("FX_BASE_CURRENCY"); v != "" {
		return v
	}
	return "JPY"
}

// BackupSchedule returns the backup cron schedule from BACKUP_SCHEDULE env var, default "0 3 * * *".
func BackupSchedule() string {
	if v := os.Getenv("BACKUP_SCHEDULE"); v != "" {
		return v
	}
	return "0 3 * * *"
}

// BackupRetentionDaily returns daily retention count from BACKUP_RETENTION_DAILY env var, default 7.
func BackupRetentionDaily() int {
	return retentionEnv("BACKUP_RETENTION_DAILY", 7)
}

// BackupRetentionWeekly returns weekly retention count from BACKUP_RETENTION_WEEKLY env var, default 4.
func BackupRetentionWeekly() int {
	return retentionEnv("BACKUP_RETENTION_WEEKLY", 4)
}

// BackupRetentionMonthly returns monthly retention count from BACKUP_RETENTION_MONTHLY env var, default 6.
func BackupRetentionMonthly() int {
	return retentionEnv("BACKUP_RETENTION_MONTHLY", 6)
}

func retentionEnv(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

// FullDBPath returns the full path to the database file.
func FullDBPath(dataDir string) string {
	return filepath.Join(dataDir, "kiroku.db")
}

// GetEnv returns the value of the named environment variable, or def if unset.
func GetEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// VAPIDPublicKey returns the VAPID public key from VAPID_PUBLIC_KEY env var.
func VAPIDPublicKey() string {
	return os.Getenv("VAPID_PUBLIC_KEY")
}

// VAPIDPrivateKey returns the VAPID private key from VAPID_PRIVATE_KEY env var.
func VAPIDPrivateKey() string {
	return os.Getenv("VAPID_PRIVATE_KEY")
}

// VAPIDSubject returns the VAPID subject from VAPID_SUBJECT env var.
func VAPIDSubject() string {
	if v := os.Getenv("VAPID_SUBJECT"); v != "" {
		return v
	}
	return "mailto:admin@kiroku.local"
}

// OCRSpaceAPIKey returns the OCR.space API key from OCR_SPACE_API_KEY env var.
func OCRSpaceAPIKey() string {
	return os.Getenv("OCR_SPACE_API_KEY")
}

// NineRouterBaseURL returns the 9Router base URL from NINE_ROUTER_URL env var.
func NineRouterBaseURL() string {
	if v := os.Getenv("NINE_ROUTER_URL"); v != "" {
		return v
	}
	return "http://localhost:20128/v1"
}

// NineRouterAPIKey returns the 9Router API key from NINE_ROUTER_API_KEY env var.
func NineRouterAPIKey() string {
	return os.Getenv("NINE_ROUTER_API_KEY")
}

// NineRouterModelComboName returns the 9Router model combo name from NINE_ROUTER_MODEL_COMBO_NAME env var.
func NineRouterModelComboName() string {
	if v := os.Getenv("NINE_ROUTER_MODEL_COMBO_NAME"); v != "" {
		return v
	}
	return "main"
}

// MockOCREnabled returns true if KIROKU_MOCK_OCR_ENABLED is "true".
func MockOCREnabled() bool {
	return os.Getenv("KIROKU_MOCK_OCR_ENABLED") == "true"
}

// TelegramBotToken returns the Telegram bot token from TELEGRAM_BOT_TOKEN env var.
func TelegramBotToken() string {
	return os.Getenv("TELEGRAM_BOT_TOKEN")
}

// TelegramBotUsername returns the Telegram bot username from TELEGRAM_BOT_USERNAME env var.
func TelegramBotUsername() string {
	if v := os.Getenv("TELEGRAM_BOT_USERNAME"); v != "" {
		return v
	}
	return "nauva_kiroku_bot"
}

// TelegramWebhookSecret returns the Telegram webhook secret from TELEGRAM_WEBHOOK_SECRET env var.
func TelegramWebhookSecret() string {
	if v := os.Getenv("TELEGRAM_WEBHOOK_SECRET"); v != "" {
		return v
	}
	// Auto-generate random 32-hex-char secret if not set
	b := make([]byte, 16)
	if _, err := rand.Read(b); err == nil {
		return hex.EncodeToString(b)
	}
	return ""
}

// UpstashRedisRestURL returns the Upstash Redis REST URL.
func UpstashRedisRestURL() string {
	return GetEnv("UPSTASH_REDIS_REST_URL", "")
}

func UpstashRedisRestToken() string {
	return os.Getenv("BACKEND_UPSTASH_REDIS_REST_TOKEN")
}

// LogRetention returns the log retention from LOG_RETENTION env var, default "90d".
func LogRetention() string {
	return GetEnv("LOG_RETENTION", "90d")
}

// TunnelEnabled returns true unless CLOUDFLARE_TUNNEL_ENABLED is "false".
func TunnelEnabled() bool {
	return os.Getenv("CLOUDFLARE_TUNNEL_ENABLED") != "false"
}

// BackupPassphrase returns the backup passphrase from BACKUP_PASSPHRASE env var.
func BackupPassphrase() string {
	return os.Getenv("BACKUP_PASSPHRASE")
}

// BackupRepoPath returns the backup repo path from BACKUP_REPO_PATH env var.
func BackupRepoPath() string {
	if v := os.Getenv("BACKUP_REPO_PATH"); v != "" {
		return v
	}
	return "backup/kiroku-backup-storage"
}

// BackupGitToken returns the backup git token.
func BackupGitToken() string {
	return GetEnv("KIROKU_BACKUP_GIT_TOKEN", GetEnv("BACKUP_GIT_TOKEN", ""))
}

// BackupProviders returns the backup providers from BACKUP_PROVIDERS env var.
func BackupProviders() string {
	return os.Getenv("BACKUP_PROVIDERS")
}

// DataDir returns the data directory from DATA_DIR env var.
func DataDir() string {
	return os.Getenv("DATA_DIR")
}

// SeedDir returns the seed data directory path under DATA_DIR/.seed.
func SeedDir() string {
	return filepath.Join(DataDir(), ".seed")
}

// SeedPINPath returns the path to seed-pin.txt.
func SeedPINPath() string {
	return filepath.Join(SeedDir(), "seed-pin.txt")
}

// SeedURLPath returns the path to backend-url.txt.
func SeedURLPath() string {
	return filepath.Join(SeedDir(), "backend-url.txt")
}

// JWTSecret returns the JWT secret from JWT_SECRET env var.
func JWTSecret() string {
	return os.Getenv("JWT_SECRET")
}

// JWTAccessTTL returns the JWT access TTL from JWT_ACCESS_TTL env var, default 15m.
func JWTAccessTTL() time.Duration {
	d, err := time.ParseDuration(os.Getenv("JWT_ACCESS_TTL"))
	if err != nil {
		return 15 * time.Minute
	}
	return d
}

// JWTRefreshTTL returns the JWT refresh TTL from JWT_REFRESH_TTL env var, default 168h.
func JWTRefreshTTL() time.Duration {
	d, err := time.ParseDuration(os.Getenv("JWT_REFRESH_TTL"))
	if err != nil {
		return 168 * time.Hour
	}
	return d
}

// Argon2Time returns the Argon2 time cost from ARGON2_TIME env var, default 3.
func Argon2Time() uint32 {
	v := retentionEnv("ARGON2_TIME", 3)
	return uint32(v)
}

// Argon2Memory returns the Argon2 memory cost from ARGON2_MEMORY env var, default 65536.
func Argon2Memory() uint32 {
	v := retentionEnv("ARGON2_MEMORY", 65536)
	return uint32(v)
}

// Argon2Threads returns the Argon2 parallelism from ARGON2_THREADS env var, default 4.
func Argon2Threads() uint8 {
	v := retentionEnv("ARGON2_THREADS", 4)
	return uint8(v)
}

// KIROKUDev returns true if KIROKU_DEV is "true".
func KIROKUDev() bool {
	return os.Getenv("KIROKU_DEV") == "true"
}

// AllowedCORSOrigins returns the allowed CORS origins from CORS_ALLOWED_ORIGINS env var.
func AllowedCORSOrigins() []string {
    if v := os.Getenv("CORS_ALLOWED_ORIGINS"); v != "" {
        var origins []string
        parts := strings.Split(v, ",")
        for _, p := range parts {
            if trimmed := strings.TrimSpace(p); trimmed != "" {
                origins = append(origins, trimmed)
            }
        }
        return origins
    }
    return []string{}
}
