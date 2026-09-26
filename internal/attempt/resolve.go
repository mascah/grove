package attempt

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/update"
	"github.com/mascah/grove/internal/versions"
)

// Resolve updates a candidate in review that conflicts with the target
// through one bounded attempt (G-178): it records feedback naming the target
// commit and the conflicting files, which is the attempt's whole mandate,
// then starts one attempt on the candidate's branch in its checkout, with the
// launch defaults. req names one work and any launch flags; where it runs is
// Resolve's choice. shown, when not nil, is the fact the caller showed, and a
// candidate or target commit that differs from it now is refused. Every
// refusal Start would make is asked before the feedback is written too; a
// launch that still fails after it says the feedback stands and how to
// launch.
func Resolve(req Request, shown *versions.Merge, now time.Time, report func(string)) (*Launch, error) {
	switch {
	case len(req.IDs) != 1:
		return nil, errors.New("resolve takes one work ID; the records sharing its candidate go with it")
	case req.Until != "" || req.Branch != "" || req.Worktree != "":
		return nil, errors.New("resolve runs on the candidate's branch in its checkout, through to the handoff; --until, --branch and --worktree do not apply")
	}
	id := req.IDs[0]
	p, ds := project.Load(req.Root, req.Root)
	if len(ds) != 0 {
		return nil, fmt.Errorf("the project is not valid; fix it before resolving:\n%s", ds[0].String())
	}
	if p.Target == "" {
		return nil, errors.New("resolve needs target: BRANCH in grove.yaml, the branch the candidate conflicts with")
	}
	res, err := versions.Inspect(p.Root, "")
	if err != nil {
		return nil, err
	}
	from, r, err := inReview(res, id, p.Target)
	if err != nil {
		return nil, err
	}
	branch := strings.TrimPrefix(from.Ref, "refs/heads/")
	var checkout *versions.Source
	for _, s := range res.Sources {
		if s.Kind == "live" && s.Ref == from.Ref {
			checkout = s
		}
	}
	if checkout == nil {
		return nil, fmt.Errorf("no checkout is on branch %s, where the feedback is committed and the attempt runs; git worktree add one", branch)
	}
	dir := filepath.Join(checkout.Worktree, filepath.FromSlash(res.Prefix))

	ms, err := versions.PredictContext(context.Background(), p.Root, "refs/heads/"+p.Target, []string{r.Candidate})
	if err != nil {
		return nil, fmt.Errorf("the merge of candidate %s into %s could not be predicted: %v", short(r.Candidate), p.Target, err)
	}
	m := ms[0]
	switch {
	case shown != nil && !update.SameCommit(shown.Commit, r.Candidate):
		return nil, fmt.Errorf("%s's candidate is %s now, not %s, whose conflict was shown; look again", id, short(r.Candidate), short(shown.Commit))
	case shown != nil && shown.Target != m.Target:
		return nil, fmt.Errorf("%s moved from %s to %s since the conflict was shown, and candidate %s now %s; look again", p.Target, short(shown.Target), short(m.Target), short(r.Candidate), m.Text(p.Target))
	case m.Outcome != "conflict":
		return nil, fmt.Errorf("candidate %s of %s %s: there is no conflict to resolve", short(r.Candidate), id, m.Text(p.Target))
	}

	// The group reopens with the feedback (G-188), so it runs together.
	ids := []string{id}
	for _, o := range update.Group(branchRecords(res, from), r) {
		if o.ID != id && o.Status == "review" {
			ids = append(ids, o.ID)
		}
	}
	for _, w := range ids {
		views, err := List(p.Root, w)
		if err != nil {
			return nil, err
		}
		for _, v := range views {
			if v.Status == Running || v.Status == Orphaned {
				return nil, fmt.Errorf("attempt %s of %s is %s; stop it or wait for its result before resolving", v.Launch.Attempt, w, v.Status)
			}
		}
	}
	d := Defaulted(req, checkout.Run)
	if d.BudgetUSD == "" || d.PermissionMode == "" {
		return nil, ErrUnsupplied
	}
	// What Start would refuse once the feedback reopened the group is asked
	// now, of the branch's checkout with those records active: a wait, the
	// entrypoints, the provider.
	wp, wds := project.Load(dir, dir)
	if len(wds) != 0 {
		return nil, fmt.Errorf("the project on %s at %s is not valid: %s", branch, checkout.Worktree, wds[0].String())
	}
	for _, o := range wp.Records {
		if slices.Contains(ids, o.ID) {
			o.Status = "active"
		}
	}
	s, err := selectionOf(wp, ids, "", containsIn(dir, checkout.Commit))
	if err != nil {
		return nil, err
	}
	if err := s.startable(); err != nil {
		return nil, err
	}
	if err := entrypoints(checkout.Worktree, filepath.FromSlash(res.Prefix), branch); err != nil {
		return nil, err
	}
	if _, _, _, err := provider(); err != nil {
		return nil, err
	}

	text := Mandate(&m, p.Target)
	if req.Policy != "" {
		text = fmt.Sprintf("delegated under %s, budget %s USD: %s", req.Policy, d.BudgetUSD, text)
	}
	fb, err := update.Feedback(dir, id, text, now)
	if err != nil {
		return nil, err
	}
	report(fmt.Sprintf("feedback: %s is active again on branch %s, commit %s, to merge %s at %s and resolve the conflict %s", id, branch, short(fb.Commit), p.Target, short(m.Target), m.Where()))
	for _, o := range fb.Reopened {
		report(fmt.Sprintf("reopened: %s shared the candidate and is active again, commit %s", o.ID, short(o.Commit)))
	}
	run := req
	run.Root, run.IDs, run.Branch, run.Worktree = dir, ids, branch, checkout.Worktree
	l, err := Start(run, now, report)
	if err != nil {
		return nil, fmt.Errorf("the feedback stands, but the attempt did not start: %v; launch it with grove run %s --branch %s --worktree %s", err, strings.Join(ids, " "), branch, checkout.Worktree)
	}
	return l, nil
}

