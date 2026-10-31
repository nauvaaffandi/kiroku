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
	"time"

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
	cmd := execCommandContext(ctx, "git", args...)
	cmd.Dir = gp.root
	out, err := cmd.CombinedOutput()
	s := scrub(string(out))
	if err != nil {
		return s, fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return s, nil
}

// helper to allow mocking or standard exec
var execCommandContext = execCommand

func execCommand(ctx context.Context, name string, arg ...string) *exec.Cmd {
	return exec.CommandContext(ctx, name, arg...)
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

	// Read content fully into memory to avoid read/write race condition on same file
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if err := os.WriteFile(destPath, data, 0o644); err != nil {
		return err
	}
	return nil
}

func (gp *GitProvider) CommitAndPush(ctx context.Context, msg string) (bool, error) {
	if err := gp.ensureRepo(); err != nil {
		return false, err
	}
	_, _ = gp.run(ctx, "add", "--", "*.xlsx")
	staged := false
	if out, err := gp.run(ctx, "status", "--porcelain", "--", "*.xlsx"); err == nil {
		staged = strings.TrimSpace(out) != ""
	}
	if !staged {
		return false, nil
	}
	if _, err := gp.run(ctx, "commit", "-m", msg); err != nil {
		return false, err
	}
	
	// Push if origin exists
	if gp.HasOrigin() {
		_ = gp.Push()
	}
	return true, nil
}

func (gp *GitProvider) HasOrigin() bool {
	ctx := context.Background()
	_, err := gp.run(ctx, "remote", "get-url", "origin")
	return err == nil
}

func (gp *GitProvider) Push() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := gp.run(ctx, "push", "origin", "HEAD")
	return err
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
	_, _ = gp.run(context.Background(), "add", "--", key)
	_, _ = gp.run(context.Background(), "commit", "-m", "backup: remove "+key)
	if gp.HasOrigin() {
		_ = gp.Push()
	}
	return nil
}

func (gp *GitProvider) HealthCheck(ctx context.Context) error {
	if err := gp.ensureRepo(); err != nil {
		return fmt.Errorf("git repo inaccessible: %w", err)
	}
	return nil
}
