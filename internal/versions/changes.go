package versions

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/mascah/grove/internal/repo"
)

// Change is one file a candidate changes against its merge base with the
// target, with the lines added and removed, or -1 each for a binary file.
type Change struct {
	Path           string
	Added, Removed int
}

// Changes is what the review view shows about a candidate (G-044): its files
// against the target, whether the target already contains it, what merging
// it into the target would do (G-177), and what the branch tip changed after
// it besides the record itself, which makes the tip a new candidate.
type Changes struct {
	Base     string // the merge base of the target and the candidate; "" without a target
	OnTarget bool   // the target contains the candidate
	Merge    *Merge // nil without a target, or when it could not be predicted
	// Unpredicted says why a merge with a target could not be predicted,
	// as on a Git before 2.38, whose merge-tree cannot write a tree.
	Unpredicted string
	Files       []Change
	After       []string
}

// ChangesContext reads a candidate's changes in root's repository, on demand,
// never during a board load: every fact against one resolved target commit,
// with a merge performed in objects only when ancestry does not already
// answer it. target is the branch grove.yaml names, or "" when none applies,
// in which case Base, Merge and Files are empty. After leaves out
// recordPaths, the files of the records sharing the candidate. Once ctx is
// done the error is ctx.Err().
func ChangesContext(ctx context.Context, root, target, candidate, tip string, recordPaths ...string) (*Changes, error) {
	c := &Changes{}
	if target != "" {
		resolved, err := resolveCommits(ctx, root, "refs/heads/"+target, candidate)
		if err != nil {
			return nil, err
		}
		at, full := resolved[0], resolved[1]
		base, err := repo.GitContext(ctx, root, "merge-base", at, full)
		if err != nil {
			return nil, err
		}
		c.Base = strings.TrimSpace(base)
		c.OnTarget = c.Base == full
		if m, _, err := predict(ctx, root, at, c.Base, full); ctx.Err() != nil {
			return nil, ctx.Err()
		} else if err != nil {
			c.Unpredicted = err.Error()
		} else {
			m.Target = at
			c.Merge = &m
		}
		out, err := repo.GitContext(ctx, root, "diff", "--numstat", "-z", c.Base, candidate)
		if err != nil {
			return nil, err
		}
		if c.Files, err = numstat(out); err != nil {
			return nil, err
		}
	}
	after, err := Others(ctx, root, candidate, tip, recordPaths...)
	if err != nil {
		return nil, err
	}
	c.After = after
	return c, nil
}

// Others lists the files other than the records' that differ between two
// commits, as paths from the repository's top, where git diff prints them;
// recordPaths, several for a group sharing a candidate, are relative to the
// project, which may sit under a prefix.
func Others(ctx context.Context, root, from, to string, recordPaths ...string) ([]string, error) {
	_, _, prefix, err := repo.IdentifyContext(ctx, root)
	if err != nil {
		return nil, err
	}
	out, err := repo.GitContext(ctx, root, "diff", "--name-only", "-z", from, to)
	if err != nil {
		return nil, err
	}
	records := map[string]bool{}
	for _, r := range recordPaths {
		records[path.Join(prefix, r)] = true
	}
	var others []string
	for _, p := range strings.Split(strings.TrimSuffix(out, "\x00"), "\x00") {
		if p != "" && !records[p] {
			others = append(others, p)
		}
	}
	return others, nil
}

// numstat parses git diff --numstat -z: "added\tremoved\tpath" per entry, or
// for a rename "added\tremoved\t" then the old and the new path as their own
// entries; a binary file counts "-".
func numstat(out string) ([]Change, error) {
	fields := strings.Split(strings.TrimSuffix(out, "\x00"), "\x00")
	var files []Change
	for i := 0; i < len(fields); i++ {
		if fields[i] == "" {
			continue
		}
		parts := strings.SplitN(fields[i], "\t", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("git diff --numstat: unexpected entry %q", fields[i])
		}
		ch := Change{Path: parts[2]}
		if ch.Path == "" { // a rename: the old path, then the new
			if i+2 >= len(fields) {
				return nil, fmt.Errorf("git diff --numstat: truncated rename entry")
			}
			ch.Path = fields[i+1] + " → " + fields[i+2]
			i += 2
		}
		var err error
		for j, n := range []*int{&ch.Added, &ch.Removed} {
			if parts[j] == "-" {
				*n = -1
			} else if *n, err = strconv.Atoi(parts[j]); err != nil {
				return nil, fmt.Errorf("git diff --numstat: unexpected count in %q", fields[i])
			}
		}
		files = append(files, ch)
	}
	return files, nil
}

// DiffContext is one file's diff between two commits, as git diff prints it
// without colour: text from the repository, to be escaped before display.
func DiffContext(ctx context.Context, root, from, to, path string) (string, error) {
	// The path is one Changes listed, from the repository's top.
	return repo.GitContext(ctx, root, "diff", "--no-color", "--no-ext-diff", from, to, "--", ":(top,literal)"+path)
}
