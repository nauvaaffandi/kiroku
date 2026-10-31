package sqlite

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"kiroku/backend/internal/domain"
)

func rowToImport(headers, row []string) domain.ImportRow {
	idx := func(name string) int {
		for i, h := range headers {
			if strings.EqualFold(h, name) {
				return i
			}
		}
		return -1
	}
	get := func(name string) string {
		i := idx(name)
		if i >= 0 && i < len(row) {
			return strings.TrimSpace(row[i])
		}
		return ""
	}
	amount := 0.0
	if v := get("amount"); v != "" {
		fmt.Sscanf(v, "%f", &amount)
	}
	return domain.ImportRow{
		DateTime:      get("date_time"),
		Type:          get("type"),
		Amount:        amount,
		Currency:      get("currency"),
		FromAccount:   get("from_account"),
		ToAccount:     get("to_account"),
		Category:      get("category"),
		Labels:        splitLabels(get("labels")),
		Note:          get("note"),
		PayeePayer:    get("payee_payer"),
		PaymentType:   get("payment_type"),
		WarrantyUntil: get("warranty_until"),
		Status:        get("status"),
		Place:         get("place"),
		Attachments:   splitLabels(get("attachments")),
	}
}

func splitLabels(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func (r *ImpexpRepo) storePreview(id string, rows []domain.ImportRow, matches struct {
	Categories []domain.CategoryMatch `json:"categories"`
	Accounts   []domain.AccountMatch  `json:"accounts"`
}) {
	data, _ := json.Marshal(map[string]interface{}{
		"rows":    rows,
		"matches": matches,
	})
	r.db.Exec(`INSERT OR REPLACE INTO import_log (id, data, created_at) VALUES (?, ?, ?)`, id, string(data), time.Now())
}

func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)
	d := make([][]int, la+1)
	for i := range d {
		d[i] = make([]int, lb+1)
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			d[i][j] = min3(d[i-1][j]+1, d[i][j-1]+1, d[i-1][j-1]+cost)
		}
	}
	return d[la][lb]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func rowHash(r domain.ImportRow) string {
	h := sha256.New()
	h.Write([]byte(r.DateTime + r.Type + r.Category + r.Note))
	return hex.EncodeToString(h.Sum(nil))
}

func parseTime(s string) time.Time {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return t
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t
	}
	return time.Now()
}

func defaultStatus(s string) string {
	if s == "" {
		return "cleared"
	}
	return s
}