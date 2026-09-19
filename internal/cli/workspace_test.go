package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const selector = "live:.:refs/heads/main@0123456789ab:W-001@0123456789ab:0123456789abcdef"

func TestWorkspaceUsage(t *testing.T) {
	for _, args := range [][]string{
		{"workspace"}, {"workspace", "W-001", "--source", selector}, {"workspace", "--source"}, {"workspace", "--source", ""},
		{"workspace", "--source", "W-001"}, {"workspace", "--source", selector, "--source", selector},
		{"workspace", "--source", "committed:main@0123456789ab:W-001@0123456789ab:0123456789abcdef"},
		{"versions", "--source", selector}, {"list", "--source", selector},
	} {
		var out, errOut bytes.Buffer
		if code := Run(args, t.TempDir(), &out, &errOut); code != 2 || out.Len() != 0 || !strings.Contains(errOut.String(), "Usage:") {
			t.Fatalf("%v: code=%d stderr=%s", args, code, errOut.String())
		}
	}
	var out, errOut bytes.Buffer
	if code := Run([]string{"workspace", "--help"}, t.TempDir(), &out, &errOut); code != 0 || !strings.Contains(out.String(), "--source") {
		t.Fatalf("help must work without a project: %d %s", code, out.String())
	}
}

// versionSelector returns the selector of id in the source whose table cell
// is where, from the versions command run in dir.
func versionSelector(t *testing.T, dir, id, where string) string {
	t.Helper()
	var out, errOut bytes.Buffer
	if code := Run([]string{"versions", id}, dir, &out, &errOut); code != 0 {
		t.Fatal(errOut.String())
	}
	for _, row := range rowsOf(out.String()) {
		if len(row) == 5 && row[2] == where {
			return row[4]
		}
	}
	t.Fatalf("no %s row for %s in\n%s", where, id, out.String())
	return ""
}

func TestWorkspaceCLI(t *testing.T) {
	root, wt := featureFixture(t)
	odd := filepath.Join(filepath.Dir(root), "odd\nwt")
	gitIn(t, root, "worktree", "add", "-q", "--detach", odd, "feature")
	before := map[string]map[string][32]byte{root: hashes(t, root), wt: hashes(t, wt), odd: hashes(t, odd)}

	live := versionSelector(t, root, "W-001", "live feature-wt refs/heads/feature")
	var out, errOut bytes.Buffer
	if code := Run([]string{"workspace", "--source", live}, root, &out, &errOut); code != 0 || out.String() != wt+"\n" {
		t.Fatalf("code=%d stdout=%q stderr=%s", code, out.String(), errOut.String())
	}
	for _, want := range []string{"Project: " + root, "Checkout: " + wt, "Branch: refs/heads/feature at " + gitIn(t, wt, "rev-parse", "HEAD"), "Record: " + filepath.Join(wt, "docs/records/work/renamed.md"), "Revision: sha256:"} {
		if !strings.Contains(errOut.String(), want) {
			t.Fatalf("missing %q in stderr:\n%s", want, errOut.String())
		}
	}
	out.Reset()
	errOut.Reset()
	detached := versionSelector(t, root, "W-001", "live odd-wt detached")
	if code := Run([]string{"workspace", "--source", detached, "--json"}, root, &out, &errOut); code != 0 {
		t.Fatal(errOut.String())
	}
	var got map[string]any
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	want := map[string]any{
		"checkout": odd, "project": odd, "record": filepath.Join(odd, "docs/records/work/renamed.md"), "ref": nil,
		"head": gitIn(t, odd, "rev-parse", "HEAD"), "revision": showJSON(t, odd, "W-001")["revision"], "selector": detached,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("json:\n%v\nwant\n%v", got, want)
	}
	if !strings.Contains(errOut.String(), "Detached: HEAD at ") || !strings.Contains(errOut.String(), `odd\nwt`) || strings.Contains(errOut.String(), "\nwt") {
		t.Fatalf("stderr escapes control characters: %s", errOut.String())
	}
	for dir, hashed := range before {
		if !reflect.DeepEqual(hashed, hashes(t, dir)) {
			t.Fatalf("workspace must leave %s unchanged", dir)
		}
	}
	// Refusals exit 1 with a grove: diagnostic and nothing on stdout.
	write(t, wt, "docs/records/work/renamed.md", strings.Replace(work, "status: proposed", "status: done", 1))
	for _, c := range []struct{ selector, want string }{
		{live, "W-001 changed since it was selected"},
		{versionSelector(t, root, "W-001", "committed refs/heads/feature"), "select the live observation instead"},
		{strings.Replace(live, "feature-wt", "elsewhere", 1), "worktree elsewhere is no longer registered"},
	} {
		out.Reset()
		errOut.Reset()
		if code := Run([]string{"workspace", "--source", c.selector}, root, &out, &errOut); code != 1 || out.Len() != 0 || !strings.Contains(errOut.String(), "\ngrove: ") || !strings.Contains(errOut.String(), c.want) {
			t.Fatalf("%s: code=%d stdout=%q stderr=%s", c.selector, code, out.String(), errOut.String())
		}
	}
	// An invalid current checkout still resolves a selection elsewhere.
	write(t, root, "grove.yaml", "schema_version: 2\nrecords: docs/records\n")
	out.Reset()
	errOut.Reset()
	if code := Run([]string{"workspace", "--source", detached}, root, &out, &errOut); code != 0 || out.String() != odd+"\n" {
		t.Fatalf("code=%d stdout=%q stderr=%s", code, out.String(), errOut.String())
	}
	plain := projectFixture(t)
	out.Reset()
	errOut.Reset()
	if code := Run([]string{"workspace", "--source", selector}, plain, &out, &errOut); code != 1 || !strings.Contains(errOut.String(), "requires a Git repository") {
		t.Fatalf("plain directories: %d %s", code, errOut.String())
	}
}

