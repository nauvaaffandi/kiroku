package backup

import (
	"strings"

	"github.com/xuri/excelize/v2"
)

const formatVersion = "1.0"

const (
	sheetMetadata     = "Metadata"
	sheetTransactions = "Transactions"
	sheetAccounts     = "Accounts"
	sheetCategories   = "Categories"
	sheetAttachments  = "Attachments"
)

var txColumns = []string{
	"id", "type", "date_time", "amount", "currency",
	"from_account_id", "to_account_id", "category_id", "labels", "note",
	"payee_payer", "payment_type", "warranty_until", "status", "place_name",
	"place_lat", "place_lng", "attachments", "template_id", "created_at", "updated_at",
}

var accountColumns = []string{
	"id", "name", "currency", "type", "icon", "color", "archived", "created_at", "updated_at",
}

var categoryColumns = []string{
	"id", "parent_id", "name", "icon", "color", "kind", "sort_order", "created_at", "updated_at",
}

type attachment struct {
	TransactionID string
	Path          string
	ContentType   string
	Size          int64
	SHA256        string
	Data          []byte
	chunks        []string
}

func readSheet(f *excelize.File, name string) ([][]string, error) {
	rows, err := f.GetRows(name)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func headerIndex(header []string) map[string]int {
	m := make(map[string]int, len(header))
	for i, h := range header {
		m[strings.ToLower(strings.TrimSpace(h))] = i
	}
	return m
}

func cellAt(row []string, idx map[string]int, name string) string {
	i, ok := idx[strings.ToLower(name)]
	if !ok || i >= len(row) {
		return ""
	}
	return row[i]
}
