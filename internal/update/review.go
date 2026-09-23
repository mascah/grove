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

// Approve records the owner's verdict on the record's candidate: approved is
// set to the candidate and the verdict is appended to the body, committed
// alone. The checkout must hold the candidate with nothing but the record
// changed since it, and the record's file must match HEAD.
func Approve(root, id, verdict string, now time.Time) (Result, error) {
	if verdict = strings.TrimSpace(verdict); verdict == "" {
		return Result{}, errors.New("a verdict is required: the owner's words on the candidate")
	}
	r, err := inReview(root, id)
	if err != nil {
		return Result{}, err
	}
	if r.Approved != "" {
		return Result{}, fmt.Errorf("%s: candidate %s is already approved; integrate it", id, r.Candidate)
	}
	if err := holdsCandidate(root, r, true); err != nil {
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
// candidate stays, so the earlier reviews still compare to it.
func Feedback(root, id, text string, now time.Time) (Result, error) {
	if text = strings.TrimSpace(text); text == "" {
		return Result{}, errors.New("feedback text is required: what the next attempt must change")
	}
	r, err := inReview(root, id)
	if err != nil {
		return Result{}, err
	}
	if err := holdsCandidate(root, r, false); err != nil {
		return Result{}, err
	}
	req := Request{
		ID: id, Expect: project.Revision(r.Source), Commit: true,
		Set:    []Field{{"status", "active"}},
		Append: fmt.Sprintf("Feedback on candidate %s, %s: %s", short(r.Candidate), now.UTC().Format("2006-01-02"), text),
	}
	if r.Approved != "" {
		req.Unset = []string{"approved"}
	}
	return Apply(root, req, now, nil)
}

// inReview loads the project and returns the work record in review.
func inReview(root, id string) (*project.Record, error) {
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		return nil, fmt.Errorf("the project is not valid; fix it before judging a candidate:\n%s", diagnostics(ds))
	}
	i := slices.IndexFunc(p.Records, func(r *project.Record) bool { return r.ID == id })
	if i < 0 {
		return nil, fmt.Errorf("record %s not found in this project", id)
	}
	r := p.Records[i]
	if r.Type != "work" {
		return nil, fmt.Errorf("%s is a %s, not work: only work has a candidate", id, r.Type)
	}
	if r.Status != "review" {
		return nil, fmt.Errorf("%s is %s, not in review: there is no candidate awaiting judgment", id, r.Status)
	}
	return r, nil
}

// holdsCandidate checks that this checkout is the one to judge the record in:
// HEAD contains the candidate, the record's file matches HEAD (these commands
// commit that file, and never someone's uncommitted edit), and, when
// exact, nothing but the record changed since the candidate, since a later
// commit is a new candidate that needs its own judgment.
func holdsCandidate(root string, r *project.Record, exact bool) error {
	if _, err := repo.Git(root, "merge-base", "--is-ancestor", r.Candidate, "HEAD"); err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) && exit.ExitCode() == 1 {
			return fmt.Errorf("this checkout does not hold candidate %s; run this in the checkout of the branch that has %s in review", r.Candidate, r.ID)
		}
		return fmt.Errorf("candidate %s could not be checked against this checkout's HEAD: %v", r.Candidate, err)
	}
	status, err := repo.Git(root, "status", "--porcelain", "-z", "--", ":(literal)"+filepath.FromSlash(r.Path))
	if err != nil {
		return err
	}
	if status != "" {
		return fmt.Errorf("%s has uncommitted changes in this checkout; commit or discard them first, since judging commits that file", r.Path)
	}
	if !exact {
		return nil
	}
	others, tip, err := changedSince(root, r.Candidate, r.Path)
	if err != nil {
		return err
	}
	if len(others) != 0 {
		return fmt.Errorf("commits after candidate %s change %s: the tip %s is a new candidate; set candidate=%s, review it, then judge that", short(r.Candidate), strings.Join(others, ", "), short(tip), short(tip))
	}
	return nil
}

// changedSince lists the files other than the record's that differ between
// the candidate and HEAD, with HEAD itself.
func changedSince(root, candidate, recordPath string) (others []string, tip string, err error) {
	if others, err = versions.Others(context.Background(), root, candidate, "HEAD", recordPath); err != nil {
		return nil, "", err
	}
	tip, err = repo.Git(root, "rev-parse", "HEAD")
	return others, strings.TrimSpace(tip), err
}

func short(commit string) string { return commit[:min(len(commit), 7)] }
