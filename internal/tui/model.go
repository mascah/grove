// Package tui is Grove's terminal interface: a Kanban board of the project's
// current work across every branch and checkout (G-042), or of one checkout's
// files, whose cards open a record's differing versions with the branches and
// checkouts holding each, the focused one's history of commits, and explicit
// selection of one existing workspace. It reads through Backend, and writes
// only through Backend's three actions on a record in review (G-044), each
// behind a prompt: approve, feedback, and integrate.
package tui

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/versions"
)

// Backend is every effect the interface has besides drawing. The reads are
// the board's; the three actions are G-044's, each run in the checkout the
// review view names and returning the facts to show.
type Backend struct {
	Inspect func(ctx context.Context, root, id string) (*versions.Result, error)
	Resolve func(ctx context.Context, root, selector string) (*versions.Workspace, error)
	// History lists the commits behind an open card's focused version. It is
	// read when a card is open, never for the board; nil leaves the section out.
	History func(ctx context.Context, root, commit, path string) ([]versions.Commit, error)
	// Changes reads a candidate's files against the target and what the tip
	// changed after it, when a detail with a candidate is open; nil leaves
	// the section out. Diff reads one of those files.
	Changes func(ctx context.Context, root, target, candidate, tip, recordPath string) (*versions.Changes, error)
	Diff    func(ctx context.Context, root, from, to, path string) (string, error)
	// Approve, Feedback and Integrate write: in root, the checkout of the
	// branch judged, or of the target. Nil leaves the action out.
	Approve   func(ctx context.Context, root, id, verdict string) ([]string, error)
	Feedback  func(ctx context.Context, root, id, text string) ([]string, error)
	Integrate func(ctx context.Context, root, id string, cleanup bool) ([]string, error)
}

type screen int

const (
	boardScreen screen = iota
	detailScreen
	versionsScreen
	chooserScreen
	sourcesScreen
	searchScreen
	resultScreen
)

var statuses = [5]string{"proposed", "active", "review", "done", "abandoned"}

// sourceKey identifies the board's checkout across refreshes. A checkout that
// moved, was re-registered, or switched branch is a different context, which
// the person must choose again.
type sourceKey struct{ locator, worktree, ref string }

func keyOf(s *versions.Source) sourceKey { return sourceKey{s.Locator, s.Worktree, s.Ref} }

type card struct {
	id, title string
	versions  int    // distinct contents, not the places holding them
	tag       string // what the card notes beside its ID
	meta      string // kind, size and priority, or when Done and its candidate
	rec       *project.Record
}

// meta is a card's last line: what a glance at the board needs beyond the
// title. Done work shows when it was last written and its candidate.
func meta(r *project.Record) string {
	if r == nil {
		return ""
	}
	var parts []string
	if r.Status == "done" {
		if r.Updated != nil {
			parts = append(parts, "done "+r.Updated.Format("2006-01-02"))
		}
		if r.Candidate != "" {
			parts = append(parts, r.Candidate[:min(len(r.Candidate), 7)])
		}
		return strings.Join(parts, " · ")
	}
	for _, v := range []string{r.Kind, r.Size} {
		if v != "" {
			parts = append(parts, v)
		}
	}
	if r.Priority != nil {
		parts = append(parts, fmt.Sprintf("P%d", *r.Priority))
	}
	return strings.Join(parts, " · ")
}

// row is one line of a card's version list: a fold standing for several
// places that hold the same bytes, or one version. A fold selects nothing;
// Enter opens it into its members, which stay separate explicit choices.
type row struct {
	key    string
	fold   []*versions.Version // more than one place holds these bytes
	v      *versions.Version
	inFold bool // v is a member shown beneath its open fold
}

type inspectMsg struct {
	gen int
	res *versions.Result
	err error
}

type resolveMsg struct {
	gen      int
	selector string
	ws       *versions.Workspace
	err      error
}

type historyMsg struct {
	gen     int
	key     string
	commits []versions.Commit
	err     error
}

// lineage is one finished history read, kept until the next inspection.
type lineage struct {
	commits []versions.Commit
	err     error
}

// reads tracks backend calls so Run can collect them after the program ends.
// A command the runtime starts after close never begins.
type reads struct {
	mu     sync.Mutex
	closed bool
	wg     sync.WaitGroup
}

