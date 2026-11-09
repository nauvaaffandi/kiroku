package domain

import (
	"context"
	"fmt"
	"math"
	"time"
)

const BudgetPeriodWeekly = "weekly"
const BudgetPeriodMonthly = "monthly"
const BudgetPeriodYearly = "yearly"
const BudgetPeriodOneTime = "one_time"

const (
	BudgetPeriodApproachingEvent = "budget_approaching"
	BudgetPeriodOverEvent        = "budget_over"
)

func IsValidBudgetPeriod(period string) bool {
	switch period {
	case BudgetPeriodWeekly, BudgetPeriodMonthly, BudgetPeriodYearly, BudgetPeriodOneTime:
		return true
	}
	return false
}

func FormatCents(cents int64) string {
	sign := ""
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	return fmt.Sprintf("%s%d.%02d", sign, cents/100, cents%100)
}

func Cents(amount float64) int64 { return int64(math.Round(amount * 100)) }

type Budget struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	Name         string     `json:"name"`
	Period       string     `json:"period"`
	Amount       int64      `json:"-"`
	AmountStr    string     `json:"amount"`
	Currency     string     `json:"currency"`
	AccountID    *string    `json:"account_id,omitempty"`
	Categories   []Category `json:"categories,omitempty"`
	Labels       []Label    `json:"labels,omitempty"`
	TrendingOver bool       `json:"trending_over"`
	OverBudget   bool       `json:"over_budget"`
	StartsAt     *time.Time `json:"starts_at,omitempty"`
	EndsAt       *time.Time `json:"ends_at,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type BudgetResult struct {
	BudgetID       string    `json:"budget_id"`
	Name           string    `json:"name"`
	Period         string    `json:"period"`
	Amount         string    `json:"amount"`
	Currency       string    `json:"currency"`
	Spent          string    `json:"spent"`
	Remaining      string    `json:"remaining"`
	PercentageUsed int       `json:"percentage_used"`
	IsOverBudget   bool      `json:"is_over_budget"`
	PeriodStart    time.Time `json:"period_start"`
	PeriodEnd      time.Time `json:"period_end"`
}

type CreateBudgetInput struct {
	Name         string   `json:"name"`
	Period       string   `json:"period"`
	Amount       float64  `json:"amount"`
	Currency     string   `json:"currency"`
	AccountID    *string  `json:"account_id,omitempty"`
	CategoryIDs  []string `json:"category_ids"`
	LabelIDs     []string `json:"label_ids"`
	TrendingOver bool     `json:"trending_over"`
	OverBudget   bool     `json:"over_budget"`
	StartsAt     *string  `json:"starts_at,omitempty"`
	EndsAt       *string  `json:"ends_at,omitempty"`
}

type UpdateBudgetInput struct {
	Name         *string  `json:"name,omitempty"`
	Period       *string  `json:"period,omitempty"`
	Amount       *float64 `json:"amount,omitempty"`
	Currency     *string  `json:"currency,omitempty"`
	AccountID    *string  `json:"account_id,omitempty"`
	CategoryIDs  []string `json:"category_ids,omitempty"`
	LabelIDs     []string `json:"label_ids,omitempty"`
	TrendingOver *bool    `json:"trending_over,omitempty"`
	OverBudget   *bool    `json:"over_budget,omitempty"`
	StartsAt     *string  `json:"starts_at,omitempty"`
	EndsAt       *string  `json:"ends_at,omitempty"`
}

type BudgetRepository interface {
	List(ctx context.Context, userID string) ([]Budget, error)
	Get(ctx context.Context, id, userID string) (*Budget, error)
	Create(ctx context.Context, input CreateBudgetInput, userID string) (*Budget, error)
	Update(ctx context.Context, id, userID string, input UpdateBudgetInput) (*Budget, error)
	Delete(ctx context.Context, id, userID string) error
	Calculate(ctx context.Context, budget *Budget) (*BudgetResult, error)
	UserIDsWithBudgets(ctx context.Context) ([]string, error)
}
