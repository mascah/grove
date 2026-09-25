package tui

import (
	"context"
	"slices"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/mascah/grove/internal/deps"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/versions"
)

// backlog is G-165's synthetic unfinished backlog on main, in the checkout
// and on the branch: a chain, a shared prerequisite (W-02), a convergence
// (W-05), unrelated work (W-09, W-11), a review candidate not on main (W-07),
// an abandoned prerequisite (W-10), a done one with its candidate (W-01) and
// an open question blocking W-11.
func (f fixture) backlog() *versions.Result {
	type w struct{ id, status, title, candidate string }
	items := []w{
		{"W-01", "done", "Scaffold the ledger project and its test harness", "c1"},
		{"W-02", "proposed", "Store accounts, balances and currencies in one shared ledger file", ""},
		{"W-03", "active", "Import bank statements from CSV and OFX exports into the ledger", ""},
		{"W-04", "proposed", "Categorize transactions with editable, ordered matching rules", ""},
		{"W-05", "proposed", "Monthly budget report with category totals, carry-over and warnings", ""},
		{"W-06", "proposed", "Export the monthly report as CSV and printable HTML", ""},
		{"W-07", "review", "Reconcile imported balances against statement closing balances", "c7"},
		{"W-08", "proposed", "Detect recurring transactions and predict next month's bills", ""},
		{"W-09", "proposed", "Dark theme for the report viewer", ""},
		{"W-10", "abandoned", "Legacy OFX 1.x parser", ""},
		{"W-11", "active", "Fix rounding of foreign-currency amounts on import", ""},
		{"W-13", "proposed", "Send budget alerts when a category passes its limit", ""},
	}
	needs := map[string][]string{"W-02": {"W-01"}, "W-03": {"W-02"}, "W-04": {"W-02"}, "W-05": {"W-03", "W-04"},
		"W-06": {"W-05"}, "W-07": {"W-03"}, "W-08": {"W-04", "W-10"}, "W-13": {"W-05"}}
	var vs []versions.Version
	for _, s := range []*versions.Source{f.cMain, f.main} {
		for _, it := range items {
			v := version(s, it.id, it.title, it.status)
			v.Record.DependsOn, v.Record.Candidate = needs[it.id], it.candidate
			v.OnTarget = it.id != "W-07" // its candidate is on a branch
			vs = append(vs, v)
		}
		q := version(s, "Q-12", "Which currencies round half-even?", "open")
		q.Record.Blocks, q.OnTarget = []string{"W-11"}, true
		vs = append(vs, q)
	}
	res := result(f.main, []*versions.Source{f.cMain, f.main}, vs...)
	res.Target = "main"
	return res
}

// ancestry answers from a table, logging each question: c1 is on main and
// in HEAD, c7 in neither.
func ancestry(log *[]string) func(context.Context, string) func(string, string) (bool, error) {
	return func(_ context.Context, root string) func(string, string) (bool, error) {
		return func(commit, ref string) (bool, error) {
			*log = append(*log, root+" "+commit+" "+ref)
			return commit == "c1", nil
		}
	}
}

func openDeps(t *testing.T, w, h int) (*Model, *fake, *[]string) {
	t.Helper()
	f := &fake{res: newFixture().backlog()}
	var asked []string
	b := f.backend()
	b.Ancestry = ancestry(&asked)
	m := New(t.Context(), "/repo/.", b)
	m.Update(tea.WindowSizeMsg{Width: w, Height: h})
	deliver(m, m.Init())
	press(m, "g")
	return m, f, &asked
}

// focusOn moves the list's focus to id.
func focusOn(t *testing.T, m *Model, id string) {
	t.Helper()
	for _, k := range []string{"down", "up"} {
		for range 20 {
			if m.depsAt == id {
				return
			}
			press(m, k)
		}
	}
	t.Fatalf("never focused %s; at %s", id, m.depsAt)
}

