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
)

func TestBoardInvocation(t *testing.T) {
	t.Parallel()
	// grove [--project DIR] [--json], in either order, and nothing else.
	for _, args := range [][]string{nil, {"--project", "d"}, {"--project=d"}, {"--json"}, {"--project", "d", "--json"}, {"--json", "--project", "d"}, {"--"}} {
		a, err := parseArgs(args)
		if err != nil || a.command != "" || a.help {
			t.Fatalf("%v: command=%q help=%v err=%v", args, a.command, a.help, err)
		}
	}
	for _, args := range [][]string{
		{"board"}, {"G-009"}, {"--json", "--json"}, {"--project", "a", "--project", "b"}, {"--wat"}, {"--project"}, {"--project="},
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
	t.Parallel()
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
	m := tui.New(t.Context(), root, tui.Live())
	m.Update(tea.WindowSizeMsg{Width: 160, Height: 40}) // five columns: 27-character titles need 32 cells each
	s := boardSession{t, m}
	s.run(m.Init())
	return s
}

func (s boardSession) run(cmd tea.Cmd) (quit bool) {
	for cmd != nil {
		msg := cmd()
		switch msg := msg.(type) {
		case tea.QuitMsg:
			return true
		case tea.BatchMsg: // the attempts read beside another (G-046)
			for _, c := range msg {
				quit = s.run(c) || quit
			}
			return quit
		}
		_, cmd = s.m.Update(msg)
	}
	return false
}

func (s boardSession) press(keys ...string) (quit bool) {
	named := map[string]tea.KeyPressMsg{"enter": {Code: tea.KeyEnter}, "down": {Code: tea.KeyDown}, "up": {Code: tea.KeyUp}, "esc": {Code: tea.KeyEscape}, "tab": {Code: tea.KeyTab}}
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

// The connected workflow on a real main/feature repository: the current view,
// a checkout's own columns through b, a card's differing versions, an explicit
// selection, and show reading exactly the selected bytes; then a target that
// changes before selection. Repeated with an unrelated invalid source. No file
// anywhere, including Git's, changes because of the board.
func TestBoardConnectedWorkflow(t *testing.T) {
	t.Parallel()
	for _, broken := range []bool{false, true} {
		t.Run(fmt.Sprintf("invalid source %v", broken), func(t *testing.T) {
			t.Parallel()
			boardWorkflow(t, broken)
		})
	}
}

func boardWorkflow(t *testing.T, broken bool) {
	{
		root, wt := featureFixture(t)
		record := filepath.Join(wt, "docs/records/work/renamed.md")
		onFeature := strings.NewReplacer("status: proposed", "status: active", "Inspect records", "Inspect records, on feature").Replace(work)
		write(t, wt, "docs/records/work/renamed.md", onFeature)
		write(t, wt, "docs/records/work/only.md", strings.NewReplacer("G-001", "G-003", "Inspect records", "Only on feature").Replace(work))
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

		// The current view: feature changed G-001 after main, and G-003 is
		// only there, so both show in feature's state.
		s := openBoard(t, root)
		s.want("Board: current view", "Proposed 1", "Active 1", "Inspect records, on feature", "Only on feature", "Deleted: none")
		s.lacks("G-002", "⑂", "uncommitted")
		incomplete(s)

		s.press("b")
		s.want("current view", "checkout . (main)", "checkout feature-wt (feature)")
		s.press("down", "enter")
		s.want("Board: checkout . (main)", "Proposed 1", "Active 0", "Inspect records", "2 versions", "Elsewhere (1", "): G-003 ")
		s.lacks("on feature", "Only on feature", "G-002")
		incomplete(s)

		s.press("b", "down")
		if broken { // sorted between main and feature-wt, and not choosable
			s.want("checkout broken-wt (broken)   UNAVAILABLE: invalid:")
			s.press("down", "enter")
			s.want("cannot fill the board", "Choose what the board shows")
		}
		s.press("down", "enter")
		s.want("Board: checkout feature-wt (feature)", "Proposed 1", "Active 1", "Inspect records, on feature", "Only on feature", "Elsewhere: none")
		incomplete(s)
		s.press("b", "down", "enter") // back to main's board
		s.want("Board: checkout . (main)", "Active 0")

		// Enter opens the record's detail: its content, rendered, its linked
		// records and the timeline from real Git; v is the version list.
		s.press("enter")
		s.want("G-001 · proposed", "Inspect records", "▶ Content", "An outcome.", "Linked", "related    G-002", "Timeline on checkout . (main)", "  proposed   ", "  init", "Sources", "current: active \"Inspect records, on feature\"", "places hold an earlier state")
		s.lacks("▸ proposed", "retitle and add")
		incomplete(s)
		s.press("v")
		s.want("G-001   2 versions differ", "▸ proposed   older  same on ", "▸ active     same on 1 branch, 1 checkout")
		s.lacks("checkout feature-wt (feature)  unchanged")
		incomplete(s)
		// The card opens on the lineage of the board's checkout, from real Git.
		s.want("History on checkout . (main):", "  proposed   ", "  init")
		s.lacks("retitle and add", "uncommitted")
		if s.press("enter") || s.m.Workspace != nil {
			t.Fatal("the ID header selected something")
		}
		s.focus("▸ active")
		s.want("Same on:  branch feature", "checkout feature-wt (feature)  unchanged", "Inspect records, on feature")
		s.want("History on branch feature:", "  retitle and add", "  feature", "  init")
		if s.press("enter") || s.m.Workspace != nil {
			t.Fatal("a fold of identical versions selected something")
		}
		s.focus("checkout feature-wt (feature)")
		s.want("Inspect records, on feature", "Status:   active", "Checkout: /", "Branch:   refs/heads/feature", "Selector: live:feature-wt:refs/heads/feature@")
		if !s.press("enter") || s.m.Workspace == nil {
			t.Fatalf("selecting the live feature version should resolve it:\n%s", s.screen())
		}
		ws := s.m.Workspace
		var out, errOut bytes.Buffer
		if code := Run([]string{"--project", ws.Project, "show", "G-001"}, root, &out, &errOut); code != 0 || out.String() != onFeature || ws.Project != wt || ws.Record != record {
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

		// The same walk from the current view, where G-001 is active, but the
		// target changes after it was displayed.
		s = openBoard(t, root)
		s.press("l", "enter", "v")
		s.focus("▸ active")
		s.press("enter")
		s.focus("checkout feature-wt (feature)")
		changed := onFeature + "Edited after the board read it.\n"
		write(t, wt, "docs/records/work/renamed.md", changed)
		before = all()
		if s.press("enter") || s.m.Workspace != nil {
			t.Fatal("a stale selection resolved")
		}
		s.want("REFUSED: G-001 changed since it was selected", "Nothing was opened. Press r to refresh")
		s.press("up", "down")
		s.want("REFUSED: G-001 changed since it was selected")
		s.press("r")
		s.lacks("REFUSED")
		s.want("active     checkout feature-wt (feature)  modified", "3 versions differ")
		if s.press("enter") || s.m.Workspace != nil {
			t.Fatal("refresh must not reselect the old row")
		}
		s.focus("checkout feature-wt (feature)")
		s.want("History on checkout feature-wt (feature):", "uncommitted       active     modified in this checkout's files", "  retitle and add")
		if !s.press("enter") || s.m.Workspace == nil {
			t.Fatalf("explicit reselection after refresh should resolve:\n%s", s.screen())
		}
		out.Reset()
		if code := Run([]string{"--project", s.m.Workspace.Project, "show", "G-001"}, root, &out, &errOut); code != 0 || out.String() != changed {
			t.Fatalf("show after reselection: %d\n%s", code, out.String())
		}
		// A committed version routes only while its checkout still matches it.
		s = openBoard(t, root)
		s.press("l", "enter", "v")
		s.focus("branch feature")
		if s.press("enter") || s.m.Workspace != nil {
			t.Fatal("a committed version whose checkout differs must be refused")
		}
		s.want("REFUSED: the live G-001 in worktree feature-wt differs from the committed version selected")
		if !reflect.DeepEqual(before, all()) {
			t.Fatal("refusal, refresh, and reselection changed files")
		}
	}
}

// The review workflow on real Git: G-001 is in review on feature with a
// candidate that adds a file. From the board in main's checkout the owner
// reads the candidate's standing, its changed files and a diff, approves it
// (written on feature), then integrates it into main, where it is done.
func TestBoardReviewWorkflow(t *testing.T) {
	t.Parallel()
	root := gitFixture(t)
	for _, kv := range [][2]string{{"user.name", "t"}, {"user.email", "t@t"}, {"commit.gpgsign", "false"}, {"maintenance.auto", "false"}} {
		gitIn(t, root, "config", kv[0], kv[1])
	}
	write(t, root, "grove.yaml", "schema_version: 3\nrecords: docs/records\ntarget: main\n")
	gitIn(t, root, "commit", "-qam", "chore: target")
	wt := filepath.Join(filepath.Dir(root), "feature-wt")
	gitIn(t, root, "worktree", "add", "-q", "-b", "feature", wt)
	write(t, wt, "code.txt", "hello <b>\x1b]0;evil\a\n")
	write(t, wt, "docs/records/work/renamed.md", strings.Replace(work, "status: proposed", "status: active", 1))
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-qm", "feat: the work")
	candidate := gitIn(t, wt, "rev-parse", "HEAD")
	var out, errOut bytes.Buffer
	if code := Run([]string{"--project", wt, "update", "G-001", "--set", "status=review", "--set", "candidate=" + candidate, "--commit"}, root, &out, &errOut); code != 0 {
		t.Fatal(errOut.String())
	}
	short := candidate[:7]

	s := openBoard(t, root)
	s.want("Board: current view, target main", "Review 1")
	s.press("l", "l", "enter")
	s.want("G-001 · review", "candidate "+short+" · not on main",
		"Review: candidate "+short+" · not yet approved · only the record changed since it · not on main",
		"a approve and f feedback run on branch feature in "+wt, "integrate runs into main in "+root,
		"Changes against main from ", "code.txt  +1 −0", "docs/records/work/renamed.md  +1 −1", "a approve   f feedback   i integrate")
	s.press("tab")
	for !strings.Contains(s.screen(), "> code.txt") {
		s.press("down")
	}
	s.press("enter")
	s.want("Diff of code.txt (Esc returns to the content)", `+hello <b>\x1b]0;evil\a`)
	if strings.Contains(s.m.View().Content, "\x1b]0;") {
		t.Fatal("the diff's control sequence reached the screen")
	}
	s.press("esc")
	s.lacks("Diff of")

	// Approval is written on feature, and the board re-read shows it.
	s.press("a")
	s.want("Approve G-001 on branch feature · verdict")
	for _, c := range "Ship it" {
		s.press(string(c))
	}
	s.press("enter")
	s.want("Approved G-001", "approved: G-001's candidate, in commit ", "The board has been re-read. Esc returns to the record.")
	if got := gitIn(t, wt, "show", "HEAD:docs/records/work/renamed.md"); !strings.Contains(got, "approved: \""+candidate+"\"") || !strings.Contains(got, "Verdict on candidate "+short+", ") || !strings.HasSuffix(got, ": Ship it") {
		t.Fatalf("feature's record after approval:\n%s", got)
	}
	s.press("esc")
	s.want("Review: candidate " + short + " · approved · only the record changed since it · not on main")

	// Integration merges feature into main and writes done there.
	s.press("i")
	s.want("Merge branch feature into main and mark G-001 done? y/n   (runs in ") // the temp path is truncated at 160 columns
	s.press("y")
	s.want("Also delete branch feature and remove its worktree? y/n   (")
	s.press("n")
	s.want("Integration of G-001", "approval: candidate "+short+" of G-001 approved on branch feature at ", "merge: fast-forward main from ", "done: G-001 done at commit ")
	if got := gitIn(t, root, "show", "HEAD:docs/records/work/renamed.md"); !strings.Contains(got, "status: done") || !strings.Contains(got, "approved: \""+candidate+"\"") {
		t.Fatalf("main's record after integration:\n%s", got)
	}
	if _, err := os.Stat(wt); err != nil {
		t.Fatal("n should keep the worktree")
	}
	s.press("esc")
	s.want("G-001 · done", "candidate "+short+" · on main")
	s.lacks("a approve")
}
