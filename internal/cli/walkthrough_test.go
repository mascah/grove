package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
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
		for _, w := range []string{"explore", "prototype", "propose", "plan", "implement", "archive"} {
			if skill := read(t, filepath.Join(dir, h, "skills", "grove-"+w, "SKILL.md")); !strings.Contains(skill, "name: grove-"+w+"\n") || !strings.Contains(skill, " instructions "+w) || w == "implement" && !strings.Contains(skill, "Implement an approved or drafted change") {
				t.Errorf("%s skill grove-%s:\n%s", h, w, skill)
			}
		}
	}
	if skill := read(t, filepath.Join(dir, ".agents/skills/grove-propose/SKILL.md")); !strings.Contains(skill, ` instructions propose "<args>"`+"`") || skill != read(t, filepath.Join(dir, ".claude/skills/grove-propose/SKILL.md")) {
		t.Errorf("propose skill, which should be one body for every harness:\n%s", skill)
	}
	yaml := strings.Replace(read(t, filepath.Join(dir, "grove.yaml")), `context: ""`, "context: Run python3 -m unittest before saying a task is done.", 1)
	yaml = strings.Replace(yaml, "rules: {}", "rules:\n  documents:\n    proposal:\n      - Name the owner for every decision.\n  skills:\n    propose:\n      - Commit the change folder when check passes.", 1)
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
	for _, s := range []string{"## Project context\n\nRun python3 -m unittest before saying a task is done.\n\n## What Grove requires of propose\n", "The project's own rules for this document:\n\n- Name the owner for every decision.\n",
		"The project's own rules for this skill:\n\n- Commit the change folder when check passes.\n"} {
		if !strings.Contains(out, s) {
			t.Errorf("propose instructions lack %q", s)
		}
	}
}

