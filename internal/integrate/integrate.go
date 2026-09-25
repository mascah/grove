// Package integrate merges an approved candidate into the configured target
// and marks its work done there, as a sequence of separately reported facts
// (G-044): approval found, merge made or refused, done written, cleanup done
// or kept. Every refusal happens before anything changes, and nothing after
// the merge undoes it.
package integrate

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
	"github.com/mascah/grove/internal/update"
	"github.com/mascah/grove/internal/versions"
)

// Request names the work to integrate from Root, the target's checkout. Cwd
// is the process's directory, which cleanup keeps out of any removed
// worktree.
type Request struct {
	Root, ID, Cwd string
	Cleanup       bool
}

// Run integrates the work, reporting each fact to report as it holds. The
// error is the refusal or failure that stopped the sequence; facts already
// reported stand.
func Run(req Request, now time.Time, report func(fact string)) error {
	root := req.Root
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		return fmt.Errorf("the project is not valid; fix it before integrating:\n%s", diagnostics(ds))
	}
	if p.Target == "" {
		return errors.New("integration needs target: BRANCH in grove.yaml, the branch approved work is merged into")
	}
	branch, err := update.Branch(root)
	if err != nil {
		return fmt.Errorf("this checkout's branch could not be read: %v", err)
	}
	if branch != p.Target {
		return fmt.Errorf("integrate runs in the checkout of the target %s; this one is on %s", p.Target, orNoBranch(branch))
	}
	if dirty, err := repo.Git(root, "status", "--porcelain", "--untracked-files=no"); err != nil {
		return err
	} else if dirty != "" {
		return fmt.Errorf("the checkout of %s has uncommitted changes; commit or set them aside before merging", p.Target)
	}
	res, err := versions.Inspect(root, req.ID)
	if err != nil {
		return err
	}
	from, r, err := approved(res, req.ID, p.Target)
	if err != nil {
		return err
	}
	name := strings.TrimPrefix(from.Ref, "refs/heads/")
	if _, err := repo.Git(root, "merge-base", "--is-ancestor", r.Candidate, from.Commit); err != nil {
		return fmt.Errorf("branch %s does not contain candidate %s, which it names; repair the record before integrating", name, r.Candidate)
	}
	if others, err := versions.Others(context.Background(), root, r.Candidate, from.Commit, r.Path); err != nil {
		return err
	} else if len(others) != 0 {
		return fmt.Errorf("commits after candidate %s on %s change %s: the tip %s is a new candidate; approve it before integrating", short(r.Candidate), name, strings.Join(others, ", "), short(from.Commit))
	}
	report(fmt.Sprintf("approval: candidate %s of %s approved on branch %s at %s%s", short(r.Candidate), req.ID, name, short(from.Commit), verdict(r)))

	before, err := head(root)
	if err != nil {
		return err
	}
	// A conflict is refused before the merge starts (G-177): merge-tree
	// performs it in objects only, against the commit that would be merged
	// into, and names the files. A prediction that fails, as on a Git
	// before 2.38, leaves the refusal to the merge below.
	if ms, err := versions.PredictContext(context.Background(), root, before, []string{from.Commit}); err == nil && ms[0].Outcome == "conflict" {
		conflict := ms[0].Text(p.Target)
		where := worktreeOf(res, from.Ref)
		if where == "" {
			where = "a checkout of " + name
		}
		return fmt.Errorf("merge of %s into %s refused: it %s; nothing was merged, %s is unchanged at %s and %s stays in review. Next: in %s, git merge %s, resolve the conflicts and commit, then hand that commit to review as the new candidate; or there, grove feedback %s %q returns it to an implementer",
			name, p.Target, conflict, p.Target, short(before), req.ID, where, p.Target, req.ID, conflict+"; merge "+p.Target+" and resolve")
	}
	// The commit the checks above read is what is merged, not the name: the
	// branch may move meanwhile, and a tag of the same name would win the
	// name. Git prints its CONFLICT lines on stdout, so both streams are read.
	if out, err := repo.Command(context.Background(), root, "merge", "--no-edit", "-m", "Merge branch '"+name+"'", from.Commit).CombinedOutput(); err != nil {
		if _, aborted := repo.Git(root, "rev-parse", "-q", "--verify", "MERGE_HEAD"); aborted == nil {
			if _, err := repo.Git(root, "merge", "--abort"); err != nil {
				return fmt.Errorf("merge of %s into %s failed and could not be aborted: %v; resolve it by hand", name, p.Target, err)
			}
		}
		why := strings.TrimSpace(string(out))
		if i := strings.Index(why, "CONFLICT"); i >= 0 {
			why = strings.ReplaceAll(why[i:], "\n", "; ")
		}
		return fmt.Errorf("merge of %s into %s refused: %s; %s is unchanged at %s and %s stays in review", name, p.Target, why, p.Target, short(before), req.ID)
	}
	after, err := head(root)
	if err != nil {
		return err
	}
	switch {
	case after == before:
		report(fmt.Sprintf("merge: nothing to merge; %s is already in %s at %s", name, p.Target, short(before)))
	case after == from.Commit:
		report(fmt.Sprintf("merge: fast-forward %s from %s to %s", p.Target, short(before), short(after)))
	default:
		report(fmt.Sprintf("merge: merge commit %s on %s (was %s)", short(after), p.Target, short(before)))
	}

	done, err := update.Apply(root, update.Request{ID: req.ID, Set: []update.Field{{Name: "status", Value: "done"}}, Commit: true}, now, nil)
	if err != nil {
		// update says whether the file was written before the commit failed.
		return fmt.Errorf("merged as %s, but marking %s done failed: %v; once that is repaired, commit the staged record here: git commit -m 'docs(%s): set status=done' -- %s (or, if it was not written, grove update %s --set status=done --commit)", short(after), req.ID, err, req.ID, r.Path, req.ID)
	}
	if done.Changed {
		report(fmt.Sprintf("done: %s done at commit %s", req.ID, short(done.Commit)))
	} else {
		report(fmt.Sprintf("done: %s was already done here", req.ID))
	}
	if !req.Cleanup {
		return nil
	}
	return cleanup(root, req.Cwd, name, worktreeOf(res, from.Ref), report)
}

