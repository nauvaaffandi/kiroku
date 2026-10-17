package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_DefaultsAndFallbacks(t *testing.T) {
	os.Setenv("TELEGRAM_BOT_TOKEN", "test-bot-token-123")
	os.Setenv("TELEGRAM_BOT_USERNAME", "test_bot")
	os.Setenv("UPSTASH_REDIS_REST_URL", "https://test.upstash.io")
	os.Setenv("UPSTASH_REDIS_REST_TOKEN", "test-token")
	defer func() {
		os.Unsetenv("TELEGRAM_BOT_TOKEN")
		os.Unsetenv("TELEGRAM_BOT_USERNAME")
		os.Unsetenv("UPSTASH_REDIS_REST_URL")
		os.Unsetenv("UPSTASH_REDIS_REST_TOKEN")
	}()

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, ":8080", cfg.Addr)
	assert.Equal(t, "test-bot-token-123", cfg.TelegramBotToken)
	assert.Equal(t, "test_bot", cfg.TelegramBotUsername)
	assert.NotEmpty(t, cfg.TelegramWebhookSecret)
	assert.Equal(t, "https://test.upstash.io", cfg.UpstashRedisRestURL)
	assert.Equal(t, "test-token", cfg.UpstashRedisRestToken)
	assert.True(t, cfg.TunnelEnabled)
}

func TestConfig_BackendUpstashTokenFallback(t *testing.T) {
	os.Setenv("BACKEND_UPSTASH_REDIS_REST_TOKEN", "backend-token-xyz")
	defer os.Unsetenv("BACKEND_UPSTASH_REDIS_REST_TOKEN")

	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "backend-token-xyz", cfg.UpstashRedisRestToken)
}
