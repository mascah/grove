package integrate

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
	"github.com/mascah/grove/internal/update"
)

const (
	config = "schema_version: 3\nrecords: grove\ntarget: main\n"
	work   = "---\nid: \"G-001\"\ntype: work\ntitle: First\nstatus: %s\n---\n\n## Outcome\n\nBody.\n"
	review = "---\nid: \"G-005\"\ntype: review\ntitle: Review of G-001\nstatus: current\nwork: [\"G-001\"]\nexamined: \"%s\"\n---\n\nFindings.\n"
)

var now = time.Date(2026, 9, 22, 18, 30, 0, 0, time.UTC)

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	full := append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-c", "maintenance.auto=false"}, args...)
	out, err := repo.Command(context.Background(), dir, full...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func write(t *testing.T, root, rel, content string) {
	t.Helper()
	full := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func record(t *testing.T, root string) *project.Record {
	t.Helper()
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		t.Fatalf("%v", ds)
	}
	return p.Records[0]
}

// fixture is a main checkout with target main, and a linked worktree on
// branch feature where G-001 is in review, approved, with a review record in
// its candidate. It returns the main root, the worktree, and the candidate.
func fixture(t *testing.T, approve bool, sub ...string) (root, wt, candidate string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	root, wt = filepath.Join(dir, "repo"), filepath.Join(dir, "feat")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, root, "init", "-q", "-b", "main")
	for _, kv := range [][2]string{{"user.name", "t"}, {"user.email", "t@t"}, {"commit.gpgsign", "false"}, {"maintenance.auto", "false"}} {
		git(t, root, "config", kv[0], kv[1]) // the product's commits use the repository's own identity
	}
	// sub places the project under a prefix, with the code above it.
	project, wtProject := filepath.Join(append([]string{root}, sub...)...), filepath.Join(append([]string{wt}, sub...)...)
	write(t, project, "grove.yaml", config)
	write(t, project, "grove/G-001-first.md", strings.Replace(work, "%s", "proposed", 1))
	write(t, root, "code.txt", "before\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "init")
	git(t, root, "worktree", "add", "-q", "-b", "feature", wt)
	write(t, wtProject, "grove/G-001-first.md", strings.Replace(work, "%s", "active", 1))
	write(t, wt, "code.txt", "the change\n")
	git(t, wt, "add", "-A")
	git(t, wt, "commit", "-qm", "feat: implement")
	write(t, wtProject, "grove/G-005-review.md", strings.Replace(review, "%s", git(t, wt, "rev-parse", "HEAD"), 1))
	git(t, wt, "add", "-A")
	git(t, wt, "commit", "-qm", "docs: evidence and review")
	candidate = git(t, wt, "rev-parse", "HEAD")
	if _, err := update.Apply(wtProject, update.Request{ID: "G-001", Set: []update.Field{{Name: "status", Value: "review"}, {Name: "candidate", Value: candidate}}, Commit: true}, now, nil); err != nil {
		t.Fatal(err)
	}
	if approve {
		if _, err := update.Approve(wtProject, "G-001", "Ship it.", now); err != nil {
			t.Fatal(err)
		}
	}
	return project, wtProject, candidate
}

func run(t *testing.T, root, cwd string, cleanup bool) ([]string, error) {
	t.Helper()
	var facts []string
	err := Run(Request{Root: root, ID: "G-001", Cwd: cwd, Cleanup: cleanup}, now, func(f string) { facts = append(facts, f) })
	return facts, err
}

// refused asserts a refusal that leaves main, its record, and its files as they were.
func refused(t *testing.T, root string, cleanup bool, want string) []string {
	t.Helper()
	head := git(t, root, "rev-parse", "HEAD")
	before := record(t, root).Source
	facts, err := run(t, root, root, cleanup)
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("got %v, want %q; facts %q", err, want, facts)
	}
	if git(t, root, "rev-parse", "HEAD") != head || string(record(t, root).Source) != string(before) || git(t, root, "status", "--porcelain", "--untracked-files=no") != "" {
		t.Fatal("a refusal changed the target")
	}
	if _, err := repo.Git(root, "rev-parse", "-q", "--verify", "MERGE_HEAD"); err == nil {
		t.Fatal("a merge was left in progress")
	}
	return facts
}

