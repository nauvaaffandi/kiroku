package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kiroku/backend/internal/db"
	"kiroku/backend/internal/domain"
)

func newBudgetTestEnv(t *testing.T) (*db.DB, *BudgetsRepo) {
	t.Helper()
	dataDir := t.TempDir() + "/data"
	d, err := db.Open(dataDir)
	require.NoError(t, err)

	_, err = d.ExecContext(context.Background(), `INSERT INTO users (id, email, pin_hash) VALUES (?, ?, ?)`, "u1", "u1@example.com", "hash")
	require.NoError(t, err)
	_, err = d.ExecContext(context.Background(), `INSERT INTO users (id, email, pin_hash) VALUES (?, ?, ?)`, "u2", "u2@example.com", "hash")
	require.NoError(t, err)

	_, err = d.ExecContext(context.Background(), `INSERT INTO accounts (id, name, type, currency) VALUES ('acc1', 'Acc 1', 'bank', 'USD')`)
	require.NoError(t, err)
	_, err = d.ExecContext(context.Background(), `INSERT INTO accounts (id, name, type, currency) VALUES ('acc2', 'Acc 2', 'bank', 'USD')`)
	require.NoError(t, err)

	_, err = d.ExecContext(context.Background(), `INSERT INTO categories (id, name, kind) VALUES ('cat1', 'Food', 'must')`)
	require.NoError(t, err)
	_, err = d.ExecContext(context.Background(), `INSERT INTO categories (id, name, kind) VALUES ('cat2', 'Transport', 'must')`)
	require.NoError(t, err)

	_, err = d.ExecContext(context.Background(), `INSERT INTO labels (id, name) VALUES ('lbl1', 'Office')`)
	require.NoError(t, err)

	return d, NewBudgetsRepo(d)
}

func TestBudgetUpdateIDOR(t *testing.T) {
	d, repo := newBudgetTestEnv(t)
	defer d.Close()

	ctx := context.Background()

	b, err := repo.Create(ctx, domain.CreateBudgetInput{
		Name:         "User 1 Budget",
		Period:       domain.BudgetPeriodMonthly,
		Amount:       100.0,
		Currency:     "USD",
		CategoryIDs:  []string{"cat1"},
		LabelIDs:     []string{"lbl1"},
		TrendingOver: true,
		OverBudget:   true,
	}, "u1")
	require.NoError(t, err)
	require.NotNil(t, b)

	// User 2 tries to update User 1's budget
	newName := "Hacked Name"
	newAmount := 999.0
	res, err := repo.Update(ctx, b.ID, "u2", domain.UpdateBudgetInput{
		Name:        &newName,
		Amount:      &newAmount,
		CategoryIDs: []string{"cat2"},
	})
	assert.Error(t, err)
	assert.Nil(t, res)

	// Verify original budget and categories untouched
	original, err := repo.Get(ctx, b.ID, "u1")
	require.NoError(t, err)
	require.NotNil(t, original)
	assert.Equal(t, "User 1 Budget", original.Name)
	assert.Equal(t, "100.00", original.AmountStr)
	var categoryID string
	err = d.QueryRowContext(ctx, `SELECT category_id FROM budget_categories WHERE budget_id = ?`, b.ID).Scan(&categoryID)
	require.NoError(t, err)
	assert.Equal(t, "cat1", categoryID)

	var labelID string
	err = d.QueryRowContext(ctx, `SELECT label_id FROM budget_labels WHERE budget_id = ?`, b.ID).Scan(&labelID)
	require.NoError(t, err)
	assert.Equal(t, "lbl1", labelID)
}

func TestBudgetUpdateCategoriesAndLabelsScopedToOwner(t *testing.T) {
	d, repo := newBudgetTestEnv(t)
	defer d.Close()

	ctx := context.Background()

	b1, err := repo.Create(ctx, domain.CreateBudgetInput{
		Name:         "User 1 Budget",
		Period:       domain.BudgetPeriodMonthly,
		Amount:       100.0,
		Currency:     "USD",
		CategoryIDs:  []string{"cat1"},
		LabelIDs:     []string{"lbl1"},
		TrendingOver: true,
		OverBudget:   true,
	}, "u1")
	require.NoError(t, err)

	b2, err := repo.Create(ctx, domain.CreateBudgetInput{
		Name:        "User 2 Budget",
		Period:      domain.BudgetPeriodMonthly,
		Amount:      50.0,
		Currency:    "USD",
		CategoryIDs: []string{"cat2"},
	}, "u2")
	require.NoError(t, err)

	// User 2 replaces User 1's categories/labels - must fail entirely
	res, err := repo.Update(ctx, b1.ID, "u2", domain.UpdateBudgetInput{
		CategoryIDs: []string{"cat2"},
		LabelIDs:    []string{"lbl1"},
	})
	assert.Error(t, err)
	assert.Nil(t, res)

	// User 1 budget still has cat1/lbl1
	var catID string
	err = d.QueryRowContext(ctx, `SELECT category_id FROM budget_categories WHERE budget_id = ?`, b1.ID).Scan(&catID)
	require.NoError(t, err)
	assert.Equal(t, "cat1", catID)

	var labelID string
	err = d.QueryRowContext(ctx, `SELECT label_id FROM budget_labels WHERE budget_id = ?`, b1.ID).Scan(&labelID)
	require.NoError(t, err)
	assert.Equal(t, "lbl1", labelID)

	// User 2 budget still has cat2
	err = d.QueryRowContext(ctx, `SELECT category_id FROM budget_categories WHERE budget_id = ?`, b2.ID).Scan(&catID)
	require.NoError(t, err)
	assert.Equal(t, "cat2", catID)
}

