package tui

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/charmbracelet/x/ansi"

	"github.com/mascah/grove/internal/attempt"
	"github.com/mascah/grove/internal/handoff"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/versions"
)

// The review view (G-044): a work record in review shows the facts a judgment
// needs, its changed files and their diffs on demand, and three explicit
// actions, each confirmed, run through the same functions the CLI offers in
// the checkout the facts name. Nothing else the board does writes.

type changesRead struct {
	c   *versions.Changes
	err error
}

type diffRead struct {
	text string
	err  error
}

type changesMsg struct {
	gen int
	key string
	c   *versions.Changes
	err error
}

type diffMsg struct {
	gen  int
	key  string
	text string
	err  error
}

type actMsg struct {
	gen   int
	kind  string
	about string // the record, or the attempt a stop was for
	facts []string
	err   error
}

// prompt is the open question on the last row: text to type for a verdict or
// feedback, or y/n for the merge and its cleanup.
type prompt struct {
	kind             string // approve, feedback, integrate, cleanup, launch, stop, resolve
	text             string
	id, root, branch string // the record, the checkout the action runs in, the branch judged
	target, wt       string // integrate: the target and the branch's checkout, if any
	cleanup          bool
	req              *attempt.Request    // launch: what Start is asked for, resolved on Enter
	run              project.RunDefaults // launch: grove.yaml's defaults in this checkout
	attempt          string              // stop: the attempt
	expect           string              // resolve: the question's revision after the editor
}

// outcome is what an action returned, shown on the result screen until Esc.
type outcome struct {
	title string
	facts []string
	err   string
}

// reviewable reports the detail of a work record in review: where the
// actions apply.
func (m *Model) reviewable() bool {
	if m.screen != detailScreen {
		return false
	}
	r := m.openRecord()
	return r != nil && r.Type == "work" && r.Status == "review"
}

// openRecord is the record the detail shows now, if any.
func (m *Model) openRecord() *project.Record {
	g := m.group()
	if g == nil {
		return nil
	}
	if v := m.shown(g); v != nil {
		return v.Record
	}
	return nil
}

// changesKey names one changes read: the candidate, the tip it is judged
// against, the target, and the record's path.
func changesKey(target string, v *versions.Version) string {
	return v.Record.Candidate + "\x00" + v.Source.Commit + "\x00" + target + "\x00" + v.Path
}

// wantChanges starts reading the shown record's changes when its detail is
// open, it has a candidate, no other read is pending, and they are not held.
func (m *Model) wantChanges() tea.Cmd {
	if m.backend.Changes == nil || m.done || m.screen != detailScreen || m.pending != "" && m.pending != "changes" {
		return nil
	}
	g := m.group()
	if g == nil {
		return nil
	}
	v := m.shown(g)
	if v == nil || v.Record == nil || v.Record.Candidate == "" {
		return nil
	}
	key := changesKey(m.res.Target, v)
	if _, held := m.changes[key]; held || key == m.reading {
		return nil
	}
	target, candidate, tip, path := m.res.Target, v.Record.Candidate, v.Source.Commit, v.Path
	cmd := m.read("changes", func(ctx context.Context, gen int) tea.Msg {
		c, err := m.backend.Changes(ctx, m.root, target, candidate, tip, path)
		return changesMsg{gen, key, c, err}
	})
	m.reading = key
	return cmd
}

// diffKey names one diff read.
func diffKey(from, to, path string) string { return from + "\x00" + to + "\x00" + path }

// diffOf is the diff the content pane shows, if one is chosen: the file
// against the changes' base.
func (m *Model) diffOf(v *versions.Version) (from, to, path string, ok bool) {
	if m.diff == "" || v == nil || v.Record == nil {
		return "", "", "", false
	}
	read, held := m.changes[changesKey(m.res.Target, v)]
	if !held || read.c == nil || read.c.Base == "" {
		return "", "", "", false
	}
	return read.c.Base, v.Record.Candidate, m.diff, true
}

