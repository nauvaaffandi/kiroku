package config

import (
	"os"
)

// UpstashRedisRestURL returns the Upstash Redis REST URL.
func UpstashRedisRestURL() string {
	return GetEnv("UPSTASH_REDIS_REST_URL", "")
}

func UpstashRedisRestToken() string {
	return os.Getenv("BACKEND_UPSTASH_REDIS_REST_TOKEN")
}

// LogRetention returns the log retention from LOG_RETENTION env var, default "90d".
func LogRetention() string {
	return GetEnv("LOG_RETENTION", "90d")
}

// TunnelEnabled returns true unless CLOUDFLARE_TUNNEL_ENABLED is "false".
func TunnelEnabled() bool {
	return os.Getenv("CLOUDFLARE_TUNNEL_ENABLED") != "false"
}
