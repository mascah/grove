package tui

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mascah/grove/internal/versions"
)

const (
	minWidth, minHeight = 40, 10
	wideWidth           = 100 // every column side by side, and rows beside details
	cardHeight          = 6   // a bordered box: ID and tag, two title rows, metadata
)

// accents colour each status column from the ANSI 16 palette, which follows
// the terminal's theme. Colour is an accent only: focus has its heavy border
// and marker, and every tag is text.
var accents = [len(statuses)]lipgloss.Style{
	lipgloss.NewStyle().Foreground(lipgloss.Color("4")),
	lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
	lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
	lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
	lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
}

// runningAccent borders a card whose work has an attempt that may be running,
// in a colour no status column uses; its tag says the same in text.
var runningAccent = cyan

// cardBox draws one card as a bordered box of rows rows and w cells: its ID
// with any tag at the right, the title on two rows, and the metadata row
// unless rows leaves no room. The focused card has a heavy border and a
// marker before its ID.
func cardBox(c card, focused bool, accent lipgloss.Style, rows, w int) []string {
	iw := max(w-2, 1)
	b := lipgloss.RoundedBorder()
	if focused {
		b, accent = lipgloss.ThickBorder(), accent.Bold(true)
	}
	edge := func(l, mid, r string) string { return accent.Render(l + strings.Repeat(mid, iw) + r) }
	side := accent.Render(b.Left)
	id := " " + c.id
	if focused {
		id = "▶" + c.id
	}
	first := line(id, iw)
	if tw := ansi.StringWidth(safe(c.tag)); c.tag != "" && ansi.StringWidth(safe(id))+2+tw <= iw {
		first = line(id, iw-tw-1) + " " + line(c.tag, tw)
	} else if c.tag != "" {
		first = line(id+"  "+c.tag, iw)
	}
	if focused {
		first = bold(first)
	}
	title := wrap(c.title, iw-1)
	if rows < cardHeight {
		title = title[:1]
	} else if len(title) > 2 {
		title[1] = line(strings.TrimRight(title[1], " ")+"…", iw-1)
	}
	inner := []string{first}
	for _, t := range fit(title, min(rows-3, 2), iw-1) {
		inner = append(inner, " "+t)
	}
	if rows >= cardHeight {
		inner = append(inner, " "+line(c.meta, iw-1))
	}
	out := []string{edge(b.TopLeft, b.Top, b.TopRight)}
	for _, r := range inner {
		out = append(out, side+r+side)
	}
	return append(out, edge(b.BottomLeft, b.Bottom, b.BottomRight))
}

// safe makes text from records, paths, and Git inert for a terminal. Controls
// (ESC, C0, C1), format characters such as bidirectional overrides, and
// invalid bytes become visible escapes; letters, marks, and wide characters
// stay. The model's bytes are never changed, only this display copy.
// ponytail: also escapes the joiners inside emoji sequences and non-ASCII
// spaces; allow those individually if titles need them.
func safe(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		r, n := utf8.DecodeRuneInString(s[i:])
		switch {
		case r == utf8.RuneError && n == 1:
			fmt.Fprintf(&b, `\x%02x`, s[i])
		case strconv.IsPrint(r):
			b.WriteRune(r)
		default:
			quoted := strconv.QuoteRune(r)
			b.WriteString(quoted[1 : len(quoted)-1])
		}
		i += n
	}
	return b.String()
}

// line is the only way text reaches the screen: escaped first, then clipped
// and padded by display cells. Styles are added around its result, so nothing
// from a file can be read as one.
func line(s string, w int) string {
	s = ansi.Truncate(safe(s), w, "…")
	return s + strings.Repeat(" ", max(w-ansi.StringWidth(s), 0))
}

// wrap escapes one logical line and breaks it into rows of at most w cells,
// between words where it can. exact breaks only at the edge instead, for
// values such as paths and selectors whose every character matters.
func wrap(s string, w int) []string  { return wrapped(s, w, false) }
func exact(s string, w int) []string { return wrapped(s, w, true) }

func wrapped(s string, w int, hard bool) []string {
	s, w = safe(strings.ReplaceAll(s, "\t", "    ")), max(w, 1)
	if hard {
		s = ansi.Hardwrap(s, w, true)
	} else {
		s = ansi.Wrap(s, w, "")
	}
	rows := strings.Split(s, "\n")
	for i := range rows {
		rows[i] += strings.Repeat(" ", max(w-ansi.StringWidth(rows[i]), 0))
	}
	return rows
}

