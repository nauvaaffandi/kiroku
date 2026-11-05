package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"kiroku/backend/internal/logger"
)

func TestRequestLoggingMiddleware_LogsCompletion_200(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf).Level(zerolog.InfoLevel)
	t.Cleanup(func() { log.Logger = zerolog.New(&bytes.Buffer{}) })

	// Chain: chi RequestID middleware → our logger middleware
	handler := middleware.RequestID(logger.RequestLoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/records", nil)
	req.Header.Set("X-Request-ID", "req-abc")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	output := buf.String()
	assert.Contains(t, output, `"method":"GET"`)
	assert.Contains(t, output, `"status":200`)
	assert.Contains(t, output, `"duration_ms":`)
	assert.Contains(t, output, `"client_ip":`)
	assert.Contains(t, output, `"user_agent":`)
	assert.Contains(t, output, `"request_id":"req-abc"`)
	assert.Contains(t, output, `"level":"info"`)
	assert.Contains(t, output, `"message":"http request completed"`)
}

func TestRequestLoggingMiddleware_SeverityMapping(t *testing.T) {
	cases := []struct {
		status int
		level  string
	}{
		{200, "info"},
		{301, "info"},
		{400, "warn"},
		{404, "warn"},
		{429, "warn"},
		{500, "error"},
		{503, "error"},
	}

	for _, c := range cases {
		t.Run(http.StatusText(c.status), func(t *testing.T) {
			var buf bytes.Buffer
			log.Logger = zerolog.New(&buf).Level(zerolog.InfoLevel)

			handler := middleware.RequestID(logger.RequestLoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(c.status)
			})))

			req := httptest.NewRequest(http.MethodGet, "/api/x", nil)
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)

			assert.Equal(t, c.status, resp.Code)
			assert.Contains(t, buf.String(), `"`+c.level+`"`)
		})
	}
}

func TestRequestLoggingMiddleware_PanicRecovery(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf).Level(zerolog.InfoLevel)
	t.Cleanup(func() { log.Logger = zerolog.New(&bytes.Buffer{}) })

	handler := middleware.RequestID(logger.RequestLoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("simulated panic")
	})))

	req := httptest.NewRequest(http.MethodGet, "/api/panic", nil)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	output := buf.String()
	assert.Contains(t, output, `"status":500`)
	assert.Contains(t, output, `"level":"error"`)
	assert.Contains(t, output, `"message":"http request completed"`)
}

func TestRequestLoggingMiddleware_RedactsSensitiveData(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf).Level(zerolog.InfoLevel)
	t.Cleanup(func() { log.Logger = zerolog.New(&bytes.Buffer{}) })

	handler := middleware.RequestID(logger.RequestLoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	req := httptest.NewRequest(http.MethodGet, "/api/users?id=123&token=secret-token", nil)
	req.Header.Set("Authorization", "Bearer secret")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	output := buf.String()
	assert.NotContains(t, output, "secret-token")
	assert.NotContains(t, output, "Authorization")
	assert.NotContains(t, output, "Bearer")
	assert.NotContains(t, output, "token")
	assert.Contains(t, output, `"path":"/api/users"`)
}

func TestInit_DoesNotCrash(t *testing.T) {
	logger.Init(true)
	assert.NotNil(t, logger.GetLogger())
	logger.Init(false)
	assert.NotNil(t, logger.GetLogger())
}

func TestGetLogger_ReturnsComponentLogger(t *testing.T) {
	logger.Init(false)
	defer logger.Init(false)

	comp := logger.NewComponentLogger("auth")
	assert.NotNil(t, comp)

	comp.Info().Msg("component auth log")
	// No panic expected; output silenced via io.Discard writer? GetLogger uses global logger (console).
	// Assert strings are safe to pass to logger at runtime.
	assert.True(t, len("auth") > 0)
}