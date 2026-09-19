package cli

import (
	"fmt"
	"io"

	"github.com/mascah/grove/internal/versions"
)

// runWorkspace resolves an explicitly selected version to its existing
// checkout. The project directory goes to stdout so a caller can pass it to
// --project; checkout, branch, record, and revision context go to stderr.
func runWorkspace(root string, a invocation, out, errOut io.Writer) int {
	w, err := versions.Resolve(root, a.source)
	if err != nil {
		report(errOut, err)
		return 1
	}
	return writeWorkspace(w, a.json, out, errOut)
}

// writeWorkspace is the one result contract for a resolved workspace, shared
// by the workspace command and the board.
func writeWorkspace(w *versions.Workspace, asJSON bool, out, errOut io.Writer) int {
	branch := "Branch: " + w.Ref
	if w.Ref == "" {
		branch = "Detached: HEAD"
	}
	fmt.Fprintf(errOut, "Checkout: %s\n%s at %s\nRecord: %s\nRevision: %s\n", visible(w.Checkout), visible(branch), w.Head, visible(w.Record), w.Revision)
	if asJSON {
		var ref any
		if w.Ref != "" {
			ref = w.Ref
		}
		return writeResult(out, errOut, marshal(map[string]any{
			"checkout": w.Checkout, "project": w.Project, "record": w.Record,
			"ref": ref, "head": w.Head, "revision": w.Revision, "selector": w.Selector,
		}))
	}
	return writeResult(out, errOut, []byte(w.Project+"\n"))
}