// wantDiff starts reading the chosen file's diff when nothing else is pending.
func (m *Model) wantDiff() tea.Cmd {
	if m.backend.Diff == nil || m.done || m.screen != detailScreen || m.pending != "" && m.pending != "diff" {
		return nil
	}
	g := m.group()
	if g == nil {
		return nil
	}
	from, to, path, ok := m.diffOf(m.shown(g))
	if !ok {
		return nil
	}
	key := diffKey(from, to, path)
	if _, held := m.diffs[key]; held || key == m.reading {
		return nil
	}
	cmd := m.read("diff", func(ctx context.Context, gen int) tea.Msg {
		text, err := m.backend.Diff(ctx, m.root, from, to, path)
		return diffMsg{gen, key, text, err}
	})
	m.reading = key
	return cmd
}

// branchOf names the branch a version stands on, or "".
func branchOf(v *versions.Version) string {
	if v == nil || v.Source.Ref == "" {
		return ""
	}
	return strings.TrimPrefix(v.Source.Ref, "refs/heads/")
}

// projectDir is the project directory inside a live source's checkout.
func (m *Model) projectDir(s *versions.Source) string {
	return filepath.Join(s.Worktree, filepath.FromSlash(m.res.Prefix))
}

// judgeRoot finds the checkout in which the shown version can be judged,
// since approval and feedback commit the record's file there, or the reason
// there is none.
func (m *Model) judgeRoot(g *versions.Group, v *versions.Version) (root, branch, why string) {
	lv, branch, why := m.checkoutOf(g, v, "judging")
	if lv == nil {
		return "", branch, why
	}
	return m.projectDir(lv.Source), branch, ""
}

// checkoutOf finds the checkout that holds the shown version to write it:
// the one valid checkout on its branch holding the record. It returns that
// copy, or the reason there is none. doing is "judging", which needs the
// copy to match HEAD, or "answering", whose uncommitted copy is the owner's
// answer so far.
func (m *Model) checkoutOf(g *versions.Group, v *versions.Version, doing string) (lv *versions.Version, branch, why string) {
	judging := doing == "judging"
	if v == nil || v.Record == nil {
		return nil, "", "the current state holds no record to " + map[bool]string{true: "judge", false: "answer"}[judging]
	}
	branch = branchOf(v)
	if branch == "" {
		return nil, "", "the current state is on no branch; " + map[bool]string{true: "a and f need", false: "e needs"}[judging] + " a branch checkout"
	}
	where := func(s *versions.Source) string {
		if s.Locator == "." {
			return "checkout ."
		}
		return "checkout " + s.Locator
	}
	dirty := func(s *versions.Source) string {
		return "the record has uncommitted changes in " + where(s) + "; commit or discard them before " + doing + " it"
	}
	if v.Source.Kind == "live" {
		if judging && v.Change != "unchanged" {
			return nil, branch, dirty(v.Source)
		}
		return v, branch, ""
	}
	var found []*versions.Version
	for _, s := range m.res.Sources {
		if s.Kind != "live" || s.Ref != v.Source.Ref || !s.Valid {
			continue
		}
		for i := range g.Versions {
			if cv := &g.Versions[i]; cv.Source == s {
				found = append(found, cv)
			}
		}
	}
	switch {
	case len(found) > 1:
		return nil, branch, fmt.Sprintf("%d checkouts are on branch %s, so which one to write is ambiguous", len(found), branch)
	case len(found) == 0 && judging:
		return nil, branch, "no checkout is on branch " + branch + "; git worktree add one, or run grove approve there"
	case len(found) == 0:
		return nil, branch, "no checkout is on branch " + branch + "; git worktree add one"
	case found[0].Record == nil || judging && found[0].Change != "unchanged":
		return nil, branch, dirty(found[0].Source)
	}
	return found[0], branch, ""
}

