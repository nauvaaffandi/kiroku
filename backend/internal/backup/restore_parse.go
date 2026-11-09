package backup

import (
	"strconv"
	"strings"
)

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
