package cli

import (
	"fmt"
	"io"
	"strings"
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
		waits := ""
		if len(c.WaitsOn) > 0 {
			waits = "  waits on " + strings.Join(c.WaitsOn, ", ")
		}
		fmt.Fprintf(stdout, "%s  %s  approved%s\n", c.ID, c.Meta.Title, waits)
		problems(stdout, c.Diagnostics)
	}
	return nil
}
