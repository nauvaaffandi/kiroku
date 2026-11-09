package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Addr string
}

func loadDotEnv() {
	paths := []string{".env", "../.env"}
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(b), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
				if _, exists := os.LookupEnv(key); !exists {
					os.Setenv(key, val)
				}
			}
		}
	}
}

func Load() (*Config, error) {
	loadDotEnv()
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}
	if os.Getenv("JWT_SECRET") == "" {
		return nil, fmt.Errorf("JWT_SECRET is required; set a strong non-empty signing secret")
	}
	return &Config{
		Addr: ":" + addr,
	}, nil
}

func GetEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// FXBaseCurrency returns the base currency from FX_BASE_CURRENCY env var, default "JPY".
func FXBaseCurrency() string {
	if v := os.Getenv("FX_BASE_CURRENCY"); v != "" {
		return v
	}
	return "JPY"
}

// BackupSchedule returns the backup cron schedule from BACKUP_SCHEDULE env var, default "0 3 * * *".
func BackupSchedule() string {
	if v := os.Getenv("BACKUP_SCHEDULE"); v != "" {
		return v
	}
	return "0 3 * * *"
}

// BackupRetentionDaily returns daily retention count from BACKUP_RETENTION_DAILY env var, default 7.
func BackupRetentionDaily() int {
	return retentionEnv("BACKUP_RETENTION_DAILY", 7)
}

// BackupRetentionWeekly returns weekly retention count from BACKUP_RETENTION_WEEKLY env var, default 4.
func BackupRetentionWeekly() int {
	return retentionEnv("BACKUP_RETENTION_WEEKLY", 4)
}

// BackupRetentionMonthly returns monthly retention count from BACKUP_RETENTION_MONTHLY env var, default 6.
func BackupRetentionMonthly() int {
	return retentionEnv("BACKUP_RETENTION_MONTHLY", 6)
}

func retentionEnv(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

// FullDBPath returns the full path to the database file.
func FullDBPath(dataDir string) string {
	return filepath.Join(dataDir, "kiroku.db")
}
