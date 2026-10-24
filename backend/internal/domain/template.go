package domain

import "time"

type Template struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Type          string       `json:"type"`
	Amount        float64      `json:"amount"`
	Currency      string       `json:"currency"`
	FromAccountID *string      `json:"from_account_id,omitempty"`
	ToAccountID   *string      `json:"to_account_id,omitempty"`
	CategoryID    *string      `json:"category_id,omitempty"`
	Labels        []string     `json:"labels,omitempty"`
	Note          string       `json:"note,omitempty"`
	PayeePayer    string       `json:"payee_payer,omitempty"`
	PaymentType   string       `json:"payment_type,omitempty"`
	WarrantyUntil *time.Time   `json:"warranty_until,omitempty"`
	Status        string       `json:"status,omitempty"`
	PlaceName     string       `json:"place_name,omitempty"`
	PlaceLat      *float64     `json:"place_lat,omitempty"`
	PlaceLng      *float64     `json:"place_lng,omitempty"`
	Attachments   []Attachment `json:"attachments,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type CreateTemplateInput struct {
	Name          string
	Type          string
	Amount        float64
	Currency      string
	FromAccountID *string
	ToAccountID   *string
	CategoryID    *string
	Labels        []string
	Note          string
	PayeePayer    string
	PaymentType   string
	WarrantyUntil *time.Time
	Status        string
	PlaceName     string
	PlaceLat      *float64
	PlaceLng      *float64
	Attachments   []Attachment
}

type UpdateTemplateInput struct {
	Name          *string
	Type          *string
	Amount        *float64
	Currency      *string
	FromAccountID *string
	ToAccountID   *string
	CategoryID    *string
	Labels        []string
	Note          *string
	PayeePayer    *string
	PaymentType   *string
	WarrantyUntil *time.Time
	Status        *string
	PlaceName     *string
	PlaceLat      *float64
	PlaceLng      *float64
	Attachments   []Attachment
}
