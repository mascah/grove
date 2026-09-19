package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mascah/grove/internal/tui"
	"github.com/mascah/grove/internal/versions"
)

func TestBoardInvocation(t *testing.T) {
	// grove [--project DIR] [--json], in either order, and nothing else.
	for _, args := range [][]string{nil, {"--project", "d"}, {"--project=d"}, {"--json"}, {"--project", "d", "--json"}, {"--json", "--project", "d"}, {"--"}} {
		a, err := parseArgs(args)
		if err != nil || a.command != "" || a.help {
			t.Fatalf("%v: command=%q help=%v err=%v", args, a.command, a.help, err)
		}
	}
	for _, args := range [][]string{
		{"board"}, {"W-009"}, {"--json", "--json"}, {"--project", "a", "--project", "b"}, {"--wat"}, {"--project"}, {"--project="},
		{"--source", selector}, {"--slug", "x"}, {"--expect", rev}, {"--set", "status=done"}, {"--unset", "kind"}, {"--json", "board"},
	} {
		var out, errOut bytes.Buffer
		if code := Run(args, t.TempDir(), &out, &errOut); code != 2 || out.Len() != 0 || !strings.Contains(errOut.String(), "Usage:") {
			t.Fatalf("%v: code=%d stderr=%s", args, code, errOut.String())
		}
	}
	// Help never needs a project or a terminal and never opens the board.
	for _, args := range [][]string{{"--help"}, {"-h"}, {"help"}, {"--json", "--help"}} {
		var out, errOut bytes.Buffer
		if code := Run(args, t.TempDir(), &out, &errOut); code != 0 || !strings.Contains(out.String(), "Open the terminal board") || errOut.Len() != 0 {
			t.Fatalf("%v: code=%d stdout=%s stderr=%s", args, code, out.String(), errOut.String())
		}
	}
}

