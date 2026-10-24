package sqlite

import (
	"context"
	"database/sql"

	"kiroku/backend/internal/domain"
	"kiroku/backend/internal/db"
)

type TemplatesRepo struct{ db *db.DB }

func NewTemplatesRepo(database *db.DB) *TemplatesRepo { return &TemplatesRepo{db: database} }

func (r *TemplatesRepo) List(ctx context.Context) ([]domain.Template, error) {
	return r.query(ctx, `SELECT id, name, type, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, place_lat, place_lng, attachments, created_at, updated_at FROM templates ORDER BY created_at DESC`)
}

func (r *TemplatesRepo) Get(ctx context.Context, id string) (*domain.Template, error) {
	list, err := r.query(ctx, `SELECT id, name, type, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, place_lat, place_lng, attachments, created_at, updated_at FROM templates WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

func (r *TemplatesRepo) Create(ctx context.Context, input domain.CreateTemplateInput) (*domain.Template, error) {
	id := genID()
	_, err := r.db.ExecContext(ctx, `INSERT INTO templates (id, name, type, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, place_lat, place_lng, attachments) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, id, input.Name, input.Type, input.Amount, input.Currency, input.FromAccountID, input.ToAccountID, input.CategoryID, input.Labels, input.Note, input.PayeePayer, input.PaymentType, input.WarrantyUntil, input.Status, input.PlaceName, input.PlaceLat, input.PlaceLng, input.Attachments)
	if err != nil {
		return nil, err
	}
	return &domain.Template{ID: id, Name: input.Name, Type: input.Type, Amount: input.Amount, Currency: input.Currency, FromAccountID: input.FromAccountID, ToAccountID: input.ToAccountID, CategoryID: input.CategoryID, Labels: input.Labels, Note: input.Note, PayeePayer: input.PayeePayer, PaymentType: input.PaymentType, WarrantyUntil: input.WarrantyUntil, Status: input.Status, PlaceName: input.PlaceName, PlaceLat: input.PlaceLat, PlaceLng: input.PlaceLng, Attachments: input.Attachments}, nil
}

func (r *TemplatesRepo) Update(ctx context.Context, id string, input domain.UpdateTemplateInput) (*domain.Template, error) {
	setParts := []string{}
	args := []interface{}{}
	if input.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *input.Name)
	}
	if input.Type != nil {
		setParts = append(setParts, "type = ?")
		args = append(args, *input.Type)
	}
	if input.Amount != nil {
		setParts = append(setParts, "amount = ?")
		args = append(args, *input.Amount)
	}
	if input.Currency != nil {
		setParts = append(setParts, "currency = ?")
		args = append(args, *input.Currency)
	}
	if input.FromAccountID != nil {
		setParts = append(setParts, "from_account_id = ?")
		args = append(args, *input.FromAccountID)
	}
	if input.ToAccountID != nil {
		setParts = append(setParts, "to_account_id = ?")
		args = append(args, *input.ToAccountID)
	}
	if input.CategoryID != nil {
		setParts = append(setParts, "category_id = ?")
		args = append(args, *input.CategoryID)
	}
	if input.Labels != nil {
		setParts = append(setParts, "labels = ?")
		args = append(args, input.Labels)
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
		args = append(args, *input.WarrantyUntil)
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
		setParts = append(setParts, "attachments = ?")
		args = append(args, input.Attachments)
	}
	if len(setParts) == 0 {
		return r.Get(ctx, id)
	}
	args = append(args, id)
	_, err := r.db.ExecContext(ctx, "UPDATE templates SET "+joinStrings(setParts, ",")+" WHERE id = ?", args...)
	if err != nil {
		return nil, err
	}
	return r.Get(ctx, id)
}

func (r *TemplatesRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM templates WHERE id = ?", id)
	return err
}

func (r *TemplatesRepo) query(ctx context.Context, query string, args ...interface{}) ([]domain.Template, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Template
	for rows.Next() {
		var t domain.Template
		var amount sql.NullInt64
		var currency, fromAcc, toAcc, catID, paymentType sql.NullString
		var warranty sql.NullTime
		var placeLat, placeLng sql.NullFloat64
		if err := rows.Scan(&t.ID, &t.Name, &t.Type, &amount, &currency, &fromAcc, &toAcc, &catID, &t.Labels, &t.Note, &t.PayeePayer, &paymentType, &warranty, &t.Status, &t.PlaceName, &placeLat, &placeLng, &t.Attachments, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		if amount.Valid {
			v := float64(amount.Int64)
			t.Amount = v
		}
		if currency.Valid {
			t.Currency = currency.String
		}
		if fromAcc.Valid {
			t.FromAccountID = &fromAcc.String
		}
		if toAcc.Valid {
			t.ToAccountID = &toAcc.String
		}
		if catID.Valid {
			t.CategoryID = &catID.String
		}
		if paymentType.Valid {
			t.PaymentType = paymentType.String
		}
		if warranty.Valid {
			t.WarrantyUntil = &warranty.Time
		}
		if placeLat.Valid {
			t.PlaceLat = &placeLat.Float64
		}
		if placeLng.Valid {
			t.PlaceLng = &placeLng.Float64
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
