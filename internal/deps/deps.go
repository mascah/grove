// Package deps interprets work dependencies in one checkout's records
// (G-161): the order a selection must follow, the layers and groups of the
// work shown, what each item needs and unlocks, and, on demand, what Git says
// about delivery and what the current view says about other versions. Only
// depends_on orders; members, priority and questions never do. The command
// and the board render the same View. Nothing here writes, starts, or adds
// work to a selection.
package deps

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"maps"
	"os/exec"
	"slices"
	"strings"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
	"github.com/mascah/grove/internal/versions"
)

// Item is one work record in a View: a row (selected, or unfinished in an
// overview) or, when Outside, a prerequisite of the rows that is not one.
type Item struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Status    string   `json:"status"`
	Revision  string   `json:"revision"`
	Candidate string   `json:"candidate,omitempty"`
	Outside   bool     `json:"outside"`   // a prerequisite listed, never added
	Layer     int      `json:"layer"`     // rows: the longest chain of rows beneath it
	Group     int      `json:"group"`     // rows: rows connected through dependencies share one, from 1
	Needs     []string `json:"needs"`     // its depends_on
	Unlocks   []string `json:"unlocks"`   // unfinished work whose depends_on names it
	NeededBy  []string `json:"needed_by"` // outside: the rows that need it
	Delivery  string   `json:"delivery"`  // what Deliver found; "" until it runs
}

// Question is an open question blocking a row or an outside prerequisite.
type Question struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Blocks []string `json:"blocks"`
}

// View is an overview (Selected nil) or a selection preview.
type View struct {
	Selected  []string   `json:"selected"` // as given; null for an overview
	Order     []string   `json:"order"`    // the rows, prerequisites first
	Items     []Item     `json:"items"`    // the rows in order, then the outside prerequisites by ID
	Questions []Question `json:"questions"`
	Notes     []string   `json:"notes"`
}

var unfinished = []string{"proposed", "active", "review"}

// Overview is the unfinished work in records, ordered by group, layer and
// ID, with the prerequisites of that work which are not unfinished.
func Overview(records []*project.Record) *View {
	var rows []string
	for _, r := range records {
		if r.Type == "work" && slices.Contains(unfinished, r.Status) {
			rows = append(rows, r.ID)
		}
	}
	slices.Sort(rows)
	return build(records, rows, nil)
}

// Preview is the explicit selection ids, in the order Order gives.
func Preview(records []*project.Record, ids []string) (*View, error) {
	order, _, err := Order(index(records), ids)
	if err != nil {
		return nil, err
	}
	return build(records, order, ids), nil
}

func index(records []*project.Record) map[string]*project.Record {
	byID := map[string]*project.Record{}
	for _, r := range records {
		byID[r.ID] = r
	}
	return byID
}

