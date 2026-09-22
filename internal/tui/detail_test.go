package tui

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"

	"github.com/mascah/grove/internal/versions"
)

// linkedFixture: W-001 is work with a plan, a review, a prerequisite, a
// dependant, a blocking question, a parent, a member and a related page;
// G-020 is a page linked from nothing but W-001.
func linkedFixture(fx fixture) *fake {
	var vs []versions.Version
	for _, s := range []*versions.Source{fx.cMain, fx.main} {
		w := version(s, "W-001", "Inspect records", "active")
		p := 2
		w.Record.Kind, w.Record.Size, w.Record.Priority, w.Record.Candidate = "feature", "small", &p, "abcdef1"
		w.Record.DependsOn, w.Record.Members, w.Record.RelatesTo = []string{"W-002"}, []string{"W-004"}, []string{"G-020"}
		w.Record.Source = []byte("---\nid: W-001\n---\n\n## Outcome\n\nA **clear** board.\n\n- one\n- two\n" + strings.Repeat("\nfiller paragraph\n", 40) + "\nTHE END\n")
		u := time.Date(2026, 9, 22, 12, 0, 0, 0, time.UTC)
		w.Record.Updated, w.OnTarget = &u, true
		typed := func(id, title, kind, status string) versions.Version {
			v := version(s, id, title, status)
			v.Record.Type = kind
			return v
		}
		plan := typed("W-005", "W-001 plan", "plan", "current")
		plan.Record.Work = []string{"W-001"}
		review := typed("W-006", "W-001 review", "review", "current")
		review.Record.Work, review.Record.Examined = []string{"W-001"}, "abcdef1234"
		q := typed("Q-001", "Which first?", "question", "open")
		q.Record.Blocks = []string{"W-001"}
		parent := version(s, "W-007", "Milestone", "active")
		parent.Record.Members = []string{"W-001"}
		dependant := version(s, "W-003", "Later work", "proposed")
		dependant.Record.DependsOn = []string{"W-001"}
		page := typed("G-020", "How boards work", "page", "")
		page.Record.Source = []byte("---\nid: G-020\n---\n\n# Boards\n\nA page.\n")
		vs = append(vs, w, version(s, "W-002", "Create records", "done"), version(s, "W-004", "A member", "proposed"),
			plan, review, q, parent, dependant, page)
	}
	res := result(fx.main, fx.sources(), vs...)
	res.Target = "main"
	return &fake{res: res, history: func(ctx context.Context, commit, path string) ([]versions.Commit, error) {
		return []versions.Commit{
			{ID: "bbbbbbb" + strings.Repeat("1", 33), Subject: "docs: set active", Status: "active", When: time.Date(2026, 9, 22, 10, 0, 0, 0, time.UTC), Source: []byte("---\nid: W-001\n---\n\n# As it was\n\nOLDER BODY\n")},
			{ID: "aaaaaaa" + strings.Repeat("2", 33), Subject: "docs: propose", Status: "proposed", When: time.Date(2026, 9, 21, 9, 0, 0, 0, time.UTC), Source: []byte("---\nid: W-001\n---\n\nOLDEST BODY\n")},
			{ID: "0000000" + strings.Repeat("3", 33), Subject: "chore: remove", Status: "-", When: time.Date(2026, 9, 20, 9, 0, 0, 0, time.UTC)},
		}, nil
	}}
}