func TestBudgetDeleteIDOR(t *testing.T) {
	d, repo := newBudgetTestEnv(t)
	defer d.Close()

	ctx := context.Background()

	b, err := repo.Create(ctx, domain.CreateBudgetInput{
		Name:     "User 1 Budget",
		Period:   domain.BudgetPeriodMonthly,
		Amount:   100.0,
		Currency: "USD",
	}, "u1")
	require.NoError(t, err)

	// User 2 cannot delete User 1's budget
	err = repo.Delete(ctx, b.ID, "u2")
	assert.NoError(t, err)

	// Budget still exists for u1
	orig, err := repo.Get(ctx, b.ID, "u1")
	require.NoError(t, err)
	require.NotNil(t, orig)
	assert.False(t, orig.DeletedAt != nil && !orig.DeletedAt.IsZero())
}

func TestBudgetCalculate_AccountMatchingOnlyFromAccount(t *testing.T) {
	d, repo := newBudgetTestEnv(t)
	defer d.Close()

	ctx := context.Background()
	acc1 := "acc1"

	b, err := repo.Create(ctx, domain.CreateBudgetInput{
		Name:      "Account-specific Budget",
		Period:    domain.BudgetPeriodMonthly,
		Amount:    200.0,
		Currency:  "USD",
		AccountID: &acc1,
	}, "u1")
	require.NoError(t, err)

	now := time.Now()
	// Insert expense transaction where acc1 is from_account_id (should be counted)
	_, err = d.ExecContext(ctx, `
		INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id)
		VALUES (?, 'expense', ?, 5000, 'USD', 'acc1')
	`, uuid.New().String(), now)
	require.NoError(t, err)

	// Insert expense transaction where acc1 is to_account_id (should NOT be counted)
	_, err = d.ExecContext(ctx, `
		INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id, to_account_id)
		VALUES (?, 'expense', ?, 8000, 'USD', 'acc2', 'acc1')
	`, uuid.New().String(), now)
	require.NoError(t, err)

	calc, err := repo.Calculate(ctx, b)
	require.NoError(t, err)
	assert.Equal(t, "50.00", calc.Spent)
	assert.Equal(t, "150.00", calc.Remaining)
	assert.Equal(t, 25, calc.PercentageUsed)
}

func TestBudgetCalculate_ExcludesDeletedTransactions(t *testing.T) {
	d, repo := newBudgetTestEnv(t)
	defer d.Close()

	ctx := context.Background()

	b, err := repo.Create(ctx, domain.CreateBudgetInput{
		Name:      "Budget",
		Period:    domain.BudgetPeriodMonthly,
		Amount:    1000.0,
		Currency:  "USD",
		AccountID: strPtr("acc1"),
	}, "u1")
	require.NoError(t, err)

	now := time.Now()
	// Active expense
	_, err = d.ExecContext(ctx, `
		INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id)
		VALUES (?, 'expense', ?, 3000, 'USD', 'acc1')
	`, uuid.New().String(), now)
	require.NoError(t, err)

	// Soft-deleted expense
	_, err = d.ExecContext(ctx, `
		INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id, deleted_at)
		VALUES (?, 'expense', ?, 7000, 'USD', 'acc1', datetime('now'))
	`, uuid.New().String(), now)
	require.NoError(t, err)

	calc, err := repo.Calculate(ctx, b)
	require.NoError(t, err)
	assert.Equal(t, "30.00", calc.Spent)
	assert.Equal(t, 3, calc.PercentageUsed)
}

