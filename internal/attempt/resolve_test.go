package attempt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/mascah/grove/internal/versions"
)

// conflicted is fixture's repository with G-001 in review on worktree-G-001,
// in its checkout, whose candidate changes shared.txt as main does too since
// the branch left it. It returns the main checkout, the branch's checkout,
// the candidate and main's tip.
func conflicted(t *testing.T) (root, wt, candidate, tip string) {
	t.Helper()
	root = fixture(t)
	write(t, root, "shared.txt", "base\n")
	git(t, root, "add", "shared.txt")
	git(t, root, "commit", "-qm", "shared")
	wt = filepath.Join(root, ".claude", "worktrees", "worktree-G-001")
	git(t, root, "worktree", "add", "-q", "-b", "worktree-G-001", wt)
	write(t, wt, "shared.txt", "branch\n")
	git(t, wt, "commit", "-qam", "the change")
	candidate = git(t, wt, "rev-parse", "HEAD")
	write(t, wt, "grove/G-001-first.md", strings.Replace(fmt.Sprintf(work, "review"), "\n---\n\n## Outcome", "\ncandidate: \""+candidate+"\"\n---\n\n## Outcome", 1))
	git(t, wt, "commit", "-qam", "handoff")
	write(t, root, "shared.txt", "main\n")
	git(t, root, "commit", "-qam", "main moves")
	return root, wt, candidate, git(t, root, "rev-parse", "HEAD")
}

// resolving is a fake agent that does what the mandate says: merge the
// commit the feedback names, resolve shared.txt, and hand the merge off.
const resolving = `T=$(grep -o 'git merge [0-9a-f]\{40\}' grove/G-001-first.md | head -1 | cut -d' ' -f3)
G="git -c user.name=f -c user.email=f@f -c commit.gpgsign=false"
$G merge -q "$T" >/dev/null 2>&1
echo resolved > shared.txt
$G add shared.txt
$G commit -q --no-edit
M=$($G rev-parse HEAD)
sed -e 's/^status: active/status: review/' -e "s/^candidate: .*/candidate: \"$M\"/" grove/G-001-first.md > r.tmp && mv r.tmp grove/G-001-first.md
$G commit -qam handoff
`

func resolve(root string, shown *versions.Merge, at time.Time) (*Launch, []string, error) {
	var facts []string
	l, err := Resolve(Request{Root: root, IDs: []string{"G-001"}, BudgetUSD: "1", PermissionMode: "acceptEdits"}, shown, at, func(f string) { facts = append(facts, f) })
	return l, facts, err
}

// recordOn reads G-001 as a branch holds it.
func recordOn(t *testing.T, root, ref string) string {
	return git(t, root, "show", ref+":grove/G-001-first.md")
}

