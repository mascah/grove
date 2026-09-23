package versions

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

// TestChangesAgainstTheTarget reads a candidate's files from the merge base,
// what the tip changed after it, and whether the target holds it.
func TestChangesAgainstTheTarget(t *testing.T) {
	t.Parallel()
	root := repoFixture(t)
	wt := addWorktree(t, root, "feature", "", "-b", "feature")
	write(t, wt, "grove/work/G-001-first.md", record("G-001", "work", "active", "Feature body.\n"))
	write(t, wt, "code.txt", "one\ntwo\n")
	candidate := commit(t, wt, "feat: implement")
	write(t, wt, "grove/work/G-001-first.md", record("G-001", "work", "review", "Feature body.\n"))
	tip := commit(t, wt, "docs: review")
	write(t, root, "notes.txt", "elsewhere\n")
	base := git(t, root, "rev-parse", "HEAD")
	commit(t, root, "docs: main moved")
	ctx := context.Background()
	c, err := ChangesContext(ctx, root, "main", candidate, tip, "grove/work/G-001-first.md")
	if err != nil {
		t.Fatal(err)
	}
	if c.Base != base || c.OnTarget || len(c.After) != 0 {
		t.Fatalf("%+v", c)
	}
	if len(c.Files) != 2 || c.Files[0] != (Change{"code.txt", 2, 0}) || c.Files[1].Path != "grove/work/G-001-first.md" || c.Files[1].Added != 2 || c.Files[1].Removed != 2 { // status and body
		t.Fatalf("files: %+v", c.Files)
	}
	diff, err := DiffContext(ctx, root, c.Base, candidate, "code.txt")
	if err != nil || !strings.Contains(diff, "+one\n+two\n") || strings.Contains(diff, "\x1b") {
		t.Fatalf("diff: %v\n%s", err, diff)
	}
	// A commit after the candidate that changes anything but the record
	// shows as After; a binary file and a rename are read too.
	write(t, wt, "code.txt", "one\ntwo\nthree\n")
	write(t, wt, "blob.bin", "\x00\x01\x02")
	git(t, wt, "mv", "grove/work/G-001-first.md", "grove/work/G-001-moved.md")
	tip = commit(t, wt, "fix: later")
	c, err = ChangesContext(ctx, root, "main", candidate, tip, "grove/work/G-001-moved.md")
	if err != nil || strings.Join(c.After, ",") != "blob.bin,code.txt" { // a rename shows its new path, which is the record
		t.Fatalf("after: %+v %v", c, err)
	}
	c, err = ChangesContext(ctx, root, "main", tip, tip, "grove/work/G-001-moved.md")
	if got := c.Files; err != nil || len(c.After) != 0 || len(got) != 3 || got[0] != (Change{"blob.bin", -1, -1}) || got[1] != (Change{"code.txt", 3, 0}) || got[2].Path != "grove/work/G-001-first.md → grove/work/G-001-moved.md" {
		t.Fatalf("files: %+v %v", c, err)
	}
	// Without a target there is no base and no file list, only the tip check.
	c, err = ChangesContext(ctx, root, "", candidate, tip, "grove/work/G-001-moved.md")
	if err != nil || c.Base != "" || c.Files != nil || len(c.After) != 2 {
		t.Fatalf("no target: %+v %v", c, err)
	}
	// Once merged, the target holds the candidate.
	git(t, root, "merge", "-q", "--no-edit", "feature")
	if c, err = ChangesContext(ctx, root, "main", candidate, tip, ""); err != nil || !c.OnTarget {
		t.Fatalf("after the merge: %+v %v", c, err)
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := ChangesContext(cancelled, root, "main", candidate, tip, ""); err != context.Canceled {
		t.Fatalf("cancelled: %v", err)
	}
	if _, err := DiffContext(cancelled, root, base, candidate, "code.txt"); err != context.Canceled {
		t.Fatalf("cancelled diff: %v", err)
	}
}

// A project under a prefix: the record's own commits are not "other files",
// code above the project is, and a diff reads the top-relative path.
func TestChangesUnderAPrefix(t *testing.T) {
	t.Parallel()
	root, wt := nestedFixture(t)
	project := filepath.Join(root, "sub")
	candidate := git(t, wt, "rev-parse", "HEAD")
	write(t, wt, "sub/grove/work/G-001-first.md", record("G-001", "work", "review", "Feature.\n"))
	tip := commit(t, wt, "docs: review")
	ctx := context.Background()
	c, err := ChangesContext(ctx, project, "main", candidate, tip, "grove/work/G-001-first.md")
	if err != nil || len(c.After) != 0 || len(c.Files) != 1 || c.Files[0].Path != "sub/grove/work/G-001-first.md" {
		t.Fatalf("%+v %v", c, err)
	}
	if diff, err := DiffContext(ctx, project, c.Base, candidate, c.Files[0].Path); err != nil || !strings.Contains(diff, "+status: active") {
		t.Fatalf("diff under a prefix: %v\n%s", err, diff)
	}
	write(t, wt, "code.txt", "above the project\n")
	tip = commit(t, wt, "feat: code")
	if others, err := Others(ctx, project, candidate, tip, "grove/work/G-001-first.md"); err != nil || strings.Join(others, ",") != "code.txt" {
		t.Fatalf("others: %v %v", others, err)
	}
}
