package httpdelivery

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"kiroku/backend/internal/domain"
)

type GoalsHandler struct {
	repo domain.GoalRepository
}

func NewGoalsHandler(repo domain.GoalRepository) *GoalsHandler {
	return &GoalsHandler{repo: repo}
}

func (h *GoalsHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	r.Post("/{id}/contribute", h.Contribute)
	r.Post("/{id}/reach", h.Reach)
	r.Get("/{id}/contributions", h.Contributions)
	r.Get("/{id}/projection", h.Projection)
	return r
}

func (h *GoalsHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.List(r.Context())
	if err != nil {
		log.Error().Err(err).Str("component", "http").Str("method", "ListGoals").Str("path", r.URL.Path).Msg("goals/list failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if list == nil {
		list = []domain.Goal{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *GoalsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateGoalInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	if req.TargetAmount <= 0 {
		http.Error(w, "target_amount must be greater than 0", http.StatusBadRequest)
		return
	}
	if req.Currency == "" {
		http.Error(w, "currency required", http.StatusBadRequest)
		return
	}
	g, err := h.repo.Create(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Str("component", "http").Str("method", "CreateGoal").Str("path", r.URL.Path).Msg("goals/create failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

func (h *GoalsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	g, err := h.repo.Get(r.Context(), id)
	if err != nil || g == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

func (h *GoalsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req domain.UpdateGoalInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	_, err := h.repo.Update(r.Context(), id, req)
	if err != nil {
		log.Error().Err(err).Str("component", "http").Str("method", "UpdateGoal").Str("path", r.URL.Path).Str("id", id).Msg("goals/update failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *GoalsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.Delete(r.Context(), id); err != nil {
		log.Error().Err(err).Str("component", "http").Str("method", "DeleteGoal").Str("path", r.URL.Path).Str("id", id).Msg("goals/delete failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *GoalsHandler) Contribute(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req domain.ContributeGoalInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.Amount <= 0 {
		http.Error(w, "amount must be greater than 0", http.StatusBadRequest)
		return
	}
	g, err := h.repo.Contribute(r.Context(), id, req)
	if err != nil {
		if g == nil && err == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		log.Error().Err(err).Str("component", "http").Str("method", "ContributeGoal").Str("path", r.URL.Path).Str("id", id).Msg("goals/contribute failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

func (h *GoalsHandler) Reach(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	g, err := h.repo.Reach(r.Context(), id)
	if err != nil {
		if g == nil && err == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		log.Error().Err(err).Str("component", "http").Str("method", "ReachGoal").Str("path", r.URL.Path).Str("id", id).Msg("goals/reach failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

func (h *GoalsHandler) Contributions(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	list, err := h.repo.Contributions(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("component", "http").Str("method", "ListGoalContributions").Str("path", r.URL.Path).Msg("goals/contributions failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if list == nil {
		list = []domain.GoalContribution{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *GoalsHandler) Projection(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	g, err := h.repo.Get(r.Context(), id)
	if err != nil || g == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	remaining := g.TargetAmount - g.SavedAlready
	if remaining <= 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(domain.GoalProjection{
			CurrentSaved:       g.SavedAlready,
			TargetAmount:       g.TargetAmount,
			Remaining:          0,
			WeeklyRate:         0,
			EstimatedWeeks:     0,
			EstimatedReachDate: formatDesiredDate(g.DesiredDate),
			Currency:           g.Currency,
		})
		return
	}

	// Calculate weekly contribution rate from contribution history
	contribs, err := h.repo.Contributions(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("component", "http").Str("method", "GoalProjection").Str("path", r.URL.Path).Msg("goals/projection contributions failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var weeklyRate float64
	if len(contribs) > 0 {
		now := time.Now()
		oldest := contribs[0].AddedAt
		if len(contribs) > 1 {
			oldest = contribs[len(contribs)-1].AddedAt
		}
		days := now.Sub(oldest).Hours() / 24
		if days > 0 {
			var total float64
			for _, c := range contribs {
				total += c.Amount
			}
			weeklyRate = total / days * 7
		}
	}

	estimatedWeeks := 0
	var reachDate string
	if weeklyRate > 0 {
		estimatedWeeks = int(remaining / weeklyRate)
		if estimatedWeeks < 1 {
			estimatedWeeks = 1
		}
		reachDate = time.Now().AddDate(0, 0, estimatedWeeks*7).Format("02-01-2006")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain.GoalProjection{
		CurrentSaved:       g.SavedAlready,
		TargetAmount:       g.TargetAmount,
		Remaining:          remaining,
		WeeklyRate:         weeklyRate,
		EstimatedWeeks:     estimatedWeeks,
		EstimatedReachDate: reachDate,
		Currency:           g.Currency,
	})
}

func formatDesiredDate(s string) string {
	if s == "" {
		return ""
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return s
	}
	return t.Format("02-01-2006")
}
