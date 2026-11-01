package auth

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kiroku/backend/internal/config"
	"kiroku/backend/internal/db"
)

func setupFreshHandler(t *testing.T, dataDir string) *Handler {
	t.Helper()
	os.Setenv("DATA_DIR", dataDir)
	t.Cleanup(func() { os.Unsetenv("DATA_DIR") })

	database, err := db.Open(dataDir)
	require.NoError(t, err)
	t.Cleanup(func() { database.Close() })

	return NewHandler(database, &config.Config{}, nil)
}

func TestGetOrCreateUser_GeneratesRandomPIN(t *testing.T) {
	dataDir := t.TempDir()
	h := setupFreshHandler(t, dataDir)

	userID, err := h.getOrCreateUser()
	require.NoError(t, err)
	assert.NotEmpty(t, userID)

	pinPath := filepath.Join(dataDir, ".seed", "seed-pin.txt")
	content, err := os.ReadFile(pinPath)
	require.NoError(t, err)

	pin := string(content)
	assert.Regexp(t, regexp.MustCompile(`^\d{6}$`), pin, "seed PIN must be exactly 6 digits")
}

func TestGetOrCreateUser_ExistingUser(t *testing.T) {
	dataDir := t.TempDir()

	database, err := db.Open(dataDir)
	require.NoError(t, err)
	defer database.Close()

	hash, err := HashPin("654321", config.Argon2Time(), config.Argon2Memory(), config.Argon2Threads())
	require.NoError(t, err)
	_, err = database.Exec("INSERT INTO users (id, email, pin_hash) VALUES (?, ?, ?)", "existing-user", "existing@example.com", hash)
	require.NoError(t, err)

	h := NewHandler(database, &config.Config{}, nil)
	userID, err := h.getOrCreateUser()
	require.NoError(t, err)
	assert.Equal(t, "existing-user", userID)

	pinPath := filepath.Join(dataDir, ".seed", "seed-pin.txt")
	_, err = os.ReadFile(pinPath)
	assert.Error(t, err, "seed-pin.txt should not exist for existing user")
	assert.True(t, os.IsNotExist(err))
}

func TestGetOrCreateUser_PinIsSixDigits(t *testing.T) {
	dataDir := t.TempDir()
	h := setupFreshHandler(t, dataDir)

	_, err := h.getOrCreateUser()
	require.NoError(t, err)

	pinPath := filepath.Join(dataDir, ".seed", "seed-pin.txt")
	content, err := os.ReadFile(pinPath)
	require.NoError(t, err)

	pin := string(content)
	assert.Equal(t, 6, len(pin))
	assert.Regexp(t, regexp.MustCompile(`^\d{6}$`), pin)
}

func TestSeedPINFileWrittenOnlyOnce(t *testing.T) {
	dataDir := t.TempDir()
	h := setupFreshHandler(t, dataDir)

	userID1, err := h.getOrCreateUser()
	require.NoError(t, err)

	pinPath := filepath.Join(dataDir, ".seed", "seed-pin.txt")
	content1, err := os.ReadFile(pinPath)
	require.NoError(t, err)

	userID2, err := h.getOrCreateUser()
	require.NoError(t, err)
	assert.Equal(t, userID1, userID2)

	content2, err := os.ReadFile(pinPath)
	require.NoError(t, err)
	assert.Equal(t, string(content1), string(content2))
}