// wrapAll wraps text whose own line breaks are kept as lines.
func wrapAll(text string, w int) []string {
	var rows []string
	for _, l := range strings.Split(strings.TrimRight(text, "\n"), "\n") {
		rows = append(rows, wrap(l, w)...)
	}
	return rows
}

// Focus is a marker and reverse video; headings are bold. Meaning never
// depends on either.
func hot(s string) string  { return "\x1b[7m" + s + "\x1b[m" }
func bold(s string) string { return "\x1b[1m" + s + "\x1b[m" }

func mark(focused bool, text string, w int) string {
	if focused {
		return hot(line("> "+text, w))
	}
	return line("  "+text, w)
}

// fit returns exactly n rows of width w.
func fit(rows []string, n, w int) []string {
	for len(rows) < n {
		rows = append(rows, strings.Repeat(" ", w))
	}
	return rows[:max(n, 0)]
}

// window shows the page holding item at, where each item is unit rows and a
// page is per items. It keeps one row at each end to count what lies beyond,
// so every item stays reachable.
func window(rows []string, at, per, unit, n, w int) []string {
	per = max(per, 1)
	at = min(max(at, 0), max(len(rows)/unit-1, 0)) // a refresh may have shortened the list
	start := at / per * per
	end := min(start+per, len(rows)/unit)
	var out []string
	if start > 0 {
		out = append(out, line(fmt.Sprintf("  ↑ %d above", start), w))
	}
	out = append(out, rows[start*unit:end*unit]...)
	if end < len(rows)/unit {
		out = append(fit(out, n-1, w), line(fmt.Sprintf("  ↓ %d below", len(rows)/unit-end), w))
	}
	return fit(out, n, w)
}

// pick returns the longest key hint that fits.
func pick(w int, hints ...string) string {
	for _, h := range hints {
		if ansi.StringWidth(h) <= w {
			return h
		}
	}
	return hints[len(hints)-1]
}

func (m *Model) View() tea.View {
	v := tea.NewView(m.render())
	v.AltScreen = true
	v.ReportFocus = true // a terminal that reports focus re-reads the board on return (G-124)
	return v
}

