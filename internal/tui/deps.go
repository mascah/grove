package tui

import (
	"cmp"
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/mascah/grove/internal/deps"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/versions"
)

// The dependency view (G-161, layout B of G-166): the board's work as a list
// grouped by what connects it and indented by layer, the focused item's
// prerequisites and what it unlocks as trees beside it, or behind Tab on a
// narrow terminal, and a preview of an explicit selection bound to one
// checkout. It renders deps' interpretation, which grove deps prints too.
// Nothing here writes, starts work, or adds to a selection, and only the
// preview reads Git, when it opens or the board is re-read.

type previewMsg struct {
	gen  int
	view *deps.View
}

// depsRecords are the records the list shows: each record's current state in
// the current view, whose divergent cards the list marks as the board does,
// or a checkout's own records on its board.
func (m *Model) depsRecords() []*project.Record {
	var out []*project.Record
	src := m.boardSource()
	for i := range m.res.Groups {
		g := &m.res.Groups[i]
		if m.current() {
			if r := m.record(g); r != nil {
				out = append(out, r)
			}
			continue
		}
		for _, v := range g.Versions {
			if src != nil && v.Source == src && v.Record != nil {
				out = append(out, v.Record)
			}
		}
	}
	return out
}

// ponytail: rebuilt per key and frame, as the board's cards are; cache per
// result if projects grow large.
func (m *Model) depsOverview() (*deps.View, map[string]*project.Record) {
	records := m.depsRecords()
	byID := map[string]*project.Record{}
	for _, r := range records {
		byID[r.ID] = r
	}
	return deps.Overview(records, m.depsAll), byID
}

// depsRows lists the rows of v, the unfinished or every work, as the list
// shows them: each connected group in deps' order, then the work connected to
// no other row. size counts each group's rows.
func depsRows(v *deps.View) (rows []deps.Item, size map[int]int) {
	size = map[int]int{}
	for _, it := range v.Items {
		if !it.Outside {
			rows = append(rows, it)
			size[it.Group]++
		}
	}
	alone := func(it deps.Item) int { return min(2-size[it.Group], 1) } // 1 for a group of one
	slices.SortStableFunc(rows, func(a, b deps.Item) int { return cmp.Compare(alone(a), alone(b)) })
	return rows, size
}

// depsFocus is the focused row: the remembered one while it is listed,
// otherwise the first.
func (m *Model) depsFocus(rows []deps.Item) int {
	at := slices.IndexFunc(rows, func(it deps.Item) bool { return it.ID == m.depsAt })
	if at < 0 && len(rows) != 0 {
		at = 0
	}
	return at
}

// openDeps shows the dependency view, focused on the board's card when the
// list holds it.
func (m *Model) openDeps() {
	if m.res == nil {
		return
	}
	m.screen, m.depsAt, m.scroll = depsScreen, m.cardID, 0
}

func (m *Model) depsKey(k string) tea.Cmd {
	if m.previewing {
		m.scrollKey(k)
		return nil
	}
	v, _ := m.depsOverview()
	rows, _ := depsRows(v)
	at := m.depsFocus(rows)
	switch k {
	case "up", "k", "down", "j", "pgup", "pgdown":
		if len(rows) != 0 {
			m.depsAt = rows[min(max(m.moved(at, k), 0), len(rows)-1)].ID
		}
	case "space":
		if at < 0 {
			return nil
		}
		id := rows[at].ID
		if i := slices.Index(m.depsPicked, id); i >= 0 {
			m.depsPicked = slices.Delete(m.depsPicked, i, i+1)
		} else {
			m.depsPicked = append(m.depsPicked, id)
		}
	case "p":
		if len(m.depsPicked) == 0 {
			m.notice = "Space selects the work to preview"
			return nil
		}
		m.previewing, m.preview, m.previewErr, m.previewRev, m.scroll = true, nil, "", nil, 0
	case "h":
		m.depsAll = !m.depsAll
	case "tab":
		m.depsTree = !m.depsTree
	case "enter":
		if at >= 0 {
			m.openWork(rows[at].ID) // Esc from the record returns here
		}
	}
	return nil
}

