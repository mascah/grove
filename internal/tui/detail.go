package tui

import (
	"fmt"
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/versions"
)

// The detail is a record's first useful view: its content rendered from
// Markdown, beside the records linked to it and the timeline of commits that
// changed it. Sources and versions are secondary, one key away.

// entry is one selectable line of the detail's sidebar: a linked record, or
// a commit of the timeline.
type entry struct {
	role, id, title, status string
	note                    string // a review: what it examined, against the candidate
	change                  *versions.Change
	commit                  *versions.Commit
}

// openDetail shows a record above whatever is open: a card from the board,
// or a linked record from another detail. From a detail, a record already on
// the path is returned to, cutting what was opened above it, rather than
// opened again; o from an attempt always opens a layer, so Esc still returns
// to the attempt.
func (m *Model) openDetail(id string) {
	if len(m.stack) == 0 {
		m.workDepth = 0 // a record opened afresh returns to the board
	}
	if i := lastIndex(m.stack, id); i >= 0 && m.screen == detailScreen {
		m.stack = m.stack[:i+1]
		if len(m.stack) < m.workDepth {
			m.workDepth = 0 // the record o opened was cut, and its return with it
		}
	} else {
		m.stack = append(m.stack, id)
	}
	m.screen, m.side, m.dscroll, m.asOf, m.diff = detailScreen, -1, 0, "", ""
	m.leaveVersions()
	m.startAtEvidence()
}

// lastIndex is where id last stands on a path. o from an attempt can put a
// record on it twice, and the later one keeps the return to the attempt.
func lastIndex(path []string, id string) int {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == id {
			return i
		}
	}
	return -1
}

// crumbs is the detail's path, from where it was entered to the open record,
// with the attempts a record was opened from by o; long paths lose their
// start, never the open record.
func (m *Model) crumbs(w int) string {
	parts := []string{"board"}
	for i, id := range m.stack {
		if m.workDepth != 0 && i == m.workDepth-1 {
			parts = append(parts, map[bool]string{true: "dependencies", false: "attempts"}[m.workBack == depsScreen])
		}
		parts = append(parts, id)
	}
	s := safe(strings.Join(parts, " › "))
	if over := ansi.StringWidth(s) - w; over > 0 {
		s = ansi.TruncateLeft(s, over+1, "…")
	}
	return line(s, w)
}

// leaveDetail steps back: from a diff or a commit's content to now, then to
// the record opened before, then to the board.
func (m *Model) leaveDetail() {
	switch {
	case m.diff != "":
		m.diff, m.dscroll = "", 0
	case m.asOf != "":
		m.asOf, m.dscroll = "", 0
	case m.workDepth != 0 && len(m.stack) == m.workDepth:
		m.stack = m.stack[:len(m.stack)-1]
		m.screen, m.workDepth, m.side, m.dscroll = m.workBack, 0, -1, 0
	case len(m.stack) > 1:
		m.stack = m.stack[:len(m.stack)-1]
		m.side, m.dscroll = -1, 0
	default:
		m.stack, m.screen = nil, boardScreen
	}
}

