package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/mascah/grove-next2/internal/doc"
	"github.com/mascah/grove-next2/internal/gitx"
)

// review shows what merging the current branch into the target would agree
// to: each change it approves, each document it adds, changes or removes
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
	head, err := doc.Load(fsys)
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
	if approves == 0 {
		fmt.Fprintf(stdout, "Merging %s into %s approves no new change.\n", branch, p.Target)
	}
	span := p.Target + "...HEAD"
	names, err := gitx.Run(p.Top, "diff", "--name-status", "--no-renames", span, "--", p.Root)
	if err != nil {
		return err
	}
	fmt.Fprintln(stdout)
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
	}
	return "maintained spec"
}
