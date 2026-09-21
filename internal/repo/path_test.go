package repo

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	full := append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-C", dir}, args...)
	if out, err := exec.Command("git", full...).CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// oddRepos builds, below one parent: a main checkout "new\nline" holding the
// directory "sub\nproject", its linked checkout "feature\tline ", and a
// checkout "sep" whose separate Git directory is named "git\tdir\nmid ". Git
// itself cannot reopen a separate Git directory whose name ends in a newline:
// its .git file format drops that terminator.
func oddRepos(t *testing.T) (parent, main, linked, sep, sepGit string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	parent, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	main, linked = filepath.Join(parent, "new\nline"), filepath.Join(parent, "feature\tline ")
	sep, sepGit = filepath.Join(parent, "sep"), filepath.Join(parent, "git\tdir\nmid ")
	if err := os.MkdirAll(filepath.Join(main, "sub\nproject"), 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, main, "init", "-q", "-b", "main")
	git(t, main, "commit", "-q", "--allow-empty", "-m", "init")
	git(t, main, "worktree", "add", "-q", "-b", "feature", linked)
	git(t, parent, "init", "-q", "-b", "main", "--separate-git-dir", sepGit, sep)
	return parent, main, linked, sep, sepGit
}

func TestLocateExactPaths(t *testing.T) {
	t.Parallel()
	_, main, linked, sep, sepGit := oddRepos(t)
	if err := os.MkdirAll(filepath.Join(linked, "sub\nproject"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct{ root, common, prefix string }{
		{main, filepath.Join(main, ".git"), ""},
		{filepath.Join(main, "sub\nproject"), filepath.Join(main, ".git"), "sub\nproject/"},
		{linked, filepath.Join(main, ".git"), ""},
		{filepath.Join(linked, "sub\nproject"), filepath.Join(main, ".git"), "sub\nproject/"},
		{sep, sepGit, ""},
	} {
		common, prefix, err := Locate(c.root)
		if err != nil || common != c.common || prefix != c.prefix {
			t.Fatalf("Locate(%q) = %q, %q, %v; want %q, %q", c.root, common, prefix, err, c.common, c.prefix)
		}
	}
	if got, err := GitPath(linked, "--git-dir"); err != nil || !strings.HasPrefix(got, filepath.Join(main, ".git", "worktrees")+"/") {
		t.Fatalf("linked Git directory: %q %v", got, err)
	}
}

// Coordination state is created under the real common directory only; the
// demonstrated stray sibling (".../new/grove") must not appear.
func TestCommonDirCreatesStateOnlyUnderTheRealCommonDirectory(t *testing.T) {
	t.Parallel()
	parent, main, linked, sep, sepGit := oddRepos(t)
	for _, root := range []string{main, linked, sep} {
		if _, _, err := CommonDir(root); err != nil {
			t.Fatal(err)
		}
	}
	var state []string
	filepath.WalkDir(parent, func(path string, entry os.DirEntry, err error) error {
		if err == nil && entry.IsDir() && entry.Name() == "grove" {
			state = append(state, path)
		}
		return nil
	})
	if want := []string{filepath.Join(sepGit, "grove"), filepath.Join(main, ".git", "grove")}; !reflect.DeepEqual(state, want) {
		t.Fatalf("coordination folders: %q, want %q", state, want)
	}
	entries, _ := os.ReadDir(parent)
	if len(entries) != 4 {
		t.Fatalf("nothing may appear beside the fixtures: %q", entries)
	}
}

func TestParseWorktrees(t *testing.T) {
	t.Parallel()
	out := "worktree /r/new\nline\x00HEAD 1111\x00branch refs/heads/main\x00\x00" +
		"worktree /r/tab\there \\ \"q\" \x00HEAD 2222\x00detached\x00locked why\nnot\x00\x00" +
		"worktree /r/gone\x00HEAD 3333\x00branch refs/heads/gone\x00prunable gitdir file points to non-existent location\x00\x00" +
		"worktree /r/bare.git\x00bare\x00\x00" +
		"worktree /r/bare-reason\x00HEAD 4444\x00branch refs/heads/x\x00prunable\x00\x00"
	want := []Worktree{
		{Path: "/r/new\nline", Head: "1111", Branch: "refs/heads/main"},
		{Path: "/r/tab\there \\ \"q\" ", Head: "2222"},
		{Path: "/r/gone", Head: "3333", Branch: "refs/heads/gone", Prunable: "gitdir file points to non-existent location"},
		{Path: "/r/bare.git", Bare: true},
		{Path: "/r/bare-reason", Head: "4444", Branch: "refs/heads/x", Prunable: "no reason given"},
	}
	if got := parseWorktrees(out); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v\nwant %+v", got, want)
	}
}

func TestWorktreesExactPaths(t *testing.T) {
	t.Parallel()
	_, main, linked, _, _ := oddRepos(t)
	for _, root := range []string{main, linked} {
		got, err := Worktrees(root)
		if err != nil || len(got) != 2 || got[0].Path != main || got[1].Path != linked || got[1].Branch != "refs/heads/feature" {
			t.Fatalf("from %q: %+v %v", root, got, err)
		}
	}
}

// Several paths come from one process only when none holds a newline; the
// answers are the same either way.
func TestGitPathsMatchSingleAnswers(t *testing.T) {
	t.Parallel()
	_, main, linked, sep, _ := oddRepos(t)
	for _, dir := range []string{main, sep, filepath.Join(main, "sub\nproject"), linked, t.TempDir()} {
		options := []string{"--git-dir", "--show-prefix", "--git-common-dir"}
		got, err := GitPathsContext(context.Background(), dir, options...)
		for i, option := range options {
			want, wantErr := GitPath(dir, option)
			if (err == nil) != (wantErr == nil) || err == nil && got[i] != want {
				t.Fatalf("%q %s: got %q, %v; want %q, %v", dir, option, got, err, want, wantErr)
			}
		}
	}
}
