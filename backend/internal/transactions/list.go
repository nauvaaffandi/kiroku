package transactions

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/models"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	query := `SELECT id, type, date_time, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, place_lat, place_lng, attachments, template_id, created_at, updated_at FROM transactions WHERE deleted_at IS NULL`
	args := []interface{}{}

	conditions := []string{}
	if aid := r.URL.Query().Get("account_id"); aid != "" {
		conditions = append(conditions, "(from_account_id = ? OR to_account_id = ?)")
		args = append(args, aid, aid)
	}
	if cid := r.URL.Query().Get("category_id"); cid != "" {
		conditions = append(conditions, "category_id = ?")
		args = append(args, cid)
	}
	if typ := r.URL.Query().Get("type"); typ != "" {
		conditions = append(conditions, "type = ?")
		args = append(args, typ)
	}
	if status := r.URL.Query().Get("status"); status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}
	if df := r.URL.Query().Get("date_from"); df != "" {
		conditions = append(conditions, "date_time >= ?")
		args = append(args, df)
	}
	if dt := r.URL.Query().Get("date_to"); dt != "" {
		conditions = append(conditions, "date_time <= ?")
		args = append(args, dt)
	}
	if amin := r.URL.Query().Get("amount_min"); amin != "" {
		conditions = append(conditions, "amount >= ?")
		args = append(args, amin)
	}
	if amax := r.URL.Query().Get("amount_max"); amax != "" {
		conditions = append(conditions, "amount <= ?")
		args = append(args, amax)
	}
	if q := r.URL.Query().Get("q"); q != "" {
		conditions = append(conditions, "id IN (SELECT id FROM transactions_fts WHERE transactions_fts MATCH ?)")
		args = append(args, q)
	}
	if label := r.URL.Query().Get("label"); label != "" {
		conditions = append(conditions, "labels LIKE ?")
		args = append(args, "%"+label+"%")
	}
	if payee := r.URL.Query().Get("payee"); payee != "" {
		conditions = append(conditions, "payee_payer LIKE ?")
		args = append(args, "%"+payee+"%")
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		sortBy = "date_time"
	}
	order := r.URL.Query().Get("order")
	if order == "" {
		order = "DESC"
	}
	query += " ORDER BY " + sortBy + " " + order

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		log.Error().Err(err).Msg("transactions: list")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		var fromAcc, toAcc, catID, templateID sql.NullString
		var warranty sql.NullTime
		var placeLat, placeLng sql.NullFloat64
		if err := rows.Scan(&t.ID, &t.Type, &t.DateTime, &t.Amount, &t.Currency, &fromAcc, &toAcc, &catID, &t.Labels, &t.Note, &t.PayeePayer, &t.PaymentType, &warranty, &t.Status, &t.PlaceName, &placeLat, &placeLng, &t.Attachments, &templateID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			log.Error().Err(err).Msg("transactions: scan")
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
		transactions = append(transactions, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "q parameter required", http.StatusBadRequest)
		return
	}

	query := `SELECT id, type, date_time, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, place_lat, place_lng, attachments, template_id, created_at, updated_at FROM transactions WHERE id IN (SELECT id FROM transactions_fts WHERE transactions_fts MATCH ?) AND deleted_at IS NULL ORDER BY date_time DESC LIMIT 50`
	rows, err := h.db.Query(query, q)
	if err != nil {
		log.Error().Err(err).Msg("transactions: search")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		var fromAcc, toAcc, catID, templateID sql.NullString
		var warranty sql.NullTime
		var placeLat, placeLng sql.NullFloat64
		if err := rows.Scan(&t.ID, &t.Type, &t.DateTime, &t.Amount, &t.Currency, &fromAcc, &toAcc, &catID, &t.Labels, &t.Note, &t.PayeePayer, &t.PaymentType, &warranty, &t.Status, &t.PlaceName, &placeLat, &placeLng, &t.Attachments, &templateID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			log.Error().Err(err).Msg("transactions: search scan")
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
		transactions = append(transactions, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (h *Handler) Autocomplete(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "q parameter required", http.StatusBadRequest)
		return
	}

	query := `SELECT DISTINCT note FROM transactions WHERE note LIKE ? AND deleted_at IS NULL ORDER BY date_time DESC LIMIT 10`
	rows, err := h.db.Query(query, q+"%")
	if err != nil {
		log.Error().Err(err).Msg("transactions: autocomplete")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var suggestions []string
	for rows.Next() {
		var note string
		if err := rows.Scan(&note); err != nil {
			continue
		}
		if note != "" {
			suggestions = append(suggestions, note)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	query := `SELECT currency, SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END) as income, SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END) as expense, SUM(CASE WHEN type = 'income' THEN amount ELSE -amount END) as net FROM transactions WHERE deleted_at IS NULL`
	args := []interface{}{}

	conditions := []string{}
	if cid := r.URL.Query().Get("category_id"); cid != "" {
		conditions = append(conditions, "category_id = ?")
		args = append(args, cid)
	}
	if aid := r.URL.Query().Get("account_id"); aid != "" {
		conditions = append(conditions, "(from_account_id = ? OR to_account_id = ?)")
		args = append(args, aid, aid)
	}
	if typ := r.URL.Query().Get("type"); typ != "" {
		conditions = append(conditions, "type = ?")
		args = append(args, typ)
	}
	if df := r.URL.Query().Get("date_from"); df != "" {
		conditions = append(conditions, "date_time >= ?")
		args = append(args, df)
	}
	if dt := r.URL.Query().Get("date_to"); dt != "" {
		conditions = append(conditions, "date_time <= ?")
		args = append(args, dt)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}
	query += " GROUP BY currency"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		log.Error().Err(err).Msg("transactions: summary")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var summary []models.SummaryItem
	for rows.Next() {
		var s models.SummaryItem
		if err := rows.Scan(&s.Currency, &s.Income, &s.Expense, &s.Net); err != nil {
			log.Error().Err(err).Msg("transactions: summary scan")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		summary = append(summary, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
