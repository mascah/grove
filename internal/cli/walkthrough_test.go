package cli

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
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
			if skill := read(t, filepath.Join(dir, h, "skills", "grove-"+w, "SKILL.md")); !strings.Contains(skill, "name: grove-"+w+"\n") || !strings.Contains(skill, " instructions "+w) {
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
	os.WriteFile(filepath.Join(dir, "grove", "specs", "record-editing", ".grove.yaml"), []byte("description: Editing records.\nowner: me\n"), 0o644)
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
	want = "grove/specs/record-editing/.grove.yaml:2: unknown field \"owner\"; this file takes title and description\n" +
		"grove/specs/record-editing/.grove.yaml: title: required field is missing\n" +
		"grove/specs/record-editing/spec.md:5: requirement \"Keep a category\" has no scenario\n" +
		"grove/specs/record-editing/spec.md:9: requirement \"Delete a record\" has no \"shall\"\n" +
		"grove/specs/record-editing/spec.md:3: this heading belongs in a change spec; a maintained spec lists its requirements under ## Requirements\n" +
		"grove/changes/261015-k3m9q-no-proposal/proposal.md: missing; it is the change's main document in schema default\n" +
		"grove/changes/261015-p4x8d-extra-field/.grove.yaml:3: unknown field \"owner\"; this file takes title, created, schema, description, depends_on and base\n" +
		"grove/changes/261015-q7r2s-drifted/proposal.md:1: the first heading should be \"# New title\", the title in .grove.yaml\n" +
		"grove/changes/no-id/: the folder name must start with a change ID, such as 261015-k3m9q\n"
	if code != 1 || stdout.String() != want || stderr.String() != "grove: 9 problems found\n" {
		t.Errorf("check: status %d\n%s%s", code, stdout.String(), stderr.String())
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
	if out := grove(t, "instructions", "propose"); !strings.Contains(out, "citing the prototype branch and the commit") {
		t.Error("propose instructions do not say to cite a prototype")
	}
}

func TestPropose(t *testing.T) {
	dir := setUp(t)
	cid, _ := propose(t, dir, "older-edits", "Allow editing older records", "")
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
	out := grove(t, "instructions", "propose", "Editing a missing record should print an error")
	req, how := strings.Index(out, "## What Grove requires of propose\n\n- A change describes future work."), strings.Index(out, "## How to propose\n\n**Propose says")
	if !strings.Contains(out, "## Owner's description\n\nEditing a missing record should print an error\n") || req < 0 || how < req ||
		!strings.Contains(out, "on a `Files:` line after `Result:`") || !strings.Contains(out, "\nFiles: [`path/to/file`") || strings.Contains(out, "file-by-file") ||
		!strings.Contains(out, "before you\n  write any document, ask one round") || !strings.Contains(out, "worded so that \"yes\"\n  accepts it") || !strings.Contains(out, "Before you tell the owner the change is written, check") || !strings.Contains(out, "### 5. Check the change") || strings.Contains(out, "Assumptions to confirm") {
		t.Errorf("propose from a description:\n%.3000s", out)
	}

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

	var stdout, stderr strings.Builder
	for _, args := range [][]string{{"implement", cid}, {"implement", cid, "--unattended"}} {
		stderr.Reset()
		if code := run(append([]string{"instructions"}, args...), &stdout, &stderr); code != 1 || stderr.String() != "grove: main holds no approved change "+cid+"\n" {
			t.Errorf("%v before the merge: status %d, %q", args, code, stderr.String())
		}
	}
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Approve", "older-edits")
	k := "261022-k0000"
	stderr.Reset()
	if code := run([]string{"instructions", "apply", k}, &stdout, &stderr); code != 2 || !strings.HasPrefix(stderr.String(), "usage: grove") {
		t.Errorf("instructions apply: status %d, %q", code, stderr.String())
	}

	shared := []string{"## What Grove requires of implement\n\n", "Before you finish, run the project's tests, re-read each done task's", "each unverified scenario. Commit",
		"The project's own rules for this skill:\n\n- Commit after each task, naming its task ID.\n"}
	owner := []string{"ask the owner and\n  wait", "with their agreement add it to that task's `Files:`", "tell the owner how to try it"}
	alone := []string{"ask no one and wait for no one's answer", "A file no task names stays as it is", "mark that task `Status: blocked`, naming the file", "whatever its\n  `needs_human` says"}
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
	has("implement instructions", out, append(append([]string{"You are implementing approved change " + cid + ", \"Allow editing older records\", as main@", "with the owner, in their own working tree.", "### change-spec: ", "### Requirement: Edit older records"}, shared...), owner...), alone)
	out = grove(t, "instructions", "implement", cid, "--unattended")
	has("unattended instructions", out, append(append([]string{"asking no one, in your own working tree.", "## Project context\n\nRun python3 -m unittest before saying a task is done.\n\n## What Grove requires of implement\n"}, shared...), alone...), owner)

	// A member of a cycle, not marked, that needs a person, gets its plan.
	planCycle(t, dir, k, "2026-10-22T09:00:00Z", false, cid)
	cycle := read(t, filepath.Join(dir, "grove/cycles", k+"-cycle.yaml")) + "needs_human:\n  " + cid + ": the owner judges it\n"
	os.WriteFile(filepath.Join(dir, "grove/cycles", k+"-cycle.yaml"), []byte(cycle), 0o644)
	commitAll(t, dir, "Needs a person")
	for _, args := range [][]string{{cid}, {cid, "--unattended"}} {
		out := grove(t, append([]string{"instructions", "implement"}, args...)...)
		has(strings.Join(args, " "), out, []string{", change 1 of 1 in cycle " + k + ".", "## The cycle\n\nCycle " + k + ", \"Cycle " + k + "\", in build order:\n\n1. " + cid + " \"Allow editing older records\" — this change\n"}, nil)
	}
	out = grove(t, "instructions", "implement", k)
	has("cycle instructions", out, []string{"You are implementing cycle " + k + ", \"Cycle " + k + "\", as main@", "with the owner", "run Grove's `instructions implement <change-id>` and", "\n1. " + cid + " \"Allow editing older records\"\n", "## The cycle\n\n<file path=\"grove/cycles/" + k + "-cycle.yaml\">\n", shared[3]}, alone)
	out = grove(t, "instructions", "implement", k, "--unattended")
	has("unattended cycle instructions", out, append([]string{"asking no one", "run Grove's `instructions implement <change-id> --unattended` and", "\n1. " + cid + " \"Allow editing older records\"\n", shared[3]}, alone...), owner)
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
	os.WriteFile(filepath.Join(dir, folder, "implementation.md"), []byte("## T-01\n\nStatus: in progress\n"), 0o644)
	var stdout, stderr strings.Builder
	if code := run([]string{"archive", cid}, &stdout, &stderr); code != 1 || stderr.String() != "grove: T-01 is not done in implementation.md; archive a change when every task is done or the file says Status: dropped\n" {
		t.Errorf("archiving an unfinished change: status %d, %q", code, stderr.String())
	}
	os.WriteFile(filepath.Join(dir, folder, "implementation.md"), []byte("## T-01\n\nStatus: done\n"), 0o644)

	// The session reconciles the change spec into a new maintained spec,
	// from the template its instructions give, then archives the change.
	out := grove(t, "instructions", "archive", cid)
	for _, s := range []string{"Reconcile the specs", "word for word", "run `archive` again", "marks `unverified`; it does\n  not stop the archive", "Status: done", "### Requirement: Edit older records"} {
		if !strings.Contains(out, s) {
			t.Errorf("archive instructions lack %q", s)
		}
	}
	os.MkdirAll(filepath.Join(dir, "grove", "specs", "record-editing"), 0o755)
	for _, name := range []string{".grove.yaml", "spec.md"} {
		m := regexp.MustCompile(`(?s)<file path="grove/specs/<capability>/` + regexp.QuoteMeta(name) + `">\n(.*?)\n</file>`).FindStringSubmatch(out)
		if m == nil {
			t.Fatalf("no %s template in:\n%s", name, out)
		}
		filled := regexp.MustCompile(`(?s)### Requirement: \[.*?\n\n## Known`).ReplaceAllString(m[1], olderEdits+"\n## Known")
		filled = strings.NewReplacer("<capability>", "record-editing", "<title>", "Record editing").Replace(filled)
		os.WriteFile(filepath.Join(dir, "grove", "specs", "record-editing", name), []byte(filled), 0o644)
	}
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
	out := grove(t, "instructions", "propose", "--schema", "bugfix", "Editing a missing record crashes")
	for _, s := range []string{", with change schema bugfix, from the owner's description", "## How to propose\n\n**Propose a bug fix.**", "/proposal.md\">\n# <title>\n\n## Observed\n",
		"### change-spec: ", "### tasks: ", "/.grove.yaml\">\ntitle: <title>\ncreated: ", "schema: bugfix\n"} {
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

	// The walkthrough's "A mistake in a schema", a change naming no schema,
	// and a schema without tasks.
	schemas := filepath.Join(dir, "grove/schemas")
	os.Rename(filepath.Join(schemas, "bugfix/documents/proposal"), filepath.Join(schemas, "bugfix/documents/propsal"))
	os.WriteFile(filepath.Join(schemas, "bugfix/schema.yaml"), []byte("description: A fix for behavior that is wrong.\ninstrution: x\n"), 0o644)
	os.WriteFile(filepath.Join(folder, ".grove.yaml"), []byte("title: Allow editing older records\ncreated: 2026-10-14T00:00:00Z\nschema: bugfx\n"), 0o644)
	os.MkdirAll(filepath.Join(schemas, "notasks"), 0o755)
	os.WriteFile(filepath.Join(schemas, "notasks/schema.yaml"), []byte("description: N.\ndocuments:\n  - {kind: proposal, path: changes/<change>/proposal.md, template: documents/proposal/template.md, guide: documents/proposal/guide.md, metadata: changes/<change>/.grove.yaml}\n"), 0o644)
	stdout.Reset()
	stderr.Reset()
	code := run([]string{"check"}, &stdout, &stderr)
	for _, s := range []string{
		"grove/schemas/bugfix/schema.yaml:2: unknown field \"instrution\"; schema.yaml takes description, documents and skills\n",
		"grove/schemas/bugfix/documents/propsal/template.md: schema bugfix does not use this file\n",
		"grove/schemas/notasks/schema.yaml: tasks must be at changes/<change>/tasks.md; Grove reads it to archive\n",
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
	os.WriteFile(filepath.Join(dir, "grove/specs/record-editing/.grove.yaml"), []byte("title: Record editing\n"), 0o644)
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

	// Implementing a cycle unattended lists its changes in dependency order.
	out = grove(t, "instructions", "implement", k, "--unattended")
	for _, s := range []string{"## What Grove requires of implement\n\n- Work in this working tree on a branch other than the target",
		"## How to implement\n\nBuild these changes in this order", "\n1. " + unknown + " \"Editing an unknown ID prints an error\"\n2. " + del + " \"Add ledger delete ID\"\n"} {
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
		"\ncreated: 20", "\nplanned: " + main + "\n", "## Owner's description\n\n0.4 release\n"} {
		if m == nil || !strings.Contains(out, s) {
			t.Errorf("plan instructions lack %q:\n%s", s, out)
		}
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
	checks(file + ":7: unattended: " + a + " needs a person, as needs_human says, so grove run skips this cycle; remove unattended, or build the cycle with /grove-implement " + k + "\n")
	git(t, dir, "checkout", "-q", "--", ".")

	// A marked cycle whose member changed on main since its plan.
	os.WriteFile(filepath.Join(aDir, "proposal.md"), []byte("# Editing an unknown ID prints an error\n\nRevised.\n"), 0o644)
	commitAll(t, dir, "Revise "+a)
	checks(file + ": planned: " + a + " changed on main since this plan, so grove run skips this cycle; re-plan it with /grove-plan " + k + "\n")

	// A malformed cycle, and one change in two cycles not done.
	os.WriteFile(cycle, []byte(strings.Replace(plan, "planned: "+main, "planned: "+strings.TrimSpace(git(t, dir, "rev-parse", "main")), 1)), 0o644)
	os.WriteFile(filepath.Join(dir, "grove/cycles/261099-zzzzz-later.yaml"), []byte("title: Later\ncreated: 2026-10-23T09:00:00Z\nchanges: ["+a+", 260101-zzzzz]\nstart: tomorrow\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "grove/cycles/261099-zzzzy-bare.yaml"), []byte("title: Bare\n"), 0o644)
	checks("grove/cycles/261099-zzzzy-bare.yaml: created: required field is missing\n" +
		"grove/cycles/261099-zzzzy-bare.yaml: changes: required field is missing\n" +
		"grove/cycles/261099-zzzzy-bare.yaml: planned: required field is missing\n" +
		"grove/cycles/261099-zzzzz-later.yaml:4: unknown field \"start\"; this file takes title, created, changes, planned, goal, needs_human and unattended\n" +
		"grove/cycles/261099-zzzzz-later.yaml: planned: required field is missing\n" +
		"grove/cycles/261099-zzzzz-later.yaml:3: changes: " + a + " is also in cycle " + k + ", which is not done\n" +
		"grove/cycles/261099-zzzzz-later.yaml:3: changes: 260101-zzzzz is no change, open or archived\n")
}

// TestRun follows the walkthrough's "Running unattended" with a
// stand-in claude on PATH, in an operator's clone of the fixture's remote.
func TestRun(t *testing.T) {
	owner := setUp(t)
	remote, checkout, stub := filepath.Join(t.TempDir(), "remote.git"), filepath.Join(t.TempDir(), "checkout"), t.TempDir()
	git(t, owner, "clone", "-q", "--bare", owner, remote)
	git(t, owner, "remote", "add", "origin", remote)
	git(t, owner, "clone", "-q", remote, checkout)
	t.Setenv("PATH", stub+":/usr/bin:/bin")
	// standIn writes claude to record its arguments and process group, then run body.
	standIn := func(body string) {
		os.WriteFile(filepath.Join(stub, "claude"), []byte("#!/bin/sh\necho \"$@\" > "+stub+"/args\necho $$ > "+stub+"/pgid\n"+body+"\n"), 0o755)
	}
	gone := func(when string) { // no process is left in the session's group
		pgid, _ := strconv.Atoi(strings.TrimSpace(read(t, filepath.Join(stub, "pgid"))))
		if err := syscall.Kill(-pgid, 0); err != syscall.ESRCH {
			t.Errorf("%s, the session's process group is still running: %v", when, err)
		}
	}
	runs := func(want string, args ...string) {
		t.Helper()
		t.Chdir(checkout)
		var stdout, stderr strings.Builder
		run(append([]string{"run"}, args...), &stdout, &stderr)
		if got := stdout.String() + stderr.String(); !regexp.MustCompile(want).MatchString(got) {
			t.Errorf("grove run %v printed %q, want %s", args, got, want)
		}
	}
	runs(`^Nothing to apply on main at [0-9a-f]{7}\.\n$`)

	// Nothing to apply: changes in no cycle, a cycle not marked, a marked
	// cycle gone stale.
	t.Chdir(owner)
	a, aDir := approve(t, owner, "a", "Editing an unknown ID prints an error", "")
	b, bDir := approve(t, owner, "b", "Add ledger delete ID", "depends_on: ["+a+"]\n")
	c, cDir := approve(t, owner, "c", "Monthly budgets in report", "")
	nothing := func() {
		t.Helper()
		git(t, owner, "push", "-q", "origin", "main")
		runs(`^Nothing to apply on main at [0-9a-f]{7}\.\n$`)
	}
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

	config := read(t, filepath.Join(owner, "grove.yaml"))
	os.WriteFile(filepath.Join(owner, "grove.yaml"), []byte(strings.Replace(config, "  claude: skills\n", "", 1)), 0o644)
	commitAll(t, owner, "Codex only")
	git(t, owner, "push", "-q", "origin", "main")
	runs(`^grove: run needs Claude Code listed under agents in grove.yaml, as skills or commands\n$`)
	os.WriteFile(filepath.Join(owner, "grove.yaml"), []byte(config), 0o644)
	commitAll(t, owner, "Claude Code again")
	git(t, owner, "push", "-q", "origin", "main")

	// Refusals: a local target of its own, an unclean checkout.
	os.WriteFile(filepath.Join(checkout, "notes.txt"), []byte("mine\n"), 0o644)
	commitAll(t, checkout, "A note of the operator's")
	runs(`^grove: the local main has commits origin's does not, so it cannot be fast-forwarded\n$`)
	git(t, checkout, "reset", "-q", "--hard", "HEAD~1")
	os.WriteFile(filepath.Join(checkout, "README.md"), []byte("edited\n"), 0o644)
	runs(`^grove: the checkout has uncommitted changes; commit, stash or discard them first\n$`)
	git(t, checkout, "checkout", "-q", "--", ".")

	// An ordinary run, with a second run while it goes, its session leaving a
	// child behind and emitting events, which the log holds before it ends.
	events := `{"type":"system","subtype":"init","model":"claude-opus"}
{"type":"assistant","message":{"content":[{"type":"text","text":"Building the cycle.\nFirst A."},{"type":"tool_use","name":"Bash","input":{"description":"Read","command":"grove instructions implement x --unattended\necho"}}]}}
{"type":"user","message":{"content":[{"type":"tool_result","content":"# Grove instructions: implement\nmore"}]}}
{"type":"assistant","message":{"content":[{"type":"thinking","thinking":"a"}]}}
{"type":"assistant","message":{"content":[{"type":"thinking","thinking":"b"}]}}
{"type":"assistant","message":{"content":[{"type":"thinking","thinking":"c"}]}}
{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Read","input":{"file_path":"ledger.py"}}]}}
{"type":"user","message":{"content":[{"type":"tool_result","is_error":true,"content":[{"type":"text","text":"File does not exist.\n"}]}]}}
{"type":"result","subtype":"success","result":"Done.","num_turns":3,"total_cost_usd":0.4213}
`
	os.WriteFile(filepath.Join(stub, "events"), []byte(events), 0o644)
	standIn("sleep 30 &\ncat " + stub + "/events\ntouch " + stub + "/started\nwhile [ ! -f " + stub + "/release ]; do sleep 0.05; done")
	first := make(chan string)
	go func() {
		var stdout, stderr strings.Builder
		run([]string{"run"}, &stdout, &stderr)
		first <- stdout.String() + stderr.String()
	}()
	for _, err := os.Stat(filepath.Join(stub, "started")); err != nil; _, err = os.Stat(filepath.Join(stub, "started")) {
		time.Sleep(20 * time.Millisecond)
	}
	runs(`^grove: a run is in progress in this checkout\n$`)
	logged := func() string {
		logs, _ := filepath.Glob(filepath.Join(checkout, ".git/grove/*.log"))
		if len(logs) == 0 {
			return ""
		}
		return read(t, logs[0])
	}
	for end := time.Now().Add(5 * time.Second); logged() != events && time.Now().Before(end); {
		time.Sleep(20 * time.Millisecond)
	}
	if got := logged(); got != events {
		t.Errorf("while the session ran, its log held %q", got)
	}
	os.WriteFile(filepath.Join(stub, "release"), nil, 0o644)
	if got := <-first; !regexp.MustCompile(`^[0-9:]{8} ▶ started: model claude-opus\n[0-9:]{8} • Building the cycle\.\n[0-9:]{8} › Bash grove instructions implement x --unattended\n` +
		`[0-9:]{8} · thinking \(×3\)\n[0-9:]{8} › Read ledger\.py\n[0-9:]{8} ✗ File does not exist\.\n[0-9:]{8} ■ result: success, 3 turns, \$0\.42\nAttempt [0-9]{6}-[0-9a-z]{5} ended: exit status 0; its output is .*/\.git/grove/[0-9]{6}-[0-9a-z]{5}\.log\.\n$`).MatchString(got) {
		t.Errorf("the first run printed %q", got)
	}
	gone("after the session exited")
	main := strings.TrimSpace(git(t, owner, "rev-parse", "main"))
	var rec attempt
	recs, _ := filepath.Glob(filepath.Join(checkout, ".git/grove/*.json"))
	json.Unmarshal([]byte(read(t, recs[0])), &rec)
	if read(t, filepath.Join(stub, "args")) != "-p /grove-implement "+k+" --unattended --output-format stream-json --verbose\n" || logged() != events || rec.Cycle != k || strings.Join(rec.Changes, " ") != a+" "+b || rec.Commit != main || rec.Outcome != "exit status 0" ||
		strings.TrimSpace(git(t, checkout, "branch", "--show-current")) != rec.Branch || strings.TrimSpace(git(t, checkout, "rev-parse", "main")) != main || git(t, checkout, "status", "--porcelain") != "" {
		t.Errorf("after the first run: %+v, %s", rec, git(t, checkout, "status"))
	}

	// Paused, though another cycle is marked, until the cycle is archived.
	runs(`^grove: the project is paused after attempt ` + rec.ID + `, whose output is ` + regexp.QuoteMeta(rec.Output) + `; once you have dealt with it, run grove run --retry\n$`)
	t.Chdir(owner)
	os.WriteFile(filepath.Join(aDir, "implementation.md"), []byte("## T-01\n\nStatus: done\n"), 0o644)
	os.WriteFile(filepath.Join(bDir, "implementation.md"), []byte("Status: dropped\n"), 0o644)
	maintain(t, owner, olderEdits)
	grove(t, "archive", a)
	grove(t, "archive", b)
	commitAll(t, owner, "Archive the batch")
	git(t, owner, "push", "-q", "origin", "main")
	standIn("")
	runs(`^Attempt .* ended: exit status 0;`)
	if recs, _ = filepath.Glob(filepath.Join(checkout, ".git/grove/*.json")); len(recs) != 2 {
		t.Fatalf("records: %v", recs)
	}
	if last, _ := latest(filepath.Join(checkout, ".git/grove")); last.Cycle != l || strings.Join(last.Changes, " ") != c {
		t.Errorf("the run after the merge selected %s %v, want %s %s", last.Cycle, last.Changes, l, c)
	}

	// Failed to start, then paused.
	os.Remove(filepath.Join(stub, "claude"))
	runs(`^grove: attempt .* ended: failed to start;`, "--retry")
	runs(`^grove: the project is paused after attempt `)

	// Cancelled after a visible edit; retry refuses until it is dealt with.
	standIn("echo edit >> README.md\ntouch " + stub + "/started\nsleep 30")
	os.Remove(filepath.Join(stub, "started"))
	go func() {
		for _, err := os.Stat(filepath.Join(stub, "started")); err != nil; _, err = os.Stat(filepath.Join(stub, "started")) {
			time.Sleep(20 * time.Millisecond)
		}
		syscall.Kill(os.Getpid(), syscall.SIGINT)
	}()
	runs(`^grove: attempt .* ended: cancelled;`, "--retry")
	gone("after cancelling")
	runs(`^grove: the checkout has uncommitted changes;`, "--retry")
	if !strings.HasSuffix(read(t, filepath.Join(checkout, "README.md")), "edit\n") {
		t.Error("the interrupted edit is gone")
	}
	git(t, checkout, "checkout", "-q", "--", ".")

	// Timed out, with a child that ignores the request to stop.
	standIn("sh -c 'trap \"\" TERM; sleep 30' &\nsleep 30")
	runs(`^grove: attempt .* ended: timed out;`, "--retry", "--timeout", "1s")
	gone("after the time limit")

	// A cycle given: built though not marked while l is, and refused, naming
	// why, when there is no such cycle or a member needs a person.
	t.Chdir(owner)
	d, _ := approve(t, owner, "d", "Show saved records", "")
	e, _ := approve(t, owner, "e", "Export records", "")
	m, n := "261023-m0000", "261023-n0000"
	planCycle(t, owner, m, "2026-10-23T09:00:00Z", false, d)
	planCycle(t, owner, n, "2026-10-23T10:00:00Z", false, e)
	nFile := filepath.Join(owner, "grove/cycles", n+"-cycle.yaml")
	os.WriteFile(nFile, []byte(read(t, nFile)+"needs_human:\n  "+e+": judging the export\n"), 0o644)
	os.WriteFile(filepath.Join(owner, "grove.yaml"), []byte(strings.Replace(config, "claude: skills", "claude: commands", 1)), 0o644)
	commitAll(t, owner, "Export needs a person, and Claude Code as commands")
	git(t, owner, "push", "-q", "origin", "main")
	branches := git(t, checkout, "branch")
	runs(`^grove: main at [0-9a-f]{7} holds no cycle 261099-zzzzz\n$`, "--retry", "261099-zzzzz")
	runs(`^grove: `+e+` needs a person, as cycle `+n+`'s needs_human says; build the cycle with /grove-implement `+n+`\n$`, "--retry", n)
	if got := git(t, checkout, "branch"); got != branches {
		t.Errorf("a refused run left branches %q, want %q", got, branches)
	}
	runs(`^Attempt .* ended: exit status 0;`, "--retry", m)
	if last, _ := latest(filepath.Join(checkout, ".git/grove")); last.Cycle != m || read(t, filepath.Join(stub, "args")) != "-p /grove:implement "+m+" --unattended --output-format stream-json --verbose\n" {
		t.Errorf("the run given %s selected %s", m, last.Cycle)
	}

	// No remote: the local target, commits of its own and all, is the fetched
	// one, and merging the attempt's branch into it ends the pause.
	git(t, checkout, "remote", "remove", "origin")
	git(t, checkout, "switch", "-q", "main")
	os.WriteFile(filepath.Join(checkout, "notes.txt"), []byte("mine\n"), 0o644)
	commitAll(t, checkout, "A note of the operator's")
	local := strings.TrimSpace(git(t, checkout, "rev-parse", "main"))
	standIn("")
	runs(`^Attempt .* ended: exit status 0;`, "--retry")
	if last, _ := latest(filepath.Join(checkout, ".git/grove")); last.Commit != local || last.Cycle != l {
		t.Errorf("the run with no remote recorded %s %s, want %s %s", last.Cycle, last.Commit, l, local)
	}
	runs(`^grove: the project is paused after attempt `)
	cRel, _ := filepath.Rel(owner, cDir)
	os.WriteFile(filepath.Join(checkout, cRel, "implementation.md"), []byte("Status: dropped\n"), 0o644)
	grove(t, "archive", c)
	commitAll(t, checkout, "Drop "+c)
	branch := strings.TrimSpace(git(t, checkout, "branch", "--show-current"))
	git(t, checkout, "switch", "-q", "main")
	git(t, checkout, "merge", "-q", "--no-ff", "-m", "Merge the attempt", branch)
	runs(`^Nothing to apply on main at [0-9a-f]{7}\.\n$`)
}

// TestMovedFiles follows the walkthrough's steps on a task's named file
// changing on main: a change alone reported against its base, a member of
// a marked cycle against its plan, and a Files: line Grove cannot read.
func TestMovedFiles(t *testing.T) {
	dir := setUp(t)
	base := strings.TrimSpace(git(t, dir, "rev-parse", "main"))
	cid, folder := propose(t, dir, "older-edits", "Allow editing older records", "base: "+base+"\n")
	tasks := "## T-01 — Allow editing older records\n\nResult: older records can be edited.\n\nFiles: `ledger/cli.py`,\n`tests/test_cli.py`\n\nAcceptance:\n- An older record saves.\n"
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

	// Unchanged on main, or changed only on the branch building it.
	checks("No problems in grove/.\n")
	git(t, dir, "switch", "-q", "-c", "older-edits-impl")
	edit("ledger/cli.py", "Build T-01")
	checks("No problems in grove/.\n")
	git(t, dir, "switch", "-q", "main")

	// Changed on main since base.
	edit("tests/test_cli.py", "Another change")
	rel := "grove/changes/" + filepath.Base(folder) + "/tasks.md"
	checks(rel + ": T-01's tests/test_cli.py changed on main since base " + base[:7] + "; revise the change with propose\n")
	for _, args := range [][]string{{cid}, {cid, "--unattended"}} {
		if out := grove(t, append([]string{"instructions", "implement"}, args...)...); !strings.Contains(out, "\n- Changed on main since base "+base[:7]+": T-01's `tests/test_cli.py`.\n") {
			t.Errorf("implement %v does not name the moved file:\n%.3000s", args, out)
		}
	}

	// A member of a marked cycle: reported against the plan, not base.
	k := "261023-k0000"
	planCycle(t, dir, k, "2026-10-23T09:00:00Z", true, cid)
	checks("No problems in grove/.\n")
	edit("ledger/cli.py", "Yet another change")
	checks("grove/cycles/" + k + "-cycle.yaml: planned: " + cid + " changed on main since this plan, so grove run skips this cycle; re-plan it with /grove-plan " + k + "\n")
	p, _ := openProject()
	if ip, err := p.agreed(); err != nil || func() bool { cy, _ := eligible(p, ip); return cy != nil }() {
		t.Errorf("grove run would build a cycle whose member's file moved: %v", err)
	}
	if out := grove(t, "instructions", "implement", cid); strings.Contains(out, "since base") {
		t.Errorf("implement instructions for a member name a moved file against base:\n%.3000s", out)
	}

	// A Files: line Grove cannot read.
	os.WriteFile(filepath.Join(folder, "tasks.md"), []byte(strings.Replace(tasks, "`ledger/cli.py`,\n`tests/test_cli.py`", "ledger/cli.py", 1)), 0o644)
	os.Remove(filepath.Join(dir, "grove/cycles", k+"-cycle.yaml"))
	checks(rel + ":5: T-01's Files: line should list repo-relative paths in backticks, separated by commas\n")
}
