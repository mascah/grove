package versions

import (
	"bytes"
	"cmp"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/mascah/grove/internal/repo"
)

// Merge predicts what merging a commit into the target would do now (G-177),
// read-only: nothing is checked out or written to a ref, and at most objects
// are written, as a merge would. It is a fact about the target commit it
// names, stale once the target moves, and says nothing about whether the
// merged changes work together.
type Merge struct {
	Target    string   `json:"target"`    // the target commit it was computed against
	Commit    string   `json:"commit"`    // the commit merged
	Outcome   string   `json:"outcome"`   // integrated, fast-forward, clean, or conflict
	Conflicts []string `json:"conflicts"` // the conflicting files, from the repository's top
}

// Text is the fact as every surface prints it, target being the branch name.
func (m *Merge) Text(target string) string {
	at := target + " at " + m.Target[:min(len(m.Target), 7)]
	switch m.Outcome {
	case "integrated":
		return "integrated: " + at + " holds it"
	case "fast-forward":
		return "merges into " + at + " as a fast-forward"
	case "clean":
		return "merges cleanly into " + at + ", which moved since the branch left it"
	}
	return "conflicts with " + at + " " + m.Where()
}

// Where names a conflict's files; Git names none for some, such as a split
// directory rename.
func (m *Merge) Where() string {
	if len(m.Conflicts) == 0 {
		return "where Git names no file"
	}
	return "in " + strings.Join(m.Conflicts, ", ")
}

// PredictContext resolves target to a commit once and predicts merging the
// commits into it in the order given, each into the result of the ones before
// it, stopping after the first conflict. A clean result that a later merge
// reads becomes an unreferenced commit object, which Git's garbage collection
// removes. One commit is the plain question: does it merge into the target?
func PredictContext(ctx context.Context, root, target string, commits []string) ([]Merge, error) {
	prefix, resolved, err := resolveCommits(ctx, root, append([]string{target}, commits...)...)
	if err != nil {
		return nil, err
	}
	tip := resolved[0]
	var out []Merge
	ours := tip
	for i, c := range resolved[1:] {
		base, err := repo.GitContext(ctx, root, "merge-base", ours, c)
		if err != nil {
			return nil, err
		}
		m, tree, err := predict(ctx, root, prefix, ours, strings.TrimSpace(base), c)
		if err != nil {
			return nil, err
		}
		m.Target = tip
		out = append(out, m)
		switch {
		case m.Outcome == "conflict" || i == len(commits)-1:
			return out, nil
		case m.Outcome == "fast-forward":
			ours = c
		case m.Outcome == "clean":
			// A fixed identity: the object is never shown, and a repository
			// without user.name must not refuse it.
			made, err := repo.GitContext(ctx, root, "-c", "user.name=grove", "-c", "user.email=grove@invalid",
				"commit-tree", tree, "-p", ours, "-p", c, "-m", "grove merge prediction")
			if err != nil {
				return nil, err
			}
			ours = strings.TrimSpace(made)
		}
	}
	return out, nil
}

// resolveCommits resolves names, which may abbreviate, to full commits in
// one process, so that ancestry compares commits by their full names, and
// reads root's prefix in the repository with them.
func resolveCommits(ctx context.Context, root string, names ...string) (prefix string, commits []string, err error) {
	args := []string{"rev-parse", "--show-prefix"} // which, without --verify, would echo --end-of-options
	for _, n := range names {
		if strings.HasPrefix(n, "-") {
			return "", nil, fmt.Errorf("%q is not a commit", n)
		}
		args = append(args, n+"^{commit}")
	}
	out, err := repo.GitContext(ctx, root, args...)
	if err != nil {
		return "", nil, err
	}
	// ponytail: a prefix holding a newline misreads; such directories are not supported here.
	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	if len(lines) != len(names)+1 {
		return "", nil, fmt.Errorf("git rev-parse: expected %d commits, read %d", len(names), len(lines)-1)
	}
	return lines[0], lines[1:], nil
}

// predict classifies merging commit into ours, whose merge base is base: an
// ancestry answer needs no merge, and otherwise git merge-tree performs it in
// objects only. tree is the merged tree of a clean merge. merge-tree names
// files from root, whose prefix in the repository turns them into paths
// from its top.
func predict(ctx context.Context, root, prefix, ours, base, commit string) (m Merge, tree string, err error) {
	m = Merge{Commit: commit, Conflicts: []string{}}
	switch base {
	case commit:
		m.Outcome = "integrated"
		return m, "", nil
	case ours:
		m.Outcome = "fast-forward"
		return m, "", nil
	}
	// Exit 1 is Git's answer, a conflict, with the files on stdout, which
	// GitContext drops on failure.
	cmd := repo.Command(ctx, root, "merge-tree", "--write-tree", "--name-only", "--no-messages", "-z", ours, commit)
	cmd.WaitDelay = repo.WaitDelay(ctx)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if ctx.Err() != nil {
		return m, "", ctx.Err()
	}
	var exit *exec.ExitError
	conflict := errors.As(err, &exit) && exit.ExitCode() == 1
	if err != nil && !conflict {
		return m, "", fmt.Errorf("git merge-tree: %s", cmp.Or(strings.TrimSpace(stderr.String()), err.Error()))
	}
	fields := strings.Split(strings.TrimSuffix(string(out), "\x00"), "\x00")
	if conflict {
		m.Outcome = "conflict"
		for _, f := range fields[1:] {
			if f != "" {
				m.Conflicts = append(m.Conflicts, path.Join(prefix, f))
			}
		}
		return m, "", nil
	}
	m.Outcome = "clean"
	return m, fields[0], nil
}
