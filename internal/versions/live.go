package versions

import (
	"context"
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
// ponytail: a live checkout costs three rev-parse processes per inspection for
// a project at the repository root and five for a nested one (entering, the
// prefix, the record folder, and the second inventory's re-entry), run one
// after another. Read checkouts concurrently if dozens of worktrees make an
// inspection slow; branches cost no processes.
func identity(ctx context.Context, dir string) (gitDir, prefix, common string, err error) {
	paths, err := repo.GitPathsContext(ctx, dir, "--git-dir", "--show-prefix", "--git-common-dir")
	if err != nil {
		return "", "", "", err
	}
	return paths[0], paths[1], paths[2], nil
}

// enterWorktree starts a live source for a registered worktree. The source
// gets a Git directory and locator only when the path is an enterable checkout
// of the repository at common, judged from the path itself and not from its
// registration alone.
func enterWorktree(ctx context.Context, w repo.Worktree, common string) *Source {
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
	gitDir, itsPrefix, itsCommon, err := identity(ctx, w.Path)
	if err != nil {
		s.fail("cannot enter worktree: " + err.Error())
		return s
	}
	// The shape of the Git directory is not enough: a separate repository can
	// keep its Git directory below <common>/worktrees, and only its common
	// directory gives it away. A plain directory inside another checkout
	// answers with that checkout's identity and a prefix; a worktree's own
	// root has none.
	switch rest, linked := strings.CutPrefix(gitDir, filepath.Join(common, "worktrees")+string(filepath.Separator)); {
	case itsCommon != common || itsPrefix != "":
		s.fail("worktree path no longer belongs to this repository")
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
func (s *Source) locate(ctx context.Context, prefix string) (dir string, ok bool) {
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
	if err := s.owns(ctx, dir, prefix); err != nil {
		s.fail(err.Error())
		return "", false
	}
	return dir, true
}

// owns checks that Git places dir in this source's worktree at prefix.
func (s *Source) owns(ctx context.Context, dir, prefix string) error {
	gitDir, itsPrefix, _, err := identity(ctx, dir)
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
// A source loaded under a cancelled ctx means nothing; callers check ctx.
func (s *Source) load(ctx context.Context, prefix string, committed *objects) {
	if s.Locator == "" {
		return
	}
	dir, ok := s.locate(ctx, prefix)
	if !ok || ctx.Err() != nil {
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
	if err := s.owns(ctx, filepath.Join(dir, filepath.FromSlash(p.RecordDir)), prefix+strings.Trim(filepath.ToSlash(filepath.Clean(p.RecordDir)), "/")+"/"); err != nil {
		s.fail(err.Error())
		return
	}
	s.project, s.Valid, s.ConfigRevision = p, true, project.Revision(p.Config)
	if strings.Trim(s.Commit, "0") == "" { // unborn branch: nothing is committed yet
		s.baseline = &tree{}
	} else {
		s.baseline = committed.loadTree(s.Commit)
	}
	if s.baseline.err != nil {
		s.fail("cannot read HEAD " + s.Commit + ": " + s.baseline.err.Error())
	} else if s.baseline.present && !s.baseline.valid() {
		s.Note = "HEAD " + s.Commit + " does not validate, so changes against it are unknown"
	}
}
