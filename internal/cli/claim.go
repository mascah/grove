package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/mascah/grove/internal/doc"
	"github.com/mascah/grove/internal/gitx"
	"github.com/mascah/grove/internal/instruct"
)

// runs is the branch holding one commit per attempt, with an empty tree
// and the attempt's record as its message.
const runs = "refs/heads/grove/runs"

// attempt is the record of one claim.
type attempt struct {
	ID       string    `json:"id"`
	Started  time.Time `json:"started"`
	Cycle    string    `json:"cycle,omitempty"` // none for a loose change
	Changes  []string  `json:"changes"`         // the cycle's, in build order, or the loose change
	Commit   string    `json:"commit"`
	Host     string    `json:"host"`
	Checkout string    `json:"checkout"`
	Branch   string    `json:"branch"`
}

// claim records an attempt at piece id, a cycle or a loose change, unless it
// is in flight or overlaps a change in flight, or, when id is empty, at the
// oldest eligible cycle, unless any change is in flight; retry claims despite
// what is in flight. Then it creates the attempt's branch from the fetched
// target and prints the piece's ID.
func claim(stdout io.Writer, given string, retry bool) error {
	id := given
	p, err := openProject()
	if err != nil {
		return err
	}
	if out, err := gitx.Run(p.Top, "status", "--porcelain"); err != nil {
		return err
	} else if out != "" {
		return errors.New("the checkout has uncommitted changes; commit, stash or discard them first")
	}
	local, err := gitx.LocalBranch(p.Top, p.Target)
	if err != nil {
		return err
	}
	remotes, err := gitx.Run(p.Top, "remote")
	if err != nil {
		return err
	}
	origin := slices.Contains(strings.Fields(remotes), "origin")
	commit, tip, err := fetch(p, origin)
	if err != nil {
		return err
	}
	if _, err := gitx.Run(p.Top, "merge-base", "--is-ancestor", local, commit); err != nil {
		return fmt.Errorf("the local %s has commits origin's does not, so it cannot be fast-forwarded", p.Target)
	}
	config, err := gitx.Run(p.Top, "show", commit+":grove.yaml")
	if err != nil {
		return err
	}
	if p, err = parseProject(p.Top, []byte(config)); err != nil {
		return err
	}
	ip, err := p.at(commit)
	if err != nil {
		return err
	}
	recs, err := records(p.Top, tip)
	if err != nil {
		return err
	}
	flight := inFlight(recs, ip.Agreed)
	var cy *doc.Cycle
	var ordered []*doc.Change
	if id == "" {
		if last := newest(recs, flight); last != nil && !retry {
			return fmt.Errorf("the project is paused after attempt %s, which %s started at %s on branch %s; once you have dealt with it, run grove claim --retry", last.ID, last.Host, last.Started.Format(time.RFC3339), last.Branch)
		}
		if cy, ordered = eligible(p, ip); cy == nil {
			return fmt.Errorf("nothing to claim on %s at %s", p.Target, commit[:7])
		}
		id = cy.ID
	} else if cy, ordered, err = piece(p, ip, id); err != nil {
		return err
	} else if err = busy(p, ip.Agreed, id, ordered, flight); err != nil && !retry {
		return err
	}

	now := time.Now()
	a := attempt{ID: doc.NewID(now), Started: now, Commit: commit, Checkout: p.Top}
	if cy != nil {
		a.Cycle = cy.ID
	}
	if a.Host, err = os.Hostname(); err != nil {
		return err
	}
	a.Branch = "grove/" + a.ID
	for _, ch := range ordered {
		a.Changes = append(a.Changes, ch.ID)
	}
	msg, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}
	tree, err := gitx.Run(p.Top, "hash-object", "-t", "tree", "-w", os.DevNull)
	if err != nil {
		return err
	}
	args := []string{"commit-tree", strings.TrimSpace(tree), "-m", string(msg)}
	if tip != "" {
		args = append(args, "-p", tip)
	}
	rec, err := gitx.Run(p.Top, args...)
	if err != nil {
		return err
	}
	// Each update succeeds only while grove/runs is at tip, so of two claims one wins.
	if origin {
		_, err = gitx.Run(p.Top, "push", "-q", "origin", strings.TrimSpace(rec)+":"+runs, "--force-with-lease="+runs+":"+tip)
	} else {
		_, err = gitx.Run(p.Top, "update-ref", runs, strings.TrimSpace(rec), tip)
	}
	if err != nil {
		if _, moved, ferr := fetch(p, origin); ferr == nil && moved != tip {
			return claim(stdout, given, retry) // another claim recorded first; decide again on its record
		}
		return fmt.Errorf("the record was refused, so nothing is recorded: %w", err)
	}
	if _, err := gitx.Run(p.Top, "switch", "-q", "-c", a.Branch, commit); err != nil {
		return err
	}
	// The implement instructions read the local target, so it moves too.
	if _, err := gitx.Run(p.Top, "update-ref", "refs/heads/"+p.Target, commit, local); err != nil {
		return err
	}
	fmt.Fprintln(stdout, id)
	return nil
}

