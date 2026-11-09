package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Load(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("JWT_SECRET", "test-secret-do-not-use-in-production")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("JWT_SECRET")

	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, ":9090", cfg.Addr)
}

func TestConfig_LoadDefaultAddr(t *testing.T) {
	os.Unsetenv("PORT")
	os.Setenv("JWT_SECRET", "test-secret-do-not-use-in-production")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("JWT_SECRET")

	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, ":8080", cfg.Addr)
}

func TestFXBaseCurrency(t *testing.T) {
	os.Setenv("FX_BASE_CURRENCY", "USD")
	defer os.Unsetenv("FX_BASE_CURRENCY")
	assert.Equal(t, "USD", FXBaseCurrency())
}

func TestFXBaseCurrencyDefault(t *testing.T) {
	os.Unsetenv("FX_BASE_CURRENCY")
	defer os.Unsetenv("FX_BASE_CURRENCY")
	assert.Equal(t, "JPY", FXBaseCurrency())
}

func TestBackupSchedule(t *testing.T) {
	os.Setenv("BACKUP_SCHEDULE", "0 0 * * *")
	defer os.Unsetenv("BACKUP_SCHEDULE")
	assert.Equal(t, "0 0 * * *", BackupSchedule())
}

func TestBackupScheduleDefault(t *testing.T) {
	os.Unsetenv("BACKUP_SCHEDULE")
	defer os.Unsetenv("BACKUP_SCHEDULE")
	assert.Equal(t, "0 3 * * *", BackupSchedule())
}

func TestBackupRetentionDaily(t *testing.T) {
	os.Setenv("BACKUP_RETENTION_DAILY", "14")
	defer os.Unsetenv("BACKUP_RETENTION_DAILY")
	assert.Equal(t, 14, BackupRetentionDaily())
}

func TestBackupRetentionDailyDefault(t *testing.T) {
	os.Unsetenv("BACKUP_RETENTION_DAILY")
	defer os.Unsetenv("BACKUP_RETENTION_DAILY")
	assert.Equal(t, 7, BackupRetentionDaily())
}

func TestBackupRetentionWeekly(t *testing.T) {
	os.Setenv("BACKUP_RETENTION_WEEKLY", "8")
	defer os.Unsetenv("BACKUP_RETENTION_WEEKLY")
	assert.Equal(t, 8, BackupRetentionWeekly())
}

func TestBackupRetentionWeeklyDefault(t *testing.T) {
	os.Unsetenv("BACKUP_RETENTION_WEEKLY")
	defer os.Unsetenv("BACKUP_RETENTION_WEEKLY")
	assert.Equal(t, 4, BackupRetentionWeekly())
}

func TestBackupRetentionMonthly(t *testing.T) {
	os.Setenv("BACKUP_RETENTION_MONTHLY", "12")
	defer os.Unsetenv("BACKUP_RETENTION_MONTHLY")
	assert.Equal(t, 12, BackupRetentionMonthly())
}

func TestBackupRetentionMonthlyDefault(t *testing.T) {
	os.Unsetenv("BACKUP_RETENTION_MONTHLY")
	defer os.Unsetenv("BACKUP_RETENTION_MONTHLY")
	assert.Equal(t, 6, BackupRetentionMonthly())
}
