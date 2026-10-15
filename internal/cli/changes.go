package cli

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/mascah/grove-next2/internal/gitx"
	"github.com/mascah/grove-next2/internal/instruct"
)

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
	fsys, err := gitx.TreeFS(p.Top, ip.Commit, p.Root)
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
	fmt.Fprint(stdout, instruct.Implement(ip, c, files).Markdown())
	return nil
}
