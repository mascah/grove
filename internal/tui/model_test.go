package tui

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"testing"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/versions"
)

// --- fixtures: fake committed/live rows with conflicting statuses and titles.

func source(kind, locator, branch string) *versions.Source {
	s := &versions.Source{Kind: kind, Commit: strings.Repeat("a", 40), Present: true, Valid: true}
	if branch != "" {
		s.Ref = "refs/heads/" + branch
	}
	if kind == "live" {
		s.Locator, s.Worktree, s.GitDir = locator, "/repo/"+locator, "/repo/.git/"+locator
	}
	return s
}

func version(s *versions.Source, id, title, status string) versions.Version {
	kind := map[byte]string{'W': "work", 'Q': "question", 'D': "decision"}[id[0]]
	path := "grove/" + kind + "/" + id + ".md"
	body := fmt.Sprintf("---\nid: %s\ntitle: %s\nstatus: %s\n---\n\nBody of %s.\n", id, title, status, id)
	v := versions.Version{
		Source: s, Path: path, Revision: project.Revision([]byte(body)),
		Record:   &project.Record{ID: id, Type: kind, Title: title, Status: status, Path: path, Source: []byte(body)},
		Selector: s.Kind + ":" + s.Locator + ":" + s.Ref + ":" + id,
	}
	if s.Kind == "live" {
		v.Change = "unchanged"
	}
	return v
}

func result(here *versions.Source, sources []*versions.Source, vs ...versions.Version) *versions.Result {
	res := &versions.Result{Project: "/repo/.", Repository: "/repo/.git", Sources: sources, Complete: true}
	if here != nil {
		res.GitDir = here.GitDir
	}
	for _, s := range sources {
		if len(s.Diagnostics) != 0 {
			res.Complete = false
		}
	}
	for _, v := range vs {
		id := v.Path[strings.LastIndex(v.Path, "/")+1 : len(v.Path)-3]
		i := slices.IndexFunc(res.Groups, func(g versions.Group) bool { return g.ID == id })
		if i < 0 {
			res.Groups, i = append(res.Groups, versions.Group{ID: id}), len(res.Groups)
		}
		res.Groups[i].Versions = append(res.Groups[i].Versions, v)
	}
	return res
}

type fixture struct {
	cMain, cFeat, main, feat *versions.Source
}

func newFixture() fixture {
	return fixture{source("committed", "", "main"), source("committed", "", "feature"), source("live", ".", "main"), source("live", "feat", "feature")}
}

func (f fixture) sources() []*versions.Source {
	return []*versions.Source{f.cFeat, f.cMain, f.main, f.feat}
}

// twoBranches: W-001 is active on main and done, retitled, on feature; W-002
// is proposed everywhere; W-010 exists only on feature; Q-001 and D-001 are
// not work.
func (f fixture) twoBranches() *versions.Result {
	var vs []versions.Version
	for _, s := range []*versions.Source{f.cMain, f.main} {
		vs = append(vs, version(s, "W-001", "Inspect records", "active"), version(s, "W-002", "Create records", "proposed"),
			version(s, "Q-001", "Which version?", "resolved"), version(s, "D-001", "Sequential IDs", "accepted"))
	}
	for _, s := range []*versions.Source{f.cFeat, f.feat} {
		vs = append(vs, version(s, "W-001", "Inspect records, finished", "done"), version(s, "W-002", "Create records", "proposed"),
			version(s, "W-010", "Only on feature", "proposed"), version(s, "Q-001", "Which version?", "resolved"))
	}
	return result(f.main, f.sources(), vs...)
}

// --- a fake backend and a driver that runs commands by hand.

type fake struct {
	mu       sync.Mutex
	res      *versions.Result
	err      error
	ws       *versions.Workspace
	refuse   error
	inspects int
	resolved []string
	// history, when set, makes the backend offer History and logs its calls.
	history   func(ctx context.Context, commit, path string) ([]versions.Commit, error)
	histories []string
}

func (f *fake) backend() Backend {
	var history func(ctx context.Context, root, commit, path string) ([]versions.Commit, error)
	if f.history != nil {
		history = func(ctx context.Context, _, commit, path string) ([]versions.Commit, error) {
			f.mu.Lock()
			f.histories = append(f.histories, commit[:1]+" "+path)
			f.mu.Unlock()
			return f.history(ctx, commit, path)
		}
	}
	return Backend{
		History: history,
		Inspect: func(ctx context.Context, root, id string) (*versions.Result, error) {
			f.mu.Lock()
			defer f.mu.Unlock()
			f.inspects++
			return f.res, f.err
		},
		Resolve: func(ctx context.Context, root, selector string) (*versions.Workspace, error) {
			f.mu.Lock()
			defer f.mu.Unlock()
			f.resolved = append(f.resolved, selector)
			if f.refuse != nil {
				return nil, f.refuse
			}
			return f.ws, nil
		},
	}
}

func keyMsg(k string) tea.KeyPressMsg {
	named := map[string]tea.KeyPressMsg{
		"enter": {Code: tea.KeyEnter}, "esc": {Code: tea.KeyEscape}, "tab": {Code: tea.KeyTab},
		"up": {Code: tea.KeyUp}, "down": {Code: tea.KeyDown}, "left": {Code: tea.KeyLeft}, "right": {Code: tea.KeyRight},
		"pgup": {Code: tea.KeyPgUp}, "pgdown": {Code: tea.KeyPgDown}, "ctrl+c": {Code: 'c', Mod: tea.ModCtrl},
	}
	if msg, ok := named[k]; ok {
		return msg
	}
	return tea.KeyPressMsg{Code: rune(k[0]), Text: k}
}

// press sends keys and returns the last command without running it.
func press(m *Model, keys ...string) (cmd tea.Cmd) {
	for _, k := range keys {
		_, cmd = m.Update(keyMsg(k))
	}
	return cmd
}

// deliver runs a command and feeds its message back, returning what follows.
func deliver(m *Model, cmd tea.Cmd) tea.Cmd {
	if cmd == nil {
		return nil
	}
	_, next := m.Update(cmd())
	return next
}

func open(t *testing.T, f *fake, w, h int) *Model {
	t.Helper()
	m := New(t.Context(), "/repo/.", f.backend())
	m.Update(tea.WindowSizeMsg{Width: w, Height: h})
	if next := deliver(m, m.Init()); next != nil {
		t.Fatal("loading must not start another command")
	}
	return m
}

