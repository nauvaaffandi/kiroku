package models

type ImportRow struct {
	DateTime      string `csv:"date_time"`
	Type          string `csv:"type"`
	Amount        string `csv:"amount"`
	Currency      string `csv:"currency"`
	FromAccount   string `csv:"from_account"`
	ToAccount     string `csv:"to_account"`
	Category      string `csv:"category"`
	Labels        string `csv:"labels"`
	Note          string `csv:"note"`
	PayeePayer    string `csv:"payee_payer"`
	PaymentType   string `csv:"payment_type"`
	WarrantyUntil string `csv:"warranty_until"`
	Status        string `csv:"status"`
	Place         string `csv:"place"`
	Attachments   string `csv:"attachments"`
}

type CategoryMatch struct {
	Input    string `json:"input"`
	Matched  string `json:"matched"`
	MatchID  string `json:"match_id"`
	Distance int    `json:"distance"`
	Created  bool   `json:"created"`
}

type AccountMatch struct {
	Input    string `json:"input"`
	Matched  string `json:"matched"`
	MatchID  string `json:"match_id"`
	Distance int    `json:"distance"`
	Created  bool   `json:"created"`
}

type ImportPreview struct {
	Rows    []ImportRow `json:"rows"`
	Matches struct {
		Categories []CategoryMatch `json:"categories"`
		Accounts   []AccountMatch  `json:"accounts"`
	} `json:"matches"`
	Errors    []string `json:"errors"`
	PreviewID string   `json:"preview_id"`
}

type ImportConfirmReq struct {
	PreviewID string `json:"preview_id"`
	Confirm   bool   `json:"confirm"`
}

type ImportConfirmResponse struct {
	Created int      `json:"created"`
	Skipped int      `json:"skipped"`
	Errors  []string `json:"errors"`
}