func TestResolveRefusals(t *testing.T) {
	root, wt, candidate, tip := conflicted(t)
	fake(t, "exit 0")
	before := git(t, root, "rev-parse", "worktree-G-001")
	unchanged := func(t *testing.T) {
		t.Helper()
		if now := git(t, root, "rev-parse", "worktree-G-001"); now != before || git(t, wt, "status", "--porcelain") != "" {
			t.Fatal("a refusal wrote something")
		}
		if views, _ := List(root, "G-001"); len(views) != 0 {
			t.Fatalf("a refusal started %d attempts", len(views))
		}
	}
	for name, c := range map[string]struct {
		req   Request
		shown *versions.Merge
		want  string
	}{
		"a selection":          {Request{IDs: []string{"G-001", "G-002"}}, nil, "resolve takes one work ID"},
		"a bound":              {Request{IDs: []string{"G-001"}, Until: "plan"}, nil, "do not apply"},
		"a branch":             {Request{IDs: []string{"G-001"}, Branch: "x"}, nil, "do not apply"},
		"not in review":        {Request{IDs: []string{"G-009"}}, nil, "no branch holds G-009 in review"},
		"another candidate":    {Request{IDs: []string{"G-001"}}, &versions.Merge{Commit: tip, Target: tip}, "look again"},
		"a moved target":       {Request{IDs: []string{"G-001"}}, &versions.Merge{Commit: candidate, Target: candidate}, "moved from"},
		"no budget or mode":    {Request{IDs: []string{"G-001"}}, nil, "run requires --budget"},
		"no provider":          {Request{IDs: []string{"G-001"}, BudgetUSD: "1", PermissionMode: "x"}, nil, "provider executable is not available"},
		"a shown fact matches": {Request{IDs: []string{"G-001"}, BudgetUSD: "1", PermissionMode: "x"}, &versions.Merge{Commit: candidate[:7], Target: tip}, "provider executable is not available"},
	} {
		t.Run(name, func(t *testing.T) {
			if strings.Contains(c.want, "provider") {
				t.Setenv(ClaudeEnv, filepath.Join(t.TempDir(), "missing"))
			}
			c.req.Root = root
			if _, err := Resolve(c.req, c.shown, now, func(string) {}); err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("got %v, want %q", err, c.want)
			}
			unchanged(t)
		})
	}

	// An attempt of the work that may be running.
	dir, err := Dir(root)
	if err != nil {
		t.Fatal(err)
	}
	adir := filepath.Join(dir, "G-001.20260922T170000Z")
	if err := os.MkdirAll(adir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := writeJSON(filepath.Join(adir, "attempt.json"), &Launch{Attempt: "G-001.20260922T170000Z", Work: "G-001"}); err != nil {
		t.Fatal(err)
	}
	lock, err := os.OpenFile(filepath.Join(adir, "owner.lock"), os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		t.Fatal(err)
	}
	if _, _, err := resolve(root, nil, now); err == nil || !strings.Contains(err.Error(), "is running; stop it or wait") {
		t.Fatalf("a running attempt: %v", err)
	}
	lock.Close()
	os.RemoveAll(adir)
	unchanged(t)

	// Resolved by hand: the candidate merges cleanly, so there is nothing to do.
	git(t, root, "revert", "--no-edit", "HEAD")
	if _, _, err := resolve(root, nil, now); err == nil || !strings.Contains(err.Error(), "there is no conflict to resolve") {
		t.Fatalf("a clean merge: %v", err)
	}
	unchanged(t)

	// Without a target nothing conflicts with anything.
	write(t, root, "grove.yaml", "schema_version: 3\nrecords: grove\n")
	if _, _, err := resolve(root, nil, now); err == nil || !strings.Contains(err.Error(), "resolve needs target") {
		t.Fatalf("no target: %v", err)
	}
	unchanged(t)
}