// fetch updates the target and grove/runs from origin, when the checkout
// has it, and returns the fetched target's commit and the tip of grove/runs.
func fetch(p project, origin bool) (commit, tip string, err error) {
	target, recs := "refs/heads/"+p.Target, runs
	if origin {
		target, recs = "refs/remotes/origin/"+p.Target, "refs/remotes/origin/grove/runs"
		// A pattern, so that origin having no grove/runs yet is no error.
		if _, err := gitx.Run(p.Top, "fetch", "-q", "origin", "+refs/heads/"+p.Target+":"+target, "+"+runs+"*:"+recs+"*"); err != nil {
			return "", "", err
		}
	}
	if commit, err = gitx.Run(p.Top, "rev-parse", "--verify", target+"^{commit}"); err != nil {
		return "", "", err
	}
	return strings.TrimSpace(commit), tipOf(p.Top, recs), nil
}

// tipOf returns the commit at ref, or "" when there is no such ref.
func tipOf(dir, ref string) string {
	out, _ := gitx.Run(dir, "rev-parse", "-q", "--verify", ref) // fails only when there is none
	return strings.TrimSpace(out)
}

// inFlight returns each change that a record of recs, newest first, names
// and c does not hold archived, with its attempt, the newest record naming it.
func inFlight(recs []attempt, c *doc.Collection) map[string]*attempt {
	flight := map[string]*attempt{}
	for i := range recs {
		for _, id := range recs[i].Changes {
			if _, ok := flight[id]; !ok && !c.InArchive(id) {
				flight[id] = &recs[i]
			}
		}
	}
	return flight
}

// newest returns the newest of recs that is the attempt of a change in
// flight, or nil.
func newest(recs []attempt, flight map[string]*attempt) *attempt {
	for i := range recs {
		if slices.Contains(slices.Collect(maps.Values(flight)), &recs[i]) {
			return &recs[i]
		}
	}
	return nil
}

// piece returns the cycle id names, if it names one, and the changes grove
// claim given id would record in build order, or why it cannot build them.
func piece(p project, ip instruct.Project, id string) (*doc.Cycle, []*doc.Change, error) {
	c := ip.Agreed
	if cy := c.Cycle(id); cy != nil {
		ordered, err := runnable(p, ip, cy)
		return cy, ordered, err
	}
	ch := c.Change(id)
	if ch == nil && c.InArchive(id) {
		return nil, nil, fmt.Errorf("change %s is archived", id)
	} else if ch == nil {
		return nil, nil, fmt.Errorf("%s at %s holds no cycle or change %s", p.Target, ip.Commit[:7], id)
	}
	for _, cy := range c.Cycles {
		if !cy.Done && slices.Contains(cy.Meta.Changes, id) {
			return nil, nil, fmt.Errorf("%s is in cycle %s; claim that with grove claim %s", id, cy.ID, cy.ID)
		}
	}
	if len(ch.WaitsOn) > 0 {
		return nil, nil, fmt.Errorf("%s waits on %s, which is not archived as finished", id, strings.Join(ch.WaitsOn, " and "))
	}
	return nil, []*doc.Change{ch}, nil
}