func build(records []*project.Record, rows, selected []string) *View {
	byID := index(records)
	reach := map[string]map[string]bool{}
	var visit func(from string, r *project.Record)
	visit = func(from string, r *project.Record) {
		for _, next := range r.DependsOn {
			if !reach[from][next] {
				reach[from][next] = true
				visit(from, byID[next])
			}
		}
	}
	for _, id := range rows {
		reach[id] = map[string]bool{}
		visit(id, byID[id])
	}
	// Layers and groups are among the rows: a row's layer is one more than
	// the deepest row it needs, and rows either of which needs the other,
	// directly or through anything, share a group.
	layer := map[string]int{}
	var depth func(id string) int
	depth = func(id string) int {
		if l, ok := layer[id]; ok {
			return l
		}
		l := 0
		for _, other := range rows {
			if reach[id][other] {
				l = max(l, depth(other)+1)
			}
		}
		layer[id] = l
		return l
	}
	group := map[string]int{}
	var join func(id string, g int)
	join = func(id string, g int) {
		group[id] = g
		for _, other := range rows {
			if group[other] == 0 && (reach[id][other] || reach[other][id]) {
				join(other, g)
			}
		}
	}
	groups := 0
	for _, id := range slices.Sorted(slices.Values(rows)) {
		if group[id] == 0 {
			groups++
			join(id, groups)
		}
	}
	if selected == nil { // rows arrive sorted by ID, which ties keep
		slices.SortStableFunc(rows, func(a, b string) int {
			return cmp.Or(cmp.Compare(group[a], group[b]), cmp.Compare(depth(a), depth(b)))
		})
	}

	v := &View{Selected: selected, Order: rows, Items: []Item{}, Questions: []Question{}, Notes: []string{}}
	item := func(r *project.Record) Item {
		it := Item{
			ID: r.ID, Title: r.Title, Status: r.Status, Revision: project.Revision(r.Source), Candidate: r.Candidate,
			Needs: append([]string{}, r.DependsOn...), Unlocks: []string{}, NeededBy: []string{},
		}
		for _, o := range records {
			if o.Type == "work" && slices.Contains(unfinished, o.Status) && slices.Contains(o.DependsOn, r.ID) {
				it.Unlocks = append(it.Unlocks, o.ID)
			}
		}
		return it
	}
	for _, id := range rows {
		it := item(byID[id])
		it.Layer, it.Group = depth(id), group[id]
		v.Items = append(v.Items, it)
	}
	// Outside: every prerequisite of a selection, or in an overview the
	// direct prerequisites that are not themselves rows (done, abandoned).
	neededBy := map[string][]string{}
	for _, id := range rows {
		needs := byID[id].DependsOn
		if selected != nil {
			needs = slices.Sorted(maps.Keys(reach[id]))
		}
		for _, p := range needs {
			if !slices.Contains(rows, p) && !slices.Contains(neededBy[p], id) {
				neededBy[p] = append(neededBy[p], id)
			}
		}
	}
	for _, p := range slices.Sorted(maps.Keys(neededBy)) {
		it := item(byID[p])
		it.Outside, it.NeededBy = true, neededBy[p]
		v.Items = append(v.Items, it)
	}

	listed := map[string]bool{}
	for _, it := range v.Items {
		listed[it.ID] = true
	}
	for _, r := range records {
		if r.Type != "question" || r.Status != "open" {
			continue
		}
		var blocks []string
		for _, id := range r.Blocks {
			if listed[id] {
				blocks = append(blocks, id)
			}
		}
		if blocks != nil {
			v.Questions = append(v.Questions, Question{ID: r.ID, Title: r.Title, Blocks: blocks})
		}
	}

	var waiting []string
	for _, it := range v.Items {
		if it.Outside && slices.Contains(unfinished, it.Status) {
			waiting = append(waiting, it.ID)
		}
		if it.Status == "abandoned" {
			var dependants []string
			for _, o := range v.Items {
				if slices.Contains(o.Needs, it.ID) {
					dependants = append(dependants, o.ID)
				}
			}
			if dependants != nil {
				v.Notes = append(v.Notes, fmt.Sprintf("%s is abandoned and will not be delivered: %s still needs it, which needs a decision.", it.ID, strings.Join(dependants, ", ")))
			}
		}
	}
	if waiting != nil {
		v.Notes = append(v.Notes, fmt.Sprintf("%s: unfinished and not selected, so the selection cannot finish until delivered; nothing was added.", strings.Join(waiting, ", ")))
	}
	if selected != nil {
		var pairs []string
		for i, a := range rows {
			for _, b := range rows[i+1:] {
				if !reach[a][b] && !reach[b][a] {
					pairs = append(pairs, a+" and "+b)
				}
			}
		}
		if pairs != nil {
			v.Notes = append(v.Notes, "No declared order between "+strings.Join(pairs, "; ")+". That is not evidence that they can proceed in parallel.")
		}
	}
	return v
}

// Deliver describes each item's delivery from its status and, for a
// candidate, Git ancestry: whether HEAD of the checkout, the base, contains
// it and, when target names a branch, whether the target does. contains
// answers one such question; Ancestry is the real one.
func (v *View) Deliver(target string, contains func(commit, ref string) (bool, error)) {
	for i := range v.Items {
		it := &v.Items[i]
		c := it.Candidate
		where := func() string {
			in, err := contains(c, "HEAD")
			if err != nil {
				return "candidate " + short(c) + " cannot be read here"
			}
			text := "candidate " + short(c) + map[bool]string{true: " in HEAD", false: " not in HEAD"}[in]
			if target != "" {
				on, err := contains(c, "refs/heads/"+target)
				switch {
				case err != nil:
					text += ", " + target + " unreadable"
				case on:
					text += ", on " + target
				default:
					text += ", not on " + target
				}
			}
			return text
		}
		switch {
		case it.Status == "proposed" || it.Status == "active":
			it.Delivery = "awaiting implementation"
		case it.Status == "review": // which requires a candidate
			it.Delivery = "awaiting review; " + where()
		case it.Status == "done" && c != "":
			it.Delivery = where()
		case it.Status == "done":
			it.Delivery = "done without a candidate: delivery unrecorded; establish it by Git ancestry or behaviour"
		case it.Status == "abandoned":
			it.Delivery = "abandoned: will not be delivered"
		}
	}
}

// Ancestry answers Deliver's question with git merge-base --is-ancestor in
// root's repository.
// ponytail: one Git process per question, on demand, for the few candidates
// a view names; batch them if a view ever names hundreds.
func Ancestry(ctx context.Context, root string) func(commit, ref string) (bool, error) {
	return func(commit, ref string) (bool, error) {
		_, err := repo.GitContext(ctx, root, "merge-base", "--is-ancestor", commit, ref)
		var exit *exec.ExitError
		if errors.As(err, &exit) && exit.ExitCode() == 1 { // Git's answer: not an ancestor
			return false, nil
		}
		return err == nil, err
	}
}