// targetRoot finds the target's checkout, where integrate runs.
func (m *Model) targetRoot() (root, why string) {
	if m.res.Target == "" {
		return "", "grove.yaml names no target branch, so nothing can be integrated"
	}
	for _, s := range m.res.Sources {
		if s.Kind == "live" && s.Ref == "refs/heads/"+m.res.Target && s.Valid {
			return m.projectDir(s), ""
		}
	}
	return "", "no checkout is on the target " + m.res.Target + "; i needs one"
}

// targetTip is the target branch's commit as the board read it, or "".
func (m *Model) targetTip() string {
	for _, s := range m.res.Sources {
		if s.Kind == "committed" && s.Ref == "refs/heads/"+m.res.Target {
			return s.Commit
		}
	}
	return ""
}

// reviewRows are the header's Review block: the candidate's standing and
// where each action would run. Facts only; every action is a key away.
func (m *Model) reviewRows(g *versions.Group, v *versions.Version) []string {
	r := v.Record
	parts := []string{"Review: candidate " + short7(r.Candidate)}
	if r.Approved != "" {
		parts = append(parts, "approved")
	} else {
		parts = append(parts, "not yet approved")
	}
	switch read, held := m.changes[changesKey(m.res.Target, v)]; {
	case m.backend.Changes == nil:
	case !held:
		parts = append(parts, "reading its changes…")
	case read.err != nil:
		parts = append(parts, "changes unreadable (r retries)")
	case len(read.c.After) != 0:
		parts = append(parts, fmt.Sprintf("%d other %s changed since it: the tip is a new candidate", len(read.c.After), plural(len(read.c.After), "file", "files")))
	default:
		parts = append(parts, "only the record changed since it")
	}
	if read, held := m.changes[changesKey(m.res.Target, v)]; held && read.c != nil && read.c.Merge != nil {
		// The prediction names the target commit it read, which the board's
		// own reading of the target may no longer be (G-177).
		text := read.c.Merge.Text(m.res.Target)
		if tip := m.targetTip(); tip != "" && tip != read.c.Merge.Target {
			text += "; the board read " + m.res.Target + " at " + short7(tip) + " (r re-reads)"
		}
		parts = append(parts, text)
	}
	rows := []string{strings.Join(parts, " · ")}
	var where []string
	if root, branch, why := m.judgeRoot(g, v); why != "" {
		where = append(where, "a approve and f feedback: "+why)
	} else {
		where = append(where, "a approve and f feedback run on branch "+branch+" in "+root)
	}
	if m.backend.Integrate != nil {
		if root, why := m.targetRoot(); why != "" {
			where = append(where, "i integrate: "+why)
		} else {
			where = append(where, "i integrate runs into "+m.res.Target+" in "+root)
		}
	}
	return append(rows, strings.Join(where, " · "))
}

// changeEntries lists the changed files as sidebar entries, in Git's order.
func (m *Model) changeEntries(v *versions.Version) []entry {
	if v == nil || v.Record == nil || v.Record.Candidate == "" {
		return nil
	}
	read, held := m.changes[changesKey(m.res.Target, v)]
	if !held || read.c == nil {
		return nil
	}
	var out []entry
	for i := range read.c.Files {
		out = append(out, entry{change: &read.c.Files[i]})
	}
	return out
}

// changesSection is the sidebar's Changes heading and rows, with item
// drawing each selectable file.
func (m *Model) changesSection(v *versions.Version, w int, heading func(string), item func(string), plain func(string)) {
	if m.backend.Changes == nil || v == nil || v.Record == nil || v.Record.Candidate == "" {
		return
	}
	read, held := m.changes[changesKey(m.res.Target, v)]
	switch {
	case !held:
		heading("Changes")
		plain("  reading…")
	case read.err != nil:
		heading("Changes")
		plain("  The changes could not be read (r retries): " + read.err.Error())
	case m.res.Target == "":
		heading("Changes")
		plain("  grove.yaml names no target, so there is no base to list files against")
	default:
		heading(fmt.Sprintf("Changes against %s from %s", m.res.Target, short7(read.c.Base)))
		if len(read.c.Files) == 0 {
			plain("  none")
		}
		for _, f := range read.c.Files {
			counts := "binary"
			if f.Added >= 0 {
				counts = fmt.Sprintf("+%d −%d", f.Added, f.Removed)
			}
			item(ansi.Truncate(safe(f.Path), max(w-2-ansi.StringWidth(counts)-2, 8), "…") + "  " + counts)
			plain("    " + m.describedBy(f.Path))
		}
		if len(read.c.After) != 0 {
			plain("  after the candidate: " + strings.Join(read.c.After, ", "))
		}
	}
}

