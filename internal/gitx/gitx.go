// Package gitx is the only place Git runs. Commands are run by argv, never
// through a shell, always for the directory the caller names.
package gitx

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Run runs git with args in dir and returns its standard output. A failure
// names the subcommand and carries Git's standard error.
func Run(dir string, args ...string) (string, error) {
	out, err := run(command(dir, args...))
	return string(out), err
}

func command(dir string, args ...string) *exec.Cmd {
	return exec.Command("git", append([]string{"-C", dir}, args...)...)
}

func run(cmd *exec.Cmd) ([]byte, error) {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git %s: %s: %w", cmd.Args[3], strings.TrimSpace(stderr.String()), err)
	}
	return out, nil
}

// TopLevel returns the root of the worktree that holds dir.
func TopLevel(dir string) (string, error) {
	out, err := Run(dir, "rev-parse", "--show-toplevel")
	return strings.TrimSuffix(out, "\n"), err
}

// ResolveBranch returns the newest commit that both the local branch name and
// origin's hold, or the one of them that exists: a commit origin has that is
// never ahead of this checkout, so what is recorded from it stays resolvable.
func ResolveBranch(dir, name string) (string, error) {
	remote, err := resolve(dir, "refs/remotes/origin/"+name)
	if err != nil {
		return LocalBranch(dir, name)
	}
	local, err := LocalBranch(dir, name)
	if err != nil {
		return remote, nil
	}
	out, err := Run(dir, "merge-base", local, remote)
	return strings.TrimSpace(out), err
}

// LocalBranch returns the commit at the tip of the local branch name.
func LocalBranch(dir, name string) (string, error) {
	return resolve(dir, "refs/heads/"+name)
}

func resolve(dir, ref string) (string, error) {
	out, err := Run(dir, "rev-parse", "--verify", "--end-of-options", ref+"^{commit}")
	return strings.TrimSpace(out), err
}

// Shallow reports whether the repository at dir is a shallow clone, which
// lacks the commits before the ones it fetched.
func Shallow(dir string) bool {
	out, err := Run(dir, "rev-parse", "--is-shallow-repository")
	return err == nil && strings.TrimSpace(out) == "true"
}
