package domain

import "time"

type ProjectionPoint struct {
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
}

type ScheduledItem struct {
	Date     time.Time `json:"date"`
	Amount   float64   `json:"amount"`
	Currency string    `json:"currency"`
	Type     string    `json:"type"`
}

type CalendarDay struct {
	Date        string            `json:"date"`
	Income      float64           `json:"income"`
	Expense     float64           `json:"expense"`
	Net         float64           `json:"net"`
	Transactions []CalendarTxn    `json:"transactions"`
}

type CalendarTxn struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	Type     string  `json:"type"`
	Category string  `json:"category"`
	Currency string  `json:"currency"`
}

type AccountBalance struct {
	AccountID   string  `json:"account_id"`
	AccountName string  `json:"account_name"`
	Currency    string  `json:"currency"`
	Balance     float64 `json:"balance"`
}

type BalanceTrendPoint struct {
	Date    time.Time `json:"date"`
	Balance float64   `json:"balance"`
}

type BudgetItem struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Currency     string  `json:"currency"`
	Planned      float64 `json:"planned"`
	Actual       float64 `json:"actual"`
	Remaining    float64 `json:"remaining"`
}

type GoalItem struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	TargetAmount float64 `json:"target_amount"`
	CurrentAmount float64 `json:"current_amount"`
	Currency     string  `json:"currency"`
}

type CategoryStats struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Currency     string  `json:"currency"`
	Amount       float64 `json:"amount"`
	Count        int     `json:"count"`
}

type CurrencyStats struct {
	Currency string  `json:"currency"`
	Income   float64 `json:"income"`
	Expense  float64 `json:"expense"`
	Net      float64 `json:"net"`
}

type AccountStats struct {
	AccountID   string  `json:"account_id"`
	AccountName string  `json:"account_name"`
	Currency    string  `json:"currency"`
	Balance     float64 `json:"balance"`
	Income      float64 `json:"income"`
	Expense     float64 `json:"expense"`
}

type MerchantStats struct {
	Merchant string  `json:"merchant"`
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
	Count    int     `json:"count"`
}

type CashFlowSeriesPoint struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Net     float64 `json:"net"`
}

type NetWorthAccount struct {
	AccountID   string  `json:"account_id"`
	AccountName string  `json:"account_name"`
	Currency    string  `json:"currency"`
	Balance     float64 `json:"balance"`
}

type TimelinePoint struct {
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
}

type ForecastPoint struct {
	Date     time.Time `json:"date"`
	Predicted float64  `json:"predicted"`
	Lower    float64   `json:"lower"`
	Upper    float64   `json:"upper"`
}

type CategoryPlanned struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Planned      float64 `json:"planned"`
	Actual       float64 `json:"actual"`
}

type AnalyticsFilter struct {
	DateFrom *time.Time
	DateTo   *time.Time
	Currency string
	Period   string
	Months   int
	AccountID *string
	CategoryID *string
}
