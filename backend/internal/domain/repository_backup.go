package domain

import "context"

type BackupRepository interface {
	List(ctx context.Context) ([]Backup, error)
	GetByID(ctx context.Context, id string) (*Backup, error)
	Create(ctx context.Context, backup *Backup) error
	Delete(ctx context.Context, id string) error
	GetState(ctx context.Context, userID string) (*BackupState, error)
	SaveState(ctx context.Context, state *BackupState) error
	GetBackupDates(ctx context.Context, userID string) ([]string, error)
	GetBackupMonths(ctx context.Context, userID string) ([]string, error)
	GetRechecks(ctx context.Context, userID string) ([]string, error)
}

type BackupState struct {
	UserID        string `json:"user_id"`
	LastDailyRun  *string `json:"last_daily_run,omitempty"`
	LastWeeklyRun *string `json:"last_weekly_run,omitempty"`
	LastMonthlyRun *string `json:"last_monthly_run,omitempty"`
}