func (m *Model) detailKey(k string) tea.Cmd {
	entries := m.entries()
	linked := slices.IndexFunc(entries, func(e entry) bool { return e.commit == nil && e.change == nil })
	changes := slices.IndexFunc(entries, func(e entry) bool { return e.change != nil })
	timeline := slices.IndexFunc(entries, func(e entry) bool { return e.commit != nil })
	switch k {
	case "tab":
		// Content, then the linked records, the changes, the timeline, then content.
		var next int
		switch {
		case m.side < 0:
			next = -1
			for _, first := range []int{linked, changes, timeline} {
				if first >= 0 {
					next = first
					break
				}
			}
			if next < 0 && !m.split() {
				next = 0 // a narrow terminal shows the sidebar even with nothing to select
			}
		case changes > m.side:
			next = changes
		case timeline > m.side:
			next = timeline
		default:
			next = -1
		}
		m.side = next
	case "up", "k", "down", "j", "pgup", "pgdown":
		if m.side < 0 {
			m.dscroll = m.moved(m.dscroll, k)
			m.clampScroll()
		} else if len(entries) != 0 {
			m.side = min(max(m.moved(m.side, k), 0), len(entries)-1)
		}
	case "enter":
		if m.side < 0 || m.side >= len(entries) {
			return nil
		}
		switch e := entries[m.side]; {
		case e.commit != nil && e.commit.Source != nil:
			m.asOf, m.diff, m.dscroll = e.commit.ID, "", 0
		case e.change != nil:
			// A rename's diff is read at its new path.
			path := e.change.Path
			if i := strings.LastIndex(path, " → "); i >= 0 {
				path = path[i+len(" → "):]
			}
			m.diff, m.asOf, m.dscroll = path, "", 0
		case e.commit == nil && e.id != "":
			m.openDetail(e.id)
		}
	case "w":
		// Hidden, the content takes the width, and a mouse selection of it
		// takes no sidebar text; Tab still reaches the sidebar, full width.
		if m.width >= wideWidth {
			m.sideHidden = !m.sideHidden
			m.clampScroll()
		}
	case "v":
		m.screen = versionsScreen
		m.leaveVersions()
	case "a", "f", "i":
		m.action(k)
	case "A":
		work := ""
		if r := m.openRecord(); r != nil && r.Type == "work" {
			work = m.openID()
		}
		m.openAttempts(work)
	case "R":
		m.launch()
	case "m":
		m.resolveConflict()
	case "e":
		return m.answer()
	}
	return nil
}

// moved applies a movement key to a position.
func (m *Model) moved(at int, k string) int {
	page := max(m.height-6, 1)
	switch k {
	case "up", "k":
		return at - 1
	case "down", "j":
		return at + 1
	case "pgup":
		return at - page
	}
	return at + page
}

// record returns the record a group shows on this board, if it has one.
func (m *Model) record(g *versions.Group) *project.Record {
	if v := m.shown(g); v != nil {
		return v.Record
	}
	return nil
}

// roles orders the sidebar's linked records.
var roles = []string{"plan", "review", "work", "needs", "needed by", "blocked by", "blocks", "part of", "member", "related"}

// linked lists the records connected to the open one, each once, by the
// first role that applies: its plans and reviews (records whose work names
// it), the work a plan or review belongs to, prerequisites both ways,
// blocking questions both ways, membership both ways, and relates_to
// either way. Roles come from fields, never from an ID or a folder.
func (m *Model) linked(id string, r *project.Record) []entry {
	var out []entry
	seen := map[string]bool{}
	add := func(role string, o *project.Record) {
		if seen[o.ID] {
			return
		}
		seen[o.ID] = true
		e := entry{role: role, id: o.ID, title: o.Title, status: o.Status}
		// A review says what it examined, and whether that is the work's
		// candidate: the fact G-044's review view starts from.
		if o.Type == "review" && o.Examined != "" {
			e.note = "examined " + o.Examined[:min(len(o.Examined), 7)]
			switch {
			case r == nil || r.Candidate == "":
			case strings.HasPrefix(o.Examined, r.Candidate) || strings.HasPrefix(r.Candidate, o.Examined):
				e.note += " = candidate"
			default:
				e.note += ", not the candidate"
			}
		}
		out = append(out, e)
	}
	has := func(list []string, id string) bool { return slices.Contains(list, id) }
	for i := range m.res.Groups {
		g := &m.res.Groups[i]
		o := m.record(g)
		if o == nil || g.ID == id {
			continue
		}
		var own project.Record
		if r != nil {
			own = *r
		}
		switch {
		case has(o.Work, id):
			add(o.Type, o)
		case has(own.Work, g.ID):
			add("work", o)
		case has(own.DependsOn, g.ID):
			add("needs", o)
		case has(o.DependsOn, id):
			add("needed by", o)
		case o.Type == "question" && has(o.Blocks, id):
			add("blocked by", o)
		case has(own.Blocks, g.ID):
			add("blocks", o)
		case has(o.Members, id):
			add("part of", o)
		case has(own.Members, g.ID):
			add("member", o)
		case has(own.RelatesTo, g.ID) || has(o.RelatesTo, id):
			add("related", o)
		}
	}
	slices.SortStableFunc(out, func(a, b entry) int { return slices.Index(roles, a.role) - slices.Index(roles, b.role) })
	return out
}

// histItem is one row of a version's history: a commit, or a note about the
// checkout's files or the record's status here.
type histItem struct {
	when, status, rest string
	commit             *versions.Commit
}

