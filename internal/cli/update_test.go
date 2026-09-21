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

const rev = "sha256:" + "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func TestUpdateUsage(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{
		{"update"}, {"update", "W-001"}, {"update", "W-001", "extra", "--expect", rev, "--set", "status=done"},
		{"update", "W-001", "--set", "status=done"}, {"update", "W-001", "--expect", rev},
		{"update", "W-001", "--expect", "abc", "--set", "status=done"},
		{"update", "W-001", "--expect", strings.ToUpper(rev), "--set", "status=done"},
		{"update", "W-001", "--expect", rev, "--expect", rev, "--set", "status=done"},
		{"update", "W-001", "--expect", rev, "--set", "status=done", "--set", "status=active"},
		{"update", "W-001", "--expect", rev, "--set", "size=small", "--unset", "size"},
		{"update", "W-001", "--expect", rev, "--unset", "size", "--unset", "size"},
		{"update", "W-001", "--expect", rev, "--set", "status"}, {"update", "W-001", "--expect", rev, "--set", "=x"},
		{"update", "W-001", "--expect", rev, "--set"}, {"update", "W-001", "--expect", rev, "--unset", ""},
		{"list", "--expect", rev}, {"show", "W-001", "--set", "a=b"}, {"check", "--unset", "size"},
	} {
		var out, errOut bytes.Buffer
		if code := Run(args, t.TempDir(), &out, &errOut); code != 2 || out.Len() != 0 || !strings.Contains(errOut.String(), "Usage:") {
			t.Fatalf("%v: code=%d stderr=%s", args, code, errOut.String())
		}
	}
	var out, errOut bytes.Buffer
	if code := Run([]string{"update", "--help"}, t.TempDir(), &out, &errOut); code != 0 || !strings.Contains(out.String(), "--expect") {
		t.Fatalf("help must work without a project: %d %s", code, out.String())
	}
}

func showJSON(t *testing.T, root, id string) map[string]any {
	t.Helper()
	var out, errOut bytes.Buffer
	if code := Run([]string{"show", id, "--json"}, root, &out, &errOut); code != 0 {
		t.Fatal(errOut.String())
	}
	var got map[string]any
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	return got
}

