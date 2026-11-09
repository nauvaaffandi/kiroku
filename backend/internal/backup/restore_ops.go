package backup

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

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
