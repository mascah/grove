// Package repo locates the Git common directory and takes the advisory locks
// that Grove's mutating commands share across linked worktrees.
package repo

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
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
	if _, err := exec.LookPath("git"); err != nil {
		return "", "", errors.New("this command requires Git on PATH; coordination state lives in the repository's common directory")
	}
	if common, err = GitPath(root, "--git-common-dir"); err == nil {
		prefix, err = GitPath(root, "--show-prefix")
	}
	if err != nil {
		return "", "", fmt.Errorf("this command requires a Git repository; coordination state lives in its common directory (%w)", err)
	}
	return common, prefix, nil
}

// GitPath answers one path-producing rev-parse option for dir, such as
// --git-dir, --git-common-dir, or --show-prefix. A path may itself contain
// newlines and trailing blanks, so each path is asked for on its own and only
// the one terminator Git appends is removed.
func GitPath(dir, option string) (string, error) {
	out, err := Git(dir, "rev-parse", "--path-format=absolute", option)
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(out, "\n") {
		return "", fmt.Errorf("git rev-parse %s: missing output terminator", option)
	}
	return strings.TrimSuffix(out, "\n"), nil
}

// Git runs one git command in dir and returns its stdout.
func Git(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git %s: %s", args[0], strings.TrimSpace(stderr.String()))
	}
	return string(out), nil
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