// historyItems lists what the history of the shown version holds, and any
// message about it (still reading, an error, nothing to list).
func (m *Model) historyItems(v *versions.Version) (items []histItem, message string) {
	status := "-" // deleted from the checkout's files
	if v.Record != nil {
		status = v.Record.Status
	}
	switch {
	case committed(v):
	case v.Change == "unknown":
		items = append(items, histItem{"uncommitted?", status, "whether these files differ from the commit is unknown", nil})
	default:
		items = append(items, histItem{"uncommitted", status, v.Change + " in this checkout's files", nil})
	}
	commit, path := historyAt(v)
	read, held := m.hist[commit+"\x00"+path]
	switch {
	case commit == "" && v.Source.Kind == "committed":
		message = "This branch deleted the record; an older row has its history."
	case commit == "":
		message = "No commit of this checkout holds the record yet."
	case !held:
		message = "reading…"
	case read.err != nil:
		message = "The history could not be read (r retries): " + read.err.Error()
	case len(read.commits) == 0:
		message = "No commit here changed this file."
	}
	// Merges are not listed, so the newest row need not be the record as it is
	// here: a merge may have set the status, or combined it with a newer commit
	// from another line. Say what is known, not which.
	if committed(v) && len(read.commits) != 0 && read.commits[0].Status != status {
		items = append(items, histItem{"here", status, "the record's status here; the newest commit below differs because merges are not listed", nil})
	}
	for i := range read.commits {
		c := &read.commits[i]
		items = append(items, histItem{c.When.Format("2006-01-02 15:04"), c.Status, c.ID[:min(len(c.ID), 7)] + "  " + c.Subject, c})
	}
	return items, message
}

// entries lists the sidebar's selectable lines in order: linked records,
// then the changed files, then the timeline's commits.
func (m *Model) entries() []entry {
	g := m.group()
	if g == nil {
		return nil
	}
	v := m.shown(g)
	var r *project.Record
	if v != nil {
		r = v.Record
	}
	out := m.linked(g.ID, r)
	out = append(out, m.changeEntries(v)...)
	if v != nil {
		items, _ := m.historyItems(v)
		for _, it := range items {
			if it.commit != nil {
				out = append(out, entry{commit: it.commit})
			}
		}
	}
	return out
}

// split reports the detail's content beside its sidebar: on a wide terminal
// that has not hidden the sidebar with w.
func (m *Model) split() bool { return m.width >= wideWidth && !m.sideHidden }

// detailBody draws the open record: its header, then its content beside
// the sidebar, or one of them at a time on a narrow terminal or with the
// sidebar hidden.
func (m *Model) detailBody(w, n int) []string {
	g := m.group()
	v := m.shown(g)
	head := m.detailHead(v, w)
	n -= len(head)
	var rows []string
	switch {
	case !m.split() && m.side >= 0:
		rows = m.sidebar(v, w, n)
	case !m.split():
		rows = m.content(v, w, n)
	default:
		lw := w * 11 / 20
		left, right := m.content(v, lw, n), m.sidebar(v, w-lw-3, n)
		for i := range n {
			rows = append(rows, left[i]+" │ "+right[i])
		}
	}
	return append(head, rows...)
}