// leaveDeps steps back: from the preview to the list, then to the board.
func (m *Model) leaveDeps() {
	if m.pending == "preview" {
		m.stop()
	}
	if m.previewing {
		m.previewing, m.preview, m.previewErr, m.scroll = false, nil, "", 0
		return
	}
	m.screen = boardScreen
}

// bound is the one checkout a preview reads: the board's, or in the current
// view the checkout Grove was started in. Edges from different sources are
// never combined into one order.
func (m *Model) bound() *versions.Source {
	if !m.current() {
		return m.boardSource()
	}
	for _, s := range m.live() {
		if s.GitDir == m.res.GitDir {
			return s
		}
	}
	return nil
}

// wantPreview computes the open preview from the bound checkout's records
// and reads what Git says of the candidates it names, when no other read is
// pending. A re-read of the board clears it, so it is computed again.
func (m *Model) wantPreview() tea.Cmd {
	if !m.previewing || m.preview != nil || m.previewErr != "" || m.pending != "" || m.done || m.res == nil {
		return nil
	}
	src := m.bound()
	switch {
	case src == nil:
		m.previewErr = "No checkout to bind the preview to: the board's checkout changed or was removed. Press b to choose one."
		return nil
	case !src.Valid:
		m.previewErr = label(src) + " cannot be read: " + sourceProblem(src)
		return nil
	}
	var records []*project.Record
	for _, g := range m.res.Groups {
		for _, v := range g.Versions {
			if v.Source == src && v.Record != nil {
				records = append(records, v.Record)
			}
		}
	}
	view, err := deps.Preview(records, m.depsPicked)
	if err != nil {
		m.previewErr = err.Error() + ". The preview reads " + label(src) + " only; b chooses another checkout."
		return nil
	}
	view.Compare(m.res, src)
	m.previewOn = fmt.Sprintf("%s at %s", label(src), src.Commit[:min(len(src.Commit), 7)])
	target, ancestry, root := m.res.Target, m.backend.Ancestry, src.Worktree
	return m.read("preview", func(ctx context.Context, gen int) tea.Msg {
		if ancestry != nil {
			view.Deliver(target, ancestry(ctx, root))
		}
		return previewMsg{gen, view}
	})
}

// gotPreview shows a preview, saying which selected or listed records changed
// since the one it replaces.
func (m *Model) gotPreview(msg previewMsg) {
	if msg.gen != m.gen || m.pending != "preview" {
		return
	}
	m.pending, m.cancel = "", nil
	rev := map[string]string{}
	var changed []string
	for _, it := range msg.view.Items {
		rev[it.ID] = it.Revision
		if old, ok := m.previewRev[it.ID]; ok && old != it.Revision {
			changed = append(changed, it.ID)
		}
	}
	if changed != nil {
		msg.view.Notes = append([]string{strings.Join(changed, ", ") + " changed since this preview was last read; what follows is the new reading."}, msg.view.Notes...)
	}
	m.preview, m.previewRev = msg.view, rev
	m.clampScroll()
}

