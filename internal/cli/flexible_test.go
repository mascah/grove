package cli

import (
	"strings"
	"testing"
)

// The fixture holds G-001 (work) and G-002 (question), so creation order fixes
// the IDs below.
func TestPagesAndConvertThroughTheCLI(t *testing.T) {
	t.Parallel()
	root := knowledgeFixture(t)
	if code, out, errOut := run(t, root, "new", "page", "Probe synthesis"); code != 0 || out != "docs/records/G-003-probe-synthesis.md\n" {
		t.Fatalf("new page: code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	// A page has no status, and is never selectable work.
	if code, out, errOut := run(t, root, "list"); code != 0 || !strings.Contains(out, "G-003  page      -") {
		t.Fatalf("list: code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	if code, out, errOut := run(t, root, "context", "G-003"); code != 1 || out != "" || !strings.Contains(errOut, "G-003 is a page; only work can be selected") {
		t.Fatalf("context: code=%d stdout=%q stderr=%s", code, out, errOut)
	}

	write(t, root, "docs/legacy-note.md", "# Legacy note\n\nBody.\n")
	want := `{"from":"docs/legacy-note.md","from_path":"docs/legacy-note.md","id":"G-004","path":"docs/records/G-004-note.md"}` + "\n"
	if code, out, errOut := run(t, root, "convert", "docs/legacy-note.md", "--type", "plan", "--title", "Legacy note", "--slug", "note"); code != 0 || out != want {
		t.Fatalf("convert: code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	if source := showJSON(t, root, "G-004")["source"].(string); !strings.Contains(source, "type: plan\ntitle: \"Legacy note\"\nstatus: current\nformerly: \"docs/legacy-note.md\"\n") || !strings.HasSuffix(source, "# Legacy note\n\nBody.\n") {
		t.Fatalf("converted record:\n%s", source)
	}
	// A rerun is refused and reserves nothing: the next record is G-005.
	if code, out, errOut := run(t, root, "convert", "docs/legacy-note.md", "--type", "plan", "--title", "Legacy note"); code != 1 || out != "" || !strings.Contains(errOut, "already converted to G-004") {
		t.Fatalf("rerun: code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	if code, out, errOut := run(t, root, "new", "page", "Next"); code != 0 || out != "docs/records/G-005-next.md\n" {
		t.Fatalf("new after refusal: code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	if code, _, errOut := run(t, root, "check"); code != 0 {
		t.Fatal(errOut)
	}
}

func TestConvertUsageErrors(t *testing.T) {
	t.Parallel()
	root := knowledgeFixture(t)
	for _, args := range [][]string{
		{"convert"},
		{"convert", "a.md", "b.md", "--type", "plan", "--title", "T"},
	} {
		if code, out, errOut := run(t, root, args...); code != 2 || out != "" || !strings.Contains(errOut, "convert requires exactly one document path") {
			t.Fatalf("%v: code=%d stdout=%q stderr=%s", args, code, out, errOut)
		}
	}
}
