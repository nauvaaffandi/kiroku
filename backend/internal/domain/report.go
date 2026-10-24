package domain

import "time"

type SpendingGroup struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Color    string  `json:"color"`
	Kind     *string `json:"kind"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Count    int     `json:"count"`
}

type TrendPoint struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Net     float64 `json:"net"`
}

type SpendingFilter struct {
	Period    string
	GroupBy   string
	DateFrom  *time.Time
	DateTo    *time.Time
	Currency  string
	AccountID *string
}

type TrendFilter struct {
	Months   int
	Currency string
	DateFrom *time.Time
	DateTo   *time.Time
	AccountID *string
}