// Enter on a card opens the record's content first, with its linked records
// by role and its timeline beside it; the version list is a key away and
// Esc comes back.
func TestDetailLeadsWithContentAndLinks(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := linkedFixture(fx)
	m := open(t, f, 120, 36)
	press(m, "right") // W-001 in Active
	deliverAll(m, press(m, "enter"))
	s := plain(m)
	for _, want := range []string{"W-001 · active", "Inspect records", "feature · small · P2 · candidate abcdef1 · on main · same on 1 branch, 1 checkout · updated 2026-09-22",
		"▶ Content", "## Outcome", "A **clear** board.", "• one",
		"plan       W-005  W-001 plan  current", "review     W-006  W-001 review  current", "             examined abcdef1 = candidate", "needs      W-002  Create records  done",
		"needed by  W-003  Later work  proposed", "blocked by Q-001  Which first?  open", "part of    W-007  Milestone  active",
		"member     W-004  A member  proposed", "related    G-020  How boards work  -",
		"Timeline on branch main", "2026-09-22 10:00  active     bbbbbbb", "2026-09-20 09:00  -          0000000",
		"Sources", `current: active "Inspect records" on branch main,`, "v lists every version"} {
		if !strings.Contains(s, want) {
			t.Fatalf("detail lacks %q:\n%s", want, s)
		}
	}
	if i, j := strings.Index(s, "plan       W-005"), strings.Index(s, "related    G-020"); i > j {
		t.Fatalf("roles are ordered:\n%s", s)
	}
	if strings.Contains(s, "▸") || strings.Contains(s, "Details") {
		t.Fatalf("the version list is not the first view:\n%s", s)
	}
	// The content scrolls; the sidebar does not move with it.
	press(m, "pgdown", "pgdown", "pgdown")
	if s = plain(m); !strings.Contains(s, "THE END") || !strings.Contains(s, "plan       W-005") {
		t.Fatalf("content should scroll to its end:\n%s", s)
	}
	press(m, "pgup", "pgup", "pgup")
	// v is the version list, unchanged, and Esc returns to the detail.
	press(m, "v")
	if s = plain(m); m.screen != versionsScreen || !strings.Contains(s, "W-001   same everywhere") || !strings.Contains(s, "History on branch main") {
		t.Fatalf("v should show the versions:\n%s", s)
	}
	if press(m, "esc"); m.screen != detailScreen {
		t.Fatalf("Esc from the versions returns to the detail, not the board: screen %d", m.screen)
	}
	if press(m, "esc"); m.screen != boardScreen || m.cardID != "W-001" || len(m.stack) != 0 {
		t.Fatalf("Esc from the detail returns to the board on the same card: screen %d card %s", m.screen, m.cardID)
	}
	if f.inspects != 1 || len(f.resolved) != 0 {
		t.Fatal("browsing must not inspect again or resolve")
	}
}

// Tab reaches the linked records and the timeline; Enter on a linked record
// opens it, of any type, and Esc comes back; Enter on a commit shows the
// record as it was, and Esc returns to now.
func TestDetailNavigatesLinksAndTimeline(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	m := open(t, linkedFixture(fx), 120, 36)
	press(m, "right")
	deliverAll(m, press(m, "enter"))
	press(m, "tab") // the first linked record
	if s := plain(m); m.side != 0 || !strings.Contains(s, "> plan       W-005") || strings.Contains(s, "▶ Content") {
		t.Fatalf("Tab should focus the first linked record:\n%s", s)
	}
	for range 7 {
		press(m, "down")
	}
	if s := plain(m); !strings.Contains(s, "> related    G-020") {
		t.Fatalf("↓ walks the linked records:\n%s", s)
	}
	press(m, "enter")
	if s := plain(m); m.openID() != "G-020" || !strings.Contains(s, "G-020 · page") || !strings.Contains(s, "# Boards") || !strings.Contains(s, "related    W-001") {
		t.Fatalf("Enter opens the page's own detail with the link back:\n%s", s)
	}
	if press(m, "esc"); m.openID() != "W-001" || m.screen != detailScreen {
		t.Fatalf("Esc returns to the record that linked it: %s", m.openID())
	}
	press(m, "tab", "tab") // content, then the linked records, then the timeline
	if s := plain(m); !strings.Contains(s, "> 2026-09-22 10:00  active") {
		t.Fatalf("Tab should reach the timeline:\n%s", s)
	}
	press(m, "enter")
	if s := plain(m); m.asOf == "" || !strings.Contains(s, "Content as of bbbbbbb") || !strings.Contains(s, "OLDER BODY") || strings.Contains(s, "A **clear** board") {
		t.Fatalf("Enter on a commit shows the record as it was:\n%s", s)
	}
	press(m, "down", "down", "enter") // the commit that deleted the file has nothing to show
	if s := plain(m); m.asOf[:7] != "bbbbbbb" || !strings.Contains(s, "OLDER BODY") {
		t.Fatalf("a deletion commit changes nothing:\n%s", s)
	}
	if press(m, "esc"); m.asOf != "" || m.screen != detailScreen || !strings.Contains(plain(m), "A **clear** board") {
		t.Fatalf("Esc returns to now first: asOf %q screen %d", m.asOf, m.screen)
	}
	press(m, "tab")
	if m.side != -1 {
		t.Fatalf("Tab from the timeline returns to the content: side %d", m.side)
	}
}

