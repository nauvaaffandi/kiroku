package impexp

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"kiroku/backend/internal/models"
)

type readSeeker struct {
	data []byte
	pos  int
}

func (r *readSeeker) Read(p []byte) (n int, err error) {
	n = copy(p, r.data[r.pos:])
	r.pos += n
	if n == 0 {
		return 0, io.EOF
	}
	return n, nil
}

func (r *readSeeker) Seek(offset int64, whence int) (int64, error) {
	var newPos int64
	switch whence {
	case io.SeekStart:
		newPos = offset
	case io.SeekCurrent:
		newPos = int64(r.pos) + offset
	case io.SeekEnd:
		newPos = int64(len(r.data)) + offset
	}
	if newPos < 0 {
		return 0, fmt.Errorf("invalid seek")
	}
	r.pos = int(newPos)
	return newPos, nil
}

func (h *Handler) parseCSV(r io.Reader) ([][]string, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1
	return reader.ReadAll()
}

func (h *Handler) parseExcel(r io.Reader) ([][]string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	f, err := excelize.OpenReader(&readSeeker{data: data})
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (h *Handler) rowToImport(headers, row []string) models.ImportRow {
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

	return models.ImportRow{
		DateTime:      get("date_time"),
		Type:          get("type"),
		Amount:        get("amount"),
		Currency:      get("currency"),
		FromAccount:   get("from_account"),
		ToAccount:     get("to_account"),
		Category:      get("category"),
		Labels:        get("labels"),
		Note:          get("note"),
		PayeePayer:    get("payee_payer"),
		PaymentType:   get("payment_type"),
		WarrantyUntil: get("warranty_until"),
		Status:        get("status"),
		Place:         get("place"),
		Attachments:   get("attachments"),
	}
}

func (h *Handler) storePreview(id string, rows []models.ImportRow, matches struct {
	Categories []models.CategoryMatch `json:"categories"`
	Accounts   []models.AccountMatch  `json:"accounts"`
}) {
	data, _ := json.Marshal(map[string]interface{}{
		"rows":    rows,
		"matches": matches,
	})
	h.db.Exec("INSERT OR REPLACE INTO import_log (id, data, created_at) VALUES (?, ?, ?)", id, string(data), time.Now())
}
