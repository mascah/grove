package versions

import (
	"bytes"
	"cmp"
	"container/heap"
	"fmt"
	"slices"
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
	path            string // where the record is, for reading it at a base
	sources         []*Source
	older           string // why another node is newer
	by              *node  // that newer node
	index           int    // in source order, which names pairs in notes
}

// project marks each group's versions current or older, adds a row for a
// branch whose current state is the record's deletion, and notes pairs that
// could not be ordered. Only groups whose observations differ cost reads.
func (o *objects) project(res *Result) {
	t := res.target()
	for i := range res.Groups {
		g := &res.Groups[i]
		o.projectGroup(res, g)
		if t == nil {
			continue
		}
		held := "" // the target's bytes, or absence
		for _, v := range g.Versions {
			if v.Source == t && v.Record != nil {
				held = v.Revision
			}
		}
		for j := range g.Versions {
			v := &g.Versions[j]
			v.OnTarget = held == v.Revision // a deletion's revision is ""
		}
	}
}

// target finds the integration target: the branch that every valid source
// naming one in grove.yaml agrees on. A source that names none has no say, so
// the answer is the same from every checkout, and a branch adding the key
// works before it merges. It returns the target's source, or nil with a note.
func (res *Result) target() *Source {
	var named []string
	first := map[string]*Source{}
	for _, s := range res.Sources {
		if t := s.target(); t != "" && first[t] == nil {
			first[t] = s
			named = append(named, t)
		}
	}
	if len(named) == 0 {
		return nil
	}
	if len(named) > 1 {
		var parts []string
		for _, t := range named {
			parts = append(parts, t+" on "+name(first[t]))
		}
		res.Notes = append(res.Notes, "grove.yaml names different targets ("+strings.Join(parts, ", ")+"), so none is used")
		return nil
	}
	t := named[0]
	i := slices.IndexFunc(res.Sources, func(s *Source) bool { return s.Kind == "committed" && s.Ref == "refs/heads/"+t })
	switch {
	case i < 0:
		res.Notes = append(res.Notes, "grove.yaml names target "+t+", which is not a local branch, so none is used")
		return nil
	case len(res.Sources[i].Diagnostics) != 0:
		res.Notes = append(res.Notes, "target branch "+t+" could not be read, so none is used")
		return nil
	}
	// A target without a project yet, as while Grove is adopted on a branch,
	// lacks every record.
	res.Target = t
	return res.Sources[i]
}

