package tui

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"

	"github.com/mascah/grove/internal/versions"
)

// reviewFixture: W-001 is in review on branch feature, held by the checkout
// feat, with a review record; main still has it proposed, superseded by
// feature's change. The target is main, whose checkout is ".".
func reviewFixture(fx fixture, approved bool) *fake {
	var vs []versions.Version
	for _, s := range []*versions.Source{fx.cMain, fx.main} {
		v := version(s, "W-001", "Inspect records", "proposed")
		v.Older = "branch feature changed it since"
		vs = append(vs, v)
	}
	for _, s := range []*versions.Source{fx.cFeat, fx.feat} {
		w := version(s, "W-001", "Inspect records", "review")
		w.Record.Candidate = "abcdef1"
		if approved {
			w.Record.Approved = "abcdef1"
		}
		w.Record.Source = []byte("---\nid: W-001\n---\n\n## Outcome\n\nA board.\n" + strings.Repeat("\nfiller\n", 30) + "\n## Evidence\n\nIt works.\n\n## Next\n\nJudge it.\n")
		review := version(s, "W-006", "W-001 review", "current")
		review.Record.Type, review.Record.Work, review.Record.Examined = "review", []string{"W-001"}, "abcdef1"
		vs = append(vs, w, review)
	}
	res := result(fx.main, fx.sources(), vs...)
	res.Target = "main"
	return &fake{res: res, actions: true,
		history: func(context.Context, string, string) ([]versions.Commit, error) { return nil, nil },
		changes: func(target, candidate, tip, path string) (*versions.Changes, error) {
			return &versions.Changes{Base: "base000", Files: []versions.Change{{Path: "internal/x.go", Added: 12, Removed: 3}, {Path: "grove/work/W-001.md", Added: 5, Removed: 1}, {Path: "bin.dat", Added: -1, Removed: -1}}}, nil
		},
		diff: func(from, to, path string) (string, error) {
			return "diff --git a/" + path + " b/" + path + "\n@@ -1,2 +1,3 @@\n-old\n+new \x1b]0;evil\a\n+\tmore\n", nil
		},
	}
}

// openReview opens W-001's detail from the Review column and delivers the
// history and changes reads.
func openReview(t *testing.T, f *fake, w, h int) *Model {
	t.Helper()
	m := open(t, f, w, h)
	press(m, "right", "right")
	if m.cardID != "W-001" || m.col != 2 {
		t.Fatalf("W-001 should be the Review column's card: col %d card %s", m.col, m.cardID)
	}
	deliverAll(m, press(m, "enter"))
	return m
}

