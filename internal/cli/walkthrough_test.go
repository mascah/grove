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
		for _, w := range []string{"explore", "propose", "plan", "apply", "archive"} {
			if skill := read(t, filepath.Join(dir, h, "skills", "grove-"+w, "SKILL.md")); !strings.Contains(skill, "name: grove-"+w+"\n") || !strings.Contains(skill, " instructions "+w) {
				t.Errorf("%s skill grove-%s:\n%s", h, w, skill)
			}
		}
	}
	if skill := read(t, filepath.Join(dir, ".agents/skills/grove-propose/SKILL.md")); !strings.Contains(skill, ` instructions propose "<args>"`+"`") {
		t.Errorf("codex propose skill:\n%s", skill)
	}
	yaml := strings.Replace(read(t, filepath.Join(dir, "grove.yaml")), `context: ""`, "context: Python 3 standard library only.", 1)
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
	for _, s := range []string{"## Project context\n\nPython 3 standard library only.\n", "The project's own rules for this document:\n\n- Name the owner for every decision.\n",
		"The project's own rules for this skill:\n\n- Commit the change folder when check passes.\n"} {
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
	config := read(t, filepath.Join(dir, "grove.yaml"))
	os.WriteFile(filepath.Join(dir, "grove.yaml"), []byte(strings.Replace(config, "rules: {}", "rules:\n  derp: [x]\n  documents:\n    propsal: [x]\nrule: {}", 1)), 0o644)
	var stdout, stderr strings.Builder
	code := run([]string{"check"}, &stdout, &stderr)
	want := "grove.yaml:25: unknown rules key \"derp\"; rules takes documents and skills\n" +
		"grove.yaml:27: unknown rules.documents key \"propsal\"; kinds are proposal, change-spec, tasks, implementation, spec\n" +
		"grove.yaml:28: unknown field \"rule\"; grove.yaml takes root, target, agents, schema, context and rules\n"
	if code != 1 || stdout.String() != want || stderr.String() != "grove: 3 problems found\n" {
		t.Errorf("check of grove.yaml: status %d\n%s%s", code, stdout.String(), stderr.String())
	}
	grove(t, "changes")
	os.WriteFile(filepath.Join(dir, "grove.yaml"), []byte(config), 0o644)

	os.MkdirAll(filepath.Join(dir, "grove", "specs", "record-editing"), 0o755)
	os.WriteFile(filepath.Join(dir, "grove", "specs", "record-editing", "spec.md"), []byte("# Record editing\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "grove", "specs", "record-editing", ".grove.yaml"), []byte("description: Editing records.\nowner: me\n"), 0o644)
	for folder, meta := range map[string]string{"261015-k3m9q-no-proposal": "title: T\ncreated: x\n", "261015-p4x8d-extra-field": "title: T\ncreated: x\nowner: me\n",
		"261015-q7r2s-drifted": "title: New title\ncreated: x\n", "no-id": "title: T\ncreated: x\n"} {
		os.MkdirAll(filepath.Join(dir, "grove", "changes", folder), 0o755)
		os.WriteFile(filepath.Join(dir, "grove", "changes", folder, ".grove.yaml"), []byte(meta), 0o644)
		if folder != "261015-k3m9q-no-proposal" {
			os.WriteFile(filepath.Join(dir, "grove", "changes", folder, "proposal.md"), []byte("# T\n"), 0o644)
		}
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"check"}, &stdout, &stderr)
	want = "grove/specs/record-editing/.grove.yaml:2: unknown field \"owner\"; this file takes title and description\n" +
		"grove/specs/record-editing/.grove.yaml: title: required field is missing\n" +
		"grove/changes/261015-k3m9q-no-proposal/proposal.md: missing; it is the change's main document in schema default\n" +
		"grove/changes/261015-p4x8d-extra-field/.grove.yaml:3: unknown field \"owner\"; this file takes title, created, schema, description, depends_on, cycle and base\n" +
		"grove/changes/261015-q7r2s-drifted/proposal.md:1: the first heading should be \"# New title\", the title in .grove.yaml\n" +
		"grove/changes/no-id/: the folder name must start with a change ID, such as 261015-k3m9q\n"
	if code != 1 || stdout.String() != want || stderr.String() != "grove: 6 problems found\n" {
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
	os.WriteFile(filepath.Join(folder, "specs", "record-editing", "spec.md"), []byte("# Changes to record-editing\n\n## ADDED Requirements\n\n### Requirement: Edit older records\n"), 0o644)
	os.WriteFile(filepath.Join(folder, "tasks.md"), []byte("## T-01 — "+title+"\n"), 0o644)
	return m[1], folder
}

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
	if !strings.Contains(out, "## Owner's description\n\nEditing a missing record should print an error\n") || req < 0 || how < req {
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

func TestApplyInstructions(t *testing.T) {
	dir := fixture(t)
	grove(t, "init")
	if !strings.Contains(read(t, filepath.Join(dir, ".claude/skills/grove-apply/SKILL.md")), "Run `grove instructions apply $ARGUMENTS`") {
		t.Error("the skill does not run instructions apply")
	}
	commitAll(t, dir, "Set up Grove")
	cid, _ := propose(t, dir, "older-edits", "Allow editing older records", "")
	commitAll(t, dir, "Propose")

	var stdout, stderr strings.Builder
	if code := run([]string{"instructions", "apply", cid}, &stdout, &stderr); code != 1 || stderr.String() != "grove: main holds no approved change "+cid+"\n" {
		t.Errorf("before the merge: status %d, %q", code, stderr.String())
	}
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Approve", "older-edits")
	out := grove(t, "instructions", "apply", cid)
	for _, s := range []string{"/implementation.md\">\n# Implementation: Allow editing older records", "### Requirement: Edit older records"} {
		if !strings.Contains(out, s) {
			t.Errorf("apply instructions lack %q", s)
		}
	}
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
	for _, s := range []string{"Reconcile the specs", "Status: done", "### Requirement: Edit older records"} {
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
		filled := strings.NewReplacer("<capability>", "record-editing", "<title>", "Record editing").Replace(m[1])
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
	for skill, s := range map[string]string{"apply": "One task at a time", "archive": "Reconcile the specs"} {
		if out := grove(t, "instructions", skill, cid); !strings.Contains(out, ", with change schema bugfix, ") || !strings.Contains(out, s) {
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
// merges it.
func finish(t *testing.T, dir, branch, id, folder, spec string) {
	t.Helper()
	git(t, dir, "switch", "-q", "-c", branch)
	os.WriteFile(filepath.Join(folder, "implementation.md"), []byte("## T-01\n\nStatus: done\n"), 0o644)
	os.MkdirAll(filepath.Join(dir, "grove/specs/record-editing"), 0o755)
	os.WriteFile(filepath.Join(dir, "grove/specs/record-editing/.grove.yaml"), []byte("title: Record editing\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "grove/specs/record-editing/spec.md"), []byte("# Record editing\n\n## Requirements\n\n"+spec), 0o644)
	if out := grove(t, "check"); out != "No problems in grove/.\n" {
		t.Errorf("check on %s printed %q", branch, out)
	}
	grove(t, "archive", id)
	commitAll(t, dir, "Implement and archive "+id)
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Merge "+branch, branch)
}

// TestCoordinating follows the walkthrough's "Coordinating changes".
func TestCoordinating(t *testing.T) {
	dir := setUp(t)
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

	// Apply refuses an unmet dependency.
	var stdout, stderr strings.Builder
	if code := run([]string{"instructions", "apply", del}, &stdout, &stderr); code != 1 || stderr.String() != "grove: "+del+" waits on "+unknown+", which main has not archived as finished\n" {
		t.Errorf("apply before its dependency is archived: status %d, %q", code, stderr.String())
	}
	keep := "### Requirement: Keep a category\n\nWhen a user edits a record and leaves its category unchanged, the system shall save the edit.\n"
	finish(t, dir, "unknown-id-impl", unknown, unknownDir, keep)
	grove(t, "instructions", "apply", del)

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

	// Cycles and /grove-plan.
	budgetsMeta := filepath.Join(dir, "grove/changes", budgets+"-monthly-budgets-in-report/.grove.yaml")
	os.WriteFile(budgetsMeta, []byte("title: Monthly budgets in report\ncreated: 2026-10-14T00:00:00Z\n"), 0o644)
	commitAll(t, dir, "Propose monthly budgets")
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Approve budgets", "budgets")
	out = grove(t, "instructions", "plan", "0.4: the README's planned list")
	m := regexp.MustCompile("### cycle: `grove/cycles/([0-9]{6}-[0-9a-z]{5})-<slug>.yaml`").FindStringSubmatch(out)
	if m == nil || !strings.Contains(out, "## What Grove requires of plan\n\n- A cycle groups approved changes.") || !strings.Contains(out, "## How to plan\n\n**Plan groups") ||
		!strings.Contains(out, "\">\ntitle: <title>\n") || !strings.Contains(out, "## Owner's description\n\n0.4: the README's planned list\n") {
		t.Fatalf("plan instructions:\n%s", out)
	}
	cycle := m[1]
	git(t, dir, "switch", "-q", "-c", "plan-04")
	cycleFile := filepath.Join(dir, "grove/cycles", cycle+"-0-4.yaml")
	os.MkdirAll(filepath.Dir(cycleFile), 0o755)
	os.WriteFile(cycleFile, []byte("title: \"0.4\"\ngoal: The README's planned list.\n"), 0o644)
	for _, meta := range []string{delMeta, budgetsMeta} {
		os.WriteFile(meta, []byte(read(t, meta)+"cycle: "+cycle+"\n"), 0o644)
	}
	if got := grove(t, "check"); got != "No problems in grove/.\n" {
		t.Errorf("check printed %q", got)
	}
	commitAll(t, dir, "Plan 0.4")
	if out := grove(t, "review"); !strings.Contains(out, "  adds     cycle            grove/cycles/"+cycle+"-0-4.yaml\n") {
		t.Errorf("review of the plan:\n%.1500s", out)
	}
	if out := grove(t, "instructions", "plan", cycle); !strings.Contains(out, "You are revising cycle "+cycle+", \"0.4\"") ||
		!strings.Contains(out, "## The cycle\n\n<file path=\"grove/cycles/"+cycle+"-0-4.yaml\">\ntitle: \"0.4\"") || !strings.Contains(out, "<file path=\"grove/changes/"+filepath.Base(delDir)+"/.grove.yaml\">") {
		t.Errorf("plan instructions for revising the cycle:\n%s", out)
	}
	os.WriteFile(cycleFile, []byte("goal: x\nowner: me\n"), 0o644)
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"check"}, &stdout, &stderr)
	if want := "grove/cycles/" + cycle + "-0-4.yaml:2: unknown field \"owner\"; this file takes title, goal, start and end\ngrove/cycles/" + cycle + "-0-4.yaml: title: required field is missing\n"; code != 1 || stdout.String() != want {
		t.Errorf("check of a cycle file: status %d\n%s", code, stdout.String())
	}
	git(t, dir, "checkout", "--", ".")
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Accept the 0.4 plan", "plan-04")

	// Lists by cycle; propose asks for a cycle.
	if got := grove(t, "changes"); !strings.HasPrefix(got, cycle+"  0.4  0 of 2 archived\n  ") || !strings.Contains(got, "\n  "+del+"  Add ledger delete ID  approved\n") {
		t.Errorf("changes printed %q", got)
	}
	if out := grove(t, "instructions", "propose", "x"); !strings.Contains(out, "Open cycles (planned work, not yet all archived):\n- "+cycle+" \"0.4\": `grove/cycles/"+cycle+"-0-4.yaml`\n  - goal: ") ||
		!strings.Contains(out, "ask the owner\n  which, if any, the change joins") {
		t.Errorf("propose instructions do not offer the cycle:\n%s", out)
	}
	os.WriteFile(budgetsMeta, []byte(strings.Replace(read(t, budgetsMeta), "cycle: "+cycle, "cycle: 260101-zzzzz", 1)), 0o644)
	stdout.Reset()
	if code := run([]string{"check"}, &stdout, &stderr); code != 1 || stdout.String() != "grove/changes/"+budgets+"-monthly-budgets-in-report/.grove.yaml:3: cycle: 260101-zzzzz names no file in cycles/\n" {
		t.Errorf("check of a cycle naming no file: status %d\n%s", code, stdout.String())
	}
	git(t, dir, "checkout", "--", ".")

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
	modify(showDir, "When a user edits a record and leaves its category unchanged, the system shall save the edit and print what was saved.\n")
	commitAll(t, dir, "Propose showing what was saved")
	git(t, dir, "switch", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Approve showing what was saved", "show-saved")
	finish(t, dir, "recategorize-impl", recat, recatDir, "### Requirement: Keep a category\n\n"+recatText)
	stdout.Reset()
	want := "grove/changes/" + filepath.Base(showDir) + "/specs/record-editing/spec.md:5: requirement \"Keep a category\" was rewritten on main since base " + base[:7] + "; revise the change with propose\n"
	if code := run([]string{"check"}, &stdout, &stderr); code != 1 || stdout.String() != want {
		t.Errorf("check of a stale change: status %d\n%s", code, stdout.String())
	}
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
}
