// Package sweep acts on candidates in review under the owner's standing
// policy (G-180, G-182): it starts one resolution attempt for a candidate that
// conflicts with the target, and approves and integrates one that meets the
// policy's conditions once the merged result passed its verification. Each
// act is attributed to the policy's grove.yaml revision. Everything the
// policy does not name waits for the owner, and no policy means nothing
// happens. Nothing runs between sweeps: the owner or a scheduler runs one.
package sweep

import (
	"bytes"
	"cmp"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/mascah/grove/internal/attempt"
	"github.com/mascah/grove/internal/integrate"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
	"github.com/mascah/grove/internal/update"
	"github.com/mascah/grove/internal/versions"
)

// ClosingLine is how a review record says no finding, knowledge findings
// included, remains open on the commit it examined: the reviewer's last
// round ends with it and the author copies it into the record.
const ClosingLine = "Open findings: none"

// Acts an Item can plan.
const (
	Wait      = "wait"      // the owner judges it
	Skip      = "skip"      // the target already holds it
	Resolve   = "resolve"   // one resolution attempt
	Approve   = "approve"   // verify, then approve
	Integrate = "integrate" // verify, then approve and integrate
)

// Item is what the policy does to one candidate in review, and why.
type Item struct {
	ID, Branch, Candidate string
	Act, Why              string

	tip, checkout string // the branch's commit and its checkout's project directory
	merge         versions.Merge
	review        *project.Record // Approve: the review the policy consumes
	budget        string          // Resolve: the attempt's budget
}

// Sweep is one sweep's plan, read from the target's checkout.
type Sweep struct {
	Root, Target string
	Policy       *project.Policy
	Attribution  string // "policy grove.yaml sha256:…"
	Items        []Item
	prefix       string // the project's directory in the repository, "" or ending in /
	records      string // the record root, project-relative
}

// Plan reads every candidate in review on a branch other than the target and
// decides, without writing anything, what the policy does to each.
func Plan(root string) (*Sweep, error) {
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		return nil, fmt.Errorf("the project is not valid; fix it before sweeping:\n%s", ds[0].String())
	}
	if p.Target == "" {
		return nil, errors.New("sweep needs target: BRANCH in grove.yaml, the branch candidates are integrated into")
	}
	if branch, err := update.Branch(root); err != nil {
		return nil, err
	} else if branch != p.Target {
		return nil, fmt.Errorf("sweep runs in the checkout of the target %s, whose grove.yaml holds the policy; this one is on %s", p.Target, cmp.Or(branch, "no branch"))
	}
	if p.Policy == nil {
		return nil, errors.New("grove.yaml has no policy: nothing is automatic, and every candidate in review waits for the owner")
	}
	if status, err := repo.Git(root, "status", "--porcelain", "--", "grove.yaml"); err != nil {
		return nil, err
	} else if status != "" {
		return nil, errors.New("grove.yaml has uncommitted changes; commit the policy first, since every act names the revision it ran under")
	}
	res, err := versions.Inspect(p.Root, "")
	if err != nil {
		return nil, err
	}
	s := &Sweep{Root: p.Root, Target: p.Target, Policy: p.Policy, Attribution: "policy grove.yaml " + project.Revision(p.Config), prefix: res.Prefix, records: p.RecordDir}
	found := map[string][]*versions.Version{}
	var ids []string
	for _, g := range res.Groups {
		for i := range g.Versions {
			v := &g.Versions[i]
			if v.Source.Kind == "committed" && v.Source.Ref != "refs/heads/"+p.Target && v.Record != nil && v.Record.Type == "work" && v.Record.Status == "review" && v.Record.Candidate != "" {
				if found[g.ID] == nil {
					ids = append(ids, g.ID)
				}
				found[g.ID] = append(found[g.ID], v)
			}
		}
	}
	slices.Sort(ids)
	spent := 0.0
	for _, id := range ids {
		it := s.plan(res, p, found[id], &spent)
		s.Items = append(s.Items, it)
	}
	return s, nil
}

