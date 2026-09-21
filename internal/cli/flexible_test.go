package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The connected schema-3 workflow: a deliberate one-line migration keeps the
// old record valid, pages and neutral IDs work through every command, an
// older-schema branch stays readable beside them, and convert maps an old ID.
func TestSchema3PagesNeutralIDsAndConversion(t *testing.T) {
	t.Parallel()
	root := gitFixture(t)
	gitIn(t, root, "branch", "old") // stays schema 1
	write(t, root, "grove.yaml", "schema_version: 3\nrecords: docs/records\nbrief: docs/records/brief.md\n")
	write(t, root, "docs/records/brief.md", "# Brief\n")
	if code, out, errOut := run(t, root, "check"); code != 0 || out != "OK: 2 records\n" {
		t.Fatalf("schema-1 records must stay valid under schema 3: code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	if code, out, errOut := run(t, root, "context", "W-001"); code != 0 || out == "" {
		t.Fatalf("code=%d stderr=%s", code, errOut)
	}
	if _, err := os.Stat(filepath.Join(root, ".git", "grove")); err == nil {
		t.Fatal("reads need no allocator state")
	}
	page := "docs/records/G-001-release-synthesis.md"
	if code, out, errOut := run(t, root, "new", "page", "Release synthesis"); code != 0 || out != page+"\n" {
		t.Fatalf("code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	if code, out, errOut := run(t, root, "new", "work", "Next thing"); code != 0 || out != "docs/records/G-002-next-thing.md\n" {
		t.Fatalf("code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	if code, out, _ := run(t, root, "list"); code != 0 || !strings.Contains(out, "G-001  page      -  ") || !strings.Contains(out, "G-002  work      proposed") {
		t.Fatalf("list:\n%s", out)
	}
	revision := func(id string) string { return showJSON(t, root, id)["revision"].(string) }
	if code, out, errOut := run(t, root, "update", "G-001", "--expect", revision("G-001"), "--set", "title=Release notes, synthesized", "--set", `relates_to=["W-001"]`); code != 0 || !strings.Contains(out, `"path":"`+page+`"`) {
		t.Fatalf("a title change must keep the path: code=%d stdout=%s stderr=%s", code, out, errOut)
	}
	if code, _, errOut := run(t, root, "update", "G-001", "--expect", revision("G-001"), "--set", "status=done"); code != 1 || !strings.Contains(errOut, "status is not a field that update accepts on page records") {
		t.Fatalf("code=%d stderr=%s", code, errOut)
	}
	if code, _, errOut := run(t, root, "context", "G-001"); code != 1 || !strings.Contains(errOut, "G-001 is a page; only work can be selected") {
		t.Fatalf("code=%d stderr=%s", code, errOut)
	}
	if code, _, errOut := run(t, root, "update", "W-001", "--expect", revision("W-001"), "--set", `relates_to=["G-001"]`); code != 0 {
		t.Fatal(errOut)
	}
	code, out, errOut := run(t, root, "context", "W-001")
	if code != 0 || !strings.Contains(out, "G-001  page  -  listed  related to W-001") || strings.Contains(out, "Source: "+page) {
		t.Fatalf("a related page is listed, never preloaded: code=%d stderr=%s\n%s", code, errOut, out)
	}
	if code, out, errOut = run(t, root, "context", "W-001", "--include", page); code != 0 || !strings.Contains(out, "G-001  page  -  included") || !strings.Contains(out, "Source: "+page) {
		t.Fatalf("code=%d stderr=%s\n%s", code, errOut, out)
	}

	gitIn(t, root, "add", "-A")
	gitIn(t, root, "-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "commit", "-q", "-m", "schema 3")
	code, out, errOut = run(t, root, "versions")
	if code != 0 || strings.Contains(errOut, "invalid") {
		t.Fatalf("an older-schema branch reads beside schema 3: code=%d stderr=%s", code, errOut)
	}
	for _, want := range []string{"G-001  -         committed refs/heads/main", "W-001  proposed  committed refs/heads/old", "W-001  proposed  committed refs/heads/main"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q:\n%s", want, out)
		}
	}
	if code, out, errOut := run(t, root, "workspace", "--source", versionSelector(t, root, "G-001", "committed refs/heads/main")); code != 0 || out != root+"\n" {
		t.Fatalf("a neutral selector must resolve: code=%d stdout=%q stderr=%s", code, out, errOut)
	}

	code, out, errOut = run(t, root, "convert", "W-001")
	if code != 0 || !strings.HasPrefix(out, `{"from":"W-001","from_path":"docs/records/work/`) || !strings.Contains(out, `"id":"G-003","path":"docs/records/G-003-`) {
		t.Fatalf("code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	if source := showJSON(t, root, "G-001")["source"].(string); !strings.Contains(source, `relates_to: ["G-003"]`) {
		t.Fatalf("the page's reference was not rebuilt:\n%s", source)
	}
	if code, out, errOut := run(t, root, "convert", "W-001"); code != 1 || out != "" || !strings.Contains(errOut, "W-001 was already converted to G-003") {
		t.Fatalf("code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	if code, out, errOut := run(t, root, "check"); code != 0 || out != "OK: 4 records\n" {
		t.Fatalf("code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	for _, args := range [][]string{{"convert"}, {"convert", "a", "b"}, {"list", "--type", "page"}, {"new", "page", "T", "--title", "x"}} {
		if code, _, _ := run(t, root, args...); code != 2 {
			t.Fatalf("%v: code=%d, want usage error", args, code)
		}
	}
}

func TestNewPageNeedsSchema3WithoutAllocating(t *testing.T) {
	t.Parallel()
	root := knowledgeFixture(t)
	if code, out, errOut := run(t, root, "new", "page", "Notes"); code != 1 || out != "" || !strings.Contains(errOut, "page records need schema_version 3 in grove.yaml; this project is schema 2") {
		t.Fatalf("code=%d stdout=%q stderr=%s", code, out, errOut)
	}
	if code, _, errOut := run(t, root, "convert", "W-001"); code != 1 || !strings.Contains(errOut, "convert needs schema_version 3") {
		t.Fatalf("code=%d stderr=%s", code, errOut)
	}
	if _, err := os.Stat(filepath.Join(root, ".git", "grove")); err == nil {
		t.Fatal("a refusal must not create allocator state")
	}
}
