package domain

import "context"

type TransactionRepository interface {
	List(ctx context.Context, filter TransactionFilter) ([]Transaction, error)
	Count(ctx context.Context, filter TransactionFilter) (int, error)
	Get(ctx context.Context, id string) (*Transaction, error)
	Create(ctx context.Context, input CreateTransactionInput) (*Transaction, error)
	Update(ctx context.Context, id string, input UpdateTransactionInput) (*Transaction, error)
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, query string, limit int) ([]Transaction, error)
	Autocomplete(ctx context.Context, query string, limit int) ([]string, error)
	Summary(ctx context.Context, filter TransactionFilter) ([]SummaryItem, error)
	Batch(ctx context.Context, input BatchInput) (int, error)
}
