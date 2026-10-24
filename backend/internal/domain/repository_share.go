package domain

import "context"

type ShareRepository interface {
	Create(ctx context.Context, input CreateShareInput) (*Share, error)
	List(ctx context.Context, accountID string) ([]Share, error)
	GetByToken(ctx context.Context, token string) (*Share, error)
	Delete(ctx context.Context, id string) error
	TouchAccessed(ctx context.Context, id string) error
}
