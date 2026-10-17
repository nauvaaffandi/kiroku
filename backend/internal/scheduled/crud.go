package scheduled

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/teambition/rrule-go"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/models"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateScheduledReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.RRule == "" {
		http.Error(w, "rrule required", http.StatusBadRequest)
		return
	}

	// Validate RRULE
	_, err := rrule.StrToRRule(req.RRule)
	if err != nil {
		http.Error(w, "invalid rrule: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate template JSON has required fields
	var tmpl map[string]interface{}
	if err := json.Unmarshal(req.TransactionTemplateJSON, &tmpl); err != nil {
		http.Error(w, "invalid transaction_template_json", http.StatusBadRequest)
		return
	}
	if tmpl["type"] == nil {
		http.Error(w, "transaction_template_json must include type", http.StatusBadRequest)
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	id := uuid.New().String()

	// Compute first next_run from rrule
	nextRun := computeNextRun(req.RRule, time.Now())

	_, err = h.db.Exec(`
		INSERT INTO scheduled (id, rrule, max_occurrences, transaction_template_json, enabled, next_run, occurrences_done)
		VALUES (?, ?, ?, ?, ?, ?, 0)
	`, id, req.RRule, req.MaxOccurrences, string(req.TransactionTemplateJSON), enabled, nextRun)
	if err != nil {
		log.Error().Err(err).Msg("scheduled: create")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *Handler) CreatePreset(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.PresetReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var rrule string
	switch req.Preset {
	case "monthly":
		day := 1
		if d, ok := req.Params["day"].(float64); ok {
			day = int(d)
		}
		rrule = "FREQ=MONTHLY;BYMONTHDAY=" + string(rune(day+'0'))
	case "weekly":
		weekday := "MO"
		if w, ok := req.Params["weekday"].(string); ok {
			weekday = w
		}
		rrule = "FREQ=WEEKLY;BYDAY=" + weekday
	case "yearly":
		month := 1
		day := 1
		if m, ok := req.Params["month"].(float64); ok {
			month = int(m)
		}
		if d, ok := req.Params["day"].(float64); ok {
			day = int(d)
		}
		rrule = "FREQ=YEARLY;BYMONTH=" + string(rune(month+'0')) + ";BYMONTHDAY=" + string(rune(day+'0'))
	default:
		http.Error(w, "unknown preset", http.StatusBadRequest)
		return
	}

	tmplJSON, _ := json.Marshal(req.Template)

	// Reuse Create logic
	id := uuid.New().String()
	enabled := true
	nextRun := computeNextRun(rrule, time.Now())
	_, err := h.db.Exec(`
		INSERT INTO scheduled (id, rrule, max_occurrences, transaction_template_json, enabled, next_run, occurrences_done)
		VALUES (?, ?, ?, ?, ?, ?, 0)
	`, id, rrule, nil, string(tmplJSON), enabled, nextRun)
	if err != nil {
		log.Error().Err(err).Msg("scheduled: preset create")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id, "rrule": rrule})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	var req models.UpdateScheduledReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.RRule != nil {
		_, err := rrule.StrToRRule(*req.RRule)
		if err != nil {
			http.Error(w, "invalid rrule: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	setParts := []string{}
	args := []interface{}{}
	if req.RRule != nil {
		setParts = append(setParts, "rrule = ?")
		args = append(args, *req.RRule)
		// Recompute next_run
		nextRun := computeNextRun(*req.RRule, time.Now())
		setParts = append(setParts, "next_run = ?")
		args = append(args, nextRun)
	}
	if req.MaxOccurrences != nil {
		setParts = append(setParts, "max_occurrences = ?")
		args = append(args, *req.MaxOccurrences)
	}
	if req.TransactionTemplateJSON != nil {
		setParts = append(setParts, "transaction_template_json = ?")
		args = append(args, string(req.TransactionTemplateJSON))
	}
	if req.Enabled != nil {
		setParts = append(setParts, "enabled = ?")
		if *req.Enabled {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
	}

	if len(setParts) == 0 {
		http.Error(w, "no fields to update", http.StatusBadRequest)
		return
	}

	args = append(args, id)
	_, err := h.db.Exec("UPDATE scheduled SET "+join(setParts, ", ")+" WHERE id = ?", args...)
	if err != nil {
		log.Error().Err(err).Msg("scheduled: update")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
