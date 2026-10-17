package analytics

import (
	"github.com/go-chi/chi/v5"

	"kiroku/backend/internal/db"
	"kiroku/backend/internal/models"
)

type Handler struct {
	db  *db.DB
	cfg *models.Config
}

func NewHandler(db *db.DB, baseCurrency string) *Handler {
	return &Handler{
		db: db,
		cfg: &models.Config{
			FXBaseCurrency: baseCurrency,
		},
	}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/cashflow/projection", h.CashFlowProjection)
	r.Get("/calendar", h.Calendar)

	r.Get("/balance/trend", h.GetBalanceTrend)
	r.Get("/balance/accounts", h.GetBalanceByAccounts)
	r.Get("/balance/currencies", h.GetBalanceByCurrencies)

	r.Get("/outlook/category", h.GetPlannedByCategory)
	r.Get("/outlook/timeline", h.GetPlannedByTimeline)
	r.Get("/outlook/forecast", h.GetBalanceForecast)

	r.Get("/budget/summary", h.GetBudgetSummary)
	r.Get("/budget/vs-actual", h.GetBudgetVsActual)

	r.Get("/goals/progress", h.GetGoalsProgress)
	r.Get("/goals/projection", h.GetGoalProjection)

	r.Get("/stats/summary", h.GetStatsSummary)
	r.Get("/stats/by-category", h.GetStatsByCategory)
	r.Get("/stats/by-currency", h.GetStatsByCurrency)
	r.Get("/stats/by-account", h.GetStatsByAccount)
	r.Get("/stats/top-merchants", h.GetTopMerchants)
	r.Get("/stats/cashflow", h.GetCashFlowStats)
	r.Get("/stats/net-worth", h.GetNetWorth)

	return r
}
