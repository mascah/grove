// Package repo locates the Git common directory and takes the advisory locks
// that Grove's mutating commands share across linked worktrees.
package repo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// CommonDir returns the repository's common Git directory for root, plus the
// path of root inside the worktree (git's --show-prefix), creating the
// Grove-owned coordination folder under it. Read-only commands use Locate.
func CommonDir(root string) (common, prefix string, err error) {
	common, prefix, err = Locate(root)
	if err != nil {
		return "", "", err
	}
	if err := os.MkdirAll(filepath.Join(common, "grove"), 0o755); err != nil {
		return "", "", err
	}
	return common, prefix, nil
}

// Locate returns the common Git directory and worktree prefix for root without
// creating anything. It is derived from Git rather than a worktree's .git file
// because linked worktrees have private metadata.
func Locate(root string) (common, prefix string, err error) {
	return LocateContext(context.Background(), root)
}

// LocateContext is Locate with cancellation, reported as ctx.Err().
func LocateContext(ctx context.Context, root string) (common, prefix string, err error) {
	_, common, prefix, err = IdentifyContext(ctx, root)
	return common, prefix, err
}

// IdentifyContext is LocateContext with root's own Git directory as well,
// which is unique to one worktree of one repository, from the same one
// process.
func IdentifyContext(ctx context.Context, root string) (gitDir, common, prefix string, err error) {
	if _, err := exec.LookPath("git"); err != nil {
		return "", "", "", errors.New("this command requires Git on PATH; coordination state lives in the repository's common directory")
	}
	var paths []string
	if paths, err = GitPathsContext(ctx, root, "--git-dir", "--git-common-dir", "--show-prefix"); err == nil {
		gitDir, common, prefix = paths[0], paths[1], paths[2]
	}
	if ctx.Err() != nil {
		return "", "", "", ctx.Err()
	}
	if err != nil {
		return "", "", "", fmt.Errorf("this command requires a Git repository; coordination state lives in its common directory (%w)", err)
	}
	return gitDir, common, prefix, nil
}

// GitPath answers one path-producing rev-parse option for dir, such as
// --git-dir, --git-common-dir, or --show-prefix. A path may itself contain
// newlines and trailing blanks, so each path is asked for on its own and only
// the one terminator Git appends is removed.
func GitPath(dir, option string) (string, error) {
	return GitPathContext(context.Background(), dir, option)
}

// GitPathContext is GitPath with cancellation.
func GitPathContext(ctx context.Context, dir, option string) (string, error) {
	out, err := GitContext(ctx, dir, "rev-parse", "--path-format=absolute", option)
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(out, "\n") {
		return "", fmt.Errorf("git rev-parse %s: missing output terminator", option)
	}
	return strings.TrimSuffix(out, "\n"), nil
}

// GitPathsContext answers several GitPath options with one process when it
// can. Each answer ends in one terminator, so a reply holding exactly one per
// option holds no path with a newline in it and splits safely; any other reply
// is discarded and each option is asked for on its own.
func GitPathsContext(ctx context.Context, dir string, options ...string) ([]string, error) {
	out, err := GitContext(ctx, dir, append([]string{"rev-parse", "--path-format=absolute"}, options...)...)
	if err != nil {
		return nil, err
	}
	if strings.Count(out, "\n") == len(options) && strings.HasSuffix(out, "\n") {
		return strings.Split(strings.TrimSuffix(out, "\n"), "\n"), nil
	}
	paths := make([]string, len(options))
	for i, option := range options {
		if paths[i], err = GitPathContext(ctx, dir, option); err != nil {
			return nil, err
		}
	}
	return paths, nil
}

// Git runs one git command in dir and returns its stdout.
func Git(dir string, args ...string) (string, error) {
	return GitContext(context.Background(), dir, args...)
}

// WaitDelay bounds the wait for a killed git's output pipes, which a
// descendant process may still hold open. A context that cannot be cancelled
// gets none: the commands that predate cancellation keep waiting for Git and
// whatever it started, however long a hook or helper holds the pipe.
func WaitDelay(ctx context.Context) time.Duration {
	if ctx.Done() == nil {
		return 0
	}
	return waitDelay
}

// waitDelay is a variable so the test for the expired wait need not spend it.
var waitDelay = 2 * time.Second

// GitContext is Git with cancellation: once ctx is done the process is killed
// and collected, and the error is ctx.Err() instead of a Git diagnostic.
func GitContext(ctx context.Context, dir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", dir}, args...)...)
	cmd.WaitDelay = WaitDelay(ctx)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	if err != nil {
		return "", gitError(args[0], stderr.String(), err)
	}
	return string(out), nil
}

// gitError reports Git's own words, or the failure itself when Git said nothing.
func gitError(command, stderr string, err error) error {
	if stderr = strings.TrimSpace(stderr); stderr != "" {
		return fmt.Errorf("git %s: %s", command, stderr)
	}
	return fmt.Errorf("git %s: %w", command, err)
}

// AllocatorLock serializes ID reservation; WriteLock serializes record
// publication. No command holds both at once.
func AllocatorLock(common string) (func(), error) {
	return Lock(filepath.Join(common, "grove", "lock"))
}
func WriteLock(common string) (func(), error) {
	return Lock(filepath.Join(common, "grove", "write.lock"))
}

// Lock takes an exclusive advisory lock that the kernel releases when the
// holding process exits, so a crash never leaves a stale lock behind. The
// file is never unlinked: unlinking would let a later opener lock a different
// inode while an earlier holder still believes it is exclusive.
func Lock(path string) (func(), error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		f.Close()
		return nil, err
	}
	return func() { f.Close() }, nil
}
