package models

import "time"

type Template struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Type          string     `json:"type"`
	Amount        *int64     `json:"amount,omitempty"`
	Currency      string     `json:"currency,omitempty"`
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
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CreateTemplateReq struct {
	Name          string     `json:"name"`
	Type          string     `json:"type"`
	Amount        *int64     `json:"amount,omitempty"`
	Currency      string     `json:"currency,omitempty"`
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
}

type UpdateTemplateReq struct {
	Name          *string    `json:"name,omitempty"`
	Type          *string    `json:"type,omitempty"`
	Amount        *int64     `json:"amount,omitempty"`
	Currency      *string    `json:"currency,omitempty"`
	FromAccountID *string    `json:"from_account_id,omitempty"`
	ToAccountID   *string    `json:"to_account_id,omitempty"`
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