func plain(m *Model) string { return ansi.Strip(m.render()) }

// chooseCheckout shows the i-th live checkout's own board; the chooser lists
// the current view first.
func chooseCheckout(m *Model, i int) {
	press(m, "b")
	for range i + 1 {
		press(m, "down")
	}
	press(m, "enter")
}

func ids(cards []card) string {
	var out []string
	for _, c := range cards {
		out = append(out, c.id)
	}
	return strings.Join(out, ",")
}

// onRow reports a card's ID row holding text after the ID, inside the card's
// own box: what follows the ID up to the next border.
func onRow(screen, id, text string) bool {
	for _, r := range strings.Split(screen, "\n") {
		if i := strings.Index(r, id); i >= 0 {
			cell := r[i:]
			if j := strings.IndexAny(cell, "│┃"); j >= 0 {
				cell = cell[:j]
			}
			if strings.Contains(cell, text) {
				return true
			}
		}
	}
	return false
}

func board(m *Model) string {
	columns, shelf := m.cards()
	var parts []string
	for i, c := range columns {
		parts = append(parts, statuses[i]+"="+ids(c))
	}
	return strings.Join(parts, " ") + " shelf=" + ids(shelf)
}

// --- explicit selection

func TestSelectionIsExplicit(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := &fake{res: fx.twoBranches(), ws: &versions.Workspace{Project: "/repo/feat"}}
	m := open(t, f, 120, 30)
	if m.cardID != "W-002" || m.col != 0 {
		t.Fatalf("initial focus = column %d card %q, want the first proposed card", m.col, m.cardID)
	}
	press(m, "right") // W-001 is active on main
	if cmd := press(m, "enter"); cmd != nil || m.screen != detailScreen || m.openID() != "W-001" {
		t.Fatalf("Enter on a card must open its detail: screen %d open %q cmd %v", m.screen, m.openID(), cmd != nil)
	}
	if cmd := press(m, "v"); cmd != nil || m.screen != versionsScreen || m.verKey != "" || m.cardID != "W-001" {
		t.Fatalf("v must open the versions on the ID header: screen %d key %q cmd %v", m.screen, m.verKey, cmd != nil)
	}
	if cmd := press(m, "enter"); cmd != nil {
		t.Fatal("Enter on the ID header must not start anything")
	}
	if len(f.resolved) != 0 {
		t.Fatalf("resolved %v before any version was chosen", f.resolved)
	}
	// Rows follow the inspection's order. Two of them hold identical bytes
	// (committed and live main): they show as one fold, which selects nothing,
	// and opening it keeps them separate choices.
	g := m.group()
	if string(g.Versions[0].Record.Source) != string(g.Versions[1].Record.Source) || g.Versions[0].Selector == g.Versions[1].Selector {
		t.Fatal("fixture should hold identical content under different selectors")
	}
	if rows := m.rows(); len(rows) != 2 || len(rows[0].fold) != 2 || len(rows[1].fold) != 2 {
		t.Fatalf("four versions of two contents should be two folds: %+v", rows)
	}
	press(m, "down")
	if s := plain(m); !strings.Contains(s, "> ▸ active     same on 1 branch, 1 checkout") || strings.Contains(s, "checkout . (main)  unchanged") && !strings.Contains(s, "Same on") {
		t.Fatalf("identical content should be one row:\n%s", s)
	}
	if cmd := press(m, "enter"); cmd != nil || len(f.resolved) != 0 || len(m.rows()) != 4 {
		t.Fatalf("Enter on a fold must only list its places: %d rows, resolved %v", len(m.rows()), f.resolved)
	}
	if press(m, "enter"); len(m.rows()) != 2 {
		t.Fatal("Enter on an open fold should close it")
	}
	press(m, "enter", "down", "down")
	if m.verKey != g.Versions[1].Selector {
		t.Fatalf("cursor on %q, want the fold's second place", m.verKey)
	}
	cmd := press(m, "enter")
	if cmd == nil || m.pending != "resolve" {
		t.Fatal("Enter on a version must start one resolve")
	}
	// While that read is pending neither Enter nor r starts a second one.
	if press(m, "enter") != nil || press(m, "r") != nil {
		t.Fatal("a second operation started while one was pending")
	}
	next := deliver(m, cmd)
	if want := []string{g.Versions[1].Selector}; !slices.Equal(f.resolved, want) {
		t.Fatalf("resolved %v, want exactly %v", f.resolved, want)
	}
	if m.Workspace != f.ws || next == nil {
		t.Fatal("a fresh resolution must end the session with its workspace")
	}
	if _, ok := next().(tea.QuitMsg); !ok {
		t.Fatal("expected quit after resolution")
	}
	if f.inspects != 1 {
		t.Fatalf("inspected %d times, want 1", f.inspects)
	}
}

func TestRefusalStaysVisibleUntilRefresh(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := &fake{res: fx.twoBranches(), refuse: errors.New("branch refs/heads/main moved from aaa to bbb; run versions and reselect")}
	m := open(t, f, 120, 30)
	press(m, "right", "enter", "v", "down", "enter", "down")
	if next := deliver(m, press(m, "enter")); next != nil || m.Workspace != nil || m.screen != versionsScreen {
		t.Fatal("a refused resolution must stay in the version view without a workspace")
	}
	for _, k := range []string{"down", "up", "tab", "tab"} {
		press(m, k)
		if !strings.Contains(plain(m), "REFUSED: branch refs/heads/main moved") {
			t.Fatalf("refusal disappeared after %s:\n%s", k, plain(m))
		}
	}
	before := m.gen
	cmd := press(m, "r")
	if cmd == nil || m.gen == before || m.verKey != "" || m.refusal != "" || m.pending != "inspect" {
		t.Fatalf("refresh must clear the version focus and refusal and request an inspection: key %q refusal %q", m.verKey, m.refusal)
	}
	deliver(m, cmd)
	if f.inspects != 2 || m.screen != versionsScreen || m.verKey != "" {
		t.Fatal("after refresh the card's versions show again with nothing selected")
	}
}

