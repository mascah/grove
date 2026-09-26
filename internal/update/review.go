package update

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
	"github.com/mascah/grove/internal/versions"
)

// The owner's two dispositions of a candidate in review (G-044), each an
// ordinary update of the record committed alone in the checkout of the branch
// that holds it: approval binds the candidate through the approved field and
// quotes the verdict in the body; feedback reopens the work with the text in
// the body and no approval left behind. Neither touches any other file, ref,
// or worktree.
//
// A candidate can be shared (G-188): when an explicitly selected set of work
// is handed off together, the records on the branch whose candidate is the
// same commit are one group. Approval stays per record, and ignores the
// other members' own record commits; feedback on any member reopens the
// group, since the next candidate replaces the shared one.

// Approve records the owner's verdict on the record's candidate: approved is
// set to the candidate and the verdict is appended to the body, committed
// alone. The checkout must hold the candidate with nothing but the record
// changed since it, and the record's file must match HEAD.
func Approve(root, id, verdict string, now time.Time) (Result, error) {
	if verdict = strings.TrimSpace(verdict); verdict == "" {
		return Result{}, errors.New("a verdict is required: the owner's words on the candidate")
	}
	p, r, err := inReview(root, id)
	if err != nil {
		return Result{}, err
	}
	if r.Approved != "" {
		return Result{}, fmt.Errorf("%s: candidate %s is already approved; integrate it", id, r.Candidate)
	}
	if err := holdsCandidate(root, r, Group(p.Records, r)); err != nil {
		return Result{}, err
	}
	return Apply(root, Request{
		ID: id, Expect: project.Revision(r.Source), Commit: true,
		Set:    []Field{{"approved", r.Candidate}},
		Append: fmt.Sprintf("Verdict on candidate %s, %s: %s", short(r.Candidate), now.UTC().Format("2006-01-02"), verdict),
	}, now, nil)
}

// Feedback returns the record to active with the owner's text appended to the
// body, committed alone; an approval is removed in the same update and the
// candidate stays, so the earlier reviews still compare to it. Every other
// member of its group in review is returned to active the same way, with a
// line naming this feedback instead of the text, each committed alone and
// listed in Reopened.
func Feedback(root, id, text string, now time.Time) (Result, error) {
	if text = strings.TrimSpace(text); text == "" {
		return Result{}, errors.New("feedback text is required: what the next attempt must change")
	}
	p, r, err := inReview(root, id)
	if err != nil {
		return Result{}, err
	}
	var others []*project.Record
	for _, o := range Group(p.Records, r) {
		if o != r && o.Status == "review" {
			others = append(others, o)
		}
	}
	// Every file is checked before the first is written.
	for _, o := range append([]*project.Record{r}, others...) {
		if err := clean(root, o); err != nil {
			return Result{}, err
		}
	}
	if err := holdsCandidate(root, r, nil); err != nil {
		return Result{}, err
	}
	reopen := func(o *project.Record, line string) (Result, error) {
		req := Request{ID: o.ID, Expect: project.Revision(o.Source), Commit: true, Set: []Field{{"status", "active"}}, Append: line}
		if o.Approved != "" {
			req.Unset = []string{"approved"}
		}
		return Apply(root, req, now, nil)
	}
	day := now.UTC().Format("2006-01-02")
	res, err := reopen(r, fmt.Sprintf("Feedback on candidate %s, %s: %s", short(r.Candidate), day, text))
	if err != nil {
		return res, err
	}
	for _, o := range others {
		done, err := reopen(o, fmt.Sprintf("Reopened with %s's feedback on candidate %s, %s: the next candidate replaces the one this group shared.", id, short(r.Candidate), day))
		if err != nil {
			return res, fmt.Errorf("%s is active again at commit %s, but reopening %s, which shares its candidate, failed: %v; run grove feedback %s with the same text", id, short(res.Commit), o.ID, err, o.ID)
		}
		res.Reopened = append(res.Reopened, done)
	}
	return res, nil
}

