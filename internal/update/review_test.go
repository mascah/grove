package update

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/project"
)

const review = "---\nid: \"G-005\"\ntype: review\ntitle: Review of G-001\nstatus: current\nwork: [\"G-001\"]\nexamined: \"%s\"\n---\n\nFindings.\n"

// reviewFixture is gitProject with an identity for commits, a feature branch
// on which G-001 is in review with its candidate, and a review record of it.
func reviewFixture(t *testing.T) (root, candidate string) {
	t.Helper()
	root = gitProject(t)
	for _, kv := range [][2]string{{"user.name", "t"}, {"user.email", "t@t"}, {"commit.gpgsign", "false"}, {"maintenance.auto", "false"}} {
		git(t, root, "config", kv[0], kv[1])
	}
	git(t, root, "checkout", "-q", "-b", "feature")
	write(t, root, "grove/work/G-001-first.md", strings.Replace(work, "status: proposed", "status: active", 1))
	write(t, root, "code.txt", "the change\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "feat: implement")
	// The review examined the implementation commit; the candidate is the
	// evidence commit that adds it, as the guide's handoff does, so that only
	// the work record changes after the candidate.
	write(t, root, "grove/reviews/G-005-review.md", strings.Replace(review, "%s", git(t, root, "rev-parse", "HEAD"), 1))
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "docs: review")
	candidate = git(t, root, "rev-parse", "HEAD")
	if _, err := Apply(root, Request{ID: "G-001", Set: []Field{{"status", "review"}, {"candidate", candidate}}, Commit: true}, now, nil); err != nil {
		t.Fatal(err)
	}
	return root, candidate
}

// TestApproveAndFeedback covers G-044's two dispositions: approval binds the
// candidate and quotes the verdict, feedback reopens the work with the text
// and no approval, each as one commit of the record alone, and the earlier
// review record is left as it was.
func TestApproveAndFeedback(t *testing.T) {
	t.Parallel()
	root, candidate := reviewFixture(t)
	reviewBefore := read(t, root, "grove/reviews/G-005-review.md")
	tip := git(t, root, "rev-parse", "HEAD")
	if _, err := Approve(root, "G-001", "  ", now); err == nil || !strings.Contains(err.Error(), "a verdict is required") {
		t.Fatalf("blank verdict: %v", err)
	}
	if _, err := Approve(root, "G-002", "fine", now); err == nil || !strings.Contains(err.Error(), "G-002 is proposed, not in review") {
		t.Fatalf("not in review: %v", err)
	}
	if _, err := Feedback(root, "G-003", "more", now); err == nil || !strings.Contains(err.Error(), "G-003 is a question, not work") {
		t.Fatalf("not work: %v", err)
	}
	res, err := Approve(root, "G-001", "Ship it.\n", now)
	if err != nil || !res.Changed || res.Commit == "" || res.Commit != git(t, root, "rev-parse", "HEAD") {
		t.Fatalf("approve: %+v %v", res, err)
	}
	r := record(t, root, "G-001")
	if r.Approved != candidate || r.Status != "review" || !strings.HasSuffix(string(r.Source), "Body --- stays.\n\nVerdict on candidate "+candidate[:7]+", 2026-09-19: Ship it.\n") {
		t.Fatalf("approved record:\n%s", r.Source)
	}
	if files := git(t, root, "show", "--format=", "--name-only", "HEAD"); files != "grove/work/G-001-first.md" {
		t.Fatalf("the approval commit must hold the record alone: %q", files)
	}
	if msg := git(t, root, "log", "-1", "--format=%s"); msg != "docs(G-001): set approved="+candidate+" note" {
		t.Fatalf("message: %q", msg)
	}
	if git(t, root, "status", "--porcelain") != "" {
		t.Fatal("approve left the checkout dirty")
	}
	if _, err := Approve(root, "G-001", "again", now); err == nil || !strings.Contains(err.Error(), "already approved; integrate it") {
		t.Fatalf("second approval: %v", err)
	}
	res, err = Feedback(root, "G-001", "Needs a test for the empty case.", now)
	if err != nil || !res.Changed || res.Commit != git(t, root, "rev-parse", "HEAD") {
		t.Fatalf("feedback: %+v %v", res, err)
	}
	r = record(t, root, "G-001")
	if r.Status != "active" || r.Approved != "" || r.Candidate != candidate || !strings.HasSuffix(string(r.Source), "Ship it.\n\nFeedback on candidate "+candidate[:7]+", 2026-09-19: Needs a test for the empty case.\n") {
		t.Fatalf("record after feedback:\n%s", r.Source)
	}
	if msg := git(t, root, "log", "-1", "--format=%s"); msg != "docs(G-001): set status=active unset approved note" {
		t.Fatalf("message: %q", msg)
	}
	if read(t, root, "grove/reviews/G-005-review.md") != reviewBefore {
		t.Fatal("feedback changed the review record")
	}
	if _, err := Feedback(root, "G-001", "more", now); err == nil || !strings.Contains(err.Error(), "G-001 is active, not in review") {
		t.Fatalf("feedback on active work: %v", err)
	}
	// Both candidates and every disposition stay in history.
	if log := git(t, root, "log", "--format=%s", tip+"..HEAD"); log != "docs(G-001): set status=active unset approved note\ndocs(G-001): set approved="+candidate+" note" {
		t.Fatalf("history: %q", log)
	}
}

