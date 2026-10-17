package analytics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/models"
)

func (h *Handler) CashFlowProjection(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	dateFromStr := r.URL.Query().Get("date_from")
	dateToStr := r.URL.Query().Get("date_to")
	accountID := r.URL.Query().Get("account_id")
	currency := r.URL.Query().Get("currency")

	if dateFromStr == "" || dateToStr == "" {
		http.Error(w, "date_from and date_to required", http.StatusBadRequest)
		return
	}

	dateFrom, err := time.Parse("2006-01-02", dateFromStr)
	if err != nil {
		http.Error(w, "invalid date_from format", http.StatusBadRequest)
		return
	}
	dateTo, err := time.Parse("2006-01-02", dateToStr)
	if err != nil {
		http.Error(w, "invalid date_to format", http.StatusBadRequest)
		return
	}

	if currency == "" {
		currency = h.cfg.FXBaseCurrency
		if currency == "" {
			currency = "JPY"
		}
	}

	projection, err := h.calculateProjection(userID, dateFrom, dateTo, accountID, currency)
	if err != nil {
		log.Error().Err(err).Msg("analytics: cashflow projection")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, projection)
}

func (h *Handler) Calendar(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")
	accountID := r.URL.Query().Get("account_id")

	if yearStr == "" || monthStr == "" {
		http.Error(w, "year and month required", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "invalid year", http.StatusBadRequest)
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		http.Error(w, "invalid month", http.StatusBadRequest)
		return
	}

	currency := h.cfg.FXBaseCurrency
	if currency == "" {
		currency = "JPY"
	}

	dateFrom := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	dateTo := dateFrom.AddDate(0, 1, -1)

	calendar, err := h.calculateCalendar(userID, dateFrom, dateTo, accountID, currency)
	if err != nil {
		log.Error().Err(err).Msg("analytics: calendar")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, calendar)
}

func (h *Handler) calculateProjection(userID string, dateFrom, dateTo time.Time, accountID, baseCurrency string) (*models.CashFlowProjectionResponse, error) {
	accounts, err := h.getAccountBalances(userID, accountID, baseCurrency)
	if err != nil {
		return nil, err
	}

	scheduledItems, err := h.getScheduledTransactions(userID, dateFrom, dateTo, accountID, baseCurrency)
	if err != nil {
		return nil, err
	}

	futureItems, err := h.getFutureTransactions(userID, dateFrom, dateTo, accountID, baseCurrency)
	if err != nil {
		return nil, err
	}

	dailyFlow := make(map[string]struct{ income, expense float64 })

	for _, item := range scheduledItems {
		amt := item.Amount
		if item.Type == "income" {
			dailyFlow[item.Date] = struct{ income, expense float64 }{
				income:  dailyFlow[item.Date].income + amt,
				expense: dailyFlow[item.Date].expense,
			}
		} else {
			dailyFlow[item.Date] = struct{ income, expense float64 }{
				income:  dailyFlow[item.Date].income,
				expense: dailyFlow[item.Date].expense + amt,
			}
		}
	}

	for _, item := range futureItems {
		amt := item.Amount
		if item.Type == "income" {
			dailyFlow[item.Date] = struct{ income, expense float64 }{
				income:  dailyFlow[item.Date].income + amt,
				expense: dailyFlow[item.Date].expense,
			}
		} else {
			dailyFlow[item.Date] = struct{ income, expense float64 }{
				income:  dailyFlow[item.Date].income,
				expense: dailyFlow[item.Date].expense + amt,
			}
		}
	}

	var runningBalance float64
	for _, acc := range accounts {
		runningBalance += acc.Balance
	}

	var projected []models.ProjectionPoint
	current := dateFrom
	for !current.After(dateTo) {
		dateStr := current.Format("2006-01-02")
		flow := dailyFlow[dateStr]
		net := flow.income - flow.expense
		runningBalance += net

		projected = append(projected, models.ProjectionPoint{
			Date:           dateStr,
			Income:         flow.income,
			Expense:        flow.expense,
			Net:            net,
			RunningBalance: runningBalance,
			Currency:       baseCurrency,
		})
		current = current.AddDate(0, 0, 1)
	}

	return &models.CashFlowProjectionResponse{
		ProjectedBalance:      projected,
		ScheduledTransactions: scheduledItems,
		FutureTransactions:    futureItems,
	}, nil
}

func (h *Handler) calculateCalendar(userID string, dateFrom, dateTo time.Time, accountID, baseCurrency string) (*models.CalendarResponse, error) {
	accounts, err := h.getAccountBalances(userID, accountID, baseCurrency)
	if err != nil {
		return nil, err
	}

	scheduledItems, err := h.getScheduledTransactions(userID, dateFrom, dateTo, accountID, baseCurrency)
	if err != nil {
		return nil, err
	}

	futureItems, err := h.getFutureTransactions(userID, dateFrom, dateTo, accountID, baseCurrency)
	if err != nil {
		return nil, err
	}

	historicalItems, err := h.getHistoricalTransactions(userID, dateFrom, dateTo, accountID, baseCurrency)
	if err != nil {
		return nil, err
	}

	dayMap := make(map[string]*models.CalendarDay)
	current := dateFrom
	for !current.After(dateTo) {
		dateStr := current.Format("2006-01-02")
		dayMap[dateStr] = &models.CalendarDay{
			Date:             dateStr,
			Transactions:     []models.CalendarTxn{},
			Scheduled:        []models.CalendarTxn{},
			BalanceDelta:     0,
			ProjectedBalance: 0,
		}
		current = current.AddDate(0, 0, 1)
	}

	var runningBalance float64
	for _, acc := range accounts {
		runningBalance += acc.Balance
	}

	for _, item := range historicalItems {
		if day, ok := dayMap[item.Date]; ok {
			day.Transactions = append(day.Transactions, models.CalendarTxn{
				ID:       item.ID,
				Name:     item.Name,
				Amount:   item.Amount,
				Type:     item.Type,
				Category: item.Category,
				Currency: item.Currency,
			})
			if item.Type == "income" {
				day.BalanceDelta += item.Amount
			} else {
				day.BalanceDelta -= item.Amount
			}
		}
	}

	for _, item := range scheduledItems {
		if day, ok := dayMap[item.Date]; ok {
			day.Scheduled = append(day.Scheduled, models.CalendarTxn{
				ID:       item.ID,
				Name:     item.Name,
				Amount:   item.Amount,
				Type:     item.Type,
				Category: item.Category,
				Currency: item.Currency,
			})
			if item.Type == "income" {
				day.BalanceDelta += item.Amount
			} else {
				day.BalanceDelta -= item.Amount
			}
		}
	}

	for _, item := range futureItems {
		if day, ok := dayMap[item.Date]; ok {
			day.Transactions = append(day.Transactions, models.CalendarTxn{
				ID:       item.ID,
				Name:     item.Name,
				Amount:   item.Amount,
				Type:     item.Type,
				Category: item.Category,
				Currency: item.Currency,
			})
			if item.Type == "income" {
				day.BalanceDelta += item.Amount
			} else {
				day.BalanceDelta -= item.Amount
			}
		}
	}

	dates := make([]string, 0, len(dayMap))
	for d := range dayMap {
		dates = append(dates, d)
	}

	for _, dateStr := range dates {
		day := dayMap[dateStr]
		runningBalance += day.BalanceDelta
		day.ProjectedBalance = runningBalance
	}

	days := make([]models.CalendarDay, 0, len(dayMap))
	for _, dateStr := range dates {
		days = append(days, *dayMap[dateStr])
	}

	return &models.CalendarResponse{Days: days}, nil
}