// describedBy names the records, other than the open one, that link a
// changed file or name it in a code span (G-153), each once at its first
// tier, in the inspection's order. It reads the loaded records only: no Git
// process, nothing stored. A file is a path from the repository's top, or a
// rename's two; the prefix, which ends in a slash, makes it a project path,
// and a file outside the project can only be named by a code span.
// ponytail: matched on every frame, about 20 ms for 60 files over 150
// records; cache by changes key if reviews grow larger.
func (m *Model) describedBy(file string) string {
	var found []string
	for i := range m.res.Groups {
		g := &m.res.Groups[i]
		if g.ID == m.openID() {
			continue
		}
		best := -1
		vs, _ := m.searched(g)
		for _, v := range vs {
			ms := m.mentionsOf(v)
			for _, p := range strings.Split(file, " → ") {
				tier := -1
				if rest, in := strings.CutPrefix(p, m.res.Prefix); in {
					tier, _ = pathTier(rest, ms)
				} else if slices.ContainsFunc(ms, func(mn handoff.Mention) bool { return mn.Span != "" && names(mn.Span, p) }) {
					tier = 2
				}
				if tier >= 0 && (best < 0 || tier < best) {
					best = tier
				}
			}
		}
		if best >= 0 {
			found = append(found, g.ID+" "+tiers[best])
		}
	}
	if len(found) == 0 {
		return "no record names it"
	}
	return "described by " + strings.Join(found, ", ")
}

// diffRows renders a diff for the content pane: every line escaped like
// record text, then an accent by its first character that nothing depends
// on. ponytail: rendered on every frame; cache by key and width if diffs
// grow large.
func diffRows(text string, w int) []string {
	var rows []string
	for _, l := range strings.Split(strings.TrimRight(text, "\n"), "\n") {
		row := line(strings.ReplaceAll(l, "\t", "    "), w)
		switch {
		case strings.HasPrefix(l, "+") && !strings.HasPrefix(l, "+++"):
			row = "\x1b[32m" + row + "\x1b[m"
		case strings.HasPrefix(l, "-") && !strings.HasPrefix(l, "---"):
			row = "\x1b[31m" + row + "\x1b[m"
		case strings.HasPrefix(l, "@@"):
			row = "\x1b[36m" + row + "\x1b[m"
		}
		rows = append(rows, row)
	}
	return rows
}

// diffContent is what the content pane shows for the chosen file.
func (m *Model) diffContent(v *versions.Version, w int) []string {
	from, to, path, ok := m.diffOf(v)
	if !ok {
		return wrapAll("The changes are not held any more; r rereads them.", w)
	}
	read, held := m.diffs[diffKey(from, to, path)]
	switch {
	case !held:
		return wrapAll("reading the diff…", w)
	case read.err != nil:
		return wrapAll("The diff could not be read (r retries): "+read.err.Error(), w)
	case strings.TrimSpace(read.text) == "":
		return wrapAll("No textual difference.", w)
	}
	return diffRows(read.text, w)
}

// startAtEvidence scrolls a review's content to its Evidence heading, so the
// first screen is the handoff.
func (m *Model) startAtEvidence() {
	g := m.group()
	if g == nil {
		return
	}
	v := m.shown(g)
	if v == nil || v.Record == nil || v.Record.Type != "work" || v.Record.Status != "review" {
		return
	}
	w := m.width
	if w >= wideWidth {
		w = w * 11 / 20
	}
	for i, row := range m.contentRows(v, w) {
		if strings.HasPrefix(strings.TrimSpace(ansi.Strip(row)), "## Evidence") {
			m.dscroll = i
			m.clampScroll()
			return
		}
	}
}

