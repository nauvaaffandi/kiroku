package models

type Account struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Type     string `json:"type"`
	Icon     string `json:"icon,omitempty"`
	Color    string `json:"color,omitempty"`
	Archived bool   `json:"archived"`
	Balance  int64  `json:"balance"`
}

type CreateAccountReq struct {
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Type     string `json:"type"`
	Icon     string `json:"icon,omitempty"`
	Color    string `json:"color,omitempty"`
}

type UpdateAccountReq struct {
	Name     *string `json:"name,omitempty"`
	Currency *string `json:"currency,omitempty"`
	Type     *string `json:"type,omitempty"`
	Icon     *string `json:"icon,omitempty"`
	Color    *string `json:"color,omitempty"`
	Archived *bool   `json:"archived,omitempty"`
}
