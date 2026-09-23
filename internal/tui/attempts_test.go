package tui

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mascah/grove/internal/attempt"
	"github.com/mascah/grove/internal/versions"
)

// runs is a fake of the attempts part of the backend: what it lists, the
// one attempt it shows, and a log of launches and stops.
type runs struct {
	mu        sync.Mutex
	views     []attempt.View
	activity  attempt.Activity
	launchErr error
	readErr   error // the next reads of one attempt fail after finding it
	launches  []attempt.Request
	stops     []string
	reads     int
}

func (r *runs) add(b Backend) Backend {
	b.Attempts = func(context.Context, string) ([]attempt.View, error) {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.reads++
		return append([]attempt.View(nil), r.views...), nil
	}
	b.Attempt = func(_ context.Context, _, id string) (*attempt.View, attempt.Activity, error) {
		r.mu.Lock()
		defer r.mu.Unlock()
		for _, v := range r.views {
			if v.Launch.Attempt == id {
				if r.readErr != nil {
					return &v, attempt.Activity{}, r.readErr
				}
				return &v, r.activity, nil
			}
		}
		return nil, attempt.Activity{}, errors.New("attempt " + id + " does not exist in this repository")
	}
	b.Launch = func(_ context.Context, req attempt.Request) ([]string, error) {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.launches = append(r.launches, req)
		if r.launchErr != nil {
			return nil, r.launchErr
		}
		return []string{"worktree: created", "attempt: " + req.ID + ".20260923T010000Z started"}, nil
	}
	b.Stop = func(_ context.Context, root, id string) ([]string, error) {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.stops = append(r.stops, root+" "+id)
		return []string{"stop: SIGTERM sent"}, nil
	}
	return b
}

func (r *runs) set(views ...attempt.View) {
	r.mu.Lock()
	r.views = views
	r.mu.Unlock()
}

func view(work, stamp string, status attempt.Status, res *attempt.Result) attempt.View {
	return attempt.View{
		Launch: attempt.Launch{Attempt: work + "." + stamp, Work: work, Branch: "worktree-" + work, Worktree: "/repo/.claude/worktrees/worktree-" + work,
			Base: strings.Repeat("b", 40), BudgetUSD: "2", PermissionMode: "auto", Started: time.Date(2026, 9, 23, 1, 0, 0, 0, time.UTC)},
		Status: status, Result: res, EventsPath: "/repo/.git/grove/attempts/" + work + "." + stamp + "/events.jsonl",
	}
}

// settle runs a command to rest: batches are opened, every message is fed
// back, and a poll's tick is dropped, so a test asks for each poll itself.
func settle(m *Model, cmd tea.Cmd) {
	if cmd == nil {
		return
	}
	switch msg := cmd().(type) {
	case tea.BatchMsg:
		for _, c := range msg {
			settle(m, c)
		}
	case attemptTick, nil:
	default:
		_, next := m.Update(msg)
		settle(m, next)
	}
}

func openRuns(t *testing.T, f *fake, r *runs, w, h int) *Model {
	t.Helper()
	m := New(t.Context(), "/repo/.", r.add(f.backend()))
	m.every = time.Millisecond
	m.clock = func() time.Time { return time.Date(2026, 9, 23, 3, 0, 0, 0, time.UTC) }
	cmd := m.Init() // the runtime's order: Init, then the size
	m.Update(tea.WindowSizeMsg{Width: w, Height: h})
	settle(m, cmd)
	return m
}

