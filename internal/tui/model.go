// Package tui is Grove's terminal interface: a Kanban board scoped to one live
// checkout, whose cards open every observed version of a record, and explicit
// selection of one version's existing workspace. It reads through Backend and
// changes nothing but the terminal: no records, refs, index, or worktrees.
package tui

import (
	"context"
	"errors"
	"slices"
	"strings"
	"sync"

	tea "charm.land/bubbletea/v2"

	"github.com/mascah/grove/internal/versions"
)

// Backend is every effect the interface has besides drawing.
type Backend struct {
	Inspect func(ctx context.Context, root, id string) (*versions.Result, error)
	Resolve func(ctx context.Context, root, selector string) (*versions.Workspace, error)
}

type screen int

const (
	boardScreen screen = iota
	versionsScreen
	chooserScreen
	sourcesScreen
)

var statuses = [4]string{"proposed", "active", "done", "abandoned"}

// sourceKey identifies the board's checkout across refreshes. A checkout that
// moved, was re-registered, or switched branch is a different context, which
// the person must choose again.
type sourceKey struct{ locator, worktree, ref string }

func keyOf(s *versions.Source) sourceKey { return sourceKey{s.Locator, s.Worktree, s.Ref} }

type card struct {
	id, title string
	versions  int
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
	pending   string // "", "inspect", or "resolve": one read at a time
	resolving string // the exact selector a pending resolve was asked for
	cancel    context.CancelFunc

	board    sourceKey
	hasBoard bool
	lost     bool // the chosen context changed identity; b must choose again

	screen, back screen
	col          int
	onShelf      bool
	cardID       string
	verKey       string // "" is the ID header, which selects nothing
	detail       bool   // the detail pane has focus
	scroll       int    // detail pane, or sources screen
	choice       int    // chooser row
	refusal      string

	// Workspace is the explicitly selected, freshly resolved result, if any.
	Workspace *versions.Workspace
}

// New returns a model that starts by inspecting every record of root.
func New(ctx context.Context, root string, backend Backend) *Model {
	return &Model{ctx: ctx, root: root, backend: backend}
}

func (m *Model) Init() tea.Cmd { return m.read("inspect", "") }

// read starts the one allowed backend call under its own cancellable context.
func (m *Model) read(kind, selector string) tea.Cmd {
	m.stop()
	gen := m.gen
	ctx, cancel := context.WithCancel(m.ctx)
	m.pending, m.resolving, m.cancel = kind, selector, cancel
	return func() tea.Msg {
		defer cancel()
		if !m.reads.begin() {
			return nil
		}
		defer m.reads.wg.Done()
		if kind == "resolve" {
			ws, err := m.backend.Resolve(ctx, m.root, selector)
			return resolveMsg{gen, selector, ws, err}
		}
		res, err := m.backend.Inspect(ctx, m.root, "")
		return inspectMsg{gen, res, err}
	}
}

// stop cancels any read in flight and outdates its reply.
func (m *Model) stop() {
	if m.cancel != nil {
		m.cancel()
	}
	m.gen++
	m.pending, m.resolving, m.cancel = "", "", nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.clampScroll()
	case inspectMsg:
		if msg.gen != m.gen || m.pending != "inspect" {
			return m, nil
		}
		m.pending, m.cancel = "", nil
		if msg.err != nil {
			m.res, m.failure = nil, msg.err.Error()
			m.screen, m.cardID, m.verKey = boardScreen, "", ""
			return m, nil
		}
		m.res, m.failure = msg.res, ""
		m.choice = min(m.choice, max(len(m.live())-1, 0))
		m.settleBoard()
		m.settleFocus()
	case resolveMsg:
		if msg.gen != m.gen || m.pending != "resolve" || msg.selector != m.resolving {
			return m, nil
		}
		m.pending, m.resolving, m.cancel = "", "", nil
		if msg.err != nil || msg.ws == nil {
			m.refusal = "no workspace was returned"
			if msg.err != nil {
				m.refusal = msg.err.Error()
			}
			return m, nil
		}
		m.Workspace = msg.ws
		return m, tea.Quit
	case tea.KeyPressMsg:
		return m, m.key(msg.String())
	}
	return m, nil
}

func (m *Model) key(k string) tea.Cmd {
	m.notice = ""
	switch k {
	case "ctrl+c":
		m.stop()
		return tea.Interrupt
	case "q":
		m.stop()
		return tea.Quit
	case "r":
		if m.pending != "" {
			m.notice = "a read is already in progress"
			return nil
		}
		// A refresh forgets the version focus: the next selection is made
		// against rows the person has seen since.
		m.verKey, m.detail, m.scroll, m.refusal = "", false, 0, ""
		return m.read("inspect", "")
	case "s":
		if m.screen != sourcesScreen && m.res != nil {
			m.back, m.screen, m.scroll = m.screen, sourcesScreen, 0
		}
		return nil
	case "esc":
		switch m.screen {
		case boardScreen:
			m.stop()
			return tea.Quit
		case versionsScreen:
			if m.pending == "resolve" {
				m.stop()
			}
			m.screen, m.verKey, m.detail, m.scroll, m.refusal = boardScreen, "", false, 0, ""
		default:
			m.screen, m.scroll = m.back, 0
		}
		return nil
	}
	switch m.screen {
	case boardScreen:
		return m.boardKey(k)
	case versionsScreen:
		return m.versionsKey(k)
	case chooserScreen:
		m.chooserKey(k)
	case sourcesScreen:
		m.scrollKey(k)
	}
	return nil
}