// detailHead is the boxed header: ID and status or type, the title, and one
// line of metadata, wrapped as the width allows.
func (m *Model) detailHead(v *versions.Version, w int) []string {
	g := m.group()
	iw := max(w-4, 1)
	accent := lipgloss.NewStyle().Bold(true)
	var r *project.Record
	if v != nil {
		r = v.Record
	}
	var first, title string
	switch {
	case r == nil:
		first, title = g.ID+" · deleted", "The current state removes this record; v lists the places that still hold an older copy."
	case r.Type == "work":
		first, title = g.ID+" · "+r.Status, r.Title
		if i := slices.Index(statuses[:], r.Status); i >= 0 {
			accent = accents[i].Bold(true)
		}
	case r.Status != "":
		first, title = g.ID+" · "+r.Type+" · "+r.Status, r.Title
	default:
		first, title = g.ID+" · "+r.Type, r.Title
	}
	// A short terminal keeps the header to the ID and one title row, so
	// the content still has rows beneath it.
	rows, compact := 2, m.height < 16
	if compact {
		rows = 1
	}
	var inner []string
	if !compact {
		inner = append(inner, m.crumbs(iw))
	}
	inner = append(inner, bold(line(first, iw)))
	inner = append(inner, wrap(title, iw)[:min(len(wrap(title, iw)), rows)]...)
	if meta := m.detailMeta(g, v); meta != "" && !compact {
		inner = append(inner, wrap(meta, iw)[:min(len(wrap(meta, iw)), rows)]...)
	}
	// A candidate in review shows its standing and where the actions run.
	if !compact && r != nil && r.Type == "work" && r.Status == "review" {
		for _, row := range m.reviewRows(g, v) {
			inner = append(inner, wrap(row, iw)[:min(len(wrap(row, iw)), rows)]...)
		}
	}
	// An open question says where e would write its answer.
	if !compact && r != nil && r.Type == "question" && r.Status == "open" && m.backend.Edit != nil {
		row := m.answerRow()
		inner = append(inner, wrap(row, iw)[:min(len(wrap(row, iw)), rows)]...)
	}
	// Work shows its latest attempt, so a person back after closing the
	// board finds a run from the card they know.
	if !compact && r != nil && r.Type == "work" && m.backend.Attempts != nil {
		row := m.attemptRow(g.ID, r.Status)
		inner = append(inner, wrap(row, iw)[:min(len(wrap(row, iw)), rows)]...)
	}
	b := lipgloss.ThickBorder()
	edge := func(l, mid, r string) string { return accent.Render(l + strings.Repeat(mid, iw+2) + r) }
	side := accent.Render(b.Left)
	out := []string{edge(b.TopLeft, b.Top, b.TopRight)}
	for _, row := range inner {
		out = append(out, side+" "+row+" "+side)
	}
	return append(out, edge(b.BottomLeft, b.Bottom, b.BottomRight))
}

// detailMeta is the header's metadata line: the record's planning fields,
// its candidate, where its current state stands against the target and
// which places hold it, and when it was last written.
func (m *Model) detailMeta(g *versions.Group, v *versions.Version) string {
	if v == nil || v.Record == nil {
		return ""
	}
	r := v.Record
	var parts []string
	for _, p := range []string{r.Kind, r.Size} {
		if p != "" {
			parts = append(parts, p)
		}
	}
	if r.Priority != nil {
		parts = append(parts, fmt.Sprintf("P%d", *r.Priority))
	}
	if r.Candidate != "" {
		parts = append(parts, "candidate "+r.Candidate[:min(len(r.Candidate), 7)])
	}
	if len(r.Work) != 0 {
		parts = append(parts, "for "+strings.Join(r.Work, " "))
	}
	if r.Examined != "" {
		parts = append(parts, "examined "+r.Examined[:min(len(r.Examined), 7)])
	}
	if len(r.Blocks) != 0 {
		parts = append(parts, "blocks "+strings.Join(r.Blocks, " "))
	}
	states := currentStates(*g)
	if m.current() && len(states) > 1 {
		parts = append(parts, fmt.Sprintf("⑂ %d states (v)", len(states)))
	}
	for _, state := range states {
		if state[0] != v && !slices.Contains(state, v) {
			continue
		}
		switch {
		case !slices.ContainsFunc(state, committed):
			parts = append(parts, "uncommitted")
		case m.res.Target == "":
		case v.OnTarget:
			parts = append(parts, "on "+m.res.Target)
		default:
			parts = append(parts, "not on "+m.res.Target)
		}
		parts = append(parts, "same on "+held(state))
		break
	}
	if !m.current() {
		parts = append(parts, "on "+label(v.Source))
	}
	if r.Updated != nil {
		parts = append(parts, "updated "+r.Updated.Format("2006-01-02"))
	}
	return strings.Join(parts, " · ")
}

// content is the record's body rendered from Markdown, at a timeline
// commit when one is chosen, under a heading that counts the rows.
func (m *Model) content(v *versions.Version, w, n int) []string {
	all := m.contentRows(v, w)
	head := fmt.Sprintf("Content  %d-%d of %d", min(m.dscroll+1, len(all)), min(m.dscroll+n-1, len(all)), len(all))
	switch {
	case m.diff != "":
		head = fmt.Sprintf("Diff of %s (Esc returns to the content)  %d-%d of %d", m.diff, min(m.dscroll+1, len(all)), min(m.dscroll+n-1, len(all)), len(all))
	case m.asOf != "":
		head = fmt.Sprintf("Content as of %s (Esc returns to now)  %d-%d of %d", m.asOf[:min(len(m.asOf), 7)], min(m.dscroll+1, len(all)), min(m.dscroll+n-1, len(all)), len(all))
	}
	if m.side < 0 {
		head = hot(line("▶ "+head, w))
	} else {
		head = bold(line("  "+head, w))
	}
	return append([]string{head}, fit(all[min(m.dscroll, len(all)):], n-1, w)...)
}

