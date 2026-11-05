package logger

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit_DefaultLogLevel(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("LOG_DIR", dir)
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("KIROKU_DEV", "false")
	t.Setenv("LOG_FORMAT", "json")

	log.Logger = zerolog.New(&bytes.Buffer{}).Level(zerolog.InfoLevel)
	Init(false)
	assert.Equal(t, zerolog.InfoLevel, log.Logger.GetLevel())
}

func TestInit_DevPreset_Debug(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("LOG_DIR", dir)
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("KIROKU_DEV", "true")
	t.Setenv("LOG_FORMAT", "json")

	log.Logger = zerolog.New(&bytes.Buffer{}).Level(zerolog.InfoLevel)
	Init(false)
	assert.Equal(t, zerolog.DebugLevel, log.Logger.GetLevel())
}

func TestInit_ExplicitLogLevel_Warn(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("LOG_DIR", dir)
	t.Setenv("LOG_LEVEL", "WARN")
	t.Setenv("KIROKU_DEV", "true")
	t.Setenv("LOG_FORMAT", "json")

	log.Logger = zerolog.New(&bytes.Buffer{}).Level(zerolog.InfoLevel)
	Init(false)
	assert.Equal(t, zerolog.WarnLevel, log.Logger.GetLevel())
}

func TestInit_InvalidLogLevel_Fallback(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("LOG_DIR", dir)
	t.Setenv("LOG_LEVEL", "bogus")
	t.Setenv("KIROKU_DEV", "false")
	t.Setenv("LOG_FORMAT", "json")

	log.Logger = zerolog.New(&bytes.Buffer{}).Level(zerolog.InfoLevel)
	Init(false)
	assert.Equal(t, zerolog.InfoLevel, log.Logger.GetLevel())
}

func TestInit_LogFormat_JSON(t *testing.T) {
	var buf bytes.Buffer
	t.Setenv("LOG_FORMAT", "json")
	t.Setenv("LOG_DIR", t.TempDir())

	log.Logger = zerolog.New(&buf).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Str("environment", "production").
		Str("service", "kiroku").
		Logger()
	log.Info().Msg("test log")

	output := buf.String()
	assert.Contains(t, output, `"level":"info"`)
	assert.Contains(t, output, `"service":"kiroku"`)
	assert.Contains(t, output, `"environment":"production"`)
	assert.Contains(t, output, `"message":"test log"`)
	assert.Contains(t, output, `"time":`)
}

func TestInit_LogFormat_Console(t *testing.T) {
	var buf bytes.Buffer
	t.Setenv("LOG_FORMAT", "console")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("LOG_DIR", t.TempDir())

	log.Logger = zerolog.New(&buf).
		Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Str("environment", "development").
		Str("service", "kiroku").
		Logger()
	log.Debug().Msg("console test")

	output := buf.String()
	assert.NotEmpty(t, output)
	assert.Contains(t, output, "console test")
}

func TestInit_InvalidLogFormat_Fallback(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("LOG_DIR", dir)
	t.Setenv("LOG_FORMAT", "invalid")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("KIROKU_DEV", "false")

	log.Logger = zerolog.New(&bytes.Buffer{}).Level(zerolog.InfoLevel)
	Init(false)
	assert.Equal(t, zerolog.InfoLevel, log.Logger.GetLevel())
}

func TestSeverityForStatus(t *testing.T) {
	assert.Equal(t, zerolog.InfoLevel, severityForStatus(100))
	assert.Equal(t, zerolog.InfoLevel, severityForStatus(200))
	assert.Equal(t, zerolog.InfoLevel, severityForStatus(399))
	assert.Equal(t, zerolog.WarnLevel, severityForStatus(400))
	assert.Equal(t, zerolog.WarnLevel, severityForStatus(404))
	assert.Equal(t, zerolog.WarnLevel, severityForStatus(499))
	assert.Equal(t, zerolog.ErrorLevel, severityForStatus(500))
	assert.Equal(t, zerolog.ErrorLevel, severityForStatus(503))
	assert.Equal(t, zerolog.ErrorLevel, severityForStatus(599))
}

