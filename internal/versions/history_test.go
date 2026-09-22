package versions

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/repo"
)

func lineage(t *testing.T, root, commit, path string) string {
	t.Helper()
	commits, err := HistoryContext(context.Background(), root, commit, path)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, c := range commits {
		if len(c.ID) != len(commit) || c.When.IsZero() {
			t.Fatalf("commit without an ID or a date: %+v", c)
		}
		out = append(out, c.Status+": "+c.Subject)
	}
	return strings.Join(out, "\n")
}

func TestHistoryFollowsOneLineOfCommits(t *testing.T) {
	t.Parallel()
	root := repoFixture(t) // G-001 proposed at "init"
	first, renamed := "grove/work/G-001-first.md", "grove/work/G-001-[re]*named.md"
	start := git(t, root, "rev-parse", "HEAD")
	write(t, root, first, record("G-001", "work", "proposed", "Main body, edited.\n"))
	commit(t, root, "edit the body")
	write(t, root, "unrelated.txt", "x\n")
	commit(t, root, "touch something else")
	write(t, root, first, record("G-001", "work", "active", "Main body, edited.\n"))
	commit(t, root, "start \x1b[31mred\u202e")
	git(t, root, "mv", first, renamed)
	commit(t, root, "rename only")
	// A status today's schema rejects is still what the record said.
	write(t, root, renamed, record("G-001", "work", "paused", "Main body, edited.\n"))
	commit(t, root, "pause")
	write(t, root, renamed, record("G-001", "work", "done", "Main body, edited.\n"))
	main := commit(t, root, "finish")

	git(t, root, "branch", "feature", start)
	wt := addWorktree(t, root, "feature-wt", "feature")
	write(t, wt, first, record("G-001", "work", "abandoned", "Main body.\n"))
	feature := commit(t, wt, "abandon on feature")

	want := "done: finish\npaused: pause\nactive: rename only\nactive: start \x1b[31mred\u202e\nproposed: edit the body\nproposed: init"
	if got := lineage(t, root, main, renamed); got != want {
		t.Errorf("main lineage:\n%q\nwant\n%q", got, want)
	}
	// The other branch has its own lineage, read from any checkout's root.
	want = "abandoned: abandon on feature\nproposed: init"
	for _, from := range []string{root, wt} {
		if got := lineage(t, from, feature, first); got != want {
			t.Errorf("feature lineage from %s:\n%q\nwant\n%q", from, got, want)
		}
	}
	// A pathspec is literal: the renamed file's name is not a pattern.
	write(t, root, "grove/work/G-001-rXnamed.md", "decoy\n")
	if got := lineage(t, root, commit(t, root, "decoy"), renamed); !strings.HasPrefix(got, "done: finish\n") {
		t.Errorf("a pattern-like name matched another file:\n%q", got)
	}
	if got := lineage(t, root, main, "grove/work/G-404-absent.md"); got != "" {
		t.Errorf("an untouched path has a lineage: %q", got)
	}
	// A commit that deleted the record has no status to show.
	git(t, wt, "rm", "-q", first)
	commit(t, wt, "drop it")
	write(t, wt, first, record("G-001", "work", "active", "Main body.\n"))
	if got := lineage(t, wt, commit(t, wt, "bring it back"), first); !strings.HasPrefix(got, "active: bring it back\n-: drop it\nabandoned: abandon on feature\n") {
		t.Errorf("delete and re-add:\n%q", got)
	}
}

// Merges are never rows, whatever they did, and they do not disturb the rows
// of the commits they brought in: the repository's own workflow of a record
// renamed and advanced on a branch, then merged without fast-forward. Dates are
// the author's, and no commit is listed above one made from it.
func TestHistoryAcrossMergesAndDates(t *testing.T) {
	root := repoFixture(t) // "init" is committed now
	first, renamed := "grove/work/G-001-first.md", "grove/work/G-001-renamed.md"
	wt := addWorktree(t, root, "feature-wt", "", "-b", "feature")
	t.Setenv("GIT_COMMITTER_DATE", "1700000100 +0000") // a clock behind init's
	git(t, wt, "mv", first, renamed)
	commit(t, wt, "rename on the branch")
	write(t, wt, renamed, record("G-001", "work", "active", "Main body.\n"))
	git(t, wt, "add", "-A")
	git(t, wt, "commit", "-q", "--date", "1700000000 +0000", "-m", "activate on the branch")
	write(t, root, "unrelated.txt", "x\n")
	commit(t, root, "main moves on")
	git(t, root, "merge", "-q", "--no-ff", "-m", "Merge branch 'feature'", "feature")
	merged := git(t, root, "rev-parse", "HEAD")
	want := "active: activate on the branch\nproposed: rename on the branch\nproposed: init"
	if got := lineage(t, root, merged, renamed); got != want {
		t.Errorf("renamed on a branch, then merged:\n%q\nwant\n%q", got, want)
	}
	commits, err := HistoryContext(context.Background(), root, merged, renamed)
	if err != nil || commits[0].When.Unix() != 1700000000 {
		t.Fatalf("the author's date: %+v %v", commits, err)
	}
	// Both sides change the status. The resolution is a merge, so it is no row;
	// the caller sees that the first row is not the record's status.
	write(t, wt, renamed, record("G-001", "work", "done", "Main body.\n"))
	commit(t, wt, "finish on feature")
	t.Setenv("GIT_COMMITTER_DATE", "1700000200 +0000")
	write(t, root, renamed, record("G-001", "work", "abandoned", "Main body.\n"))
	commit(t, root, "abandon on main")
	out, err := repo.Command(context.Background(), root, "-c", "user.name=t", "-c", "user.email=t@t", "merge", "-q", "feature").CombinedOutput()
	if err == nil {
		t.Fatalf("expected a conflict: %s", out)
	}
	if _, statErr := os.Stat(filepath.Join(root, ".git", "MERGE_HEAD")); statErr != nil {
		t.Fatalf("the merge failed without a conflict to resolve: %v\n%s", err, out)
	}
	write(t, root, renamed, record("G-001", "work", "proposed", "Resolved.\n"))
	resolved := commit(t, root, "resolve by reopening")
	want = "abandoned: abandon on main\ndone: finish on feature\n" + want
	if got := lineage(t, root, resolved, renamed); got != want {
		t.Errorf("after a conflict resolution:\n%q\nwant\n%q", got, want)
	}
}

func TestHistoryUnderPrefixAndCancellation(t *testing.T) {
	root, _ := nestedFixture(t)
	project, path := filepath.Join(root, "sub"), "grove/work/G-001-first.md"
	tip := git(t, root, "rev-parse", "feature")
	if got := lineage(t, project, tip, path); !strings.HasPrefix(got, "active: feature\nproposed: nested") {
		t.Errorf("nested lineage:\n%q", got)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if commits, err := HistoryContext(ctx, project, tip, path); !errors.Is(err, context.Canceled) || commits != nil {
		t.Fatalf("cancelled before start: %+v %v", commits, err)
	}
	block, fifo := blockingGit(t)
	for _, arg := range []string{"log", "cat-file"} {
		t.Run("blocked in "+arg, func(t *testing.T) {
			must(t, os.WriteFile(block, []byte(arg), 0o644))
			defer os.Remove(block)
			cancelDuring(t, fifo, func(ctx context.Context) (bool, error) {
				commits, err := HistoryContext(ctx, project, tip, path)
				return commits != nil, err
			})
		})
	}
}
