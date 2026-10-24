package sqlite

import (
	"context"
	"database/sql"

	"kiroku/backend/internal/domain"
	"kiroku/backend/internal/db"
)

type ScheduledRepo struct{ db *db.DB }

func NewScheduledRepo(database *db.DB) *ScheduledRepo { return &ScheduledRepo{db: database} }

func (r *ScheduledRepo) List(ctx context.Context) ([]domain.ScheduledTransaction, error) {
	return r.query(ctx, `SELECT id, rrule, max_occurrences, transaction_template_json, enabled, next_run, occurrences_done, created_at, updated_at FROM scheduled ORDER BY created_at DESC`)
}

func (r *ScheduledRepo) Get(ctx context.Context, id string) (*domain.ScheduledTransaction, error) {
	list, err := r.query(ctx, `SELECT id, rrule, max_occurrences, transaction_template_json, enabled, next_run, occurrences_done, created_at, updated_at FROM scheduled WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

func (r *ScheduledRepo) Create(ctx context.Context, input domain.CreateScheduledInput) (*domain.ScheduledTransaction, error) {
	id := genID()
	tmplJSON, err := jsonMarshal(input.TransactionTemplate)
	if err != nil {
		return nil, err
	}
	enabled := input.Enabled
	nextRun := computeNextRun(input.RRule)
	_, err = r.db.ExecContext(ctx, `INSERT INTO scheduled (id, rrule, max_occurrences, transaction_template_json, enabled, next_run, occurrences_done) VALUES (?, ?, ?, ?, ?, ?, 0)`, id, input.RRule, input.MaxOccurrences, string(tmplJSON), enabled, nextRun)
	if err != nil {
		return nil, err
	}
	return &domain.ScheduledTransaction{ID: id, RRule: input.RRule, MaxOccurrences: input.MaxOccurrences, TransactionTemplate: tmplJSON, Enabled: enabled, NextRun: nextRun}, nil
}

func (r *ScheduledRepo) Update(ctx context.Context, id string, input domain.UpdateScheduledInput) (*domain.ScheduledTransaction, error) {
	setParts := []string{}
	args := []interface{}{}
	if input.RRule != nil {
		setParts = append(setParts, "rrule = ?")
		args = append(args, *input.RRule)
		nextRun := computeNextRun(*input.RRule)
		setParts = append(setParts, "next_run = ?")
		args = append(args, nextRun)
	}
	if input.MaxOccurrences != nil {
		setParts = append(setParts, "max_occurrences = ?")
		args = append(args, *input.MaxOccurrences)
	}
	if input.TransactionTemplate != nil {
		tmplJSON, err := jsonMarshal(input.TransactionTemplate)
		if err != nil {
			return nil, err
		}
		setParts = append(setParts, "transaction_template_json = ?")
		args = append(args, string(tmplJSON))
	}
	if input.Enabled != nil {
		setParts = append(setParts, "enabled = ?")
		if *input.Enabled {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
	}
	if len(setParts) == 0 {
		return r.Get(ctx, id)
	}
	args = append(args, id)
	_, err := r.db.ExecContext(ctx, "UPDATE scheduled SET "+joinStrings(setParts, ",")+" WHERE id = ?", args...)
	if err != nil {
		return nil, err
	}
	return r.Get(ctx, id)
}

func (r *ScheduledRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM scheduled WHERE id = ?", id)
	return err
}

func (r *ScheduledRepo) Toggle(ctx context.Context, id string, enabled bool) (*domain.ScheduledTransaction, error) {
	_, err := r.db.ExecContext(ctx, "UPDATE scheduled SET enabled = NOT enabled WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return r.Get(ctx, id)
}

func (r *ScheduledRepo) query(ctx context.Context, query string, args ...interface{}) ([]domain.ScheduledTransaction, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.ScheduledTransaction
	for rows.Next() {
		var s domain.ScheduledTransaction
		var maxOcc sql.NullInt64
		var nextRun sql.NullTime
		if err := rows.Scan(&s.ID, &s.RRule, &maxOcc, &s.TransactionTemplate, &s.Enabled, &nextRun, &s.OccurrencesDone, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		if maxOcc.Valid {
			v := int(maxOcc.Int64)
			s.MaxOccurrences = &v
		}
		if nextRun.Valid {
			s.NextRun = &nextRun.Time
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
