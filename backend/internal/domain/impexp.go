package domain

type ImportRow struct {
	DateTime      string   `json:"date_time"`
	Type          string   `json:"type"`
	Amount        float64  `json:"amount"`
	Currency      string   `json:"currency"`
	FromAccount   string   `json:"from_account"`
	ToAccount     string   `json:"to_account"`
	Category      string   `json:"category"`
	Labels        []string `json:"labels"`
	Note          string   `json:"note"`
	PayeePayer    string   `json:"payee_payer"`
	PaymentType   string   `json:"payment_type"`
	WarrantyUntil string   `json:"warranty_until"`
	Status        string   `json:"status"`
	Place         string   `json:"place"`
	Attachments   []string `json:"attachments"`
}

type CategoryMatch struct {
	Input    string `json:"input"`
	Matched  string `json:"matched"`
	Category string `json:"category"`
	Score    int    `json:"score"`
	MatchID  string `json:"match_id,omitempty"`
	Distance int    `json:"distance,omitempty"`
}

type AccountMatch struct {
	Input    string `json:"input"`
	Matched  string `json:"matched"`
	Account  string `json:"account"`
	Score    int    `json:"score"`
	MatchID  string `json:"match_id,omitempty"`
	Distance int    `json:"distance,omitempty"`
}

type ImportPreview struct {
	PreviewID  string          `json:"preview_id"`
	Rows       []ImportRow     `json:"rows"`
	Total      int             `json:"total"`
	Categories []CategoryMatch `json:"categories,omitempty"`
	Accounts   []AccountMatch  `json:"accounts,omitempty"`
}
