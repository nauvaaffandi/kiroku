package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"kiroku/backend/internal/domain"
)

func (r *TransactionsRepo) Create(ctx context.Context, input domain.CreateTransactionInput) (*domain.Transaction, error) {
	id := genID()
	status := input.Status
	if status == "" {
		status = "cleared"
	}
	labelsStr := strings.Join(input.Labels, ",")
	attachmentsStr, _ := jsonMarshal(input.Attachments)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, place_lat, place_lng, attachments, template_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, id, input.Type, input.DateTime, input.Amount, input.Currency, input.FromAccountID, input.ToAccountID, input.CategoryID, labelsStr, input.Note, input.PayeePayer, input.PaymentType, input.WarrantyUntil, status, input.PlaceName, input.PlaceLat, input.PlaceLng, string(attachmentsStr), input.TemplateID)
	if err != nil {
		return nil, err
	}

	if len(input.LabelIDs) > 0 {
		for _, lid := range input.LabelIDs {
			if _, err := tx.ExecContext(ctx, `INSERT INTO transaction_labels (transaction_id, label_id) VALUES (?, ?)`, id, lid); err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.Get(ctx, id)
}

func (r *TransactionsRepo) Update(ctx context.Context, id string, input domain.UpdateTransactionInput) (*domain.Transaction, error) {
	setParts := []string{}
	args := []interface{}{}
	if input.DateTime != nil {
		setParts = append(setParts, "date_time = ?")
		args = append(args, input.DateTime.Format(time.RFC3339))
	}
	if input.Amount != nil {
		setParts = append(setParts, "amount = ?")
		args = append(args, *input.Amount)
	}
	if input.CategoryID != nil {
		setParts = append(setParts, "category_id = ?")
		args = append(args, *input.CategoryID)
	}
	if input.Labels != nil {
		setParts = append(setParts, "labels = ?")
		args = append(args, strings.Join(input.Labels, ","))
	}
	if input.Note != nil {
		setParts = append(setParts, "note = ?")
		args = append(args, *input.Note)
	}
	if input.PayeePayer != nil {
		setParts = append(setParts, "payee_payer = ?")
		args = append(args, *input.PayeePayer)
	}
	if input.PaymentType != nil {
		setParts = append(setParts, "payment_type = ?")
		args = append(args, *input.PaymentType)
	}
	if input.WarrantyUntil != nil {
		setParts = append(setParts, "warranty_until = ?")
		args = append(args, input.WarrantyUntil.Format("2006-01-02"))
	}
	if input.Status != nil {
		setParts = append(setParts, "status = ?")
		args = append(args, *input.Status)
	}
	if input.PlaceName != nil {
		setParts = append(setParts, "place_name = ?")
		args = append(args, *input.PlaceName)
	}
	if input.PlaceLat != nil {
		setParts = append(setParts, "place_lat = ?")
		args = append(args, *input.PlaceLat)
	}
	if input.PlaceLng != nil {
		setParts = append(setParts, "place_lng = ?")
		args = append(args, *input.PlaceLng)
	}
	if input.Attachments != nil {
		attStr, _ := jsonMarshal(input.Attachments)
		setParts = append(setParts, "attachments = ?")
		args = append(args, string(attStr))
	}
	if len(setParts) == 0 && input.LabelIDs == nil {
		return r.Get(ctx, id)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if len(setParts) > 0 {
		args = append(args, id)
		_, err := tx.ExecContext(ctx, "UPDATE transactions SET "+strings.Join(setParts, ", ")+" WHERE id = ?", args...)
		if err != nil {
			return nil, err
		}
	}

	if input.LabelIDs != nil {
		_, err := tx.ExecContext(ctx, `DELETE FROM transaction_labels WHERE transaction_id = ?`, id)
		if err != nil {
			return nil, err
		}
		for _, lid := range *input.LabelIDs {
			_, err := tx.ExecContext(ctx, `INSERT INTO transaction_labels (transaction_id, label_id) VALUES (?, ?)`, id, lid)
			if err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.Get(ctx, id)
}

func (r *TransactionsRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE transactions SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?", id)
	return err
}

func (r *TransactionsRepo) Batch(ctx context.Context, input domain.BatchInput) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	affected := 0
	for _, id := range input.IDs {
		var res sql.Result
		var execErr error
		switch input.Action {
		case domain.BatchDelete:
			res, execErr = tx.ExecContext(ctx, "UPDATE transactions SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL", id)
		case domain.BatchRecategorize:
			res, execErr = tx.ExecContext(ctx, "UPDATE transactions SET category_id = ? WHERE id = ? AND deleted_at IS NULL", input.CategoryID, id)
		case domain.BatchRelabel:
			labels := strings.Join(input.Labels, ",")
			res, execErr = tx.ExecContext(ctx, "UPDATE transactions SET labels = ? WHERE id = ? AND deleted_at IS NULL", labels, id)
		default:
			tx.Rollback()
			return 0, fmt.Errorf("unknown action: %s", input.Action)
		}
		if execErr != nil {
			continue
		}
		if n, _ := res.RowsAffected(); n > 0 {
			affected++
		}
	}
	return affected, tx.Commit()
}