// TestHarnesses follows the walkthrough's agents from the repository,
// commands, dropping a harness, and a harness Grove does not have.
func TestHarnesses(t *testing.T) {
	dir := fixture(t)
	for _, d := range []string{".claude", ".opencode", ".agents/skills/notes"} {
		os.MkdirAll(filepath.Join(dir, d), 0o755)
	}
	if got, want := grove(t, "init"), "Wrote grove.yaml and the Grove skills for claude and opencode; nothing is committed.\n"; got != want {
		t.Errorf("init said %q, want %q", got, want)
	}
	yaml := read(t, filepath.Join(dir, "grove.yaml"))
	if !strings.Contains(yaml, "agents:\n  claude: skills\n  opencode: skills\n#") || read(t, filepath.Join(dir, ".claude/skills/grove-implement/SKILL.md")) != read(t, filepath.Join(dir, ".agents/skills/grove-implement/SKILL.md")) {
		t.Errorf("init chose other agents, or skills that differ:\n%s", yaml)
	}

	os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte("# Agents\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "grove.yaml"), []byte(strings.Replace(yaml, "claude: skills", "claude: commands", 1)), 0o644)
	if got, want := grove(t, "init"), "Kept grove.yaml and wrote the Grove skills for opencode and commands for claude, and removed grove-archive, grove-explore, grove-implement, grove-plan, grove-propose and grove-prototype from .claude/skills; nothing is committed.\n"; got != want {
		t.Errorf("init with claude's commands said %q, want %q", got, want)
	}
	if cmd := read(t, filepath.Join(dir, ".claude/commands/grove/propose.md")); !strings.HasPrefix(cmd, "---\ndescription: ") || !strings.Contains(cmd, "Run `grove instructions propose \"$ARGUMENTS\"` and follow what it prints, in this conversation.\n") {
		t.Errorf("claude's propose command:\n%s", cmd)
	}

	os.WriteFile(filepath.Join(dir, "grove.yaml"), []byte(strings.Replace(yaml, "  claude: skills\n  opencode: skills\n", "  claude: commands\n", 1)), 0o644)
	if got, want := grove(t, "init"), "Kept grove.yaml and wrote the Grove commands for claude, and removed grove-archive, grove-explore, grove-implement, grove-plan, grove-propose and grove-prototype from .agents/skills; nothing is committed.\n"; got != want {
		t.Errorf("init without opencode said %q, want %q", got, want)
	}
	if _, err := os.Stat(filepath.Join(dir, ".agents/skills/notes")); err != nil {
		t.Errorf("init removed a skill that is not Grove's: %v", err)
	}

	commitAll(t, dir, "Set up Grove")
	os.WriteFile(filepath.Join(dir, "grove.yaml"), []byte(strings.Replace(yaml, "  claude: skills\n  opencode: skills\n", "  codex: commands\n  zed: skills\n", 1)), 0o644)
	var stdout, stderr strings.Builder
	want := "grove.yaml:9: harness \"codex\" has no form \"commands\"; it has skills\n" +
		"grove.yaml:10: unknown harness \"zed\"; agents takes claude, codex, cursor, hermes, opencode and pi\n"
	if code := run([]string{"check"}, &stdout, &stderr); code != 1 || stdout.String() != want {
		t.Errorf("check of agents: status %d\n%s%s", code, stdout.String(), stderr.String())
	}
	stderr.Reset()
	if code := run([]string{"init"}, &stdout, &stderr); code != 1 || stderr.String() != "grove: grove.yaml: harness \"codex\" has no form \"commands\"; it has skills\n" {
		t.Errorf("init with a form codex has not: status %d, %q", code, stderr.String())
	}
}

func TestVersion(t *testing.T) {
	version = "0.1.0"
	t.Cleanup(func() { version = "" })
	if got := grove(t, "version"); got != "grove 0.1.0\n" {
		t.Errorf("version said %q", got)
	}
}

func TestCheck(t *testing.T) {
	dir := setUp(t)
	if got := grove(t, "check"); got != "No problems in grove/.\n" {
		t.Errorf("check printed %q", got)
	}
	config := read(t, filepath.Join(dir, "grove.yaml"))
	os.WriteFile(filepath.Join(dir, "grove.yaml"), []byte(strings.Replace(config, "rules: {}", "rules:\n  derp: [x]\n  documents:\n    propsal: [x]\n  skills:\n    apply: [x]\nrule: {}", 1)), 0o644)
	var stdout, stderr strings.Builder
	code := run([]string{"check"}, &stdout, &stderr)
	want := "grove.yaml:30: unknown rules key \"derp\"; rules takes documents and skills\n" +
		"grove.yaml:32: unknown rules.documents key \"propsal\"; kinds are proposal, change-spec, tasks, implementation, spec\n" +
		"grove.yaml:34: unknown rules.skills key \"apply\"; its lines go under rules.skills.implement\n" +
		"grove.yaml:35: unknown field \"rule\"; grove.yaml takes root, target, agents, schema, context and rules\n"
	if code != 1 || stdout.String() != want || stderr.String() != "grove: 4 problems found\n" {
		t.Errorf("check of grove.yaml: status %d\n%s%s", code, stdout.String(), stderr.String())
	}
	grove(t, "changes")
	os.WriteFile(filepath.Join(dir, "grove.yaml"), []byte(config), 0o644)

	os.MkdirAll(filepath.Join(dir, "grove", "specs", "record-editing"), 0o755)
	os.WriteFile(filepath.Join(dir, "grove", "specs", "record-editing", "spec.md"), []byte("# Record editing\n\n## ADDED Requirements\n\n"+
		"### Requirement: Keep a category\n\nThe system shall keep it.\n\n### Requirement: Delete a record\n\nUsers can delete a record.\n\n#### Scenario: Delete\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "grove", "specs", "record-editing", ".grove.yaml"), []byte("title: Record editing\nowner: me\n"), 0o644)
	for file, body := range map[string]string{"notes/ideas.md": "# Ideas\n", "old/.grove.yaml": "title: Old\n", "untitled/spec.md": "## Purpose\n\nUntitled.\n"} {
		os.MkdirAll(filepath.Join(dir, "grove", "specs", filepath.Dir(file)), 0o755)
		os.WriteFile(filepath.Join(dir, "grove", "specs", file), []byte(body), 0o644)
	}
	for folder, meta := range map[string]string{"261015-k3m9q-no-proposal": "title: T\ncreated: x\n", "261015-p4x8d-extra-field": "title: T\ncreated: x\nowner: me\n",
		"261015-q7r2s-drifted": "title: New title\ncreated: x\n", "no-id": "title: T\ncreated: x\n"} {
		os.MkdirAll(filepath.Join(dir, "grove", "changes", folder), 0o755)
		os.WriteFile(filepath.Join(dir, "grove", "changes", folder, ".grove.yaml"), []byte(meta), 0o644)
		if folder != "261015-k3m9q-no-proposal" {
			os.WriteFile(filepath.Join(dir, "grove", "changes", folder, "proposal.md"), []byte("# T\n"), 0o644)
		}
	}
	for _, folder := range []string{"261015-p4x8d-extra-field", "archive/261014-a0000-old"} {
		os.MkdirAll(filepath.Join(dir, "grove", "changes", folder, "specs", "record-editing"), 0o755)
		os.WriteFile(filepath.Join(dir, "grove", "changes", folder, "specs", "record-editing", "spec.md"), []byte("## REMOVED Requirements\n\n### Requirement: Old rule\n\nNot wanted.\n"), 0o644)
	}
	os.WriteFile(filepath.Join(dir, "grove", "changes", "archive/261014-a0000-old", ".grove.yaml"), []byte("title: Old\ncreated: x\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "grove", "changes", "archive/261014-a0000-old", "proposal.md"), []byte("# Old\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "grove", "changes", "archive/261014-a0000-old", "specs", "record-editing", "spec.md"), []byte("## ADDED Requirements\n\n### Requirement: Old rule\n\nNot wanted.\n"), 0o644)
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"check"}, &stdout, &stderr)
	want = "grove/specs/old/.grove.yaml: no longer read, since a spec's title and description come from its spec.md; delete it\n" +
		"grove/specs/record-editing/.grove.yaml: no longer read, since a spec's title and description come from its spec.md; delete it\n" +
		"grove/specs/record-editing/spec.md:5: requirement \"Keep a category\" has no scenario\n" +
		"grove/specs/record-editing/spec.md:9: requirement \"Delete a record\" has no \"shall\"\n" +
		"grove/specs/record-editing/spec.md:3: this heading belongs in a change spec; a maintained spec lists its requirements under ## Requirements\n" +
		"grove/specs/untitled/spec.md: no \"# \" heading to give the spec its title\n" +
		"grove/changes/261015-k3m9q-no-proposal/proposal.md: missing; it is the change's main document in schema default\n" +
		"grove/changes/261015-p4x8d-extra-field/.grove.yaml:3: unknown field \"owner\"; this file takes title, created, schema, description, depends_on and base\n" +
		"grove/changes/261015-q7r2s-drifted/proposal.md:1: the first heading should be \"# New title\", the title in .grove.yaml\n" +
		"grove/changes/no-id/: the folder name must start with a change ID, such as 261015-k3m9q\n"
	if code != 1 || stdout.String() != want || stderr.String() != "grove: 10 problems found\n" {
		t.Errorf("check: status %d\n%s%s", code, stdout.String(), stderr.String())
	}
	commitAll(t, dir, "Specs")
	if got := grove(t, "specs"); strings.Contains(got, "notes") || !strings.Contains(got, "untitled    grove/specs/untitled/spec.md\n") {
		t.Errorf("specs printed %q", got)
	}
}

// propose plays a session that proposes a change on a new branch: it asks
// for the propose instructions and writes the change folder they name,
// with meta's lines added to its .grove.yaml.
func propose(t *testing.T, dir, branch, title, meta string, args ...string) (changeID, folder string) {
	t.Helper()
	git(t, dir, "switch", "-q", "-c", branch)
	out := grove(t, append([]string{"instructions", "propose"}, args...)...)
	m := regexp.MustCompile("### proposal: `grove/changes/([0-9]{6}-[0-9a-z]{5})-<slug>/proposal.md`").FindStringSubmatch(out)
	if m == nil {
		t.Fatalf("propose instructions name no new change:\n%.3000s", out)
	}
	folder = filepath.Join(dir, "grove", "changes", m[1]+"-"+strings.ToLower(strings.ReplaceAll(title, " ", "-")))
	os.MkdirAll(filepath.Join(folder, "specs", "record-editing"), 0o755)
	os.WriteFile(filepath.Join(folder, ".grove.yaml"), []byte("title: "+title+"\ncreated: 2026-10-14T00:00:00Z\n"+meta), 0o644)
	os.WriteFile(filepath.Join(folder, "proposal.md"), []byte("# "+title+"\n"), 0o644)
	os.WriteFile(filepath.Join(folder, "specs", "record-editing", "spec.md"), []byte("# Changes to record-editing\n\n## ADDED Requirements\n\n"+olderEdits), 0o644)
	os.WriteFile(filepath.Join(folder, "tasks.md"), []byte("## T-01 — "+title+"\n"), 0o644)
	return m[1], folder
}

// olderEdits is the requirement propose writes.
const olderEdits = "### Requirement: Edit older records\n\nWhen a user edits a record made before the current categories, the system shall save the edit.\n\n#### Scenario: A legacy record\n\n- THEN the edit is saved\n"

// setUp runs init on a fresh fixture and commits it, as getting started §1.
func setUp(t *testing.T) string {
	t.Helper()
	dir := fixture(t)
	grove(t, "init")
	commitAll(t, dir, "Set up Grove")
	return dir
}

// approve proposes a change as propose does, commits it on its branch and
// merges it into main.
func approve(t *testing.T, dir, branch, title, meta string) (changeID, folder string) {
	t.Helper()
	cid, folder := propose(t, dir, branch, title, meta)
	commitAll(t, dir, "Propose "+title)
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Approve "+cid, branch)
	return cid, folder
}

func TestExplore(t *testing.T) {
	setUp(t)
	out := grove(t, "instructions", "explore")
	for _, s := range []string{"# Grove instructions: explore", "Interview before you conclude anything",
		"Maintained specs (agreed behavior):\n- None.", "`instructions prototype \"<question>\"`"} {
		if !strings.Contains(out, s) {
			t.Errorf("explore instructions lack %q", s)
		}
	}
	if strings.Contains(out, "## Documents") {
		t.Error("explore instructions offer a document to fill in")
	}
}

func TestPrototype(t *testing.T) {
	setUp(t)
	out := grove(t, "instructions", "prototype", "Does a budget line read clearly?")
	for _, s := range []string{"You are prototyping the owner's question below", "## What Grove requires of prototype\n", "branch `prototype/<slug>`",
		"## Owner's description\n\nDoes a budget line read clearly?\n", "Maintained specs (agreed behavior):"} {
		if !strings.Contains(out, s) {
			t.Errorf("prototype instructions lack %q", s)
		}
	}
	if strings.Contains(out, "## Documents") {
		t.Error("prototype instructions offer a document to fill in")
	}
	if out := grove(t, "instructions", "propose"); !strings.Contains(out, "citing the prototype branch and the commit") || strings.Contains(out, "## Project context") {
		t.Error("propose instructions do not say to cite a prototype, or have a context the project has not")
	}
}

func TestPropose(t *testing.T) {
	dir := setUp(t)
	cid, _ := propose(t, dir, "older-edits", "Allow editing older records", "")
	if got := grove(t, "check"); got != "No problems in grove/.\n" {
		t.Errorf("check printed %q", got)
	}
	revise := grove(t, "instructions", "propose", cid)
	for _, s := range []string{"You are revising drafted change " + cid, "### proposal: `grove/changes/" + cid + "-allow-editing-older-records/proposal.md`", "## The change", "\n- `grove/changes/" + cid + "-allow-editing-older-records/specs/record-editing/spec.md`\n"} {
		if !strings.Contains(revise, s) {
			t.Errorf("revise instructions lack %q", s)
		}
	}
	if strings.Contains(revise, "### Requirement: Edit older records") {
		t.Errorf("revise instructions print the change's text:\n%.3000s", revise)
	}
	var stdout, stderr strings.Builder
	if code := run([]string{"instructions", "propose", "261015-zzzzz"}, &stdout, &stderr); code != 1 || stderr.String() != "grove: no change 261015-zzzzz in this working tree\n" {
		t.Errorf("an unknown change: status %d, %q", code, stderr.String())
	}
	out := grove(t, "instructions", "propose", "Editing a missing record should print an error")
	req, how := strings.Index(out, "## What Grove requires of propose\n\n- Write the documents below"), strings.Index(out, "## How to propose\n\n**Propose says")
	if !strings.Contains(out, "## Owner's description\n\nEditing a missing record should print an error\n") || req < 0 || how < req ||
		!strings.Contains(out, "under `Files:` after `Result:`") || !strings.Contains(out, "\nFiles:\n- [`path/to/file`]\n") || strings.Contains(out, "file-by-file") ||
		!strings.Contains(out, "before you\n  write any document, ask one round") || !strings.Contains(out, "worded so that \"yes\"\n  accepts it") || !strings.Contains(out, "Before you tell the owner the change is written, check") || !strings.Contains(out, "### 5. Check the change") || strings.Contains(out, "Assumptions to confirm") || strings.Contains(out, "implementing unattended ignores it") {
		t.Errorf("propose from a description:\n%.3000s", out)
	}
	for _, s := range []string{"say so and ask\n  whether to split it", "`depends_on`", "Order it by consequence: first what changes behavior, scope\n  or compatibility or loses data", "as\n  recommended defaults on one line, not a list, that a single \"yes\"\n  accepts",
		"Write a task's `Delegated:` line", "only where the\n  owner agreed to it in the round", "A task may have a `Delegated:` line, right after `Files:`", "Implementing honors it,\n  with the owner and unattended.", "Delegated: test layout and private helpers"} {
		if !strings.Contains(out, s) {
			t.Errorf("propose instructions lack %q", s)
		}
	}
	if !strings.Contains(out, "offer both next steps\n  and choose neither") || !strings.Contains(out, "`/grove-implement <change-id>` on this branch so one merge\n  approves it with its code") ||
		!strings.Contains(out, "For a change it lists as approved, the next\n  step is merging the revision") {
		t.Errorf("propose's finish does not offer one merge or two:\n%.3000s", out)
	}
	// A title holding ": ", filled into the printed template, is valid YAML.
	m := regexp.MustCompile("(?s)`grove/changes/([0-9]{6}-[0-9a-z]{5})-<slug>/proposal.md`.*/.grove.yaml\">\n(.*?)</file>").FindStringSubmatch(out)
	if m == nil {
		t.Fatalf("propose instructions print no .grove.yaml template:\n%.3000s", out)
	}
	colon := filepath.Join(dir, "grove/changes", m[1]+"-records-delete-by-id")
	os.MkdirAll(colon, 0o755)
	os.WriteFile(filepath.Join(colon, ".grove.yaml"), []byte(strings.Replace(m[2], "<title>", "Records: delete by ID", 1)), 0o644)
	os.WriteFile(filepath.Join(colon, "proposal.md"), []byte("# Records: delete by ID\n"), 0o644)
	os.WriteFile(filepath.Join(colon, "tasks.md"), []byte("## T-01 — Delete by ID\n"), 0o644)
	if got := grove(t, "check"); got != "No problems in grove/.\n" {
		t.Errorf("check of a title with a colon printed %q", got)
	}
	os.RemoveAll(colon)

	os.MkdirAll(filepath.Join(dir, "grove/schemas/default/documents/tasks"), 0o755)
	os.WriteFile(filepath.Join(dir, "grove/schemas/default/documents/tasks/template.md"), []byte("## T-01 — [title]\n\nRollback: [how to undo this task]\n"), 0o644)
	if out := grove(t, "instructions", "propose", "x"); !strings.Contains(out, "/tasks.md\">\n## T-01 — [title]\n\nRollback: [how to undo this task]\n</file>") {
		t.Errorf("propose with the project's tasks template:\n%s", out)
	}
	grove(t, "check")
}

func TestApprove(t *testing.T) {
	dir := setUp(t)
	cid, _ := propose(t, dir, "older-edits", "Allow editing older records", "")
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

// TestPlanned follows the walkthrough's "Planned against actual": review
// marks each file a branch changes outside grove/ against the change's
// tasks, before and after archiving.
func TestPlanned(t *testing.T) {
	dir := setUp(t)
	cid, folder := propose(t, dir, "c", "Allow editing older records", "")
	os.WriteFile(filepath.Join(folder, "tasks.md"), []byte("## T-01 — Allow editing older records\n\nFiles:\n- `ledger/records.py`\n- `tests/fixtures/`\n"), 0o644)
	commitAll(t, dir, "Propose")
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Approve", "c")
	git(t, dir, "switch", "-q", "-c", "c-impl")
	os.WriteFile(filepath.Join(folder, "implementation.md"), []byte("## T-01\n\nStatus: done\n\nAlso changed: `README.md` — T-01 makes its usage example false\n"), 0o644)
	os.MkdirAll(filepath.Join(dir, "tests/fixtures"), 0o755)
	for _, f := range []string{"ledger/records.py", "tests/fixtures/records.json", "README.md", "Cargo.lock"} {
		os.WriteFile(filepath.Join(dir, f), []byte("changed\n"), 0o644)
	}
	commitAll(t, dir, "Build")
	want := "\n  named nowhere  Cargo.lock\n  also changed " + cid + ":T-01: T-01 makes its usage example false  README.md\n  planned " + cid + ":T-01  ledger/records.py\n  planned " + cid + ":T-01  tests/fixtures/records.json\n\n"
	if out := grove(t, "review"); !strings.Contains(out, want) {
		t.Errorf("review does not mark the files against the plan:\n%.2000s", out)
	}
	maintain(t, dir, olderEdits)
	grove(t, "archive", cid)
	commitAll(t, dir, "Archive")
	if out := grove(t, "review"); !strings.Contains(out, want) {
		t.Errorf("review after archiving does not mark the files:\n%.2000s", out)
	}
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Implement", "c-impl")
	git(t, dir, "switch", "-q", "-c", "other")
	os.WriteFile(filepath.Join(dir, "Cargo.lock"), []byte("again\n"), 0o644)
	commitAll(t, dir, "Other")
	if out := grove(t, "review"); strings.Contains(out, "named nowhere") {
		t.Errorf("review of a branch with no change marks files:\n%.2000s", out)
	}
}

// TestOneMerge follows the walkthrough's "One merge for a small change": a
// change drafted, built and archived on one branch, approved by one merge.
func TestOneMerge(t *testing.T) {
	dir := setUp(t)
	cid, folder := propose(t, dir, "rename", "Allow editing older records", "")
	commitAll(t, dir, "Propose")
	grove(t, "instructions", "implement", cid)
	os.WriteFile(filepath.Join(folder, "implementation.md"), []byte("## T-01\n\nStatus: done\n"), 0o644)
	maintain(t, dir, olderEdits)
	grove(t, "archive", cid)
	commitAll(t, dir, "Implement and archive "+cid)
	if out := grove(t, "review"); !strings.HasPrefix(out, "Merging rename into main approves change "+cid+" with its code.\n\n") {
		t.Errorf("review of a draft built on its branch:\n%.1000s", out)
	}
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Merge rename", "rename")
	if changes, specs := grove(t, "changes"), grove(t, "specs"); changes != "" || !strings.Contains(specs, "record-editing") {
		t.Errorf("after one merge, changes printed %q and specs %q", changes, specs)
	}
}

func TestImplementInstructions(t *testing.T) {
	dir := fixture(t)
	os.MkdirAll(filepath.Join(dir, ".claude/skills/grove-apply"), 0o755)
	os.MkdirAll(filepath.Join(dir, ".claude/commands/grove"), 0o755)
	os.WriteFile(filepath.Join(dir, ".claude/skills/grove-apply/SKILL.md"), []byte("old"), 0o644)
	os.WriteFile(filepath.Join(dir, ".claude/commands/grove/apply.md"), []byte("old"), 0o644)
	if got := grove(t, "init"); !strings.Contains(got, ", and removed apply.md from .claude/commands/grove and grove-apply from .claude/skills;") {
		t.Errorf("init after upgrading said %q", got)
	}
	if !strings.Contains(read(t, filepath.Join(dir, ".claude/skills/grove-implement/SKILL.md")), "Run `grove instructions implement <args>`") {
		t.Error("the skill does not run instructions implement")
	}
	yaml := strings.Replace(read(t, filepath.Join(dir, "grove.yaml")), `context: ""`, "context: Run python3 -m unittest before saying a task is done.", 1)
	os.WriteFile(filepath.Join(dir, "grove.yaml"), []byte(strings.Replace(yaml, "rules: {}", "rules:\n  skills:\n    implement:\n      - Commit after each task, naming its task ID.", 1)), 0o644)
	commitAll(t, dir, "Set up Grove")
	cid, _ := propose(t, dir, "older-edits", "Allow editing older records", "")
	commitAll(t, dir, "Propose")

	// A draft: built with the owner on its branch, refused unattended or
	// while it waits.
	var stdout, stderr strings.Builder
	if code := run([]string{"instructions", "implement", cid, "--unattended"}, &stdout, &stderr); code != 1 || stderr.String() != "grove: main holds no approved change "+cid+"\n" {
		t.Errorf("unattended before the merge: status %d, %q", code, stderr.String())
	}
	draft := grove(t, "instructions", "implement", cid)
	for _, s := range []string{"You are implementing drafted change " + cid + ", \"Allow editing older records\", as this working tree holds it, with the owner", "Merging this branch into the\n  target approves it with its code, so edit its proposal, change specs and\n  tasks only with the owner's agreement", "\n- `grove/changes/" + cid + "-allow-editing-older-records/tasks.md`\n"} {
		if !strings.Contains(draft, s) {
			t.Errorf("draft instructions lack %q:\n%.3000s", s, draft)
		}
	}
	later, _ := propose(t, dir, "later", "Show saved records", "depends_on: ["+cid+"]\n")
	stderr.Reset()
	if code := run([]string{"instructions", "implement", later}, &stdout, &stderr); code != 1 || stderr.String() != "grove: "+later+" waits on "+cid+", which this working tree has not archived as finished\n" {
		t.Errorf("a draft that waits: status %d, %q", code, stderr.String())
	}
	git(t, dir, "checkout", "-q", "--", ".")
	git(t, dir, "clean", "-qfd")
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Approve", "older-edits")
	k := "261022-k0000"
	stderr.Reset()
	if code := run([]string{"instructions", "apply", k}, &stdout, &stderr); code != 2 || !strings.HasPrefix(stderr.String(), "usage: grove") {
		t.Errorf("instructions apply: status %d, %q", code, stderr.String())
	}

	shared := []string{"## What Grove requires of implement\n\n", "Before you finish, run the project's tests, re-read each done task's", "each unverified scenario. When every task is done", "Commit only when the owner or the project's\nrules ask.",
		"The project's own rules for this skill:\n\n- Commit after each task, naming its task ID.\n"}
	// Both modes edit a file a task needs or makes false, and say why.
	shared = append(shared, "changed to be built as the change says, or makes", "`` Also changed: `<path>` — <why> `` under that task")
	// Both modes build every task but a checkpoint's dependents, and take answers.
	shared = append(shared, "A checkpoint is a task whose acceptance no check you can run shows", "and mark it `Status: to check`", "build nothing of it and mark it `Status: blocked`", "on a checkpoint unbuilt, naming the checkpoint under it", "`Answer: \"<their words>\"`", "Finish: report the checkpoints first")
	gone := []string{"with their agreement edit it", "stays as it is", "ignore\n  `Delegated:` lines", "write no\n  `Also changed:` line", "tell the owner how to try it", "an open question in the change decides"}
	owner := []string{"stop it at once, ask the owner and wait", "without asking and without blocking the task", "by running Grove's\n  `instructions propose <change-id>` and following it, never editing a\n  `Files:` line or `.grove.yaml` by hand", "brought the file, recommend propose, which\n  re-bases the change", "Do not stop between tasks.", "stop once: put each checkpoint to the owner",
		"Make the choices and change the paths a task's `Delegated:` line names\n  without asking the owner", "stop that task and build nothing more of it", "in `grove.yaml`, from one correction, quoting\n  it"}
	alone := []string{"ask no one and wait for no one's answer", "checkpoint for a later session and finish", "naming the revision the change needs", "Make the choices and change the paths a task's `Delegated:` line names.\n", "mark that task `Status: blocked`, naming the file", "whatever its\n  `needs_human` says", "Do not add it to `grove.yaml`."}
	has := func(what, out string, want, not []string) {
		t.Helper()
		for _, s := range want {
			if !strings.Contains(out, s) {
				t.Errorf("%s lack %q:\n%.4000s", what, s, out)
			}
		}
		for _, s := range not {
			if strings.Contains(out, s) {
				t.Errorf("%s have %q", what, s)
			}
		}
	}
	out := grove(t, "instructions", "implement", cid)
	listed := "\n- `grove/changes/" + cid + "-allow-editing-older-records/specs/record-editing/spec.md`\n"
	has("implement instructions", out, append(append([]string{"You are implementing approved change " + cid + ", \"Allow editing older records\", as main@", "with the owner, in their own working tree.", "### implementation: ", listed}, shared...), owner...),
		append(append(alone, gone...), "### Requirement: Edit older records", ": differs from main", "### proposal: ", "### change-spec: ", "### tasks: ", "/.grove.yaml\">"))
	// A document edited on the branch is named, with the approved text.
	tasks := filepath.Join(dir, "grove/changes", cid+"-allow-editing-older-records/tasks.md")
	os.WriteFile(tasks, []byte(read(t, tasks)+"\nResult: edited.\n"), 0o644)
	main := strings.TrimSpace(git(t, dir, "rev-parse", "main"))
	has("implement after an edit", grove(t, "instructions", "implement", cid, "--unattended"), []string{"\n- `grove/changes/" + cid + "-allow-editing-older-records/tasks.md`: differs from main; `git show " + main + ":grove/changes/" + cid + "-allow-editing-older-records/tasks.md` prints the approved text\n", listed}, nil)
	git(t, dir, "checkout", "-q", "--", ".")
	out = grove(t, "instructions", "implement", cid, "--unattended")
	has("unattended instructions", out, append(append([]string{"asking no one, in your own working tree.", "## Project context\n\nRun python3 -m unittest before saying a task is done.\n\n## What Grove requires of implement\n"}, shared...), alone...), append(owner, gone...))

	// A member of a cycle, not marked, that needs a person, gets its plan.
	planCycle(t, dir, k, "2026-10-22T09:00:00Z", false, cid)
	cycle := read(t, filepath.Join(dir, "grove/cycles", k+"-cycle.yaml")) + "needs_human:\n  " + cid + ": the owner judges it\n"
	os.WriteFile(filepath.Join(dir, "grove/cycles", k+"-cycle.yaml"), []byte(cycle), 0o644)
	commitAll(t, dir, "Needs a person")
	for _, args := range [][]string{{cid}, {cid, "--unattended"}} {
		out := grove(t, append([]string{"instructions", "implement"}, args...)...)
		has(strings.Join(args, " "), out, []string{", change 1 of 1 in cycle " + k + ".", "## The cycle\n\nCycle " + k + ", \"Cycle " + k + "\", in build order:\n\n1. " + cid + " \"Allow editing older records\" — this change\n", "<file path=\"grove/cycles/" + k + "-cycle.yaml\">\n" + cycle}, nil)
	}
	// Both modes build each member in a fresh context, unattended, and only
	// the session archives and commits.
	contexts := []string{"in a fresh context of your own, such as a\n  subagent", "`instructions implement <change-id> --unattended` and follow what it\n  prints, but neither archive nor commit",
		"Start at once every change whose `after` changes are\n  archived and whose tasks' `Files:`", "Only you archive and commit, one returned change at a time", "run the\n  project's tests once both have returned",
		"build the other again on its own", "put the checkpoints each returned change holds to the\n  owner at once", "Unattended, leave checkpoints\n  for your report", "Where you cannot start a fresh context, build each change yourself"}
	out = grove(t, "instructions", "implement", k)
	has("cycle instructions", out, append([]string{"You are implementing cycle " + k + ", \"Cycle " + k + "\", as main@", "with the owner", "run Grove's `instructions implement <change-id>` for each", "\n1. " + cid + " \"Allow editing older records\"\n", "## The cycle\n\n<file path=\"grove/cycles/" + k + "-cycle.yaml\">\n", shared[3], "Finish: report the checkpoints of every change first"}, contexts...), alone)
	out = grove(t, "instructions", "implement", k, "--unattended")
	has("unattended cycle instructions", out, append(append([]string{"asking no one", "run Grove's `instructions implement <change-id> --unattended` for each", "\n1. " + cid + " \"Allow editing older records\"\n", shared[3], "Finish: report the checkpoints of every change first"}, contexts...), alone...), owner)
}

func TestDropChange(t *testing.T) {
	dir := setUp(t)
	cid, _ := approve(t, dir, "older-edits", "Allow editing older records", "")
	folder := "grove/changes/" + cid + "-allow-editing-older-records"
	os.WriteFile(filepath.Join(dir, folder, "implementation.md"), []byte("# Implementation\n\nStatus: dropped\n\nOwner: \"not needed\".\n"), 0o644)
	grove(t, "archive", cid)
}

func TestArchive(t *testing.T) {
	dir := setUp(t)
	cid, _ := approve(t, dir, "older-edits", "Allow editing older records", "")
	git(t, dir, "switch", "-q", "-c", "older-edits-impl")
	folder := "grove/changes/" + cid + "-allow-editing-older-records"
	var stdout, stderr strings.Builder
	// A checkpoint's guess is not done until the owner answers it.
	for _, status := range []string{"in progress", "to check"} {
		os.WriteFile(filepath.Join(dir, folder, "implementation.md"), []byte("## T-01\n\nStatus: "+status+"\n"), 0o644)
		stderr.Reset()
		if code := run([]string{"archive", cid}, &stdout, &stderr); code != 1 || stderr.String() != "grove: T-01 is not done in implementation.md; archive a change when every task is done or the file says Status: dropped\n" {
			t.Errorf("archiving a change whose T-01 is %s: status %d, %q", status, code, stderr.String())
		}
	}
	os.WriteFile(filepath.Join(dir, folder, "implementation.md"), []byte("## T-01\n\nStatus: done\n"), 0o644)

	// The session reconciles the change spec into a new maintained spec,
	// from the template its instructions give, then archives the change.
	out := grove(t, "instructions", "archive", cid)
	for _, s := range []string{"Reconcile the specs", "word for word", "run `archive` again", "marks `unverified`; it does\n  not stop the archive", "Status: done", "\n- `" + folder + "/implementation.md`\n", "\n- `" + folder + "/specs/record-editing/spec.md`\n", "`## History`", "`- [<change-id>](../../changes/archive/<folder>/) <title>`"} {
		if !strings.Contains(out, s) {
			t.Errorf("archive instructions lack %q", s)
		}
	}
	if strings.Contains(out, "### Requirement: Edit older records") {
		t.Errorf("archive instructions print the change's text:\n%.3000s", out)
	}
	os.MkdirAll(filepath.Join(dir, "grove", "specs", "record-editing"), 0o755)
	m := regexp.MustCompile(`(?s)<file path="grove/specs/<capability>/spec.md">\n(.*?)\n</file>`).FindStringSubmatch(out)
	if m == nil || strings.Contains(out, "grove/specs/<capability>/.grove.yaml") || !regexp.MustCompile(`(?s)\n## History\n\n<!--.*?-->$`).MatchString(m[1]) {
		t.Fatalf("no spec.md template ending with History, or a .grove.yaml one, in:\n%s", out)
	}
	purpose := "Editing a record in `records.json` keeps its category.\nA record is one line of the ledger.\n"
	filled := regexp.MustCompile(`(?s)### Requirement: \[.*?\n\n## Known`).ReplaceAllString(m[1], olderEdits+"\n## Known")
	filled = regexp.MustCompile(`(?s)(## Purpose\n\n)<!--.*?-->\n`).ReplaceAllString(filled, "${1}"+purpose)
	filled = regexp.MustCompile(`(?s)(## History\n\n)<!--.*?-->`).ReplaceAllString(filled, "${1}- ["+cid+"](../../changes/archive/"+cid+"-allow-editing-older-records/) Allow editing older records")
	filled = strings.NewReplacer("<capability>", "record-editing", "<title>", "Record editing").Replace(filled)
	os.WriteFile(filepath.Join(dir, "grove", "specs", "record-editing", "spec.md"), []byte(filled), 0o644)
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
		"Maintained specs (agreed behavior):\n- record-editing \"Record editing\": `grove/specs/record-editing/spec.md` — Editing a record in `records.json` keeps its category.\n",
		"Approved changes (future work; implementation is not established):\n- None.",
	} {
		if !strings.Contains(next, s) {
			t.Errorf("fresh explore instructions lack %q", s)
		}
	}
	if len(next) >= 12<<10 {
		t.Errorf("explore instructions are %d bytes, want under 12 KiB", len(next))
	}

	os.WriteFile(filepath.Join(dir, "grove/specs/record-editing/spec.md"), []byte(strings.Replace(filled, "## Purpose\n\n"+purpose, "", 1)), 0o644)
	commitAll(t, dir, "No purpose")
	if out, got := grove(t, "instructions", "explore"), grove(t, "check"); !strings.Contains(out, "`grove/specs/record-editing/spec.md`\n") || got != "No problems in grove/.\n" {
		t.Errorf("a spec with no Purpose: check printed %q, and explore's index:\n%s", got, out)
	}
}

// TestFold follows the walkthrough's "Archive checks the fold": archive
// compares the maintained spec with where the branch left main, the changes
// archived on the branch since, and the change archived.
func TestFold(t *testing.T) {
	dir := setUp(t)
	req := func(name, rule string) string {
		return "### Requirement: " + name + "\n\nWhen a user " + rule + ", the system shall save it.\n\n#### Scenario: " + name + "\n\n- THEN it is saved\n\n"
	}
	keep, show, old := req("Keep a category", "edits a record"), req("Show a record", "asks for a record"), req("Old rule", "types")
	maintain(t, dir, keep+show+old)
	commitAll(t, dir, "Record editing")
	changeSpec := func(branch, title, body string) (string, string) {
		cid, folder := propose(t, dir, branch, title, "")
		os.WriteFile(filepath.Join(folder, "specs/record-editing/spec.md"), []byte("# Changes to record-editing\n\n"+body), 0o644)
		commitAll(t, dir, "Propose "+title)
		git(t, dir, "switch", "-q", "main")
		git(t, dir, "merge", "-q", "--no-ff", "-m", "Approve "+cid, branch)
		os.WriteFile(filepath.Join(folder, "implementation.md"), []byte("## T-01\n\nStatus: done\n"), 0o644)
		return cid, folder
	}
	keep2 := req("Keep a category", "edits an older record")
	a, aDir := changeSpec("a", "Keep older categories", "## MODIFIED Requirements\n\n"+keep2)
	del := req("Delete a record", "deletes a record")
	c, _ := changeSpec("c", "Delete records", "## ADDED Requirements\n\n"+del+"## REMOVED Requirements\n\n### Requirement: Old rule\n\nNo one types.\n")
	d, dDir := changeSpec("d", "Show more", "## MODIFIED Requirements\n\n"+req("Show a record", "asks for more"))
	git(t, dir, "switch", "-q", "-c", "impl")
	git(t, dir, "switch", "-q", "main")
	maintain(t, dir, keep+req("Show a record", "asks for a saved record")+old)
	commitAll(t, dir, "The target moves")
	git(t, dir, "switch", "-q", "impl")

	refuses := func(id, want string) {
		t.Helper()
		var stdout, stderr strings.Builder
		if code := run([]string{"archive", id}, &stdout, &stderr); code != 1 || stdout.String() != want || !strings.HasSuffix(stderr.String(), "so "+id+" stays where it is; fix the maintained specs and run grove archive again\n") {
			t.Errorf("archive %s: status %d\n%s%s", id, code, stdout.String(), stderr.String())
		}
	}
	maintain(t, dir, strings.Replace(keep2, "older", "old", 1)+req("Show a record", "asks")+old)
	refuses(a, "grove/specs/record-editing/spec.md:5: requirement \"Keep a category\" is not as "+a+" gives it\n"+
		"grove/specs/record-editing/spec.md:13: requirement \"Show a record\" changed though no change names it\n")
	maintain(t, dir, old+strings.ReplaceAll(keep2, "edits an ", "edits\nan ")+show)
	grove(t, "archive", a)

	maintain(t, dir, keep2+show+old+del)
	refuses(c, "grove/specs/record-editing/spec.md:21: requirement \"Old rule\" is still here though "+c+" removes it\n")
	maintain(t, dir, keep2+show+del)
	grove(t, "archive", c)

	os.WriteFile(filepath.Join(dDir, "implementation.md"), []byte("# Implementation\n\nStatus: dropped\n"), 0o644)
	grove(t, "archive", d)
	if _, err := os.Stat(filepath.Join(dir, "grove/changes/archive", filepath.Base(aDir))); err != nil {
		t.Error(err)
	}
}

// TestSchemas follows the walkthrough's "A bug": a project schema chosen
// for a new change, and kept by the change.
func TestSchemas(t *testing.T) {
	dir := setUp(t)
	// Grove's bugfix, with no schemas of the project's.
	if out := grove(t, "instructions", "propose", "Editing a missing record crashes"); !strings.Contains(out, "- `default`: Software changes") ||
		!strings.Contains(out, "\n- `bugfix`: A fix for behavior that is wrong, whether or not a maintained spec describes it.\n") {
		t.Errorf("propose instructions do not list Grove's schemas:\n%s", out)
	}
	if out := grove(t, "instructions", "propose", "Editing a missing record crashes"); !strings.Contains(out, "describes behavior that is wrong, suggest the `bugfix` change schema") || !strings.Contains(out, "Once the owner confirms it") {
		t.Errorf("propose instructions do not suggest bugfix:\n%.3000s", out)
	}
	if out := grove(t, "instructions", "explore"); !strings.Contains(out, "\n- `bugfix`: A fix for behavior that is wrong") || !strings.Contains(out, "describes behavior that is wrong, suggest the `bugfix`\n  change schema") || !strings.Contains(out, "If the owner has\n  confirmed `bugfix`") {
		t.Errorf("explore instructions do not suggest bugfix:\n%.3000s", out)
	}
	own := func(name, content string) {
		os.MkdirAll(filepath.Dir(filepath.Join(dir, "grove/schemas", name)), 0o755)
		os.WriteFile(filepath.Join(dir, "grove/schemas", name), []byte(content), 0o644)
	}
	own("default/documents/tasks/template.md", "## T-01 — [title]\n\nRollback: [how to undo this task]\n")
	out := grove(t, "instructions", "propose", "--schema", "bugfix", "x")
	for _, s := range []string{"## How to propose\n\n**Propose a bug fix: \"This is wrong; make it right.\"**", "fails before the\n  fix", "Write no change spec",
		"## Observed\n", "## Expected\n", "## Reproduction\n", "## The requirement it breaks\n", "## Open questions\n", "Rollback: [how to undo this task]", "schema: bugfix\n"} {
		if !strings.Contains(out, s) {
			t.Errorf("Grove's bugfix propose instructions lack %q:\n%.3000s", s, out)
		}
	}
	own("bugfix/documents/proposal/template.md", "# <title>\n\n## Seen\n")
	if out := grove(t, "instructions", "propose", "--schema", "bugfix", "x"); !strings.Contains(out, "/proposal.md\">\n# <title>\n\n## Seen\n</file>") || !strings.Contains(out, "**Propose a bug fix: \"This is wrong") {
		t.Errorf("bugfix with the project's proposal template:\n%.3000s", out)
	}
	bug, _ := propose(t, dir, "bug", "Editing a missing record prints an error", "schema: bugfix\n")
	if got := grove(t, "check"); got != "No problems in grove/.\n" {
		t.Errorf("check with change %s of schema bugfix printed %q", bug, got)
	}
	git(t, dir, "checkout", "-q", "--", ".")
	git(t, dir, "clean", "-qfd")
	git(t, dir, "switch", "-q", "main")
	os.RemoveAll(filepath.Join(dir, "grove/schemas"))

	for name, content := range map[string]string{
		"bugfix/schema.yaml":                    "description: A fix for behavior that is wrong.\n",
		"bugfix/skills/propose.md":              "**Propose a bug fix.**\n",
		"bugfix/documents/proposal/template.md": "# <title>\n\n## Observed\n",
		"short/schema.yaml":                     "description: Short.\nskills:\n  propose:\n    writes: [proposal, tasks]\n",
	} {
		os.MkdirAll(filepath.Dir(filepath.Join(dir, "grove/schemas", name)), 0o755)
		os.WriteFile(filepath.Join(dir, "grove/schemas", name), []byte(content), 0o644)
	}
	if out := grove(t, "instructions", "propose", "Editing a missing record crashes"); !strings.Contains(out, "these instructions follow `default`.\n\n- `default`: Software changes") ||
		!strings.Contains(out, "- `bugfix`: A fix for behavior that is wrong.\n- `short`: Short.\n") {
		t.Errorf("propose instructions do not list the schemas:\n%s", out)
	}
	out = grove(t, "instructions", "propose", "--schema", "bugfix", "Editing a missing record crashes")
	for _, s := range []string{", with change schema bugfix, from the owner's description", "## How to propose\n\n**Propose a bug fix.**", "/proposal.md\">\n# <title>\n\n## Observed\n",
		"### change-spec: ", "### tasks: ", "/.grove.yaml\">\ntitle: \"<title>\"\ncreated: ", "schema: bugfix\n"} {
		if !strings.Contains(out, s) {
			t.Errorf("bugfix propose instructions lack %q", s)
		}
	}
	if out := grove(t, "instructions", "propose", "--schema", "short"); strings.Contains(out, "### change-spec: ") || !strings.Contains(out, "**Propose says") {
		t.Errorf("short propose instructions:\n%s", out)
	}
	var stdout, stderr strings.Builder
	if code := run([]string{"instructions", "propose", "--schema", "bugfx"}, &stdout, &stderr); code != 1 || stderr.String() != "grove: schema \"bugfx\" is not a change schema; schemas are default, bugfix and short\n" {
		t.Errorf("an unknown schema: status %d, %q", code, stderr.String())
	}

	cid, folder := propose(t, dir, "older-edits", "Allow editing older records", "schema: bugfix\n")
	stderr.Reset()
	if code := run([]string{"instructions", "propose", "--schema", "default", cid}, &stdout, &stderr); code != 1 || !strings.Contains(stderr.String(), "a change keeps its schema") {
		t.Errorf("--schema with a change ID: status %d, %q", code, stderr.String())
	}
	if out := grove(t, "instructions", "propose", cid); !strings.Contains(out, "\"Allow editing older records\", with change schema bugfix, in") || !strings.Contains(out, "**Propose a bug fix.**") {
		t.Errorf("bugfix revise instructions:\n%.2000s", out)
	}
	commitAll(t, dir, "Propose")
	if out := grove(t, "review"); !strings.Contains(out, "  adds     change schema    grove/schemas/bugfix/schema.yaml\n") {
		t.Errorf("review of a schema:\n%.1000s", out)
	}
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Approve", "older-edits")
	planCycle(t, dir, "261022-k0000", "2026-10-22T09:00:00Z", false, cid)
	for skill, s := range map[string]string{"implement": "One task at a time", "archive": "Reconcile the specs"} {
		if out := grove(t, append(append([]string{"instructions"}, strings.Fields(skill)...), cid)...); !strings.Contains(out, ", with change schema bugfix, ") || !strings.Contains(out, s) {
			t.Errorf("bugfix %s instructions:\n%.2000s", skill, out)
		}
	}
	if out := grove(t, "instructions", "implement", cid); strings.Contains(out, "### proposal: ") || !strings.Contains(out, "### implementation: ") {
		t.Errorf("bugfix implement instructions print templates other than implementation's:\n%.3000s", out)
	}

	// The walkthrough's "A mistake in a schema", a change naming no schema,
	// and a schema without tasks.
	schemas := filepath.Join(dir, "grove/schemas")
	os.Rename(filepath.Join(schemas, "bugfix/documents/proposal"), filepath.Join(schemas, "bugfix/documents/propsal"))
	os.WriteFile(filepath.Join(schemas, "bugfix/schema.yaml"), []byte("description: A fix for behavior that is wrong.\ninstrution: x\n"), 0o644)
	os.WriteFile(filepath.Join(folder, ".grove.yaml"), []byte("title: Allow editing older records\ncreated: 2026-10-14T00:00:00Z\nschema: bugfx\n"), 0o644)
	os.MkdirAll(filepath.Join(schemas, "notasks"), 0o755)
	os.WriteFile(filepath.Join(schemas, "notasks/schema.yaml"), []byte("description: N.\ndocuments:\n  - {kind: proposal, path: changes/<change>/proposal.md, template: documents/proposal/template.md, guide: documents/proposal/guide.md, metadata: changes/<change>/.grove.yaml}\n  - {kind: spec, path: specs/<capability>/spec.md, template: t.md, guide: g.md, metadata: specs/<capability>/.grove.yaml}\n"), 0o644)
	stdout.Reset()
	stderr.Reset()
	code := run([]string{"check"}, &stdout, &stderr)
	for _, s := range []string{
		"grove/schemas/bugfix/schema.yaml:2: unknown field \"instrution\"; schema.yaml takes description, documents and skills\n",
		"grove/schemas/bugfix/documents/propsal/template.md: schema bugfix does not use this file\n",
		"grove/schemas/notasks/schema.yaml: tasks must be at changes/<change>/tasks.md; Grove reads it to archive\n",
		"grove/schemas/notasks/schema.yaml: spec: metadata must be changes/<change>/.grove.yaml\n",
		"grove/changes/" + cid + "-allow-editing-older-records/.grove.yaml:3: schema \"bugfx\" is not a change schema; schemas are default, bugfix, notasks and short\n",
	} {
		if code != 1 || !strings.Contains(stdout.String(), s) {
			t.Errorf("check lacks %q: status %d\n%s", s, code, stdout.String())
		}
	}
	grove(t, "changes")
	grove(t, "instructions", "propose", "--schema", "bugfix") // an unknown field or unused file stops no skill
	stderr.Reset()
	if code := run([]string{"instructions", "propose", "--schema", "notasks"}, &stdout, &stderr); code != 1 || !strings.Contains(stderr.String(), "grove: grove/schemas/notasks/schema.yaml: ") {
		t.Errorf("propose with a broken schema: status %d, %q", code, stderr.String())
	}
}

// finish plays the sessions that apply and archive change id on a new
// branch, archiving writing spec as the record-editing maintained spec, then
// merges it. It returns what check printed on the branch, before archive
// moved the change.
func finish(t *testing.T, dir, branch, id, folder, spec string) string {
	t.Helper()
	git(t, dir, "switch", "-q", "-c", branch)
	os.WriteFile(filepath.Join(folder, "implementation.md"), []byte("## T-01\n\nStatus: done\n"), 0o644)
	maintain(t, dir, spec)
	var stdout, stderr strings.Builder
	run([]string{"check"}, &stdout, &stderr)
	grove(t, "archive", id)
	commitAll(t, dir, "Implement and archive "+id)
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Merge "+branch, branch)
	return stdout.String()
}

// maintain writes reqs as the requirements of the record-editing
// maintained spec, as archive's fold does.
func maintain(t *testing.T, dir, reqs string) {
	t.Helper()
	os.MkdirAll(filepath.Join(dir, "grove/specs/record-editing"), 0o755)
	os.WriteFile(filepath.Join(dir, "grove/specs/record-editing/spec.md"), []byte("# Record editing\n\n## Requirements\n\n"+reqs), 0o644)
}

// TestCoordinating follows the walkthrough's "Coordinating changes".
func TestCoordinating(t *testing.T) {
	dir := setUp(t)
	keep := "### Requirement: Keep a category\n\nWhen a user edits a record and leaves its category unchanged, the system shall save the edit.\n\n#### Scenario: A note\n"
	maintain(t, dir, keep)
	commitAll(t, dir, "Record editing")
	unknown, unknownDir := approve(t, dir, "unknown-id", "Editing an unknown ID prints an error", "")

	// Propose records depends_on and base; lists show what waits.
	main := strings.TrimSpace(git(t, dir, "rev-parse", "main"))
	out := grove(t, "instructions", "propose", "Add ledger delete ID")
	for _, s := range []string{"ask the owner whether it needs another open change", "\nbase: " + main + "\n",
		"- " + unknown + " \"Editing an unknown ID prints an error\": `grove/changes/" + unknown + "-editing-an-unknown-id-prints-an-error/`\n  - capabilities: record-editing\n"} {
		if !strings.Contains(out, s) {
			t.Errorf("propose instructions lack %q:\n%s", s, out)
		}
	}
	del, delDir := approve(t, dir, "delete", "Add ledger delete ID", "depends_on: ["+unknown+"]\nbase: "+main+"\n")
	if got := grove(t, "changes"); !strings.Contains(got, unknown+"  Editing an unknown ID prints an error  approved\n") || !strings.Contains(got, del+"  Add ledger delete ID  approved  waits on "+unknown+"\n") {
		t.Errorf("changes printed %q", got)
	}
	if out := grove(t, "instructions", "explore"); !strings.Contains(out, "  - waits on: "+unknown+", not archived as finished\n") {
		t.Errorf("explore's index does not say what delete waits on:\n%s", out)
	}

	// Implement refuses an unmet dependency, with the owner or not.
	k := "261022-k0000"
	planCycle(t, dir, k, "2026-10-22T09:00:00Z", false, del, unknown)
	var stdout, stderr strings.Builder
	for _, args := range [][]string{{"implement", del, "--unattended"}, {"implement", del}} {
		stderr.Reset()
		if code := run(append([]string{"instructions"}, args...), &stdout, &stderr); code != 1 || stderr.String() != "grove: "+del+" waits on "+unknown+", which this working tree has not archived as finished\n" {
			t.Errorf("%v before its dependency is archived: status %d, %q", args, code, stderr.String())
		}
	}

	// Implementing a cycle unattended lists its changes in dependency order,
	// each after the members it depends on.
	out = grove(t, "instructions", "implement", k, "--unattended")
	for _, s := range []string{"## What Grove requires of implement\n\n- Work in this working tree on a branch other than the target",
		"## How to implement\n\nThe changes in build order", "\n1. " + unknown + " \"Editing an unknown ID prints an error\"\n2. " + del + " \"Add ledger delete ID\" after " + unknown + "\n"} {
		if !strings.Contains(out, s) {
			t.Errorf("cycle instructions lack %q:\n%s", s, out)
		}
	}

	// A dependency archived on this branch, not yet on main, is met.
	git(t, dir, "switch", "-q", "-c", "batch")
	os.WriteFile(filepath.Join(unknownDir, "implementation.md"), []byte("## T-01\n\nStatus: done\n"), 0o644)
	maintain(t, dir, keep+olderEdits)
	grove(t, "archive", unknown)
	if out := grove(t, "instructions", "implement", del, "--unattended"); !strings.Contains(out, "You are implementing approved change "+del) || !strings.Contains(out, ", change 2 of 2 in cycle "+k+".") ||
		!strings.Contains(out, "1. "+unknown+" \"Editing an unknown ID prints an error\" — archived in this working tree\n2. "+del+" \"Add ledger delete ID\" — this change\n") {
		t.Errorf("implement after its dependency is archived on the branch:\n%s", out)
	}
	if out := grove(t, "instructions", "implement", k, "--unattended"); !strings.Contains(out, "\n1. "+del+" \"Add ledger delete ID\"\n\n## Project index") {
		t.Errorf("cycle after archiving on the branch:\n%s", out)
	}
	git(t, dir, "checkout", "-q", "--", ".")
	git(t, dir, "clean", "-qfd")
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "branch", "-q", "-D", "batch")
	if out := finish(t, dir, "unknown-id-impl", unknown, unknownDir, keep+olderEdits); out != "No problems in grove/.\n" {
		t.Errorf("check on unknown-id-impl printed %q", out)
	}
	grove(t, "instructions", "implement", del, "--unattended")

	// check reports broken dependencies (Other paths).
	budgets, _ := propose(t, dir, "budgets", "Monthly budgets in report", "depends_on: ["+del+"]\n")
	delMeta := filepath.Join(delDir, ".grove.yaml")
	os.WriteFile(delMeta, []byte("title: Add ledger delete ID\ncreated: 2026-10-14T00:00:00Z\ndepends_on: [260101-zzzzz, "+budgets+"]\n"), 0o644)
	unknownImpl := filepath.Join(dir, "grove/changes/archive", filepath.Base(unknownDir), "implementation.md")
	os.WriteFile(unknownImpl, []byte("Status: dropped\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "grove/changes", budgets+"-monthly-budgets-in-report/.grove.yaml"), []byte("title: Monthly budgets in report\ncreated: 2026-10-14T00:00:00Z\ndepends_on: ["+del+", "+unknown+"]\n"), 0o644)
	stdout.Reset()
	stderr.Reset()
	code := run([]string{"check"}, &stdout, &stderr)
	for _, s := range []string{
		"grove/changes/" + filepath.Base(delDir) + "/.grove.yaml:3: depends_on: 260101-zzzzz is no change, open or archived\n",
		"grove/changes/" + filepath.Base(delDir) + "/.grove.yaml:3: depends_on: " + budgets + " waits on this change in turn, so neither can be applied\n",
		"grove/changes/" + budgets + "-monthly-budgets-in-report/.grove.yaml:3: depends_on: " + del + " waits on this change in turn, so neither can be applied\n",
		"grove/changes/" + budgets + "-monthly-budgets-in-report/.grove.yaml:3: depends_on: " + unknown + " was dropped, so this change could never be applied\n",
	} {
		if code != 1 || !strings.Contains(stdout.String(), s) {
			t.Errorf("check lacks %q: status %d\n%s", s, code, stdout.String())
		}
	}
	grove(t, "changes")
	git(t, dir, "checkout", "--", ".")

	// A change that waits on nothing.
	budgetsMeta := filepath.Join(dir, "grove/changes", budgets+"-monthly-budgets-in-report/.grove.yaml")
	os.WriteFile(budgetsMeta, []byte("title: Monthly budgets in report\ncreated: 2026-10-14T00:00:00Z\n"), 0o644)
	commitAll(t, dir, "Propose monthly budgets")
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Approve budgets", "budgets")
	if got := grove(t, "changes"); got != del+"  Add ledger delete ID  approved  in cycle "+k+"\n"+budgets+"  Monthly budgets in report  approved\n" && got != budgets+"  Monthly budgets in report  approved\n"+del+"  Add ledger delete ID  approved  in cycle "+k+"\n" {
		t.Errorf("changes printed %q", got)
	}

	// The stale-plan check.
	modify := func(folder, text string) {
		os.WriteFile(filepath.Join(folder, "specs/record-editing/spec.md"), []byte("# Changes to record-editing\n\n## MODIFIED Requirements\n\n### Requirement: Keep a category\n\n"+text), 0o644)
	}
	base := strings.TrimSpace(git(t, dir, "rev-parse", "main"))
	recat, recatDir := propose(t, dir, "recategorize", "Recategorize older records", "base: "+base+"\n")
	recatText := "When a user edits a record and leaves its category unchanged, the system shall save the edit.\n\n#### Scenario: Change an older record's category to a current one\n"
	modify(recatDir, recatText)
	commitAll(t, dir, "Propose recategorizing")
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Approve recategorizing", "recategorize")
	show, showDir := propose(t, dir, "show-saved", "Show what was saved", "base: "+base+"\n")
	modify(showDir, "When a user edits a record and leaves its category unchanged, the system shall save the edit and print what was saved.\n\n#### Scenario: A note\n")
	commitAll(t, dir, "Propose showing what was saved")
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Approve showing what was saved", "show-saved")
	want := "grove/changes/" + filepath.Base(showDir) + "/specs/record-editing/spec.md:5: requirement \"Keep a category\" was rewritten since base " + base[:7] + "; revise the change with propose\n"
	if out := finish(t, dir, "recategorize-impl", recat, recatDir, olderEdits+"### Requirement: Keep a category\n\n"+recatText); out != want {
		t.Errorf("check while archiving recategorize printed %q", out)
	}
	stdout.Reset()
	if code := run([]string{"check"}, &stdout, &stderr); code != 1 || stdout.String() != want {
		t.Errorf("check of a stale change: status %d\n%s", code, stdout.String())
	}
	// Had show-saved depended on recategorize, its text would be expected.
	os.WriteFile(filepath.Join(showDir, ".grove.yaml"), []byte("title: Show what was saved\ncreated: 2026-10-14T00:00:00Z\ndepends_on: ["+recat+"]\nbase: "+base+"\n"), 0o644)
	if got := grove(t, "check"); got != "No problems in grove/.\n" {
		t.Errorf("check with depends_on on the rewriting change printed %q", got)
	}
	git(t, dir, "checkout", "-q", "--", ".")
	main = strings.TrimSpace(git(t, dir, "rev-parse", "main"))
	git(t, dir, "switch", "-q", "-c", "show-saved-revise")
	if out := grove(t, "instructions", "propose", show); !strings.Contains(out, "\nbase: "+main+"\n") || !strings.Contains(out, "never clear the report by hand") {
		t.Errorf("revise instructions for a stale change:\n%s", out)
	}
	modify(showDir, strings.Replace(recatText, "save the edit", "save the edit and print what was saved", 1))
	os.WriteFile(filepath.Join(showDir, ".grove.yaml"), []byte("title: Show what was saved\ncreated: 2026-10-14T00:00:00Z\nbase: "+main+"\n"), 0o644)
	if got := grove(t, "check"); got != "No problems in grove/.\n" {
		t.Errorf("check after revising printed %q", got)
	}

	// grove changes --json.
	var v struct {
		Commit  string
		Changes []struct {
			ID, State, Path string
			WaitsOn         []string `json:"waits_on"`
			Tasks           []struct{ ID, Status string }
		}
	}
	if err := json.Unmarshal([]byte(grove(t, "changes", "--json")), &v); err != nil {
		t.Fatal(err)
	}
	states := map[string]string{}
	for _, c := range v.Changes {
		states[c.ID] = c.State
		if c.ID == del && (len(c.WaitsOn) > 0 || c.Path != "changes/"+filepath.Base(delDir)+"/proposal.md") || c.ID == unknown && (len(c.Tasks) != 1 || c.Tasks[0].Status != "done") {
			t.Errorf("changes --json gives %s as %+v", c.ID, c)
		}
	}
	if v.Commit != main || states[unknown] != "archived" || states[recat] != "archived" || states[del] != "approved" {
		t.Errorf("changes --json: %+v", v)
	}

	// Implementing a cycle leaves out a change that waits on a change that is
	// neither being built nor finished.
	imports, _ := approve(t, dir, "imports", "Import a CSV", "depends_on: [260101-zzzzz]\n")
	planCycle(t, dir, "261022-m0000", "2026-10-22T10:00:00Z", false, imports, budgets, show)
	if out := grove(t, "instructions", "implement", "261022-m0000", "--unattended"); !strings.Contains(out, "\nLeave out "+imports+", \"Import a CSV\": it waits on 260101-zzzzz, which is neither being built nor archived as finished.") || !strings.Contains(out, " "+budgets+" ") || !strings.Contains(out, " "+show+" ") {
		t.Errorf("cycle with a change that cannot be built:\n%s", out)
	}
}

// planCycle commits cycle id, created at created, listing changes and
// planned at the current commit, on the current branch of dir, as a merged
// plan; mark adds unattended.
func planCycle(t *testing.T, dir, id, created string, mark bool, changes ...string) {
	t.Helper()
	body := "title: Cycle " + id + "\ncreated: " + created + "\nchanges: [" + strings.Join(changes, ", ") + "]\nplanned: " + strings.TrimSpace(git(t, dir, "rev-parse", "HEAD")) + "\n"
	if mark {
		body += "unattended: true\n"
	}
	os.MkdirAll(filepath.Join(dir, "grove/cycles"), 0o755)
	os.WriteFile(filepath.Join(dir, "grove/cycles", id+"-cycle.yaml"), []byte(body), 0o644)
	commitAll(t, dir, "Plan "+id)
}

// TestCycles follows the walkthrough's "Planning cycles": plan
// instructions, a cycle file merged and marked, and check's reports on it.
func TestCycles(t *testing.T) {
	dir := setUp(t)
	a, aDir := approve(t, dir, "a", "Editing an unknown ID prints an error", "")
	b, _ := approve(t, dir, "b", "Add ledger delete ID", "")
	main := strings.TrimSpace(git(t, dir, "rev-parse", "main"))
	out := grove(t, "instructions", "plan", "0.4 release")
	m := regexp.MustCompile("### cycle: `grove/cycles/([0-9]{6}-[0-9a-z]{5})-<slug>.yaml`").FindStringSubmatch(out)
	for _, s := range []string{"You are planning new cycle ", "## What Grove requires of plan\n", "Never write or remove `unattended`", "## How to plan\n",
		"\ntitle: \"<title>\"\ncreated: 20", "\nplanned: " + main + "\n", "## Owner's description\n\n0.4 release\n",
		"\n# depends_on: [", "ask the owner\n  whether it follows any of them",
		"A new cycle: before writing its file, run Grove's `check`.", "tell the owner this edits approved work, bring the\n  member up to date on this branch as Grove's\n  `instructions propose <change-id>` says", "`grove run` builds\n    nothing of the cycle until it is"} {
		if m == nil || !strings.Contains(out, s) {
			t.Errorf("plan instructions lack %q:\n%s", s, out)
		}
	}
	if strings.Contains(out, "the only thing you write") {
		t.Errorf("plan instructions still say the cycle file is the only thing written")
	}
	// The finish names implement, then run given the cycle, then the mark.
	if i, r, u := strings.Index(out, "say `/grove-implement <cycle-id>`"), strings.Index(out, "say whether `grove run <cycle-id>`"), strings.Index(out, "say `unattended: true`"); i < 0 || i > r || r > u {
		t.Errorf("plan's finish is not in order:\n%s", out)
	}
	if out := grove(t, "instructions", "plan", a); !strings.Contains(out, "with approved change "+a+", \"Editing an unknown ID prints an error\", as its only member") {
		t.Errorf("plan of one change:\n%.1000s", out)
	}

	// The plan merged, then marked in a merge of its own.
	k := m[1]
	cycle := filepath.Join(dir, "grove/cycles", k+"-0-4-release.yaml")
	plan := "title: 0.4 release\ncreated: 2026-10-22T10:00:00Z\nchanges: [" + a + "]\nplanned: " + main + "\n"
	git(t, dir, "switch", "-q", "-c", "plan")
	os.MkdirAll(filepath.Dir(cycle), 0o755)
	os.WriteFile(cycle, []byte(plan), 0o644)
	if got := grove(t, "check"); got != "No problems in grove/.\n" {
		t.Errorf("check of the plan printed %q", got)
	}
	commitAll(t, dir, "Plan 0.4")
	if out := grove(t, "review"); !strings.Contains(out, "  adds     cycle            grove/cycles/"+k+"-0-4-release.yaml\n") {
		t.Errorf("review of a cycle:\n%s", out)
	}
	if out := grove(t, "instructions", "plan", k); !strings.Contains(out, "You are revising cycle "+k+", \"0.4 release\"") || !strings.Contains(out, "## The cycle\n\n<file path=\"grove/cycles/"+k+"-0-4-release.yaml\">\n"+plan) {
		t.Errorf("revise instructions for a cycle:\n%s", out)
	}
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Plan 0.4", "plan")
	if got := grove(t, "changes"); !strings.Contains(got, a+"  Editing an unknown ID prints an error  approved  in cycle "+k+"\n") || !strings.Contains(got, b+"  Add ledger delete ID  approved\n") {
		t.Errorf("changes printed %q", got)
	}
	if out := grove(t, "instructions", "explore"); !strings.Contains(out, "- "+k+" \"0.4 release\": `grove/cycles/"+k+"-0-4-release.yaml`\n  - changes: "+a+"\n") {
		t.Errorf("the project index lacks the cycle:\n%s", out)
	}
	os.WriteFile(cycle, []byte(plan+"unattended: true\n"), 0o644)
	commitAll(t, dir, "Mark 0.4 unattended")
	if got := grove(t, "changes"); !strings.Contains(got, "  in cycle "+k+", unattended\n") {
		t.Errorf("changes of a marked cycle printed %q", got)
	}
	// Revising a member re-plans its cycle on the same branch; revising a
	// change in no cycle does not.
	marked, rel := strings.TrimSpace(git(t, dir, "rev-parse", "main")), "grove/cycles/"+k+"-0-4-release.yaml"
	revise := grove(t, "instructions", "propose", a)
	for _, s := range []string{", and re-planning cycle " + k + ", \"0.4 release\", which lists it.", "### cycle: `" + rel + "`", "\nplanned: " + marked + "\n",
		"never write or remove `unattended`", "## The cycle\n\n<file path=\"" + rel + "\">\n" + plan + "unattended: true\n</file>"} {
		if !strings.Contains(revise, s) {
			t.Errorf("revising a member lacks %q:\n%.4000s", s, revise)
		}
	}
	if out := grove(t, "instructions", "propose", b); strings.Contains(out, "## The cycle") || strings.Contains(out, "### cycle: ") {
		t.Errorf("revising a change in no cycle prints a cycle:\n%.3000s", out)
	}
	if out := grove(t, "instructions", "plan"); !strings.Contains(out, "that re-plans the cycle on\n    the same branch, and one merge removes the obstacle") {
		t.Errorf("plan's finish does not name the shorter way:\n%.3000s", out)
	}
	checks := func(want string) {
		t.Helper()
		var stdout, stderr strings.Builder
		if code := run([]string{"check"}, &stdout, &stderr); code != 1 || stdout.String() != want {
			t.Errorf("check: status %d\n%s%s\nwant\n%s", code, stdout.String(), stderr.String(), want)
		}
	}
	file := "grove/cycles/" + k + "-0-4-release.yaml"

	// A marked cycle with a change that needs a person.
	os.WriteFile(cycle, []byte(plan+"needs_human:\n  "+a+": judging the error's wording\nunattended: true\n"), 0o644)
	checks(file + ":7: unattended: " + a + " needs a person, as needs_human says, so grove claim skips this cycle; remove unattended, or build the cycle with /grove-implement " + k + "\n")
	git(t, dir, "checkout", "-q", "--", ".")

	// A marked cycle whose member changed on main since its plan.
	os.WriteFile(filepath.Join(aDir, "proposal.md"), []byte("# Editing an unknown ID prints an error\n\nRevised.\n"), 0o644)
	commitAll(t, dir, "Revise "+a)
	checks(file + ": planned: " + a + " changed on main since this plan, so grove claim skips this cycle; re-plan it with /grove-plan " + k + "\n")

	// A malformed cycle, and one change in two cycles not done.
	os.WriteFile(cycle, []byte(strings.Replace(plan, "planned: "+main, "planned: "+strings.TrimSpace(git(t, dir, "rev-parse", "main")), 1)), 0o644)
	os.WriteFile(filepath.Join(dir, "grove/cycles/261099-zzzzz-later.yaml"), []byte("title: Later\ncreated: 2026-10-23T09:00:00Z\nchanges: ["+a+", 260101-zzzzz]\nstart: tomorrow\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "grove/cycles/261099-zzzzy-bare.yaml"), []byte("title: Bare\n"), 0o644)
	checks("grove/cycles/261099-zzzzy-bare.yaml: created: required field is missing\n" +
		"grove/cycles/261099-zzzzy-bare.yaml: changes: required field is missing\n" +
		"grove/cycles/261099-zzzzy-bare.yaml: planned: required field is missing\n" +
		"grove/cycles/261099-zzzzz-later.yaml:4: unknown field \"start\"; this file takes title, created, changes, planned, goal, depends_on, needs_human and unattended\n" +
		"grove/cycles/261099-zzzzz-later.yaml: planned: required field is missing\n" +
		"grove/cycles/261099-zzzzz-later.yaml:3: changes: " + a + " is also in cycle " + k + ", which is not done\n" +
		"grove/cycles/261099-zzzzz-later.yaml:3: changes: 260101-zzzzz is no change, open or archived\n")
	git(t, dir, "checkout", "-q", "--", ".")
	git(t, dir, "clean", "-qfd")

	// A cycle that follows another, one that follows no cycle, and two that
	// follow each other.
	y, now := "261099-y0000", strings.TrimSpace(git(t, dir, "rev-parse", "main"))
	fresh, next := strings.Replace(plan, main, now, 1), filepath.Join(dir, "grove/cycles/261099-y0000-next.yaml")
	follow := func(k, y string) {
		os.WriteFile(next, []byte("title: Next\ncreated: 2026-10-23T09:00:00Z\nchanges: ["+b+"]\nplanned: "+now+"\ndepends_on: ["+y+"]\n"), 0o644)
		os.WriteFile(cycle, []byte(fresh+k), 0o644)
	}
	follow("", k)
	if got := grove(t, "check"); got != "No problems in grove/.\n" {
		t.Errorf("check of a cycle that follows another printed %q", got)
	}
	commitAll(t, dir, "Plan next")
	if got := grove(t, "changes"); !strings.Contains(got, b+"  Add ledger delete ID  approved  in cycle "+y+", waits on cycle "+k+"\n") {
		t.Errorf("changes of a cycle that waits on a cycle printed %q", got)
	}
	if out := grove(t, "instructions", "explore"); !strings.Contains(out, "`grove/cycles/261099-y0000-next.yaml`\n  - changes: "+b+"\n  - follows "+k+"\n") {
		t.Errorf("the project index does not name the cycle followed:\n%s", out)
	}
	follow("", "261099-zzzzz")
	checks("grove/cycles/261099-y0000-next.yaml:5: depends_on: 261099-zzzzz is no cycle\n")
	follow("depends_on: ["+y+"]\n", k)
	checks(file + ":5: depends_on: " + y + " waits on this cycle in turn, so neither can run\n" +
		"grove/cycles/261099-y0000-next.yaml:5: depends_on: " + k + " waits on this cycle in turn, so neither can run\n")
	git(t, dir, "checkout", "-q", "--", ".")
	finish(t, dir, "a-impl", a, aDir, olderEdits)
	if got := grove(t, "changes"); !strings.Contains(got, b+"  Add ledger delete ID  approved  in cycle "+y+"\n") {
		t.Errorf("changes once the cycle followed is done printed %q", got)
	}
}

// TestMain runs the test binary as grove when a test starts it so, to
// claim from a second checkout at the same moment as the test's own claim.
func TestMain(m *testing.M) {
	if os.Getenv("GROVE_TEST_MAIN") != "" {
		os.Exit(Main(os.Args[1:]))
	}
	os.Exit(m.Run())
}

// TestClaim follows the walkthrough's "Running unattended" with a bare
// origin and two operators' clones of it, h1 and h2.
func TestClaim(t *testing.T) {
	owner := setUp(t)
	remote, h1, h2 := filepath.Join(t.TempDir(), "remote.git"), filepath.Join(t.TempDir(), "h1"), filepath.Join(t.TempDir(), "h2")
	git(t, owner, "clone", "-q", "--bare", owner, remote)
	git(t, owner, "remote", "add", "origin", remote)
	git(t, owner, "clone", "-q", remote, h1)
	git(t, owner, "clone", "-q", remote, h2)
	for _, k := range []string{"GIT_AUTHOR_NAME", "GIT_COMMITTER_NAME", "GIT_AUTHOR_EMAIL", "GIT_COMMITTER_EMAIL"} {
		t.Setenv(k, "operator@example.com")
	}
	host, _ := os.Hostname()
	recorded := func() int {
		if git(t, remote, "for-each-ref", "refs/heads/grove/runs") == "" {
			return 0
		}
		return strings.Count(git(t, remote, "log", "--format=%H", "grove/runs"), "\n")
	}
	latest := func(dir, ref string) (a attempt) {
		json.Unmarshal([]byte(git(t, dir, "log", "-1", "--format=%B", ref)), &a)
		return
	}
	// claims runs grove claim in checkout, wanting standard output out and
	// standard error matching want; refused, it records nothing and
	// creates no branch.
	claims := func(checkout, out, want string, args ...string) {
		t.Helper()
		t.Chdir(checkout)
		branches, n := git(t, checkout, "for-each-ref", "refs/heads"), recorded()
		var stdout, stderr strings.Builder
		code := run(append([]string{"claim"}, args...), &stdout, &stderr)
		if stdout.String() != out || !regexp.MustCompile(want).MatchString(stderr.String()) || (code == 0) != (out != "") ||
			out == "" && (git(t, checkout, "for-each-ref", "refs/heads") != branches || recorded() != n) {
			t.Errorf("grove claim %v in %s: status %d, %q, %q; want %q, %s", args, filepath.Base(checkout), code, stdout.String(), stderr.String(), out, want)
		}
	}
	if code := run([]string{"run"}, io.Discard, io.Discard); code != 2 {
		t.Errorf("grove run: status %d, want the usage's 2", code)
	}
	nothing := func() {
		t.Helper()
		git(t, owner, "push", "-q", "origin", "main")
		claims(h1, "", `^grove: nothing to claim on main at [0-9a-f]{7}\n$`)
	}
	nothing()

	// Nothing to claim: changes in no cycle, a cycle not marked, a marked
	// cycle gone stale. Claude Code need not be listed under agents.
	t.Chdir(owner)
	config := read(t, filepath.Join(owner, "grove.yaml"))
	os.WriteFile(filepath.Join(owner, "grove.yaml"), []byte(strings.Replace(config, "  claude: skills\n", "", 1)), 0o644)
	commitAll(t, owner, "Codex only")
	a, aDir := approve(t, owner, "a", "Editing an unknown ID prints an error", "")
	b, bDir := approve(t, owner, "b", "Add ledger delete ID", "depends_on: ["+a+"]\n")
	c, cDir := approve(t, owner, "c", "Monthly budgets in report", "")
	nothing()
	k, l := "261022-z0000", "261022-00000" // l the lower ID, k created first the same day
	planCycle(t, owner, k, "2026-10-22T09:00:00Z", false, a, b)
	nothing()
	planCycle(t, owner, k, "2026-10-22T09:00:00Z", true, a, b)
	os.WriteFile(filepath.Join(aDir, "proposal.md"), []byte("# Editing an unknown ID prints an error\n\nRevised.\n"), 0o644)
	commitAll(t, owner, "Revise "+a)
	nothing()
	planCycle(t, owner, k, "2026-10-22T09:00:00Z", true, a, b)
	planCycle(t, owner, l, "2026-10-22T10:00:00Z", true, c)
	git(t, owner, "push", "-q", "origin", "main")

	// Refusals: a local target of its own, an unclean checkout, even with
	// --retry, and a push origin refuses.
	os.WriteFile(filepath.Join(h1, "notes.txt"), []byte("mine\n"), 0o644)
	commitAll(t, h1, "A note of the operator's")
	claims(h1, "", `^grove: the local main has commits origin's does not, so it cannot be fast-forwarded\n$`)
	git(t, h1, "reset", "-q", "--hard", "HEAD~1")
	os.WriteFile(filepath.Join(h1, "README.md"), []byte("edited\n"), 0o644)
	claims(h1, "", `^grove: the checkout has uncommitted changes; commit, stash or discard them first\n$`, "--retry")
	if read(t, filepath.Join(h1, "README.md")) != "edited\n" {
		t.Error("the interrupted edit is gone")
	}
	git(t, h1, "checkout", "-q", "--", ".")
	hook := filepath.Join(remote, "hooks", "pre-receive")
	os.WriteFile(hook, []byte("#!/bin/sh\necho not allowed >&2\nexit 1\n"), 0o755)
	claims(h1, "", `^grove: git push: .*not allowed`)
	if got := git(t, h1, "for-each-ref", "refs/remotes/origin/grove"); got != "" {
		t.Errorf("a refused push left %q", got)
	}
	os.Remove(hook)

	// Two hosts at once: h1 claims while h2's push is under way, so h2's is
	// refused, and h2, reading again, finds the project paused.
	bin, _ := os.Executable()
	h1Out := filepath.Join(t.TempDir(), "h1.out")
	os.WriteFile(filepath.Join(h2, ".git/hooks/pre-push"), []byte("#!/bin/sh\nunset GIT_DIR\nrm \"$0\"\ncd "+h1+" && GROVE_TEST_MAIN=1 "+bin+" claim > "+h1Out+"\n"), 0o755)
	t.Chdir(h2)
	var stdout, stderr strings.Builder
	h2Branches := git(t, h2, "for-each-ref", "refs/heads")
	if code := run([]string{"claim"}, &stdout, &stderr); code != 1 || stdout.String() != "" || git(t, h2, "for-each-ref", "refs/heads") != h2Branches ||
		!regexp.MustCompile(`^grove: the project is paused after attempt [0-9]{6}-[0-9a-z]{5}, which `+regexp.QuoteMeta(host)+` started at \S+ on branch grove/[0-9]{6}-[0-9a-z]{5}; once you have dealt with it, run grove claim --retry\n$`).MatchString(stderr.String()) {
		t.Errorf("h2's claim at the same moment: status %d, %q, %q", code, stdout.String(), stderr.String())
	}
	main := strings.TrimSpace(git(t, owner, "rev-parse", "main"))
	rec := latest(remote, "grove/runs")
	if read(t, h1Out) != k+"\n" || recorded() != 1 || git(t, remote, "ls-tree", "grove/runs") != "" || rec.Cycle != k || strings.Join(rec.Changes, " ") != a+" "+b ||
		rec.Commit != main || rec.Host != host || rec.Started.IsZero() || strings.TrimSpace(git(t, h1, "branch", "--show-current")) != rec.Branch ||
		strings.TrimSpace(git(t, h1, "rev-parse", "HEAD", "main")) != main+"\n"+main || git(t, h1, "status", "--porcelain") != "" {
		t.Errorf("after h1's claim: %s, %+v, %s", read(t, h1Out), rec, git(t, h1, "status"))
	}
	if _, err := os.Stat(filepath.Join(h1, ".git/grove")); err == nil {
		t.Error("claim wrote .git/grove/")
	}

	// Paused, though another cycle is marked, until the cycle is archived:
	// a finished, b dropped.
	claims(h1, "", `^grove: the project is paused after attempt `+rec.ID+`, `)
	t.Chdir(owner)
	os.WriteFile(filepath.Join(aDir, "implementation.md"), []byte("## T-01\n\nStatus: done\n"), 0o644)
	os.WriteFile(filepath.Join(bDir, "implementation.md"), []byte("Status: dropped\n"), 0o644)
	maintain(t, owner, olderEdits)
	grove(t, "archive", a)
	grove(t, "archive", b)
	commitAll(t, owner, "Archive the batch")
	git(t, owner, "push", "-q", "origin", "main")
	claims(h2, l+"\n", "^$")
	if last := latest(remote, "grove/runs"); recorded() != 2 || last.Cycle != l || strings.Join(last.Changes, " ") != c || latest(remote, "grove/runs~1").ID != rec.ID {
		t.Errorf("the claim after the merge recorded %+v", last)
	}

	// A cycle given: refused, naming why, when there is no such cycle, a
	// member needs a person, or the project is paused; with --retry, claimed
	// though not marked while l is, l's record and branch kept.
	t.Chdir(owner)
	d, _ := approve(t, owner, "d", "Show saved records", "")
	e, _ := approve(t, owner, "e", "Export records", "")
	m, n := "261023-m0000", "261023-n0000"
	planCycle(t, owner, m, "2026-10-23T09:00:00Z", false, d)
	planCycle(t, owner, n, "2026-10-23T10:00:00Z", false, e)
	nFile := filepath.Join(owner, "grove/cycles", n+"-cycle.yaml")
	os.WriteFile(nFile, []byte(read(t, nFile)+"needs_human:\n  "+e+": judging the export\n"), 0o644)
	commitAll(t, owner, "Export needs a person")
	git(t, owner, "push", "-q", "origin", "main")
	claims(h2, "", `^grove: main at [0-9a-f]{7} holds no cycle 261099-zzzzz\n$`, "--retry", "261099-zzzzz")
	claims(h2, "", `^grove: `+e+` needs a person, as cycle `+n+`'s needs_human says; build the cycle with /grove-implement `+n+`\n$`, "--retry", n)
	claims(h2, "", `^grove: the project is paused after attempt `, m)
	lBranch := strings.TrimSpace(git(t, h2, "branch", "--show-current"))
	claims(h2, m+"\n", "^$", "--retry", m)
	if recorded() != 3 || latest(remote, "grove/runs").Cycle != m || latest(remote, "grove/runs~1").Cycle != l || git(t, h2, "branch", "--list", lBranch) == "" {
		t.Errorf("the claim given %s recorded %+v", m, latest(remote, "grove/runs"))
	}

	// No remote: the local target, commits of its own and all, is the fetched
	// one, the record goes on the local grove/runs, and merging the attempt's
	// branch into the local target ends the pause.
	git(t, h1, "fetch", "-q", "origin", "main:main")
	git(t, h1, "remote", "remove", "origin")
	git(t, h1, "switch", "-q", "main")
	os.WriteFile(filepath.Join(h1, "notes.txt"), []byte("mine\n"), 0o644)
	commitAll(t, h1, "A note of the operator's")
	local := strings.TrimSpace(git(t, h1, "rev-parse", "main"))
	claims(h1, l+"\n", "^$")
	if last := latest(h1, "grove/runs"); last.Commit != local || last.Cycle != l || recorded() != 3 {
		t.Errorf("the claim with no remote recorded %+v, want %s %s", last, l, local)
	}
	claims(h1, "", `^grove: the project is paused after attempt `)
	cRel, _ := filepath.Rel(owner, cDir)
	os.WriteFile(filepath.Join(h1, cRel, "implementation.md"), []byte("Status: dropped\n"), 0o644)
	grove(t, "archive", c)
	commitAll(t, h1, "Drop "+c)
	branch := strings.TrimSpace(git(t, h1, "branch", "--show-current"))
	git(t, h1, "switch", "-q", "main")
	git(t, h1, "merge", "-q", "--no-ff", "-m", "Merge the attempt", branch)
	claims(h1, "", `^grove: nothing to claim on main at [0-9a-f]{7}\n$`)
}

// TestReplan follows the walkthrough's planning step 5 and running step 6:
// a member revised with its re-plan in one squash or merge commit leaves
// the plan fresh; as two commits on main, or before a later mark, stale.
func TestReplan(t *testing.T) {
	dir := setUp(t)
	b, bDir := approve(t, dir, "b", "Add ledger delete ID", "")
	k := "261023-k0000"
	planCycle(t, dir, k, "2026-10-23T09:00:00Z", true, b)
	stale := "grove/cycles/" + k + "-cycle.yaml: planned: " + b + " changed on main since this plan, so grove claim skips this cycle; re-plan it with /grove-plan " + k + "\n"
	fresh := func(want bool, how string) {
		t.Helper()
		var stdout, stderr strings.Builder
		run([]string{"check"}, &stdout, &stderr)
		p, _ := openProject()
		ip, _ := p.agreed()
		cy, _ := eligible(p, ip)
		wantOut := map[bool]string{true: "No problems in grove/.\n", false: stale}[want]
		if got := stdout.String(); got != wantOut || (cy != nil) != want {
			t.Errorf("%s: check printed %q, grove claim would claim the cycle: %v", how, got, cy != nil)
		}
	}
	revise := func(branch, words string, commits int) {
		t.Helper()
		main := strings.TrimSpace(git(t, dir, "rev-parse", "main")) // the target commit propose names
		git(t, dir, "switch", "-q", "-c", branch)
		os.WriteFile(filepath.Join(bDir, "proposal.md"), []byte("# Add ledger delete ID\n\n"+words+"\n"), 0o644)
		if commits == 2 {
			commitAll(t, dir, "Revise "+b)
		}
		os.WriteFile(filepath.Join(dir, "grove/cycles", k+"-cycle.yaml"), []byte("title: Cycle "+k+"\ncreated: 2026-10-23T09:00:00Z\nchanges: ["+b+"]\nplanned: "+main+"\nunattended: true\n"), 0o644)
		commitAll(t, dir, "Re-plan "+k)
		git(t, dir, "switch", "-q", "main")
	}
	fresh(true, "planned")
	revise("squash", "Squashed.", 1)
	git(t, dir, "merge", "-q", "--squash", "squash")
	commitAll(t, dir, "Revise and re-plan, squashed")
	fresh(true, "a squash merge")
	revise("merge", "Merged.", 2)
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Revise and re-plan", "merge")
	fresh(true, "a merge commit")
	revise("ff", "Fast-forwarded.", 2)
	git(t, dir, "merge", "-q", "--ff-only", "ff")
	fresh(false, "two commits fast-forwarded")

	// A member changed in one merge, the cycle marked in a later one.
	planCycle(t, dir, k, "2026-10-23T09:00:00Z", false, b)
	git(t, dir, "switch", "-q", "-c", "later")
	os.WriteFile(filepath.Join(bDir, "proposal.md"), []byte("# Add ledger delete ID\n\nLater.\n"), 0o644)
	commitAll(t, dir, "Revise "+b)
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Revise", "later")
	cycle := filepath.Join(dir, "grove/cycles", k+"-cycle.yaml")
	os.WriteFile(cycle, []byte(read(t, cycle)+"unattended: true\n"), 0o644)
	commitAll(t, dir, "Mark")
	fresh(false, "marked after a member changed")
}

// TestMovedFiles follows the walkthrough's steps on a task's named file
// changing on main: a change alone reported against its base, a member of
// a marked cycle against its plan, and a Files: line Grove cannot read.
func TestMovedFiles(t *testing.T) {
	dir := setUp(t)
	base := strings.TrimSpace(git(t, dir, "rev-parse", "main"))
	cid, folder := propose(t, dir, "older-edits", "Allow editing older records", "base: "+base+"\n")
	// One path per bullet, ledger/ naming every file under it.
	tasks := "## T-01 — Allow editing older records\n\nResult: older records can be edited.\n\nFiles:\n- `ledger/`\n- `tests/test_cli.py`\nDelegated: test layout; `tests/test_records.py`\n\nAcceptance:\n- An older record saves.\n"
	os.WriteFile(filepath.Join(folder, "tasks.md"), []byte(tasks), 0o644)
	commitAll(t, dir, "Propose")
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Approve", "older-edits")
	checks := func(want string) {
		t.Helper()
		var stdout, stderr strings.Builder
		if code := run([]string{"check"}, &stdout, &stderr); stdout.String() != want || (code == 0) != (want == "No problems in grove/.\n") {
			t.Errorf("check: status %d\n%s%s\nwant\n%s", code, stdout.String(), stderr.String(), want)
		}
	}
	edit := func(file, msg string) {
		t.Helper()
		os.WriteFile(filepath.Join(dir, file), []byte(read(t, filepath.Join(dir, file))+"# "+msg+"\n"), 0o644)
		commitAll(t, dir, msg)
	}

	// Unchanged on main, or changed only on the branch building it; a path
	// only the Delegated: line names is not watched.
	checks("No problems in grove/.\n")
	os.WriteFile(filepath.Join(dir, "tests/test_records.py"), []byte("# Records\n"), 0o644)
	commitAll(t, dir, "Add test_records")
	checks("No problems in grove/.\n")
	git(t, dir, "switch", "-q", "-c", "older-edits-impl")
	edit("ledger/cli.py", "Build T-01")
	checks("No problems in grove/.\n")
	git(t, dir, "switch", "-q", "main")

	// Changed on main since base, one file under the folder T-01 names.
	edit("tests/test_cli.py", "Another change")
	edit("ledger/report.py", "A change under ledger/")
	rel := "grove/changes/" + filepath.Base(folder) + "/tasks.md"
	moved := "; revise the change with propose, which also re-bases it after a dependency lands\n"
	checks(rel + ": T-01's ledger/report.py changed on main since base " + base[:7] + moved + rel + ": T-01's tests/test_cli.py changed on main since base " + base[:7] + moved)
	for _, args := range [][]string{{cid}, {cid, "--unattended"}} {
		if out := grove(t, append([]string{"instructions", "implement"}, args...)...); !strings.Contains(out, "\n- Changed on main since base "+base[:7]+": T-01's `ledger/report.py`, T-01's `tests/test_cli.py`.\n") {
			t.Errorf("implement %v does not name the moved file:\n%.3000s", args, out)
		}
	}

	// A member of a cycle: reported against the plan, not base, whether or
	// not the cycle is marked.
	k, plan := "261023-k0000", strings.TrimSpace(git(t, dir, "rev-parse", "HEAD"))
	planCycle(t, dir, k, "2026-10-23T09:00:00Z", false, cid)
	checks("No problems in grove/.\n")
	edit("ledger/cli.py", "Yet another change")
	stale, cycle := "grove/cycles/"+k+"-cycle.yaml: planned: "+cid+" changed on main since this plan", filepath.Join(dir, "grove/cycles", k+"-cycle.yaml")
	checks(stale + "; re-plan it with /grove-plan " + k + "\n")
	for _, args := range [][]string{{cid}, {cid, "--unattended"}} {
		if out := grove(t, append([]string{"instructions", "implement"}, args...)...); !strings.Contains(out, "\n- Changed on main since cycle "+k+"'s plan "+plan[:7]+": T-01's `ledger/cli.py`.\n") || strings.Contains(out, "since base") {
			t.Errorf("implement %v does not name the file moved since the plan, alone:\n%.3000s", args, out)
		}
	}
	planned := read(t, cycle)
	os.WriteFile(cycle, []byte(regexp.MustCompile(`planned: \w+`).ReplaceAllString(planned, "planned: 0000000")), 0o644)
	checks("grove/cycles/" + k + "-cycle.yaml: planned: 0000000 is not a commit in this repository\n")
	os.WriteFile(cycle, []byte(planned+"unattended: true\n"), 0o644)
	commitAll(t, dir, "Mark")
	checks(stale + ", so grove claim skips this cycle; re-plan it with /grove-plan " + k + "\n")
	p, _ := openProject()
	if ip, err := p.agreed(); err != nil || func() bool { cy, _ := eligible(p, ip); return cy != nil }() {
		t.Errorf("grove claim would claim a cycle whose member's file moved: %v", err)
	}

	// A comma line, wrapped, still reads; a Files: line Grove cannot read.
	files := "Files:\n- `ledger/`\n- `tests/test_cli.py`"
	os.WriteFile(filepath.Join(folder, "tasks.md"), []byte(strings.Replace(tasks, files, "Files: `ledger/cli.py`,\n`tests/test_cli.py`", 1)), 0o644)
	os.Remove(filepath.Join(dir, "grove/cycles", k+"-cycle.yaml"))
	since := " changed on main since base " + base[:7] + moved
	checks(rel + ": T-01's ledger/cli.py" + since + rel + ": T-01's tests/test_cli.py" + since)
	for _, bad := range []string{"Files: ledger/cli.py", "Files: `../ledger/cli.py`", "Files:\n- `ledger/cli.py`, `tests/test_cli.py`"} {
		os.WriteFile(filepath.Join(folder, "tasks.md"), []byte(strings.Replace(tasks, files, bad, 1)), 0o644)
		checks(rel + ":5: T-01's Files: line should list repo-relative paths in backticks, one per bullet or separated by commas, and a folder ends in /\n")
	}
}

// TestFollows follows the walkthrough's steps on a cycle that follows
// another: run waits on it, and the merges archiving its changes, or those
// of a cycle it follows in turn, leave the plan fresh, and the implement
// instructions naming no moved file; nothing else does, but a re-plan.
func TestFollows(t *testing.T) {
	dir := setUp(t)
	b, bDir := approve(t, dir, "b", "Add ledger delete ID", "")
	os.WriteFile(filepath.Join(bDir, "tasks.md"), []byte("## T-01 — Add ledger delete ID\n\nFiles: `ledger/cli.py`\n"), 0o644)
	commitAll(t, dir, "Name B's file")
	c, cDir := approve(t, dir, "c", "Monthly budgets in report", "")
	d, dDir := approve(t, dir, "d", "Show saved records", "")
	os.WriteFile(filepath.Join(dDir, "tasks.md"), []byte("## T-01 — Show saved records\n\nFiles: `ledger/cli.py`\n"), 0o644)
	commitAll(t, dir, "Name D's file")
	e, eDir := approve(t, dir, "e", "Export records", "")
	f, fDir := approve(t, dir, "f", "Import a CSV", "")
	k, l, m, n := "261023-k0000", "261023-p0000", "261023-q0000", "261023-r0000" // K follows L, L follows M
	planCycle(t, dir, m, "2026-10-23T11:00:00Z", false, e)
	planCycle(t, dir, l, "2026-10-23T10:00:00Z", true, c, d)
	planCycle(t, dir, k, "2026-10-23T09:00:00Z", true, b)
	planCycle(t, dir, n, "2026-10-23T12:00:00Z", false, f)
	for id, follows := range map[string]string{k: l, l: m} {
		file := filepath.Join(dir, "grove/cycles", id+"-cycle.yaml")
		os.WriteFile(file, []byte(read(t, file)+"depends_on: ["+follows+"]\n"), 0o644)
	}
	commitAll(t, dir, "Order the cycles")

	staleOf := func(cycle, change string) string {
		return "grove/cycles/" + cycle + "-cycle.yaml: planned: " + change + " changed on main since this plan, so grove claim skips this cycle; re-plan it with /grove-plan " + cycle + "\n"
	}
	clean, stale := "No problems in grove/.\n", staleOf(k, b)
	state := func(want, selects, how string) {
		t.Helper()
		var stdout, stderr strings.Builder
		run([]string{"check"}, &stdout, &stderr)
		p, _ := openProject()
		ip, _ := p.agreed()
		got := ""
		if cy, _ := eligible(p, ip); cy != nil {
			got = cy.ID
		}
		if stdout.String() != want || got != selects {
			t.Errorf("%s: check printed %q, grove claim selects %q, want %q", how, stdout.String(), got, selects)
		}
	}
	namesIn := func(change, cycle string, want bool, how string) {
		t.Helper()
		if out := grove(t, "instructions", "implement", change); strings.Contains(out, "\n- Changed on main since cycle "+cycle+"'s plan ") != want {
			t.Errorf("%s: implement %s names a moved file: %v, want %v", how, change, !want, want)
		}
	}
	names := func(want bool, how string) { t.Helper(); namesIn(b, k, want, how) }
	edit := func(msg string) {
		os.WriteFile(filepath.Join(dir, "ledger/cli.py"), []byte(read(t, filepath.Join(dir, "ledger/cli.py"))+"# "+msg+"\n"), 0o644)
		commitAll(t, dir, msg)
	}
	// land builds change id on a branch, archiving it and editing B's file, and merges it.
	land := func(id, folder string, squash bool) {
		t.Helper()
		git(t, dir, "switch", "-q", "-c", "land-"+id)
		os.WriteFile(filepath.Join(folder, "implementation.md"), []byte("Status: dropped\n"), 0o644)
		grove(t, "archive", id)
		edit("Build " + id)
		git(t, dir, "switch", "-q", "main")
		if squash {
			git(t, dir, "merge", "-q", "--squash", "land-"+id)
			commitAll(t, dir, "Land "+id+", squashed")
		} else {
			git(t, dir, "merge", "-q", "--no-ff", "-m", "Land "+id, "land-"+id)
		}
	}

	// K waits on L, which waits on M; given K, run refuses it, and implement
	// builds it with the owner as it would without depends_on.
	state(clean, "", "planned")
	var stdout, stderr strings.Builder
	branches := git(t, dir, "branch")
	if code := run([]string{"claim", k}, &stdout, &stderr); code != 1 || stderr.String() != "grove: cycle "+k+" waits on cycle "+l+", which is not done\n" || git(t, dir, "branch") != branches {
		t.Errorf("grove claim %s: status %d, %q", k, code, stderr.String())
	}
	if out := grove(t, "instructions", "implement", k); !strings.Contains(out, "\n1. "+b+" \"Add ledger delete ID\"\n") {
		t.Errorf("implement %s does not list its member:\n%.3000s", k, out)
	}

	// M lands in a merge commit, then L as two squashes, the last also as a
	// merge commit. L's own archive leaves it fresh, though it edits a file
	// d names; a commit archiving nothing does not. The older marked K waits
	// until L is done.
	land(e, eDir, false)
	state(clean, l, "M's change archived in a merge commit")
	land(c, cDir, true)
	state(clean, l, "one of L's changes archived, editing a file its other change names")
	namesIn(d, l, false, "one of L's changes archived")
	edit("Another change")
	state(stale+staleOf(l, d), "", "a commit archiving nothing after one of L's changes")
	namesIn(d, l, true, "a commit archiving nothing after one of L's changes")
	git(t, dir, "reset", "-q", "--hard", "HEAD~1")
	land(d, dDir, true)
	state(clean, k, "L's last change archived in a squash")
	git(t, dir, "reset", "-q", "--hard", "HEAD~1")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Land "+d, "land-"+d)
	state(clean, k, "L's last change archived in a merge commit")
	names(false, "L's last change archived")

	// Any other commit still counts, until K is re-planned.
	edit("Another change")
	state(stale, "", "a commit archiving nothing")
	names(true, "a commit archiving nothing")
	git(t, dir, "reset", "-q", "--hard", "HEAD~1")
	land(f, fDir, false)
	state(stale, "", "an archive of a cycle K does not follow")
	file := filepath.Join(dir, "grove/cycles", k+"-cycle.yaml")
	os.WriteFile(file, []byte(regexp.MustCompile(`planned: \w+`).ReplaceAllString(read(t, file), "planned: "+strings.TrimSpace(git(t, dir, "rev-parse", "HEAD")))), 0o644)
	commitAll(t, dir, "Re-plan "+k)
	state(clean, k, "re-planned")
	names(false, "re-planned")
}

// TestMisnamedCycle follows the walkthrough's planning step 6, a cycle
// file named without an ID: check reports it, and changes names no cycle
// for its member, while a cycle named well is shown as before.
func TestMisnamedCycle(t *testing.T) {
	dir := setUp(t)
	a, _ := approve(t, dir, "a", "Editing an unknown ID prints an error", "")
	b, _ := approve(t, dir, "b", "Add ledger delete ID", "")
	planCycle(t, dir, "261105-proto-unknown-id", "2026-11-05T09:00:00Z", false, a) // IDs leave out o
	planCycle(t, dir, "261105-pr0t0", "2026-11-05T10:00:00Z", false, b)
	var stdout, stderr strings.Builder
	if code := run([]string{"check"}, &stdout, &stderr); code == 0 || !strings.Contains(stdout.String(), "grove/cycles/261105-proto-unknown-id-cycle.yaml: the file name must start with a cycle ID, such as 261015-k3m9q\n") {
		t.Errorf("check of a misnamed cycle: status %d, %q", code, stdout.String())
	}
	out := grove(t, "changes")
	if !strings.Contains(out, a+"  Editing an unknown ID prints an error  approved\n") || !strings.Contains(out, b+"  Add ledger delete ID  approved  in cycle 261105-pr0t0\n") {
		t.Errorf("changes with a misnamed cycle:\n%s", out)
	}
}

// TestShallowClone follows the walkthrough's "A shallow clone": a base or
// planned commit the clone did not fetch is named as missing from it.
func TestShallowClone(t *testing.T) {
	dir := setUp(t)
	base := strings.TrimSpace(git(t, dir, "rev-parse", "HEAD"))
	a, aDir := approve(t, dir, "a", "Editing an unknown ID prints an error", "base: "+base+"\n")
	planCycle(t, dir, "261105-pr0t0", "2026-11-05T09:00:00Z", false, a) // planned is the commit before this one
	planned := strings.TrimSpace(git(t, dir, "rev-parse", "HEAD~1"))
	checks := func(want string) {
		t.Helper()
		var stdout, stderr strings.Builder
		if run([]string{"check"}, &stdout, &stderr); !strings.Contains(stdout.String(), want) {
			t.Errorf("check lacks %q:\n%s", want, stdout.String())
		}
	}
	clone := t.TempDir()
	git(t, dir, "clone", "-q", "--depth", "1", "file://"+dir, clone)
	t.Chdir(clone)
	checks("/.grove.yaml: base: " + base + " is not in this shallow clone; fetch its full history, such as with git fetch --unshallow\n")
	checks("grove/cycles/261105-pr0t0-cycle.yaml: planned: " + planned + " is not in this shallow clone; fetch its full history, such as with git fetch --unshallow\n")
	t.Chdir(dir)
	file := filepath.Join(aDir, ".grove.yaml")
	os.WriteFile(file, []byte(strings.Replace(read(t, file), "base: "+base, "base: 0000000", 1)), 0o644)
	checks("/.grove.yaml: base: 0000000 is not a commit in this repository\n")
}

// serving starts grove serve with args in the working directory and
// returns what it printed first, and its error once it stops.
func serving(t *testing.T, args ...string) (string, <-chan string) {
	t.Helper()
	r, w := io.Pipe()
	stopped := make(chan string, 1)
	go func() {
		var stderr strings.Builder
		run(append([]string{"serve"}, args...), w, &stderr)
		w.Close()
		stopped <- stderr.String()
	}()
	line, _ := bufio.NewReader(r).ReadString('\n')
	go io.Copy(io.Discard, r)
	return line, stopped
}

func get(t *testing.T, url string) (int, string) {
	t.Helper()
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(b)
}

// TestServe follows the walkthrough's "Watching the project": the fixture
// with cycles ready, needing a person and stale, branches of every kind,
// an uncommitted edit and attempts recorded on origin, served and read.
func TestServe(t *testing.T) {
	dir := setUp(t)
	remote := filepath.Join(t.TempDir(), "remote.git")
	git(t, dir, "init", "-q", "--bare", remote)
	git(t, dir, "remote", "add", "origin", remote)
	write := func(rel, s string) { os.WriteFile(filepath.Join(dir, rel), []byte(s), 0o644) }
	maintain(t, dir, olderEdits)
	commitAll(t, dir, "Record editing")
	a, _ := approve(t, dir, "a", "Editing an unknown ID prints an error", "")
	b, bDir := approve(t, dir, "b", "Add ledger delete ID", "depends_on: ["+a+"]\n")
	os.WriteFile(filepath.Join(bDir, "tasks.md"), []byte("## T-01 — Delete\n\n## T-02 — Confirm\n"), 0o644)
	commitAll(t, dir, "Tasks for "+b)
	c, cDir := approve(t, dir, "c", "Export the ledger as CSV", "")
	d, dDir := approve(t, dir, "d", "Monthly budgets in report", "")
	f, fDir := approve(t, dir, "f", "Editing a record prints what was saved", "")
	e, _ := approve(t, dir, "e", "Import records from a CSV", "depends_on: ["+f+"]\n")
	os.WriteFile(filepath.Join(fDir, "proposal.md"), []byte("# Editing a record prints what was saved\n\nRevised.\n"), 0o644)
	commitAll(t, dir, "Revise "+f)
	git(t, dir, "switch", "-q", "-c", "apply/f")
	os.WriteFile(filepath.Join(fDir, "implementation.md"), []byte("### T-01\n\nStatus: done\n"), 0o644)
	maintain(t, dir, olderEdits+"\nBuilt.\n")
	commitAll(t, dir, "Build "+f)
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--squash", "apply/f")
	commitAll(t, dir, "Build "+f+" (squashed)")
	n, k, l, m := "261101-n0000", "261102-k0000", "261103-p0000", "261104-m0000"
	planCycle(t, dir, n, "2026-11-01T00:00:00Z", false, f)
	os.MkdirAll(filepath.Join(dir, "grove/changes/archive"), 0o755)
	git(t, dir, "mv", fDir, filepath.Join(dir, "grove/changes/archive", filepath.Base(fDir)))
	commitAll(t, dir, "Archive "+f)
	maintain(t, dir, olderEdits+"\nEdited again.\n")
	commitAll(t, dir, "Edit record-editing")
	planCycle(t, dir, k, "2026-11-02T00:00:00Z", true, a, b)
	write("grove/cycles/"+l+"-cycle.yaml", "title: Cycle L\ncreated: 2026-11-03T00:00:00Z\nchanges: ["+c+"]\nplanned: "+strings.TrimSpace(git(t, dir, "rev-parse", "HEAD"))+"\nneeds_human:\n  "+c+": a person judges the file\n")
	commitAll(t, dir, "Plan "+l)
	planCycle(t, dir, m, "2026-11-04T00:00:00Z", true, d)
	os.WriteFile(filepath.Join(dDir, "proposal.md"), []byte("# Monthly budgets in report\n\nRevised.\n"), 0o644)
	commitAll(t, dir, "Revise "+d)
	write("ledger/records.py", "# records\n")
	commitAll(t, dir, "Records only")
	git(t, dir, "push", "-q", "origin", "main")
	git(t, dir, "switch", "-q", "-c", "build-b")
	os.WriteFile(filepath.Join(bDir, "implementation.md"), []byte("### T-01\n\nStatus: done\n\n### T-02\n\nStatus: blocked\n"), 0o644)
	commitAll(t, dir, "Build "+b)
	git(t, dir, "push", "-q", "origin", "build-b")
	git(t, dir, "switch", "-q", "-c", "teammate", "main")
	os.WriteFile(filepath.Join(bDir, "implementation.md"), []byte("### T-01\n\nStatus: to check\n"), 0o644)
	commitAll(t, dir, "Build "+b+" elsewhere")
	git(t, dir, "push", "-q", "origin", "teammate")
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "branch", "-q", "-D", "teammate")
	g, gDir := propose(t, dir, "draft-notes", "Notes on records", "")
	os.WriteFile(filepath.Join(gDir, "proposal.md"), []byte("# Notes on records\n\n<script>alert(1)</script>\n"), 0o644)
	commitAll(t, dir, "Draft "+g)
	git(t, dir, "switch", "-q", "main")
	os.WriteFile(filepath.Join(cDir, "implementation.md"), []byte("### T-01\n\nStatus: in progress\n"), 0o644)

	// The default address, then that address taken, then one given.
	if url, _ := serving(t); url != "http://127.0.0.1:7070/\n" {
		t.Fatalf("grove serve printed %q", url)
	}
	if code, page := get(t, "http://127.0.0.1:7070/"); code != 200 || !strings.Contains(page, "<title>Grove</title>") {
		t.Errorf("/ answered %d:\n%.300s", code, page)
	}
	if url, stopped := serving(t); url != "" || !strings.Contains(<-stopped, "address already in use") {
		t.Errorf("a second grove serve printed %q", url)
	}
	if url, stopped := serving(t, "--addr", "0.0.0.0:0"); url != "" || !strings.Contains(<-stopped, "--addr must be a loopback address") {
		t.Errorf("grove serve on every address printed %q", url)
	}
	url, _ := serving(t, "--addr", "127.0.0.1:0")
	url = strings.TrimSpace(url)
	if !regexp.MustCompile(`^http://127\.0\.0\.1:[0-9]+/$`).MatchString(url) {
		t.Fatalf("grove serve --addr printed %q", url)
	}

	req, _ := http.NewRequest("GET", url+"api/project", nil)
	req.Host = "rebound.example:7070" // a name a site pointed at loopback
	if res, err := http.DefaultClient.Do(req); err != nil || res.StatusCode != 403 {
		t.Errorf("a request for another host answered %v, %v", res, err)
	}
	_, body := get(t, url+"api/project")
	var v value
	if err := json.Unmarshal([]byte(body), &v); err != nil {
		t.Fatalf("%v\n%.500s", err, body)
	}
	head := strings.TrimSpace(git(t, dir, "rev-parse", "main"))
	if v.Root != "grove" || v.Target != "main" || v.Commit != head || len(v.Specs) != 1 {
		t.Errorf("value names %s %s %s with %d specs", v.Root, v.Target, v.Commit, len(v.Specs))
	}
	for _, ch := range v.Changes {
		if ch.ID == e && (ch.Cycle != "" || !ch.Ready) || ch.ID == a && ch.Cycle != k {
			t.Errorf("change %s in cycle %q, ready %v", ch.ID, ch.Cycle, ch.Ready)
		}
	}
	verdicts := map[string]pageCycle{}
	for _, cy := range v.Cycles {
		verdicts[cy.ID] = cy
	}
	if cy := verdicts[k]; strings.Join(cy.Changes, " ") != a+" "+b || !cy.Runs || cy.Verdict != "grove claim would claim it" || len(v.Runs) != 0 {
		t.Errorf("cycle K: %+v", cy)
	}
	if cy := verdicts[l]; cy.Runs || !strings.HasPrefix(cy.Verdict, c+" needs a person") {
		t.Errorf("cycle L: %+v", cy)
	}
	if cy := verdicts[m]; strings.Join(cy.Stale, " ") != d || !strings.Contains(cy.Verdict, "stale, since "+d+" changed") {
		t.Errorf("cycle M: %+v", cy)
	}
	if cy := verdicts[n]; !cy.Done || strings.Join(cy.Archived, " ") != f {
		t.Errorf("cycle N: %+v", cy)
	}

	// Attempts another host recorded on origin, none in this checkout: K
	// claimed, then the project paused after L's.
	record := func(id, cycle, changes, parent string) string {
		tree := strings.TrimSpace(git(t, dir, "hash-object", "-t", "tree", "-w", "/dev/null"))
		args := []string{"commit-tree", tree, "-m", `{"id":"` + id + `","started":"2026-11-05T09:00:00Z","cycle":"` + cycle + `","changes":[` + changes + `],"host":"h2.example","branch":"grove/` + id + `"}`}
		if parent != "" {
			args = append(args, "-p", parent)
		}
		rec, work := strings.TrimSpace(git(t, dir, args...)), strings.TrimSpace(git(t, dir, "commit-tree", "main^{tree}", "-p", "main", "-m", "Work on "+cycle))
		git(t, dir, "push", "-q", "origin", rec+":refs/heads/grove/runs", work+":refs/heads/grove/"+id)
		return rec
	}
	claimed := func() (pageCycle, []pageRun) {
		var w value
		_, body := get(t, url+"api/project")
		json.Unmarshal([]byte(body), &w)
		for _, cy := range w.Cycles {
			if cy.ID == k {
				return cy, w.Runs
			}
		}
		return pageCycle{}, w.Runs
	}
	first := record("261105-r0000", k, `"`+a+`","`+b+`"`, "")
	cy, recs := claimed()
	if cy.Runs || cy.Verdict != "claimed by attempt 261105-r0000, which h2.example started at 2026-11-05T09:00:00Z on branch grove/261105-r0000" ||
		len(recs) != 1 || recs[0].Host != "h2.example" || recs[0].Cycle != k || recs[0].Branch != "grove/261105-r0000" || recs[0].Started.IsZero() || recs[0].Moved <= "2026-11-05" ||
		git(t, dir, "for-each-ref", "refs/heads/grove") != "" {
		t.Errorf("cycle K claimed elsewhere: %+v, runs %+v", cy, recs)
	}
	git(t, dir, "branch", "grove/runs", first) // a local copy, not a branch of work
	record("261106-r0000", l, `"`+c+`"`, first)
	if cy, recs := claimed(); cy.Runs || cy.Verdict != "grove claim --retry would claim it; the project is paused after attempt 261106-r0000 of cycle "+l || len(recs) != 2 || recs[0].ID != "261106-r0000" {
		t.Errorf("cycle K while paused after L: %+v, runs %+v", cy, recs)
	}

	var rows []string
	for _, r := range v.Branches {
		row := fmt.Sprintf("%s origin_only=%v pushed=%v", r.Name, r.OriginOnly, r.Pushed)
		for _, ch := range r.Changes {
			row += fmt.Sprintf(" %s draft=%v", ch.ID, ch.Draft)
			for _, t := range ch.Tasks {
				row += " " + t.ID + ":" + t.Status
			}
			row += " new=" + strings.Join(ch.New, ",")
		}
		rows = append(rows, row)
	}
	cRel, bRel, gRel := "changes/"+filepath.Base(cDir)+"/", "changes/"+filepath.Base(bDir)+"/", "changes/"+filepath.Base(gDir)+"/"
	want := []string{
		"working tree origin_only=false pushed=false " + c + " draft=false T-01:in progress new=" + cRel + "implementation.md",
		"draft-notes origin_only=false pushed=false " + g + " draft=true T-01: new=" + strings.Join([]string{gRel + ".grove.yaml", gRel + "proposal.md", gRel + "tasks.md"}, ","),
		"origin/teammate origin_only=true pushed=false " + b + " draft=false T-01:to check T-02: new=" + bRel + "implementation.md",
		"build-b origin_only=false pushed=true " + b + " draft=false T-01:done T-02:blocked new=" + bRel + "implementation.md",
	}
	// The branches' commits share a second, so only the working tree, newest, has its place.
	slices.Sort(rows[1:])
	slices.Sort(want[1:])
	if strings.Join(rows, "\n") != strings.Join(want, "\n") {
		t.Errorf("branches:\n%s\nwant:\n%s", strings.Join(rows, "\n"), strings.Join(want, "\n"))
	}

	marks := map[string]string{}
	var order []string
	for _, ev := range v.History {
		for _, mk := range ev.Marks {
			marks[ev.Subject] += mk.Kind + " " + mk.ID + "; "
		}
		order = append(order, ev.Subject)
	}
	for subject, want := range map[string]string{
		"Approve " + f: "approved " + f + "; ", "Revise " + f: "revised " + f + "; ",
		"Build " + f + " (squashed)": "revised " + f + "; spec record-editing; ", "Plan " + n: "planned " + n + "; ",
		"Archive " + f: "archived " + f + "; ", "Plan " + m: "planned " + m + "; ", "Records only": "",
	} {
		if marks[subject] != want {
			t.Errorf("history marks %q as %q, want %q", subject, marks[subject], want)
		}
	}
	if i, j := slices.Index(order, "Archive "+f), slices.Index(order, "Approve "+f); i < 0 || i > j {
		t.Errorf("history is not newest first: %q", order)
	}
	if len(v.Problems) != 1 || !strings.Contains(v.Problems[0], d+" changed on main since this plan") {
		t.Errorf("problems: %q", v.Problems)
	}

	// Documents: a proposal, a cycle file, a branch's file, raw HTML, a path outside the root.
	tip := strings.TrimSpace(git(t, dir, "rev-parse", "build-b"))
	draft := strings.TrimSpace(git(t, dir, "rev-parse", "draft-notes"))
	for _, tc := range []struct{ query, want string }{
		{"path=changes/" + filepath.Base(dDir) + "/proposal.md", "<h1>Monthly budgets in report</h1>"},
		{"path=cycles/" + k + "-cycle.yaml", "unattended: true\n"},
		{"ref=" + tip + "&path=" + bRel + "implementation.md", "<p>Status: blocked</p>"},
		{"ref=worktree&path=" + cRel + "implementation.md", "<p>Status: in progress</p>"},
		{"ref=" + draft + "&path=" + gRel + "proposal.md", "<!-- raw HTML omitted -->"},
	} {
		if code, got := get(t, url+"api/doc?"+tc.query); code != 200 || !strings.Contains(got, tc.want) || strings.Contains(got, "<script") {
			t.Errorf("/api/doc?%s answered %d:\n%s", tc.query, code, got)
		}
	}
	for _, query := range []string{"path=../grove.yaml", "path=notes.txt", "ref=main&path=cycles/" + k + "-cycle.yaml"} {
		if code, _ := get(t, url+"api/doc?"+query); code != 400 {
			t.Errorf("/api/doc?%s answered %d, want 400", query, code)
		}
	}
}