func TestDepsListAndTree(t *testing.T) {
	m, _, asked := openDeps(t, 120, 40)
	screen := plain(m)
	for _, want := range []string{
		"Dependencies: 10 unfinished work, 1 connected group, 2 unconnected · 2 done or abandoned prerequisites only in the trees",
		"Connected · 8 work", "Unconnected · 2 work",
		"│ W-02 proposed · layer 0 · connected with 7 other listed work", // the focus starts on the first row
		"indent: layer (equal: no declared order)",
	} {
		if !strings.Contains(screen, want) {
			t.Errorf("missing %q:\n%s", want, screen)
		}
	}
	// Indentation is the layer: W-05 sits two layers below W-02.
	row := func(id string) string {
		for _, r := range strings.Split(plain(m), "\n") {
			if i := strings.Index(r, " "+id+" "); i >= 0 && i < 50 {
				return r[:strings.Index(r, "│")]
			}
		}
		return ""
	}
	if a, b := strings.Index(row("W-02"), "W-02"), strings.Index(row("W-05"), "W-05"); b-a != 4 {
		t.Errorf("W-05 is not indented two layers under W-02 (%d, %d)", a, b)
	}
	if !strings.Contains(row("W-07"), "[not on main]") {
		t.Errorf("the review candidate off main is not marked: %q", row("W-07"))
	}

	focusOn(t, m, "W-05")
	screen = plain(m)
	for _, want := range []string{
		"← Needs", "├─ W-03 active", "│  └─ W-02 proposed", "│     └─ ✓ W-01 done", "└─ W-04 proposed", "W-02 (shown above)",
		"→ Unlocks", "├─ W-06 proposed", "└─ W-13 proposed",
	} {
		if !strings.Contains(screen, want) {
			t.Errorf("W-05's trees lack %q:\n%s", want, screen)
		}
	}
	focusOn(t, m, "W-08")
	if screen = plain(m); !strings.Contains(screen, "✗ W-10 abandoned") {
		t.Errorf("the abandoned prerequisite is not in W-08's tree:\n%s", screen)
	}
	focusOn(t, m, "W-11")
	if screen = plain(m); !strings.Contains(screen, "W-11 active ? Q-12") || !strings.Contains(screen, "Blocked by open questions, not by work") {
		t.Errorf("the blocking question is not kept apart from work edges:\n%s", screen)
	}
	if len(*asked) != 0 {
		t.Errorf("the list read Git: %q", *asked)
	}

	// h shows every work: the done and abandoned prerequisites become rows.
	press(m, "h")
	if screen = plain(m); !strings.Contains(screen, "every work, 12, 1 connected group") || !strings.Contains(screen, "W-10 abandoned") {
		t.Errorf("h did not show every work:\n%s", screen)
	}
	press(m, "h")

	// Enter opens the record, and Esc returns to the list, not the board.
	press(m, "enter")
	if m.screen != detailScreen || m.openID() != "W-11" || !strings.Contains(plain(m), "board › dependencies › W-11") {
		t.Fatalf("Enter did not open W-11 from the dependencies: %v %q", m.screen, m.stack)
	}
	press(m, "esc")
	if m.screen != depsScreen || m.depsAt != "W-11" {
		t.Fatalf("Esc from the record did not return to the list: %v %s", m.screen, m.depsAt)
	}
	press(m, "esc")
	if m.screen != boardScreen {
		t.Fatalf("Esc from the list did not return to the board: %v", m.screen)
	}
}

func TestDepsNarrowSwapsListAndTree(t *testing.T) {
	m, _, _ := openDeps(t, 80, 24)
	focusOn(t, m, "W-05")
	screen := plain(m)
	if strings.Contains(screen, "← Needs") || !strings.Contains(screen, "Connected · 8 work") || !strings.Contains(screen, "Tab tree") {
		t.Fatalf("80 columns should show the list alone:\n%s", screen)
	}
	press(m, "tab")
	if screen = plain(m); !strings.Contains(screen, "← Needs") || strings.Contains(screen, "Connected · 8 work") {
		t.Fatalf("Tab should show W-05's tree:\n%s", screen)
	}
	m.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
	if screen = plain(m); !strings.Contains(screen, "← Needs") || !strings.Contains(screen, "Connected · 8 work") {
		t.Fatalf("a wide terminal shows both:\n%s", screen)
	}
}

