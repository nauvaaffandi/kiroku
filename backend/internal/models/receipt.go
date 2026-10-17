package models

import "time"

type ReceiptJob struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	Status           string    `json:"status"`
	Stage            string    `json:"stage"`
	OriginalSize     int64     `json:"original_size"`
	CompressedSize   int64     `json:"compressed_size"`
	OCRText          string    `json:"ocr_text,omitempty"`
	AIRawOutput      string    `json:"ai_raw_output,omitempty"`
	TransactionDraft string    `json:"transaction_draft,omitempty"`
	TransactionID    *string   `json:"transaction_id,omitempty"`
	ErrorMessage     string    `json:"error_message,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ReceiptJobEvent struct {
	ID        string    `json:"id"`
	JobID     string    `json:"job_id"`
	EventType string    `json:"event_type"`
	Stage     string    `json:"stage"`
	Attempt   int       `json:"attempt"`
	Status    string    `json:"status"`
	Message   string    `json:"message,omitempty"`
	Metadata  string    `json:"metadata,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateReceiptJobReq struct {
	ImageBase64    string `json:"image_base64"`
	OriginalSize   int64  `json:"original_size"`
	CompressedSize int64  `json:"compressed_size"`
}

type ReceiptJobResponse struct {
	JobID  string            `json:"job_id"`
	Job    ReceiptJob        `json:"job,omitempty"`
	Events []ReceiptJobEvent `json:"events,omitempty"`
}

type AINormalizationResult struct {
	Type          string   `json:"type"`
	DateTime      string   `json:"date_time"`
	Amount        int64    `json:"amount"`
	Currency      string   `json:"currency"`
	FromAccountID *string  `json:"from_account_id,omitempty"`
	ToAccountID   *string  `json:"to_account_id,omitempty"`
	CategoryID    *string  `json:"category_id,omitempty"`
	Note          string   `json:"note,omitempty"`
	PayeePayer    string   `json:"payee_payer,omitempty"`
	PaymentType   string   `json:"payment_type,omitempty"`
	WarrantyUntil *string  `json:"warranty_until,omitempty"`
	Status        string   `json:"status,omitempty"`
	PlaceName     string   `json:"place_name,omitempty"`
	PlaceLat      *float64 `json:"place_lat,omitempty"`
	PlaceLng      *float64 `json:"place_lng,omitempty"`
	Labels        string   `json:"labels,omitempty"`
}

type OCRResponse struct {
	ParsedResults []struct {
		ParsedText   string `json:"ParsedText"`
		ErrorMessage string `json:"ErrorMessage"`
		ErrorDetails string `json:"ErrorDetails"`
	} `json:"ParsedResults"`
	OCRExitCode           int    `json:"OCRExitCode"`
	IsErroredOnProcessing bool   `json:"IsErroredOnProcessing"`
	ErrorMessage          string `json:"ErrorMessage"`
	ErrorDetails          string `json:"ErrorDetails"`
}
