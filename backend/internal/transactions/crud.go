package transactions

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/models"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateTransactionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Type == "" || req.Amount <= 0 || req.Currency == "" || req.DateTime.IsZero() {
		http.Error(w, "type, amount, currency, date_time required", http.StatusBadRequest)
		return
	}

	if req.Type == "transfer" {
		if req.FromAccountID == nil || req.ToAccountID == nil {
			http.Error(w, "transfer requires from_account_id and to_account_id", http.StatusBadRequest)
			return
		}
		if *req.FromAccountID == *req.ToAccountID {
			http.Error(w, "from_account_id and to_account_id must differ", http.StatusBadRequest)
			return
		}
	} else if req.Type == "income" {
		if req.ToAccountID == nil {
			http.Error(w, "income requires to_account_id", http.StatusBadRequest)
			return
		}
	} else if req.Type == "expense" {
		if req.FromAccountID == nil {
			http.Error(w, "expense requires from_account_id", http.StatusBadRequest)
			return
		}
	}

	if req.Status == "" {
		req.Status = "cleared"
	}

	id := uuid.New().String()

	_, err := h.db.Exec(`
		INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, place_lat, place_lng, attachments, template_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, id, req.Type, req.DateTime, req.Amount, req.Currency, req.FromAccountID, req.ToAccountID, req.CategoryID, req.Labels, req.Note, req.PayeePayer, req.PaymentType, req.WarrantyUntil, req.Status, req.PlaceName, req.PlaceLat, req.PlaceLng, req.Attachments, req.TemplateID)
	if err != nil {
		log.Error().Err(err).Msg("transactions: create")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	query := `SELECT id, type, date_time, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, place_lat, place_lng, attachments, template_id, created_at, updated_at FROM transactions WHERE id = ? AND deleted_at IS NULL`
	var t models.Transaction
	var fromAcc, toAcc, catID, templateID sql.NullString
	var warranty sql.NullTime
	var placeLat, placeLng sql.NullFloat64
	err := h.db.QueryRow(query, id).Scan(&t.ID, &t.Type, &t.DateTime, &t.Amount, &t.Currency, &fromAcc, &toAcc, &catID, &t.Labels, &t.Note, &t.PayeePayer, &t.PaymentType, &warranty, &t.Status, &t.PlaceName, &placeLat, &placeLng, &t.Attachments, &templateID, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("transactions: get")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if fromAcc.Valid {
		t.FromAccountID = &fromAcc.String
	}
	if toAcc.Valid {
		t.ToAccountID = &toAcc.String
	}
	if catID.Valid {
		t.CategoryID = &catID.String
	}
	if templateID.Valid {
		t.TemplateID = &templateID.String
	}
	if warranty.Valid {
		t.WarrantyUntil = &warranty.Time
	}
	if placeLat.Valid {
		t.PlaceLat = &placeLat.Float64
	}
	if placeLng.Valid {
		t.PlaceLng = &placeLng.Float64
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	var req models.UpdateTransactionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var txType string
	err := h.db.QueryRow("SELECT type FROM transactions WHERE id = ?", id).Scan(&txType)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("transactions: update check")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	setParts := []string{}
	args := []interface{}{}
	if req.DateTime != nil {
		setParts = append(setParts, "date_time = ?")
		args = append(args, *req.DateTime)
	}
	if req.Amount != nil {
		setParts = append(setParts, "amount = ?")
		args = append(args, *req.Amount)
	}
	if req.CategoryID != nil {
		setParts = append(setParts, "category_id = ?")
		args = append(args, *req.CategoryID)
	}
	if req.Labels != nil {
		setParts = append(setParts, "labels = ?")
		args = append(args, *req.Labels)
	}
	if req.Note != nil {
		setParts = append(setParts, "note = ?")
		args = append(args, *req.Note)
	}
	if req.PayeePayer != nil {
		setParts = append(setParts, "payee_payer = ?")
		args = append(args, *req.PayeePayer)
	}
	if req.PaymentType != nil {
		setParts = append(setParts, "payment_type = ?")
		args = append(args, *req.PaymentType)
	}
	if req.WarrantyUntil != nil {
		setParts = append(setParts, "warranty_until = ?")
		args = append(args, *req.WarrantyUntil)
	}
	if req.Status != nil {
		setParts = append(setParts, "status = ?")
		args = append(args, *req.Status)
	}
	if req.PlaceName != nil {
		setParts = append(setParts, "place_name = ?")
		args = append(args, *req.PlaceName)
	}
	if req.PlaceLat != nil {
		setParts = append(setParts, "place_lat = ?")
		args = append(args, *req.PlaceLat)
	}
	if req.PlaceLng != nil {
		setParts = append(setParts, "place_lng = ?")
		args = append(args, *req.PlaceLng)
	}
	if req.Attachments != nil {
		setParts = append(setParts, "attachments = ?")
		args = append(args, *req.Attachments)
	}

	if len(setParts) == 0 {
		http.Error(w, "no fields to update", http.StatusBadRequest)
		return
	}

	args = append(args, id)
	_, err = h.db.Exec("UPDATE transactions SET "+join(setParts, ", ")+" WHERE id = ?", args...)
	if err != nil {
		log.Error().Err(err).Msg("transactions: update")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	_, err := h.db.Exec("UPDATE transactions SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?", id)
	if err != nil {
		log.Error().Err(err).Msg("transactions: delete")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