func TestDeletedRowCannotResolve(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	res := fx.twoBranches()
	g := &res.Groups[slices.IndexFunc(res.Groups, func(g versions.Group) bool { return g.ID == "W-010" })]
	g.Versions = append(g.Versions, versions.Version{Source: fx.main, Path: "grove/work/W-010.md", Change: "deleted"})
	f := &fake{res: res}
	m := open(t, f, 120, 30)
	chooseCheckout(m, 0)
	if got := board(m); !strings.HasSuffix(got, "shelf=W-010") {
		t.Fatalf("work deleted from the board's live files belongs on the shelf: %s", got)
	}
	press(m, "tab", "enter", "v", "down", "down", "down")
	if v := m.focused(); v == nil || v.Record != nil {
		t.Fatal("the deleted row should stay visible and focusable")
	}
	if press(m, "enter") != nil || len(f.resolved) != 0 || !strings.Contains(plain(m), "REFUSED: this record was deleted") {
		t.Fatalf("a deleted row must not resolve:\n%s", plain(m))
	}
}

func TestStaleAndCancelledReplies(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := &fake{res: fx.twoBranches(), ws: &versions.Workspace{Project: "/repo/."}}
	m := open(t, f, 120, 30)
	press(m, "right", "enter", "v", "down", "enter", "down")
	old := press(m, "enter")
	if old == nil {
		t.Fatal("the fixture should start a resolve")
	}
	press(m, "esc", "esc") // leaves the versions and the detail, abandoning the resolve
	if m.pending != "" || m.screen != boardScreen {
		t.Fatal("Esc must cancel a pending resolve")
	}
	if next := deliver(m, old); next != nil || m.Workspace != nil {
		t.Fatal("an abandoned resolve selected a workspace")
	}
	// An older inspection must not replace newer state, in either order.
	cmd := press(m, "r")
	stale := *fx.twoBranches()
	stale.Groups = nil
	m.Update(inspectMsg{gen: m.gen - 1, res: &stale})
	if m.pending != "inspect" || len(m.res.Groups) == 0 {
		t.Fatal("an old-generation inspection was applied")
	}
	deliver(m, cmd)
	m.Update(inspectMsg{gen: m.gen - 1, res: &stale})
	m.Update(resolveMsg{gen: m.gen, selector: "never requested", ws: f.ws})
	if len(m.res.Groups) == 0 || m.Workspace != nil {
		t.Fatal("a late or unrequested reply changed state")
	}
	// Even a current, pending resolve accepts only the selector it asked for.
	press(m, "enter", "v", "down", "enter", "down") // still on W-001's column
	press(m, "enter")
	if m.pending != "resolve" || m.resolving == "" {
		t.Fatal("expected a pending resolve")
	}
	if _, cmd := m.Update(resolveMsg{gen: m.gen, selector: "another row", ws: f.ws}); cmd != nil || m.Workspace != nil || m.pending != "resolve" {
		t.Fatal("a reply for a different selector was accepted")
	}
}

func TestQuitAndInterruptCancelTheRead(t *testing.T) {
	t.Parallel()
	for key, want := range map[string]tea.Msg{"q": tea.QuitMsg{}, "ctrl+c": tea.InterruptMsg{}, "esc": tea.QuitMsg{}} {
		started, stopped := make(chan struct{}), make(chan error, 1)
		m := New(t.Context(), "/repo/.", Backend{Inspect: func(ctx context.Context, _, _ string) (*versions.Result, error) {
			close(started)
			<-ctx.Done()
			stopped <- ctx.Err()
			return nil, ctx.Err()
		}})
		cmd := m.Init()
		reply := make(chan tea.Msg, 1)
		go func() { reply <- cmd() }()
		<-started
		end := press(m, key)
		if end == nil || end() != want {
			t.Fatalf("%s: want %T", key, want)
		}
		if err := <-stopped; !errors.Is(err, context.Canceled) {
			t.Fatalf("%s: the read was not cancelled: %v", key, err)
		}
		m.Update(<-reply)
		m.reads.close() // returns because the read was collected
		if m.Workspace != nil || m.res != nil {
			t.Fatalf("%s: a cancelled read produced state", key)
		}
	}
}

func TestReadsNeverBeginAfterClose(t *testing.T) {
	t.Parallel()
	f := &fake{res: newFixture().twoBranches()}
	m := New(t.Context(), "/repo/.", f.backend())
	cmd := m.Init()
	m.reads.close()
	if msg := cmd(); msg != nil || f.inspects != 0 {
		t.Fatal("a command started after the session was collected")
	}
}

// --- board derivation

func TestBoardComesFromOneLiveSource(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := &fake{res: fx.twoBranches()}
	m := open(t, f, 120, 30)
	chooseCheckout(m, 0)
	if got, want := board(m), "proposed=W-002 active=W-001 review= done= abandoned= shelf=W-010"; got != want {
		t.Fatalf("main board:\n got %s\nwant %s", got, want)
	}
	screen := plain(m)
	for _, want := range []string{"Board: checkout . (main)", "read 2 branches, 2 checkouts", "Inspect records", "Elsewhere (1", "): W-010 ", "Done 0", "(none)"} {
		if !strings.Contains(screen, want) {
			t.Fatalf("board lacks %q:\n%s", want, screen)
		}
	}
	if !onRow(screen, "W-001", "2 versions") || onRow(screen, "W-002", "versions") {
		t.Fatalf("W-001 alone has differing versions:\n%s", screen)
	}
	for _, unwanted := range []string{"Q-001", "D-001", "Which version", "finished", "Only on feature", "W-010 ["} {
		if strings.Contains(screen, unwanted) {
			t.Fatalf("main's board shows %q:\n%s", unwanted, screen)
		}
	}
	// b lists live checkouts; choosing one changes only what is displayed.
	chooseCheckout(m, 1)
	if got, want := board(m), "proposed=W-002,W-010 active= review= done=W-001 abandoned= shelf="; got != want {
		t.Fatalf("feature board:\n got %s\nwant %s", got, want)
	}
	if screen = plain(m); !strings.Contains(screen, "Board: checkout feat (feature)") || !strings.Contains(screen, "Inspect records, finished") {
		t.Fatalf("the feature board should carry feature's label and titles:\n%s", screen)
	}
	if f.inspects != 1 || len(f.resolved) != 0 {
		t.Fatal("choosing a context must not call the backend")
	}
	// Each version keeps its own status in the card.
	press(m, "right", "right", "right", "enter", "v") // W-001 sits in Done, past Review
	screen = plain(m)
	for _, want := range []string{"W-001   2 versions differ", "(2 branches, 2 checkouts)", "▸ active     same on 1 branch, 1 checkout", "▸ done       same on 1 branch, 1 checkout", "Diverging: 2 current states"} {
		if !strings.Contains(screen, want) {
			t.Fatalf("versions lack %q:\n%s", want, screen)
		}
	}
}

