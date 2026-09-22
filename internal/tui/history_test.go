package tui

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/mascah/grove/internal/versions"
)

var (
	mainTime = time.Date(2026, 9, 19, 16, 7, 0, 0, time.Local)
	mainLog  = []versions.Commit{
		{ID: strings.Repeat("1", 40), Subject: "close \x1b]0;owned\a\u202eit", Status: "active", When: mainTime},
		{ID: strings.Repeat("2", 40), Subject: "propose it", Status: "proposed", When: mainTime.Add(-time.Hour)},
		{ID: strings.Repeat("3", 40), Subject: "before statuses", When: mainTime.Add(-2 * time.Hour)},
	}
	featLog = []versions.Commit{{ID: strings.Repeat("4", 40), Subject: "finish on feature", Status: "done", When: mainTime}}
)

// lineageFixture is twoBranches with feature at its own commit, and a history
// per commit.
func lineageFixture() (fixture, *fake) {
	fx := newFixture()
	fx.cFeat.Commit, fx.feat.Commit = strings.Repeat("b", 40), strings.Repeat("b", 40)
	f := &fake{res: fx.twoBranches(), ws: &versions.Workspace{Project: "/repo/feat"}}
	f.history = func(_ context.Context, commit, _ string) ([]versions.Commit, error) {
		if commit[0] == 'b' {
			return featLog, nil
		}
		return mainLog, nil
	}
	return fx, f
}

func TestHistoryFollowsTheFocusedVersion(t *testing.T) {
	t.Parallel()
	_, f := lineageFixture()
	m := open(t, f, 120, 30)
	chooseCheckout(m, 0) // under the ID header, history follows the board's checkout
	press(m, "right", "left", "right", "b", "esc", "s", "esc")
	if len(f.histories) != 0 || m.pending != "" {
		t.Fatalf("the board read histories: %v", f.histories)
	}
	cmd := press(m, "enter") // W-001, active on main's checkout
	if cmd == nil || m.pending != "history" || !strings.Contains(plain(m), "reading…") {
		t.Fatalf("opening a card should start reading its history: pending %q\n%s", m.pending, plain(m))
	}
	if next := deliver(m, cmd); next != nil || m.pending != "" {
		t.Fatal("a delivered history must not start another read")
	}
	press(m, "v")
	screen := plain(m)
	for _, want := range []string{
		"History on checkout . (main): commits that changed",
		"2026-09-19 16:07  active     1111111  close \\x1b]0;owned\\a\\u202eit",
		"2026-09-19 15:07  proposed   2222222  propose it",
		"2026-09-19 14:07  ?          3333333  before statuses",
		"A version is this record's exact content",
	} {
		if !strings.Contains(screen, want) {
			t.Fatalf("the open card lacks %q:\n%s", want, screen)
		}
	}
	if raw := m.render(); strings.Contains(raw, "\x1b]0;") || strings.Contains(raw, "\u202e") || strings.Contains(raw, "\a") {
		t.Fatal("a commit subject reached the terminal unescaped")
	}
	if strings.Contains(screen, "finish on feature") || strings.Index(screen, "active     1111111") > strings.Index(screen, "proposed   2222222") {
		t.Fatalf("expected main's lineage alone, newest first:\n%s", screen)
	}
	// The same commit and path is one read, however many places hold it.
	if cmd := press(m, "down"); cmd != nil || !strings.Contains(plain(m), "History on branch main:") {
		t.Fatalf("main's branch tip is the commit already read:\n%s", plain(m))
	}
	// The other branch has its own lineage.
	if screen = plain(deliverAll(m, press(m, "down"))); !strings.Contains(screen, "History on branch feature:") ||
		!strings.Contains(screen, "done       4444444  finish on feature") || strings.Contains(screen, "propose it") {
		t.Fatalf("expected feature's lineage alone:\n%s", screen)
	}
	if cmd := press(m, "up", "up", "down", "down"); cmd != nil {
		t.Fatal("a history already held was read again")
	}
	if got := strings.Join(f.histories, ","); got != "a grove/work/W-001.md,b grove/work/W-001.md" {
		t.Fatalf("history reads: %s", got)
	}
	// A refresh forgets them: the branches may have moved.
	deliverAll(m, press(m, "r"))
	if len(f.histories) != 3 || f.inspects != 2 {
		t.Fatalf("after a refresh: %d inspections, histories %v", f.inspects, f.histories)
	}
}

// Merges are not listed, so the newest row may not be the record's status
// here. The card says so rather than let the first row stand for it.
func TestHistorySaysWhenTheNewestCommitIsNotTheRecord(t *testing.T) {
	t.Parallel()
	_, f := lineageFixture()
	log := []versions.Commit{featLog[0], mainLog[0]} // newest: done; then active
	f.history = func(context.Context, string, string) ([]versions.Commit, error) { return log, nil }
	m := open(t, f, 120, 30)
	chooseCheckout(m, 0)
	deliverAll(m, press(m, "right", "enter")) // W-001 is active here
	press(m, "v")
	screen := plain(m)
	note, newest := strings.Index(screen, "here              active     the record's status here"), strings.Index(screen, "done       4444444")
	if note < 0 || newest < note {
		t.Fatalf("expected the record's status above the newest commit:\n%s", screen)
	}
	if screen = plain(deliverAll(m, press(m, "down", "down"))); strings.Contains(screen, "merges are not listed") {
		t.Fatalf("feature is done, as its newest commit says:\n%s", screen)
	}
	// A checkout with uncommitted changes has its own first row instead.
	f.res.Groups[0].Versions[1].Change = "modified"
	press(m, "up", "up")
	if screen = plain(m); strings.Contains(screen, "merges are not listed") || !strings.Contains(screen, "uncommitted       active     modified") {
		t.Fatalf("a modified checkout:\n%s", screen)
	}
}

