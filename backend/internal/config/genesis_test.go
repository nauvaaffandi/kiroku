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
	expectedRuntimeDir := filepath.Join(dataDir, "runtime")
	assert.Equal(t, expectedDir, GenesisDir())
	assert.Equal(t, filepath.Join(expectedDir, "genesis-pin.txt"), GenesisPinPath())
	assert.Equal(t, expectedRuntimeDir, RuntimeDir())
	assert.Equal(t, filepath.Join(expectedRuntimeDir, "backend-url.txt"), BackendURLPath())
}

func TestGenesisPaths_EmptyDataDir(t *testing.T) {
	os.Unsetenv("DATA_DIR")
	t.Cleanup(func() { os.Setenv("DATA_DIR", "") })

	dir := GenesisDir()
	assert.Equal(t, "persistent", dir)

	pinPath := GenesisPinPath()
	assert.Equal(t, filepath.Join("persistent", "genesis-pin.txt"), pinPath)

	urlPath := BackendURLPath()
	assert.Equal(t, filepath.Join("runtime", "backend-url.txt"), urlPath)
}