func (m *Model) render() string {
	w, h := m.width, m.height
	if w <= 0 || h <= 0 {
		return ""
	}
	if w < minWidth || h < minHeight {
		return strings.Join(fit([]string{
			line(fmt.Sprintf("Grove needs %dx%d; this is %dx%d.", minWidth, minHeight, w, h), w),
			line("Resize, or q to quit.", w),
		}, min(h, 2), w), "\n")
	}
	body := h - 3
	var rows []string
	var hints string
	side, sideKey := "", "" // w hides the detail's sidebar where there is one beside the content
	if w >= wideWidth {
		side, sideKey = "w sidebar   ", "w sidebar  "
	}
	switch {
	case m.res == nil:
		rows, hints = m.emptyBody(w), "r retry   q quit"
	case m.screen == detailScreen && m.group() != nil && m.backend.Attempts != nil && m.openRecord() != nil && m.openRecord().Type == "work":
		// Work adds its attempts: A lists them, R launches one where it can.
		run, runKey, judge, keys := "", "", "", ""
		if s := m.openRecord().Status; m.backend.Launch != nil && (s == "proposed" || s == "active") {
			run, runKey = "R launch   ", "R  "
		}
		if m.reviewable() {
			judge, keys = "a approve   f feedback   i integrate   ", "a  f  i  "
		}
		rows, hints = m.detailBody(w, body), pick(w,
			"↑/↓ PgUp/PgDn scroll or move   Tab pane   Enter open   "+run+"A attempts   "+judge+side+"v versions   s   r   Esc back   q quit",
			"↑↓ PgUp/PgDn  Tab pane  Enter open  "+run+"A attempts  "+keys+sideKey+"v  Esc  q",
			"↑↓  Tab  Enter  "+runKey+"A  "+keys+"v  Esc  q quit")
	case m.screen == detailScreen && m.group() != nil && m.reviewable():
		rows, hints = m.detailBody(w, body), pick(w,
			"↑/↓ PgUp/PgDn scroll or move   Tab pane   Enter open   a approve   f feedback   i integrate   "+side+"v versions   s   r   Esc back   q quit",
			"↑↓ PgUp/PgDn  Tab pane  Enter open  a approve  f feedback  i integrate  "+sideKey+"v  Esc  q",
			"↑↓  Tab  Enter  a  f  i  v  Esc  q quit")
	case m.screen == detailScreen && m.group() != nil && m.backend.Edit != nil && m.openRecord() != nil && m.openRecord().Type == "question" && m.openRecord().Status == "open":
		rows, hints = m.detailBody(w, body), pick(w,
			"↑/↓ PgUp/PgDn scroll or move   Tab pane   Enter open   e answer in your editor   "+side+"v versions   s   r   Esc back   q quit",
			"↑↓ PgUp/PgDn  Tab pane  Enter open  e answer  "+sideKey+"v  Esc  q",
			"↑↓  Tab  Enter  e answer  v  Esc  q quit")
	case m.screen == detailScreen && m.group() != nil:
		rows, hints = m.detailBody(w, body), pick(w,
			"↑/↓ PgUp/PgDn scroll or move   Tab content, linked, changes, timeline   Enter open   "+side+"v versions and places   s sources   r refresh   Esc back   q quit",
			"↑↓ PgUp/PgDn  Tab pane  Enter open  "+sideKey+"v versions  s  r  Esc back  q quit",
			"↑↓  Tab  Enter open  v  Esc back  q quit")
	case m.screen == resultScreen && m.result != nil:
		rows, hints = m.scrolled(m.resultRows(w), body, w), pick(w, "↑/↓ PgUp/PgDn scroll   Esc back to the record   q quit", "↑↓ scroll  Esc back  q quit")
		if m.pending == "inspect" {
			hints = pick(w, "↑/↓ PgUp/PgDn scroll   Esc waits for the re-read   q quit", "↑↓ scroll  Esc waits  q quit")
		}
	case m.screen == versionsScreen && m.group() != nil:
		rows, hints = m.versionsBody(w, body), pick(w,
			"↑/↓ rows   Enter list places, or select a workspace   Tab details   PgUp/PgDn scroll   s what was read   r refresh   Esc board   q quit",
			"↑↓ move  Enter list places or select workspace  Tab details  PgUp/PgDn  s  r  Esc  q",
			"↑↓  Enter list or select  Tab details  r  Esc  q quit",
			"Enter select  Tab  r  Esc back  q quit")
	case m.screen == chooserScreen:
		rows, hints = m.chooserBody(w, body), pick(w, "↑/↓ choices   Enter show it   Esc back   q quit", "Enter choose  Esc back  q quit")
	case m.screen == searchScreen:
		rows, hints = m.searchBody(w, body), pick(w, "type to filter   Backspace   ↑/↓ move   Enter open   Esc close   Ctrl-C quit", "type  ↑↓  Enter open  Esc close  ^C quit")
	case m.screen == sourcesScreen:
		rows, hints = m.scrolled(m.sourceRows(w), body, w), pick(w, "↑/↓ PgUp/PgDn scroll   r refresh   Esc back   q quit", "↑↓ scroll  Esc back  q quit")
	case m.screen == attemptsScreen:
		answer, key := "", ""
		if list, _, at := m.listed(); at >= 0 && m.waits(&list[at]) {
			answer, key = "e answer   ", "e answer  "
		}
		rows, hints = m.attemptsBody(w, body), pick(w,
			"↑/↓ attempts   Enter show it   "+answer+"x stop it   o open its record   r refresh   Esc back   q quit",
			"↑↓  Enter show  "+key+"x stop  o record  Esc back  q quit")
	case m.screen == depsScreen && m.previewing:
		rows, hints = m.depsBody(w, body), pick(w, "↑/↓ PgUp/PgDn scroll   r re-read   Esc back to the list   q quit", "↑↓ scroll  r  Esc back  q quit")
	case m.screen == depsScreen:
		move, tab := "↑/↓ move", "Tab trees"
		if m.depsTree {
			move, tab = "↑/↓ PgUp/PgDn scroll the trees", "Tab list"
		}
		rows, hints = m.depsBody(w, body), pick(w,
			fmt.Sprintf("%s  Space select  p preview (%d selected)  c clear  Enter open  h history  %s  b checkout  Esc board  q", move, len(m.depsPicked), tab),
			fmt.Sprintf("↑↓  Space select  p preview (%d)  c  Enter open  h history  %s  Esc  q", len(m.depsPicked), tab),
			fmt.Sprintf("↑↓  Space  p (%d)  Enter  h  %s  Esc  q", len(m.depsPicked), tab))
	case m.screen == attemptScreen:
		answer, key := "", ""
		if v := m.attemptOf(m.runID); v != nil && m.waits(v) {
			answer, key = "e answer the question   ", "e answer  "
		}
		rows, hints = m.scrolled(m.attemptRows(w), body, w), pick(w,
			"↑/↓ PgUp/PgDn scroll   "+answer+"d details   x stop it   o open its record   r refresh   Esc back   q quit",
			"↑↓ scroll  "+key+"d details  x stop  o work  Esc back  q quit")
	default:
		shelf := "elsewhere"
		if m.current() {
			shelf = "deleted"
		}
		rows, hints = m.boardBody(w, body), pick(w,
			"←/→ columns   ↑/↓ cards   Enter open   / search   g dependencies   a abandoned   A attempts   Tab "+shelf+"   b view or checkout   s sources   r refresh   q quit",
			"←→↑↓ move  Enter open  / search  g dependencies  a abandoned  Tab "+shelf+"  b view or checkout  s  r  q quit",
			"←→↑↓  Enter open  / search  g deps  a  Tab "+shelf+"  b  s  r  q quit",
			"Enter open  / g a Tab b s r  q quit")
	}
	last := line(hints, w)
	if m.prompt != nil {
		last = m.promptRow(w)
	}
	out := append([]string{bold(line(m.header(), w)), line(m.banner(), w)}, fit(rows, body, w)...)
	return strings.Join(append(out, last), "\n")
}

