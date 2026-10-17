package backup

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"kiroku/backend/internal/models"
)

func (h *Handler) CreateBackup(tier string) (*models.Backup, error) {
	backupDir := "data/backups"
	os.MkdirAll(backupDir, 0755)

	timestamp := time.Now().Format("20060102_150405")
	filename := "kiroku_" + timestamp + ".tar.gz"
	if h.cfg.BackupPassphrase != "" {
		filename += ".enc"
	}
	filePath := filepath.Join(backupDir, filename)

	// Create temp tar.gz
	tmpTar := filepath.Join(backupDir, "temp_"+timestamp+".tar.gz")
	outFile, err := os.Create(tmpTar)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpTar)

	gw := gzip.NewWriter(outFile)
	tw := tar.NewWriter(gw)

	// 1. Manifest
	manifest := map[string]interface{}{
		"version":           "1.0",
		"timestamp":         time.Now().Format(time.RFC3339),
		"db_schema_version": 1,
	}
	manifestBytes, _ := json.Marshal(manifest)
	tw.WriteHeader(&tar.Header{
		Name:    "manifest.json",
		Size:    int64(len(manifestBytes)),
		Mode:    0644,
		ModTime: time.Now(),
	})
	tw.Write(manifestBytes)

	// 2. Database backup
	dbPath := "data/kiroku.db"
	if _, err := os.Stat(dbPath); err == nil {
		dbFile, err := os.Open(dbPath)
		if err == nil {
			stat, _ := dbFile.Stat()
			tw.WriteHeader(&tar.Header{
				Name:    "kiroku.db",
				Size:    stat.Size(),
				Mode:    0644,
				ModTime: stat.ModTime(),
			})
			io.Copy(tw, dbFile)
			dbFile.Close()
		}
	}

	// 3. Attachments
	attachmentsDir := "data/attachments"
	if _, err := os.Stat(attachmentsDir); err == nil {
		filepath.Walk(attachmentsDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			rel, _ := filepath.Rel("data", path)
			f, err := os.Open(path)
			if err != nil {
				return nil
			}
			defer f.Close()

			tw.WriteHeader(&tar.Header{
				Name:    rel,
				Size:    info.Size(),
				Mode:    0644,
				ModTime: info.ModTime(),
			})
			io.Copy(tw, f)
			return nil
		})
	}

	tw.Close()
	gw.Close()
	outFile.Close()

	var finalSize int64
	if h.cfg.BackupPassphrase != "" {
		// Encrypt file
		plaintext, err := os.ReadFile(tmpTar)
		if err != nil {
			return nil, err
		}
		ciphertext, err := encrypt(plaintext, h.cfg.BackupPassphrase)
		if err != nil {
			return nil, err
		}
		os.WriteFile(filePath, ciphertext, 0644)
		stat, _ := os.Stat(filePath)
		finalSize = stat.Size()
	} else {
		os.Rename(tmpTar, filePath)
		stat, _ := os.Stat(filePath)
		finalSize = stat.Size()
	}

	id := uuid.New().String()
	_, err = h.db.Exec(
		"INSERT INTO backups (id, filename, size, retention_tier, created_at) VALUES (?, ?, ?, ?, ?)",
		id, filename, finalSize, tier, time.Now(),
	)
	if err != nil {
		return nil, err
	}

	return &models.Backup{
		ID:            id,
		Filename:      filename,
		Size:          finalSize,
		RetentionTier: tier,
		CreatedAt:     time.Now(),
	}, nil
}