func TestUpdateWorkflowCreateUpdateCloseReopenCheck(t *testing.T) {
	t.Parallel()
	root := gitFixture(t)
	var out, errOut bytes.Buffer
	if code := Run([]string{"new", "work", "Workflow record", "--slug", "workflow"}, root, &out, &errOut); code != 0 {
		t.Fatal(errOut.String())
	}
	path := strings.TrimSpace(out.String())
	shown := showJSON(t, root, "W-002")
	created := shown["source"].(string)
	steps := []struct {
		args    []string
		changed bool
		status  string
	}{
		{[]string{"--set", "status=active", "--set", "kind=feature", "--set=priority=2", "--set", `depends_on=["W-001"]`, "--set", "title=Renamed: 版本 \"quoted\""}, true, "active"},
		{[]string{"--set", "status=active", "--unset", "size"}, false, "active"},
		{[]string{"--set", "status=done", "--unset", "priority"}, true, "done"},
		{[]string{"--set", "status=active"}, true, "active"},
	}
	for _, step := range steps {
		expect := showJSON(t, root, "W-002")["revision"].(string)
		out.Reset()
		errOut.Reset()
		args := append([]string{"update", "W-002", "--expect", expect}, step.args...)
		if code := Run(args, filepath.Join(root, "docs"), &out, &errOut); code != 0 {
			t.Fatalf("%v: %s", args, errOut.String())
		}
		var result map[string]any
		if err := json.Unmarshal(out.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		after := showJSON(t, root, "W-002")
		want := map[string]any{"id": "W-002", "path": path, "revision": after["revision"], "changed": step.changed}
		if !reflect.DeepEqual(result, want) {
			t.Fatalf("got %v want %v", result, want)
		}
		if !step.changed && after["source"] != created && after["revision"] != expect {
			t.Fatal("no-op changed the file")
		}
		if step.changed && after["revision"] == expect {
			t.Fatal("change did not alter the revision")
		}
		if !strings.Contains(after["source"].(string), "status: "+step.status+"\n") {
			t.Fatalf("status not %s:\n%s", step.status, after["source"])
		}
		created = after["source"].(string)
	}
	final := showJSON(t, root, "W-002")["source"].(string)
	createdLine := strings.Split(strings.SplitN(showJSON(t, root, "W-002")["source"].(string), "created: ", 2)[1], "\n")[0]
	if !strings.Contains(created, "created: "+createdLine) || !strings.HasSuffix(final, "## Outcome\n\n## Constraints\n\n## Acceptance\n\n## Next\n") {
		t.Fatalf("created or body changed:\n%s", final)
	}
	if strings.Contains(final, "priority") || !strings.Contains(final, "title: \"Renamed: 版本 \\\"quoted\\\"\"\n") || !strings.Contains(final, "depends_on: [\"W-001\"]\n") {
		t.Fatalf("final source:\n%s", final)
	}
	out.Reset()
	if code := Run([]string{"check"}, root, &out, &errOut); code != 0 || !strings.Contains(out.String(), "3 records") {
		t.Fatalf("check: %s %s", out.String(), errOut.String())
	}
	out.Reset()
	if code := Run([]string{"list"}, root, &out, &errOut); code != 0 || !strings.Contains(out.String(), "W-002  work      active    Renamed: 版本 \\\"quoted\\\"") {
		t.Fatalf("list: %s", out.String())
	}
	if entries, _ := os.ReadDir(filepath.Join(root, "docs/records/work")); len(entries) != 2 {
		t.Fatalf("no files may be added or renamed: %v", entries)
	}
}

func TestUpdateOperationErrorsAndOutputFailure(t *testing.T) {
	t.Parallel()
	root := gitFixture(t)
	expect := showJSON(t, root, "W-001")["revision"].(string)
	before := hashes(t, filepath.Join(root, "docs"))
	for _, args := range [][]string{
		{"update", "W-001", "--expect", rev, "--set", "status=active"},
		{"update", "W-404", "--expect", expect, "--set", "status=active"},
		{"update", "W-001", "--expect", expect, "--set", "status=bogus"},
		{"update", "W-001", "--expect", expect, "--set", "id=W-009"},
		{"update", "W-001", "--expect", expect, "--set", "blocks=[]"},
		{"update", "Q-001", "--expect", showJSON(t, root, "Q-001")["revision"].(string), "--set", `blocks=["Q-001"]`},
	} {
		var out, errOut bytes.Buffer
		if code := Run(args, root, &out, &errOut); code != 1 || out.Len() != 0 || !strings.HasPrefix(errOut.String(), "Project: ") || !strings.Contains(errOut.String(), "grove: ") {
			t.Fatalf("%v: code=%d stdout=%q stderr=%q", args, code, out.String(), errOut.String())
		}
	}
	if !reflect.DeepEqual(before, hashes(t, filepath.Join(root, "docs"))) {
		t.Fatal("failed updates changed project files")
	}
	var errOut bytes.Buffer
	if code := Run([]string{"update", "W-001", "--expect", expect, "--set", "status=active"}, root, brokenWriter{}, &errOut); code != 1 || !strings.Contains(errOut.String(), "the update was applied to docs/records/work/renamed.md; revision sha256:") {
		t.Fatalf("output failure must report the applied update: code=%d stderr=%s", code, errOut.String())
	}
	after := showJSON(t, root, "W-001")
	if !strings.Contains(errOut.String(), after["revision"].(string)) || !strings.Contains(after["source"].(string), "status: active") {
		t.Fatal("reported revision must describe the published file")
	}
	errOut.Reset()
	if code := Run([]string{"update", "W-001", "--expect", after["revision"].(string), "--set", "status=active"}, root, brokenWriter{}, &errOut); code != 1 || !strings.Contains(errOut.String(), "no change was needed") {
		t.Fatalf("output failure after a no-op: %s", errOut.String())
	}
	root = projectFixture(t)
	errOut.Reset()
	var out bytes.Buffer
	if code := Run([]string{"update", "W-001", "--expect", expect, "--set", "status=active"}, root, &out, &errOut); code != 1 || !strings.Contains(errOut.String(), "Git") {
		t.Fatalf("outside Git: code=%d stderr=%s", code, errOut.String())
	}
}