func deliverAll(m *Model, cmd tea.Cmd) *Model {
	for cmd != nil {
		cmd = deliver(m, cmd)
	}
	return m
}

func TestHistoryOfUncommittedChanges(t *testing.T) {
	t.Parallel()
	fx, f := lineageFixture()
	vs := f.res.Groups[0].Versions // W-001: branch main, checkout ., branch feature, checkout feat
	vs[1].Change, vs[1].HeadPath = "renamed", "grove/work/W-001-old.md"
	vs[3].Change = "added"
	f.history = func(context.Context, string, string) ([]versions.Commit, error) {
		return nil, errors.New("git log: \x1b[31mbroken")
	}
	m := deliverAll(open(t, f, 120, 40), nil)
	chooseCheckout(m, 0)
	deliverAll(m, press(m, "right", "enter"))
	press(m, "v")
	screen := plain(m)
	for _, want := range []string{"uncommitted       active     renamed in this checkout's files", "could not be read (r retries): git log: \\x1b[31mbroken"} {
		if !strings.Contains(screen, want) {
			t.Fatalf("a renamed checkout lacks %q:\n%s", want, screen)
		}
	}
	if got := strings.Join(f.histories, ","); got != "a grove/work/W-001-old.md" {
		t.Fatalf("a checkout is followed from its path at HEAD: %s", got)
	}
	// An added record has no commit to read. feat is the last row.
	press(m, "esc", "esc")
	chooseCheckout(m, 1) // the feature checkout's board
	if fx.feat != m.boardSource() {
		t.Fatal("expected the feature board")
	}
	before := len(f.histories)
	cmd := press(m, "right", "right", "right", "enter") // past Active and Review to Done
	for cmd != nil {
		cmd = deliver(m, cmd)
	}
	press(m, "v")
	if screen = plain(m); len(f.histories) != before || !strings.Contains(screen, "uncommitted       done       added in this checkout's files") ||
		!strings.Contains(screen, "No commit of this checkout holds the record yet.") {
		t.Fatalf("an added record: reads %v\n%s", f.histories[before:], screen)
	}
}

// A history read never makes a key wait: whatever comes next replaces it.
func TestHistoryReadYieldsToEveryKey(t *testing.T) {
	t.Parallel()
	_, f := lineageFixture()
	var cancelled []error
	block := true
	f.history = func(ctx context.Context, commit, _ string) ([]versions.Commit, error) {
		if block {
			<-ctx.Done()
			cancelled = append(cancelled, ctx.Err())
			return nil, ctx.Err()
		}
		return mainLog, nil
	}
	m := open(t, f, 120, 30)
	run := func(cmd tea.Cmd) chan tea.Msg {
		reply := make(chan tea.Msg, 1)
		go func() { reply <- cmd() }()
		return reply
	}
	// Moving to another version's history cancels the first and ignores its reply.
	first := run(press(m, "right", "enter"))
	press(m, "v")
	second := press(m, "down", "down")
	if second == nil || m.pending != "history" {
		t.Fatal("moving to another history should start its read")
	}
	if _, cmd := m.Update(<-first); cmd != nil || len(m.hist) != 0 || len(cancelled) != 1 || !errors.Is(cancelled[0], context.Canceled) {
		t.Fatalf("the superseded read: cancelled %v, held %d", cancelled, len(m.hist))
	}
	// A workspace selection replaces a history read rather than waiting for it.
	reply := run(second)
	press(m, "enter", "down") // open the fold; its first place shares the history being read
	f.refuse = errors.New("the branch moved")
	resolve := press(m, "enter")
	if resolve == nil || m.pending != "resolve" || m.notice != "" {
		t.Fatalf("selection during a history read: pending %q notice %q", m.pending, m.notice)
	}
	m.Update(<-reply)
	if len(cancelled) != 2 || len(m.hist) != 0 {
		t.Fatal("the history read was not cancelled by the selection")
	}
	// The refusal returns to the card, which still wants its history.
	again := deliver(m, resolve)
	if again == nil || m.pending != "history" || !strings.Contains(plain(m), "REFUSED: the branch moved") {
		t.Fatalf("after a refusal: pending %q\n%s", m.pending, plain(m))
	}
	// So does a refresh, and Esc, and quitting.
	reply = run(again)
	refresh := press(m, "r")
	if m.Update(<-reply); refresh == nil || m.pending != "inspect" || m.notice != "" || len(cancelled) != 3 {
		t.Fatalf("refresh during a history read: pending %q notice %q", m.pending, m.notice)
	}
	for i, key := range []string{"esc", "q"} {
		history := press(m, "enter") // the card again, after Esc
		if i == 0 {
			history = deliver(m, refresh) // the refreshed card
		}
		reply = run(history)
		if key == "esc" {
			press(m, "esc") // the versions close onto the detail, which asks for its history again
		}
		end := press(m, key)
		if _, cmd := m.Update(<-reply); cmd != nil || m.pending != "" || len(cancelled) != 4+i || len(m.hist) != 0 {
			t.Fatalf("%s during a history read: pending %q cancelled %d", key, m.pending, len(cancelled))
		}
		if key == "q" && (end == nil || end() != tea.QuitMsg{}) {
			t.Fatal("q must quit")
		}
	}
	m.reads.close() // every read was collected
	if len(f.resolved) != 1 || m.Workspace != nil {
		t.Fatalf("resolved %v", f.resolved)
	}
}