// contentRows renders what the content pane shows: the record here, or at
// the chosen commit, whose bytes the history read holds.
func (m *Model) contentRows(v *versions.Version, w int) []string {
	if v == nil {
		return nil
	}
	if m.diff != "" {
		return m.diffContent(v, w)
	}
	if m.asOf != "" {
		commit, path := historyAt(v)
		if read, held := m.hist[commit+"\x00"+path]; held {
			if i := slices.IndexFunc(read.commits, func(c versions.Commit) bool { return c.ID == m.asOf }); i >= 0 && read.commits[i].Source != nil {
				return m.rendered(m.asOf+"\x00"+path, body(read.commits[i].Source), w)
			}
		}
		return wrapAll("That commit's content is no longer held; r rereads the history.", w)
	}
	if v.Record == nil {
		return wrapAll("Deleted here. Press v for the places that hold an older copy.", w)
	}
	return m.rendered(v.Revision, body(v.Record.Source), w)
}

// sidebar lists the linked records, the timeline and the sources, keeping
// the cursor's entry in view.
func (m *Model) sidebar(v *versions.Version, w, n int) []string {
	var rows []string
	at := -1 // the cursor's row
	entries := 0
	heading := func(text string) { rows = append(rows, bold(line(text, w))) }
	item := func(text string) {
		if entries == m.side {
			at = len(rows)
			rows = append(rows, hot(line("> "+text, w)))
		} else {
			rows = append(rows, line("  "+text, w))
		}
		entries++
	}
	heading("Linked")
	var r *project.Record
	if v != nil {
		r = v.Record
	}
	linked := m.linked(m.openID(), r)
	if len(linked) == 0 {
		rows = append(rows, line("  none", w))
	}
	for _, e := range linked {
		status := e.status
		if status == "" {
			status = "-"
		}
		title := ansi.Truncate(safe(e.title), max(w-2-max(10, ansi.StringWidth(e.role))-1-len(e.id)-2-ansi.StringWidth(safe(status))-2, 8), "…")
		item(fmt.Sprintf("%-10s %s  %s  %s", e.role, e.id, title, status))
		if e.note != "" {
			rows = append(rows, line("             "+e.note, w))
		}
	}
	m.changesSection(v, w, heading, item, func(text string) { rows = append(rows, wrap(text, w)...) })
	if v != nil {
		heading("Timeline on " + label(v.Source))
		items, message := m.historyItems(v)
		for _, it := range items {
			status := it.status
			if status == "" {
				status = "?"
			}
			text := fmt.Sprintf("%-16s  %-9s  %s", it.when, status, it.rest)
			if it.commit != nil {
				item(text)
			} else {
				rows = append(rows, line("  "+text, w))
			}
		}
		if message != "" {
			rows = append(rows, wrap("  "+message, w)...)
		}
	}
	heading("Sources")
	if g := m.group(); g != nil {
		rows = append(rows, wrapAll(sourcesText(g, m.res.Target)+"\n  v lists every version and place, and selects a workspace", w)...)
	}
	// Scroll the pane so the cursor's row is visible.
	off := 0
	if at >= n {
		off = at - n + 1
	}
	return fit(rows[min(off, len(rows)):], n, w)
}

// sourcesText names each current state and counts the older places.
func sourcesText(g *versions.Group, target string) string {
	var lines []string
	for _, state := range currentStates(*g) {
		lines = append(lines, "  current: "+stateText(state, target))
	}
	n := 0
	for _, v := range g.Versions {
		if v.Older != "" {
			n++
		}
	}
	if n != 0 {
		lines = append(lines, fmt.Sprintf("  older: %d %s an earlier state", n, plural(n, "place holds", "places hold")))
	}
	return strings.Join(lines, "\n")
}
