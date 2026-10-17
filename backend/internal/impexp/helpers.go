package impexp

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"

	"kiroku/backend/internal/models"
)

func (h *Handler) resolveAccount(name string) *string {
	if name == "" {
		return nil
	}
	var id string
	err := h.db.QueryRow("SELECT id FROM accounts WHERE name = ? AND archived = 0", name).Scan(&id)
	if err != nil {
		return nil
	}
	return &id
}

func (h *Handler) resolveCategory(name string) *string {
	if name == "" {
		return nil
	}
	var id string
	err := h.db.QueryRow("SELECT id FROM categories WHERE name = ? AND archived = 0", name).Scan(&id)
	if err != nil {
		return nil
	}
	return &id
}

func rowHash(row models.ImportRow) string {
	data := row.DateTime + "|" + row.Type + "|" + row.Amount + "|" + row.Currency + "|" + row.FromAccount + "|" + row.ToAccount + "|" + row.Category + "|" + row.Labels + "|" + row.Note
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func parseTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	return &t
}

func defaultStatus(s string) string {
	if s == "" {
		return "cleared"
	}
	return s
}

func randomPastel() string {
	return "#" + fmt.Sprintf("%06x", uuid.New().ClockSequence()%0xFFFFFF)
}
