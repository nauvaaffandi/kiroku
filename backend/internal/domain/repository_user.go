package domain

import (
	"context"
	"time"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	UpdatePinHash(ctx context.Context, id, pinHash string) error
}

type PinChangeTokenRepository interface {
	Create(ctx context.Context, token, userID, purpose string, expiresAt time.Time) error
	Get(ctx context.Context, token string) (userID string, purpose string, err error)
	Delete(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
}

type TelegramBindingRepository interface {
	Get(ctx context.Context, chatID int64) (*TelegramBinding, error)
	Upsert(ctx context.Context, binding TelegramBinding) error
}
