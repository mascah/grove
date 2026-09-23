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
