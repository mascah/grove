package versions

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// Each mutation lands after the inspection and before the final check of the
// target about to be returned.
func TestResolveFinalCheck(t *testing.T) {
	record1 := "sub/grove/work/W-001-first.md"
	cases := []struct {
		name, kind string // kind of selection: live or committed
		mutate     func(t *testing.T, root, wt string)
		want       string
	}{
		{"record bytes", "live", func(t *testing.T, root, wt string) {
			write(t, wt, record1, record("W-001", "work", "done", "Edited.\n"))
		}, "changed while its workspace was being resolved"},
		{"configuration bytes", "committed", func(t *testing.T, root, wt string) {
			write(t, wt, "sub/grove.yaml", config+"# note\n")
		}, "changed while its workspace was being resolved"},
		{"record path", "live", func(t *testing.T, root, wt string) {
			must(t, os.Rename(filepath.Join(wt, record1), filepath.Join(wt, "sub/grove/work/W-001-moved.md")))
		}, "changed while its workspace was being resolved"},
		{"record deleted", "committed", func(t *testing.T, root, wt string) {
			must(t, os.Remove(filepath.Join(wt, record1)))
		}, "changed while its workspace was being resolved"},
		{"source made invalid", "live", func(t *testing.T, root, wt string) {
			write(t, wt, "sub/grove/work/W-002-second.md", "---\nid: \"W-002\"\ntype: work\ntitle: T\nstatus: bogus\n---\n")
		}, "stopped being a valid source"},
		{"project removed", "live", func(t *testing.T, root, wt string) {
			must(t, os.RemoveAll(filepath.Join(wt, "sub")))
		}, "lost its project"},
		{"checkout deleted", "live", func(t *testing.T, root, wt string) {
			must(t, os.RemoveAll(wt))
		}, "worktree is prunable"},
		{"checkout deleted on the committed route", "committed", func(t *testing.T, root, wt string) {
			must(t, os.RemoveAll(wt))
		}, "worktree is prunable"},
		{"project made foreign", "live", func(t *testing.T, root, wt string) {
			must(t, os.RemoveAll(filepath.Join(wt, "sub")))
			foreignRepo(t, filepath.Join(wt, "sub"), "")
		}, "belongs to another repository or worktree"},
		{"project symlinked", "committed", func(t *testing.T, root, wt string) {
			must(t, os.Rename(filepath.Join(wt, "sub"), filepath.Join(wt, "aside")))
			must(t, os.Symlink("aside", filepath.Join(wt, "sub")))
		}, "is a symlink"},
		{"registration moved", "live", func(t *testing.T, root, wt string) {
			git(t, root, "worktree", "move", wt, wt+"-moved")
		}, "was removed or moved"},
		{"registration removed", "committed", func(t *testing.T, root, wt string) {
			git(t, root, "worktree", "remove", "--force", wt)
		}, "was removed or moved"},
		{"HEAD advanced", "live", func(t *testing.T, root, wt string) {
			commit(t, wt, "empty advance")
		}, "changed while its workspace was being resolved"},
		{"branch advanced", "committed", func(t *testing.T, root, wt string) {
			commit(t, wt, "empty advance")
		}, "changed while its workspace was being resolved"},
		{"detached", "live", func(t *testing.T, root, wt string) {
			git(t, wt, "checkout", "-q", "--detach")
		}, "changed while its workspace was being resolved"},
		{"branch tip moved under a detached checkout", "committed", func(t *testing.T, root, wt string) {
			git(t, wt, "checkout", "-q", "--detach")
			git(t, root, "branch", "-q", "-f", "feature", "main")
		}, "changed while its workspace was being resolved"},
		{"second checkout of the branch", "committed", func(t *testing.T, root, wt string) {
			addWorktree(t, root, "second", "feature", "-f")
		}, "worktree second also checked out refs/heads/feature"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root, wt := nestedFixture(t)
			project := filepath.Join(root, "sub")
			where := "feature"
			if c.kind == "committed" {
				where = "refs/heads/feature"
			}
			selected := selectorFor(t, project, "W-001", c.kind, where)
			var after map[string][32]byte
			w, err := resolveWith(t.Context(), project, selected, func() {
				c.mutate(t, root, wt)
				after = treeHashes(t, filepath.Dir(root))
			})
			if err == nil || w != nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("expected %q, got %+v %v", c.want, w, err)
			}
			if !reflect.DeepEqual(after, treeHashes(t, filepath.Dir(root))) {
				t.Fatal("the final check and refusal must not write anything")
			}
		})
	}
}

// Unchanged targets still resolve through the final check, a second checkout
// does not disturb an explicit live selection, unrelated invalid sources do
// not block it, and nothing in any checkout or Git directory is written.
func TestResolveFinalCheckAdmits(t *testing.T) {
	root, wt := nestedFixture(t)
	project := filepath.Join(root, "sub")
	bad := addWorktree(t, root, "bad", "", "-b", "bad")
	write(t, bad, "sub/grove/work/W-002-second.md", "---\nid: \"W-002\"\ntype: work\ntitle: [unclosed\nstatus: proposed\n---\n")
	write(t, wt, "sub/dirty.txt", "unrelated\n")
	live := selectorFor(t, project, "W-001", "live", "feature")
	committed := selectorFor(t, project, "W-001", "committed", "refs/heads/feature")
	before := treeHashes(t, filepath.Dir(root))
	for _, selected := range []string{live, committed} {
		w, err := resolveWith(t.Context(), project, selected, func() { write(t, wt, "sub/dirty.txt", "unrelated\n") })
		if err != nil || w.Checkout != wt || w.Selector != live {
			t.Fatalf("%s: %+v %v", selected, w, err)
		}
	}
	if !reflect.DeepEqual(before, treeHashes(t, filepath.Dir(root))) {
		t.Fatal("resolution must not write anything")
	}
	w, err := resolveWith(t.Context(), project, live, func() { addWorktree(t, root, "second", "feature", "-f") })
	if err != nil || w.Checkout != wt {
		t.Fatalf("an explicit live selection survives a new duplicate checkout: %+v %v", w, err)
	}
}