// The detail of a candidate in review leads with its standing and where the
// actions run, starts at the Evidence, lists the changed files, and shows a
// file's diff escaped; a narrow terminal keeps every row its width.
func TestReviewDetailShowsStandingChangesAndDiffs(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := reviewFixture(fx, false)
	m := openReview(t, f, 120, 36)
	s := plain(m)
	for _, want := range []string{
		"W-001 · review", "candidate abcdef1 · not on main",
		"Review: candidate abcdef1 · not yet approved · only the record changed since it · not on main",
		"a approve and f feedback run on branch feature in /repo/feat · i integrate runs into main in /repo",
		"## Evidence", "It works.",
		"review     W-006  W-001 review  current", "examined abcdef1 = candidate",
		"Changes against main from base000", "internal/x.go  +12 −3", "grove/work/W-001.md  +5 −1", "bin.dat  binary",
		"a approve  f feedback  i integrate",
	} {
		if !strings.Contains(s, want) {
			t.Fatalf("review detail lacks %q:\n%s", want, s)
		}
	}
	if strings.Contains(s, "## Outcome") {
		t.Fatalf("the content should start at the Evidence:\n%s", s)
	}
	if strings.Join(f.reads, ";") != "changes main abcdef1 a grove/work/W-001.md" {
		t.Fatalf("reads: %v", f.reads)
	}
	// Tab reaches the changes after the linked records; Enter shows the diff.
	press(m, "tab", "down")
	s = plain(m)
	if !strings.Contains(s, "> internal/x.go  +12 −3") {
		t.Fatalf("Tab and ↓ should reach the changed file:\n%s", s)
	}
	deliverAll(m, press(m, "enter"))
	s = plain(m)
	for _, want := range []string{"Diff of internal/x.go (Esc returns to the content)", "-old", `+new \x1b]0;evil\a`, "+    more"} {
		if !strings.Contains(s, want) {
			t.Fatalf("diff view lacks %q:\n%s", want, s)
		}
	}
	if strings.Contains(m.render(), "\x1b]0;") {
		t.Fatal("the diff's control sequence reached the screen")
	}
	if strings.Join(f.reads, ";") != "changes main abcdef1 a grove/work/W-001.md;diff base000 abcdef1 internal/x.go" {
		t.Fatalf("reads: %v", f.reads)
	}
	press(m, "esc")
	if s = plain(m); m.diff != "" || !strings.Contains(s, "▶ Content") && !strings.Contains(s, "  Content") || strings.Contains(s, "Diff of") {
		t.Fatalf("Esc returns to the content:\n%s", s)
	}
	press(m, "esc")
	if m.screen != boardScreen {
		t.Fatal("Esc then returns to the board")
	}
	// Narrow: every row is the terminal's width with the review block.
	for _, size := range [][2]int{{99, 30}, {80, 24}, {40, 12}} {
		m := openReview(t, reviewFixture(fx, true), size[0], size[1])
		for _, row := range strings.Split(m.render(), "\n") {
			if got := ansi.StringWidth(row); got != size[0] {
				t.Fatalf("%v: row is %d cells: %q", size, got, ansi.Strip(row))
			}
		}
		if s := plain(m); size[1] >= 16 && !strings.Contains(s, "Review: candidate abcdef1 · approved") {
			t.Fatalf("%v: review block missing:\n%s", size, s)
		}
	}
}

// Esc from the detail cancels a changes or diff read as it cancels history.
func TestReviewReadsStopWithTheDetail(t *testing.T) {
	t.Parallel()
	f := reviewFixture(newFixture(), false)
	m := open(t, f, 120, 36)
	press(m, "right", "right", "enter")
	deliverAll(m, press(m, "esc")) // the history read is cancelled: the detail's own
	if m.pending != "" && m.pending != "inspect" {
		t.Fatalf("pending %q", m.pending)
	}
	m = openReview(t, f, 120, 36)
	press(m, "tab", "down")
	if cmd := press(m, "enter"); cmd == nil || m.pending != "diff" {
		t.Fatalf("Enter on a change starts the diff read: pending %q", m.pending)
	}
	press(m, "esc")
	if m.pending != "" || m.cancel != nil {
		t.Fatalf("Esc must cancel the diff read: pending %q", m.pending)
	}
}

