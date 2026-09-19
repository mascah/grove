package versions

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// deepFixture puts the project at outer/sub and adds a feature worktree.
func deepFixture(t *testing.T) (root, wt string) {
	t.Helper()
	root = repoFixture(t)
	git(t, root, "rm", "-q", "-r", "grove.yaml", "grove")
	write(t, root, "outer/sub/grove.yaml", config)
	write(t, root, "outer/sub/grove/work/W-001-first.md", record("W-001", "work", "proposed", "Main.\n"))
	commit(t, root, "nested project")
	return root, addWorktree(t, root, "feature", "", "-b", "feature")
}

// foreignRepo is a separate repository whose valid project is at below.
func foreignRepo(t *testing.T, dir, below string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, dir, "init", "-q", "-b", "main")
	write(t, dir, filepath.Join(below, "grove.yaml"), config)
	write(t, dir, filepath.Join(below, "grove/work/W-001-first.md"), record("W-001", "work", "done", "Foreign.\n"))
	commit(t, dir, "foreign")
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func TestInspectProjectLocation(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(t *testing.T, wt string)
		want   string // diagnostic; "" means the project is absent there
	}{
		{"external symlink", func(t *testing.T, wt string) {
			outside := filepath.Join(filepath.Dir(wt), "outside")
			foreignRepo(t, outside, "")
			must(t, os.RemoveAll(filepath.Join(wt, "outer/sub")))
			must(t, os.Symlink(outside, filepath.Join(wt, "outer/sub")))
		}, "is a symlink"},
		{"middle symlink", func(t *testing.T, wt string) {
			outside := filepath.Join(filepath.Dir(wt), "outside")
			foreignRepo(t, outside, "sub")
			must(t, os.RemoveAll(filepath.Join(wt, "outer")))
			must(t, os.Symlink(outside, filepath.Join(wt, "outer")))
		}, "is a symlink"},
		{"internal symlink", func(t *testing.T, wt string) {
			must(t, os.Rename(filepath.Join(wt, "outer/sub"), filepath.Join(wt, "elsewhere")))
			must(t, os.Symlink("../elsewhere", filepath.Join(wt, "outer/sub")))
		}, "is a symlink"},
		{"dangling symlink", func(t *testing.T, wt string) {
			must(t, os.RemoveAll(filepath.Join(wt, "outer/sub")))
			must(t, os.Symlink("nowhere", filepath.Join(wt, "outer/sub")))
		}, "is a symlink"},
		{"regular file", func(t *testing.T, wt string) {
			must(t, os.RemoveAll(filepath.Join(wt, "outer/sub")))
			write(t, wt, "outer/sub", "not a directory\n")
		}, "is not a directory"},
		{"nested repository", func(t *testing.T, wt string) {
			must(t, os.RemoveAll(filepath.Join(wt, "outer/sub")))
			foreignRepo(t, filepath.Join(wt, "outer/sub"), "")
		}, "belongs to another repository or worktree"},
		{"nested repository at the record folder", func(t *testing.T, wt string) {
			must(t, os.RemoveAll(filepath.Join(wt, "outer/sub/grove")))
			write(t, wt, "outer/sub/grove/work/W-001-first.md", record("W-001", "work", "done", "Foreign.\n"))
			git(t, filepath.Join(wt, "outer/sub/grove"), "init", "-q")
		}, "belongs to another repository or worktree"},
		{"missing directory", func(t *testing.T, wt string) {
			must(t, os.RemoveAll(filepath.Join(wt, "outer/sub")))
		}, ""},
		{"missing parent", func(t *testing.T, wt string) {
			must(t, os.RemoveAll(filepath.Join(wt, "outer")))
		}, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root, wt := deepFixture(t)
			project := filepath.Join(root, "outer/sub")
			earlier := selectorFor(t, project, "W-001", "live", "feature")
			c.mutate(t, wt)
			before := treeHashes(t, filepath.Dir(root))

			res := mustInspect(t, project, "W-001")
			s := source(t, res, "live", "feature")
			if s.Valid || res.Complete != (c.want == "") || (c.want == "") != (len(s.Diagnostics) == 0 && !s.Present) {
				t.Fatalf("feature source: complete=%v %+v", res.Complete, s)
			}
			if c.want != "" && !strings.Contains(strings.Join(s.Diagnostics, "\n"), c.want) {
				t.Fatalf("expected %q in %v", c.want, s.Diagnostics)
			}
			g := group(t, res, "W-001")
			for _, v := range g.Versions {
				if v.Source == s {
					t.Fatalf("the feature checkout must contribute no version: %+v", v)
				}
			}
			if main := find(t, g, "live", "."); main.Record.Status != "proposed" || main.Selector == "" {
				t.Fatalf("main's healthy records stay visible: %+v", main)
			}
			want := "worktree feature is not a valid source"
			if c.want == "" {
				want = "the project is absent there"
			}
			refuse(t, project, earlier, want)
			if !reflect.DeepEqual(before, treeHashes(t, filepath.Dir(root))) {
				t.Fatal("inspection and refusal must not write anything")
			}
		})
	}
}

