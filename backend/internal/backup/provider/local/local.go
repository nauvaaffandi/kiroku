package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"kiroku/backend/internal/backup/provider"
)

type LocalProvider struct {
	root string
}

func New(root string) *LocalProvider {
	return &LocalProvider{root: root}
}

func (lp *LocalProvider) Name() string {
	return "local"
}

func (lp *LocalProvider) Upload(ctx context.Context, key string, r io.Reader, meta provider.Metadata) error {
	if err := os.MkdirAll(lp.root, 0o755); err != nil {
		return err
	}
	destPath := filepath.Join(lp.root, key)
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return err
	}
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, r); err != nil {
		return err
	}
	return nil
}

func (lp *LocalProvider) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	destPath := filepath.Join(lp.root, key)
	f, err := os.Open(destPath)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (lp *LocalProvider) List(ctx context.Context) ([]provider.Metadata, error) {
	var metas []provider.Metadata
	err := filepath.Walk(lp.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(lp.root, path)
		if err != nil {
			return err
		}
		metas = append(metas, provider.Metadata{
			Key:          rel,
			Size:         info.Size(),
			LastModified: info.ModTime(),
		})
		return nil
	})
	return metas, err
}

func (lp *LocalProvider) Delete(ctx context.Context, key string) error {
	destPath := filepath.Join(lp.root, key)
	return os.Remove(destPath)
}

func (lp *LocalProvider) HealthCheck(ctx context.Context) error {
	if err := os.MkdirAll(lp.root, 0o755); err != nil {
		return fmt.Errorf("local root inaccessible: %w", err)
	}
	return nil
}
