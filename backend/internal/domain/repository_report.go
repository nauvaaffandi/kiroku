package domain

import "context"

type SpendingRepository interface {
	SpendingByCategory(ctx context.Context, filter SpendingFilter) ([]SpendingGroup, error)
	SpendingByParent(ctx context.Context, filter SpendingFilter) ([]SpendingGroup, error)
	Trend(ctx context.Context, filter TrendFilter) ([]TrendPoint, error)
}
