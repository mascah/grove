package cli

import (
	"github.com/mascah/grove-next2/internal/doc"
	"github.com/mascah/grove-next2/internal/instruct"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

// onlyWorkspace returns the only workspace's directory.
func onlyWorkspace(t *testing.T, dir string) string {
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
	ws := onlyWorkspace(t, dir)
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

func TestComingBack(t *testing.T) {
	dir := fixture(t)
	grove(t, "init")
	commitAll(t, dir, "Set up Grove")
	_, args := harness(t, `for f in grove/explorations/*.md; do echo "- Owner: only the misc records." >> "$f"; done`)
	grove(t, "explore", "Allow editing older records")
	id := filepath.Base(onlyWorkspace(t, dir))[:12]

	if got, want := grove(t, "workspaces"), id+"  Allow editing older records  no change yet, 1 file uncommitted\n"; got != want {
		t.Errorf("workspaces printed %q, want %q", got, want)
	}
	grove(t, "explore", id)
	if !strings.Contains(args(), "- Owner: only the misc records.") {
		t.Error("a new conversation was not given the saved exploration")
	}
	if out := grove(t, "explore", id, "--resume"); args() != "--continue" || !strings.Contains(out, "Draft: exploration "+id) {
		t.Errorf("--resume ran claude %q and printed %q", args(), out)
	}
}

func TestBaselineTemplateParses(t *testing.T) {
	src := instruct.Spec("record-editing", "Record editing")
	if d := doc.Parse("specs/record-editing/spec.md", []byte(src), []string{"id", "title"}); len(d.Diagnostics) > 0 {
		t.Errorf("%v\n%s", d.Diagnostics, src)
	}
}

func TestCheck(t *testing.T) {
	dir := fixture(t)
	grove(t, "init")
	commitAll(t, dir, "Set up Grove")
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

// propose plays a session that switches to propose: it writes the change the
// instructions name, as a session following them would.
func propose(t *testing.T, ws, id string) (changeID string) {
	t.Helper()
	t.Chdir(ws)
	out := grove(t, "instructions", "propose", id)
	m := regexp.MustCompile("Create the change in `grove/(changes/([0-9]{6}-[0-9a-z]{5})-[a-z0-9-]+)/`").FindStringSubmatch(out)
	if m == nil {
		t.Fatalf("propose instructions name no new change:\n%.2000s", out)
	}
	dir := filepath.Join(ws, "grove", m[1])
	os.MkdirAll(filepath.Join(dir, "specs", "record-editing"), 0o755)
	os.WriteFile(filepath.Join(dir, "proposal.md"), []byte("---\nid: "+m[2]+"\ntitle: Allow editing older records\ncreated: \"2026-10-14T00:00:00Z\"\nexploration: "+id+"\n---\n\n# Allow editing older records\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "specs", "record-editing", "spec.md"), []byte("# Changes to record-editing\n\n## ADDED Requirements\n\n### Requirement: Edit older records\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "tasks.md"), []byte("## T-01 — Edit older records\n"), 0o644)
	return m[2]
}

func TestPropose(t *testing.T) {
	dir := fixture(t)
	grove(t, "init")
	commitAll(t, dir, "Set up Grove")
	harness(t, "")
	grove(t, "explore", "Allow editing older records")
	ws := onlyWorkspace(t, dir)
	id := filepath.Base(ws)[:12]

	cid := propose(t, ws, id)
	if got := grove(t, "check"); got != "No problems in grove/.\n" {
		t.Errorf("check printed %q", got)
	}
	again := grove(t, "instructions", "propose", id)
	if strings.Contains(again, "## Templates") || !strings.Contains(again, "## This change, as drafted so far") || !strings.Contains(again, "### Requirement: Edit older records") {
		t.Errorf("a second propose did not continue the drafted change:\n%.3000s", again)
	}
	t.Chdir(dir)
	if got := grove(t, "workspaces"); !strings.Contains(got, "change "+cid+" draft, 4 files uncommitted") {
		t.Errorf("workspaces printed %q", got)
	}
	if got := grove(t, "explore", id, "--resume"); !strings.HasPrefix(got, "Draft: change "+cid+" in ") {
		t.Errorf("report: %q", got)
	}
}

func TestCommitAndMerge(t *testing.T) {
	dir := fixture(t)
	grove(t, "init")
	commitAll(t, dir, "Set up Grove")
	harness(t, "")
	grove(t, "explore", "Allow editing older records")
	ws := onlyWorkspace(t, dir)
	cid := propose(t, ws, filepath.Base(ws)[:12])
	commitAll(t, ws, "Propose allowing edits to older records")
	t.Chdir(dir)

	out := grove(t, "review", cid)
	branch := "grove/" + filepath.Base(ws)
	for _, s := range []string{
		"Merging " + branch + " into main approves change " + cid + " as future work; it does not implement it.\n",
		"  adds     change           grove/changes/" + cid + "-",
		"  adds     exploration      grove/explorations/",
		"+### Requirement: Edit older records",
	} {
		if !strings.Contains(out, s) {
			t.Errorf("review lacks %q:\n%s", s, out)
		}
	}
	if got := grove(t, "changes"); got != "" {
		t.Errorf("before the merge, changes printed %q", got)
	}
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Approve "+cid, branch)
	if got, want := grove(t, "changes"), cid+"  Allow editing older records  approved\n"; got != want {
		t.Errorf("changes printed %q, want %q", got, want)
	}
}

func TestImplementInstructions(t *testing.T) {
	dir := fixture(t)
	grove(t, "init")
	if !strings.Contains(read(t, filepath.Join(dir, ".claude/skills/grove-implement/SKILL.md")), " instructions implement $ARGUMENTS`") {
		t.Error("the skill does not run instructions implement")
	}
	commitAll(t, dir, "Set up Grove")
	harness(t, "")
	grove(t, "explore", "Allow editing older records")
	ws := onlyWorkspace(t, dir)
	cid := propose(t, ws, filepath.Base(ws)[:12])
	commitAll(t, ws, "Propose")
	t.Chdir(dir)

	var stdout, stderr strings.Builder
	if code := run([]string{"instructions", "implement", cid}, &stdout, &stderr); code != 1 || stderr.String() != "grove-next2: main holds no approved change "+cid+"\n" {
		t.Errorf("before the merge: status %d, %q", code, stderr.String())
	}
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Approve", "grove/"+filepath.Base(ws))
	out := grove(t, "instructions", "implement", cid)
	for _, s := range []string{"Reconcile the specs", "/implementation.md\">", "<file path=\"grove/specs/record-editing/spec.md\">\n---\nid: record-editing", "### Requirement: Edit older records"} {
		if !strings.Contains(out, s) {
			t.Errorf("implement instructions lack %q", s)
		}
	}
}
