package tui

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/mascah/grove/internal/attempt"
)

// Managed attempts (G-046): R on proposed or active work launches one
// bounded attempt through G-045's attempt.Start, after a typed budget and
// permission mode; A lists attempts and Enter opens one, showing its facts,
// an outcome derived from them, the final report and recent activity; x
// stops one. The board only observes: each attempt's owner process and files
// are G-045's, so closing the board changes nothing and reopening it reads
// the same attempts again.

// attemptsMsg is one read of the repository's attempts and, when the
// attempt screen is open, of that attempt in full.
type attemptsMsg struct {
	views    []attempt.View
	err      error
	open     string
	run      *attempt.View
	activity attempt.Activity
	runErr   error
}

// attemptTick asks for the next read while an attempt runs.
type attemptTick struct{}

// pollEvery is how often attempts are re-read while one runs.
const pollEvery = 2 * time.Second

// wantAttempts starts a read of the attempts when one is due and none is in
// flight. It is separate from the one-at-a-time Git reads, so it never
// cancels or waits for them.
func (m *Model) wantAttempts() tea.Cmd {
	// The list is read from the common directory the board found, so it
	// starts no Git process; only the open attempt's read does, cancellably.
	if m.backend.Attempts == nil || m.done || m.attemptsReading || !m.attemptsStale || m.res == nil {
		return nil
	}
	m.attemptsReading, m.attemptsStale = true, false
	root, dir, open := m.root, filepath.Join(m.res.Repository, "grove", "attempts"), ""
	if m.screen == attemptScreen {
		open = m.runID
	}
	ctx, backend := m.ctx, m.backend
	return func() tea.Msg {
		if !m.reads.begin() {
			return nil
		}
		defer m.reads.wg.Done()
		msg := attemptsMsg{open: open}
		msg.views, msg.err = backend.Attempts(ctx, dir)
		if open != "" && backend.Attempt != nil {
			msg.run, msg.activity, msg.runErr = backend.Attempt(ctx, root, open)
		}
		return msg
	}
}

// gotAttempts takes a read in. An attempt that was running and no longer is
// has changed its branch, so the board is re-read unless a read that must
// finish is under way; while any attempt runs, the next read is scheduled.
func (m *Model) gotAttempts(msg attemptsMsg) tea.Cmd {
	m.attemptsReading = false
	ended := false
	if msg.err != nil {
		m.attemptsErr = msg.err.Error()
	} else {
		for _, v := range msg.views {
			if old := m.attemptOf(v.Launch.Attempt); old != nil && live(old) && !live(&v) {
				ended = true
			}
		}
		m.attempts, m.attemptsErr = msg.views, ""
	}
	if msg.open != "" && msg.open == m.runID {
		m.run, m.activity, m.runErr = msg.run, msg.activity, ""
		if msg.runErr != nil {
			m.runErr = msg.runErr.Error()
		}
	}
	m.clampScroll()
	var cmds []tea.Cmd
	if ended && !m.busy() {
		cmds = append(cmds, m.inspect())
	}
	if !m.ticking && slices.ContainsFunc(m.attempts, func(v attempt.View) bool { return live(&v) }) {
		m.ticking = true
		cmds = append(cmds, tea.Tick(m.every, func(time.Time) tea.Msg { return attemptTick{} }))
	}
	return tea.Batch(cmds...)
}

// live reports an attempt whose process may still be running.
func live(v *attempt.View) bool { return v.Status == attempt.Running || v.Status == attempt.Orphaned }

func (m *Model) attemptOf(id string) *attempt.View {
	i := slices.IndexFunc(m.attempts, func(v attempt.View) bool { return v.Launch.Attempt == id })
	if i < 0 {
		return nil
	}
	return &m.attempts[i]
}

// attemptsOf lists one work's attempts, newest first; "" lists every one.
func (m *Model) attemptsOf(work string) []attempt.View {
	if work == "" {
		return m.attempts
	}
	var out []attempt.View
	for _, v := range m.attempts {
		if v.Launch.Work == work {
			out = append(out, v)
		}
	}
	return out
}

