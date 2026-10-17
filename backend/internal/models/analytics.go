package models

type ProjectionPoint struct {
	Date           string  `json:"date"`
	Income         float64 `json:"income"`
	Expense        float64 `json:"expense"`
	Net            float64 `json:"net"`
	RunningBalance float64 `json:"running_balance"`
	Currency       string  `json:"currency"`
}

type ScheduledItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	Category    string  `json:"category"`
	AccountID   string  `json:"account_id"`
	AccountName string  `json:"account_name"`
	Date        string  `json:"date"`
	Currency    string  `json:"currency"`
	Recurring   bool    `json:"recurring"`
}

type CashFlowProjectionResponse struct {
	ProjectedBalance      []ProjectionPoint `json:"projected_balance"`
	ScheduledTransactions []ScheduledItem   `json:"scheduled_transactions"`
	FutureTransactions    []ScheduledItem   `json:"future_transactions"`
}

type CalendarDay struct {
	Date             string        `json:"date"`
	Transactions     []CalendarTxn `json:"transactions"`
	Scheduled        []CalendarTxn `json:"scheduled"`
	BalanceDelta     float64       `json:"balance_delta"`
	ProjectedBalance float64       `json:"projected_balance"`
}

type CalendarTxn struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	Type     string  `json:"type"`
	Category string  `json:"category"`
	Currency string  `json:"currency"`
}

type CalendarResponse struct {
	Days []CalendarDay `json:"days"`
}

type AccountBalance struct {
	ID       string
	Name     string
	Balance  float64
	Currency string
}

type BalanceTrendPoint struct {
	Period         string  `json:"period"`
	OpeningBalance float64 `json:"opening_balance"`
	ClosingBalance float64 `json:"closing_balance"`
	NetChange      float64 `json:"net_change"`
	Currency       string  `json:"currency"`
}

type AccountBalanceSummary struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Currency       string  `json:"currency"`
	OpeningBalance float64 `json:"opening_balance"`
	ClosingBalance float64 `json:"closing_balance"`
	NetChange      float64 `json:"net_change"`
	Color          string  `json:"color"`
}

type CurrencyBalanceSummary struct {
	Currency       string  `json:"currency"`
	OpeningBalance float64 `json:"opening_balance"`
	ClosingBalance float64 `json:"closing_balance"`
	NetChange      float64 `json:"net_change"`
	FXRateToBase   float64 `json:"fx_rate_to_base"`
	ConvertedTotal float64 `json:"converted_total"`
}

type BudgetItem struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Icon         string  `json:"icon"`
	Color        string  `json:"color"`
	BudgetAmount float64 `json:"budget_amount"`
	BudgetPeriod string  `json:"budget_period"`
	Spent        float64 `json:"spent"`
	Remaining    float64 `json:"remaining"`
	ProgressPct  float64 `json:"progress_pct"`
	Status       string  `json:"status"`
}

type BudgetVsActual struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Icon         string  `json:"icon"`
	Color        string  `json:"color"`
	BudgetAmount float64 `json:"budget_amount"`
	BudgetPeriod string  `json:"budget_period"`
	ActualSpent  float64 `json:"actual_spent"`
	Variance     float64 `json:"variance"`
	VariancePct  float64 `json:"variance_pct"`
}

type GoalItem struct {
	CategoryID   string
	CategoryName string
	Icon         string
	Color        string
	TargetAmount float64
	TargetPeriod string
}

type GoalProgress struct {
	CategoryID          string  `json:"category_id"`
	CategoryName        string  `json:"category_name"`
	Icon                string  `json:"icon"`
	Color               string  `json:"color"`
	TargetAmount        float64 `json:"target_amount"`
	TargetPeriod        string  `json:"target_period"`
	CurrentAmount       float64 `json:"current_amount"`
	ProgressPct         float64 `json:"progress_pct"`
	Remaining           float64 `json:"remaining"`
	ProjectedCompletion string  `json:"projected_completion"`
	Status              string  `json:"status"`
}

type GoalProjectionPoint struct {
	Period           string  `json:"period"`
	ActualCumulative float64 `json:"actual_cumulative"`
	TargetCumulative float64 `json:"target_cumulative"`
	TargetTotal      float64 `json:"target_total"`
}

type Config struct {
	FXBaseCurrency string
}

type CategoryStats struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Icon         string  `json:"icon"`
	Color        string  `json:"color"`
	TotalAmount  float64 `json:"total_amount"`
	Count        int     `json:"count"`
	Percentage   float64 `json:"percentage"`
}

type CurrencyStats struct {
	Currency string  `json:"currency"`
	Income   float64 `json:"income"`
	Expense  float64 `json:"expense"`
	Net      float64 `json:"net"`
	Count    int     `json:"count"`
}

type AccountStats struct {
	AccountID   string  `json:"account_id"`
	AccountName string  `json:"account_name"`
	Currency    string  `json:"currency"`
	Income      float64 `json:"income"`
	Expense     float64 `json:"expense"`
	Net         float64 `json:"net"`
	Count       int     `json:"count"`
}

type MerchantStats struct {
	Merchant    string  `json:"merchant"`
	TotalAmount float64 `json:"total_amount"`
	Count       int     `json:"count"`
}

type CashFlowSeriesPoint struct {
	Period  string  `json:"period"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Net     float64 `json:"net"`
}

type NetWorthAccount struct {
	AccountID string  `json:"account_id"`
	Name      string  `json:"name"`
	Currency  string  `json:"currency"`
	Balance   float64 `json:"balance"`
	Type      string  `json:"type"`
}

type TimelinePoint struct {
	Period  string  `json:"period"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Net     float64 `json:"net"`
}

type ForecastPoint struct {
	Date             string  `json:"date"`
	ProjectedBalance float64 `json:"projected_balance"`
	Delta            float64 `json:"delta"`
}

type CategoryPlanned struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Icon         string  `json:"icon"`
	Color        string  `json:"color"`
	TotalAmount  float64 `json:"total_amount"`
}
