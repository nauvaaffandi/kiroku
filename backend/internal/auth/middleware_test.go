package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kiroku/backend/internal/config"
)

func newConfig() *config.Config { return &config.Config{} }

func TestAuthMiddleware_MissingToken(t *testing.T) {
	middleware := AuthMiddleware(newConfig())
	called := false
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	middleware := AuthMiddleware(newConfig())
	called := false
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "garbage"})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	_, _ = setupTestHandler(t)
	token, err := GenerateAccessToken("user-1", config.JWTSecret(), time.Hour)
	require.NoError(t, err)

	middleware := AuthMiddleware(newConfig())
	called := false
	var gotUserID string
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		gotUserID = GetUserID(r)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.True(t, called)
	assert.Equal(t, "user-1", gotUserID)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestOptionalAuthMiddleware_ValidTokenInjectsContext(t *testing.T) {
	_, _ = setupTestHandler(t)
	token, err := GenerateAccessToken("user-2", config.JWTSecret(), time.Hour)
	require.NoError(t, err)

	called := false
	var gotUserID string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		gotUserID = GetUserID(r)
		w.WriteHeader(http.StatusOK)
	})
	middleware := OptionalAuthMiddleware(newConfig())
	wrapped := middleware(next)

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)
	assert.True(t, called)
	assert.Equal(t, "user-2", gotUserID)
}

func TestOptionalAuthMiddleware_NoTokenPassesThrough(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Empty(t, GetUserID(r))
		w.WriteHeader(http.StatusOK)
	})
	wrapped := OptionalAuthMiddleware(newConfig())(next)
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	assert.True(t, called)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestExtractToken_BearerPrefix(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer abc123")
	assert.Equal(t, "abc123", extractToken(req))
}

func TestExtractToken_CookiePreferred(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "cookieval"})
	assert.Equal(t, "cookieval", extractToken(req))
}

func TestExtractToken_None(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	assert.Equal(t, "", extractToken(req))
}

func TestGetUserID_MissingFromContext(t *testing.T) {
	assert.Equal(t, "", GetUserID(httptest.NewRequest("GET", "/", nil)))
}

func TestWithUserID_ReturnsRequestWithContext(t *testing.T) {
	original := httptest.NewRequest("GET", "/", strings.NewReader(""))
	out := WithUserID(original, "uid-9")
	assert.NotSame(t, original, out)
	assert.Equal(t, "uid-9", GetUserID(out))
}
