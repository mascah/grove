package versions

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	good := map[string]Selection{
		"committed:refs/heads/main@0123456789ab:W-001@0123456789ab:0123456789abcdef":         {Kind: "committed", Ref: "refs/heads/main", Commit: "0123456789ab", ID: "W-001", Revision: "0123456789ab", Binding: "0123456789abcdef"},
		"live:.:refs/heads/a@b@0123456789ab:Q-1000@0123456789ab:0123456789abcdef":            {Kind: "live", Locator: ".", Ref: "refs/heads/a@b", Commit: "0123456789ab", ID: "Q-1000", Revision: "0123456789ab", Binding: "0123456789abcdef"},
		"live:odd-name-with-space:detached@0123456789ab:D-001@0123456789ab:0123456789abcdef": {Kind: "live", Locator: "odd-name-with-space", Commit: "0123456789ab", ID: "D-001", Revision: "0123456789ab", Binding: "0123456789abcdef"},
	}
	for selector, want := range good {
		if got, err := Parse(selector); err != nil || got != want {
			t.Fatalf("%s: %+v %v", selector, got, err)
		}
	}
	for _, bad := range []string{
		"", "W-001", "committed:refs/heads/main@0123456789ab:W-001@0123456789ab", "live:refs/heads/main@0123456789ab:W-001@0123456789ab:0123456789abcdef",
		"committed:main@0123456789ab:W-001@0123456789ab:0123456789abcdef", "committed:refs/heads/@0123456789ab:W-001@0123456789ab:0123456789abcdef",
		"committed:detached@0123456789ab:W-001@0123456789ab:0123456789abcdef", "committed:refs/heads/main@0123456789AB:W-001@0123456789ab:0123456789abcdef",
		"committed:refs/heads/main@0123456789ab:W-01@0123456789ab:0123456789abcdef", "committed:refs/heads/main@0123456789ab:W-001@0123456789ab:0123456789abcde",
		"live::refs/heads/main@0123456789ab:W-001@0123456789ab:0123456789abcdef", "live:a/b:refs/heads/main@0123456789ab:W-001@0123456789ab:0123456789abcdef",
		"other:refs/heads/main@0123456789ab:W-001@0123456789ab:0123456789abcdef", "committed:refs/heads/main@0123456789ab:W-001:0123456789abcdef",
	} {
		if _, err := Parse(bad); err == nil {
			t.Fatalf("%q must not parse", bad)
		}
	}
}

// selectorFor returns the current selector of id in the named source.
func selectorFor(t *testing.T, root, id, kind, where string) string {
	t.Helper()
	return find(t, group(t, mustInspect(t, root, id), id), kind, where).Selector
}

func resolve(t *testing.T, root, selector string) *Workspace {
	t.Helper()
	w, err := Resolve(root, selector)
	if err != nil {
		t.Fatalf("%s: %v", selector, err)
	}
	return w
}

func refuse(t *testing.T, root, selector, want string) {
	t.Helper()
	w, err := Resolve(root, selector)
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("%s: expected %q, got %+v %v", selector, want, w, err)
	}
}

func treeHashes(t *testing.T, dir string) map[string][32]byte {
	t.Helper()
	result := map[string][32]byte{}
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type().IsRegular() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			result[path] = sha256.Sum256(data)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}

// nestedFixture puts the project below the repository root and adds a feature
// worktree whose W-001 is active.
func nestedFixture(t *testing.T) (root, wt string) {
	t.Helper()
	root = repoFixture(t)
	git(t, root, "rm", "-q", "-r", "grove.yaml", "grove")
	write(t, root, "sub/grove.yaml", config)
	write(t, root, "sub/grove/work/W-001-first.md", record("W-001", "work", "proposed", "Main.\n"))
	commit(t, root, "nested")
	wt = addWorktree(t, root, "feature", "", "-b", "feature")
	write(t, wt, "sub/grove/work/W-001-first.md", record("W-001", "work", "active", "Feature.\n"))
	commit(t, wt, "feature")
	return root, wt
}

