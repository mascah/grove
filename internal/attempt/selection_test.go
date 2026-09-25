package attempt

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"
)

// member writes work id needing deps, in status.
func member(t *testing.T, root, id, status string, deps ...string) {
	t.Helper()
	needs := ""
	if deps != nil {
		needs = fmt.Sprintf("depends_on: [\"%s\"]\n", strings.Join(deps, "\", \""))
	}
	write(t, root, "grove/"+id+".md", fmt.Sprintf("---\nid: %q\ntype: work\ntitle: Work %s\nstatus: %s\ncreated: \"2026-09-22T10:00:00Z\"\nupdated: \"2026-09-22T10:00:00Z\"\n%s---\n\n## Outcome\n\nA thing.\n", id, id, status, needs))
}

// chain is the fixture plus G-003 needing G-001, G-004 needing G-003, and
// G-005 needing G-001: a chain and a branch from G-001.
func chain(t *testing.T) string {
	t.Helper()
	root := fixture(t)
	member(t, root, "G-003", "proposed", "G-001")
	member(t, root, "G-004", "proposed", "G-003")
	member(t, root, "G-005", "active", "G-001")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "chain")
	return root
}

func preview(t *testing.T, root string, ids ...string) *Launch {
	t.Helper()
	l, err := Preview(Request{Root: root, IDs: ids, BudgetUSD: "3", PermissionMode: "auto"})
	if err != nil {
		t.Fatal(err)
	}
	return l
}

func waits(l *Launch) map[string]string {
	out := map[string]string{}
	for _, m := range l.Members() {
		out[m.ID] = m.Wait
	}
	return out
}

func TestSelectionChainAndBranch(t *testing.T) {
	t.Parallel()
	root := chain(t)
	head := git(t, root, "rev-parse", "HEAD")
	l := preview(t, root, "G-004", "G-005", "G-003", "G-001")
	s := l.Selection
	if !slices.Equal(s.Order, []string{"G-001", "G-005", "G-003", "G-004"}) || !slices.Equal(s.Selected, []string{"G-004", "G-005", "G-003", "G-001"}) {
		t.Fatalf("order %v selected %v", s.Order, s.Selected)
	}
	if l.Work != "G-004" || l.Branch != "worktree-G-004-G-005-G-003-G-001" || l.Base != head || l.WorktreeReused {
		t.Fatalf("%+v", l)
	}
	for id, w := range waits(l) {
		if w != "" {
			t.Fatalf("%s waits: %s", id, w)
		}
	}
	if !strings.HasPrefix(s.Digest, "sha256:") || preview(t, root, "G-004", "G-005", "G-003", "G-001").Selection.Digest != s.Digest {
		t.Fatalf("digest %q is not stable", s.Digest)
	}
	if preview(t, root, "G-001", "G-003", "G-004", "G-005").Selection.Digest == s.Digest {
		t.Fatal("the IDs as given are part of the assignment")
	}
	// A preview writes nothing: no branch, worktree or attempt.
	if out := git(t, root, "branch", "--list", "worktree-*"); out != "" {
		t.Fatalf("a preview made a branch: %q", out)
	}
	if views, err := List(root, ""); err != nil || len(views) != 0 {
		t.Fatalf("%v %v", views, err)
	}
	text := strings.Join(Explain(l, func(v string) string { return v }), "\n")
	for _, want := range []string{"Selected: G-004 G-005 G-003 G-001; order G-001, G-005, G-003, G-004", "Member G-005: active at grove/G-005.md", "can start", "Review boundary: one review at the end", "Continuation: members run in order"} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in\n%s", want, text)
		}
	}

	// A question on the chain's middle stops it and what needs it; the
	// branch beside it still starts.
	write(t, root, "grove/G-002-q.md", strings.Replace(question, `blocks: ["G-001"]`, `blocks: ["G-003"]`, 1))
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "question")
	w := waits(preview(t, root, "G-001", "G-003", "G-004", "G-005"))
	if w["G-001"] != "" || w["G-005"] != "" || !strings.Contains(w["G-003"], "blocked by open question G-002") || w["G-004"] != "needs G-003, which waits" {
		t.Fatalf("%q", w)
	}
	// With nothing able to start, the launch is refused before any write.
	_, err := Start(Request{Root: root, IDs: []string{"G-003", "G-004"}, BudgetUSD: "1", PermissionMode: "auto"}, now, func(string) {})
	if err == nil || !strings.Contains(err.Error(), "nothing in the selection can start: G-003 blocked by open question G-002 (Which colour?); needs G-001") || !strings.Contains(err.Error(), "G-004 needs G-001, which is proposed and not selected; needs G-003, which waits") {
		t.Fatal(err)
	}
}

