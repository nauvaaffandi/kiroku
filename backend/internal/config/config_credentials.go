package config

import (
	"crypto/rand"
	"encoding/hex"
	"os"
)

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
