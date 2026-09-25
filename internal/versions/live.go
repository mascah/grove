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
// ponytail: an inspection costs one rev-parse for the root plus one per live
// checkout whose project reads cleanly (readWorktree), run one after another;
// a checkout that fails a step costs one more per step, and branches cost no
// processes. Read checkouts concurrently if dozens of worktrees make an
// inspection slow.
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
	switch locator, ok := locatorOf(gitDir, common); {
	case itsCommon != common || itsPrefix != "":
		s.fail("worktree path no longer belongs to this repository")
	case ok:
		s.GitDir, s.Locator = gitDir, locator
	default:
		s.fail("unrecognized worktree Git directory " + gitDir)
	}
	return s
}

// locatorOf names the worktree whose Git directory is gitDir: "." for the
// main worktree, the administrative name for a linked one, and ok false for
// a directory of neither shape.
func locatorOf(gitDir, common string) (locator string, ok bool) {
	if gitDir == common {
		return ".", true
	}
	rest, linked := strings.CutPrefix(gitDir, filepath.Join(common, "worktrees")+string(filepath.Separator))
	if linked && rest != "" && !strings.ContainsRune(rest, filepath.Separator) {
		return rest, true
	}
	return "", false
}

// readWorktree enters w and reads its project at prefix, proving both with
// one Git process where it can. Git's prefix in the deepest directory the
// loader used (the record folder, or the project directory when the project
// does not validate), with the Git and common directories it reports, says
// that this directory, every one above it, and w itself are one worktree of
// this repository: discovery would have answered for the first nested or
// foreign repository on the way up instead. Anything else, including a
// prunable entry, a symlink or missing project, or a reply that does not
// match, takes the step-by-step path, whose diagnostics name the first thing
// that fails, at one process per step.
func readWorktree(ctx context.Context, w repo.Worktree, common, prefix string, committed *objects) *Source {
	s := &Source{Kind: "live", Ref: w.Branch, Commit: w.Head, Worktree: w.Path}
	if dir, ok := s.walk(prefix); w.Prunable == "" && isDir(w.Path) && ok {
		if _, err := os.Lstat(filepath.Join(dir, "grove.yaml")); err == nil {
			s.dotGit = dotGit(w.Path)
			p, ds := project.Load(dir, dir)
			deep, want := dir, prefix
			if len(ds) == 0 {
				deep, want = filepath.Join(dir, filepath.FromSlash(p.RecordDir)), recordPrefix(prefix, p)
			}
			gitDir, itsPrefix, itsCommon, err := identity(ctx, deep)
			if locator, ok := locatorOf(gitDir, common); err == nil && ok && itsCommon == common && itsPrefix == want {
				s.GitDir, s.Locator, s.Present = gitDir, locator, true
				s.admitLive(p, ds, committed)
				return s
			}
		}
	}
	s = enterWorktree(ctx, w, common)
	s.load(ctx, prefix, committed)
	return s
}

// unchanged reports, without a Git process, that a checkout read by
// readWorktree's one-process path still has the identity and project it was
// read with: its registration is not prunable, its root and the .git entry
// there are as they were and still lead to this repository's common
// directory, no repository has begun in a directory between its root and the
// project (Git discovery stops at a .git entry or a HEAD file), the project is
// still reached through plain directories, and grove.yaml has the same bytes.
// A checkout read step by step, or one where anything differs, is re-entered
// through Git for its diagnostic.
// ponytail: a .git directory is compared by shape, not contents, and a
// filesystem boundary appearing between the root and the project is not
// looked for; re-enter through Git always if such mid-read edits of Git
// metadata ever matter.
func (s *Source) unchanged(w repo.Worktree, common, prefix string) bool {
	if s.dotGit == "" || w.Prunable != "" || !isDir(w.Path) || dotGit(w.Path) != s.dotGit || commonOf(s.GitDir) != common {
		return false
	}
	dir, ok := (&Source{Worktree: w.Path}).walk(prefix)
	if !ok {
		return false
	}
	for d := dir; d != w.Path; d = filepath.Dir(d) {
		for _, name := range []string{".git", "HEAD"} {
			if _, err := os.Lstat(filepath.Join(d, name)); !errors.Is(err, fs.ErrNotExist) {
				return false
			}
		}
	}
	current, _ := os.ReadFile(filepath.Join(dir, "grove.yaml"))
	return project.Revision(current) == s.ConfigRevision
}

// isDir reports a directory at path itself, not through a symlink.
func isDir(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.IsDir()
}

// dotGit describes dir's .git entry for comparison: "dir" for a directory,
// the bytes of a gitfile, and "" for anything else.
func dotGit(dir string) string {
	info, err := os.Lstat(filepath.Join(dir, ".git"))
	switch {
	case err != nil:
		return ""
	case info.IsDir():
		return "dir"
	case info.Mode().IsRegular():
		if b, err := os.ReadFile(filepath.Join(dir, ".git")); err == nil {
			return "file:" + string(b)
		}
	}
	return ""
}

// commonOf resolves gitDir's common directory as Git does: the commondir
// file inside it, relative to it, or gitDir itself without one. "" means it
// could not be read.
func commonOf(gitDir string) string {
	b, err := os.ReadFile(filepath.Join(gitDir, "commondir"))
	if errors.Is(err, fs.ErrNotExist) {
		return gitDir
	}
	if err != nil {
		return ""
	}
	p := strings.TrimSpace(string(b))
	if !filepath.IsAbs(p) {
		p = filepath.Join(gitDir, p)
	}
	return filepath.Clean(p)
}

// recordPrefix is Git's prefix for the record folder of the project at prefix.
func recordPrefix(prefix string, p *project.Project) string {
	return prefix + strings.Trim(filepath.ToSlash(filepath.Clean(p.RecordDir)), "/") + "/"
}

// locate returns the project directory at prefix below an entered checkout.
// ok is false without a diagnostic when a component is genuinely missing (the
// project is absent), and false with one when the location is a symlink, is
// not a directory, cannot be read, or belongs to another repository or
// worktree.
func (s *Source) locate(ctx context.Context, prefix string) (dir string, ok bool) {
	if dir, ok = s.walk(prefix); !ok || prefix == "" {
		return dir, ok
	}
	if err := s.owns(ctx, dir, prefix); err != nil {
		s.fail(err.Error())
		return "", false
	}
	return dir, true
}

// walk is locate without asking Git. Components are checked one by one
// without following symlinks, so a link cannot lead out of the checkout or
// stand in for the project.
func (s *Source) walk(prefix string) (dir string, ok bool) {
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
	// The loader already refuses symlinks in the record folder's path; a nested
	// repository there would be foreign in the same way as one at the prefix.
	if len(ds) == 0 {
		if err := s.owns(ctx, filepath.Join(dir, filepath.FromSlash(p.RecordDir)), recordPrefix(prefix, p)); err != nil {
			s.fail(err.Error())
			return
		}
	}
	s.admitLive(p, ds, committed)
}

// admitLive takes a located and owned project as this source's, or its
// loader diagnostics, and reads its HEAD baseline.
func (s *Source) admitLive(p *project.Project, ds []project.Diagnostic, committed *objects) {
	if len(ds) != 0 {
		for _, d := range ds {
			s.fail(d.String())
		}
		return
	}
	s.project, s.Valid, s.ConfigRevision, s.Run = p, true, project.Revision(p.Config), p.Run
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