// R launches only proposed or active work, after a checked budget and a
// permission mode, with the revision this checkout showed; the outcome shows
// and the attempts are read again.
func TestLaunchFromTheDetail(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := &fake{res: fx.twoBranches()}
	f.res.Target = "main"
	onTarget(f.res, "W-002")
	r := &runs{}
	m := openRuns(t, f, r, 120, 36)
	if r.reads != 1 || f.inspects != 1 {
		t.Fatalf("opening reads the attempts once beside the board: %d, %d", r.reads, f.inspects)
	}
	m.openDetail("W-002")
	if s := plain(m); !strings.Contains(s, "Attempts: none · R launches one") || !strings.Contains(s, "R launch   A attempts") {
		t.Fatalf("a work detail names its attempts and the keys:\n%s", s)
	}
	press(m, "R")
	if s := plain(m); m.prompt == nil || !strings.Contains(s, "Launch W-002: budget in USD, required") || !strings.Contains(s, "on branch worktree-W-002") {
		t.Fatalf("R opens the budget prompt:\n%s", s)
	}
	for _, bad := range []string{"", "0", "-1", "two", "NaN"} {
		typeText(m, bad)
		press(m, "enter")
		if m.prompt == nil || m.prompt.kind != "budget" || !strings.Contains(plain(m), "the budget is a positive dollar amount") {
			t.Fatalf("budget %q is refused:\n%s", bad, plain(m))
		}
		for range bad {
			press(m, "backspace")
		}
	}
	typeText(m, "2.5")
	press(m, "enter")
	if s := plain(m); m.prompt == nil || m.prompt.kind != "mode" || !strings.Contains(s, "Launch W-002 for 2.5 USD: permission mode, required") {
		t.Fatalf("the budget leads to the permission mode:\n%s", s)
	}
	press(m, "enter")
	if !strings.Contains(plain(m), "type one permission mode") || len(r.launches) != 0 {
		t.Fatal("an empty mode is refused and nothing launches")
	}
	typeText(m, "auto")
	cmd := press(m, "enter")
	if m.pending != "act" || !strings.Contains(plain(m), "Launching the attempt…") {
		t.Fatalf("Enter launches: %q", m.pending)
	}
	settle(m, cmd)
	want := attempt.Request{Root: "/repo/.", ID: "W-002", Expect: fx.twoBranches().Groups[1].Versions[1].Revision, BudgetUSD: "2.5", PermissionMode: "auto"}
	if len(r.launches) != 1 || r.launches[0] != want {
		t.Fatalf("launched %+v, want %+v", r.launches, want)
	}
	if s := plain(m); m.screen != resultScreen || !strings.Contains(s, "Launch of an attempt of W-002") || !strings.Contains(s, "attempt: W-002.20260923T010000Z started") {
		t.Fatalf("the launch's facts show:\n%s", s)
	}
	if r.reads != 2 || f.inspects != 2 {
		t.Fatalf("a launch re-reads the attempts and the board: %d, %d", r.reads, f.inspects)
	}
	press(m, "esc")
	if m.screen != detailScreen || m.openID() != "W-002" {
		t.Fatal("Esc returns to the record")
	}
	// Esc cancels, and a refused launch shows why with nothing claimed.
	press(m, "R", "esc")
	if !strings.Contains(plain(m), "cancelled; nothing was launched") {
		t.Fatal(plain(m))
	}
	r.launchErr = errors.New("W-002 changed since it was read: grove/work/W-002.md is sha256:new here, not sha256:old")
	press(m, "R")
	typeText(m, "1")
	press(m, "enter")
	typeText(m, "auto")
	settle(m, press(m, "enter"))
	if s := plain(m); !strings.Contains(s, "NOT DONE: W-002 changed since it was read") {
		t.Fatalf("a refusal shows:\n%s", s)
	}
}

// onTarget marks every version of id as the target's bytes, as Inspect does
// for content the target holds.
func onTarget(res *versions.Result, id string) {
	for i := range res.Groups {
		if res.Groups[i].ID == id {
			for j := range res.Groups[i].Versions {
				res.Groups[i].Versions[j].OnTarget = true
			}
		}
	}
}

// Where a launch runs follows the current state, not the order the places
// were read in: the same bytes on the target and on another branch start
// afresh, whichever branch sorts first; work only on a branch continues
// there, in its checkout, with the revision of the checkout Grove was opened
// in, a linked worktree included; without a target, this checkout is the base.
func TestLaunchPlace(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	where := func(res *versions.Result, id string) attempt.Request {
		t.Helper()
		r := &runs{}
		m := openRuns(t, &fake{res: res}, r, 120, 36)
		m.root = "/repo/" + strings.TrimPrefix(strings.TrimPrefix(res.GitDir, "/repo/.git/"), "/repo/.git")
		m.openDetail(id)
		press(m, "R")
		if m.prompt == nil {
			t.Fatalf("%s: no prompt:\n%s", id, plain(m))
		}
		typeText(m, "1")
		press(m, "enter")
		typeText(m, "auto")
		settle(m, press(m, "enter"))
		return r.launches[0]
	}
	// Inspect orders committed branches by name: feature before main.
	res := fx.twoBranches()
	res.Target = "main"
	onTarget(res, "W-002")
	for i := range res.Groups {
		if g := &res.Groups[i]; g.ID == "W-002" {
			g.Versions = append(g.Versions[2:], g.Versions[:2]...)
		}
	}
	if got := where(res, "W-002"); got.Branch != "" || got.Worktree != "" || got.Expect != res.Groups[1].Versions[3].Revision {
		t.Fatalf("the target holds W-002, so it starts afresh: %+v", got)
	}
	// Opened in the feature checkout, a linked worktree: W-010 is only there.
	res = fx.twoBranches()
	res.Target, res.GitDir = "main", fx.feat.GitDir
	var feat *versions.Version
	for i := range res.Groups {
		for j := range res.Groups[i].Versions {
			if v := &res.Groups[i].Versions[j]; res.Groups[i].ID == "W-010" && v.Source == fx.feat {
				feat = v
			}
		}
	}
	if got := where(res, "W-010"); got.Branch != "feature" || got.Worktree != "/repo/feat" || got.Expect != feat.Revision || got.Root != "/repo/feat" {
		t.Fatalf("W-010 continues on feature from this checkout: %+v", got)
	}
	// A fresh start reuses the default branch's checkout wherever it is.
	res = fx.twoBranches()
	res.Target = "main"
	onTarget(res, "W-002")
	elsewhere := source("live", "wt2", "worktree-W-002")
	elsewhere.Worktree = "/elsewhere/W-002"
	res.Sources = append(res.Sources, elsewhere)
	if got := where(res, "W-002"); got.Branch != "worktree-W-002" || got.Worktree != "/elsewhere/W-002" {
		t.Fatalf("the default branch's checkout is reused: %+v", got)
	}
	// No target: this checkout's own state is the base.
	res = fx.twoBranches()
	if got := where(res, "W-002"); got.Branch != "" {
		t.Fatalf("without a target the checkout's state starts afresh: %+v", got)
	}
}

