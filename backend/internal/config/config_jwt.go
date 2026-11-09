package config

import (
	"os"
	"time"
)

// JWTSecret returns the JWT secret from JWT_SECRET env var.
func JWTSecret() string {
	return os.Getenv("JWT_SECRET")
}

// JWTAccessTTL returns the JWT access TTL from JWT_ACCESS_TTL env var, default 15m.
func JWTAccessTTL() time.Duration {
	d, err := time.ParseDuration(os.Getenv("JWT_ACCESS_TTL"))
	if err != nil {
		return 15 * time.Minute
	}
	return d
}

// JWTRefreshTTL returns the JWT refresh TTL from JWT_REFRESH_TTL env var, default 168h.
func JWTRefreshTTL() time.Duration {
	d, err := time.ParseDuration(os.Getenv("JWT_REFRESH_TTL"))
	if err != nil {
		return 168 * time.Hour
	}
	return d
}

// Argon2Time returns the Argon2 time cost from ARGON2_TIME env var, default 3.
func Argon2Time() uint32 {
	v := retentionEnv("ARGON2_TIME", 3)
	return uint32(v)
}

// Argon2Memory returns the Argon2 memory cost from ARGON2_MEMORY env var, default 65536.
func Argon2Memory() uint32 {
	v := retentionEnv("ARGON2_MEMORY", 65536)
	return uint32(v)
}

// Argon2Threads returns the Argon2 parallelism from ARGON2_THREADS env var, default 4.
func Argon2Threads() uint8 {
	v := retentionEnv("ARGON2_THREADS", 4)
	return uint8(v)
}

// KIROKUDev returns true if KIROKU_DEV is "true".
func KIROKUDev() bool {
	return os.Getenv("KIROKU_DEV") == "true"
}
