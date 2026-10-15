package gitx

import (
	"fmt"
	"strings"
)

// Worktree is one entry of a repository's worktree list.
type Worktree struct {
	Path   string // absolute
	Head   string // checked-out commit
	Branch string // short name of the checked-out branch; empty when detached
}

// Worktrees lists the worktrees of dir's repository, the main one first. It
// reads the NUL-delimited porcelain format, which keeps paths exact.
func Worktrees(dir string) ([]Worktree, error) {
	out, err := Run(dir, "worktree", "list", "--porcelain", "-z")
	if err != nil {
		return nil, err
	}
	var worktrees []Worktree
	var w Worktree
	for field := range strings.SplitSeq(out, "\x00") {
		key, value, _ := strings.Cut(field, " ")
		switch key {
		case "":
			if w.Path != "" {
				worktrees = append(worktrees, w)
			}
			w = Worktree{}
		case "worktree":
			w.Path = value
		case "HEAD":
			w.Head = value
		case "branch":
			w.Branch = strings.TrimPrefix(value, "refs/heads/")
		}
	}
	return worktrees, nil
}

// AddWorktree creates branch at commit and checks it out in a new linked
// worktree at path. Git refuses when the branch exists or path is not empty.
func AddWorktree(dir, path, branch, commit string) error {
	if _, err := Run(dir, "worktree", "add", "--quiet", "--no-track", "-b", branch, "--end-of-options", path, commit); err != nil {
		return fmt.Errorf("add worktree %s: %w", path, err)
	}
	return nil
}