// TestJudgingNeedsTheCandidateCheckoutClean covers the refusals: a checkout
// without the candidate, a record with uncommitted changes, and, for
// approval, commits after the candidate that make the tip a new one.
func TestJudgingNeedsTheCandidateCheckoutClean(t *testing.T) {
	t.Parallel()
	root, candidate := reviewFixture(t)
	refuse := func(judge func() (Result, error), want string) {
		t.Helper()
		before := read(t, root, "grove/work/G-001-first.md")
		head := git(t, root, "rev-parse", "HEAD")
		_, err := judge()
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Fatalf("got %v, want %q", err, want)
		}
		if read(t, root, "grove/work/G-001-first.md") != before || git(t, root, "rev-parse", "HEAD") != head {
			t.Fatal("a refusal changed the record or the branch")
		}
	}
	approve := func() (Result, error) { return Approve(root, "G-001", "ok", now) }
	feedback := func() (Result, error) { return Feedback(root, "G-001", "no", now) }
	// An uncommitted edit of the record is never committed by a judgment.
	src := read(t, root, "grove/work/G-001-first.md")
	write(t, root, "grove/work/G-001-first.md", src+"\nA local edit.\n")
	refuse(approve, "grove/work/G-001-first.md has uncommitted changes in this checkout")
	refuse(feedback, "grove/work/G-001-first.md has uncommitted changes in this checkout")
	git(t, root, "checkout", "-q", "--", ".")
	// A commit after the candidate that changes anything else is a new candidate.
	write(t, root, "code.txt", "a later change\n")
	git(t, root, "commit", "-qam", "fix: later")
	tip := git(t, root, "rev-parse", "HEAD")
	refuse(approve, "commits after candidate "+candidate[:7]+" change code.txt: the tip "+tip[:7]+" is a new candidate; set candidate="+tip[:7])
	if _, err := Feedback(root, "G-001", "still no", now); err != nil { // feedback on a moved tip is still feedback
		t.Fatalf("feedback after a later commit: %v", err)
	}
	if _, err := Apply(root, Request{ID: "G-001", Set: []Field{{"status", "review"}, {"candidate", tip}}, Commit: true}, now, nil); err != nil {
		t.Fatal(err)
	}
	// The target's checkout does not hold the candidate before the merge.
	git(t, root, "checkout", "-q", "main")
	write(t, root, "grove/work/G-001-first.md", strings.Replace(work, "status: proposed", "status: review\ncandidate: \""+tip+"\"", 1))
	git(t, root, "commit", "-qam", "docs: pretend review on main")
	refuse(approve, "this checkout does not hold candidate "+tip+"; run this in the checkout of the branch that has G-001 in review")
	refuse(feedback, "this checkout does not hold candidate "+tip)
}

// TestAppendKeepsEveryByte covers the one body edit: the paragraph lands
// after the body in the file's own line ending, and nothing before it moves.
func TestAppendKeepsEveryByte(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, source, want string }{
		{"lf", "---\nid: \"G-001\"\ntype: work\ntitle: T\nstatus: proposed\n---\n\nBody.\n", "---\nid: \"G-001\"\ntype: work\ntitle: T\nstatus: proposed\nupdated: \"2026-09-19T18:30:00Z\"\n---\n\nBody.\n\nA note.\nTwo lines.\n"},
		{"no final newline", "---\nid: \"G-001\"\ntype: work\ntitle: T\nstatus: proposed\n---\nBody", "---\nid: \"G-001\"\ntype: work\ntitle: T\nstatus: proposed\nupdated: \"2026-09-19T18:30:00Z\"\n---\nBody\n\nA note.\nTwo lines.\n"},
		{"crlf", "---\r\nid: \"G-001\"\r\ntype: work\r\ntitle: T\r\nstatus: proposed\r\n---\r\n\r\nBody.\r\n", "---\r\nid: \"G-001\"\r\ntype: work\r\ntitle: T\r\nstatus: proposed\r\nupdated: \"2026-09-19T18:30:00Z\"\r\n---\r\n\r\nBody.\r\n\r\nA note.\r\nTwo lines.\r\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := gitProject(t)
			write(t, root, "grove/work/G-001-first.md", tc.source)
			res, err := Apply(root, Request{ID: "G-001", Append: " A note.\nTwo lines. \n"}, now, nil)
			if err != nil || !res.Changed {
				t.Fatalf("%+v %v", res, err)
			}
			if got := read(t, root, "grove/work/G-001-first.md"); got != tc.want {
				t.Fatalf("got:\n%q\nwant:\n%q", got, tc.want)
			}
		})
	}
	root := gitProject(t)
	before := read(t, root, "grove/work/G-001-first.md")
	if _, err := Apply(root, Request{ID: "G-001", Append: "\xff"}, now, nil); err == nil || !strings.Contains(err.Error(), "nonempty valid UTF-8") || read(t, root, "grove/work/G-001-first.md") != before {
		t.Fatalf("invalid append: %v", err)
	}
	if res, err := Apply(root, Request{ID: "G-001", Expect: project.Revision([]byte(before))}, now, nil); err != nil || res.Changed {
		t.Fatalf("nothing requested is a no-op: %+v %v", res, err)
	}
	_ = bytes.Equal
}
