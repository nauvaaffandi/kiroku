package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeedPaths(t *testing.T) {
	dataDir := t.TempDir()
	os.Setenv("DATA_DIR", dataDir)
	t.Cleanup(func() { os.Unsetenv("DATA_DIR") })

	expectedDir := filepath.Join(dataDir, ".seed")
	assert.Equal(t, expectedDir, SeedDir())
	assert.Equal(t, filepath.Join(expectedDir, "seed-pin.txt"), SeedPINPath())
	assert.Equal(t, filepath.Join(expectedDir, "backend-url.txt"), SeedURLPath())
}

func TestSeedPaths_EmptyDataDir(t *testing.T) {
	os.Unsetenv("DATA_DIR")
	t.Cleanup(func() { os.Setenv("DATA_DIR", "") })

	dir := SeedDir()
	assert.Equal(t, ".seed", dir)

	pinPath := SeedPINPath()
	assert.Equal(t, filepath.Join(".seed", "seed-pin.txt"), pinPath)

	urlPath := SeedURLPath()
	assert.Equal(t, filepath.Join(".seed", "backend-url.txt"), urlPath)
}