func (m *Model) header() string {
	// The board's label comes first: a long project path must not clip it.
	text := "Grove    Board: "
	switch s := m.boardSource(); {
	case m.res == nil:
		return "Grove    " + m.root
	case m.current():
		text += "current view"
		if m.res.Target != "" {
			text += ", target " + m.res.Target
		}
	case s == nil:
		text += "no checkout selected"
	case !s.Valid:
		text += label(s) + " [INVALID]"
	default:
		text += label(s)
	}
	return fmt.Sprintf("%s    read %s at %s    %s", text, places(m.res.Sources), m.readAt.Format("15:04:05"), m.root)
}

// banner is always on screen: what is being read, and whether the result is
// incomplete, whatever else is showing.
func (m *Model) banner() string {
	// The incomplete warning leads: a narrow terminal clips the row's end.
	var parts []string
	if m.res != nil && !m.res.Complete {
		bad := 0
		for _, s := range m.res.Sources {
			if len(s.Diagnostics) != 0 {
				bad++
			}
		}
		parts = append(parts, fmt.Sprintf("INCOMPLETE: %d of %d branches and checkouts could not be read (s shows why)", bad, len(m.res.Sources)))
	}
	switch m.pending {
	case "inspect":
		parts = append(parts, "Reading branches and checkouts…")
	case "resolve":
		parts = append(parts, "Resolving the selected workspace…")
	case "act":
		parts = append(parts, actingText(m.acting))
	}
	if m.notice != "" {
		parts = append(parts, m.notice)
	}
	return strings.Join(parts, "   ")
}

func (m *Model) emptyBody(w int) []string {
	if m.failure == "" {
		return nil
	}
	rows := []string{bold(line("The repository's branches and checkouts could not be listed, so there is nothing to select:", w))}
	return append(rows, wrapAll(m.failure, w)...)
}

