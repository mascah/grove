package sweep

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mascah/grove/internal/attempt"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
	"github.com/mascah/grove/internal/update"
)

// TestMain lets the test binary be the owner a resolution attempt starts, as
// cmd/grove does.
func TestMain(m *testing.M) {
	if dir := os.Getenv(attempt.OwnerEnv); dir != "" {
		os.Exit(attempt.Own(dir))
	}
	os.Exit(m.Run())
}

const (
	work   = "---\nid: \"G-001\"\ntype: work\ntitle: First\nstatus: %s\n---\n\n## Outcome\n\nBody.\n"
	review = "---\nid: \"G-005\"\ntype: review\ntitle: Review of G-001\nstatus: current\nwork: [\"G-001\"]\nexamined: \"%s\"\n---\n\nFindings: none.\n\n%s\n"
	policy = "schema_version: 3\nrecords: grove\ntarget: main\npolicy:\n  budget: 5\n  resolve:\n    budget: 1\n  approve:\n    verify: [%s]\n    max_lines: 30\n    never: [grove.yaml, secret/**]\n  integrate: true\nrun:\n  permission_mode: acceptEdits\n"
)

var now = time.Date(2026, 9, 25, 18, 30, 0, 0, time.UTC)

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

// fixture is main with config, and branch worktree-G-001 in its checkout
// where G-001 is in review: its candidate writes files over code.txt's
// "base", and a review record, added after it, examined it and closes with
// closing. It returns main's checkout and the branch's.
func fixture(t *testing.T, config string, files map[string]string, closing string) (root, wt string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	root, wt = filepath.Join(dir, "repo"), filepath.Join(dir, "repo", ".claude", "worktrees", "worktree-G-001")
	git(t, dir, "init", "-q", "-b", "main", root)
	for _, kv := range [][2]string{{"user.name", "t"}, {"user.email", "t@t"}, {"commit.gpgsign", "false"}, {"maintenance.auto", "false"}} {
		git(t, root, "config", kv[0], kv[1])
	}
	write(t, root, ".gitignore", ".claude/worktrees/\n")
	write(t, root, "grove.yaml", config)
	write(t, root, "grove/G-001-first.md", strings.Replace(work, "%s", "active", 1))
	write(t, root, attempt.SkillPath, "---\nname: grove-work\n---\n")
	write(t, root, "code.txt", "base\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "init")
	git(t, root, "worktree", "add", "-q", "-b", "worktree-G-001", wt)
	for name, content := range files {
		write(t, wt, name, content)
	}
	git(t, wt, "add", "-A")
	git(t, wt, "commit", "-qm", "feat: the change")
	examined := git(t, wt, "rev-parse", "HEAD")
	write(t, wt, "grove/G-005-review.md", strings.Replace(strings.Replace(review, "%s", examined, 1), "%s", closing, 1))
	git(t, wt, "add", "-A")
	git(t, wt, "commit", "-qm", "docs: review")
	candidate := git(t, wt, "rev-parse", "HEAD")
	if _, err := update.Apply(wt, update.Request{ID: "G-001", Set: []update.Field{{Name: "status", Value: "review"}, {Name: "candidate", Value: candidate}}, Commit: true}, now, nil); err != nil {
		t.Fatal(err)
	}
	return root, wt
}

func record(t *testing.T, dir string) *project.Record {
	t.Helper()
	p, ds := project.Load(dir, dir)
	if len(ds) != 0 {
		t.Fatalf("%v", ds)
	}
	for _, r := range p.Records {
		if r.ID == "G-001" {
			return r
		}
	}
	t.Fatal("no G-001")
	return nil
}

func sweep(t *testing.T, root string) (*Sweep, []string) {
	t.Helper()
	s, err := Plan(root)
	if err != nil {
		t.Fatal(err)
	}
	var facts []string
	s.Run(now, func(f string) { facts = append(facts, f) })
	return s, facts
}

func TestSweepIntegratesACandidateInsideThePolicy(t *testing.T) {
	t.Parallel()
	root, _ := fixture(t, strings.Replace(policy, "%s", "grep -q change code.txt, test -f grove/G-005-review.md", 1), map[string]string{"code.txt": "the change\n"}, ClosingLine)
	s, err := Plan(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Items) != 1 || s.Items[0].Act != Integrate || !strings.Contains(s.Items[0].Why, "review G-005 has no open finding; 14 changed lines") {
		t.Fatalf("plan: %+v", s.Items)
	}
	if !strings.HasPrefix(s.Attribution, "policy grove.yaml sha256:") {
		t.Fatalf("attribution %q", s.Attribution)
	}
	var facts []string
	s.Run(now, func(f string) { facts = append(facts, f) })
	joined := strings.Join(facts, "\n")
	for _, want := range []string{"G-001: verified: the merge of", "G-001: approved under policy grove.yaml sha256:", "G-001: merge: fast-forward main", "G-001: done: G-001 done at commit"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing %q in:\n%s", want, joined)
		}
	}
	r := record(t, root)
	if r.Status != "done" || !update.Delegated(r) {
		t.Fatalf("record on main:\n%s", r.Source)
	}
	for _, want := range []string{
		": delegated under " + s.Attribution + ": review G-005 examined ",
		"verification passed (grep -q change code.txt; test -f grove/G-005-review.md); no Grove attempt is recorded as producing it",
		"Integrated under " + s.Attribution + " by fast-forwarding main from ",
	} {
		if !strings.Contains(string(r.Source), want) {
			t.Fatalf("missing %q in:\n%s", want, r.Source)
		}
	}
	if out := git(t, root, "worktree", "list", "--porcelain"); strings.Count(out, "worktree ") != 2 {
		t.Fatalf("the verification worktree remains:\n%s", out)
	}
	// Integrated, it is no longer a candidate in review anywhere but on its
	// branch, which the target now holds.
	if _, facts := sweep(t, root); len(facts) != 1 || !strings.Contains(facts[0], "G-001: skipped: ") {
		t.Fatalf("second sweep: %q", facts)
	}
}