func (r *reads) begin() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.closed {
		r.wg.Add(1)
	}
	return !r.closed
}

func (r *reads) close() {
	r.mu.Lock()
	r.closed = true
	r.mu.Unlock()
	r.wg.Wait()
}

// Model is the whole interface state. Focus is held as identities (a record
// ID, a version's selector), never as a row position that a refresh could
// hand to a different record.
type Model struct {
	ctx     context.Context
	root    string
	backend Backend
	reads   reads

	width, height int

	res     *versions.Result
	failure string // the inventory itself failed: no rows, retry or quit
	notice  string // one-shot message, cleared by the next key

	gen       int    // the newest request; older replies are ignored
	pending   string // "", "inspect", "resolve", "history", "changes", "diff" or "act": one at a time
	resolving string // the exact selector a pending resolve was asked for
	reading   string // the key a pending history, changes or diff read was asked for
	acting    string // the running action, for the banner
	cancel    context.CancelFunc
	hist      map[string]lineage     // by commit and path, for the current result only
	changes   map[string]changesRead // by candidate, tip, target and path, likewise
	diffs     map[string]diffRead    // by base, candidate and path, likewise
	md        map[string][]string    // rendered Markdown by key and width, for the current result only
	done      bool                   // the session is ending: start nothing more

	board         sourceKey
	hasBoard      bool // a checkout's own board; otherwise the current view
	lost          bool // the chosen checkout changed identity; b must choose again
	showAbandoned bool // the Abandoned column is hidden until asked for

	screen, back screen
	col          int
	onShelf      bool
	cardID       string   // the board's focus
	stack        []string // open records, the last showing in the detail
	side         int      // the detail's sidebar cursor; -1 while the content has focus
	dscroll      int      // the detail's content scroll
	asOf         string   // a timeline commit whose content the detail shows; "" is now
	diff         string   // a changed file whose diff the detail shows; "" is the content
	prompt       *prompt  // the open question on the last row, if any
	result       *outcome // the last action's outcome, on the result screen
	query        string   // the search's text
	hit          int      // the search's cursor
	verKey       string   // a row's key; "" is the ID header, which selects nothing
	unfolded     string   // the key of the one fold showing its members
	detail       bool     // the detail pane has focus
	scroll       int      // detail pane, or sources screen
	choice       int      // chooser row
	refusal      string

	// Workspace is the explicitly selected, freshly resolved result, if any.
	Workspace *versions.Workspace
}

// New returns a model that starts by inspecting every record of root.
func New(ctx context.Context, root string, backend Backend) *Model {
	return &Model{ctx: ctx, root: root, backend: backend}
}

func (m *Model) Init() tea.Cmd { return m.inspect() }

// read starts the one allowed backend call under its own cancellable context,
// cancelling and outdating whichever came before.
func (m *Model) read(kind string, call func(ctx context.Context, gen int) tea.Msg) tea.Cmd {
	m.stop()
	gen := m.gen
	ctx, cancel := context.WithCancel(m.ctx)
	m.pending, m.cancel = kind, cancel
	return func() tea.Msg {
		defer cancel()
		if !m.reads.begin() {
			return nil
		}
		defer m.reads.wg.Done()
		return call(ctx, gen)
	}
}

func (m *Model) inspect() tea.Cmd {
	return m.read("inspect", func(ctx context.Context, gen int) tea.Msg {
		res, err := m.backend.Inspect(ctx, m.root, "")
		return inspectMsg{gen, res, err}
	})
}

func (m *Model) resolve(selector string) tea.Cmd {
	cmd := m.read("resolve", func(ctx context.Context, gen int) tea.Msg {
		ws, err := m.backend.Resolve(ctx, m.root, selector)
		return resolveMsg{gen, selector, ws, err}
	})
	m.resolving = selector
	return cmd
}

// busy reports a read or action that a key must wait for. A history,
// changes or diff read is not one: whatever the person asks for next
// replaces it.
func (m *Model) busy() bool {
	return m.pending == "inspect" || m.pending == "resolve" || m.pending == "act"
}

