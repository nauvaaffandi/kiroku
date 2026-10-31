package backup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/backup/provider"
	"kiroku/backend/internal/backup/provider/gdrive"
	"kiroku/backend/internal/backup/provider/git"
	"kiroku/backend/internal/backup/provider/local"
	"kiroku/backend/internal/config"
	"kiroku/backend/internal/db"
)

type Scheduler struct {
	db        *db.DB
	cfg       *config.Config
	root      string
	providers []provider.StorageProvider

	mu                   sync.Mutex
	lastCreatedProcessed string
	lastCommit           string
}

func NewScheduler(database *db.DB, cfg *config.Config, root, token string) *Scheduler {
	var providers []provider.StorageProvider
	
	providerNames := provider.ParseProviders(cfg.BackupProviders)
	if len(providerNames) == 0 && (token != "" || cfg.BackupRepoPath != "") {
		providerNames = []string{"git"}
	}

	for _, name := range providerNames {
		switch name {
		case "git":
			providers = append(providers, git.New(root, token))
		case "gdrive":
			providers = append(providers, gdrive.New("", "", "", ""))
		}
	}

	return &Scheduler{
		db:        database,
		cfg:       cfg,
		root:      root,
		providers: providers,
	}
}

func (s *Scheduler) RunJob() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	started := time.Now()
	defer func() {
		status := "ok"
		msg := ""
		if err != nil {
			status = "failed"
			msg = err.Error()
		}
		s.saveState(started, status, msg)
	}()

	// 1. Mandatory Local-First Backup Execution
	localProvider := local.New(s.root)
	if err := localProvider.HealthCheck(context.Background()); err != nil {
		return fmt.Errorf("local backup storage initialization failed: %w", err)
	}

	cutoff := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	s.lastCreatedProcessed = cutoff

	groups, err := s.loadEligible(cutoff)
	if err != nil {
		return fmt.Errorf("load eligible transactions: %w", err)
	}
	accounts, err := s.loadAccounts()
	if err != nil {
		return fmt.Errorf("load accounts: %w", err)
	}
	categories, err := s.loadCategories()
	if err != nil {
		return fmt.Errorf("load categories: %w", err)
	}
	rechecks, err := s.takeRechecks()
	if err != nil {
		return fmt.Errorf("load recheck dates: %w", err)
	}
	existing, err := s.loadDates()
	if err != nil {
		return fmt.Errorf("load backup dates: %w", err)
	}
	months, err := s.loadMonths()
	if err != nil {
		return fmt.Errorf("load backup months: %w", err)
	}

	allDates := map[string]bool{}
	for d := range groups {
		allDates[d] = true
	}
	for d := range existing {
		allDates[d] = true
	}
	for _, d := range rechecks {
		allDates[d] = true
	}

	var changedKeys []string

	changedMonths := map[string]bool{}
	for date := range allDates {
		txs := groups[date]
		if len(txs) == 0 {
			if err := s.removeDaily(date); err != nil {
				return fmt.Errorf("remove daily %s: %w", date, err)
			}
			if err := s.deleteDate(date); err != nil {
				return err
			}
			changedMonths[monthOf(date)] = true
			continue
		}
		atts := collectAttachments(s.cfg, txs)
		fp := fingerprint(txs, accounts, categories, atts)
		if existing[date] == fp {
			continue
		}
		if err := s.writeDaily(date, txs, accounts, categories, atts); err != nil {
			return fmt.Errorf("write daily %s: %w", date, err)
		}
		if err := s.upsertDate(date, fp, true); err != nil {
			return err
		}
		changedMonths[monthOf(date)] = true
		changedKeys = append(changedKeys, fmt.Sprintf("%s/%s/%s.xlsx", date[:4], date[5:7], date[8:10]))
	}

	for month := range changedMonths {
		var monthTxs [][]string
		for date, txs := range groups {
			if strings.HasPrefix(date, month) {
				monthTxs = append(monthTxs, txs...)
			}
		}
		if len(monthTxs) == 0 {
			if err := s.removeMonthly(month); err != nil {
				return fmt.Errorf("remove monthly %s: %w", month, err)
			}
			if err := s.deleteMonth(month); err != nil {
				return err
			}
			continue
		}
		atts := collectAttachments(s.cfg, monthTxs)
		fp := fingerprint(monthTxs, accounts, categories, atts)
		if months[month] == fp {
			continue
		}
		if err := s.writeMonthly(month, monthTxs, accounts, categories, atts); err != nil {
			return fmt.Errorf("write monthly %s: %w", month, err)
		}
		if err := s.upsertMonth(month, fp); err != nil {
			return err
		}
		changedKeys = append(changedKeys, fmt.Sprintf("%s/%s/full_%s.xlsx", month[:4], month[5:7], month[5:7]))
	}

	// 2. Dispatch to Configured Pluggable Remote Providers with Failure Isolation
	ctx := context.Background()
	for _, prov := range s.providers {
		if err := prov.HealthCheck(ctx); err != nil {
			log.Warn().Str("provider", prov.Name()).Err(err).Msg("backup: provider health check failed")
			s.recordSyncState(started, prov.Name(), "failed", err.Error())
			continue
		}

		providerFailed := false
		for _, key := range changedKeys {
			fullPath := filepath.Join(s.root, key)
			f, err := os.Open(fullPath)
			if err != nil {
				continue
			}
			fi, _ := f.Stat()
			meta := provider.Metadata{
				Key:          key,
				Size:         fi.Size(),
				LastModified: fi.ModTime(),
			}
			if err := prov.Upload(ctx, key, f, meta); err != nil {
				f.Close()
				log.Warn().Str("provider", prov.Name()).Str("key", key).Err(err).Msg("backup: provider upload failed")
				providerFailed = true
				break
			}
			f.Close()
		}

		if providerFailed {
			s.recordSyncState(started, prov.Name(), "failed", "one or more artifact uploads failed")
		} else {
			s.recordSyncState(started, prov.Name(), "success", "")
		}
	}

	// 3. Git provider: batched commit-and-push after all provider uploads
	for _, prov := range s.providers {
		if gitProv, ok := prov.(*git.GitProvider); ok {
			if committed, err := gitProv.CommitAndPush(ctx, "backup: reconcile through "+cutoff); err != nil {
				log.Warn().Err(err).Msg("backup: git commit failed, will retry on next run")
				s.recordSyncState(started, prov.Name(), "failed", err.Error())
			} else if committed {
				s.recordSyncState(started, prov.Name(), "success", "")
			}
		}
	}

	s.lastCommit = time.Now().Format(time.RFC3339)
	return nil
}

