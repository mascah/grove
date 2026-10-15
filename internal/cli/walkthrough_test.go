package cli

import (
	"github.com/mascah/grove-next2/internal/doc"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// fixture creates the ledger fixture as a fresh repository, as `just fixture`
// does, and makes it the working directory.
func fixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.CopyFS(dir, os.DirFS("../../testdata/ledger")); err != nil {
		t.Fatal(err)
	}
	git(t, dir, "init", "-q", "-b", "main")
	commitAll(t, dir, "Ledger 0.3.0")
	t.Chdir(dir)
	return dir
}

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=Fixture Owner", "GIT_AUTHOR_EMAIL=fixture@example.com",
		"GIT_COMMITTER_NAME=Fixture Owner", "GIT_COMMITTER_EMAIL=fixture@example.com")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

func commitAll(t *testing.T, dir, msg string) {
	t.Helper()
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-q", "-m", msg)
}

// grove runs one command and returns its standard output, failing the test
// on a non-zero status.
func grove(t *testing.T, args ...string) string {
	t.Helper()
	var stdout, stderr strings.Builder
	if code := run(args, &stdout, &stderr); code != 0 {
		t.Fatalf("grove-next2 %v: status %d\n%s", args, code, stderr.String())
	}
	return stdout.String()
}

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestInit(t *testing.T) {
	dir := fixture(t)
	out := grove(t, "init")
	if !strings.Contains(out, "nothing is committed") {
		t.Errorf("init said %q", out)
	}
	if got := read(t, filepath.Join(dir, ".gitignore")); !strings.HasSuffix(got, "*.egg-info/\n.grove/\n") {
		t.Errorf(".gitignore is %q", got)
	}
	if !strings.Contains(read(t, filepath.Join(dir, "grove.yaml")), "target: main") {
		t.Error("grove.yaml does not name the target")
	}
	commitAll(t, dir, "Set up Grove")
}

// harness puts a stand-in for claude on PATH. It saves the directory it ran
// in and its arguments, then runs script there, as a session's edits.
func harness(t *testing.T, script string) (cwd, args func() string) {
	t.Helper()
	bin, out := t.TempDir(), t.TempDir()
	stub := "#!/bin/sh\npwd -P > " + out + "/cwd\nprintf '%s' \"$*\" > " + out + "/args\n" + script + "\n"
	if err := os.WriteFile(filepath.Join(bin, "claude"), []byte(stub), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	return func() string { return strings.TrimSpace(read(t, out+"/cwd")) }, func() string { return read(t, out+"/args") }
}

// workspace returns the only workspace's directory.
func workspace(t *testing.T, dir string) string {
	t.Helper()
	m, _ := filepath.Glob(filepath.Join(dir, ".grove", "worktrees", "*"))
	if len(m) != 1 {
		t.Fatalf("workspaces: %v", m)
	}
	return m[0]
}

func TestExplore(t *testing.T) {
	dir := fixture(t)
	grove(t, "init")
	commitAll(t, dir, "Set up Grove")
	cwd, args := harness(t, "")

	out := grove(t, "explore", "Allow editing older records")
	ws := workspace(t, dir)
	name := filepath.Base(ws)
	id := name[:12]
	if got, _ := filepath.EvalSymlinks(ws); cwd() != got {
		t.Errorf("session ran in %s, want %s", cwd(), got)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	want := []string{
		"Workspace .grove/worktrees/" + name + " on branch grove/" + name + ", from main@",
		"Draft: exploration " + id + " in .grove/worktrees/" + name + ", 1 file uncommitted.",
		"Continue: ",
	}
	if len(lines) != 3 || !strings.HasPrefix(lines[0], want[0]) || lines[1] != want[1] || !strings.HasPrefix(lines[2], want[2]) || !strings.HasSuffix(lines[2], " explore "+id+" --resume") {
		t.Errorf("explore printed:\n%s", out)
	}
	src := read(t, filepath.Join(ws, "grove", "explorations", name+".md"))
	if d := doc.Parse("x.md", []byte(src), []string{"id", "title", "created"}); len(d.Diagnostics) > 0 || d.Meta.ID != id {
		t.Errorf("exploration: %v\n%s", d.Diagnostics, src)
	}
	prompt := args()
	for _, s := range []string{"Interview before you conclude", "## Owner's description\n\nAllow editing older records", "The project has no brief yet.", "Maintained specs (agreed behavior):\n- None."} {
		if !strings.Contains(prompt, s) {
			t.Errorf("instructions lack %q", s)
		}
	}
}