// outcomeOf says where an attempt stands, from its files and the records the
// board read. A clean exit alone is never a candidate: only the record in
// review with a candidate on the attempt's branch is.
func (m *Model) outcomeOf(v *attempt.View) string {
	switch v.Status {
	case attempt.Running:
		return "running"
	case attempt.Orphaned:
		return "orphaned: its owner is gone and the provider still runs (x stops it)"
	case attempt.Interrupted:
		return "interrupted: the owner and the provider are gone without a result"
	}
	r := v.Result
	if r == nil {
		return string(v.Status)
	}
	work := v.Launch.Work
	if r.Record != nil && r.Record.Status == "review" && r.Record.Candidate != "" {
		return fmt.Sprintf("candidate ready: %s in review on %s with candidate %s", work, v.Launch.Branch, short7(r.Record.Candidate))
	}
	if q := m.blockingQuestion(work); q != "" {
		return "waiting on question " + q
	}
	exit := fmt.Sprintf("exit %d", r.ExitCode)
	if r.Signal != "" {
		exit = r.Signal
	}
	switch f := r.Events.Result; {
	case r.Stopped:
		return "stopped (" + exit + ")"
	case f == nil:
		return "failed: no result event (" + exit + ")"
	case f.IsError || r.ExitCode != 0:
		return "failed: " + f.Subtype + " (" + exit + ")"
	}
	status := "unreadable"
	if r.Record != nil {
		status = r.Record.Status
	}
	return fmt.Sprintf("ended without a handoff: %s is %s on %s, with no candidate", work, status, v.Launch.Branch)
}

// blockingQuestion names an open question that blocks work in the records
// the board shows, or "".
func (m *Model) blockingQuestion(work string) string {
	if m.res == nil {
		return ""
	}
	for i := range m.res.Groups {
		if q := m.record(&m.res.Groups[i]); q != nil && q.Type == "question" && q.Status == "open" && slices.Contains(q.Blocks, work) {
			return q.ID + " (" + q.Title + ")"
		}
	}
	return ""
}

// attemptTag marks a card whose work has an attempt that may be running.
func (m *Model) attemptTag(work string) string {
	for _, v := range m.attempts {
		if v.Launch.Work == work && live(&v) {
			return map[attempt.Status]string{attempt.Running: "● running", attempt.Orphaned: "● orphaned"}[v.Status]
		}
	}
	return ""
}

// attemptRow is a work detail's header row: its latest attempt and the keys.
func (m *Model) attemptRow(work, status string) string {
	mine := m.attemptsOf(work)
	var parts []string
	if len(mine) == 0 {
		parts = append(parts, "Attempts: none")
	} else {
		parts = append(parts, fmt.Sprintf("Attempts: %d, latest %s %s", len(mine), mine[0].Launch.Attempt, m.outcomeOf(&mine[0])), "A lists them")
	}
	if m.backend.Launch != nil && (status == "proposed" || status == "active") {
		parts = append(parts, "R launches one")
	}
	return strings.Join(parts, " · ")
}

// openAttempts shows the attempts screen for work, or for every work.
func (m *Model) openAttempts(work string) {
	if m.backend.Attempts == nil {
		return
	}
	m.listBack, m.screen, m.listFor, m.scroll, m.attemptsStale = m.screen, attemptsScreen, work, 0, true
}

// listed is the attempts screen's list and its cursor, the first row when
// the remembered one is not listed.
func (m *Model) listed() ([]attempt.View, int) {
	list := m.attemptsOf(m.listFor)
	at := slices.IndexFunc(list, func(v attempt.View) bool { return v.Launch.Attempt == m.listAt })
	if at < 0 && len(list) != 0 {
		at = 0
	}
	return list, at
}

// openAttempt shows one attempt, read again now.
func (m *Model) openAttempt(id string) {
	if m.runID != id {
		m.run, m.activity, m.runErr = nil, attempt.Activity{}, ""
	}
	m.runBack, m.screen, m.runID, m.scroll, m.attemptsStale = m.screen, attemptScreen, id, 0, true
}

func (m *Model) attemptsKey(k string) {
	list, at := m.listed()
	switch k {
	case "up", "k", "down", "j":
		if len(list) == 0 {
			return
		}
		if k == "up" || k == "k" {
			at--
		} else {
			at++
		}
		m.listAt = list[min(max(at, 0), len(list)-1)].Launch.Attempt
	case "enter":
		if at >= 0 {
			m.openAttempt(list[at].Launch.Attempt)
		}
	case "x":
		if at >= 0 {
			m.askStop(&list[at])
		}
	case "o":
		if at >= 0 {
			m.openWork(list[at].Launch.Work)
		}
	}
}

