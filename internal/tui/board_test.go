package tui

import (
	"fmt"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/mascah/grove/internal/versions"
)

// manyDone is a current view with two open cards, twelve done ones written
// on successive days (out of ID order), and one abandoned.
func manyDone(fx fixture) *versions.Result {
	var vs []versions.Version
	for _, s := range []*versions.Source{fx.cMain, fx.main} {
		vs = append(vs, version(s, "W-001", "Inspect records", "active"), version(s, "W-002", "Create records", "proposed"))
		for i := 1; i <= 12; i++ {
			v := version(s, fmt.Sprintf("W-%03d", 200+i), fmt.Sprintf("Finished %d", i), "done")
			u := time.Date(2026, 9, (i*5)%12+1, 0, 0, 0, 0, time.UTC)
			p := 3
			v.Record.Updated, v.Record.Kind, v.Record.Size, v.Record.Priority, v.Record.Candidate = &u, "feature", "medium", &p, "abcdef0123"
			vs = append(vs, v)
		}
		vs = append(vs, version(s, "W-300", "Gave up", "abandoned"))
	}
	res := result(fx.main, fx.sources(), vs...)
	res.Target = "main"
	return res
}

// The Done column shows the page of most recently written cards and counts
// the rest, which focus never reaches; the same rule bounds a checkout's
// board.
func TestDoneIsBoundedToItsPageNewestFirst(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	for _, size := range [][2]int{{120, 36}, {80, 24}, {40, 10}} {
		m := open(t, &fake{res: manyDone(fx)}, size[0], size[1])
		columns, _, older := m.bounded()
		per := m.pageSize(doneColumn)
		if len(columns[doneColumn]) != per || older != 12-per {
			t.Fatalf("%v: %d shown, %d older, page %d", size, len(columns[doneColumn]), older, per)
		}
		// Newest first: day 12 (W-207), day 11 (W-202), day 10 (W-209).
		if got := ids(columns[doneColumn][:3]); got != "W-207,W-202,W-209" {
			t.Fatalf("%v: order %s", size, got)
		}
		press(m, "right", "right", "right")
		s := plain(m)
		if !strings.Contains(s, fmt.Sprintf("+ %d older · / to search", older)) || !strings.Contains(s, "Done 12") {
			t.Fatalf("%v: the cut is not counted:\n%s", size, s)
		}
		if size[0] >= wideWidth && !strings.Contains(s, fmt.Sprintf("Done 12 · %d recent", per)) {
			t.Fatalf("%v: the heading says how many show:\n%s", size, s)
		}
		for range 20 {
			press(m, "down")
		}
		if m.cardID != columns[doneColumn][per-1].id || !strings.Contains(plain(m), "▶"+m.cardID) {
			t.Fatalf("%v: focus went past the page to %s", size, m.cardID)
		}
		if size[0] >= wideWidth && !strings.Contains(s, "done 2026-09-11 · abcdef0") {
			t.Fatalf("%v: a done card shows its date and candidate:\n%s", size, s)
		}
		// A smaller terminal shrinks the page: the focus stays on it.
		m.Update(tea.WindowSizeMsg{Width: size[0], Height: 10})
		columns, _, _ = m.bounded()
		if m.cardID != columns[doneColumn][m.pageSize(doneColumn)-1].id || !strings.Contains(plain(m), "▶"+m.cardID) {
			t.Fatalf("%v shrunk: focus stayed past the page on %s", size, m.cardID)
		}
	}
}

// Abandoned work is hidden until a shows it, and the count is on screen
// either way; hiding it moves focus out of its column.
func TestAbandonedIsHiddenUntilAsked(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	m := open(t, &fake{res: manyDone(fx)}, 120, 30)
	if s := plain(m); strings.Contains(s, "W-300") || !strings.Contains(s, "Abandoned 1 hidden · a shows") {
		t.Fatalf("abandoned shown by default:\n%s", s)
	}
	press(m, "right", "right", "right", "right", "right")
	if m.col != doneColumn {
		t.Fatalf("→ reached the hidden column %d", m.col)
	}
	press(m, "a", "right")
	if s := plain(m); m.col != 4 || m.cardID != "W-300" || !strings.Contains(s, "▶W-300") || !strings.Contains(s, "Abandoned 1") || strings.Contains(s, "hidden") {
		t.Fatalf("a should show the column and let → reach it (col %d, card %s):\n%s", m.col, m.cardID, s)
	}
	press(m, "a")
	if s := plain(m); m.col != doneColumn || strings.Contains(s, "W-300") || !strings.Contains(s, "hidden") {
		t.Fatalf("a again hides it and moves focus (col %d):\n%s", m.col, s)
	}
	// A refresh keeps the choice, and a hidden focus is re-settled.
	press(m, "a", "right")
	f := &fake{res: manyDone(fx)}
	m.backend = f.backend()
	deliver(m, press(m, "r"))
	if !m.showAbandoned || m.cardID != "W-300" {
		t.Fatalf("after refresh: shown %v, card %s", m.showAbandoned, m.cardID)
	}
}

// Cards keep their metadata row until the terminal is too short for a full
// card and Done's footer, and their titles wrap to two rows at most.
func TestCardsFitTheirColumn(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	long := version(fx.main, "W-001", strings.Repeat("word ", 30), "proposed")
	long.Record.Kind = "feature"
	m := open(t, &fake{res: result(fx.main, fx.sources(), long)}, 120, 30)
	s := plain(m)
	if !strings.Contains(s, "word…") || !strings.Contains(s, "feature") || strings.Count(s, "word") > 12 {
		t.Fatalf("a long title is cut at two rows with its metadata below:\n%s", s)
	}
	for _, h := range []int{10, 11, 12} {
		m.width, m.height = 40, h
		rows := strings.Split(m.render(), "\n")
		if len(rows) != h {
			t.Fatalf("height %d: %d rows", h, len(rows))
		}
		if want := h >= 12; strings.Contains(plain(m), "feature") != want {
			t.Fatalf("height %d: metadata shown = %v:\n%s", h, !want, plain(m))
		}
	}
}

// ←/→ skip empty columns, which are still drawn: with nothing active or in
// review, Done is one key from Proposed, and past the last column with
// cards the focus stays.
func TestColumnKeysSkipEmptyColumns(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	var vs []versions.Version
	for _, s := range []*versions.Source{fx.cMain, fx.main} {
		vs = append(vs, version(s, "W-001", "Open work", "proposed"), version(s, "W-002", "Finished", "done"))
	}
	for _, w := range []int{120, 80} {
		m := open(t, &fake{res: result(fx.main, fx.sources(), vs...)}, w, 24)
		if press(m, "right"); m.col != doneColumn || m.cardID != "W-002" {
			t.Fatalf("width %d: → from Proposed lands on Done: column %d card %s", w, m.col, m.cardID)
		}
		if press(m, "l"); m.col != doneColumn {
			t.Fatalf("width %d: past the last column with cards the focus stays: column %d", w, m.col)
		}
		if press(m, "h"); m.col != 0 || m.cardID != "W-001" {
			t.Fatalf("width %d: ← returns to Proposed: column %d", w, m.col)
		}
		if s := plain(m); w >= wideWidth && (!strings.Contains(s, "Active 0") || !strings.Contains(s, "Review 0")) {
			t.Fatalf("empty columns are still drawn:\n%s", s)
		}
	}
}
