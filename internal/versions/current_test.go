package versions

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
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
	if why := find(t, group(t, res, "G-010"), "committed", "refs/heads/stale").Older; why != "branch main changed it since their common history" {
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

// TestCurrentViewUnorderedPair shows ambiguity explicitly: two branches whose
// common commit holds a project that does not validate cannot be ordered, so
// both stay current and a note names them.
func TestCurrentViewUnorderedPair(t *testing.T) {
	t.Parallel()
	root := repoFixture(t)
	write(t, root, "grove/G-020.md", "---\nid: [broken\n---\n")
	write(t, root, "grove/G-021.md", record("G-021", "work", "proposed", "Start.\n"))
	base := commit(t, root, "base does not validate")
	write(t, root, "grove/G-020.md", record("G-020", "work", "proposed", "Fixed.\n"))
	commit(t, root, "fix on main")
	git(t, root, "checkout", "-q", "-b", "other", base)
	write(t, root, "grove/G-020.md", record("G-020", "work", "proposed", "Fixed.\n"))
	write(t, root, "grove/G-021.md", record("G-021", "work", "active", "Start.\n"))
	commit(t, root, "fix and change on other")
	git(t, root, "checkout", "-q", "main")

	res := mustInspect(t, root, "G-021")
	want := fmt.Sprintf("branch main and branch other could not be ordered: their common commit %s holds a project that does not validate", base[:12])
	expectStanding(t, res, "G-021",
		"branch main proposed current", "branch other active current", "checkout . proposed current", want)
}

// TestMergeBases checks the walk against Git's own answers, including a
// criss-cross history with two bases and commits sharing one timestamp.
func TestMergeBases(t *testing.T) {
	t.Setenv("GIT_COMMITTER_DATE", "1700000000 +0000")
	root := repoFixture(t)
	git(t, root, "checkout", "-q", "-b", "a")
	commit(t, root, "a1")
	git(t, root, "checkout", "-q", "-b", "b", "main")
	commit(t, root, "b1")
	git(t, root, "checkout", "-q", "a")
	git(t, root, "merge", "-q", "--no-ff", "-m", "a merges b", "b")
	git(t, root, "checkout", "-q", "b")
	git(t, root, "merge", "-q", "--no-ff", "-m", "b merges a", "a~1")
	commit(t, root, "b2")
	git(t, root, "checkout", "-q", "main")
	commit(t, root, "m1")

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
