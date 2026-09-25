package cli

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mascah/grove"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

const defaultConfig = "schema_version: 3\nrecords: grove\nbrief: grove/brief.md\n"

const placeholderBrief = "# Brief\n\n" +
	"This brief is not written yet. `grove init` created it so that the\n" +
	"configuration validates; it states no purpose, constraint, or direction, and\n" +
	"nothing here is a decision. Develop it in a shaping session (`/grove-shape` in\n" +
	"Claude Code, `$grove-shape` in Codex) and replace this text.\n"

// initStep is one planned path: what init will say about it, and the write it
// still has to do, if any.
type initStep struct {
	path, verdict, note string
	content             []byte // nil: nothing to write
	dir                 bool
}

// runInit sets up the Git checkout at cwd or --project: grove.yaml, the record
// root, a placeholder brief, and the managed harness entrypoints. It plans
// every path first and writes nothing when any path conflicts.
func runInit(cwd string, a invocation, out, errOut io.Writer) int {
	root, err := filepath.Abs(cwd)
	if err == nil && a.project != "" {
		if root = a.project; !filepath.IsAbs(root) {
			root = filepath.Join(cwd, root)
		}
		root = filepath.Clean(root)
	}
	if err == nil {
		var info os.FileInfo
		if info, err = os.Stat(root); err == nil && !info.IsDir() {
			err = fmt.Errorf("%s is not a directory", root)
		}
	}
	if err != nil {
		report(errOut, err)
		return 1
	}
	if _, err := fmt.Fprintf(errOut, "Project: %s\n", visible(root)); err != nil {
		return 1
	}
	_, prefix, err := repo.Locate(root)
	if err != nil {
		report(errOut, fmt.Errorf("init needs the top of a Git checkout: %w", err))
		return 1
	}
	if prefix != "" {
		report(errOut, fmt.Errorf("init needs the top of a Git checkout, and %s is below it, at %s; a nested grove.yaml would end discovery there", visible(root), visible(strings.TrimSuffix(prefix, "/"))))
		return 1
	}
	steps, conflicts := planInit(root)
	if len(conflicts) != 0 {
		for _, c := range conflicts {
			fmt.Fprintf(errOut, "grove: conflict %s\n", visible(c))
		}
		fmt.Fprintln(errOut, "grove: nothing was written; resolve the conflicts and run init again")
		return 1
	}
	for _, s := range steps {
		if s.dir && s.content == nil && s.verdict == "created" {
			err = os.MkdirAll(filepath.Join(root, filepath.FromSlash(s.path)), 0o755)
		} else if s.content != nil {
			full := filepath.Join(root, filepath.FromSlash(s.path))
			if err = os.MkdirAll(filepath.Dir(full), 0o755); err == nil {
				err = os.WriteFile(full, s.content, 0o644)
			}
		}
		if err != nil {
			fmt.Fprintf(errOut, "grove: %s: %s (the paths above were written; run init again after fixing this)\n", visible(s.path), err)
			return 1
		}
		if _, err := fmt.Fprintf(out, "%s %s%s\n", s.verdict, visible(s.path), s.note); err != nil {
			fmt.Fprintf(errOut, "grove: write output: %s (init wrote through %s)\n", err, visible(s.path))
			return 1
		}
	}
	fmt.Fprintln(errOut, "Next: grove check, then commit what init wrote: an attempt's worktree holds only committed\n"+
		"files, and grove run refuses one without the grove-work skill. The entrypoints run `grove`\n"+
		"from PATH and let the agent name its branches; say otherwise in AGENTS.md or CLAUDE.md, which\n"+
		"they defer to for how grove is invoked and how work and proposal branches are named.")
	return 0
}