// plan decides one candidate. spent is the aggregate budget the resolutions
// planned before it take.
func (s *Sweep) plan(res *versions.Result, p *project.Project, vs []*versions.Version, spent *float64) Item {
	v := vs[0]
	r := v.Record
	it := Item{ID: r.ID, Branch: strings.TrimPrefix(v.Source.Ref, "refs/heads/"), Candidate: r.Candidate, tip: v.Source.Commit}
	wait := func(format string, args ...any) Item {
		it.Act, it.Why = Wait, fmt.Sprintf(format, args...)
		return it
	}
	if len(vs) > 1 {
		var names []string
		for _, o := range vs {
			names = append(names, strings.TrimPrefix(o.Source.Ref, "refs/heads/"))
		}
		it.Branch = strings.Join(names, ", ")
		return wait("in review on several branches")
	}
	records := branchRecords(res, v.Source)
	if group := update.Group(records, r); len(group) > 1 {
		var others []string
		for _, o := range group {
			if o != r {
				others = append(others, o.ID)
			}
		}
		return wait("shares its candidate with %s; a shared candidate waits for the owner", strings.Join(others, ", "))
	}
	ms, err := versions.PredictContext(context.Background(), s.Root, "refs/heads/"+s.Target, []string{it.tip})
	if err != nil {
		return wait("its merge into %s could not be predicted: %v", s.Target, err)
	}
	if it.merge = ms[0]; it.merge.Outcome == "integrated" {
		it.Act, it.Why = Skip, it.merge.Text(s.Target)
		return it
	}
	if r.Approved != "" {
		who := "the owner"
		if update.Delegated(r) {
			who = "the policy"
		}
		return wait("approved by %s; integrating it is the owner's: grove integrate %s", who, r.ID)
	}
	for _, q := range records {
		if q.Type == "question" && q.Status == "open" && slices.Contains(q.Blocks, r.ID) {
			return wait("open question %s blocks it", q.ID)
		}
	}
	for _, src := range res.Sources {
		if src.Kind == "live" && src.Ref == v.Source.Ref {
			it.checkout = filepath.Join(src.Worktree, filepath.FromSlash(res.Prefix))
		}
	}
	if it.checkout == "" {
		return wait("no checkout is on branch %s", it.Branch)
	}
	views, err := attempt.List(s.Root, r.ID)
	if err != nil {
		return wait("its attempts could not be read: %v", err)
	}
	for _, a := range views {
		if a.Status == attempt.Running || a.Status == attempt.Orphaned {
			return wait("attempt %s is %s", a.Launch.Attempt, a.Status)
		}
	}
	if it.merge.Outcome == "conflict" {
		return s.planResolve(it, r, p, spent)
	}
	return s.planApprove(it, r, records)
}

func (s *Sweep) planResolve(it Item, r *project.Record, p *project.Project, spent *float64) Item {
	conflict := it.merge.Text(s.Target)
	it.Act = Wait
	if !s.Policy.Resolve {
		it.Why = conflict + "; the policy does not resolve"
		return it
	}
	// Once per target movement: the resolution feedback names the target
	// commit in full.
	if strings.Contains(string(r.Source), "git merge "+it.merge.Target) {
		it.Why = fmt.Sprintf("%s again after a resolution of %s at %s was attempted; a repeated conflict waits for the owner", conflict, s.Target, short(it.merge.Target))
		return it
	}
	it.budget = cmp.Or(s.Policy.ResolveBudgetUSD, p.Run.BudgetUSD)
	if it.budget == "" {
		it.Why = conflict + "; no budget: set policy.resolve.budget or run.budget"
		return it
	}
	cost, _ := strconv.ParseFloat(it.budget, 64)
	limit, _ := strconv.ParseFloat(s.Policy.BudgetUSD, 64)
	if *spent+cost > limit {
		it.Why = fmt.Sprintf("%s; a %s USD attempt would pass the policy's %s USD for one sweep, %g USD of it already planned", conflict, it.budget, s.Policy.BudgetUSD, *spent)
		return it
	}
	*spent += cost
	it.Act, it.Why = Resolve, fmt.Sprintf("%s; one resolution attempt, budget %s USD", conflict, it.budget)
	return it
}

func (s *Sweep) planApprove(it Item, r *project.Record, records []*project.Record) Item {
	merge := it.merge.Text(s.Target)
	it.Act = Wait
	if !s.Policy.Approve {
		it.Why = merge + "; the policy does not approve"
		return it
	}
	if others, err := versions.Others(context.Background(), s.Root, r.Candidate, it.tip, r.Path); err != nil || len(others) != 0 {
		it.Why = fmt.Sprintf("commits after candidate %s change %s: the tip is a new candidate", short(r.Candidate), cmp.Or(strings.Join(others, ", "), fmt.Sprint(err)))
		return it
	}
	review, why := s.review(r, records)
	if review == nil {
		it.Why = why
		return it
	}
	it.review = review
	lines, why := s.scope(r)
	if why != "" {
		it.Why = why
		return it
	}
	it.Act = Approve
	then := "approve"
	if s.Policy.Integrate {
		it.Act, then = Integrate, "approve and integrate"
	}
	commands := "commands"
	if len(s.Policy.Verify) == 1 {
		commands = "command"
	}
	it.Why = fmt.Sprintf("%s; review %s has no open finding; %d changed lines, no never path; verify with %d %s on the merged result, then %s", merge, review.ID, lines, len(s.Policy.Verify), commands, then)
	return it
}

