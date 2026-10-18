package backup

import (
	"database/sql"
	"time"

	"kiroku/backend/internal/db"
)

type backupState struct {
	LastCreatedProcessed string
	LastRunAt            time.Time
	LastStatus           string
	LastError            string
	LastPushedAt         *time.Time
	LastCommit           string
	FormatVersion        int
}

func addRecheckDates(database *db.DB, dates []string) error {
	for _, d := range dates {
		if d == "" {
			continue
		}
		if _, err := database.Exec(`INSERT OR IGNORE INTO backup_rechecks (date) VALUES (?)`, d); err != nil {
			return err
		}
	}
	return nil
}

func AddRecheckDates(database *db.DB, dates []string) error {
	return addRecheckDates(database, dates)
}

var recheckHook func([]string)

func RegisterRecheckHook(fn func([]string)) {
	recheckHook = fn
}

func NeedRecheck(dates ...string) {
	if recheckHook != nil && len(dates) > 0 {
		recheckHook(dates)
	}
}

func (s *Scheduler) getState() (backupState, error) {
	var st backupState
	var lastCreated, lastRun, lastStatus, lastError, lastPushed, lastCommit sql.NullString
	var formatVersion sql.NullInt64
	err := s.db.QueryRow(`
		SELECT last_created_processed, last_run_at, last_status, last_error, last_pushed_at, last_commit, format_version
		FROM backup_state WHERE id = 'singleton'
	`).Scan(&lastCreated, &lastRun, &lastStatus, &lastError, &lastPushed, &lastCommit, &formatVersion)
	if err == sql.ErrNoRows {
		return backupState{FormatVersion: 1}, nil
	}
	if err != nil {
		return st, err
	}
	st.LastCreatedProcessed = lastCreated.String
	st.LastStatus = lastStatus.String
	st.LastError = lastError.String
	st.LastCommit = lastCommit.String
	st.FormatVersion = int(formatVersion.Int64)
	if lastRun.Valid && lastRun.String != "" {
		st.LastRunAt = parseTimeValue(lastRun.String)
	}
	if lastPushed.Valid && lastPushed.String != "" {
		t := parseTimeValue(lastPushed.String)
		st.LastPushedAt = &t
	}
	return st, nil
}

func (s *Scheduler) saveState(started time.Time, status, errMsg string) {
	var lastCreated any
	if s.lastCreatedProcessed != "" {
		lastCreated = s.lastCreatedProcessed
	}
	_, _ = s.db.Exec(`
		INSERT INTO backup_state (id, last_created_processed, last_run_at, last_status, last_error, last_pushed_at, last_commit, format_version)
		VALUES ('singleton', ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			last_created_processed = excluded.last_created_processed,
			last_run_at = excluded.last_run_at,
			last_status = excluded.last_status,
			last_error = excluded.last_error,
			last_pushed_at = excluded.last_pushed_at,
			last_commit = excluded.last_commit,
			format_version = excluded.format_version
	`, lastCreated, started.Format(time.RFC3339), status, errMsg, nil, s.lastCommit, formatVersion)
}

func parseTimeValue(s string) time.Time {
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05.999999999-07:00", "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

func (s *Scheduler) loadDates() (map[string]string, error) {
	rows, err := s.db.Query(`SELECT date, fingerprint FROM backup_dates`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var d, fp string
		if err := rows.Scan(&d, &fp); err != nil {
			return nil, err
		}
		out[d] = fp
	}
	return out, rows.Err()
}

func (s *Scheduler) upsertDate(date, fingerprint string, hasData bool) error {
	_, err := s.db.Exec(`
		INSERT INTO backup_dates (date, fingerprint, has_data, status, checked_at)
		VALUES (?, ?, ?, 'processed', CURRENT_TIMESTAMP)
		ON CONFLICT(date) DO UPDATE SET fingerprint = excluded.fingerprint, has_data = excluded.has_data, status = 'processed', checked_at = CURRENT_TIMESTAMP
	`, date, fingerprint, boolInt(hasData))
	return err
}

func (s *Scheduler) deleteDate(date string) error {
	_, err := s.db.Exec(`DELETE FROM backup_dates WHERE date = ?`, date)
	return err
}

func (s *Scheduler) loadMonths() (map[string]string, error) {
	rows, err := s.db.Query(`SELECT month, fingerprint FROM backup_months`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var m, fp string
		if err := rows.Scan(&m, &fp); err != nil {
			return nil, err
		}
		out[m] = fp
	}
	return out, rows.Err()
}

func (s *Scheduler) upsertMonth(month, fingerprint string) error {
	_, err := s.db.Exec(`
		INSERT INTO backup_months (month, fingerprint, status, updated_at)
		VALUES (?, ?, 'processed', CURRENT_TIMESTAMP)
		ON CONFLICT(month) DO UPDATE SET fingerprint = excluded.fingerprint, status = 'processed', updated_at = CURRENT_TIMESTAMP
	`, month, fingerprint)
	return err
}

func (s *Scheduler) deleteMonth(month string) error {
	_, err := s.db.Exec(`DELETE FROM backup_months WHERE month = ?`, month)
	return err
}

func (s *Scheduler) takeRechecks() ([]string, error) {
	rows, err := s.db.Query(`SELECT date FROM backup_rechecks ORDER BY date`)
	if err != nil {
		return nil, err
	}
	var dates []string
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			rows.Close()
			return nil, err
		}
		dates = append(dates, d)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if _, err := s.db.Exec(`DELETE FROM backup_rechecks`); err != nil {
		return nil, err
	}
	return dates, nil
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
