package domain

import "context"

type AccountRepository interface {
	List(ctx context.Context) ([]AccountWithBalance, error)
	Get(ctx context.Context, id string) (*AccountWithBalance, error)
	Create(ctx context.Context, input CreateAccountInput) (*Account, error)
	Update(ctx context.Context, id string, input UpdateAccountInput) (*Account, error)
	Delete(ctx context.Context, id string) error
	Summary(ctx context.Context) ([]CurrencyBalanceSummary, error)
	CountTransactions(ctx context.Context, accountID string) (int, error)
}

type CurrencyBalanceSummary struct {
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
}
