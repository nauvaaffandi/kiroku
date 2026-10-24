package domain

import (
	"encoding/json"
	"time"
)

type ScheduledTransaction struct {
	ID                   string          `json:"id"`
	RRule                string          `json:"rrule"`
	MaxOccurrences       *int            `json:"max_occurrences,omitempty"`
	TransactionTemplate  json.RawMessage `json:"transaction_template"`
	Enabled              bool            `json:"enabled"`
	NextRun              *time.Time      `json:"next_run,omitempty"`
	OccurrencesDone      int             `json:"occurrences_done"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

type CreateScheduledInput struct {
	RRule              string
	MaxOccurrences     *int
	TransactionTemplate map[string]interface{}
	Enabled            bool
}

type UpdateScheduledInput struct {
	RRule              *string
	MaxOccurrences     *int
	TransactionTemplate map[string]interface{}
	Enabled            *bool
}

type PresetParams struct {
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	FromAccountID *string `json:"from_account_id"`
	ToAccountID   *string `json:"to_account_id"`
	CategoryID    *string `json:"category_id"`
	Labels        []string `json:"labels"`
	Note          string  `json:"note"`
	PayeePayer    string  `json:"payee_payer"`
	PaymentType   string  `json:"payment_type"`
	Status        string  `json:"status"`
	Frequency     string  `json:"frequency"`
	Interval      int     `json:"interval"`
	DayOfMonth    int     `json:"day_of_month"`
	DayOfWeek     int     `json:"day_of_week"`
}
