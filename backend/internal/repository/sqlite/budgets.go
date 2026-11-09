package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"kiroku/backend/internal/db"
	"kiroku/backend/internal/domain"
)

// BudgetsRepo persists budgets and their category/label associations.
type BudgetsRepo struct{ db *db.DB }

func NewBudgetsRepo(database *db.DB) *BudgetsRepo { return &BudgetsRepo{db: database} }

// UserIDsWithBudgets returns distinct user_ids that have non-deleted budgets.
func (r *BudgetsRepo) UserIDsWithBudgets(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT DISTINCT user_id FROM budgets WHERE deleted_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		out = append(out, uid)
	}
	return out, rows.Err()
}

func (r *BudgetsRepo) List(ctx context.Context, userID string) ([]domain.Budget, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, name, period, amount, currency, account_id,
		       trending_over, over_budget, starts_at, ends_at,
		       deleted_at, created_at, updated_at
		FROM budgets
		WHERE user_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Budget
	for rows.Next() {
		b, err := r.scanBudget(rows)
		if err != nil {
			return nil, err
		}
		if err := r.loadRelations(ctx, &b); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	if out == nil {
		out = []domain.Budget{}
	}
	return out, rows.Err()
}

func (r *BudgetsRepo) Calculate(ctx context.Context, budget *domain.Budget) (*domain.BudgetResult, error) {
	return r.budgetCalc(ctx, budget)
}

func (r *BudgetsRepo) Get(ctx context.Context, id, userID string) (*domain.Budget, error) {
	var b domain.Budget
	var amount int64
	var startsAt, endsAt, deletedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, name, period, amount, currency, account_id,
		       trending_over, over_budget, starts_at, ends_at,
		       deleted_at, created_at, updated_at
		FROM budgets
		WHERE id = ? AND user_id = ? AND deleted_at IS NULL
	`, id, userID).Scan(
		&b.ID, &b.UserID, &b.Name, &b.Period, &amount, &b.Currency, &b.AccountID,
		&b.TrendingOver, &b.OverBudget, &startsAt, &endsAt,
		&deletedAt, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	b.Amount = amount
	b.AmountStr = domain.FormatCents(amount)
	if startsAt.Valid {
		b.StartsAt = &startsAt.Time
	}
	if endsAt.Valid {
		b.EndsAt = &endsAt.Time
	}
	if deletedAt.Valid {
		t := deletedAt.Time
		b.DeletedAt = &t
	}
	if err := r.loadRelations(ctx, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BudgetsRepo) scanBudget(scanner interface {
	Scan(...interface{}) error
}) (domain.Budget, error) {
	var b domain.Budget
	var amount int64
	var startsAt, endsAt, deletedAt sql.NullTime
	err := scanner.Scan(
		&b.ID, &b.UserID, &b.Name, &b.Period, &amount, &b.Currency, &b.AccountID,
		&b.TrendingOver, &b.OverBudget, &startsAt, &endsAt,
		&deletedAt, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return b, err
	}
	b.Amount = amount
	b.AmountStr = domain.FormatCents(amount)
	if startsAt.Valid {
		b.StartsAt = &startsAt.Time
	}
	if endsAt.Valid {
		b.EndsAt = &endsAt.Time
	}
	if deletedAt.Valid {
		t := deletedAt.Time
		b.DeletedAt = &t
	}
	return b, nil
}

func (r *BudgetsRepo) loadRelations(ctx context.Context, b *domain.Budget) error {
	cRows, err := r.db.QueryContext(ctx, `SELECT category_id FROM budget_categories WHERE budget_id = ?`, b.ID)
	if err != nil {
		return err
	}
	defer cRows.Close()
	for cRows.Next() {
		var cid string
		if err := cRows.Scan(&cid); err != nil {
			return err
		}
		cat, err := scanCategory(r.db.QueryRowContext(ctx, `SELECT id, parent_id, name, icon, color, kind, sort_order, archived FROM categories WHERE id = ?`, cid))
		if err != nil || cat == nil {
			continue
		}
		b.Categories = append(b.Categories, *cat)
	}

	lRows, err := r.db.QueryContext(ctx, `SELECT label_id FROM budget_labels WHERE budget_id = ?`, b.ID)
	if err != nil {
		return err
	}
	defer lRows.Close()
	for lRows.Next() {
		var lid string
		if err := lRows.Scan(&lid); err != nil {
			return err
		}
		label, err := scanLabel(r.db.QueryRowContext(ctx, `SELECT id, name, color, icon, description, sort_order FROM labels WHERE id = ?`, lid))
		if err != nil || label == nil {
			continue
		}
		b.Labels = append(b.Labels, *label)
	}
	return nil
}

func scanCategory(scanner interface{ Scan(...interface{}) error }) (*domain.Category, error) {
	var c domain.Category
	var parentID sql.NullString
	var kind sql.NullString
	var archived int
	if err := scanner.Scan(&c.ID, &parentID, &c.Name, &c.Icon, &c.Color, &kind, &c.SortOrder, &archived); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if parentID.Valid {
		c.ParentID = &parentID.String
	}
	if kind.Valid {
		c.Kind = &kind.String
	}
	c.Archived = archived != 0
	return &c, nil
}

func scanLabel(scanner interface{ Scan(...interface{}) error }) (*domain.Label, error) {
	var l domain.Label
	var description sql.NullString
	var icon sql.NullString
	var sortOrder int
	if err := scanner.Scan(&l.ID, &l.Name, &l.Color, &icon, &description, &sortOrder); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if icon.Valid {
		l.Icon = icon.String
	}
	if description.Valid {
		l.Description = description.String
	}
	l.SortOrder = sortOrder
	return &l, nil
}

func (r *BudgetsRepo) Create(ctx context.Context, input domain.CreateBudgetInput, userID string) (*domain.Budget, error) {
	id := uuid.NewString()
	now := time.Now()
	amount := domain.Cents(input.Amount)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	startsAt, endsAt, err := oneTimeWindow(input.Period, input.StartsAt, input.EndsAt)
	if err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO budgets (id, user_id, name, period, amount, currency, account_id,
		           trending_over, over_budget, starts_at, ends_at, deleted_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, id, userID, input.Name, input.Period, amount, input.Currency, input.AccountID,
		input.TrendingOver, input.OverBudget, startsAt, endsAt, nil, now, now); err != nil {
		return nil, err
	}

	if err := insertRelations(tx, ctx, id, input.CategoryIDs, `budget_categories`, `category_id`); err != nil {
		return nil, err
	}
	if err := insertRelations(tx, ctx, id, input.LabelIDs, `budget_labels`, `label_id`); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.Get(ctx, id, userID)
}

func (r *BudgetsRepo) Update(ctx context.Context, id, userID string, input domain.UpdateBudgetInput) (*domain.Budget, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var owned bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM budgets WHERE id = ? AND user_id = ? AND deleted_at IS NULL)`, id, userID).Scan(&owned); err != nil {
		return nil, err
	}
	if !owned {
		return nil, sql.ErrNoRows
	}

	setParts := []string{}
	args := []interface{}{}

	if input.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *input.Name)
	}
	if input.Period != nil {
		if *input.Period == domain.BudgetPeriodOneTime {
			if input.StartsAt == nil || input.EndsAt == nil {
				tx.Rollback()
				return nil, fmt.Errorf("one_time budget requires starts_at and ends_at")
			}
			s, e, err := oneTimeWindow(*input.Period, input.StartsAt, input.EndsAt)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			setParts = append(setParts, "period = ?, starts_at = ?, ends_at = ?")
			args = append(args, *input.Period, s, e)
		} else {
			setParts = append(setParts, "period = ?, starts_at = NULL, ends_at = NULL")
			args = append(args, *input.Period)
		}
	}
	if input.Amount != nil {
		setParts = append(setParts, "amount = ?")
		args = append(args, domain.Cents(*input.Amount))
	}
	if input.Currency != nil {
		setParts = append(setParts, "currency = ?")
		args = append(args, *input.Currency)
	}
	if input.AccountID != nil {
		setParts = append(setParts, "account_id = ?")
		args = append(args, *input.AccountID)
	}
	if input.TrendingOver != nil {
		setParts = append(setParts, "trending_over = ?")
		args = append(args, *input.TrendingOver)
	}
	if input.OverBudget != nil {
		setParts = append(setParts, "over_budget = ?")
		args = append(args, *input.OverBudget)
	}

	if len(setParts) > 0 {
		args = append(args, id, userID)
		if _, err := tx.ExecContext(ctx, "UPDATE budgets SET "+strings.Join(setParts, ", ")+" WHERE id = ? AND user_id = ?", args...); err != nil {
			return nil, err
		}
	}

	if input.CategoryIDs != nil {
		if _, err := tx.ExecContext(ctx, `DELETE FROM budget_categories WHERE budget_id = ? AND budget_id IN (SELECT id FROM budgets WHERE id = ? AND user_id = ?)`, id, id, userID); err != nil {
			return nil, err
		}
		if err := insertRelations(tx, ctx, id, input.CategoryIDs, `budget_categories`, `category_id`); err != nil {
			return nil, err
		}
	}
	if input.LabelIDs != nil {
		if _, err := tx.ExecContext(ctx, `DELETE FROM budget_labels WHERE budget_id = ? AND budget_id IN (SELECT id FROM budgets WHERE id = ? AND user_id = ?)`, id, id, userID); err != nil {
			return nil, err
		}
		if err := insertRelations(tx, ctx, id, input.LabelIDs, `budget_labels`, `label_id`); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.Get(ctx, id, userID)
}

