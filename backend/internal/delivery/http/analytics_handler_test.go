package httpdelivery

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"kiroku/backend/internal/domain"
)

type mockAnalyticsRepo struct {
	calendarFunc         func(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CalendarDay, error)
	balanceTrendFunc     func(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.BalanceTrendPoint, error)
	balanceAccountsFunc  func(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.AccountBalance, error)
	balanceCurrenciesFunc func(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CurrencyBalanceSummary, error)
}

func (m *mockAnalyticsRepo) Calendar(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CalendarDay, error) {
	if m.calendarFunc != nil {
		return m.calendarFunc(ctx, filter)
	}
	return nil, nil
}

func (m *mockAnalyticsRepo) CashFlowProjection(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.ProjectionPoint, []domain.ScheduledItem, error) {
	return nil, nil, nil
}
func (m *mockAnalyticsRepo) BalanceTrend(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.BalanceTrendPoint, error) {
	return nil, nil
}
func (m *mockAnalyticsRepo) BalanceAccounts(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.AccountBalance, error) {
	return nil, nil
}
func (m *mockAnalyticsRepo) BalanceCurrencies(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CurrencyBalanceSummary, error) {
	return nil, nil
}
func (m *mockAnalyticsRepo) OutlookCategory(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CategoryStats, error) {
	return nil, nil
}
func (m *mockAnalyticsRepo) OutlookTimeline(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.TimelinePoint, error) {
	return nil, nil
}
func (m *mockAnalyticsRepo) OutlookForecast(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.ForecastPoint, error) {
	return nil, nil
}
func (m *mockAnalyticsRepo) BudgetSummary(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.BudgetItem, error) {
	return nil, nil
}
func (m *mockAnalyticsRepo) BudgetVsActual(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CategoryPlanned, error) {
	return nil, nil
}
func (m *mockAnalyticsRepo) GoalsProgress(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.GoalItem, error) {
	return nil, nil
}
func (m *mockAnalyticsRepo) GoalsProjection(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.ProjectionPoint, error) {
	return nil, nil
}
func (m *mockAnalyticsRepo) StatsSummary(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CurrencyStats, error) {
	return nil, nil
}
func (m *mockAnalyticsRepo) StatsByCategory(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CategoryStats, error) {
	return nil, nil
}
func (m *mockAnalyticsRepo) StatsByCurrency(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CurrencyStats, error) {
	return nil, nil
}
func (m *mockAnalyticsRepo) StatsByAccount(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.AccountStats, error) {
	return nil, nil
}
func (m *mockAnalyticsRepo) StatsTopMerchants(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.MerchantStats, error) {
	return nil, nil
}
func (m *mockAnalyticsRepo) StatsCashflow(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CashFlowSeriesPoint, error) {
	return nil, nil
}
func (m *mockAnalyticsRepo) StatsNetWorth(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.NetWorthAccount, error) {
	return nil, nil
}

func TestAnalyticsHandler_Calendar(t *testing.T) {
	repo := &mockAnalyticsRepo{
		calendarFunc: func(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CalendarDay, error) {
			if filter.DateFrom == nil || filter.DateTo == nil {
				t.Fatalf("expected date_from and date_to to be set")
			}
			return []domain.CalendarDay{
				{
					Date:    "2026-10-01",
					Income:  100,
					Expense: 20,
					Net:     80,
					Transactions: []domain.CalendarTxn{
						{
							ID:       "txn-1",
							Name:     "Salary",
							Type:     "income",
							Amount:   100,
							Category: "Income",
							Currency: "USD",
						},
						{
							ID:       "txn-2",
							Name:     "Lunch",
							Type:     "expense",
							Amount:   20,
							Category: "Food",
							Currency: "USD",
						},
					},
				},
			}, nil
		},
	}

	handler := NewAnalyticsHandler(repo)
	router := handler.Routes()

	req := httptest.NewRequest("GET", "/calendar?year=2026&month=10", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var res struct {
		Days []domain.CalendarDayBackend `json:"days"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(res.Days) != 1 {
		t.Fatalf("expected 1 day, got %d", len(res.Days))
	}
	day := res.Days[0]
	if day.Date != "2026-10-01" {
		t.Errorf("expected date 2026-10-01, got %s", day.Date)
	}
	if day.BalanceDelta != 80 {
		t.Errorf("expected delta 80, got %f", day.BalanceDelta)
	}
	if day.ProjectedBalance != 80 {
		t.Errorf("expected projected 80, got %f", day.ProjectedBalance)
	}
	if len(day.Transactions) != 2 {
		t.Errorf("expected 2 transactions, got %d", len(day.Transactions))
	}
	if day.Transactions[0].Name != "Salary" || day.Transactions[0].Category != "Income" || day.Transactions[0].Currency != "USD" {
		t.Errorf("unexpected transaction fields: %+v", day.Transactions[0])
	}
}

func TestAnalyticsHandler_Calendar_Validation(t *testing.T) {
	repo := &mockAnalyticsRepo{}
	handler := NewAnalyticsHandler(repo)
	router := handler.Routes()

	tests := []struct {
		name       string
		query      string
		wantStatus int
	}{
		{"invalid year string", "?year=abc&month=10", http.StatusBadRequest},
		{"negative year", "?year=-2026&month=10", http.StatusBadRequest},
		{"zero year", "?year=0&month=10", http.StatusBadRequest},
		{"invalid month string", "?year=2026&month=foo", http.StatusBadRequest},
		{"month zero", "?year=2026&month=0", http.StatusBadRequest},
		{"month 13", "?year=2026&month=13", http.StatusBadRequest},
		{"empty query params ok", "", http.StatusOK},
		{"valid query params ok", "?year=2026&month=10", http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/calendar"+tc.query, nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if rec.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, rec.Code)
			}
		})
	}
}

func TestAnalyticsHandler_Calendar_RunningBalanceOrder(t *testing.T) {
	repo := &mockAnalyticsRepo{
		calendarFunc: func(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.CalendarDay, error) {
			// Even if days arrive in arbitrary or non-chronological order,
			// or sorted order, handler processes them in order.
			return []domain.CalendarDay{
				{Date: "2026-10-01", Net: 100},
				{Date: "2026-10-02", Net: -30},
				{Date: "2026-10-03", Net: 50},
			}, nil
		},
	}

	handler := NewAnalyticsHandler(repo)
	router := handler.Routes()

	req := httptest.NewRequest("GET", "/calendar?year=2026&month=10", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var res struct {
		Days []domain.CalendarDayBackend `json:"days"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(res.Days) != 3 {
		t.Fatalf("expected 3 days, got %d", len(res.Days))
	}
	expectedProjected := []float64{100, 70, 120}
	for i, d := range res.Days {
		if d.ProjectedBalance != expectedProjected[i] {
			t.Errorf("day %d (%s): expected projected balance %f, got %f", i, d.Date, expectedProjected[i], d.ProjectedBalance)
		}
	}
}
