package attempt

import (
	"context"
	"crypto/sha256"
	"fmt"
	"slices"
	"strings"

	"github.com/mascah/grove/internal/deps"
	"github.com/mascah/grove/internal/project"
)

// Selection is the explicit set of work one attempt runs (G-162, decision
// G-188): the IDs as given, the order deps gives them, each member as the
// launching checkout held it and whether it can start, the prerequisites
// outside it with their delivery at the base, and a digest over what the
// launch would run. The attempt implements the members in order under one
// budget; nothing outside the selection is added.
type Selection struct {
	Selected []string  `json:"selected"`
	Order    []string  `json:"order"`
	Members  []Member  `json:"members"` // in order
	Outside  []Outside `json:"outside,omitempty"`
	Notes    []string  `json:"notes,omitempty"`
	Digest   string    `json:"digest"`
}

// Member is one selected work record at launch.
type Member struct {
	ID       string `json:"id"`
	Path     string `json:"record_path"`     // project-relative
	Revision string `json:"record_revision"` // as the launching checkout held it
	Status   string `json:"status"`          // where the attempt starts from
	Wait     string `json:"wait,omitempty"`  // why it cannot start; "" when it can
}

// Outside is a prerequisite of the selection that is not in it: listed,
// never implemented or added.
type Outside struct {
	ID       string   `json:"id"`
	Status   string   `json:"status"`
	Delivery string   `json:"delivery"`
	NeededBy []string `json:"needed_by"`
}

// Continuation is the policy every selection runs under, shown by the
// preview and the attempt; the work guide carries it to the agent.
const Continuation = "members run in order, one at a time; a member that waits is not started, and a new question, an outside blocker or a failure stops that member and every selected member that needs it while the rest continue; budget exhaustion or Stop ends the attempt with what is committed kept"

// Boundary is the review boundary every selection has (G-188).
const Boundary = "one review at the end: the complete members enter review together on one shared candidate, each judged on its own acceptance and integrated as a group; a started member left incomplete holds the whole branch out of review"

// Members lists the attempt's selected work in order; an attempt launched
// before selections is a selection of its one work.
func (l *Launch) Members() []Member {
	if l.Selection != nil {
		return l.Selection.Members
	}
	return []Member{{ID: l.Work, Path: l.RecordPath, Revision: l.RecordRevision}}
}

// Includes reports whether id is one of the attempt's members.
func (l *Launch) Includes(id string) bool {
	return slices.ContainsFunc(l.Members(), func(m Member) bool { return m.ID == id })
}

// selectionOf interprets ids in p, where the attempt starts from base: each
// member must be proposed or active work, and waits on an open question, on
// an outside prerequisite base does not definitely hold, or on a selected
// prerequisite that waits. contains answers whether base holds a commit.
func selectionOf(p *project.Project, ids []string, contains func(commit string) (bool, error)) (*Selection, error) {
	byID := map[string]*project.Record{}
	for _, r := range p.Records {
		byID[r.ID] = r
	}
	order, reach, err := deps.Order(byID, ids)
	if err != nil {
		return nil, err
	}
	v, err := deps.Preview(p.Records, ids)
	if err != nil {
		return nil, err
	}
	s := &Selection{Selected: ids, Order: order, Notes: v.Notes}
	waits := map[string][]string{}
	for _, q := range v.Questions {
		for _, id := range q.Blocks {
			waits[id] = append(waits[id], fmt.Sprintf("blocked by open question %s (%s)", q.ID, q.Title))
		}
	}
	for _, it := range v.Items {
		if !it.Outside {
			continue
		}
		o := Outside{ID: it.ID, Status: it.Status, NeededBy: it.NeededBy}
		wait := ""
		switch {
		case it.Status == "":
			o.Delivery, wait = "not among the records read", "needs "+it.ID+", which is not in this checkout"
		case it.Status == "done" && it.Candidate == "":
			o.Delivery = "done without a candidate: delivery unrecorded"
		case it.Status == "done":
			// A candidate Git cannot read here is not in the base either.
			in, err := contains(it.Candidate)
			switch {
			case err != nil:
				o.Delivery = "done, but candidate " + short(it.Candidate) + " cannot be read here"
				wait = "needs " + it.ID + ", whose candidate " + short(it.Candidate) + " cannot be read here"
			case in:
				o.Delivery = "delivered: candidate " + short(it.Candidate) + " is in the base"
			default:
				o.Delivery = "done, but candidate " + short(it.Candidate) + " is not in the base"
				wait = "needs " + it.ID + ", whose candidate " + short(it.Candidate) + " the base lacks"
			}
		default:
			o.Delivery = "not delivered: " + it.Status
			wait = "needs " + it.ID + ", which is " + it.Status + " and not selected"
		}
		if wait != "" {
			for _, id := range it.NeededBy {
				waits[id] = append(waits[id], wait)
			}
		}
		s.Outside = append(s.Outside, o)
	}
	for _, id := range order {
		r := byID[id]
		if r.Status != "proposed" && r.Status != "active" {
			return nil, fmt.Errorf("%s is %s; only proposed or active work runs", id, r.Status)
		}
		for _, prior := range s.Members { // earlier in order, so already decided
			if reach[id][prior.ID] && prior.Wait != "" {
				waits[id] = append(waits[id], "needs "+prior.ID+", which waits")
			}
		}
		s.Members = append(s.Members, Member{ID: id, Path: r.Path, Revision: project.Revision(r.Source), Status: r.Status, Wait: strings.Join(waits[id], "; ")})
	}
	return s, nil
}

