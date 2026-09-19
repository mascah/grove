// Package versions inspects one project's records across local branch tips
// and registered worktrees. It reads Git objects and live files and changes
// none of them: no refs, index, worktrees, records, or coordination state.
package versions

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

// Source is one place records were read from: a branch tip or a worktree.
type Source struct {
	Kind           string // "committed" or "live"
	Ref            string // full branch ref; "" when a live source is detached
	Commit         string // the branch tip, or the worktree's HEAD
	Worktree       string // live: absolute checkout path
	Locator        string // live: "." for the main worktree, else the administrative name
	GitDir         string // live: absolute Git directory, bound by selectors
	Present        bool   // grove.yaml exists at the project prefix
	Valid          bool   // present, validated, and stable while being read
	ConfigRevision string
	Note           string // live: why changes against HEAD are unknown
	Diagnostics    []string
	project        *project.Project
	baseline       *tree // live: the HEAD commit's project
}

// Detached reports a live source with no branch.
func (s *Source) Detached() bool { return s.Kind == "live" && s.Ref == "" }

func (s *Source) fail(message string) {
	s.Valid = false
	s.Diagnostics = append(s.Diagnostics, message)
}

// Version is one observation of one record in one source.
type Version struct {
	Source   *Source
	Record   *project.Record // nil for a record deleted from live files
	Path     string
	Revision string
	Change   string // live: unchanged, modified, renamed, added, deleted, unknown
	HeadPath string // live: the record's path at HEAD when it differs
	Selector string // "" for deleted rows, which cannot be opened
}

// Group holds every observation of one record ID.
type Group struct {
	ID       string
	Versions []Version
}

// Result is the inventory and observations of one inspection.
type Result struct {
	Project, Repository, Prefix string
	Complete                    bool // every source is valid or absent
	Sources                     []*Source
	Groups                      []Group
}

// Inspect reads every local branch tip and registered worktree of root's
// repository at root's prefix. id restricts Groups to one record; an empty
// Groups with an id means the record is in no valid source.
func Inspect(root, id string) (*Result, error) { return inspect(root, id, nil) }

