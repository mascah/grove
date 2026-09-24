package tui

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/mascah/grove/internal/attempt"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/versions"
)

// editor fakes the owner's editor and the resolve: the editor runs as a
// command whose exit is written, and Answer logs what it was asked.
type editor struct {
	mu      sync.Mutex
	write   string // appended to the file by the editor; "" saves nothing
	lock    bool   // the editor leaves the file read-only
	fail    error  // the editor's exit
	answer  error  // Answer's failure
	edits   []string
	answers []string
}

func (e *editor) add(b Backend) Backend {
	b.Edit = func(path string, done func(error) tea.Msg) tea.Cmd {
		e.mu.Lock()
		e.edits = append(e.edits, path)
		e.mu.Unlock()
		return func() tea.Msg {
			if e.write != "" {
				f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
				if err != nil {
					return done(err)
				}
				f.WriteString(e.write)
				f.Close()
			}
			if e.lock {
				os.Chmod(path, 0o444)
			}
			return done(e.fail)
		}
	}
	b.Answer = func(_ context.Context, root, id, expect string) ([]string, error) {
		e.mu.Lock()
		defer e.mu.Unlock()
		e.answers = append(e.answers, root+" "+id+" "+expect)
		if e.answer != nil {
			return nil, e.answer
		}
		return []string{"resolved: " + id + " in " + root}, nil
	}
	return b
}

const question = "---\nid: Q-001\ntype: question\nstatus: open\n---\n\n## Question\n\nRed or blue?\n"

// answerFixture: Q-001, open, blocks W-001 and is on branch feature, whose
// checkout feat is a real directory holding the file as the board read it.
func answerFixture(t *testing.T) (fx fixture, f *fake, path string) {
	fx = newFixture()
	fx.feat.Worktree = t.TempDir()
	path = filepath.Join(fx.feat.Worktree, "grove", "question", "Q-001.md")
	os.MkdirAll(filepath.Dir(path), 0o755)
	if err := os.WriteFile(path, []byte(question), 0o644); err != nil {
		t.Fatal(err)
	}
	var vs []versions.Version
	for _, s := range []*versions.Source{fx.cFeat, fx.feat} {
		q := version(s, "Q-001", "Red or blue?", "open")
		q.Record.Blocks, q.Record.Source, q.Revision = []string{"W-001"}, []byte(question), project.Revision([]byte(question))
		vs = append(vs, q, version(s, "W-001", "Inspect records", "active"))
	}
	return fx, &fake{res: result(fx.main, fx.sources(), vs...)}, path
}

func openQuestion(t *testing.T, f *fake, e *editor) *Model {
	t.Helper()
	m := open(t, &fake{}, 120, 36)
	m.backend = e.add(f.backend())
	m.backend.Launch = func(context.Context, attempt.Request) ([]string, error) { return nil, nil }
	deliver(m, press(m, "r"))
	m.openDetail("Q-001")
	return m
}

func fileIs(t *testing.T, path, want string) {
	t.Helper()
	if got, _ := os.ReadFile(path); string(got) != want {
		t.Fatalf("the file holds:\n%q\nwant\n%q", got, want)
	}
}

// e opens the question in the editor under an Answer heading; the edit is
// offered for resolving, and y resolves it at the revision the editor left,
// in the branch's checkout, then names the work to launch.
func TestAnswerResolvesInTheBranchCheckout(t *testing.T) {
	t.Parallel()
	fx, f, path := answerFixture(t)
	e := &editor{write: "Blue.\n"}
	m := openQuestion(t, f, e)
	if s := plain(m); !strings.Contains(s, "e answers it in your editor, then resolves and commits it on branch feature in") || !strings.Contains(s, fx.feat.Worktree) || !strings.Contains(s, "e answer") {
		t.Fatalf("the question's detail should say where e writes:\n%s", s)
	}
	cmd := press(m, "e")
	answered := question + "\n## Answer\n\nBlue.\n"
	if len(e.edits) != 1 || e.edits[0] != path {
		t.Fatalf("edits: %v", e.edits)
	}
	fileIs(t, path, question+"\n## Answer\n\n")
	deliverAll(m, cmd)
	fileIs(t, path, answered)
	if m.prompt == nil || m.prompt.kind != "resolve" || !strings.Contains(plain(m), "Resolve Q-001 and commit it with your answer on branch feature? y/n") {
		t.Fatalf("the resolve prompt should be open:\n%s", plain(m))
	}
	deliverAll(m, press(m, "y"))
	if want := fx.feat.Worktree + " Q-001 " + project.Revision([]byte(answered)); len(e.answers) != 1 || e.answers[0] != want {
		t.Fatalf("answers: %v, want %s", e.answers, want)
	}
	s := plain(m)
	for _, want := range []string{"Answer to Q-001", "resolved: Q-001 in", "next: R on W-001 launches its next attempt"} {
		if !strings.Contains(s, want) {
			t.Fatalf("the result lacks %q:\n%s", want, s)
		}
	}
	if press(m, "esc"); m.screen != detailScreen || m.openID() != "Q-001" {
		t.Fatal("Esc returns to the question")
	}
}