// review finds a current review of the work whose closing line says no
// finding is open, and that examined the candidate or an earlier commit from
// which only records changed.
func (s *Sweep) review(r *project.Record, records []*project.Record) (*project.Record, string) {
	why := fmt.Sprintf("no current review of candidate %s ends with %q", short(r.Candidate), ClosingLine)
	for _, o := range slices.Backward(records) {
		if o.Type != "review" || o.Status != "current" || !slices.Contains(o.Work, r.ID) || o.Examined == "" {
			continue
		}
		closing := ""
		for l := range strings.SplitSeq(string(o.Source), "\n") {
			if l = strings.TrimSpace(l); strings.HasPrefix(l, "Open findings:") {
				closing = l
			}
		}
		if closing != ClosingLine {
			why = fmt.Sprintf("review %s does not end with %q", o.ID, ClosingLine)
			continue
		}
		if update.SameCommit(o.Examined, r.Candidate) {
			return o, ""
		}
		changed, err := repo.Git(s.Root, "diff", "--name-only", "-z", "--no-relative", o.Examined, r.Candidate)
		if err != nil {
			why = fmt.Sprintf("review %s examined %s, which could not be compared with the candidate: %v", o.ID, short(o.Examined), err)
			continue
		}
		records := path.Join(s.prefix, s.records)
		if !slices.ContainsFunc(split(changed), func(f string) bool { return !strings.HasPrefix(f, records+"/") }) {
			return o, ""
		}
		why = fmt.Sprintf("review %s examined %s, and more than records changed before candidate %s", o.ID, short(o.Examined), short(r.Candidate))
	}
	return nil, why
}

// scope counts the lines the candidate changes against its merge base with
// the target, and refuses a never path, a path outside the project, a binary
// change, or more than max_lines.
func (s *Sweep) scope(r *project.Record) (int, string) {
	out, err := repo.Git(s.Root, "diff", "--numstat", "-z", "--no-relative", "--no-renames", "refs/heads/"+s.Target+"..."+r.Candidate)
	if err != nil {
		return 0, fmt.Sprintf("its changes could not be read: %v", err)
	}
	lines := 0
	for _, entry := range split(out) {
		fields := strings.SplitN(entry, "\t", 3)
		if len(fields) != 3 {
			return 0, "its changes could not be read: " + entry
		}
		file := fields[2]
		rel, ok := strings.CutPrefix(file, s.prefix)
		switch {
		case !ok:
			return 0, fmt.Sprintf("changes %s, outside the project", file)
		case rel == r.Path:
			continue
		}
		if pattern, ok := s.Policy.Matches(rel); ok {
			return 0, fmt.Sprintf("changes %s, which never matches (%s)", rel, pattern)
		}
		added, err1 := strconv.Atoi(fields[0])
		removed, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			return 0, fmt.Sprintf("changes the binary file %s", rel)
		}
		lines += added + removed
	}
	if s.Policy.MaxLines > 0 && lines > s.Policy.MaxLines {
		return 0, fmt.Sprintf("changes %d lines, over the policy's %d", lines, s.Policy.MaxLines)
	}
	return lines, ""
}

// Run acts on the plan, reporting one line per candidate and per fact. An
// act that fails leaves that candidate waiting with the reason; the sweep
// goes on to the next.
func (s *Sweep) Run(now time.Time, report func(string)) {
	for _, it := range s.Items {
		say := func(format string, args ...any) { report(it.ID + ": " + fmt.Sprintf(format, args...)) }
		switch it.Act {
		case Wait:
			say("waits: %s", it.Why)
		case Skip:
			say("skipped: %s", it.Why)
		case Resolve:
			// Resolve predicts the candidate again and refuses what no
			// longer conflicts.
			l, err := attempt.Resolve(attempt.Request{Root: s.Root, IDs: []string{it.ID}, BudgetUSD: it.budget, Policy: s.Attribution}, nil, now, func(f string) { say("%s", f) })
			if err != nil {
				say("waits: resolution refused: %v", err)
				continue
			}
			say("resolution attempt %s started under %s, budget %s USD", l.Attempt, s.Attribution, l.BudgetUSD)
		case Approve, Integrate:
			s.approve(it, now, say)
		}
	}
}

