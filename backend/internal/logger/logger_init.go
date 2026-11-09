package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// consoleWriter returns the log writer for the given format.
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

// Init initializes the logger with configuration from environment.
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
