package cli

import (
	"bytes"
	"fmt"
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
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", v.Launch.Attempt, v.Launch.Work, v.Status, v.Launch.Started.UTC().Format(time.RFC3339), visible(v.Launch.Branch), exit, cost)
	}
	w.Flush()
	return buffer.Bytes()
}

// attemptText prints one attempt as facts, without any provider text.
func attemptText(v *attempt.View) string {
	return strings.Join(attempt.Facts(v, visible), "\n") + "\n"
}
