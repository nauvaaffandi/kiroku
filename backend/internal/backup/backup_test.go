package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"

	"kiroku/backend/internal/config"
	"kiroku/backend/internal/db"
)

func testDB(t *testing.T) *db.DB {
	t.Helper()
	dir := t.TempDir()
	os.Setenv("DATA_DIR", dir)
	t.Cleanup(func() { os.Unsetenv("DATA_DIR") })
	database, err := db.Open(dir)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { database.Close() })
	log.Info().Str("path", filepath.Join(dir, "kiroku.db")).Msg("test db opened")
	return database
}

func testScheduler(t *testing.T, database *db.DB) *Scheduler {
	t.Helper()
	root := filepath.Join(t.TempDir(), "backup", "kiroku_backup_storage")
	dataDir := filepath.Join(t.TempDir(), "data")
	os.Setenv("DATA_DIR", dataDir)
	os.Setenv("BACKUP_PROVIDERS", "git")
	t.Cleanup(func() {
		os.Unsetenv("DATA_DIR")
		os.Unsetenv("BACKUP_PROVIDERS")
	})
	return NewScheduler(database, &config.Config{}, root, "")
}

func insertTx(t *testing.T, database *db.DB, id, typ string, dateTime, created time.Time, amount int64, currency string) {
	t.Helper()
	_, err := database.Exec(`
		INSERT INTO transactions (id, type, date_time, amount, currency, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, 'cleared', ?, ?)
	`, id, typ, dateTime.Format(time.RFC3339), amount, currency, created.Format(time.RFC3339), created.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("insert tx: %v", err)
	}
}