// wantHistory starts reading the focused version's history when a card is
// showing, no other read is pending, and that history is not already held or
// being read. The board never asks for one.
func (m *Model) wantHistory() tea.Cmd {
	if m.backend.History == nil || m.done || m.screen != versionsScreen && m.screen != detailScreen || m.pending != "" && m.pending != "history" {
		return nil
	}
	commit, path := historyAt(m.historyOf())
	key := commit + "\x00" + path
	if _, held := m.hist[key]; commit == "" || held || key == m.reading {
		return nil
	}
	cmd := m.read("history", func(ctx context.Context, gen int) tea.Msg {
		commits, err := m.backend.History(ctx, m.root, commit, path)
		return historyMsg{gen, key, commits, err}
	})
	m.reading = key
	return cmd
}

// historyOf returns the version whose history the details show: the focused
// row's, a fold's first place, and otherwise the group's shown version.
func (m *Model) historyOf() *versions.Version {
	g := m.group()
	if g == nil || len(g.Versions) == 0 {
		return nil
	}
	if r := m.focusedRow(); r != nil {
		if r.fold != nil {
			return r.fold[0]
		}
		return r.v
	}
	return m.shown(g)
}

// shown returns the version that stands for a group on this board: the
// board's checkout's version, the current view's first current record, a
// current deletion, or the group's first.
func (m *Model) shown(g *versions.Group) *versions.Version {
	if len(g.Versions) == 0 {
		return nil
	}
	for i := range g.Versions {
		if v := &g.Versions[i]; m.hasBoard && v.Source == m.boardSource() || m.current() && v.Older == "" && v.Record != nil {
			return v
		}
	}
	if i := slices.IndexFunc(g.Versions, func(v versions.Version) bool { return m.current() && v.Older == "" }); i >= 0 {
		return &g.Versions[i]
	}
	return &g.Versions[0]
}

// historyAt names the commit and the path there that a version's history
// starts from. A checkout's record is followed from its HEAD, under the name
// it has there; one added since HEAD, or deleted on a branch, has no commit
// to read.
func historyAt(v *versions.Version) (commit, path string) {
	if v == nil || v.Change == "added" || v.Path == "" {
		return "", ""
	}
	if v.HeadPath != "" {
		return v.Source.Commit, v.HeadPath
	}
	return v.Source.Commit, v.Path
}

// stop cancels any read in flight and outdates its reply.
func (m *Model) stop() {
	if m.cancel != nil {
		m.cancel()
	}
	m.gen++
	m.pending, m.resolving, m.reading, m.acting, m.cancel = "", "", "", "", nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := m.update(msg)
	// One read at a time, in this order: a detail's history, then its
	// changes, then a chosen diff; the next starts when the last delivered.
	if cmd == nil {
		cmd = m.wantHistory()
	}
	if cmd == nil {
		cmd = m.wantChanges()
	}
	if cmd == nil {
		cmd = m.wantDiff()
	}
	return m, cmd
}

