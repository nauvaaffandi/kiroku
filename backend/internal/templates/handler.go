package templates

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/db"
	"kiroku/backend/internal/models"
)

type Handler struct {
	db *db.DB
}

func NewHandler(database *db.DB) *Handler {
	return &Handler{db: database}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	r.Post("/{id}/apply", h.Apply)
	return r
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	query := `SELECT id, name, type, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, place_lat, place_lng, attachments, created_at, updated_at FROM templates ORDER BY created_at DESC`
	rows, err := h.db.Query(query)
	if err != nil {
		log.Error().Err(err).Msg("templates: list")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var templates []models.Template
	for rows.Next() {
		var t models.Template
		var amount sql.NullInt64
		var currency, fromAcc, toAcc, catID, paymentType sql.NullString
		var warranty sql.NullTime
		var placeLat, placeLng sql.NullFloat64
		if err := rows.Scan(&t.ID, &t.Name, &t.Type, &amount, &currency, &fromAcc, &toAcc, &catID, &t.Labels, &t.Note, &t.PayeePayer, &paymentType, &warranty, &t.Status, &t.PlaceName, &placeLat, &placeLng, &t.Attachments, &t.CreatedAt, &t.UpdatedAt); err != nil {
			log.Error().Err(err).Msg("templates: scan")
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
		templates = append(templates, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	query := `SELECT id, name, type, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, place_lat, place_lng, attachments, created_at, updated_at FROM templates WHERE id = ?`
	var t models.Template
	var amount sql.NullInt64
	var currency, fromAcc, toAcc, catID, paymentType sql.NullString
	var warranty sql.NullTime
	var placeLat, placeLng sql.NullFloat64
	err := h.db.QueryRow(query, id).Scan(&t.ID, &t.Name, &t.Type, &amount, &currency, &fromAcc, &toAcc, &catID, &t.Labels, &t.Note, &t.PayeePayer, &paymentType, &warranty, &t.Status, &t.PlaceName, &placeLat, &placeLng, &t.Attachments, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("templates: get")
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	_, err := h.db.Exec("DELETE FROM templates WHERE id = ?", id)
	if err != nil {
		log.Error().Err(err).Msg("templates: delete")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