func (m *Model) attemptKey(k string) {
	switch k {
	case "x":
		if v := m.attemptOf(m.runID); v != nil {
			m.askStop(v)
		}
	case "o":
		if v := m.attemptOf(m.runID); v != nil {
			m.openWork(v.Launch.Work)
		}
	default:
		m.scrollKey(k)
	}
}

// openWork opens an attempt's work record, where its reviews and evidence are.
func (m *Model) openWork(id string) {
	if m.groupOf(id) == nil {
		m.notice = id + " is not on any readable branch or checkout"
		return
	}
	m.openDetail(id)
}

// askStop asks before stopping an attempt that may be running.
func (m *Model) askStop(v *attempt.View) {
	switch {
	case m.backend.Stop == nil:
	case !live(v):
		m.notice = v.Launch.Attempt + " is " + string(v.Status) + "; nothing to stop"
	default:
		m.prompt = &prompt{kind: "stop", id: v.Launch.Work, attempt: v.Launch.Attempt, root: m.root}
	}
}

// launch opens the launch prompt for the open work record, or says why not.
// Where it runs is decided now, from what the board shows: the branch the
// record stands on and that branch's checkout, so a continuation after
// feedback runs on its candidate's branch; and the revision of the record in
// this checkout, which Start refuses to launch from if it changed since.
func (m *Model) launch() {
	g := m.group()
	if g == nil || m.backend.Launch == nil {
		return
	}
	v := m.shown(g)
	switch {
	case v == nil || v.Record == nil || v.Record.Type != "work":
		m.notice = g.ID + " is not work: only work is launched"
		return
	case v.Record.Status == "review":
		m.notice = g.ID + " is in review: judge its candidate (a approve, f feedback) before another attempt"
		return
	case v.Record.Status != "proposed" && v.Record.Status != "active":
		m.notice = g.ID + " is " + v.Record.Status + "; only proposed or active work is launched"
		return
	}
	if q := m.blockingQuestion(g.ID); q != "" {
		m.notice = g.ID + " is blocked by open question " + q + "; resolve it before another attempt"
		return
	}
	for _, a := range m.attemptsOf(g.ID) {
		if live(&a) {
			m.notice = fmt.Sprintf("attempt %s of %s is %s; A shows it, x stops it", a.Launch.Attempt, g.ID, a.Status)
			return
		}
	}
	// The base is the target, or without one this checkout's branch; work
	// standing on any other branch continues there.
	req, base := attempt.Request{Root: m.root, ID: g.ID}, m.res.Target
	for i := range g.Versions {
		if h := &g.Versions[i]; h.Source.Kind == "live" && h.Source.Locator == "." {
			if h.Record != nil {
				req.Expect = h.Revision
			}
			if base == "" {
				base = branchOf(h)
			}
		}
	}
	if b := branchOf(v); b != "" && b != base {
		req.Branch = b
		for _, s := range m.res.Sources {
			if s.Kind == "live" && s.Ref == v.Source.Ref {
				req.Worktree = s.Worktree
			}
		}
	}
	if req.Expect == "" {
		m.notice = g.ID + " is not in this checkout (" + m.root + "); open Grove in a checkout that holds it to launch"
		return
	}
	m.prompt = &prompt{kind: "budget", id: g.ID, root: m.root, req: &req}
}

// where says where a launch will run, for its prompt.
func (p *prompt) where() string {
	switch {
	case p.req.Worktree != "":
		return "on branch " + p.req.Branch + " in " + p.req.Worktree
	case p.req.Branch != "":
		return "on branch " + p.req.Branch
	}
	return "on branch worktree-" + p.id
}

// launchKey takes Enter in the launch prompt: the budget, checked, then the
// permission mode, then the launch. Neither has a default (G-045).
func (m *Model) launchKey(p *prompt) tea.Cmd {
	text := strings.TrimSpace(p.text)
	switch {
	case p.kind == "budget" && !attempt.ValidBudget(text):
		m.notice = "the budget is a positive dollar amount, such as 2 or 0.5"
	case p.kind == "budget":
		p.req.BudgetUSD, p.kind, p.text = text, "mode", ""
	case text == "" || strings.ContainsAny(text, " \t"):
		m.notice = "type one permission mode, such as acceptEdits or auto, or Esc"
	default:
		p.req.PermissionMode = text
		return m.act(p)
	}
	return nil
}

