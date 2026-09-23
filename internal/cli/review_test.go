package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestApproveAndFeedbackUsage(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{
		{"approve"}, {"approve", "G-001"}, {"approve", "G-001", "fine", "extra"},
		{"feedback"}, {"feedback", "G-001"}, {"feedback", "G-001", "more", "--set", "status=active"},
		{"approve", "G-001", "fine", "--commit"}, {"feedback", "G-001", "more", "--expect", rev},
		{"approve", "G-001", "fine", "--json"}, {"approve", "G-001", "fine", "--cleanup"},
	} {
		var out, errOut bytes.Buffer
		if code := Run(args, t.TempDir(), &out, &errOut); code != 2 || out.Len() != 0 || !strings.Contains(errOut.String(), "Usage:") {
			t.Fatalf("%v: code=%d stderr=%s", args, code, errOut.String())
		}
	}
}

// TestApproveAndFeedbackCommands drives G-044's two dispositions through the
// CLI on a work branch: approve prints what update prints and commits,
// feedback reopens the work and says where to continue.
func TestApproveAndFeedbackCommands(t *testing.T) {
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
				t.Fatalf("%v: %v\n%s", args, err, out.String())
			}
		}
		return code, result, errOut.String()
	}
	gitIn(t, root, "checkout", "-q", "-b", "feature")
	run("update", "G-001", "--set", "status=active", "--commit")
	candidate := gitIn(t, root, "rev-parse", "HEAD")
	if code, _, stderr := run("update", "G-001", "--set", "status=review", "--set", "candidate="+candidate, "--commit"); code != 0 {
		t.Fatal(stderr)
	}
	code, result, stderr := run("approve", "G-001", "Good enough to ship.")
	if code != 0 || result["changed"] != true || result["commit"] != gitIn(t, root, "rev-parse", "HEAD") || result["id"] != "G-001" {
		t.Fatalf("approve: code=%d result=%v stderr=%s", code, result, stderr)
	}
	if src := showJSON(t, root, "G-001")["source"].(string); !strings.Contains(src, "approved: \""+candidate+"\"\n") || !strings.HasSuffix(src, "Verdict on candidate "+candidate[:7]+", "+today()+": Good enough to ship.\n") {
		t.Fatalf("approved record:\n%s", src)
	}
	code, result, stderr = run("feedback", "G-001", "Add the empty case.")
	if code != 0 || result["changed"] != true || result["commit"] != gitIn(t, root, "rev-parse", "HEAD") {
		t.Fatalf("feedback: code=%d result=%v stderr=%s", code, result, stderr)
	}
	if !strings.Contains(stderr, "Next: G-001 is active on branch feature in "+root+"; continue there with /grove-work G-001\n") {
		t.Fatalf("feedback must say where to continue:\n%s", stderr)
	}
	if src := showJSON(t, root, "G-001")["source"].(string); strings.Contains(src, "approved:") || !strings.Contains(src, "status: active\n") || !strings.Contains(src, "candidate: \""+candidate+"\"\n") || !strings.HasSuffix(src, "Feedback on candidate "+candidate[:7]+", "+today()+": Add the empty case.\n") {
		t.Fatalf("record after feedback:\n%s", src)
	}
	code, result, stderr = run("approve", "G-001", "again")
	if code != 1 || result != nil || !strings.Contains(stderr, "grove: G-001 is active, not in review") {
		t.Fatalf("approve on active work: code=%d stderr=%s", code, stderr)
	}
}

func today() string { return time.Now().UTC().Format("2006-01-02") }

func TestIntegrateUsage(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{
		{"integrate"}, {"integrate", "G-001", "extra"}, {"integrate", "G-001", "--cleanup", "--cleanup"},
		{"integrate", "G-001", "--json"}, {"integrate", "G-001", "--commit"}, {"list", "--cleanup"}, {"update", "G-001", "--set", "status=done", "--cleanup"},
	} {
		var out, errOut bytes.Buffer
		if code := Run(args, t.TempDir(), &out, &errOut); code != 2 || out.Len() != 0 || !strings.Contains(errOut.String(), "Usage:") {
			t.Fatalf("%v: code=%d stderr=%s", args, code, errOut.String())
		}
	}
}

// TestIntegrateCommand runs the whole loop through the CLI on one checkout:
// review and approval on a branch, then integration from main, with one
// fact per line, and a refusal reported with exit 1 and no facts.
func TestIntegrateCommand(t *testing.T) {
	t.Parallel()
	root := gitFixture(t)
	for _, kv := range [][2]string{{"user.name", "t"}, {"user.email", "t@t"}, {"commit.gpgsign", "false"}, {"maintenance.auto", "false"}} {
		gitIn(t, root, "config", kv[0], kv[1])
	}
	run := func(args ...string) (int, string, string) {
		t.Helper()
		var out, errOut bytes.Buffer
		code := Run(args, root, &out, &errOut)
		return code, out.String(), errOut.String()
	}
	write(t, root, "grove.yaml", "schema_version: 3\nrecords: docs/records\ntarget: main\n")
	gitIn(t, root, "commit", "-qam", "chore: target")
	code, out, stderr := run("integrate", "G-001")
	if code != 1 || out != "" || !strings.Contains(stderr, "grove: no branch holds G-001 in review; nothing to integrate") {
		t.Fatalf("code=%d out=%q stderr=%s", code, out, stderr)
	}
	gitIn(t, root, "checkout", "-q", "-b", "feature")
	run("update", "G-001", "--set", "status=active", "--commit")
	candidate := gitIn(t, root, "rev-parse", "HEAD")
	run("update", "G-001", "--set", "status=review", "--set", "candidate="+candidate, "--commit")
	if code, _, stderr := run("approve", "G-001", "Yes"); code != 0 {
		t.Fatal(stderr)
	}
	tip := gitIn(t, root, "rev-parse", "HEAD")
	gitIn(t, root, "checkout", "-q", "main")
	before := gitIn(t, root, "rev-parse", "HEAD")
	code, out, stderr = run("integrate", "G-001", "--cleanup")
	head := gitIn(t, root, "rev-parse", "HEAD")
	want := "approval: candidate " + candidate[:7] + " of G-001 approved on branch feature (Verdict on candidate " + candidate[:7] + ", " + today() + ": Yes)\n" +
		"merge: fast-forward main from " + before[:7] + " to " + tip[:7] + "\n" +
		"done: G-001 done at commit " + head[:7] + "\n" +
		"cleanup: deleted branch feature\n"
	if code != 0 || out != want {
		t.Fatalf("code=%d stderr=%s\nout:\n%s\nwant:\n%s", code, stderr, out, want)
	}
	if src := showJSON(t, root, "G-001")["source"].(string); !strings.Contains(src, "status: done\n") || !strings.Contains(src, "approved: \""+candidate+"\"\n") {
		t.Fatalf("record on main:\n%s", src)
	}
}
