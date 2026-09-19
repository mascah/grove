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

// identity asks Git where dir is: its Git directory, its common directory,
// and its path inside the worktree.
func identity(dir string) (gitDir, common, prefix string, err error) {
	out, err := repo.Git(dir, "rev-parse", "--path-format=absolute", "--git-dir", "--git-common-dir", "--show-prefix")
	if err != nil {
		return "", "", "", err
	}
	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	if len(lines) != 3 {
		return "", "", "", fmt.Errorf("git rev-parse: unexpected reply %q", out)
	}
	return lines[0], lines[1], lines[2], nil
}

// enterWorktree starts a live source for a registered worktree. The source
// gets a Git directory and locator only when the path is an enterable checkout
// of the repository at common, judged from the path itself and not from its
// registration alone.
func enterWorktree(w worktree, common string) *Source {
	s := &Source{Kind: "live", Ref: w.branch, Commit: w.head, Worktree: w.path}
	if w.prunable != "" {
		s.fail("worktree is prunable: " + w.prunable)
		return s
	}
	if info, err := os.Lstat(w.path); err != nil {
		s.fail("cannot enter worktree: " + err.Error())
		return s
	} else if !info.IsDir() {
		s.fail("worktree path is a symlink or not a directory")
		return s
	}
	gitDir, itsCommon, itsPrefix, err := identity(w.path)
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
	if err := s.owns(dir, common, prefix); err != nil {
		s.fail(err.Error())
		return "", false
	}
	return dir, true
}

// owns checks that Git places dir in this source's worktree at prefix.
func (s *Source) owns(dir, common, prefix string) error {
	gitDir, itsCommon, itsPrefix, err := identity(dir)
	if err != nil {
		return fmt.Errorf("cannot identify %s: %w", dir, err)
	}
	if gitDir != s.GitDir || itsCommon != common || itsPrefix != prefix {
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
	if err := s.owns(filepath.Join(dir, filepath.FromSlash(p.RecordDir)), common, prefix+strings.Trim(filepath.ToSlash(filepath.Clean(p.RecordDir)), "/")+"/"); err != nil {
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
