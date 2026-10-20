package db

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"kiroku/backend/internal/config"
	"kiroku/backend/internal/logger"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type DB struct {
	*sql.DB
}

func Open(cfg *config.Config) (*DB, error) {
	dbPath := cfg.FullDBPath()

	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)&_pragma=journal_mode(WAL)", dbPath)
	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	db := &DB{sqlDB}
	if err := db.runMigrations(); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	logger.GetLogger().Info().Str("path", dbPath).Msg("Database opened with WAL mode")
	return db, nil
}

func (d *DB) runMigrations() error {
	_, err := d.Exec(`
		CREATE TABLE IF NOT EXISTS schema_version (
			version INTEGER PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	var currentVersion int
	err = d.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_version").Scan(&currentVersion)
	if err != nil {
		return err
	}

	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}

		var version int
		_, err := fmt.Sscanf(entry.Name(), "%d_", &version)
		if err != nil || version <= currentVersion {
			continue
		}

		content, err := migrationFiles.ReadFile(filepath.Join("migrations", entry.Name()))
		if err != nil {
			return err
		}

		if version == 13 {
			applied, err := d.columnExists("users", "pin_hash")
			if err != nil {
				return err
			}
			if applied {
				if _, err := d.Exec("INSERT INTO schema_version (version) VALUES (?)", version); err != nil {
					return err
				}
				logger.GetLogger().Info().Int("version", version).Str("file", entry.Name()).Msg("Migration already applied, skipping")
				currentVersion = version
				continue
			}
		}

		tx, err := d.Begin()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %s: %w", entry.Name(), err)
		}

		if _, err := tx.Exec("INSERT INTO schema_version (version) VALUES (?)", version); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		logger.GetLogger().Info().Int("version", version).Str("file", entry.Name()).Msg("Applied migration")
		currentVersion = version
	}

	return nil
}

func (d *DB) columnExists(table, column string) (bool, error) {
	rows, err := d.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var cid, notnull, pk int
		var name, ctype string
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			return false, err
		}
		if name == column {
			return true, nil
		}
	}
	return false, nil
}
