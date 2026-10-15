package cli

import (
	"fmt"
	"io"
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
