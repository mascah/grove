# Check names each moved file once per change

## Problem

`grove check` reports a change in no open cycle once for every file of
every task whose named file changed on the target since the change's
`base`: `movedSinceBase` in `internal/cli/check.go` at 999ab56 appends one
diagnostic per task per file. A file several tasks name is reported once
per task, and a change touching a hub file fills the output.

The owner's brief, `memex/docs/efficiency-2026-11-14.md` (2026-11-14),
measured it on a project using Grove: "`grove check` on memex main prints
53 problems, all drift on hub files (`App.tsx`, `api.ts`, `main.rs`,
`README.md`, `CLAUDE.md`)."

## Intended outcome

`grove check` prints one line for such a change, on its `tasks.md`, naming
each changed file once and saying, as today, to revise the change with
propose, which also re-bases it after a dependency lands. The count of
problems is the number of drifted changes, not of files times tasks.

## Scope

### In scope

- The message and count of `check`'s report of a task's named file changed
  on the target since `base` (`grove/specs/propose/spec.md`).
- The docs and walkthrough steps that quote it.

### Out of scope

- The lines `grove instructions implement` prints naming each task's
  changed file: the implementing agent needs them per task, and they are
  one line already.
- The stale-plan report of a cycle, already one line per member.

## Decisions

The owner's brief, 2026-11-14, set the direction: "`grove check` prints
drift once per change, not once per file per task." It also says "Do not
ask a round. Take the recommended choice, write it down where the reviewer
will see it, and go on", so each item below is the recommended answer,
taken under that direction.

1. **One line per change, on `tasks.md`, naming files only.** The line
   lists each changed file once, in path order, joined by "and", and names
   no task: the implement instructions name the tasks, and a file several
   tasks name would otherwise repeat. Example: `grove/changes/<id>/tasks.md:
   ledger/report.py and tests/test_cli.py changed on main since base
   4f1c2a9; revise the change with propose, which also re-bases it after a
   dependency lands`.
2. **Not breaking.** No document clean before reports a problem now, and
   none is renamed; the title is `fix(cli)`.

## Open questions

None.

## Related documents

- Requirements: [specs/propose/spec.md](specs/propose/spec.md)
- Tasks: [tasks.md](tasks.md)