func (m *Model) boardBody(w, n int) []string {
	columns, shelf, older := m.bounded()
	shelfName, shelfWhy := "Elsewhere", "not in this checkout"
	if m.current() {
		shelfName, shelfWhy = "Deleted", "the current state removes the record"
	}
	if m.onShelf {
		rows := make([]string, len(shelf))
		at := -1
		for i, c := range shelf {
			if c.id == m.cardID {
				at = i
			}
			rows[i] = mark(c.id == m.cardID, strings.TrimRight(c.id+"   "+c.tag, " "), w)
		}
		head := "Elsewhere: work on other branches or checkouts that this board's checkout lacks, so no status here"
		if m.current() {
			head = "Deleted: work whose current state removes its record; older copies remain on other branches or checkouts"
		}
		return append([]string{bold(line(head, w))}, window(rows, at, n-3, 1, n-1, w)...)
	}
	shelfRow := shelfName + ": none"
	if len(shelf) != 0 {
		var items []string
		for _, c := range shelf {
			item := c.id
			if c.tag != "" {
				item += " [" + c.tag + "]"
			}
			items = append(items, item)
		}
		shelfRow = fmt.Sprintf("%s (%d, %s; Tab): %s", shelfName, len(shelf), shelfWhy, strings.Join(items, "  "))
	}
	// The hidden column is counted at the shelf row's end, so nothing is
	// silently missing; a long shelf gives way to it.
	note := "a hides Abandoned"
	if !m.showAbandoned {
		note = fmt.Sprintf("Abandoned %d hidden · a shows", len(columns[len(statuses)-1]))
	}
	shelfRow = line(shelfRow, w-ansi.StringWidth(note)-2) + "  " + note
	area := n - 2 // a heading row above, the shelf row below
	var rows []string
	if s := m.boardSource(); !m.current() && (s == nil || !s.Valid) {
		rows = append([]string{bold(line("No board: "+m.contextProblem(s), w))}, wrapAll("This is not an empty board. Press b to view another checkout, or s for what went wrong in each branch and checkout.", w)...)
		if s != nil {
			for _, d := range s.Diagnostics {
				rows = append(rows, wrapAll("  "+d, w)...)
			}
		}
		rows = fit(rows, n-1, w)
	} else if vis := m.visible(); w >= wideWidth {
		heads := make([]string, len(vis))
		cells := make([][]string, len(vis))
		for j, i := range vis {
			cw := w / len(vis)
			if j == len(vis)-1 {
				cw = w - j*cw
			}
			heads[j] = accents[i].Bold(true).Render(line(m.heading(i, len(columns[i]), older), cw-1)) + " "
			cells[j] = m.column(columns[i], i, i == m.col, older, cw-1, area)
		}
		rows = []string{strings.Join(heads, "")}
		for r := range area {
			var b strings.Builder
			for i := range cells {
				b.WriteString(cells[i][r] + " ")
			}
			rows = append(rows, b.String())
		}
	} else {
		var tabs []string
		for _, i := range vis {
			name := title(statuses[i])
			if w < 60 {
				name = shortNames[i]
			}
			tab := fmt.Sprintf("%s %d", name, len(columns[i]))
			if i == doneColumn {
				tab = fmt.Sprintf("%s %d", name, len(columns[i])+older)
			}
			if i == m.col {
				tab = "[" + tab + "]"
			}
			tabs = append(tabs, tab)
		}
		rows = append([]string{bold(line(strings.Join(tabs, " "), w))}, m.column(columns[m.col], m.col, true, older, w, area)...)
	}
	return append(rows, line(shelfRow, w))
}

// heading names a column with its count; Done also says how many of its
// cards the page shows.
func (m *Model) heading(status, n, older int) string {
	if status == doneColumn && older > 0 {
		return fmt.Sprintf("%s %d · %d recent", title(statuses[status]), n+older, n)
	}
	return fmt.Sprintf("%s %d", title(statuses[status]), n)
}

// shortNames abbreviates the column names for narrow terminals, one per status.
var shortNames = [len(statuses)]string{"Prop", "Act", "Rev", "Done", "Aban"}

// column renders one status column to exactly n rows: a page of cards, with
// the counts beyond it, or for Done the count its bound cut off and how to
// reach them.
func (m *Model) column(cards []card, status int, focused bool, older, w, n int) []string {
	if len(cards) == 0 {
		return fit([]string{line("  (none)", w)}, n, w)
	}
	var rows []string
	at := -1
	for i, c := range cards {
		on := focused && c.id == m.cardID
		if on {
			at = i
		}
		accent := accents[status]
		if c.running {
			accent = runningAccent
		}
		rows = append(rows, cardBox(c, on, accent, m.cardRows(), w)...)
	}
	out := window(rows, at, m.pageSize(status), m.cardRows(), n, w)
	if status == doneColumn && older > 0 {
		out = append(fit(out, n-1, w), line(fmt.Sprintf("  + %d older · / to search", older), w))
	}
	return out
}

func (m *Model) versionsBody(w, n int) []string {
	g, top := m.group(), m.refusalRows(w, n)
	n -= len(top)
	list := func(w int) []string {
		rows := []string{mark(m.verKey == "" && !m.detail, g.ID+"   "+summary(g), w)}
		at := 0
		for i, r := range m.rows() {
			var text string
			switch {
			case r.fold != nil && m.unfolded == r.key:
				text = fmt.Sprintf("▾ %-9s  %ssame on %s", r.fold[0].Record.Status, older(r.fold...), held(r.fold))
			case r.fold != nil:
				text = fmt.Sprintf("▸ %-9s  %ssame on %s", r.fold[0].Record.Status, older(r.fold...), held(r.fold))
			case r.inFold:
				text = fmt.Sprintf("      %s%s  %s", older(r.v), label(r.v.Source), r.v.Change)
			case r.v.Record == nil:
				text = fmt.Sprintf("  %-9s  %s%s  %s", "-", older(r.v), label(r.v.Source), r.v.Change)
			default:
				text = fmt.Sprintf("  %-9s  %s%s  %s", r.v.Record.Status, older(r.v), label(r.v.Source), r.v.Change)
			}
			if r.key == m.verKey {
				at = i + 1
			}
			rows = append(rows, mark(r.key == m.verKey && !m.detail, text, w))
		}
		return window(rows, at, n-2, 1, n, w)
	}
	details := func(w int) []string {
		all := m.detailRows(w)
		head := fmt.Sprintf("Details  %d-%d of %d", min(m.scroll+1, len(all)), min(m.scroll+n-1, len(all)), len(all))
		if m.detail {
			head = hot(line("> "+head, w))
		} else {
			head = bold(line("  "+head, w))
		}
		return append([]string{head}, fit(all[min(m.scroll, len(all)):], n-1, w)...)
	}
	var rows []string
	switch {
	case w < wideWidth && m.detail:
		rows = details(w)
	case w < wideWidth:
		rows = list(w)
	default:
		lw := w * 2 / 5
		left, right := list(lw), details(w-lw-3)
		for i := range n {
			rows = append(rows, left[i]+" │ "+right[i])
		}
	}
	return append(top, rows...)
}

