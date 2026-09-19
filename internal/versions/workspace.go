package versions

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

// Selection is a parsed selector: which source and record version a caller
// chose explicitly, with the readable commit and revision prefixes and the
// binding digest that decides whether the version is still current.
type Selection struct {
	Kind, Locator, Ref, Commit, ID, Revision, Binding string // Ref "" means detached
}

var (
	hex12   = regexp.MustCompile(`^[0-9a-f]{12}$`)
	hex16   = regexp.MustCompile(`^[0-9a-f]{16}$`)
	idForm  = regexp.MustCompile(`^[WQD]-[0-9]{3,}$`)
	locForm = regexp.MustCompile(`^[^/:\x00-\x1f]+$`)
)

// Parse checks a selector's grammar. It does not consult the repository.
func Parse(selector string) (Selection, error) {
	parts := strings.Split(selector, ":")
	var sel Selection
	fail := func() (Selection, error) {
		return Selection{}, errors.New("selector must look like committed:REF@COMMIT:ID@REVISION:BINDING or live:WORKTREE:REF@COMMIT:ID@REVISION:BINDING as printed by versions")
	}
	switch {
	case len(parts) == 4 && parts[0] == "committed":
		sel.Kind = "committed"
	case len(parts) == 5 && parts[0] == "live" && locForm.MatchString(parts[1]):
		sel.Kind, sel.Locator = "live", parts[1]
		parts = parts[1:]
	default:
		return fail()
	}
	at := strings.LastIndex(parts[1], "@")
	if at < 0 {
		return fail()
	}
	sel.Ref, sel.Commit = parts[1][:at], parts[1][at+1:]
	var ok bool
	if sel.ID, sel.Revision, ok = strings.Cut(parts[2], "@"); !ok {
		return fail()
	}
	sel.Binding = parts[3]
	if !hex12.MatchString(sel.Commit) || !hex12.MatchString(sel.Revision) || !hex16.MatchString(sel.Binding) || !idForm.MatchString(sel.ID) {
		return fail()
	}
	switch {
	case sel.Ref == "detached" && sel.Kind == "live":
		sel.Ref = ""
	case !strings.HasPrefix(sel.Ref, "refs/heads/") || len(sel.Ref) == len("refs/heads/"):
		return fail()
	}
	return sel, nil
}

// Workspace is the checkout holding a selected version, validated at the
// moment of resolution. Another process may change it afterwards; a later
// mutation must check again.
type Workspace struct {
	Checkout, Project, Record string // absolute paths
	Ref, Head, Revision       string // Ref "" means detached
	Selector                  string // the live observation that was resolved
}

// Resolve locates the existing checkout for selector without creating,
// switching, claiming, or editing anything. A committed selection routes to
// the one checkout of its branch whose live record still has the committed
// bytes; otherwise the caller must refresh and select a live observation.
// ponytail: re-runs the whole inspection (one ls-tree and one cat-file per
// branch) to answer one selector; inspect only the selected source if
// repositories with many branches make this slow.
func Resolve(root, selector string) (*Workspace, error) { return resolveWith(root, selector, nil) }

// resolveWith is Resolve with a hook that runs before the final check of the
// selected target, so tests can change it meanwhile.
func resolveWith(root, selector string, before func()) (*Workspace, error) {
	sel, err := Parse(selector)
	if err != nil {
		return nil, err
	}
	res, err := Inspect(root, sel.ID)
	if err != nil {
		return nil, err
	}
	var s *Source
	for _, candidate := range res.Sources {
		if candidate.Kind == sel.Kind && (sel.Kind == "live" && candidate.Locator == sel.Locator || sel.Kind == "committed" && candidate.Ref == sel.Ref) {
			s = candidate
		}
	}
	reselect := "; run versions and reselect"
	if s == nil {
		if sel.Kind == "live" {
			return nil, fmt.Errorf("worktree %s is no longer registered in this repository (moved or removed)%s", sel.Locator, reselect)
		}
		return nil, fmt.Errorf("branch %s no longer exists%s", sel.Ref, reselect)
	}
	v, err := observation(res, s, sel.ID)
	if err != nil {
		return nil, err
	}
	if v.Selector != selector {
		return nil, attribute(sel, v)
	}
	final := func(lv Version) (*Workspace, error) {
		if before != nil {
			before()
		}
		if err := recheck(root, res, sel, s.Commit, lv); err != nil {
			return nil, err
		}
		return workspace(res, lv), nil
	}
	if sel.Kind == "live" {
		return final(v)
	}
	var checkouts []*Source
	for _, candidate := range res.Sources {
		// An entry without a locator could not be entered (prunable or
		// foreign); it is not a checkout anyone can be sent to.
		if candidate.Kind == "live" && candidate.Ref == sel.Ref && candidate.Locator != "" {
			checkouts = append(checkouts, candidate)
		}
	}
	switch len(checkouts) {
	case 0:
		return nil, fmt.Errorf("no registered worktree has %s checked out; this command does not create one", sel.Ref)
	case 1:
	default:
		paths := make([]string, len(checkouts))
		for i, c := range checkouts {
			paths[i] = c.Locator
		}
		return nil, fmt.Errorf("%d worktrees have %s checked out (%s); select one of their live versions instead", len(checkouts), sel.Ref, strings.Join(paths, ", "))
	}
	live := checkouts[0]
	if live.Commit != s.Commit {
		return nil, fmt.Errorf("worktree %s is at %s, not the selected %s tip %s%s", live.Locator, live.Commit, sel.Ref, s.Commit, reselect)
	}
	lv, err := observation(res, live, sel.ID)
	if err != nil {
		return nil, err
	}
	if live.ConfigRevision != s.ConfigRevision {
		return nil, fmt.Errorf("the live grove.yaml in worktree %s differs from the committed configuration selected; run versions and select the live observation instead", live.Locator)
	}
	if lv.Path != v.Path || lv.Revision != v.Revision {
		return nil, fmt.Errorf("the live %s in worktree %s differs from the committed version selected (%s at %s); run versions and select the live observation instead", sel.ID, live.Locator, lv.Change, lv.Path)
	}
	return final(lv)
}