// inspect is Inspect with a hook that runs after the reads and before the
// worktree inventory is compared, so tests can change identities meanwhile.
func inspect(root, id string, between func()) (*Result, error) {
	common, prefix, err := repo.Locate(root)
	if err != nil {
		return nil, err
	}
	result := &Result{Project: root, Repository: common, Prefix: prefix}
	branches, err := listBranches(root)
	if err != nil {
		return nil, err
	}
	first, err := listWorktrees(root)
	if err != nil {
		return nil, err
	}
	trees := map[string]*tree{}
	for _, b := range branches {
		s := &Source{Kind: "committed", Ref: b.ref, Commit: b.commit}
		s.admit(loadTree(root, b.commit, trees))
		result.Sources = append(result.Sources, s)
	}
	for _, w := range first {
		if w.bare {
			continue
		}
		s := &Source{Kind: "live", Ref: w.branch, Commit: w.head, Worktree: w.path}
		result.Sources = append(result.Sources, s)
		if w.prunable != "" {
			s.fail("worktree is prunable: " + w.prunable)
			continue
		}
		out, err := repo.Git(w.path, "rev-parse", "--path-format=absolute", "--git-dir", "--git-common-dir")
		if err != nil {
			s.fail("cannot enter worktree: " + err.Error())
			continue
		}
		dirs := strings.Split(strings.TrimRight(out, "\n"), "\n")
		if len(dirs) != 2 || dirs[1] != common {
			s.fail("worktree path no longer belongs to this repository")
			continue
		}
		s.GitDir = dirs[0]
		switch rest, linked := strings.CutPrefix(s.GitDir, filepath.Join(common, "worktrees")+string(filepath.Separator)); {
		case s.GitDir == common:
			s.Locator = "."
		case linked && rest != "" && !strings.ContainsRune(rest, filepath.Separator):
			s.Locator = rest
		default:
			s.fail("unrecognized worktree Git directory " + s.GitDir)
			continue
		}
		dir := filepath.Join(w.path, filepath.FromSlash(prefix))
		if _, err := os.Lstat(filepath.Join(dir, "grove.yaml")); errors.Is(err, fs.ErrNotExist) {
			continue
		} else if err != nil {
			s.fail(err.Error())
			continue
		}
		s.Present = true
		p, ds := project.Load(dir, dir)
		if len(ds) != 0 {
			for _, d := range ds {
				s.fail(d.String())
			}
			continue
		}
		s.project, s.Valid, s.ConfigRevision = p, true, project.Revision(p.Config)
		if strings.Trim(w.head, "0") == "" { // unborn branch: nothing is committed yet
			s.baseline = &tree{}
		} else {
			s.baseline = loadTree(root, w.head, trees)
		}
		if s.baseline.err != nil {
			s.fail("cannot read HEAD " + w.head + ": " + s.baseline.err.Error())
		} else if s.baseline.present && !s.baseline.valid() {
			s.Note = "HEAD " + w.head + " does not validate, so changes against it are unknown"
		}
	}
	if between != nil {
		between()
	}
	second, err := listWorktrees(root)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	for _, s := range result.Sources {
		if s.Kind != "live" {
			continue
		}
		seen[s.Worktree] = true
		i := slices.IndexFunc(second, func(w worktree) bool { return w.path == s.Worktree })
		if i < 0 {
			s.fail("worktree was removed while being read")
		} else if w := second[i]; w.head != s.Commit || w.branch != s.Ref {
			s.fail(fmt.Sprintf("worktree changed while being read: %s to %s", describe(s.Ref, s.Commit), describe(w.branch, w.head)))
		}
	}
	for _, w := range second {
		if !w.bare && !seen[w.path] {
			s := &Source{Kind: "live", Ref: w.branch, Commit: w.head, Worktree: w.path}
			s.fail("worktree appeared while being read")
			result.Sources = append(result.Sources, s)
		}
	}
	slices.SortFunc(result.Sources, compareSources)
	result.Complete = true
	groups := map[string]*Group{}
	group := func(recordID string) *Group {
		if groups[recordID] == nil {
			groups[recordID] = &Group{ID: recordID}
		}
		return groups[recordID]
	}
	for _, s := range result.Sources {
		if len(s.Diagnostics) != 0 {
			result.Complete = false
		}
		if !s.Valid {
			continue
		}
		for _, r := range s.project.Records {
			if id != "" && r.ID != id {
				continue
			}
			v := Version{Source: s, Record: r, Path: r.Path, Revision: project.Revision(r.Source)}
			if s.Kind == "live" {
				v.Change, v.HeadPath = s.change(r)
			}
			v.Selector = selector(common, prefix, s, r.ID, r.Path, v.Revision)
			g := group(r.ID)
			g.Versions = append(g.Versions, v)
		}
		if s.Kind == "live" && s.baseline.valid() {
			for _, h := range s.baseline.project.Records {
				if (id == "" || h.ID == id) && !slices.ContainsFunc(s.project.Records, func(r *project.Record) bool { return r.ID == h.ID }) {
					g := group(h.ID)
					g.Versions = append(g.Versions, Version{Source: s, Path: h.Path, Change: "deleted"})
				}
			}
		}
	}
	for _, g := range groups {
		result.Groups = append(result.Groups, *g)
	}
	slices.SortFunc(result.Groups, func(a, b Group) int { return compareIDs(a.ID, b.ID) })
	return result, nil
}

func (s *Source) admit(t *tree) {
	switch {
	case t.err != nil:
		s.fail(t.err.Error())
	case !t.present:
	case len(t.ds) != 0:
		s.Present = true
		for _, d := range t.ds {
			s.fail(d.String())
		}
	default:
		s.Present, s.Valid, s.project, s.ConfigRevision = true, true, t.project, project.Revision(t.project.Config)
	}
}

