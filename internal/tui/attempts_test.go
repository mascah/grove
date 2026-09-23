package tui

import (
	"context"
	"errors"
	"fmt"
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
	// W-001's attempt branch was read: its tip holds the record in review.
	tip := version(source("committed", "", "worktree-W-001"), "W-001", "Inspect records", "review")
	res.Groups[0].Versions = append(res.Groups[0].Versions, tip)
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
		{view("W-002", "11", attempt.Finished, &attempt.Result{ExitCode: 1, Dirty: true, Events: attempt.Events{Result: &attempt.Final{Subtype: "error_max_budget_usd", IsError: true}}, Record: &attempt.State{Status: "review", Candidate: "c0ffee12"}}),
			"failed: error_max_budget_usd (exit 1); its record says review, uncommitted"},
		{view("W-002", "12", attempt.Finished, &attempt.Result{Dirty: true, Events: attempt.Events{Result: ok}, Record: &attempt.State{Status: "review", Candidate: "c0ffee12"}}),
			"ended without a handoff: W-002 is review on worktree-W-002, with no candidate; its record says review, uncommitted"},
		// A question open now does not make a stopped attempt, or an earlier one, a wait.
		{view("W-009", "13", attempt.Finished, &attempt.Result{ExitCode: 130, Stopped: true, Events: attempt.Events{Result: ok}}), "stopped (exit 130)"},
		{view("W-009", "0", attempt.Finished, &attempt.Result{Events: attempt.Events{Result: ok}, Record: &attempt.State{Status: "active"}}),
			"ended without a handoff: W-009 is active on worktree-W-009, with no candidate"},
		// On a branch the board read, the tip decides: leftover untracked files do not.
		{view("W-001", "14", attempt.Finished, &attempt.Result{Dirty: true, Events: attempt.Events{Result: ok}, Record: &attempt.State{Status: "review", Candidate: "c0ffee12", Revision: tip.Revision}}),
			"candidate ready: W-001 in review on worktree-W-001 with candidate c0ffee1"},
		{view("W-001", "15", attempt.Finished, &attempt.Result{Events: attempt.Events{Result: ok}, Record: &attempt.State{Status: "review", Candidate: "c0ffee12", Revision: "sha256:other"}}),
			"ended without a handoff: W-001 is review on worktree-W-001, with no candidate; its record says review, uncommitted"},
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
	older := view("W-001", "20260922T010000Z", attempt.Finished, &attempt.Result{ExitCode: 0, Events: attempt.Events{Result: &attempt.Final{Subtype: "success", CostUSD: 0.5}}, Record: &attempt.State{Status: "review", Candidate: "abcdef1"}})
	r.set(running, older)
	r.activity = attempt.Activity{Lines: []string{"started: model m", "tool: Bash go test", "evil \x1b]0;title\a text"}, Report: "## Done\n\nIt \x1b[31mworks."}
	m := openRuns(t, f, r, 120, 36) // a new session: the attempts are the same files
	if !m.ticking {
		t.Fatal("a running attempt schedules the next read")
	}
	settle(m, press(m, "A"))
	s := plain(m)
	for _, want := range []string{"Attempts in this repository  (2)", "> W-002.20260923T010000Z  running  worktree-W-002", "W-001.20260922T010000Z  candidate ready: W-001 in review on worktree-W-001 with candidate abcdef1  $0.50"} {
		if !strings.Contains(s, want) {
			t.Fatalf("the list lacks %q:\n%s", want, s)
		}
	}
	settle(m, press(m, "enter"))
	s = plain(m)
	for _, want := range []string{"Attempt W-002.20260923T010000Z of W-002", "Outcome: running", "Bounds: budget 2 USD, permission mode auto", "Final report", "## Done", `evil \x1b]0;title\a text`, "Activity, newest first (3)"} {
		if !strings.Contains(s, want) {
			t.Fatalf("the attempt lacks %q:\n%s", want, s)
		}
	}
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
	if m.gen == gen || f.inspects != inspects+2 || m.pending != "" || m.ticking || !strings.Contains(plain(m), "Outcome: stopped (exit 130)") {
		t.Fatalf("gen %d→%d inspects %d→%d pending %q ticking %v\n%s", gen, m.gen, inspects, f.inspects, m.pending, m.ticking, plain(m))
	}
	press(m, "x")
	if m.prompt != nil || !strings.Contains(plain(m), "is finished; nothing to stop") {
		t.Fatal(plain(m))
	}
	press(m, "o")
	if m.screen != detailScreen || m.openID() != "W-002" || !strings.Contains(plain(m), "Attempts: 1, latest W-002.20260923T010000Z stopped (exit 130) · A lists them") {
		t.Fatalf("o opens the work's record:\n%s", plain(m))
	}
	press(m, "A")
	if s := plain(m); !strings.Contains(s, "Attempts of W-002  (1)") || strings.Contains(s, "W-001.") {
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
		r.activity.Lines = append(r.activity.Lines, fmt.Sprintf("tool: Bash %s %d", strings.Repeat("長", 80), i))
	}
	r.activity.Cut = true
	for _, w := range []int{40, 80, 160} {
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
		press(m, "pgdown", "pgdown", "pgdown")
		if m.scroll == 0 {
			t.Fatal("PgDn scrolls")
		}
		for range 50 {
			press(m, "pgdown")
		}
		if s := plain(m); !strings.Contains(s, "earlier events are in") {
			t.Fatalf("width %d: the end says where the rest is:\n%s", w, s)
		}
		press(m, "esc")
		if m.screen != boardScreen {
			t.Fatal("Esc leaves")
		}
	}
}
