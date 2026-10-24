package sqlite

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"kiroku/backend/internal/domain"
	"kiroku/backend/internal/db"
)

type ImpexpRepo struct{ db *db.DB }

func NewImpexpRepo(database *db.DB) *ImpexpRepo { return &ImpexpRepo{db: database} }

func (r *ImpexpRepo) ExportCSV(ctx context.Context, dateFrom, dateTo, typ, categoryID, accountID string) ([][]string, error) {
	query := `SELECT t.date_time, t.type, t.amount, t.currency,
		fa.name as from_account, ta.name as to_account,
		c.name as category, t.labels, t.note, t.payee_payer,
		t.payment_type, t.warranty_until, t.status,
		t.place_name, t.attachments
		FROM transactions t
		LEFT JOIN accounts fa ON t.from_account_id = fa.id
		LEFT JOIN accounts ta ON t.to_account_id = ta.id
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.deleted_at IS NULL`
	args := []interface{}{}
	conditions := []string{}
	if dateFrom != "" {
		conditions = append(conditions, "t.date_time >= ?")
		args = append(args, dateFrom)
	}
	if dateTo != "" {
		conditions = append(conditions, "t.date_time <= ?")
		args = append(args, dateTo)
	}
	if typ != "" {
		conditions = append(conditions, "t.type = ?")
		args = append(args, typ)
	}
	if categoryID != "" {
		conditions = append(conditions, "t.category_id = ?")
		args = append(args, categoryID)
	}
	if accountID != "" {
		conditions = append(conditions, "(t.from_account_id = ? OR t.to_account_id = ?)")
		args = append(args, accountID, accountID)
	}
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY t.date_time DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	headers := []string{"date_time", "type", "amount", "currency", "from_account", "to_account", "category", "labels", "note", "payee_payer", "payment_type", "warranty_until", "status", "place", "attachments"}
	records := [][]string{headers}

	for rows.Next() {
		var dateTime time.Time
		var typStr, currency, fromAcc, toAcc, category, labels, note, payeePayer, paymentType, warrantyUntil, status, placeName, attachments sql.NullString
		var amount int64
		if err := rows.Scan(&dateTime, &typStr, &amount, &currency, &fromAcc, &toAcc, &category, &labels, &note, &payeePayer, &paymentType, &warrantyUntil, &status, &placeName, &attachments); err != nil {
			continue
		}
		record := []string{
			dateTime.Format(time.RFC3339),
			typStr.String,
			fmt.Sprintf("%d", amount),
			currency.String,
			fromAcc.String,
			toAcc.String,
			category.String,
			labels.String,
			note.String,
			payeePayer.String,
			paymentType.String,
			warrantyUntil.String,
			status.String,
			placeName.String,
			attachments.String,
		}
		records = append(records, record)
	}
	return records, rows.Err()
}

func (r *ImpexpRepo) ImportPreview(ctx context.Context, headers []string, dataRows [][]string) (*domain.ImportPreview, error) {
	if len(dataRows) > 10 {
		dataRows = dataRows[:10]
	}

	preview := &domain.ImportPreview{PreviewID: genID()}
	catCache := map[string]string{}
	accCache := map[string]string{}

	for _, row := range dataRows {
		if len(row) < len(headers) {
			continue
		}
		rec := rowToImport(headers, row)
		preview.Rows = append(preview.Rows, rec)

		if rec.Category != "" {
			if _, ok := catCache[rec.Category]; !ok {
				match := r.fuzzyMatchCategory(rec.Category)
				preview.Categories = append(preview.Categories, match)
				catCache[rec.Category] = match.MatchID
			}
		}
		if rec.FromAccount != "" {
			if _, ok := accCache[rec.FromAccount]; !ok {
				match := r.fuzzyMatchAccount(rec.FromAccount)
				preview.Accounts = append(preview.Accounts, match)
				accCache[rec.FromAccount] = match.MatchID
			}
		}
		if rec.ToAccount != "" {
			if _, ok := accCache[rec.ToAccount]; !ok {
				match := r.fuzzyMatchAccount(rec.ToAccount)
				preview.Accounts = append(preview.Accounts, match)
				accCache[rec.ToAccount] = match.MatchID
			}
		}
	}

	matches := struct {
		Categories []domain.CategoryMatch `json:"categories"`
		Accounts   []domain.AccountMatch  `json:"accounts"`
	}{
		Categories: preview.Categories,
		Accounts:   preview.Accounts,
	}
	r.storePreview(preview.PreviewID, preview.Rows, matches)
	return preview, nil
}

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