func (m *Model) update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.settleFocus() // the Done bound may have moved past the focused card
		m.clampScroll()
	case inspectMsg:
		if msg.gen != m.gen || m.pending != "inspect" {
			return nil
		}
		m.pending, m.cancel, m.hist, m.md, m.asOf = "", nil, map[string]lineage{}, nil, ""
		m.changes, m.diffs, m.diff = map[string]changesRead{}, map[string]diffRead{}, ""
		if msg.err != nil {
			m.res, m.failure = nil, msg.err.Error()
			m.screen, m.cardID = boardScreen, ""
			m.leaveVersions()
			return nil
		}
		m.res, m.failure = msg.res, ""
		m.choice = min(m.choice, len(m.live()))
		m.settleBoard()
		m.settleFocus()
	case resolveMsg:
		if msg.gen != m.gen || m.pending != "resolve" || msg.selector != m.resolving {
			return nil
		}
		m.pending, m.resolving, m.cancel = "", "", nil
		if msg.err != nil || msg.ws == nil {
			m.refusal = "no workspace was returned"
			if msg.err != nil {
				m.refusal = msg.err.Error()
			}
			return nil
		}
		m.Workspace, m.done = msg.ws, true
		return tea.Quit
	case historyMsg:
		if msg.gen != m.gen || m.pending != "history" {
			return nil
		}
		m.pending, m.reading, m.cancel = "", "", nil
		m.hist[msg.key] = lineage{msg.commits, msg.err}
		m.clampScroll()
	case changesMsg:
		if msg.gen != m.gen || m.pending != "changes" {
			return nil
		}
		m.pending, m.reading, m.cancel = "", "", nil
		m.changes[msg.key] = changesRead{msg.c, msg.err}
		m.clampScroll()
	case diffMsg:
		if msg.gen != m.gen || m.pending != "diff" {
			return nil
		}
		m.pending, m.reading, m.cancel = "", "", nil
		m.diffs[msg.key] = diffRead{msg.text, msg.err}
		m.clampScroll()
	case actMsg:
		if msg.gen != m.gen || m.pending != "act" {
			return nil
		}
		m.pending, m.acting, m.cancel = "", "", nil
		title := map[string]string{"approve": "Approved " + m.openID(), "feedback": "Feedback recorded on " + m.openID(), "integrate": "Integration of " + m.openID()}[msg.kind]
		m.result = &outcome{title: title, facts: msg.facts}
		if msg.err != nil {
			m.result.err = msg.err.Error()
		}
		m.screen, m.scroll = resultScreen, 0
		// Whatever was written, the board is re-read; the outcome stays up.
		return m.inspect()
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m.key("ctrl+c")
		}
		m.notice = ""
		if m.pending == "act" {
			m.notice = actingText(m.acting) + " Keys wait for it."
			return nil
		}
		if m.prompt != nil {
			return m.promptKey(msg)
		}
		if m.screen == searchScreen {
			m.searchKey(msg)
			return nil
		}
		return m.key(msg.String())
	}
	return nil
}

func (m *Model) key(k string) tea.Cmd {
	m.notice = ""
	switch k {
	case "ctrl+c":
		m.stop()
		m.done = true
		return tea.Interrupt
	case "q":
		m.stop()
		m.done = true
		return tea.Quit
	case "r":
		if m.busy() {
			m.notice = "a read is already in progress"
			return nil
		}
		m.leaveVersions()
		return m.inspect()
	case "s":
		if m.screen != sourcesScreen && m.res != nil {
			m.back, m.screen, m.scroll = m.screen, sourcesScreen, 0
		}
		return nil
	case "esc":
		switch m.screen {
		case boardScreen:
			m.stop()
			m.done = true
			return tea.Quit
		case detailScreen:
			// The board reads no history, changes or diff; whatever the
			// detail returns to asks again for its own.
			if m.pending != "" && m.pending != "inspect" && m.pending != "act" {
				m.stop()
			}
			m.leaveDetail()
		case versionsScreen:
			if m.pending != "inspect" {
				m.stop()
			}
			m.screen = detailScreen
			m.leaveVersions()
		case resultScreen:
			if m.pending == "inspect" {
				return nil // the record shown next must be the re-read one
			}
			m.result, m.scroll = nil, 0
			if m.screen = boardScreen; len(m.stack) != 0 {
				m.screen = detailScreen
			}
		default:
			m.screen, m.scroll = m.back, 0
		}
		return nil
	}
	switch m.screen {
	case boardScreen:
		return m.boardKey(k)
	case detailScreen:
		m.detailKey(k)
	case versionsScreen:
		return m.versionsKey(k)
	case chooserScreen:
		m.chooserKey(k)
	case sourcesScreen, resultScreen:
		m.scrollKey(k)
	}
	return nil
}

