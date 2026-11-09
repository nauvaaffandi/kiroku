package sqlite

import (
	"context"
	"fmt"

	"kiroku/backend/internal/domain"
)

func (r *AnalyticsRepo) CashFlowProjection(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.ProjectionPoint, []domain.ScheduledItem, error) {
	return nil, nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) OutlookCategory(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CategoryStats, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) OutlookTimeline(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.TimelinePoint, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) OutlookForecast(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.ForecastPoint, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) StatsSummary(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CurrencyStats, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) StatsByCategory(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CategoryStats, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) StatsByCurrency(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CurrencyStats, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) StatsByAccount(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.AccountStats, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) StatsTopMerchants(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.MerchantStats, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) StatsCashflow(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CashFlowSeriesPoint, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *AnalyticsRepo) StatsNetWorth(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.NetWorthAccount, error) {
	return nil, fmt.Errorf("not implemented")
}
