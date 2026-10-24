package sqlite

import (
	"context"
	"database/sql"

	"kiroku/backend/internal/domain"
	"kiroku/backend/internal/db"
)

type CategoriesRepo struct{ db *db.DB }

func NewCategoriesRepo(database *db.DB) *CategoriesRepo { return &CategoriesRepo{db: database} }

func (r *CategoriesRepo) List(ctx context.Context, includeArchived bool) ([]domain.Category, error) {
	query := `SELECT id, parent_id, name, icon, color, kind, sort_order FROM categories`
	args := []interface{}{}
	if !includeArchived {
		query += ` WHERE archived = 0`
	}
	query += ` ORDER BY sort_order, name`
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Category
	for rows.Next() {
		var c domain.Category
		var parentID sql.NullString
		var kind sql.NullString
		if err := rows.Scan(&c.ID, &parentID, &c.Name, &c.Icon, &c.Color, &kind, &c.SortOrder); err != nil {
			return nil, err
		}
		if parentID.Valid {
			c.ParentID = &parentID.String
		}
		if kind.Valid {
			c.Kind = &kind.String
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *CategoriesRepo) Tree(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, parent_id, name, icon, color, kind, sort_order FROM categories WHERE archived = 0 ORDER BY sort_order, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	catMap := map[string]*domain.Category{}
	var roots []domain.Category

	for rows.Next() {
		var c domain.Category
		var parentID sql.NullString
		var kind sql.NullString
		if err := rows.Scan(&c.ID, &parentID, &c.Name, &c.Icon, &c.Color, &kind, &c.SortOrder); err != nil {
			return nil, err
		}
		if parentID.Valid {
			c.ParentID = &parentID.String
		}
		if kind.Valid {
			c.Kind = &kind.String
		}
		c.Children = []*domain.Category{}
		catMap[c.ID] = &c
	}

	for _, c := range catMap {
		if c.ParentID == nil {
			roots = append(roots, *c)
		} else if parent, ok := catMap[*c.ParentID]; ok {
			parent.Children = append(parent.Children, c)
		}
	}
	return roots, rows.Err()
}

func (r *CategoriesRepo) Get(ctx context.Context, id string) (*domain.Category, error) {
	var c domain.Category
	var parentID sql.NullString
	var kind sql.NullString
	err := r.db.QueryRowContext(ctx, `SELECT id, parent_id, name, icon, color, kind, sort_order FROM categories WHERE id = ?`, id).Scan(&c.ID, &parentID, &c.Name, &c.Icon, &c.Color, &kind, &c.SortOrder)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if parentID.Valid {
		c.ParentID = &parentID.String
	}
	if kind.Valid {
		c.Kind = &kind.String
	}
	return &c, nil
}

func (r *CategoriesRepo) Create(ctx context.Context, input domain.CreateCategoryInput) (*domain.Category, error) {
	id := genID()
	sortOrder, err := r.NextSortOrder(ctx, input.ParentID)
	if err != nil {
		return nil, err
	}

	color := input.Color
	if color == "" {
		color = randomPastel()
	}

	_, err = r.db.ExecContext(ctx, `INSERT INTO categories (id, parent_id, name, icon, color, kind, sort_order) VALUES (?, ?, ?, ?, ?, ?, ?)`, id, input.ParentID, input.Name, input.Icon, color, input.Kind, sortOrder+1)
	if err != nil {
		return nil, err
	}
	return &domain.Category{
		ID:        id,
		ParentID:  input.ParentID,
		Name:      input.Name,
		Icon:      input.Icon,
		Color:     color,
		Kind:      input.Kind,
		SortOrder: sortOrder + 1,
	}, nil
}

func (r *CategoriesRepo) Update(ctx context.Context, id string, input domain.UpdateCategoryInput) (*domain.Category, error) {
	setParts := []string{}
	args := []interface{}{}
	if input.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *input.Name)
	}
	if input.ParentID != nil {
		setParts = append(setParts, "parent_id = ?")
		args = append(args, *input.ParentID)
	}
	if input.Icon != nil {
		setParts = append(setParts, "icon = ?")
		args = append(args, *input.Icon)
	}
	if input.Color != nil {
		setParts = append(setParts, "color = ?")
		args = append(args, *input.Color)
	}
	if input.Kind != nil {
		setParts = append(setParts, "kind = ?")
		args = append(args, *input.Kind)
	}
	if input.Archived != nil {
		setParts = append(setParts, "archived = ?")
		if *input.Archived {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
	}
	if input.SortOrder != nil {
		setParts = append(setParts, "sort_order = ?")
		args = append(args, *input.SortOrder)
	}

	if len(setParts) == 0 {
		return r.Get(ctx, id)
	}
	args = append(args, id)
	_, err := r.db.ExecContext(ctx, "UPDATE categories SET "+joinStrings(setParts, ",")+" WHERE id = ?", args...)
	if err != nil {
		return nil, err
	}
	return r.Get(ctx, id)
}

func (r *CategoriesRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE categories SET archived = 1 WHERE id = ? OR id IN (WITH RECURSIVE descendants AS (SELECT id FROM categories WHERE parent_id = ? UNION ALL SELECT c.id FROM categories c JOIN descendants d ON c.parent_id = d.id) SELECT id FROM descendants)`, id, id)
	return err
}

func (r *CategoriesRepo) Reorder(ctx context.Context, items []domain.ReorderItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	for _, item := range items {
		if _, err := tx.ExecContext(ctx, "UPDATE categories SET sort_order = ? WHERE id = ?", item.SortOrder, item.ID); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *CategoriesRepo) CountChildren(ctx context.Context, id string) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM categories WHERE parent_id = ?`, id).Scan(&n)
	return n, err
}

func (r *CategoriesRepo) NextSortOrder(ctx context.Context, parentID *string) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COALESCE(MAX(sort_order), 0) FROM categories WHERE parent_id IS ?`, parentID).Scan(&n)
	return n, err
}
