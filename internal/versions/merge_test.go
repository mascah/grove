package versions

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/repo"
)

// TestPredictMerges covers each outcome, alone and in a stated order, and
// shows that nothing but objects is written.
func TestPredictMerges(t *testing.T) {
	t.Parallel()
	root := repoFixture(t)
	ctx := context.Background()
	write(t, root, "shared.txt", "one\ntwo\n")
	commit(t, root, "shared")
	branch := func(name, file, content string) string {
		wt := addWorktree(t, root, name, "main", "-b", name)
		write(t, wt, file, content)
		return commit(t, wt, name)
	}
	ff := branch("ff", "ff.txt", "ff\n")
	left := branch("left", "shared.txt", "LEFT\ntwo\n")
	right := branch("right", "shared.txt", "RIGHT\ntwo\n")
	other := branch("other", "other.txt", "other\n")
	state := func() string {
		return git(t, root, "for-each-ref") + git(t, root, "rev-parse", "HEAD") + git(t, root, "status", "--porcelain")
	}
	one := func(target, c string) Merge {
		t.Helper()
		ms, err := PredictContext(ctx, root, target, []string{c})
		if err != nil || len(ms) != 1 {
			t.Fatalf("predict %s: %v %+v", c, err, ms)
		}
		return ms[0]
	}
	before := state()
	start := git(t, root, "rev-parse", "HEAD")
	if m := one("main", ff[:7]); m.Outcome != "fast-forward" || m.Target != start || m.Commit != ff {
		t.Fatalf("fast-forward: %+v", m)
	}

	// The target moves: a clean merge, and a textual conflict named by file.
	write(t, root, "shared.txt", "MAIN\ntwo\n")
	moved := commit(t, root, "main moved")
	before = state()
	if m := one("refs/heads/main", ff); m.Outcome != "clean" || m.Target != moved {
		t.Fatalf("clean: %+v", m)
	}
	m := one("main", left)
	if m.Outcome != "conflict" || !reflect.DeepEqual(m.Conflicts, []string{"shared.txt"}) {
		t.Fatalf("conflict: %+v", m)
	}
	if got := m.Text("main"); got != "conflicts with main at "+moved[:7]+" in shared.txt" {
		t.Fatalf("text: %q", got)
	}
	// The same answer through the review view's read.
	c, err := ChangesContext(ctx, root, "main", left, left, "")
	if err != nil || c.OnTarget || c.Merge == nil || !reflect.DeepEqual(*c.Merge, m) {
		t.Fatalf("changes: %+v %v", c, err)
	}

	// In a stated order the second merge reads the first: other then ff are
	// both clean, and left after right conflicts although each alone would
	// conflict only with main. The sequence stops at the first conflict.
	git(t, root, "reset", "-q", "--hard", start)
	before = state()
	ms, err := PredictContext(ctx, root, "main", []string{other, right, left, ff})
	if err != nil || len(ms) != 3 {
		t.Fatalf("sequence: %v %+v", err, ms)
	}
	if ms[0].Outcome != "fast-forward" || ms[1].Outcome != "clean" || ms[2].Outcome != "conflict" || !reflect.DeepEqual(ms[2].Conflicts, []string{"shared.txt"}) || ms[2].Target != start {
		t.Fatalf("sequence: %+v", ms)
	}
	if m := one("main", left); m.Outcome != "fast-forward" {
		t.Fatalf("left alone: %+v", m)
	}
	if after := state(); after != before {
		t.Fatalf("refs, HEAD or the checkout changed:\n%s\n%s", before, after)
	}

	// Once merged, the target holds it.
	git(t, root, "merge", "-q", "--ff-only", "left")
	if m := one("main", left); m.Outcome != "integrated" || m.Text("main") != "integrated: main at "+left[:7]+" holds it" {
		t.Fatalf("integrated: %+v", m)
	}
	c, err = ChangesContext(ctx, root, "main", left[:9], left, "")
	if err != nil || !c.OnTarget || c.Merge.Outcome != "integrated" {
		t.Fatalf("changes, abbreviated candidate: %+v %v", c, err)
	}
	// A split directory rename conflicts in no file Git names.
	write(t, root, "x/a", "a\n")
	write(t, root, "x/b", "b\n")
	commit(t, root, "x")
	split := addWorktree(t, root, "split", "main", "-b", "split")
	write(t, split, "y/a", "a\n")
	write(t, split, "z/b", "b\n")
	git(t, split, "rm", "-q", "x/a", "x/b")
	splitTip := commit(t, split, "split x")
	write(t, root, "x/c", "c\n")
	now := commit(t, root, "x/c")
	if m := one("main", splitTip); m.Outcome != "conflict" || len(m.Conflicts) != 0 || m.Text("main") != "conflicts with main at "+now[:7]+" where Git names no file" {
		t.Fatalf("split rename: %+v %q", m, m.Text("main"))
	}
	if _, err := PredictContext(ctx, root, "nope", []string{left}); err == nil {
		t.Fatal("an unknown target predicted")
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := PredictContext(cancelled, root, "main", []string{right}); err != context.Canceled {
		t.Fatalf("cancelled: %v", err)
	}
}

// TestResolution finds the latest target commit merged into a candidate's
// branch, the files whose merged content is neither side's, and the
// candidate the record named before, and ignores a merge of another branch.
func TestResolution(t *testing.T) {
	t.Parallel()
	root := repoFixture(t)
	ctx := context.Background()
	const path = "grove/work/G-001-first.md"
	write(t, root, "shared.txt", "one\n")
	commit(t, root, "shared")
	wt := addWorktree(t, root, "work", "main", "-b", "work")
	write(t, wt, "shared.txt", "work\n")
	previous := commit(t, wt, "work")
	write(t, wt, path, strings.Replace(record("G-001", "work", "active", "Body.\n"), "status: active\n", "status: active\ncandidate: \""+previous+"\"\n", 1))
	commit(t, wt, "feedback")
	changes := func(candidate string) *Changes {
		t.Helper()
		c, err := ChangesContext(ctx, root, "main", candidate, candidate, path)
		if err != nil {
			t.Fatal(err)
		}
		return c
	}
	if c := changes(previous); c.Resolution != nil {
		t.Fatalf("no merge yet: %+v", c.Resolution)
	}
	write(t, root, "shared.txt", "main\n")
	write(t, root, "main.txt", "main\n")
	target := commit(t, root, "main moves")
	repo.Command(ctx, wt, "merge", "-q", "main").Run() // conflicts in shared.txt
	write(t, wt, "shared.txt", "both\n")
	merge := commit(t, wt, "resolved")
	want := &Resolution{Merge: merge, Target: target, Previous: previous, Files: []Resolved{{"shared.txt", ""}}}
	if c := changes(merge); !reflect.DeepEqual(c.Resolution, want) {
		t.Fatalf("resolution %+v, want %+v", c.Resolution, want)
	}
	// A conflict settled by taking the target's side is named as such, and a
	// file Git merged by itself, each side changing its own line, is not one.
	write(t, root, "lines.txt", "a\nb\nc\nd\ne\n")
	commit(t, root, "lines")
	git(t, wt, "merge", "-q", "--no-edit", "main")
	write(t, wt, "shared.txt", "work again\n")
	write(t, wt, "lines.txt", "A\nb\nc\nd\ne\n")
	commit(t, wt, "more work")
	write(t, root, "shared.txt", "main again\n")
	write(t, root, "lines.txt", "a\nb\nc\nd\nE\n")
	target = commit(t, root, "main again")
	repo.Command(ctx, wt, "merge", "-q", "main").Run()
	git(t, wt, "checkout", "--theirs", "shared.txt")
	merge = commit(t, wt, "took main's")
	if c := changes(merge); c.Resolution == nil || c.Resolution.Target != target || !reflect.DeepEqual(c.Resolution.Files, []Resolved{{"shared.txt", "target"}}) {
		t.Fatalf("taking a side: %+v", c.Resolution)
	}
	// Read from a project under a prefix, the paths are still the top's.
	if c, err := ChangesContext(ctx, filepath.Join(root, "grove"), "main", merge, merge); err != nil || c.Resolution == nil || !reflect.DeepEqual(c.Resolution.Files, []Resolved{{"shared.txt", "target"}}) {
		t.Fatalf("under a prefix: %+v %v", c.Resolution, err)
	}
	// Renamed apart on each side, the branch's name kept: the target's name
	// is gone, which is the branch's side, whatever diff.renames says.
	write(t, root, "old.txt", "one\ntwo\nthree\nfour\n")
	commit(t, root, "old")
	git(t, wt, "merge", "-q", "--no-edit", "main")
	git(t, wt, "mv", "old.txt", "ours.txt")
	commit(t, wt, "ours")
	git(t, root, "mv", "old.txt", "theirs.txt")
	commit(t, root, "theirs")
	repo.Command(ctx, wt, "merge", "-q", "main").Run()
	git(t, wt, "rm", "-q", "--cached", "--ignore-unmatch", "old.txt", "theirs.txt")
	os.Remove(filepath.Join(wt, "theirs.txt"))
	git(t, wt, "add", "ours.txt")
	merge = commit(t, wt, "kept ours")
	got := map[string]string{}
	for _, f := range changes(merge).Resolution.Files {
		got[f.Path] = f.Kept
	}
	if got["theirs.txt"] != "branch" || got["ours.txt"] != "branch" || got["old.txt"] != "" {
		t.Fatalf("renamed apart: %v", got)
	}
	// Only the latest merge is read: one of another branch, or of unrelated
	// history, is no target's.
	side := addWorktree(t, root, "side", previous, "-b", "side")
	write(t, side, "side.txt", "side\n")
	commit(t, side, "side")
	git(t, wt, "merge", "-q", "--no-edit", "side")
	if c := changes(git(t, wt, "rev-parse", "HEAD")); c.Resolution != nil {
		t.Fatalf("a merge of side: %+v", c.Resolution)
	}
	orphan := addWorktree(t, root, "orphan", "", "--orphan", "-b", "orphan")
	write(t, orphan, "alone.txt", "alone\n")
	commit(t, orphan, "alone")
	git(t, wt, "merge", "-q", "--no-edit", "--allow-unrelated-histories", "orphan")
	if c := changes(git(t, wt, "rev-parse", "HEAD")); c.Resolution != nil {
		t.Fatalf("a merge of unrelated history: %+v", c.Resolution)
	}
}
