package impexp

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	r.Get("/export/csv", h.ExportCSV)
	r.Post("/import/csv", h.ImportCSVPreview)
	r.Post("/import/csv/confirm", h.ImportCSVConfirm)
	return r
}

func (h *Handler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	query := `
		SELECT t.date_time, t.type, t.amount, t.currency,
		       fa.name as from_account, ta.name as to_account,
		       c.name as category, t.labels, t.note, t.payee_payer,
		       t.payment_type, t.warranty_until, t.status,
		       t.place_name, t.attachments
		FROM transactions t
		LEFT JOIN accounts fa ON t.from_account_id = fa.id
		LEFT JOIN accounts ta ON t.to_account_id = ta.id
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.deleted_at IS NULL
	`
	args := []interface{}{}

	conditions := []string{}
	if df := r.URL.Query().Get("date_from"); df != "" {
		conditions = append(conditions, "t.date_time >= ?")
		args = append(args, df)
	}
	if dt := r.URL.Query().Get("date_to"); dt != "" {
		conditions = append(conditions, "t.date_time <= ?")
		args = append(args, dt)
	}
	if typ := r.URL.Query().Get("type"); typ != "" {
		conditions = append(conditions, "t.type = ?")
		args = append(args, typ)
	}
	if cid := r.URL.Query().Get("category_id"); cid != "" {
		conditions = append(conditions, "t.category_id = ?")
		args = append(args, cid)
	}
	if aid := r.URL.Query().Get("account_id"); aid != "" {
		conditions = append(conditions, "(t.from_account_id = ? OR t.to_account_id = ?)")
		args = append(args, aid, aid)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY t.date_time DESC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		log.Error().Err(err).Msg("export: query")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="transactions_`+time.Now().Format("20060102")+`.csv"`)

	writer := csv.NewWriter(w)
	defer writer.Flush()

	headers := []string{"date_time", "type", "amount", "currency", "from_account", "to_account", "category", "labels", "note", "payee_payer", "payment_type", "warranty_until", "status", "place", "attachments"}
	writer.Write(headers)

	for rows.Next() {
		var dateTime time.Time
		var typ, currency, fromAcc, toAcc, category, labels, note, payeePayer, paymentType, warrantyUntil, status, placeName, attachments sql.NullString
		var amount int64

		if err := rows.Scan(&dateTime, &typ, &amount, &currency, &fromAcc, &toAcc, &category, &labels, &note, &payeePayer, &paymentType, &warrantyUntil, &status, &placeName, &attachments); err != nil {
			log.Error().Err(err).Msg("export: scan")
			continue
		}

		record := []string{
			dateTime.Format(time.RFC3339),
			typ.String,
			strconv.FormatInt(amount, 10),
			currency.String,
			fromAcc.String,
			toAcc.String,
			category.String,
			labels.String,
			note.String,
			payeePayer.String,
			paymentType.String,
			warrantyUntil.String,
			status.String,
			placeName.String,
			attachments.String,
		}
		writer.Write(record)
	}
}

func (h *Handler) ImportCSVPreview(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := r.FormValue("content_type")
	if contentType == "" {
		contentType = r.Header.Get("Content-Type")
	}

	createMissing := r.FormValue("create_missing") == "true"

	var records [][]string
	if strings.Contains(contentType, "spreadsheetml") || strings.Contains(contentType, "excel") || strings.HasSuffix(r.FormValue("filename"), ".xlsx") {
		records, err = h.parseExcel(file)
	} else {
		records, err = h.parseCSV(file)
	}
	if err != nil {
		http.Error(w, "parse error: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(records) < 2 {
		http.Error(w, "CSV must have header and at least one row", http.StatusBadRequest)
		return
	}

	headers := records[0]
	dataRows := records[1:]
	if len(dataRows) > 10 {
		dataRows = dataRows[:10]
	}

	var preview models.ImportPreview
	preview.PreviewID = uuid.New().String()

	catCache := make(map[string]string)
	accCache := make(map[string]string)

	for _, row := range dataRows {
		if len(row) < len(headers) {
			continue
		}
		rec := h.rowToImport(headers, row)
		preview.Rows = append(preview.Rows, rec)

		if rec.Category != "" {
			if _, ok := catCache[rec.Category]; !ok {
				match := h.fuzzyMatchCategory(rec.Category, createMissing)
				preview.Matches.Categories = append(preview.Matches.Categories, match)
				catCache[rec.Category] = match.MatchID
			}
		}
		if rec.FromAccount != "" {
			if _, ok := accCache[rec.FromAccount]; !ok {
				match := h.fuzzyMatchAccount(rec.FromAccount, createMissing)
				preview.Matches.Accounts = append(preview.Matches.Accounts, match)
				accCache[rec.FromAccount] = match.MatchID
			}
		}
		if rec.ToAccount != "" {
			if _, ok := accCache[rec.ToAccount]; !ok {
				match := h.fuzzyMatchAccount(rec.ToAccount, createMissing)
				preview.Matches.Accounts = append(preview.Matches.Accounts, match)
				accCache[rec.ToAccount] = match.MatchID
			}
		}
	}

	h.storePreview(preview.PreviewID, preview.Rows, preview.Matches)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(preview)
}

func (h *Handler) ImportCSVConfirm(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.ImportConfirmReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if !req.Confirm {
		http.Error(w, "confirm must be true", http.StatusBadRequest)
		return
	}

	var data string
	err := h.db.QueryRow("SELECT data FROM import_log WHERE id = ?", req.PreviewID).Scan(&data)
	if err == sql.ErrNoRows {
		http.Error(w, "preview not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var preview models.ImportPreview
	json.Unmarshal([]byte(data), &preview)

	created := 0
	skipped := 0
	var errors []string

	for _, row := range preview.Rows {
		hash := rowHash(row)
		var exists bool
		h.db.QueryRow("SELECT 1 FROM import_log WHERE id = ?", hash).Scan(&exists)
		if exists {
			skipped++
			continue
		}

		dateTime, err := time.Parse(time.RFC3339, row.DateTime)
		if err != nil {
			dateTime = time.Now()
		}

		amount, _ := strconv.ParseInt(row.Amount, 10, 64)
		if amount == 0 {
			errors = append(errors, row.DateTime+": amount required")
			continue
		}

		fromAccID := h.resolveAccount(row.FromAccount)
		toAccID := h.resolveAccount(row.ToAccount)
		catID := h.resolveCategory(row.Category)

		if row.Type == "expense" && fromAccID == nil {
			errors = append(errors, row.DateTime+": expense requires from_account")
			continue
		}
		if row.Type == "income" && toAccID == nil {
			errors = append(errors, row.DateTime+": income requires to_account")
			continue
		}

		id := uuid.New().String()
		_, err = h.db.Exec(`
			INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, attachments)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, id, row.Type, dateTime, amount, row.Currency, fromAccID, toAccID, catID, row.Labels, row.Note, row.PayeePayer, row.PaymentType, parseTime(row.WarrantyUntil), defaultStatus(row.Status), row.Place, row.Attachments)
		if err != nil {
			errors = append(errors, row.DateTime+": "+err.Error())
			continue
		}

		h.db.Exec("INSERT INTO import_log (id, data, created_at) VALUES (?, ?, ?)", hash, string(data), time.Now())
		created++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"created": created,
		"skipped": skipped,
		"errors":  errors,
	})
}