func TestSelectionThroughUnselectedWork(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	member(t, root, "G-003", "proposed", "G-001") // unselected, between the two
	member(t, root, "G-004", "proposed", "G-003")
	// G-006 is done with a candidate on a branch the base does not contain.
	git(t, root, "checkout", "-q", "-b", "elsewhere")
	git(t, root, "commit", "-q", "--allow-empty", "-m", "delivered elsewhere")
	elsewhere := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-q", "main")
	member(t, root, "G-006", "done", "G-001")
	write(t, root, "grove/G-006.md", strings.Replace(readFile(t, root, "grove/G-006.md"), "status: done\n", "status: done\ncandidate: \""+elsewhere+"\"\n", 1))
	member(t, root, "G-007", "proposed", "G-006")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "path")
	l := preview(t, root, "G-004", "G-001", "G-007")
	s := l.Selection
	if !slices.Equal(s.Order, []string{"G-001", "G-004", "G-007"}) {
		t.Fatalf("order %v", s.Order)
	}
	w := waits(l)
	if w["G-001"] != "" || w["G-004"] != "needs G-003, which is proposed and not selected" || w["G-007"] != "needs G-006, whose candidate "+elsewhere[:7]+" the base lacks" {
		t.Fatalf("%q", w)
	}
	var outside []string
	for _, o := range s.Outside {
		outside = append(outside, o.ID+":"+strings.Join(o.NeededBy, ","))
	}
	if !slices.Equal(outside, []string{"G-003:G-004", "G-006:G-007"}) || slices.ContainsFunc(s.Members, func(m Member) bool { return m.ID == "G-003" }) {
		t.Fatalf("outside %v members %+v", outside, s.Members)
	}
}

