package budgets

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/domain"
	"kiroku/backend/internal/notification"
)

const (
	approachingThreshold = 80
	overBudgetThreshold  = 100
)

type notifier interface {
	Create(ctx context.Context, input domain.CreateNotificationInput) (*domain.Notification, error)
}

// Service ties budget calculation to the existing domain-agnostic notification
// package. It is wired into the transaction lifecycle so committed transaction
// changes trigger budget recalculation and threshold notifications.
type Service struct {
	repo     domain.BudgetRepository
	notifier notifier
}

func NewService(repo domain.BudgetRepository, notifier notifier) *Service {
	return &Service{repo: repo, notifier: notifier}
}

// Check recalculates every active budget for the user and creates threshold
// notifications when a budget crosses the approaching or over-budget limits.
// Notification creation is idempotent per budget/period/event via the
// notification package.
func (s *Service) Check(ctx context.Context, userID string) error {
	if userID == "" {
		return nil
	}

	items, err := s.repo.List(ctx, userID)
	if err != nil {
		return err
	}

	for _, item := range items {
		result, err := s.repo.Calculate(ctx, &item)
		if err != nil {
			return err
		}

		if item.TrendingOver && result.PercentageUsed >= approachingThreshold {
			if err := s.notify(ctx, userID, item.ID, result, domain.BudgetPeriodApproachingEvent, approachingThreshold); err != nil {
				return err
			}
		}
		if item.OverBudget && result.PercentageUsed >= overBudgetThreshold {
			if err := s.notify(ctx, userID, item.ID, result, domain.BudgetPeriodOverEvent, overBudgetThreshold); err != nil {
				return err
			}
		}
	}
	return nil
}

// ReconcileAll periodically recalculates all budgets for users who have at
// least one budget and evaluates threshold notifications. It is a lightweight
// fallback for missed post-commit budget checks and relies on the notification
// idempotency keys to avoid duplicate alerts.
func (s *Service) ReconcileAll(ctx context.Context) error {
	users, err := s.repo.UserIDsWithBudgets(ctx)
	if err != nil {
		return err
	}
	for _, userID := range users {
		if err := s.Check(ctx, userID); err != nil {
			log.Warn().Err(err).Str("user_id", userID).Msg("budget: reconcile check failed")
		}
	}
	return nil
}

func (s *Service) notify(ctx context.Context, userID, budgetID string, result *domain.BudgetResult, event string, threshold int) error {
	periodStart := result.PeriodStart.Format("2006-01-02")
	key := fmt.Sprintf("budget:%s:%s:%s", budgetID, periodStart, event)

	var title, body string
	if event == domain.BudgetPeriodOverEvent {
		title = "Budget exceeded"
		body = fmt.Sprintf("Your %s budget has exceeded its %s spending limit (%s of %s %s).", result.Name, result.Period, result.Spent, result.Currency, result.Amount)
	} else {
		title = "Budget almost reached"
		body = fmt.Sprintf("Your %s budget has used %d%% of its %s limit (%s of %s %s).", result.Name, result.PercentageUsed, result.Period, result.Spent, result.Currency, result.Amount)
	}

	_, err := s.notifier.Create(ctx, domain.CreateNotificationInput{
		RecipientID:    userID,
		Type:           event,
		Title:          title,
		Body:           body,
		ActionURL:      "/budgets/" + budgetID,
		IdempotencyKey: key,
		Data: map[string]any{
			"budget_id":   budgetID,
			"event":       event,
			"period_start": periodStart,
			"period":      result.Period,
			"spent":       result.Spent,
			"amount":      result.Amount,
			"currency":    result.Currency,
			"percentage":  result.PercentageUsed,
		},
	})
	return err
}

var _ = time.Time{}
var _ = notification.Service{}