// Work in review, done, not work, or with an attempt that may be running is
// refused before any prompt, with what to do instead.
func TestLaunchRefusals(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := reviewFixture(fx, false)
	r := &runs{}
	m := openRuns(t, f, r, 120, 36)
	m.openDetail("W-001")
	press(m, "R")
	if m.prompt != nil || !strings.Contains(plain(m), "W-001 is in review: judge its candidate (a approve, f feedback) before another attempt") {
		t.Fatalf("review refuses:\n%s", plain(m))
	}
	press(m, "esc")
	m.openDetail("W-006")
	press(m, "R")
	if m.prompt != nil || !strings.Contains(plain(m), "W-006 is not work") {
		t.Fatal(plain(m))
	}
	// A running attempt: the duplicate is refused here and by Start; the card is tagged.
	f2 := &fake{res: fx.twoBranches()}
	r2 := &runs{}
	r2.set(view("W-002", "20260923T005900Z", attempt.Running, nil))
	m = openRuns(t, f2, r2, 120, 36)
	if !onRow(plain(m), "W-002", "● running") {
		t.Fatalf("a running attempt tags its card:\n%s", plain(m))
	}
	m.openDetail("W-002")
	press(m, "R")
	if m.prompt != nil || !strings.Contains(plain(m), "attempt W-002.20260923T005900Z of W-002 is running; A shows it, x stops it") {
		t.Fatalf("a duplicate start is refused:\n%s", plain(m))
	}
	if len(r2.launches) != 0 {
		t.Fatal("nothing launched")
	}
	// An open question that blocks the work is a wait, not a launch.
	res := fx.twoBranches()
	q := version(fx.cMain, "Q-002", "Red or blue?", "open")
	q.Record.Blocks = []string{"W-002"}
	res.Groups = append(res.Groups, versions.Group{ID: "Q-002", Versions: []versions.Version{q}})
	m = openRuns(t, &fake{res: res}, &runs{}, 120, 36)
	m.openDetail("W-002")
	press(m, "R")
	if m.prompt != nil || !strings.Contains(plain(m), "W-002 is blocked by open question Q-002 (Red or blue?); resolve it before another attempt") {
		t.Fatalf("a blocking question refuses:\n%s", plain(m))
	}
}

// After feedback the record is active on its candidate's branch: R
// continues there, in that branch's checkout, and feedback's outcome says so.
func TestFeedbackContinuesOnTheCandidateBranch(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := reviewFixture(fx, false)
	r := &runs{}
	m := openRuns(t, f, r, 120, 36)
	press(m, "right", "right")
	settle(m, press(m, "enter"))
	press(m, "f")
	typeText(m, "Handle the empty case")
	settle(m, press(m, "enter"))
	if s := plain(m); !strings.Contains(s, "or: R on W-001 launches a bounded attempt on its branch") {
		t.Fatalf("feedback names the relaunch:\n%s", s)
	}
	// The re-read board: feature now holds W-001 active.
	for _, g := range f.res.Groups {
		for i := range g.Versions {
			if v := &g.Versions[i]; v.Record.ID == "W-001" && v.Source.Ref == "refs/heads/feature" {
				v.Record.Status = "active"
			}
		}
	}
	settle(m, press(m, "r"))
	press(m, "esc")
	press(m, "R")
	if m.prompt == nil || !strings.Contains(plain(m), "on branch feature in /repo/feat") {
		t.Fatalf("the continuation runs on the candidate's branch:\n%s", plain(m))
	}
	typeText(m, "1")
	press(m, "enter")
	typeText(m, "acceptEdits")
	settle(m, press(m, "enter"))
	if got := r.launches[0]; got.Branch != "feature" || got.Worktree != "/repo/feat" || got.Expect == "" {
		t.Fatalf("%+v", got)
	}
}