// refusalRows keeps a refused selection and its reason above the versions
// until the person refreshes or leaves.
func (m *Model) refusalRows(w, n int) []string {
	if m.refusal == "" {
		return nil
	}
	top := append(wrapAll("REFUSED: "+m.refusal, w), line("Nothing was opened. Press r to refresh, then select a version again.", w))
	if limit := max(n/2, 2); len(top) > limit {
		top = append(top[:limit-1], line("  … (s shows each branch and checkout's diagnostics)", w))
	}
	for i := range top {
		top[i] = bold(top[i])
	}
	return top
}

// detailRows is the details pane: the history first, then the description.
func (m *Model) detailRows(w int) []string {
	rows := m.historyRows(w)
	if len(rows) != 0 {
		rows = append(rows, line(strings.Repeat("─", w), w))
	}
	return append(rows, m.describeRows(w)...)
}

// historyRows lists the commits behind one version, newest first, under a
// heading naming the branch or checkout they were read from. They say what
// happened there, and nothing about any other branch.
func (m *Model) historyRows(w int) []string {
	v := m.historyOf()
	if v == nil || m.backend.History == nil {
		return nil
	}
	var rows []string
	for _, r := range wrap("History on "+label(v.Source)+": commits that changed this record's file, newest first", w) {
		rows = append(rows, bold(r))
	}
	items, message := m.historyItems(v)
	for _, it := range items {
		status := it.status
		if status == "" {
			status = "?"
		}
		// Continuation rows are indented under the first.
		for i, r := range wrap(fmt.Sprintf("%-16s  %-9s  %s", it.when, status, it.rest), w-2) {
			if i == 0 {
				rows = append(rows, r+"  ")
			} else {
				rows = append(rows, "  "+r)
			}
		}
	}
	if message != "" {
		rows = append(rows, wrapAll(message, w)...)
	}
	return rows
}

// describeRows describes the focused row, or the group when the ID header has
// focus. The header and a fold describe; they select nothing.
func (m *Model) describeRows(w int) []string {
	g, r := m.group(), m.focusedRow()
	if g == nil {
		return nil
	}
	fields := [][2]string{}
	add := func(name, value string) {
		if value != "" {
			fields = append(fields, [2]string{name, value})
		}
	}
	finish := func(source []byte) []string {
		var rows []string
		for _, f := range fields {
			name := f[0] + ":"
			if f[0] == "" { // continues the field above
				name = ""
			}
			rows = append(rows, exact(fmt.Sprintf("%-10s%s", name, f[1]), w)...)
		}
		if source != nil {
			rows = append(rows, line(strings.Repeat("─", w), w))
			rows = append(rows, wrapAll(string(source), w)...)
		}
		return rows
	}
	switch {
	case r == nil:
		all := make([]*versions.Version, len(g.Versions))
		for i := range g.Versions {
			all[i] = &g.Versions[i]
		}
		return wrapAll(g.ID+": "+summary(g)+" ("+held(all)+").\n\n"+currentText(g, m.res.Target)+"\n\nA version is this record's exact content. Grove read it on every local branch (its committed tip) and in every checkout (its files on disk, committed or not) and lists each differing content once, current ones first.\n\nMove to a row. Enter on ▸ lists the branches and checkouts holding that content; Enter on one of them, or on a row naming a single place, returns the path of the existing checkout to work in.", w)
	case r.fold != nil:
		first := r.fold[0]
		add("Title", first.Record.Title)
		add("Status", first.Record.Status)
		standing := "current"
		if older(r.fold...) != "" {
			standing = "older: " + first.Older
		}
		add("Standing", standing)
		add("Revision", first.Revision)
		for i, v := range r.fold {
			name := ""
			if i == 0 {
				name = "Same on"
			}
			fields = append(fields, [2]string{name, strings.TrimSpace(label(v.Source) + "  " + v.Change)})
		}
		add("Note", "Enter lists these places so one can be selected.")
		return finish(first.Record.Source)
	}
	v := r.v
	s := v.Source
	if v.Record != nil {
		add("Title", v.Record.Title)
		add("Status", v.Record.Status)
	} else if s.Kind == "committed" {
		add("Status", "deleted on this branch; cannot be opened")
	} else {
		add("Status", "deleted from this checkout's live files; cannot be opened")
	}
	standing := "current"
	if v.Older != "" {
		standing = "older: " + v.Older
	}
	add("Standing", standing)
	if s.Kind == "live" {
		add("Seen in", "a checkout's files on disk")
	} else {
		add("Seen in", "a branch's committed tip")
	}
	if s.Ref != "" {
		add("Branch", s.Ref)
	} else {
		add("Branch", "detached HEAD")
	}
	add("Commit", s.Commit)
	add("Checkout", s.Worktree)
	add("Change", v.Change)
	add("Path", v.Path)
	add("HEAD path", v.HeadPath)
	add("Revision", v.Revision)
	add("Selector", v.Selector)
	add("Note", s.Note)
	if v.Record == nil {
		return finish(nil)
	}
	return finish(v.Record.Source)
}