// approved finds the one branch holding the record in review with its
// candidate approved. The target itself never counts: what it holds is what
// integration produces, not a candidate.
func approved(res *versions.Result, id, target string) (*versions.Source, *project.Record, error) {
	var review, ok []*versions.Version
	for i := range res.Groups {
		if res.Groups[i].ID != id {
			continue
		}
		for j := range res.Groups[i].Versions {
			v := &res.Groups[i].Versions[j]
			if v.Source.Kind != "committed" || v.Source.Ref == "refs/heads/"+target || v.Record == nil || v.Record.Type != "work" || v.Record.Status != "review" {
				continue
			}
			review = append(review, v)
			if v.Record.Approved != "" {
				ok = append(ok, v)
			}
		}
	}
	names := func(vs []*versions.Version) string {
		var out []string
		for _, v := range vs {
			out = append(out, strings.TrimPrefix(v.Source.Ref, "refs/heads/"))
		}
		return strings.Join(out, ", ")
	}
	switch {
	case len(ok) == 1:
		return ok[0].Source, ok[0].Record, nil
	case len(ok) > 1:
		return nil, nil, fmt.Errorf("%s is approved in review on several branches (%s); integrate needs one", id, names(ok))
	case len(review) != 0:
		v := review[0]
		where := "its checkout"
		if w := worktreeOf(res, v.Source.Ref); w != "" {
			where = w
		}
		return nil, nil, fmt.Errorf("%s is in review on %s but not approved: run grove approve %s VERDICT in %s first", id, names(review), id, where)
	}
	return nil, nil, fmt.Errorf("no branch holds %s in review; nothing to integrate", id)
}

// worktreeOf is the registered checkout of a branch, or "".
func worktreeOf(res *versions.Result, ref string) string {
	for _, s := range res.Sources {
		if s.Kind == "live" && s.Ref == ref {
			return s.Worktree
		}
	}
	return ""
}

// verdict quotes the verdict paragraph the approval appended, when the body
// still holds one for this candidate.
func verdict(r *project.Record) string {
	prefix := "Verdict on candidate " + short(r.Candidate) + ","
	found := ""
	for _, l := range strings.Split(string(r.Source), "\n") {
		if strings.HasPrefix(l, prefix) {
			found = strings.TrimRight(l, "\r")
		}
	}
	if found == "" {
		return ""
	}
	return " (" + found + ")"
}

// cleanup removes the branch's worktree, then the branch, through Git's own
// refusals: a worktree with changes or untracked files and a branch the
// target does not contain are kept. A worktree holding cwd is kept too, and
// one holding ignored files, which git worktree remove would delete.
// Whatever is kept is reported and makes the result an error, since the
// caller asked for a cleanup that did not fully happen.
func cleanup(root, cwd, name, worktree string, report func(string)) error {
	kept := false
	keep := func(what, reason string) {
		kept = true
		report(fmt.Sprintf("cleanup: kept %s: %s", what, reason))
	}
	if worktree != "" {
		if rel, err := filepath.Rel(worktree, cwd); cwd != "" && err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			keep("worktree "+worktree, "it holds this process's working directory")
		} else if ignored, err := ignored(worktree); err != nil {
			keep("worktree "+worktree, err.Error())
		} else if len(ignored) != 0 {
			keep("worktree "+worktree, "it holds ignored files ("+strings.Join(ignored, ", ")+"); remove them or the worktree by hand")
		} else if _, err := repo.Git(root, "worktree", "remove", "--", worktree); err != nil {
			keep("worktree "+worktree, err.Error())
		} else {
			report("cleanup: removed worktree " + worktree)
		}
	}
	if _, err := repo.Git(root, "branch", "-d", "--", name); err != nil {
		keep("branch "+name, err.Error())
	} else {
		report("cleanup: deleted branch " + name)
	}
	if kept {
		return errors.New("cleanup incomplete; the integration stands")
	}
	return nil
}

// ignored lists a checkout's ignored files, which git worktree remove deletes
// without --force; a .env or a build directory is not Grove's to delete.
func ignored(worktree string) ([]string, error) {
	out, err := repo.Git(worktree, "status", "--porcelain", "--ignored", "-z")
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range strings.Split(strings.TrimSuffix(out, "\x00"), "\x00") {
		if strings.HasPrefix(entry, "!! ") {
			files = append(files, entry[3:])
		}
	}
	return files, nil
}

func head(root string) (string, error) {
	out, err := repo.Git(root, "rev-parse", "HEAD")
	return strings.TrimSpace(out), err
}

func orNoBranch(branch string) string {
	if branch == "" {
		return "no branch"
	}
	return branch
}

func short(commit string) string { return commit[:min(len(commit), 7)] }

func diagnostics(ds []project.Diagnostic) string {
	lines := make([]string, len(ds))
	for i, d := range ds {
		lines[i] = d.String()
	}
	return strings.Join(lines, "\n")
}