func TestSweepLeavesAFailedVerificationUnchanged(t *testing.T) {
	t.Parallel()
	root, wt := fixture(t, strings.Replace(policy, "%s", "'true', 'echo broken; exit 3'", 1), map[string]string{"code.txt": "the change\n"}, ClosingLine)
	main, branch := git(t, root, "rev-parse", "main"), git(t, root, "rev-parse", "worktree-G-001")
	_, facts := sweep(t, root)
	if len(facts) != 1 || !strings.Contains(facts[0], "G-001: waits: verification of the merge with main at ") || !strings.Contains(facts[0], "echo broken; exit 3: exit status 3; last output: broken; nothing was approved") {
		t.Fatalf("facts %q", facts)
	}
	if git(t, root, "rev-parse", "main") != main || git(t, root, "rev-parse", "worktree-G-001") != branch || record(t, wt).Approved != "" {
		t.Fatal("a failed verification changed something")
	}
}

func TestSweepWaitsOutsideThePolicy(t *testing.T) {
	t.Parallel()
	for name, c := range map[string]struct {
		files   map[string]string
		closing string
		want    string
	}{
		"a never path":    {map[string]string{"secret/key": "x\n"}, ClosingLine, "changes secret/key, which never matches (secret/**)"},
		"the policy":      {map[string]string{"grove.yaml": strings.Replace(policy, "%s", "'true'", 1) + "# more\n"}, ClosingLine, "changes grove.yaml, which never matches (grove.yaml)"},
		"too many lines":  {map[string]string{"code.txt": strings.Repeat("line\n", 30)}, ClosingLine, "lines, over the policy's 30"},
		"an open finding": {map[string]string{"code.txt": "the change\n"}, "Open findings: 1", `review G-005 does not end with "Open findings: none"`},
		"a binary":        {map[string]string{"blob": "\x00\x01"}, ClosingLine, "changes the binary file blob"},
		"no closing line": {map[string]string{"code.txt": "the change\n"}, "", `does not end with "Open findings: none"`},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			root, _ := fixture(t, strings.Replace(policy, "%s", "'true'", 1), c.files, c.closing)
			main := git(t, root, "rev-parse", "main")
			_, facts := sweep(t, root)
			if len(facts) != 1 || !strings.HasPrefix(facts[0], "G-001: waits: ") || !strings.Contains(facts[0], c.want) {
				t.Fatalf("facts %q, want %q", facts, c.want)
			}
			if git(t, root, "rev-parse", "main") != main {
				t.Fatal("main moved")
			}
		})
	}
	t.Run("approved by the owner", func(t *testing.T) {
		t.Parallel()
		root, wt := fixture(t, strings.Replace(policy, "%s", "'true'", 1), map[string]string{"code.txt": "x\n"}, ClosingLine)
		if _, err := update.Approve(wt, "G-001", "Mine.", now); err != nil {
			t.Fatal(err)
		}
		if _, facts := sweep(t, root); len(facts) != 1 || !strings.Contains(facts[0], "waits: approved by the owner; integrating it is the owner's: grove integrate G-001") {
			t.Fatalf("facts %q", facts)
		}
	})
	t.Run("no approve section", func(t *testing.T) {
		t.Parallel()
		root, _ := fixture(t, "schema_version: 3\nrecords: grove\ntarget: main\npolicy:\n  budget: 1\n  resolve: {}\n", map[string]string{"code.txt": "x\n"}, ClosingLine)
		if _, facts := sweep(t, root); len(facts) != 1 || !strings.Contains(facts[0], "the policy does not approve") {
			t.Fatalf("facts %q", facts)
		}
	})
}

