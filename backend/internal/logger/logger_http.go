package logger

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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
