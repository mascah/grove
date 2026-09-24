package tui

import (
	"bytes"
	"cmp"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/mascah/grove/internal/attempt"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/update"
)

// Answering a question (G-125): e on an open question's detail, or on an
// attempt waiting on one, suspends the board for the owner's editor on the
// question's file in the checkout holding the shown version; on return it
// offers to resolve the question and commit it there. The board writes only
// an Answer heading, which it takes back if the editor leaves it unused, and
// the one update behind the prompt.

// editing is an edit the editor has: where, and the file's bytes as read and
// as handed over.
type editing struct {
	id, root, branch, path string
	before, given          []byte
	earlier                bool // the file already held an uncommitted edit
}

// editedMsg is the editor's exit.
type editedMsg struct{ err error }

var answerHeading = regexp.MustCompile(`(?m)^##[ \t]+Answer[ \t]*\r?$`)

// answer opens the open question the detail shows in the owner's editor, or
// says why it cannot. The file must be what the board read.
func (m *Model) answer() tea.Cmd {
	g := m.group()
	if m.backend.Edit == nil || g == nil {
		return nil
	}
	v := m.shown(g)
	switch {
	case v == nil || v.Record == nil || v.Record.Type != "question":
		m.notice = "e answers a question; " + g.ID + " is not one"
		return nil
	case v.Record.Status != "open":
		m.notice = g.ID + " is " + v.Record.Status + "; e answers an open question"
		return nil
	}
	lv, branch, why := m.checkoutOf(g, v, "answering")
	if why != "" {
		m.notice = why + "; nothing was written"
		return nil
	}
	root := m.projectDir(lv.Source)
	path := filepath.Join(root, filepath.FromSlash(lv.Path))
	before, err := os.ReadFile(path)
	switch {
	case err != nil:
		m.notice = err.Error() + "; nothing was written"
		return nil
	case project.Revision(before) != lv.Revision:
		m.notice = g.ID + " changed in " + root + " since the board read it; r re-reads it, and nothing was written"
		return nil
	}
	given := before
	if !answerHeading.Match(before) {
		// ponytail: written without the write lock, like the editor's own
		// save; only an update of this file in this instant could interleave.
		if given = withHeading(before); os.WriteFile(path, given, 0o644) != nil {
			given = before // the editor still opens; the answer finds its own place
		}
	}
	m.editing = &editing{g.ID, root, branch, path, before, given, lv.Change != "unchanged"}
	return m.backend.Edit(path, func(err error) tea.Msg { return editedMsg{err} })
}

// withHeading appends an Answer heading in the file's own line ending.
func withHeading(b []byte) []byte {
	nl := "\n"
	if bytes.Contains(b, []byte("\r\n")) {
		nl = "\r\n"
	}
	out := bytes.Clone(b)
	if len(out) != 0 && out[len(out)-1] != '\n' {
		out = append(out, nl...)
	}
	return append(out, nl+"## Answer"+nl+nl...)
}

// edited takes the editor's exit: a failure or no change says so, and an
// edit, or an earlier one still uncommitted, opens the prompt to resolve and
// commit it. The heading the board added is taken back when nothing was
// written beneath it.
func (m *Model) edited(msg editedMsg) tea.Cmd {
	e := m.editing
	m.editing = nil
	if e == nil {
		return nil
	}
	after, err := os.ReadFile(e.path)
	unused := err == nil && bytes.Equal(after, e.given)
	if unused {
		if werr := e.takeBack(); werr != nil {
			m.notice = "the Answer heading could not be taken back (" + werr.Error() + "); nothing was committed, and it stays in " + e.path
			return m.refresh()
		}
		after = e.before // what an earlier edit's resolve expects
	}
	kept := "; nothing was committed, and the edit stays uncommitted in " + e.root
	switch {
	case msg.err != nil && unused:
		m.notice = "the editor failed (" + msg.err.Error() + "); nothing was written"
		return nil
	case msg.err != nil:
		m.notice = "the editor failed (" + msg.err.Error() + ")" + kept
	case err != nil:
		m.notice = err.Error() + "; nothing was committed"
	case unused && !e.earlier:
		m.notice = "no change was saved to " + e.id + "; nothing was written"
		return nil
	default:
		m.prompt = &prompt{kind: "resolve", id: e.id, root: e.root, branch: e.branch, expect: project.Revision(after)}
	}
	// The detail shows the file as it is now.
	return m.refresh()
}

// takeBack restores the file as read when it is still exactly as handed to
// the editor: the Answer heading the board added is removed.
func (e *editing) takeBack() error {
	if bytes.Equal(e.given, e.before) {
		return nil
	}
	if now, err := os.ReadFile(e.path); err != nil || !bytes.Equal(now, e.given) {
		return err
	}
	return os.WriteFile(e.path, e.before, 0o644)
}

// waits reports an attempt waiting on a question that e can answer.
func (m *Model) waits(v *attempt.View) bool {
	kind, _ := m.outcome(v)
	return m.backend.Edit != nil && kind == "question"
}

// answerFor opens the question an attempt waits on as o opens work, so Esc
// returns to the attempt, and starts its edit.
func (m *Model) answerFor(v *attempt.View) tea.Cmd {
	if m.backend.Edit == nil {
		return nil
	}
	if !m.waits(v) {
		m.notice = v.Launch.Attempt + " is not waiting on a question; e answers one"
		return nil
	}
	q, _, _ := strings.Cut(m.blockingQuestion(v.Launch.Work), " (")
	if m.openWork(q); m.openID() != q {
		return nil // openWork said why
	}
	return m.answer()
}

// answerRow is an open question's header row: where e writes, or why not.
func (m *Model) answerRow() string {
	g := m.group()
	lv, branch, why := m.checkoutOf(g, m.shown(g), "answering")
	if lv == nil {
		return "e answer: " + why
	}
	return "e answers it in your editor, then resolves and commits it on branch " + branch + " in " + m.projectDir(lv.Source)
}

// liveAnswer is the answering part of the backend over the real repository:
// the owner's editor on the terminal itself, since the program's output is a
// wrapper the editor would see as a pipe, and update as the CLI runs it.
func liveAnswer(b *Backend, input, screen *os.File) {
	b.Edit = func(path string, done func(error) tea.Msg) tea.Cmd {
		// As Git runs it: the variable may hold arguments and quotes.
		editor := cmp.Or(strings.TrimSpace(os.Getenv("VISUAL")), strings.TrimSpace(os.Getenv("EDITOR")), "vi")
		c := exec.Command("sh", "-c", editor+` "$@"`, editor, path)
		c.Stdin, c.Stdout, c.Stderr = input, screen, screen
		return tea.ExecProcess(c, done)
	}
	b.Answer = func(_ context.Context, root, id, expect string) ([]string, error) {
		res, err := update.Apply(root, update.Request{ID: id, Expect: expect, Set: []update.Field{{Name: "status", Value: "resolved"}}, Commit: true}, time.Now(), nil)
		switch {
		case err != nil:
			return nil, err
		case !res.Changed:
			return nil, fmt.Errorf("%s already says resolved in %s, so nothing was committed; commit %s there yourself", id, root, res.Path)
		}
		branch, _ := update.Branch(root)
		return []string{fmt.Sprintf("resolved: %s, with the answer, in commit %s on branch %s in %s", id, short7(res.Commit), branch, root)}, nil
	}
}
