package telegram

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kiroku/backend/internal/config"
	"kiroku/backend/internal/db"
	"kiroku/backend/internal/models"
)

func setupTestDB(t *testing.T) *db.DB {
	t.Helper()
	dataDir := t.TempDir()
	os.Setenv("DATA_DIR", dataDir)
	t.Cleanup(func() { os.Unsetenv("DATA_DIR") })
	database, err := db.Open(dataDir)
	require.NoError(t, err)
	t.Cleanup(func() {
		database.Close()
	})
	return database
}

func TestHandleWebhook_Unauthorized(t *testing.T) {
	os.Setenv("TELEGRAM_BOT_TOKEN", "test-token")
	os.Setenv("TELEGRAM_WEBHOOK_SECRET", "super-secret-token")
	t.Cleanup(func() {
		os.Unsetenv("TELEGRAM_BOT_TOKEN")
		os.Unsetenv("TELEGRAM_WEBHOOK_SECRET")
	})
	svc := NewTelegramService(nil, &config.Config{})

	body := []byte(`{"update_id": 1}`)

	// Request without header
	req := httptest.NewRequest("POST", "/webhook", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	svc.HandleWebhook(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	// Request with wrong header
	req = httptest.NewRequest("POST", "/webhook", bytes.NewReader(body))
	req.Header.Set("X-Telegram-Bot-Api-Secret-Token", "wrong-token")
	rec = httptest.NewRecorder()
	svc.HandleWebhook(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHandleWebhook_ValidSecret(t *testing.T) {
	os.Setenv("TELEGRAM_BOT_TOKEN", "test-token")
	os.Setenv("TELEGRAM_WEBHOOK_SECRET", "super-secret-token")
	t.Cleanup(func() {
		os.Unsetenv("TELEGRAM_BOT_TOKEN")
		os.Unsetenv("TELEGRAM_WEBHOOK_SECRET")
	})
	svc := NewTelegramService(nil, &config.Config{})

	body := []byte(`{"update_id": 1}`)
	req := httptest.NewRequest("POST", "/webhook", bytes.NewReader(body))
	req.Header.Set("X-Telegram-Bot-Api-Secret-Token", "super-secret-token")
	rec := httptest.NewRecorder()

	svc.HandleWebhook(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandleCommand_StartRegister(t *testing.T) {
	database := setupTestDB(t)

	// Insert test user and register token
	userID := "user-123"
	_, err := database.Exec(`INSERT INTO users (id, email, pin_hash) VALUES (?, 'test@example.com', 'hash')`, userID)
	require.NoError(t, err)

	token := "valid-token-123"
	expiresAt := time.Now().Add(10 * time.Minute).Format("2006-01-02 15:04:05")
	_, err = database.Exec(`
		INSERT INTO pin_change_tokens (id, user_id, token, type, expires_at)
		VALUES ('tok-1', ?, ?, 'register', ?)
	`, userID, token, expiresAt)
	require.NoError(t, err)

	// Mock server for Telegram sendMessage
	sendMessageCalled := false
	mockTg := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sendMessageCalled = true
		assert.Contains(t, r.URL.Path, "/sendMessage")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok": true, "result": {}}`))
	}))
	defer mockTg.Close()

	os.Setenv("TELEGRAM_BOT_TOKEN", "mock-token")
	os.Setenv("TELEGRAM_WEBHOOK_SECRET", "sec-123")
	t.Cleanup(func() {
		os.Unsetenv("TELEGRAM_BOT_TOKEN")
		os.Unsetenv("TELEGRAM_WEBHOOK_SECRET")
	})
	svc := NewTelegramService(database, &config.Config{})
	// Use a transport that redirects requests to the mock server
	svc.httpClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			req.URL.Scheme = "http"
			req.URL.Host = mockTg.Listener.Addr().String()
			return http.DefaultTransport.RoundTrip(req)
		}),
	}

	msg := &models.Message{
		MessageID: 1,
		Chat:      models.Chat{ID: 98765},
		From:      &models.User{ID: 98765},
		Text:      "/start register_" + token,
	}

	svc.handleCommand(msg)

	assert.True(t, sendMessageCalled)

	// Verify binding was created
	var boundUserID string
	err = database.QueryRow(`SELECT user_id FROM telegram_bindings WHERE telegram_chat_id = '98765'`).Scan(&boundUserID)
	require.NoError(t, err)
	assert.Equal(t, userID, boundUserID)

	// Verify token was deleted
	var count int
	database.QueryRow(`SELECT COUNT(*) FROM pin_change_tokens WHERE token = ?`, token).Scan(&count)
	assert.Equal(t, 0, count)
}

func TestSetWebhook_Payload(t *testing.T) {
	var capturedPayload map[string]interface{}
	mockTg := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Contains(t, r.URL.Path, "/botmy-bot-token/setWebhook")
		json.NewDecoder(r.Body).Decode(&capturedPayload)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok": true, "result": true}`))
	}))
	defer mockTg.Close()

	os.Setenv("TELEGRAM_BOT_TOKEN", "my-bot-token")
	os.Setenv("TELEGRAM_WEBHOOK_SECRET", "secret-xyz")
	t.Cleanup(func() {
		os.Unsetenv("TELEGRAM_BOT_TOKEN")
		os.Unsetenv("TELEGRAM_WEBHOOK_SECRET")
	})
	svc := NewTelegramService(nil, &config.Config{})
	// Point httpClient transport to the mock server
	svc.httpClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			req.URL.Scheme = "http"
			req.URL.Host = mockTg.Listener.Addr().String()
			return http.DefaultTransport.RoundTrip(req)
		}),
	}

	err := svc.SetWebhook("https://example.trycloudflare.com/api/v1/telegram/webhook")
	require.NoError(t, err)

	assert.Equal(t, "https://example.trycloudflare.com/api/v1/telegram/webhook", capturedPayload["url"])
	assert.Equal(t, "secret-xyz", capturedPayload["secret_token"])
	assert.Equal(t, true, capturedPayload["drop_pending_updates"])
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