// busy says why grove claim given piece id, whose changes are those, records
// nothing while flight holds what it holds: a change of the piece is in
// flight, or plans a path a change in flight plans.
func busy(p project, c *doc.Collection, id string, changes []*doc.Change, flight map[string]*attempt) error {
	for _, ch := range changes {
		if a := flight[ch.ID]; a != nil {
			return fmt.Errorf("%s is in flight in attempt %s, which %s started at %s on branch %s", id, a.ID, a.Host, a.Started.Format(time.RFC3339), a.Branch)
		}
	}
	for _, f := range slices.Sorted(maps.Keys(flight)) {
		for _, ch := range changes {
			if other := c.Change(f); other != nil {
				if shared := doc.Overlap(ch.Planned(p.Root), other.Planned(p.Root)); shared != "" {
					return fmt.Errorf("%s shares %s with %s, in flight in attempt %s", id, shared, f, flight[f].ID)
				}
			}
		}
	}
	return nil
}

// records returns the attempts on grove/runs from tip back, newest first.
func records(dir, tip string) ([]attempt, error) {
	if tip == "" {
		return nil, nil
	}
	out, err := gitx.Run(dir, "log", "--format=%B%x00", tip)
	if err != nil {
		return nil, err
	}
	var all []attempt
	for msg := range strings.SplitSeq(out, "\x00") {
		if msg = strings.TrimSpace(msg); msg == "" {
			continue
		}
		var a attempt
		if err := json.Unmarshal([]byte(msg), &a); err != nil {
			return all, fmt.Errorf("a record on grove/runs is not JSON: %w", err)
		}
		all = append(all, a)
	}
	return all, nil
}

// eligible returns the oldest cycle the target ip holds that is marked
// unattended and runnable, and its members in build order.
func eligible(p project, ip instruct.Project) (*doc.Cycle, []*doc.Change) {
	for _, cy := range ip.Agreed.Cycles {
		if !cy.Meta.Unattended {
			continue
		}
		// A planned that is no commit is check's to report; the cycle is skipped.
		if ordered, err := runnable(p, ip, cy); err == nil {
			return cy, ordered
		}
	}
	return nil, nil
}

// runnable returns cycle cy's members in build order, or why grove claim
// cannot build it: it is done, waits on a cycle not done, has a member waiting on a change neither in it nor
// archived as finished, or has a member that is not as its plan read it.
func runnable(p project, ip instruct.Project, cy *doc.Cycle) ([]*doc.Change, error) {
	ordered, waiting := buildable(ip.Agreed, ip.Agreed, cy)
	if len(ordered)+len(waiting) == 0 {
		return nil, fmt.Errorf("cycle %s is done", cy.ID)
	}
	if len(cy.WaitsOn) > 0 {
		return nil, fmt.Errorf("cycle %s waits on cycle %s, which is not done", cy.ID, strings.Join(cy.WaitsOn, " and "))
	}
	if len(waiting) > 0 {
		return nil, fmt.Errorf("%s waits on %s, neither in cycle %s nor archived as finished", waiting[0].ID, strings.Join(waiting[0].WaitsOn, " and "), cy.ID)
	}
	changed, err := changedSincePlan(p, ip.Commit, ip.Agreed, cy)
	if err != nil {
		return nil, absent(p, cy.Path, "planned", cy.Meta.Planned)
	}
	if len(changed) > 0 {
		return nil, fmt.Errorf("cycle %s's plan is stale, since %s changed on %s after it; re-plan it with /grove-plan %s", cy.ID, strings.Join(changed, " and "), p.Target, cy.ID)
	}
	return ordered, nil
}