func (s *Source) target() string {
	if !s.Valid {
		return ""
	}
	return s.project.Target
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
			n.content, n.path = v.Revision, v.Path
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
			n.index = i
			nodes = append(nodes, &n)
		}
		nodes[i].sources = append(nodes[i].sources, s)
		at[s] = nodes[i]
	}
	// A pair can only mark one of its nodes older, so a node already older
	// needs no more comparisons of its own. Candidates for current are found
	// first, keeping every pair's order: each node meets the candidates so far
	// and drops those it is newer than. Then each remaining candidate meets
	// every node, so nothing newer than a current node goes unseen, without
	// assuming that older is transitive.
	type pair struct{ older, newer *node }
	rels := map[[2]*node]pair{}
	// rel orders a pair once: the older node and the newer, or neither when
	// their bytes match or their common commit cannot be read (a note).
	rel := func(a, b *node) pair {
		if a.content == b.content {
			return pair{}
		}
		if b.index < a.index {
			a, b = b, a
		}
		if p, ok := rels[[2]*node{a, b}]; ok {
			return p
		}
		var p pair
		base, err := o.baseOf(a, b, g.ID)
		switch {
		case err != nil:
			g.Notes = append(g.Notes, fmt.Sprintf("%s and %s could not be ordered: %v", name(a.sources[0]), name(b.sources[0]), err))
		case base == a.content:
			p = pair{a, b}
		case base == b.content:
			p = pair{b, a}
		}
		rels[[2]*node{a, b}] = p
		return p
	}
	order := func(a, b *node) {
		if a.older != "" && b.older != "" {
			return
		}
		if p := rel(a, b); p.older != nil && p.older.older == "" {
			p.older.older, p.older.by = why(p.older, p.newer), p.newer
		}
	}
	// Newer commits tend to supersede, so meeting them first drops the rest
	// after one comparison each. The date orders the work, never the result.
	byDate := slices.Clone(nodes)
	when := func(n *node) int64 {
		if unborn(n.commit) {
			return 0
		}
		c, err := o.commit(n.commit)
		if err != nil { // a cancelled read; the result is discarded
			return 0
		}
		return c.when
	}
	slices.SortStableFunc(byDate, func(a, b *node) int {
		if a.live != b.live {
			return map[bool]int{true: -1, false: 1}[a.live]
		}
		return cmp.Compare(when(b), when(a)) // ties keep the sources' order
	})
	// ponytail: n diverging states cost n² merge-base walks (1,001 branches
	// editing one record: about 1 s more per load); group sources by content
	// before comparing if real repositories diverge that widely.
	var candidates []*node
	for _, n := range byDate {
		for i := 0; i < len(candidates) && n.older == ""; {
			if order(candidates[i], n); candidates[i].older != "" {
				candidates = slices.Delete(candidates, i, i+1)
			} else {
				i++
			}
		}
		if n.older == "" {
			candidates = append(candidates, n)
		}
	}
	for _, c := range candidates {
		for _, n := range byDate {
			if c.older == "" {
				order(c, n)
			}
		}
	}
	// Reverts carried across merges can make older a cycle. A node whose
	// reasons lead back to it rather than to a current node may belong to a
	// cycle that nothing outside is newer than: then the whole relation
	// decides, and a node is current when everything newer than it, through
	// any chain, is also older than it.
	loops := func(n *node) bool {
		seen := map[*node]bool{}
		for ; n.older != ""; n = n.by {
			if seen[n] {
				return true
			}
			seen[n] = true
		}
		return false
	}
	if slices.ContainsFunc(nodes, loops) {
		// ponytail: orders every pair and walks each node's reach, O(n³) at
		// worst; only records whose history holds such a cycle pay it.
		newer := map[*node][]*node{}
		for i, a := range nodes {
			for _, b := range nodes[i+1:] {
				if p := rel(a, b); p.older != nil {
					newer[p.older] = append(newer[p.older], p.newer)
				}
			}
		}
		reach := map[*node]map[*node]bool{}
		for _, n := range nodes {
			r := map[*node]bool{}
			for stack := []*node{n}; len(stack) != 0; {
				x := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				for _, y := range newer[x] {
					if !r[y] {
						r[y] = true
						stack = append(stack, y)
					}
				}
			}
			reach[n] = r
		}
		restored := false
		for _, n := range nodes {
			sink := true
			for m := range reach[n] {
				sink = sink && reach[m][n]
			}
			if sink && n.older != "" {
				n.older, n.by, restored = "", nil, true
			}
		}
		if restored {
			g.Notes = append(g.Notes, "some versions could not be ordered: each is older than another, through changes and reverts that merges carried across branches")
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
	var bases []string
	switch {
	case a.commit != b.commit:
		var err error
		if bases, err = o.mergeBases(a.commit, b.commit); err != nil {
			return "", err
		}
	case !unborn(a.commit): // an unborn HEAD has no history
		bases = []string{a.commit}
	}
	content := ""
	for i, c := range bases {
		got, err := o.recordAt(c, id, a, b)
		switch {
		case err != nil:
			return "", err
		case i > 0 && got != content:
			return "", fmt.Errorf("their common commits %s and %s hold different versions", short(bases[0]), short(c))
		}
		content = got
	}
	return content, nil
}

// recordAt returns id's revision at commit, "" when it is absent there. It
// reads the paths the compared observations have the record at, where it
// almost always was; only when neither holds it does it load the whole
// project there, which tells absence from a record that moved.
func (o *objects) recordAt(commit, id string, a, b *node) (string, error) {
	key := [2]string{commit, id}
	if r, ok := o.records[key]; ok {
		return r.revision, r.err
	}
	revision, err := o.readRecordAt(commit, id, a, b)
	o.records[key] = recordRead{revision, err}
	return revision, err
}

type recordRead struct {
	revision string
	err      error
}

func (o *objects) readRecordAt(commit, id string, a, b *node) (string, error) {
	for _, n := range []*node{a, b} {
		if n.path == "" {
			continue
		}
		data, err := o.fileAt(commit, n.path)
		if err != nil {
			return "", fmt.Errorf("their common commit %s cannot be read: %v", short(commit), err)
		}
		if data == nil {
			continue
		}
		// Bytes equal to a side's are that record; others need their ID read.
		revision := project.Revision(data)
		if revision == a.content || revision == b.content {
			return revision, nil
		}
		if r, _ := project.ParseRecord(n.path, data); r.ID == id {
			return revision, nil
		}
	}
	t := o.loadTree(commit)
	got, known := t.revisionOf(id)
	switch {
	case t.err != nil:
		return "", fmt.Errorf("their common commit %s cannot be read: %v", short(commit), t.err)
	case !known:
		return "", fmt.Errorf("their common commit %s holds a project that does not validate", short(commit))
	}
	return got, nil
}

// fileAt returns the regular file at path below the prefix in commit, or nil.
func (o *objects) fileAt(commit, path string) ([]byte, error) {
	es, err := o.entries(commit + "^{tree}")
	if err != nil {
		return nil, err
	}
	parts := strings.Split(strings.Trim(o.prefix+path, "/"), "/")
	for i, part := range parts {
		e, ok := entryNamed(es, part)
		switch {
		case !ok:
			return nil, nil
		case i == len(parts)-1:
			if e.mode != "100644" && e.mode != "100755" {
				return nil, nil
			}
			return o.blob(e.id)
		case !e.isDir():
			return nil, nil
		}
		if es, err = o.entries(e.id); err != nil {
			return nil, err
		}
	}
	return nil, nil
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

// commitNode is one commit read for merge bases, its parents linked as they are
// read. flags belong to the walk numbered epoch, so a walk starts clean
// without clearing the graph.
type commitNode struct {
	id      string
	read    bool
	parents []*commitNode
	when    int64
	epoch   int
	flags   int
}

// commit returns id's node, reading its parents and committer time once.
func (o *objects) commit(id string) (*commitNode, error) {
	c := o.node(id)
	if c.read {
		return c, nil
	}
	data, err := o.read(id, "commit")
	if err != nil {
		return nil, err
	}
	for line := range bytes.Lines(data) {
		line = bytes.TrimSuffix(line, []byte("\n"))
		if len(line) == 0 {
			break // the message follows
		}
		if p, ok := bytes.CutPrefix(line, []byte("parent ")); ok {
			c.parents = append(c.parents, o.node(string(p)))
		} else if rest, ok := bytes.CutPrefix(line, []byte("committer ")); ok {
			// "<name> <email> <seconds> <zone>"
			fields := bytes.Fields(rest)
			if len(fields) >= 2 {
				c.when, _ = strconv.ParseInt(string(fields[len(fields)-2]), 10, 64)
			}
		}
	}
	c.read = true
	return c, nil
}

func (o *objects) node(id string) *commitNode {
	c := o.graph[id]
	if c == nil {
		c = &commitNode{id: id}
		o.graph[id] = c
	}
	return c
}

// flag returns c's flags in the current walk.
func (o *objects) flag(c *commitNode) int {
	if c.epoch != o.epoch {
		return 0
	}
	return c.flags
}

func (o *objects) mark(c *commitNode, f int) {
	c.flags, c.epoch = o.flag(c)|f, o.epoch
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
	o.epoch++
	var q commitQueue
	for i, id := range []string{a, b} {
		c, err := o.commit(id)
		if err != nil {
			return nil, err
		}
		o.mark(c, []int{fromA, fromB}[i])
		heap.Push(&q, c)
	}
	var bases []*commitNode
	for q.nonStale(o, stale) {
		c := heap.Pop(&q).(*commitNode)
		f := o.flag(c) & (fromA | fromB | stale)
		if f == fromA|fromB {
			if o.flag(c)&found == 0 {
				o.mark(c, found)
				bases = append(bases, c)
			}
			f |= stale
		}
		for _, p := range c.parents {
			if o.flag(p)&f == f {
				continue
			}
			if _, err := o.commit(p.id); err != nil {
				return nil, err
			}
			o.mark(p, f)
			heap.Push(&q, p)
		}
	}
	var best []*commitNode
	for _, c := range bases {
		if o.flag(c)&stale == 0 { // not below another base
			best = append(best, c)
		}
	}
	best, err := o.independent(best)
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(best))
	for i, c := range best {
		ids[i] = c.id
	}
	o.bases[key] = ids
	return ids, nil
}

// independent drops each base that is an ancestor of another, as Git does
// after its walk: commits sharing a timestamp, or clock skew, can let the walk
// stop before marking one. It runs only when there are several bases.
// ponytail: each check may walk all history below a base; use the walk's
// generation order if repositories with criss-cross merges make it slow.
func (o *objects) independent(bases []*commitNode) ([]*commitNode, error) {
	if len(bases) < 2 {
		return bases, nil
	}
	var keep []*commitNode
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
func (o *objects) reaches(y, x *commitNode) (bool, error) {
	seen := map[*commitNode]bool{y: true}
	for stack := []*commitNode{y}; len(stack) != 0; {
		c := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if c == x {
			return true, nil
		}
		if _, err := o.commit(c.id); err != nil {
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

// commitQueue pops the newest commit first, ties by ID, so walks repeat.
type commitQueue []*commitNode

func (q commitQueue) Len() int { return len(q) }
func (q commitQueue) Less(i, j int) bool {
	if q[i].when != q[j].when {
		return q[i].when > q[j].when
	}
	return q[i].id < q[j].id
}
func (q commitQueue) Swap(i, j int) { q[i], q[j] = q[j], q[i] }
func (q *commitQueue) Push(x any)   { *q = append(*q, x.(*commitNode)) }
func (q *commitQueue) Pop() any {
	old := *q
	x := old[len(old)-1]
	*q = old[:len(old)-1]
	return x
}

// nonStale reports a queued commit that is not yet known to be below a base.
func (q commitQueue) nonStale(o *objects, stale int) bool {
	for _, c := range q {
		if o.flag(c)&stale == 0 {
			return true
		}
	}
	return false
}

// Name is how notes name a source: "branch B" or "checkout L (B)".
func (s *Source) Name() string { return name(s) }