func TestContextStates(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	empty := result(fx.main, fx.sources(), version(fx.main, "Q-001", "Which version?", "resolved"))
	m := open(t, &fake{res: empty}, 120, 30)
	if s := plain(m); !strings.Contains(s, "Proposed 0") || strings.Contains(s, "No board") {
		t.Fatalf("a valid checkout without work is an empty board:\n%s", s)
	}

	// A chosen checkout that stops validating is not an empty board.
	bad := newFixture()
	bad.main.Valid = false
	bad.main.Diagnostics = []string{"grove/work/W-001.md: status: expected one of proposed, active"}
	res := result(bad.main, bad.sources(), version(bad.feat, "W-010", "Only on feature", "proposed"))
	f := &fake{res: fx.twoBranches()}
	m = open(t, f, 120, 30)
	chooseCheckout(m, 0)
	f.res = res
	deliver(m, press(m, "r"))
	s := plain(m)
	for _, want := range []string{"[INVALID]", "No board: checkout . (main) is not usable", "not an empty board", "expected one of proposed", "INCOMPLETE: 1 of 4", "): W-010"} {
		if !strings.Contains(s, want) {
			t.Fatalf("invalid context lacks %q:\n%s", want, s)
		}
	}
	if strings.Contains(s, "Proposed ") {
		t.Fatalf("an invalid context must not draw columns:\n%s", s)
	}
	// The valid subset stays selectable, and the invalid source cannot be chosen.
	press(m, "b", "down", "enter")
	if m.screen != chooserScreen || !strings.Contains(plain(m), "cannot fill the board") {
		t.Fatalf("an invalid checkout was accepted:\n%s", plain(m))
	}
	press(m, "down", "enter")
	if got := board(m); got != "proposed=W-010 active= review= done= abandoned= shelf=" || !strings.Contains(plain(m), "INCOMPLETE") {
		t.Fatalf("valid subset: %s\n%s", got, plain(m))
	}

	// The current view still shows the valid sources, with the warning.
	press(m, "b", "enter")
	if got := board(m); got != "proposed=W-010 active= review= done= abandoned= shelf=" || !strings.Contains(plain(m), "Board: current view") || !strings.Contains(plain(m), "INCOMPLETE") {
		t.Fatalf("current view of a valid subset: %s\n%s", got, plain(m))
	}

	// Detached checkouts are contexts too, named by commit.
	det := newFixture()
	det.main.Ref = ""
	m = open(t, &fake{res: result(det.main, det.sources(), version(det.main, "W-001", "Detached work", "active"))}, 120, 30)
	chooseCheckout(m, 0)
	if s := plain(m); !strings.Contains(s, "Board: checkout . (detached at aaaaaaaaaaaa)") || !strings.Contains(s, "Detached work") {
		t.Fatalf("detached context:\n%s", s)
	}

	// The invoking checkout plays no part in the current view, even when it
	// is not a source at all.
	none := newFixture()
	m = open(t, &fake{res: result(nil, none.sources(), version(none.feat, "W-010", "Only on feature", "proposed"))}, 120, 30)
	if s := plain(m); !strings.Contains(s, "Board: current view") || board(m) != "proposed=W-010 active= review= done= abandoned= shelf=" {
		t.Fatalf("current view without an invoking checkout:\n%s", s)
	}

	// A failed inventory is not an empty board and offers nothing to select.
	f = &fake{err: errors.New("git worktree: fatal: not a git repository")}
	m = open(t, f, 120, 30)
	if s := plain(m); !strings.Contains(s, "could not be listed") || !strings.Contains(s, "fatal: not a git repository") || !strings.Contains(s, "r retry") {
		t.Fatalf("inventory failure:\n%s", s)
	}
	press(m, "enter", "b", "s", "tab")
	if m.screen != boardScreen || len(f.resolved) != 0 {
		t.Fatal("nothing may be selectable after an inventory failure")
	}
	f.err, f.res = nil, fx.twoBranches()
	if deliver(m, press(m, "r")); m.res == nil || m.failure != "" {
		t.Fatal("retry should recover")
	}
}

func TestRefreshFollowsIdentityNotPosition(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := &fake{res: fx.twoBranches()}
	m := open(t, f, 120, 30)
	chooseCheckout(m, 0)
	press(m, "right") // W-001, active
	// W-001 moves to done and a new card takes its old place.
	next := newFixture()
	f.res = result(next.main, next.sources(), version(next.main, "W-001", "Inspect records", "done"), version(next.main, "W-003", "Newcomer", "active"))
	deliver(m, press(m, "r"))
	if m.cardID != "W-001" || m.col != 3 {
		t.Fatalf("focus went to column %d card %q instead of following W-001", m.col, m.cardID)
	}
	// The open card disappears entirely.
	press(m, "enter", "v", "down")
	f.res = result(next.main, next.sources(), version(next.main, "W-003", "Newcomer", "active"))
	deliver(m, press(m, "r"))
	if m.screen != boardScreen || m.verKey != "" || !strings.Contains(plain(m), "W-001 is no longer on any readable branch or checkout") {
		t.Fatalf("a vanished card must return to the board with a reason:\n%s", plain(m))
	}
	// The context switches branch: the old choice is void until b chooses again.
	moved := newFixture()
	moved.main.Ref = "refs/heads/other"
	f.res = result(moved.main, moved.sources(), version(moved.main, "W-003", "Newcomer", "active"))
	deliver(m, press(m, "r"))
	if s := plain(m); m.hasBoard || !strings.Contains(s, "changed branch, moved, or was removed") || strings.Contains(s, "Active ") {
		t.Fatalf("a changed context identity must require a new choice:\n%s", s)
	}
	if deliver(m, press(m, "r")); m.hasBoard {
		t.Fatal("a lost context must not be re-adopted silently")
	}
	chooseCheckout(m, 0)
	if !m.hasBoard || board(m) != "proposed= active=W-003 review= done= abandoned= shelf=" {
		t.Fatalf("after choosing again: %s", board(m))
	}
}

// --- rendering