func TestBudgetCalculate_CategoryAndLabelFilters(t *testing.T) {
	d, repo := newBudgetTestEnv(t)
	defer d.Close()

	ctx := context.Background()

	b, err := repo.Create(ctx, domain.CreateBudgetInput{
		Name:        "Category Budget",
		Period:      domain.BudgetPeriodMonthly,
		Amount:      100.0,
		Currency:    "USD",
		CategoryIDs: []string{"cat1"},
		LabelIDs:    []string{"lbl1"},
	}, "u1")
	require.NoError(t, err)

	now := time.Now()
	// Txn with cat1 + lbl1 (should count)
	t1 := uuid.New().String()
	_, err = d.ExecContext(ctx, `
		INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id, category_id)
		VALUES (?, 'expense', ?, 4000, 'USD', 'acc1', 'cat1')
	`, t1, now)
	require.NoError(t, err)
	_, err = d.ExecContext(ctx, `INSERT INTO transaction_labels (transaction_id, label_id) VALUES (?, 'lbl1')`, t1)
	require.NoError(t, err)

	// Txn with cat2 only (should not count)
	_, err = d.ExecContext(ctx, `
		INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id, category_id)
		VALUES (?, 'expense', ?, 6000, 'USD', 'acc1', 'cat2')
	`, uuid.New().String(), now)
	require.NoError(t, err)

	// Txn with cat1 but no lbl1 (should not count - label filter required)
	t3 := uuid.New().String()
	_, err = d.ExecContext(ctx, `
		INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id, category_id)
		VALUES (?, 'expense', ?, 2000, 'USD', 'acc1', 'cat1')
	`, t3, now)
	require.NoError(t, err)

	calc, err := repo.Calculate(ctx, b)
	require.NoError(t, err)
	assert.Equal(t, "40.00", calc.Spent)
}

func TestBudgetCalculate_FXPerTransactionDate(t *testing.T) {
	d, repo := newBudgetTestEnv(t)
	defer d.Close()

	ctx := context.Background()

	b, err := repo.Create(ctx, domain.CreateBudgetInput{
		Name:     "EUR Budget",
		Period:   domain.BudgetPeriodMonthly,
		Amount:   500.0,
		Currency: "EUR",
	}, "u1")
	require.NoError(t, err)

	// Add exchange rates for two different dates
	_, err = d.ExecContext(ctx, `
		INSERT INTO fx_rates (base, quote, rate, date)
		VALUES ('USD', 'EUR', 0.80, '2026-11-01'),
		       ('USD', 'EUR', 0.90, '2026-11-05')
	`)
	require.NoError(t, err)

	t1, err := time.Parse("2006-01-02 15:04:05", "2026-11-01 12:00:00")
	require.NoError(t, err)
	t2, err := time.Parse("2006-01-02 15:04:05", "2026-11-05 12:00:00")
	require.NoError(t, err)

	// Txn 1: 100 USD on Nov 1 -> 100 * 0.80 = 80 EUR = 8000 cents
	_, err = d.ExecContext(ctx, `
		INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id)
		VALUES (?, 'expense', ?, 10000, 'USD', 'acc1')
	`, uuid.New().String(), t1)
	require.NoError(t, err)

	// Txn 2: 100 USD on Nov 5 -> 100 * 0.90 = 90 EUR = 9000 cents
	_, err = d.ExecContext(ctx, `
		INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id)
		VALUES (?, 'expense', ?, 10000, 'USD', 'acc1')
	`, uuid.New().String(), t2)
	require.NoError(t, err)

	// Custom period window so calculation matches transactions
	b.Period = domain.BudgetPeriodOneTime
	start := t1.Add(-1 * time.Hour)
	end := t2.Add(1 * time.Hour)
	b.StartsAt = &start
	b.EndsAt = &end

	calc, err := repo.Calculate(ctx, b)
	require.NoError(t, err)

	// Total spent = 80.00 + 90.00 = 170.00 EUR
	assert.Equal(t, "170.00", calc.Spent)
}

func TestBudgetCalculate_MissingFXRateReturnsError(t *testing.T) {
	d, repo := newBudgetTestEnv(t)
	defer d.Close()

	ctx := context.Background()

	b, err := repo.Create(ctx, domain.CreateBudgetInput{
		Name:     "GBP Budget",
		Period:   domain.BudgetPeriodMonthly,
		Amount:   100.0,
		Currency: "GBP",
	}, "u1")
	require.NoError(t, err)

	now := time.Now()
	_, err = d.ExecContext(ctx, `
		INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id)
		VALUES (?, 'expense', ?, 10000, 'USD', 'acc1')
	`, uuid.New().String(), now)
	require.NoError(t, err)

	calc, err := repo.Calculate(ctx, b)
	assert.Error(t, err)
	assert.Nil(t, calc)
}