// n keeps the edit uncommitted and says where; an unsaved editor takes back
// the heading and writes nothing; a failed editor says so.
func TestAnswerDeclinedUnsavedOrFailed(t *testing.T) {
	t.Parallel()
	_, f, path := answerFixture(t)
	e := &editor{write: "Blue.\n"}
	m := openQuestion(t, f, e)
	deliverAll(m, press(m, "e"))
	if press(m, "n"); m.prompt != nil || !strings.Contains(plain(m), "Q-001 is not resolved; your edit stays uncommitted in ") || len(e.answers) != 0 {
		t.Fatalf("n should decline:\n%s", plain(m))
	}
	edited := question + "\n## Answer\n\nBlue.\n"
	fileIs(t, path, edited)
	// Read again, the checkout holds the uncommitted answer: e reopens it,
	// and with nothing more saved still offers the resolve.
	for i := range f.res.Groups[0].Versions {
		if v := &f.res.Groups[0].Versions[i]; v.Source.Kind == "live" {
			v.Change, v.Revision, v.Record.Source = "modified", project.Revision([]byte(edited)), []byte(edited)
		}
	}
	deliverAll(m, press(m, "r"))
	e.write = ""
	deliverAll(m, press(m, "e"))
	if m.prompt == nil || m.prompt.kind != "resolve" || m.prompt.expect != project.Revision([]byte(edited)) {
		t.Fatalf("e on the uncommitted answer should offer the resolve again:\n%s", plain(m))
	}
	fileIs(t, path, edited)

	// An earlier answer without the heading: the heading is taken back, and
	// the resolve expects the file as it is again.
	_, f, path = answerFixture(t)
	byHand := question + "Blue, by hand.\n"
	os.WriteFile(path, []byte(byHand), 0o644)
	for i := range f.res.Groups[0].Versions {
		if v := &f.res.Groups[0].Versions[i]; v.Source.Kind == "live" {
			v.Change, v.Revision, v.Record.Source = "modified", project.Revision([]byte(byHand)), []byte(byHand)
		}
	}
	e = &editor{}
	m = openQuestion(t, f, e)
	deliverAll(m, press(m, "e"))
	if m.prompt == nil || m.prompt.expect != project.Revision([]byte(byHand)) {
		t.Fatalf("the resolve should expect the answer as written by hand:\n%s", plain(m))
	}
	fileIs(t, path, byHand)

	for _, c := range []struct {
		fail   error
		notice string
	}{
		{nil, "no change was saved to Q-001; nothing was written"},
		{errors.New("exit status 1"), "the editor failed (exit status 1); nothing was written"},
	} {
		_, f, path := answerFixture(t)
		e := &editor{fail: c.fail}
		m := openQuestion(t, f, e)
		deliverAll(m, press(m, "e"))
		if m.prompt != nil || !strings.Contains(plain(m), c.notice) {
			t.Fatalf("want %q:\n%s", c.notice, plain(m))
		}
		fileIs(t, path, question)
	}

	// A failed editor that saved leaves its edit, and no prompt.
	_, f, path = answerFixture(t)
	e = &editor{write: "Half.\n", fail: errors.New("signal: killed")}
	m = openQuestion(t, f, e)
	deliverAll(m, press(m, "e"))
	if m.prompt != nil || !strings.Contains(plain(m), "the editor failed (signal: killed); nothing was committed, and the edit stays uncommitted in ") {
		t.Fatalf("a failed editor that saved:\n%s", plain(m))
	}
	fileIs(t, path, question+"\n## Answer\n\nHalf.\n")

	// A heading that cannot be taken back is reported, and nothing offered.
	_, f, path = answerFixture(t)
	e = &editor{lock: true}
	m = openQuestion(t, f, e)
	deliverAll(m, press(m, "e"))
	if m.prompt != nil || !strings.Contains(plain(m), "the Answer heading could not be taken back") {
		t.Fatalf("a failed take-back:\n%s", plain(m))
	}
	fileIs(t, path, question+"\n## Answer\n\n")
}

