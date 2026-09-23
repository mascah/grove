package tui

import (
	"cmp"
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mascah/grove/internal/attempt"
	"github.com/mascah/grove/internal/versions"
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
	// An inspection already under way may predate the attempt's last commit,
	// so it is started again; a resolve or an action is left to finish.
	if ended && m.pending != "resolve" && m.pending != "act" {
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

// latest reports whether v is its work's newest attempt, or its work has
// none listed; it copies nothing, since every row of the list asks.
func (m *Model) latest(v *attempt.View) bool {
	for i := range m.attempts {
		if m.attempts[i].Launch.Work == v.Launch.Work {
			return m.attempts[i].Launch.Attempt == v.Launch.Attempt
		}
	}
	return true
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

// outcome says where an attempt stood when it ended, from its files and the
// records the board read, as a kind and a sentence. A clean exit alone is
// never a candidate: only the record in review with a candidate on the
// attempt's branch is.
func (m *Model) outcome(v *attempt.View) (kind, text string) {
	switch v.Status {
	case attempt.Running:
		return "running", "running"
	case attempt.Orphaned:
		return "orphaned", "orphaned: its owner is gone and the provider still runs (x stops it)"
	case attempt.Interrupted:
		return "interrupted", "interrupted: the owner and the provider are gone without a result"
	}
	r := v.Result
	if r == nil {
		return string(v.Status), string(v.Status)
	}
	work := v.Launch.Work
	review := r.Record != nil && r.Record.Status == "review" && r.Record.Candidate != ""
	// The record is read from the worktree's files, so it is the handoff only
	// if the owner found it committed when the process ended. Later commits
	// to the branch, such as an approval, change nothing about that.
	if review && !r.RecordUncommitted {
		return "candidate", fmt.Sprintf("candidate ready: %s in review on %s with candidate %s", work, v.Launch.Branch, short7(r.Record.Candidate))
	}
	unsaid := ""
	if review {
		unsaid = "; its record says review with candidate " + short7(r.Record.Candidate) + ", uncommitted"
	}
	exit := fmt.Sprintf("exit %d", r.ExitCode)
	if r.Signal != "" {
		exit = r.Signal
	}
	switch f := r.Events.Result; {
	case r.Stopped:
		return "stopped", "stopped (" + exit + ")" + unsaid
	case f == nil:
		return "failed", "failed: no result event (" + exit + ")" + unsaid
	case f.IsError || r.ExitCode != 0:
		return "failed", "failed: " + f.Subtype + " (" + exit + ")" + unsaid
	}
	// A question open now is this attempt's wait only if it is the work's
	// latest: an earlier one ended before whatever came after.
	if m.latest(v) {
		if q := m.blockingQuestion(work); q != "" {
			return "question", "waiting on question " + q + unsaid
		}
	}
	status, none := "unreadable", ", with no candidate"
	if r.Record != nil {
		status = r.Record.Status
	}
	if review {
		none = ""
	}
	return "unhanded", fmt.Sprintf("ended without a handoff: %s is %s on %s%s%s", work, status, v.Launch.Branch, none, unsaid)
}

func (m *Model) outcomeOf(v *attempt.View) string {
	_, text := m.outcome(v)
	return text
}

// Attempts fall in three groups (G-117): those that need the owner, those
// running, and the settled rest, each of which says why nothing is needed.
const (
	needsYou = iota
	runningNow
	settled
)

var groupNames = [...]string{"Needs you", "Running", "Settled"}

// tones colour each group from the ANSI 16 palette, as accents only: every
// group is also named, and every state is text.
var (
	tones = [...]lipgloss.Style{
		lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("4")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
	}
	faint = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	green = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	red   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	cyan  = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	pink  = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
)

// standing is where an attempt stands now: its group, the list's short
// state, and the attempt screen's State and Next lines (next "" is none).
type standing struct {
	group       int
	short       string
	state, next string
}

// standingOf weighs an attempt's outcome against the work's current record
// and its later attempts. Only the latest attempt of work still proposed,
// active or in review needs the owner, and an orphan always does, since its
// process runs unowned; a stopped attempt was the owner's own act.
func (m *Model) standingOf(v *attempt.View) standing {
	kind, text := m.outcome(v)
	work := v.Launch.Work
	s := standing{group: settled, state: sentence(text)}
	status, current := "", ""
	if g := m.groupOf(work); g != nil {
		if r := m.record(g); r != nil {
			status, current = r.Status, r.Candidate
		}
	}
	cand := ""
	if kind == "candidate" {
		cand = v.Result.Record.Candidate
	}
	latest := m.latest(v)
	switch {
	case kind == "running":
		s.group, s.short, s.next = runningNow, "running", "x stops it; the report arrives when it ends"
	case kind == "orphaned":
		s.group, s.short, s.next = needsYou, "orphaned: x stops it", "x stops it"
	case kind == "stopped":
		s.short = "stopped by x"
		if latest && (status == "proposed" || status == "active") {
			s.next = "R on " + work + " launches again"
		}
	case status == "done" && sameCommit(cand, current):
		s.short = "done: candidate " + short7(cand)
		s.state += " Candidate " + short7(cand) + " was integrated: " + work + " is done."
	case cand != "" && (!latest || status == "done" || status == "abandoned"):
		s.short = "candidate " + short7(cand) + ", superseded"
		s.state += " " + work + " has moved on: it is " + orUnread(status) + " with candidate " + short7(orUnread(current)) + "."
	case status == "done" || status == "abandoned":
		s.short = "work " + status + " since"
		s.state += " " + work + " is " + status + " now; nothing needs you."
	case !latest:
		s.short = shortOf(kind, v) + ", superseded"
		s.state += " A later attempt of " + work + " followed this one."
	default:
		s.group = needsYou
		switch kind {
		case "candidate":
			if status == "review" || status == "" {
				s.short, s.next = "judge candidate "+short7(cand), "o opens "+work+": a approves, f gives feedback"
			} else {
				s.short, s.next = "feedback given: R again", "o opens "+work+": R launches the next attempt"
				s.state += " " + work + " is " + status + " now: it had feedback."
			}
		case "question":
			q, _, _ := strings.Cut(m.blockingQuestion(work), " (")
			s.short, s.next = "answer question "+q, "o opens "+work+", whose detail lists the question"
		case "failed", "interrupted":
			s.short, s.next = shortOf(kind, v), "d shows the details and the raw log; R on "+work+" launches again"
		default:
			s.short, s.next = "ended, no handoff", "o opens "+work+"; the report says why"
		}
	}
	return s
}

// shortOf is an outcome kind in a few words.
func shortOf(kind string, v *attempt.View) string {
	switch kind {
	case "failed":
		f := v.Result.Events.Result
		switch {
		case f == nil:
			return "failed: no result"
		case f.Subtype == "error_max_budget_usd":
			return "failed: budget exhausted"
		case f.Subtype == "error_max_turns":
			return "failed: turn limit"
		}
		return "failed: " + strings.TrimPrefix(f.Subtype, "error_")
	case "question":
		return "waited on a question"
	case "unhanded":
		return "ended, no handoff"
	}
	return kind
}

// sameCommit compares two spellings of a commit, either possibly short.
func sameCommit(a, b string) bool {
	return a != "" && b != "" && (strings.HasPrefix(a, b) || strings.HasPrefix(b, a))
}

func orUnread(s string) string {
	if s == "" {
		return "unread"
	}
	return s
}

// sentence makes an outcome a sentence for the State line.
func sentence(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:] + "."
}

// when is how long an attempt has run, or how long ago it ended.
func (m *Model) when(v *attempt.View) string {
	switch now := m.clock(); {
	case v.Status == attempt.Running:
		return age(now.Sub(v.Launch.Started)) + " so far"
	case v.Result != nil && !v.Result.Finished.IsZero():
		return age(now.Sub(v.Result.Finished)) + " ago"
	default:
		return age(now.Sub(v.Launch.Started)) + " ago"
	}
}

func age(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "<1m"
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 48*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}

// titleOf is the title of the work an attempt concerns, as the board read it.
func (m *Model) titleOf(work string) string {
	if g := m.groupOf(work); g != nil {
		if r := m.record(g); r != nil {
			return r.Title
		}
	}
	return "(title unread)"
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
		parts = append(parts, fmt.Sprintf("Attempts %d · latest: %s, %s", len(mine), m.standingOf(&mine[0]).short, m.when(&mine[0])), "A lists them")
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

// listed is the attempts screen's list in its groups, newest first within
// each, their standings, and the cursor, the first row when the remembered
// one is not listed.
func (m *Model) listed() ([]attempt.View, []standing, int) {
	list := m.attemptsOf(m.listFor)
	order := make([]int, len(list))
	st := make([]standing, len(list))
	for i := range list {
		order[i], st[i] = i, m.standingOf(&list[i])
	}
	slices.SortStableFunc(order, func(a, b int) int { return st[a].group - st[b].group })
	views, standings := make([]attempt.View, len(list)), make([]standing, len(list))
	for i, j := range order {
		views[i], standings[i] = list[j], st[j]
	}
	at := slices.IndexFunc(views, func(v attempt.View) bool { return v.Launch.Attempt == m.listAt })
	if at < 0 && len(views) != 0 {
		at = 0
	}
	return views, standings, at
}

// openAttempt shows one attempt, read again now.
func (m *Model) openAttempt(id string) {
	if m.runID != id {
		m.run, m.activity, m.runErr = nil, attempt.Activity{}, ""
	}
	m.runBack, m.screen, m.runID, m.scroll, m.attemptsStale = m.screen, attemptScreen, id, 0, true
}

func (m *Model) attemptsKey(k string) {
	list, _, at := m.listed()
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
	case "d":
		m.facts = !m.facts
		m.clampScroll()
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
	back := m.screen
	m.openDetail(id)
	m.workBack, m.workDepth = back, len(m.stack) // Esc from this record returns to the attempt
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
	// This checkout is the live source in the Git directory the board was
	// read from: Start reads its record and checks it against Expect.
	var here *versions.Version
	for i := range g.Versions {
		if h := &g.Versions[i]; h.Source.Kind == "live" && h.Source.GitDir == m.res.GitDir {
			here = h
		}
	}
	if here == nil || here.Record == nil {
		m.notice = g.ID + " is not in this checkout (" + m.root + "); open Grove in a checkout that holds it to launch"
		return
	}
	req := attempt.Request{Root: m.root, ID: g.ID, Expect: here.Revision}
	// Work whose current state the base does not hold continues on that
	// state's branch. The base is the target, or without one this checkout;
	// the same bytes on several branches are one state, which starts afresh.
	onBase := v.OnTarget
	if m.res.Target == "" {
		onBase = v.Revision == here.Revision
	}
	ref := "refs/heads/worktree-" + g.ID // a fresh start reuses the default branch's checkout, wherever it is
	if b := branchOf(v); b != "" && !onBase {
		req.Branch, ref = b, v.Source.Ref
	}
	for _, s := range m.res.Sources {
		if s.Kind == "live" && s.Ref == ref {
			req.Branch, req.Worktree = strings.TrimPrefix(ref, "refs/heads/"), s.Worktree
		}
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

// attemptsBody is the attempts screen: one row per attempt in its group,
// with the work's ID and title, its state and its time. The title shrinks
// first and goes below 60 columns; the state is cut last.
func (m *Model) attemptsBody(w, n int) []string {
	list, st, cursor := m.listed()
	head := "Attempts in this repository"
	if m.listFor != "" {
		head = "Attempts of " + m.listFor
	}
	var count [len(groupNames)]int
	for _, s := range st {
		count[s.group]++
	}
	rows := []string{bold(line(fmt.Sprintf("%s (%d) · %d need you · %d running · %d settled", head, len(list), count[needsYou], count[runningNow], count[settled]), w))}
	if m.attemptsErr != "" {
		rows = append(rows, wrap("The attempts could not be read (r retries): "+m.attemptsErr, w)...)
	}
	if len(list) == 0 {
		rows = append(rows, line("  none yet; R on a proposed or active work record launches one", w))
	}
	const tw = 11 // "14m so far"
	idw, sw := 0, 0
	for i, v := range list {
		idw = max(idw, ansi.StringWidth(safe(v.Launch.Work)))
		sw = max(sw, ansi.StringWidth(safe(st[i].short)))
	}
	idw, sw = min(idw, 12), min(sw, 32)
	titleW := min(w-2-idw-2-sw-2-tw-2, 48) // wider than that, the rest is margin
	if w < 60 || titleW < 10 {
		titleW, sw = 0, max(w-2-idw-2-tw-2, 1)
	}
	at, group := 0, -1
	for i, v := range list {
		if st[i].group != group {
			group = st[i].group
			rows = append(rows, tones[group].Bold(true).Render(line(" "+groupNames[group], w)))
		}
		text := line(v.Launch.Work, idw) + "  "
		if titleW > 0 {
			text += line(m.titleOf(v.Launch.Work), titleW) + "  "
		}
		when := m.when(&v)
		text += line(st[i].short, sw) + "  " + line(strings.Repeat(" ", max(tw-len(when), 0))+when, tw)
		if i == cursor {
			at = len(rows)
			rows = append(rows, hot(clip("> "+text, w)))
		} else {
			rows = append(rows, clip("  "+text, w))
		}
	}
	off := 0
	if at >= n {
		off = at - n + 1
	}
	return rows[min(off, len(rows)):]
}

// attemptRows is the attempt screen (G-117): the work and a coloured state,
// the run's configuration and what it has used, State and Next, the
// details folded behind d, then the final report beside the activity,
// newest first, or above it where the terminal is narrow.
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
	st := m.standingOf(v)
	l, a := &v.Launch, &m.activity
	rows := []string{bold(line(l.Work+"  "+m.titleOf(l.Work), w))}
	glyph := [...]string{"◆", "●", "✓"}[st.group]
	badge := line(glyph+" "+st.short, min(ansi.StringWidth(safe(glyph+" "+st.short)), w))
	rows = append(rows, tones[st.group].Bold(true).Render(badge)+faint.Render(line("  "+m.when(v), w-ansi.StringWidth(badge))))
	if m.runErr != "" {
		rows = append(rows, wrap("The last read failed (r retries): "+m.runErr, w)...)
	}

	// The run's configuration, one labelled row each.
	blank := line("", w)
	rows = append(rows, blank)
	field := func(label string, value string, tone lipgloss.Style) {
		rows = append(rows, faint.Render(line("  "+label, 11))+tone.Render(line(value, w-11)))
	}
	field("Attempt", l.Attempt, lipgloss.NewStyle())
	model := cmp.Or(a.Model, l.Model, "not reported yet")
	if l.ClaudeVersion != "" {
		model += " · " + l.ClaudeVersion
	}
	field("Model", model, pink)
	spent, budget := "", l.BudgetUSD
	if r := v.Result; r != nil && r.Events.Result != nil {
		spent = fmt.Sprintf("$%.2f of ", r.Events.Result.CostUSD)
	}
	if spent == "" {
		spent = "spend known at the end, of "
	}
	field("Budget", spent+"$"+budget+" · permission mode "+l.PermissionMode, green)
	reuse := "created"
	if l.WorktreeReused {
		reuse = "reused"
	}
	field("Branch", l.Branch+" from "+l.Base[:min(len(l.Base), 12)]+" ("+reuse+")", cyan)
	started := l.Started.Local().Format("15:04:05 Mon 2 Jan")
	if r := v.Result; r != nil {
		started += " · ended " + r.Finished.Local().Format("15:04:05") + " after " + age(r.Finished.Sub(l.Started))
	}
	field("Started", started, lipgloss.NewStyle())
	rows = append(rows, m.metricRows(w)...)

	rows = append(rows, blank)
	for i, r := range wrap(st.state, w-7) {
		label := "       "
		if i == 0 {
			label = "State  "
		}
		rows = append(rows, faint.Render(label)+tones[st.group].Render(r))
	}
	if st.next != "" {
		for i, r := range wrap(st.next, w-7) {
			label := "       "
			if i == 0 {
				label = "Next   "
			}
			rows = append(rows, faint.Render(label)+r)
		}
	}
	if m.facts {
		rows = append(rows, blank, bold(line("Details · d hides them", w)))
		for _, f := range attempt.Facts(m.run, func(s string) string { return s }) {
			rows = append(rows, wrap(f, w)...)
		}
	} else {
		rows = append(rows, faint.Render(line("Details: bounds, provenance, events and raw files · d shows them", w)))
	}
	rows = append(rows, blank)

	report := func(w int) []string {
		out := []string{bold(line("Final report", w))}
		switch {
		case a.Report != "":
			return append(out, m.rendered("report\x00"+l.Attempt, a.Report, w)...)
		case v.Status == attempt.Running:
			return append(out, wrap("None yet: it arrives when the attempt ends.", w)...)
		}
		return append(out, wrap("None: the log's end holds no result event with a report.", w)...)
	}
	activity := func(w int) []string {
		out := []string{bold(line(fmt.Sprintf("Activity, newest first (%d)", len(a.Entries)), w))}
		if len(a.Entries) == 0 {
			out = append(out, line("  none yet", w))
		}
		for i := len(a.Entries) - 1; i >= 0; i-- {
			out = append(out, entryRow(a.Entries[i], w))
		}
		if a.Cut {
			for _, r := range wrap("earlier events are in "+m.run.EventsPath, w) {
				out = append(out, faint.Render(r))
			}
		}
		return out
	}
	if w < wideWidth {
		rows = append(rows, report(w)...)
		return append(append(rows, blank), activity(w)...)
	}
	lw := (w - 3) / 2
	left, right := report(lw), activity(w-3-lw)
	for i := range max(len(left), len(right)) {
		lr, rr := strings.Repeat(" ", lw), strings.Repeat(" ", w-3-lw)
		if i < len(left) {
			lr = clip(left[i], lw)
		}
		if i < len(right) {
			rr = right[i]
		}
		rows = append(rows, lr+faint.Render(" │ ")+rr)
	}
	return rows
}

// metricRows is what the run has used, as labelled values flowed into rows
// of w cells. A count taken from a cut window is a lower bound, shown ≥;
// an unknown value is –, never 0.
func (m *Model) metricRows(w int) []string {
	a := &m.activity
	mt := a.Metrics
	count := func(n int) string {
		if a.Cut {
			return "≥" + number(n)
		}
		return number(n)
	}
	tokens := count(mt.InputTokens) + " in · – out"
	if mt.Total {
		tokens = number(mt.InputTokens) + " in · " + number(mt.OutputTokens) + " out"
	}
	ctx := "–"
	if mt.Context > 0 {
		ctx = number(mt.Context)
		if mt.Window > 0 {
			ctx += " of " + number(mt.Window) + " " + gauge(float64(mt.Context)/float64(mt.Window))
		}
	}
	chips := [][2]string{
		{"Turns", count(mt.Turns)},
		{"Tokens", tokens},
		{"Context", ctx},
		{"Subagents", count(mt.Subagents)},
		{"Compactions", count(mt.Compactions)},
		{"Tools", count(mt.Tools)},
		{"Errors", count(mt.ToolErrors)},
	}
	var rows []string
	row, used := "  ", 2
	for _, c := range chips {
		cw := ansi.StringWidth(c[0]) + 1 + ansi.StringWidth(c[1]) + 3
		if used+cw > w && used > 2 {
			rows = append(rows, clip(row, w))
			row, used = "  ", 2
		}
		row += faint.Render(c[0]) + " " + bold(c[1]) + "   "
		used += cw
	}
	return append(rows, clip(row, w))
}

// number is a count in at most four characters and a unit.
func number(n int) string {
	switch {
	case n >= 1e6:
		return strings.TrimSuffix(fmt.Sprintf("%.1f", float64(n)/1e6), ".0") + "M"
	case n >= 1e4:
		return fmt.Sprintf("%dk", n/1000)
	case n >= 1e3:
		return strings.TrimSuffix(fmt.Sprintf("%.1f", float64(n)/1e3), ".0") + "k"
	}
	return fmt.Sprint(n)
}

// gauge is a fraction as ten cells and a percentage.
func gauge(f float64) string {
	f = min(max(f, 0), 1)
	n := int(f*10 + 0.5)
	return strings.Repeat("▰", n) + strings.Repeat("▱", 10-n) + fmt.Sprintf(" %d%%", int(f*100+0.5))
}

// kinds mark each kind of activity with a glyph and a colour.
var kinds = map[string]struct {
	glyph string
	tone  lipgloss.Style
}{
	"start":  {"▶", green},
	"text":   {"•", lipgloss.NewStyle()},
	"tool":   {"›", cyan},
	"error":  {"✗", red},
	"notice": {"·", faint},
	"result": {"■", green},
}

// entryRow is one activity entry in w cells: its local time, or blanks when
// the provider gave none, then the entry with its repeats counted.
func entryRow(e attempt.Entry, w int) string {
	stamp := strings.Repeat(" ", 8)
	if !e.Time.IsZero() {
		stamp = e.Time.Local().Format("15:04:05")
	}
	k, ok := kinds[e.Kind]
	if !ok {
		k = kinds["notice"]
	}
	if e.Kind == "result" && e.Text != "result: success" {
		k.tone = red
	}
	times := "" // kept whole: the text is cut before it
	if e.Count > 1 {
		times = fmt.Sprintf(" (×%d)", e.Count)
	}
	room, tw := max(w-9, 1), ansi.StringWidth(times)
	text := ansi.Truncate(safe(k.glyph+" "+e.Text), max(room-tw, 1), "…")
	pad := max(room-ansi.StringWidth(text)-tw, 0)
	return faint.Render(stamp) + " " + k.tone.Render(text) + faint.Render(times) + strings.Repeat(" ", pad)
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
		v, err := attempt.ShowContext(ctx, root, id, false) // the activity is the bounded read of the events
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
