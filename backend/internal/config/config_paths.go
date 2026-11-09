package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// projectRoot is the directory containing backend/go.mod, resolved once at init
// by walking up from the source file location (runtime.Caller) until a go.mod is
// found, then falling back to walking up from the executable. Relative DATA_DIR
// values are resolved against it. Works with `go run` because the source tree is
// present. Assumes this file lives at backend/internal/config/config.go relative
// to the project root.
var projectRoot string

func init() {
	for _, dir := range sourceDirs() {
		if root, ok := findGoModDir(dir); ok {
			projectRoot = root
			return
		}
	}
}

// sourceDirs returns candidate ancestor directories to search for go.mod,
// in priority order: the directory of this source file, then the executable
// directory (for compiled binaries deployed near the project).
func sourceDirs() []string {
	var dirs []string
	if _, file, _, ok := runtime.Caller(0); ok {
		if dir := filepath.Dir(file); dir != "" {
			dirs = append(dirs, dir)
		}
	}
	if exe, err := os.Executable(); err == nil {
		if dir := filepath.Dir(exe); dir != "" {
			dirs = append(dirs, dir)
		}
	}
	return dirs
}

// findGoModDir walks up from dir until a directory containing go.mod is found.
func findGoModDir(dir string) (string, bool) {
	cur := dir
	for {
		if _, err := os.Stat(filepath.Join(cur, "go.mod")); err == nil {
			return cur, true
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return "", false
		}
		cur = parent
	}
}

// DataDir returns the data directory from DATA_DIR env var. Relative paths are
// resolved against the project root (backend/). Absolute paths and empty values
// are returned as-is (empty skips directory init).
func DataDir() string {
	v := os.Getenv("DATA_DIR")
	if v == "" {
		return ""
	}
	if filepath.IsAbs(v) {
		return v
	}
	if projectRoot == "" {
		return v
	}
	return filepath.Join(projectRoot, v)
}

// ValidateDataDir ensures DATA_DIR is safe: absolute path, no traversal.
func ValidateDataDir(path string) error {
	if !filepath.IsAbs(path) {
		return errors.New("DATA_DIR must be an absolute path")
	}
	if strings.Contains(path, "..") {
		return errors.New("DATA_DIR must not contain '..'")
	}
	return nil
}

// InitDirectories creates DATA_DIR/persistent and DATA_DIR/runtime with 0o700.
// Skips when DATA_DIR is empty. Returns generic error on MkdirAll or Chmod failure
// to avoid leaking internal paths.
func InitDirectories() error {
	dataDir := DataDir()
	if dataDir == "" {
		return nil
	}
	if err := ValidateDataDir(dataDir); err != nil {
		return err
	}
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		return errors.New("failed to create data directory")
	}
	if err := os.Chmod(dataDir, 0o700); err != nil {
		return errors.New("failed to set data directory permissions")
	}
	for _, d := range []string{"persistent", "runtime"} {
		path := filepath.Join(dataDir, d)
		if err := os.MkdirAll(path, 0o700); err != nil {
			return errors.New("failed to create data subdirectory")
		}
		if err := os.Chmod(path, 0o700); err != nil {
			return errors.New("failed to set data subdirectory permissions")
		}
	}
	return nil
}

// GenesisDir returns the genesis data directory path under DATA_DIR/persistent.
func GenesisDir() string {
	return filepath.Join(DataDir(), "persistent")
}

// GenesisPinPath returns the path to genesis-pin.txt.
func GenesisPinPath() string {
	return filepath.Join(GenesisDir(), "genesis-pin.txt")
}

// RuntimeDir returns the runtime data directory path under DATA_DIR/runtime.
func RuntimeDir() string {
	return filepath.Join(DataDir(), "runtime")
}

// BackendURLPath returns the path to backend-url.txt under the runtime data directory.
func BackendURLPath() string {
	return filepath.Join(RuntimeDir(), "backend-url.txt")
}
