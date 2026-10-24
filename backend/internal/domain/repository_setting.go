package domain

import "context"

type SettingsRepository interface {
	GetAll(ctx context.Context) (map[string]string, error)
	GetOne(ctx context.Context, key string) (string, error)
	SetOne(ctx context.Context, key, value string) error
	UpdateAll(ctx context.Context, settings map[string]string) error
}
