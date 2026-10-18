package backup

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func restoreTarGz(data []byte) error {
	gzr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("invalid backup format")
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
			return fmt.Errorf("read archive: %w", err)
		}

		switch {
		case header.Name == "manifest.json":
			hasManifest = true
		case header.Name == "kiroku.db":
			if err := writeAll("data/kiroku.db", tr); err != nil {
				return err
			}
		case strings.HasPrefix(header.Name, "attachments/"):
			if err := writeAll(filepath.Join("data", header.Name), tr); err != nil {
				return err
			}
		}
	}

	if !hasManifest {
		return fmt.Errorf("missing manifest")
	}
	return nil
}

func writeAll(path string, r io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if _, err := io.Copy(tmp, r); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}
