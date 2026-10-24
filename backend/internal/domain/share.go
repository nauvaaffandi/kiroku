package domain

import "time"

type Share struct {
	ID             string     `json:"id"`
	AccountID      string     `json:"account_id"`
	AccountName    string     `json:"account_name,omitempty"`
	Token          string     `json:"token"`
	ExpiresAt      *time.Time `json:"expires_at"`
	CreatedAt      time.Time  `json:"created_at"`
	LastAccessedAt *time.Time `json:"last_accessed_at,omitempty"`
}

type CreateShareInput struct {
	AccountID string
	ExpiresAt *time.Time
}

type SharedView struct {
	Account        *AccountWithBalance `json:"account"`
	Transactions   []Transaction       `json:"transactions"`
	Attribution    string              `json:"attribution"`
}