func (m *Model) boardKey(k string) tea.Cmd {
	columns, shelf := m.cards()
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
			if k == "left" || k == "h" {
				m.col = max(m.col-1, 0)
			} else {
				m.col = min(m.col+1, len(statuses)-1)
			}
			focus(columns[m.col], at)
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
	case "enter":
		// Opening a card shows its versions. It never resolves a workspace,
		// even when only one version exists.
		if at >= 0 {
			m.screen, m.verKey, m.detail, m.scroll, m.refusal = versionsScreen, "", false, 0, ""
		}
	}
	return nil
}

func (m *Model) versionsKey(k string) tea.Cmd {
	g := m.group()
	if g == nil {
		return nil
	}
	at := slices.IndexFunc(g.Versions, func(v versions.Version) bool { return rowKey(v) == m.verKey }) // -1 is the header
	move := func(i int) {
		m.verKey, m.scroll = "", 0
		if i = min(i, len(g.Versions)-1); i >= 0 {
			m.verKey = rowKey(g.Versions[i])
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
			// The header and the detail pane select no source.
		case m.pending != "":
			m.notice = "a read is in progress; wait for it before selecting"
		case g.Versions[at].Selector == "":
			m.refusal = "this record was deleted from that checkout's live files, so there is nothing to open there"
		default:
			m.refusal = ""
			return m.read("resolve", g.Versions[at].Selector)
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
		m.choice = min(m.choice+1, max(len(live)-1, 0))
	case "enter":
		if m.choice >= len(live) {
			return
		}
		// Choosing a context changes what the board displays. It switches no
		// branch and no directory.
		if s := live[m.choice]; s.Valid {
			m.board, m.hasBoard, m.lost = keyOf(s), true, false
			m.screen, m.col, m.onShelf, m.cardID = boardScreen, 0, false, ""
			m.settleFocus()
		} else {
			m.notice = "that checkout cannot be a board context: " + sourceProblem(s)
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

// settleBoard keeps the board context only while its identity is unchanged.
// Before any choice it is the invocation's own checkout.
func (m *Model) settleBoard() {
	if m.hasBoard {
		if m.boardSource() == nil {
			m.hasBoard, m.lost = false, true
		}
		return
	}
	if m.lost {
		return
	}
	for _, s := range m.live() {
		if s.GitDir != "" && s.GitDir == m.res.GitDir {
			m.board, m.hasBoard = keyOf(s), true
		}
	}
}

// settleFocus follows the focused card to wherever the new result places it.
func (m *Model) settleFocus() {
	columns, shelf := m.cards()
	if m.cardID != "" {
		for i, column := range columns {
			if slices.ContainsFunc(column, func(c card) bool { return c.id == m.cardID }) {
				m.col, m.onShelf = i, false
				return
			}
		}
		if slices.ContainsFunc(shelf, func(c card) bool { return c.id == m.cardID }) {
			m.onShelf = true
			return
		}
		// The versions may be open beneath the sources screen.
		if m.screen == versionsScreen || m.back == versionsScreen {
			m.notice = m.cardID + " is no longer in any valid source"
			if m.back = boardScreen; m.screen == versionsScreen {
				m.screen = boardScreen
			}
		}
	}
	m.cardID, m.onShelf = "", m.onShelf && len(shelf) != 0
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

// cards derives the board from the current result: one card per work record
// live in the board source, in that source's own status, and a shelf of work
// groups with no live record there. No status is combined across sources.
// ponytail: recomputed per key and frame; cache per result if boards grow large.
func (m *Model) cards() (columns [4][]card, shelf []card) {
	if m.res == nil {
		return
	}
	src := m.boardSource()
	for _, g := range m.res.Groups {
		if !isWork(g) {
			continue
		}
		placed := false
		for _, v := range g.Versions {
			if src == nil || v.Source != src || v.Record == nil {
				continue
			}
			if i := slices.Index(statuses[:], v.Record.Status); i >= 0 {
				columns[i] = append(columns[i], card{g.ID, v.Record.Title, len(g.Versions)})
				placed = true
			}
		}
		if !placed {
			shelf = append(shelf, card{id: g.ID, versions: len(g.Versions)})
		}
	}
	return
}

// isWork reports a work group. A group of only deleted rows carries no record,
// and validation ties the ID prefix to the type.
func isWork(g versions.Group) bool {
	for _, v := range g.Versions {
		if v.Record != nil {
			return v.Record.Type == "work"
		}
	}
	return strings.HasPrefix(g.ID, "W-")
}

func (m *Model) group() *versions.Group {
	if m.res != nil {
		for i := range m.res.Groups {
			if m.res.Groups[i].ID == m.cardID {
				return &m.res.Groups[i]
			}
		}
	}
	return nil
}

// focused returns the version under the cursor, or nil on the ID header.
func (m *Model) focused() *versions.Version {
	if g := m.group(); g != nil && m.verKey != "" {
		for i := range g.Versions {
			if rowKey(g.Versions[i]) == m.verKey {
				return &g.Versions[i]
			}
		}
	}
	return nil
}

// rowKey identifies a version row. A deleted row has no selector; a checkout
// contributes at most one to a group.
func rowKey(v versions.Version) string {
	if v.Selector != "" {
		return v.Selector
	}
	return "deleted\x00" + v.Source.Worktree
}

// interrupted reports the ways a session ends without the person's consent to
// an ordinary exit.
func interrupted(err error) bool {
	return errors.Is(err, tea.ErrInterrupted) || errors.Is(err, context.Canceled)
}