// action opens the prompt for a, f or i on a work record in review, or says
// why it cannot.
func (m *Model) action(k string) {
	g := m.group()
	v := m.shown(g)
	if !m.reviewable() {
		if v != nil && v.Record != nil {
			m.notice = g.ID + " is not in review: nothing to approve, give feedback on, or integrate"
		}
		return
	}
	r := v.Record
	switch k {
	case "a", "f":
		root, branch, why := m.judgeRoot(g, v)
		if why != "" {
			m.notice = why
			return
		}
		if k == "a" && r.Approved != "" {
			m.notice = "candidate " + short7(r.Candidate) + " is already approved; i integrates it"
			return
		}
		m.prompt = &prompt{kind: map[string]string{"a": "approve", "f": "feedback"}[k], id: g.ID, root: root, branch: branch}
	case "i":
		if m.backend.Integrate == nil {
			return
		}
		if r.Approved == "" {
			m.notice = "approve candidate " + short7(r.Candidate) + " first (a)"
			return
		}
		root, why := m.targetRoot()
		if why != "" {
			m.notice = why
			return
		}
		wt, branch, _ := m.judgeRoot(g, v)
		m.prompt = &prompt{kind: "integrate", id: g.ID, root: root, branch: branch, target: m.res.Target, wt: wt}
	}
}

// promptKey handles every key while a prompt is open: text goes into it,
// Enter and y/n answer it, Esc cancels it.
func (m *Model) promptKey(msg tea.KeyPressMsg) tea.Cmd {
	p := m.prompt
	switch k := msg.String(); {
	case k == "esc":
		m.prompt = nil
		m.notice = "cancelled; nothing was written"
		if p.req != nil {
			m.notice = "cancelled; nothing was launched"
		} else if p.kind == "stop" {
			m.notice = "cancelled; nothing was stopped"
		} else if p.kind == "resolve" {
			m.notice = unresolved(p)
		}
	case msg.Text == "y" && p.kind == "stop":
		return m.act(p)
	case msg.Text == "n" && p.kind == "stop":
		m.prompt = nil
		m.notice = "cancelled; nothing was stopped"
	case msg.Text == "y" && p.kind == "resolve":
		return m.act(p)
	case msg.Text == "n" && p.kind == "resolve":
		m.prompt = nil
		m.notice = unresolved(p)
	case p.kind == "approve" || p.kind == "feedback" || p.req != nil:
		switch k {
		case "enter":
			if p.req != nil {
				return m.launchKey(p)
			}
			if strings.TrimSpace(p.text) == "" {
				m.notice = "type the " + map[string]string{"approve": "verdict", "feedback": "feedback"}[p.kind] + " first, or Esc"
				return nil
			}
			return m.act(p)
		case "backspace":
			if n := len(p.text); n != 0 {
				_, size := lastRune(p.text)
				p.text = p.text[:n-size]
			}
		default:
			if msg.Text != "" {
				p.text += msg.Text
			}
		}
	case msg.Text == "y" && p.kind == "integrate":
		if p.wt == "" {
			return m.act(p) // no worktree: nothing to ask about
		}
		p.kind = "cleanup"
	case msg.Text == "y" && p.kind == "cleanup":
		p.cleanup = true
		return m.act(p)
	case msg.Text == "n" && p.kind == "cleanup":
		return m.act(p)
	case msg.Text == "n":
		m.prompt = nil
		m.notice = "cancelled; nothing was merged"
	}
	return nil
}

