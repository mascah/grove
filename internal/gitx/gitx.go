// Package gitx is the only place Git runs. Commands are run by argv, never
// through a shell, always for the directory the caller names.
package gitx

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Run runs git with args in dir and returns its standard output. A failure
// names the subcommand and carries Git's standard error.
func Run(dir string, args ...string) (string, error) {
	out, err := run(command(dir, args...))
	return string(out), err
}

// Input runs git with args in dir, given stdin as its standard input, and
// returns its standard output.
func Input(dir, stdin string, args ...string) (string, error) {
	cmd := command(dir, args...)
	cmd.Stdin = strings.NewReader(stdin)
	out, err := run(cmd)
	return string(out), err
}

// Object is an object git cat-file --batch read: its ID and contents, or
// neither for a name that is no object.
type Object struct {
	ID   string
	Data []byte
}

// Objects reads each of names, a revision such as <commit>:<path>, in one
// git cat-file --batch, in order.
func Objects(dir string, names []string) ([]Object, error) {
	if len(names) == 0 {
		return nil, nil
	}
	out, err := Input(dir, strings.Join(names, "\n")+"\n", "cat-file", "--batch")
	if err != nil {
		return nil, err
	}
	objects := make([]Object, len(names))
	// Each reply is "<id> <type> <size>\n<contents>\n", or "<name> missing\n"
	// or "<name> ambiguous\n".
	for i := range objects {
		header, rest, _ := strings.Cut(out, "\n")
		f := strings.Fields(header)
		if out = rest; f[len(f)-1] == "missing" || f[len(f)-1] == "ambiguous" {
			continue
		}
		size, _ := strconv.Atoi(f[len(f)-1])
		objects[i] = Object{f[0], []byte(rest[:size])}
		out = rest[size+1:]
	}
	return objects, nil
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
