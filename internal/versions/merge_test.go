package versions

import (
	"context"
	"reflect"
	"testing"
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