func TestLayoutAtEverySize(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	for _, size := range [][2]int{{120, 30}, {100, 24}, {99, 24}, {80, 24}, {40, 10}} {
		w, h := size[0], size[1]
		m := open(t, &fake{res: fx.twoBranches(), refuse: errors.New("stale\nselection")}, w, h)
		check := func(where string) {
			t.Helper()
			rows := strings.Split(m.render(), "\n")
			if len(rows) != h {
				t.Fatalf("%dx%d %s: %d rows:\n%s", w, h, where, len(rows), plain(m))
			}
			for i, r := range rows {
				if got := ansi.StringWidth(r); got != w {
					t.Fatalf("%dx%d %s: row %d is %d cells: %q", w, h, where, i, got, ansi.Strip(r))
				}
			}
			if last := ansi.Strip(rows[h-1]); !strings.Contains(last, "q") || strings.Contains(last, "…") {
				t.Fatalf("%dx%d %s: key hints must fit and keep quit: %q", w, h, where, last)
			}
		}
		check("board")
		s := plain(m)
		if wide := strings.Contains(s, "Proposed 2") && strings.Contains(s, "Done 0") && !strings.Contains(s, "[Proposed"); wide != (w >= wideWidth) {
			t.Fatalf("%dx%d: columns side by side = %v:\n%s", w, h, wide, s)
		}
		if strings.Count(s, "Aban") != 1 || !strings.Contains(s, "Abandoned 0 hidden") {
			t.Fatalf("%dx%d: Abandoned is hidden until asked for:\n%s", w, h, s)
		}
		press(m, "a")
		check("abandoned shown")
		if s := plain(m); !(strings.Contains(s, "Abandoned 0") || strings.Contains(s, "Aban 0")) || strings.Contains(s, "hidden") {
			t.Fatalf("%dx%d: a should add the Abandoned column:\n%s", w, h, s)
		}
		press(m, "a")
		if w < wideWidth {
			tabs := "[Proposed 2] Active 1 Review 0 Done 0"
			if w < 60 {
				tabs = "[Prop 2] Act 1 Rev 0 Done 0"
			}
			if !strings.Contains(s, tabs) || strings.Contains(s, "Inspect records") {
				t.Fatalf("%dx%d: want tabs %q and only the focused column:\n%s", w, h, tabs, s)
			}
			if press(m, "right"); !strings.Contains(plain(m), "Inspect records") {
				t.Fatalf("%dx%d: the next tab should show its cards:\n%s", w, h, plain(m))
			}
			press(m, "left")
		}
		press(m, "tab")
		check("shelf")
		press(m, "tab", "enter")
		check("detail")
		press(m, "tab")
		check("detail sidebar")
		press(m, "v", "down")
		check("versions")
		press(m, "enter", "down")
		check("open fold")
		deliver(m, press(m, "enter"))
		check("refusal")
		if !strings.Contains(plain(m), "REFUSED: stale") {
			t.Fatalf("%dx%d: refusal hidden:\n%s", w, h, plain(m))
		}
		press(m, "tab")
		check("details")
		if s := plain(m); !strings.Contains(s, "> Details") || (w < wideWidth && strings.Contains(s, "same everywhere")) {
			t.Fatalf("%dx%d: narrow terminals show only the focused pane:\n%s", w, h, s)
		}
		press(m, "s")
		check("sources")
		press(m, "esc", "esc", "esc", "b")
		check("chooser")
	}
	m := open(t, &fake{res: fx.twoBranches()}, 39, 30)
	if s := plain(m); !strings.Contains(s, "Grove needs 40x10; this is 39x30") || press(m, "q") == nil {
		t.Fatalf("below the minimum:\n%s", s)
	}
	m.Update(tea.WindowSizeMsg{Width: 120, Height: 9})
	if s := plain(m); !strings.Contains(s, "this is 120x9") {
		t.Fatalf("below the minimum height:\n%s", s)
	}
	m.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
	if s := plain(m); !strings.Contains(s, "Proposed 2") {
		t.Fatalf("resizing back should redraw the board:\n%s", s)
	}
}

func TestEverythingStaysReachable(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	var vs []versions.Version
	for i := 1; i <= 30; i++ {
		vs = append(vs, version(fx.main, fmt.Sprintf("W-%03d", i), fmt.Sprintf("Card %d", i), "proposed"))
		vs = append(vs, version(fx.feat, fmt.Sprintf("W-%03d", 100+i), "Elsewhere", "active"))
	}
	for i := range 25 {
		s := source("committed", "", fmt.Sprintf("branch-%02d", i))
		fx2 := version(s, "W-001", "Card 1", "done")
		vs = append(vs, fx2)
	}
	res := result(fx.main, fx.sources(), vs...)
	res.Groups[0].Versions[0].Record.Source = []byte(strings.Repeat("filler line\n", 200) + "THE LAST LINE\n")
	fx.cFeat.Diagnostics = make([]string, 60)
	for i := range fx.cFeat.Diagnostics {
		fx.cFeat.Diagnostics[i] = fmt.Sprintf("diagnostic %d", i)
	}
	fx.cFeat.Valid, res.Complete = false, false

	for _, size := range [][2]int{{120, 30}, {80, 24}, {40, 10}} {
		m := open(t, &fake{res: res}, size[0], size[1])
		chooseCheckout(m, 0) // the shelf of a checkout's board
		focusedRow := func(id string) bool {
			for _, r := range strings.Split(plain(m), "\n") {
				if strings.Contains(r, "> "+id) || strings.Contains(r, "▶"+id) {
					return true
				}
			}
			return false
		}
		for i := 1; i <= 30; i++ {
			if id := fmt.Sprintf("W-%03d", i); m.cardID != id || !focusedRow(id) {
				t.Fatalf("%v: card %s not focused and visible (focus %s):\n%s", size, id, m.cardID, plain(m))
			}
			press(m, "down")
		}
		press(m, "tab")
		for i := 1; i <= 30; i++ {
			if id := fmt.Sprintf("W-%03d", 100+i); !focusedRow(id) {
				t.Fatalf("%v: shelf item %s not visible:\n%s", size, id, plain(m))
			}
			press(m, "down")
		}
		press(m, "tab", "enter", "v")
		// main's own content, then one fold of 25 branches, opened.
		if press(m, "down"); !focusedRow("  proposed   checkout . (main)") {
			t.Fatalf("%v: the differing version is not visible:\n%s", size, plain(m))
		}
		if press(m, "down"); !focusedRow("▸ done       same on 25 branches") {
			t.Fatalf("%v: identical versions should fold into one row:\n%s", size, plain(m))
		}
		press(m, "enter")
		for i := range 25 {
			press(m, "down")
			if want := fmt.Sprintf("      branch branch-%02d", i); !focusedRow(want) {
				t.Fatalf("%v: place %d (%s) not focused and visible:\n%s", size, i, want, plain(m))
			}
		}
		press(m, "esc", "v", "down", "tab")
		for range 300 {
			press(m, "pgdown")
		}
		if !strings.Contains(plain(m), "THE LAST LINE") {
			t.Fatalf("%v: the end of a large body is unreachable:\n%s", size, plain(m))
		}
		if press(m, "pgup"); strings.Contains(plain(m), "THE LAST LINE") {
			t.Fatalf("%v: scrolling back should move", size)
		}
		press(m, "s")
		for range 100 {
			press(m, "down")
		}
		if !strings.Contains(plain(m), "diagnostic 59") || !strings.Contains(plain(m), "INCOMPLETE") {
			t.Fatalf("%v: the last diagnostic is unreachable:\n%s", size, plain(m))
		}
	}
}