func (m *Model) boardKey(k string) tea.Cmd {
	columns, shelf, _ := m.bounded()
	list := shelf
	if !m.onShelf {
		list = columns[m.col]
	}
	at := slices.IndexFunc(list, func(c card) bool { return c.id == m.cardID })
	focus := func(list []card, i int) {
		m.cardID = ""
		if len(list) != 0 {
			m.cardID = list[min(max(i, 0), len(list)-1)].id
		}
	}
	switch k {
	case "up", "k":
		focus(list, at-1)
	case "down", "j":
		focus(list, at+1)
	case "left", "h", "right", "l":
		if !m.onShelf {
			vis := m.visible()
			i := slices.Index(vis, m.col)
			if k == "left" || k == "h" {
				i = max(i-1, 0)
			} else {
				i = min(i+1, len(vis)-1)
			}
			m.col = vis[i]
			focus(columns[m.col], at)
		}
	case "a":
		// Abandoned work is hidden by default; showing it adds its column.
		if m.showAbandoned = !m.showAbandoned; !m.showAbandoned && m.col == len(statuses)-1 {
			m.col = len(statuses) - 2
			focus(columns[m.col], 0)
		}
	case "tab":
		if m.onShelf = !m.onShelf && len(shelf) != 0; m.onShelf {
			focus(shelf, 0)
		} else {
			focus(columns[m.col], 0)
		}
	case "b":
		if m.res != nil {
			m.back, m.screen, m.choice = boardScreen, chooserScreen, 0
		}
	case "/":
		if m.res != nil {
			m.back, m.screen, m.query, m.hit = boardScreen, searchScreen, "", 0
		}
	case "enter":
		// Opening a card shows its detail. It never resolves a workspace,
		// even when only one version exists.
		if at >= 0 {
			m.openDetail(m.cardID)
		}
	}
	return nil
}

func (m *Model) versionsKey(k string) tea.Cmd {
	g := m.group()
	if g == nil {
		return nil
	}
	rows := m.rows()
	at := slices.IndexFunc(rows, func(r row) bool { return r.key == m.verKey }) // -1 is the header
	move := func(i int) {
		m.verKey, m.scroll = "", 0
		if i = min(i, len(rows)-1); i >= 0 {
			m.verKey = rows[i].key
		}
	}
	switch k {
	case "tab":
		m.detail = !m.detail
	case "pgup", "pgdown":
		m.scrollKey(k)
	case "up", "k", "down", "j":
		switch {
		case m.detail:
			m.scrollKey(k)
		case k == "up" || k == "k":
			move(at - 1)
		default:
			move(at + 1)
		}
	case "enter":
		switch {
		case at < 0 || m.detail:
			// The header and the detail pane select nothing.
		case rows[at].fold != nil:
			// Nor does a fold: it shows or hides the places to choose from.
			if m.unfolded == rows[at].key {
				m.unfolded = ""
			} else {
				m.unfolded = rows[at].key
			}
		case m.busy():
			m.notice = "a read is in progress; wait for it before selecting"
		case rows[at].v.Selector == "" && rows[at].v.Source.Kind == "committed":
			m.refusal = "this record was deleted on that branch, so there is nothing to open there"
		case rows[at].v.Selector == "":
			m.refusal = "this record was deleted from that checkout's live files, so there is nothing to open there"
		default:
			m.refusal = ""
			return m.resolve(rows[at].v.Selector)
		}
	}
	return nil
}

func (m *Model) chooserKey(k string) {
	live := m.live()
	switch k {
	case "up", "k":
		m.choice = max(m.choice-1, 0)
	case "down", "j":
		m.choice = min(m.choice+1, len(live)) // the current view comes first
	case "enter":
		// Choosing a context changes what the board displays. It switches no
		// branch and no directory.
		if m.choice == 0 {
			m.hasBoard, m.lost = false, false
			m.screen, m.col, m.onShelf, m.cardID = boardScreen, 0, false, ""
			m.settleFocus()
			return
		}
		if m.choice > len(live) {
			return
		}
		if s := live[m.choice-1]; s.Valid {
			m.board, m.hasBoard, m.lost = keyOf(s), true, false
			m.screen, m.col, m.onShelf, m.cardID = boardScreen, 0, false, ""
			m.settleFocus()
		} else {
			m.notice = "that checkout cannot fill the board: " + sourceProblem(s)
		}
	}
}

func (m *Model) scrollKey(k string) {
	page := max(m.height-6, 1)
	switch k {
	case "up", "k":
		m.scroll--
	case "down", "j":
		m.scroll++
	case "pgup":
		m.scroll -= page
	case "pgdown":
		m.scroll += page
	}
	m.clampScroll()
}

// settleBoard keeps a chosen checkout's board only while its identity is
// unchanged. Before any choice the board is the current view, which is the
// same from every checkout.
func (m *Model) settleBoard() {
	if m.hasBoard && m.boardSource() == nil {
		m.hasBoard, m.lost = false, true
	}
}