func (r *ImpexpRepo) fuzzyMatchCategory(input string) domain.CategoryMatch {
	input = strings.TrimSpace(input)
	if input == "" {
		return domain.CategoryMatch{Input: input}
	}
	best := domain.CategoryMatch{Input: input, Distance: 100}
	rows, err := r.db.Query(`SELECT id, name FROM categories WHERE archived = 0`)
	if err != nil {
		return best
	}
	defer rows.Close()
	for rows.Next() {
		var id, name string
		if rows.Scan(&id, &name) != nil {
			continue
		}
		dist := levenshtein(input, name)
		if dist < best.Distance {
			best = domain.CategoryMatch{Input: input, Matched: name, MatchID: id, Distance: dist}
		}
	}
	if best.Distance <= 2 {
		return best
	}
	id := genID()
	_, _ = r.db.Exec(`INSERT INTO categories (id, name, color, sort_order) VALUES (?, ?, ?, 0)`, id, input, randomPastel())
	return domain.CategoryMatch{Input: input, Matched: input, MatchID: id, Distance: 0}
}

func (r *ImpexpRepo) fuzzyMatchAccount(input string) domain.AccountMatch {
	input = strings.TrimSpace(input)
	if input == "" {
		return domain.AccountMatch{Input: input}
	}
	best := domain.AccountMatch{Input: input, Distance: 100}
	rows, err := r.db.Query(`SELECT id, name FROM accounts WHERE archived = 0`)
	if err != nil {
		return best
	}
	defer rows.Close()
	for rows.Next() {
		var id, name string
		if rows.Scan(&id, &name) != nil {
			continue
		}
		dist := levenshtein(input, name)
		if dist < best.Distance {
			best = domain.AccountMatch{Input: input, Matched: name, MatchID: id, Distance: dist}
		}
	}
	if best.Distance <= 2 {
		return best
	}
	id := genID()
	_, _ = r.db.Exec(`INSERT INTO accounts (id, name, type, currency, sort_order) VALUES (?, ?, 'cash', 'JPY', 0)`, id, input)
	return domain.AccountMatch{Input: input, Matched: input, MatchID: id, Distance: 0}
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

func (r *ImpexpRepo) GetPreview(ctx context.Context, previewID string) (*domain.ImportPreview, error) {
	var data string
	err := r.db.QueryRowContext(ctx, `SELECT data FROM import_log WHERE id = ?`, previewID).Scan(&data)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var preview domain.ImportPreview
	if err := json.Unmarshal([]byte(data), &preview); err != nil {
		return nil, err
	}
	return &preview, nil
}

func (r *ImpexpRepo) Confirm(ctx context.Context, previewID string) (int, int, []string, error) {
	preview, err := r.GetPreview(ctx, previewID)
	if err != nil {
		return 0, 0, nil, err
	}
	if preview == nil {
		return 0, 0, nil, sql.ErrNoRows
	}

	created := 0
	skipped := 0
	var errors []string

	for _, row := range preview.Rows {
		hash := rowHash(row)
		var exists bool
		r.db.QueryRowContext(ctx, `SELECT 1 FROM import_log WHERE id = ?`, hash).Scan(&exists)
		if exists {
			skipped++
			continue
		}

		dateTime, err := time.Parse(time.RFC3339, row.DateTime)
		if err != nil {
			dateTime = time.Now()
		}

		if row.Amount == 0 {
			errors = append(errors, row.DateTime+": amount required")
			continue
		}

		fromAccID := resolveAccountByName(r.db, row.FromAccount)
		toAccID := resolveAccountByName(r.db, row.ToAccount)
		catID := resolveCategoryByName(r.db, row.Category)

		if row.Type == "expense" && fromAccID == nil {
			errors = append(errors, row.DateTime+": expense requires from_account")
			continue
		}
		if row.Type == "income" && toAccID == nil {
			errors = append(errors, row.DateTime+": income requires to_account")
			continue
		}

		id := genID()
		warranty := parseTime(row.WarrantyUntil)
		status := defaultStatus(row.Status)
		_, err = r.db.ExecContext(ctx, `
			INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, attachments)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, id, row.Type, dateTime, row.Amount, row.Currency, fromAccID, toAccID, catID, row.Labels, row.Note, row.PayeePayer, row.PaymentType, warranty, status, row.Place, row.Attachments)
		if err != nil {
			errors = append(errors, row.DateTime+": "+err.Error())
			continue
		}

		r.db.ExecContext(ctx, `INSERT INTO import_log (id, data, created_at) VALUES (?, ?, ?)`, hash, "{}", time.Now())
		created++
	}

	return created, skipped, errors, nil
}

func resolveAccountByName(database *db.DB, name string) *string {
	if name == "" {
		return nil
	}
	var id string
	err := database.QueryRow(`SELECT id FROM accounts WHERE name = ? AND archived = 0`, name).Scan(&id)
	if err != nil {
		return nil
	}
	return &id
}

func resolveCategoryByName(database *db.DB, name string) *string {
	if name == "" {
		return nil
	}
	var id string
	err := database.QueryRow(`SELECT id FROM categories WHERE name = ? AND archived = 0`, name).Scan(&id)
	if err != nil {
		return nil
	}
	return &id
}
