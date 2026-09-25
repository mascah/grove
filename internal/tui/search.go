package tui

import (
	"fmt"
	"path"
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/mascah/grove/internal/handoff"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/versions"
)

// Search lists every record of the project in its current state, of every
// type, including what the board hides: Abandoned, Done beyond its page,
// and knowledge that is never a card (pages, terms, decisions, questions,
// plans, reviews). Typed text is matched by tier (G-153): the ID, type,
// status and title; a link in the body to that project path or under it; a
// code span naming it; then the body's text. A record hits once, at its
// first tier, with the line that matched; hits sort by tier, then in the
// inspection's order, with no score. In the current view each current
// state of a diverging record is matched and shown on its own, since none
// replaces the others. Enter opens the detail.

type hit struct {
	id, kind, status, title string
	tier                    int
	where, snippet          string // where: the places holding a diverging state
}

var tiers = [...]string{"title", "link", "code span", "text"}

// mentionsOf is v's links and code spans, parsed once per revision and path.
func (m *Model) mentionsOf(v *versions.Version) []handoff.Mention {
	key := v.Revision + "\x00" + v.Path
	if ms, ok := m.mentions[key]; ok {
		return ms
	}
	if m.mentions == nil {
		m.mentions = map[string][]handoff.Mention{}
	}
	ms := handoff.Mentions(v.Record)
	m.mentions[key] = ms
	return ms
}

// pathTier finds a mention of the project path p: a link to p or under it
// (tier 1), else a code span naming it (tier 2); -1 when none does.
func pathTier(p string, ms []handoff.Mention) (int, string) {
	if p = path.Clean(strings.TrimPrefix(p, "./")); p == "." || p == ".." || strings.HasPrefix(p, "../") || strings.HasPrefix(p, "/") {
		return -1, ""
	}
	for _, mn := range ms {
		if mn.Path != "" && (mn.Path == p || strings.HasPrefix(mn.Path, p+"/")) {
			return 1, mn.Line
		}
	}
	for _, mn := range ms {
		if mn.Span != "" && names(mn.Span, p) {
			return 2, mn.Line
		}
	}
	return -1, ""
}

// names reports whether a code span names the path p: it is p, or a
// trailing run of p's components, so `ws.rs` names crates/server/src/ws.rs.
func names(span, p string) bool {
	s := path.Clean(strings.TrimPrefix(span, "./"))
	if s == "." || s == ".." || strings.HasPrefix(s, "../") || strings.HasPrefix(s, "/") {
		return false
	}
	return p == s || strings.HasSuffix(p, "/"+s)
}

// match is the first tier at which r answers the trimmed query, with the
// line that matched, or -1.
func match(q string, r *project.Record, ms []handoff.Mention) (int, string) {
	lq := strings.ToLower(q)
	status := r.Status
	if status == "" {
		status = "-"
	}
	if strings.Contains(strings.ToLower(r.ID+" "+r.Type+" "+status+" "+r.Title), lq) {
		return 0, ""
	}
	if tier, snippet := pathTier(q, ms); tier >= 0 {
		return tier, snippet
	}
	for _, l := range strings.Split(body(r.Source), "\n") {
		if strings.Contains(strings.ToLower(l), lq) {
			return 3, strings.TrimSpace(l)
		}
	}
	return -1, ""
}

// searched is what search and the review listing look through: every
// current state in the current view, or the record one checkout's board
// shows, each with where it is held when the record diverges.
func (m *Model) searched(g *versions.Group) (vs []*versions.Version, where []string) {
	if !m.current() {
		if v := m.shown(g); v != nil && v.Record != nil {
			return []*versions.Version{v}, []string{""}
		}
		return nil, nil
	}
	states := currentStates(*g)
	for _, state := range states {
		if state[0].Record == nil {
			continue
		}
		vs = append(vs, state[0])
		if len(states) > 1 {
			where = append(where, "on "+heldBy(state))
		} else {
			where = append(where, "")
		}
	}
	return vs, where
}

// hits lists the records matching the query, by tier, then in the
// inspection's order.
func (m *Model) hits() []hit {
	q := strings.TrimSpace(m.query)
	var out []hit
	for i := range m.res.Groups {
		g := &m.res.Groups[i]
		vs, where := m.searched(g)
		for j, v := range vs {
			r := v.Record
			tier, snippet := match(q, r, m.mentionsOf(v))
			if tier < 0 {
				continue
			}
			h := hit{g.ID, r.Type, r.Status, r.Title, tier, where[j], snippet}
			if h.status == "" {
				h.status = "-"
			}
			out = append(out, h)
		}
	}
	slices.SortStableFunc(out, func(a, b hit) int { return a.tier - b.tier })
	return out
}

// searchKey handles every key while the search shows: printable text goes
// into the query, so the letters that mean something elsewhere do not.
func (m *Model) searchKey(msg tea.KeyPressMsg) {
	switch k := msg.String(); k {
	case "esc":
		m.screen, m.query = m.back, ""
	case "enter":
		if hits := m.hits(); m.hit < len(hits) {
			m.openDetail(hits[m.hit].id)
			m.query = ""
		}
	case "up", "down", "pgup", "pgdown":
		m.hit = m.moved(m.hit, k)
	case "backspace":
		if n := len(m.query); n != 0 {
			_, size := lastRune(m.query)
			m.query = m.query[:n-size]
		}
		m.hit = 0
	default:
		if msg.Text != "" {
			m.query += msg.Text
			m.hit = 0
		}
	}
	m.hit = min(max(m.hit, 0), max(len(m.hits())-1, 0))
}

func lastRune(s string) (rune, int) {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i]&0xC0 != 0x80 {
			return rune(s[i]), len(s) - i
		}
	}
	return 0, len(s)
}

// searchBody is the query line with the count, then the matching records;
// with a query each hit has a second row saying why it matched.
func (m *Model) searchBody(w, n int) []string {
	hits := m.hits()
	total := 0
	for i := range m.res.Groups {
		if m.record(&m.res.Groups[i]) != nil {
			total++
		}
	}
	found := map[string]bool{}
	for _, h := range hits {
		found[h.id] = true
	}
	count := fmt.Sprintf("%d of %d records", len(found), total)
	query := line("/ "+m.query+"▏", w-len(count)-2) + "  " + count
	head := []string{bold(query), line("Every record in its current state, by ID, type, status or title, then by link, code span or body text.", w)}
	unit := 1
	if strings.TrimSpace(m.query) != "" {
		unit = 2
	}
	var rows []string
	for i, h := range hits {
		rows = append(rows, mark(i == m.hit, fmt.Sprintf("%-6s  %-9s  %-9s  %s", h.id, h.kind, h.status, h.title), w))
		if unit == 2 {
			why := tiers[h.tier]
			if h.where != "" {
				why += " · " + h.where
			}
			if h.snippet != "" {
				why += ": " + h.snippet
			}
			rows = append(rows, line("          "+why, w))
		}
	}
	if len(hits) == 0 {
		rows, unit = []string{line("  no record matches", w)}, 1
	}
	n -= len(head)
	return append(head, window(rows, m.hit, (n-2)/unit, unit, n, w)...)
}