func TestIntegrateFastForwardAndCleanup(t *testing.T) {
	t.Parallel()
	root, wt, candidate := fixture(t, true)
	before, tip := git(t, root, "rev-parse", "HEAD"), git(t, wt, "rev-parse", "HEAD")
	facts, err := run(t, root, root, true)
	if err != nil {
		t.Fatalf("%v; facts %q", err, facts)
	}
	head := git(t, root, "rev-parse", "HEAD")
	want := []string{
		"approval: candidate " + candidate[:7] + " of G-001 approved on branch feature at " + tip[:7] + " (Verdict on candidate " + candidate[:7] + ", 2026-09-22: Ship it.)",
		"merge: fast-forward main from " + before[:7] + " to " + tip[:7],
		"done: G-001 done at commit " + head[:7],
		"cleanup: removed worktree " + wt,
		"cleanup: deleted branch feature",
	}
	if strings.Join(facts, "\n") != strings.Join(want, "\n") {
		t.Fatalf("facts:\n%s\nwant:\n%s", strings.Join(facts, "\n"), strings.Join(want, "\n"))
	}
	r := record(t, root)
	if r.Status != "done" || r.Candidate != candidate || r.Approved != candidate || !strings.Contains(string(r.Source), "Verdict on candidate") {
		t.Fatalf("record on main: %+v", r)
	}
	if files := git(t, root, "show", "--format=", "--name-only", "HEAD"); files != "grove/G-001-first.md" {
		t.Fatalf("the done commit must hold the record alone: %q", files)
	}
	if msg := git(t, root, "log", "-1", "--format=%s"); msg != "docs(G-001): set status=done" {
		t.Fatalf("message: %q", msg)
	}
	if _, err := os.Stat(wt); !os.IsNotExist(err) {
		t.Fatal("the worktree remains")
	}
	if git(t, root, "branch", "--list", "feature") != "" {
		t.Fatal("the branch remains")
	}
	// Running again finds no candidate: the branch is gone.
	refused(t, root, false, "no branch holds G-001 in review; nothing to integrate")
}