func TestResolveCleanly(t *testing.T) {
	skipShort(t)
	root, wt, candidate, tip := conflicted(t)
	mark, _ := fake(t, initLine+"\n"+resolving+resultLine("success", false))
	shown := &versions.Merge{Commit: candidate, Target: tip}
	l, facts, err := resolve(root, shown, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(facts) < 2 || !strings.HasPrefix(facts[0], "feedback: G-001 is active again on branch worktree-G-001") || !strings.Contains(facts[0], "resolve shared.txt") {
		t.Fatalf("facts %q", facts)
	}
	if l.Branch != "worktree-G-001" || !samePath(l.Worktree, wt) || !l.WorktreeReused || strings.Join(l.Selection.Selected, " ") != "G-001" {
		t.Fatalf("launch %+v", l)
	}
	// The mandate is in the record, where the headless guide reads it.
	rec := recordOn(t, root, l.Base)
	want := fmt.Sprintf("Feedback on candidate %s, 2026-09-22: conflicts with main at %s in shared.txt. Resolve only that (grove resolve): in this branch, git merge %s,", candidate[:7], tip[:7], tip)
	if !strings.Contains(rec, want) || !strings.Contains(rec, "status: active") {
		t.Fatalf("the record the attempt starts from lacks %q:\n%s", want, rec)
	}
	v := await(t, root, l.Attempt, Finished)
	if v.Result.Record == nil || v.Result.Record.Status != "review" {
		t.Fatalf("result %+v", v.Result.Record)
	}
	next := v.Result.Record.Candidate
	for _, ancestor := range []string{candidate, tip} {
		git(t, root, "merge-base", "--is-ancestor", ancestor, next)
	}
	if n := starts(t, mark); n != 1 {
		t.Fatalf("the provider started %d times", n)
	}
	ms, err := versions.PredictContext(t.Context(), root, "refs/heads/main", []string{next})
	if err != nil || ms[0].Outcome != "fast-forward" {
		t.Fatalf("the new candidate %+v %v", ms, err)
	}
	// What the owner judges: the merge, what it merged, and the file it resolved.
	c, err := versions.ChangesContext(t.Context(), root, "main", next, next, "grove/G-001-first.md")
	if err != nil {
		t.Fatal(err)
	}
	if r := c.Resolution; r == nil || r.Merge != next || r.Target != tip || r.Previous != candidate || strings.Join(r.Files, ",") != "shared.txt" {
		t.Fatalf("resolution %+v", c.Resolution)
	}
	// In review again, with an attempt no longer running: a second resolve
	// finds no conflict, not a duplicate to start.
	if _, _, err := resolve(root, nil, now.Add(time.Minute)); err == nil || !strings.Contains(err.Error(), "no conflict") {
		t.Fatalf("again: %v", err)
	}
}

func TestResolveNeedsAChoice(t *testing.T) {
	skipShort(t)
	root, _, candidate, _ := conflicted(t)
	// The agent stops at a choice the record does not settle: it commits
	// nothing, and the clean exit marks nothing.
	fake(t, initLine+"\n"+resultLine("success", false))
	l, _, err := resolve(root, nil, now)
	if err != nil {
		t.Fatal(err)
	}
	v := await(t, root, l.Attempt, Finished)
	if r := v.Result.Record; r == nil || r.Status != "active" || r.Candidate != candidate || v.Result.Head != l.Base {
		t.Fatalf("result %+v %+v", v.Result, r)
	}
	if _, _, err := resolve(root, nil, now.Add(time.Minute)); err == nil || !strings.Contains(err.Error(), "no branch holds G-001 in review") {
		t.Fatalf("active work is not resolved again: %v", err)
	}
}

func TestResolveWhileTheTargetMoves(t *testing.T) {
	skipShort(t)
	root, _, candidate, tip := conflicted(t)
	_, release := fake(t, initLine+"\nwhile [ ! -e \"$RELEASE\" ]; do sleep 0.05; done\n"+resolving+resultLine("success", false))
	l, _, err := resolve(root, nil, now)
	if err != nil {
		t.Fatal(err)
	}
	await(t, root, l.Attempt, Running)
	write(t, root, "other.txt", "later\n")
	git(t, root, "add", "other.txt")
	git(t, root, "commit", "-qm", "main moves again")
	later := git(t, root, "rev-parse", "HEAD")
	os.WriteFile(release, nil, 0o644)
	v := await(t, root, l.Attempt, Finished)
	next := v.Result.Record.Candidate
	if v.Result.Record.Status != "review" || git(t, root, "rev-parse", next+"^2") != tip {
		t.Fatalf("the attempt merges the commit its feedback names, %s: %+v", tip[:7], v.Result.Record)
	}
	git(t, root, "merge-base", "--is-ancestor", candidate, next)
	ms, err := versions.PredictContext(t.Context(), root, "refs/heads/main", []string{next})
	if err != nil || ms[0].Target != later || ms[0].Outcome != "clean" {
		t.Fatalf("the prediction names the moved target: %+v %v", ms, err)
	}
}

func TestResolveStopped(t *testing.T) {
	skipShort(t)
	root, _, candidate, _ := conflicted(t)
	fake(t, waiting+resultLine("success", false))
	l, _, err := resolve(root, nil, now)
	if err != nil {
		t.Fatal(err)
	}
	await(t, root, l.Attempt, Running)
	if err := Stop(root, l.Attempt, func(string) {}); err != nil {
		t.Fatal(err)
	}
	v := await(t, root, l.Attempt, Finished)
	if r := v.Result.Record; !v.Result.Stopped || r == nil || r.Status != "active" || r.Candidate != candidate {
		t.Fatalf("result %+v %+v", v.Result, r)
	}
}