func readFile(t *testing.T, root, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestSelectionRefusals(t *testing.T) {
	t.Parallel()
	root := chain(t)
	l := preview(t, root, "G-001", "G-003")
	// A changed member changes the digest, and the launch names what it
	// would run now.
	member(t, root, "G-003", "active", "G-001")
	git(t, root, "commit", "-qam", "moved on")
	_, err := Start(Request{Root: root, IDs: []string{"G-001", "G-003"}, BudgetUSD: "3", PermissionMode: "auto", Digest: l.Selection.Digest}, now, func(string) {})
	if err == nil || !strings.Contains(err.Error(), "the assignment changed since its preview") || !strings.Contains(err.Error(), "Member G-003: active") {
		t.Fatal(err)
	}
	if _, err := Preview(Request{Root: root, IDs: []string{"G-001", "G-001"}, BudgetUSD: "1", PermissionMode: "auto"}); err == nil || !strings.Contains(err.Error(), "G-001 is selected more than once") {
		t.Fatal(err)
	}
	if _, err := Preview(Request{Root: root, IDs: []string{"G-001", "G-003"}, BudgetUSD: "1", PermissionMode: "auto", Expect: "sha256:x"}); err == nil || !strings.Contains(err.Error(), "a selection is checked by its digest") {
		t.Fatal(err)
	}

	// An attempt running over G-003 owns it: an overlapping selection is
	// refused, a disjoint one is not.
	dir, err := Dir(root)
	if err != nil {
		t.Fatal(err)
	}
	adir := filepath.Join(dir, "G-003.20260922T100000Z")
	if err := os.MkdirAll(adir, 0o755); err != nil {
		t.Fatal(err)
	}
	running := Launch{Attempt: "G-003.20260922T100000Z", Work: "G-003", Started: now.Add(-time.Hour), Selection: &Selection{Members: []Member{{ID: "G-003"}, {ID: "G-004"}}}}
	if err := writeJSON(filepath.Join(adir, "attempt.json"), running); err != nil {
		t.Fatal(err)
	}
	lock, err := os.OpenFile(filepath.Join(adir, "owner.lock"), os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Close()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		t.Fatal(err)
	}
	_, err = Start(Request{Root: root, IDs: []string{"G-001", "G-004"}, BudgetUSD: "1", PermissionMode: "auto"}, now, func(string) {})
	if err == nil || !strings.Contains(err.Error(), "attempt G-003.20260922T100000Z of G-004 is running since") {
		t.Fatal(err)
	}
	if views, _ := List(root, "G-004"); len(views) != 1 || views[0].Launch.Attempt != running.Attempt {
		t.Fatalf("attempts of a member: %+v", views)
	}

	// A reused branch where one member is already in review: judge it first.
	wt := filepath.Join(root, ".claude", "worktrees", "worktree-G-001-G-005")
	git(t, root, "worktree", "add", "-q", "-b", "worktree-G-001-G-005", wt)
	write(t, wt, "grove/G-005.md", strings.Replace(readFile(t, wt, "grove/G-005.md"), "status: active\n", "status: review\ncandidate: \""+git(t, wt, "rev-parse", "HEAD")+"\"\n", 1))
	git(t, wt, "commit", "-qam", "review")
	_, err = Preview(Request{Root: root, IDs: []string{"G-001", "G-005"}, BudgetUSD: "1", PermissionMode: "auto"})
	if err == nil || !strings.Contains(err.Error(), "G-005 is review on worktree-G-001-G-005 at "+wt+"; judge that candidate") {
		t.Fatal(err)
	}
}

// TestMemberResults reconciles a selection's worktree as the owner and Stop
// do, and reads an attempt from before selections as a selection of one.
func TestMemberResults(t *testing.T) {
	t.Parallel()
	root := chain(t)
	head := git(t, root, "rev-parse", "HEAD")
	write(t, root, "grove/G-001-first.md", strings.Replace(fmt.Sprintf(work, "review"), "---\n\n## Outcome", "candidate: \""+head+"\"\n---\n\n## Outcome", 1))
	member(t, root, "G-003", "active", "G-001")
	write(t, root, "grove/G-002-q.md", strings.Replace(question, `blocks: ["G-001"]`, `blocks: ["G-005"]`, 1))
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "progress")
	write(t, root, "grove/G-003.md", readFile(t, root, "grove/G-003.md")+"\n## Next\n\nHalf done.\n") // uncommitted
	l := &Launch{Attempt: "G-001.20260922T183000Z", Work: "G-001", Worktree: root, Selection: &Selection{
		Selected: []string{"G-001", "G-003", "G-004", "G-005"}, Order: []string{"G-001", "G-003", "G-004", "G-005"},
		Members: []Member{{ID: "G-001"}, {ID: "G-003"}, {ID: "G-004", Wait: "needs G-003, which waits"}, {ID: "G-005"}},
	}}
	dir := t.TempDir()
	res := Result{Stopped: true}
	reconcile(dir, l, &res, t.Logf)
	if len(res.Members) != 4 || res.Record == nil || res.Record.Status != "review" || res.RecordUncommitted {
		t.Fatalf("%+v", res)
	}
	v := &View{Dir: dir, Launch: *l, Status: Finished, Result: &res}
	text := strings.Join(Facts(v, func(s string) string { return s }), "\n")
	for _, want := range []string{
		"Member G-001 on the branch: awaiting judgment: review, candidate " + head[:12],
		"Member G-003 on the branch: active; its Next holds the checkpoint, uncommitted",
		"Member G-004 on the branch: not started: needs G-003, which waits",
		"Member G-005 on the branch: active, waiting on question G-002 (Which colour?)",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in\n%s", want, text)
		}
	}

	// An attempt file from before selections: one work, found by its ID.
	dir, err := Dir(root)
	if err != nil {
		t.Fatal(err)
	}
	old := filepath.Join(dir, "G-001.20260101T000000Z")
	if err := os.MkdirAll(old, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(old, "attempt.json"), []byte(`{"attempt":"G-001.20260101T000000Z","work":"G-001","record_path":"grove/G-001-first.md","record_revision":"sha256:old"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(old, "result.json"), []byte(`{"exit_code":0,"record":{"status":"active","revision":"sha256:x"}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	views, err := List(root, "G-001")
	if err != nil || len(views) != 1 {
		t.Fatalf("%+v %v", views, err)
	}
	if m := views[0].Launch.Members(); len(m) != 1 || m[0].Path != "grove/G-001-first.md" {
		t.Fatalf("%+v", m)
	}
	text = strings.Join(Facts(&views[0], func(s string) string { return s }), "\n")
	if !strings.Contains(text, "Work: G-001 at grove/G-001-first.md, record sha256:old") || !strings.Contains(text, "Record on the branch: G-001 active") {
		t.Fatal(text)
	}
}

// TestRunSelection launches a selection through the real owner: one
// process, the IDs as given in its prompt, one branch named for them.
func TestRunSelection(t *testing.T) {
	skipShort(t)
	root := chain(t)
	fake(t, `printf '%s\n' "$@" > argv.txt
`+initLine+`
`+resultLine("success", false))
	l, err := Start(Request{Root: root, IDs: []string{"G-003", "G-001"}, BudgetUSD: "2", PermissionMode: "auto"}, now, func(string) {})
	if err != nil {
		t.Fatal(err)
	}
	v := await(t, root, l.Attempt, Finished)
	if l.Attempt != "G-003.20260922T183000Z" || l.Branch != "worktree-G-003-G-001" || !slices.Contains(l.Command, "/grove-work G-003 G-001 --interaction headless") {
		t.Fatalf("%+v", l)
	}
	if len(v.Result.Members) != 2 || v.Result.Members[0].ID != "G-001" || v.Result.Members[1].Record.Status != "proposed" {
		t.Fatalf("%+v", v.Result.Members)
	}
	if !strings.Contains(strings.Join(Facts(v, func(s string) string { return s }), "\n"), "Member G-003 on the branch: not started") {
		t.Fatal("no member line")
	}
}

// A wait in the launching checkout stands on a reused branch too, since the
// agent in the branch's worktree would never see it; bounded at plans, a
// prerequisite outside the selection stops nothing; an unreadable attempt
// fails only listings that could include it.
func TestSelectionWaitsBothPlaces(t *testing.T) {
	t.Parallel()
	root := chain(t)
	wt := filepath.Join(root, ".claude", "worktrees", "worktree-G-001")
	git(t, root, "worktree", "add", "-q", "-b", "worktree-G-001", wt)
	write(t, root, "grove/G-002-q.md", question) // blocks G-001, on main only
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "question on main")
	_, err := Preview(Request{Root: root, IDs: []string{"G-001"}, BudgetUSD: "1", PermissionMode: "auto"})
	if err == nil || !strings.Contains(err.Error(), "G-001 blocked by open question G-002 (Which colour?)") {
		t.Fatal(err)
	}
	// Bounded at its plan, G-003 is not stopped by G-001 being unfinished;
	// through to the handoff it is.
	if _, err := Preview(Request{Root: root, IDs: []string{"G-003"}, BudgetUSD: "1", PermissionMode: "auto"}); err == nil || !strings.Contains(err.Error(), "G-003 needs G-001, which is proposed and not selected") {
		t.Fatal(err)
	}
	if l, err := Preview(Request{Root: root, IDs: []string{"G-003"}, BudgetUSD: "1", PermissionMode: "auto", Until: "plan"}); err != nil || l.Selection.Members[0].Wait != "" || l.Selection.Outside[0].ID != "G-001" {
		t.Fatalf("%+v %v", l, err)
	}
	dir, err := Dir(root)
	if err != nil {
		t.Fatal(err)
	}
	write(t, dir, "G-004.20260101T000000Z/attempt.json", "{not json")
	if _, err := List(root, "G-003"); err != nil {
		t.Fatalf("another work's broken attempt: %v", err)
	}
	if _, err := List(root, "G-004"); err == nil {
		t.Fatal("its own work's listing must say it is broken")
	}
}

// After feedback reopens a group, relaunching part of it is refused: the
// next candidate is the group's.
func TestSelectionReopenedGroupRunsTogether(t *testing.T) {
	t.Parallel()
	root := chain(t)
	wt := filepath.Join(root, ".claude", "worktrees", "worktree-G-001-G-003")
	git(t, root, "worktree", "add", "-q", "-b", "worktree-G-001-G-003", wt)
	shared := git(t, wt, "rev-parse", "HEAD")
	write(t, wt, "grove/G-001-first.md", strings.Replace(fmt.Sprintf(work, "active"), "---\n\n## Outcome", "candidate: \""+shared+"\"\n---\n\n## Outcome", 1))
	write(t, wt, "grove/G-003.md", strings.Replace(readFile(t, wt, "grove/G-003.md"), "status: proposed\n", "status: active\ncandidate: \""+shared+"\"\n", 1))
	git(t, wt, "commit", "-qam", "reopened by feedback")
	_, err := Preview(Request{Root: root, IDs: []string{"G-001"}, Branch: "worktree-G-001-G-003", BudgetUSD: "1", PermissionMode: "auto"})
	if err == nil || !strings.Contains(err.Error(), "G-003 shares candidate "+shared[:7]+" with G-001 on worktree-G-001-G-003 and was reopened with it; select them together (grove run G-001 G-003)") {
		t.Fatal(err)
	}
	l, err := Preview(Request{Root: root, IDs: []string{"G-001", "G-003"}, BudgetUSD: "1", PermissionMode: "auto"})
	if err != nil || !l.WorktreeReused {
		t.Fatalf("%+v %v", l, err)
	}
	seen := map[string]bool{}
	for _, n := range l.Selection.Notes {
		if seen[n] {
			t.Fatalf("note repeated: %q", n)
		}
		seen[n] = true
	}
}
