package logger

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRetentionDays(t *testing.T) {
	assert.Equal(t, 90, parseRetentionDays("90d"))
	assert.Equal(t, 30, parseRetentionDays("30D"))
	assert.Equal(t, 0, parseRetentionDays("abc"))
	assert.Equal(t, 0, parseRetentionDays("-5d"))
}

func TestCleanupOldLogs(t *testing.T) {
	dir := logsDir()
	require.NoError(t, os.MkdirAll(dir, 0755))

	old := time.Now().Add(-100 * 24 * time.Hour).Format("2006-01-02") + ".log"
	recent := time.Now().Add(-10 * 24 * time.Hour).Format("2006-01-02") + ".log"
	weird := "not-a-date.log"

	require.NoError(t, os.WriteFile(filepath.Join(dir, old), []byte("{}"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, recent), []byte("{}"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, weird), []byte("{}"), 0644))

	cleanupOldLogs(90)

	assert.False(t, fileExists(filepath.Join(dir, old)), "old log should be removed")
	assert.True(t, fileExists(filepath.Join(dir, recent)), "recent log should be kept")
	assert.True(t, fileExists(filepath.Join(dir, weird)), "non-matching name should be kept")

	os.Remove(filepath.Join(dir, recent))
	os.Remove(filepath.Join(dir, weird))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// setupErrorLogWriter replaces the package-level zerolog logger with one that
// writes to an in-memory buffer. Returns the buffer to inspect.
func setupErrorLogWriter(t *testing.T) *bytes.Buffer {
	var buf bytes.Buffer
	newLogger := zerolog.New(&buf).Level(zerolog.ErrorLevel)
	oldLogger := log.Logger
	log.Logger = newLogger
	t.Cleanup(func() { log.Logger = oldLogger })
	return &buf
}

// TestLogHTTPServerError_LogsAtErrorLevel verifies that LogHTTPServerError
// logs at Error level with the expected fields when status >= 500.
func TestLogHTTPServerError_LogsAtErrorLevel(t *testing.T) {
	buf := setupErrorLogWriter(t)

	LogHTTPServerError(http.StatusInternalServerError, "GET", "/api/transactions", "req-123")

	output := buf.String()
	require.NotEmpty(t, output, "expected log output to not be empty")
	require.Contains(t, output, `"status":500`)
	require.Contains(t, output, `"method":"GET"`)
	require.Contains(t, output, `"path":"/api/transactions"`)
	require.Contains(t, output, `"request_id":"req-123"`)
}

// TestLogHTTPServerError_DoesNotLogWhenStatusLessThan500 verifies that
// LogHTTPServerError does not log when status < 500.
func TestLogHTTPServerError_DoesNotLogWhenStatusLessThan500(t *testing.T) {
	buf := setupErrorLogWriter(t)

	LogHTTPServerError(http.StatusOK, "GET", "/api/transactions", "req-123")

	output := buf.String()
	assert.Empty(t, output, "expected no log output when status < 500")
}
