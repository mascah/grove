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

Status: done

- `internal/cli/check.go`: `sincePlan` counts the commits since a cycle's
  plan, as "Count a member changed with its plan" says, and returns the
  paths they changed. `changedSincePlan`, which check and `grove run`
  call, derives the stale members from those paths, as before. `naming`
  filters a change's tasks by changed paths, and `moved` uses it too.
- `internal/cli/instructions.go`: for a member of a cycle the target
  holds that is not done, `instructionsChange` adds
  ``- Changed on main since cycle <cycle-id>'s plan <sha7>: T-01's `<file>`.``
  It does this with the owner and with `--unattended`. A change in no
  cycle still gets the line against `base`. `grove instructions implement
  <cycle-id>` is unchanged.
- `internal/instruct/requires/implement-owner.md`,
  `implement-unattended.md`: the rule reads "changed on the target",
  without "since the change's base", so it covers both lines. That is
  three fewer words in each file.
- `internal/cli/walkthrough_test.go`: `TestMovedFiles` checks that an
  unmarked member's moved file is named against the plan, with and without
  `--unattended`, and not against `base`. `TestFollows` checks three
  things: after the merges that archive L's and M's changes, which edit
  B's file, B's instructions name nothing. After a commit that archives
  nothing, they name the file. After K is re-planned, they name nothing.
  Both tests failed with `instructions.go` reverted, and pass with it.
- `docs/coordinating.md`: says a cycle member's file that changed since
  the plan is named as one changed since `base` is, quotes the line, and
  points to "Re-plan a cycle" instead of "Mark it unattended". T-01 moved
  the stale report into that section.
- `dev/walkthrough.md`: Planning cycles step 5 shows
  `grove instructions implement <budgets-id>` naming T-01's
  `ledger/report.py` against the 0.4 cycle's plan before the re-plan.
- `just check`: passed (3902 Go lines, 1603 test lines, 7697 instruction
  words).

## Scenarios

<!-- Before finishing, one line per scenario of the change specs: the check
that shows it and its result, or `unverified` and why:

    - record-editing "Update the note on a legacy record": tests/test_cli.py test_edit_legacy, passed
-->

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
