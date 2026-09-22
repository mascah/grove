package repo

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCommandIgnoresAmbientRepository is the G-089 regression: with GIT_DIR
// and GIT_WORK_TREE naming another repository, as a Git hook exports them,
// Command still acts on the directory it was given.
func TestCommandIgnoresAmbientRepository(t *testing.T) {
	decoy, target := t.TempDir(), t.TempDir()
	for _, dir := range []string{decoy, target} {
		git(t, dir, "init", "-q", "-b", "main")
	}
	t.Setenv("GIT_DIR", filepath.Join(decoy, ".git"))
	t.Setenv("GIT_WORK_TREE", decoy)
	gitDir := func(cmd *exec.Cmd) string {
		t.Helper()
		out, err := cmd.Output()
		if err != nil {
			t.Fatal(err)
		}
		got, err := filepath.EvalSymlinks(strings.TrimSpace(string(out)))
		if err != nil {
			t.Fatal(err)
		}
		return got
	}
	want, err := filepath.EvalSymlinks(filepath.Join(target, ".git"))
	if err != nil {
		t.Fatal(err)
	}
	if got := gitDir(Command(context.Background(), target, "rev-parse", "--absolute-git-dir")); got != want {
		t.Fatalf("Command acted on %s, want %s", got, want)
	}
	// The unguarded child must land in the decoy, or this test proves nothing.
	if got := gitDir(exec.Command("git", "-C", target, "rev-parse", "--absolute-git-dir")); got == want {
		t.Fatalf("GIT_DIR did not reach an unguarded child; the guard was not exercised")
	}
}
