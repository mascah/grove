package cli

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/mascah/grove/internal/doc"
	"github.com/mascah/grove/internal/gitx"
)

// review shows what merging the current branch into the target would agree
// to: each change it approves, as future work or with its code, each document it adds, changes or removes
// with its meaning, and the diff.
func review(stdout io.Writer) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ip, err := p.agreed()
	if err != nil {
		return err
	}
	branch, err := gitx.Run(p.Top, "branch", "--show-current")
	if err != nil {
		return err
	}
	fsys, err := gitx.TreeFS(p.Top, "HEAD", p.Root)
	if err != nil {
		return err
	}
	head, err := doc.Load(fsys, ip.Layout)
	if err != nil {
		return err
	}
	branch = strings.TrimSpace(branch)
	approves := 0
	for _, c := range head.Changes {
		if ip.Agreed.Change(c.ID) == nil {
			fmt.Fprintf(stdout, "Merging %s into %s approves change %s as future work; it does not implement it.\n", branch, p.Target, c.ID)
			approves++
		}
	}
	for _, c := range head.Archived { // drafted and built on this branch (rule 4)
		if head.Finished(c.ID) && ip.Agreed.Change(c.ID) == nil && !ip.Agreed.InArchive(c.ID) {
			fmt.Fprintf(stdout, "Merging %s into %s approves change %s with its code.\n", branch, p.Target, c.ID)
			approves++
		}
	}
	if approves == 0 {
		fmt.Fprintf(stdout, "Merging %s into %s approves no new change.\n", branch, p.Target)
	}
	span := p.Target + "...HEAD"
	names, err := gitx.Run(p.Top, "diff", "--name-status", "--no-renames", span, "--", p.Root)
	if err != nil {
		return err
	}
	fmt.Fprintln(stdout)
	if err := planned(stdout, p, ip.Agreed, head, span); err != nil {
		return err
	}
	for line := range strings.Lines(names) {
		status, path, _ := strings.Cut(strings.TrimSpace(line), "\t")
		fmt.Fprintf(stdout, "  %-8s %-16s %s\n", map[string]string{"A": "adds", "M": "changes", "D": "removes"}[status], meaning(p.Root, path), path)
	}
	diff, err := gitx.Run(p.Top, "diff", span, "--", p.Root)
	fmt.Fprintf(stdout, "\n%s", diff)
	return err
}

// meaning names what a document under the root stands for once merged.
func meaning(root, path string) string {
	switch {
	case strings.HasPrefix(path, root+"/changes/archive/"):
		return "archived change"
	case strings.HasPrefix(path, root+"/changes/"):
		return "change"
	case strings.HasPrefix(path, root+"/schemas/"):
		return "change schema"
	case strings.HasPrefix(path, root+"/cycles/"):
		return "cycle"
	}
	return "maintained spec"
}

// planned lists each file the branch changes outside the root against the
// plan of the changes it holds open, or archived where agreed does not.
func planned(stdout io.Writer, p project, agreed, head *doc.Collection, span string) error {
	changes := slices.Clone(head.Changes)
	for _, c := range head.Archived {
		if !agreed.InArchive(c.ID) {
			changes = append(changes, c)
		}
	}
	if len(changes) == 0 {
		return nil
	}
	names, err := gitx.Run(p.Top, "diff", "--name-only", "--no-renames", span)
	if err != nil {
		return err
	}
	listed := false
	for path := range strings.Lines(names) {
		if path = strings.TrimSpace(path); path != "" && !strings.HasPrefix(path, p.Root+"/") {
			fmt.Fprintf(stdout, "  %s  %s\n", mark(changes, path), path)
			listed = true
		}
	}
	if listed {
		fmt.Fprintln(stdout)
	}
	return nil
}

// mark says whether a task of changes plans path on its Files: line, or
// failing that agrees to it on an Also changed: line, or none names it.
func mark(changes []*doc.Change, path string) string {
	for _, c := range changes {
		for _, t := range c.Tasks {
			if slices.ContainsFunc(t.Files, func(f string) bool { return doc.Names(f, path) }) {
				return "planned " + c.ID + ":" + t.ID
			}
		}
	}
	for _, c := range changes {
		for _, t := range c.Tasks {
			if reason, ok := t.Also[path]; ok {
				return "agreed " + c.ID + ":" + t.ID + ": " + reason
			}
		}
	}
	return "named nowhere"
}
