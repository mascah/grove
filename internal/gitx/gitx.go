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
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git %s: %s: %w", args[0], strings.TrimSpace(stderr.String()), err)
	}
	return string(out), nil
}

// TopLevel returns the root of the worktree that holds dir.
func TopLevel(dir string) (string, error) {
	out, err := Run(dir, "rev-parse", "--show-toplevel")
	return strings.TrimSuffix(out, "\n"), err
}