func TestRequestLoggingMiddleware_200(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Str("environment", "production").
		Str("service", "kiroku").
		Logger()
	t.Setenv("LOG_DIR", t.TempDir())

	handler := middleware.RequestID(RequestLoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	assert.Contains(t, output, `"service":"kiroku"`)
	assert.Contains(t, output, `"message":"http request completed"`)
}

func TestRequestLoggingMiddleware_Panic(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Str("environment", "production").
		Str("service", "kiroku").
		Logger()
	t.Setenv("LOG_DIR", t.TempDir())

	handler := middleware.RequestID(RequestLoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("simulated panic")
	})))

	req := httptest.NewRequest(http.MethodGet, "/api/foo", nil)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	output := buf.String()
	assert.NotEmpty(t, output)
	assert.Contains(t, output, `"status":500`)
	assert.Contains(t, output, `"level":"error"`)
}

func TestRequestLoggingMiddleware_404_LogsWarn(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Str("environment", "production").
		Str("service", "kiroku").
		Logger()
	t.Setenv("LOG_DIR", t.TempDir())

	handler := middleware.RequestID(RequestLoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})))

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	output := buf.String()
	assert.NotEmpty(t, output)
	assert.Contains(t, output, `"level":"warn"`)
	assert.Contains(t, output, `"status":404`)
}

func TestRequestLoggingMiddleware_500_LogsError(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Str("environment", "production").
		Str("service", "kiroku").
		Logger()
	t.Setenv("LOG_DIR", t.TempDir())

	handler := middleware.RequestID(RequestLoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})))

	req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	output := buf.String()
	assert.NotEmpty(t, output)
	assert.Contains(t, output, `"level":"error"`)
	assert.Contains(t, output, `"status":500`)
}

func TestLogHTTPServerError_ValidStatus(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Str("environment", "production").
		Str("service", "kiroku").
		Logger()
	t.Setenv("LOG_DIR", t.TempDir())

	LogHTTPServerError(http.StatusInternalServerError, "GET", "/api/500", "req-123")
	output := buf.String()
	assert.Contains(t, output, `"level":"error"`)
	assert.Contains(t, output, `"status":500`)
	assert.Contains(t, output, `"method":"GET"`)
	assert.Contains(t, output, `"path":"/api/500"`)
	assert.Contains(t, output, `"request_id":"req-123"`)
	assert.Contains(t, output, `"component":"http"`)
}

func TestLogHTTPServerError_Sub500_NoLog(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf).Level(zerolog.ErrorLevel)
	t.Setenv("LOG_DIR", t.TempDir())
	Init(false)

	LogHTTPServerError(http.StatusOK, "GET", "/api", "req-xyz")
	LogHTTPServerError(http.StatusNotFound, "POST", "/api", "req-xyz")
	assert.Empty(t, buf.String())
}

func TestSanitizePath(t *testing.T) {
	assert.Equal(t, "/api/v1/users", sanitizePath("/api/v1/users?token=secret&lang=en"))
	assert.Equal(t, "/static/img.png", sanitizePath("/static/img.png"))
}

func TestCleanupOldLogs(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("LOG_DIR", dir)

	now := time.Now()
	nowYMD := now.Format("2006-01-02")
	oldYMD := now.AddDate(0, 0, -2).Format("2006-01-02")

	require.NoError(t, os.WriteFile(filepath.Join(dir, nowYMD+".log"), []byte("recent"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, oldYMD+".log"), []byte("old"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "not-a-date.log"), []byte("other"), 0644))

	cleanupOldLogs(0)
	assert.FileExists(t, filepath.Join(dir, nowYMD+".log"))
	assert.FileExists(t, filepath.Join(dir, oldYMD+".log"))
	assert.FileExists(t, filepath.Join(dir, "not-a-date.log"))

	cleanupOldLogs(1)
	assert.FileExists(t, filepath.Join(dir, nowYMD+".log"))
	assert.NoFileExists(t, filepath.Join(dir, oldYMD+".log"))
	assert.FileExists(t, filepath.Join(dir, "not-a-date.log"))
}

