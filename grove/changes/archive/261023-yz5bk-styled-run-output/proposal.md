# Styled grove run output

## Problem

The owner's usage feedback of 2026-10-23: "grove run stream output
definitely needs to be more readable when interactive. ... The way we're
parsing things right now is too verbose. We should take some inspiration
from how ../grove-legacy handled displaying background attempts in the
TUI. I especially prefer how we stack multiple thinking events."

- `relay` in `internal/cli/run.go` at 2f8f44d prints, plainly, the first
  line of each assistant message, each tool call with the first line of its
  input, and, indented, the first line of every tool result, so most lines
  are tool output. Thinking and system events are not shown at all.
- grove-legacy's attempt view (`internal/tui/attempts.go` lines 1368-1401
  and `internal/harness/activity.go` lines 189-196, in `../grove-legacy`)
  showed one row per event: a time, a mark and color for each kind, the
  row cut to the width, and a repeat of the previous row counted as
  `(×N)` instead of printed again.

`CONSTITUTION.md` rule 16, as amended in #32, lists Lip Gloss.

## Intended outcome

While a session runs, a terminal shows one short, colored row per event
in the legacy layout, with repeated rows such as thinking stacked into one
counted row, and tool results only when they are errors. A timer's log
gets the same rows as plain text. The attempt's `.log` is unchanged.

## Scope

### In scope

- What `grove run` prints to standard output while the session runs, in a
  terminal and not.
- The documentation and walkthrough steps showing it.

### Out of scope

- The attempt's `.log` file and `.json` record, which stay as they are.
- A board listing attempts, as grove-legacy's TUI had.
- Token-by-token output (`--include-partial-messages`).
- `grove run <cycle-id>`: change 261023-3f9h6.

## Decisions

All from the owner, 2026-10-23, "yes to all" to the round, after "Yes rule
16 can be amended, so lets make it look nice".

- **Lip Gloss styles the rows; charm log is not used.** Reason: charm log
  is a leveled key-value logger, not a display; Lip Gloss is already in the
  build through Glamour.
- **One row per event, `15:04:05 › Bash go test`,** marked `▶` session
  start with its model, `•` text (its first line), `›` tool call, `✗` tool
  error, `·` thinking and other notices, `■` the final result, green on
  success and red otherwise. Reason: grove-legacy's layout, which the owner
  prefers.
- **A tool result is shown only when it is an error,** as `✗` and its
  first line. Reason: tool output is most of today's verbosity.
- **Consecutive identical rows collapse into one with `(×N)`;** in a
  terminal the count updates in place, otherwise the row is printed once
  the run of repeats ends. Reason: the owner's "stack multiple thinking
  events".
- **The time is when Grove reads the event,** so every row has one.
- **In a terminal rows are cut to its width; otherwise no color, no cut
  and no in-place update.** Reason: a timer's log stays plain text.
- **The final `■` row gives the session's turn count and cost** from its
  result event.
- **The `.log` stays every event as JSON Lines, and the display stays in
  `internal/cli/run.go`,** `cli` being the only package that writes to the
  terminal.
- **No `depends_on`.** This change and 261023-3f9h6 both edit
  `internal/cli/run.go` and `docs/running.md`.

## Open questions

None.

## Related documents

- Requirements: [specs/runs/spec.md](specs/runs/spec.md) (maintained spec:
  `grove/specs/runs/spec.md`)
- Tasks: [tasks.md](tasks.md)
