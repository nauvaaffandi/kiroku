package backup

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

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