// The outcome is derived from the attempt's files and the records: only a
// record in review with a candidate is ready, whatever the exit says.
func TestAttemptOutcomes(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	res := fx.twoBranches()
	q := version(fx.cMain, "Q-002", "Red or blue?", "open")
	q.Record.Blocks = []string{"W-009"}
	res.Groups = append(res.Groups, versions.Group{ID: "Q-002", Versions: []versions.Version{q}})
	m := openRuns(t, &fake{res: res}, &runs{}, 120, 36)
	ok := &attempt.Final{Subtype: "success"}
	cases := []struct {
		v    attempt.View
		want string
	}{
		{view("W-002", "1", attempt.Running, nil), "running"},
		{view("W-002", "2", attempt.Orphaned, nil), "orphaned: its owner is gone and the provider still runs (x stops it)"},
		{view("W-002", "3", attempt.Interrupted, nil), "interrupted: the owner and the provider are gone without a result"},
		{view("W-002", "4", attempt.Finished, &attempt.Result{Events: attempt.Events{Result: ok}, Record: &attempt.State{Status: "review", Candidate: "c0ffee12"}}),
			"candidate ready: W-002 in review on worktree-W-002 with candidate c0ffee1"},
		{view("W-002", "5", attempt.Finished, &attempt.Result{Events: attempt.Events{Result: ok}, Record: &attempt.State{Status: "active"}}),
			"ended without a handoff: W-002 is active on worktree-W-002, with no candidate"},
		{view("W-009", "6", attempt.Finished, &attempt.Result{Events: attempt.Events{Result: ok}, Record: &attempt.State{Status: "active"}}),
			"waiting on question Q-002 (Red or blue?)"},
		{view("W-002", "7", attempt.Finished, &attempt.Result{ExitCode: 130, Stopped: true, Events: attempt.Events{Result: ok}}), "stopped (exit 130)"},
		{view("W-002", "8", attempt.Finished, &attempt.Result{ExitCode: -1, Signal: "killed"}), "failed: no result event (killed)"},
		{view("W-002", "9", attempt.Finished, &attempt.Result{ExitCode: 1, Events: attempt.Events{Result: &attempt.Final{Subtype: "error_max_budget_usd", IsError: true}}}),
			"failed: error_max_budget_usd (exit 1)"},
		{view("W-002", "10", attempt.Finished, &attempt.Result{Events: attempt.Events{Result: &attempt.Final{Subtype: "success"}}, Record: &attempt.State{Status: "review"}}),
			"ended without a handoff: W-002 is review on worktree-W-002, with no candidate"},
		// The record says review in files never committed, and the run failed.
		{view("W-002", "11", attempt.Finished, &attempt.Result{ExitCode: 1, Dirty: true, RecordUncommitted: true, Events: attempt.Events{Result: &attempt.Final{Subtype: "error_max_budget_usd", IsError: true}}, Record: &attempt.State{Status: "review", Candidate: "c0ffee12"}}),
			"failed: error_max_budget_usd (exit 1); its record says review with candidate c0ffee1, uncommitted"},
		{view("W-002", "12", attempt.Finished, &attempt.Result{Dirty: true, RecordUncommitted: true, Events: attempt.Events{Result: ok}, Record: &attempt.State{Status: "review", Candidate: "c0ffee12"}}),
			"ended without a handoff: W-002 is review on worktree-W-002; its record says review with candidate c0ffee1, uncommitted"},
		// A question open now does not make a stopped attempt, or an earlier one, a wait.
		{view("W-009", "13", attempt.Finished, &attempt.Result{ExitCode: 130, Stopped: true, Events: attempt.Events{Result: ok}}), "stopped (exit 130)"},
		{view("W-009", "0", attempt.Finished, &attempt.Result{Events: attempt.Events{Result: ok}, Record: &attempt.State{Status: "active"}}),
			"ended without a handoff: W-009 is active on worktree-W-009, with no candidate"},
		// The record's own file decides, as the owner found it at the end:
		// leftover untracked files do not, and nothing committed later does.
		{view("W-001", "14", attempt.Finished, &attempt.Result{Dirty: true, Events: attempt.Events{Result: ok}, Record: &attempt.State{Status: "review", Candidate: "c0ffee12", Revision: "sha256:then"}}),
			"candidate ready: W-001 in review on worktree-W-001 with candidate c0ffee1"},
	}
	m.attempts = []attempt.View{view("W-009", "6", attempt.Finished, nil), view("W-009", "0", attempt.Finished, nil)}
	for _, c := range cases {
		if got := m.outcomeOf(&c.v); got != c.want {
			t.Errorf("%s: %q, want %q", c.v.Launch.Attempt, got, c.want)
		}
	}
}