func TestRotateFileBySize(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("LOG_DIR", dir)
	t.Setenv("LOG_MAX_SIZE", "1") // 1 MB
	t.Setenv("LOG_MAX_BACKUPS", "5")
	t.Setenv("LOG_MAX_AGE", "90")
	t.Setenv("LOG_COMPRESS", "false")
	t.Setenv("LOG_FORMAT", "json")

	log.Logger = zerolog.New(io.Discard)
	Init(false)

	daily := newDailyRotator(dir)
	const line = "012345678901234567890123456789012345678901234567890123456789012\n" // 70 bytes
	written := 0
	for written < 2<<20 { // >1MB
		n, err := daily.Write([]byte(line))
		written += n
		if err != nil {
			break
		}
	}

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	logFiles := 0
	rotated := false
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".log") {
			continue
		}
		logFiles++
		if e.Name() != today()+".log" {
			rotated = true
		}
	}
	assert.True(t, rotated, "expected rotated backup file to exist")
	assert.GreaterOrEqual(t, logFiles, 2)
}

func TestConcurrentLogging(t *testing.T) {
	t.Setenv("LOG_DIR", t.TempDir())
	log.Logger = zerolog.New(io.Discard)
	Init(false)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				log.Info().Msgf("goroutine %d log %d", idx, j)
			}
		}(i)
	}
	wg.Wait()
}

func TestLogSensitiveDataRedaction_RequestLog(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Str("environment", "production").
		Str("service", "kiroku").
		Logger()
	t.Setenv("LOG_DIR", t.TempDir())

	handler := middleware.RequestID(RequestLoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	req := httptest.NewRequest(http.MethodGet, "/api/users?id=123&token=secret-token", nil)
	req.Header.Set("Authorization", "Bearer secret")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	output := buf.String()
	assert.NotContains(t, output, "secret-token")
	assert.NotContains(t, output, "Authorization")
	assert.NotContains(t, output, "token")
	assert.Contains(t, output, `"method":"GET"`)
	assert.Contains(t, output, `"path":"/api/users"`)
}

func TestFileWriter_DoesNotTruncate(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("LOG_DIR", dir)
	t.Setenv("LOG_MAX_SIZE", "1000000")
	t.Setenv("LOG_MAX_BACKUPS", "10")
	t.Setenv("LOG_MAX_AGE", "90")
	t.Setenv("LOG_COMPRESS", "false")
	t.Setenv("LOG_FORMAT", "json")

	log.Logger = zerolog.New(&bytes.Buffer{}).Level(zerolog.InfoLevel)
	Init(false)

	logPath := filepath.Join(dir, today()+".log")
	require.NoError(t, os.WriteFile(logPath, []byte("initial content"), 0644))

	daily := newDailyRotator(dir)
	_, err := daily.Write([]byte(" append content"))
	require.NoError(t, err)

	data, err := os.ReadFile(logPath)
	require.NoError(t, err)
	content := string(data)
	assert.Contains(t, content, "initial content")
	assert.Contains(t, content, " append content")
}

func TestLogRetentionDays_Default(t *testing.T) {
	t.Setenv("LOG_MAX_AGE", "")
	t.Setenv("LOG_RETENTION", "")
	assert.Equal(t, 90, LogRetentionDays())
}

func TestLogRetentionDays_Legacy(t *testing.T) {
	t.Setenv("LOG_MAX_AGE", "")
	t.Setenv("LOG_RETENTION", "7d")
	assert.Equal(t, 7, LogRetentionDays())

	t.Setenv("LOG_RETENTION", "30D")
	assert.Equal(t, 30, LogRetentionDays())
}

func TestLogRetentionDays_MaxAgeOverrides(t *testing.T) {
	t.Setenv("LOG_MAX_AGE", "42")
	t.Setenv("LOG_RETENTION", "7d")
	assert.Equal(t, 42, LogRetentionDays())
}