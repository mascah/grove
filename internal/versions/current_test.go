package versions

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/repo"
)

// standing summarizes one group as "place status current|older" lines, with
// deleted rows as status "-", and its notes.
func standing(g Group) []string {
	var lines []string
	for _, v := range g.Versions {
		place := "branch " + strings.TrimPrefix(v.Source.Ref, "refs/heads/")
		if v.Source.Kind == "live" {
			place = "checkout " + v.Source.Locator
		}
		status := "-"
		if v.Record != nil {
			status = v.Record.Status
		}
		state := "current"
		if v.Older != "" {
			state = "older"
		}
		lines = append(lines, place+" "+status+" "+state)
	}
	return append(lines, g.Notes...)
}

func expectStanding(t *testing.T, res *Result, id string, want ...string) {
	t.Helper()
	if got := standing(group(t, res, id)); !slices.Equal(got, want) {
		t.Errorf("%s:\n got %q\nwant %q", id, got, want)
	}
}

// TestCurrentView builds one repository holding every case of G-042's
// projection and checks each record's current and older observations.
func TestCurrentView(t *testing.T) {
	t.Parallel()
	root := repoFixture(t)
	rec := func(id, status, body string) { write(t, root, "grove/"+id+".md", record(id, "work", status, body)) }
	for _, id := range []string{"G-010", "G-011", "G-013", "G-014", "G-015", "G-016", "G-017"} {
		rec(id, "proposed", "Start.\n")
	}
	c0 := commit(t, root, "records")

	// A stale branch with its own commits, none touching the record folder:
	// the G-030/G-023 shape, where its old Proposed copy must not obscure main.
	git(t, root, "branch", "stale")
	stale := addWorktree(t, root, "stale", "stale")
	write(t, stale, "notes.txt", "unrelated\n")
	commit(t, stale, "unrelated work")
	git(t, root, "worktree", "remove", stale)

	// A feature branch with unmerged progress, work only it has, a change
	// that diverges from main's, and a deletion.
	feature := addWorktree(t, root, "feature", c0, "-b", "feature")
	write(t, feature, "grove/G-011.md", record("G-011", "work", "active", "Start.\n"))
	write(t, feature, "grove/G-012.md", record("G-012", "work", "proposed", "Only here.\n"))
	write(t, feature, "grove/G-013.md", record("G-013", "work", "active", "Start.\n"))
	if err := os.Remove(filepath.Join(feature, "grove/G-014.md")); err != nil {
		t.Fatal(err)
	}
	commit(t, feature, "feature progress")

	// A branch merged into main, whose record main then moved on.
	git(t, root, "checkout", "-q", "-b", "merged", c0)
	rec("G-016", "active", "Start.\n")
	commit(t, root, "merged progress")
	git(t, root, "checkout", "-q", "main")
	git(t, root, "merge", "-q", "--no-ff", "-m", "merge", "merged")

	rec("G-010", "done", "Start.\n")
	rec("G-013", "proposed", "Main edited this.\n")
	rec("G-015", "active", "Start.\n")
	commit(t, root, "main moves on")
	git(t, root, "branch", "mid")
	rec("G-015", "proposed", "Start.\n") // a revert to c0's exact bytes
	rec("G-016", "done", "Start.\n")
	commit(t, root, "revert G-015, finish G-016")

	// A branch from before schema 3 is an invalid source, not an observation.
	git(t, root, "checkout", "-q", "-b", "oldschema", c0)
	write(t, root, "grove.yaml", "schema_version: 2\nrecords: grove\n")
	commit(t, root, "old schema")
	git(t, root, "checkout", "-q", "main")

	addWorktree(t, root, "old", c0, "--detach")
	// Uncommitted: an edit and a new record in the feature checkout,
	// a deletion in the main checkout.
	write(t, feature, "grove/G-011.md", record("G-011", "work", "active", "Working on it.\n"))
	write(t, feature, "grove/G-019.md", record("G-019", "work", "proposed", "New.\n"))
	if err := os.Remove(filepath.Join(root, "grove/G-017.md")); err != nil {
		t.Fatal(err)
	}

	res := mustInspect(t, root, "")
	if res.Complete || source(t, res, "committed", "refs/heads/oldschema").Valid {
		t.Fatalf("the old-schema branch should make the result incomplete: %s", dump(res))
	}
	for _, s := range res.Sources {
		if s.Kind == "live" && !s.Valid {
			t.Fatalf("checkout %s: %v", s.Worktree, s.Diagnostics)
		}
	}
	// Sources: branches feature, main, merged, mid, oldschema (invalid),
	// stale; checkouts . (main), feature, old (detached at c0).
	expectStanding(t, res, "G-010",
		"branch feature proposed older", "branch main done current", "branch merged proposed older",
		"branch mid done current", "branch stale proposed older",
		"checkout . done current", "checkout feature proposed older", "checkout old proposed older")
	expectStanding(t, res, "G-011",
		"branch feature active older", "branch main proposed older", "branch merged proposed older",
		"branch mid proposed older", "branch stale proposed older",
		"checkout . proposed older", "checkout feature active current", "checkout old proposed older")
	expectStanding(t, res, "G-012", "branch feature proposed current", "checkout feature proposed current")
	expectStanding(t, res, "G-013",
		"branch feature active current", "branch main proposed current", "branch merged proposed older",
		"branch mid proposed current", "branch stale proposed older",
		"checkout . proposed current", "checkout feature active current", "checkout old proposed older")
	expectStanding(t, res, "G-014",
		"branch feature - current", "branch main proposed older", "branch merged proposed older",
		"branch mid proposed older", "branch stale proposed older",
		"checkout . proposed older", "checkout old proposed older")
	// The revert: main's bytes equal c0's, and still main's state is current
	// while mid's later-looking active is older. The stale branch's identical
	// bytes are older than mid's change, since it never made one.
	expectStanding(t, res, "G-015",
		"branch feature proposed older", "branch main proposed current", "branch merged proposed older",
		"branch mid active older", "branch stale proposed older",
		"checkout . proposed current", "checkout feature proposed older", "checkout old proposed older")
	expectStanding(t, res, "G-016",
		"branch feature proposed older", "branch main done current", "branch merged active older",
		"branch mid active older", "branch stale proposed older",
		"checkout . done current", "checkout feature proposed older", "checkout old proposed older")
	expectStanding(t, res, "G-017",
		"branch feature proposed older", "branch main proposed older", "branch merged proposed older",
		"branch mid proposed older", "branch stale proposed older",
		"checkout . - current", "checkout feature proposed older", "checkout old proposed older")
	expectStanding(t, res, "G-019", "checkout feature proposed current")

	if why := find(t, group(t, res, "G-011"), "committed", "refs/heads/feature").Older; why != "checkout feature (feature) has an uncommitted change to it on top of this commit" {
		t.Errorf("reason for the committed copy under an uncommitted edit: %q", why)
	}
	// Main and mid hold the same current bytes; either may be named.
	if why := find(t, group(t, res, "G-010"), "committed", "refs/heads/stale").Older; why != "branch main changed it since their common history" && why != "branch mid changed it since their common history" {
		t.Errorf("reason for the stale copy: %q", why)
	}

	// Any checkout sees the same view.
	other := mustInspect(t, feature, "")
	for i := range res.Groups {
		if a, b := standing(res.Groups[i]), standing(other.Groups[i]); !slices.Equal(a, b) {
			t.Errorf("%s differs by invoking checkout:\n%q\n%q", res.Groups[i].ID, a, b)
		}
	}
	if one := mustInspect(t, root, "G-013"); !slices.Equal(standing(one.Groups[0]), standing(group(t, res, "G-013"))) {
		t.Errorf("one record's view should equal its part of the whole")
	}
}

