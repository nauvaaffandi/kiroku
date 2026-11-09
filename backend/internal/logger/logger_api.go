package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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
