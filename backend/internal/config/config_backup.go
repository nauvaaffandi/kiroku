package config

import (
	"os"
)

// BackupPassphrase returns the backup passphrase from BACKUP_PASSPHRASE env var.
func BackupPassphrase() string {
	return os.Getenv("BACKUP_PASSPHRASE")
}

// BackupRepoPath returns the backup repo path from BACKUP_REPO_PATH env var.
func BackupRepoPath() string {
	if v := os.Getenv("BACKUP_REPO_PATH"); v != "" {
		return v
	}
	return "backup/kiroku-backup-storage"
}

// BackupGitToken returns the backup git token.
func BackupGitToken() string {
	return GetEnv("KIROKU_BACKUP_GIT_TOKEN", GetEnv("BACKUP_GIT_TOKEN", ""))
}

// BackupProviders returns the backup providers from BACKUP_PROVIDERS env var.
func BackupProviders() string {
	return os.Getenv("BACKUP_PROVIDERS")
}