func TestHostileTextIsInert(t *testing.T) {
	t.Parallel()
	const osc, erase, c1 = "\x1b]52;c;aGk=\x07", "\x1b[2J\x1b[H", "\u009b31m"
	hostile := "T" + osc + erase + c1 + "\x9b\r\u202e日本語e\u0301\ttab"
	fx := newFixture()
	fx.feat.Worktree = "/repo/new\nline\t" + osc
	fx.feat.Note = "note " + erase
	fx.cFeat.Valid, fx.cFeat.Diagnostics = false, []string{"bad " + hostile}
	v := version(fx.main, "W-001", hostile, "active")
	v.Path, v.HeadPath = "grove/work/W-001"+osc+"\n.md", "old"+erase
	body := "---\ntitle: x\n---\n" + hostile + "\nsecond line 日本語 e\u0301\n"
	v.Record.Source = []byte(body)
	v.Path = "grove/work/W-001.md"
	w := version(fx.feat, "W-001", "plain", "done")
	res := result(fx.main, fx.sources(), v, w)
	res.Groups[0].Versions[0].Path = "grove/work/W-001" + osc + "\n.md"
	res.Project = "/repo/" + erase
	f := &fake{res: res, refuse: errors.New("refused " + hostile)}

	for _, size := range [][2]int{{120, 30}, {80, 24}, {40, 10}} {
		m := New(t.Context(), "/repo/"+erase+osc, f.backend())
		m.Update(tea.WindowSizeMsg{Width: size[0], Height: size[1]})
		deliver(m, m.Init())
		check := func(where string) {
			t.Helper()
			out := sgr.ReplaceAllString(m.render(), "") // the interface's own styles
			if !utf8.ValidString(out) {
				t.Fatalf("%v %s: invalid UTF-8 reached the screen", size, where)
			}
			for _, r := range out {
				if r != '\n' && (r < 0x20 || r == 0x7f || (r >= 0x80 && r <= 0x9f) || r == '\u202e') {
					t.Fatalf("%v %s: control %U reached the screen:\n%q", size, where, r, out)
				}
			}
			for i, row := range strings.Split(out, "\n") {
				if got := ansi.StringWidth(row); got != size[0] {
					t.Fatalf("%v %s: row %d is %d cells: %q", size, where, i, got, row)
				}
			}
		}
		check("board")
		press(m, "right", "enter")
		check("detail")
		press(m, "tab")
		check("detail sidebar")
		press(m, "v", "down")
		check("versions")
		deliver(m, press(m, "enter"))
		check("refusal")
		press(m, "tab")
		check("details")
		if size[0] >= 80 {
			s := plain(m)
			for _, want := range []string{`\x1b]52;c;aGk=\a`, `\x1b[2J`, `\u009b31m`, `\x9b`, `\r`, `\u202e`, "日本語e\u0301"} {
				if !strings.Contains(s, want) {
					t.Fatalf("%v: details should show %q visibly:\n%s", size, want, s)
				}
			}
		}
		press(m, "s")
		check("sources")
		press(m, "esc", "esc", "esc", "b")
		check("chooser")
		m.notice = hostile
		check("notice")
	}
	if string(v.Record.Source) != body {
		t.Fatal("display escaping changed the record's bytes")
	}
	// Line structure of a source is kept; other controls are shown, not obeyed.
	rows := wrapAll("a\tb\r\nc\x1b[2J", 20)
	if len(rows) != 2 || strings.TrimSpace(rows[0]) != `a    b\r` || strings.TrimSpace(rows[1]) != `c\x1b[2J` {
		t.Fatalf("wrapAll = %q", rows)
	}
	// Clipping counts display cells after escaping.
	if got := line("日本語テキスト", 5); ansi.StringWidth(got) != 5 || !strings.HasPrefix(got, "日本…") {
		t.Fatalf("line = %q", got)
	}
	if got := line("e\u0301e\u0301e\u0301", 3); got != "e\u0301e\u0301e\u0301" {
		t.Fatalf("combining marks take no cells: %q", got)
	}
}

