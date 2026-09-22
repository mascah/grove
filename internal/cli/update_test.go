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
		{"update"}, {"update", "G-001"}, {"update", "G-001", "extra", "--expect", rev, "--set", "status=done"},
		{"update", "G-001", "--expect", rev}, {"update", "G-001", "--commit"},
		{"update", "G-001", "--commit", "--commit", "--set", "status=done"},
		{"update", "G-001", "--expect", "abc", "--set", "status=done"},
		{"update", "G-001", "--expect", strings.ToUpper(rev), "--set", "status=done"},
		{"update", "G-001", "--expect", rev, "--expect", rev, "--set", "status=done"},
		{"update", "G-001", "--expect", rev, "--set", "status=done", "--set", "status=active"},
		{"update", "G-001", "--expect", rev, "--set", "size=small", "--unset", "size"},
		{"update", "G-001", "--expect", rev, "--unset", "size", "--unset", "size"},
		{"update", "G-001", "--expect", rev, "--set", "status"}, {"update", "G-001", "--expect", rev, "--set", "=x"},
		{"update", "G-001", "--expect", rev, "--set"}, {"update", "G-001", "--expect", rev, "--unset", ""},
		{"list", "--expect", rev}, {"show", "G-001", "--set", "a=b"}, {"check", "--unset", "size"}, {"list", "--commit"},
	} {
		var out, errOut bytes.Buffer
		if code := Run(args, t.TempDir(), &out, &errOut); code != 2 || out.Len() != 0 || !strings.Contains(errOut.String(), "Usage:") {
			t.Fatalf("%v: code=%d stderr=%s", args, code, errOut.String())
		}
	}
	var out, errOut bytes.Buffer
	if code := Run([]string{"update", "--help"}, t.TempDir(), &out, &errOut); code != 0 || !strings.Contains(out.String(), "[--expect REVISION]") || !strings.Contains(out.String(), "[--commit]") {
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
	shown := showJSON(t, root, "G-003")
	created := shown["source"].(string)
	steps := []struct {
		args    []string
		changed bool
		status  string
	}{
		{[]string{"--set", "status=active", "--set", "kind=feature", "--set=priority=2", "--set", `depends_on=["G-001"]`, "--set", "title=Renamed: 版本 \"quoted\""}, true, "active"},
		{[]string{"--set", "status=active", "--unset", "size"}, false, "active"},
		{[]string{"--set", "status=review", "--set", "candidate=" + gitIn(t, root, "rev-parse", "HEAD"), "--unset", "priority"}, true, "review"},
		{[]string{"--set", "status=done"}, true, "done"}, // HEAD contains the candidate, so this checkout may close it
		{[]string{"--set", "status=active"}, true, "active"},
	}
	for _, step := range steps {
		expect := showJSON(t, root, "G-003")["revision"].(string)
		out.Reset()
		errOut.Reset()
		args := append([]string{"update", "G-003", "--expect", expect}, step.args...)
		if code := Run(args, filepath.Join(root, "docs"), &out, &errOut); code != 0 {
			t.Fatalf("%v: %s", args, errOut.String())
		}
		var result map[string]any
		if err := json.Unmarshal(out.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		after := showJSON(t, root, "G-003")
		want := map[string]any{"id": "G-003", "path": path, "revision": after["revision"], "changed": step.changed}
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
	final := showJSON(t, root, "G-003")["source"].(string)
	createdLine := strings.Split(strings.SplitN(showJSON(t, root, "G-003")["source"].(string), "created: ", 2)[1], "\n")[0]
	if !strings.Contains(created, "created: "+createdLine) || !strings.HasSuffix(final, "## Outcome\n\n## Constraints\n\n## Acceptance\n\n## Next\n") {
		t.Fatalf("created or body changed:\n%s", final)
	}
	if strings.Contains(final, "priority") || !strings.Contains(final, "title: \"Renamed: 版本 \\\"quoted\\\"\"\n") || !strings.Contains(final, "depends_on: [\"G-001\"]\n") {
		t.Fatalf("final source:\n%s", final)
	}
	out.Reset()
	if code := Run([]string{"check"}, root, &out, &errOut); code != 0 || !strings.Contains(out.String(), "3 records") {
		t.Fatalf("check: %s %s", out.String(), errOut.String())
	}
	out.Reset()
	if code := Run([]string{"list"}, root, &out, &errOut); code != 0 || !strings.Contains(out.String(), "G-003  work      active    Renamed: 版本 \\\"quoted\\\"") {
		t.Fatalf("list: %s", out.String())
	}
	// The new record lands flat under the record root, not in this type
	// folder, which update must leave holding only the original fixture file.
	if entries, _ := os.ReadDir(filepath.Join(root, "docs/records/work")); len(entries) != 1 {
		t.Fatalf("no files may be added or renamed: %v", entries)
	}
}

func TestUpdateOperationErrorsAndOutputFailure(t *testing.T) {
	t.Parallel()
	root := gitFixture(t)
	expect := showJSON(t, root, "G-001")["revision"].(string)
	before := hashes(t, filepath.Join(root, "docs"))
	for _, args := range [][]string{
		{"update", "G-001", "--expect", rev, "--set", "status=active"},
		{"update", "G-404", "--expect", expect, "--set", "status=active"},
		{"update", "G-001", "--expect", expect, "--set", "status=bogus"},
		{"update", "G-001", "--expect", expect, "--set", "id=G-009"},
		{"update", "G-001", "--expect", expect, "--set", "blocks=[]"},
		{"update", "G-002", "--expect", showJSON(t, root, "G-002")["revision"].(string), "--set", `blocks=["G-002"]`},
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
	if code := Run([]string{"update", "G-001", "--expect", expect, "--set", "status=active"}, root, brokenWriter{}, &errOut); code != 1 || !strings.Contains(errOut.String(), "the update was applied to docs/records/work/renamed.md; revision sha256:") {
		t.Fatalf("output failure must report the applied update: code=%d stderr=%s", code, errOut.String())
	}
	after := showJSON(t, root, "G-001")
	if !strings.Contains(errOut.String(), after["revision"].(string)) || !strings.Contains(after["source"].(string), "status: active") {
		t.Fatal("reported revision must describe the published file")
	}
	errOut.Reset()
	if code := Run([]string{"update", "G-001", "--expect", after["revision"].(string), "--set", "status=active"}, root, brokenWriter{}, &errOut); code != 1 || !strings.Contains(errOut.String(), "no change was needed") {
		t.Fatalf("output failure after a no-op: %s", errOut.String())
	}
	root = projectFixture(t)
	errOut.Reset()
	var out bytes.Buffer
	if code := Run([]string{"update", "G-001", "--expect", expect, "--set", "status=active"}, root, &out, &errOut); code != 1 || !strings.Contains(errOut.String(), "Git") {
		t.Fatalf("outside Git: code=%d stderr=%s", code, errOut.String())
	}
}

// TestUpdateCommitResultAndFailure covers the CLI side of G-079: no --expect,
// commit in the result (null for a no-op), and a failed commit reported as an
// applied update with exit 1.
func TestUpdateCommitResultAndFailure(t *testing.T) {
	t.Parallel()
	root := gitFixture(t)
	for _, kv := range [][2]string{{"user.name", "t"}, {"user.email", "t@t"}, {"commit.gpgsign", "false"}, {"maintenance.auto", "false"}} {
		gitIn(t, root, "config", kv[0], kv[1])
	}
	run := func(args ...string) (int, map[string]any, string) {
		t.Helper()
		var out, errOut bytes.Buffer
		code := Run(args, root, &out, &errOut)
		var result map[string]any
		if out.Len() != 0 {
			if err := json.Unmarshal(out.Bytes(), &result); err != nil {
				t.Fatal(err)
			}
		}
		return code, result, errOut.String()
	}
	code, result, stderr := run("update", "G-001", "--set", "status=active", "--commit")
	if code != 0 || result["changed"] != true || result["commit"] != gitIn(t, root, "rev-parse", "HEAD") {
		t.Fatalf("code=%d result=%v stderr=%s", code, result, stderr)
	}
	if files := gitIn(t, root, "show", "--format=", "--name-only", "HEAD"); files != "docs/records/work/renamed.md" {
		t.Fatalf("commit must hold the record alone: %q", files)
	}
	code, result, _ = run("update", "G-001", "--set", "status=active", "--commit")
	if commit, present := result["commit"]; code != 0 || result["changed"] != false || !present || commit != nil {
		t.Fatalf("no-op: code=%d result=%v", code, result)
	}
	code, result, _ = run("update", "G-001", "--set", "status=proposed")
	if _, present := result["commit"]; code != 0 || present || gitIn(t, root, "status", "--porcelain") != "M docs/records/work/renamed.md" { // gitIn trims the leading space
		t.Fatalf("without --commit nothing is committed and no key is printed: %v", result)
	}
	hooks := filepath.Join(root, "hooks")
	write(t, root, "hooks/pre-commit", "#!/bin/sh\nexit 1\n")
	if err := os.Chmod(filepath.Join(hooks, "pre-commit"), 0o755); err != nil {
		t.Fatal(err)
	}
	gitIn(t, root, "config", "core.hooksPath", hooks)
	code, result, stderr = run("update", "G-001", "--set", "status=active", "--commit")
	if code != 1 || result != nil || !strings.Contains(stderr, "nothing was committed (the update was applied to docs/records/work/renamed.md; revision "+showJSON(t, root, "G-001")["revision"].(string)+")") {
		t.Fatalf("failed commit: code=%d stderr=%s", code, stderr)
	}
}
