package receipts

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/config"
	"kiroku/backend/internal/db"
	"kiroku/backend/internal/models"
)

type Handler struct {
	db  *db.DB
	cfg *config.Config
}

func NewHandler(database *db.DB, cfg *config.Config) *Handler {
	return &Handler{db: database, cfg: cfg}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/scan", h.ScanReceipt)
	r.Get("/jobs/{id}", h.GetJob)
	r.Post("/jobs/{id}/confirm", h.ConfirmTransaction)
	return r
}

func (h *Handler) logEvent(jobID, eventType, stage string, attempt int, status, message string) {
	id := uuid.New().String()
	_, err := h.db.Exec(
		`INSERT INTO receipt_job_events (id, job_id, event_type, stage, attempt, status, message) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, jobID, eventType, stage, attempt, status, message,
	)
	if err != nil {
		log.Error().Err(err).Msg("receipts: log event error")
	}
}

func (h *Handler) updateJobStatus(jobID, status, stage, errorMessage string) {
	_, err := h.db.Exec(
		`UPDATE receipt_jobs SET status = ?, stage = ?, error_message = ? WHERE id = ?`,
		status, stage, errorMessage, jobID,
	)
	if err != nil {
		log.Error().Err(err).Msg("receipts: update job status error")
	}
}

func (h *Handler) ScanReceipt(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateReceiptJobReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.ImageBase64 == "" {
		http.Error(w, "image_base64 required", http.StatusBadRequest)
		return
	}

	jobID := uuid.New().String()
	_, err := h.db.Exec(
		`INSERT INTO receipt_jobs (id, user_id, status, stage, original_size, compressed_size) VALUES (?, ?, 'processing', 'ocr_pending', ?, ?)`,
		jobID, userID, req.OriginalSize, req.CompressedSize,
	)
	if err != nil {
		log.Error().Err(err).Msg("receipts: create job")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	h.logEvent(jobID, "job_created", "init", 1, "success", "Receipt processing job created")

	// Run processing asynchronously
	go h.processReceiptAsync(jobID, req.ImageBase64)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"job_id": jobID})
}

func (h *Handler) processReceiptAsync(jobID string, imageBase64 string) {
	ctx := context.Background()

	// 1. OCR Stage with retry
	h.logEvent(jobID, "ocr_started", "ocr", 1, "processing", "Reading receipt with OCR...")
	var ocrText string
	var ocrErr error
	maxRetries := 3

	for attempt := 1; attempt <= maxRetries; attempt++ {
		ocrText, ocrErr = CallOCRSpace(ctx, h.cfg.OCRSpaceAPIKey, imageBase64, h.cfg.MockOCREnabled)
		if ocrErr == nil {
			break
		}
		log.Warn().Err(ocrErr).Int("attempt", attempt).Msg("receipts: OCR attempt failed")
		h.logEvent(jobID, "ocr_retry", "ocr", attempt, "warning", "OCR request failed; retrying...")
		time.Sleep(time.Duration(attempt) * 2 * time.Second)
	}

	if ocrErr != nil {
		h.updateJobStatus(jobID, "failed", "ocr_failed", ocrErr.Error())
		h.logEvent(jobID, "ocr_failed", "ocr", maxRetries, "error", "OCR failed: "+ocrErr.Error())
		return
	}

	h.db.Exec(`UPDATE receipt_jobs SET ocr_text = ? WHERE id = ?`, ocrText, jobID)
	h.logEvent(jobID, "ocr_succeeded", "ocr", 1, "success", "Receipt text extracted successfully")

	// Fetch categories and accounts for AI prompt
	categories := h.fetchCategories(ctx)
	accounts := h.fetchAccounts(ctx)

	// 2. AI Normalization Stage with schema validation loop
	h.logEvent(jobID, "ai_started", "normalization", 1, "processing", "Converting receipt text to transaction draft...")
	aiClient := NewNineRouterClient(h.cfg.NineRouterBaseURL, h.cfg.NineRouterAPIKey, h.cfg.NineRouterModelComboName)

	var draft *models.AINormalizationResult
	var aiErr error
	var rawOutput string

	draft, aiErr = aiClient.NormalizeReceipt(ctx, ocrText, categories, accounts, h.cfg.FXBaseCurrency)
	if aiErr == nil {
		draftBytes, _ := json.Marshal(draft)
		rawOutput = string(draftBytes)
	}

	// Schema validation check & retry if invalid
	if aiErr != nil || draft.Amount <= 0 || draft.DateTime == "" {
		validationErr := "Missing or invalid amount/date"
		if aiErr != nil {
			validationErr = aiErr.Error()
		}
		h.logEvent(jobID, "ai_validation_failed", "normalization", 1, "warning", "Validation failed: "+validationErr+". Requesting correction...")

		draft, aiErr = aiClient.RetryWithCorrection(ctx, ocrText, categories, accounts, h.cfg.FXBaseCurrency, rawOutput, validationErr)
		if aiErr == nil {
			draftBytes, _ := json.Marshal(draft)
			rawOutput = string(draftBytes)
		}
	}

	if aiErr != nil || draft == nil || draft.Amount <= 0 {
		errMsg := "AI normalization failed"
		if aiErr != nil {
			errMsg = aiErr.Error()
		}
		h.updateJobStatus(jobID, "failed", "ai_failed", errMsg)
		h.logEvent(jobID, "ai_failed", "normalization", 2, "error", errMsg)
		return
	}

	h.db.Exec(`UPDATE receipt_jobs SET status = 'draft_ready', stage = 'review', ai_raw_output = ?, transaction_draft = ? WHERE id = ?`, rawOutput, rawOutput, jobID)
	h.logEvent(jobID, "draft_ready", "review", 1, "success", "Transaction draft is ready for review")
}

func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	jobID := chi.URLParam(r, "id")
	var job models.ReceiptJob
	var origSize, compSize sql.NullInt64
	var ocrText, aiRaw, draft, errMsg, txID sql.NullString

	err := h.db.QueryRow(
		`SELECT id, user_id, status, stage, original_size, compressed_size, ocr_text, ai_raw_output, transaction_draft, transaction_id, error_message, created_at, updated_at FROM receipt_jobs WHERE id = ? AND user_id = ?`,
		jobID, userID,
	).Scan(&job.ID, &job.UserID, &job.Status, &job.Stage, &origSize, &compSize, &ocrText, &aiRaw, &draft, &txID, &errMsg, &job.CreatedAt, &job.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("receipts: get job")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if origSize.Valid {
		job.OriginalSize = origSize.Int64
	}
	if compSize.Valid {
		job.CompressedSize = compSize.Int64
	}
	if ocrText.Valid {
		job.OCRText = ocrText.String
	}
	if aiRaw.Valid {
		job.AIRawOutput = aiRaw.String
	}
	if draft.Valid {
		job.TransactionDraft = draft.String
	}
	if txID.Valid {
		job.TransactionID = &txID.String
	}
	if errMsg.Valid {
		job.ErrorMessage = errMsg.String
	}

	// Fetch events
	rows, err := h.db.Query(`SELECT id, job_id, event_type, stage, attempt, status, message, created_at FROM receipt_job_events WHERE job_id = ? ORDER BY created_at ASC`, jobID)
	var events []models.ReceiptJobEvent
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ev models.ReceiptJobEvent
			var msg sql.NullString
			rows.Scan(&ev.ID, &ev.JobID, &ev.EventType, &ev.Stage, &ev.Attempt, &ev.Status, &msg, &ev.CreatedAt)
			if msg.Valid {
				ev.Message = msg.String
			}
			events = append(events, ev)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ReceiptJobResponse{
		JobID:  jobID,
		Job:    job,
		Events: events,
	})
}

func (h *Handler) ConfirmTransaction(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	jobID := chi.URLParam(r, "id")
	var status, draftStr string
	err := h.db.QueryRow(`SELECT status, transaction_draft FROM receipt_jobs WHERE id = ? AND user_id = ?`, jobID, userID).Scan(&status, &draftStr)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil || status != "draft_ready" {
		http.Error(w, "job not ready or already confirmed", http.StatusBadRequest)
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

	txID := uuid.New().String()
	_, err = h.db.Exec(`
		INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, place_lat, place_lng)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, txID, req.Type, req.DateTime, req.Amount, req.Currency, req.FromAccountID, req.ToAccountID, req.CategoryID, req.Labels, req.Note, req.PayeePayer, req.PaymentType, req.WarrantyUntil, req.Status, req.PlaceName, req.PlaceLat, req.PlaceLng)

	if err != nil {
		log.Error().Err(err).Msg("receipts: confirm transaction create")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	h.db.Exec(`UPDATE receipt_jobs SET status = 'confirmed', stage = 'completed', transaction_id = ? WHERE id = ?`, txID, jobID)
	h.logEvent(jobID, "transaction_saved", "completed", 1, "success", "Transaction saved successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": txID})
}

func (h *Handler) fetchCategories(ctx context.Context) []models.Category {
	rows, err := h.db.QueryContext(ctx, `SELECT id, parent_id, name, icon, color, kind, sort_order FROM categories WHERE archived = 0 ORDER BY sort_order, name`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var cats []models.Category
	for rows.Next() {
		var c models.Category
		var parentID, kind sql.NullString
		if err := rows.Scan(&c.ID, &parentID, &c.Name, &c.Icon, &c.Color, &kind, &c.SortOrder); err != nil {
			continue
		}
		if parentID.Valid {
			c.ParentID = &parentID.String
		}
		if kind.Valid {
			c.Kind = &kind.String
		}
		cats = append(cats, c)
	}
	return cats
}

func (h *Handler) fetchAccounts(ctx context.Context) []models.Account {
	rows, err := h.db.QueryContext(ctx, `SELECT id, name, currency, type, icon, color FROM accounts WHERE archived = 0`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var accs []models.Account
	for rows.Next() {
		var a models.Account
		rows.Scan(&a.ID, &a.Name, &a.Currency, &a.Type, &a.Icon, &a.Color)
		accs = append(accs, a)
	}
	return accs
}