// The joint W-004/W-005 fixture: list versions, select one explicitly,
// resolve its workspace, and read exactly that version's bytes there.
func TestJointWorkflow(t *testing.T) {
	root, wt := featureFixture(t)
	var out, errOut bytes.Buffer
	if code := Run([]string{"versions", "W-001", "--json"}, root, &out, &errOut); code != 0 {
		t.Fatal(errOut.String())
	}
	var listed struct {
		Records []struct {
			Versions []map[string]any
		}
	}
	if err := json.Unmarshal(out.Bytes(), &listed); err != nil {
		t.Fatal(err)
	}
	var chosen map[string]any
	for _, v := range listed.Records[0].Versions {
		if v["kind"] == "live" && v["ref"] == "refs/heads/feature" {
			chosen = v
		}
	}
	if chosen == nil || chosen["status"] != "active" {
		t.Fatalf("no live feature version among %v", listed.Records[0].Versions)
	}
	out.Reset()
	errOut.Reset()
	if code := Run([]string{"workspace", "--source", chosen["selector"].(string), "--json"}, root, &out, &errOut); code != 0 {
		t.Fatal(errOut.String())
	}
	var located map[string]any
	if err := json.Unmarshal(out.Bytes(), &located); err != nil {
		t.Fatal(err)
	}
	if located["project"] != wt || located["revision"] != chosen["revision"] {
		t.Fatalf("workspace: %v", located)
	}
	out.Reset()
	errOut.Reset()
	if code := Run([]string{"--project", located["project"].(string), "show", "W-001"}, t.TempDir(), &out, &errOut); code != 0 {
		t.Fatal(errOut.String())
	}
	if out.String() != chosen["source"] || !strings.Contains(out.String(), "status: active") {
		t.Fatalf("show in the located workspace must return the selected bytes:\n%s", out.String())
	}
	if data, err := os.ReadFile(filepath.Join(root, "docs/records/work/renamed.md")); err != nil || string(data) != work || gitIn(t, root, "rev-parse", "--abbrev-ref", "HEAD") != "main" {
		t.Fatal("main's copy and branch are untouched")
	}
}
