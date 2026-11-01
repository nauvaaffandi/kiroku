package logger

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRetentionDays(t *testing.T) {
	assert.Equal(t, 90, parseRetentionDays("90d"))
	assert.Equal(t, 30, parseRetentionDays("30D"))
	assert.Equal(t, 0, parseRetentionDays("abc"))
	assert.Equal(t, 0, parseRetentionDays("-5d"))
}

func TestCleanupOldLogs(t *testing.T) {
	dir := logsDir()
	require.NoError(t, os.MkdirAll(dir, 0755))

	old := time.Now().Add(-100 * 24 * time.Hour).Format("2006-01-02") + ".log"
	recent := time.Now().Add(-10 * 24 * time.Hour).Format("2006-01-02") + ".log"
	weird := "not-a-date.log"

	require.NoError(t, os.WriteFile(filepath.Join(dir, old), []byte("{}"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, recent), []byte("{}"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, weird), []byte("{}"), 0644))

	cleanupOldLogs(90)

	assert.False(t, fileExists(filepath.Join(dir, old)), "old log should be removed")
	assert.True(t, fileExists(filepath.Join(dir, recent)), "recent log should be kept")
	assert.True(t, fileExists(filepath.Join(dir, weird)), "non-matching name should be kept")

	os.Remove(filepath.Join(dir, recent))
	os.Remove(filepath.Join(dir, weird))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
