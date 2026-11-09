package config

import (
	"os"
	"strings"
)

// AllowedCORSOrigins returns the allowed CORS origins from CORS_ALLOWED_ORIGINS env var.
func AllowedCORSOrigins() []string {
	if v := os.Getenv("CORS_ALLOWED_ORIGINS"); v != "" {
		var origins []string
		parts := strings.Split(v, ",")
		for _, p := range parts {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				origins = append(origins, trimmed)
			}
		}
		return origins
	}
	return []string{}
}