// A lists attempts, Enter shows one with its facts, final report and
// activity, all escaped; x asks, then stops through the backend; while one
// runs the list is polled, and its end re-reads the board.
func TestAttemptScreensReconnectAndStop(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := &fake{res: fx.twoBranches()}
	r := &runs{}
	running := view("W-002", "20260923T010000Z", attempt.Running, nil)
	older := view("W-001", "20260922T010000Z", attempt.Finished, &attempt.Result{Finished: time.Date(2026, 9, 22, 2, 0, 0, 0, time.UTC), ExitCode: 0, Events: attempt.Events{Result: &attempt.Final{Subtype: "success", CostUSD: 0.5}}, Record: &attempt.State{Status: "review", Candidate: "abcdef1"}})
	r.set(running, older)
	r.activity = attempt.Activity{Entries: []attempt.Entry{{Kind: "start", Text: "started: model m"}, {Kind: "tool", Text: "Bash go test", Time: time.Date(2026, 9, 23, 1, 2, 3, 0, time.UTC)}, {Kind: "notice", Text: "system: thinking_tokens", Count: 14}, {Kind: "text", Text: "evil \x1b]0;title\a text"}}, Report: "## Done\n\nIt \x1b[31mworks."}
	m := openRuns(t, f, r, 120, 36) // a new session: the attempts are the same files
	if !m.ticking {
		t.Fatal("a running attempt schedules the next read")
	}
	settle(m, press(m, "A"))
	s := plain(m)
	// W-001 is active on the board: its candidate had feedback, so the
	// next attempt is the owner's to launch.
	for _, want := range []string{"Attempts in this repository (2) · 1 need you · 1 running · 0 settled", " Needs you", "> W-001  Inspect records", "feedback given: R again      25h ago", "  W-002  Create records", "running                    2h so far"} {
		if !strings.Contains(s, want) {
			t.Fatalf("the list lacks %q:\n%s", want, s)
		}
	}
	press(m, "down")
	settle(m, press(m, "enter"))
	s = plain(m)
	for _, want := range []string{"W-002  Create records", "● running  2h so far", "Attempt  W-002.20260923T010000Z", "Budget   spend known at the end, of $2 · permission mode auto", "Model    not reported yet", "State  Running.", "Next   x stops it", "d shows them", "Final report", "## Done", `evil \x1b]0;title\a text`, "Activity, newest first (4)", time.Date(2026, 9, 23, 1, 2, 3, 0, time.UTC).Local().Format("15:04:05") + " › Bash go test", "         · system: thinking_tokens (×14)"} {
		if !strings.Contains(s, want) {
			t.Fatalf("the attempt lacks %q:\n%s", want, s)
		}
	}
	if strings.Contains(s, "Bounds:") {
		t.Fatal("the details are folded")
	}
	press(m, "d")
	if s := plain(m); !strings.Contains(s, "Bounds: budget 2 USD, permission mode auto") || !strings.Contains(s, "d hides them") {
		t.Fatalf("d shows the details:\n%s", s)
	}
	press(m, "d")
	if raw := m.render(); strings.Contains(raw, "\x1b]0;") || strings.Contains(raw, "\x1b[31mworks") {
		t.Fatal("provider text reached the terminal unescaped")
	}
	if strings.Index(s, "evil") > strings.Index(s, "started: model m") {
		t.Fatal("the newest activity comes first")
	}
	press(m, "x")
	if m.prompt == nil || !strings.Contains(plain(m), "Stop attempt W-002.20260923T010000Z of W-002? Its partial work stays. y/n") {
		t.Fatal(plain(m))
	}
	press(m, "n")
	if len(r.stops) != 0 || !strings.Contains(plain(m), "cancelled; nothing was stopped") {
		t.Fatal("n stops nothing")
	}
	press(m, "x")
	settle(m, press(m, "y"))
	if len(r.stops) != 1 || r.stops[0] != "/repo/. W-002.20260923T010000Z" || !strings.Contains(plain(m), "Stop of attempt W-002.20260923T010000Z") {
		t.Fatalf("y stops: %v\n%s", r.stops, plain(m))
	}
	press(m, "esc")
	if m.screen != attemptScreen {
		t.Fatalf("Esc returns to the attempt: %d", m.screen)
	}
	// The poll sees the attempt end: the board is re-read, and polling stops.
	// An inspection already under way when it ends may predate its last
	// commit, so it is started again.
	stopped := view("W-002", "20260923T010000Z", attempt.Finished, &attempt.Result{ExitCode: 130, Stopped: true})
	r.set(stopped, older)
	early := press(m, "r")
	gen, inspects := m.gen, f.inspects
	m.ticking = false
	_, cmd := m.Update(attemptTick{})
	settle(m, cmd)
	settle(m, early) // its reply is outdated by the restart
	if m.gen == gen || f.inspects != inspects+2 || m.pending != "" || m.ticking || !strings.Contains(plain(m), "State  Stopped (exit 130).") || !strings.Contains(plain(m), "· stopped by x") {
		t.Fatalf("gen %d→%d inspects %d→%d pending %q ticking %v\n%s", gen, m.gen, inspects, f.inspects, m.pending, m.ticking, plain(m))
	}
	press(m, "x")
	if m.prompt != nil || !strings.Contains(plain(m), "is finished; nothing to stop") {
		t.Fatal(plain(m))
	}
	press(m, "o")
	if m.screen != detailScreen || m.openID() != "W-002" || !strings.Contains(plain(m), "Attempts 1 · latest: stopped by x, 2h ago · A lists them") {
		t.Fatalf("o opens the work's record:\n%s", plain(m))
	}
	press(m, "A")
	if s := plain(m); !strings.Contains(s, "Attempts of W-002 (1)") || strings.Contains(s, "W-001.") {
		t.Fatalf("A on a work lists its attempts:\n%s", s)
	}
	press(m, "esc")
	if m.screen != detailScreen {
		t.Fatal("Esc returns to the record")
	}
	press(m, "esc")
	if m.screen != attemptScreen || len(m.stack) != 0 {
		t.Fatalf("Esc from the record o opened returns to the attempt: screen %d stack %v", m.screen, m.stack)
	}
}