func TestBudgetCalculate_Rounding(t *testing.T) {
	d, repo := newBudgetTestEnv(t)
	defer d.Close()

	ctx := context.Background()

	b, err := repo.Create(ctx, domain.CreateBudgetInput{
		Name:     "Rounding Budget",
		Period:   domain.BudgetPeriodMonthly,
		Amount:   100.0,
		Currency: "USD",
	}, "u1")
	require.NoError(t, err)

	now := time.Now()
	_, err = d.ExecContext(ctx, `
		INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id)
		VALUES (?, 'expense', ?, 3334, 'USD', 'acc1')
	`, uuid.New().String(), now)
	require.NoError(t, err)

	calc, err := repo.Calculate(ctx, b)
	require.NoError(t, err)
	assert.Equal(t, "33.34", calc.Spent)
	assert.Equal(t, "66.66", calc.Remaining)
	assert.Equal(t, 33, calc.PercentageUsed)
}

func TestBudgetCalculate_EmptyPeriodReturnsZero(t *testing.T) {
	d, repo := newBudgetTestEnv(t)
	defer d.Close()

	ctx := context.Background()

	b, err := repo.Create(ctx, domain.CreateBudgetInput{
		Name:     "Empty Budget",
		Period:   domain.BudgetPeriodMonthly,
		Amount:   100.0,
		Currency: "USD",
	}, "u1")
	require.NoError(t, err)

	calc, err := repo.Calculate(ctx, b)
	require.NoError(t, err)
	assert.Equal(t, "0.00", calc.Spent)
	assert.Equal(t, "100.00", calc.Remaining)
	assert.Equal(t, 0, calc.PercentageUsed)
	assert.False(t, calc.IsOverBudget)
}

func TestBudgetCalculate_WeeklyPeriodBoundary(t *testing.T) {
	d, repo := newBudgetTestEnv(t)
	defer d.Close()

	ctx := context.Background()

	b, err := repo.Create(ctx, domain.CreateBudgetInput{
		Name:      "Weekly Budget",
		Period:    domain.BudgetPeriodWeekly,
		Amount:    1000.0,
		Currency:  "USD",
		AccountID: strPtr("acc1"),
	}, "u1")
	require.NoError(t, err)

	now := time.Now()
	// Transaction today inside weekly window
	_, err = d.ExecContext(ctx, `
		INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id)
		VALUES (?, 'expense', ?, 5000, 'USD', 'acc1')
	`, uuid.New().String(), now)
	require.NoError(t, err)

	// Transaction exactly at period end boundary (exclusive) - should NOT count
	_, periodEnd, err := periodWindowForTest(b)
	require.NoError(t, err)
	_, err = d.ExecContext(ctx, `
		INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id)
		VALUES (?, 'expense', ?, 50000, 'USD', 'acc1')
`, uuid.New().String(), periodEnd.Add(time.Hour))

	require.NoError(t, err)

	calc, err := repo.Calculate(ctx, b)
	require.NoError(t, err)
	assert.Equal(t, "50.00", calc.Spent)
}

func TestBudgetCalculate_MonthlyPeriodBoundary(t *testing.T) {
	d, repo := newBudgetTestEnv(t)
	defer d.Close()

	ctx := context.Background()

	b, err := repo.Create(ctx, domain.CreateBudgetInput{
		Name:     "Monthly Budget",
		Period:   domain.BudgetPeriodMonthly,
		Amount:   1000.0,
		Currency: "USD",
	}, "u1")
	require.NoError(t, err)

	now := time.Now()
	_, err = d.ExecContext(ctx, `
		INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id)
		VALUES (?, 'expense', ?, 5000, 'USD', 'acc1')
	`, uuid.New().String(), now)
	require.NoError(t, err)

	// Transaction exactly at period end boundary (exclusive) - should NOT count
	_, periodEnd, err := periodWindowForTest(b)
	require.NoError(t, err)
	_, err = d.ExecContext(ctx, `
		INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id)
		VALUES (?, 'expense', ?, 50000, 'USD', 'acc1')
`, uuid.New().String(), periodEnd.Add(time.Hour))

	require.NoError(t, err)

	calc, err := repo.Calculate(ctx, b)
	require.NoError(t, err)
	assert.Equal(t, "50.00", calc.Spent)
}

// periodWindowForTest is a small helper to avoid duplicating logic; it calls
// the unexported periodWindow within this package.
func periodWindowForTest(b *domain.Budget) (time.Time, time.Time, error) {
	return periodWindow(b)
}

func strPtr(s string) *string { return &s }
