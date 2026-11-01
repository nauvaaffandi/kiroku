package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kiroku/backend/internal/config"
	"kiroku/backend/internal/db"
)

func setupTestHandler(t *testing.T) (*Handler, *config.Config) {
	t.Helper()
	dataDir := t.TempDir()
	os.Setenv("DATA_DIR", dataDir)
	os.Setenv("JWT_SECRET", "test-secret-key-32-bytes-long!!!!!")
	os.Setenv("JWT_ACCESS_TTL", "15m")
	os.Setenv("JWT_REFRESH_TTL", "24h")
	os.Setenv("ARGON2_TIME", "1")
	os.Setenv("ARGON2_MEMORY", "65536")
	os.Setenv("ARGON2_THREADS", "2")
	t.Cleanup(func() {
		os.Unsetenv("DATA_DIR")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("JWT_ACCESS_TTL")
		os.Unsetenv("JWT_REFRESH_TTL")
		os.Unsetenv("ARGON2_TIME")
		os.Unsetenv("ARGON2_MEMORY")
		os.Unsetenv("ARGON2_THREADS")
	})

	cfg := &config.Config{}

	database, err := db.Open(dataDir)
	require.NoError(t, err)
	t.Cleanup(func() { database.Close() })

	hash, err := HashPin("654321",
		config.Argon2Time(),
		config.Argon2Memory(),
		config.Argon2Threads(),
	)
	require.NoError(t, err)
	_, err = database.Exec("INSERT INTO users (id, email, pin_hash) VALUES (?, ?, ?)", "test-user", "test@example.com", hash)
	require.NoError(t, err)
	h := NewHandler(database, cfg, nil)
	return h, cfg
}

func TestAuthFlow_LoginAndMe(t *testing.T) {
	h, _ := setupTestHandler(t)
	router := h.Routes()

	// 1. Unauthenticated GET /me should return 401
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	// 2. Successful login with test PIN "654321"
	loginBody, _ := json.Marshal(map[string]string{"pin": "654321"})
	req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Extract access_token cookie
	cookies := rec.Result().Cookies()
	var accessCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "access_token" {
			accessCookie = c
			break
		}
	}
	require.NotNil(t, accessCookie, "access_token cookie should be set")

	// 3. Authenticated GET /me using cookie should return 200 with user_id
	req = httptest.NewRequest(http.MethodGet, "/me", nil)
	req.AddCookie(accessCookie)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var meResp map[string]string
	err := json.NewDecoder(rec.Body).Decode(&meResp)
	require.NoError(t, err)
	assert.NotEmpty(t, meResp["user_id"])

	// 4. Authenticated GET /me using Authorization: Bearer header should also return 200
	req = httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+accessCookie.Value)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var meBearerResp map[string]string
	err = json.NewDecoder(rec.Body).Decode(&meBearerResp)
	require.NoError(t, err)
	assert.Equal(t, meResp["user_id"], meBearerResp["user_id"])

	// 5. POST /logout should clear cookies
	req = httptest.NewRequest(http.MethodPost, "/logout", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}
