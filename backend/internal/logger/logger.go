package logger

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/config"
)

const (
	defaultRetention = "90d"
	logsDirName      = "logs"
	logFileExt       = ".log"
)

func logsDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return filepath.Join("backend", logsDirName)
	}
	dir := filepath.Dir(file) // backend/internal/logger
	backend := filepath.Join(dir, "..", "..")
	return filepath.Join(backend, logsDirName)
}

func Init(devMode bool) {
	zerolog.TimeFieldFormat = time.RFC3339

	cfg, err := config.Load()
	retentionDays := parseRetentionDays(defaultRetention)
	if err == nil && cfg.LogRetention != "" {
		if d := parseRetentionDays(cfg.LogRetention); d > 0 {
			retentionDays = d
		}
	}

	var writers []io.Writer

	if devMode {
		writers = append(writers, zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		})
	} else {
		writers = append(writers, os.Stdout)
	}

	if dirErr := os.MkdirAll(logsDir(), 0755); dirErr == nil {
		today := time.Now().Format("2006-01-02")
		path := filepath.Join(logsDir(), today+logFileExt)
		if f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil {
			writers = append(writers, f)
		}
	}

	mw := io.MultiWriter(writers...)

	base := zerolog.New(mw).Level(zerolog.InfoLevel).With().Timestamp().Logger()
	if devMode {
		base = base.Level(zerolog.DebugLevel)
	}
	log.Logger = base

	go cleanupOldLogs(retentionDays)
}

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

func cleanupOldLogs(retentionDays int) {
	if retentionDays <= 0 {
		return
	}
	entries, err := os.ReadDir(logsDir())
	if err != nil {
		return
	}
	now := time.Now()
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), logFileExt) {
			continue
		}
		dateStr := strings.TrimSuffix(e.Name(), logFileExt)
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		if now.Sub(t) > time.Duration(retentionDays)*24*time.Hour {
			os.Remove(filepath.Join(logsDir(), e.Name()))
		}
	}
}

func GetLogger() *zerolog.Logger {
	l := log.Logger
	return &l
}

// LogFilePath returns the path of today's log file, used by the CLI tool.
func LogFilePath() string {
	return filepath.Join(logsDir(), time.Now().Format("2006-01-02")+logFileExt)
}

// LogsDir returns the logs directory path.
func LogsDir() string {
	return logsDir()
}
