package backup

import (
	"testing"
)

func TestBackupExportWithNullKind(t *testing.T) {
	database := testDB(t)

	_, err := database.Exec(`INSERT INTO categories (id, name, sort_order) VALUES ('cat-nul', 'NullKind', 0)`)
	if err != nil {
		t.Fatalf("insert null kind: %v", err)
	}
	_, err = database.Exec(`INSERT INTO categories (id, name, kind, sort_order) VALUES ('cat-must', 'Must', 'must', 1)`)
	if err != nil {
		t.Fatalf("insert must kind: %v", err)
	}

	rows, err := database.Query(`SELECT id, parent_id, name, icon, color, kind, sort_order, created_at, updated_at FROM categories ORDER BY sort_order`)
	if err != nil {
		t.Fatalf("query categories: %v", err)
	}

	result, err := scanCategories(rows)
	if err != nil {
		t.Fatalf("scanCategories: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(result))
	}

	if result[0][5] != "" {
		t.Fatalf("expected empty kind for null, got %q", result[0][5])
	}
	if result[1][5] != "must" {
		t.Fatalf("expected 'must' kind, got %q", result[1][5])
	}
}
