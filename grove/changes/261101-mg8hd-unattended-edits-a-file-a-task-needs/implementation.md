# Implementation: Implement edits a file a task needs and says why

<!-- Progress on change 261101-mg8hd, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Implement edits a file a task needs and says why

Status: done

- `internal/instruct/requires/implement-owner.md` and
  `implement-unattended.md`: both say to make the choices and change the
  paths a `Delegated:` line names and, beyond them and the `Files:` lines,
  to edit any other file a task needs changed or makes false, without
  asking and without blocking the task, writing
  `` Also changed: `<path>` — <why> `` under it, and never `grove.yaml`
  for it. The owner text no longer says to ask first; the unattended text
  no longer says such a file stays as it is, to ignore `Delegated:` lines,
  or to write no `Also changed:` line. The bullet saying not to edit the
  change's files other than `implementation.md` stays.
- `internal/instruct/schemas/default/documents/implementation/guide.md`:
  an `Also changed:` line records a file the implementing agent changed,
  and why. `.../tasks/guide.md`: `Delegated:` is honored in both modes,
  and a file a task needs or makes false is recorded with `Also changed:`.
- `internal/cli/review.go`: `mark` labels such a file
  `also changed <change-id>:T-NN: <why>`.
- `internal/cli/walkthrough_test.go`: `TestImplementInstructions` checks
  the new text in both modes and the absence of the old;
  `TestPropose` checks the tasks guide; `TestPlanned` checks
  `also changed` and that `Cargo.lock` is still `named nowhere`.
- `docs/coordinating.md`, `docs/getting-started.md`, `docs/reference.md`:
  both modes and the label, with no "once you agree".
- `dev/walkthrough.md`: "Delegated with me, ignored unattended" is now
  "Delegated, with me and unattended", expecting `grove run <cycle-id>` to
  create `tests/test_records.py`; "A file no task names, with me" expects
  the edit without asking and review's `also changed`; "Planned against
  actual" shows `also changed`.
- Would have changed `internal/instruct/requires/implement.md`, which no
  `Files:` line names: its shared line "Change only the files the tasks
  name on their `Files:` lines." still prints above both modes' new
  bullet, which reads as the exception ("Beyond them and the `Files:`
  lines, ..."). Left as it is.
- `just check`: passed (4031 Go lines, 1697 test lines, 8048 words).
