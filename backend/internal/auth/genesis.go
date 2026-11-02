package auth

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/config"
	"kiroku/backend/internal/db"
)

// EnsureGenesisUser memastikan ada user default di database saat startup.
// Jika user sudah ada → skip (tidak generate ulang, tidak menulis genesis-pin.txt).
// Jika user tidak ada → generate user + random 6-digit PIN, INSERT user dulu,
// baru tulis PIN ke temp file dan rename secara atomic ke genesis-pin.txt.
// Urutan ini mencegah file PIN orphan jika crash sebelum commit.
// Error → return error (dihandle log.Fatal di main.go).
func EnsureGenesisUser(database *db.DB) error {
	dataDir := config.DataDir()
	if dataDir == "" {
		return nil
	}

	if err := config.ValidateDataDir(dataDir); err != nil {
		return err
	}

	tx, err := database.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Cek apakah user sudah ada di DB
	var userID string
	err = tx.QueryRow("SELECT id FROM users LIMIT 1").Scan(&userID)
	if err == nil {
		// user ada, commit dan skip
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit: %w", err)
		}
		return nil
	}
	if err != sql.ErrNoRows {
		tx.Rollback()
		return err // DB error
	}

	// User tidak ada → generate PIN dan hash
	pin, err := generateRandomPIN()
	if err != nil {
		tx.Rollback()
		return err
	}
	hash, err := HashPin(pin,
		config.Argon2Time(),
		config.Argon2Memory(),
		config.Argon2Threads(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// INSERT user terlebih dahulu sebelum menulis file PIN.
	// Jika INSERT gagal, tidak ada file yang tercipta (tanpa orphan).
	userID = uuid.New().String()
	_, err = tx.Exec(
		"INSERT INTO users (id, email, pin_hash) VALUES (?, ?, ?)",
		userID, "user@kiroku.local", hash,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction sebelum menulis PIN.
	// Jika commit gagal, tidak ada file yang tercipta (tanpa orphan).
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	// Tulis PIN ke temp file, lalu rename secara atomic.
	// Jika write atau rename gagal, server tidak start (return error),
	// tapi tidak ada file orphan.
	if err := writeGenesisPin(pin); err != nil {
		return err
	}

	return nil
}

// writeGenesisPIN menulis plaintext PIN ke {DATA_DIR}/persistent/genesis-pin.txt
// menggunakan temp file + atomic rename untuk mencegah file orphan jika crash.
// Mengembalikan error pada kegagalan sehingga caller bisa abort.
func writeGenesisPin(pin string) error {
	dataDir := config.DataDir()
	if dataDir == "" {
		log.Warn().Msg("auth: DATA_DIR not set, skipping genesis-pin.txt write")
		return nil
	}
	dir := config.GenesisDir()
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create genesis dir: %w", err)
	}

	// Tulis ke temp file di direktori yang sama.
	tmpPath := config.GenesisPinPath() + ".tmp"
	if err := os.WriteFile(tmpPath, []byte(pin), 0o600); err != nil {
		return fmt.Errorf("write genesis-pin.txt.tmp: %w", err)
	}

	// Atomic rename ke path akhir.
	if err := os.Rename(tmpPath, config.GenesisPinPath()); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("rename genesis-pin.txt: %w", err)
	}

	return nil
}