// planInit decides each path's verdict without writing. Conflicts are the
// reasons init must not proceed.
func planInit(root string) (steps []initStep, conflicts []string) {
	lstat := func(relative string) (os.FileInfo, error) {
		return os.Lstat(filepath.Join(root, filepath.FromSlash(relative)))
	}
	conflict := func(relative, reason string) { conflicts = append(conflicts, relative+": "+reason) }
	// regular reports whether the path is absent (nil, false), a regular file
	// (info, true), or something init cannot replace, which is a conflict.
	// symlinkedParent reports a conflict where a component above relative is a
	// symlink, which would carry the write outside the checkout.
	symlinkedParent := func(relative string) bool {
		for i, c := range relative {
			if c == '/' {
				if parent, err := lstat(relative[:i]); err == nil && parent.Mode()&fs.ModeSymlink != 0 {
					conflict(relative, relative[:i]+" is a symlink")
					return true
				}
			}
		}
		return false
	}
	regular := func(relative string) (os.FileInfo, bool) {
		if symlinkedParent(relative) {
			return nil, false
		}
		info, err := lstat(relative)
		switch {
		case errors.Is(err, fs.ErrNotExist):
			return nil, false
		case err != nil:
			conflict(relative, err.Error())
		case info.Mode()&fs.ModeSymlink != 0:
			conflict(relative, "is a symlink")
		case !info.Mode().IsRegular():
			conflict(relative, "is not a regular file")
		default:
			return info, true
		}
		return nil, false
	}

	recordDir, brief := "grove", "grove/brief.md"
	if _, exists := regular("grove.yaml"); exists {
		source, err := os.ReadFile(filepath.Join(root, "grove.yaml"))
		if err != nil {
			conflict("grove.yaml", err.Error())
		} else if dir, b, ds := project.ParseConfig(source); len(ds) != 0 {
			for _, d := range ds {
				conflict("grove.yaml", "exists but is not a schema 3 configuration: "+strings.TrimPrefix(d.String(), "grove.yaml: "))
			}
		} else {
			recordDir, brief = dir, b
			steps = append(steps, initStep{path: "grove.yaml", verdict: "kept", note: " (exists and validates)"})
		}
	} else {
		steps = append(steps, initStep{path: "grove.yaml", verdict: "created", content: []byte(defaultConfig)})
	}

	if symlinkedParent(recordDir) {
	} else if info, err := lstat(recordDir); errors.Is(err, fs.ErrNotExist) {
		steps = append(steps, initStep{path: recordDir, verdict: "created", dir: true})
	} else if err != nil {
		conflict(recordDir, err.Error())
	} else if info.Mode()&fs.ModeSymlink != 0 || !info.IsDir() {
		conflict(recordDir, "must be a directory for grove.yaml: records")
	} else {
		steps = append(steps, initStep{path: recordDir, verdict: "kept", dir: true})
	}

	if brief != "" {
		if _, exists := regular(brief); exists {
			steps = append(steps, initStep{path: brief, verdict: "kept", note: " (never rewritten)"})
		} else {
			steps = append(steps, initStep{path: brief, verdict: "created", content: []byte(placeholderBrief), note: " (a placeholder that states no intent)"})
		}
	}

	files := grove.Entrypoints()
	for _, relative := range sortedKeys(files) {
		want := []byte(files[relative])
		_, exists := regular(relative)
		if !exists {
			steps = append(steps, initStep{path: relative, verdict: "created", content: want})
			continue
		}
		have, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relative)))
		switch {
		case err != nil:
			conflict(relative, err.Error())
		case !strings.Contains(string(have), grove.ManagedMarker):
			steps = append(steps, initStep{path: relative, verdict: "kept", note: " (not managed by grove init; delete it to get the managed version)"})
		case string(have) == string(want):
			steps = append(steps, initStep{path: relative, verdict: "unchanged"})
		default:
			steps = append(steps, initStep{path: relative, verdict: "updated", content: want})
		}
	}
	return steps, conflicts
}

func sortedKeys(m map[string]string) []string {
	return slices.Sorted(maps.Keys(m))
}
