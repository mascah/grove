# Implementation: Tasks name their files

<!-- Progress on change 261023-6j5bw, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Tasks name their files

Status: done

- The tasks guide says a task names its files on a `Files:` line after
  `Result:` (backticked repo-relative paths, comma-separated) and its
  interfaces on an `Interfaces:` line, with no code, commands or steps; the
  tasks template shows both lines in its example and skeleton.
- The propose guide's "no file-by-file plan" line is gone: each task names
  its files and interfaces, so where the change goes is settled at propose.
- `requires/implement.md` and `requires/apply.md` say to change only the
  files the tasks name: implement asks the owner to add a file to the task,
  apply leaves it untouched and names it in the task's evidence. The schema
  implement and apply guides drop their old "Files no task names" line, now
  said in what Grove requires.
- Docs: `docs/concepts.md`'s table, `docs/reference.md` (what a schema cannot
  change, and the bugfix example guide naming its files with both lines),
  `docs/coordinating.md` ("Applying a cycle"), and `dev/walkthrough.md`'s
  propose and implement steps of the main line.
- `internal/cli/walkthrough_test.go`: `TestPropose` checks the propose
  instructions print the `Files:` rule and template line and no
  "file-by-file"; `TestApplyInstructions` checks apply and implement print
  their rule. `just check` passes (3294 Go lines, 1079 test lines).
- Not run: the acceptance's hand check on the ledger fixture in Claude Code
  (`/grove-propose`, then `/grove-implement`); this unattended session has
  no owner to answer propose's questions.
- A file no task names that looks like it should change, left as it is:
  `docs/getting-started.md` line 92 says "`tasks.md`: the work, each task
  with how to check it is done", which no longer mentions the files a task
  names.

### T-02 — Check reports a named file that moved

Status: done

- `internal/doc/doc.go`: `doc.Task` gains `Files`, read from the task's
  `Files:` line, which may wrap until a blank line or the next field; a
  line that is not backticked paths separated by commas, or names a path
  that is not repo-relative (`fs.ValidPath`), is a diagnostic on
  `tasks.md` at its line.
- `internal/cli/check.go`: `movedSinceBase` reports, for each change in no
  cycle that is not done, each task's named file that `git diff
  --name-only <base> <target> -- <files>` lists, one Git call per change:
  "T-01's ledger/cli.py changed on main since base abc1234; revise the
  change with propose". `changedSincePlan` adds each open member's named
  files to its one diff, so the cycle report keeps its words and `grove
  run`'s `eligible`, which calls it, skips the cycle; `run.go` needed no
  edit.
- `internal/cli/instructions.go`: `instructions implement` for a change in
  no cycle that is not done appends "Changed on main since base …: T-01's
  `file`." to what Grove requires; `requires/implement.md` says to tell the
  owner what changed and ask before building that task.
- Docs: `docs/coordinating.md` ("When the plan is stale"),
  `docs/planning.md` (stale cycle, re-plan), `docs/reference.md`
  (`grove check` row), `dev/walkthrough.md` (coordinating step 5, planning
  step 5).
- `internal/cli/walkthrough_test.go`, `TestMovedFiles`: unreported when
  unchanged and when only the branch changed the file; reported against
  base once main changed it; `instructions implement` names T-01 and the
  file; as a member of a marked cycle, unreported against base while the
  cycle is reported stale, `eligible` (grove run) selects no cycle, and
  implement names nothing; an unreadable `Files:` line reported at its
  line. Removing the cycle condition from `movedSinceBase` fails the test.
  `just check` passes (3396 Go lines, 1139 test lines); `grove check` on
  this repository is clean.
- Not run: the walkthrough's hand steps in Claude Code.

## Reconciliation

- cycles: applied to specs/cycles/spec.md (Purpose's stale definition, and
  the MODIFIED "Report a cycle that cannot run unattended" and "Revise a
  cycle")
- implement: applied to specs/implement/spec.md
- propose: applied to a new specs/propose/spec.md, with its .grove.yaml
