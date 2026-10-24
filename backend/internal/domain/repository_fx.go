package domain

import (
	"context"
	"time"
)

type FXRepository interface {
	GetQuotes(ctx context.Context, base string, date time.Time) ([]FXQuote, error)
	Upsert(ctx context.Context, quotes []FXQuote) error
	Count(ctx context.Context) (int, error)
	GetRate(ctx context.Context, base, quote string, date time.Time) (*float64, error)
}
