package domain

import "context"

type ReceiptRepository interface {
	CreateJob(ctx context.Context, job *ReceiptJob) error
	GetJob(ctx context.Context, id, userID string) (*ReceiptJob, error)
	UpdateJob(ctx context.Context, job *ReceiptJob) error
	CreateEvent(ctx context.Context, event *ReceiptJobEvent) error
	GetCategories(ctx context.Context) ([]Category, error)
	GetAccounts(ctx context.Context) ([]Account, error)
}

type TransactionCreationForReceipt interface {
	Create(ctx context.Context, input CreateTransactionInput) (*Transaction, error)
}