// current reports the current view: no checkout chosen, and none lost.
func (m *Model) current() bool { return !m.hasBoard && !m.lost }

// settleFocus follows the focused card to wherever the new result places it.
func (m *Model) settleFocus() {
	// An open record that vanished leaves the stack, with the reason; when
	// none is left, its detail closes, and the versions or sources screen
	// above it.
	if gone := slices.IndexFunc(m.stack, func(id string) bool { return m.groupOf(id) == nil }); gone >= 0 {
		m.notice = m.stack[gone] + " is no longer on any readable branch or checkout"
		m.stack = slices.DeleteFunc(slices.Clone(m.stack), func(id string) bool { return m.groupOf(id) == nil })
		if len(m.stack) == 0 {
			m.leaveVersions()
			if m.back = boardScreen; m.screen == detailScreen || m.screen == versionsScreen {
				m.screen = boardScreen
			}
		}
	}
	columns, shelf, _ := m.bounded()
	if m.cardID != "" {
		for _, i := range m.visible() {
			if slices.ContainsFunc(columns[i], func(c card) bool { return c.id == m.cardID }) {
				m.col, m.onShelf = i, false
				return
			}
		}
		if slices.ContainsFunc(shelf, func(c card) bool { return c.id == m.cardID }) {
			m.onShelf = true
			return
		}
	}
	m.cardID, m.onShelf = "", m.onShelf && len(shelf) != 0
	if !slices.Contains(m.visible(), m.col) {
		m.col = len(statuses) - 2
	}
	if list := columns[m.col]; !m.onShelf && len(list) != 0 {
		m.cardID = list[0].id
	} else if m.onShelf {
		m.cardID = shelf[0].id
	}
}

func (m *Model) live() []*versions.Source {
	var live []*versions.Source
	if m.res != nil {
		for _, s := range m.res.Sources {
			if s.Kind == "live" {
				live = append(live, s)
			}
		}
	}
	return live
}

func (m *Model) boardSource() *versions.Source {
	if m.hasBoard {
		for _, s := range m.live() {
			if s.Locator != "" && keyOf(s) == m.board {
				return s
			}
		}
	}
	return nil
}

// cards derives the board from the current result. In the current view each
// work record is one card, placed by its current state; see currentCards. A
// checkout's board has one card per work record live there, in that source's
// own status, and a shelf of work groups with no live record there. No status
// is combined across sources.
// ponytail: recomputed per key and frame; cache per result if boards grow large.
func (m *Model) cards() (columns [len(statuses)][]card, shelf []card) {
	if m.res == nil {
		return
	}
	if m.current() {
		return m.currentCards()
	}
	src := m.boardSource()
	for _, g := range m.res.Groups {
		if !isWork(g, src) {
			continue
		}
		placed := false
		for _, v := range g.Versions {
			if src == nil || v.Source != src || v.Record == nil {
				continue
			}
			if i := slices.Index(statuses[:], v.Record.Status); i >= 0 {
				columns[i] = append(columns[i], card{g.ID, v.Record.Title, distinct(g), count("", distinct(g), ""), meta(v.Record), v.Record})
				placed = true
			}
		}
		if !placed {
			shelf = append(shelf, card{id: g.ID, versions: distinct(g), tag: count("", distinct(g), "")})
		}
	}
	newestFirst(columns[doneColumn])
	return
}

const doneColumn = 3

// newestFirst orders Done by when each record was last written, so the
// bounded column shows the latest work: updated, then created, then the ID.
func newestFirst(cards []card) {
	when := func(t *time.Time) int64 {
		if t == nil {
			return 0
		}
		return t.Unix()
	}
	slices.SortStableFunc(cards, func(a, b card) int {
		if c := cmp.Compare(when(b.rec.Updated), when(a.rec.Updated)); c != 0 {
			return c
		}
		if c := cmp.Compare(when(b.rec.Created), when(a.rec.Created)); c != 0 {
			return c
		}
		return cmp.Compare(len(b.id), len(a.id))*2 + cmp.Compare(b.id, a.id)
	})
}