func TestWorkspaceLive(t *testing.T) {
	root, wt := nestedFixture(t)
	project := filepath.Join(root, "sub")
	write(t, root, "sub/notes.txt", "dirty in main\n")
	write(t, wt, "sub/staged.txt", "staged in feature\n")
	git(t, wt, "add", "sub/staged.txt")
	write(t, wt, "sub/unstaged.txt", "unstaged in feature\n")
	before := map[string]map[string][32]byte{root: treeHashes(t, root), wt: treeHashes(t, wt)}
	branch := git(t, root, "rev-parse", "--abbrev-ref", "HEAD")

	selector := selectorFor(t, project, "W-001", "live", "feature")
	w := resolve(t, project, selector)
	want := &Workspace{
		Checkout: wt, Project: filepath.Join(wt, "sub"), Record: filepath.Join(wt, "sub", "grove", "work", "W-001-first.md"),
		Ref: "refs/heads/feature", Head: git(t, wt, "rev-parse", "HEAD"), Selector: selector,
	}
	data, err := os.ReadFile(want.Record)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	want.Revision = "sha256:" + hex.EncodeToString(sum[:])
	if !reflect.DeepEqual(w, want) {
		t.Fatalf("got %+v\nwant %+v", w, want)
	}
	if got := resolve(t, project, selectorFor(t, project, "W-001", "live", ".")); got.Checkout != root || got.Project != project || got.Ref != "refs/heads/main" {
		t.Fatalf("main's own live version resolves to main: %+v", got)
	}
	if !reflect.DeepEqual(before[root], treeHashes(t, root)) || !reflect.DeepEqual(before[wt], treeHashes(t, wt)) || git(t, root, "rev-parse", "--abbrev-ref", "HEAD") != branch {
		t.Fatal("resolution must leave both checkouts, their dirty files, and Git state unchanged")
	}
	if git(t, wt, "status", "--porcelain") != "A  sub/staged.txt\n?? sub/unstaged.txt" {
		t.Fatalf("staged and unstaged files must survive: %q", git(t, wt, "status", "--porcelain"))
	}
}

func TestWorkspaceCommitted(t *testing.T) {
	root, wt := nestedFixture(t)
	project := filepath.Join(root, "sub")
	committed := selectorFor(t, project, "W-001", "committed", "refs/heads/feature")
	w := resolve(t, project, committed)
	live := selectorFor(t, project, "W-001", "live", "feature")
	if w.Checkout != wt || w.Selector != live || w.Ref != "refs/heads/feature" {
		t.Fatalf("a matching committed selection resolves to its checkout and reports the live selector: %+v", w)
	}
	write(t, wt, "sub/grove/work/W-001-first.md", record("W-001", "work", "done", "Feature edited.\n"))
	refuse(t, project, committed, "differs from the committed version selected (modified at grove/work/W-001-first.md); run versions and select the live observation")
	refuse(t, project, live, "W-001 changed since it was selected")
	fresh := selectorFor(t, project, "W-001", "live", "feature")
	if got := resolve(t, project, fresh); got.Checkout != wt || got.Revision == w.Revision {
		t.Fatalf("the fresh live selection resolves: %+v", got)
	}
	if err := os.Remove(filepath.Join(wt, "sub/grove/work/W-001-first.md")); err != nil {
		t.Fatal(err)
	}
	refuse(t, project, committed, "W-001 was deleted from the live files of worktree feature")
	refuse(t, project, fresh, "W-001 was deleted from the live files of worktree feature")
	git(t, wt, "checkout", "-q", "--", "sub/grove/work/W-001-first.md")
	if got := resolve(t, project, committed); got.Checkout != wt {
		t.Fatalf("restored bytes resolve again: %+v", got)
	}
	if got := resolve(t, project, selectorFor(t, project, "W-001", "committed", "refs/heads/main")); got.Checkout != root || got.Project != project {
		t.Fatalf("main's committed version resolves to the main checkout: %+v", got)
	}
	// A valid but different live grove.yaml refuses the committed route even
	// though the record bytes and path still match.
	write(t, wt, "sub/grove.yaml", "schema_version: 1\nrecords: grove # local note\n")
	refuse(t, project, committed, "the live grove.yaml in worktree feature differs from the committed configuration selected; run versions and select the live observation")
	git(t, wt, "checkout", "-q", "--", "sub/grove.yaml")
	// A checkout with no project at the prefix says so instead of reporting an
	// empty list of problems.
	if err := os.Rename(filepath.Join(wt, "sub/grove.yaml"), filepath.Join(wt, "grove.yaml.aside")); err != nil {
		t.Fatal(err)
	}
	refuse(t, project, committed, "worktree feature has no grove.yaml at the selected project location; the project is absent there")
	refuse(t, project, live, "worktree feature has no grove.yaml at the selected project location; the project is absent there")
}