// TestCurrentViewUnorderedPair shows ambiguity explicitly: where the record
// cannot be read at two branches' common commit, they cannot be ordered, so
// both stay current and a note names them. A project that does not validate
// there does not matter while the record itself reads at its path.
func TestCurrentViewUnorderedPair(t *testing.T) {
	t.Parallel()
	root := repoFixture(t)
	write(t, root, "grove/G-020.md", "---\nid: [broken\n---\n")
	write(t, root, "grove/G-021.md", record("G-021", "work", "proposed", "Start.\n"))
	base := commit(t, root, "base does not validate")
	write(t, root, "grove/G-020.md", record("G-020", "work", "proposed", "Fixed on main.\n"))
	commit(t, root, "fix on main")
	git(t, root, "checkout", "-q", "-b", "other", base)
	write(t, root, "grove/G-020.md", record("G-020", "work", "proposed", "Fixed on other.\n"))
	write(t, root, "grove/G-021.md", record("G-021", "work", "active", "Start.\n"))
	commit(t, root, "fix and change on other")
	git(t, root, "checkout", "-q", "main")

	res := mustInspect(t, root, "")
	expectStanding(t, res, "G-021", "branch main proposed older", "branch other active current", "checkout . proposed older")
	want := fmt.Sprintf("branch main and branch other could not be ordered: their common commit %s holds a project that does not validate", base[:12])
	expectStanding(t, res, "G-020",
		"branch main proposed current", "branch other proposed current", "checkout . proposed current", want)
}

