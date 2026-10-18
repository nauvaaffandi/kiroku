package backup

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"kiroku/backend/internal/db"
)

func restoreXLSX(database *db.DB, data []byte) error {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("open workbook: %w", err)
	}
	defer f.Close()

	meta := sheetRows(f, sheetMetadata)
	if len(meta) < 2 {
		return fmt.Errorf("invalid backup: missing metadata sheet")
	}
	hasVersion := false
	for _, r := range meta[1:] {
		if len(r) >= 1 && r[0] == "format_version" {
			hasVersion = true
			break
		}
	}
	if !hasVersion {
		return fmt.Errorf("invalid backup: missing format version")
	}

	txRows := sheetRows(f, sheetTransactions)
	if len(txRows) == 0 {
		return fmt.Errorf("invalid backup: missing transactions sheet")
	}
	txIdx := headerIndex(txRows[0])
	for _, col := range []string{"id", "type", "date_time", "amount", "currency"} {
		if _, ok := txIdx[col]; !ok {
			return fmt.Errorf("invalid backup: missing column %q", col)
		}
	}

	accounts := parseAccountsSheet(sheetRows(f, sheetAccounts))
	categories := parseCategoriesSheet(sheetRows(f, sheetCategories))
	atts, err := parseAttachmentsSheet(sheetRows(f, sheetAttachments))
	if err != nil {
		return err
	}

	tx, err := database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := restoreAccounts(tx, accounts); err != nil {
		return err
	}
	if err := restoreCategories(tx, categories); err != nil {
		return err
	}
	if err := restoreTransactions(tx, txRows, txIdx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return writeAttachments(atts)
}

func sheetRows(f *excelize.File, name string) [][]string {
	if idx, err := f.GetSheetIndex(name); err != nil || idx < 0 {
		return nil
	}
	rows, err := f.GetRows(name)
	if err != nil {
		return nil
	}
	return rows
}

type accountRecord struct {
	ID, Name, Currency, Type, Icon, Color string
	Archived                              bool
	CreatedAt, UpdatedAt                  string
}

func parseAccountsSheet(rows [][]string) []accountRecord {
	if len(rows) < 2 {
		return nil
	}
	idx := headerIndex(rows[0])
	var out []accountRecord
	for _, r := range rows[1:] {
		if cellAt(r, idx, "id") == "" {
			continue
		}
		out = append(out, accountRecord{
			ID:        cellAt(r, idx, "id"),
			Name:      cellAt(r, idx, "name"),
			Currency:  cellAt(r, idx, "currency"),
			Type:      cellAt(r, idx, "type"),
			Icon:      cellAt(r, idx, "icon"),
			Color:     cellAt(r, idx, "color"),
			Archived:  cellAt(r, idx, "archived") == "1" || strings.EqualFold(cellAt(r, idx, "archived"), "true"),
			CreatedAt: cellAt(r, idx, "created_at"),
			UpdatedAt: cellAt(r, idx, "updated_at"),
		})
	}
	return out
}

type categoryRecord struct {
	ID, ParentID, Name, Icon, Color, Kind string
	SortOrder                             int64
	CreatedAt, UpdatedAt                  string
}

func parseCategoriesSheet(rows [][]string) []categoryRecord {
	if len(rows) < 2 {
		return nil
	}
	idx := headerIndex(rows[0])
	var out []categoryRecord
	for _, r := range rows[1:] {
		if cellAt(r, idx, "id") == "" {
			continue
		}
		so, _ := strconv.ParseInt(cellAt(r, idx, "sort_order"), 10, 64)
		out = append(out, categoryRecord{
			ID:        cellAt(r, idx, "id"),
			ParentID:  cellAt(r, idx, "parent_id"),
			Name:      cellAt(r, idx, "name"),
			Icon:      cellAt(r, idx, "icon"),
			Color:     cellAt(r, idx, "color"),
			Kind:      cellAt(r, idx, "kind"),
			SortOrder: so,
			CreatedAt: cellAt(r, idx, "created_at"),
			UpdatedAt: cellAt(r, idx, "updated_at"),
		})
	}
	return out
}

func restoreAccounts(tx *sql.Tx, records []accountRecord) error {
	for _, a := range records {
		_, err := tx.Exec(`
			INSERT INTO accounts (id, name, currency, type, icon, color, archived, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name = excluded.name, currency = excluded.currency, type = excluded.type,
				icon = excluded.icon, color = excluded.color, archived = excluded.archived,
				updated_at = excluded.updated_at
		`, a.ID, a.Name, a.Currency, a.Type, nullIfEmpty(a.Icon), nullIfEmpty(a.Color),
			boolInt(a.Archived), timeOrNow(a.CreatedAt), timeOrNow(a.UpdatedAt))
		if err != nil {
			return fmt.Errorf("restore account %s: %w", a.ID, err)
		}
	}
	return nil
}

func restoreCategories(tx *sql.Tx, records []categoryRecord) error {
	for _, c := range records {
		_, err := tx.Exec(`
			INSERT INTO categories (id, name, icon, color, kind, sort_order, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name = excluded.name, icon = excluded.icon, color = excluded.color,
				kind = excluded.kind, sort_order = excluded.sort_order, updated_at = excluded.updated_at
		`, c.ID, c.Name, nullIfEmpty(c.Icon), nullIfEmpty(c.Color), nullIfEmpty(c.Kind),
			c.SortOrder, timeOrNow(c.CreatedAt), timeOrNow(c.UpdatedAt))
		if err != nil {
			return fmt.Errorf("restore category %s: %w", c.ID, err)
		}
	}
	for _, c := range records {
		if c.ParentID == "" {
			continue
		}
		if _, err := tx.Exec(`UPDATE categories SET parent_id = ? WHERE id = ?`, c.ParentID, c.ID); err != nil {
			return fmt.Errorf("restore category parent %s: %w", c.ID, err)
		}
	}
	return nil
}