// Without a terminal the board refuses at once: exit 1, nothing on stdout, no
// other format, and the noninteractive commands named. It needs no project to
// say so. A file that is not a terminal is refused like a buffer.
func TestBoardRefusesWithoutTerminal(t *testing.T) {
	file, err := os.Create(filepath.Join(t.TempDir(), "stderr"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	root := projectFixture(t)
	for _, args := range [][]string{nil, {"--json"}, {"--project", root}} {
		var out, errOut bytes.Buffer
		if code := Run(args, t.TempDir(), &out, &errOut); code != 1 || out.Len() != 0 {
			t.Fatalf("%v: code=%d stdout=%q", args, code, out.String())
		}
		for _, want := range []string{"needs a terminal on stdin and stderr", "grove list", "grove versions", "grove workspace", "grove --help"} {
			if !strings.Contains(errOut.String(), want) {
				t.Fatalf("%v: refusal lacks %q: %s", args, want, errOut.String())
			}
		}
		if code := Run(args, root, &out, file); code != 1 || out.Len() != 0 {
			t.Fatalf("%v with a plain file: code=%d stdout=%q", args, code, out.String())
		}
	}
}

// boardSession drives the real model against real Git, running each command
// by hand as the terminal runtime would.
type boardSession struct {
	t *testing.T
	m *tui.Model
}

func openBoard(t *testing.T, root string) boardSession {
	t.Helper()
	m := tui.New(t.Context(), root, tui.Backend{Inspect: versions.InspectContext, Resolve: versions.ResolveContext})
	m.Update(tea.WindowSizeMsg{Width: 140, Height: 40})
	s := boardSession{t, m}
	s.run(m.Init())
	return s
}

func (s boardSession) run(cmd tea.Cmd) (quit bool) {
	for cmd != nil {
		msg := cmd()
		if _, ok := msg.(tea.QuitMsg); ok {
			return true
		}
		_, cmd = s.m.Update(msg)
	}
	return false
}

func (s boardSession) press(keys ...string) (quit bool) {
	named := map[string]tea.KeyPressMsg{"enter": {Code: tea.KeyEnter}, "down": {Code: tea.KeyDown}, "up": {Code: tea.KeyUp}, "esc": {Code: tea.KeyEscape}}
	for _, k := range keys {
		msg, ok := named[k]
		if !ok {
			msg = tea.KeyPressMsg{Code: rune(k[0]), Text: k}
		}
		_, cmd := s.m.Update(msg)
		quit = s.run(cmd)
	}
	return quit
}

// focus moves down to the row showing label, as a person reading it would.
func (s boardSession) focus(label string) {
	s.t.Helper()
	for range 12 {
		for _, row := range strings.Split(s.screen(), "\n") {
			if strings.HasPrefix(row, "> ") && strings.Contains(row, label) {
				return
			}
		}
		s.press("down")
	}
	s.t.Fatalf("no row %q to focus:\n%s", label, s.screen())
}

func (s boardSession) screen() string { return ansi.Strip(s.m.View().Content) }

func (s boardSession) want(parts ...string) {
	s.t.Helper()
	for _, p := range parts {
		if !strings.Contains(s.screen(), p) {
			s.t.Fatalf("screen lacks %q:\n%s", p, s.screen())
		}
	}
}

func (s boardSession) lacks(parts ...string) {
	s.t.Helper()
	for _, p := range parts {
		if strings.Contains(s.screen(), p) {
			s.t.Fatalf("screen shows %q:\n%s", p, s.screen())
		}
	}
}

// The connected workflow on a real main/feature repository: scoped columns,
// another checkout through b, a card's differing versions, an explicit
// selection, and show reading exactly the selected bytes; then a target that
// changes before selection. Repeated with an unrelated invalid source. No file
// anywhere, including Git's, changes because of the board.
func TestBoardConnectedWorkflow(t *testing.T) {
	for _, broken := range []bool{false, true} {
		t.Run(fmt.Sprintf("invalid source %v", broken), func(t *testing.T) { boardWorkflow(t, broken) })
	}
}

func boardWorkflow(t *testing.T, broken bool) {
	{
		root, wt := featureFixture(t)
		record := filepath.Join(wt, "docs/records/work/renamed.md")
		onFeature := strings.NewReplacer("status: proposed", "status: active", "Inspect records", "Inspect records, on feature").Replace(work)
		write(t, wt, "docs/records/work/renamed.md", onFeature)
		write(t, wt, "docs/records/work/only.md", strings.NewReplacer("W-001", "W-002", "Inspect records", "Only on feature").Replace(work))
		gitIn(t, wt, "add", "-A")
		gitIn(t, wt, "commit", "-q", "-m", "retitle and add")
		if broken {
			gitIn(t, root, "branch", "broken")
			gitIn(t, root, "worktree", "add", "-q", filepath.Join(filepath.Dir(root), "broken-wt"), "broken")
			write(t, filepath.Dir(root), "broken-wt/docs/records/work/renamed.md", strings.Replace(work, "status: proposed", "status: nonsense", 1))
		}
		incomplete := func(s boardSession) {
			if t.Helper(); broken {
				s.want("INCOMPLETE: 1 of")
			} else {
				s.lacks("INCOMPLETE")
			}
		}
		all := func() []map[string][32]byte {
			return []map[string][32]byte{hashes(t, root), hashes(t, wt)}
		}
		before := all()

		s := openBoard(t, root)
		s.want("Board: live . main", "Proposed (1)", "Active (0)", "Inspect records", "Other sources (1", "W-002 [2 versions]")
		s.lacks("on feature", "Only on feature", "Q-001")
		incomplete(s)

		s.press("b")
		s.want("live . main", "live feature-wt feature")
		if broken { // sorted between main and feature-wt, and not choosable
			s.want("live broken-wt broken   UNAVAILABLE: invalid:")
			s.press("down", "enter")
			s.want("cannot be a board context", "Choose the checkout")
		}
		s.press("down", "enter")
		s.want("Board: live feature-wt feature", "Proposed (1)", "Active (1)", "Inspect records, on feature", "Only on feature", "Other sources: none")
		incomplete(s)
		s.press("b", "enter") // back to main's board
		s.want("Board: live . main", "Active (0)")

		s.press("enter")
		s.want("versions: select one explicitly", "proposed   committed main", "proposed   live . main  unchanged", "active     committed feature", "active     live feature-wt feature  unchanged")
		incomplete(s)
		if s.press("enter") || s.m.Workspace != nil {
			t.Fatal("the ID header selected something")
		}
		s.focus("live feature-wt feature")
		s.want("Inspect records, on feature", "Status:   active", "Checkout: /", "Branch:   refs/heads/feature", "Selector: live:feature-wt:refs/heads/feature@")
		if !s.press("enter") || s.m.Workspace == nil {
			t.Fatalf("selecting the live feature version should resolve it:\n%s", s.screen())
		}
		ws := s.m.Workspace
		var out, errOut bytes.Buffer
		if code := Run([]string{"--project", ws.Project, "show", "W-001"}, root, &out, &errOut); code != 0 || out.String() != onFeature || ws.Project != wt || ws.Record != record {
			t.Fatalf("show through the selected workspace: code=%d project=%s\n%s", code, ws.Project, out.String())
		}
		// The board's result goes through workspace's own writer.
		out.Reset()
		errOut.Reset()
		if code := writeWorkspace(ws, false, &out, &errOut); code != 0 || out.String() != wt+"\n" || !strings.Contains(errOut.String(), "Checkout: "+wt) {
			t.Fatalf("result: %q %q", out.String(), errOut.String())
		}
		if !reflect.DeepEqual(before, all()) {
			t.Fatal("browsing and selecting changed files")
		}

		// The same walk, but the target changes after it was displayed.
		s = openBoard(t, root)
		s.press("enter")
		s.focus("live feature-wt feature")
		changed := onFeature + "Edited after the board read it.\n"
		write(t, wt, "docs/records/work/renamed.md", changed)
		before = all()
		if s.press("enter") || s.m.Workspace != nil {
			t.Fatal("a stale selection resolved")
		}
		s.want("REFUSED: W-001 changed since it was selected", "Nothing was opened. Press r to refresh")
		s.press("up", "down")
		s.want("REFUSED: W-001 changed since it was selected")
		s.press("r")
		s.lacks("REFUSED")
		s.want("active     live feature-wt feature  modified")
		if s.press("enter") || s.m.Workspace != nil {
			t.Fatal("refresh must not reselect the old row")
		}
		if s.focus("live feature-wt feature"); !s.press("enter") || s.m.Workspace == nil {
			t.Fatalf("explicit reselection after refresh should resolve:\n%s", s.screen())
		}
		out.Reset()
		if code := Run([]string{"--project", s.m.Workspace.Project, "show", "W-001"}, root, &out, &errOut); code != 0 || out.String() != changed {
			t.Fatalf("show after reselection: %d\n%s", code, out.String())
		}
		// A committed version routes only while its checkout still matches it.
		s = openBoard(t, root)
		s.press("enter")
		s.focus("committed feature")
		if s.press("enter") || s.m.Workspace != nil {
			t.Fatal("a committed version whose checkout differs must be refused")
		}
		s.want("REFUSED: the live W-001 in worktree feature-wt differs from the committed version selected")
		if !reflect.DeepEqual(before, all()) {
			t.Fatal("refusal, refresh, and reselection changed files")
		}
	}
}
