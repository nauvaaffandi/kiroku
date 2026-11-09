package backup

import (
	"bytes"
	"fmt"

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