func countRowsInSheet(t *testing.T, path, sheet string) int {
	t.Helper()
	f, err := excelize.OpenFile(path)
	if err != nil {
		t.Fatalf("open xlsx %s: %v", path, err)
	}
	defer f.Close()
	rows, err := f.GetRows(sheet)
	if err != nil {
		t.Fatalf("get rows %s: %v", sheet, err)
	}
	if len(rows) == 0 {
		return 0
	}
	return len(rows) - 1
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func today() string { return time.Now().Format("2006-01-02") }

func TestSevenDayWaitingPeriod(t *testing.T) {
	database := testDB(t)
	s := testScheduler(t, database)

	now := time.Now()
	insertTx(t, database, "tx-old", "expense", now, now.AddDate(0, 0, -10), 1000, "JPY")     // layak
	insertTx(t, database, "tx-boundary", "expense", now, now.AddDate(0, 0, -7), 2000, "JPY") // tepat 7 hari, layak
	insertTx(t, database, "tx-young", "expense", now, now.AddDate(0, 0, -2), 3000, "JPY")    // belum layak

	if err := s.RunJob(); err != nil {
		t.Fatalf("run job: %v", err)
	}

	day := today()
	path := filepath.Join(s.root, day[:4], day[5:7], day[8:10]+".xlsx")
	if !fileExists(path) {
		t.Fatalf("expected daily backup file %s to exist", path)
	}
	n := countRowsInSheet(t, path, "Transactions")
	if n < 2 {
		t.Fatalf("expected at least 2 transactions (old + boundary) in backup, got %d", n)
	}

	// Jalankan lagi: idempoten, tidak ada berkas baru/dubah.
	before := gitCommitCount(t, s.root)
	if err := s.RunJob(); err != nil {
		t.Fatalf("run job second: %v", err)
	}
	after := gitCommitCount(t, s.root)
	if after != before {
		t.Fatalf("second run should not create a new commit (%d vs %d)", after, before)
	}
}

func TestEmptyDateNoFile(t *testing.T) {
	database := testDB(t)
	s := testScheduler(t, database)

	now := time.Now()
	insertTx(t, database, "tx1", "expense", now, now.AddDate(0, 0, -10), 1000, "JPY")

	if err := s.RunJob(); err != nil {
		t.Fatalf("run job: %v", err)
	}

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	path := filepath.Join(s.root, yesterday[:4], yesterday[5:7], yesterday[8:10]+".xlsx")
	if fileExists(path) {
		t.Fatalf("empty date should not produce a backup file: %s", path)
	}
}

func TestEditReconciliation(t *testing.T) {
	database := testDB(t)
	s := testScheduler(t, database)

	now := time.Now()
	insertTx(t, database, "tx-edit", "expense", now, now.AddDate(0, 0, -10), 1000, "JPY")
	if err := s.RunJob(); err != nil {
		t.Fatalf("run job: %v", err)
	}

	day := today()
	path := filepath.Join(s.root, day[:4], day[5:7], day[8:10]+".xlsx")
	n1 := countRowsInSheet(t, path, "Transactions")
	if n1 != 1 {
		t.Fatalf("expected 1 transaction, got %d", n1)
	}

	// Ubah jumlah dan pindahkan tanggal transaksi ke hari lain.
	tomorrow := now.AddDate(0, 0, 1)
	_, err := database.Exec(`UPDATE transactions SET amount = 5000, date_time = ? WHERE id = ?`, tomorrow.Format(time.RFC3339), "tx-edit")
	if err != nil {
		t.Fatalf("update tx: %v", err)
	}
	NeedRecheck(day, tomorrow.Format("2006-01-02"))

	if err := s.RunJob(); err != nil {
		t.Fatalf("run job after edit: %v", err)
	}

	if fileExists(path) {
		t.Fatalf("old date file should be removed after transaction moved away: %s", path)
	}
	tomorrowPath := filepath.Join(s.root, tomorrow.Format("2006-01-02")[:4], tomorrow.Format("2006-01-02")[5:7], tomorrow.Format("2006-01-02")[8:10]+".xlsx")
	if !fileExists(tomorrowPath) {
		t.Fatalf("new date file should exist: %s", tomorrowPath)
	}
	if n := countRowsInSheet(t, tomorrowPath, "Transactions"); n != 1 {
		t.Fatalf("expected 1 transaction in new date file, got %d", n)
	}
}

func TestDeleteReconciliation(t *testing.T) {
	database := testDB(t)
	s := testScheduler(t, database)

	now := time.Now()
	insertTx(t, database, "tx-del", "expense", now, now.AddDate(0, 0, -10), 1000, "JPY")
	if err := s.RunJob(); err != nil {
		t.Fatalf("run job: %v", err)
	}

	day := today()
	path := filepath.Join(s.root, day[:4], day[5:7], day[8:10]+".xlsx")

	_, err := database.Exec(`UPDATE transactions SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?`, "tx-del")
	if err != nil {
		t.Fatalf("delete tx: %v", err)
	}
	NeedRecheck(day)

	if err := s.RunJob(); err != nil {
		t.Fatalf("run job after delete: %v", err)
	}
	if fileExists(path) {
		t.Fatalf("empty daily file should be removed after deletion: %s", path)
	}
}

func TestMonthlyAggregate(t *testing.T) {
	database := testDB(t)
	s := testScheduler(t, database)

	now := time.Now()
	d1 := now.AddDate(0, 0, -3)
	d2 := now.AddDate(0, 0, -2)
	insertTx(t, database, "tx-m1", "expense", d1, d1.AddDate(0, 0, -10), 1000, "JPY")
	insertTx(t, database, "tx-m2", "expense", d2, d2.AddDate(0, 0, -10), 2000, "JPY")
	if err := s.RunJob(); err != nil {
		t.Fatalf("run job: %v", err)
	}

	month := now.Format("2006-01")
	fullPath := filepath.Join(s.root, month[:4], month[5:7], "full_"+month[5:7]+".xlsx")
	if !fileExists(fullPath) {
		t.Fatalf("expected monthly aggregate %s", fullPath)
	}
	n := countRowsInSheet(t, fullPath, "Transactions")
	if n != 2 {
		t.Fatalf("expected 2 transactions in monthly aggregate, got %d", n)
	}
}

func TestGitCommitOnlyOnChange(t *testing.T) {
	database := testDB(t)
	s := testScheduler(t, database)

	now := time.Now()
	insertTx(t, database, "tx-git", "expense", now, now.AddDate(0, 0, -10), 1000, "JPY")
	if err := s.RunJob(); err != nil {
		t.Fatalf("run job: %v", err)
	}

	first := gitCommitCount(t, s.root)
	if first != 1 {
		t.Fatalf("expected exactly 1 commit after first run, got %d", first)
	}

	// Data tidak berubah: run kedua tidak boleh menambah commit.
	if err := s.RunJob(); err != nil {
		t.Fatalf("run job second: %v", err)
	}
	if c := gitCommitCount(t, s.root); c != first {
		t.Fatalf("commit count changed (%d -> %d) despite no data change", first, c)
	}

	// Data baru: commit baru.
	insertTx(t, database, "tx-git2", "income", now, now.AddDate(0, 0, -10), 500, "JPY")
	if err := s.RunJob(); err != nil {
		t.Fatalf("run job third: %v", err)
	}
	if c := gitCommitCount(t, s.root); c != first+1 {
		t.Fatalf("expected commit count %d, got %d", first+1, c)
	}
}

func TestCatchUpAfterGap(t *testing.T) {
	database := testDB(t)
	s := testScheduler(t, database)

	// Simulasikan backup pernah dijalankan 20 hari lalu (state tertinggal).
	_, err := database.Exec(`UPDATE backup_state SET last_created_processed = ? WHERE id = 'singleton'`, time.Now().AddDate(0, 0, -20).Format("2006-01-02"))
	if err != nil {
		// baris belum ada; buat dulu
		_, err = database.Exec(`INSERT INTO backup_state (id, last_created_processed) VALUES ('singleton', ?)`, time.Now().AddDate(0, 0, -20).Format("2006-01-02"))
		if err != nil {
			t.Fatalf("seed state: %v", err)
		}
	}

	now := time.Now()
	insertTx(t, database, "tx-catchup", "expense", now.AddDate(0, 0, -15), now.AddDate(0, 0, -15), 700, "JPY")
	if err := s.RunJob(); err != nil {
		t.Fatalf("run job catch-up: %v", err)
	}

	day := now.AddDate(0, 0, -15).Format("2006-01-02")
	path := filepath.Join(s.root, day[:4], day[5:7], day[8:10]+".xlsx")
	if !fileExists(path) {
		t.Fatalf("expected catch-up backup file %s", path)
	}
}

func TestRestoreXLSXRoundTrip(t *testing.T) {
	source := testDB(t)
	s := testScheduler(t, source)

	now := time.Now()
	insertTx(t, source, "tx-restore", "expense", now, now.AddDate(0, 0, -10), 1234, "IDR")
	if err := s.RunJob(); err != nil {
		t.Fatalf("run job: %v", err)
	}

	day := today()
	path := filepath.Join(s.root, day[:4], day[5:7], day[8:10]+".xlsx")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read backup: %v", err)
	}

	target := testDB(t)
	if err := restoreXLSX(target, data); err != nil {
		t.Fatalf("restore xlsx: %v", err)
	}

	var id string
	err = target.QueryRow(`SELECT id FROM transactions WHERE id = ?`, "tx-restore").Scan(&id)
	if err != nil {
		t.Fatalf("restored transaction missing: %v", err)
	}
	if id != "tx-restore" {
		t.Fatalf("restored id mismatch: %s", id)
	}
}

func runGit(root string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	return cmd.CombinedOutput()
}

func gitCommitCount(t *testing.T, root string) int {
	t.Helper()
	out, err := runGit(root, "rev-list", "--count", "HEAD")
	if err != nil {
		return 0
	}
	var n int
	_, _ = fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &n)
	return n
}
