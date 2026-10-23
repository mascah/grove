package cli

import (
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mascah/grove/internal/doc"
	"github.com/mascah/grove/internal/gitx"
)

// archive moves change id's folder in the working tree to changes/archive/,
// once it is dropped, or finished with the maintained specs holding what it
// gives them.
func archive(stdout io.Writer, id string) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	root := filepath.Join(p.Top, p.Root)
	c, _, err := change(os.DirFS(root), p.instruct().Layout, id)
	if err != nil {
		return err
	}
	for _, t := range c.Tasks {
		if !c.Dropped && t.Status != "done" {
			return fmt.Errorf("%s is not done in implementation.md; archive a change when every task is done or the file says Status: dropped", t.ID)
		}
	}
	if !c.Dropped {
		ds, err := fold(p, os.DirFS(root), c)
		if err != nil {
			return err
		}
		for _, d := range ds {
			fmt.Fprintf(stdout, "%s/%s\n", p.Root, d)
		}
		if len(ds) > 0 {
			return fmt.Errorf("%s found, so %s stays where it is; fix the maintained specs and run grove archive again", plural(len(ds), "problem"), id)
		}
	}
	to := "changes/archive/" + filepath.Base(c.Dir)
	if err := os.MkdirAll(filepath.Join(root, "changes", "archive"), 0o755); err != nil {
		return err
	}
	if err := os.Rename(filepath.Join(root, c.Dir), filepath.Join(root, to)); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "Archived %s to %s/%s/; nothing is committed.\n", id, p.Root, to)
	return nil
}

// given is what a maintained requirement should be, and who gives it: a
// change's ID, or "" for where this branch left the target.
type given struct {
	doc.Requirement
	path, by string
}

// fold returns a problem for each requirement of a maintained spec c has a
// change spec for, in root, that is not what it should be, whitespace aside:
// its text where this branch left the target, then what each change archived
// on this branch since then gives it, then what c gives it.
func fold(p project, root fs.FS, c *doc.Change) ([]doc.Diagnostic, error) {
	left, err := gitx.Run(p.Top, "merge-base", "HEAD", "refs/heads/"+p.Target)
	if err != nil {
		return nil, err
	}
	base, err := gitx.TreeFS(p.Top, strings.TrimSpace(left), p.Root)
	if err != nil {
		return nil, err
	}
	now, err := doc.Load(root, p.instruct().Layout)
	if err != nil {
		return nil, err
	}
	var since []*doc.Change
	for _, ch := range now.Archived {
		if _, err := fs.Stat(base, ch.Dir); !ch.Dropped && err != nil {
			since = append(since, ch)
		}
	}
	// A dependency archived before this branch left the target is met.
	ordered, waiting := doc.Order(since, func(id string) bool {
		return !slices.ContainsFunc(since, func(ch *doc.Change) bool { return ch.ID == id })
	})
	var out []doc.Diagnostic
	for _, capability := range c.Capabilities() {
		spec := "specs/" + capability + "/spec.md"
		want := map[string]given{}
		b, _ := fs.ReadFile(base, spec)
		for _, r := range doc.Requirements(string(b)) {
			want[r.Name] = given{r, spec, ""}
		}
		for _, ch := range slices.Concat(ordered, waiting, []*doc.Change{c}) {
			for _, file := range slices.Sorted(maps.Keys(ch.Touches)) {
				if ch.Touches[file] != capability {
					continue
				}
				b, _ := fs.ReadFile(root, file)
				for _, r := range doc.Requirements(string(b)) {
					if r.Section != "" {
						want[r.Name] = given{r, file, ch.ID}
					}
				}
			}
		}
		b, _ = fs.ReadFile(root, spec)
		for _, r := range doc.Requirements(string(b)) {
			g, ok := want[r.Name]
			delete(want, r.Name)
			same := ok && strings.Join(strings.Fields(r.Text), " ") == strings.Join(strings.Fields(g.Text), " ")
			msg := ""
			switch {
			case !ok:
				msg = fmt.Sprintf("requirement %q is here though no change adds it", r.Name)
			case g.Section == "REMOVED":
				msg = fmt.Sprintf("requirement %q is still here though %s removes it", r.Name, g.by)
			case same:
				continue
			case g.by == "":
				msg = fmt.Sprintf("requirement %q changed though no change names it", r.Name)
			default:
				msg = fmt.Sprintf("requirement %q is not as %s gives it", r.Name, g.by)
			}
			out = append(out, doc.Diagnostic{Path: spec, Line: r.Line, Message: msg})
		}
		for _, name := range slices.Sorted(maps.Keys(want)) {
			switch g := want[name]; {
			case g.Section == "REMOVED":
			case g.by == "":
				out = append(out, doc.Diagnostic{Path: spec, Message: fmt.Sprintf("requirement %q is gone though no change removes it", name)})
			default:
				out = append(out, doc.Diagnostic{Path: g.path, Line: g.Line, Message: fmt.Sprintf("requirement %q is not in %s though %s gives it", name, spec, g.by)})
			}
		}
	}
	return out, nil
}
