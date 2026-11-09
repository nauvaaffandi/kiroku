package budgets

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kiroku/backend/internal/db"
	"kiroku/backend/internal/domain"
	"kiroku/backend/internal/notification"
	"kiroku/backend/internal/repository/sqlite"
)

type mockBudgetRepo struct {
	items    []domain.Budget
	users    []string
	calcs    map[string]*domain.BudgetResult
	listErr  map[string]error
	calcErr  error
	usersErr error
}

func (m *mockBudgetRepo) List(ctx context.Context, userID string) ([]domain.Budget, error) {
	if err := m.listErr[userID]; err != nil {
		return nil, err
	}
	var res []domain.Budget
	for _, b := range m.items {
		if b.UserID == userID {
			res = append(res, b)
		}
	}
	return res, nil
}

func (m *mockBudgetRepo) Get(context.Context, string, string) (*domain.Budget, error) {
	return nil, nil
}

func (m *mockBudgetRepo) Create(context.Context, domain.CreateBudgetInput, string) (*domain.Budget, error) {
	return nil, nil
}

func (m *mockBudgetRepo) Update(context.Context, string, string, domain.UpdateBudgetInput) (*domain.Budget, error) {
	return nil, nil
}

func (m *mockBudgetRepo) Delete(context.Context, string, string) error { return nil }

func (m *mockBudgetRepo) Calculate(ctx context.Context, budget *domain.Budget) (*domain.BudgetResult, error) {
	if m.calcErr != nil {
		return nil, m.calcErr
	}
	return m.calcs[budget.ID], nil
}

func (m *mockBudgetRepo) UserIDsWithBudgets(context.Context) ([]string, error) {
	return m.users, m.usersErr
}

type mockNotifier struct {
	mu        sync.Mutex
	created   []domain.CreateNotificationInput
	createErr error
}

func (m *mockNotifier) Create(_ context.Context, input domain.CreateNotificationInput) (*domain.Notification, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.createErr != nil {
		return nil, m.createErr
	}
	m.created = append(m.created, input)
	return &domain.Notification{ID: "notif-1"}, nil
}

func (m *mockNotifier) Created() []domain.CreateNotificationInput {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]domain.CreateNotificationInput(nil), m.created...)
}

func TestBudgetService_ThresholdEdgesAndFlags(t *testing.T) {
	cases := []struct {
		name        string
		percentage  int
		approaching bool
		overEnabled bool
		wantEvents  []string
	}{
		{name: "below both", percentage: 79, approaching: true, overEnabled: true},
		{name: "at approaching", percentage: 80, approaching: true, overEnabled: true, wantEvents: []string{domain.BudgetPeriodApproachingEvent}},
		{name: "below over", percentage: 99, approaching: true, overEnabled: true, wantEvents: []string{domain.BudgetPeriodApproachingEvent}},
		{name: "at over both", percentage: 100, approaching: true, overEnabled: true, wantEvents: []string{domain.BudgetPeriodApproachingEvent, domain.BudgetPeriodOverEvent}},
		{name: "direct jump both", percentage: 150, approaching: true, overEnabled: true, wantEvents: []string{domain.BudgetPeriodApproachingEvent, domain.BudgetPeriodOverEvent}},
		{name: "approaching disabled", percentage: 90, approaching: false, overEnabled: true},
		{name: "over disabled", percentage: 110, approaching: true, overEnabled: false, wantEvents: []string{domain.BudgetPeriodApproachingEvent}},
		{name: "both disabled", percentage: 150},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			budget := domain.Budget{ID: "b1", UserID: "u1", TrendingOver: tc.approaching, OverBudget: tc.overEnabled}
			repo := &mockBudgetRepo{
				items: []domain.Budget{budget},
				calcs: map[string]*domain.BudgetResult{"b1": {
					BudgetID: "b1", Name: "Budget", Period: "monthly", Amount: "100.00", Spent: "100.00", Currency: "USD",
					PercentageUsed: tc.percentage, PeriodStart: time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
				}},
			}
			notifier := &mockNotifier{}
			require.NoError(t, NewService(repo, notifier).Check(context.Background(), "u1"))
			var events []string
			for _, n := range notifier.Created() {
				events = append(events, n.Type)
			}
			assert.ElementsMatch(t, tc.wantEvents, events)
		})
	}
}