// A directory that merely sits where a locked worktree used to be, inside
// another checkout, is not that worktree.
func TestInspectWorktreeReplacedByPlainDirectory(t *testing.T) {
	root := repoFixture(t)
	inner := filepath.Join(root, "inner")
	git(t, root, "worktree", "add", "-q", "--lock", "-b", "inner", inner)
	must(t, os.RemoveAll(inner))
	must(t, os.Mkdir(inner, 0o755))
	res := mustInspect(t, root, "W-001")
	for _, s := range res.Sources {
		if s.Worktree == inner && (s.Valid || s.Locator != "" || len(s.Diagnostics) == 0) {
			t.Fatalf("plain directory admitted as a checkout: %+v", s)
		}
	}
	if res.Complete || !source(t, res, "live", ".").Valid {
		t.Fatalf("incomplete, with main intact: %s", dump(res))
	}
}

func TestInspectPrunableDuringRead(t *testing.T) {
	root := repoFixture(t)
	wt := addWorktree(t, root, "feature", "", "-b", "feature")
	res, err := inspect(root, "W-001", func() {
		must(t, os.RemoveAll(wt))
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Complete || source(t, res, "live", "feature").Valid {
		t.Fatal("deleted checkout remained selectable")
	}
	for _, v := range group(t, res, "W-001").Versions {
		if v.Source.Worktree == wt {
			t.Fatalf("a vanished checkout must not contribute: %+v", v)
		}
	}
}

func TestInspectForeignDuringRead(t *testing.T) {
	root, wt := deepFixture(t)
	project := filepath.Join(root, "outer/sub")
	res, err := inspect(project, "W-001", func() {
		must(t, os.Rename(filepath.Join(wt, "outer/sub"), filepath.Join(wt, "elsewhere")))
		must(t, os.Symlink("../elsewhere", filepath.Join(wt, "outer/sub")))
	})
	if err != nil {
		t.Fatal(err)
	}
	if s := source(t, res, "live", "feature"); res.Complete || s.Valid || !strings.Contains(strings.Join(s.Diagnostics, "\n"), "is a symlink") {
		t.Fatalf("a project swapped for a symlink during the read: %+v", s)
	}
}

// W-008: newlines, tabs, and trailing blanks in the main checkout, the project
// prefix, and a linked checkout keep their exact identity through inspection
// and both kinds of resolution, and reading creates no coordination state.
func TestOddGitPathsRoundTrip(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	parent, err := filepath.EvalSymlinks(t.TempDir())
	must(t, err)
	root, wt := filepath.Join(parent, "new\nline"), filepath.Join(parent, "feature\tline ")
	must(t, os.Mkdir(root, 0o755))
	git(t, root, "init", "-q", "-b", "main")
	write(t, root, "sub\nproject/grove.yaml", config)
	write(t, root, "sub\nproject/grove/work/W-001-first.md", record("W-001", "work", "proposed", "Main.\n"))
	commit(t, root, "init")
	git(t, root, "worktree", "add", "-q", "-b", "feature", wt)
	write(t, wt, "sub\nproject/grove/work/W-001-first.md", record("W-001", "work", "active", "Feature.\n"))
	commit(t, wt, "feature")
	project := filepath.Join(root, "sub\nproject")
	before := treeHashes(t, parent)

	res := mustInspect(t, project, "W-001")
	if !res.Complete || res.Repository != filepath.Join(root, ".git") || res.Prefix != "sub\nproject/" {
		t.Fatalf("repository %q prefix %q: %s", res.Repository, res.Prefix, dump(res))
	}
	var linked *Source
	for _, s := range res.Sources {
		if s.Kind == "live" && s.Worktree == wt {
			linked = s
		}
	}
	if linked == nil || !linked.Valid || linked.GitDir != filepath.Join(root, ".git", "worktrees", linked.Locator) {
		t.Fatalf("linked source: %+v", linked)
	}
	g := group(t, res, "W-001")
	for _, c := range []struct{ kind, where, checkout string }{
		{"live", ".", root}, {"live", linked.Locator, wt},
		{"committed", "refs/heads/main", root}, {"committed", "refs/heads/feature", wt},
	} {
		w := resolve(t, project, find(t, g, c.kind, c.where).Selector)
		if w.Checkout != c.checkout || w.Project != filepath.Join(c.checkout, "sub\nproject") || w.Record != filepath.Join(w.Project, "grove/work/W-001-first.md") {
			t.Fatalf("%s %s: %+v", c.kind, c.where, w)
		}
	}
	// The same picture from the linked checkout.
	if again := mustInspect(t, filepath.Join(wt, "sub\nproject"), "W-001"); !reflect.DeepEqual(selectorsOf(res), selectorsOf(again)) {
		t.Fatal("selectors depend on the repository, not the invoking checkout")
	}
	if !reflect.DeepEqual(before, treeHashes(t, parent)) {
		t.Fatal("reads must not write anything below the fixtures' parent")
	}
	if entries, _ := os.ReadDir(parent); len(entries) != 2 {
		t.Fatalf("nothing may appear beside the checkouts: %q", entries)
	}
}

func TestInspectConfigurationRemovedDuringRead(t *testing.T) {
	root, wt := deepFixture(t)
	project := filepath.Join(root, "outer/sub")
	res, err := inspect(project, "W-001", func() {
		must(t, os.Remove(filepath.Join(wt, "outer/sub/grove.yaml")))
	})
	must(t, err)
	if s := source(t, res, "live", "feature"); res.Complete || s.Valid || !strings.Contains(strings.Join(s.Diagnostics, "\n"), "grove.yaml changed or disappeared") {
		t.Fatalf("a project whose configuration vanished during the read: %+v", s)
	}
}
