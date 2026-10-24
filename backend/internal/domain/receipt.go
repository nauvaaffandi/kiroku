package domain

import (
	"encoding/json"
	"time"
)

type ReceiptJobStatus string

const (
	ReceiptStatusPending    ReceiptJobStatus = "pending"
	ReceiptStatusProcessing ReceiptJobStatus = "processing"
	ReceiptStatusOCRDone    ReceiptJobStatus = "ocr_done"
	ReceiptStatusAIDone     ReceiptJobStatus = "ai_done"
	ReceiptStatusConfirmed  ReceiptJobStatus = "confirmed"
	ReceiptStatusFailed     ReceiptJobStatus = "failed"
)

type ReceiptJob struct {
	ID                string          `json:"id"`
	UserID            string          `json:"user_id"`
	Status            string          `json:"status"`
	Stage             string          `json:"stage"`
	OriginalSize      int64           `json:"original_size"`
	CompressedSize    int64           `json:"compressed_size"`
	OCRText           string          `json:"ocr_text"`
	AIRawOutput       string          `json:"ai_raw_output"`
	TransactionDraft  json.RawMessage `json:"transaction_draft"`
	TransactionID     *string         `json:"transaction_id,omitempty"`
	ErrorMessage      string          `json:"error_message,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

type ReceiptJobEvent struct {
	ID        string    `json:"id"`
	JobID     string    `json:"job_id"`
	EventType string    `json:"event_type"`
	Stage     string    `json:"stage"`
	Attempt   int       `json:"attempt"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Metadata  string    `json:"metadata,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type OCRResponse struct {
	Text  string `json:"text"`
	Parsed *struct {
		Merchant string  `json:"merchant"`
		Date     string  `json:"date"`
		Total    float64 `json:"total"`
		Currency string  `json:"currency"`
		Items    []struct {
			Name   string  `json:"name"`
			Amount float64 `json:"amount"`
		} `json:"items"`
	} `json:"parsed"`
}