// visible lists the columns on the board: every status but Abandoned, which
// a hides until asked for.
func (m *Model) visible() []int {
	vis := []int{0, 1, 2, 3}
	if m.showAbandoned {
		vis = append(vis, 4)
	}
	return vis
}

// bounded is the board as shown: the columns with Done cut to the cards that
// fit its page, newest first, how many that cut off, and the shelf. Focus
// never lands on a cut card; search still reaches it.
func (m *Model) bounded() (columns [len(statuses)][]card, shelf []card, older int) {
	columns, shelf = m.cards()
	if per := m.pageSize(doneColumn); len(columns[doneColumn]) > per {
		older = len(columns[doneColumn]) - per
		columns[doneColumn] = columns[doneColumn][:per]
	}
	return
}

// columnArea is the rows a column has under the header, banner and column
// heading and above the shelf row.
func (m *Model) columnArea() int { return m.height - 3 - 2 }

// pageSize is how many cards one column page holds: two rows are kept for
// the counts of cards beyond the page, or one for Done's footer, since its
// bound is its page.
func (m *Model) pageSize(status int) int {
	if status == doneColumn {
		return max((m.columnArea()-1)/m.cardRows(), 1)
	}
	return max((m.columnArea()-2)/m.cardRows(), 1)
}

// cardRows is a card's height on this terminal: a bordered box of four
// rows, or of two (the ID and one title row) where the area is too short
// for a full card and Done's footer.
func (m *Model) cardRows() int {
	if m.columnArea() < cardHeight+1 {
		return cardHeight - 2
	}
	return cardHeight
}

// currentCards places each work record by its current states (G-042). One
// state puts the card in its status. Diverging states make one card, marked,
// in the earliest status among them: the owner's choice, so that work is not
// shown further along until its branches agree. A state held only by
// uncommitted files is marked, and so is a card none of whose committed
// current states the integration target holds. The target places nothing.
// The shelf holds work whose current state deletes it. A group is work when
// a current record says so.
func (m *Model) currentCards() (columns [len(statuses)][]card, shelf []card) {
	for _, g := range m.res.Groups {
		states := currentStates(g)
		best, work, deleted, uncommitted := -1, false, false, false
		var rec *project.Record
		for _, state := range states {
			uncommitted = uncommitted || !slices.ContainsFunc(state, committed)
			r := state[0].Record
			if r == nil {
				deleted = true
				continue
			}
			if r.Type != "work" {
				continue
			}
			work = true
			if i := slices.Index(statuses[:], r.Status); i >= 0 && (best < 0 || i < best) {
				best, rec = i, r
			}
		}
		var tags []string
		if len(states) > 1 {
			tags = append(tags, fmt.Sprintf("⑂ %d states", len(states)))
		}
		if uncommitted {
			tags = append(tags, "uncommitted")
		}
		if t := m.res.Target; t != "" && !slices.ContainsFunc(states, func(state []*versions.Version) bool { return state[0].OnTarget }) &&
			slices.ContainsFunc(states, func(state []*versions.Version) bool { return slices.ContainsFunc(state, committed) }) {
			tags = append(tags, "not on "+t)
		}
		tag := strings.Join(tags, " ")
		switch {
		case best >= 0:
			columns[best] = append(columns[best], card{g.ID, rec.Title, len(states), tag, meta(rec), rec})
		case !work && deleted && isWork(g, nil):
			shelf = append(shelf, card{id: g.ID, versions: len(states), tag: tag})
		}
	}
	newestFirst(columns[doneColumn])
	return
}

// committed reports a version held by a commit: a branch's, or a checkout's
// file that matches its HEAD.
func committed(v *versions.Version) bool {
	return v.Source.Kind == "committed" || v.Change == "unchanged"
}

// currentStates lists a group's current states: each distinct current content
// once, with the current versions holding it, in the group's order. Every
// current deletion is one state, whose versions have no record.
func currentStates(g versions.Group) [][]*versions.Version {
	var states [][]*versions.Version
	for i := range g.Versions {
		v := &g.Versions[i]
		if v.Older != "" {
			continue
		}
		if at := slices.IndexFunc(states, func(s []*versions.Version) bool { return contentKey(s[0]) == contentKey(v) }); at >= 0 {
			states[at] = append(states[at], v)
		} else {
			states = append(states, []*versions.Version{v})
		}
	}
	return states
}

