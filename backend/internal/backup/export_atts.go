package backup

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"kiroku/backend/internal/config"
)

func collectAttachments(cfg *config.Config, txs [][]string) []attachment {
	if cfg == nil {
		return nil
	}
	seen := map[string]bool{}
	var atts []attachment
	for _, tx := range txs {
		id := tx[0]
		raw := tx[17]
		if raw == "" {
			continue
		}
		for _, p := range parseAttachmentPaths(raw) {
			if seen[p] {
				continue
			}
			seen[p] = true
			full := resolveAttachmentPath(p)
			data, err := os.ReadFile(full)
			if err != nil {
				continue
			}
			sum := sha256.Sum256(data)
			atts = append(atts, attachment{
				TransactionID: id,
				Path:          p,
				ContentType:   contentTypeFor(p),
				Size:          int64(len(data)),
				SHA256:        hex.EncodeToString(sum[:]),
				Data:          data,
			})
		}
	}
	return atts
}

func parseAttachmentPaths(raw string) []string {
	var arr []string
	if err := json.Unmarshal([]byte(raw), &arr); err == nil {
		return arr
	}
	if strings.TrimSpace(raw) != "" {
		return []string{strings.TrimSpace(raw)}
	}
	return nil
}

func resolveAttachmentPath(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	if strings.HasPrefix(p, "attachments"+string(filepath.Separator)) || p == "attachments" {
		return filepath.Join(config.DataDir(), filepath.FromSlash(p))
	}
	return filepath.Join(config.DataDir(), "attachments", filepath.Base(p))
}

func contentTypeFor(p string) string {
	switch strings.ToLower(filepath.Ext(p)) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}
