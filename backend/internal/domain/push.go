package domain

import (
	"context"
	"time"
)

type PushSubscription struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Endpoint     string    `json:"endpoint"`
	P256DH       string    `json:"p256dh"`
	Auth         string    `json:"auth"`
	UserAgent    string    `json:"user_agent"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreatePushSubscriptionInput struct {
	UserID    string `json:"user_id"`
	Endpoint  string `json:"endpoint"`
	Keys      Keys   `json:"keys"`
	UserAgent string `json:"user_agent"`
}

type Keys struct {
	P256DH string `json:"p256dh"`
	Auth   string `json:"auth"`
}

type PushSubscriptionRepository interface {
	Create(ctx context.Context, input CreatePushSubscriptionInput) (*PushSubscription, error)
	Get(ctx context.Context, id string) (*PushSubscription, error)
	GetByEndpoint(ctx context.Context, endpoint, userID string) (*PushSubscription, error)
	ListByUser(ctx context.Context, userID string) ([]*PushSubscription, error)
	Delete(ctx context.Context, id string) error
	DeleteByEndpoint(ctx context.Context, endpoint, userID string) error
	DeleteByUser(ctx context.Context, userID string) error
	Update(ctx context.Context, id string, update UpdatePushSubscriptionInput) (*PushSubscription, error)
}

type UpdatePushSubscriptionInput struct {
	UserAgent *string    `json:"user_agent,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