func TestIntegrateMergeCommitWhenTheTargetMoved(t *testing.T) {
	t.Parallel()
	root, _, candidate := fixture(t, true)
	write(t, root, "notes.txt", "elsewhere\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "docs: notes")
	before := git(t, root, "rev-parse", "HEAD")
	facts, err := run(t, root, root, false)
	if err != nil || len(facts) != 3 {
		t.Fatalf("%v; facts %q", err, facts)
	}
	merge := git(t, root, "rev-parse", "HEAD~1")
	if facts[1] != "merge: merge commit "+merge[:7]+" on main (was "+before[:7]+")" || git(t, root, "rev-list", "--parents", "-1", merge) == merge {
		t.Fatalf("facts %q; parents %q", facts, git(t, root, "rev-list", "--parents", "-1", merge))
	}
	if r := record(t, root); r.Status != "done" || r.Candidate != candidate {
		t.Fatalf("record on main: %+v", r)
	}
	// A second integration has nothing to merge and nothing to close.
	facts, err = run(t, root, root, false)
	if err != nil || facts[1] != "merge: nothing to merge; feature is already in main at "+git(t, root, "rev-parse", "--short=7", "HEAD") || facts[2] != "done: G-001 was already done here" {
		t.Fatalf("%v; facts %q", err, facts)
	}
}

func TestIntegrateRefusesConflictsBeforeAnythingChanges(t *testing.T) {
	t.Parallel()
	t.Run("a file conflict", func(t *testing.T) {
		t.Parallel()
		root, _, _ := fixture(t, true)
		write(t, root, "code.txt", "a different change\n")
		git(t, root, "commit", "-qam", "fix: on main")
		facts := refused(t, root, false, "merge of feature into main refused: CONFLICT (content): Merge conflict in code.txt; Automatic merge failed; fix conflicts and then commit the result.; main is unchanged at")
		if len(facts) != 1 || !strings.HasPrefix(facts[0], "approval: ") {
			t.Fatalf("facts %q", facts)
		}
	})
	t.Run("the record changed on the target", func(t *testing.T) {
		t.Parallel()
		root, _, _ := fixture(t, true)
		write(t, root, "grove/G-001-first.md", strings.Replace(work, "%s", "abandoned", 1))
		git(t, root, "commit", "-qam", "docs: abandon on main")
		refused(t, root, false, "main is unchanged at")
	})
}

func TestIntegrateRefusesWhatIsNotReadyToMerge(t *testing.T) {
	t.Parallel()
	t.Run("stale approval", func(t *testing.T) {
		t.Parallel()
		root, wt, candidate := fixture(t, true)
		write(t, wt, "code.txt", "later\n")
		git(t, wt, "commit", "-qam", "fix: after the handoff")
		tip := git(t, wt, "rev-parse", "HEAD")
		refused(t, root, false, "commits after candidate "+candidate[:7]+" on feature change code.txt: the tip "+tip[:7]+" is a new candidate; approve it before integrating")
	})
	t.Run("dirty target", func(t *testing.T) {
		t.Parallel()
		root, _, _ := fixture(t, true)
		write(t, root, "code.txt", "edited\n")
		head := git(t, root, "rev-parse", "HEAD")
		if _, err := run(t, root, root, false); err == nil || !strings.Contains(err.Error(), "the checkout of main has uncommitted changes") || git(t, root, "rev-parse", "HEAD") != head {
			t.Fatalf("%v", err)
		}
		write(t, root, "untracked.txt", "fine\n") // untracked files do not block a merge
		git(t, root, "checkout", "-q", "--", "code.txt")
		if _, err := run(t, root, root, false); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("wrong branch", func(t *testing.T) {
		t.Parallel()
		_, wt, _ := fixture(t, true)
		if _, err := run(t, wt, wt, false); err == nil || !strings.Contains(err.Error(), "integrate runs in the checkout of the target main; this one is on feature") {
			t.Fatalf("%v", err)
		}
	})
	t.Run("no target", func(t *testing.T) {
		t.Parallel()
		root, _, _ := fixture(t, true)
		write(t, root, "grove.yaml", "schema_version: 3\nrecords: grove\n")
		git(t, root, "commit", "-qam", "chore: no target")
		refused(t, root, false, "integration needs target: BRANCH in grove.yaml")
	})
	t.Run("unapproved", func(t *testing.T) {
		t.Parallel()
		root, wt, _ := fixture(t, false)
		refused(t, root, false, "G-001 is in review on feature but not approved: run grove approve G-001 VERDICT in "+wt+" first")
	})
	t.Run("two branches", func(t *testing.T) {
		t.Parallel()
		root, _, _ := fixture(t, true)
		git(t, root, "branch", "feature2", "feature")
		facts := refused(t, root, false, "G-001 is approved in review on several branches (feature, feature2); integrate needs one")
		if len(facts) != 0 {
			t.Fatalf("facts %q", facts)
		}
	})
	t.Run("not in review", func(t *testing.T) {
		t.Parallel()
		root, wt, _ := fixture(t, true)
		if _, err := update.Feedback(wt, "G-001", "no", now); err != nil {
			t.Fatal(err)
		}
		refused(t, root, false, "no branch holds G-001 in review")
	})
}

func TestIntegrateCleanupKeepsWhatGitOrTheSessionHolds(t *testing.T) {
	t.Parallel()
	t.Run("a dirty worktree", func(t *testing.T) {
		t.Parallel()
		root, wt, _ := fixture(t, true)
		write(t, wt, "scratch.txt", "not committed\n")
		facts, err := run(t, root, root, true)
		if err == nil || !strings.Contains(err.Error(), "cleanup incomplete; the integration stands") {
			t.Fatalf("%v; facts %q", err, facts)
		}
		if len(facts) != 5 || !strings.HasPrefix(facts[3], "cleanup: kept worktree "+wt+": git worktree: ") || !strings.Contains(facts[3], "--force") ||
			!strings.HasPrefix(facts[4], "cleanup: kept branch feature: git branch: ") {
			t.Fatalf("facts %q", facts)
		}
		if r := record(t, root); r.Status != "done" {
			t.Fatalf("the integration must stand: %+v", r)
		}
		if _, err := os.Stat(filepath.Join(wt, "scratch.txt")); err != nil {
			t.Fatal("the worktree's file is gone")
		}
	})
	t.Run("ignored files", func(t *testing.T) {
		t.Parallel()
		root, wt, _ := fixture(t, true)
		write(t, root, ".git/info/exclude", "*.log\n")
		write(t, wt, "junk.log", "ignored, so git worktree remove would delete it\n")
		facts, err := run(t, root, root, true)
		if err == nil || len(facts) != 5 || facts[3] != "cleanup: kept worktree "+wt+": it holds ignored files (junk.log); remove them or the worktree by hand" {
			t.Fatalf("%v; facts %q", err, facts)
		}
		if _, err := os.Stat(filepath.Join(wt, "junk.log")); err != nil {
			t.Fatal("the ignored file is gone")
		}
	})
	t.Run("the session's directory", func(t *testing.T) {
		t.Parallel()
		root, wt, _ := fixture(t, true)
		facts, err := run(t, root, filepath.Join(wt, "grove"), true)
		if err == nil || len(facts) != 5 || facts[3] != "cleanup: kept worktree "+wt+": it holds this process's working directory" || !strings.HasPrefix(facts[4], "cleanup: kept branch feature: ") {
			t.Fatalf("%v; facts %q", err, facts)
		}
	})
}

// TestIntegrateReportsADoneWriteThatFailedAfterTheMerge: the merge stands, the
// record is not done, and the error says how to finish.
func TestIntegrateReportsADoneWriteThatFailedAfterTheMerge(t *testing.T) {
	t.Parallel()
	root, wt, _ := fixture(t, true)
	tip := git(t, wt, "rev-parse", "HEAD")
	hooks := filepath.Join(root, "hooks")
	write(t, root, "hooks/pre-commit", "#!/bin/sh\nexit 1\n")
	if err := os.Chmod(filepath.Join(hooks, "pre-commit"), 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, root, "config", "core.hooksPath", hooks)
	facts, err := run(t, root, root, true)
	if err == nil || !strings.Contains(err.Error(), "merged as "+tip[:7]+", but marking G-001 done failed: ") || !strings.Contains(err.Error(), "commit the staged record here: git commit -m 'docs(G-001): set status=done' -- grove/G-001-first.md") {
		t.Fatalf("%v; facts %q", err, facts)
	}
	if len(facts) != 2 || git(t, root, "rev-parse", "HEAD") != tip {
		t.Fatalf("the merge must stand: facts %q", facts)
	}
	if r := record(t, root); r.Status != "done" || !strings.HasPrefix(git(t, root, "status", "--porcelain"), "M  grove/G-001-first.md") { // staged, uncommitted
		t.Fatalf("the applied but uncommitted done write: %+v %q", r, git(t, root, "status", "--porcelain"))
	}
	if _, err := os.Stat(wt); err != nil {
		t.Fatal("cleanup must not run after a failed done write")
	}
}
