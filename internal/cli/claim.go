package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	Cycle    string    `json:"cycle"`
	Changes  []string  `json:"changes"` // the cycle's, in build order
	Commit   string    `json:"commit"`
	Host     string    `json:"host"`
	Checkout string    `json:"checkout"`
	Branch   string    `json:"branch"`
}

// claim records an attempt at cycle id, or, when id is empty, at the oldest
// eligible cycle, unless the latest attempt's changes are still open and
// retry is not set; then it creates the attempt's branch from the fetched
// target and prints the cycle's ID.
func claim(stdout io.Writer, id string, retry bool) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	if out, err := gitx.Run(p.Top, "status", "--porcelain"); err != nil {
		return err
	} else if out != "" {
		return errors.New("the checkout has uncommitted changes; commit, stash or discard them first")
	}
	local, err := gitx.ResolveBranch(p.Top, p.Target)
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
	if last := pausing(recs, ip.Agreed); last != nil && !retry {
		return fmt.Errorf("the project is paused after attempt %s, which %s started at %s on branch %s; once you have dealt with it, run grove claim --retry", last.ID, last.Host, last.Started.Format(time.RFC3339), last.Branch)
	}
	cy, ordered := eligible(p, ip)
	if id != "" {
		if cy = ip.Agreed.Cycle(id); cy == nil {
			return fmt.Errorf("%s at %s holds no cycle %s", p.Target, commit[:7], id)
		}
		if ordered, err = runnable(p, ip, cy); err != nil {
			return err
		}
	}
	if cy == nil {
		return fmt.Errorf("nothing to claim on %s at %s", p.Target, commit[:7])
	}

	now := time.Now()
	a := attempt{ID: doc.NewID(now), Started: now, Cycle: cy.ID, Commit: commit, Checkout: p.Top}
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
			return claim(stdout, id, retry) // another claim recorded first; decide again on its record
		}
		return err
	}
	if _, err := gitx.Run(p.Top, "switch", "-q", "-c", a.Branch, commit); err != nil {
		return err
	}
	// The implement instructions read the local target, so it moves too.
	if _, err := gitx.Run(p.Top, "update-ref", "refs/heads/"+p.Target, commit, local); err != nil {
		return err
	}
	fmt.Fprintln(stdout, cy.ID)
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

// pausing returns the latest of recs while a change it names is not
// archived in c, or else nil.
func pausing(recs []attempt, c *doc.Collection) *attempt {
	if len(recs) > 0 && slices.ContainsFunc(recs[0].Changes, func(id string) bool { return !c.InArchive(id) }) {
		return &recs[0]
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
// cannot build it: it is done, names a member in needs_human, waits on a
// cycle not done, has a member waiting on a change neither in it nor
// archived as finished, or has a member that is not as its plan read it.
func runnable(p project, ip instruct.Project, cy *doc.Cycle) ([]*doc.Change, error) {
	ordered, waiting := buildable(ip.Agreed, ip.Agreed, cy)
	if len(ordered)+len(waiting) == 0 {
		return nil, fmt.Errorf("cycle %s is done", cy.ID)
	}
	for _, id := range cy.Meta.Changes {
		if _, ok := cy.Meta.NeedsHuman[id]; ok {
			return nil, fmt.Errorf("%s needs a person, as cycle %s's needs_human says; build the cycle with /grove-implement %s", id, cy.ID, cy.ID)
		}
	}
	if len(cy.WaitsOn) > 0 {
		return nil, fmt.Errorf("cycle %s waits on cycle %s, which is not done", cy.ID, strings.Join(cy.WaitsOn, " and "))
	}
	if len(waiting) > 0 {
		return nil, fmt.Errorf("%s waits on %s, neither in cycle %s nor archived as finished", waiting[0].ID, strings.Join(waiting[0].WaitsOn, " and "), cy.ID)
	}
	changed, err := changedSincePlan(p, ip.Commit, ip.Agreed, cy)
	if err != nil {
		return nil, err
	}
	if len(changed) > 0 {
		return nil, fmt.Errorf("cycle %s's plan is stale, since %s changed on %s after it; re-plan it with /grove-plan %s", cy.ID, strings.Join(changed, " and "), p.Target, cy.ID)
	}
	return ordered, nil
}
