package cli

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
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
		if out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput(); err != nil {
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
		"created .claude/skills/grove-shape/SKILL.md\ncreated .claude/skills/grove-work/SKILL.md\n"
	if out != want || !strings.Contains(errOut, "Project: "+root) {
		t.Fatalf("stdout=%q stderr=%q", out, errOut)
	}
	adapter, err := os.ReadFile(filepath.Join(root, ".claude/skills/grove-work/SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, needle := range []string{"disable-model-invocation: true", managedMarker, "Assignment: $ARGUMENTS", "grove guide work", "stop and say so"} {
		if !strings.Contains(string(adapter), needle) {
			t.Fatalf("the Claude adapter lacks %q:\n%s", needle, adapter)
		}
	}
	portable := map[string]string{".claude/skills/grove-work/SKILL.md": string(adapter)}
	for _, relative := range []string{".claude/skills/grove-shape/SKILL.md", ".agents/skills/grove-work/SKILL.md", ".agents/skills/grove-shape/SKILL.md", ".agents/skills/grove-work/agents/openai.yaml"} {
		source, err := os.ReadFile(filepath.Join(root, relative))
		if err != nil {
			t.Fatal(err)
		}
		portable[relative] = string(source)
	}
	for _, name := range []string{"work", "shape"} {
		var guide, guideErr bytes.Buffer
		if code := Run([]string{"guide", name}, t.TempDir(), &guide, &guideErr); code != 0 {
			t.Fatal(guideErr.String())
		}
		portable["guide "+name] = guide.String()
	}
	for name, text := range portable {
		for _, forbidden := range []string{"go run", "docs/work-execution.md", "docs/work-shaping.md", "../grove/", "../.claude/", "../.agents/", ".claude/worktrees", "worktree-G-", "AGENTS.md` here"} {
			if strings.Contains(text, forbidden) {
				t.Fatalf("%s must not depend on Grove's own repository, but mentions %q", name, forbidden)
			}
		}
	}
	if !strings.Contains(errOut, "Next: grove check") {
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
	edited := strings.Replace(string(adapter), "stop and say so", "stop", 1)
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
		{[]string{"version"}, "grove "}, // ends with the guide digest, checked below
	} {
		var out, errOut bytes.Buffer
		if code := Run(c.args, t.TempDir(), &out, &errOut); code != 0 || !strings.HasPrefix(out.String(), c.prefix) {
			t.Fatalf("%v: code=%d stdout=%q stderr=%q", c.args, code, out.String(), errOut.String())
		}
	}
	var version, versionErr bytes.Buffer
	if code := Run([]string{"version"}, t.TempDir(), &version, &versionErr); code != 0 || !regexp.MustCompile(`^grove \S.* guides sha256:[0-9a-f]{12}\n$`).MatchString(version.String()) {
		t.Fatalf("version=%q", version.String())
	}
	for _, args := range [][]string{{"guide"}, {"guide", "both"}, {"guide", "work", "shape"}, {"version", "x"}, {"init", "here"}} {
		var out, errOut bytes.Buffer
		if code := Run(args, t.TempDir(), &out, &errOut); code != 2 || out.Len() != 0 {
			t.Fatalf("%v: code=%d stdout=%q", args, code, out.String())
		}
	}
}
