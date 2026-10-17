package models

type Group struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Color    string  `json:"color"`
	Kind     string  `json:"kind,omitempty"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Count    int     `json:"count"`
	Children []Group `json:"children,omitempty"`
}

type MonthData struct {
	Month    string  `json:"month"`
	Income   float64 `json:"income"`
	Expense  float64 `json:"expense"`
	Net      float64 `json:"net"`
	Currency string  `json:"currency"`
}
