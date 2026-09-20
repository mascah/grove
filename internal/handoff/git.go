package handoff

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mascah/grove/internal/repo"
)

// gitIdentity reports the checkout holding root, or nil when no .git entry
// encloses it. Once a repository is detected, a failing or missing Git is an
// error rather than a quiet "not a repository".
func gitIdentity(ctx context.Context, root string) (*Git, error) {
	for dir := root; ; dir = filepath.Dir(dir) {
		_, err := os.Lstat(filepath.Join(dir, ".git"))
		if err == nil {
			break
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		if dir == filepath.Dir(dir) {
			return nil, nil
		}
	}
	paths, err := repo.GitPathsContext(ctx, root, "--show-toplevel", "--git-common-dir")
	if err != nil {
		return nil, err
	}
	g := &Git{Checkout: paths[0], CommonDir: paths[1]}
	// Both questions exit 1 without a word when the answer is "none": a
	// detached HEAD has no ref, and a branch without commits has no HEAD.
	for _, q := range []struct {
		into *string
		args []string
	}{
		{&g.Ref, []string{"symbolic-ref", "-q", "HEAD"}},
		{&g.Head, []string{"rev-parse", "-q", "--verify", "HEAD^{commit}"}},
	} {
		out, err := repo.GitContext(ctx, root, q.args...)
		var exit *exec.ExitError
		if err != nil && !(errors.As(err, &exit) && exit.ExitCode() == 1) {
			return nil, err
		}
		*q.into = strings.TrimSuffix(out, "\n")
	}
	return g, nil
}
