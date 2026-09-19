package repo

import "strings"

// Worktree is one entry of the repository's worktree inventory, with the raw
// path and fields Git reports.
type Worktree struct {
	Path, Head, Branch string // Branch is "" when detached
	Prunable           string // the reason when Git could prune the entry
	Bare               bool
}

// Worktrees lists the registered worktrees of root's repository from the
// NUL-delimited porcelain format, the only one that keeps paths Git would
// otherwise quote for display (newlines, tabs, quotes) exact.
func Worktrees(root string) ([]Worktree, error) {
	out, err := Git(root, "worktree", "list", "--porcelain", "-z")
	if err != nil {
		return nil, err
	}
	return parseWorktrees(out), nil
}

func parseWorktrees(out string) []Worktree {
	var worktrees []Worktree
	var w Worktree
	for _, field := range strings.Split(out, "\x00") {
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
			w.Branch = value
		case "bare":
			w.Bare = true
		case "prunable":
			w.Prunable = value
			if w.Prunable == "" {
				w.Prunable = "no reason given"
			}
		}
	}
	return worktrees
}
