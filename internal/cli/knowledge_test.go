package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// knowledgeFixture is the Git fixture with a brief added.
func knowledgeFixture(t *testing.T) string {
	t.Helper()
	root := gitFixture(t)
	write(t, root, "grove.yaml", "schema_version: 3\nrecords: docs/records\nbrief: docs/brief.md\n")
	write(t, root, "docs/brief.md", "# Brief\n")
	return root
}

func run(t *testing.T, root string, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code = Run(args, root, &out, &errOut)
	return code, out.String(), errOut.String()
}

func TestKnowledgeRecordsThroughNewUpdateAndContext(t *testing.T) {
	t.Parallel()
	root := knowledgeFixture(t)
	// The counter is shared, so creation order fixes the IDs: G-001 (work) and
	// G-002 (question) already exist from the fixture.
	for _, tc := range []struct{ kind, title, want string }{
		{"term", "Attempt", "docs/records/G-003-attempt.md\n"},
		{"plan", "Shared plan", "docs/records/G-004-shared-plan.md\n"},
		{"review", "First review", "docs/records/G-005-first-review.md\n"},
	} {
		if code, out, errOut := run(t, root, "new", tc.kind, tc.title); code != 0 || out != tc.want {
			t.Fatalf("new %s: code=%d stdout=%q stderr=%s", tc.kind, code, out, errOut)
		}
	}
	// A second record for one term is refused before anything is reserved or
	// written, so the project stays loadable and the next ID stays free.
	if code, out, errOut := run(t, root, "new", "term", " attempt "); code != 1 || out != "" || !strings.Contains(errOut, "the term attempt is already defined by G-003 in docs/records/G-003-attempt.md") {
		t.Fatalf("code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	if code, out, errOut := run(t, root, "new", "term", "Candidate"); code != 0 || out != "docs/records/G-006-candidate.md\n" {
		t.Fatalf("code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	revision := func(id string) string { return showJSON(t, root, id)["revision"].(string) }
	if code, _, errOut := run(t, root, "update", "G-004", "--expect", revision("G-004"), "--set", `work=["G-001"]`); code != 0 {
		t.Fatal(errOut)
	}
	// An all-letter or all-digit commit must stay a YAML string.
	for _, commit := range []string{"abcdefa", "1234567"} {
		if code, _, errOut := run(t, root, "update", "G-005", "--expect", revision("G-005"), "--set", `work=["G-001"]`, "--set", "examined="+commit); code != 0 {
			t.Fatalf("examined=%s: %s", commit, errOut)
		}
		if source := showJSON(t, root, "G-005")["source"].(string); !strings.Contains(source, `examined: "`+commit+`"`) {
			t.Fatalf("examined must be quoted:\n%s", source)
		}
	}
	for _, c := range []struct{ id, set, want string }{
		{"G-005", "examined=main", "examined: expected a quoted Git commit"},
		{"G-004", `work=["G-404"]`, "work: unresolved target G-404"},
		{"G-004", `work=["G-003"]`, "work: target G-003 must be work"},
		{"G-004", "examined=abcdefa", "examined is not a field that update accepts on plan records"},
		{"G-003", `work=["G-001"]`, "work is not a field that update accepts on term records"},
	} {
		before := revision(c.id)
		if code, _, errOut := run(t, root, "update", c.id, "--expect", before, "--set", c.set); code != 1 || !strings.Contains(errOut, c.want) || revision(c.id) != before {
			t.Fatalf("%s %s: code=%d stderr=%s", c.id, c.set, code, errOut)
		}
	}
	if code, out, errOut := run(t, root, "check"); code != 0 || out != "OK: 6 records\n" {
		t.Fatalf("code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	code, out, errOut := run(t, root, "context", "G-001")
	if code != 0 {
		t.Fatal(errOut)
	}
	for _, want := range []string{"G-004  plan  current  listed  plan for G-001", "G-005  review  current  listed  review of G-001"} {
		if !strings.Contains(out, want) {
			t.Fatalf("context must list %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "Source: docs/records/G-004") || strings.Contains(out, "Source: docs/brief.md") || strings.Contains(out, "G-003") {
		t.Fatalf("context must not include plans or the brief, or list unrelated terms:\n%s", out)
	}
	// One plan shared by several work items is listed for each, never copied.
	if code, _, errOut := run(t, root, "new", "work", "Second"); code != 0 {
		t.Fatal(errOut)
	}
	if code, _, errOut := run(t, root, "update", "G-004", "--expect", revision("G-004"), "--set", `work=["G-001", "G-007"]`); code != 0 {
		t.Fatal(errOut)
	}
	if code, out, errOut := run(t, root, "context", "G-007", "G-001"); code != 0 || !strings.Contains(out, "G-004  plan  current  listed  plan for G-007; plan for G-001") {
		t.Fatalf("code=%d stderr=%s\n%s", code, errOut, out)
	}
	if code, out, errOut = run(t, root, "context", "G-001", "--include", "docs/records/G-004-shared-plan.md"); code != 0 || !strings.Contains(out, "G-004  plan  current  included") {
		t.Fatalf("a caller can include the plan: code=%d stderr=%s\n%s", code, errOut, out)
	}
}

func TestBriefCommand(t *testing.T) {
	t.Parallel()
	root := knowledgeFixture(t)
	if code, out, errOut := run(t, root, "brief"); code != 0 || out != "# Brief\n" || !strings.Contains(errOut, "File: docs/brief.md\n") {
		t.Fatalf("code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	if code, out, _ := run(t, root, "brief", "--json"); code != 0 || !strings.Contains(out, `"path":"docs/brief.md"`) || !strings.Contains(out, `"revision":"sha256:`) {
		t.Fatalf("code=%d stdout=%q", code, out)
	}
	if code, _, _ := run(t, root, "brief", "extra"); code != 2 {
		t.Fatalf("brief takes no arguments: %d", code)
	}
	if err := os.Remove(filepath.Join(root, "docs/brief.md")); err != nil {
		t.Fatal(err)
	}
	for _, command := range []string{"brief", "check", "list"} {
		if code, out, errOut := run(t, root, command); code != 1 || out != "" || !strings.Contains(errOut, "grove.yaml: brief: ") {
			t.Fatalf("%s with a missing brief: code=%d stdout=%q stderr=%s", command, code, out, errOut)
		}
	}
	plain := gitFixture(t)
	if code, _, errOut := run(t, plain, "brief"); code != 1 || !strings.Contains(errOut, "grove.yaml names no brief") {
		t.Fatalf("code=%d stderr=%s", code, errOut)
	}
}

// A committed tree goes through the versions tree reader, which gives the
// loader only grove.yaml and the record folder.
func TestVersionsReadsCommittedKnowledgeRecords(t *testing.T) {
	t.Parallel()
	root := gitFixture(t)
	write(t, root, "grove.yaml", "schema_version: 3\nrecords: docs/records\nbrief: docs/records/brief.md\n")
	write(t, root, "docs/records/brief.md", "# Brief\n")
	for _, kind := range []string{"term", "plan", "review"} {
		if code, _, errOut := run(t, root, "new", kind, "A "+kind); code != 0 {
			t.Fatal(errOut)
		}
	}
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "commit", "-q", "-m", "knowledge records")
	code, out, errOut := run(t, root, "versions")
	if code != 0 || strings.Contains(errOut, "invalid") {
		t.Fatalf("code=%d stderr=%s", code, errOut)
	}
	for _, want := range []string{"G-003  proposed  committed refs/heads/main", "G-004  current   committed refs/heads/main", "G-005  current   committed refs/heads/main"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "brief") {
		t.Fatalf("the brief is not a record:\n%s", out)
	}
	selector := versionSelector(t, root, "G-003", "committed refs/heads/main")
	if code, out, errOut := run(t, root, "workspace", "--source", selector); code != 0 || out != root+"\n" {
		t.Fatalf("a term selector must resolve: code=%d stdout=%q stderr=%s", code, out, errOut)
	}
}
