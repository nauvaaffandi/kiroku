package domain

import "context"

type ScheduledRepository interface {
	List(ctx context.Context) ([]ScheduledTransaction, error)
	Get(ctx context.Context, id string) (*ScheduledTransaction, error)
	Create(ctx context.Context, input CreateScheduledInput) (*ScheduledTransaction, error)
	Update(ctx context.Context, id string, input UpdateScheduledInput) (*ScheduledTransaction, error)
	Delete(ctx context.Context, id string) error
	Toggle(ctx context.Context, id string, enabled bool) (*ScheduledTransaction, error)
}
