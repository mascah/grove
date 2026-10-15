package cli

import (
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
		t.Fatalf("grove %v: status %d\n%s", args, code, stderr.String())
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
	if got, want := grove(t, "init"), "Wrote grove.yaml and the Grove skills for claude and codex; nothing is committed.\n"; got != want {
		t.Errorf("init said %q, want %q", got, want)
	}
	for _, h := range []string{".claude", ".agents"} {
		for _, w := range []string{"explore", "propose", "implement", "archive"} {
			if skill := read(t, filepath.Join(dir, h, "skills", "grove-"+w, "SKILL.md")); !strings.Contains(skill, "name: grove-"+w+"\n") || !strings.Contains(skill, " instructions "+w) {
				t.Errorf("%s skill grove-%s:\n%s", h, w, skill)
			}
		}
	}
	if skill := read(t, filepath.Join(dir, ".agents/skills/grove-propose/SKILL.md")); !strings.Contains(skill, ` instructions propose "<args>"`+"`") {
		t.Errorf("codex propose skill:\n%s", skill)
	}
	yaml := strings.Replace(read(t, filepath.Join(dir, "grove.yaml")), `context: ""`, "context: Python 3 standard library only.", 1)
	yaml = strings.Replace(yaml, "rules: {}", "rules:\n  proposal:\n    - Name the owner for every decision.", 1)
	os.WriteFile(filepath.Join(dir, "grove.yaml"), []byte(yaml), 0o644)
	os.WriteFile(filepath.Join(dir, ".claude/skills/grove-explore/SKILL.md"), []byte("edited"), 0o644)
	if got, want := grove(t, "init"), "Kept grove.yaml and wrote the Grove skills for claude and codex; nothing is committed.\n"; got != want {
		t.Errorf("init again said %q, want %q", got, want)
	}
	if read(t, filepath.Join(dir, "grove.yaml")) != yaml || read(t, filepath.Join(dir, ".claude/skills/grove-explore/SKILL.md")) == "edited" {
		t.Error("init again did not keep grove.yaml and rewrite the skills")
	}
	commitAll(t, dir, "Set up Grove")
	out := grove(t, "instructions", "propose")
	for _, s := range []string{"## Project context\n\nPython 3 standard library only.\n", "The project's own rules for this document:\n\n- Name the owner for every decision.\n"} {
		if !strings.Contains(out, s) {
			t.Errorf("propose instructions lack %q", s)
		}
	}

	os.WriteFile(filepath.Join(dir, "grove.yaml"), []byte("agents: [cursor]\n"), 0o644)
	var stdout, stderr strings.Builder
	if code := run([]string{"init"}, &stdout, &stderr); code != 1 || stderr.String() != "grove: grove.yaml lists agent \"cursor\"; Grove knows claude and codex\n" {
		t.Errorf("an unknown agent: status %d, %q", code, stderr.String())
	}
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
	if code != 1 || stdout.String() != want || stderr.String() != "grove: 2 problems found\n" {
		t.Errorf("check: status %d\n%s%s", code, stdout.String(), stderr.String())
	}
}

