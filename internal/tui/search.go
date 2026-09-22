package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// Search lists every record of the project in its current state, of every
// type, including what the board hides: Abandoned, Done beyond its page,
// and knowledge that is never a card (pages, terms, decisions, questions,
// plans, reviews). It matches typed text against the ID, type, status and
// title, never the body, and Enter opens the detail.

type hit struct{ id, kind, status, title string }

// hits lists the records matching the query, in the inspection's order.
func (m *Model) hits() []hit {
	q := strings.ToLower(strings.TrimSpace(m.query))
	var out []hit
	for i := range m.res.Groups {
		g := &m.res.Groups[i]
		r := m.record(g)
		if r == nil {
			continue
		}
		h := hit{g.ID, r.Type, r.Status, r.Title}
		if h.status == "" {
			h.status = "-"
		}
		if q == "" || strings.Contains(strings.ToLower(g.ID+" "+r.Type+" "+h.status+" "+r.Title), q) {
			out = append(out, h)
		}
	}
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

// searchBody is the query line with the count, then the matching records.
func (m *Model) searchBody(w, n int) []string {
	hits := m.hits()
	total := 0
	for i := range m.res.Groups {
		if m.record(&m.res.Groups[i]) != nil {
			total++
		}
	}
	count := fmt.Sprintf("%d of %d records", len(hits), total)
	query := line("/ "+m.query+"▏", w-len(count)-2) + "  " + count
	head := []string{bold(query), line("Every record, of every type, in its current state; type to filter by ID, type, status or title.", w)}
	rows := make([]string, len(hits))
	for i, h := range hits {
		rows[i] = mark(i == m.hit, fmt.Sprintf("%-6s  %-9s  %-9s  %s", h.id, h.kind, h.status, h.title), w)
	}
	if len(hits) == 0 {
		rows = []string{line("  no record matches", w)}
	}
	n -= len(head)
	return append(head, window(rows, m.hit, n-2, 1, n, w)...)
}
