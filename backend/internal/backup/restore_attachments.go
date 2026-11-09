package backup

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func parseAttachmentsSheet(rows [][]string) ([]attachment, error) {
	if len(rows) < 2 {
		return nil, nil
	}
	idx := headerIndex(rows[0])
	grouped := map[string]*attachment{}
	var order []string
	for _, r := range rows[1:] {
		txID := cellAt(r, idx, "transaction_id")
		path := cellAt(r, idx, "path")
		if txID == "" || path == "" {
			continue
		}
		key := txID + "\x00" + path
		a, ok := grouped[key]
		if !ok {
			size, _ := strconv.ParseInt(cellAt(r, idx, "size"), 10, 64)
			a = &attachment{
				TransactionID: txID,
				Path:          path,
				ContentType:   cellAt(r, idx, "content_type"),
				Size:          size,
				SHA256:        cellAt(r, idx, "sha256"),
			}
			grouped[key] = a
			order = append(order, key)
		}
		chunkIdx, _ := strconv.Atoi(cellAt(r, idx, "chunk_index"))
		for len(a.chunks) <= chunkIdx {
			a.chunks = append(a.chunks, "")
		}
		a.chunks[chunkIdx] = cellAt(r, idx, "data")
	}
	var out []attachment
	for _, key := range order {
		a := grouped[key]
		encoded := strings.Join(a.chunks, "")
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return nil, fmt.Errorf("attachment %s: invalid base64: %w", a.Path, err)
		}
		if a.SHA256 != "" {
			sum := sha256.Sum256(decoded)
			if !strings.EqualFold(a.SHA256, hex.EncodeToString(sum[:])) {
				return nil, fmt.Errorf("attachment %s: checksum mismatch", a.Path)
			}
		}
		a.Data = decoded
		out = append(out, *a)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out, nil
}

func writeAttachments(atts []attachment) error {
	for _, a := range atts {
		rel := filepath.Clean(a.Path)
		if filepath.IsAbs(rel) || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return fmt.Errorf("attachment %s: unsafe path", a.Path)
		}
		full := filepath.Join("data", filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(full, a.Data, 0o644); err != nil {
			return err
		}
	}
	return nil
}
