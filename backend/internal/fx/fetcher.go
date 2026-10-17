package fx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/db"
)

type Fetcher struct {
	db     *db.DB
	base   string
	client *http.Client
}

func NewFetcher(database *db.DB, baseCurrency string) *Fetcher {
	return &Fetcher{
		db:     database,
		base:   baseCurrency,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (f *Fetcher) FetchAndStore() error {
	if f.base == "" {
		f.base = "JPY"
	}

	url := fmt.Sprintf("https://api.exchangerate.host/latest?base=%s", f.base)
	resp, err := f.client.Get(url)
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

	tx, err := f.db.Begin()
	if err != nil {
		return fmt.Errorf("tx begin: %w", err)
	}

	for quote, rate := range result.Rates {
		if quote == f.base {
			continue
		}
		_, err = tx.Exec(`
			INSERT OR REPLACE INTO fx_rates (base, quote, rate, date, source, updated_at)
			VALUES (?, ?, ?, ?, 'exchangerate.host', CURRENT_TIMESTAMP)
		`, f.base, quote, rate, date)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("insert rate: %w", err)
		}
	}

	return tx.Commit()
}

func (f *Fetcher) GetRates(base, date string) (map[string]float64, string, error) {
	if base == "" {
		base = f.base
	}
	if date == "" || date == "latest" {
		date = time.Now().Format("2006-01-02")
	}

	query := `SELECT quote, rate FROM fx_rates WHERE base = ? AND date = ?`
	rows, err := f.db.Query(query, base, date)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	rates := make(map[string]float64)
	for rows.Next() {
		var quote string
		var rate float64
		rows.Scan(&quote, &rate)
		rates[quote] = rate
	}

	if len(rates) == 0 {
		// Try to fetch and store
		if err := f.FetchAndStore(); err != nil {
			log.Warn().Err(err).Msg("fx: fetch on demand failed")
		}
		// Retry query
		rows, err := f.db.Query(query, base, date)
		if err != nil {
			return nil, "", err
		}
		defer rows.Close()
		rates = make(map[string]float64)
		for rows.Next() {
			var quote string
			var rate float64
			rows.Scan(&quote, &rate)
			rates[quote] = rate
		}
	}

	rates[base] = 1.0
	return rates, date, nil
}

func (f *Fetcher) Convert(from, to string, amount float64) (float64, error) {
	rates, _, err := f.GetRates(from, "latest")
	if err != nil {
		return 0, err
	}
	fromRate, ok1 := rates[from]
	toRate, ok2 := rates[to]
	if !ok1 || !ok2 {
		return 0, fmt.Errorf("missing rate for %s or %s", from, to)
	}
	return amount * toRate / fromRate, nil
}

func StartScheduler(database *db.DB, baseCurrency string, schedule string) {
	fetcher := NewFetcher(database, baseCurrency)

	// Initial fetch if table empty
	var count int
	database.QueryRow("SELECT COUNT(*) FROM fx_rates").Scan(&count)
	if count == 0 {
		log.Info().Msg("fx: table empty, fetching initial rates")
		if err := fetcher.FetchAndStore(); err != nil {
			log.Warn().Err(err).Msg("fx: initial fetch failed")
		}
	}

	// Parse cron schedule
	// For simplicity, we'll run at 04:00 daily
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), 4, 0, 0, 0, now.Location())
			if next.Before(now) {
				next = next.AddDate(0, 0, 1)
			}
			time.Sleep(time.Until(next))

			log.Info().Msg("fx: scheduled fetch started")
			if err := fetcher.FetchAndStore(); err != nil {
				log.Warn().Err(err).Msg("fx: scheduled fetch failed")
			}
		}
	}()
}
