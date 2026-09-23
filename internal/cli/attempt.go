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
	l := &v.Launch
	var b strings.Builder
	line := func(format string, args ...any) { fmt.Fprintf(&b, format+"\n", args...) }
	line("Attempt: %s (%s)", l.Attempt, v.Status)
	line("Work: %s at %s, record %s", l.Work, visible(l.RecordPath), l.RecordRevision)
	if v.InputsChanged != "" {
		line("Inputs changed: %s", visible(v.InputsChanged))
	}
	reuse := "created"
	if l.WorktreeReused {
		reuse = "reused"
	}
	line("Worktree: %s on %s from %s (%s)", visible(l.Worktree), visible(l.Branch), short(l.Base), reuse)
	line("Started: %s by %s with %s (%s), owner pid %d", l.Started.UTC().Format(time.RFC3339), l.GroveVersion, visible(l.Executable), visible(l.ClaudeVersion), l.Owner)
	line("Bounds: budget %s USD, permission mode %s, prompts none; one process, no retries; subagents share the budget", l.BudgetUSD, l.PermissionMode)
	line("Session: %s", l.SessionID)
	line("Command: %s", visible(strings.Join(l.Command, " ")))
	ev := v.Events
	if v.Result != nil {
		ev = &v.Result.Events
	}
	if ev != nil {
		var types []string
		for _, k := range []string{"system", "assistant", "user", "stream_event", "rate_limit_event", "result"} {
			if ev.Types[k] != 0 {
				types = append(types, fmt.Sprintf("%s %d", k, ev.Types[k]))
			}
		}
		line("Events: %d lines, %d bytes (%s); unknown %d, malformed %d, oversized %d, partial %v", ev.Lines, ev.Bytes, strings.Join(types, ", "), ev.Unknown, ev.Malformed, ev.Oversized, ev.Partial)
		if i := ev.Init; i != nil {
			line("Provider: %s, model %s, permission mode %s, %d tools, capabilities %s", visible(i.Version), visible(i.Model), visible(i.PermissionMode), i.Tools, visible(strings.Join(i.Capabilities, " ")))
		}
		if f := ev.Result; f != nil {
			line("Result event: %s, is_error %v, %d turns, %.4f USD, %d permission denials, %d bytes of text, session %s", visible(f.Subtype), f.IsError, f.Turns, f.CostUSD, f.PermissionDenials, ev.ResultText, visible(f.SessionID))
		} else {
			line("Result event: none")
		}
	}
	if r := v.Result; r != nil {
		exit := fmt.Sprintf("exit %d", r.ExitCode)
		if r.Signal != "" {
			exit = r.Signal
		}
		how := ""
		if r.Stopped {
			how = ", stopped"
		}
		if r.ReconciledBy != "" {
			how += ", reconciled by " + r.ReconciledBy
		}
		line("Finished: %s, %s%s", r.Finished.UTC().Format(time.RFC3339), exit, how)
		dirty := "clean"
		if r.Dirty {
			dirty = "uncommitted or untracked changes"
		}
		line("Worktree after: HEAD %s, %s", short(r.Head), dirty)
		if r.Record != nil {
			line("Record on the branch: %s %s, candidate %s, revision %s", l.Work, visible(r.Record.Status), visible(orNone(r.Record.Candidate)), r.Record.Revision)
			if r.RecordError != "" {
				line("Record problems: %s", visible(r.RecordError))
			}
		} else {
			line("Record on the branch: unreadable: %s", visible(r.RecordError))
		}
	}
	line("Files: %s", visible(v.Dir))
	return b.String()
}

func orNone(s string) string {
	if s == "" {
		return "none"
	}
	return s
}