// isWork reports a work group for the board of src. Sources can disagree
// about a record's type, so src's own record decides; a group src does not
// hold is work when any source says so, which is what the shelf is for. A
// group of only deleted rows carries no record anywhere, so nothing says it
// was work: it is not shelved. Existing limit: a work item deleted from every
// branch and worktree drops off the board rather than staying visible as a
// ghost card.
func isWork(g versions.Group, src *versions.Source) bool {
	elsewhere := false
	for _, v := range g.Versions {
		if v.Record == nil {
			continue
		}
		if v.Source == src {
			return v.Record.Type == "work"
		}
		elsewhere = elsewhere || v.Record.Type == "work"
	}
	return elsewhere
}

// openID is the record the detail and versions screens are about: the open
// detail's, or the board's focus.
func (m *Model) openID() string {
	if len(m.stack) != 0 {
		return m.stack[len(m.stack)-1]
	}
	return m.cardID
}

func (m *Model) group() *versions.Group { return m.groupOf(m.openID()) }

func (m *Model) groupOf(id string) *versions.Group {
	if m.res != nil {
		for i := range m.res.Groups {
			if m.res.Groups[i].ID == id {
				return &m.res.Groups[i]
			}
		}
	}
	return nil
}

// leaveVersions forgets the version focus: the next selection is made against
// rows the person has seen since.
func (m *Model) leaveVersions() {
	m.verKey, m.unfolded, m.detail, m.scroll, m.refusal = "", "", false, 0, ""
}

// contentKey is what rows fold on: the exact bytes. A deleted row has none
// and never folds.
func contentKey(v *versions.Version) string {
	if v.Record == nil {
		return ""
	}
	return v.Revision
}

// distinct counts a group's differing versions; each deleted row is its own.
func distinct(g versions.Group) int {
	seen := map[string]bool{}
	n := 0
	for i := range g.Versions {
		if k := contentKey(&g.Versions[i]); k == "" || !seen[k] {
			seen[k], n = true, n+1
		}
	}
	return n
}

// rows lists the open card: one row per distinct content, current contents
// first and otherwise in the inspection's order, a fold where several places
// hold it, and the open fold's members.
func (m *Model) rows() []row {
	g := m.group()
	if g == nil {
		return nil
	}
	var order []*versions.Version
	for _, older := range []bool{false, true} {
		for i := range g.Versions {
			if v := &g.Versions[i]; (v.Older != "") == older {
				order = append(order, v)
			}
		}
	}
	same := map[string][]*versions.Version{}
	for _, v := range order {
		if contentKey(v) != "" {
			same[contentKey(v)] = append(same[contentKey(v)], v)
		}
	}
	var rows []row
	for _, v := range order {
		members := same[contentKey(v)]
		switch {
		case len(members) < 2:
			rows = append(rows, row{key: rowKey(*v), v: v})
		case members[0] == v:
			fold := row{key: "fold\x00" + v.Revision, fold: members}
			rows = append(rows, fold)
			if m.unfolded == fold.key {
				for _, member := range members {
					rows = append(rows, row{key: rowKey(*member), v: member, inFold: true})
				}
			}
		}
	}
	return rows
}

// focusedRow returns the row under the cursor, or nil on the ID header.
func (m *Model) focusedRow() *row {
	if m.verKey != "" {
		for _, r := range m.rows() {
			if r.key == m.verKey {
				return &r
			}
		}
	}
	return nil
}

// focused returns the one version under the cursor, if the row is one.
func (m *Model) focused() *versions.Version {
	if r := m.focusedRow(); r != nil {
		return r.v
	}
	return nil
}

// rowKey identifies a version row. A deleted row has no selector; a branch or
// checkout contributes at most one to a group.
func rowKey(v versions.Version) string {
	if v.Selector != "" {
		return v.Selector
	}
	return "deleted\x00" + v.Source.Kind + "\x00" + v.Source.Ref + "\x00" + v.Source.Worktree
}

// interrupted reports the ways a session ends without the person's consent to
// an ordinary exit.
func interrupted(err error) bool {
	return errors.Is(err, tea.ErrInterrupted) || errors.Is(err, context.Canceled)
}
