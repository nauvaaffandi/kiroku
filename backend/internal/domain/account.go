package domain

import "time"

type AccountType string

const (
	AccountTypeBank     AccountType = "bank"
	AccountTypeCash     AccountType = "cash"
	AccountTypeCard     AccountType = "card"
	AccountTypeInvest   AccountType = "investment"
	AccountTypeOther    AccountType = "other"
)

type Account struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Currency  string     `json:"currency"`
	Type      string     `json:"type"`
	Icon      string     `json:"icon"`
	Color     string     `json:"color"`
	Archived  bool       `json:"archived"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type AccountWithBalance struct {
	Account
	Balance float64 `json:"balance"`
}

type CreateAccountInput struct {
	Name     string
	Currency string
	Type     string
	Icon     string
	Color    string
}

type UpdateAccountInput struct {
	Name     *string
	Currency *string
	Type     *string
	Icon     *string
	Color    *string
	Archived *bool
}
