package sqlite

import (
	"context"
	"database/sql"

	"kiroku/backend/internal/domain"
)

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

var _ interface {
	List(context.Context, string) ([]domain.Budget, error)
	Get(context.Context, string, string) (*domain.Budget, error)
} = (*BudgetsRepo)(nil)
