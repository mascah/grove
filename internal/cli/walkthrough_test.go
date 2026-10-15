package cli

import (
	"github.com/mascah/grove-next2/internal/doc"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
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
	if !strings.Contains(read(t, filepath.Join(dir, "grove.yaml")), "target: main") {
		t.Error("grove.yaml does not name the target")
	}
	commitAll(t, dir, "Set up Grove")
}

func TestCheck(t *testing.T) {
	dir := setUp(t)
	if got := grove(t, "check"); got != "No problems in grove/.\n" {
		t.Errorf("check printed %q", got)
	}
	os.MkdirAll(filepath.Join(dir, "grove", "specs", "record-editing"), 0o755)
	os.WriteFile(filepath.Join(dir, "grove", "specs", "record-editing", "spec.md"), []byte("---\nid: record-editing\nowner: me\n---\n"), 0o644)
	var stdout, stderr strings.Builder
	code := run([]string{"check"}, &stdout, &stderr)
	want := "grove/specs/record-editing/spec.md:3: unknown field \"owner\"; this document takes id, title and description\n" +
		"grove/specs/record-editing/spec.md:1: title: required field is missing\n"
	if code != 1 || stdout.String() != want || stderr.String() != "grove-next2: 2 problems found\n" {
		t.Errorf("check: status %d\n%s%s", code, stdout.String(), stderr.String())
	}
}

// propose plays a session that proposes a change: it writes the change
// folder in the working tree, on a new branch, and returns its ID.
func propose(t *testing.T, dir string) (changeID string) {
	t.Helper()
	git(t, dir, "switch", "-q", "-c", "older-edits")
	id := doc.NewID(time.Now())
	change := filepath.Join(dir, "grove", "changes", id+"-allow-editing-older-records")
	os.MkdirAll(filepath.Join(change, "specs", "record-editing"), 0o755)
	os.WriteFile(filepath.Join(change, "proposal.md"), []byte("---\nid: "+id+"\ntitle: Allow editing older records\ncreated: \"2026-10-14T00:00:00Z\"\n---\n\n# Allow editing older records\n"), 0o644)
	os.WriteFile(filepath.Join(change, "specs", "record-editing", "spec.md"), []byte("# Changes to record-editing\n\n## ADDED Requirements\n\n### Requirement: Edit older records\n"), 0o644)
	os.WriteFile(filepath.Join(change, "tasks.md"), []byte("## T-01 — Edit older records\n"), 0o644)
	return id
}

// setUp runs init on a fresh fixture and commits it, as getting started §1.
func setUp(t *testing.T) string {
	t.Helper()
	dir := fixture(t)
	grove(t, "init")
	commitAll(t, dir, "Set up Grove")
	return dir
}

// approve commits the proposed change on its branch and merges it into main.
func approve(t *testing.T, dir string) (changeID string) {
	t.Helper()
	cid := propose(t, dir)
	commitAll(t, dir, "Propose allowing edits to older records")
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Approve "+cid, "older-edits")
	return cid
}

func TestApprove(t *testing.T) {
	dir := setUp(t)
	propose(t, dir)
	commitAll(t, dir, "Propose allowing edits to older records")
	if got := grove(t, "changes"); got != "" {
		t.Errorf("before the merge, changes printed %q", got)
	}
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Approve", "older-edits")
	if got := grove(t, "changes"); !strings.HasSuffix(got, "  Allow editing older records  approved\n") {
		t.Errorf("changes printed %q", got)
	}
}

func TestImplementInstructions(t *testing.T) {
	dir := fixture(t)
	grove(t, "init")
	if !strings.Contains(read(t, filepath.Join(dir, ".claude/skills/grove-implement/SKILL.md")), " instructions implement $ARGUMENTS`") {
		t.Error("the skill does not run instructions implement")
	}
	commitAll(t, dir, "Set up Grove")
	cid := propose(t, dir)
	commitAll(t, dir, "Propose")

	var stdout, stderr strings.Builder
	if code := run([]string{"instructions", "implement", cid}, &stdout, &stderr); code != 1 || stderr.String() != "grove-next2: main holds no approved change "+cid+"\n" {
		t.Errorf("before the merge: status %d, %q", code, stderr.String())
	}
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Approve", "older-edits")
	out := grove(t, "instructions", "implement", cid)
	for _, s := range []string{"Reconcile the specs", "/implementation.md\">", "<file path=\"grove/specs/<capability>/spec.md\">\n---\nid: <capability>", "### Requirement: Edit older records"} {
		if !strings.Contains(out, s) {
			t.Errorf("implement instructions lack %q", s)
		}
	}
}

func TestSeeTheResult(t *testing.T) {
	dir := setUp(t)
	cid := approve(t, dir)
	// The implementer reconciles the change spec into a new maintained spec,
	// from the template its instructions give.
	out := grove(t, "instructions", "implement", cid)
	m := regexp.MustCompile(`(?s)<file path="grove/specs/<capability>/spec.md">\n(.*?)\n</file>`).FindStringSubmatch(out)
	if m == nil {
		t.Fatalf("no spec template in:\n%s", out)
	}
	spec := strings.NewReplacer("<capability>", "record-editing", "<title>", "Record editing").Replace(m[1])
	os.MkdirAll(filepath.Join(dir, "grove", "specs", "record-editing"), 0o755)
	os.WriteFile(filepath.Join(dir, "grove", "specs", "record-editing", "spec.md"), []byte(spec), 0o644)
	if got := grove(t, "check"); got != "No problems in grove/.\n" {
		t.Errorf("check printed %q", got)
	}
	commitAll(t, dir, "Implement")

	if got, want := grove(t, "specs"), "record-editing  Record editing  grove/specs/record-editing/spec.md\n"; got != want {
		t.Errorf("specs printed %q, want %q", got, want)
	}
	if got := grove(t, "changes"); !strings.HasPrefix(got, cid+"  ") || !strings.HasSuffix(got, "  approved\n") {
		t.Errorf("changes printed %q", got)
	}
}