// Many activity lines and a narrow terminal: every row keeps its width, and
// scrolling and leaving stay immediate.
func TestAttemptActivityIsBounded(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	r := &runs{}
	r.set(view("W-002", "20260923T010000Z", attempt.Running, nil))
	for i := range 200 {
		r.activity.Entries = append(r.activity.Entries, attempt.Entry{Kind: "tool", Text: fmt.Sprintf("Bash %s %d", strings.Repeat("長", 80), i), Count: (i + 1) % 3})
	}
	r.activity.Cut = true
	r.activity.Report = strings.Repeat("A long report line with `code` and **bold** text that wraps.\n\n", 40)
	r.activity.Metrics = attempt.Metrics{Turns: 205, OutputTokens: 93069, Context: 412345, Window: 1000000, Subagents: 3, Tools: 88}
	for _, w := range []int{40, 80, 99, 100, 160} {
		m := openRuns(t, &fake{res: fx.twoBranches()}, r, w, 20)
		m.openAttempt("W-002.20260923T010000Z")
		settle(m, m.wantAttempts()) // opened directly, so the read is asked for here
		rows := strings.Split(m.render(), "\n")
		if len(rows) != 20 {
			t.Fatalf("width %d: %d rows", w, len(rows))
		}
		for _, row := range rows {
			if ansi.StringWidth(row) != w {
				t.Fatalf("width %d: a row of %d cells: %q", w, ansi.StringWidth(row), ansi.Strip(row))
			}
		}
		if s := plain(m); w == 160 && (!strings.Contains(s, "Turns ≥205   Tokens ≥0 in · – out   Context 412k of 1M ▰▰▰▰▱▱▱▱▱▱ 41%   Subagents ≥3") || !strings.Contains(s, "│          › Bash") || !strings.Contains(s, "(×2)")) {
			t.Fatalf("the metrics and the timeline beside the report:\n%s", s)
		}
		press(m, "pgdown", "pgdown", "pgdown")
		if m.scroll == 0 {
			t.Fatal("PgDn scrolls")
		}
		m.scroll = 1 << 20 // the end, as PgDn would reach it, without paging there
		m.clampScroll()
		if s := plain(m); !strings.Contains(s, "earlier events are in") {
			t.Fatalf("width %d: the end says where the rest is:\n%s", w, s)
		}
		press(m, "esc")
		if m.screen != boardScreen {
			t.Fatal("Esc leaves")
		}
	}
}