// Mandate is the feedback a resolution attempt works from, and its whole
// assignment: merge the named target commit, resolve the named files,
// verify, and hand off, or stop at a choice the record does not settle.
func Mandate(m *versions.Merge, target string) string {
	return fmt.Sprintf("%s. Resolve only that (grove resolve): in this branch, git merge %s, that commit of %s even if %s has moved since, never a rebase; "+
		"resolve those files keeping both sides' intent; rerun the repository's verification; and hand off the merge as the new candidate, "+
		"with the previous candidate %s, the merged commit and the resolved files in Evidence. Change nothing else. "+
		"If a resolution needs a choice this record does not settle, stop with a checkpoint naming it.",
		m.Text(target), m.Target, target, target, short(m.Commit))
}

// inReview finds the one committed branch other than the target holding the
// record in review with a candidate.
func inReview(res *versions.Result, id, target string) (*versions.Source, *project.Record, error) {
	var found []*versions.Version
	for _, g := range res.Groups {
		for i := range g.Versions {
			v := &g.Versions[i]
			if g.ID == id && v.Source.Kind == "committed" && v.Source.Ref != "refs/heads/"+target && v.Record != nil && v.Record.Type == "work" && v.Record.Status == "review" && v.Record.Candidate != "" {
				found = append(found, v)
			}
		}
	}
	switch len(found) {
	case 0:
		return nil, nil, fmt.Errorf("no branch holds %s in review with a candidate; nothing to resolve", id)
	case 1:
		return found[0].Source, found[0].Record, nil
	}
	var names []string
	for _, v := range found {
		names = append(names, strings.TrimPrefix(v.Source.Ref, "refs/heads/"))
	}
	return nil, nil, fmt.Errorf("%s is in review on several branches (%s); resolve needs one", id, strings.Join(names, ", "))
}

// branchRecords is every record the committed source holds.
func branchRecords(res *versions.Result, from *versions.Source) []*project.Record {
	var out []*project.Record
	for _, g := range res.Groups {
		for _, v := range g.Versions {
			if v.Source == from && v.Record != nil && !slices.Contains(out, v.Record) {
				out = append(out, v.Record)
			}
		}
	}
	return out
}
