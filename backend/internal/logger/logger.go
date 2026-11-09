package logger

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
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
		"trace": zerolog.TraceLevel,
		"debug": zerolog.DebugLevel,
		"info":  zerolog.InfoLevel,
		"warn":  zerolog.WarnLevel,
		"error": zerolog.ErrorLevel,
		"fatal": zerolog.FatalLevel,
		"panic": zerolog.PanicLevel,
	}
)

func initLevel() (zerolog.Level, bool) {
	v := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if v == "" {
		if isDev() {
			return zerolog.DebugLevel, false
		}
		return zerolog.InfoLevel, false
	}
	lvl, ok := levels[v]
	if !ok {
		return zerolog.InfoLevel, true
	}
	return lvl, false
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