func (s *Sweep) approve(it Item, now time.Time, say func(string, ...any)) {
	if it.Act == Integrate {
		if dirty, err := repo.Git(s.Root, "status", "--porcelain", "--untracked-files=no"); err != nil || dirty != "" {
			say("waits: the checkout of %s has uncommitted changes, so it could not be integrated; nothing was verified or approved", s.Target)
			return
		}
	}
	if err := s.verify(it); err != nil {
		say("waits: verification of the merge with %s at %s failed: %v; nothing was approved and %s is unchanged", s.Target, short(it.merge.Target), err, s.Target)
		return
	}
	say("verified: the merge of %s with %s at %s passed %s", short(it.tip), s.Target, short(it.merge.Target), strings.Join(s.Policy.Verify, "; "))
	verdict := fmt.Sprintf("delegated under %s: review %s examined %s with no open finding; merged with %s at %s, verification passed (%s); %s",
		s.Attribution, it.review.ID, short(it.review.Examined), s.Target, short(it.merge.Target), strings.Join(s.Policy.Verify, "; "), s.produced(it))
	res, err := update.Approve(it.checkout, it.ID, verdict, now)
	if err != nil {
		say("waits: approval refused: %v", err)
		return
	}
	say("approved under %s, commit %s on %s", s.Attribution, short(res.Commit), it.Branch)
	if it.Act != Integrate {
		return
	}
	err = integrate.Run(integrate.Request{Root: s.Root, ID: it.ID, Expect: it.merge.Target, Policy: s.Attribution}, now, func(f string) { say("%s", f) })
	if err != nil {
		say("waits: approved, but integration refused: %v; the owner integrates it with grove integrate %s", err, it.ID)
	}
}

// verify merges the branch's tip into the target commit the prediction read,
// in a temporary worktree outside every checkout, and runs the policy's
// commands there, in order, until one fails.
func (s *Sweep) verify(it Item) error {
	dir, err := os.MkdirTemp("", "grove-sweep-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	wt := filepath.Join(dir, "merge")
	if _, err := repo.Git(s.Root, "worktree", "add", "-q", "--detach", wt, it.merge.Target); err != nil {
		return err
	}
	defer repo.Git(s.Root, "worktree", "remove", "--force", wt)
	// --no-commit leaves the merged tree in the worktree without an identity.
	if _, err := repo.Git(wt, "merge", "-q", "--no-commit", "--no-ff", it.tip); err != nil {
		return fmt.Errorf("the merge failed: %v", err)
	}
	drop := append([]string{}, repo.GitLocation...)
	var env []string
	for _, e := range os.Environ() {
		if name, _, _ := strings.Cut(e, "="); !slices.Contains(drop, name) {
			env = append(env, e)
		}
	}
	for _, c := range s.Policy.Verify {
		cmd := exec.Command("sh", "-c", c)
		cmd.Dir, cmd.Env = filepath.Join(wt, filepath.FromSlash(s.prefix)), env
		var out bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &out
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s: %v%s", c, err, tail(out.String()))
		}
	}
	return nil
}

// produced names the latest finished attempt whose result held the
// candidate, with its cost.
func (s *Sweep) produced(it Item) string {
	views, _ := attempt.List(s.Root, it.ID)
	for _, v := range views {
		if v.Result == nil {
			continue
		}
		states := []*attempt.State{v.Result.Record}
		for _, m := range v.Result.Members {
			states = append(states, m.Record)
		}
		if !slices.ContainsFunc(states, func(st *attempt.State) bool { return st != nil && update.SameCommit(st.Candidate, it.Candidate) }) {
			continue
		}
		if f := v.Result.Events.Result; f != nil {
			return fmt.Sprintf("attempt %s produced it for %.2f USD", v.Launch.Attempt, f.CostUSD)
		}
		return fmt.Sprintf("attempt %s produced it; its cost was not reported", v.Launch.Attempt)
	}
	return "no Grove attempt is recorded as producing it"
}

// branchRecords is every record the committed source holds.
func branchRecords(res *versions.Result, from *versions.Source) []*project.Record {
	var out []*project.Record
	for _, g := range res.Groups {
		for _, v := range g.Versions {
			if v.Source == from && v.Record != nil {
				out = append(out, v.Record)
			}
		}
	}
	return out
}

func split(z string) []string {
	var out []string
	for _, f := range strings.Split(z, "\x00") {
		if f = strings.TrimPrefix(f, "\n"); f != "" {
			out = append(out, f)
		}
	}
	return out
}

// tail is the last lines of a command's output, for a failure's reason.
func tail(out string) string {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return ""
	}
	return "; last output: " + strings.Join(lines[max(0, len(lines)-5):], " | ")
}

func short(commit string) string { return commit[:min(len(commit), 7)] }