func restoreTransactions(tx *sql.Tx, rows [][]string, idx map[string]int) error {
	placeholderTemplates := map[string]bool{}
	for _, r := range rows[1:] {
		id := cellAt(r, idx, "id")
		if id == "" {
			continue
		}
		templateID := cellAt(r, idx, "template_id")
		if templateID != "" && !placeholderTemplates[templateID] {
			placeholderTemplates[templateID] = true
			if _, err := tx.Exec(`
				INSERT OR IGNORE INTO templates (id, name, type, amount, currency)
				VALUES (?, 'restored', 'expense', 0, 'JPY')
			`, templateID); err != nil {
				return fmt.Errorf("restore template %s: %w", templateID, err)
			}
		}

		amount, err := strconv.ParseInt(strings.TrimSpace(cellAt(r, idx, "amount")), 10, 64)
		if err != nil {
			return fmt.Errorf("restore transaction %s: invalid amount: %w", id, err)
		}

		_, err = tx.Exec(`
			INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, place_lat, place_lng, attachments, template_id, created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NULL)
			ON CONFLICT(id) DO UPDATE SET
				type = excluded.type, date_time = excluded.date_time, amount = excluded.amount,
				currency = excluded.currency, from_account_id = excluded.from_account_id,
				to_account_id = excluded.to_account_id, category_id = excluded.category_id,
				labels = excluded.labels, note = excluded.note, payee_payer = excluded.payee_payer,
				payment_type = excluded.payment_type, warranty_until = excluded.warranty_until,
				status = excluded.status, place_name = excluded.place_name, place_lat = excluded.place_lat,
				place_lng = excluded.place_lng, attachments = excluded.attachments,
				template_id = excluded.template_id, created_at = excluded.created_at,
				updated_at = excluded.updated_at, deleted_at = NULL
		`,
			id,
			cellAt(r, idx, "type"),
			cellAt(r, idx, "date_time"),
			amount,
			cellAt(r, idx, "currency"),
			nullIfEmpty(cellAt(r, idx, "from_account_id")),
			nullIfEmpty(cellAt(r, idx, "to_account_id")),
			nullIfEmpty(cellAt(r, idx, "category_id")),
			nullIfEmpty(cellAt(r, idx, "labels")),
			nullIfEmpty(cellAt(r, idx, "note")),
			nullIfEmpty(cellAt(r, idx, "payee_payer")),
			nullIfEmpty(cellAt(r, idx, "payment_type")),
			nullIfEmpty(cellAt(r, idx, "warranty_until")),
			orDefault(cellAt(r, idx, "status"), "cleared"),
			nullIfEmpty(cellAt(r, idx, "place_name")),
			nullFloatOrNil(cellAt(r, idx, "place_lat")),
			nullFloatOrNil(cellAt(r, idx, "place_lng")),
			nullIfEmpty(cellAt(r, idx, "attachments")),
			nullIfEmpty(templateID),
			timeOrNow(cellAt(r, idx, "created_at")),
			timeOrNow(cellAt(r, idx, "updated_at")),
		)
		if err != nil {
			return fmt.Errorf("restore transaction %s: %w", id, err)
		}
	}
	return nil
}

func parseAttachmentsSheet(rows [][]string) ([]attachment, error) {
	if len(rows) < 2 {
		return nil, nil
	}
	idx := headerIndex(rows[0])
	grouped := map[string]*attachment{}
	var order []string
	for _, r := range rows[1:] {
		txID := cellAt(r, idx, "transaction_id")
		path := cellAt(r, idx, "path")
		if txID == "" || path == "" {
			continue
		}
		key := txID + "\x00" + path
		a, ok := grouped[key]
		if !ok {
			size, _ := strconv.ParseInt(cellAt(r, idx, "size"), 10, 64)
			a = &attachment{
				TransactionID: txID,
				Path:          path,
				ContentType:   cellAt(r, idx, "content_type"),
				Size:          size,
				SHA256:        cellAt(r, idx, "sha256"),
			}
			grouped[key] = a
			order = append(order, key)
		}
		chunkIdx, _ := strconv.Atoi(cellAt(r, idx, "chunk_index"))
		for len(a.chunks) <= chunkIdx {
			a.chunks = append(a.chunks, "")
		}
		a.chunks[chunkIdx] = cellAt(r, idx, "data")
	}
	var out []attachment
	for _, key := range order {
		a := grouped[key]
		encoded := strings.Join(a.chunks, "")
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return nil, fmt.Errorf("attachment %s: invalid base64: %w", a.Path, err)
		}
		if a.SHA256 != "" {
			sum := sha256.Sum256(decoded)
			if !strings.EqualFold(a.SHA256, hex.EncodeToString(sum[:])) {
				return nil, fmt.Errorf("attachment %s: checksum mismatch", a.Path)
			}
		}
		a.Data = decoded
		out = append(out, *a)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out, nil
}

func writeAttachments(atts []attachment) error {
	for _, a := range atts {
		rel := filepath.Clean(a.Path)
		if filepath.IsAbs(rel) || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return fmt.Errorf("attachment %s: unsafe path", a.Path)
		}
		full := filepath.Join("data", filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(full, a.Data, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

func nullFloatOrNil(s string) interface{} {
	if s == "" {
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return v
}

func timeOrNow(s string) string {
	if s == "" {
		return time.Now().Format("2006-01-02 15:04:05")
	}
	return s
}
