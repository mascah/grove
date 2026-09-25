package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/mascah/grove"
	"github.com/mascah/grove/internal/repo"
)

// emptyRepo is a Git checkout with one commit and no Grove content.
func emptyRepo(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	write(t, root, "README.md", "# A project\n")
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"add", "-A"}, {"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-c", "maintenance.auto=false", "commit", "-q", "-m", "init"}} {
		if out, err := repo.Command(context.Background(), root, args...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	return root
}

func runInitAt(t *testing.T, root string) (int, string, string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code := Run([]string{"--project", root, "init"}, t.TempDir(), &out, &errOut)
	return code, out.String(), errOut.String()
}

func TestInitCreatesAProjectAndRerunsWithoutTouchingUserFiles(t *testing.T) {
	t.Parallel()
	root := emptyRepo(t)
	code, out, errOut := runInitAt(t, root)
	if code != 0 {
		t.Fatal(errOut)
	}
	want := "created grove.yaml\ncreated grove\ncreated grove/brief.md (a placeholder that states no intent)\n" +
		"created .agents/skills/grove-shape/SKILL.md\ncreated .agents/skills/grove-shape/agents/openai.yaml\n" +
		"created .agents/skills/grove-work/SKILL.md\ncreated .agents/skills/grove-work/agents/openai.yaml\n" +
		"created .claude/agents/grove-reviewer.md\n" +
		"created .claude/skills/grove-shape/SKILL.md\ncreated .claude/skills/grove-work/SKILL.md\n"
	if out != want || !strings.Contains(errOut, "Project: "+root) {
		t.Fatalf("stdout=%q stderr=%q", out, errOut)
	}
	adapter, err := os.ReadFile(filepath.Join(root, ".claude/skills/grove-work/SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, needle := range []string{"disable-model-invocation: true", grove.ManagedMarker, "Assignment: $ARGUMENTS", "grove guide work --entrypoint 2", "grove entrypoint revision 2", "stop and say"} {
		if !strings.Contains(string(adapter), needle) {
			t.Fatalf("the Claude adapter lacks %q:\n%s", needle, adapter)
		}
	}
	portable := map[string]string{".claude/skills/grove-work/SKILL.md": string(adapter)}
	for _, relative := range []string{".claude/agents/grove-reviewer.md", ".claude/skills/grove-shape/SKILL.md", ".agents/skills/grove-work/SKILL.md", ".agents/skills/grove-shape/SKILL.md", ".agents/skills/grove-work/agents/openai.yaml"} {
		source, err := os.ReadFile(filepath.Join(root, relative))
		if err != nil {
			t.Fatal(err)
		}
		portable[relative] = string(source)
	}
	for _, name := range []string{"work", "shape", "review", "model"} {
		var guide, guideErr bytes.Buffer
		if code := Run([]string{"guide", name}, t.TempDir(), &guide, &guideErr); code != 0 {
			t.Fatal(guideErr.String())
		}
		portable["guide "+name] = guide.String()
	}
	for name, text := range portable {
		for _, forbidden := range []string{"go run", "docs/work-execution.md", "docs/work-shaping.md", "docs/work-review.md", "../grove/", "../.claude/", "../.agents/", ".claude/worktrees", "worktree-G-", "AGENTS.md` here"} {
			if strings.Contains(text, forbidden) {
				t.Fatalf("%s must not depend on Grove's own repository, but mentions %q", name, forbidden)
			}
		}
	}
	if !strings.Contains(errOut, "Next: grove check, then commit what init wrote") {
		t.Fatalf("init must say what the adopter's instructions may add: %q", errOut)
	}
	var checkOut, checkErr bytes.Buffer
	if code := Run([]string{"--project", root, "check"}, t.TempDir(), &checkOut, &checkErr); code != 0 || checkOut.String() != "OK: 0 records\n" {
		t.Fatalf("the initialized project must validate: %d %s %s", code, checkOut.String(), checkErr.String())
	}
	checkOut.Reset()
	if code := Run([]string{"--project", root, "new", "work", "First"}, t.TempDir(), &checkOut, &checkErr); code != 0 || checkOut.String() != "grove/G-001-first.md\n" {
		t.Fatalf("new must work in the initialized project: %d %s %s", code, checkOut.String(), checkErr.String())
	}

	// The user develops the brief and edits one managed file without giving
	// up the marker; another they take over; the rest stay as written.
	write(t, root, "grove/brief.md", "# Real brief\n\nWritten by a person.\n")
	edited := strings.Replace(string(adapter), "stop and say", "stop", 1)
	write(t, root, ".claude/skills/grove-work/SKILL.md", edited)
	write(t, root, ".agents/skills/grove-work/SKILL.md", "my own instructions\n")
	before := hashes(t, root)
	code, out, errOut = runInitAt(t, root)
	if code != 0 {
		t.Fatal(errOut)
	}
	want = "kept grove.yaml (exists and validates)\nkept grove\nkept grove/brief.md (never rewritten)\n" +
		"unchanged .agents/skills/grove-shape/SKILL.md\nunchanged .agents/skills/grove-shape/agents/openai.yaml\n" +
		"kept .agents/skills/grove-work/SKILL.md (not managed by grove init; delete it to get the managed version)\n" +
		"unchanged .agents/skills/grove-work/agents/openai.yaml\n" +
		"unchanged .claude/agents/grove-reviewer.md\n" +
		"unchanged .claude/skills/grove-shape/SKILL.md\nupdated .claude/skills/grove-work/SKILL.md\n"
	if out != want {
		t.Fatalf("second run:\n%s", out)
	}
	after := hashes(t, root)
	updated := filepath.Join(root, ".claude/skills/grove-work/SKILL.md")
	if after[updated] == before[updated] {
		t.Fatal("the marked file must be rewritten from the template")
	}
	delete(before, updated)
	delete(after, updated)
	if !reflect.DeepEqual(before, after) {
		t.Fatal("a rerun must change nothing but the managed file whose template changed")
	}
	if restored, _ := os.ReadFile(updated); string(restored) != string(adapter) {
		t.Fatal("the managed update must restore the template exactly")
	}
}

func TestInitRespectsAnExistingConfiguration(t *testing.T) {
	t.Parallel()
	root := emptyRepo(t)
	write(t, root, "grove.yaml", "schema_version: 3\nrecords: docs/records\n")
	code, out, errOut := runInitAt(t, root)
	if code != 0 || !strings.HasPrefix(out, "kept grove.yaml (exists and validates)\ncreated docs/records\ncreated .agents/") {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if _, err := os.Stat(filepath.Join(root, "grove")); err == nil {
		t.Fatal("init must follow the configured record root, not its default")
	}
}

func TestInitRefusesConflictsWithoutWriting(t *testing.T) {
	t.Parallel()
	cases := map[string]func(root string){
		"unsupported configuration": func(root string) { write(t, root, "grove.yaml", "schema_version: 2\nrecords: work\n") },
		"record root is a file":     func(root string) { write(t, root, "grove", "not a directory\n") },
		"managed path is a directory": func(root string) {
			write(t, root, ".claude/skills/grove-work/SKILL.md/inner", "x")
		},
	}
	for name, arrange := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			root := emptyRepo(t)
			arrange(root)
			before := hashes(t, root)
			code, out, errOut := runInitAt(t, root)
			if code != 1 || out != "" || !strings.Contains(errOut, "grove: conflict ") || !strings.Contains(errOut, "nothing was written") {
				t.Fatalf("code=%d stdout=%q stderr=%q", code, out, errOut)
			}
			if !reflect.DeepEqual(before, hashes(t, root)) {
				t.Fatal("a conflict must leave every file as it was")
			}
		})
	}
	t.Run("below the checkout top, from the working directory", func(t *testing.T) {
		t.Parallel()
		root := emptyRepo(t)
		nested := filepath.Join(root, "sub")
		if err := os.Mkdir(nested, 0o755); err != nil {
			t.Fatal(err)
		}
		before := hashes(t, root)
		var out, errOut bytes.Buffer
		if code := Run([]string{"init"}, nested, &out, &errOut); code != 1 || out.Len() != 0 || !strings.Contains(errOut.String(), "below it, at sub") || !reflect.DeepEqual(before, hashes(t, root)) {
			t.Fatalf("code=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
		}
		out.Reset()
		errOut.Reset()
		if code := Run([]string{"init"}, root, &out, &errOut); code != 0 || !strings.HasPrefix(out.String(), "created grove.yaml\n") {
			t.Fatalf("init from the checkout top without --project: code=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
		}
	})
	t.Run("symlinked parent of the record root", func(t *testing.T) {
		t.Parallel()
		root := emptyRepo(t)
		elsewhere := t.TempDir()
		write(t, root, "grove.yaml", "schema_version: 3\nrecords: docs/records\n")
		if err := os.Symlink(elsewhere, filepath.Join(root, "docs")); err != nil {
			t.Skip("symlinks unavailable")
		}
		code, out, errOut := runInitAt(t, root)
		if code != 1 || out != "" || !strings.Contains(errOut, "docs/records: docs is a symlink") {
			t.Fatalf("code=%d stdout=%q stderr=%q", code, out, errOut)
		}
		if entries, _ := os.ReadDir(elsewhere); len(entries) != 0 {
			t.Fatal("the record root may not be created through the symlink")
		}
	})
	t.Run("symlinked parent of a managed path", func(t *testing.T) {
		t.Parallel()
		root := emptyRepo(t)
		elsewhere := t.TempDir()
		if err := os.Symlink(elsewhere, filepath.Join(root, ".claude")); err != nil {
			t.Skip("symlinks unavailable")
		}
		code, out, errOut := runInitAt(t, root)
		if code != 1 || out != "" || !strings.Contains(errOut, ".claude is a symlink") {
			t.Fatalf("code=%d stdout=%q stderr=%q", code, out, errOut)
		}
		if entries, _ := os.ReadDir(elsewhere); len(entries) != 0 {
			t.Fatal("nothing may be written through the symlink")
		}
		if _, err := os.Stat(filepath.Join(root, "grove.yaml")); err == nil {
			t.Fatal("a conflict must write nothing")
		}
	})
	t.Run("outside Git", func(t *testing.T) {
		t.Parallel()
		code, out, errOut := runInitAt(t, t.TempDir())
		if code != 1 || out != "" || !strings.Contains(errOut, "top of a Git checkout") {
			t.Fatalf("code=%d stdout=%q stderr=%q", code, out, errOut)
		}
	})
}

func TestGuideAndVersionNeedNoProject(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		args   []string
		prefix string
	}{
		{[]string{"guide", "work"}, "# Executing assigned Grove work\n"},
		{[]string{"guide", "shape"}, "# Shaping Grove work\n"},
		{[]string{"guide", "model"}, "# Record model\n"},
		{[]string{"guide", "review"}, "# Reviewing Grove work\n"},
		{[]string{"guide", "work", "--entrypoint", "2"}, "# Executing assigned Grove work\n"},
		{[]string{"guide", "shape", "--entrypoint=2"}, "# Shaping Grove work\n"},
		{[]string{"version"}, "grove "}, // ends with the content digests, checked below
	} {
		var out, errOut bytes.Buffer
		if code := Run(c.args, t.TempDir(), &out, &errOut); code != 0 || !strings.HasPrefix(out.String(), c.prefix) {
			t.Fatalf("%v: code=%d stdout=%q stderr=%q", c.args, code, out.String(), errOut.String())
		}
	}
	var version, versionErr bytes.Buffer
	if code := Run([]string{"version"}, t.TempDir(), &version, &versionErr); code != 0 || version.String() != grove.Identity().String()+"\n" || !regexp.MustCompile(`^grove \S.* guides sha256:[0-9a-f]{12} content sha256:[0-9a-f]{12}\n$`).MatchString(version.String()) {
		t.Fatalf("version=%q", version.String())
	}
	// The shipped documents reach projects that have no docs folder and live
	// G- IDs of their own, so none links outside itself, names a G- ID beyond
	// its own examples, or points at Grove's repository or the predecessor.
	shipped := map[string]string{}
	examples := map[string][]string{"work": {"G-030", "G-031"}, "shape": {"G-037"}, "review": nil, "model": {"G-001", "G-003", "G-1000"}}
	for name := range examples {
		var out, errOut bytes.Buffer
		Run([]string{"guide", name}, t.TempDir(), &out, &errOut)
		shipped[name] = out.String()
	}
	for name, text := range shipped {
		for _, link := range regexp.MustCompile(`\]\(([^)]*)\)`).FindAllStringSubmatch(text, -1) {
			if !strings.HasPrefix(link[1], "#") && !strings.HasPrefix(link[1], "https://") {
				t.Errorf("%s links outside itself: %s", name, link[1])
			}
		}
		for _, id := range regexp.MustCompile(`G-[0-9]+`).FindAllString(text, -1) {
			if !slices.Contains(examples[name], id) {
				t.Errorf("%s names %s, which is a live ID in an adopting project", name, id)
			}
		}
		// Hard-wrapped prose splits a phrase across lines as often as not.
		prose := strings.ReplaceAll(strings.Join(strings.Fields(strings.ToLower(text)), " "), "’", "'")
		for _, phrase := range []string{"grove's own repository", "grove's repository", "grove's own records", "command reference", "predecessor"} {
			if strings.Contains(prose, phrase) {
				t.Errorf("%s mentions %q, which an adopting project lacks", name, phrase)
			}
		}
	}
	for _, args := range [][]string{{"guide"}, {"guide", "both"}, {"guide", "work", "shape"}, {"version", "x"}, {"init", "here"}} {
		var out, errOut bytes.Buffer
		if code := Run(args, t.TempDir(), &out, &errOut); code != 2 || out.Len() != 0 {
			t.Fatalf("%v: code=%d stdout=%q", args, code, out.String())
		}
	}
}

// revisionOne is what init wrote before entrypoint revisions, read from
// testdata/entrypoints-revision-1.txtar, which the binary of that commit
// generated: an actual old install, not this binary's templates edited.
func revisionOne(t *testing.T) map[string]string {
	t.Helper()
	source, err := os.ReadFile(filepath.Join("testdata", "entrypoints-revision-1.txtar"))
	if err != nil {
		t.Fatal(err)
	}
	files := map[string]string{}
	parts := regexp.MustCompile(`(?m)^-- (\S+) --\n`).Split(string(source), -1)
	names := regexp.MustCompile(`(?m)^-- (\S+) --\n`).FindAllStringSubmatch(string(source), -1)
	for i, name := range names {
		files[name[1]] = parts[i+1]
	}
	if len(files) != len(grove.Entrypoints()) {
		t.Fatalf("the fixture holds %d files, init manages %d", len(files), len(grove.Entrypoints()))
	}
	return files
}

func runInitCheck(t *testing.T, root string) (int, string, string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code := Run([]string{"--project", root, "init", "--check"}, t.TempDir(), &out, &errOut)
	return code, out.String(), errOut.String()
}

// An install from before entrypoint revisions is diagnosed without a write as
// legacy, which this binary does not serve, refreshed by init with the
// project's own file kept, and then current; a newer, older-than-served or
// missing entrypoint fails the check, and other text at a served revision
// does not.
func TestInitCheckDiagnosesAnOldInstallAndInitRefreshesIt(t *testing.T) {
	t.Parallel()
	root := emptyRepo(t)
	old := revisionOne(t)
	for relative, text := range old {
		write(t, root, relative, text)
	}
	custom := ".agents/skills/grove-shape/SKILL.md"
	write(t, root, custom, "our own shaping skill\n")
	before := hashes(t, root)
	code, out, errOut := runInitCheck(t, root)
	legacy := " (no revision line: written before entrypoint revisions, with a grammar and review brief of its own that this grove cannot vouch for; init rewrites it)\n"
	want := "custom " + custom + " (no init marker: the project's own, not judged; delete it to get the managed version)\n" +
		"legacy .agents/skills/grove-shape/agents/openai.yaml" + legacy +
		"legacy .agents/skills/grove-work/SKILL.md" + legacy +
		"legacy .agents/skills/grove-work/agents/openai.yaml" + legacy +
		"legacy .claude/agents/grove-reviewer.md" + legacy +
		"legacy .claude/skills/grove-shape/SKILL.md" + legacy +
		"legacy .claude/skills/grove-work/SKILL.md" + legacy
	if code != 1 || out != want || !strings.Contains(errOut, "grove: 6 entrypoints are missing or unusable with this grove; nothing was written.") {
		t.Fatalf("check of an old install: %d\n%s%s", code, out, errOut)
	}
	if !reflect.DeepEqual(hashes(t, root), before) {
		t.Fatal("init --check wrote")
	}

	code, out, errOut = runInitAt(t, root)
	if code != 0 || strings.Count(out, "\nupdated ") != 6 || !strings.Contains(out, "\nkept "+custom+" (not managed") {
		t.Fatalf("refresh: %d\n%s%s", code, out, errOut)
	}
	if text, _ := os.ReadFile(filepath.Join(root, custom)); string(text) != "our own shaping skill\n" {
		t.Fatalf("init rewrote the project's own file: %q", text)
	}
	_, out, _ = runInitCheck(t, root)
	if strings.Count(out, "current ") != 6 || strings.Count(out, fmt.Sprintf(" (revision %d)\n", grove.EntrypointRevision)) != 6 {
		t.Fatalf("after refresh:\n%s", out)
	}
	before = hashes(t, root)
	if _, out, _ = runInitAt(t, root); strings.Count(out, "\nunchanged ") != 6 || !reflect.DeepEqual(hashes(t, root), before) {
		t.Fatalf("a second refresh must change nothing:\n%s", out)
	}

	// Other text at a served revision is compatible; content alone never
	// makes a file incompatible.
	work := ".claude/skills/grove-work/SKILL.md"
	current := fmt.Sprintf("grove entrypoint revision %d", grove.EntrypointRevision)
	write(t, root, work, strings.Replace(grove.Entrypoints()[work], "Assignment:", "Work:", 1))
	if code, out, _ := runInitCheck(t, root); code != 0 || !strings.Contains(out, fmt.Sprintf("compatible %s (revision %d, which this grove serves, in other text; init rewrites it)\n", work, grove.EntrypointRevision)) {
		t.Fatalf("%d\n%s", code, out)
	}
	for _, revision := range []string{"99", "0", "two"} {
		write(t, root, work, strings.Replace(grove.Entrypoints()[work], current, "grove entrypoint revision "+revision, 1))
		code, out, errOut := runInitCheck(t, root)
		if code != 1 || !strings.Contains(out, fmt.Sprintf("incompatible %s (revision %s; this grove serves %s; init rewrites it)\n", work, revision, grove.ServedEntrypoints())) || !strings.Contains(errOut, "grove: 1 entrypoints are missing or unusable with this grove; nothing was written.") {
			t.Fatalf("revision %s: %d\n%s%s", revision, code, out, errOut)
		}
	}
	if err := os.Remove(filepath.Join(root, ".claude/agents/grove-reviewer.md")); err != nil {
		t.Fatal(err)
	}
	if code, out, _ := runInitCheck(t, root); code != 1 || !strings.Contains(out, "missing .claude/agents/grove-reviewer.md (init writes it)\n") {
		t.Fatalf("%d\n%s", code, out)
	}
}

// Each generated entrypoint states its revision and asks for its guide with
// it; guide serves exactly the revisions this binary names and refuses the
// rest with the repair, where an entrypoint stops.
func TestGuideServesEntrypointRevisions(t *testing.T) {
	t.Parallel()
	current := fmt.Sprintf("grove entrypoint revision %d", grove.EntrypointRevision)
	for relative, text := range grove.Entrypoints() {
		if verdict, _ := grove.Diagnose(relative, text); !strings.Contains(text, current) || verdict != "current" {
			t.Errorf("%s: %s, or no revision line", relative, verdict)
		}
		if strings.HasSuffix(relative, ".md") && !regexp.MustCompile(fmt.Sprintf("`grove guide (work|shape|review) --entrypoint %d`", grove.EntrypointRevision)).MatchString(text) {
			t.Errorf("%s does not load its guide with its revision", relative)
		}
	}
	for _, revision := range []string{"0", "1", "3", "99", "02", "two"} {
		var out, errOut bytes.Buffer
		if code := Run([]string{"guide", "work", "--entrypoint", revision}, t.TempDir(), &out, &errOut); code != 1 || out.Len() != 0 || !strings.Contains(errOut.String(), "is revision "+revision+", and this grove serves "+grove.ServedEntrypoints()+"; nothing was printed.\nRerun `grove init`") {
			t.Fatalf("%s: %d %q %q", revision, code, out.String(), errOut.String())
		}
	}
	for _, args := range [][]string{{"guide", "work", "--entrypoint"}, {"version", "--entrypoint", "2"}, {"check", "--check"}, {"init", "--check", "--check"}} {
		var out, errOut bytes.Buffer
		if code := Run(args, t.TempDir(), &out, &errOut); code != 2 {
			t.Fatalf("%v: %d %q", args, code, errOut.String())
		}
	}
}