// currentText describes a record's current state or states, how many places
// hold an older copy, and any pair that could not be ordered.
func currentText(g *versions.Group, target string) string {
	states := currentStates(*g)
	var b strings.Builder
	if len(states) == 1 {
		fmt.Fprintf(&b, "Current: %s.", stateText(states[0], target))
	} else {
		fmt.Fprintf(&b, "Diverging: %d current states. None is known to replace the others: each changed this record since they split from a common commit, or their order could not be read (below). The card sits in the earliest status among them until one side takes the other's change, by a merge or an edit.", len(states))
		for _, state := range states {
			b.WriteString("\n  - " + stateText(state, target))
		}
	}
	n := 0
	for _, v := range g.Versions {
		if v.Older != "" {
			n++
		}
	}
	if n != 0 {
		fmt.Fprintf(&b, "\nOlder: %d %s an earlier state that another place has changed since (marked older below).", n, plural(n, "place holds", "places hold"))
	}
	for _, note := range g.Notes {
		b.WriteString("\nCould not order: " + note)
	}
	return b.String()
}

// stateText names one current state, where it is held, and whether the
// integration target holds it.
func stateText(state []*versions.Version, target string) string {
	text := "deleted"
	if r := state[0].Record; r != nil {
		text = r.Status + " \"" + r.Title + "\""
	}
	names := held(state)
	if len(state) <= 3 {
		var labels []string
		for _, v := range state {
			labels = append(labels, label(v.Source))
		}
		names = strings.Join(labels, ", ")
	}
	text += " on " + names
	switch {
	case !slices.ContainsFunc(state, committed):
		text += ", uncommitted"
	case target == "":
	case state[0].OnTarget:
		text += ", on " + target
	default:
		text += ", not on " + target
	}
	return text
}

func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

func (m *Model) chooserBody(w, n int) []string {
	rows := []string{mark(m.choice == 0, "current view   every branch and checkout, each record in its current state", w)}
	for i, s := range m.live() {
		state := "valid"
		if !s.Valid {
			state = "UNAVAILABLE: " + sourceProblem(s)
		}
		// The path comes last: clipping a long one must not hide the state.
		rows = append(rows, mark(i+1 == m.choice, fmt.Sprintf("%s   %s   %s", label(s), state, s.Worktree), w))
	}
	head := wrapAll("Choose what the board shows: the current view, or one checkout's own files. Only this display changes: Git switches no branch, and your shell stays where it is.", w)
	for i := range head {
		head[i] = bold(head[i])
	}
	n -= len(head)
	return append(head, window(rows, m.choice, n-2, 1, n, w)...)
}