// act runs the prompt's action in the background. Keys wait for it, and a
// key never cancels its Git commands.
func (m *Model) act(p *prompt) tea.Cmd {
	m.prompt = nil
	kind, root, id, text, cleanup, req, about, expect := p.kind, p.root, p.id, strings.TrimSpace(p.text), p.cleanup, p.req, p.id, p.expect
	switch {
	case kind == "cleanup":
		kind = "integrate"
	case req != nil:
		kind = "launch"
	case kind == "stop":
		about = p.attempt
	}
	m.resultBack = m.screen
	cmd := m.read("act", func(ctx context.Context, gen int) tea.Msg {
		var facts []string
		var err error
		switch kind {
		case "launch":
			facts, err = m.backend.Launch(ctx, *req)
		case "stop":
			facts, err = m.backend.Stop(ctx, root, about)
		case "approve":
			facts, err = m.backend.Approve(ctx, root, id, text)
		case "feedback":
			facts, err = m.backend.Feedback(ctx, root, id, text)
		case "resolve":
			facts, err = m.backend.Answer(ctx, root, id, expect)
		default:
			facts, err = m.backend.Integrate(ctx, root, id, cleanup)
		}
		return actMsg{gen, kind, about, facts, err}
	})
	m.acting = kind
	return cmd
}

// acting names the running action for the banner.
func actingText(kind string) string {
	return map[string]string{"approve": "Approving…", "feedback": "Recording the feedback…", "integrate": "Integrating…",
		"launch": "Launching the attempt…", "stop": "Stopping the attempt…", "resolve": "Resolving and committing…"}[kind]
}

// promptRow is the last row while a prompt is open.
func (m *Model) promptRow(w int) string {
	p := m.prompt
	var text string
	switch p.kind {
	case "approve":
		text = fmt.Sprintf("Approve %s on branch %s · verdict (Enter records it, Esc cancels): %s▏", p.id, p.branch, p.text)
	case "feedback":
		text = fmt.Sprintf("Feedback on %s, returning it to active on branch %s (Enter records it, Esc cancels): %s▏", p.id, p.branch, p.text)
	case "integrate":
		// The question first: a long checkout path is what truncation drops.
		text = fmt.Sprintf("Merge branch %s into %s and mark %s done? y/n   (runs in %s)", p.branch, p.target, p.id, p.root)
	case "launch":
		req, _ := p.resolved() // a line that does not parse yet shows what it has so far
		// What is typed comes first: truncation drops the help at the end.
		text = fmt.Sprintf("Launch %s %s▏ · %s, %s · Enter launches; flags such as --until plan or --effort xhigh change it; Esc cancels", p.id, p.text, launchText(req), p.where())
	case "stop":
		text = fmt.Sprintf("Stop attempt %s of %s? Its partial work stays. y/n", p.attempt, p.id)
	case "resolve":
		text = fmt.Sprintf("Resolve %s and commit it with your answer on branch %s? y/n   (runs in %s)", p.id, p.branch, p.root)
	default:
		text = fmt.Sprintf("Also delete branch %s and remove its worktree? y/n   (%s)", p.branch, p.wt)
	}
	return hot(line(text, w))
}

// resultRows is the result screen: the action's facts, and the refusal or
// failure when there was one.
func (m *Model) resultRows(w int) []string {
	o := m.result
	rows := []string{bold(line(o.title, w))}
	for _, f := range o.facts {
		rows = append(rows, wrap("  "+f, w)...)
	}
	if o.err != "" {
		failed := wrap("NOT DONE: "+o.err, w)
		failed[0] = bold(failed[0])
		rows = append(append(rows, line("", w)), failed...)
	}
	if m.pending == "inspect" {
		return append(rows, line("", w), line("Re-reading the board…", w))
	}
	return append(rows, line("", w), line("The board has been re-read. Esc returns to the record.", w))
}

// unresolved says what declining to resolve leaves.
func unresolved(p *prompt) string {
	return p.id + " is not resolved; your edit stays uncommitted in " + p.root + ", and e reopens it to resolve"
}

func short7(commit string) string { return commit[:min(len(commit), 7)] }
