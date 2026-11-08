package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"kiroku/backend/internal/db"
	"kiroku/backend/internal/domain"
)

type GoalsRepo struct{ db *db.DB }

func NewGoalsRepo(database *db.DB) *GoalsRepo { return &GoalsRepo{db: database} }

func (r *GoalsRepo) List(ctx context.Context) ([]domain.Goal, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, target_amount, icon, desired_date, color_theme, note, saved_already, currency, deleted_at, created_at, updated_at FROM goals WHERE deleted_at IS NULL ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Goal
	for rows.Next() {
		var g domain.Goal
		var desiredDate sql.NullString
		var deletedAt sql.NullTime
		if err := rows.Scan(&g.ID, &g.Name, &g.TargetAmount, &g.Icon, &desiredDate, &g.ColorTheme, &g.Note, &g.SavedAlready, &g.Currency, &deletedAt, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		g.TargetAmount = g.TargetAmount / 100
		g.SavedAlready = g.SavedAlready / 100
		if desiredDate.Valid {
			g.DesiredDate = desiredDate.String
		}
		if deletedAt.Valid {
			g.DeletedAt = &deletedAt.Time
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

func (r *GoalsRepo) Get(ctx context.Context, id string) (*domain.Goal, error) {
	var g domain.Goal
	var desiredDate sql.NullString
	var deletedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `SELECT id, name, target_amount, icon, desired_date, color_theme, note, saved_already, currency, deleted_at, created_at, updated_at FROM goals WHERE id = ? AND deleted_at IS NULL`, id).Scan(&g.ID, &g.Name, &g.TargetAmount, &g.Icon, &desiredDate, &g.ColorTheme, &g.Note, &g.SavedAlready, &g.Currency, &deletedAt, &g.CreatedAt, &g.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	g.TargetAmount = g.TargetAmount / 100
	g.SavedAlready = g.SavedAlready / 100
	if desiredDate.Valid {
		g.DesiredDate = desiredDate.String
	}
	if deletedAt.Valid {
		g.DeletedAt = &deletedAt.Time
	}
	return &g, nil
}

func (r *GoalsRepo) Create(ctx context.Context, input domain.CreateGoalInput) (*domain.Goal, error) {
	id := genID()
	var desiredDate interface{}
	if input.DesiredDate != "" {
		dd, err := parseDDMMYYYY(input.DesiredDate)
		if err != nil {
			return nil, err
		}
		desiredDate = dd
	}
	_, err := r.db.ExecContext(ctx, `INSERT INTO goals (id, name, target_amount, icon, desired_date, color_theme, note, saved_already, currency) VALUES (?, ?, ?, ?, ?, ?, ?, 0, ?)`, id, input.Name, int(input.TargetAmount*100), input.Icon, desiredDate, input.ColorTheme, input.Note, input.Currency)
	if err != nil {
		return nil, err
	}
	return &domain.Goal{
		ID:           id,
		Name:         input.Name,
		TargetAmount: input.TargetAmount,
		Icon:         input.Icon,
		DesiredDate:  input.DesiredDate,
		ColorTheme:   input.ColorTheme,
		Note:         input.Note,
		Currency:     input.Currency,
	}, nil
}

func (r *GoalsRepo) Update(ctx context.Context, id string, input domain.UpdateGoalInput) (*domain.Goal, error) {
	setParts := []string{}
	args := []interface{}{}
	if input.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *input.Name)
	}
	if input.TargetAmount != nil {
		setParts = append(setParts, "target_amount = ?")
		args = append(args, int(*input.TargetAmount*100))
	}
	if input.Icon != nil {
		setParts = append(setParts, "icon = ?")
		args = append(args, *input.Icon)
	}
	if input.DesiredDate != nil {
		setParts = append(setParts, "desired_date = ?")
		if *input.DesiredDate != "" {
			dd, err := parseDDMMYYYY(*input.DesiredDate)
			if err != nil {
				return nil, err
			}
			args = append(args, dd)
		} else {
			args = append(args, nil)
		}
	}
	if input.ColorTheme != nil {
		setParts = append(setParts, "color_theme = ?")
		args = append(args, *input.ColorTheme)
	}
	if input.Note != nil {
		setParts = append(setParts, "note = ?")
		args = append(args, *input.Note)
	}
	if input.Currency != nil {
		setParts = append(setParts, "currency = ?")
		args = append(args, *input.Currency)
	}

	if len(setParts) == 0 {
		return r.Get(ctx, id)
	}
	args = append(args, id)
	_, err := r.db.ExecContext(ctx, "UPDATE goals SET "+joinStrings(setParts, ",")+" WHERE id = ?", args...)
	if err != nil {
		return nil, err
	}
	return r.Get(ctx, id)
}

func (r *GoalsRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE goals SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
	return err
}

func (r *GoalsRepo) Contribute(ctx context.Context, goalID string, input domain.ContributeGoalInput) (*domain.Goal, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var current int
	err = tx.QueryRowContext(ctx, `SELECT saved_already FROM goals WHERE id = ? AND deleted_at IS NULL`, goalID).Scan(&current)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	contribID := genID()
	if _, err := tx.ExecContext(ctx, `INSERT INTO goal_contributions (id, goal_id, amount, note) VALUES (?, ?, ?, ?)`, contribID, goalID, int(input.Amount*100), input.Note); err != nil {
		return nil, err
	}

	newSaved := current + int(input.Amount*100)
	if _, err := tx.ExecContext(ctx, `UPDATE goals SET saved_already = ? WHERE id = ?`, newSaved, goalID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	g, err := r.Get(ctx, goalID)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (r *GoalsRepo) Reach(ctx context.Context, id string) (*domain.Goal, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var target int
	err = tx.QueryRowContext(ctx, `SELECT target_amount FROM goals WHERE id = ? AND deleted_at IS NULL`, id).Scan(&target)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `UPDATE goals SET saved_already = target_amount WHERE id = ?`, id); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	g, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (r *GoalsRepo) Contributions(ctx context.Context, goalID string) ([]domain.GoalContribution, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, goal_id, amount, added_at, note FROM goal_contributions WHERE goal_id = ? ORDER BY added_at DESC`, goalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.GoalContribution
	for rows.Next() {
		var c domain.GoalContribution
		if err := rows.Scan(&c.ID, &c.GoalID, &c.Amount, &c.AddedAt, &c.Note); err != nil {
			return nil, err
		}
		c.Amount = c.Amount / 100
		out = append(out, c)
	}
	return out, rows.Err()
}

func parseDDMMYYYY(s string) (string, error) {
	t, err := time.Parse("02-01-2006", s)
	if err != nil {
		return "", fmt.Errorf("invalid desired_date (expected DD-MM-YYYY): %w", err)
	}
	return t.Format("2006-01-02"), nil
}

func formatDDMMYYYY(s string) string {
	if s == "" {
		return ""
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return s
	}
	return t.Format("02-01-2006")
}