// change classifies a live record against the same ID at the checkout's HEAD.
func (s *Source) change(r *project.Record) (change, headPath string) {
	t := s.baseline
	if !t.present {
		return "added", ""
	}
	if !t.valid() {
		return "unknown", ""
	}
	i := slices.IndexFunc(t.project.Records, func(h *project.Record) bool { return h.ID == r.ID })
	if i < 0 {
		return "added", ""
	}
	h := t.project.Records[i]
	if h.Path != r.Path {
		headPath = h.Path
	}
	switch {
	case !bytes.Equal(h.Source, r.Source):
		return "modified", headPath
	case headPath != "":
		return "renamed", headPath
	}
	return "unchanged", ""
}

// selector binds one observation to its repository, prefix, source identity
// (kind, ref, worktree path and Git directory), commit, configuration, record
// path, and content. The readable commit and
// revision prefixes attribute common staleness; the binding digest decides.
func selector(common, prefix string, s *Source, id, path, revision string) string {
	ref := s.Ref
	if ref == "" {
		ref = "detached"
	}
	sum := sha256.Sum256([]byte(strings.Join([]string{common, prefix, s.Kind, ref, s.Worktree, s.GitDir, s.Commit, s.ConfigRevision, id, path, revision}, "\x00")))
	tail := fmt.Sprintf("%s@%s:%s@%s:%s", ref, s.Commit[:12], id, strings.TrimPrefix(revision, "sha256:")[:12], hex.EncodeToString(sum[:8]))
	if s.Kind == "live" {
		return "live:" + s.Locator + ":" + tail
	}
	return "committed:" + tail
}

func describe(ref, commit string) string {
	if ref == "" {
		return "detached at " + commit
	}
	return ref + " at " + commit
}

func compareSources(a, b *Source) int {
	if a.Kind != b.Kind { // committed before live
		return strings.Compare(a.Kind, b.Kind)
	}
	if a.Kind == "committed" {
		return strings.Compare(a.Ref, b.Ref)
	}
	if (a.Locator == ".") != (b.Locator == ".") {
		if a.Locator == "." {
			return -1
		}
		return 1
	}
	if c := strings.Compare(a.Locator, b.Locator); c != 0 {
		return c
	}
	return strings.Compare(a.Worktree, b.Worktree)
}

// compareIDs orders W, then Q, then D, numerically within each prefix.
func compareIDs(a, b string) int {
	rank := func(id string) int { return strings.Index("WQD", id[:1]) }
	if c := rank(a) - rank(b); c != 0 {
		return c
	}
	if len(a) != len(b) {
		return len(a) - len(b)
	}
	return strings.Compare(a, b)
}

type branch struct{ ref, commit string }

func listBranches(root string) ([]branch, error) {
	out, err := repo.Git(root, "for-each-ref", "--format=%(refname)%00%(objectname)", "refs/heads/")
	if err != nil {
		return nil, err
	}
	var branches []branch
	for _, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
		if ref, commit, ok := strings.Cut(line, "\x00"); ok {
			branches = append(branches, branch{ref, commit})
		}
	}
	return branches, nil
}

type worktree struct {
	path, head, branch, prunable string
	bare                         bool
}

// listWorktrees parses the NUL-terminated porcelain format, so paths with
// spaces or newlines survive.
func listWorktrees(root string) ([]worktree, error) {
	out, err := repo.Git(root, "worktree", "list", "--porcelain", "-z")
	if err != nil {
		return nil, err
	}
	var worktrees []worktree
	var w worktree
	for _, line := range strings.Split(out, "\x00") {
		key, value, _ := strings.Cut(line, " ")
		switch key {
		case "":
			if w.path != "" {
				worktrees = append(worktrees, w)
			}
			w = worktree{}
		case "worktree":
			w.path = value
		case "HEAD":
			w.head = value
		case "branch":
			w.branch = value
		case "bare":
			w.bare = true
		case "prunable":
			w.prunable = value
			if w.prunable == "" {
				w.prunable = "no reason given"
			}
		}
	}
	return worktrees, nil
}
