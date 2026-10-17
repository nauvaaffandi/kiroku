package templates

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

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

	var req models.CreateTemplateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Type == "" {
		http.Error(w, "name and type required", http.StatusBadRequest)
		return
	}

	if req.Status == "" {
		req.Status = "cleared"
	}

	id := uuid.New().String()
	_, err := h.db.Exec(`
		INSERT INTO templates (id, name, type, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, place_lat, place_lng, attachments)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, id, req.Name, req.Type, req.Amount, req.Currency, req.FromAccountID, req.ToAccountID, req.CategoryID, req.Labels, req.Note, req.PayeePayer, req.PaymentType, req.WarrantyUntil, req.Status, req.PlaceName, req.PlaceLat, req.PlaceLng, req.Attachments)
	if err != nil {
		log.Error().Err(err).Msg("templates: create")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	var req models.UpdateTemplateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	setParts := []string{}
	args := []interface{}{}
	if req.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *req.Name)
	}
	if req.Type != nil {
		setParts = append(setParts, "type = ?")
		args = append(args, *req.Type)
	}
	if req.Amount != nil {
		setParts = append(setParts, "amount = ?")
		args = append(args, *req.Amount)
	}
	if req.Currency != nil {
		setParts = append(setParts, "currency = ?")
		args = append(args, *req.Currency)
	}
	if req.FromAccountID != nil {
		setParts = append(setParts, "from_account_id = ?")
		args = append(args, *req.FromAccountID)
	}
	if req.ToAccountID != nil {
		setParts = append(setParts, "to_account_id = ?")
		args = append(args, *req.ToAccountID)
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
	_, err := h.db.Exec("UPDATE templates SET "+join(setParts, ", ")+" WHERE id = ?", args...)
	if err != nil {
		log.Error().Err(err).Msg("templates: update")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) Apply(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	query := `SELECT id, name, type, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, place_lat, place_lng, attachments FROM templates WHERE id = ?`
	var t models.Template
	var amount sql.NullInt64
	var currency, fromAcc, toAcc, catID, paymentType sql.NullString
	var warranty sql.NullTime
	var placeLat, placeLng sql.NullFloat64
	err := h.db.QueryRow(query, id).Scan(&t.ID, &t.Name, &t.Type, &amount, &currency, &fromAcc, &toAcc, &catID, &t.Labels, &t.Note, &t.PayeePayer, &paymentType, &warranty, &t.Status, &t.PlaceName, &placeLat, &placeLng, &t.Attachments)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("templates: apply")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if amount.Valid {
		t.Amount = &amount.Int64
	}
	if currency.Valid {
		t.Currency = currency.String
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
	if paymentType.Valid {
		t.PaymentType = paymentType.String
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

	// Return as transaction-like object with date_time = now, amount from template or 0
	response := map[string]interface{}{
		"type":            t.Type,
		"date_time":       time.Now().Format(time.RFC3339),
		"amount":          0,
		"currency":        t.Currency,
		"from_account_id": t.FromAccountID,
		"to_account_id":   t.ToAccountID,
		"category_id":     t.CategoryID,
		"labels":          t.Labels,
		"note":            t.Note,
		"payee_payer":     t.PayeePayer,
		"payment_type":    t.PaymentType,
		"warranty_until":  t.WarrantyUntil,
		"status":          t.Status,
		"place_name":      t.PlaceName,
		"place_lat":       t.PlaceLat,
		"place_lng":       t.PlaceLng,
		"attachments":     t.Attachments,
		"template_id":     t.ID,
	}
	if t.Amount != nil {
		response["amount"] = *t.Amount
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