// attemptsBody is the attempts screen: one row per attempt, newest first.
func (m *Model) attemptsBody(w, n int) []string {
	list, cursor := m.listed()
	head := "Attempts in this repository"
	if m.listFor != "" {
		head = "Attempts of " + m.listFor
	}
	rows := []string{bold(line(fmt.Sprintf("%s  (%d)", head, len(list)), w))}
	if m.attemptsErr != "" {
		rows = append(rows, wrap("The attempts could not be read (r retries): "+m.attemptsErr, w)...)
	}
	if len(list) == 0 {
		rows = append(rows, line("  none yet; R on a proposed or active work record launches one", w))
	}
	at := 0
	for i, v := range list {
		cost := ""
		if r := v.Result; r != nil && r.Events.Result != nil {
			cost = fmt.Sprintf("  $%.2f", r.Events.Result.CostUSD)
		}
		text := fmt.Sprintf("%s  %s%s  %s", v.Launch.Attempt, m.outcomeOf(&v), cost, v.Launch.Branch)
		if i == cursor {
			at = len(rows)
			rows = append(rows, hot(line("> "+text, w)))
		} else {
			rows = append(rows, line("  "+text, w))
		}
	}
	off := 0
	if at >= n {
		off = at - n + 1
	}
	return rows[min(off, len(rows)):]
}

// attemptRows is the attempt screen: facts, the outcome, the final report,
// then recent activity newest first, so the report never scrolls away and
// the newest activity is always nearest the facts.
func (m *Model) attemptRows(w int) []string {
	v := m.run
	if v == nil || v.Launch.Attempt != m.runID {
		if m.runErr != "" {
			return wrapAll("The attempt "+m.runID+" could not be read (r retries): "+m.runErr, w)
		}
		return []string{line("Reading attempt "+m.runID+"…", w)}
	}
	if cur := m.attemptOf(v.Launch.Attempt); cur != nil && cur.Status != v.Status {
		v = cur // the list read is newer; its status leads until the next full read
	}
	rows := []string{bold(line("Attempt "+v.Launch.Attempt+" of "+v.Launch.Work, w))}
	rows = append(rows, wrap("Outcome: "+m.outcomeOf(v), w)...)
	if m.runErr != "" {
		rows = append(rows, wrap("The last read failed (r retries): "+m.runErr, w)...)
	}
	for _, f := range attempt.Facts(m.run, func(s string) string { return s }) {
		rows = append(rows, wrap(f, w)...)
	}
	if report := m.activity.Report; report != "" {
		rows = append(rows, line("", w), bold(line("Final report", w)))
		rows = append(rows, m.rendered("report\x00"+v.Launch.Attempt, report, w)...)
	}
	lines := m.activity.Lines
	rows = append(rows, line("", w), bold(line(fmt.Sprintf("Activity, newest first (%d)", len(lines)), w)))
	if len(lines) == 0 {
		rows = append(rows, line("  none yet", w))
	}
	for i := len(lines) - 1; i >= 0; i-- {
		rows = append(rows, line("  "+lines[i], w))
	}
	if m.activity.Cut {
		rows = append(rows, wrap("  earlier events are in "+m.run.EventsPath, w)...)
	}
	return rows
}

// launched is what a launch reports once it is running.
func launched(l *attempt.Launch) []string {
	return []string{
		fmt.Sprintf("attempt: %s started; owner pid %d, session %s, budget %s USD, permission mode %s", l.Attempt, l.Owner, l.SessionID, l.BudgetUSD, l.PermissionMode),
		"it runs without this board: closing Grove leaves it running, A lists attempts, x stops one",
	}
}

// liveAttempts is the attempts part of the backend over the real repository.
func liveAttempts(b *Backend) {
	b.Attempts = func(_ context.Context, dir string) ([]attempt.View, error) { return attempt.ListDir(dir, "") }
	b.Attempt = func(ctx context.Context, root, id string) (*attempt.View, attempt.Activity, error) {
		v, err := attempt.ShowContext(ctx, root, id)
		if err != nil {
			return nil, attempt.Activity{}, err
		}
		a, err := attempt.ReadActivity(v.EventsPath, attempt.ActivityWindow)
		return v, a, err
	}
	b.Launch = func(_ context.Context, req attempt.Request) ([]string, error) {
		var facts []string
		l, err := attempt.Start(req, time.Now(), func(f string) { facts = append(facts, f) })
		if err != nil {
			return facts, err
		}
		return append(facts, launched(l)...), nil
	}
	b.Stop = func(_ context.Context, root, id string) ([]string, error) {
		var facts []string
		err := attempt.Stop(root, id, func(f string) { facts = append(facts, f) })
		return facts, err
	}
}