func TestBudgetService_CheckErrorsPropagate(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		repo := &mockBudgetRepo{listErr: map[string]error{"u1": errors.New("list failed")}}
		err := NewService(repo, &mockNotifier{}).Check(context.Background(), "u1")
		require.ErrorContains(t, err, "list failed")
	})
	t.Run("calculate", func(t *testing.T) {
		repo := &mockBudgetRepo{items: []domain.Budget{{ID: "b1", UserID: "u1"}}, calcErr: errors.New("calculate failed")}
		err := NewService(repo, &mockNotifier{}).Check(context.Background(), "u1")
		require.ErrorContains(t, err, "calculate failed")
	})
	t.Run("notify", func(t *testing.T) {
		repo := &mockBudgetRepo{
			items: []domain.Budget{{ID: "b1", UserID: "u1", TrendingOver: true}},
			calcs: map[string]*domain.BudgetResult{"b1": {PercentageUsed: 80, PeriodStart: time.Now()}},
		}
		err := NewService(repo, &mockNotifier{createErr: errors.New("notify failed")}).Check(context.Background(), "u1")
		require.ErrorContains(t, err, "notify failed")
	})
}

func TestBudgetService_EmptyUserDoesNothing(t *testing.T) {
	notifier := &mockNotifier{}
	require.NoError(t, NewService(&mockBudgetRepo{}, notifier).Check(context.Background(), ""))
	assert.Empty(t, notifier.Created())
}

func TestBudgetService_ReconcileContinuesAfterUserFailure(t *testing.T) {
	repo := &mockBudgetRepo{
		users:   []string{"u1", "u2"},
		listErr: map[string]error{"u1": errors.New("u1 unavailable")},
		items:   []domain.Budget{{ID: "b2", UserID: "u2", TrendingOver: true}},
		calcs:   map[string]*domain.BudgetResult{"b2": {PercentageUsed: 80, PeriodStart: time.Now()}},
	}
	notifier := &mockNotifier{}
	require.NoError(t, NewService(repo, notifier).ReconcileAll(context.Background()))
	require.Len(t, notifier.Created(), 1)
	assert.Equal(t, "u2", notifier.Created()[0].RecipientID)
}

func newBudgetNotificationFixture(t *testing.T) (*db.DB, *sqlite.BudgetsRepo, *notification.Service) {
	t.Helper()
	database, err := db.Open(t.TempDir() + "/data")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close()) })

	for _, id := range []string{"u1", "u2"} {
		_, err = database.ExecContext(context.Background(), `INSERT INTO users (id, email, pin_hash) VALUES (?, ?, 'hash')`, id, id+"@example.test")
		require.NoError(t, err)
	}
	_, err = database.ExecContext(context.Background(), `INSERT INTO accounts (id, name, type, currency) VALUES ('acc1', 'Cash', 'cash', 'USD')`)
	require.NoError(t, err)

	budgetRepo := sqlite.NewBudgetsRepo(database)
	notifRepo := sqlite.NewNotificationRepo(database)
	deliveryRepo := sqlite.NewDeliveryRepo(database)
	subRepo := sqlite.NewPushSubscriptionRepo(database)
	notifService := notification.NewService(notifRepo, deliveryRepo, subRepo, nil, notification.ServiceConfig{})
	return database, budgetRepo, notifService
}

func TestBudgetService_ReconciliationRecoversAndDeduplicates(t *testing.T) {
	database, repo, notifications := newBudgetNotificationFixture(t)
	ctx := context.Background()
	budget, err := repo.Create(ctx, domain.CreateBudgetInput{
		Name: "Monthly", Period: domain.BudgetPeriodMonthly, Amount: 100, Currency: "USD",
		TrendingOver: true, OverBudget: true,
	}, "u1")
	require.NoError(t, err)
	_, err = database.ExecContext(ctx, `INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id) VALUES (?, 'expense', ?, 12000, 'USD', 'acc1')`, uuid.NewString(), time.Now())
	require.NoError(t, err)

	service := NewService(repo, notifications)
	require.NoError(t, service.ReconcileAll(ctx))
	require.NoError(t, service.ReconcileAll(ctx))
	count, err := notifications.Count(ctx, "u1", false)
	require.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.NotEmpty(t, budget.ID)
}

func TestBudgetService_ConcurrentReconciliationIsIdempotent(t *testing.T) {
	database, repo, notifications := newBudgetNotificationFixture(t)
	ctx := context.Background()
	_, err := repo.Create(ctx, domain.CreateBudgetInput{
		Name: "Monthly", Period: domain.BudgetPeriodMonthly, Amount: 100, Currency: "USD",
		TrendingOver: true, OverBudget: true,
	}, "u1")
	require.NoError(t, err)
	_, err = database.ExecContext(ctx, `INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id) VALUES (?, 'expense', ?, 12000, 'USD', 'acc1')`, uuid.NewString(), time.Now())
	require.NoError(t, err)

	service := NewService(repo, notifications)
	var wg sync.WaitGroup
	errs := make(chan error, 8)
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs <- service.ReconcileAll(ctx)
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}
	count, err := notifications.Count(ctx, "u1", false)
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}