func TestWorkspaceStale(t *testing.T) {
	root, wt := nestedFixture(t)
	project := filepath.Join(root, "sub")
	live := selectorFor(t, project, "W-001", "live", "feature")
	committed := selectorFor(t, project, "W-001", "committed", "refs/heads/feature")

	write(t, wt, "sub/grove/work/W-002-second.md", record("W-002", "work", "proposed", "x\n"))
	commit(t, wt, "advance")
	refuse(t, project, live, "worktree feature's HEAD moved from")
	refuse(t, project, committed, "branch refs/heads/feature moved from")

	live = selectorFor(t, project, "W-001", "live", "feature")
	git(t, wt, "checkout", "-q", "--detach")
	refuse(t, project, live, "worktree feature is now detached at")
	detached := selectorFor(t, project, "W-001", "live", "feature")
	if w := resolve(t, project, detached); w.Ref != "" || w.Checkout != wt {
		t.Fatalf("a detached live selection resolves: %+v", w)
	}
	git(t, wt, "checkout", "-q", "feature")
	refuse(t, project, detached, "worktree feature is now refs/heads/feature at")

	live = selectorFor(t, project, "W-001", "live", "feature")
	if err := os.Rename(filepath.Join(wt, "sub/grove/work/W-001-first.md"), filepath.Join(wt, "sub/grove/work/W-001-moved.md")); err != nil {
		t.Fatal(err)
	}
	refuse(t, project, live, "its worktree path, configuration, record path, or project location changed")
	git(t, wt, "checkout", "-q", "--", "sub/grove")
	if err := os.Remove(filepath.Join(wt, "sub/grove/work/W-001-moved.md")); err != nil {
		t.Fatal(err)
	}
	if w := resolve(t, project, live); w.Checkout != wt {
		t.Fatalf("restoring the path restores the selection: %+v", w)
	}

	write(t, wt, "sub/grove.yaml", "schema_version: 1\nrecords: grove # moved\n")
	refuse(t, project, live, "configuration, record path, or project location changed")
	git(t, wt, "checkout", "-q", "--", "sub/grove.yaml")

	moved := filepath.Join(filepath.Dir(root), "feature-moved")
	git(t, root, "worktree", "move", wt, moved)
	refuse(t, project, live, "its worktree path, configuration, record path, or project location changed")
	if w := resolve(t, project, selectorFor(t, project, "W-001", "live", "feature")); w.Checkout != moved {
		t.Fatalf("a fresh selection follows the moved worktree: %+v", w)
	}
	git(t, root, "worktree", "move", moved, wt)

	live = selectorFor(t, project, "W-001", "live", "feature")
	git(t, root, "worktree", "remove", "--force", wt)
	refuse(t, project, live, "worktree feature is no longer registered in this repository")
	refuse(t, project, selectorFor(t, project, "W-001", "committed", "refs/heads/feature"), "no registered worktree has refs/heads/feature checked out; this command does not create one")
	git(t, root, "branch", "-q", "-D", "feature")
	refuse(t, project, committed, "branch refs/heads/feature no longer exists")

	write(t, root, "sub/grove/work/W-002-second.md", "---\nid: \"W-002\"\ntype: work\ntitle: T\nstatus: bogus\n---\n")
	refuse(t, project, selectorFor(t, project, "W-001", "committed", "refs/heads/main"), "worktree . is not a valid source:")
	if err := os.Remove(filepath.Join(root, "sub/grove/work/W-002-second.md")); err != nil {
		t.Fatal(err)
	}
	refuse(t, project, "committed:refs/heads/main@0123456789ab:W-001@0123456789ab:0123456789abcdef", "branch refs/heads/main moved from 0123456789ab")
	refuse(t, project, "live:.:refs/heads/main@0123456789ab:W-009@0123456789ab:0123456789abcdef", "W-009 is not present in worktree .")
	refuse(t, project, "nonsense", "selector must look like")
}

func TestWorkspaceMissingAndAmbiguous(t *testing.T) {
	root, wt := nestedFixture(t)
	project := filepath.Join(root, "sub")
	temp := addWorktree(t, root, "temp", "", "-b", "orphan")
	write(t, temp, "sub/grove/work/W-001-first.md", record("W-001", "work", "done", "Orphan.\n"))
	commit(t, temp, "orphan")
	git(t, root, "worktree", "remove", "--force", temp)
	count := git(t, root, "worktree", "list", "--porcelain")
	refuse(t, project, selectorFor(t, project, "W-001", "committed", "refs/heads/orphan"), "no registered worktree has refs/heads/orphan checked out; this command does not create one")
	if git(t, root, "worktree", "list", "--porcelain") != count {
		t.Fatal("a missing checkout must not be created")
	}
	second := addWorktree(t, root, "second", "feature", "-f")
	refuse(t, project, selectorFor(t, project, "W-001", "committed", "refs/heads/feature"), "2 worktrees have refs/heads/feature checked out (feature, second); select one of their live versions instead")
	if w := resolve(t, project, selectorFor(t, project, "W-001", "live", "second")); w.Checkout != second || w.Ref != "refs/heads/feature" {
		t.Fatalf("an explicit live selection disambiguates: %+v", w)
	}
	if w := resolve(t, project, selectorFor(t, project, "W-001", "live", "feature")); w.Checkout != wt {
		t.Fatalf("the other checkout resolves too: %+v", w)
	}
	// A prunable duplicate is not a checkout anyone can be sent to, so the
	// remaining checkout is unambiguous; the prunable entry still marks the
	// inspection incomplete on its own.
	if err := os.RemoveAll(second); err != nil {
		t.Fatal(err)
	}
	if w := resolve(t, project, selectorFor(t, project, "W-001", "committed", "refs/heads/feature")); w.Checkout != wt {
		t.Fatalf("expected the surviving checkout: %+v", w)
	}
}