// TestMergeBases checks the walk against Git's own answers, including a
// criss-cross history with two bases and commits sharing one timestamp.
func TestMergeBases(t *testing.T) {
	t.Parallel()
	root := repoFixture(t)
	// Every commit here shares one committer time.
	dated := func(args ...string) {
		t.Helper()
		cmd := repo.Command(t.Context(), root, append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-c", "maintenance.auto=false"}, args...)...)
		cmd.Env = append(cmd.Environ(), "GIT_COMMITTER_DATE=1700000000 +0000")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	datedCommit := func(message string) { dated("commit", "-q", "--allow-empty", "-m", message) }
	dated("checkout", "-q", "-b", "a")
	datedCommit("a1")
	dated("checkout", "-q", "-b", "b", "main")
	datedCommit("b1")
	dated("checkout", "-q", "a")
	dated("merge", "-q", "--no-ff", "-m", "a merges b", "b")
	dated("checkout", "-q", "b")
	dated("merge", "-q", "--no-ff", "-m", "b merges a", "a~1")
	datedCommit("b2")
	dated("checkout", "-q", "main")
	datedCommit("m1")

	o := newObjects(t.Context(), root, "")
	defer o.close()
	for _, pair := range [][2]string{{"a", "b"}, {"a", "main"}, {"main", "b"}, {"a", "a~1"}} {
		x, y := git(t, root, "rev-parse", pair[0]), git(t, root, "rev-parse", pair[1])
		got, err := o.mergeBases(x, y)
		if err != nil {
			t.Fatal(err)
		}
		want := strings.Fields(git(t, root, "merge-base", "--all", x, y))
		if pair[0] == "a" && pair[1] == "b" && len(want) != 2 {
			t.Fatalf("the fixture should be criss-cross: %v", want)
		}
		slices.Sort(got)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Errorf("merge bases of %s and %s: got %v, want %v", pair[0], pair[1], got, want)
		}
	}
}

// TestCurrentViewCycle: a revert carried across merges can make older a
// cycle. Its states must not vanish, even while an unrelated branch that
// cannot be ordered against them stays current: nothing outside the cycle is
// newer, so each state in it is current, with a note.
func TestCurrentViewCycle(t *testing.T) {
	t.Parallel()
	root := repoFixture(t)
	rec := func(body string) { write(t, root, "grove/G-030.md", record("G-030", "work", "proposed", body)) }
	rec("X0\n")
	commit(t, root, "M0")
	git(t, root, "checkout", "-q", "-b", "e")
	rec("W\n")
	commit(t, root, "unrelated edit on e")
	git(t, root, "checkout", "-q", "main")
	rec("X\n")
	m1 := commit(t, root, "M1")
	rec("Y\n")
	commit(t, root, "M3")
	git(t, root, "checkout", "-q", "-b", "b")
	write(t, root, "other.txt", "x\n")
	commit(t, root, "unrelated on b")
	git(t, root, "checkout", "-q", "-b", "d", m1)
	rec("Z\n")
	commit(t, root, "D")
	git(t, root, "checkout", "-q", "-b", "c", "main")
	git(t, root, "merge", "-q", "--no-ff", "-m", "merge d", "-s", "ours", "d")
	rec("Z\n")
	commit(t, root, "resolve to Z")
	git(t, root, "checkout", "-q", "-b", "a", "d")
	rec("X\n")
	commit(t, root, "revert to X")
	git(t, root, "checkout", "-q", "main")

	// a is older than b (base M1 has X), b than c (base M3 has Y), c than a
	// (base D has Z). d is older than a, which changed d's record, and
	// outside the cycle. e meets every other at M0, holding neither's bytes.
	expectStanding(t, mustInspect(t, root, ""), "G-030",
		"branch a proposed current", "branch b proposed current", "branch c proposed current",
		"branch d proposed older", "branch e proposed current", "branch main proposed current", "checkout . proposed current",
		"some versions could not be ordered: each is older than another, through changes and reverts that merges carried across branches")
}

// TestCurrentViewCancelled: a read failing during the projection, as a
// cancelled load's does, must not panic; the caller discards the result.
func TestCurrentViewCancelled(t *testing.T) {
	t.Parallel()
	root := repoFixture(t)
	git(t, root, "checkout", "-q", "-b", "f")
	write(t, root, "grove/work/G-001-first.md", record("G-001", "work", "active", "x\n"))
	commit(t, root, "f")
	git(t, root, "checkout", "-q", "main")
	res := mustInspect(t, root, "")
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	o := newObjects(ctx, root, "")
	defer o.close()
	o.project(res)
}

// TestBaseOfUnborn: two observations on one unborn HEAD share no history, so
// the record was absent at their base rather than unreadable.
func TestBaseOfUnborn(t *testing.T) {
	t.Parallel()
	root := repoFixture(t)
	o := newObjects(t.Context(), root, "")
	defer o.close()
	zero := strings.Repeat("0", 40)
	base, err := o.baseOf(&node{commit: zero, live: true, content: "a"}, &node{commit: zero, live: true, content: "b"}, "G-001")
	if base != "" || err != nil {
		t.Errorf("got %q, %v", base, err)
	}
}