// Refusals write nothing and open no editor: the file changed since the
// board read it, no checkout or two are on the branch, the question is
// resolved, or the record is not a question. A resolve refused at its
// revision shows as not done.
func TestAnswerRefusals(t *testing.T) {
	t.Parallel()
	refuse := func(name string, change func(fx fixture, f *fake, path string), id, want string) {
		t.Helper()
		fx, f, path := answerFixture(t)
		change(fx, f, path)
		e := &editor{write: "Blue.\n"}
		m := openQuestion(t, f, e)
		m.openDetail(id)
		press(m, "e")
		if len(e.edits) != 0 || !strings.Contains(plain(m), want) {
			t.Fatalf("%s: want %q, edits %v:\n%s", name, want, e.edits, plain(m))
		}
		if name != "stale" {
			fileIs(t, path, question)
		}
	}
	refuse("stale", func(_ fixture, _ *fake, path string) { os.WriteFile(path, []byte(question+"Later.\n"), 0o644) }, "Q-001",
		"Q-001 changed in ")
	refuse("no checkout", func(_ fixture, f *fake, _ string) {
		f.res.Sources = f.res.Sources[:3]
	}, "Q-001", "no checkout is on branch feature; git worktree add one; nothing was written")
	refuse("ambiguous", func(fx fixture, f *fake, _ string) {
		other := source("live", "feat2", "feature")
		f.res.Sources = append(f.res.Sources, other)
		v := f.res.Groups[0].Versions[1]
		v.Source = other
		f.res.Groups[0].Versions = append(f.res.Groups[0].Versions, v)
	}, "Q-001", "2 checkouts are on branch feature, so which one to write is ambiguous")
	refuse("resolved", func(_ fixture, f *fake, _ string) {
		for i := range f.res.Groups[0].Versions {
			f.res.Groups[0].Versions[i].Record.Status = "resolved"
		}
	}, "Q-001", "Q-001 is resolved; e answers an open question")
	refuse("not a question", func(fixture, *fake, string) {}, "W-001", "e answers a question; W-001 is not one")

	_, f, _ := answerFixture(t)
	e := &editor{write: "Blue.\n", answer: errors.New("Q-001 changed since the expected revision")}
	m := openQuestion(t, f, e)
	deliverAll(m, press(m, "e"))
	deliverAll(m, press(m, "y"))
	if s := plain(m); !strings.Contains(s, "NOT DONE: Q-001 changed since the expected revision") || strings.Contains(s, "next: R") {
		t.Fatalf("a refused resolve:\n%s", s)
	}
}

// From an attempt waiting on the question, e opens it in the editor above
// the attempt; once it is resolved after the attempt ended, the attempt
// asks for the next launch.
func TestAnswerFromTheAttempt(t *testing.T) {
	t.Parallel()
	_, f, path := answerFixture(t)
	ended := time.Date(2026, 9, 23, 2, 0, 0, 0, time.UTC)
	asked := ended.Add(-time.Minute)
	for i := range f.res.Groups[0].Versions {
		f.res.Groups[0].Versions[i].Record.Created, f.res.Groups[0].Versions[i].Record.Updated = &asked, &asked
	}
	r := &runs{}
	waited := view("W-001", "1", attempt.Finished, &attempt.Result{Finished: ended, Events: attempt.Events{Result: &attempt.Final{Subtype: "success"}}, Record: &attempt.State{Status: "active"}})
	r.set(waited)
	e := &editor{write: "Blue.\n"}
	m := New(t.Context(), "/repo/.", e.add(r.add(f.backend())))
	m.clock = func() time.Time { return ended.Add(time.Hour) }
	cmd := m.Init()
	m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	settle(m, cmd)
	m.openAttempt(waited.Launch.Attempt)
	settle(m, nil)
	m.run = &waited
	if s := plain(m); !strings.Contains(s, "answer question Q-001") || !strings.Contains(s, "e answers Q-001 in your editor and offers to resolve it") {
		t.Fatalf("the attempt should offer e:\n%s", s)
	}
	cmd = press(m, "e")
	if m.screen != detailScreen || m.openID() != "Q-001" || len(e.edits) != 1 || e.edits[0] != path {
		t.Fatalf("e should open Q-001 in the editor: screen %d, %s, %v", m.screen, m.openID(), e.edits)
	}
	settle(m, cmd)
	if m.prompt == nil || m.prompt.kind != "resolve" {
		t.Fatalf("the resolve prompt should be open:\n%s", plain(m))
	}
	// The resolve commits; the board read next finds Q-001 resolved.
	later := ended.Add(time.Minute)
	for i := range f.res.Groups[0].Versions {
		f.res.Groups[0].Versions[i].Record.Status, f.res.Groups[0].Versions[i].Record.Updated = "resolved", &later
	}
	settle(m, press(m, "y"))
	press(m, "esc", "esc")
	if m.screen != attemptScreen {
		t.Fatalf("Esc from the result, then from the question, returns to the attempt: screen %d", m.screen)
	}
	if s := plain(m); !strings.Contains(s, "question answered: R again") || !strings.Contains(s, "Waited on question Q-001, answered since.") || !strings.Contains(s, "o opens W-001: R launches the next attempt") {
		t.Fatalf("the answered attempt:\n%s", s)
	}
	// Only an answer written in or after the second the attempt ended, to a
	// question not asked after it, is the one it waited for.
	for _, c := range []struct {
		asked, written time.Time
		want           string
	}{
		{asked, later, "Q-001"},
		{asked, ended.Truncate(time.Second), "Q-001"},
		{asked, ended.Add(-time.Second), ""},
		{later, later, ""},
	} {
		for i := range f.res.Groups[0].Versions {
			f.res.Groups[0].Versions[i].Record.Created, f.res.Groups[0].Versions[i].Record.Updated = &c.asked, &c.written
		}
		if got := m.answeredSince("W-001", ended); got != c.want {
			t.Errorf("asked %v, written %v: %q, want %q", c.asked, c.written, got, c.want)
		}
	}
}
