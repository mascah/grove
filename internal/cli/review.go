package cli

import (
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/mascah/grove-next2/internal/gitx"
	"github.com/mascah/grove-next2/internal/instruct"
)

// review shows what merging the workspace that holds change id would agree
// to: each document the merge adds, changes or removes, with its meaning,
// and the diff.
func review(stdout io.Writer, id string) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	all, err := p.workspaces()
	if err != nil {
		return err
	}
	for _, ws := range all {
		if ws.Docs.Change(id) == nil {
			continue
		}
		span := p.Target + "..." + ws.Branch
		names, err := gitx.Run(p.Main, "diff", "--name-status", "--no-renames", span, "--", p.Root)
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "Merging %s into %s approves change %s as future work; it does not implement it.\n\n", ws.Branch, p.Target, id)
		for line := range strings.Lines(names) {
			status, path, _ := strings.Cut(strings.TrimSpace(line), "\t")
			fmt.Fprintf(stdout, "  %-8s %-16s %s\n", map[string]string{"A": "adds", "M": "changes", "D": "removes"}[status], meaning(p.Root, path), path)
		}
		diff, err := gitx.Run(p.Main, "diff", span, "--", p.Root)
		fmt.Fprintf(stdout, "\n%s", diff)
		return err
	}
	return fmt.Errorf("no workspace holds change %s", id)
}

// meaning names what a document under the root stands for once merged.
func meaning(root, path string) string {
	switch rest := strings.TrimPrefix(path, root+"/"); {
	case strings.HasPrefix(rest, "changes/"):
		return "change"
	case strings.HasPrefix(rest, "explorations/"):
		return "exploration"
	case strings.HasPrefix(rest, "specs/"):
		return "maintained spec"
	}
	return "brief"
}

func listChanges(stdout io.Writer) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ip, err := p.agreed()
	if err != nil {
		return err
	}
	for _, c := range ip.Agreed.Changes {
		fmt.Fprintf(stdout, "%s  %s  approved\n", c.Proposal.Meta.ID, c.Proposal.Meta.Title)
		problems(stdout, c.Proposal.Diagnostics)
	}
	return nil
}

// instructionsImplement prints the implement instructions for change id as
// the target holds it.
func instructionsImplement(stdout io.Writer, id string) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ip, err := p.agreed()
	if err != nil {
		return err
	}
	c := ip.Agreed.Change(id)
	if c == nil {
		return fmt.Errorf("%s holds no approved change %s", p.Target, id)
	}
	fsys, err := gitx.TreeFS(p.Main, ip.Commit, p.Root)
	if err != nil {
		return err
	}
	files := map[string]string{}
	for _, name := range c.Files {
		b, err := fs.ReadFile(fsys, name)
		if err != nil {
			return err
		}
		files[name] = string(b)
	}
	fmt.Fprint(stdout, instruct.Implement(ip, c, files))
	return nil
}
