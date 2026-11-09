package httpdelivery

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/domain"
)

type BudgetsHandler struct{ repo domain.BudgetRepository }

func NewBudgetsHandler(repo domain.BudgetRepository) *BudgetsHandler {
	return &BudgetsHandler{repo: repo}
}

func (h *BudgetsHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}/consumption", h.Consumption)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

func (h *BudgetsHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.List(r.Context(), auth.GetUserID(r))
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	results := make([]domain.BudgetResult, 0, len(items))
	for i := range items {
		result, err := h.repo.Calculate(r.Context(), &items[i])
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		results = append(results, *result)
	}
	writeBudgetJSON(w, http.StatusOK, results)
}

func (h *BudgetsHandler) Get(w http.ResponseWriter, r *http.Request) {
	b, err := h.repo.Get(r.Context(), chi.URLParam(r, "id"), auth.GetUserID(r))
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if b == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	result, err := h.repo.Calculate(r.Context(), b)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeBudgetJSON(w, http.StatusOK, map[string]any{"budget": b, "consumption": result})
}

func (h *BudgetsHandler) Consumption(w http.ResponseWriter, r *http.Request) {
	b, err := h.repo.Get(r.Context(), chi.URLParam(r, "id"), auth.GetUserID(r))
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if b == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	result, err := h.repo.Calculate(r.Context(), b)
	if err != nil {
		http.Error(w, "unable to calculate budget consumption", http.StatusUnprocessableEntity)
		return
	}
	writeBudgetJSON(w, http.StatusOK, result)
}

func (h *BudgetsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateBudgetInput
	if err := decodeBudgetRequest(r, &input); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := validateCreateBudget(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := h.repo.Create(r.Context(), input, auth.GetUserID(r))
	if err != nil {
		http.Error(w, "invalid budget", http.StatusBadRequest)
		return
	}
	writeBudgetJSON(w, http.StatusCreated, b)
}

func (h *BudgetsHandler) Update(w http.ResponseWriter, r *http.Request) {
	var input domain.UpdateBudgetInput
	if err := decodeBudgetRequest(r, &input); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := validateUpdateBudget(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := h.repo.Update(r.Context(), chi.URLParam(r, "id"), auth.GetUserID(r), input)
	if err != nil {
		http.Error(w, "invalid budget", http.StatusBadRequest)
		return
	}
	if b == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeBudgetJSON(w, http.StatusOK, b)
}

func (h *BudgetsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.Delete(r.Context(), chi.URLParam(r, "id"), auth.GetUserID(r)); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

var budgetCurrencyPattern = regexp.MustCompile(`^[A-Z]{3}$`)

func validateCreateBudget(input domain.CreateBudgetInput) error {
	if strings.TrimSpace(input.Name) == "" || len(input.Name) > 255 {
		return errors.New("name must be 1-255 characters")
	}
	if math.IsNaN(input.Amount) || math.IsInf(input.Amount, 0) || input.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if !budgetCurrencyPattern.MatchString(input.Currency) {
		return errors.New("currency must be three uppercase letters")
	}
	if !domain.IsValidBudgetPeriod(input.Period) {
		return errors.New("invalid period")
	}
	if err := validateIDs(input.CategoryIDs); err != nil {
		return errors.New("duplicate category ids")
	}
	if err := validateIDs(input.LabelIDs); err != nil {
		return errors.New("duplicate label ids")
	}
	if input.Period == domain.BudgetPeriodOneTime {
		if _, _, err := oneTimeWindow(input.StartsAt, input.EndsAt); err != nil {
			return err
		}
	}
	return nil
}

func validateUpdateBudget(input domain.UpdateBudgetInput) error {
	if input.Name != nil && (strings.TrimSpace(*input.Name) == "" || len(*input.Name) > 255) {
		return errors.New("name must be 1-255 characters")
	}
	if input.Amount != nil && (math.IsNaN(*input.Amount) || math.IsInf(*input.Amount, 0) || *input.Amount <= 0) {
		return errors.New("amount must be greater than zero")
	}
	if input.Currency != nil && !budgetCurrencyPattern.MatchString(*input.Currency) {
		return errors.New("currency must be three uppercase letters")
	}
	if input.Period != nil && !domain.IsValidBudgetPeriod(*input.Period) {
		return errors.New("invalid period")
	}
	if input.CategoryIDs != nil && validateIDs(input.CategoryIDs) != nil {
		return errors.New("duplicate category ids")
	}
	if input.LabelIDs != nil && validateIDs(input.LabelIDs) != nil {
		return errors.New("duplicate label ids")
	}
	if input.Period != nil && *input.Period == domain.BudgetPeriodOneTime {
		if _, _, err := oneTimeWindow(input.StartsAt, input.EndsAt); err != nil {
			return err
		}
	}
	return nil
}

func oneTimeWindow(start, end *string) (time.Time, time.Time, error) {
	if start == nil || end == nil {
		return time.Time{}, time.Time{}, errors.New("one_time period requires starts_at and ends_at")
	}
	s, err := time.Parse(time.RFC3339, *start)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("starts_at must be RFC3339")
	}
	e, err := time.Parse(time.RFC3339, *end)
	if err != nil || !s.Before(e) {
		return time.Time{}, time.Time{}, errors.New("ends_at must be after starts_at")
	}
	return s, e, nil
}

func validateIDs(ids []string) error {
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			return errors.New("duplicate id")
		}
		seen[id] = struct{}{}
	}
	return nil
}

func decodeBudgetRequest(r *http.Request, target any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var extra any
	if decoder.Decode(&extra) == nil {
		return errors.New("multiple JSON values")
	}
	return nil
}

func writeBudgetJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
