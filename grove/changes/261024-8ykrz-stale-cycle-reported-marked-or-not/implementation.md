# Implementation: A stale cycle is reported whether or not it is marked

<!-- Progress on change 261024-8ykrz, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

The owner agreed to build on main's text at 1b101b6. Every task file had
changed since base 2d07433, but none of those changes touched the cycle
staleness code or the steps this change edits.

## Tasks

### T-01 — Check reports a stale cycle whether or not it is marked

Status: done

- `internal/cli/check.go`: `stalePlans` reports every cycle that is not done
  and has a `planned`. Only a marked cycle's line adds `, so grove run
  skips this cycle`. A `planned` that is not a commit is reported for
  either kind. `movedSinceBase` is unchanged.
- `internal/cli/walkthrough_test.go`, `TestMovedFiles`: an unmarked cycle
  whose member's file moved gives the unmarked line and exit status 1.
  With `planned: 0000000`, it gives the not-a-commit line. Marked, it gives
  today's line. Before the fix, check printed `No problems in grove/.` with
  status 0 for the first two. After the fix the test passes.
- `docs/planning.md`: "Re-plan a cycle" now describes and quotes the stale
  report for any cycle, and names the marked cycle's extra words. "Mark it
  unattended" points to it.
- `dev/walkthrough.md`: Planning cycles step 5 cites "Re-plan a cycle". It
  touches `ledger/report.py`, which budgets names, shows the 0.4 cycle's
  unmarked line, and re-plans that cycle. Its show-saved comment uses a
  file the 0.4 cycle's changes do not name, so only show-saved's cycle is
  reported.
- `just check`: passed.

### T-02 — Implement names a member's files moved since the plan

Status: not started

## Scenarios

<!-- Before finishing, one line per scenario of the change specs: the check
that shows it and its result, or `unverified` and why:

    - record-editing "Update the note on a legacy record": tests/test_cli.py test_edit_legacy, passed
-->

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