// observation returns the selected record's current version in s.
func observation(res *Result, s *Source, id string) (Version, error) {
	where := "branch " + s.Ref
	if s.Kind == "live" {
		where = "worktree " + s.Locator
	}
	if !s.Present && len(s.Diagnostics) == 0 {
		return Version{}, fmt.Errorf("%s has no grove.yaml at the selected project location; the project is absent there; run versions and reselect", where)
	}
	if !s.Valid {
		return Version{}, fmt.Errorf("%s is not a valid source:\n%s", where, strings.Join(s.Diagnostics, "\n"))
	}
	for _, g := range res.Groups {
		for _, v := range g.Versions {
			if v.Source == s && v.Record != nil {
				return v, nil
			}
			if v.Source == s {
				return Version{}, fmt.Errorf("%s was deleted from the live files of %s (it is still at HEAD as %s); run versions and reselect", id, where, v.Path)
			}
		}
	}
	return Version{}, fmt.Errorf("%s is not present in %s; run versions and reselect", id, where)
}

// attribute explains why the current observation no longer matches the
// selection, from its readable parts first and the binding digest last.
func attribute(sel Selection, v Version) error {
	s := v.Source
	reselect := "; run versions and reselect"
	if s.Kind == "live" && s.Ref != sel.Ref {
		return fmt.Errorf("worktree %s is now %s; the selection was %s%s", s.Locator, describe(s.Ref, s.Commit), describe(sel.Ref, sel.Commit), reselect)
	}
	if s.Commit[:12] != sel.Commit {
		what := "branch " + s.Ref
		if s.Kind == "live" {
			what = "worktree " + s.Locator + "'s HEAD"
		}
		return fmt.Errorf("%s moved from %s to %s%s", what, sel.Commit, s.Commit[:12], reselect)
	}
	if strings.TrimPrefix(v.Revision, "sha256:")[:12] != sel.Revision {
		return fmt.Errorf("%s changed since it was selected (revision %s, selected %s)%s", sel.ID, strings.TrimPrefix(v.Revision, "sha256:")[:12], sel.Revision, reselect)
	}
	return fmt.Errorf("the selection no longer matches %s at %s: its worktree path, configuration, record path, or project location changed%s", sel.ID, v.Path, reselect)
}

func workspace(res *Result, v Version) *Workspace {
	project := filepath.Join(v.Source.Worktree, filepath.FromSlash(res.Prefix))
	return &Workspace{
		Checkout: v.Source.Worktree, Project: project, Record: filepath.Join(project, filepath.FromSlash(v.Path)),
		Ref: v.Source.Ref, Head: v.Source.Commit, Revision: v.Revision, Selector: v.Selector,
	}
}

// recheck confirms, from a fresh inventory and fresh reads, that the workspace
// about to be returned is still the observation lv: its registration,
// ownership, configuration, record path and bytes, and branch or detached
// HEAD. A committed route also needs its branch tip unmoved and no second
// enterable checkout of that branch. Nothing here substitutes a new selection,
// and changes after this check remain possible.
func recheck(root string, res *Result, sel Selection, tip string, lv Version) error {
	const reselect = "; run versions and reselect"
	worktrees, err := repo.Worktrees(root)
	if err != nil {
		return err
	}
	var target *Source
	for _, w := range worktrees {
		switch {
		case w.Bare:
		case w.Path == lv.Source.Worktree:
			target = enterWorktree(w, res.Repository)
			target.load(root, res.Repository, res.Prefix, map[string]*tree{})
		case sel.Kind == "committed" && w.Branch == sel.Ref:
			if other := enterWorktree(w, res.Repository); other.Locator != "" {
				return fmt.Errorf("worktree %s also checked out %s while its workspace was being resolved%s", other.Locator, sel.Ref, reselect)
			}
		}
	}
	where := "worktree " + lv.Source.Locator
	if target == nil {
		return fmt.Errorf("%s was removed or moved while its workspace was being resolved%s", where, reselect)
	}
	if !target.Valid {
		if len(target.Diagnostics) == 0 {
			return fmt.Errorf("%s lost its project while its workspace was being resolved%s", where, reselect)
		}
		return fmt.Errorf("%s stopped being a valid source while its workspace was being resolved:\n%s", where, strings.Join(target.Diagnostics, "\n"))
	}
	current := ""
	for _, r := range target.project.Records {
		if r.ID == sel.ID {
			current = selector(res.Repository, res.Prefix, target, r.ID, r.Path, project.Revision(r.Source))
		}
	}
	if current != lv.Selector {
		return fmt.Errorf("%s or its %s changed while its workspace was being resolved%s", where, sel.ID, reselect)
	}
	if sel.Kind == "committed" {
		out, err := repo.Git(root, "rev-parse", "--verify", "--quiet", sel.Ref+"^{commit}")
		if err != nil || strings.TrimSpace(out) != tip {
			return fmt.Errorf("branch %s moved while its workspace was being resolved%s", sel.Ref, reselect)
		}
	}
	return nil
}
