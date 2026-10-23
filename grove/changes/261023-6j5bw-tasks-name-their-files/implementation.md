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

Status: not started

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