// depsBody is the list with the focused tree beside it from wideWidth, or
// either alone below it (Tab swaps), under a summary and above the legend.
func (m *Model) depsBody(w, n int) []string {
	if m.previewing {
		return m.scrolled(m.previewRows(w), n, w)
	}
	if s := m.boardSource(); !m.current() && (s == nil || !s.Valid) {
		return fit(wrapAll("No board: "+m.contextProblem(s)+". Press b to view another checkout.", w), n, w)
	}
	v, byID := m.depsOverview()
	rows, size := depsRows(v)
	at := m.depsFocus(rows)
	connected, alone := 0, 0
	for g, n := range size {
		if n > 1 {
			connected++
		} else if g != 0 {
			alone++
		}
	}
	shape := fmt.Sprintf("%d connected %s, %d unconnected", connected, plural(connected, "group", "groups"), alone)
	head := fmt.Sprintf("Dependencies: %d unfinished work, %s", len(rows), shape)
	if m.depsAll {
		head = fmt.Sprintf("Dependencies: every work, %d, %s · h shows unfinished only", len(rows), shape)
	} else if collapsed := len(v.Items) - len(rows); collapsed != 0 {
		head += fmt.Sprintf(" · %d done or abandoned %s only in the trees · h shows every work", collapsed, plural(collapsed, "prerequisite", "prerequisites"))
	}
	top := wrap(head, w)
	for i := range top {
		top[i] = bold(top[i])
	}
	legend := wrap("← needs · → unlocks · ✓ done · ✗ abandoned · ? open question · ● selected · indent: layer (equal: no declared order)", w)
	area := n - len(top) - len(legend)
	blockedBy := openBlocks(byID)
	list := func(w int) []string {
		if len(rows) == 0 {
			return fit([]string{line("  No unfinished work. h shows every work.", w)}, area, w)
		}
		tags := map[string]string{}
		columns, shelf := m.cards()
		for _, c := range slices.Concat(append(columns[:], shelf)...) {
			tags[c.id] = c.tag
		}
		var lines []string
		focus, group := 0, 0
		for i, it := range rows {
			switch {
			case size[it.Group] > 1 && it.Group != group:
				lines = append(lines, bold(line(fmt.Sprintf("Connected · %d work", size[it.Group]), w)))
			case size[it.Group] == 1 && (i == 0 || size[rows[i-1].Group] > 1):
				lines = append(lines, bold(line(fmt.Sprintf("Unconnected · %d work: no edge to another row", alone), w)))
			}
			group = it.Group
			if i == at {
				focus = len(lines)
			}
			picked := " "
			if slices.Contains(m.depsPicked, it.ID) {
				picked = "●"
			}
			text := picked + " " + strings.Repeat("  ", min(it.Layer, w/16)) + it.ID + " " + it.Status
			for _, q := range blockedBy[it.ID] {
				text += " ? " + q
			}
			if t := tags[it.ID]; t != "" {
				text += " [" + t + "]"
			}
			lines = append(lines, mark(i == at, text+"  "+it.Title, w))
		}
		return window(lines, focus, area-2, 1, area, w)
	}
	tree := func(w int) []string {
		if at < 0 {
			return fit(nil, area, w)
		}
		place := "unconnected"
		if n := size[rows[at].Group]; n > 1 {
			place = fmt.Sprintf("connected with %d other listed work", n-1)
		}
		out := m.treeRows(rows[at], place, byID, blockedBy, w)
		if len(out) > area {
			out = append(out[:area-1], line(fmt.Sprintf("… %d more rows", len(out)-area+1), w))
		}
		return fit(out, area, w)
	}
	var body []string
	switch {
	case w >= wideWidth:
		lw := w * 9 / 20
		left, right := list(lw), tree(w-lw-3)
		for i := range area {
			body = append(body, left[i]+" │ "+right[i])
		}
	case m.depsTree:
		body = tree(w)
	default:
		body = list(w)
	}
	return slices.Concat(top, body, legend)
}

// openBlocks maps each record to the open questions blocking it.
func openBlocks(byID map[string]*project.Record) map[string][]string {
	blockedBy := map[string][]string{}
	for _, id := range slices.Sorted(maps.Keys(byID)) {
		if r := byID[id]; r.Type == "question" && r.Status == "open" {
			for _, b := range r.Blocks {
				blockedBy[b] = append(blockedBy[b], id)
			}
		}
	}
	return blockedBy
}

// treeRows describes one row: its title and place, then what it needs, down
// to the work with no prerequisites, what it unlocks, and the open questions
// blocking it. A record met again is written (shown above), never expanded
// twice.
func (m *Model) treeRows(it deps.Item, place string, byID map[string]*project.Record, blockedBy map[string][]string, w int) []string {
	unlocks := map[string][]string{}
	for _, id := range slices.Sorted(maps.Keys(byID)) {
		if r := byID[id]; r.Type == "work" && (m.depsAll || r.Status == "proposed" || r.Status == "active" || r.Status == "review") {
			for _, p := range r.DependsOn {
				unlocks[p] = append(unlocks[p], id)
			}
		}
	}
	name := func(id string) string {
		r := byID[id]
		if r == nil {
			return id + " (not among the records read)"
		}
		text := map[string]string{"done": "✓ ", "abandoned": "✗ "}[r.Status] + id + " " + r.Status
		for _, q := range blockedBy[id] {
			text += " ? " + q
		}
		return text + "  " + r.Title
	}
	needs := func(id string) []string {
		if r := byID[id]; r != nil {
			return r.DependsOn
		}
		return nil
	}
	out := []string{bold(line(fmt.Sprintf("%s %s · layer %d · %s", it.ID, it.Status, it.Layer, place), w))}
	out = append(out, wrap(it.Title, w)...)
	section := func(title string, next func(string) []string) {
		out = append(out, line("", w), bold(line(title, w)))
		for _, r := range branches(it.ID, next, name) {
			out = append(out, line(r, w))
		}
	}
	section("← Needs", needs)
	section("→ Unlocks", func(id string) []string { return unlocks[id] })
	if qs := blockedBy[it.ID]; qs != nil {
		out = append(out, line("", w), bold(line("? Blocked by open questions, not by work", w)))
		for _, q := range qs {
			out = append(out, line(name(q), w))
		}
	}
	return out
}