// a and f open a prompt; Enter runs the action in the branch's checkout,
// shows its facts on the result screen, and re-reads the board; keys wait
// while it runs; Esc cancels a prompt; an empty verdict is refused.
func TestReviewApproveAndFeedbackFromTheBoard(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := reviewFixture(fx, false)
	m := openReview(t, f, 120, 36)
	press(m, "a")
	if s := plain(m); m.prompt == nil || !strings.Contains(s, "Approve W-001 on branch feature · verdict (Enter records it, Esc cancels): ▏") {
		t.Fatalf("a should open the verdict prompt:\n%s", s)
	}
	press(m, "enter")
	if s := plain(m); m.prompt == nil || !strings.Contains(s, "type the verdict first, or Esc") {
		t.Fatalf("an empty verdict is refused:\n%s", s)
	}
	typeText(m, "Good ✓")
	press(m, "backspace")
	press(m, "esc")
	if s := plain(m); m.prompt != nil || !strings.Contains(s, "cancelled; nothing was written") || len(f.acts) != 0 {
		t.Fatalf("Esc cancels the prompt:\n%s", s)
	}
	press(m, "a")
	typeText(m, "Ship it")
	cmd := press(m, "enter")
	if cmd == nil || m.pending != "act" || m.prompt != nil || !strings.Contains(plain(m), "Approving…") {
		t.Fatalf("Enter should start the action: pending %q\n%s", m.pending, plain(m))
	}
	press(m, "q", "esc", "a")
	if m.done || !strings.Contains(plain(m), "Approving… Keys wait for it.") {
		t.Fatalf("keys wait while an action runs:\n%s", plain(m))
	}
	next := deliver(m, cmd)
	if m.screen != resultScreen || next == nil || m.pending != "inspect" {
		t.Fatalf("the outcome shows and the board is re-read: screen %d pending %q", m.screen, m.pending)
	}
	s := plain(m)
	for _, want := range []string{"Approved W-001", "approved: W-001 in /repo/feat", "Re-reading the board…"} {
		if !strings.Contains(s, want) {
			t.Fatalf("result lacks %q:\n%s", want, s)
		}
	}
	press(m, "esc") // waits for the re-read, so the record shown next is current
	if m.screen != resultScreen || m.pending != "inspect" || !strings.Contains(plain(m), "Esc waits for the re-read") {
		t.Fatal("Esc must wait for the re-read, and the hints say so")
	}
	deliverAll(m, next)
	if s := plain(m); m.screen != resultScreen || !strings.Contains(s, "The board has been re-read. Esc returns to the record.") {
		t.Fatalf("after the re-read:\n%s", s)
	}
	if f.inspects != 2 || strings.Join(f.acts, ";") != "approve /repo/feat W-001 Ship it" {
		t.Fatalf("inspects %d acts %v", f.inspects, f.acts)
	}
	press(m, "esc")
	if m.screen != detailScreen || m.openID() != "W-001" {
		t.Fatalf("Esc returns to the record: screen %d", m.screen)
	}
	// Feedback, with the continuation in its facts.
	press(m, "f")
	if s := plain(m); !strings.Contains(s, "Feedback on W-001, returning it to active on branch feature") {
		t.Fatalf("f should open the feedback prompt:\n%s", s)
	}
	typeText(m, "Needs the empty case")
	deliverAll(m, press(m, "enter"))
	if s := plain(m); m.screen != resultScreen || !strings.Contains(s, "Feedback recorded on W-001") || !strings.Contains(s, "next: /grove-work W-001") {
		t.Fatalf("feedback outcome:\n%s", s)
	}
	if f.acts[1] != "feedback /repo/feat W-001 Needs the empty case" {
		t.Fatalf("acts %v", f.acts)
	}
	// A failed action shows why, and the facts it returned.
	press(m, "esc", "a")
	f.fail = errors.New("this checkout does not hold candidate abcdef1")
	typeText(m, "x")
	deliverAll(m, press(m, "enter"))
	if s := plain(m); !strings.Contains(s, "NOT DONE: this checkout does not hold candidate abcdef1") {
		t.Fatalf("a failure is shown:\n%s", s)
	}
}

// i needs an approval and the target's checkout, confirms the merge and then
// the cleanup, and runs in the target's checkout.
func TestReviewIntegrateFromTheBoard(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	m := openReview(t, reviewFixture(fx, false), 120, 36)
	press(m, "i")
	if s := plain(m); m.prompt != nil || !strings.Contains(s, "approve candidate abcdef1 first (a)") {
		t.Fatalf("i before approval:\n%s", s)
	}
	f := reviewFixture(fx, true)
	m = openReview(t, f, 120, 36)
	press(m, "a")
	if s := plain(m); m.prompt != nil || !strings.Contains(s, "candidate abcdef1 is already approved; i integrates it") {
		t.Fatalf("a after approval:\n%s", s)
	}
	press(m, "i")
	if s := plain(m); m.prompt == nil || !strings.Contains(s, "Merge branch feature into main and mark W-001 done? y/n   (runs in /repo)") {
		t.Fatalf("i should confirm the merge:\n%s", s)
	}
	press(m, "n")
	if s := plain(m); m.prompt != nil || !strings.Contains(s, "cancelled; nothing was merged") || len(f.acts) != 0 {
		t.Fatalf("n cancels:\n%s", s)
	}
	press(m, "i", "y")
	if s := plain(m); m.prompt == nil || !strings.Contains(s, "Also delete branch feature and remove its worktree? y/n   (/repo/feat)") {
		t.Fatalf("y should ask about cleanup:\n%s", s)
	}
	deliverAll(m, press(m, "n"))
	if s := plain(m); m.screen != resultScreen || !strings.Contains(s, "Integration of W-001") || !strings.Contains(s, "merge: fast-forward") {
		t.Fatalf("integration outcome:\n%s", s)
	}
	press(m, "esc", "i", "y")
	deliverAll(m, press(m, "y"))
	if strings.Join(f.acts, ";") != "integrate /repo W-001 cleanup=false;integrate /repo W-001 cleanup=true" {
		t.Fatalf("acts %v", f.acts)
	}
	// Other keys during a y/n prompt do nothing.
	press(m, "esc", "i", "x", "q")
	if m.prompt == nil || m.done {
		t.Fatal("a y/n prompt ignores other keys")
	}
	press(m, "esc")
	if m.prompt != nil {
		t.Fatal("Esc cancels a y/n prompt")
	}
}

