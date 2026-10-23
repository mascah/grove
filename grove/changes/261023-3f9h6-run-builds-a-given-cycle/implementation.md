# Implementation: grove run builds a cycle it is given

<!-- Progress on change 261023-3f9h6, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — grove run builds a cycle it is given

Status: done

- `internal/cli/cli.go`: `run` takes one optional argument, passed to
  `runAttempt`; the usage line reads
  `run [--retry] [--timeout <duration>] [<cycle-id>]`.
- `internal/cli/run.go`: `runAttempt(stdout, id, retry, timeout)`. The gates
  `eligible` applied moved into `runnable`, which returns why a cycle cannot
  be built (done, `needs_human`, waiting, stale plan); `eligible` keeps only
  the `unattended` check over it. A given ID is looked up on the fetched
  target and refused with `runnable`'s reason, after the pause check and
  before any branch is created.
- `internal/cli/walkthrough_test.go`: `TestRun` gives an unknown ID and a
  cycle whose member is under `needs_human` (refusals, no new branch), then
  an unmarked cycle while marked cycle `l` is eligible, and sees the
  session on it.
- `docs/running.md` "Run a cycle you name", a `docs/reference.md` row, and
  `dev/walkthrough.md` running step 6.
- `just check`: passes; 3429 Go lines, 1163 test lines.

### T-02 — Plan's finish names /grove-implement first

Status: done

- `internal/instruct/requires/plan.md`: the finish says, in order,
  `/grove-implement <cycle-id>`, whether `grove run <cycle-id>` would build
  the cycle, and `unattended: true` for a `grove run` given no cycle; a
  member under `needs_human` makes `grove run` refuse the cycle, given or
  marked. The `unattended` bullet now says the mark hands the cycle to a
  `grove run` given no cycle.
- `docs/planning.md`: the intro and both finish examples;
  `dev/walkthrough.md` planning steps 1 and 3.
- `internal/cli/walkthrough_test.go`: `TestCycles` checks the three appear
  in that order in `grove instructions plan`.
- `just check`: passes; 3429 Go lines, 1167 test lines.

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
