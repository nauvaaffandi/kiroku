package logger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	serviceName       = "kiroku"
	defaultLevel      = zerolog.InfoLevel
	defaultFormat     = "console"
	defaultLogDir     = "logs"
	defaultMaxSize    = 100 // megabytes
	defaultMaxBackups = 10
	defaultMaxAge     = 90 // days
	defaultCompress   = true
	logFileExt        = ".log"
	envDev            = "development"
	envProd           = "production"
	componentHTTP     = "http"
	componentAuth     = "auth"
)

// getEnv mirrors the config package's GetEnv without importing it, so the
// logger stays independent of config loading order.
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

var (
	levels = map[string]zerolog.Level{
		"debug": zerolog.DebugLevel,
		"info":  zerolog.InfoLevel,
		"warn":  zerolog.WarnLevel,
		"error": zerolog.ErrorLevel,
	}
)

func initLevel() (zerolog.Level, bool) {
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		lvl, ok := levels[strings.ToLower(v)]
		if !ok {
			return zerolog.InfoLevel, true
		}
		return lvl, false
	}
	if isDev() {
		return zerolog.DebugLevel, false
	}
	return zerolog.InfoLevel, false
}

func initFormat() (string, bool) {
	if v := os.Getenv("LOG_FORMAT"); v != "" {
		switch strings.ToLower(v) {
		case "console", "json":
			return v, false
		}
		return "console", true
	}
	if isDev() {
		return "console", false
	}
	return "json", false
}

func logDir() string {
	if v := os.Getenv("LOG_DIR"); v != "" {
		return v
	}
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return filepath.Join("backend", defaultLogDir)
	}
	dir := filepath.Dir(file)
	backend := filepath.Join(dir, "..", "..")
	return filepath.Join(backend, defaultLogDir)
}

func parseMaxSize() int {
	if v := os.Getenv("LOG_MAX_SIZE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return defaultMaxSize
}

func parseMaxBackups() int {
	if v := os.Getenv("LOG_MAX_BACKUPS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			return n
		}
	}
	return defaultMaxBackups
}

func parseMaxAge() int {
	if v := os.Getenv("LOG_MAX_AGE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			return n
		}
	}
	return 0 // unset; LogRetentionDays falls back to LOG_RETENTION
}

func parseCompress() bool {
	return os.Getenv("LOG_COMPRESS") != "false"
}

func LogRetentionDays() int {
	if age := parseMaxAge(); age > 0 {
		return age
	}
	return parseRetentionDays(getEnv("LOG_RETENTION", defaultRetention))
}

const defaultRetention = "90d"

func parseRetentionDays(s string) int {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "d")
	s = strings.TrimSuffix(s, "D")
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return 0
	}
	return n
}

type dailyRotator struct {
	mu     sync.Mutex
	dir    string
	day    string
	writer io.Writer
}

func newDailyRotator(dir string) *dailyRotator {
	return &dailyRotator{dir: dir}
}

func today() string {
	return time.Now().Format("2006-01-02")
}

func (d *dailyRotator) activePath() string {
	return filepath.Join(d.dir, today()+logFileExt)
}

func (d *dailyRotator) newWriter(path string) io.Writer {
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    parseMaxSize(),
		MaxBackups: parseMaxBackups(),
		MaxAge:     parseMaxAge(),
		Compress:   parseCompress(),
		LocalTime:  false,
	}
}

func (d *dailyRotator) Write(p []byte) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	day := today()
	if d.writer == nil || day != d.day {
		d.day = day
		if err := os.MkdirAll(d.dir, 0o755); err != nil {
			d.writer = io.Discard
			return 0, err
		}
		d.writer = d.newWriter(d.activePath())
	}
	return d.writer.Write(p)
}

func consoleWriter(format string) io.Writer {
	if format == "json" {
		return os.Stderr
	}
	return zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}
}

func Init(devMode bool) {
	zerolog.TimeFieldFormat = time.RFC3339

	level, levelInvalid := initLevel()
	format, formatInvalid := initFormat()
	dir := logDir()

	var writers []io.Writer
	writers = append(writers, consoleWriter(format))

	rotator := newDailyRotator(dir)
	writers = append(writers, rotator)

	base := zerolog.New(io.MultiWriter(writers...)).
		Level(level).
		With().
		Timestamp().
		Str("environment", environment()).
		Str("service", serviceName).
		Logger()
	log.Logger = base

	if levelInvalid {
		log.Warn().
			Str("env", "LOG_LEVEL").
			Str("fallback", "INFO").
			Msg("invalid LOG_LEVEL value; falling back to INFO")
	}
	if formatInvalid {
		log.Warn().
			Str("env", "LOG_FORMAT").
			Str("fallback", "console").
			Msg("invalid LOG_FORMAT value; falling back to console")
	}

	go cleanupOldLogs(LogRetentionDays())
}

// GetLogger returns the package-level logger so callers can attach
// component-specific fields without mutating shared global state.
//
//	logger.GetLogger().With().Str("component", "auth").Logger()
func GetLogger() *zerolog.Logger {
	l := log.Logger
	return &l
}

