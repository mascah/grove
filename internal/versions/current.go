package versions

import (
	"bytes"
	"container/heap"
	"fmt"
	"strconv"
	"strings"

	"github.com/mascah/grove/internal/project"
)

// The current view (G-042). Every valid source observes a record once: its
// exact bytes there, or its absence. A branch observes at its tip. A checkout
// whose file has its HEAD's bytes observes that commit; one whose bytes differ
// observes an uncommitted change on top of HEAD. Observation A is older than B
// when their bytes differ and the record at the merge base of their commits
// has A's bytes: since the two split, only B's side changed it. Current
// observations are those nothing is newer than; several current contents are
// genuine divergence. Timestamps, statuses, and which tip is newer never
// decide, and a revert is a change like any other. The set of sources is the
// same from every checkout, and so is the view.

// node is one distinct observation: a commit, whether an uncommitted change
// sits on top of it, and the record's content there ("" for absence).
type node struct {
	commit, content string
	live            bool
	sources         []*Source
	older           string // why another node is newer
}

// project marks each group's versions current or older, adds a row for a
// branch whose current state is the record's deletion, and notes pairs that
// could not be ordered. Only groups whose observations differ cost reads.
func (o *objects) project(res *Result) {
	for i := range res.Groups {
		o.projectGroup(res, &res.Groups[i])
	}
}

func (o *objects) projectGroup(res *Result, g *Group) {
	of := map[*Source]*Version{}
	for i := range g.Versions {
		of[g.Versions[i].Source] = &g.Versions[i]
	}
	var nodes []*node
	at := map[*Source]*node{}
	for _, s := range res.Sources {
		if !s.Valid {
			continue
		}
		n := node{commit: s.Commit}
		if v := of[s]; v != nil && v.Record != nil {
			n.content = v.Revision
		}
		if s.Kind == "live" {
			head, known := s.baseline.revisionOf(g.ID)
			n.live = !known || head != n.content
		}
		i := 0
		for i < len(nodes) && (nodes[i].commit != n.commit || nodes[i].live != n.live || nodes[i].content != n.content) {
			i++
		}
		if i == len(nodes) {
			nodes = append(nodes, &n)
		}
		nodes[i].sources = append(nodes[i].sources, s)
		at[s] = nodes[i]
	}
	for i, a := range nodes {
		for _, b := range nodes[i+1:] {
			if a.content == b.content {
				continue
			}
			base, err := o.baseOf(a, b, g.ID)
			switch {
			case err != nil:
				g.Notes = append(g.Notes, fmt.Sprintf("%s and %s could not be ordered: %v", name(a.sources[0]), name(b.sources[0]), err))
			case base == a.content && a.older == "":
				a.older = why(a, b)
			case base == b.content && b.older == "":
				b.older = why(b, a)
			}
		}
	}
	var versions []Version
	for _, s := range res.Sources {
		v, n := of[s], at[s]
		switch {
		case v != nil:
			if n != nil {
				v.Older = n.older
			}
			versions = append(versions, *v)
		case n != nil && n.content == "" && n.older == "" && s.Kind == "committed":
			// The branch's current state removes the record.
			versions = append(versions, Version{Source: s, Change: "deleted"})
		}
	}
	g.Versions = versions
}

// why says what makes a older than b.
func why(a, b *node) string {
	if a.commit == b.commit {
		return name(b.sources[0]) + " has an uncommitted change to it on top of this commit"
	}
	return name(b.sources[0]) + " changed it since their common history"
}

// name is a source as the reasons above word it.
func name(s *Source) string {
	ref := strings.TrimPrefix(s.Ref, "refs/heads/")
	if s.Kind == "committed" {
		return "branch " + ref
	}
	if ref == "" {
		ref = "detached"
	}
	return "checkout " + s.Locator + " (" + ref + ")"
}

// baseOf returns the record's content where a's and b's histories meet: at
// their merge base, or at their shared commit when they sit on one. No
// common history means the record was absent there. Several bases that
// disagree, or a base whose project cannot be read, are an error.
func (o *objects) baseOf(a, b *node, id string) (string, error) {
	bases := []string{a.commit}
	if a.commit != b.commit {
		var err error
		if bases, err = o.mergeBases(a.commit, b.commit); err != nil {
			return "", err
		}
	}
	content := ""
	for i, c := range bases {
		t := o.loadTree(c)
		got, known := t.revisionOf(id)
		switch {
		case t.err != nil:
			return "", fmt.Errorf("their common commit %s cannot be read: %v", short(c), t.err)
		case !known:
			return "", fmt.Errorf("their common commit %s holds a project that does not validate", short(c))
		case i > 0 && got != content:
			return "", fmt.Errorf("their common commits %s and %s hold different versions", short(bases[0]), short(c))
		}
		content = got
	}
	return content, nil
}

// revisionOf returns the revision of id in a commit's project, "" when the
// project or the record is absent, and known false when the project cannot
// be read or does not validate.
func (t *tree) revisionOf(id string) (revision string, known bool) {
	switch {
	case t == nil || t.err != nil:
		return "", false
	case !t.present:
		return "", true
	case !t.valid():
		return "", false
	}
	if t.ids == nil {
		t.ids = map[string]string{}
		for _, r := range t.project.Records {
			t.ids[r.ID] = project.Revision(r.Source)
		}
	}
	return t.ids[id], true
}

