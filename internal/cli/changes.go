package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/mascah/grove/internal/doc"
)

// listChanges lists the target's approved changes and what each waits on;
// as JSON, it prints its changes, archived ones included (D-64).
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
			Changes []*doc.Change `json:"changes"`
		}{p.Root, p.Target, ip.Commit, append(append([]*doc.Change{}, ip.Agreed.Changes...), ip.Agreed.Archived...)})
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
