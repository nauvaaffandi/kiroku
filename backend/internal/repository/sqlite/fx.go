package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"kiroku/backend/internal/domain"
	"kiroku/backend/internal/db"
)

type FXRepo struct{ db *db.DB }

func NewFXRepo(database *db.DB) *FXRepo { return &FXRepo{db: database} }

func (r *FXRepo) GetQuotes(ctx context.Context, base string, date time.Time) ([]domain.FXQuote, error) {
	dateStr := date.Format("2006-01-02")
	rows, err := r.db.QueryContext(ctx, `SELECT quote, rate FROM fx_rates WHERE base = ? AND date = ?`, base, dateStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.FXQuote
	for rows.Next() {
		var q domain.FXQuote
		q.Base = base
		q.UpdatedAt = date
		if err := rows.Scan(&q.Quote, &q.Rate); err != nil {
			return nil, err
		}
		out = append(out, q)
	}
	return out, rows.Err()
}

func (r *FXRepo) Upsert(ctx context.Context, quotes []domain.FXQuote) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	for _, q := range quotes {
		if q.Base == q.Quote {
			continue
		}
		date := q.UpdatedAt.Format("2006-01-02")
		if _, err := tx.ExecContext(ctx, `INSERT OR REPLACE INTO fx_rates (base, quote, rate, date, source, updated_at) VALUES (?, ?, ?, ?, 'exchangerate.host', CURRENT_TIMESTAMP)`, q.Base, q.Quote, q.Rate, date); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *FXRepo) Count(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM fx_rates`).Scan(&n)
	return n, err
}

func (r *FXRepo) GetRate(ctx context.Context, base, quote string, date time.Time) (*float64, error) {
	dateStr := date.Format("2006-01-02")
	var rate float64
	err := r.db.QueryRowContext(ctx, `SELECT rate FROM fx_rates WHERE base = ? AND quote = ? AND date = ?`, base, quote, dateStr).Scan(&rate)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

func FetchAndStoreFX(database *db.DB, baseCurrency string) error {
	if baseCurrency == "" {
		baseCurrency = "JPY"
	}
	url := fmt.Sprintf("https://api.exchangerate.host/latest?base=%s", baseCurrency)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("fetch failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read failed: %w", err)
	}
	var result struct {
		Base  string             `json:"base"`
		Date  string             `json:"date"`
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("parse failed: %w", err)
	}

	date := result.Date
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	tx, err := database.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("tx begin: %w", err)
	}
	for quote, rate := range result.Rates {
		if quote == baseCurrency {
			continue
		}
		if _, err := tx.Exec(`INSERT OR REPLACE INTO fx_rates (base, quote, rate, date, source, updated_at) VALUES (?, ?, ?, ?, 'exchangerate.host', CURRENT_TIMESTAMP)`, baseCurrency, quote, rate, date); err != nil {
			tx.Rollback()
			return fmt.Errorf("insert rate: %w", err)
		}
	}
	return tx.Commit()
}