func short(commit string) string { return commit[:min(len(commit), 12)] }

// unborn reports an unborn branch's HEAD, which has no history.
func unborn(commit string) bool { return strings.Trim(commit, "0") == "" }

type commitInfo struct {
	parents []string
	when    int64
}

// commit reads one commit's parents and committer time, once.
func (o *objects) commit(id string) (commitInfo, error) {
	if c, ok := o.commitInfos[id]; ok {
		return c, nil
	}
	data, err := o.read(id, "commit")
	if err != nil {
		return commitInfo{}, err
	}
	var c commitInfo
	for line := range bytes.Lines(data) {
		line = bytes.TrimSuffix(line, []byte("\n"))
		if len(line) == 0 {
			break // the message follows
		}
		if p, ok := bytes.CutPrefix(line, []byte("parent ")); ok {
			c.parents = append(c.parents, string(p))
		} else if rest, ok := bytes.CutPrefix(line, []byte("committer ")); ok {
			// "<name> <email> <seconds> <zone>"
			fields := bytes.Fields(rest)
			if len(fields) >= 2 {
				c.when, _ = strconv.ParseInt(string(fields[len(fields)-2]), 10, 64)
			}
		}
	}
	o.commitInfos[id] = c
	return c, nil
}

// mergeBases returns the best common ancestors of a and b, as git merge-base
// --all would, by Git's paint-down-to-common walk in committer-date order,
// reading commits through the inspection's one cat-file process. Clock skew
// can leave a base that is an ancestor of another; baseOf treats disagreeing
// bases as unordered.
func (o *objects) mergeBases(a, b string) ([]string, error) {
	if unborn(a) || unborn(b) {
		return nil, nil
	}
	key := [2]string{min(a, b), max(a, b)}
	if bases, ok := o.bases[key]; ok {
		return bases, nil
	}
	const (
		fromA = 1 << iota
		fromB
		stale
		found
	)
	flags := map[string]int{a: fromA, b: fromB}
	var q commitQueue
	push := func(id string) error {
		c, err := o.commit(id)
		heap.Push(&q, queued{id, c.when})
		return err
	}
	if err := push(a); err != nil {
		return nil, err
	}
	if err := push(b); err != nil {
		return nil, err
	}
	var bases []string
	for q.nonStale(flags, stale) {
		id := heap.Pop(&q).(queued).id
		f := flags[id] & (fromA | fromB | stale)
		if f == fromA|fromB {
			if flags[id]&found == 0 {
				flags[id] |= found
				bases = append(bases, id)
			}
			f |= stale
		}
		c, err := o.commit(id)
		if err != nil {
			return nil, err
		}
		for _, p := range c.parents {
			if flags[p]&f == f {
				continue
			}
			flags[p] |= f
			if err := push(p); err != nil {
				return nil, err
			}
		}
	}
	var best []string
	for _, id := range bases {
		if flags[id]&stale == 0 { // not below another base
			best = append(best, id)
		}
	}
	best, err := o.independent(best)
	if err != nil {
		return nil, err
	}
	o.bases[key] = best
	return best, nil
}

// independent drops each base that is an ancestor of another, as Git does
// after its walk: commits sharing a timestamp, or clock skew, can let the walk
// stop before marking one. It runs only when there are several bases.
// ponytail: each check may walk all history below a base; use the walk's
// generation order if repositories with criss-cross merges make it slow.
func (o *objects) independent(bases []string) ([]string, error) {
	var keep []string
	for _, x := range bases {
		below := false
		for _, y := range bases {
			if x == y || below {
				continue
			}
			var err error
			if below, err = o.reaches(y, x); err != nil {
				return nil, err
			}
		}
		if !below {
			keep = append(keep, x)
		}
	}
	return keep, nil
}

// reaches reports whether x is y or one of its ancestors.
func (o *objects) reaches(y, x string) (bool, error) {
	seen := map[string]bool{y: true}
	for stack := []string{y}; len(stack) != 0; {
		id := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if id == x {
			return true, nil
		}
		c, err := o.commit(id)
		if err != nil {
			return false, err
		}
		for _, p := range c.parents {
			if !seen[p] {
				seen[p] = true
				stack = append(stack, p)
			}
		}
	}
	return false, nil
}

type queued struct {
	id   string
	when int64
}

// commitQueue pops the newest commit first, ties by ID, so walks repeat.
type commitQueue []queued

func (q commitQueue) Len() int { return len(q) }
func (q commitQueue) Less(i, j int) bool {
	if q[i].when != q[j].when {
		return q[i].when > q[j].when
	}
	return q[i].id < q[j].id
}
func (q commitQueue) Swap(i, j int) { q[i], q[j] = q[j], q[i] }
func (q *commitQueue) Push(x any)   { *q = append(*q, x.(queued)) }
func (q *commitQueue) Pop() any {
	old := *q
	x := old[len(old)-1]
	*q = old[:len(old)-1]
	return x
}

// nonStale reports a queued commit that is not yet known to be below a base.
func (q commitQueue) nonStale(flags map[string]int, stale int) bool {
	for _, c := range q {
		if flags[c.id]&stale == 0 {
			return true
		}
	}
	return false
}