// Group is the work in records whose candidate is r's commit, r included,
// ordered by ID (G-188): the members one shared candidate was handed off
// for. A record without a candidate is a group of none.
func Group(records []*project.Record, r *project.Record) []*project.Record {
	var group []*project.Record
	for _, o := range records {
		if o.Type == "work" && r.Candidate != "" && SameCommit(o.Candidate, r.Candidate) {
			group = append(group, o)
		}
	}
	slices.SortFunc(group, func(a, b *project.Record) int { return strings.Compare(a.ID, b.ID) })
	return group
}

// SameCommit compares two recorded commits, either of which may be
// abbreviated.
func SameCommit(a, b string) bool {
	return a != "" && b != "" && (strings.HasPrefix(a, b) || strings.HasPrefix(b, a))
}

// inReview loads the project and returns the work record in review.
func inReview(root, id string) (*project.Project, *project.Record, error) {
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		return nil, nil, fmt.Errorf("the project is not valid; fix it before judging a candidate:\n%s", diagnostics(ds))
	}
	i := slices.IndexFunc(p.Records, func(r *project.Record) bool { return r.ID == id })
	if i < 0 {
		return nil, nil, fmt.Errorf("record %s not found in this project", id)
	}
	r := p.Records[i]
	if r.Type != "work" {
		return nil, nil, fmt.Errorf("%s is a %s, not work: only work has a candidate", id, r.Type)
	}
	if r.Status != "review" {
		return nil, nil, fmt.Errorf("%s is %s, not in review: there is no candidate awaiting judgment", id, r.Status)
	}
	return p, r, nil
}

// holdsCandidate checks that this checkout is the one to judge the record in:
// HEAD contains the candidate, the record's file matches HEAD (these commands
// commit that file, and never someone's uncommitted edit), and, when group
// names the records sharing the candidate, nothing but their files changed
// since it, since a later commit is a new candidate that needs its own
// judgment.
func holdsCandidate(root string, r *project.Record, group []*project.Record) error {
	if _, err := repo.Git(root, "merge-base", "--is-ancestor", r.Candidate, "HEAD"); err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) && exit.ExitCode() == 1 {
			return fmt.Errorf("this checkout does not hold candidate %s; run this in the checkout of the branch that has %s in review", r.Candidate, r.ID)
		}
		return fmt.Errorf("candidate %s could not be checked against this checkout's HEAD: %v", r.Candidate, err)
	}
	if err := clean(root, r); err != nil {
		return err
	}
	if group == nil {
		return nil
	}
	var paths []string
	for _, o := range group {
		paths = append(paths, o.Path)
	}
	others, err := versions.Others(context.Background(), root, r.Candidate, "HEAD", paths...)
	if err != nil {
		return err
	}
	if len(others) != 0 {
		tip, err := repo.Git(root, "rev-parse", "HEAD")
		if err != nil {
			return err
		}
		tip = strings.TrimSpace(tip)
		return fmt.Errorf("commits after candidate %s change %s: the tip %s is a new candidate; set candidate=%s, review it, then judge that", short(r.Candidate), strings.Join(others, ", "), short(tip), short(tip))
	}
	return nil
}

// clean refuses a record whose file differs from HEAD in this checkout.
func clean(root string, r *project.Record) error {
	status, err := repo.Git(root, "status", "--porcelain", "-z", "--", ":(literal)"+filepath.FromSlash(r.Path))
	if err != nil {
		return err
	}
	if status != "" {
		return fmt.Errorf("%s has uncommitted changes in this checkout; commit or discard them first, since judging commits that file", r.Path)
	}
	return nil
}

func short(commit string) string { return commit[:min(len(commit), 7)] }

// Delegated reports whether the latest verdict on r's candidate was given
// under a standing policy (G-182) rather than by the owner: the sweep writes
// its verdicts beginning "delegated under policy".
func Delegated(r *project.Record) bool {
	prefix := "Verdict on candidate " + short(r.Candidate) + ", "
	delegated := false
	for l := range strings.SplitSeq(string(r.Source), "\n") {
		if strings.HasPrefix(l, prefix) {
			_, verdict, _ := strings.Cut(l, ": ")
			delegated = strings.HasPrefix(verdict, "delegated under policy ")
		}
	}
	return delegated
}
