package versions

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

// identity asks Git where dir is: its Git directory and its path inside the
// worktree. The Git directory is unique to one worktree of one repository.
// ponytail: one git process per path, because rev-parse cannot delimit paths
// that contain newlines. A live checkout costs about ten rev-parse processes
// per inspection (entering, the prefix, the record folder, and the second
// inventory's re-entry) where it used to cost one; batch them or skip the
// prefix checks for root projects if inspecting many worktrees becomes slow.
func identity(dir string) (gitDir, prefix string, err error) {
	if gitDir, err = repo.GitPath(dir, "--git-dir"); err == nil {
		prefix, err = repo.GitPath(dir, "--show-prefix")
	}
	return gitDir, prefix, err
}

// enterWorktree starts a live source for a registered worktree. The source
// gets a Git directory and locator only when the path is an enterable checkout
// of the repository at common, judged from the path itself and not from its
// registration alone.
func enterWorktree(w repo.Worktree, common string) *Source {
	s := &Source{Kind: "live", Ref: w.Branch, Commit: w.Head, Worktree: w.Path}
	if w.Prunable != "" {
		s.fail("worktree is prunable: " + w.Prunable)
		return s
	}
	if info, err := os.Lstat(w.Path); err != nil {
		s.fail("cannot enter worktree: " + err.Error())
		return s
	} else if !info.IsDir() {
		s.fail("worktree path is a symlink or not a directory")
		return s
	}
	gitDir, itsPrefix, err := identity(w.Path)
	var itsCommon string
	if err == nil {
		itsCommon, err = repo.GitPath(w.Path, "--git-common-dir")
	}
	if err != nil {
		s.fail("cannot enter worktree: " + err.Error())
		return s
	}
	// A plain directory inside another checkout answers with that checkout's
	// identity and a prefix; a worktree's own root has none.
	if itsCommon != common || itsPrefix != "" {
		s.fail("worktree path no longer belongs to this repository")
		return s
	}
	switch rest, linked := strings.CutPrefix(gitDir, filepath.Join(common, "worktrees")+string(filepath.Separator)); {
	case gitDir == common:
		s.GitDir, s.Locator = gitDir, "."
	case linked && rest != "" && !strings.ContainsRune(rest, filepath.Separator):
		s.GitDir, s.Locator = gitDir, rest
	default:
		s.fail("unrecognized worktree Git directory " + gitDir)
	}
	return s
}

// locate returns the project directory at prefix below an entered checkout.
// ok is false without a diagnostic when a component is genuinely missing (the
// project is absent), and false with one when the location is a symlink, is
// not a directory, cannot be read, or belongs to another repository or
// worktree. Components are checked one by one without following symlinks, so
// a link cannot lead out of the checkout or stand in for the project.
func (s *Source) locate(common, prefix string) (dir string, ok bool) {
	dir = s.Worktree
	if prefix == "" {
		return dir, true
	}
	for _, part := range strings.Split(strings.TrimSuffix(prefix, "/"), "/") {
		dir = filepath.Join(dir, part)
		info, err := os.Lstat(dir)
		switch {
		case errors.Is(err, fs.ErrNotExist):
			return "", false
		case err != nil:
			s.fail("cannot read the project location: " + err.Error())
			return "", false
		case info.Mode()&fs.ModeSymlink != 0:
			s.fail("project location " + dir + " is a symlink; the project must be inside its checkout")
			return "", false
		case !info.IsDir():
			s.fail("project location " + dir + " is not a directory")
			return "", false
		}
	}
	if err := s.owns(dir, prefix); err != nil {
		s.fail(err.Error())
		return "", false
	}
	return dir, true
}

// owns checks that Git places dir in this source's worktree at prefix.
func (s *Source) owns(dir, prefix string) error {
	gitDir, itsPrefix, err := identity(dir)
	if err != nil {
		return fmt.Errorf("cannot identify %s: %w", dir, err)
	}
	if gitDir != s.GitDir || itsPrefix != prefix {
		return fmt.Errorf("%s belongs to another repository or worktree, not to this checkout", dir)
	}
	return nil
}

// load reads and validates the project of an entered checkout, and its HEAD
// baseline. An absent project leaves the source not present and undiagnosed.
func (s *Source) load(root, common, prefix string, trees map[string]*tree) {
	if s.Locator == "" {
		return
	}
	dir, ok := s.locate(common, prefix)
	if !ok {
		return
	}
	if _, err := os.Lstat(filepath.Join(dir, "grove.yaml")); errors.Is(err, fs.ErrNotExist) {
		return
	} else if err != nil {
		s.fail(err.Error())
		return
	}
	s.Present = true
	p, ds := project.Load(dir, dir)
	if len(ds) != 0 {
		for _, d := range ds {
			s.fail(d.String())
		}
		return
	}
	// The loader already refuses symlinks in the record folder's path; a nested
	// repository there would be foreign in the same way as one at the prefix.
	if err := s.owns(filepath.Join(dir, filepath.FromSlash(p.RecordDir)), prefix+strings.Trim(filepath.ToSlash(filepath.Clean(p.RecordDir)), "/")+"/"); err != nil {
		s.fail(err.Error())
		return
	}
	s.project, s.Valid, s.ConfigRevision = p, true, project.Revision(p.Config)
	if strings.Trim(s.Commit, "0") == "" { // unborn branch: nothing is committed yet
		s.baseline = &tree{}
	} else {
		s.baseline = loadTree(root, s.Commit, trees)
	}
	if s.baseline.err != nil {
		s.fail("cannot read HEAD " + s.Commit + ": " + s.baseline.err.Error())
	} else if s.baseline.present && !s.baseline.valid() {
		s.Note = "HEAD " + s.Commit + " does not validate, so changes against it are unknown"
	}
}
