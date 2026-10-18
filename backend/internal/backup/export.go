package backup

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"kiroku/backend/internal/config"
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

func nullStr(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func nullFloat(nf sql.NullFloat64) string {
	if nf.Valid {
		return strconv.FormatFloat(nf.Float64, 'f', -1, 64)
	}
	return ""
}

func nullTimeStr(nt sql.NullTime) string {
	if nt.Valid {
		return nt.Time.Format(time.RFC3339)
	}
	return ""
}

func scanTransactions(rows *sql.Rows) ([][]string, error) {
	defer rows.Close()
	var out [][]string
	for rows.Next() {
		var (
			id, typ, dateTime, currency          string
			amount                               int64
			fromAcc, toAcc, catID, labels, note  sql.NullString
			payee, paymentType, warranty, status sql.NullString
			placeName, attachments, templateID   sql.NullString
			placeLat, placeLng                   sql.NullFloat64
			createdAt, updatedAt                 sql.NullTime
		)
		if err := rows.Scan(&id, &typ, &dateTime, &amount, &currency,
			&fromAcc, &toAcc, &catID, &labels, &note,
			&payee, &paymentType, &warranty, &status, &placeName,
			&placeLat, &placeLng, &attachments, &templateID, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		out = append(out, []string{
			id, typ, dateTime, strconv.FormatInt(amount, 10), currency,
			nullStr(fromAcc), nullStr(toAcc), nullStr(catID), nullStr(labels), nullStr(note),
			nullStr(payee), nullStr(paymentType), nullStr(warranty), nullStr(status), nullStr(placeName),
			nullFloat(placeLat), nullFloat(placeLng), nullStr(attachments), nullStr(templateID),
			nullTimeStr(createdAt), nullTimeStr(updatedAt),
		})
	}
	return out, rows.Err()
}

func scanAccounts(rows *sql.Rows) ([][]string, error) {
	defer rows.Close()
	var out [][]string
	for rows.Next() {
		var (
			id, name, currency, typ, icon, color string
			archived                             int64
			createdAt, updatedAt                 sql.NullTime
		)
		if err := rows.Scan(&id, &name, &currency, &typ, &icon, &color, &archived, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		out = append(out, []string{id, name, currency, typ, icon, color,
			strconv.FormatInt(archived, 10), nullTimeStr(createdAt), nullTimeStr(updatedAt)})
	}
	return out, rows.Err()
}

func scanCategories(rows *sql.Rows) ([][]string, error) {
	defer rows.Close()
	var out [][]string
	for rows.Next() {
		var (
			id, name, icon, color, kind string
			parentID                    sql.NullString
			sortOrder                   int64
			createdAt, updatedAt        sql.NullTime
		)
		if err := rows.Scan(&id, &parentID, &name, &icon, &color, &kind, &sortOrder, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		out = append(out, []string{id, nullStr(parentID), name, icon, color, kind,
			strconv.FormatInt(sortOrder, 10), nullTimeStr(createdAt), nullTimeStr(updatedAt)})
	}
	return out, rows.Err()
}

func writeSheet(f *excelize.File, name string, rows [][]string) {
	for i, row := range rows {
		cell, err := excelize.CoordinatesToCellName(1, i+1)
		if err != nil {
			continue
		}
		_ = f.SetSheetRow(name, cell, &row)
	}
}

func fingerprint(txs, accounts, categories [][]string, atts []attachment) string {
	h := sha256.New()
	writeFingerprintRows(h, txs)
	writeFingerprintRows(h, accounts)
	writeFingerprintRows(h, categories)
	for _, a := range atts {
		fmt.Fprintf(h, "%s|%s|%s|%d|%s|%x\n", a.TransactionID, a.Path, a.ContentType, a.Size, a.SHA256, a.Data)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func writeFingerprintRows(w io.Writer, rows [][]string) {
	sorted := make([][]string, len(rows))
	copy(sorted, rows)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i][0] < sorted[j][0] })
	for _, r := range sorted {
		fmt.Fprintln(w, strings.Join(r, "\x1f"))
	}
}

func buildWorkbook(scope, dateFrom, dateTo string, txs, accounts, categories [][]string, atts []attachment) (*excelize.File, error) {
	f := excelize.NewFile()
	f.SetSheetName("Sheet1", sheetMetadata)

	writeSheet(f, sheetMetadata, [][]string{
		{"key", "value"},
		{"format_version", formatVersion},
		{"scope", scope},
		{"date_from", dateFrom},
		{"date_to", dateTo},
		{"generated_at", time.Now().Format(time.RFC3339)},
		{"transaction_count", strconv.Itoa(len(txs))},
		{"account_count", strconv.Itoa(len(accounts))},
		{"category_count", strconv.Itoa(len(categories))},
		{"attachment_count", strconv.Itoa(len(atts))},
	})

	for _, sheet := range []string{sheetTransactions, sheetAccounts, sheetCategories, sheetAttachments} {
		if _, err := f.NewSheet(sheet); err != nil {
			return nil, err
		}
	}

	if len(txs) > 0 {
		writeSheet(f, sheetTransactions, append([][]string{txColumns}, txs...))
	}
	if len(accounts) > 0 {
		writeSheet(f, sheetAccounts, append([][]string{accountColumns}, accounts...))
	}
	if len(categories) > 0 {
		writeSheet(f, sheetCategories, append([][]string{categoryColumns}, categories...))
	}
	if len(atts) > 0 {
		writeSheet(f, sheetAttachments, attachmentRows(atts))
	}
	return f, nil
}

const attachmentChunkSize = 20000

func attachmentRows(atts []attachment) [][]string {
	rows := [][]string{{"transaction_id", "path", "content_type", "size", "sha256", "chunk_index", "data"}}
	sorted := make([]attachment, len(atts))
	copy(sorted, atts)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].TransactionID != sorted[j].TransactionID {
			return sorted[i].TransactionID < sorted[j].TransactionID
		}
		return sorted[i].Path < sorted[j].Path
	})
	for _, a := range sorted {
		encoded := base64.StdEncoding.EncodeToString(a.Data)
		chunks := splitChunks(encoded, attachmentChunkSize)
		if len(chunks) == 0 {
			chunks = []string{""}
		}
		for i, c := range chunks {
			rows = append(rows, []string{
				a.TransactionID, a.Path, a.ContentType,
				strconv.FormatInt(a.Size, 10), a.SHA256, strconv.Itoa(i), c,
			})
		}
	}
	return rows
}

