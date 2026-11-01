package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenesisPaths(t *testing.T) {
	dataDir := t.TempDir()
	os.Setenv("DATA_DIR", dataDir)
	t.Cleanup(func() { os.Unsetenv("DATA_DIR") })

	expectedDir := filepath.Join(dataDir, "persistent")
	assert.Equal(t, expectedDir, GenesisDir())
	assert.Equal(t, filepath.Join(expectedDir, "genesis-pin.txt"), GenesisPinPath())
	assert.Equal(t, filepath.Join(expectedDir, "backend-url.txt"), SeedURLPath())
}

func TestGenesisPaths_EmptyDataDir(t *testing.T) {
	os.Unsetenv("DATA_DIR")
	t.Cleanup(func() { os.Setenv("DATA_DIR", "") })

	dir := GenesisDir()
	assert.Equal(t, "persistent", dir)

	pinPath := GenesisPinPath()
	assert.Equal(t, filepath.Join("persistent", "genesis-pin.txt"), pinPath)

	urlPath := SeedURLPath()
	assert.Equal(t, filepath.Join("persistent", "backend-url.txt"), urlPath)
}