func TestDepsPreview(t *testing.T) {
	m, f, asked := openDeps(t, 120, 60)
	if cmd := press(m, "p"); cmd != nil || m.previewing || !strings.Contains(plain(m), "Space selects the work to preview") {
		t.Fatal("p with nothing selected must preview nothing")
	}
	// Marked out of order, W-05 then W-03 and W-07.
	focusOn(t, m, "W-05")
	press(m, "space")
	focusOn(t, m, "W-03")
	press(m, "space")
	focusOn(t, m, "W-07")
	press(m, "space", "space", "space") // unmarked and marked again: last
	if !slices.Equal(m.depsPicked, []string{"W-05", "W-03", "W-07"}) || !strings.Contains(plain(m), "●     W-05 proposed") || strings.Contains(plain(m), "●     W-04") {
		t.Fatalf("selection %q:\n%s", m.depsPicked, plain(m))
	}
	cmd := press(m, "p")
	if cmd == nil || !strings.Contains(plain(m), "Reading W-05 W-03 W-07") {
		t.Fatalf("p should start the delivery read:\n%s", plain(m))
	}
	deliver(m, cmd)
	screen := plain(m)
	for _, want := range []string{
		"Bound to checkout . (main) at aaaaaaa · target main",
		"Selected: W-05 W-03 W-07 (as marked)",
		"Order:    W-03 → W-05 → W-07",
		"1   W-03    active", "2   W-05    proposed", "3   W-07    review",
		"awaiting review; candidate c7 not in HEAD, not on main",
		"Outside the selection, not added",
		"W-01    done", "needed by W-03 W-05 W-07 · needs nothing · candidate c1 in HEAD, on main",
		"W-02    proposed", "W-04    proposed",
		"Questions: none open blocks",
		"No declared order between W-05 and W-07",
		"A preview adds no work, starts nothing, and authorizes nothing.",
	} {
		if !strings.Contains(screen, want) {
			t.Errorf("preview lacks %q:\n%s", want, screen)
		}
	}
	// The board and grove deps render one interpretation: same order, same
	// outside prerequisites.
	var records []*project.Record
	for _, g := range f.res.Groups {
		for _, v := range g.Versions {
			if v.Source.Kind == "live" {
				records = append(records, v.Record)
			}
		}
	}
	want, _ := deps.Preview(records, []string{"W-05", "W-03", "W-07"})
	if !slices.Equal(want.Order, m.preview.Order) || len(want.Items) != len(m.preview.Items) {
		t.Errorf("board %q, command %q", m.preview.Order, want.Order)
	}
	if !slices.Contains(*asked, "/repo/. c7 refs/heads/main") {
		t.Errorf("delivery was not read in the bound checkout: %q", *asked)
	}

	// A re-read recomputes the preview and says what changed.
	for i := range f.res.Groups {
		for j := range f.res.Groups[i].Versions {
			if v := &f.res.Groups[i].Versions[j]; v.Record.ID == "W-03" {
				v.Record.Status, v.Record.Candidate = "review", "c3"
				v.Record.Source = append(slices.Clone(v.Record.Source), "Moved to review.\n"...)
				v.Revision = project.Revision(v.Record.Source)
			}
		}
	}
	cmd = press(m, "r")
	cmd = deliver(m, cmd) // the inspection
	if m.preview != nil || cmd == nil {
		t.Fatal("a re-read must recompute the preview")
	}
	deliver(m, cmd)
	if screen = plain(m); !strings.Contains(screen, "W-03 changed since this preview was last read") || !strings.Contains(screen, "1   W-03    review") || strings.Contains(screen, "changed while it was being read") {
		t.Errorf("the re-read preview does not say W-03 changed:\n%s", screen)
	}
	press(m, "esc")
	if m.previewing || m.screen != depsScreen {
		t.Fatal("Esc should close the preview")
	}
}

func TestDepsPreviewBindsToOneCheckout(t *testing.T) {
	m, f, _ := openDeps(t, 120, 40)
	// W-20 is only on a feature branch: the list (current view) shows it,
	// but the preview reads this checkout and says so.
	feat := source("committed", "", "feature")
	f.res.Sources = append(f.res.Sources, feat)
	f.res.Groups = append(f.res.Groups, versions.Group{ID: "W-20", Versions: []versions.Version{version(feat, "W-20", "Only on feature", "proposed")}})
	deliver(m, press(m, "r"))
	focusOn(t, m, "W-20")
	press(m, "space")
	if cmd := press(m, "p"); cmd != nil {
		t.Fatal("no preview should be read")
	}
	if screen := plain(m); !strings.Contains(screen, "work W-20 is not in this checkout") || !strings.Contains(screen, "only; b chooses another checkout") {
		t.Errorf("the preview does not explain its binding:\n%s", screen)
	}
}