func (s *Scheduler) recordSyncState(execTime time.Time, providerName, status, errMsg string) {
	_, _ = s.db.Exec(
		`INSERT INTO backup_provider_sync_states (execution_id, provider_name, status, error_message, created_at) VALUES (?, ?, ?, ?, ?)`,
		execTime.Format(time.RFC3339), providerName, status, errMsg, time.Now(),
	)
}

func (s *Scheduler) loadEligible(cutoff string) (map[string][][]string, error) {
	rows, err := s.db.Query(
		`SELECT `+strings.Join(txColumns, ", ")+` FROM transactions WHERE deleted_at IS NULL AND substr(created_at,1,10) <= ? ORDER BY date_time, id`,
		cutoff,
	)
	if err != nil {
		return nil, err
	}
	txs, err := scanTransactions(rows)
	if err != nil {
		return nil, err
	}
	groups := map[string][][]string{}
	for _, t := range txs {
		if len(t) < 3 || len(t[2]) < 10 {
			continue
		}
		d := t[2][:10]
		groups[d] = append(groups[d], t)
	}
	return groups, nil
}

func (s *Scheduler) loadAccounts() ([][]string, error) {
	rows, err := s.db.Query(`SELECT id, name, currency, type, icon, color, archived, created_at, updated_at FROM accounts ORDER BY id`)
	if err != nil {
		return nil, err
	}
	return scanAccounts(rows)
}

func (s *Scheduler) loadCategories() ([][]string, error) {
	rows, err := s.db.Query(`SELECT id, parent_id, name, icon, color, kind, sort_order, created_at, updated_at FROM categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	return scanCategories(rows)
}

func (s *Scheduler) writeDaily(date string, txs, accounts, categories [][]string, atts []attachment) error {
	f, err := buildWorkbook("daily", date, date, txs, accounts, categories, atts)
	if err != nil {
		return err
	}
	return writeXLSX(dailyPath(s.root, date), f)
}

func (s *Scheduler) writeMonthly(month string, txs, accounts, categories [][]string, atts []attachment) error {
	f, err := buildWorkbook("monthly", month+"-01", month+"-31", txs, accounts, categories, atts)
	if err != nil {
		return err
	}
	return writeXLSX(monthlyPath(s.root, month), f)
}

func (s *Scheduler) removeDaily(date string) error {
	path := dailyPath(s.root, date)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	_ = os.Remove(filepath.Dir(path))
	return nil
}

func (s *Scheduler) removeMonthly(month string) error {
	path := monthlyPath(s.root, month)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func dailyPath(root, date string) string {
	return filepath.Join(root, date[:4], date[5:7], date[8:10]+".xlsx")
}

func monthlyPath(root, month string) string {
	return filepath.Join(root, month[:4], month[5:7], "full_"+month[5:7]+".xlsx")
}

func monthOf(date string) string {
	if len(date) >= 7 {
		return date[:7]
	}
	return date
}

func StartScheduler(s *Scheduler) {
	go func() {
		if err := s.RunJob(); err != nil {
			log.Warn().Err(err).Msg("backup: initial run failed")
		}
		for {
			next, err := nextCron(s.cfg.BackupSchedule, time.Now())
			if err != nil {
				log.Warn().Err(err).Msg("backup: invalid schedule, retrying in 1 hour")
				time.Sleep(time.Hour)
				continue
			}
			d := time.Until(next)
			if d <= 0 {
				d = time.Minute
			}
			time.Sleep(d)
			if err := s.RunJob(); err != nil {
				log.Warn().Err(err).Msg("backup: scheduled run failed")
			}
		}
	}()
}
