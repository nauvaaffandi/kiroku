package backup

import (
	"database/sql"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/config"
	"kiroku/backend/internal/db"
	"kiroku/backend/internal/models"
)

type Handler struct {
	db    *db.DB
	cfg   *config.Config
	sched *Scheduler
}

func NewHandler(database *db.DB, cfg *config.Config) *Handler {
	return &Handler{
		db:    database,
		cfg:   cfg,
		sched: NewScheduler(database, cfg, cfg.BackupRepoPath, cfg.BackupGitToken),
	}
}

func (h *Handler) Scheduler() *Scheduler {
	return h.sched
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Trigger)
	r.Post("/create", h.Trigger)
	r.Get("/{id}/download", h.Download)
	r.Post("/restore", h.Restore)
	return r
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var backups []models.Backup

	rows, err := h.db.Query(`SELECT date, checked_at FROM backup_dates WHERE has_data = 1 ORDER BY date DESC`)
	if err != nil {
		log.Error().Err(err).Msg("backup: list daily")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var date string
		var checkedAt sql.NullTime
		if err := rows.Scan(&date, &checkedAt); err != nil {
			continue
		}
		b := models.Backup{
			ID:            date,
			Filename:      date + ".xlsx",
			RetentionTier: "daily",
		}
		if fi, err := os.Stat(dailyPath(h.sched.root, date)); err == nil {
			b.Size = fi.Size()
		}
		if checkedAt.Valid {
			b.CreatedAt = checkedAt.Time
		}
		backups = append(backups, b)
	}
	rows.Close()

	legacy, err := h.db.Query("SELECT id, filename, size, retention_tier, created_at FROM backups ORDER BY created_at DESC")
	if err != nil {
		log.Error().Err(err).Msg("backup: list legacy")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer legacy.Close()
	for legacy.Next() {
		var b models.Backup
		if err := legacy.Scan(&b.ID, &b.Filename, &b.Size, &b.RetentionTier, &b.CreatedAt); err != nil {
			continue
		}
		backups = append(backups, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(backups)
}

func (h *Handler) Trigger(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.sched.RunJob(); err != nil {
		log.Error().Err(err).Msg("backup: create")
		http.Error(w, "backup failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "backup completed"})
}

func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	var hasData int
	err := h.db.QueryRow(`SELECT has_data FROM backup_dates WHERE date = ?`, id).Scan(&hasData)
	if err == nil && hasData == 1 {
		w.Header().Set("Content-Disposition", `attachment; filename="`+id+`.xlsx"`)
		http.ServeFile(w, r, dailyPath(h.sched.root, id))
		return
	}

	var filename string
	err = h.db.QueryRow("SELECT filename FROM backups WHERE id = ?", id).Scan(&filename)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)

	http.ServeFile(w, r, filepath.Join("data/backups", filename))
}

func (h *Handler) Restore(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	file, header, err := formFile(r, "file", "backup")
	if err != nil {
		http.Error(w, "file required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "read error", http.StatusBadRequest)
		return
	}

	if isXLSX(data, header) {
		if err := restoreXLSX(h.db, data); err != nil {
			log.Error().Err(err).Msg("backup: restore xlsx")
			http.Error(w, "restore failed: "+err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "restored successfully"})
		return
	}

	if h.cfg.BackupPassphrase != "" {
		data, err = decrypt(data, h.cfg.BackupPassphrase)
		if err != nil {
			http.Error(w, "decrypt error (check passphrase)", http.StatusBadRequest)
			return
		}
	}

	if err := h.restoreLegacy(data); err != nil {
		log.Error().Err(err).Msg("backup: restore legacy")
		http.Error(w, "invalid backup: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "restored successfully. please restart server."})
}

func (h *Handler) restoreLegacy(data []byte) error {
	return restoreTarGz(data)
}

func isXLSX(data []byte, header *multipart.FileHeader) bool {
	if header != nil && (filepath.Ext(header.Filename) == ".xlsx" || filepath.Ext(header.Filename) == ".xls") {
		return true
	}
	return len(data) >= 4 && data[0] == 'P' && data[1] == 'K' && data[2] == 3 && data[3] == 4
}

func formFile(r *http.Request, names ...string) (multipart.File, *multipart.FileHeader, error) {
	var lastErr error
	for _, name := range names {
		f, h, err := r.FormFile(name)
		if err == nil {
			return f, h, nil
		}
		lastErr = err
	}
	return nil, nil, lastErr
}
