package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInitDirectories menguji fungsi inisialisasi direktori DATA_DIR/persistent
// dan DATA_DIR/runtime saat startup.
//
// Table-driven test dengan tiga skenario:
//   - DATA_DIR diset: kedua direktori dibuat dengan permission 0o700
//   - DATA_DIR kosong: tidak ada direktori dibuat, tidak ada error
//   - Direktori sudah ada: tidak ada error (idempotent via MkdirAll)
func TestInitDirectories(t *testing.T) {
	tests := []struct {
		name       string
		setDataDir bool // true: isi DATA_DIR dengan t.TempDir, false: kosongkan
		preCreate  bool // true: buat direktori sebelum pemanggilan (uji idempotensi)
		wantErr    bool
	}{
		{
			name:       "DATA_DIR diset: kedua direktori dibuat dengan permission 0o700",
			setDataDir: true,
			wantErr:    false,
		},
		{
			name:       "DATA_DIR kosong: tidak ada direktori dibuat, tidak ada error",
			setDataDir: false,
			wantErr:    false,
		},
		{
			name:       "Direktori sudah ada: tidak ada error (idempotent via MkdirAll)",
			setDataDir: true,
			preCreate:  true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Isolasi tiap subtest dengan t.TempDir
			tempDir := t.TempDir()

			if tt.setDataDir {
				t.Setenv("DATA_DIR", tempDir)
			} else {
				t.Setenv("DATA_DIR", "")
			}

			// Untuk skenario idempotensi: buat direktori terlebih dahulu
			if tt.preCreate {
				require.NoError(t, os.MkdirAll(filepath.Join(tempDir, "persistent"), 0o700))
				require.NoError(t, os.MkdirAll(filepath.Join(tempDir, "runtime"), 0o700))
			}

			err := InitDirectories()

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if tt.setDataDir {
				// Verifikasi kedua direktori dibuat dengan permission 0o700
				for _, sub := range []string{"persistent", "runtime"} {
					path := filepath.Join(tempDir, sub)
					info, err := os.Stat(path)
					require.NoError(t, err, "direktori %s harus ada", sub)
					assert.Equal(t, os.FileMode(0o700), info.Mode().Perm(), "permission direktori %s harus 0o700", sub)
				}
			}
		})
	}
}
