package models

import "time"

type CreateShareReq struct {
	AccountID string  `json:"account_id"`
	ExpiresAt *string `json:"expires_at,omitempty"`
}

type Share struct {
	ID             string     `json:"id"`
	AccountID      string     `json:"account_id"`
	AccountName    string     `json:"account_name"`
	Token          string     `json:"token"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	LastAccessedAt *time.Time `json:"last_accessed_at,omitempty"`
}

type SharedViewResponse struct {
	Account struct {
		ID       string  `json:"id"`
		Name     string  `json:"name"`
		Currency string  `json:"currency"`
		Balance  float64 `json:"balance"`
	} `json:"account"`
	Transactions []struct {
		ID            string  `json:"id"`
		Date          string  `json:"date"`
		Amount        float64 `json:"amount"`
		Type          string  `json:"type"`
		Note          string  `json:"note"`
		Payee         string  `json:"payee"`
		Category      string  `json:"category"`
		CategoryIcon  string  `json:"category_icon"`
		CategoryColor string  `json:"category_color"`
	} `json:"transactions"`
	Attribution string `json:"attribution"`
}

type SharedTransaction struct {
	ID            string  `json:"id"`
	Date          string  `json:"date"`
	Amount        float64 `json:"amount"`
	Type          string  `json:"type"`
	Note          string  `json:"note"`
	Payee         string  `json:"payee"`
	Category      string  `json:"category"`
	CategoryIcon  string  `json:"category_icon"`
	CategoryColor string  `json:"category_color"`
}
