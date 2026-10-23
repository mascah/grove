# Implementation: Styled grove run output

<!-- Progress on change 261023-yz5bk, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Styled rows while a session runs

Status: done

- `internal/cli/run.go`: `relay(stdout, out, r, tty, width)` prints one
  row per event, `15:04:05 <mark> <text>`: `▶` init with its model, `•`
  text, `›` tool call, `✗` tool result with `is_error`, `·` thinking and
  other system events, `■` the result with `num_turns` and
  `total_cost_usd`. A row equal to the last (mark and text) is counted as
  `(×N)`. On a terminal: Lip Gloss colors, cut to the width with the count
  kept whole, and in-place updates (`\r` and erase-to-end, the cursor left
  on the last row). Otherwise: plain rows, each printed when the next
  differs or the session ends. `terminal` uses the `TIOCGWINSZ` ioctl from
  `syscall`, so Lip Gloss is the only dependency added. The `.log` and
  `.json` are unchanged.
- `go.mod`: `charm.land/lipgloss/v2` moved from indirect to direct;
  `go.sum` unchanged.
- `internal/cli/walkthrough_test.go`: `TestRun`'s stand-in session emits
  init, text, a successful and a failed tool result, three thinking blocks
  and a result. The plain output is checked row by row: no row for the
  success, `✗` for the error, `· thinking (×3)`, and
  `■ result: success, 3 turns, $0.42`.
- On the ledger fixture, `grove run <cycle-id>` under a 60-column pty with
  a stand-in `claude` printed colored rows (green `▶`/`■`, cyan `›`, red
  `✗`, faint `·`), the long text row cut to 60 cells, and `· thinking`
  rewritten in place from no count to `(×2)` and then `(×3)`.
- `docs/running.md` "An ordinary run", and `dev/walkthrough.md` running
  steps 2 (terminal) and 3 (output to a file).
- `just check`: passes; 3524 Go lines, 1171 test lines.

## Reconciliation

- runs: applied to specs/runs/spec.md
