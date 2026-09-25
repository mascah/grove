package attempt

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"
)

// Facts lists one attempt as lines of facts, without any provider text, for
// the CLI's attempt and the board's attempt screen. visible escapes each value
// that came from a file or a process for the caller's display.
func Facts(v *View, visible func(string) string) []string {
	l := &v.Launch
	var lines []string
	line := func(format string, args ...any) { lines = append(lines, fmt.Sprintf(format, args...)) }
	short := func(commit string) string { return commit[:min(len(commit), 12)] }
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
	switch {
	case l.Skill == "":
		line("Entrypoints: not recorded at launch")
	case len(l.Differs) == 0:
		line("Entrypoints: skill %s; the skill and any reviewer match the launching grove's templates", l.Skill)
	default:
		line("Entrypoints: skill %s; differ from the launching grove's templates: %s", l.Skill, visible(strings.Join(l.Differs, ", ")))
	}
	line("Agent's grove: not recorded; PATH or the project's instructions choose it, which may not be the launching grove")
	line("Bounds: budget %s USD, permission mode %s, prompts none; one process, no retries; subagents share the budget", l.BudgetUSD, l.PermissionMode)
	line("Requested: %s", visible(Requested(l)))
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
			if len(f.ModelCostUSD) != 0 {
				line("Cost by model: %s", visible(CostByModel(f.ModelCostUSD)))
			}
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
			uncommitted := ""
			if r.RecordUncommitted {
				uncommitted = ", uncommitted"
			}
			line("Record on the branch: %s %s, candidate %s, revision %s%s", l.Work, visible(r.Record.Status), visible(orNone(r.Record.Candidate)), r.Record.Revision, uncommitted)
			if r.RecordError != "" {
				line("Record problems: %s", visible(r.RecordError))
			}
		} else {
			line("Record on the branch: unreadable: %s", visible(r.RecordError))
		}
	}
	line("Files: %s", visible(v.Dir))
	return lines
}

// Requested is what the launch asked of the provider beyond the fixed
// command: the bound, model, effort and the reviewer definition present.
func Requested(l *Launch) string {
	until := "through to the handoff"
	if l.Until != "" {
		until = "until " + l.Until
	}
	reviewer := "reviewer " + ReviewerPath + " " + l.Reviewer
	switch l.Reviewer {
	case "none":
		reviewer = "no reviewer definition"
	case "":
		reviewer = "reviewer definition not recorded"
	}
	return fmt.Sprintf("%s, model %s, effort %s, %s", until, cmp.Or(l.Model, "default"), cmp.Or(l.Effort, "default"), reviewer)
}

// CostByModel lists the result's cost per model, most expensive first.
func CostByModel(costs map[string]float64) string {
	names := slices.SortedFunc(maps.Keys(costs), func(a, b string) int {
		return cmp.Or(cmp.Compare(costs[b], costs[a]), strings.Compare(a, b))
	})
	parts := make([]string, len(names))
	for i, n := range names {
		parts[i] = fmt.Sprintf("%s $%.2f", n, costs[n])
	}
	return strings.Join(parts, ", ")
}

func orNone(s string) string {
	if s == "" {
		return "none"
	}
	return s
}