// branches draws the tree below root that next gives, each record once.
func branches(root string, next func(string) []string, name func(string) string) []string {
	rows := []string{root}
	seen := map[string]bool{root: true}
	var walk func(id, indent string)
	walk = func(id, indent string) {
		kids := next(id)
		for i, k := range kids {
			branch, more := "├─ ", "│  "
			if i == len(kids)-1 {
				branch, more = "└─ ", "   "
			}
			if seen[k] {
				rows = append(rows, indent+branch+k+" (shown above)")
				continue
			}
			seen[k] = true
			rows = append(rows, indent+branch+name(k))
			walk(k, indent+more)
		}
	}
	walk(root, "")
	if len(rows) == 1 {
		rows[0] += ": none"
	}
	return rows
}

// previewRows is the selection preview: the selection as marked, its order
// through transitive prerequisites, the prerequisites outside it, which are
// never added, the questions, delivery and notes, bound to one checkout.
func (m *Model) previewRows(w int) []string {
	rows := []string{bold(line("Selection preview", w))}
	switch {
	case m.previewErr != "":
		return append(rows, wrapAll(m.previewErr, w)...)
	case m.preview == nil:
		return append(rows, line("Reading "+strings.Join(m.depsPicked, " ")+" and what Git says of their candidates…", w))
	}
	v := m.preview
	target := "no integration target"
	if m.res.Target != "" {
		target = "target " + m.res.Target
	}
	rows = append(rows, wrapAll("Bound to "+m.previewOn+" · "+target+"\nSelected: "+strings.Join(v.Selected, " ")+" (as marked)\nOrder:    "+strings.Join(v.Order, " → "), w)...)
	list := func(ids []string) string {
		if len(ids) == 0 {
			return "nothing"
		}
		return strings.Join(ids, " ")
	}
	item := func(n string, it deps.Item) {
		rows = append(rows, line(fmt.Sprintf(" %-3s %-7s %-9s %s", n, it.ID, it.Status, it.Title), w))
		facts := "needs " + list(it.Needs)
		if it.Outside {
			facts = "needed by " + list(it.NeededBy) + " · " + facts
		}
		if it.Delivery != "" {
			facts += " · " + it.Delivery
		} else {
			facts += " · delivery not read"
		}
		for _, r := range wrap(facts, w-5) {
			rows = append(rows, "     "+r)
		}
	}
	rows = append(rows, line("", w), bold(line(" In order", w)))
	outside := false
	for i, it := range v.Items {
		if it.Outside && !outside {
			outside = true
			rows = append(rows, line("", w), bold(line(" Outside the selection, not added", w)))
		}
		n := ""
		if !it.Outside {
			n = fmt.Sprint(i + 1)
		}
		item(n, it)
	}
	rows = append(rows, line("", w))
	if len(v.Questions) == 0 {
		rows = append(rows, line(" Questions: none open blocks the selection or its prerequisites.", w))
	} else {
		rows = append(rows, bold(line(" Open questions", w)))
		for _, q := range v.Questions {
			rows = append(rows, wrap(fmt.Sprintf(" ? %s blocks %s  %s", q.ID, list(q.Blocks), q.Title), w)...)
		}
	}
	rows = append(rows, line("", w))
	for _, n := range append(v.Notes, "Order comes from depends_on only; membership and priority never change it.", "A preview adds no work, starts nothing, and authorizes nothing.") {
		rows = append(rows, wrap(" · "+n, w)...)
	}
	return rows
}