// sameSet compares prerequisite lists, whose order means nothing.
func sameSet(a, b []string) bool {
	return slices.Equal(slices.Compact(slices.Sorted(slices.Values(a))), slices.Compact(slices.Sorted(slices.Values(b))))
}

func short(commit string) string { return commit[:min(len(commit), 7)] }

// Compare adds what the current view (G-042) says about each item beyond
// here, the source whose records the View was built from, using res, an
// inspection of every branch and checkout. The View stays that source's:
// other versions, their edges included, are described and never merged in.
func (v *View) Compare(res *versions.Result, here *versions.Source) {
	note := func(format string, args ...any) { v.Notes = append(v.Notes, fmt.Sprintf(format, args...)) }
	if !res.Complete {
		note("Some branches or checkouts could not be read (grove versions lists them), so what is said here about other versions may be incomplete.")
	}
	for _, it := range v.Items {
		i := slices.IndexFunc(res.Groups, func(g versions.Group) bool { return g.ID == it.ID })
		if i < 0 || here == nil {
			continue
		}
		g := res.Groups[i]
		j := slices.IndexFunc(g.Versions, func(o versions.Version) bool { return o.Source == here })
		if j < 0 || g.Versions[j].Record == nil {
			continue
		}
		mine := g.Versions[j]
		if mine.Revision != it.Revision {
			note("%s changed while it was being read; rerun.", it.ID)
			continue
		}
		if mine.Change != "unchanged" {
			note("%s has uncommitted changes (%s) in this checkout, which is what is read here.", it.ID, mine.Change)
		}
		var others, edges []string
		seen := map[string]int{} // another current content: its index in others
		for _, o := range g.Versions {
			if o.Older != "" || o.Revision == mine.Revision {
				continue
			}
			if k, ok := seen[o.Revision]; ok {
				others[k] += ", " + o.Source.Name()
				continue
			}
			status := "deleted"
			if o.Record != nil {
				status = o.Record.Status
				if !sameSet(o.Record.DependsOn, mine.Record.DependsOn) {
					edges = append(edges, fmt.Sprintf("%s %q", o.Source.Name(), o.Record.DependsOn))
				}
			}
			seen[o.Revision] = len(others)
			others = append(others, status+" on "+o.Source.Name())
		}
		switch {
		case mine.Older != "":
			note("%s here is older than its current version: %s.", it.ID, strings.Join(others, "; "))
		case others != nil:
			note("%s diverges; other current versions: %s.", it.ID, strings.Join(others, "; "))
		}
		if edges != nil {
			note("%s depends on other work elsewhere: %s; this uses this checkout's %q.", it.ID, strings.Join(edges, "; "), mine.Record.DependsOn)
		}
		for _, n := range g.Notes {
			note("%s: %s", it.ID, n)
		}
	}
}

// Order validates the requested IDs and orders them so that each follows
// the selected work it depends on, directly or through unselected work; ties
// keep the requested order. reach maps each ID to its transitive prerequisites.
// Load has already refused dependency cycles and non-work targets.
func Order(byID map[string]*project.Record, ids []string) (order []string, reach map[string]map[string]bool, err error) {
	if len(ids) == 0 {
		return nil, nil, errors.New("context requires at least one work ID")
	}
	reach = map[string]map[string]bool{}
	for _, id := range ids {
		r := byID[id]
		switch {
		case reach[id] != nil:
			return nil, nil, fmt.Errorf("%s is selected more than once", id)
		case r == nil:
			return nil, nil, fmt.Errorf("work %s is not in this checkout; only explicit work IDs in the selected checkout can be selected", id)
		case r.Type != "work":
			return nil, nil, fmt.Errorf("%s is a %s; only work can be selected", id, r.Type)
		}
		reach[id] = map[string]bool{}
		var visit func(*project.Record)
		visit = func(r *project.Record) {
			for _, next := range r.DependsOn {
				if !reach[id][next] {
					reach[id][next] = true
					visit(byID[next])
				}
			}
		}
		visit(r)
	}
	placed := map[string]bool{}
	for len(order) < len(ids) {
		before := len(order)
		for _, id := range ids {
			ready := !placed[id]
			for _, other := range ids {
				ready = ready && (!reach[id][other] || placed[other])
			}
			if ready {
				placed[id] = true
				order = append(order, id)
				break
			}
		}
		if len(order) == before { // Load refuses cycles; never spin if one gets here
			return nil, nil, errors.New("the selected work depends on itself in a cycle")
		}
	}
	return order, reach, nil
}
