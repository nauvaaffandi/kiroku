package domain

import "time"

type Backup struct {
	ID             string    `json:"id"`
	Filename       string    `json:"filename"`
	Size           int64     `json:"size"`
	RetentionTier  string    `json:"retention_tier"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateBackupInput struct {
	UserID string
}
