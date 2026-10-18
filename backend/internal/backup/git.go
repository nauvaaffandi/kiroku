package backup

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var secretPattern = regexp.MustCompile(`github_pat_[A-Za-z0-9_]+|ghp_[A-Za-z0-9]+|gho_[A-Za-z0-9]+|ghs_[A-Za-z0-9]+|ghr_[A-Za-z0-9]+|xox[baprs]-[A-Za-z0-9-]+`)

func scrub(s string) string {
	return secretPattern.ReplaceAllString(s, "[REDACTED]")
}

type Git struct {
	root  string
	token string
}

func NewGit(root, token string) *Git {
	return &Git{root: root, token: token}
}

func (g *Git) run(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = g.root
	out, err := cmd.CombinedOutput()
	s := scrub(string(out))
	if err != nil {
		return s, fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return s, nil
}

func (g *Git) EnsureRepo() error {
	if err := os.MkdirAll(g.root, 0o755); err != nil {
		return err
	}
	if fi, err := os.Stat(filepath.Join(g.root, ".git")); err == nil && fi.IsDir() {
		return nil
	}
	ctx := context.Background()
	if _, err := g.run(ctx, "init", "-b", "main"); err != nil {
		return err
	}
	_, _ = g.run(ctx, "config", "user.email", "kiroku-backup@localhost")
	_, _ = g.run(ctx, "config", "user.name", "Kiroku Backup")
	return nil
}

func (g *Git) HasOrigin() bool {
	ctx := context.Background()
	_, err := g.run(ctx, "remote", "get-url", "origin")
	return err == nil
}

func (g *Git) StageAndCommit(msg string) (committed bool, err error) {
	ctx := context.Background()
	_, _ = g.run(ctx, "add", "--", "*.xlsx")
	staged := false
	if out, err := g.run(ctx, "status", "--porcelain", "--", "*.xlsx"); err == nil {
		staged = strings.TrimSpace(out) != ""
	}
	if !staged {
		return false, nil
	}
	if _, err := g.run(ctx, "commit", "-m", msg); err != nil {
		return false, err
	}
	return true, nil
}

func (g *Git) Push() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := g.run(ctx, "push", "origin", "HEAD")
	return err
}
