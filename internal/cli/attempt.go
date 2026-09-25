package cli

import (
	"bytes"
	"cmp"
	"fmt"
	"slices"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/mascah/grove/internal/attempt"
)

// attemptsTable lists attempts one per line, newest first.
func attemptsTable(views []attempt.View) []byte {
	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ATTEMPT\tWORK\tSTATUS\tSTARTED\tBRANCH\tEXIT\tCOST")
	for _, v := range views {
		exit, cost := "-", "-"
		if r := v.Result; r != nil {
			exit = fmt.Sprint(r.ExitCode)
			if r.Signal != "" {
				exit = r.Signal
			}
			if r.Stopped {
				exit += " stopped"
			}
			if r.Events.Result != nil {
				cost = fmt.Sprintf("%.2f", r.Events.Result.CostUSD)
			}
		}
		var work []string
		for _, m := range v.Launch.Members() {
			work = append(work, m.ID)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", v.Launch.Attempt, strings.Join(work, ","), v.Status, v.Launch.Started.UTC().Format(time.RFC3339), visible(v.Launch.Branch), exit, cost)
	}
	w.Flush()
	return buffer.Bytes()
}

// attemptText prints one attempt as facts, without any provider text.
func attemptText(v *attempt.View) string {
	return strings.Join(attempt.Facts(v, visible), "\n") + "\n"
}

// previewText prints what run would launch, and how to launch exactly that.
func previewText(l *attempt.Launch) string {
	reuse := "a new worktree from this checkout's HEAD"
	switch {
	case l.WorktreeReused:
		reuse = "its existing worktree"
	case l.Base != "" && l.Selection != nil && slices.ContainsFunc(l.Selection.Notes, func(n string) bool { return strings.HasPrefix(n, "Branch ") }):
		reuse = "the existing branch's tip"
	}
	lines := []string{
		"Preview: nothing was written or started",
		fmt.Sprintf("Checkout: %s", visible(l.Project)),
		fmt.Sprintf("Base: %s (%s)", l.Base, reuse),
		fmt.Sprintf("Worktree: %s on %s", visible(l.Worktree), visible(l.Branch)),
		fmt.Sprintf("Bounds: budget %s USD over the whole selection, permission mode %s, prompts none; one process, no retries", l.BudgetUSD, visible(l.PermissionMode)),
		fmt.Sprintf("Requested: %s, model %s, effort %s", map[bool]string{true: "until plan", false: "through to the handoff"}[l.Until == "plan"], visible(cmp.Or(l.Model, "default")), visible(cmp.Or(l.Effort, "default"))),
	}
	lines = append(lines, attempt.Explain(l, visible)...)
	lines = append(lines, "Digest: "+l.Selection.Digest,
		"Launch exactly this: grove run "+visible(strings.Join(l.Selection.Selected, " "))+" --expect "+l.Selection.Digest+" (with the same options)")
	return strings.Join(lines, "\n") + "\n"
}
