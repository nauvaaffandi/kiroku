package domain

import "time"

type TransactionType string

const (
	TxnTypeIncome  TransactionType = "income"
	TxnTypeExpense TransactionType = "expense"
	TxnTypeTransfer TransactionType = "transfer"
)

type TransactionStatus string

const (
	TxnStatusCleared  TransactionStatus = "cleared"
	TxnStatusPending  TransactionStatus = "pending"
	TxnStatusReconciled TransactionStatus = "reconciled"
)

type Attachment struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Type     string `json:"type"`
	Size     int64  `json:"size"`
}

type Transaction struct {
	ID            string     `json:"id"`
	Type          string     `json:"type"`
	DateTime      time.Time  `json:"date_time"`
	Amount        float64    `json:"amount"`
	Currency      string     `json:"currency"`
	FromAccountID *string    `json:"from_account_id,omitempty"`
	ToAccountID   *string    `json:"to_account_id,omitempty"`
	CategoryID    *string    `json:"category_id,omitempty"`
	Labels        []string   `json:"labels,omitempty"`
	Note          string     `json:"note,omitempty"`
	PayeePayer    string     `json:"payee_payer,omitempty"`
	PaymentType   string     `json:"payment_type,omitempty"`
	WarrantyUntil *time.Time `json:"warranty_until,omitempty"`
	Status        string     `json:"status,omitempty"`
	PlaceName     string     `json:"place_name,omitempty"`
	PlaceLat      *float64   `json:"place_lat,omitempty"`
	PlaceLng      *float64   `json:"place_lng,omitempty"`
	Attachments   []Attachment `json:"attachments,omitempty"`
	TemplateID    *string    `json:"template_id,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

type CreateTransactionInput struct {
	Type          string
	DateTime      time.Time
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
	TemplateID    *string
}

type UpdateTransactionInput struct {
	Type          *string
	DateTime      *time.Time
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
	TemplateID    *string
}

type TransactionFilter struct {
	Types        []string
	AccountID    *string
	CategoryID   *string
	Currency     *string
	Status       *string
	DateFrom     *time.Time
	DateTo       *time.Time
	MinAmount    *float64
	MaxAmount    *float64
	SearchQuery  string
	Sort         string
	Order        string
	Page         int
	Limit        int
}

type SummaryItem struct {
	Currency string  `json:"currency"`
	Income   float64 `json:"income"`
	Expense  float64 `json:"expense"`
	Net      float64 `json:"net"`
}

type BatchAction string

const (
	BatchDelete       BatchAction = "delete"
	BatchRecategorize BatchAction = "recategorize"
	BatchRelabel      BatchAction = "relabel"
)

type BatchInput struct {
	IDs        []string
	Action     BatchAction
	CategoryID *string
	Labels     []string
}
