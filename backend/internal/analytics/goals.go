package analytics

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/models"
)

func (h *Handler) GetGoalsProgress(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")

	if dateFrom == "" || dateTo == "" {
		dateFrom = time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
		dateTo = time.Now().Format("2006-01-02")
	}

	baseCurrency := h.cfg.FXBaseCurrency
	if baseCurrency == "" {
		baseCurrency = "JPY"
	}

	query := `
		SELECT c.id, c.name, c.icon, c.color, c.budget_amount, c.budget_period
		FROM categories c
		WHERE c.user_id = ? AND c.budget_amount > 0 AND c.budget_period != ''
	`
	rows, err := h.db.Query(query, userID)
	if err != nil {
		log.Error().Err(err).Msg("analytics: goals categories")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []models.GoalProgress
	var goalsAchieved, goalsOnTrack, goalsBehind int
	var totalTargets, totalCurrent float64

	for rows.Next() {
		var item models.GoalProgress
		if err := rows.Scan(&item.CategoryID, &item.CategoryName, &item.Icon, &item.Color, &item.TargetAmount, &item.TargetPeriod); err != nil {
			log.Error().Err(err).Msg("analytics: goals scan")
			continue
		}

		var currentAmount float64
		h.db.QueryRow(`
			SELECT COALESCE(SUM(t.amount * COALESCE(r.rate, 1)), 0)
			FROM transactions t
			LEFT JOIN fx_rates r ON t.currency = r.base AND r.quote = ?
			WHERE t.user_id = ? AND t.category_id = ? AND t.type = 'expense'
				AND t.date >= ? AND t.date <= ? AND t.deleted_at IS NULL
		`, baseCurrency, userID, item.CategoryID, dateFrom, dateTo).Scan(&currentAmount)

		item.CurrentAmount = currentAmount
		item.Remaining = item.TargetAmount - currentAmount
		item.ProgressPct = safeDiv(currentAmount, item.TargetAmount) * 100

		if item.ProgressPct >= 100 {
			item.Status = "achieved"
			goalsAchieved++
		} else {
			daysElapsed := time.Since(time.Now().AddDate(-1, 0, 0)).Hours() / 24
			daysTotal := 365.0
			expectedProgress := daysElapsed / daysTotal
			dailyRate := item.TargetAmount / daysTotal
			if dailyRate > 0 && currentAmount < item.TargetAmount {
				remaining := item.TargetAmount - currentAmount
				daysToTarget := remaining / dailyRate
				item.ProjectedCompletion = time.Now().AddDate(0, 0, int(daysToTarget)).Format("2006-01-02")
			}
			if item.ProgressPct >= expectedProgress*100 {
				item.Status = "on_track"
				goalsOnTrack++
			} else {
				item.Status = "behind"
				goalsBehind++
			}
		}

		results = append(results, item)
		totalTargets += item.TargetAmount
		totalCurrent += currentAmount
	}

	overallProgress := safeDiv(totalCurrent, totalTargets) * 100

	writeJSON(w, map[string]interface{}{
		"goals": results,
		"summary": map[string]interface{}{
			"total_targets":    totalTargets,
			"total_current":    totalCurrent,
			"overall_progress": overallProgress,
			"goals_achieved":   goalsAchieved,
			"goals_on_track":   goalsOnTrack,
			"goals_behind":     goalsBehind,
		},
	})
}

func (h *Handler) GetGoalProjection(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	categoryID := r.URL.Query().Get("category_id")
	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")

	if categoryID == "" {
		http.Error(w, "category_id required", http.StatusBadRequest)
		return
	}

	if dateFrom == "" || dateTo == "" {
		dateFrom = time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
		dateTo = time.Now().Format("2006-01-02")
	}

	start, _ := time.Parse("2006-01-02", dateFrom)
	end, _ := time.Parse("2006-01-02", dateTo)

	var targetAmount float64
	var targetPeriod string
	var categoryName, icon, color string
	err := h.db.QueryRow(`
		SELECT name, icon, color, budget_amount, budget_period
		FROM categories
		WHERE id = ?
	`, categoryID).Scan(&categoryName, &icon, &color, &targetAmount, &targetPeriod)
	if err != nil || targetAmount == 0 {
		http.Error(w, "goal not found", http.StatusNotFound)
		return
	}

	var projection []models.GoalProjectionPoint
	currStart := start
	var cumulativeActual float64

	for currStart.Before(end) {
		currEnd := currStart.AddDate(0, 1, 0)
		if currEnd.After(end) {
			currEnd = end
		}

		var actualAmount float64
		h.db.QueryRow(`
			SELECT COALESCE(SUM(amount), 0)
			FROM transactions
			WHERE user_id = ? AND category_id = ? AND type = 'expense'
			AND date >= ? AND date < ? AND deleted_at IS NULL
		`, userID, categoryID, currStart.Format("2006-01-02"), currEnd.Format("2006-01-02")).Scan(&actualAmount)

		cumulativeActual += actualAmount

		daysElapsed := currEnd.Sub(start).Hours() / 24
		daysTotal := end.Sub(start).Hours() / 24
		expectedProgress := daysElapsed / daysTotal
		expectedAmount := targetAmount * expectedProgress

		projection = append(projection, models.GoalProjectionPoint{
			Period:           currStart.Format("2006-01"),
			ActualCumulative: cumulativeActual,
			TargetCumulative: expectedAmount,
			TargetTotal:      targetAmount,
		})

		currStart = currEnd
	}

	writeJSON(w, map[string]interface{}{
		"category_id":   categoryID,
		"category_name": categoryName,
		"icon":          icon,
		"color":         color,
		"target_amount": targetAmount,
		"target_period": targetPeriod,
		"projection":    projection,
	})
}