// propose plays a session that proposes a change on a new branch: it asks
// for the propose instructions and writes the change folder they name.
func propose(t *testing.T, dir string, args ...string) (changeID string) {
	t.Helper()
	git(t, dir, "switch", "-q", "-c", "older-edits")
	out := grove(t, append([]string{"instructions", "propose"}, args...)...)
	m := regexp.MustCompile("### proposal: `grove/changes/([0-9]{6}-[0-9a-z]{5})-<slug>/proposal.md`").FindStringSubmatch(out)
	if m == nil {
		t.Fatalf("propose instructions name no new change:\n%.3000s", out)
	}
	id := m[1]
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

func TestExplore(t *testing.T) {
	setUp(t)
	out := grove(t, "instructions", "explore")
	for _, s := range []string{"# Grove instructions: explore", "Interview before you conclude anything",
		"Maintained specs (agreed behavior):\n- None."} {
		if !strings.Contains(out, s) {
			t.Errorf("explore instructions lack %q", s)
		}
	}
	if strings.Contains(out, "## Documents") {
		t.Error("explore instructions offer a document to fill in")
	}
}

func TestPropose(t *testing.T) {
	dir := setUp(t)
	cid := propose(t, dir)
	if got := grove(t, "check"); got != "No problems in grove/.\n" {
		t.Errorf("check printed %q", got)
	}
	revise := grove(t, "instructions", "propose", cid)
	for _, s := range []string{"You are revising drafted change " + cid, "### proposal: `grove/changes/" + cid + "-allow-editing-older-records/proposal.md`", "## The change", "### Requirement: Edit older records"} {
		if !strings.Contains(revise, s) {
			t.Errorf("revise instructions lack %q", s)
		}
	}
	var stdout, stderr strings.Builder
	if code := run([]string{"instructions", "propose", "261015-zzzzz"}, &stdout, &stderr); code != 1 || stderr.String() != "grove: no change 261015-zzzzz in this working tree\n" {
		t.Errorf("an unknown change: status %d, %q", code, stderr.String())
	}
	if out := grove(t, "instructions", "propose", "Editing a missing record should print an error"); !strings.Contains(out, "## Owner's description\n\nEditing a missing record should print an error\n") {
		t.Errorf("propose from a description:\n%.3000s", out)
	}
}

func TestApprove(t *testing.T) {
	dir := setUp(t)
	cid := propose(t, dir)
	commitAll(t, dir, "Propose allowing edits to older records")
	out := grove(t, "review")
	for _, s := range []string{
		"Merging older-edits into main approves change " + cid + " as future work; it does not implement it.\n\n",
		"  adds     change           grove/changes/" + cid + "-allow-editing-older-records/proposal.md\n",
		"+### Requirement: Edit older records",
	} {
		if !strings.Contains(out, s) {
			t.Errorf("review lacks %q:\n%s", s, out)
		}
	}
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
	if !strings.Contains(read(t, filepath.Join(dir, ".claude/skills/grove-implement/SKILL.md")), "Run `grove instructions implement $ARGUMENTS`") {
		t.Error("the skill does not run instructions implement")
	}
	commitAll(t, dir, "Set up Grove")
	cid := propose(t, dir)
	commitAll(t, dir, "Propose")

	var stdout, stderr strings.Builder
	if code := run([]string{"instructions", "implement", cid}, &stdout, &stderr); code != 1 || stderr.String() != "grove: main holds no approved change "+cid+"\n" {
		t.Errorf("before the merge: status %d, %q", code, stderr.String())
	}
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Approve", "older-edits")
	out := grove(t, "instructions", "implement", cid)
	for _, s := range []string{"/implementation.md\">\n# Implementation: Allow editing older records", "### Requirement: Edit older records"} {
		if !strings.Contains(out, s) {
			t.Errorf("implement instructions lack %q", s)
		}
	}
}

func TestArchive(t *testing.T) {
	dir := setUp(t)
	cid := approve(t, dir)
	git(t, dir, "switch", "-q", "-c", "older-edits-impl")
	folder := "grove/changes/" + cid + "-allow-editing-older-records"
	os.WriteFile(filepath.Join(dir, folder, "implementation.md"), []byte("## T-01\n\nStatus: done\n"), 0o644)

	// The session reconciles the change spec into a new maintained spec,
	// from the template its instructions give, then archives the change.
	out := grove(t, "instructions", "archive", cid)
	for _, s := range []string{"Reconcile the specs", "Status: done", "### Requirement: Edit older records"} {
		if !strings.Contains(out, s) {
			t.Errorf("archive instructions lack %q", s)
		}
	}
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
	if got, want := grove(t, "archive", cid), "Archived "+cid+" to grove/changes/archive/"+cid+"-allow-editing-older-records/; nothing is committed.\n"; got != want {
		t.Errorf("archive printed %q, want %q", got, want)
	}
	read(t, filepath.Join(dir, "grove/changes/archive", cid+"-allow-editing-older-records", "implementation.md"))
	commitAll(t, dir, "Implement and archive")
	if out := grove(t, "review"); !strings.Contains(out, "approves no new change.") || !strings.Contains(out, "  adds     archived change  grove/changes/archive/") || !strings.Contains(out, "  adds     maintained spec  grove/specs/record-editing/spec.md") {
		t.Errorf("review of the implementation branch:\n%.2000s", out)
	}
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Implement", "older-edits-impl")

	if got, want := grove(t, "specs"), "record-editing  Record editing  grove/specs/record-editing/spec.md\n"; got != want {
		t.Errorf("specs printed %q, want %q", got, want)
	}
	if got := grove(t, "changes"); got != "" {
		t.Errorf("changes printed %q", got)
	}
	next := grove(t, "instructions", "explore")
	for _, s := range []string{
		"Maintained specs (agreed behavior):\n- record-editing \"Record editing\": `grove/specs/record-editing/spec.md`",
		"Approved changes (future work; implementation is not established):\n- None.",
	} {
		if !strings.Contains(next, s) {
			t.Errorf("fresh explore instructions lack %q", s)
		}
	}
	if len(next) >= 12<<10 {
		t.Errorf("explore instructions are %d bytes, want under 12 KiB", len(next))
	}
}
