package auth

import (
	"database/sql"
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

func TestEnsureGenesisUser(t *testing.T) {
	tests := []struct {
		name             string
		preExistingUser  bool
		preExistingPin   bool
		envEmpty         bool
		persistentAsFile bool
		closeDB          bool
		wantErr          bool
	}{
		{
			name:    "DB kosong, DATA_DIR diset: buat user dan genesis-pin.txt",
			wantErr: false,
		},
		{
			name:            "DB sudah ada user, genesis-pin.txt ada: skip, tidak generate ulang",
			preExistingUser: true,
			preExistingPin:  true,
			wantErr:         false,
		},
		{
			name:            "DB sudah ada user, genesis-pin.txt hilang: skip, jangan buat file ulang",
			preExistingUser: true,
			preExistingPin:  false,
			wantErr:         false,
		},
		{
			name:     "DB kosong, DATA_DIR kosong: skip, return nil",
			envEmpty: true,
			wantErr:  false,
		},
		{
			name:    "DB error (query gagal): return error",
			closeDB: true,
			wantErr: true,
		},
		{
			name:             "genesis-pin.txt gagal ditulis (persistent adalah file): return error, tidak ada file orphan",
			persistentAsFile: true,
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataDir := t.TempDir()

			if tt.envEmpty {
				t.Setenv("DATA_DIR", "")
			} else {
				t.Setenv("DATA_DIR", dataDir)
			}

			database, err := db.Open(dataDir)
			require.NoError(t, err)
			defer database.Close()

			if tt.preExistingUser {
				hash, err := HashPin("123456", config.Argon2Time(), config.Argon2Memory(), config.Argon2Threads())
				require.NoError(t, err)
				_, err = database.Exec(
					"INSERT INTO users (id, email, pin_hash) VALUES (?, ?, ?)",
					"existing-user", "existing@example.com", hash,
				)
				require.NoError(t, err)
			}

			if tt.preExistingPin {
				pinPath := filepath.Join(dataDir, "persistent", "genesis-pin.txt")
				require.NoError(t, os.MkdirAll(filepath.Dir(pinPath), 0o700))
				require.NoError(t, os.WriteFile(pinPath, []byte("123456"), 0o600))
			}

			if tt.persistentAsFile {
				persistentPath := filepath.Join(dataDir, "persistent")
				require.NoError(t, os.WriteFile(persistentPath, []byte{}, 0o644))
			}

			if tt.closeDB {
				database.Close()
			}

			err = EnsureGenesisUser(database)

			if tt.wantErr {
				require.Error(t, err)

				// Ketika writeGenesisPin gagal (persistentAsFile), tidak ada genesis-pin.txt yang ditulis
				// sehingga tidak ada file orphan yang perlu dibersihkan — rollback implicit.
				if tt.persistentAsFile {
					pinPath := filepath.Join(dataDir, "persistent", "genesis-pin.txt")
					_, err := os.ReadFile(pinPath)
					assert.Error(t, err, "tidak boleh ada genesis-pin.txt orphan saat write gagal")
				}
				return
			}
			require.NoError(t, err)

			var userID string
			err = database.QueryRow("SELECT id FROM users LIMIT 1").Scan(&userID)

			if tt.preExistingUser {
				require.NoError(t, err)
				assert.Equal(t, "existing-user", userID)

				pinPath := filepath.Join(dataDir, "persistent", "genesis-pin.txt")
				if tt.preExistingPin {
					content, err := os.ReadFile(pinPath)
					require.NoError(t, err)
					assert.Equal(t, "123456", string(content))
				} else {
					_, err := os.ReadFile(pinPath)
					assert.True(t, os.IsNotExist(err), "genesis-pin.txt tidak boleh dibuat ulang")
				}
			} else if tt.envEmpty {
				require.ErrorIs(t, err, sql.ErrNoRows, "tidak boleh ada user yang dibuat saat DATA_DIR kosong")
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, userID)

				pinPath := filepath.Join(dataDir, "persistent", "genesis-pin.txt")
				content, err := os.ReadFile(pinPath)
				require.NoError(t, err)
				assert.Regexp(t, regexp.MustCompile(`^\d{6}$`), string(content), "genesis PIN harus 6 digit angka")
			}
		})
	}
}
