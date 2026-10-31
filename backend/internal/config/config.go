package config

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Addr          string        `koanf:"addr"`
	DataDir       string        `koanf:"data_dir"`
	DBPath        string        `koanf:"db_path"`
	JWTSecret     string        `koanf:"jwt_secret"`
	JWTAccessTTL  time.Duration `koanf:"jwt_access_ttl"`
	JWTRefreshTTL time.Duration `koanf:"jwt_refresh_ttl"`
	Argon2Time    uint32        `koanf:"argon2_time"`
	Argon2Memory  uint32        `koanf:"argon2_memory"`
	Argon2Threads uint8         `koanf:"argon2_threads"`
	// WARNING: Changing Argon2 parameters will invalidate existing password hashes.
	// All existing users will need to reset their passwords.
	FXBaseCurrency           string `koanf:"fx_base_currency"`
	BackupSchedule           string `koanf:"backup_schedule"`
	BackupRetentionDaily     int    `koanf:"backup_retention_daily"`
	BackupRetentionWeekly    int    `koanf:"backup_retention_weekly"`
	BackupRetentionMonthly   int    `koanf:"backup_retention_monthly"`
	BackupPassphrase         string `koanf:"backup_passphrase"`
	BackupRepoPath           string `koanf:"backup_repo_path"`
	BackupGitToken           string `koanf:"backup_git_token"`
	BackupProviders          string `koanf:"backup_providers"`
	OCRSpaceAPIKey           string `koanf:"ocr_space_api_key"`
	NineRouterBaseURL        string `koanf:"ninerouter_base_url"`
	NineRouterAPIKey         string `koanf:"ninerouter_api_key"`
	NineRouterModelComboName string `koanf:"ninerouter_model_combo_name"`
	MockOCREnabled           bool   `koanf:"mock_ocr_enabled"`
	TelegramBotToken         string `koanf:"telegram_bot_token"`
	TelegramBotUsername      string `koanf:"telegram_bot_username"`
	TelegramWebhookSecret    string `koanf:"telegram_webhook_secret"`
	UpstashRedisRestURL      string `koanf:"upstash_redis_rest_url"`
	UpstashRedisRestToken    string `koanf:"upstash_redis_rest_token"`
	TunnelEnabled            bool   `koanf:"tunnel_enabled"`
	LogRetention             string `koanf:"log_retention"`
	VAPIDPublicKey           string `koanf:"vapid_public_key"`
	VAPIDPrivateKey          string `koanf:"vapid_private_key"`
	VAPIDSubject             string `koanf:"vapid_subject"`
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
	k := koanf.New(".")

	k.Load(file.Provider(filepath.Join(".", "config.yaml")), yaml.Parser())
	k.Load(env.Provider("KIROKU_", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(s, "KIROKU_")), "_", ".")
	}), nil)

	cfg := &Config{}
	if err := k.Unmarshal("", cfg); err != nil {
		return nil, err
	}

	if cfg.Addr == "" {
		cfg.Addr = ":8080"
	}
	if cfg.DataDir == "" {
		cfg.DataDir = "./data"
	}
	if cfg.DBPath == "" {
		cfg.DBPath = "kiroku.db"
	}
	if cfg.JWTSecret == "" {
		cfg.JWTSecret = "changeme"
	}
	if cfg.JWTAccessTTL == 0 {
		cfg.JWTAccessTTL = 15 * time.Minute
	}
	if cfg.JWTRefreshTTL == 0 {
		cfg.JWTRefreshTTL = 168 * time.Hour
	}
	if cfg.Argon2Time == 0 {
		cfg.Argon2Time = 3
	}
	if cfg.Argon2Memory == 0 {
		cfg.Argon2Memory = 65536
	}
	if cfg.Argon2Threads == 0 {
		cfg.Argon2Threads = 4
	}
	if cfg.FXBaseCurrency == "" {
		cfg.FXBaseCurrency = "JPY"
	}
	if cfg.BackupSchedule == "" {
		cfg.BackupSchedule = "0 3 * * *"
	}
	if cfg.BackupRetentionDaily == 0 {
		cfg.BackupRetentionDaily = 7
	}
	if cfg.BackupRetentionWeekly == 0 {
		cfg.BackupRetentionWeekly = 4
	}
	if cfg.BackupRetentionMonthly == 0 {
		cfg.BackupRetentionMonthly = 6
	}
	if cfg.BackupRepoPath == "" {
		cfg.BackupRepoPath = "backup/kiroku-backup-storage"
	}
	if cfg.BackupGitToken == "" {
		cfg.BackupGitToken = os.Getenv("KIROKU_BACKUP_GIT_TOKEN")
	}
	if cfg.BackupGitToken == "" {
		cfg.BackupGitToken = os.Getenv("BACKUP_GIT_TOKEN")
	}
	if cfg.NineRouterBaseURL == "" {
		cfg.NineRouterBaseURL = "http://localhost:20128/v1"
	}
	if cfg.NineRouterModelComboName == "" {
		cfg.NineRouterModelComboName = "main"
	}

	// Environment variable fallbacks for Telegram & Upstash
	if cfg.TelegramBotToken == "" {
		cfg.TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	}
	if cfg.TelegramBotUsername == "" {
		cfg.TelegramBotUsername = GetEnv("TELEGRAM_BOT_USERNAME", "nauva_kiroku_bot")
	}
	if cfg.TelegramWebhookSecret == "" {
		cfg.TelegramWebhookSecret = os.Getenv("TELEGRAM_WEBHOOK_SECRET")
	}
	if cfg.TelegramWebhookSecret == "" {
		// Auto-generate random 32-hex-char secret if not set
		bytes := make([]byte, 16)
		if _, err := rand.Read(bytes); err == nil {
			cfg.TelegramWebhookSecret = hex.EncodeToString(bytes)
		}
	}
	if cfg.UpstashRedisRestURL == "" {
		cfg.UpstashRedisRestURL = os.Getenv("UPSTASH_REDIS_REST_URL")
	}
	if cfg.UpstashRedisRestToken == "" {
		cfg.UpstashRedisRestToken = os.Getenv("UPSTASH_REDIS_REST_TOKEN")
	}
	if cfg.UpstashRedisRestToken == "" {
		cfg.UpstashRedisRestToken = os.Getenv("BACKEND_UPSTASH_REDIS_REST_TOKEN")
	}
	if cfg.LogRetention == "" {
		cfg.LogRetention = os.Getenv("KIROKU_LOG_RETENTION")
	}
	if cfg.LogRetention == "" {
		cfg.LogRetention = os.Getenv("LOG_RETENTION")
	}
	if cfg.LogRetention == "" {
		cfg.LogRetention = "90d"
	}

	// VAPID configuration: env var fallbacks; generate fresh key pair if missing
	if cfg.VAPIDPublicKey == "" {
		cfg.VAPIDPublicKey = os.Getenv("VAPID_PUBLIC_KEY")
	}
	if cfg.VAPIDPrivateKey == "" {
		cfg.VAPIDPrivateKey = os.Getenv("VAPID_PRIVATE_KEY")
	}
	if cfg.VAPIDSubject == "" {
		cfg.VAPIDSubject = os.Getenv("VAPID_SUBJECT")
		if cfg.VAPIDSubject == "" {
			cfg.VAPIDSubject = "mailto:admin@kiroku.local"
		}
	}
	if cfg.VAPIDPublicKey == "" || cfg.VAPIDPrivateKey == "" {
		pubKey, privKey, err := generateVAPIDKeys()
		if err == nil {
			if cfg.VAPIDPublicKey == "" {
				cfg.VAPIDPublicKey = pubKey
			}
			if cfg.VAPIDPrivateKey == "" {
				cfg.VAPIDPrivateKey = privKey
			}
		}
	}

	if os.Getenv("KIROKU_TUNNEL_ENABLED") == "false" || os.Getenv("CLOUDFLARE_TUNNEL_ENABLED") == "false" {
		cfg.TunnelEnabled = false
	} else {
		cfg.TunnelEnabled = true
	}
	return cfg, nil
}

func (c *Config) FullDBPath() string {
	return filepath.Join(c.DataDir, c.DBPath)
}

func GetEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func generateVAPIDKeys() (string, string, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}
	pubBytes := elliptic.Marshal(priv.Curve, priv.X, priv.Y)
	pubKey := base64.RawURLEncoding.EncodeToString(pubBytes)
	privKey := base64.RawURLEncoding.EncodeToString(priv.D.Bytes())
	return pubKey, privKey, nil
}
