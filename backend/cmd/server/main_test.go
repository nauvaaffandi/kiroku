package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStatusRecorder_WriteHeader_CapturesStatusCode verifies that
// statusRecorder.WriteHeader records the status code from wrapped handlers.
func TestStatusRecorder_WriteHeader_CapturesStatusCode(t *testing.T) {
	recorder := newStatusRecorder(httptest.NewRecorder())
	recorder.WriteHeader(http.StatusInternalServerError)
	assert.Equal(t, http.StatusInternalServerError, recorder.status)

	recorder.WriteHeader(http.StatusForbidden)
	assert.Equal(t, http.StatusForbidden, recorder.status)
}

// setupErrorLogWriter replaces the package-level zerolog logger with one that
// writes to an in-memory buffer for Error-level checks.
func setupErrorLogWriter(t *testing.T) *bytes.Buffer {
	var buf bytes.Buffer
	newLogger := zerolog.New(&buf).Level(zerolog.ErrorLevel)
	oldLogger := log.Logger
	log.Logger = newLogger
	t.Cleanup(func() { log.Logger = oldLogger })
	return &buf
}

// TestLoggingMiddleware_Passes500StatusToLogger verifies that
// loggingMiddleware calls logger.LogHTTPServerError (Error level) when
// the wrapped handler returns 500.
func TestLoggingMiddleware_Passes500StatusToLogger(t *testing.T) {
	buf := setupErrorLogWriter(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	mw := loggingMiddleware(handler)
	req := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
	rec := httptest.NewRecorder()
	mw.ServeHTTP(rec, req)

	output := buf.String()
	require.NotEmpty(t, output, "expected Error-level log output for 500 responses")
	require.Contains(t, output, `"status":500`)
	require.Contains(t, output, `"method":"GET"`)
	require.Contains(t, output, `"/api/transactions"`)
	require.Contains(t, output, `"level":"error"`)
}

// TestLoggingMiddleware_DoesNotLogErrorFor200 verifies that
// loggingMiddleware does NOT log Error level when the wrapped handler returns 200.
func TestLoggingMiddleware_DoesNotLogErrorFor200(t *testing.T) {
	buf := setupErrorLogWriter(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := loggingMiddleware(handler)
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	mw.ServeHTTP(rec, req)

	output := buf.String()
	assert.Empty(t, output, "expected no log output for 200 responses")
}

// TestLoggingMiddleware_DoesNotLogErrorFor404 verifies that
// loggingMiddleware does NOT log Error level when the wrapped handler returns 404.
func TestLoggingMiddleware_DoesNotLogErrorFor404(t *testing.T) {
	buf := setupErrorLogWriter(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	mw := loggingMiddleware(handler)
	req := httptest.NewRequest(http.MethodGet, "/unknown-route", nil)
	rec := httptest.NewRecorder()
	mw.ServeHTTP(rec, req)

	output := buf.String()
	assert.Empty(t, output, "expected no log output for 404 responses")
}