func splitChunks(s string, size int) []string {
	if s == "" {
		return nil
	}
	var out []string
	for len(s) > size {
		out = append(out, s[:size])
		s = s[size:]
	}
	out = append(out, s)
	return out
}

func writeXLSX(path string, f *excelize.File) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	buf, err := f.WriteToBuffer()
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".tmp-*.xlsx")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if _, err := tmp.Write(buf.Bytes()); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

func collectAttachments(cfg *config.Config, txs [][]string) []attachment {
	if cfg == nil {
		return nil
	}
	seen := map[string]bool{}
	var atts []attachment
	for _, tx := range txs {
		id := tx[0]
		raw := tx[17]
		if raw == "" {
			continue
		}
		for _, p := range parseAttachmentPaths(raw) {
			if seen[p] {
				continue
			}
			seen[p] = true
			full := resolveAttachmentPath(cfg, p)
			data, err := os.ReadFile(full)
			if err != nil {
				continue
			}
			sum := sha256.Sum256(data)
			atts = append(atts, attachment{
				TransactionID: id,
				Path:          p,
				ContentType:   contentTypeFor(p),
				Size:          int64(len(data)),
				SHA256:        hex.EncodeToString(sum[:]),
				Data:          data,
			})
		}
	}
	return atts
}

func parseAttachmentPaths(raw string) []string {
	var arr []string
	if err := json.Unmarshal([]byte(raw), &arr); err == nil {
		return arr
	}
	if strings.TrimSpace(raw) != "" {
		return []string{strings.TrimSpace(raw)}
	}
	return nil
}

func resolveAttachmentPath(cfg *config.Config, p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	if strings.HasPrefix(p, "attachments"+string(filepath.Separator)) || p == "attachments" {
		return filepath.Join(cfg.DataDir, filepath.FromSlash(p))
	}
	return filepath.Join(cfg.DataDir, "attachments", filepath.Base(p))
}

func contentTypeFor(p string) string {
	switch strings.ToLower(filepath.Ext(p)) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
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