// A narrow terminal shows one pane at a time and Tab cycles them; every
// row is the terminal's width at every size.
func TestDetailFitsNarrowTerminals(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	for _, size := range [][2]int{{120, 36}, {99, 24}, {80, 24}, {40, 10}} {
		m := open(t, linkedFixture(fx), size[0], size[1])
		press(m, "right")
		deliverAll(m, press(m, "enter"))
		check := func(where string) {
			t.Helper()
			rows := strings.Split(m.render(), "\n")
			if len(rows) != size[1] {
				t.Fatalf("%v %s: %d rows", size, where, len(rows))
			}
			for i, r := range rows {
				if got := ansi.StringWidth(r); got != size[0] {
					t.Fatalf("%v %s: row %d is %d cells: %q", size, where, i, got, ansi.Strip(r))
				}
			}
		}
		check("content")
		s := plain(m)
		if narrow := size[0] < wideWidth; narrow == strings.Contains(s, "plan       W-005") || !strings.Contains(s, "## Outcome") {
			t.Fatalf("%v: one pane at a time when narrow:\n%s", size, s)
		}
		press(m, "tab")
		check("linked")
		if s = plain(m); !strings.Contains(s, "> plan") || size[0] < wideWidth && strings.Contains(s, "## Outcome") {
			t.Fatalf("%v: Tab shows the sidebar:\n%s", size, s)
		}
		press(m, "tab")
		check("timeline")
		press(m, "tab")
		check("content again")
		if !strings.Contains(plain(m), "## Outcome") {
			t.Fatalf("%v: Tab cycles back to the content", size)
		}
	}
	// A lone record with no linked records and no commits still has its
	// Sources: a narrow terminal reaches them with Tab.
	fx = newFixture()
	m := open(t, &fake{res: result(fx.main, fx.sources(), version(fx.main, "W-001", "Alone", "active"))}, 80, 24)
	press(m, "right")
	deliverAll(m, press(m, "enter"))
	if s := plain(m); strings.Contains(s, "Sources") {
		t.Fatalf("narrow: the content first:\n%s", s)
	}
	press(m, "tab", "down", "enter")
	if s := plain(m); m.screen != detailScreen || !strings.Contains(s, "Linked") || !strings.Contains(s, "none") || !strings.Contains(s, "Sources") {
		t.Fatalf("narrow: Tab shows the sidebar with nothing to select:\n%s", s)
	}
	press(m, "tab")
	if s := plain(m); strings.Contains(s, "Sources") {
		t.Fatalf("narrow: Tab returns to the content:\n%s", s)
	}
}

// A refresh keeps the open detail while its record exists, drops the as-of
// view since the history is reread, and closes a vanished record with the
// reason; a current deletion opens with nothing to render.
func TestDetailSurvivesRefresh(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := linkedFixture(fx)
	m := open(t, f, 120, 36)
	press(m, "right")
	deliverAll(m, press(m, "enter"))
	press(m, "tab", "tab", "enter")
	if m.asOf == "" {
		t.Fatal("expected an as-of view")
	}
	deliverAll(m, press(m, "r"))
	if m.screen != detailScreen || m.openID() != "W-001" || m.asOf != "" || !strings.Contains(plain(m), "A **clear** board") {
		t.Fatalf("refresh should keep the detail on now: screen %d open %s asOf %q", m.screen, m.openID(), m.asOf)
	}
	// A record under the open one vanishes: it leaves the stack, and the
	// open one stays.
	for m.openID() != "W-005" { // Tab to the linked records; the plan is first
		press(m, "tab", "enter")
	}
	f.res = result(fx.main, fx.sources(), version(fx.main, "W-005", "W-001 plan", "current"), version(fx.main, "W-009", "Newcomer", "active"))
	deliverAll(m, press(m, "r"))
	if m.screen != detailScreen || m.openID() != "W-005" || len(m.stack) != 1 || !strings.Contains(plain(m), "W-001 is no longer on any readable branch or checkout") {
		t.Fatalf("a vanished record under the open one: screen %d stack %v\n%s", m.screen, m.stack, plain(m))
	}
	press(m, "esc")
	if m.screen != boardScreen {
		t.Fatal("Esc from the last open record returns to the board")
	}
	press(m, "/")
	typeText(m, "W-005")
	deliverAll(m, press(m, "enter"))
	f.res = result(fx.main, fx.sources(), version(fx.main, "W-009", "Newcomer", "active"))
	deliverAll(m, press(m, "r"))
	if m.screen != boardScreen || len(m.stack) != 0 || !strings.Contains(plain(m), "W-005 is no longer on any readable branch or checkout") {
		t.Fatalf("a vanished record closes its detail with the reason:\n%s", plain(m))
	}
	// A record whose current state is a deletion.
	gone := result(fx.main, fx.sources(), version(fx.cMain, "W-001", "Inspect records", "active"))
	gone.Groups[0].Versions[0].Older = "checkout . (main) deleted it since"
	gone.Groups[0].Versions = append(gone.Groups[0].Versions, versions.Version{Source: fx.main, Path: "grove/work/W-001.md", Change: "deleted"})
	f.res = gone
	deliverAll(m, press(m, "r"))
	press(m, "tab", "enter")
	if s := plain(m); m.screen != detailScreen || !strings.Contains(s, "W-001 · deleted") || !strings.Contains(s, "Deleted here") {
		t.Fatalf("a deleted current state opens without content:\n%s", s)
	}
}