// The actions need a clean checkout of the branch, and the target's
// checkout: the review block says what is missing, and the key says so too.
func TestReviewActionsNeedTheRightCheckout(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := reviewFixture(fx, true)
	for i := range f.res.Groups[0].Versions { // W-001's checkout copy is modified
		if v := &f.res.Groups[0].Versions[i]; v.Source == fx.feat {
			v.Change = "modified"
		}
	}
	m := openReview(t, f, 120, 36)
	if s := plain(m); !strings.Contains(s, "a approve and f feedback: the record has uncommitted changes in checkout feat; commit or discard them before judging") {
		t.Fatalf("review block:\n%s", s)
	}
	press(m, "f")
	if s := plain(m); m.prompt != nil || !strings.Contains(s, "uncommitted changes in checkout feat") {
		t.Fatalf("f on a modified record:\n%s", s)
	}
	// No checkout of the branch, and no checkout of the target.
	f = reviewFixture(fx, true)
	f.res.Sources = []*versions.Source{fx.cFeat, fx.cMain}
	f.res.Groups[0].Versions = f.res.Groups[0].Versions[:0:0]
	for _, s := range []*versions.Source{fx.cMain, fx.cFeat} {
		v := version(s, "W-001", "Inspect records", map[bool]string{true: "proposed", false: "review"}[s == fx.cMain])
		if s == fx.cMain {
			v.Older = "branch feature changed it since"
		} else {
			v.Record.Candidate, v.Record.Approved = "abcdef1", "abcdef1"
		}
		f.res.Groups[0].Versions = append(f.res.Groups[0].Versions, v)
	}
	f.res.Groups = f.res.Groups[:1]
	m = openReview(t, f, 120, 36)
	if s := plain(m); !strings.Contains(s, "a approve and f feedback: no checkout is on branch feature; git worktree add one, or run grove approve there") || !strings.Contains(s, "integrate: no checkout is on the target main; i needs one") {
		t.Fatalf("review block without checkouts:\n%s", s)
	}
	press(m, "a")
	if s := plain(m); m.prompt != nil || !strings.Contains(s, "no checkout is on branch feature") {
		t.Fatalf("a without a checkout:\n%s", s)
	}
	press(m, "i")
	if s := plain(m); m.prompt != nil || !strings.Contains(s, "no checkout is on the target main") {
		t.Fatalf("i without the target's checkout:\n%s", s)
	}
	if len(f.acts) != 0 {
		t.Fatalf("nothing should have run: %v", f.acts)
	}
	// Not in review: the keys explain.
	m = open(t, linkedFixture(newFixture()), 120, 36)
	press(m, "right")
	deliverAll(m, press(m, "enter"))
	press(m, "a")
	if s := plain(m); m.prompt != nil || !strings.Contains(s, "W-001 is not in review: nothing to approve") || strings.Contains(s, "Review: candidate") {
		t.Fatalf("a on active work:\n%s", s)
	}
}

