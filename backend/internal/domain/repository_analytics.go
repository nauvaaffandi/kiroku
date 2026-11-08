package domain

import "context"

type AnalyticsRepository interface {
	CashFlowProjection(ctx context.Context, filter AnalyticsFilter) ([]ProjectionPoint, []ScheduledItem, error)
	Calendar(ctx context.Context, filter AnalyticsFilter) ([]CalendarDay, error)
	BalanceTrend(ctx context.Context, filter AnalyticsFilter) ([]BalanceTrendPoint, error)
	BalanceAccounts(ctx context.Context, filter AnalyticsFilter) ([]AccountBalance, error)
	BalanceCurrencies(ctx context.Context, filter AnalyticsFilter) ([]CurrencyBalanceSummary, error)
	OutlookCategory(ctx context.Context, filter AnalyticsFilter) ([]CategoryStats, error)
	OutlookTimeline(ctx context.Context, filter AnalyticsFilter) ([]TimelinePoint, error)
	OutlookForecast(ctx context.Context, filter AnalyticsFilter) ([]ForecastPoint, error)
	StatsSummary(ctx context.Context, filter AnalyticsFilter) ([]CurrencyStats, error)
	StatsByCategory(ctx context.Context, filter AnalyticsFilter) ([]CategoryStats, error)
	StatsByCurrency(ctx context.Context, filter AnalyticsFilter) ([]CurrencyStats, error)
	StatsByAccount(ctx context.Context, filter AnalyticsFilter) ([]AccountStats, error)
	StatsTopMerchants(ctx context.Context, filter AnalyticsFilter) ([]MerchantStats, error)
	StatsCashflow(ctx context.Context, filter AnalyticsFilter) ([]CashFlowSeriesPoint, error)
	StatsNetWorth(ctx context.Context, filter AnalyticsFilter) ([]NetWorthAccount, error)
}