// Review regressions: a refresh under an open chooser or sources screen.
func TestRefreshUnderOverlays(t *testing.T) {
	t.Parallel()
	live := func(n int) *versions.Result {
		var sources []*versions.Source
		var vs []versions.Version
		for i := range n {
			s := source("live", fmt.Sprintf("wt%d", i), fmt.Sprintf("b%d", i))
			sources, vs = append(sources, s), append(vs, version(s, "W-001", "Card", "proposed"))
		}
		return result(sources[0], sources, vs...)
	}
	for _, left := range []int{4, 2, 1} {
		f := &fake{res: live(6)}
		m := open(t, f, 40, 10)
		press(m, "b", "down", "down", "down", "down", "down")
		f.res = live(left)
		deliver(m, press(m, "r"))
		if s := plain(m); m.choice != left || !strings.Contains(s, fmt.Sprintf("> checkout wt%d", left-1)) {
			t.Fatalf("%d checkouts left: choice %d:\n%s", left, m.choice, s)
		}
	}

	// The open card vanishes while the sources screen covers its versions.
	fx := newFixture()
	f := &fake{res: fx.twoBranches()}
	m := open(t, f, 120, 30)
	press(m, "right", "enter", "v", "down", "s")
	f.res = result(fx.main, fx.sources(), version(fx.main, "W-003", "Newcomer", "active"))
	deliver(m, press(m, "r"))
	if !strings.Contains(plain(m), "W-001 is no longer on any readable branch or checkout") {
		t.Fatalf("the reason was lost:\n%s", plain(m))
	}
	if press(m, "esc"); m.screen != boardScreen || m.cardID != "W-003" {
		t.Fatalf("Esc should return to the board, not another card's versions: screen %d", m.screen)
	}

	// The incomplete warning survives a pending read at the narrowest size.
	bad := newFixture()
	bad.cFeat.Valid, bad.cFeat.Diagnostics = false, []string{"broken"}
	m = open(t, &fake{res: result(bad.main, bad.sources(), version(bad.main, "W-001", "Card", "active"))}, 40, 10)
	press(m, "right", "enter", "v", "down")
	if press(m, "enter"); !strings.Contains(plain(m), "INCOMPLETE: 1 of 4") {
		t.Fatalf("the warning was clipped by the pending read:\n%s", plain(m))
	}
}

// Nothing ties an ID to a type, so the board follows the type
// field alone: neutral-ID work is a card, and a page, which has no status,
// never is, whatever its ID or wherever it sits.
func TestBoardFollowsTypeNotIDOrPlacement(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	retype := func(v versions.Version, kind, status string) versions.Version {
		v.Record.Type, v.Record.Status = kind, status
		return v
	}
	var vs []versions.Version
	for _, s := range []*versions.Source{fx.cMain, fx.main} {
		vs = append(vs,
			retype(version(s, "G-001", "Neutral work", "active"), "work", "active"),
			retype(version(s, "G-002", "A page about proposed work", ""), "page", ""),
			retype(version(s, "W-003", "A page under an old work ID", ""), "page", ""),
			retype(version(s, "D-004", "Reclassified into work", "proposed"), "work", "proposed"))
	}
	// Sources may disagree about a type: the board source's own record decides,
	// and committed main, which sorts first, must not. G-060 became work in the
	// live checkout; G-061 became a page there; G-062 is work only on feature.
	vs = append(vs,
		retype(version(fx.cMain, "G-060", "Was a page", ""), "page", ""), retype(version(fx.main, "G-060", "Now work", "active"), "work", "active"),
		retype(version(fx.cMain, "G-061", "Was work", "done"), "work", "done"), retype(version(fx.main, "G-061", "Now a page", ""), "page", ""),
		retype(version(fx.cFeat, "G-062", "Work elsewhere", "proposed"), "work", "proposed"))
	// G-063 is deleted everywhere it is seen, so no record says it was work.
	vs = append(vs, versions.Version{Source: fx.main, Path: "grove/G-063.md", Change: "deleted"})
	m := open(t, &fake{res: result(fx.main, fx.sources(), vs...)}, 120, 30)
	chooseCheckout(m, 0)
	if got, want := board(m), "proposed=D-004 active=G-001,G-060 review= done= abandoned= shelf=G-062"; got != want {
		t.Fatalf("board:\n got %s\nwant %s", got, want)
	}
	if screen := plain(m); strings.Contains(screen, "page") || strings.Contains(screen, "G-002") || strings.Contains(screen, "W-003") || strings.Contains(screen, "G-061") || strings.Contains(screen, "G-063") {
		t.Fatalf("a page must not appear on the board:\n%s", screen)
	}
}

