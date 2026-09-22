package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mascah/grove/internal/versions"
)

func typeText(m *Model, text string) {
	for _, r := range text {
		m.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
	}
}

// / lists every record in its current state, whatever the board hides or
// never shows, filtered as text is typed; Enter opens the detail.
func TestSearchReachesEveryRecord(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	res := manyDone(fx)
	groups := map[string][]versions.Version{}
	for _, s := range []*versions.Source{fx.cMain, fx.main} {
		page := version(s, "G-020", "How boards work", "")
		page.Record.Type = "page"
		page.Record.Source = []byte("---\nid: G-020\n---\n\n# Boards\n\nA page.\n")
		term := version(s, "G-021", "Candidate", "settled")
		term.Record.Type = "term"
		hostile := version(s, "G-022", "Sneaky \x1b[2J\x1b]52;c;aGk=\x07 title", "resolved")
		hostile.Record.Type = "question"
		for _, v := range []versions.Version{page, term, hostile} {
			groups[v.Record.ID] = append(groups[v.Record.ID], v)
		}
	}
	for _, id := range []string{"G-020", "G-021", "G-022"} {
		res.Groups = append(res.Groups, versions.Group{ID: id, Versions: groups[id]})
	}
	f := &fake{res: res}
	m := open(t, f, 120, 30)
	press(m, "/")
	s := plain(m)
	if m.screen != searchScreen || !strings.Contains(s, "18 of 18 records") || !strings.Contains(s, "> W-001   work       active     Inspect records") ||
		!strings.Contains(s, "W-300   work       abandoned  Gave up") || !strings.Contains(s, "G-020   page       -          How boards work") || !strings.Contains(s, "G-021   term       settled    Candidate") {
		t.Fatalf("search should list everything:\n%s", s)
	}
	// Letters filter; they do not quit, refresh, or open other screens.
	typeText(m, "q r s b")
	if s = plain(m); m.screen != searchScreen || f.inspects != 1 || !strings.Contains(s, "/ q r s b▏") || !strings.Contains(s, "0 of 18 records") || !strings.Contains(s, "no record matches") {
		t.Fatalf("typed letters must only filter:\n%s", s)
	}
	for range 7 {
		m.Update(tea.KeyPressMsg{Code: tea.KeyBackspace})
	}
	typeText(m, "gave")
	if s = plain(m); !strings.Contains(s, "1 of 18 records") || !strings.Contains(s, "> W-300") {
		t.Fatalf("Backspace edits and the hidden abandoned card is found:\n%s", s)
	}
	press(m, "enter")
	if s = plain(m); m.screen != detailScreen || m.openID() != "W-300" || !strings.Contains(s, "W-300 · abandoned") {
		t.Fatalf("Enter opens the detail:\n%s", s)
	}
	if press(m, "esc"); m.screen != boardScreen || m.cardID != "W-002" {
		t.Fatalf("Esc returns to the board where it was: screen %d card %s", m.screen, m.cardID)
	}
	// Done work beyond the page, a page by its type, a term by its title.
	press(m, "/")
	typeText(m, "finished 1")
	if s = plain(m); !strings.Contains(s, "4 of 18 records") { // 1, 10, 11, 12
		t.Fatalf("older done work is searchable:\n%s", s)
	}
	press(m, "down", "down", "down", "enter")
	if m.openID() != "W-212" {
		t.Fatalf("↓ walks the hits: opened %s", m.openID())
	}
	press(m, "esc", "/")
	typeText(m, "PAGE")
	press(m, "enter")
	if s = plain(m); m.openID() != "G-020" || !strings.Contains(s, "# Boards") {
		t.Fatalf("a page opens in the detail, matched by type whatever the case:\n%s", s)
	}
	press(m, "esc", "/")
	typeText(m, "candid")
	press(m, "enter")
	if m.openID() != "G-021" {
		t.Fatalf("a term opens: %s", m.openID())
	}
	// A hostile title is shown escaped.
	press(m, "esc", "/")
	typeText(m, "sneaky")
	raw := m.render()
	if strings.Contains(raw, "\x1b[2J") || strings.Contains(raw, "\x1b]52") || !strings.Contains(ansi.Strip(raw), `Sneaky \x1b[2J\x1b]52;c;aGk=\a title`) {
		t.Fatalf("a hostile title reached the terminal:\n%q", raw)
	}
	// Esc closes the search without opening anything, and clears the query.
	press(m, "esc")
	if m.screen != boardScreen || m.query != "" {
		t.Fatalf("Esc: screen %d query %q", m.screen, m.query)
	}
	press(m, "/")
	if !strings.Contains(plain(m), "18 of 18 records") {
		t.Fatal("a new search starts empty")
	}
	if press(m, "ctrl+c") == nil {
		t.Fatal("Ctrl-C must still interrupt")
	}
}
