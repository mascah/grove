package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/mascah/grove/internal/doc"
)

// listChanges lists the target's approved changes under each cycle that is
// not done, then those in no cycle; as JSON, it prints its cycles and its
// changes, archived ones included (D-64).
func listChanges(stdout io.Writer, asJSON bool) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ip, err := p.agreed()
	if err != nil {
		return err
	}
	if asJSON {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(struct {
			Root    string        `json:"root"`
			Target  string        `json:"target"`
			Commit  string        `json:"commit"`
			Cycles  []*doc.Cycle  `json:"cycles"`
			Changes []*doc.Change `json:"changes"`
		}{p.Root, p.Target, ip.Commit, append([]*doc.Cycle{}, ip.Agreed.Cycles...), append(append([]*doc.Change{}, ip.Agreed.Changes...), ip.Agreed.Archived...)})
	}
	listed := map[string]bool{}
	for _, cy := range ip.Agreed.Cycles {
		if cy.Done {
			continue
		}
		listed[cy.ID] = true
		archived, open := 0, 0
		for _, c := range ip.Agreed.Archived {
			if c.Meta.Cycle == cy.ID {
				archived++
			}
		}
		for _, c := range ip.Agreed.Changes {
			if c.Meta.Cycle == cy.ID {
				open++
			}
		}
		fmt.Fprintf(stdout, "%s  %s%s  %d of %d archived\n", cy.ID, cy.Meta.Title, dates(cy.Meta), archived, archived+open)
		problems(stdout, "  ", cy.Diagnostics)
		for _, c := range ip.Agreed.Changes {
			if c.Meta.Cycle == cy.ID {
				listChange(stdout, "  ", c)
			}
		}
	}
	for _, c := range ip.Agreed.Changes {
		if !listed[c.Meta.Cycle] {
			listChange(stdout, "", c)
		}
	}
	return nil
}

func listChange(stdout io.Writer, indent string, c *doc.Change) {
	waits := ""
	if len(c.WaitsOn) > 0 {
		waits = "  waits on " + strings.Join(c.WaitsOn, ", ")
	}
	fmt.Fprintf(stdout, "%s%s  %s  approved%s\n", indent, c.ID, c.Meta.Title, waits)
	problems(stdout, indent+"  ", c.Diagnostics)
}

// dates are a cycle's start and end, those it has.
func dates(m doc.Meta) string {
	switch {
	case m.Start != "" && m.End != "":
		return "  " + m.Start + " to " + m.End
	case m.Start != "":
		return "  from " + m.Start
	case m.End != "":
		return "  until " + m.End
	}
	return ""
}
