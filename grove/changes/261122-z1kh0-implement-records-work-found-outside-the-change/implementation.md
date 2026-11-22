# Implementation: Implement records work it finds outside the change

<!-- Progress on change 261122-z1kh0, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Implement writes and reports `Found:` lines

Status: done

- `internal/instruct/requires/implement.md`: a new rule after the
  `Also changed:` one: work the change does not ask for and no task needs
  is left as it is and written as
  `Found: <what is wrong or missing, where, and how it was seen>` under the
  task being built; a file a task needs or makes false stays
  `Also changed:`, a wrong approach a revision; Grove reads nothing from
  it. The finish line lists every `Found:` line after the `Check:` lines,
  blocked questions and unverified scenarios, before what was done.
- The implementation guide names `Found:` beside `Also changed:`; the
  template's comment shows `Also changed:` and `Found:` examples beside
  `Check:` and `Answer:`.
- `internal/cli/walkthrough_test.go`: `TestImplementInstructions` checks
  the rule, the report order and the template example in
  `grove instructions implement`; `TestArchive` adds a `Found:` line to an
  archived change on the branch and checks that `grove review`'s readout
  and `grove changes --json` stay the same.
- `bin/grove instructions implement 261122-z1kh0` prints the rule and the
  finish line. `go test ./internal/cli/ -run 'TestImplement|TestArchive'`
  passed. Instruction text is 8,382 words of 9,000.

Check: `grove review` prints the branch's diff after its readout, and the diff shows any line added to a file, a `Found:` line included. The test compares the readout above the diff, and `grove changes --json` without the time it was read and the branch tip's hash, date and subject, which a new commit always changes. The `Found:` line itself changes neither.
Check: `internal/cli/walkthrough_test.go` changed on main after this change was proposed (#189, #194: explore and propose tests). Those edits are to other tests; this task's edits sit beside them on the file as main holds it.

Found: the cycle instructions' finish line (`internal/instruct/requires/implement-cycle.md`) asks each change's subagent to report its status, `Check:` lines and blocked tasks, and the session to report those, but not `Found:` lines; seen reading it beside `implement.md`. The lines still reach each change's `implementation.md` through `instructions implement --runner`.

## Scenarios

<!-- Before finishing, one line per scenario of the change specs: the check
that shows it and its result, or `unverified` and why:

    - record-editing "Update the note on a legacy record": tests/test_cli.py test_edit_legacy, passed
-->

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
