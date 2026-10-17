package models

import "time"

type Transaction struct {
	ID            string     `json:"id"`
	Type          string     `json:"type"`
	DateTime      time.Time  `json:"date_time"`
	Amount        int64      `json:"amount"`
	Currency      string     `json:"currency"`
	FromAccountID *string    `json:"from_account_id,omitempty"`
	ToAccountID   *string    `json:"to_account_id,omitempty"`
	CategoryID    *string    `json:"category_id,omitempty"`
	Labels        string     `json:"labels,omitempty"`
	Note          string     `json:"note,omitempty"`
	PayeePayer    string     `json:"payee_payer,omitempty"`
	PaymentType   string     `json:"payment_type,omitempty"`
	WarrantyUntil *time.Time `json:"warranty_until,omitempty"`
	Status        string     `json:"status"`
	PlaceName     string     `json:"place_name,omitempty"`
	PlaceLat      *float64   `json:"place_lat,omitempty"`
	PlaceLng      *float64   `json:"place_lng,omitempty"`
	Attachments   string     `json:"attachments,omitempty"`
	TemplateID    *string    `json:"template_id,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CreateTransactionReq struct {
	Type          string     `json:"type"`
	DateTime      time.Time  `json:"date_time"`
	Amount        int64      `json:"amount"`
	Currency      string     `json:"currency"`
	FromAccountID *string    `json:"from_account_id,omitempty"`
	ToAccountID   *string    `json:"to_account_id,omitempty"`
	CategoryID    *string    `json:"category_id,omitempty"`
	Labels        string     `json:"labels,omitempty"`
	Note          string     `json:"note,omitempty"`
	PayeePayer    string     `json:"payee_payer,omitempty"`
	PaymentType   string     `json:"payment_type,omitempty"`
	WarrantyUntil *time.Time `json:"warranty_until,omitempty"`
	Status        string     `json:"status,omitempty"`
	PlaceName     string     `json:"place_name,omitempty"`
	PlaceLat      *float64   `json:"place_lat,omitempty"`
	PlaceLng      *float64   `json:"place_lng,omitempty"`
	Attachments   string     `json:"attachments,omitempty"`
	TemplateID    *string    `json:"template_id,omitempty"`
}

type UpdateTransactionReq struct {
	DateTime      *time.Time `json:"date_time,omitempty"`
	Amount        *int64     `json:"amount,omitempty"`
	CategoryID    *string    `json:"category_id,omitempty"`
	Labels        *string    `json:"labels,omitempty"`
	Note          *string    `json:"note,omitempty"`
	PayeePayer    *string    `json:"payee_payer,omitempty"`
	PaymentType   *string    `json:"payment_type,omitempty"`
	WarrantyUntil *time.Time `json:"warranty_until,omitempty"`
	Status        *string    `json:"status,omitempty"`
	PlaceName     *string    `json:"place_name,omitempty"`
	PlaceLat      *float64   `json:"place_lat,omitempty"`
	PlaceLng      *float64   `json:"place_lng,omitempty"`
	Attachments   *string    `json:"attachments,omitempty"`
}

type SummaryItem struct {
	Currency string `json:"currency"`
	Income   int64  `json:"income"`
	Expense  int64  `json:"expense"`
	Net      int64  `json:"net"`
}

type BatchReq struct {
	IDs        []string `json:"ids"`
	Action     string   `json:"action"`
	CategoryID *string  `json:"category_id,omitempty"`
	Labels     *string  `json:"labels,omitempty"`
}
