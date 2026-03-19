package gitops

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// mu serializes git operations so concurrent worker goroutines don't
// stomp on each other.
var mu sync.Mutex

// EnsureRepo ensures outputDir is a git repository pointed at remote.
// If the directory is not yet a git repo, it is initialized and the
// remote is added. If it already is a repo, the remote URL is updated.
// If remote is empty, this is a no-op.
func EnsureRepo(dir, remote string) {
	if remote == "" {
		return
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		slog.Warn("gitops: could not create output dir", "dir", dir, "err", err)
		return
	}

	if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
		// Already a repo — just make sure the remote points where we want.
		if out, err := git(dir, "remote", "set-url", "origin", remote); err != nil {
			slog.Warn("gitops: could not set remote url", "err", err, "output", string(out))
		}
		return
	}

	// Fresh init.
	if out, err := git(dir, "init"); err != nil {
		slog.Warn("gitops: git init failed", "err", err, "output", string(out))
		return
	}
	if out, err := git(dir, "remote", "add", "origin", remote); err != nil {
		slog.Warn("gitops: could not add remote", "err", err, "output", string(out))
	}

	slog.Info("gitops: initialized output repo", "dir", dir, "remote", remote)
}

// CommitAndPush stages all changes under dir, commits them with msg, and
// pushes to origin. It is a no-op if there is nothing to commit. Errors
// are logged as warnings rather than returned — git is best-effort.
func CommitAndPush(dir, msg string) {
	mu.Lock()
	defer mu.Unlock()

	out, err := git(dir, "status", "--porcelain", ".")
	if err != nil {
		slog.Warn("gitops: git status failed", "err", err, "output", string(out))
		return
	}
	if strings.TrimSpace(string(out)) == "" {
		return // nothing to commit
	}

	if out, err := git(dir, "add", "."); err != nil {
		slog.Warn("gitops: git add failed", "err", err, "output", string(out))
		return
	}

	fullMsg := fmt.Sprintf("barista: %s (%s)", msg, time.Now().UTC().Format(time.RFC3339))
	if out, err := git(dir, "commit", "-m", fullMsg); err != nil {
		slog.Warn("gitops: git commit failed", "err", err, "output", string(out))
		return
	}

	// --set-upstream handles both the first push and subsequent ones.
	if out, err := git(dir, "push", "--set-upstream", "origin", "HEAD"); err != nil {
		slog.Warn("gitops: git push failed", "err", err, "output", string(out))
		return
	}

	slog.Info("gitops: committed and pushed", "msg", fullMsg)
}

func git(dir string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	return cmd.CombinedOutput()
}