// Where each attempt stands now (G-117): only the latest attempt of work
// still proposed, active or in review needs the owner, and an orphan always
// does; every settled attempt says why.
func TestAttemptStandings(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	res := fx.twoBranches()
	add := func(id, title, status, candidate string, blocks ...string) {
		v := version(fx.cMain, id, title, status)
		v.Record.Candidate, v.Record.Blocks = candidate, blocks
		res.Groups = append(res.Groups, versions.Group{ID: id, Versions: []versions.Version{v}})
	}
	add("W-101", "Judge me", "review", "c0ffee12")
	add("W-102", "Merged", "done", "c0ffee12")
	add("W-103", "Moved on", "done", "d00d")
	add("W-104", "Retried", "active", "")
	add("W-105", "Given feedback", "active", "c0ffee12")
	add("W-106", "Dropped", "abandoned", "")
	add("Q-002", "Red or blue?", "open", "", "W-107")
	add("W-107", "Asked", "active", "")
	add("W-108", "Continued by hand", "review", "beefcafe")
	m := openRuns(t, &fake{res: res}, &runs{}, 120, 36)
	ok := &attempt.Final{Subtype: "success"}
	ready := func(id, stamp string) attempt.View {
		return view(id, stamp, attempt.Finished, &attempt.Result{Events: attempt.Events{Result: ok}, Record: &attempt.State{Status: "review", Candidate: "c0ffee12"}})
	}
	failed := view("W-104", "1", attempt.Finished, &attempt.Result{ExitCode: 1, Events: attempt.Events{Result: &attempt.Final{Subtype: "error_max_budget_usd", IsError: true}}})
	cases := []struct {
		v                  attempt.View
		group              int
		short, state, next string
	}{
		{view("W-002", "1", attempt.Running, nil), runningNow, "running", "Running.", "x stops it; the report arrives when it ends"},
		{view("W-002", "0", attempt.Orphaned, nil), needsYou, "orphaned: x stops it", "Orphaned: its owner is gone and the provider still runs (x stops it).", "x stops it"},
		{ready("W-101", "1"), needsYou, "judge candidate c0ffee1", "Candidate ready: W-101 in review on worktree-W-101 with candidate c0ffee1.", "o opens W-101: a approves, f gives feedback"},
		{ready("W-102", "1"), settled, "done: candidate c0ffee1", "Candidate ready: W-102 in review on worktree-W-102 with candidate c0ffee1. Candidate c0ffee1 was integrated: W-102 is done.", ""},
		{ready("W-103", "1"), settled, "candidate c0ffee1, superseded", "Candidate ready: W-103 in review on worktree-W-103 with candidate c0ffee1. W-103 has moved on: it is done with candidate d00d.", ""},
		{failed, settled, "failed: budget exhausted, superseded", "Failed: error_max_budget_usd (exit 1). A later attempt of W-104 followed this one.", ""},
		{view("W-104", "2", attempt.Interrupted, nil), needsYou, "interrupted", "Interrupted: the owner and the provider are gone without a result.", "d shows the details and the raw log; R on W-104 launches again"},
		{ready("W-105", "1"), needsYou, "feedback given: R again", "Candidate ready: W-105 in review on worktree-W-105 with candidate c0ffee1. W-105 is active now: it had feedback.", "o opens W-105: R launches the next attempt"},
		{view("W-106", "1", attempt.Finished, &attempt.Result{Events: attempt.Events{Result: ok}, Record: &attempt.State{Status: "active"}}), settled, "work abandoned since", "Ended without a handoff: W-106 is active on worktree-W-106, with no candidate. W-106 is abandoned now; nothing needs you.", ""},
		{view("W-107", "1", attempt.Finished, &attempt.Result{Events: attempt.Events{Result: ok}, Record: &attempt.State{Status: "active"}}), needsYou, "answer question Q-002", "Waiting on question Q-002 (Red or blue?).", "o opens W-107, whose detail lists the question"},
		{view("W-107", "0", attempt.Finished, &attempt.Result{ExitCode: 130, Stopped: true}), settled, "stopped by x", "Stopped (exit 130).", ""},
		{view("W-001", "1", attempt.Finished, &attempt.Result{Events: attempt.Events{Result: ok}, Record: &attempt.State{Status: "active"}}), needsYou, "ended, no handoff", "Ended without a handoff: W-001 is active on worktree-W-001, with no candidate.", "o opens W-001; the report says why"},
		// Continued outside an attempt to another candidate: a approves that one, not this.
		{ready("W-108", "1"), settled, "candidate c0ffee1, superseded", "Candidate ready: W-108 in review on worktree-W-108 with candidate c0ffee1. W-108 has moved on: it is in review with candidate beefcaf.", ""},
		{ready("W-999", "1"), needsYou, "judge candidate c0ffee1", "Candidate ready: W-999 in review on worktree-W-999 with candidate c0ffee1.", "o opens W-999: a approves, f gives feedback"},
	}
	for _, c := range cases {
		m.attempts = append(m.attempts, c.v)
	}
	// Newest first: each work's second attempt here is its latest.
	slices.SortStableFunc(m.attempts, func(a, b attempt.View) int { return strings.Compare(b.Launch.Attempt, a.Launch.Attempt) })
	for _, c := range cases {
		got := m.standingOf(&c.v)
		if got != (standing{c.group, c.short, c.state, c.next}) {
			t.Errorf("%s:\n got %+v\nwant %+v", c.v.Launch.Attempt, got, standing{c.group, c.short, c.state, c.next})
		}
	}
	if got := m.titleOf("W-999"); got != "(title unread)" {
		t.Fatal(got)
	}
}

