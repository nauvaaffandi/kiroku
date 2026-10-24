package domain

import (
	"context"
	"time"
)

type ImportLogRepository interface {
	StorePreview(ctx context.Context, previewID string, rows []ImportRow, expiresAt time.Time) error
	GetPreview(ctx context.Context, previewID string) ([]ImportRow, error)
	DeletePreview(ctx context.Context, previewID string) error
}

type AccountLookupForImport interface {
	List(ctx context.Context) ([]AccountWithBalance, error)
}

type CategoryLookupForImport interface {
	List(ctx context.Context, includeArchived bool) ([]Category, error)
}
