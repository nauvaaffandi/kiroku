package backup

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"kiroku/backend/internal/auth"
)

func (h *Handler) Restore(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	file, _, err := r.FormFile("file")
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

	if h.cfg.BackupPassphrase != "" {
		data, err = decrypt(data, h.cfg.BackupPassphrase)
		if err != nil {
			http.Error(w, "decrypt error (check passphrase)", http.StatusBadRequest)
			return
		}
	}

	// Extract tar.gz from data
	gzr, err := gzip.NewReader(&readSeeker{data: data})
	if err != nil {
		http.Error(w, "invalid backup format", http.StatusBadRequest)
		return
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	hasManifest := false

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		if header.Name == "manifest.json" {
			hasManifest = true
		} else if header.Name == "kiroku.db" {
			dbOut, err := os.Create("data/kiroku.db")
			if err == nil {
				io.Copy(dbOut, tr)
				dbOut.Close()
			}
		} else if strings.HasPrefix(header.Name, "attachments/") {
			os.MkdirAll(filepath.Dir(filepath.Join("data", header.Name)), 0755)
			attOut, err := os.Create(filepath.Join("data", header.Name))
			if err == nil {
				io.Copy(attOut, tr)
				attOut.Close()
			}
		}
	}

	if !hasManifest {
		http.Error(w, "invalid backup: missing manifest", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "restored successfully. please restart server."})
}

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