func (m *Model) sourceRows(w int) []string {
	rows := []string{bold(line("Branches and checkouts read in "+m.res.Repository, w))}
	if m.res.Target != "" {
		rows = append(rows, line("Target: "+m.res.Target+", named in grove.yaml", w))
	}
	for _, n := range m.res.Notes {
		rows = append(rows, wrapAll("Target: "+n, w)...)
	}
	for _, s := range m.res.Sources {
		state := "valid"
		if !s.Valid {
			state = sourceProblem(s)
		}
		text := label(s) + "  " + short(s.Commit)
		if s.Kind == "live" {
			text += "  " + s.Worktree
		}
		rows = append(rows, exact(text+"  ("+state+")", w)...)
		if s.Note != "" {
			rows = append(rows, wrapAll("    note: "+s.Note, w)...)
		}
		for _, d := range s.Diagnostics {
			rows = append(rows, wrapAll("    "+d, w)...)
		}
	}
	return rows
}

func (m *Model) scrolled(rows []string, n, w int) []string {
	return fit(rows[min(m.scroll, len(rows)):], n, w)
}

// clampScroll keeps the scroll offset inside what the current pane shows.
func (m *Model) clampScroll() {
	var rows, n int
	switch {
	case m.res == nil:
	case m.screen == sourcesScreen:
		rows, n = len(m.sourceRows(m.width)), m.height-3
	case m.screen == resultScreen && m.result != nil:
		rows, n = len(m.resultRows(m.width)), m.height-3
	case m.screen == attemptScreen:
		rows, n = len(m.attemptRows(m.width)), m.height-3
	case m.screen == depsScreen && m.previewing:
		rows, n = len(m.previewRows(m.width)), m.height-3
	case m.screen == versionsScreen:
		w := m.width
		if w >= wideWidth {
			w -= w*2/5 + 3
		}
		rows, n = len(m.detailRows(w)), m.height-4-len(m.refusalRows(m.width, m.height-3))
	case m.screen == detailScreen:
		g := m.group()
		if g == nil {
			return
		}
		w, v := m.width, m.shown(g)
		if m.split() {
			w = w * 11 / 20
		}
		rows, n = len(m.contentRows(v, w)), m.height-3-len(m.detailHead(v, m.width))-1
		m.dscroll = max(min(m.dscroll, rows-n), 0)
		return
	}
	m.scroll = max(min(m.scroll, rows-n), 0)
}

func (m *Model) contextProblem(s *versions.Source) string {
	switch {
	case s != nil:
		return label(s) + " is not usable: " + sourceProblem(s)
	}
	return "the chosen checkout changed branch, moved, or was removed"
}

func sourceProblem(s *versions.Source) string {
	switch {
	case len(s.Diagnostics) != 0:
		return "invalid: " + s.Diagnostics[0]
	case !s.Present:
		return "no project at this location"
	}
	return "invalid"
}

// label names a branch tip or a checkout the way every screen shows it.
func label(s *versions.Source) string {
	ref := strings.TrimPrefix(s.Ref, "refs/heads/")
	if s.Kind != "live" {
		return "branch " + ref
	}
	if ref == "" {
		ref = "detached at " + short(s.Commit)
	}
	locator := s.Locator
	if locator == "" {
		locator = "?"
	}
	return "checkout " + locator + " (" + ref + ")"
}

// places words how many branch tips and checkouts are among sources.
func places(sources []*versions.Source) string {
	live := 0
	for _, s := range sources {
		if s.Kind == "live" {
			live++
		}
	}
	var parts []string
	for _, p := range [][3]string{{strconv.Itoa(len(sources) - live), "branch", "branches"}, {strconv.Itoa(live), "checkout", "checkouts"}} {
		if p[0] == "1" {
			parts = append(parts, "1 "+p[1])
		} else if p[0] != "0" {
			parts = append(parts, p[0]+" "+p[2])
		}
	}
	return strings.Join(parts, ", ")
}

// older marks versions that another is newer than (G-042): a fold only when
// every place in it holds an older copy.
func older(vs ...*versions.Version) string {
	for _, v := range vs {
		if v.Older == "" {
			return ""
		}
	}
	return "older  "
}

func held(vs []*versions.Version) string {
	sources := make([]*versions.Source, len(vs))
	for i, v := range vs {
		sources[i] = v.Source
	}
	return places(sources)
}

// summary is an open card's first line, short enough for the list pane.
func summary(g *versions.Group) string {
	if n := distinct(*g); n > 1 {
		return strconv.Itoa(n) + " versions differ"
	} else if len(g.Versions) > 1 {
		return "same everywhere"
	}
	return "in one place"
}

func short(commit string) string { return commit[:min(len(commit), 12)] }

// count notes differing versions between before and after; agreement needs
// no note.
func count(before string, n int, after string) string {
	if n < 2 {
		return ""
	}
	return before + strconv.Itoa(n) + " versions" + after
}

func title(s string) string { return strings.ToUpper(s[:1]) + s[1:] }