// startable refuses a selection none of whose members can start: the wait
// is already recorded, and launching again would only spend on it.
func (s *Selection) startable() error {
	var why []string
	for _, m := range s.Members {
		if m.Wait == "" {
			return nil
		}
		why = append(why, m.ID+" "+m.Wait)
	}
	return fmt.Errorf("nothing in the selection can start: %s; resolve that before another attempt", strings.Join(why, "; "))
}

// digest is sha256 over what the launch would run: the IDs as given, each
// member's revision in order, the base, and the resolved options.
func digest(l *Launch) string {
	h := sha256.New()
	fmt.Fprintf(h, "selected %q\n", l.Selection.Selected)
	for _, m := range l.Selection.Members {
		fmt.Fprintf(h, "member %s %s\n", m.ID, m.Revision)
	}
	fmt.Fprintf(h, "base %s\nbranch %s\nworktree %s\nbudget %s\npermission %s\nmodel %s\neffort %s\nuntil %s\n",
		l.Base, l.Branch, l.Worktree, l.BudgetUSD, l.PermissionMode, l.Model, l.Effort, l.Until)
	return fmt.Sprintf("sha256:%x", h.Sum(nil))
}

// Preview checks req as Start would and returns the launch it would make,
// without writing, creating or starting anything: no attempt id, command or
// owner. A branch that exists without a worktree is read at its tip only
// for the base; its records are checked when Start checks it out.
func Preview(req Request) (*Launch, error) {
	pre, err := prepare(req)
	if err != nil {
		return nil, err
	}
	return pre.launch, nil
}

// Explain lists a selection as preview lines, one fact each; visible escapes
// each value that came from a file.
func Explain(l *Launch, visible func(string) string) []string {
	s := l.Selection
	var lines []string
	line := func(format string, args ...any) { lines = append(lines, fmt.Sprintf(format, args...)) }
	line("Selected: %s; order %s", strings.Join(s.Selected, " "), strings.Join(s.Order, ", "))
	for _, m := range s.Members {
		state := "can start"
		if m.Wait != "" {
			state = "waits: " + visible(m.Wait)
		}
		line("Member %s: %s at %s, record %s; %s", m.ID, m.Status, visible(m.Path), m.Revision, state)
	}
	for _, o := range s.Outside {
		line("Outside %s: %s; needed by %s; never added", o.ID, visible(o.Delivery), strings.Join(o.NeededBy, ", "))
	}
	for _, n := range s.Notes {
		line("Note: %s", visible(n))
	}
	line("Review boundary: %s", Boundary)
	line("Continuation: %s", Continuation)
	return lines
}

// contains answers whether base holds a commit, in root's repository.
func containsIn(root, base string) func(string) (bool, error) {
	ancestry := deps.Ancestry(context.Background(), root)
	return func(commit string) (bool, error) { return ancestry(commit, base) }
}
