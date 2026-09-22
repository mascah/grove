package cli

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// gitFixture commits the plain fixture so new can allocate IDs.
func gitFixture(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root := projectFixture(t)
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"add", "-A"}, {"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-c", "maintenance.auto=false", "commit", "-q", "-m", "init"}} {
		if out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	return root
}

func TestNewCreatesRecordAndReadCommandsLeaveNoState(t *testing.T) {
	t.Parallel()
	root := gitFixture(t)
	state := filepath.Join(root, ".git", "grove")
	for _, args := range [][]string{{"list"}, {"show", "G-001"}, {"check"}} {
		var out, errOut bytes.Buffer
		if code := Run(args, root, &out, &errOut); code != 0 {
			t.Fatal(errOut.String())
		}
	}
	if _, err := os.Stat(state); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("read commands must not create allocator state: %v", err)
	}
	var out, errOut bytes.Buffer
	if code := Run([]string{"new", "work", "Second thing", "--slug", "second"}, root, &out, &errOut); code != 0 {
		t.Fatal(errOut.String())
	}
	if out.String() != "docs/records/G-003-second.md\n" || !strings.Contains(errOut.String(), "Initialized") {
		t.Fatalf("stdout=%q stderr=%q", out.String(), errOut.String())
	}
	out.Reset()
	errOut.Reset()
	if code := Run([]string{"--project", root, "new", "question", "Why?"}, t.TempDir(), &out, &errOut); code != 0 || out.String() != "docs/records/G-004-why.md\n" {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	out.Reset()
	if code := Run([]string{"check"}, root, &out, &errOut); code != 0 || !strings.Contains(out.String(), "4 records") {
		t.Fatalf("created records must validate: %s %s", out.String(), errOut.String())
	}
	out.Reset()
	if code := Run([]string{"list"}, root, &out, &errOut); code != 0 || !strings.Contains(out.String(), "G-003  work      proposed  Second thing") {
		t.Fatalf("list must show the created record: %s", out.String())
	}
}

func TestNewUsageAndFailures(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{{"new"}, {"new", "work"}, {"new", "work", "T", "extra"}, {"new", "work", "T", "--slug"}, {"list", "--slug", "x"}, {"new", "work", "T", "--slug=a", "--slug=b"}} {
		var out, errOut bytes.Buffer
		if code := Run(args, t.TempDir(), &out, &errOut); code != 2 || out.Len() != 0 {
			t.Fatalf("%v: code=%d stderr=%s", args, code, errOut.String())
		}
	}
	root := projectFixture(t)
	before := hashes(t, root)
	for _, args := range [][]string{{"new", "work", "Outside Git"}, {"new", "release", "Bad type"}} {
		var out, errOut bytes.Buffer
		if code := Run(args, root, &out, &errOut); code != 1 || out.Len() != 0 || !strings.Contains(errOut.String(), "grove:") {
			t.Fatalf("%v: code=%d stdout=%q stderr=%q", args, code, out.String(), errOut.String())
		}
	}
	if !reflect.DeepEqual(before, hashes(t, root)) {
		t.Fatal("failed creation changed project files")
	}
}
