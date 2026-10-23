## T-01 — Styled rows while a session runs

Result: `grove run` prints the session as [specs/runs/spec.md](specs/runs/spec.md)
"Show the session as it runs" says, and the documentation specifies it.

Files: `internal/cli/run.go`, `internal/cli/walkthrough_test.go`, `go.mod`, `go.sum`, `docs/running.md`, `dev/walkthrough.md`

Interfaces: `relay` takes whether standard output is a terminal, and its width.

Constraints: Lip Gloss is the only dependency added (rule 16); the `.log`
and `.json` are unchanged; within the rule 12 budget.

Acceptance:
- `docs/running.md`'s "An ordinary run" shows rows with times, marks and a
  `(×N)` thinking row, says tool results show only on error, and says a
  file or pipe gets plain rows; `dev/walkthrough.md`'s step matches and
  runs it in a terminal and with output sent to a file (rule 8).
- `TestRun` in `internal/cli/walkthrough_test.go`, whose output is not a
  terminal, sees plain rows from a stand-in session's events: a `▶` start,
  a `›` tool call, no row for a successful result, a `✗` row for an error,
  one `·` row counting repeated thinking, and a `■` result row.
- On the ledger fixture in a terminal, the rows are colored, cut to the
  width, and the thinking count updates in place.
- `just check` passes.
