package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"
)

// listChanges lists the target's approved changes, what each waits on and
// the cycle not done that holds it; as JSON, it prints the project value
// serve answers /api/project with (D-64).
func listChanges(stdout io.Writer, asJSON bool) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	if asJSON {
		v, err := (&server{p: p, scans: map[string][]branchChange{}}).value()
		if err != nil {
			return err
		}
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(v)
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
		for _, cy := range ip.Agreed.Cycles {
			if cy.ID != "" && !cy.Done && slices.Contains(cy.Meta.Changes, c.ID) {
				waits += "  in cycle " + cy.ID + map[bool]string{true: ", unattended"}[cy.Meta.Unattended]
				if len(cy.WaitsOn) > 0 {
					waits += ", waits on cycle " + strings.Join(cy.WaitsOn, ", ")
				}
			}
		}
		fmt.Fprintf(stdout, "%s  %s  approved%s\n", c.ID, c.Meta.Title, waits)
		problems(stdout, c.Diagnostics)
	}
	return nil
}