func TestSweepRefusesWithoutAPolicy(t *testing.T) {
	t.Parallel()
	root, wt := fixture(t, "schema_version: 3\nrecords: grove\ntarget: main\n", map[string]string{"code.txt": "x\n"}, ClosingLine)
	if _, err := Plan(root); err == nil || !strings.Contains(err.Error(), "grove.yaml has no policy: nothing is automatic") {
		t.Fatalf("got %v", err)
	}
	if _, err := Plan(wt); err == nil || !strings.Contains(err.Error(), "sweep runs in the checkout of the target main") {
		t.Fatalf("from the branch: %v", err)
	}
}

func TestSweepResolvesAConflictOncePerTargetCommit(t *testing.T) {
	if testing.Short() {
		t.Skip("waits on a fake provider's attempt")
	}
	root, wt := fixture(t, strings.Replace(policy, "%s", "'true'", 1), map[string]string{"code.txt": "branch\n"}, ClosingLine)
	write(t, root, "code.txt", "main\n")
	git(t, root, "commit", "-qam", "main moves")
	fake := filepath.Join(t.TempDir(), "claude")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\n[ \"$1\" = --version ] && { echo 'fake 0.1'; exit 0; }\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(attempt.ClaudeEnv, fake)

	s, facts := sweep(t, root)
	if len(s.Items) != 1 || s.Items[0].Act != Resolve || !strings.Contains(s.Items[0].Why, "in code.txt; one resolution attempt, budget 1 USD") {
		t.Fatalf("plan %+v", s.Items)
	}
	if last := facts[len(facts)-1]; !strings.Contains(last, "resolution attempt ") || !strings.Contains(last, "started under "+s.Attribution+", budget 1 USD") {
		t.Fatalf("facts %q", facts)
	}
	r := record(t, wt)
	if r.Status != "active" || !strings.Contains(string(r.Source), ": delegated under "+s.Attribution+", budget 1 USD: conflicts with main at ") {
		t.Fatalf("record:\n%s", r.Source)
	}
	deadline := time.Now().Add(20 * time.Second)
	for {
		views, err := attempt.List(root, "G-001")
		if err != nil {
			t.Fatal(err)
		}
		if len(views) == 1 && views[0].Status == attempt.Finished {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("attempt did not finish: %+v", views)
		}
		time.Sleep(50 * time.Millisecond)
	}
	// The attempt handed the same candidate back without resolving: the
	// conflict with the same main commit waits for the owner.
	if _, err := update.Apply(wt, update.Request{ID: "G-001", Set: []update.Field{{Name: "status", Value: "review"}}, Commit: true}, now, nil); err != nil {
		t.Fatal(err)
	}
	if _, facts := sweep(t, root); len(facts) != 1 || !strings.Contains(facts[0], "waits: ") || !strings.Contains(facts[0], "again after a resolution of main at ") {
		t.Fatalf("second sweep %q", facts)
	}
}

func TestSweepKeepsResolutionsInsideTheBudget(t *testing.T) {
	t.Parallel()
	config := strings.Replace(strings.Replace(policy, "%s", "'true'", 1), "  resolve:\n    budget: 1\n", "  resolve:\n    budget: 6\n", 1)
	root, _ := fixture(t, config, map[string]string{"code.txt": "branch\n"}, ClosingLine)
	write(t, root, "code.txt", "main\n")
	git(t, root, "commit", "-qam", "main moves")
	s, err := Plan(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Items) != 1 || s.Items[0].Act != Wait || !strings.Contains(s.Items[0].Why, "a 6 USD attempt would pass the policy's 5 USD for one sweep") {
		t.Fatalf("plan %+v", s.Items)
	}
}