// NewComponentLogger returns a logger carrying the given component field.
func NewComponentLogger(component string) zerolog.Logger {
	return GetLogger().With().Str("component", component).Logger()
}

// LogFilePath returns the path of today's log file, used by the CLI tool.
func LogFilePath() string {
	return filepath.Join(logDir(), today()+logFileExt)
}

// LogsDir returns the configured logs directory path.
func LogsDir() string {
	return logDir()
}

// cleanupOldLogs removes dated log files older than retentionDays.
func cleanupOldLogs(retentionDays int) {
	if retentionDays <= 0 {
		return
	}
	entries, err := os.ReadDir(logDir())
	if err != nil {
		return
	}
	now := time.Now()
	cutoff := now.Add(-time.Duration(retentionDays) * 24 * time.Hour)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), logFileExt) {
			continue
		}
		dateStr := strings.TrimSuffix(e.Name(), logFileExt)
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		if t.Before(cutoff) {
			os.Remove(filepath.Join(logDir(), e.Name()))
		}
	}
}

// LogHTTPServerError logs critical server errors (HTTP 5xx) at Error level.
// No log is emitted for status < 500. No sensitive request data is logged.
func LogHTTPServerError(status int, method, path, requestID string) {
	if status < 500 {
		return
	}
	log.Error().
		Str("component", componentHTTP).
		Int("status", status).
		Str("method", method).
		Str("path", sanitizePath(path)).
		Str("request_id", requestID).
		Msg("critical: server error")
}

// sanitizePath strips the raw query string so sensitive query parameters are
// never written to logs.
func sanitizePath(path string) string {
	if i := strings.IndexByte(path, '?'); i >= 0 {
		return path[:i]
	}
	return path
}

// severityForStatus maps an HTTP status code to the required log level:
// 1xx-3xx INFO, 4xx WARN, 5xx ERROR.
func severityForStatus(status int) zerolog.Level {
	switch {
	case status < 400:
		return zerolog.InfoLevel
	case status < 500:
		return zerolog.WarnLevel
	default:
		return zerolog.ErrorLevel
	}
}

// statusRecorder captures the HTTP status code written by a wrapped handler.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func newStatusRecorder(w http.ResponseWriter) *statusRecorder {
	return &statusRecorder{ResponseWriter: w, status: http.StatusOK}
}

func (sr *statusRecorder) WriteHeader(statusCode int) {
	sr.status = statusCode
	sr.ResponseWriter.WriteHeader(statusCode)
}

// RequestLoggingMiddleware logs exactly one completion entry per completed
// HTTP request. It recovers panics so completion logs are never lost, maps
// status to severity, strips query strings, and never logs credentials or
// bodies. Request ID comes from the existing chi RequestID middleware.
func RequestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sr := newStatusRecorder(w)

		var panicErr any
		func() {
			defer func() { panicErr = recover() }()
			next.ServeHTTP(sr, r)
		}()

		duration := time.Since(start)
		status := sr.status
		if panicErr != nil {
			status = http.StatusInternalServerError
			sr.WriteHeader(http.StatusInternalServerError)
		}

		level := severityForStatus(status)
		evt := log.With().
			Str("component", componentHTTP).
			Str("request_id", middleware.GetReqID(r.Context())).
			Str("method", r.Method).
			Str("path", sanitizePath(r.URL.Path)).
			Int("status", status).
			Int64("duration_ms", duration.Milliseconds()).
			Str("client_ip", clientIP(r)).
			Str("user_agent", r.UserAgent()).
			Logger()

		switch level {
		case zerolog.DebugLevel:
			evt.Debug().Msg("http request completed")
		case zerolog.WarnLevel:
			evt.Warn().Msg("http request completed")
		case zerolog.ErrorLevel:
			if panicErr != nil {
				evt = evt.With().Err(fmt.Errorf("panic: %v", panicErr)).Logger()
			}
			evt.Error().Msg("http request completed")
		default:
			evt.Info().Msg("http request completed")
		}
	})
}

// clientIP extracts the client IP from the request, preferring
// X-Forwarded-For when present and falling back to RemoteAddr.
func clientIP(r *http.Request) string {
	if ip := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); ip != "" {
		return ip
	}
	if i := strings.LastIndexByte(r.RemoteAddr, ':'); i >= 0 {
		return r.RemoteAddr[:i]
	}
	return r.RemoteAddr
}

// environment reports the resolved environment name.
func environment() string {
	if isDev() {
		return envDev
	}
	return envProd
}

func isDev() bool {
	return os.Getenv("KIROKU_DEV") == "true"
}

// ErrInvalidLogLevel is returned when LOG_LEVEL contains an unrecognized value.
// (Used internally; Init produces a runtime warning instead of returning an error.)
var ErrInvalidLogLevel = fmt.Errorf("invalid LOG_LEVEL value; valid: debug, info, warn, error")
