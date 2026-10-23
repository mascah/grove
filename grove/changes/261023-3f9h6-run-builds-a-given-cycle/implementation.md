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

Status: not started

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