func (r *BudgetsRepo) Delete(ctx context.Context, id, userID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE budgets SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND user_id = ?`, id, userID)
	return err
}

// insertRelations inserts many-to-many rows for a budget.
func insertRelations(tx interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}, ctx context.Context, budgetID string, ids []string, table, column string) error {
	for _, id := range ids {
		if _, err := tx.ExecContext(ctx, `INSERT INTO `+table+` (budget_id, `+column+`) VALUES (?, ?)`, budgetID, id); err != nil {
			return err
		}
	}
	return nil
}

// oneTimeWindow returns stored starts_at/ends_at for one_time budgets, NULL otherwise.
func oneTimeWindow(period string, startsAt, endsAt *string) (interface{}, interface{}, error) {
	if period != domain.BudgetPeriodOneTime {
		return nil, nil, nil
	}
	if startsAt == nil || endsAt == nil {
		return nil, nil, fmt.Errorf("one_time budget requires starts_at and ends_at")
	}
	s, err := time.Parse(time.RFC3339, *startsAt)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid starts_at: %w", err)
	}
	e, err := time.Parse(time.RFC3339, *endsAt)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid ends_at: %w", err)
	}
	if !s.Before(e) {
		return nil, nil, fmt.Errorf("starts_at must be before ends_at")
	}
	return s.Format(time.RFC3339), e.Format(time.RFC3339), nil
}

// budgetCalc computes consumption for a budget directly from transactions.
func (r *BudgetsRepo) budgetCalc(ctx context.Context, b *domain.Budget) (*domain.BudgetResult, error) {
	periodStart, periodEnd, err := periodWindow(b)
	if err != nil {
		return nil, err
	}

	accountCond := ""
	args := []interface{}{periodStart, periodEnd}
	if b.AccountID != nil {
		accountCond = `AND t.from_account_id = ?`
		args = append(args, *b.AccountID)
	}

	categoryCond := ""
	if len(b.Categories) > 0 {
		catPlaceholders := make([]string, len(b.Categories))
		for i, c := range b.Categories {
			catPlaceholders[i] = "?"
			args = append(args, c.ID)
		}
		categoryCond = `AND EXISTS (SELECT 1 FROM budget_categories bc WHERE bc.budget_id = ? AND bc.category_id = t.category_id)`
		args = append(args, b.ID)
	}

	labelCond := ""
	if len(b.Labels) > 0 {
		labelCond = `AND EXISTS (SELECT 1 FROM transaction_labels tl JOIN budget_labels bl ON bl.label_id = tl.label_id AND bl.budget_id = ? WHERE tl.transaction_id = t.id)`
		args = append(args, b.ID)
	}

	query := `
		SELECT t.amount, t.currency, t.date_time
		FROM transactions t
		WHERE t.type = 'expense'
		  AND t.deleted_at IS NULL
		  AND t.date_time >= ?
		  AND t.date_time < ?
		  ` + accountCond + categoryCond + labelCond

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	totalCents := int64(0)
	for rows.Next() {
		var amt int64
		var currency string
		var dt time.Time
		if err := rows.Scan(&amt, &currency, &dt); err != nil {
			return nil, err
		}
		if currency == b.Currency {
			totalCents += amt
			continue
		}
		rate, err := latestRate(ctx, r.db, currency, b.Currency, dt)
		if err != nil {
			return nil, err
		}
		if rate == nil {
			return nil, fmt.Errorf("missing exchange rate for %s to %s on %s", currency, b.Currency, dt.Format("2006-01-02"))
		}
		totalCents += int64(math.Round(float64(amt) * *rate))
	}

	spentStr := domain.FormatCents(totalCents)
	remainingCents := b.Amount - totalCents
	remainingStr := domain.FormatCents(remainingCents)
	pct := 0
	if b.Amount > 0 {
		pct = int(math.Round(float64(totalCents) / float64(b.Amount) * 100))
	}

	return &domain.BudgetResult{
		BudgetID:        b.ID,
		Name:            b.Name,
		Period:          b.Period,
		Amount:          b.AmountStr,
		Currency:        b.Currency,
		Spent:           spentStr,
		Remaining:       remainingStr,
		PercentageUsed:  pct,
		IsOverBudget:    totalCents >= b.Amount,
		PeriodStart:     periodStart,
		PeriodEnd:       periodEnd,
	}, nil
}

// periodWindow returns inclusive [start, exclusive end) boundaries for the
// budget's current period in process-local time.
func periodWindow(b *domain.Budget) (time.Time, time.Time, error) {
	now := time.Now()
	loc := now.Location()
	switch b.Period {
	case domain.BudgetPeriodWeekly:
		weekday := int(now.Weekday())
		daysSinceMonday := (weekday + 6) % 7
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, -daysSinceMonday)
		return start, start.AddDate(0, 0, 7), nil
	case domain.BudgetPeriodMonthly:
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		return start, start.AddDate(0, 1, 0), nil
	case domain.BudgetPeriodYearly:
		start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, loc)
		return start, start.AddDate(1, 0, 0), nil
	case domain.BudgetPeriodOneTime:
		if b.StartsAt == nil || b.EndsAt == nil || b.StartsAt.IsZero() || b.EndsAt.IsZero() {
			return time.Time{}, time.Time{}, nil
		}
		return b.StartsAt.In(loc), b.EndsAt.In(loc), nil
	}
	return time.Time{}, time.Time{}, fmt.Errorf("unsupported period: %s", b.Period)
}

// latestRate returns the most recent fx_rates rate on or before the given date.
func latestRate(ctx context.Context, db interface {
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}, base, quote string, date time.Time) (*float64, error) {
	dateStr := date.Format("2006-01-02")
	var rate float64
	err := db.QueryRowContext(ctx, `
		SELECT rate FROM fx_rates
		WHERE base = ? AND quote = ? AND date <= ?
		ORDER BY date DESC LIMIT 1
	`, base, quote, dateStr).Scan(&rate)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

var _ interface {
	List(context.Context, string) ([]domain.Budget, error)
	Get(context.Context, string, string) (*domain.Budget, error)
	Create(context.Context, domain.CreateBudgetInput, string) (*domain.Budget, error)
	Update(context.Context, string, string, domain.UpdateBudgetInput) (*domain.Budget, error)
	Delete(context.Context, string, string) error
} = (*BudgetsRepo)(nil)