// Beside each changed file the review lists the other records that link it
// or name it in a code span, or says none does; it reads nothing more. A
// path under the project's prefix, which Git gives with its slash, is
// matched as a project path, one outside it only by a code span, and a
// rename by either side.
func TestReviewListsRecordsDescribingEachFile(t *testing.T) {
	t.Parallel()
	for _, prefix := range []string{"", "proj/"} {
		fx := newFixture()
		f := reviewFixture(fx, false)
		for _, s := range []*versions.Source{fx.cMain, fx.main, fx.cFeat, fx.feat} {
			f.res.Groups = append(f.res.Groups,
				versions.Group{ID: "Q-002", Versions: []versions.Version{withBody(version(s, "Q-002", "Where?", "open"), "In `x.go`, and [bin](../../bin.dat).\n")}},
				versions.Group{ID: "D-002", Versions: []versions.Version{withBody(version(s, "D-002", "Keep x", "accepted"), "[x](../../internal/x.go), [old](../../old/y.go), `Cargo.toml`\n")}})
		}
		for i := range f.res.Groups {
			for j := range f.res.Groups[i].Versions {
				if v := &f.res.Groups[i].Versions[j]; v.Record.ID == "W-001" {
					*v = withBody(*v, "Changes [x](../../internal/x.go).\n")
					v.Record.Candidate = "abcdef1"
				}
			}
		}
		f.res.Prefix = prefix
		changes := f.changes
		f.changes = func(target, candidate, tip, path string) (*versions.Changes, error) {
			c, err := changes(target, candidate, tip, path)
			for i := range c.Files {
				c.Files[i].Path = prefix + c.Files[i].Path
			}
			c.Files = append(c.Files, versions.Change{Path: prefix + "old/y.go → " + prefix + "new/y.go", Added: 1}, versions.Change{Path: "Cargo.toml", Added: 1})
			return c, err
		}
		m := openReview(t, f, 160, 50)
		s := plain(m)
		want := []string{"internal/x.go  +12 −3", "described by Q-002 code span, D-002 link", "grove/work/W-001.md  +5 −1", "no record names it",
			"bin.dat  binary", "described by Q-002 link", "new/y.go  +1 −0", "described by D-002 link", "Cargo.toml  +1 −0", "described by D-002 code span"}
		at := 0
		for _, w := range want {
			i := strings.Index(s[at:], w)
			if i < 0 {
				t.Fatalf("prefix %q: the review lacks %q in order:\n%s", prefix, w, s)
			}
			at += i + len(w)
		}
		if strings.Contains(s, "W-001 link") || strings.Join(f.reads, ";") != "changes main abcdef1 a grove/work/W-001.md" {
			t.Fatalf("prefix %q: the open record is not listed, and nothing more is read: %v\n%s", prefix, f.reads, s)
		}
	}
}

// Records sharing a candidate on the branch (G-188) are named where the
// action covers them, and their files are not changes after the candidate.
func TestReviewNamesTheGroupSharingACandidate(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := reviewFixture(fx, true)
	for _, s := range []*versions.Source{fx.cFeat, fx.feat} {
		o := version(s, "W-003", "Build on it", "review")
		o.Record.Candidate, o.Record.Approved = "abcdef1", "abcdef1"
		f.res.Groups = append(f.res.Groups, versions.Group{ID: "W-003", Versions: []versions.Version{o}})
	}
	m := openReview(t, f, 160, 36)
	f.mu.Lock()
	reads := strings.Join(f.reads, "\n")
	f.mu.Unlock()
	if !strings.Contains(reads, "changes main abcdef1 a grove/work/W-001.md,grove/work/W-003.md") {
		t.Fatalf("the changes read leaves out both records. files:\n%s", reads)
	}
	press(m, "f")
	if s := plain(m); !strings.Contains(s, "Feedback on W-001, returning it and W-003, which share its candidate, to active on branch feature") {
		t.Fatalf("f names the group:\n%s", s)
	}
	press(m, "esc")
	press(m, "i")
	if s := plain(m); !strings.Contains(s, "Merge branch feature into main and mark W-001 and W-003 done? y/n") {
		t.Fatalf("i names the group:\n%s", s)
	}
}