// TestCurrentViewBoard: the board opens on the current view (G-042), the same
// whichever checkout invoked it. W-001 is done on feature, which main's older
// active copy does not obscure; W-002 diverges, so it is one marked card in
// the earlier of its statuses; W-003's current state is an uncommitted edit;
// W-004's is a deletion on feature.
func TestCurrentViewBoard(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	older := func(v versions.Version, why string) versions.Version { v.Older = why; return v }
	var vs []versions.Version
	for _, s := range []*versions.Source{fx.cMain, fx.main} {
		vs = append(vs, older(version(s, "W-001", "Inspect records", "active"), "branch feature changed it since their common history"),
			version(s, "W-002", "Create records", "proposed"),
			older(version(s, "W-003", "Edit records", "proposed"), "checkout feat (feature) has an uncommitted change to it on top of this commit"),
			older(version(s, "W-004", "Drop records", "proposed"), "branch feature changed it since their common history"),
			older(version(s, "W-005", "Shelve records", "proposed"), "checkout feat (feature) has an uncommitted change to it on top of this commit"))
	}
	vs = append(vs, version(fx.cFeat, "W-001", "Inspect records, finished", "done"), version(fx.feat, "W-001", "Inspect records, finished", "done"),
		version(fx.cFeat, "W-002", "Create records, started", "active"), version(fx.feat, "W-002", "Create records, started", "active"),
		older(version(fx.cFeat, "W-003", "Edit records", "proposed"), "checkout feat (feature) has an uncommitted change to it on top of this commit"))
	edited := version(fx.feat, "W-003", "Edit records", "active")
	edited.Change, edited.Record.Source, edited.Revision = "modified", []byte("edited"), "sha256:edited"
	vs = append(vs, edited,
		older(version(fx.cFeat, "W-005", "Shelve records", "proposed"), "checkout feat (feature) has an uncommitted change to it on top of this commit"),
		versions.Version{Source: fx.feat, Path: "grove/work/W-005.md", Change: "deleted"})
	res := result(fx.main, fx.sources(), vs...)
	res.Groups[1].Notes = []string{"branch main and branch feature could not be ordered: example"}
	// The projection's row for a branch's deletion has no path; a checkout's
	// has its HEAD's.
	res.Groups[3].Versions = append(res.Groups[3].Versions, versions.Version{Source: fx.cFeat, Change: "deleted"})

	m := open(t, &fake{res: res}, 120, 30)
	want := "proposed=W-002 active=W-003 review= done=W-001 abandoned= shelf=W-004,W-005"
	if got := board(m); got != want {
		t.Fatalf("current view:\n got %s\nwant %s", got, want)
	}
	screen := plain(m)
	for _, want := range []string{"Board: current view", "Create records", "Inspect records, finished", "Deleted (2, the current state removes the record; Tab): W-004  W-005 [uncommitted]"} {
		if !strings.Contains(screen, want) {
			t.Fatalf("board lacks %q:\n%s", want, screen)
		}
	}
	if !onRow(screen, "W-002", "⑂ 2 states") || !onRow(screen, "W-003", "uncommitted") {
		t.Fatalf("cards lack their tags:\n%s", screen)
	}
	other := *res
	other.GitDir = fx.feat.GitDir
	if got := board(open(t, &fake{res: &other}, 120, 30)); got != want {
		t.Fatalf("another invoking checkout sees %s", got)
	}

	// W-001's current state leads its versions; main's copy is marked older.
	press(m, "right", "right", "right", "enter", "v")
	screen = plain(m)
	text := currentText(m.group(), "") // the details wrap it
	for _, want := range []string{"▸ done       same on 1 branch, 1 checkout", "▸ active     older  same on 1 branch",
		`Current: done "Inspect records, finished" on branch feature, checkout feat (feature).`, "Older: 2 places hold an earlier state"} {
		if !strings.Contains(screen, want) && !strings.Contains(text, want) {
			t.Fatalf("W-001 lacks %q:\n%s", want, screen)
		}
	}
	if strings.Index(screen, "▸ done") > strings.Index(screen, "▸ active") {
		t.Fatalf("the current state should come first:\n%s", screen)
	}
	press(m, "down", "down")
	if screen = strings.Join(m.describeRows(200), "\n"); !strings.Contains(screen, "Standing: older: branch feature changed it since their common history") {
		t.Fatalf("an older row should say why:\n%s", screen)
	}
	// W-002 explains its divergence and the pair it could not order.
	press(m, "esc", "esc", "left", "left", "left", "enter")
	screen = currentText(m.group(), "")
	for _, want := range []string{"Diverging: 2 current states.", `- proposed "Create records" on branch main, checkout . (main)`,
		`- active "Create records, started" on branch feature, checkout feat (feature)`, "Could not order: branch main and branch feature could not be ordered: example"} {
		if !strings.Contains(screen, want) {
			t.Fatalf("W-002 lacks %q:\n%s", want, screen)
		}
	}
	// Only W-005's deletion is uncommitted.
	press(m, "esc", "tab")
	if screen = plain(m); !strings.Contains(screen, "W-005   uncommitted") || strings.Contains(screen, "W-004   uncommitted") {
		t.Fatalf("the shelf should mark W-005's uncommitted deletion alone:\n%s", screen)
	}
	// W-004's deletion on a branch has no history to read, and is a row that
	// opens nothing.
	h := &fake{res: res, history: func(context.Context, string, string) ([]versions.Commit, error) { return nil, nil }}
	withHistory := open(t, h, 120, 30)
	press(withHistory, "tab", "enter")
	if rows := strings.Join(withHistory.historyRows(200), "\n"); withHistory.wantHistory() != nil || !strings.Contains(rows, "This branch deleted the record") {
		t.Fatalf("a branch's deletion has no history:\n%s", rows)
	}
	press(m, "enter", "v", "down")
	if press(m, "enter") != nil || !strings.Contains(plain(m), "REFUSED: this record was deleted on that branch") {
		t.Fatalf("a branch's deletion must not resolve:\n%s", plain(m))
	}
}

// TestCurrentViewTarget: with a target, a card none of whose current states
// the target holds is marked, unless its only current state is uncommitted,
// and the details say which state is on the target. A divergence that main
// holds one side of is not marked; a branch's deletion of what main holds is.
func TestCurrentViewTarget(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	older := func(v versions.Version) versions.Version {
		v.Older = "branch feature changed it since their common history"
		return v
	}
	var vs []versions.Version
	for _, s := range []*versions.Source{fx.cMain, fx.main} {
		vs = append(vs, older(version(s, "W-001", "Inspect records", "proposed")), version(s, "W-002", "Create records", "proposed"),
			older(version(s, "W-003", "Edit records", "proposed")), version(s, "W-004", "Split records", "proposed"),
			older(version(s, "W-005", "Drop records", "proposed")))
	}
	vs = append(vs, version(fx.cFeat, "W-004", "Split records", "active"), version(fx.feat, "W-004", "Split records", "active"))
	vs = append(vs, version(fx.cFeat, "W-001", "Inspect records", "active"), version(fx.feat, "W-001", "Inspect records", "active"),
		version(fx.cFeat, "W-002", "Create records", "proposed"), version(fx.feat, "W-002", "Create records", "proposed"),
		older(version(fx.cFeat, "W-003", "Edit records", "proposed")))
	edited := version(fx.feat, "W-003", "Edit records", "active")
	edited.Change, edited.Revision = "modified", "sha256:edited"
	vs = append(vs, edited)
	res := result(fx.main, fx.sources(), vs...)
	res.Target = "main"
	i := slices.IndexFunc(res.Groups, func(g versions.Group) bool { return g.ID == "W-005" })
	res.Groups[i].Versions = append(res.Groups[i].Versions, versions.Version{Source: fx.cFeat, Change: "deleted"})
	for _, g := range res.Groups { // as the projection marks them
		i := slices.IndexFunc(g.Versions, func(v versions.Version) bool { return v.Source == fx.cMain })
		for j := range g.Versions {
			g.Versions[j].OnTarget = g.Versions[j].Revision == g.Versions[i].Revision
		}
	}

	m := open(t, &fake{res: res}, 120, 30)
	screen := plain(m)
	for _, want := range []string{"Board: current view, target main", "Tab): W-005 [not on main]"} {
		if !strings.Contains(screen, want) {
			t.Fatalf("board lacks %q:\n%s", want, screen)
		}
	}
	if !onRow(screen, "W-001", "not on main") || !onRow(screen, "W-003", "uncommitted") || !onRow(screen, "W-004", "⑂ 2 states") {
		t.Fatalf("cards lack their tags:\n%s", screen)
	}
	if onRow(screen, "W-002", "not on") || strings.Contains(screen, "uncommitted not on") || strings.Contains(screen, "states not on") {
		t.Fatalf("W-002 is on main, and W-003's state is uncommitted only:\n%s", screen)
	}
	press(m, "right", "enter")
	if text := currentText(m.group(), m.res.Target); !strings.Contains(text, `Current: active "Inspect records" on branch feature, checkout feat (feature), not on main.`) {
		t.Fatalf("W-001: %s", text)
	}
}
