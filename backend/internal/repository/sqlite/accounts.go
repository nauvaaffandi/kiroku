package sqlite

import (
	"context"
	"database/sql"

	"kiroku/backend/internal/domain"
	"kiroku/backend/internal/db"
)

type AccountsRepo struct{ db *db.DB }

func NewAccountsRepo(database *db.DB) *AccountsRepo { return &AccountsRepo{db: database} }

func (r *AccountsRepo) List(ctx context.Context) ([]domain.AccountWithBalance, error) {
	query := `SELECT a.id, a.name, a.currency, a.type, a.icon, a.color, a.archived,
		COALESCE(SUM(CASE WHEN t.to_account_id = a.id THEN t.amount ELSE 0 END), 0) -
		COALESCE(SUM(CASE WHEN t.from_account_id = a.id THEN t.amount ELSE 0 END), 0) AS balance
		FROM accounts a LEFT JOIN transactions t ON (t.from_account_id = a.id OR t.to_account_id = a.id) AND t.deleted_at IS NULL
		WHERE a.archived = FALSE GROUP BY a.id ORDER BY a.created_at`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.AccountWithBalance
	for rows.Next() {
		var a domain.AccountWithBalance
		if err := rows.Scan(&a.ID, &a.Name, &a.Currency, &a.Type, &a.Icon, &a.Color, &a.Archived, &a.Balance); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *AccountsRepo) Get(ctx context.Context, id string) (*domain.AccountWithBalance, error) {
	query := `SELECT a.id, a.name, a.currency, a.type, a.icon, a.color, a.archived,
		COALESCE(SUM(CASE WHEN t.to_account_id = a.id THEN t.amount ELSE 0 END), 0) -
		COALESCE(SUM(CASE WHEN t.from_account_id = a.id THEN t.amount ELSE 0 END), 0) AS balance
		FROM accounts a LEFT JOIN transactions t ON (t.from_account_id = a.id OR t.to_account_id = a.id) AND t.deleted_at IS NULL
		WHERE a.id = ? AND a.archived = FALSE GROUP BY a.id`
	var a domain.AccountWithBalance
	err := r.db.QueryRowContext(ctx, query, id).Scan(&a.ID, &a.Name, &a.Currency, &a.Type, &a.Icon, &a.Color, &a.Archived, &a.Balance)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AccountsRepo) Create(ctx context.Context, input domain.CreateAccountInput) (*domain.Account, error) {
	id := genID()
	_, err := r.db.ExecContext(ctx, `INSERT INTO accounts (id, name, currency, type, icon, color, archived) VALUES (?, ?, ?, ?, ?, ?, 0)`, id, input.Name, input.Currency, input.Type, input.Icon, input.Color)
	if err != nil {
		return nil, err
	}
	return &domain.Account{ID: id, Name: input.Name, Currency: input.Currency, Type: input.Type, Icon: input.Icon, Color: input.Color}, nil
}

func (r *AccountsRepo) Update(ctx context.Context, id string, input domain.UpdateAccountInput) (*domain.Account, error) {
	setParts := []string{}
	args := []interface{}{}
	if input.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *input.Name)
	}
	if input.Currency != nil {
		setParts = append(setParts, "currency = ?")
		args = append(args, *input.Currency)
	}
	if input.Type != nil {
		setParts = append(setParts, "type = ?")
		args = append(args, *input.Type)
	}
	if input.Icon != nil {
		setParts = append(setParts, "icon = ?")
		args = append(args, *input.Icon)
	}
	if input.Color != nil {
		setParts = append(setParts, "color = ?")
		args = append(args, *input.Color)
	}
	if input.Archived != nil {
		setParts = append(setParts, "archived = ?")
		if *input.Archived {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
	}
	if len(setParts) == 0 {
		return r.getPlain(ctx, id)
	}
	args = append(args, id)
	_, err := r.db.ExecContext(ctx, "UPDATE accounts SET "+joinStrings(setParts, ",")+" WHERE id = ?", args...)
	if err != nil {
		return nil, err
	}
	return r.getPlain(ctx, id)
}

func (r *AccountsRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE accounts SET archived = 1 WHERE id = ?", id)
	return err
}

func (r *AccountsRepo) Summary(ctx context.Context) ([]domain.CurrencyBalanceSummary, error) {
	query := `SELECT a.currency,
		COALESCE(SUM(CASE WHEN t.to_account_id = a.id THEN t.amount ELSE 0 END), 0) -
		COALESCE(SUM(CASE WHEN t.from_account_id = a.id THEN t.amount ELSE 0 END), 0) AS balance
		FROM accounts a LEFT JOIN transactions t ON (t.from_account_id = a.id OR t.to_account_id = a.id) AND t.deleted_at IS NULL
		WHERE a.archived = FALSE GROUP BY a.currency`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.CurrencyBalanceSummary
	for rows.Next() {
		var s domain.CurrencyBalanceSummary
		if err := rows.Scan(&s.Currency, &s.Balance); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *AccountsRepo) CountTransactions(ctx context.Context, accountID string) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM transactions WHERE (from_account_id = ? OR to_account_id = ?) AND deleted_at IS NULL`, accountID, accountID).Scan(&n)
	return n, err
}

func (r *AccountsRepo) getPlain(ctx context.Context, id string) (*domain.Account, error) {
	var a domain.Account
	err := r.db.QueryRowContext(ctx, `SELECT id, name, currency, type, icon, color, archived FROM accounts WHERE id = ?`, id).Scan(&a.ID, &a.Name, &a.Currency, &a.Type, &a.Icon, &a.Color, &a.Archived)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}
