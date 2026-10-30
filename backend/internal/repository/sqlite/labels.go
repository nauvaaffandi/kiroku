package sqlite

import (
	"context"
	"database/sql"
	"strings"

	"kiroku/backend/internal/db"
	"kiroku/backend/internal/domain"
)

type LabelsRepo struct {
	db *db.DB
}

func NewLabelsRepo(database *db.DB) *LabelsRepo {
	return &LabelsRepo{db: database}
}

func (r *LabelsRepo) List(ctx context.Context) ([]domain.Label, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, color, icon, description, sort_order, created_at, updated_at FROM labels ORDER BY sort_order ASC, name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Label
	for rows.Next() {
		var l domain.Label
		var icon, desc sql.NullString
		if err := rows.Scan(&l.ID, &l.Name, &l.Color, &icon, &desc, &l.SortOrder, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		if icon.Valid {
			l.Icon = icon.String
		}
		if desc.Valid {
			l.Description = desc.String
		}
		out = append(out, l)
	}
	if out == nil {
		out = []domain.Label{}
	}
	return out, rows.Err()
}

func (r *LabelsRepo) Get(ctx context.Context, id string) (*domain.Label, error) {
	var l domain.Label
	var icon, desc sql.NullString
	err := r.db.QueryRowContext(ctx, `SELECT id, name, color, icon, description, sort_order, created_at, updated_at FROM labels WHERE id = ?`, id).
		Scan(&l.ID, &l.Name, &l.Color, &icon, &desc, &l.SortOrder, &l.CreatedAt, &l.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if icon.Valid {
		l.Icon = icon.String
	}
	if desc.Valid {
		l.Description = desc.String
	}
	return &l, nil
}

func (r *LabelsRepo) Create(ctx context.Context, input domain.CreateLabelInput) (*domain.Label, error) {
	id := genID()
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, sql.ErrNoRows // or custom validation error, let repository return what db or standard does
	}
	color := "#6366f1"
	if input.Color != nil && *input.Color != "" {
		color = *input.Color
	}

	var icon, desc interface{}
	if input.Icon != nil {
		icon = *input.Icon
	}
	if input.Description != nil {
		desc = *input.Description
	}

	sortOrder := 0
	if input.SortOrder != nil {
		sortOrder = *input.SortOrder
	} else {
		// find max sort order
		var maxOrder sql.NullInt64
		_ = r.db.QueryRowContext(ctx, `SELECT MAX(sort_order) FROM labels`).Scan(&maxOrder)
		if maxOrder.Valid {
			sortOrder = int(maxOrder.Int64) + 1
		}
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO labels (id, name, color, icon, description, sort_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, id, name, color, icon, desc, sortOrder)
	if err != nil {
		return nil, err
	}

	return r.Get(ctx, id)
}

func (r *LabelsRepo) Update(ctx context.Context, id string, input domain.UpdateLabelInput) (*domain.Label, error) {
	setParts := []string{"updated_at = CURRENT_TIMESTAMP"}
	args := []interface{}{}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		setParts = append(setParts, "name = ?")
		args = append(args, name)
	}
	if input.Color != nil {
		setParts = append(setParts, "color = ?")
		args = append(args, *input.Color)
	}
	if input.Icon != nil {
		setParts = append(setParts, "icon = ?")
		args = append(args, *input.Icon)
	}
	if input.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *input.Description)
	}
	if input.SortOrder != nil {
		setParts = append(setParts, "sort_order = ?")
		args = append(args, *input.SortOrder)
	}

	args = append(args, id)
	query := `UPDATE labels SET ` + strings.Join(setParts, ", ") + ` WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return r.Get(ctx, id)
}

func (r *LabelsRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM labels WHERE id = ?`, id)
	return err
}

func (r *LabelsRepo) Reorder(ctx context.Context, items []domain.ReorderItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, item := range items {
		_, err := tx.ExecContext(ctx, `UPDATE labels SET sort_order = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, item.SortOrder, item.ID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