// Two hundred attempts with long titles: every row keeps its width, the
// title goes before the state is cut, and the cursor stays in view.
func TestAttemptListFits(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	res := fx.twoBranches()
	for _, g := range res.Groups {
		for i := range g.Versions {
			g.Versions[i].Record.Title = strings.Repeat("A very long title 長 ", 10)
		}
	}
	r := &runs{}
	var views []attempt.View
	for i := range 200 {
		v := view("W-002", fmt.Sprintf("2026092%dT%06dZ", i%2, 200-i), attempt.Finished, &attempt.Result{ExitCode: 130, Stopped: true})
		if i%3 == 0 { // every group mixed with the cursor
			v.Status, v.Result = attempt.Orphaned, nil
		}
		views = append(views, v)
	}
	views[1].Status, views[1].Result = attempt.Running, nil
	r.set(views...)
	for _, size := range [][2]int{{40, 10}, {59, 24}, {60, 24}, {80, 24}, {160, 48}} {
		w, h := size[0], size[1]
		m := openRuns(t, &fake{res: res}, r, w, h)
		settle(m, press(m, "A"))
		for range 150 {
			press(m, "down")
		}
		rows := strings.Split(m.render(), "\n")
		if len(rows) != h {
			t.Fatalf("%dx%d: %d rows", w, h, len(rows))
		}
		for _, row := range rows {
			if ansi.StringWidth(row) != w {
				t.Fatalf("%dx%d: a row of %d cells: %q", w, h, ansi.StringWidth(row), ansi.Strip(row))
			}
		}
		s := plain(m)
		if !strings.Contains(s, "> W-002") || !strings.Contains(s, "stopped by x") {
			t.Fatalf("%dx%d: the cursor's row and its state show:\n%s", w, h, s)
		}
		for range 150 {
			press(m, "up")
		}
		if s := plain(m); !strings.Contains(s, " Needs you") || !strings.Contains(s, "orphaned") {
			t.Fatalf("%dx%d: back at the top:\n%s", w, h, s)
		}
		if strings.Contains(s, "A very long") != (w >= 60) {
			t.Fatalf("%dx%d: the title shows from 60 columns:\n%s", w, h, s)
		}
	}
}

// What the attempt screen cannot know it says: the model from the owner's
// scan when the window lacks it, spend without a result, turns counted from
// messages as approximate, and a failed read keeps the last good one.
func TestAttemptScreenHonesty(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	r := &runs{}
	stopped := view("W-002", "20260923T010000Z", attempt.Finished, &attempt.Result{ExitCode: 130, Stopped: true, Events: attempt.Events{Init: &attempt.Init{Model: "model-at-exit"}}})
	r.set(stopped)
	r.activity = attempt.Activity{Cut: true, Entries: []attempt.Entry{{Kind: "text", Text: "last words", Count: 1}}, Metrics: attempt.Metrics{Turns: 4}}
	m := openRuns(t, &fake{res: fx.twoBranches()}, r, 120, 36)
	m.openAttempt(stopped.Launch.Attempt)
	settle(m, m.wantAttempts())
	s := plain(m)
	for _, want := range []string{"Model    model-at-exit", "Budget   spend unknown: no result event, of $2", "Turns ≥4", "· stopped by x"} {
		if !strings.Contains(s, want) {
			t.Fatalf("lacks %q:\n%s", want, s)
		}
	}
	running := view("W-002", "20260923T010000Z", attempt.Running, nil)
	r.set(running)
	settle(m, press(m, "r"))
	if s := plain(m); !strings.Contains(s, "Model    not in the part of the log read") {
		t.Fatalf("a cut window without the start:\n%s", s)
	}
	r.activity.Cut = false
	settle(m, press(m, "r"))
	if s := plain(m); !strings.Contains(s, "Model    not reported yet") || !strings.Contains(s, "Turns ≈4") {
		t.Fatalf("an uncut window without a result:\n%s", s)
	}
	r.readErr = errors.New("permission denied")
	settle(m, press(m, "r"))
	if s := plain(m); !strings.Contains(s, "The last read failed (r retries): permission denied") || !strings.Contains(s, "last words") || !strings.Contains(s, "Turns ≈4") {
		t.Fatalf("a failed read keeps the last good one:\n%s", s)
	}
	m.openAttempt("W-002.20260923T020000Z") // never read: nothing to keep
	r.set(running, view("W-002", "20260923T020000Z", attempt.Running, nil))
	settle(m, m.wantAttempts())
	if s := plain(m); !strings.Contains(s, "could not be read (r retries): permission denied") || strings.Contains(s, "Turns") {
		t.Fatalf("a first read that fails shows no figures:\n%s", s)
	}
}
