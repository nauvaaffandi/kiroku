package git

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"kiroku/backend/internal/backup/provider"
)

var secretPattern = regexp.MustCompile(`github_pat_[A-Za-z0-9_]+|ghp_[A-Za-z0-9]+|gho_[A-Za-z0-9]+|ghs_[A-Za-z0-9]+|ghr_[A-Za-z0-9]+|xox[baprs]-[A-Za-z0-9-]+`)

func scrub(s string) string {
	return secretPattern.ReplaceAllString(s, "[REDACTED]")
}

type GitProvider struct {
	root  string
	token string
}

func New(root, token string) *GitProvider {
	return &GitProvider{root: root, token: token}
}

func (gp *GitProvider) Name() string {
	return "git"
}

func (gp *GitProvider) run(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = gp.root
	out, err := cmd.CombinedOutput()
	s := scrub(string(out))
	if err != nil {
		return s, fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return s, nil
}

func (gp *GitProvider) ensureRepo() error {
	if err := os.MkdirAll(gp.root, 0o755); err != nil {
		return err
	}
	if fi, err := os.Stat(filepath.Join(gp.root, ".git")); err == nil && fi.IsDir() {
		return nil
	}
	ctx := context.Background()
	if _, err := gp.run(ctx, "init", "-b", "main"); err != nil {
		return err
	}
	_, _ = gp.run(ctx, "config", "user.email", "kiroku-backup@localhost")
	_, _ = gp.run(ctx, "config", "user.name", "Kiroku Backup")
	return nil
}

func (gp *GitProvider) Upload(ctx context.Context, key string, r io.Reader, meta provider.Metadata) error {
	if err := gp.ensureRepo(); err != nil {
		return err
	}
	destPath := filepath.Join(gp.root, key)
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

	// Stage, commit, and push
	_, _ = gp.run(ctx, "add", "--", key)
	staged := false
	if output, err := gp.run(ctx, "status", "--porcelain", "--", key); err == nil {
		staged = strings.TrimSpace(output) != ""
	}
	if staged {
		if _, err := gp.run(ctx, "commit", "-m", "backup: update "+key); err != nil {
			return err
		}
		// Try pushing if remote exists
		_, _ = gp.run(ctx, "push", "origin", "HEAD")
	}
	return nil
}

func (gp *GitProvider) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	destPath := filepath.Join(gp.root, key)
	f, err := os.Open(destPath)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (gp *GitProvider) List(ctx context.Context) ([]provider.Metadata, error) {
	var metas []provider.Metadata
	err := filepath.Walk(gp.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || strings.Contains(path, ".git") {
			return nil
		}
		rel, err := filepath.Rel(gp.root, path)
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

func (gp *GitProvider) Delete(ctx context.Context, key string) error {
	destPath := filepath.Join(gp.root, key)
	if err := os.Remove(destPath); err != nil {
		return err
	}
	_, _ = gp.run(ctx, "add", "--", key)
	_, _ = gp.run(ctx, "commit", "-m", "backup: remove "+key)
	_, _ = gp.run(ctx, "push", "origin", "HEAD")
	return nil
}

func (gp *GitProvider) HealthCheck(ctx context.Context) error {
	if err := gp.ensureRepo(); err != nil {
		return fmt.Errorf("git repo inaccessible: %w", err)
	}
	return nil
}
